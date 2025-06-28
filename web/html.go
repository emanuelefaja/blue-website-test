package web

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HTMLService handles HTML page pre-rendering
type HTMLService struct {
	cache         *HTMLCache
	pagesDir      string
	layoutsDir    string
	componentsDir string
}

// NewHTMLService creates a new HTML service
func NewHTMLService(pagesDir, layoutsDir, componentsDir string) *HTMLService {
	return &HTMLService{
		cache:         NewHTMLCache(),
		pagesDir:      pagesDir,
		layoutsDir:    layoutsDir,
		componentsDir: componentsDir,
	}
}

// PreRenderAllHTMLPages pre-renders all HTML pages in the pages directory
func (hs *HTMLService) PreRenderAllHTMLPages(navigationService *NavigationService, seoService *SEOService, changelogService *ChangelogService) error {
	startTime := time.Now()
	count := 0

	// List of pages to exclude from pre-rendering (dynamic content)
	excludedPages := []string{
		"/platform/status", // Dynamic status page
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
			log.Printf("Skipping pre-rendering for excluded page: %s", urlPath)
			return nil
		}

		// Get file info for modification time
		info, err := d.Info()
		if err != nil {
			log.Printf("Warning: could not get file info for %s: %v", path, err)
			return nil // Continue processing other files
		}

		// Pre-render the HTML page
		html, err := hs.renderHTMLPage(path, urlPath, navigationService, seoService, changelogService)
		if err != nil {
			log.Printf("Warning: failed to pre-render %s: %v", path, err)
			return nil // Continue processing other files
		}

		// Cache the pre-rendered content
		cachedContent := &CachedContent{
			HTML:        html,
			Frontmatter: nil, // HTML pages don't have frontmatter
			ModTime:     info.ModTime(),
			FilePath:    path,
		}

		hs.cache.Set(urlPath, cachedContent)
		count++

		if count%5 == 0 {
			log.Printf("Pre-rendered %d HTML pages...", count)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk pages directory: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("Pre-rendered %d HTML pages in %v", count, duration)
	return nil
}

// renderHTMLPage renders a single HTML page with templates  
func (hs *HTMLService) renderHTMLPage(filePath, urlPath string, navigationService *NavigationService, seoService *SEOService, changelogService *ChangelogService) (string, error) {
	// Read the HTML file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Prepare page data
	pageData := hs.preparePageData(urlPath, "", false, nil, navigationService, seoService, changelogService)

	// Create template for page content
	contentTmpl := template.New("page-content").Funcs(templateFuncs)

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

	// Create main template
	mainTmpl := template.New("main.html").Funcs(templateFuncs)
	mainTmpl, err = mainTmpl.ParseFiles(templateFiles...)
	if err != nil {
		return "", fmt.Errorf("failed to parse main template: %w", err)
	}

	// Prepare final page data with rendered content
	finalPageData := hs.preparePageData(urlPath, template.HTML(renderedContent.String()), false, nil, navigationService, seoService, changelogService)

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

// GetAllCachedContent returns all cached HTML content (for search indexing)
func (hs *HTMLService) GetAllCachedContent() map[string]*CachedContent {
	return hs.cache.GetAll()
}

// GetCacheSize returns the number of cached HTML items
func (hs *HTMLService) GetCacheSize() int {
	return hs.cache.Size()
}

// preparePageData creates PageData with metadata for the given path
func (hs *HTMLService) preparePageData(path string, content template.HTML, isMarkdown bool, frontmatter *Frontmatter, navigationService *NavigationService, seoService *SEOService, changelogService *ChangelogService) PageData {
	// Get metadata from SEO service
	title, description, keywords, pageMeta, siteMeta := seoService.PreparePageMetadata(path, isMarkdown, frontmatter)

	// Prepare changelog data - only include if on changelog page
	var changelog []ChangelogMonth
	if path == "/changelog" && changelogService != nil {
		changelog = changelogService.GetChangelog()
		log.Printf("Loading changelog for path=%s, found %d entries", path, len(changelog))
	} else {
		log.Printf("Not loading changelog: path=%s, service=%v", path, changelogService != nil)
	}

	// Extract table of contents
	toc := make([]TOCEntry, 0)
	tocExcludedPaths := []string{"/changelog", "/roadmap", "/platform/status"}
	
	isExcluded := false
	for _, excludedPath := range tocExcludedPaths {
		if path == excludedPath {
			isExcluded = true
			break
		}
	}
	
	if !isExcluded && string(content) != "" {
		var err error
		if isMarkdown {
			// For markdown content, extract H2 elements from converted HTML
			toc, err = ExtractH2TOC(string(content))
		} else {
			// For HTML pages, extract from section elements
			toc, err = ExtractHTMLTOC(string(content))
		}

		if err != nil {
			log.Printf("Error extracting TOC for path=%s: %v", path, err)
		} else {
			log.Printf("Extracted %d TOC entries for path=%s", len(toc), path)
		}
	} else if isExcluded {
		log.Printf("TOC generation skipped for excluded path: %s", path)
	}

	// Return PageData with all components
	return PageData{
		Title:          title,
		Content:        content,
		Navigation:     navigationService.GetNavigationForPath(path),
		PageMeta:       pageMeta,
		SiteMeta:       siteMeta,
		Description:    description,
		Keywords:       keywords,
		IsMarkdown:     isMarkdown,
		Frontmatter:    frontmatter,
		Changelog:      changelog,
		TOC:            toc,
		CustomerNumber: 17000,
	}
}