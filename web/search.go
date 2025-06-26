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
}

// SearchFrontmatter represents the YAML frontmatter in markdown files
type SearchFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// GenerateSearchIndex creates a JSON search index from all content
func GenerateSearchIndex() error {
	var searchItems []SearchItem

	// Index HTML pages
	if err := indexHTMLPages(&searchItems); err != nil {
		return fmt.Errorf("failed to index HTML pages: %w", err)
	}

	// Index markdown content
	if err := indexMarkdownContent(&searchItems); err != nil {
		return fmt.Errorf("failed to index markdown content: %w", err)
	}

	// Write search index to public directory
	return writeSearchIndex(searchItems)
}

// indexHTMLPages indexes all HTML pages in the pages directory
func indexHTMLPages(items *[]SearchItem) error {
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

		// Extract title from HTML (basic extraction)
		title := extractHTMLTitle(string(content))
		if title == "" {
			// Simple title generation from filename
			title = strings.ReplaceAll(filepath.Base(path), "-", " ")
		}

		// Extract text content (remove HTML tags)
		textContent := extractTextFromHTML(string(content))

		*items = append(*items, SearchItem{
			Title:   title,
			Content: textContent,
			URL:     url,
			Type:    "page",
		})

		return nil
	})
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

		// Convert file path to URL
		url := "/" + strings.TrimSuffix(path, ".md")
		url = strings.ReplaceAll(url, "content/", "")

		// Determine section from path
		section := strings.Split(strings.TrimPrefix(path, "content/"), "/")[0]

		// Use frontmatter title or generate from filename
		title := fm.Title
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

// extractTextFromHTML removes HTML tags to get plain text (basic implementation)
func extractTextFromHTML(html string) string {
	// Simple tag removal - for production you might want a proper HTML parser
	text := html
	
	// Remove script and style content
	for {
		start := strings.Index(text, "<script")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], "</script>")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+9:]
	}
	
	for {
		start := strings.Index(text, "<style")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], "</style>")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+8:]
	}
	
	// Remove all HTML tags
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + " " + text[start+end+1:]
	}
	
	// Clean up whitespace
	text = strings.Join(strings.Fields(text), " ")
	
	return strings.TrimSpace(text)
}