package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Router handles file-based routing for HTML pages
type Router struct {
	pagesDir          string
	layoutsDir        string
	componentsDir     string
	contentDir        string
	navigationService *NavigationService
	contentService    *ContentService
	markdownService   *MarkdownService
	htmlService       *HTMLService
	seoService        *SEOService
	schemaService     *SchemaService
	statusChecker     *HealthChecker
	tocExcludedPaths  []string
	logger            *Logger
}

// loadComponentTemplates loads all component template files
func (r *Router) loadComponentTemplates() ([]string, error) {
	componentFiles, err := filepath.Glob(filepath.Join(r.componentsDir, "*.html"))
	if err != nil {
		return nil, err
	}
	return componentFiles, nil
}

// NewRouter creates a new router instance
func NewRouter(pagesDir string, logger *Logger) *Router {
	// Initialize SEO service
	seoService := NewSEOService()
	if err := seoService.LoadData(); err != nil {
		log.Printf("⚠️  Error loading SEO data: %v", err)
	}

	// Initialize services
	markdownService := NewMarkdownService()
	contentService := NewContentService("content")
	navigationService := NewNavigationService(seoService)
	htmlService := NewHTMLService(pagesDir, "layouts", "components", markdownService)
	schemaService := NewSchemaService(seoService.metadata, "https://blue.cc")
	
	// Set schema service on HTML service
	htmlService.SetSchemaService(schemaService)

	router := &Router{
		pagesDir:          pagesDir,
		layoutsDir:        "layouts",
		componentsDir:     "components",
		contentDir:        "content",
		markdownService:   markdownService,
		contentService:    contentService,
		navigationService: navigationService,
		htmlService:       htmlService,
		seoService:        seoService,
		schemaService:     schemaService,
		statusChecker:     nil, // Will be set by SetStatusChecker
		tocExcludedPaths: []string{ // These pages will not show toc
			"/changelog",
			"/roadmap",
			"/platform/status",
		},
		logger: logger,
	}

	// Pre-render all markdown content at startup
	renderStart := time.Now()
	logger.Log(LogRender, "📝", "Starting", "Markdown pre-rendering")
	if err := markdownService.PreRenderAllMarkdown(contentService, seoService); err != nil {
		logger.Log(LogError, "⚠️", "Failed", fmt.Sprintf("Pre-render markdown: %v", err))
	} else {
		fileCount := markdownService.GetCacheSize()
		logger.Log(LogRender, "✅", "Completed", fmt.Sprintf("%d Markdown files cached", fileCount), time.Since(renderStart))
	}

	// Pre-render all HTML pages at startup
	htmlStart := time.Now()
	logger.Log(LogRender, "🌐", "Starting", "HTML pre-rendering")
	if err := htmlService.PreRenderAllHTMLPages(navigationService, seoService); err != nil {
		logger.Log(LogError, "⚠️", "Failed", fmt.Sprintf("Pre-render HTML: %v", err))
	} else {
		pageCount := htmlService.GetCacheSize()
		logger.Log(LogRender, "✅", "Completed", fmt.Sprintf("%d pages across %d languages", pageCount, len(SupportedLanguages)), time.Since(htmlStart))
	}

	// Generate search index with pre-rendered content
	searchStart := time.Now()
	logger.Log(LogIndex, "🔍", "Starting", "Search index generation")
	itemCount, err := GenerateSearchIndexWithCaches(markdownService, htmlService, logger)
	if err != nil {
		logger.Log(LogError, "⚠️", "Failed", fmt.Sprintf("Generate search index: %v", err))
	} else {
		logger.Log(LogIndex, "✅", "Completed", fmt.Sprintf("%d items indexed", itemCount), time.Since(searchStart))
	}

	// Link checker moved to main.go for parallel execution

	return router
}

// SetStatusChecker sets the status checker for the router
func (r *Router) SetStatusChecker(checker *HealthChecker) {
	r.statusChecker = checker
	// Also set it on the HTML service for status page rendering
	if r.htmlService != nil {
		r.htmlService.SetStatusChecker(checker)
	}
}


// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Handle API endpoints first (before language detection)
	if strings.HasPrefix(req.URL.Path, "/api/") {
		switch req.URL.Path {
		case "/api/assistant":
			HandleAssistant(w, req)
			return
		case "/api/assistant/stream":
			HandleAssistantStream(w, req)
			return
		default:
			http.Error(w, "API endpoint not found", http.StatusNotFound)
			return
		}
	}

	// Handle health check endpoint
	if req.URL.Path == "/health" {
		HealthHandler(w, req)
		return
	}

	// Status API endpoints removed - status is now served as static pre-rendered pages

	// Skip favicon requests early
	if req.URL.Path == "/favicon.ico" {
		return
	}

	// Language detection and routing
	lang, cleanPath := extractLanguageFromPath(req.URL.Path)
	
	// If no language in URL, detect and redirect
	if lang == "" {
		detectedLang := detectPreferredLanguage(req)
		// Redirect to language-prefixed URL
		http.Redirect(w, req, "/"+detectedLang+req.URL.Path, http.StatusFound)
		return
	}
	
	// Set language cookie for future visits
	setLanguageCookie(w, lang)

	// Serve component files
	if strings.HasPrefix(cleanPath, "/components/") {
		componentPath := strings.TrimPrefix(cleanPath, "/components/")
		
		// Add .html extension if not present
		if !strings.HasSuffix(componentPath, ".html") {
			componentPath += ".html"
		}
		
		// Sanitize the path to prevent directory traversal
		filePath := filepath.Join(r.componentsDir, componentPath)
		filePath = filepath.Clean(filePath)
		
		// Ensure the resolved path is within the components directory
		absComponentsDir, err := filepath.Abs(r.componentsDir)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		// Verify the file path is within the components directory
		if !strings.HasPrefix(absFilePath, absComponentsDir+string(filepath.Separator)) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, req, filePath)
		return
	}

	// Use the clean path (without language prefix)
	path := cleanPath

	// Handle api-docs to api redirects
	if strings.HasPrefix(path, "/api-docs/") {
		newPath := strings.Replace(path, "/api-docs/", "/api/", 1)
		http.Redirect(w, req, newPath, http.StatusMovedPermanently)
		return
	}

	// Check for redirects first
	if redirectTo, statusCode, shouldRedirect := r.seoService.CheckRedirect(path); shouldRedirect {
		// Preserve language prefix in redirects
		if lang != "" {
			redirectTo = "/" + lang + redirectTo
		}
		http.Redirect(w, req, redirectTo, statusCode)
		return
	}

	// Redirect .html URLs to clean URLs
	if strings.HasSuffix(path, ".html") && path != "/" {
		cleanURL := strings.TrimSuffix(path, ".html")
		http.Redirect(w, req, cleanURL, http.StatusMovedPermanently)
		return
	}

	// Convert URL path to file path
	var filePath string
	if path == "/" {
		// Root path maps to index.html
		filePath = filepath.Join(r.pagesDir, "index.html")
	} else {
		// Remove leading slash
		cleanPath := strings.TrimPrefix(path, "/")

		if strings.HasSuffix(cleanPath, "/") {
			// Directory path, look for index.html
			filePath = filepath.Join(r.pagesDir, cleanPath, "index.html")
		} else {
			// Try as direct .html file first
			filePath = filepath.Join(r.pagesDir, cleanPath+".html")

			// If that doesn't exist, try as directory with index.html
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				indexPath := filepath.Join(r.pagesDir, cleanPath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					// Redirect to trailing slash version
					http.Redirect(w, req, path+"/", http.StatusMovedPermanently)
					return
				}
			}
		}
	}

	var contentBytes []byte
	var isMarkdown bool
	var frontmatter *Frontmatter
	var cachedContent *CachedContent

	// Check if HTML page file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// HTML page doesn't exist, try cached markdown first
		if cached, found := r.markdownService.GetCachedContentForLang(path, lang); found {
			// Found in cache - use pre-rendered content
			cachedContent = cached
			contentBytes = []byte(cached.HTML)
			frontmatter = cached.Frontmatter
			isMarkdown = true
			
			// For SPA requests, we need to pass the code highlight flag
			// Since we're continuing to render through templates, store it temporarily
			if req.Header.Get("X-Requested-With") == "XMLHttpRequest" && cached.NeedsCodeHighlight {
				w.Header().Set("X-Needs-Code-Highlight", "true")
			}
			// Cached markdown served
		} else {
			// Not in cache, try to find and process markdown file (fallback)
			markdownPath, mdErr := r.contentService.FindMarkdownFileForLang(path, lang)
			if mdErr != nil {
				// Before serving 404, check if this is a directory URL that should redirect
				if firstItemURL := r.navigationService.GetFirstItemInDirectory(path); firstItemURL != "" {
					http.Redirect(w, req, firstItemURL, http.StatusMovedPermanently)
					return
				}
				r.serve404(w, req)
				return
			}

			// Process markdown file on-demand
			htmlContent, fm, err := r.markdownService.ProcessMarkdownFile(markdownPath, r.seoService)
			if err != nil {
				http.Error(w, "Error processing markdown file", http.StatusInternalServerError)
				log.Printf("Markdown processing error: %v", err)
				return
			}

			contentBytes = []byte(htmlContent)
			frontmatter = fm
			isMarkdown = true
			// On-demand markdown served
		}
	} else {
		// HTML page exists, try cached version first
		if cachedContent, found := r.htmlService.GetCachedContentForLang(path, lang); found {
			// Found in cache - use pre-rendered content
			w.Header().Set("Content-Type", "text/html")
			
			// For SPA requests, include code highlight flag in header
			if req.Header.Get("X-Requested-With") == "XMLHttpRequest" {
				if cachedContent.NeedsCodeHighlight {
					w.Header().Set("X-Needs-Code-Highlight", "true")
				}
			}
			
			_, err := w.Write([]byte(cachedContent.HTML))
			if err != nil {
				// Check if the error is due to client disconnect
				if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
					// Client disconnected, silently return
					return
				}
				// Log other write errors
				log.Printf("Error writing cached content: %v", err)
			}
			// Cached HTML served
			return
		}

		// Not in cache, process on-demand (fallback for dynamic pages like status)
		var err error
		contentBytes, err = os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading page", http.StatusInternalServerError)
			log.Printf("Page reading error: %v", err)
			return
		}
		isMarkdown = false

		// Process all HTML pages as templates to enable template variables
		pageData := r.preparePageData(path, "", isMarkdown, frontmatter, r.navigationService.GetNavigationForPathWithLanguage(path, lang), lang)

		// Create a template for the page content with language-specific functions
		contentTmpl := template.New("page-content").Funcs(getTemplateFuncs(lang))

		// Auto-scan all component templates for page content parsing
		componentFiles, err := r.loadComponentTemplates()
		if err != nil {
			http.Error(w, "Error loading components for page content", http.StatusInternalServerError)
			log.Printf("Component scanning error for page content: %v", err)
			return
		}

		// Parse component files first
		if len(componentFiles) > 0 {
			contentTmpl, err = contentTmpl.ParseFiles(componentFiles...)
			if err != nil {
				http.Error(w, "Error parsing component templates for page content", http.StatusInternalServerError)
				log.Printf("Component template parsing error for page content: %v", err)
				return
			}
		}

		// Then parse the page content
		contentTmpl, err = contentTmpl.Parse(string(contentBytes))
		if err != nil {
			http.Error(w, "Error parsing page template", http.StatusInternalServerError)
			log.Printf("Page template parsing error: %v", err)
			return
		}

		// Execute the page template with the page data
		var renderedContent strings.Builder
		if err = contentTmpl.Execute(&renderedContent, pageData); err != nil {
			http.Error(w, "Error rendering page template", http.StatusInternalServerError)
			log.Printf("Page template execution error: %v", err)
			return
		}

		contentBytes = []byte(renderedContent.String())
		// On-demand HTML served
	}

	// Prepare template files - start with layout
	templateFiles := []string{
		filepath.Join(r.layoutsDir, "main.html"),
	}

	// Auto-scan all component templates
	componentFiles, err := r.loadComponentTemplates()
	if err != nil {
		http.Error(w, "Error loading components", http.StatusInternalServerError)
		log.Printf("Component scanning error: %v", err)
		return
	}

	// Add all component files
	templateFiles = append(templateFiles, componentFiles...)

	// Create template with custom functions including language-specific translation
	tmpl := template.New("main.html").Funcs(getTemplateFuncs(lang))

	// Parse all template files
	tmpl, err = tmpl.ParseFiles(templateFiles...)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	// Prepare page data
	pageData := r.preparePageDataWithCache(path, template.HTML(contentBytes), isMarkdown, frontmatter, r.navigationService.GetNavigationForPathWithLanguage(path, lang), lang, cachedContent)

	// Set content type and execute main layout
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.ExecuteTemplate(w, "main.html", pageData); err != nil {
		// Check if the error is due to client disconnect (broken pipe)
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
			// Client disconnected, this is expected during rapid navigation
			// Silently return without logging
			return
		}
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
		return
	}
}

