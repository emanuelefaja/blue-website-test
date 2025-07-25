package web

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// SearchItem represents a single searchable document
type SearchItem struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content"`
	URL         string `json:"url"`
	Type        string `json:"type"` // "page", "doc", "blog", etc.
	Section     string `json:"section,omitempty"`
	Category    string `json:"category,omitempty"` // Display category like "Feature", "Docs", "API"
}

// SearchFrontmatter represents the YAML frontmatter in markdown files
type SearchFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
}

// GenerateSearchIndex creates a JSON search index from all content
func GenerateSearchIndex() error {
	var searchItems []SearchItem

	// Load metadata.json for title lookup
	// TODO: Update to support multi-language metadata structure
	var metadata *Metadata = nil
	// metadata, err := loadMetadata()
	// if err != nil {
	// 	fmt.Printf("Warning: Could not load metadata.json: %v\n", err)
	// }

	// Index HTML pages
	if err := indexHTMLPages(&searchItems, metadata, nil); err != nil {
		return fmt.Errorf("failed to index HTML pages: %w", err)
	}

	// Index markdown content
	if err := indexMarkdownContent(&searchItems); err != nil {
		return fmt.Errorf("failed to index markdown content: %w", err)
	}

	// Write search index to public directory
	return writeSearchIndex(searchItems)
}

// GenerateSearchIndexWithCache creates a JSON search index using pre-rendered markdown cache
func GenerateSearchIndexWithCache(markdownService *MarkdownService) error {
	var searchItems []SearchItem

	// Load metadata.json for title lookup
	// TODO: Update to support multi-language metadata structure
	var metadata *Metadata = nil
	// metadata, err := loadMetadata()
	// if err != nil {
	// 	fmt.Printf("Warning: Could not load metadata.json: %v\n", err)
	// }

	// Index HTML pages
	if err := indexHTMLPages(&searchItems, metadata, nil); err != nil {
		return fmt.Errorf("failed to index HTML pages: %w", err)
	}

	// Index cached markdown content
	if err := indexCachedMarkdownContent(&searchItems, markdownService); err != nil {
		return fmt.Errorf("failed to index cached markdown content: %w", err)
	}

	// Write search index to public directory
	return writeSearchIndex(searchItems)
}

// GenerateSearchIndexWithCaches creates a JSON search index using both HTML and markdown caches
func GenerateSearchIndexWithCaches(markdownService *MarkdownService, htmlService *HTMLService) error {
	// Load metadata.json for title lookup
	metadata, err := loadMetadata()
	if err != nil {
		fmt.Printf("Warning: Could not load metadata.json: %v\n", err)
	}

	// Collect all items first
	allItems, err := collectAllSearchItems(markdownService, htmlService, metadata)
	if err != nil {
		return err
	}

	// Group items by language
	itemsByLanguage := make(map[string][]SearchItem)
	
	for _, item := range allItems {
		lang, _ := extractLanguageFromPath(item.URL)
		if lang == "" {
			lang = DefaultLanguage
		}
		itemsByLanguage[lang] = append(itemsByLanguage[lang], item)
	}

	// Generate search index for each language in parallel
	type languageResult struct {
		lang string
		err  error
	}
	
	resultChan := make(chan languageResult, len(itemsByLanguage))
	
	// Process each language concurrently
	for lang, items := range itemsByLanguage {
		go func(lang string, items []SearchItem) {
			filename := "public/searchIndex.json"
			if lang != DefaultLanguage {
				filename = fmt.Sprintf("public/searchIndex-%s.json", lang)
			}
			
			data, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				resultChan <- languageResult{lang, fmt.Errorf("failed to marshal search index for language %s: %w", lang, err)}
				return
			}
			
			if err := os.WriteFile(filename, data, 0644); err != nil {
				resultChan <- languageResult{lang, fmt.Errorf("failed to write search index for language %s: %w", lang, err)}
				return
			}
			
			resultChan <- languageResult{lang, nil}
		}(lang, items)
	}
	
	// Wait for all languages to complete and check for errors
	for i := 0; i < len(itemsByLanguage); i++ {
		result := <-resultChan
		if result.err != nil {
			return result.err
		}
	}

	return nil
}

// collectAllSearchItems collects all search items from various sources
func collectAllSearchItems(markdownService *MarkdownService, htmlService *HTMLService, metadata *Metadata) ([]SearchItem, error) {
	// Use channels to collect results from goroutines
	htmlItemsChan := make(chan []SearchItem, 1)
	markdownItemsChan := make(chan []SearchItem, 1)
	nonCachedItemsChan := make(chan []SearchItem, 1)
	errorChan := make(chan error, 3)

	// Index cached HTML pages concurrently
	go func() {
		var items []SearchItem
		if err := indexCachedHTMLPages(&items, htmlService, metadata); err != nil {
			errorChan <- fmt.Errorf("failed to index cached HTML pages: %w", err)
			return
		}
		htmlItemsChan <- items
	}()

	// Index cached markdown content concurrently
	go func() {
		var items []SearchItem
		if err := indexCachedMarkdownContent(&items, markdownService); err != nil {
			errorChan <- fmt.Errorf("failed to index cached markdown content: %w", err)
			return
		}
		markdownItemsChan <- items
	}()

	// Wait for HTML items to build indexed URLs map
	var searchItems []SearchItem
	select {
	case items := <-htmlItemsChan:
		searchItems = append(searchItems, items...)
	case err := <-errorChan:
		return nil, err
	}

	// Build a map of indexed URLs from cached pages
	indexedURLs := make(map[string]bool)
	for _, item := range searchItems {
		if item.Type == "page" {
			indexedURLs[item.URL] = true
		}
	}

	// Index non-cached HTML pages (like status page) concurrently
	go func() {
		var items []SearchItem
		if err := indexHTMLPages(&items, metadata, indexedURLs); err != nil {
			errorChan <- fmt.Errorf("failed to index non-cached HTML pages: %w", err)
			return
		}
		nonCachedItemsChan <- items
	}()

	// Collect remaining results
	for i := 0; i < 2; i++ {
		select {
		case items := <-markdownItemsChan:
			searchItems = append(searchItems, items...)
		case items := <-nonCachedItemsChan:
			searchItems = append(searchItems, items...)
		case err := <-errorChan:
			return nil, err
		}
	}

	return searchItems, nil
}

// loadMetadata loads the metadata.json file
func loadMetadata() (*Metadata, error) {
	// Use the SEOService to load metadata properly
	seoService := NewSEOService()
	if err := seoService.LoadData(); err != nil {
		return nil, err
	}
	return seoService.metadata, nil
}

// getPageKey converts URL path to metadata key
func getPageKey(path string) string {
	if path == "/" {
		return "home"
	}

	// Remove leading/trailing slashes
	cleanPath := strings.Trim(path, "/")
	return cleanPath
}

// indexHTMLPages indexes all HTML pages in the pages directory
func indexHTMLPages(items *[]SearchItem, metadata *Metadata, indexedURLs map[string]bool) error {
	return filepath.WalkDir("pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Convert file path to URL
		url := "/" + strings.TrimSuffix(path, ".html")
		if strings.HasSuffix(url, "/index") {
			url = strings.TrimSuffix(url, "/index")
		}
		if url == "/pages" {
			url = "/"
		}
		url = strings.TrimPrefix(url, "/pages")
		if url == "" {
			url = "/"
		}

		// Extract title using the new algorithm
		title := extractPageTitle(string(content), url, path, metadata)

		// For non-cached HTML pages, use metadata description as content if available
		var textContent string
		var description string
		if metadata != nil {
			pageKey := getPageKey(url)
			lang, _ := extractLanguageFromPath(url)
			if lang == "" {
				lang = DefaultLanguage
			}
			
			if pageData, exists := metadata.Pages[pageKey]; exists {
				// Try language-specific description first
				if langMeta, langExists := pageData[lang]; langExists && langMeta.Description != "" {
					textContent = langMeta.Description
					description = langMeta.Description
				} else if enMeta, enExists := pageData["en"]; enExists && enMeta.Description != "" {
					// Fallback to English
					textContent = enMeta.Description
					description = enMeta.Description
				}
			}
		}
		
		// Fallback to HTML extraction if no metadata description
		if textContent == "" {
			textContent = extractTextFromHTML(string(content))
		}

		// Skip if this URL was already indexed from cache
		if indexedURLs != nil && indexedURLs[url] {
			return nil
		}

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Content:     textContent,
			URL:         url,
			Type:        "page",
		})

		return nil
	})
}

// extractPageTitle implements the title extraction algorithm:
// 1. Check metadata.json first
// 2. Extract from H1 tags
// 3. Generate clean title from filename
func extractPageTitle(htmlContent, url, filePath string, metadata *Metadata) string {
	// 1. Check metadata.json first
	if metadata != nil {
		pageKey := getPageKey(url)
		lang, _ := extractLanguageFromPath(url)
		if lang == "" {
			lang = DefaultLanguage
		}
		
		if pageData, exists := metadata.Pages[pageKey]; exists {
			// Try language-specific title first
			if langMeta, langExists := pageData[lang]; langExists && langMeta.Title != "" {
				return langMeta.Title
			}
			// Fallback to English
			if enMeta, enExists := pageData["en"]; enExists && enMeta.Title != "" {
				return enMeta.Title
			}
		}
	}

	// 2. Try to extract from H1 tags
	if title := extractH1Title(htmlContent); title != "" {
		return title
	}

	// 3. Try to extract from title tags (fallback)
	if title := extractHTMLTitle(htmlContent); title != "" {
		return title
	}

	// 4. Generate clean title from filename
	return generateTitleFromFilename(filePath)
}

// stripHTMLTags removes HTML tags from a string
func stripHTMLTags(text string) string {
	// Remove all HTML tags
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			// Handle unclosed tags
			text = text[:start]
			break
		}
		text = text[:start] + " " + text[start+end+1:]
	}
	
	// Clean up whitespace
	text = strings.Join(strings.Fields(text), " ")
	return strings.TrimSpace(text)
}

// extractH1Title extracts title from the first H1 element
func extractH1Title(html string) string {
	// Look for <h1> tags
	start := strings.Index(strings.ToLower(html), "<h1")
	if start == -1 {
		return ""
	}

	// Find the end of the opening tag
	tagEnd := strings.Index(html[start:], ">")
	if tagEnd == -1 {
		return ""
	}
	contentStart := start + tagEnd + 1

	// Find the closing tag
	end := strings.Index(strings.ToLower(html[contentStart:]), "</h1>")
	if end == -1 {
		return ""
	}

	// Extract and clean the title content
	titleContent := html[contentStart : contentStart+end]

	// Remove any nested HTML tags and clean up
	title := stripHTMLTags(titleContent)
	return strings.TrimSpace(title)
}

// generateTitleFromFilename creates a clean title from filename
func generateTitleFromFilename(filePath string) string {
	// Get just the filename without extension
	filename := filepath.Base(filePath)
	filename = strings.TrimSuffix(filename, ".html")

	// Replace hyphens and underscores with spaces
	title := strings.ReplaceAll(filename, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Capitalize each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// indexMarkdownContent indexes all markdown files in the content directory
func indexMarkdownContent(items *[]SearchItem) error {

	return filepath.WalkDir("content", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Split frontmatter and content
		frontmatter, markdownContent := splitFrontmatter(string(content))

		// Parse frontmatter
		var fm SearchFrontmatter
		if frontmatter != "" {
			yaml.Unmarshal([]byte(frontmatter), &fm)
		}

		// Convert file path to URL with cleaning
		url := "/" + strings.TrimSuffix(path, ".md")
		url = strings.ReplaceAll(url, "content/", "")

		// Clean URL parts: remove number prefixes and normalize
		urlParts := strings.Split(url, "/")
		for i, part := range urlParts {
			if part != "" {
				// Remove number prefixes like "1.", "4.", "11."
				re := regexp.MustCompile(`^\d+\.?\s*`)
				part = re.ReplaceAllString(part, "")
				// Replace spaces with hyphens for clean URLs
				part = strings.ReplaceAll(part, " ", "-")
				// Convert to lowercase
				part = strings.ToLower(part)
				urlParts[i] = part
			}
		}
		url = strings.Join(urlParts, "/")

		// Determine section from path
		section := strings.Split(strings.TrimPrefix(path, "content/"), "/")[0]

		// Use frontmatter title or generate from filename
		title := fm.Title
		if title == "" {
			// Try to extract title from markdown H1
			if strings.HasPrefix(strings.TrimSpace(markdownContent), "# ") {
				lines := strings.Split(markdownContent, "\n")
				if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
					title = strings.TrimPrefix(lines[0], "# ")
					title = strings.TrimSpace(title)
				}
			}
			
			// If still no title, generate from filename
			if title == "" {
				// Clean filename: remove numbers and file extension
				title = filepath.Base(path)
				title = strings.TrimSuffix(title, ".md")
				// Remove number prefixes like "1.", "4.", "11."
				re := regexp.MustCompile(`^\d+\.?\s*`)
				title = re.ReplaceAllString(title, "")
				// Replace hyphens and underscores with spaces
				title = strings.ReplaceAll(title, "-", " ")
				title = strings.ReplaceAll(title, "_", " ")
				// Clean up multiple spaces
				title = strings.Join(strings.Fields(title), " ")
			}
		}

		*items = append(*items, SearchItem{
			Title:       title,
			Description: fm.Description,
			Content:     markdownContent,
			URL:         url,
			Type:        "content",
			Section:     section,
		})

		return nil
	})
}

// writeSearchIndex writes the search index to a JSON file
func writeSearchIndex(items []SearchItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("public/searchIndex.json", data, 0644)
}

// splitFrontmatter separates YAML frontmatter from markdown content
func splitFrontmatter(content string) (string, string) {
	// Normalize line endings to handle both \n and \r\n
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	if !strings.HasPrefix(content, "---\n") {
		return "", content
	}

	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return "", content
	}

	return parts[0], strings.TrimSpace(parts[1])
}

// extractHTMLTitle extracts title from HTML content (basic implementation)
func extractHTMLTitle(html string) string {
	start := strings.Index(html, "<title>")
	if start == -1 {
		return ""
	}
	start += 7

	end := strings.Index(html[start:], "</title>")
	if end == -1 {
		return ""
	}

	return html[start : start+end]
}

// indexCachedMarkdownContent indexes pre-rendered markdown content from cache
func indexCachedMarkdownContent(items *[]SearchItem, markdownService *MarkdownService) error {
	cachedContent := markdownService.GetAllCachedContent()

	for urlPath, content := range cachedContent {
		// Determine section from URL path
		section := ""
		if strings.HasPrefix(urlPath, "/docs") {
			section = "docs"
		} else if strings.HasPrefix(urlPath, "/api") {
			section = "api-docs"
		} else if strings.HasPrefix(urlPath, "/legal") {
			section = "legal"
		} else {
			// Extract first part of URL as section
			parts := strings.Split(strings.Trim(urlPath, "/"), "/")
			if len(parts) > 0 && parts[0] != "" {
				section = parts[0]
			}
		}

		// Get title from frontmatter or generate from URL
		title := ""
		description := ""
		category := ""
		if content.Frontmatter != nil {
			title = content.Frontmatter.Title
			description = content.Frontmatter.Description
		}

		// Set category based on URL path
		if strings.HasPrefix(urlPath, "/docs") {
			category = "Docs"
		} else if strings.HasPrefix(urlPath, "/api") {
			category = "API"
		} else if strings.HasPrefix(urlPath, "/legal") {
			category = "Legal"
		} else if strings.HasPrefix(urlPath, "/insights") {
			category = "Insights"
		}

		if title == "" {
			// Generate title from URL path
			title = generateTitleFromURL(urlPath)
		}

		// Extract clean text from pre-rendered HTML
		textContent := extractTextFromHTML(content.HTML)

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Content:     textContent,
			URL:         urlPath,
			Type:        "content",
			Section:     section,
			Category:    category,
		})
	}

	return nil
}

// generateTitleFromURL creates a clean title from URL path
func generateTitleFromURL(urlPath string) string {
	// Remove leading slash and get last segment
	cleanPath := strings.Trim(urlPath, "/")
	
	// Handle root/empty URLs
	if cleanPath == "" {
		return "Home"
	}
	
	parts := strings.Split(cleanPath, "/")

	// Use last segment as title
	lastSegment := parts[len(parts)-1]
	
	// Clean up the segment
	title := strings.ReplaceAll(lastSegment, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	
	// Capitalize each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// indexCachedHTMLPages indexes pre-rendered HTML pages from cache
func indexCachedHTMLPages(items *[]SearchItem, htmlService *HTMLService, metadata *Metadata) error {
	cachedContent := htmlService.GetAllCachedContent()

	for cacheKey, content := range cachedContent {
		// Cache key format is "lang:path", e.g., "es:/company/about"
		parts := strings.SplitN(cacheKey, ":", 2)
		if len(parts) != 2 {
			continue // Skip invalid cache keys
		}
		
		lang := parts[0]
		urlPath := parts[1]
		
		// Build the actual URL with language prefix
		actualURL := urlPath
		if lang != DefaultLanguage && lang != "" {
			actualURL = "/" + lang + urlPath
		}
		// Extract title from rendered HTML
		title := ""
		if metadata != nil {
			pageKey := getPageKey(urlPath)
			if pageData, exists := metadata.Pages[pageKey]; exists {
				// Use the language from the cache key
				if langMeta, langExists := pageData[lang]; langExists && langMeta.Title != "" {
					title = langMeta.Title
				} else if enMeta, enExists := pageData["en"]; enExists && enMeta.Title != "" {
					// Fallback to English
					title = enMeta.Title
				}
			}
		}
		// If no title from metadata, extract from HTML
		if title == "" {
			title = extractPageTitle(content.HTML, actualURL, content.FilePath, metadata)
		}

		// Get description from metadata if available
		description := ""
		if metadata != nil {
			pageKey := getPageKey(urlPath)
			
			if pageData, exists := metadata.Pages[pageKey]; exists {
				// Use the language from the cache key
				if langMeta, langExists := pageData[lang]; langExists && langMeta.Description != "" {
					description = langMeta.Description
				} else if enMeta, enExists := pageData["en"]; enExists && enMeta.Description != "" {
					// Fallback to English
					description = enMeta.Description
				}
			}
		}

		// Extract clean text from pre-rendered HTML
		textContent := extractTextFromHTML(content.HTML)

		// Determine section/type
		pageType := "page"
		section := ""
		if strings.HasPrefix(urlPath, "/platform") {
			section = "platform"
		} else if strings.HasPrefix(urlPath, "/company") {
			section = "company"
		}

		// Determine category based on URL path
		category := ""
		if strings.HasPrefix(urlPath, "/platform/features") {
			category = "Feature"
		} else if strings.HasPrefix(urlPath, "/solutions") {
			category = "Solution"
		}

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Content:     textContent,
			URL:         actualURL,  // Use the actual URL with language prefix
			Type:        pageType,
			Section:     section,
			Category:    category,
		})
	}

	return nil
}

// extractTextFromHTML extracts clean text starting from the first H1 tag and excludes script content
func extractTextFromHTML(html string) string {
	// First, completely remove all script content (including everything between <script> and </script>)
	text := html
	for {
		start := strings.Index(strings.ToLower(text), "<script")
		if start == -1 {
			break
		}
		end := strings.Index(strings.ToLower(text[start:]), "</script>")
		if end == -1 {
			// Handle unclosed script tags by removing everything after <script
			text = text[:start]
			break
		}
		// Calculate the actual end position and ensure it doesn't go out of bounds
		actualEnd := start + end + 9
		if actualEnd > len(text) {
			text = text[:start]
		} else {
			text = text[:start] + " " + text[actualEnd:]
		}
	}

	// Remove style content completely
	for {
		start := strings.Index(strings.ToLower(text), "<style")
		if start == -1 {
			break
		}
		end := strings.Index(strings.ToLower(text[start:]), "</style>")
		if end == -1 {
			text = text[:start]
			break
		}
		// Calculate the actual end position and ensure it doesn't go out of bounds
		actualEnd := start + end + 8
		if actualEnd > len(text) {
			text = text[:start]
		} else {
			text = text[:start] + " " + text[actualEnd:]
		}
	}

	// Find the first H1 tag and start content extraction from there
	h1Start := strings.Index(strings.ToLower(text), "<h1")
	if h1Start == -1 {
		// No H1 found, return empty string as per requirement
		return ""
	}

	// Start text extraction from the H1 tag onwards
	text = text[h1Start:]

	// Remove common problematic elements
	problematicTags := []string{
		"<noscript", "</noscript>",
		"<svg", "</svg>",
		"<path", "</path>",
		"<head", "</head>",
		"<meta", "</meta>",
		"<link", "</link>",
	}
	
	for i := 0; i < len(problematicTags); i += 2 {
		openTag := problematicTags[i]
		closeTag := problematicTags[i+1]
		
		for {
			start := strings.Index(strings.ToLower(text), openTag)
			if start == -1 {
				break
			}
			end := strings.Index(strings.ToLower(text[start:]), closeTag)
			if end == -1 {
				text = text[:start]
				break
			}
			text = text[:start] + " " + text[start+end+len(closeTag):]
		}
	}

	// Remove all remaining HTML tags
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			// Handle unclosed tags
			text = text[:start]
			break
		}
		text = text[:start] + " " + text[start+end+1:]
	}

	// Remove AlpineJS directives and JavaScript patterns
	jsPatterns := []string{
		"x-data=",
		"x-show=",
		"x-if=",
		"x-for=",
		"x-on:",
		"@click",
		"@keydown",
		"handleKeydown",
		"showResults",
		"hideResults",
		"function(",
		"if (",
		"for (",
		"while (",
		"return ",
		"const ",
		"let ",
		"var ",
		"===",
		"!==",
		"&&",
		"||",
		"addEventListener",
		"querySelector",
		"getElementById",
		"console.log",
		"window.",
		"document.",
	}

	for _, pattern := range jsPatterns {
		text = strings.ReplaceAll(text, pattern, " ")
	}

	// Clean up multiple spaces and special characters
	text = strings.ReplaceAll(text, "{", " ")
	text = strings.ReplaceAll(text, "}", " ")
	text = strings.ReplaceAll(text, "(", " ")
	text = strings.ReplaceAll(text, ")", " ")
	text = strings.ReplaceAll(text, "[", " ")
	text = strings.ReplaceAll(text, "]", " ")
	// Don't remove semicolons as they're part of HTML entities like &amp;
	text = strings.ReplaceAll(text, ":", " ")
	text = strings.ReplaceAll(text, "...", " ")
	text = strings.ReplaceAll(text, "=>", " ")
	text = strings.ReplaceAll(text, "=", " ")

	// Clean up whitespace and return
	text = strings.Join(strings.Fields(text), " ")
	return strings.TrimSpace(text)
}
