package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// serve404 serves a custom 404 page or falls back to default
func (r *Router) serve404(w http.ResponseWriter, req *http.Request) {
	// Extract language from the request URL
	lang, _ := extractLanguageFromPath(req.URL.Path)
	if lang == "" {
		lang = detectPreferredLanguage(req)
	}
	// Try to serve custom 404.html page
	notFoundPath := filepath.Join(r.pagesDir, "404.html")
	if _, err := os.Stat(notFoundPath); err == nil {
		// Custom 404 page exists, we'll set status when we write the response

		// Read the 404 page content
		contentBytes, err := os.ReadFile(notFoundPath)
		if err != nil {
			// Fallback to default if we can't read the file
			http.NotFound(w, req)
			return
		}

		// Prepare page data for 404 page (before processing templates)
		pageData := r.preparePageData("/404", "", false, nil, r.navigationService.GetNavigationForPath("/404"), lang)

		// Create a template for the 404 page content
		contentTmpl := template.New("404-content").Funcs(getTemplateFuncs(lang))

		// Auto-scan all component templates for 404 content parsing
		componentFiles, err := r.loadComponentTemplates()
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Parse component files first
		if len(componentFiles) > 0 {
			contentTmpl, err = contentTmpl.ParseFiles(componentFiles...)
			if err != nil {
				http.NotFound(w, req)
				return
			}
		}

		// Then parse the 404 page content
		contentTmpl, err = contentTmpl.Parse(string(contentBytes))
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Execute the 404 content template with the page data
		var renderedContent strings.Builder
		if err = contentTmpl.Execute(&renderedContent, pageData); err != nil {
			http.NotFound(w, req)
			return
		}

		// Now prepare template files for main layout
		templateFiles := []string{
			filepath.Join(r.layoutsDir, "main.html"),
		}
		templateFiles = append(templateFiles, componentFiles...)

		// Create template with custom functions for main layout
		tmpl := template.New("main.html").Funcs(getTemplateFuncs(lang))

		// Parse all template files
		tmpl, err = tmpl.ParseFiles(templateFiles...)
		if err != nil {
			http.NotFound(w, req)
			return
		}

		// Update page data with rendered content
		pageData.Content = template.HTML(renderedContent.String())

		// Set 404 status and execute main layout template
		w.WriteHeader(http.StatusNotFound)
		if err := tmpl.ExecuteTemplate(w, "main.html", pageData); err != nil {
			// Check if the error is due to client disconnect (broken pipe)
			if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
				// Client disconnected, this is expected during rapid navigation
				// Silently return without logging
				return
			}
			// Headers already written, just log the error
			log.Printf("Error executing 404 template: %v", err)
			return
		}
	} else {
		// Custom 404 page doesn't exist, use default
		http.NotFound(w, req)
	}
}