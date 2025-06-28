package web

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownService handles markdown processing
type MarkdownService struct {
	markdown goldmark.Markdown
	cache    *MarkdownCache
}

// NewMarkdownService creates a new markdown service
func NewMarkdownService() *MarkdownService {
	// Configure Goldmark with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			NewYouTubeExtension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // Allow raw HTML (including video tags)
		),
	)

	return &MarkdownService{
		markdown: md,
		cache:    NewMarkdownCache(),
	}
}

// Convert converts markdown content to HTML
func (ms *MarkdownService) Convert(markdownContent []byte) (string, error) {
	var htmlBuffer strings.Builder
	if err := ms.markdown.Convert(markdownContent, &htmlBuffer); err != nil {
		return "", err
	}
	return htmlBuffer.String(), nil
}

// ProcessMarkdownFile reads a markdown file, parses frontmatter, and converts to HTML
func (ms *MarkdownService) ProcessMarkdownFile(filePath string, seoService *SEOService) (string, *Frontmatter, error) {
	// Read file
	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	// Parse frontmatter
	var markdownContent []byte
	frontmatter, markdownContent, err := seoService.ParseFrontmatter(mdBytes)
	if err != nil {
		// Continue without frontmatter
		markdownContent = mdBytes
	}

	// Convert to HTML
	html, err := ms.Convert(markdownContent)
	if err != nil {
		return "", nil, err
	}

	return html, frontmatter, nil
}

// PreRenderAllMarkdown pre-renders all markdown files in the content directory
func (ms *MarkdownService) PreRenderAllMarkdown(contentService *ContentService, seoService *SEOService) error {
	startTime := time.Now()
	count := 0

	// Walk through all markdown files in content directory
	err := filepath.WalkDir("content", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non-markdown files
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Get file info for modification time
		info, err := d.Info()
		if err != nil {
			log.Printf("Warning: could not get file info for %s: %v", path, err)
			return nil // Continue processing other files
		}

		// Generate URL path for this file
		urlPath := ms.generateURLPath(path)
		
		// Debug log for insights files
		if strings.Contains(path, "insights") {
			log.Printf("DEBUG: Processing insights file: %s → URL: %s", path, urlPath)
		}

		// Process the markdown file
		html, frontmatter, err := ms.ProcessMarkdownFile(path, seoService)
		if err != nil {
			log.Printf("Warning: failed to process %s: %v", path, err)
			return nil // Continue processing other files
		}

		// Debug log for insights frontmatter
		if strings.Contains(path, "insights") && frontmatter != nil {
			log.Printf("DEBUG: Insights frontmatter - Title: %s, Category: %s, Tags: %v", 
				frontmatter.Title, frontmatter.Category, frontmatter.Tags)
		}

		// Cache the pre-rendered content
		cachedContent := &CachedContent{
			HTML:        html,
			Frontmatter: frontmatter,
			ModTime:     info.ModTime(),
			FilePath:    path,
		}

		ms.cache.Set(urlPath, cachedContent)
		count++

		if count%10 == 0 {
			log.Printf("Pre-rendered %d markdown files...", count)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk content directory: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("Pre-rendered %d markdown files in %v", count, duration)
	return nil
}

// generateURLPath converts a file path to a clean URL path
func (ms *MarkdownService) generateURLPath(filePath string) string {
	// Remove content/ prefix and .md suffix
	urlPath := strings.TrimPrefix(filePath, "content/")
	urlPath = strings.TrimSuffix(urlPath, ".md")

	// Handle index files
	if strings.HasSuffix(urlPath, "/index") {
		urlPath = strings.TrimSuffix(urlPath, "/index")
	}

	// Clean URL parts: remove number prefixes and normalize
	urlParts := strings.Split(urlPath, "/")
	for i, part := range urlParts {
		if part != "" {
			// Use the same cleaning logic as the utils
			urlParts[i] = CleanID(part)
		}
	}

	// Reconstruct URL with leading slash
	cleanURL := "/" + strings.Join(urlParts, "/")
	
	// Handle root case
	if cleanURL == "/" {
		return "/"
	}

	return cleanURL
}

// GetCachedContent retrieves pre-rendered content from cache
func (ms *MarkdownService) GetCachedContent(urlPath string) (*CachedContent, bool) {
	return ms.cache.Get(urlPath)
}

// GetAllCachedContent returns all cached content (for search indexing)
func (ms *MarkdownService) GetAllCachedContent() map[string]*CachedContent {
	return ms.cache.GetAll()
}

// GetCacheSize returns the number of cached items
func (ms *MarkdownService) GetCacheSize() int {
	return ms.cache.Size()
}