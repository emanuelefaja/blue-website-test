package web

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// templateFuncs defines template functions used across all templates
var templateFuncs = template.FuncMap{
	"toJSON": func(v any) template.JS {
		data, _ := json.Marshal(v)
		return template.JS(data)
	},
	"dict": dict,
	"html": func(s string) template.HTML {
		return template.HTML(s)
	},
	"parseJSON": func(s string) (interface{}, error) {
		var data interface{}
		err := json.Unmarshal([]byte(s), &data)
		return data, err
	},
}

// PageData holds data for template rendering
type PageData struct {
	Title          string
	Content        template.HTML
	Navigation     *Navigation
	PageMeta       *PageMetadata
	SiteMeta       *SiteMetadata
	Description    string
	Keywords       []string
	IsMarkdown     bool
	Frontmatter    *Frontmatter
	Changelog      []ChangelogMonth
	TOC            []TOCEntry
	CustomerNumber int
}

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
	changelogService  *ChangelogService
	tocExcludedPaths  []string
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
func NewRouter(pagesDir string) *Router {
	// Initialize SEO service
	seoService := NewSEOService()
	if err := seoService.LoadData(); err != nil {
		log.Printf("Error loading SEO data: %v", err)
	}

	// Initialize services
	markdownService := NewMarkdownService()
	contentService := NewContentService("content")
	navigationService := NewNavigationService(seoService)
	htmlService := NewHTMLService(pagesDir, "layouts", "components")

	// Initialize changelog service
	changelogService := NewChangelogService()
	if err := changelogService.LoadChangelog(); err != nil {
		log.Printf("Error loading changelog data: %v", err)
	}

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
		changelogService:  changelogService,
		tocExcludedPaths: []string{ // These pages will not show toc
			"/changelog",
			"/roadmap",
			"/platform/status",
		},
	}

	// Pre-render all markdown content at startup
	log.Printf("Pre-rendering markdown content...")
	if err := markdownService.PreRenderAllMarkdown(contentService, seoService); err != nil {
		log.Printf("Warning: failed to pre-render markdown content: %v", err)
	} else {
		log.Printf("Markdown cache initialized with %d files", markdownService.GetCacheSize())
	}

	// Pre-render all HTML pages at startup
	log.Printf("Pre-rendering HTML pages...")
	if err := htmlService.PreRenderAllHTMLPages(navigationService, seoService, changelogService); err != nil {
		log.Printf("Warning: failed to pre-render HTML pages: %v", err)
	} else {
		log.Printf("HTML cache initialized with %d pages", htmlService.GetCacheSize())
	}

	// Generate search index with pre-rendered content
	log.Printf("Generating search index with cached content...")
	if err := GenerateSearchIndexWithCaches(markdownService, htmlService); err != nil {
		log.Printf("Warning: failed to generate search index: %v", err)
	} else {
		log.Printf("Search index generated successfully")
	}

	return router
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Skip public file requests
	if strings.HasPrefix(req.URL.Path, "/public/") {
		return
	}

	// Serve component files
	if strings.HasPrefix(req.URL.Path, "/components/") {
		componentPath := strings.TrimPrefix(req.URL.Path, "/components/")
		filePath := filepath.Join(r.componentsDir, componentPath)

		// Add .html extension if not present
		if !strings.HasSuffix(filePath, ".html") {
			filePath += ".html"
		}

		http.ServeFile(w, req, filePath)
		return
	}

	// Get the requested path
	path := req.URL.Path

	// Check for redirects first
	if redirectTo, statusCode, shouldRedirect := r.seoService.CheckRedirect(path); shouldRedirect {
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

	// Check if HTML page file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// HTML page doesn't exist, try cached markdown first
		if cachedContent, found := r.markdownService.GetCachedContent(path); found {
			// Found in cache - use pre-rendered content
			contentBytes = []byte(cachedContent.HTML)
			frontmatter = cachedContent.Frontmatter
			isMarkdown = true
			log.Printf("Served cached markdown for path: %s", path)
		} else {
			// Not in cache, try to find and process markdown file (fallback)
			markdownPath, mdErr := r.contentService.FindMarkdownFile(path)
			if mdErr != nil {
				http.NotFound(w, req)
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
			log.Printf("Served on-demand markdown for path: %s", path)
		}
	} else {
		// HTML page exists, try cached version first
		if cachedContent, found := r.htmlService.GetCachedContent(path); found {
			// Found in cache - use pre-rendered content
			http.Header.Set(w.Header(), "Content-Type", "text/html")
			w.Write([]byte(cachedContent.HTML))
			log.Printf("Served cached HTML page for path: %s", path)
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
		pageData := r.preparePageData(path, "", isMarkdown, frontmatter, r.navigationService.GetNavigationForPath(path))

		// Create a template for the page content
		contentTmpl := template.New("page-content").Funcs(templateFuncs)

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
		log.Printf("Served on-demand HTML page for path: %s", path)
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

	// Create template with custom functions
	tmpl := template.New("main.html").Funcs(templateFuncs)

	// Parse all template files
	tmpl, err = tmpl.ParseFiles(templateFiles...)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	// Prepare page data
	pageData := r.preparePageData(path, template.HTML(contentBytes), isMarkdown, frontmatter, r.navigationService.GetNavigationForPath(path))

	// Set content type and execute main layout
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.ExecuteTemplate(w, "main.html", pageData); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
		return
	}
}

// isTOCExcluded checks if a path should be excluded from TOC generation
func (r *Router) isTOCExcluded(path string) bool {
	for _, excludedPath := range r.tocExcludedPaths {
		if path == excludedPath {
			return true
		}
	}
	return false
}

// preparePageData creates PageData with metadata for the given path
func (r *Router) preparePageData(path string, content template.HTML, isMarkdown bool, frontmatter *Frontmatter, navigation *Navigation) PageData {
	// Get metadata from SEO service
	title, description, keywords, pageMeta, siteMeta := r.seoService.PreparePageMetadata(path, isMarkdown, frontmatter)

	// Prepare changelog data - only include if on changelog page
	var changelog []ChangelogMonth
	if path == "/changelog" && r.changelogService != nil {
		changelog = r.changelogService.GetChangelog()
		log.Printf("Loading changelog for path=%s, found %d entries", path, len(changelog))
	} else {
		log.Printf("Not loading changelog: path=%s, service=%v", path, r.changelogService != nil)
	}

	// Extract table of contents (skip if path is excluded)
	toc := make([]TOCEntry, 0)
	if !r.isTOCExcluded(path) && string(content) != "" {
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
	} else if r.isTOCExcluded(path) {
		log.Printf("TOC generation skipped for excluded path: %s", path)
	}

	// Return PageData with all components
	return PageData{
		Title:          title,
		Content:        content,
		Navigation:     navigation,
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

// dict creates a map from key-value pairs for use in templates
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("dict requires even number of arguments")
	}
	m := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m, nil
}
