package web

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// HTMLService handles HTML page pre-rendering
type HTMLService struct {
	cache           *HTMLCache
	pagesDir        string
	layoutsDir      string
	componentsDir   string
	markdownService *MarkdownService
	schemaService   *SchemaService
	statusChecker   *HealthChecker
	
	// Template cache for pre-compiled templates
	templateCache struct {
		sync.RWMutex
		contentTemplates map[string]*template.Template // [lang]->template
		mainTemplates    map[string]*template.Template // [lang]->template
		initialized      bool
	}
	
	// Page content cache for HTML files
	pageContentCache struct {
		sync.RWMutex
		content map[string][]byte // [filepath]->content
		loaded  bool
	}
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

// SetStatusChecker sets the status checker for the HTML service
func (hs *HTMLService) SetStatusChecker(statusChecker *HealthChecker) {
	hs.statusChecker = statusChecker
}

// precompileTemplates pre-compiles templates for all languages
func (hs *HTMLService) precompileTemplates() error {
	hs.templateCache.Lock()
	defer hs.templateCache.Unlock()
	
	if hs.templateCache.initialized {
		return nil
	}
	
	// Initialize maps
	hs.templateCache.contentTemplates = make(map[string]*template.Template)
	hs.templateCache.mainTemplates = make(map[string]*template.Template)
	
	// Load component files once
	componentFiles, err := hs.loadComponentTemplates()
	if err != nil {
		return fmt.Errorf("failed to load component templates: %w", err)
	}
	
	// Pre-compile for each language
	for _, lang := range SupportedLanguages {
		// Content template with components
		contentTmpl := template.New("page-content").Funcs(getTemplateFuncs(lang))
		if len(componentFiles) > 0 {
			contentTmpl, err = contentTmpl.ParseFiles(componentFiles...)
			if err != nil {
				return fmt.Errorf("failed to parse content components for %s: %w", lang, err)
			}
		}
		hs.templateCache.contentTemplates[lang] = contentTmpl
		
		// Main layout template with components
		templateFiles := []string{filepath.Join(hs.layoutsDir, "main.html")}
		templateFiles = append(templateFiles, componentFiles...)
		
		mainTmpl := template.New("main.html").Funcs(getTemplateFuncs(lang))
		mainTmpl, err = mainTmpl.ParseFiles(templateFiles...)
		if err != nil {
			return fmt.Errorf("failed to parse main template for %s: %w", lang, err)
		}
		hs.templateCache.mainTemplates[lang] = mainTmpl
	}
	
	hs.templateCache.initialized = true
	return nil
}

// preloadPageContent loads all HTML page content into memory cache
func (hs *HTMLService) preloadPageContent() error {
	hs.pageContentCache.Lock()
	defer hs.pageContentCache.Unlock()
	
	if hs.pageContentCache.loaded {
		return nil
	}
	
	hs.pageContentCache.content = make(map[string][]byte)
	
	// Walk through all HTML files in pages directory
	err := filepath.WalkDir(hs.pagesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// Skip non-HTML files
		if !strings.HasSuffix(path, ".html") {
			return nil
		}
		
		// Read file content once
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}
		
		// Cache the content
		hs.pageContentCache.content[path] = content
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to walk pages directory: %w", err)
	}
	
	hs.pageContentCache.loaded = true
	return nil
}

// getCachedPageContent retrieves cached page content or falls back to disk
func (hs *HTMLService) getCachedPageContent(filePath string) ([]byte, error) {
	hs.pageContentCache.RLock()
	defer hs.pageContentCache.RUnlock()
	
	content, ok := hs.pageContentCache.content[filePath]
	if !ok {
		// Fallback to reading from disk if not cached
		return os.ReadFile(filePath)
	}
	
	// Return a copy to prevent mutations
	return append([]byte(nil), content...), nil
}

// getCacheKey generates a language-specific cache key
func (hs *HTMLService) getCacheKey(lang, path string) string {
	return lang + ":" + path
}

// htmlTask represents an HTML file + language combination to be rendered
type htmlTask struct {
	path     string
	urlPath  string
	lang     string
	modTime  time.Time
}

// PreRenderAllHTMLPages pre-renders all HTML pages in the pages directory
func (hs *HTMLService) PreRenderAllHTMLPages(navigationService *NavigationService, seoService *SEOService) error {
	// List of pages to exclude from pre-rendering (dynamic content)
	excludedPages := []string{
		// Note: /platform/status is now pre-rendered with status data baked in
		// Note: /insights is now pre-rendered with insights data baked in
	}

	// Collect all HTML rendering tasks
	var htmlTasks []htmlTask

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

		// Create tasks for each language
		for _, lang := range SupportedLanguages {
			htmlTasks = append(htmlTasks, htmlTask{
				path:    path,
				urlPath: urlPath,
				lang:    lang,
				modTime: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk pages directory: %w", err)
	}

	// Preload all page content into memory
	if err := hs.preloadPageContent(); err != nil {
		return fmt.Errorf("failed to preload page content: %w", err)
	}

	// Pre-compile templates for all languages once
	if err := hs.precompileTemplates(); err != nil {
		return fmt.Errorf("failed to pre-compile templates: %w", err)
	}

	// Process all tasks in parallel using worker pool
	const numWorkers = 30
	taskChan := make(chan htmlTask, len(htmlTasks))
	resultChan := make(chan int, len(htmlTasks))
	errorChan := make(chan error, numWorkers)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				// Pre-render the HTML page with the specific language using optimized method
				html, err := hs.renderHTMLPageOptimized(task.path, task.urlPath, navigationService, seoService, task.lang)
				if err != nil {
					log.Printf("Warning: failed to pre-render %s for language %s: %v", task.path, task.lang, err)
					continue // Continue processing other tasks
				}

				// Detect if the HTML contains code blocks
				needsCodeHighlight := DetectCodeBlocks(html)

				// Cache the pre-rendered content with language-specific key
				cachedContent := &CachedContent{
					HTML:               html,
					Frontmatter:        nil, // HTML pages don't have frontmatter
					ModTime:            task.modTime,
					FilePath:           task.path,
					NeedsCodeHighlight: needsCodeHighlight,
				}

				cacheKey := hs.getCacheKey(task.lang, task.urlPath)
				hs.cache.Set(cacheKey, cachedContent)
				resultChan <- 1
			}
		}()
	}

	// Send all tasks to workers
	for _, task := range htmlTasks {
		taskChan <- task
	}
	close(taskChan)

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Count processed pages
	count := 0
	for range resultChan {
		count++
	}

	// Check for any errors
	select {
	case err := <-errorChan:
		if err != nil {
			return err
		}
	default:
	}

	// Pre-rendering complete
	return nil
}

// renderHTMLPageWithComponents renders a single HTML page with pre-loaded components
func (hs *HTMLService) renderHTMLPageWithComponents(filePath, urlPath string, navigationService *NavigationService, seoService *SEOService, lang string, componentFiles []string) (string, error) {
	// Read the HTML file
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Prepare page data with the specified language
	pageData := hs.preparePageData(urlPath, "", false, nil, navigationService, seoService, lang)

	// Create template for page content with the specified language
	contentTmpl := template.New("page-content").Funcs(getTemplateFuncs(lang))

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

// renderHTMLPageWithLang renders a single HTML page with templates for a specific language
func (hs *HTMLService) renderHTMLPageWithLang(filePath, urlPath string, navigationService *NavigationService, seoService *SEOService, lang string) (string, error) {
	// Load component templates
	componentFiles, err := hs.loadComponentTemplates()
	if err != nil {
		return "", fmt.Errorf("failed to load components: %w", err)
	}

	return hs.renderHTMLPageWithComponents(filePath, urlPath, navigationService, seoService, lang, componentFiles)
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
	urlPath = strings.TrimSuffix(urlPath, "/index")

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

// renderHTMLPageOptimized renders an HTML page using pre-compiled templates
func (hs *HTMLService) renderHTMLPageOptimized(
	filePath, urlPath string,
	navigationService *NavigationService,
	seoService *SEOService,
	lang string,
) (string, error) {
	// Get cached page content instead of reading from disk
	contentBytes, err := hs.getCachedPageContent(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	
	// Get pre-compiled templates for this language
	hs.templateCache.RLock()
	contentTmpl, contentOk := hs.templateCache.contentTemplates[lang]
	mainTmpl, mainOk := hs.templateCache.mainTemplates[lang]
	hs.templateCache.RUnlock()
	
	if !contentOk || !mainOk {
		return "", fmt.Errorf("templates not pre-compiled for language: %s", lang)
	}
	
	// Clone templates for thread safety
	contentTmpl, err = contentTmpl.Clone()
	if err != nil {
		return "", fmt.Errorf("failed to clone content template: %w", err)
	}
	mainTmpl, err = mainTmpl.Clone()
	if err != nil {
		return "", fmt.Errorf("failed to clone main template: %w", err)
	}
	
	// Parse the specific page content into the cloned template
	contentTmpl, err = contentTmpl.Parse(string(contentBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse page template: %w", err)
	}
	
	// Prepare page data
	pageData := hs.preparePageData(urlPath, "", false, nil, navigationService, seoService, lang)
	
	// Execute content template
	var renderedContent strings.Builder
	if err = contentTmpl.Execute(&renderedContent, pageData); err != nil {
		return "", fmt.Errorf("failed to execute page template: %w", err)
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

// GetCachedContentByLanguage returns all cached HTML content for a specific language
func (hs *HTMLService) GetCachedContentByLanguage(lang string) map[string]*CachedContent {
	return hs.cache.GetByLanguage(lang)
}

// GetCacheSize returns the number of cached HTML items
func (hs *HTMLService) GetCacheSize() int {
	return hs.cache.Size()
}

// RegenerateStatusPages regenerates the status page for all languages with fresh data
func (hs *HTMLService) RegenerateStatusPages(router *Router) error {
	// Status page file path
	statusPagePath := filepath.Join(hs.pagesDir, "platform", "status.html")
	statusURLPath := "/platform/status"
	
	// Check if status page exists
	if _, err := os.Stat(statusPagePath); os.IsNotExist(err) {
		return fmt.Errorf("status page not found: %s", statusPagePath)
	}
	
	// Get file info for modification time
	info, err := os.Stat(statusPagePath)
	if err != nil {
		return fmt.Errorf("failed to stat status page: %w", err)
	}
	
	// Ensure page content is preloaded
	if err := hs.preloadPageContent(); err != nil {
		return fmt.Errorf("failed to preload page content: %w", err)
	}
	
	// Ensure templates are pre-compiled
	if err := hs.precompileTemplates(); err != nil {
		return fmt.Errorf("failed to pre-compile templates: %w", err)
	}
	
	// Regenerate for each language
	for _, lang := range SupportedLanguages {
		// Render the status page with current data using optimized method
		html, err := hs.renderHTMLPageOptimized(
			statusPagePath, 
			statusURLPath, 
			router.navigationService, 
			router.seoService, 
			lang,
		)
		if err != nil {
			log.Printf("Failed to regenerate status page for language %s: %v", lang, err)
			continue // Continue with other languages
		}
		
		// Detect if the HTML contains code blocks
		needsCodeHighlight := DetectCodeBlocks(html)
		
		// Cache the rendered content
		cachedContent := &CachedContent{
			HTML:               html,
			Frontmatter:        nil,
			ModTime:            info.ModTime(),
			FilePath:           statusPagePath,
			NeedsCodeHighlight: needsCodeHighlight,
		}
		
		cacheKey := hs.getCacheKey(lang, statusURLPath)
		hs.cache.Set(cacheKey, cachedContent)
		
		log.Printf("Regenerated status page for language: %s", lang)
	}
	
	return nil
}

// preparePageData creates PageData with metadata for the given path
// This is a wrapper around the shared page data preparation logic
func (hs *HTMLService) preparePageData(path string, content template.HTML, isMarkdown bool, frontmatter *Frontmatter, navigationService *NavigationService, seoService *SEOService, lang string) PageData {
	// Create a temporary router instance to use its preparePageData method
	// This avoids code duplication while maintaining the existing API
	tempRouter := &Router{
		markdownService:   hs.markdownService,
		seoService:        seoService,
		navigationService: navigationService,
		schemaService:     hs.schemaService,
		statusChecker:     hs.statusChecker,
		tocExcludedPaths:  []string{"/changelog", "/roadmap", "/platform/status"},
	}

	return tempRouter.preparePageData(path, content, isMarkdown, frontmatter, navigationService.GetNavigationForPathWithLanguage(path, lang), lang)
}
