package web

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
)

// SearchItem represents a single searchable document
type SearchItem struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords"` // Extracted keywords instead of full content
	URL         string   `json:"url"`
	Type        string   `json:"type"` // "page", "doc", "blog", etc.
	Section     string   `json:"section,omitempty"`
	Category    string   `json:"category,omitempty"` // Display category like "Feature", "Docs", "API"
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

// languageStats tracks timing for each language's processing
type languageStats struct {
	itemCount    int
	collectTime  time.Duration
	marshalTime  time.Duration
	writeTime    time.Duration
	totalTime    time.Duration
}

// GenerateSearchIndexWithCaches creates a JSON search index using both HTML and markdown caches
func GenerateSearchIndexWithCaches(markdownService *MarkdownService, htmlService *HTMLService) error {
	// Load metadata.json for title lookup
	metadataStart := time.Now()
	metadata, err := loadMetadata()
	if err != nil {
		fmt.Printf("Warning: Could not load metadata.json: %v\n", err)
	}
	fmt.Printf("   📊 Metadata loading: %v\n", time.Since(metadataStart))

	// Generate search index for each language in parallel
	type languageResult struct {
		lang string
		err  error
		stats languageStats
	}

	resultChan := make(chan languageResult, len(SupportedLanguages))

	// Process each language concurrently with true independence
	parallelStart := time.Now()
	for _, lang := range SupportedLanguages {
		go func(lang string) {
			langStart := time.Now()
			stats := languageStats{}
			
			// Collect items for this language only
			collectStart := time.Now()
			items, err := collectSearchItemsForLanguage(markdownService, htmlService, metadata, lang)
			stats.collectTime = time.Since(collectStart)
			stats.itemCount = len(items)
			
			if err != nil {
				resultChan <- languageResult{lang, err, stats}
				return
			}

			// Generate index file directly
			filename := "public/searchIndex.json"
			if lang != DefaultLanguage {
				filename = fmt.Sprintf("public/searchIndex-%s.json", lang)
			}

			marshalStart := time.Now()
			data, err := json.Marshal(items)
			stats.marshalTime = time.Since(marshalStart)
			
			if err != nil {
				resultChan <- languageResult{lang, fmt.Errorf("failed to marshal search index for language %s: %w", lang, err), stats}
				return
			}

			writeStart := time.Now()
			if err := os.WriteFile(filename, data, 0644); err != nil {
				resultChan <- languageResult{lang, fmt.Errorf("failed to write search index for language %s: %w", lang, err), stats}
				return
			}
			stats.writeTime = time.Since(writeStart)
			stats.totalTime = time.Since(langStart)

			resultChan <- languageResult{lang, nil, stats}
		}(lang)
	}

	// Wait for all languages to complete and check for errors
	var totalItems int
	var maxCollectTime, maxMarshalTime, maxWriteTime time.Duration
	
	for i := 0; i < len(SupportedLanguages); i++ {
		result := <-resultChan
		if result.err != nil {
			return result.err
		}
		
		// Track statistics
		totalItems += result.stats.itemCount
		if result.stats.collectTime > maxCollectTime {
			maxCollectTime = result.stats.collectTime
		}
		if result.stats.marshalTime > maxMarshalTime {
			maxMarshalTime = result.stats.marshalTime
		}
		if result.stats.writeTime > maxWriteTime {
			maxWriteTime = result.stats.writeTime
		}
	}
	
	fmt.Printf("   ⏱️  Search index breakdown:\n")
	fmt.Printf("      - Collection (max): %v\n", maxCollectTime)
	fmt.Printf("      - JSON Marshal (max): %v\n", maxMarshalTime)
	fmt.Printf("      - File Write (max): %v\n", maxWriteTime)
	fmt.Printf("      - Total items indexed: %d\n", totalItems)
	fmt.Printf("      - Parallel processing: %v\n", time.Since(parallelStart))

	return nil
}

// searchTask represents a single item to be indexed
type searchTask struct {
	urlPath string
	content *CachedContent
	isHTML  bool
}

// processSearchTask processes a single search task into a SearchItem
func processSearchTask(task searchTask, metadata *Metadata, lang string, extractor *KeywordExtractor) (SearchItem, error) {
	if task.isHTML {
		return processHTMLSearchTask(task, metadata, lang, extractor)
	}
	return processMarkdownSearchTask(task, lang, extractor)
}

// processHTMLSearchTask processes an HTML content task
func processHTMLSearchTask(task searchTask, metadata *Metadata, lang string, extractor *KeywordExtractor) (SearchItem, error) {
	urlPath := task.urlPath
	content := task.content
	
	// Build the actual URL with language prefix if not default
	actualURL := urlPath
	if lang != DefaultLanguage && lang != "" {
		actualURL = "/" + lang + urlPath
	}

	// Extract title from rendered HTML
	title := ""
	if metadata != nil {
		pageKey := getPageKey(urlPath)
		if pageData, exists := metadata.Pages[pageKey]; exists {
			// Use the specified language
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
		title = extractPageTitleWithLang(content.HTML, actualURL, content.FilePath, metadata, lang)
	}

	// Get description from metadata if available
	description := ""
	if metadata != nil {
		pageKey := getPageKey(urlPath)

		if pageData, exists := metadata.Pages[pageKey]; exists {
			// Use the specified language
			if langMeta, langExists := pageData[lang]; langExists && langMeta.Description != "" {
				description = langMeta.Description
			} else if enMeta, enExists := pageData["en"]; enExists && enMeta.Description != "" {
				// Fallback to English
				description = enMeta.Description
			}
		}
	}

	// Extract keywords from content
	var keywords []string
	if len(content.HTML) < 500000 { // Skip extraction for pages larger than 500KB
		textContent := extractTextFromHTML(content.HTML)
		// Combine title, description, and content for keyword extraction
		fullText := title + " " + description + " " + textContent
		keywords = extractor.ExtractKeywords(fullText, 80)
	}

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

	return SearchItem{
		Title:       title,
		Description: description,
		Keywords:    keywords,
		URL:         actualURL,
		Type:        pageType,
		Section:     section,
		Category:    category,
	}, nil
}

// processMarkdownSearchTask processes a markdown content task
func processMarkdownSearchTask(task searchTask, lang string, extractor *KeywordExtractor) (SearchItem, error) {
	urlPath := task.urlPath
	content := task.content
	
	// Build actual URL with language prefix if not default
	actualURL := urlPath
	if lang != DefaultLanguage {
		actualURL = "/" + lang + urlPath
	}

	// Determine section from URL path (using original urlPath, not actualURL)
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

	// Extract keywords from content
	var keywords []string
	if len(content.HTML) < 500000 { // Skip extraction for pages larger than 500KB
		textContent := extractTextFromHTML(content.HTML)
		// Combine title, description, and content for keyword extraction
		fullText := title + " " + description + " " + textContent
		keywords = extractor.ExtractKeywords(fullText, 80)
	}

	return SearchItem{
		Title:       title,
		Description: description,
		Keywords:    keywords,
		URL:         actualURL,
		Type:        "content",
		Section:     section,
		Category:    category,
	}, nil
}

// collectSearchItemsForLanguage collects search items for a specific language only
func collectSearchItemsForLanguage(markdownService *MarkdownService, htmlService *HTMLService, metadata *Metadata, lang string) ([]SearchItem, error) {
	// Create a shared keyword extractor for better performance
	extractor := NewKeywordExtractor()
	
	// Collect all tasks for this language
	var tasks []searchTask
	
	// Get cached HTML pages
	cacheStart := time.Now()
	htmlContent := htmlService.GetCachedContentByLanguage(lang)
	htmlCacheTime := time.Since(cacheStart)
	
	for urlPath, content := range htmlContent {
		tasks = append(tasks, searchTask{
			urlPath: urlPath,
			content: content,
			isHTML:  true,
		})
	}
	
	// Get cached markdown content
	cacheStart = time.Now()
	markdownContent := markdownService.GetCachedContentByLanguage(lang)
	mdCacheTime := time.Since(cacheStart)
	
	for urlPath, content := range markdownContent {
		tasks = append(tasks, searchTask{
			urlPath: urlPath,
			content: content,
			isHTML:  false,
		})
	}
	
	if lang == DefaultLanguage {
		fmt.Printf("      [%s] Cache retrieval: HTML %v, MD %v, Tasks: %d\n", lang, htmlCacheTime, mdCacheTime, len(tasks))
	}
	
	// Process tasks with worker pool
	processingStart := time.Now()
	const numWorkers = 30
	// Buffer channels to avoid blocking
	taskChan := make(chan searchTask, len(tasks))
	resultChan := make(chan SearchItem, len(tasks))
	errorChan := make(chan error, 1) // Only need to capture first error
	
	// Track timing for text extraction
	var extractionTime time.Duration
	var extractionCount int
	var extractionMutex sync.Mutex
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				// Time the text extraction specifically
				extractStart := time.Now()
				item, err := processSearchTask(task, metadata, lang, extractor)
				elapsed := time.Since(extractStart)
				
				extractionMutex.Lock()
				extractionTime += elapsed
				extractionCount++
				extractionMutex.Unlock()
				
				if err != nil {
					errorChan <- err
					continue
				}
				resultChan <- item
			}
		}()
	}
	
	// Send all tasks
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)
	
	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()
	
	// Collect results
	var searchItems []SearchItem
	for item := range resultChan {
		searchItems = append(searchItems, item)
	}
	
	// Check for errors
	select {
	case err := <-errorChan:
		if err != nil {
			return nil, err
		}
	default:
	}
	
	if lang == DefaultLanguage {
		avgExtraction := time.Duration(0)
		if extractionCount > 0 {
			avgExtraction = extractionTime / time.Duration(extractionCount)
		}
		fmt.Printf("      [%s] Processing: %v (avg extraction: %v × %d items)\n", lang, time.Since(processingStart), avgExtraction, extractionCount)
	}

	// Build a map of indexed URLs to avoid duplicates from non-cached pages
	indexedURLs := make(map[string]bool)
	for _, item := range searchItems {
		indexedURLs[item.URL] = true
	}

	// Handle non-cached HTML pages (like status page) for this language
	// Note: These are typically language-agnostic pages, so we only index them for the default language
	// to avoid duplication across language indexes
	if lang == DefaultLanguage {
		if err := indexHTMLPages(&searchItems, metadata, indexedURLs); err != nil {
			return nil, fmt.Errorf("failed to index non-cached HTML pages for language %s: %w", lang, err)
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
	// Create a keyword extractor for this function
	extractor := NewKeywordExtractor()

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
		url = strings.TrimSuffix(url, "/index")
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

		// Extract keywords from content
		fullText := title + " " + description + " " + textContent
		keywords := extractor.ExtractKeywords(fullText, 80)

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Keywords:    keywords,
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

// extractPageTitleWithLang implements the same title extraction algorithm but accepts
// language as a parameter to avoid redundant extractLanguageFromPath() calls
func extractPageTitleWithLang(htmlContent, url, filePath string, metadata *Metadata, lang string) string {
	// 1. Check metadata.json first
	if metadata != nil {
		pageKey := getPageKey(url)
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
	// Create a keyword extractor for this function
	extractor := NewKeywordExtractor()


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

		// Extract keywords from markdown content
		fullText := title + " " + fm.Description + " " + markdownContent
		keywords := extractor.ExtractKeywords(fullText, 80)

		*items = append(*items, SearchItem{
			Title:       title,
			Description: fm.Description,
			Keywords:    keywords,
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
	// Create a keyword extractor for this function
	extractor := NewKeywordExtractor()

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

		// Extract keywords from content
		fullText := title + " " + description + " " + textContent
		keywords := extractor.ExtractKeywords(fullText, 80)

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Keywords:    keywords,
			URL:         urlPath,
			Type:        "content",
			Section:     section,
			Category:    category,
		})
	}

	return nil
}

// indexCachedMarkdownContentForLanguage indexes pre-rendered markdown content for a specific language
func indexCachedMarkdownContentForLanguage(items *[]SearchItem, markdownService *MarkdownService, lang string) error {
	// Create a keyword extractor for this function
	extractor := NewKeywordExtractor()
	
	cachedContent := markdownService.GetCachedContentByLanguage(lang)

	for urlPath, content := range cachedContent {
		// Build actual URL with language prefix if not default
		actualURL := urlPath
		if lang != DefaultLanguage {
			actualURL = "/" + lang + urlPath
		}

		// Determine section from URL path (using original urlPath, not actualURL)
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

		// Extract keywords from content
		fullText := title + " " + description + " " + textContent
		keywords := extractor.ExtractKeywords(fullText, 80)

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Keywords:    keywords,
			URL:         actualURL,
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

// indexCachedHTMLPagesForLanguage indexes pre-rendered HTML pages for a specific language
func indexCachedHTMLPagesForLanguage(items *[]SearchItem, htmlService *HTMLService, metadata *Metadata, lang string) error {
	// Create a keyword extractor for this function
	extractor := NewKeywordExtractor()

	cachedContent := htmlService.GetCachedContentByLanguage(lang)

	for urlPath, content := range cachedContent {
		// Build the actual URL with language prefix if not default
		actualURL := urlPath
		if lang != DefaultLanguage && lang != "" {
			actualURL = "/" + lang + urlPath
		}

		// Extract title from rendered HTML
		title := ""
		if metadata != nil {
			pageKey := getPageKey(urlPath)
			if pageData, exists := metadata.Pages[pageKey]; exists {
				// Use the specified language
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
			title = extractPageTitleWithLang(content.HTML, actualURL, content.FilePath, metadata, lang)
		}

		// Get description from metadata if available
		description := ""
		if metadata != nil {
			pageKey := getPageKey(urlPath)

			if pageData, exists := metadata.Pages[pageKey]; exists {
				// Use the specified language
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

		// Extract keywords from content
		fullText := title + " " + description + " " + textContent
		keywords := extractor.ExtractKeywords(fullText, 80)

		*items = append(*items, SearchItem{
			Title:       title,
			Description: description,
			Keywords:    keywords,
			URL:         actualURL,
			Type:        pageType,
			Section:     section,
			Category:    category,
		})
	}

	return nil
}

// extractTextFromHTML extracts clean text from HTML using the x/net/html parser
func extractTextFromHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	
	var text strings.Builder
	var h1Found bool
	var inScript, inStyle bool
	
	// Walk the HTML tree and extract text
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		// Skip script and style content
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script":
				inScript = true
				defer func() { inScript = false }()
			case "style":
				inStyle = true
				defer func() { inStyle = false }()
			case "h1":
				h1Found = true
			}
		}
		
		// Only extract text after finding H1
		if h1Found && n.Type == html.TextNode && !inScript && !inStyle {
			// Clean the text
			cleaned := strings.TrimSpace(n.Data)
			if cleaned != "" {
				text.WriteString(cleaned)
				text.WriteString(" ")
			}
		}
		
		// Process children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	
	extractText(doc)
	
	// Normalize whitespace
	result := text.String()
	result = strings.Join(strings.Fields(result), " ")
	return strings.TrimSpace(result)
}
