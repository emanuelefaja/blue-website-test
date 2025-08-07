package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Handle actual API endpoints first (before language detection)
	if strings.HasPrefix(req.URL.Path, "/api/") {
		if r.handleAPI(w, req) {
			return
		}
		// If it's not an actual API endpoint, fall through to regular routing
		// This allows API documentation (/api/introduction, etc.) to be served as content
	}

	// Handle health check endpoint
	if req.URL.Path == "/health" {
		HealthHandler(w, req)
		return
	}

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
		r.handleComponent(w, req, cleanPath)
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

	// Check for redirects
	if r.checkRedirects(w, req, path, lang) {
		return
	}

	// Redirect .html URLs to clean URLs
	if strings.HasSuffix(path, ".html") && path != "/" {
		cleanURL := strings.TrimSuffix(path, ".html")
		http.Redirect(w, req, cleanURL, http.StatusMovedPermanently)
		return
	}

	// Main page handling
	r.handlePage(w, req, path, lang)
}

// handleAPI handles API endpoints, returns true if handled, false if should fall through
func (r *Router) handleAPI(w http.ResponseWriter, req *http.Request) bool {
	switch req.URL.Path {
	case "/api/assistant":
		HandleAssistant(w, req)
		return true
	case "/api/assistant/stream":
		HandleAssistantStream(w, req)
		return true
	default:
		// Not an actual API endpoint, let it fall through to content serving
		return false
	}
}

// handleComponent serves component files
func (r *Router) handleComponent(w http.ResponseWriter, req *http.Request, cleanPath string) {
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
}

// checkRedirects handles URL redirects
func (r *Router) checkRedirects(w http.ResponseWriter, req *http.Request, path, lang string) bool {
	if redirectTo, statusCode, shouldRedirect := r.seoService.CheckRedirect(path); shouldRedirect {
		// Preserve language prefix in redirects
		if lang != "" {
			redirectTo = "/" + lang + redirectTo
		}
		http.Redirect(w, req, redirectTo, statusCode)
		return true
	}
	return false
}

// handlePage routes the request to appropriate render method
func (r *Router) handlePage(w http.ResponseWriter, req *http.Request, path, lang string) {
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

	// Check if HTML page file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// HTML page doesn't exist, try to render markdown
		r.RenderMarkdown(w, req, path, lang)
	} else {
		// HTML page exists, render it
		r.RenderHTML(w, req, filePath, path, lang)
	}
}