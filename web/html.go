package web

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// HTMLService handles HTML page pre-rendering
type HTMLService struct {
	cache           *HTMLCache
	pagesDir        string
	layoutsDir      string
	componentsDir   string
	markdownService *MarkdownService
	schemaService   *SchemaService
}

// NewHTMLService creates a new HTML service
func NewHTMLService(pagesDir, layoutsDir, componentsDir string, markdownService *MarkdownService) *HTMLService {
	return &HTMLService{
		cache:           NewHTMLCache(),
		pagesDir:        pagesDir,
		layoutsDir:      layoutsDir,
		componentsDir:   componentsDir,
		markdownService: markdownService,
		schemaService:   nil,
	}
}

// SetSchemaService sets the schema service for the HTML service
func (hs *HTMLService) SetSchemaService(schemaService *SchemaService) {
	hs.schemaService = schemaService
}

// getCacheKey generates a language-specific cache key
func (hs *HTMLService) getCacheKey(lang, path string) string {
	return lang + ":" + path
}

// PreRenderAllHTMLPages pre-renders all HTML pages in the pages directory
func (hs *HTMLService) PreRenderAllHTMLPages(navigationService *NavigationService, seoService *SEOService) error {
	count := 0

	// List of pages to exclude from pre-rendering (dynamic content)
	excludedPages := []string{
		"/platform/status", // Dynamic status page (truly dynamic - status changes)
		// Note: /insights is now pre-rendered with insights data baked in
	}

	// Walk through all HTML files in pages directory
	err := filepath.WalkDir(hs.pagesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non-HTML files
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Generate URL path for this file
		urlPath := hs.generateURLPath(path)

		// Check if this page should be excluded
		if hs.isExcluded(urlPath, excludedPages) {
			// Skipping excluded page
			return nil
		}

		// Get file info for modification time
		info, err := d.Info()
		if err != nil {
			log.Printf("Warning: could not get file info for %s: %v", path, err)
			return nil // Continue processing other files
		}

		// Pre-render the HTML page for each supported language
		for _, lang := range SupportedLanguages {
			// Pre-render the HTML page with the specific language
			html, err := hs.renderHTMLPageWithLang(path, urlPath, navigationService, seoService, lang)
			if err != nil {
				log.Printf("Warning: failed to pre-render %s for language %s: %v", path, lang, err)
				continue // Continue with next language
			}

			// Cache the pre-rendered content with language-specific key
			cachedContent := &CachedContent{
				HTML:        html,
				Frontmatter: nil, // HTML pages don't have frontmatter
				ModTime:     info.ModTime(),
				FilePath:    path,
			}

			cacheKey := hs.getCacheKey(lang, urlPath)
			hs.cache.Set(cacheKey, cachedContent)
			count++
		}

		// Progress tracking removed for cleaner output

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk pages directory: %w", err)
	}

	// Pre-rendering complete
	return nil
}

// renderHTMLPageWithLang renders a single HTML page with templates for a specific language
func (hs *HTMLService) renderHTMLPageWithLang(filePath, urlPath string, navigationService *NavigationService, seoService *SEOService, lang string) (string, error) {
	// Read the HTML file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Prepare page data with the specified language
	pageData := hs.preparePageData(urlPath, "", false, nil, navigationService, seoService, lang)

	// Create template for page content with the specified language
	contentTmpl := template.New("page-content").Funcs(getTemplateFuncs(lang))

	// Load component templates
	componentFiles, err := hs.loadComponentTemplates()
	if err != nil {
		return "", fmt.Errorf("failed to load components: %w", err)
	}

	// Parse component files first
	if len(componentFiles) > 0 {
		contentTmpl, err = contentTmpl.ParseFiles(componentFiles...)
		if err != nil {
			return "", fmt.Errorf("failed to parse components: %w", err)
		}
	}

	// Parse the page content
	contentTmpl, err = contentTmpl.Parse(string(contentBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse page template: %w", err)
	}

	// Execute the page template
	var renderedContent strings.Builder
	if err = contentTmpl.Execute(&renderedContent, pageData); err != nil {
		return "", fmt.Errorf("failed to execute page template: %w", err)
	}

	// Now render with main layout
	templateFiles := []string{
		filepath.Join(hs.layoutsDir, "main.html"),
	}
	templateFiles = append(templateFiles, componentFiles...)

	// Create main template with the specified language
	mainTmpl := template.New("main.html").Funcs(getTemplateFuncs(lang))
	mainTmpl, err = mainTmpl.ParseFiles(templateFiles...)
	if err != nil {
		return "", fmt.Errorf("failed to parse main template: %w", err)
	}

	// Prepare final page data with rendered content
	finalPageData := hs.preparePageData(urlPath, template.HTML(renderedContent.String()), false, nil, navigationService, seoService, lang)

	// Execute main template
	var finalHTML strings.Builder
	if err := mainTmpl.ExecuteTemplate(&finalHTML, "main.html", finalPageData); err != nil {
		return "", fmt.Errorf("failed to execute main template: %w", err)
	}

	return finalHTML.String(), nil
}

// loadComponentTemplates loads all component template files
func (hs *HTMLService) loadComponentTemplates() ([]string, error) {
	componentFiles, err := filepath.Glob(filepath.Join(hs.componentsDir, "*.html"))
	if err != nil {
		return nil, err
	}
	return componentFiles, nil
}

// generateURLPath converts a file path to a clean URL path
func (hs *HTMLService) generateURLPath(filePath string) string {
	// Remove pages/ prefix and .html suffix
	urlPath := strings.TrimPrefix(filePath, hs.pagesDir+"/")
	urlPath = strings.TrimSuffix(urlPath, ".html")

	// Handle index files
	if strings.HasSuffix(urlPath, "/index") {
		urlPath = strings.TrimSuffix(urlPath, "/index")
	}

	// Add leading slash
	if urlPath == "" {
		return "/"
	}
	return "/" + urlPath
}

// isExcluded checks if a URL path should be excluded from pre-rendering
func (hs *HTMLService) isExcluded(urlPath string, excludedPages []string) bool {
	for _, excluded := range excludedPages {
		if urlPath == excluded {
			return true
		}
	}
	return false
}

// GetCachedContent retrieves pre-rendered HTML content from cache
func (hs *HTMLService) GetCachedContent(urlPath string) (*CachedContent, bool) {
	return hs.cache.Get(urlPath)
}

// GetCachedContentForLang retrieves pre-rendered HTML content from cache for a specific language
func (hs *HTMLService) GetCachedContentForLang(urlPath, lang string) (*CachedContent, bool) {
	cacheKey := hs.getCacheKey(lang, urlPath)
	return hs.cache.Get(cacheKey)
}

// GetAllCachedContent returns all cached HTML content (for search indexing)
func (hs *HTMLService) GetAllCachedContent() map[string]*CachedContent {
	return hs.cache.GetAll()
}

// GetCacheSize returns the number of cached HTML items
func (hs *HTMLService) GetCacheSize() int {
	return hs.cache.Size()
}

// preparePageData creates PageData with metadata for the given path
// This is a wrapper around the shared page data preparation logic
func (hs *HTMLService) preparePageData(path string, content template.HTML, isMarkdown bool, frontmatter *Frontmatter, navigationService *NavigationService, seoService *SEOService, lang string) PageData {
	// Create a temporary router instance to use its preparePageData method
	// This avoids code duplication while maintaining the existing API
	tempRouter := &Router{
		markdownService:  hs.markdownService,
		seoService:       seoService,
		navigationService: navigationService,
		schemaService:    hs.schemaService,
		tocExcludedPaths: []string{"/changelog", "/roadmap", "/platform/status"},
	}
	
	return tempRouter.preparePageData(path, content, isMarkdown, frontmatter, navigationService.GetNavigationForPath(path), lang)
}
