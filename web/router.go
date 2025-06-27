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

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// NavItem represents a navigation item
type NavItem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Href       string    `json:"href,omitempty"`
	External   bool      `json:"external,omitempty"`
	Expanded   bool      `json:"expanded,omitempty"`
	Children   []NavItem `json:"children,omitempty"`
	OriginalID string    `json:"-"` // Store original directory name for sorting (not sent to JSON)
}

// Navigation holds the complete navigation structure
type Navigation struct {
	Sections []NavItem `json:"sections"`
}

// DirMetadata represents _dir.yml file content
type DirMetadata struct {
	Title string `yaml:"title"`
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
	pagesDir         string
	layoutsDir       string
	componentsDir    string
	contentDir       string
	navigation       *Navigation
	docsNavigation   *Navigation
	apiNavigation    *Navigation
	legalNavigation  *Navigation
	seoService       *SEOService
	changelogService *ChangelogService
	markdown         goldmark.Markdown
	tocExcludedPaths []string
}

// NewRouter creates a new router instance
func NewRouter(pagesDir string) *Router {
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
		),
	)

	// Initialize SEO service
	seoService := NewSEOService()
	if err := seoService.LoadData(); err != nil {
		log.Printf("Error loading SEO data: %v", err)
	}

	// Initialize changelog service
	changelogService := NewChangelogService()
	if err := changelogService.LoadChangelog(); err != nil {
		log.Printf("Error loading changelog data: %v", err)
	}

	router := &Router{
		pagesDir:         pagesDir,
		layoutsDir:       "layouts",
		componentsDir:    "components",
		contentDir:       "content",
		markdown:         md,
		seoService:       seoService,
		changelogService: changelogService,
		tocExcludedPaths: []string{ // These pages will not show toc
			"/changelog",
			"/roadmap",
		},
	}

	// Load navigation data
	if err := router.loadNavigation(); err != nil {
		log.Printf("Error loading navigation: %v", err)
	}

	// Generate dynamic navigation for docs
	if docsNav, err := router.generateContentNavigation("content/docs", "/docs"); err != nil {
		log.Printf("Error generating docs navigation: %v", err)
	} else {
		router.docsNavigation = docsNav
	}

	// Generate dynamic navigation for API
	if apiNav, err := router.generateContentNavigation("content/api-docs", "/api"); err != nil {
		log.Printf("Error generating API navigation: %v", err)
	} else {
		router.apiNavigation = apiNav
	}

	// Generate dynamic navigation for Legal
	if legalNav, err := router.generateContentNavigation("content/legal", "/legal"); err != nil {
		log.Printf("Error generating legal navigation: %v", err)
	} else {
		router.legalNavigation = legalNav
	}

	return router
}

// loadNavigation loads navigation data from JSON file
func (r *Router) loadNavigation() error {
	data, err := os.ReadFile("data/nav.json")
	if err != nil {
		return err
	}

	r.navigation = &Navigation{}
	return json.Unmarshal(data, r.navigation)
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
		// HTML page doesn't exist, try markdown
		markdownPath, mdErr := r.findMarkdownFile(path)
		if mdErr != nil {
			http.NotFound(w, req)
			return
		}

		// Read markdown file
		mdBytes, err := os.ReadFile(markdownPath)
		if err != nil {
			http.Error(w, "Error reading markdown file", http.StatusInternalServerError)
			log.Printf("Markdown reading error: %v", err)
			return
		}

		// Parse frontmatter
		var markdownContent []byte
		frontmatter, markdownContent, err = r.seoService.ParseFrontmatter(mdBytes)
		if err != nil {
			log.Printf("Frontmatter parsing error: %v", err)
			// Continue without frontmatter
			markdownContent = mdBytes
		} else if frontmatter != nil {
			log.Printf("Successfully parsed frontmatter: title=%s, desc=%s", frontmatter.Title, frontmatter.Description)
		} else {
			log.Printf("No frontmatter found in file")
		}

		// Convert markdown to HTML
		var htmlBuffer strings.Builder
		if err := r.markdown.Convert(markdownContent, &htmlBuffer); err != nil {
			http.Error(w, "Error converting markdown", http.StatusInternalServerError)
			log.Printf("Markdown conversion error: %v", err)
			return
		}

		contentBytes = []byte(htmlBuffer.String())
		isMarkdown = true
	} else {
		// Read HTML page content
		var err error
		contentBytes, err = os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading page", http.StatusInternalServerError)
			log.Printf("Page reading error: %v", err)
			return
		}
		isMarkdown = false

		// Process all HTML pages as templates to enable template variables
		pageData := r.preparePageData(path, "", isMarkdown, frontmatter, r.getNavigationForPath(path))

		// Create a template for the page content
		contentTmpl := template.New("page-content").Funcs(template.FuncMap{
			"toJSON": func(v any) template.JS {
				data, _ := json.Marshal(v)
				return template.JS(data)
			},
			"dict": dict,
			"html": func(s string) template.HTML {
				return template.HTML(s)
			},
		})

		// Auto-scan all component templates for page content parsing
		componentFiles, err := filepath.Glob(filepath.Join(r.componentsDir, "*.html"))
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
	}

	// Prepare template files - start with layout
	templateFiles := []string{
		filepath.Join(r.layoutsDir, "main.html"),
	}

	// Auto-scan all component templates
	componentFiles, err := filepath.Glob(filepath.Join(r.componentsDir, "*.html"))
	if err != nil {
		http.Error(w, "Error loading components", http.StatusInternalServerError)
		log.Printf("Component scanning error: %v", err)
		return
	}

	// Add all component files
	templateFiles = append(templateFiles, componentFiles...)

	// Create template with custom functions
	tmpl := template.New("main.html").Funcs(template.FuncMap{
		"toJSON": func(v any) template.JS {
			data, _ := json.Marshal(v)
			return template.JS(data)
		},
		"dict": dict,
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// Parse all template files
	tmpl, err = tmpl.ParseFiles(templateFiles...)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	// Prepare page data
	pageData := r.preparePageData(path, template.HTML(contentBytes), isMarkdown, frontmatter, r.getNavigationForPath(path))

	// Set content type and execute main layout
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.ExecuteTemplate(w, "main.html", pageData); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
		return
	}
}

// getNavigationForPath returns the appropriate navigation based on the URL path
func (r *Router) getNavigationForPath(path string) *Navigation {
	// Always start with static navigation
	if r.navigation == nil {
		return &Navigation{}
	}

	// Make a copy of the static navigation
	nav := &Navigation{
		Sections: make([]NavItem, len(r.navigation.Sections)),
	}
	copy(nav.Sections, r.navigation.Sections)

	// Always add Documentation section if available
	if r.docsNavigation != nil {
		docSection := NavItem{
			ID:       "documentation",
			Name:     "Documentation",
			Expanded: strings.HasPrefix(path, "/docs"), // Only expand when on docs pages
			Children: r.docsNavigation.Sections,
		}
		nav.Sections = append(nav.Sections, docSection)
	}

	// Always add API Reference section if available
	if r.apiNavigation != nil {
		apiSection := NavItem{
			ID:       "api-reference",
			Name:     "API Reference",
			Expanded: strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/api-docs"), // Only expand when on API pages
			Children: r.apiNavigation.Sections,
		}
		nav.Sections = append(nav.Sections, apiSection)
	}

	// Always add Legal section if available (placed at end for bottom positioning)
	if r.legalNavigation != nil {
		legalSection := NavItem{
			ID:       "legal",
			Name:     "Legal",
			Expanded: strings.HasPrefix(path, "/legal"), // Only expand when on legal pages
			Children: r.legalNavigation.Sections,
		}
		nav.Sections = append(nav.Sections, legalSection)
	}

	return nav
}

// findMarkdownFile searches for a markdown file matching the given path
func (r *Router) findMarkdownFile(path string) (string, error) {
	// Convert URL path to potential file paths
	cleanPath := strings.Trim(path, "/")

	// For content paths, we need to map clean URLs back to numbered files/directories
	if strings.HasPrefix(cleanPath, "docs/") || strings.HasPrefix(cleanPath, "api/") || strings.HasPrefix(cleanPath, "api-docs/") || strings.HasPrefix(cleanPath, "legal/") {
		return r.findNumberedMarkdownFile(cleanPath)
	}

	// Try simple patterns for non-content paths
	patterns := []string{
		filepath.Join(r.contentDir, cleanPath+".md"),
		filepath.Join(r.contentDir, cleanPath, "index.md"),
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern, nil
		}
	}

	return "", os.ErrNotExist
}

// findNumberedMarkdownFile handles finding files with numeric prefixes
func (r *Router) findNumberedMarkdownFile(cleanPath string) (string, error) {
	parts := strings.Split(cleanPath, "/")
	if len(parts) < 2 {
		return "", os.ErrNotExist
	}

	// Map content type to directory
	contentType := parts[0]
	var contentDir string
	switch contentType {
	case "docs":
		contentDir = "content/docs"
	case "api", "api-docs":
		contentDir = "content/api-docs"
	case "legal":
		contentDir = "content/legal"
	default:
		return "", os.ErrNotExist
	}

	// Build path progressively, finding numbered directories/files
	currentPath := contentDir
	for i := 1; i < len(parts); i++ {
		cleanSegment := parts[i]

		if i == len(parts)-1 {
			// Last segment - look for numbered file in the current specific directory
			// Try multiple patterns to handle spaces vs hyphens vs case variations
			patterns := []string{
				"*" + cleanSegment + ".md",                               // e.g., *download-apps.md
				"*" + strings.ReplaceAll(cleanSegment, "-", " ") + ".md", // e.g., *download apps.md
				"*" + strings.ReplaceAll(cleanSegment, "-", "_") + ".md", // e.g., *download_apps.md
			}

			// Try each pattern with case-insensitive matching
			for _, pattern := range patterns {
				glob := filepath.Join(currentPath, pattern)
				matches, err := filepath.Glob(glob)
				if err == nil && len(matches) > 0 {
					return matches[0], nil
				}

				// If no matches, try case-insensitive search by reading directory
				if err := r.findFileIgnoreCase(currentPath, pattern, &matches); err == nil && len(matches) > 0 {
					return matches[0], nil
				}
			}

			// Also try as directory with index
			glob := filepath.Join(currentPath, "*"+cleanSegment, "index.md")
			matches, err := filepath.Glob(glob)
			if err == nil && len(matches) > 0 {
				return matches[0], nil
			}
		} else {
			// Intermediate segment - look for numbered directory
			// Try multiple patterns to handle spaces vs hyphens in directory names
			dirPatterns := []string{
				"*" + cleanSegment, // e.g., *start-guide
				"*" + strings.ReplaceAll(cleanSegment, "-", " "), // e.g., *start guide
				"*" + strings.ReplaceAll(cleanSegment, "-", "_"), // e.g., *start_guide
			}

			found := false
			for _, pattern := range dirPatterns {
				glob := filepath.Join(currentPath, pattern)
				matches, err := filepath.Glob(glob)
				if err == nil && len(matches) > 0 {
					currentPath = matches[0]
					found = true
					break
				}
			}

			if !found {
				return "", os.ErrNotExist
			}
		}
	}

	return "", os.ErrNotExist
}

// findFileIgnoreCase performs case-insensitive file matching
func (r *Router) findFileIgnoreCase(dir, pattern string, matches *[]string) error {
	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Extract the pattern without the wildcard and directory
	patternName := strings.TrimPrefix(pattern, "*")
	patternLower := strings.ToLower(patternName)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// Check if file matches pattern (case-insensitive)
		if strings.HasSuffix(strings.ToLower(fileName), patternLower) {
			// Also check if it has a numeric prefix (to match our numbered file pattern)
			parts := strings.SplitN(fileName, ".", 2)
			if len(parts) == 2 {
				// Check if first part starts with a number
				if len(parts[0]) > 0 && parts[0][0] >= '0' && parts[0][0] <= '9' {
					*matches = append(*matches, filepath.Join(dir, fileName))
					return nil
				}
			}
		}
	}

	return os.ErrNotExist
}

// generateContentNavigation creates navigation tree from content directory
func (r *Router) generateContentNavigation(contentDir, baseURL string) (*Navigation, error) {
	var sections []NavItem

	// Read the content directory
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}

	// Process each directory/file
	for _, entry := range entries {
		if entry.IsDir() {
			navItem, err := r.processDirectory(filepath.Join(contentDir, entry.Name()), entry.Name(), baseURL)
			if err != nil {
				log.Printf("Error processing directory %s: %v", entry.Name(), err)
				continue
			}
			if navItem != nil {
				sections = append(sections, *navItem)
			}
		} else if strings.HasSuffix(entry.Name(), ".md") {
			// Handle individual markdown files at root level
			fileName := strings.TrimSuffix(entry.Name(), ".md")
			fileTitle := r.cleanTitle(fileName)

			// Try to get title from frontmatter
			if filePath := filepath.Join(contentDir, entry.Name()); filePath != "" {
				if data, err := os.ReadFile(filePath); err == nil {
					if frontmatter, _, err := r.seoService.ParseFrontmatter(data); err == nil && frontmatter != nil && frontmatter.Title != "" {
						fileTitle = frontmatter.Title
					}
				}
			}

			href := baseURL + "/" + r.cleanID(fileName)

			sections = append(sections, NavItem{
				ID:         r.cleanID(fileName),
				Name:       fileTitle,
				Href:       href,
				OriginalID: fileName,
			})
		}
	}

	// Sort sections by numeric prefix
	r.sortNavItems(sections)

	return &Navigation{Sections: sections}, nil
}

// processDirectory processes a content directory and creates NavItem
func (r *Router) processDirectory(dirPath, dirName, baseURL string) (*NavItem, error) {
	// Read directory metadata
	title := r.cleanTitle(dirName)
	dirMetaPath := filepath.Join(dirPath, "_dir.yml")
	if data, err := os.ReadFile(dirMetaPath); err == nil {
		var dirMeta DirMetadata
		if err := yaml.Unmarshal(data, &dirMeta); err == nil && dirMeta.Title != "" {
			title = dirMeta.Title
		}
	}

	// Create nav item
	navItem := &NavItem{
		ID:         r.cleanID(dirName),
		Name:       title,
		Expanded:   false,
		OriginalID: dirName,
	}

	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return navItem, nil
	}

	var children []NavItem

	// Process subdirectories and files
	for _, entry := range entries {
		if entry.Name() == "_dir.yml" {
			continue
		}

		if entry.IsDir() {
			// Recursive subdirectory
			childNav, err := r.processDirectory(filepath.Join(dirPath, entry.Name()), entry.Name(), baseURL)
			if err != nil {
				log.Printf("Error processing subdirectory %s: %v", entry.Name(), err)
				continue
			}
			if childNav != nil {
				children = append(children, *childNav)
			}
		} else if strings.HasSuffix(entry.Name(), ".md") {
			// Markdown file
			fileName := strings.TrimSuffix(entry.Name(), ".md")
			fileTitle := r.cleanTitle(fileName)

			// Try to get title from frontmatter
			if filePath := filepath.Join(dirPath, entry.Name()); filePath != "" {
				if data, err := os.ReadFile(filePath); err == nil {
					if frontmatter, _, err := r.seoService.ParseFrontmatter(data); err == nil && frontmatter != nil && frontmatter.Title != "" {
						fileTitle = frontmatter.Title
					}
				}
			}

			// Create relative path for href
			// Remove both content dir and the specific content type (docs/api-docs)
			relDir := strings.TrimPrefix(dirPath, r.contentDir+"/")

			// Remove the content type prefix (e.g., "docs/" or "api-docs/" or "legal/")
			if strings.HasPrefix(relDir, "docs/") {
				relDir = strings.TrimPrefix(relDir, "docs/")
			} else if strings.HasPrefix(relDir, "api-docs/") {
				relDir = strings.TrimPrefix(relDir, "api-docs/")
			} else if strings.HasPrefix(relDir, "legal/") {
				relDir = strings.TrimPrefix(relDir, "legal/")
			}

			// Clean numeric prefixes from directory path
			relDir = r.cleanDirectoryPath(relDir)

			href := baseURL + "/" + relDir + "/" + r.cleanID(fileName)

			children = append(children, NavItem{
				ID:         r.cleanID(fileName),
				Name:       fileTitle,
				Href:       href,
				OriginalID: fileName,
			})
		}
	}

	// Sort children by numeric prefix
	r.sortNavItems(children)

	if len(children) > 0 {
		navItem.Children = children
	}

	return navItem, nil
}

// cleanTitle removes numeric prefixes and cleans up titles
func (r *Router) cleanTitle(name string) string {
	// Remove numeric prefix (e.g., "1.start-guide" -> "start-guide")
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		name = parts[1]
	}

	// Replace hyphens/underscores with spaces and title case
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Simple title case
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// cleanID creates clean IDs for navigation
func (r *Router) cleanID(name string) string {
	// Remove numeric prefix
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		name = parts[1]
	}

	// Convert to lowercase and replace spaces/special chars with hyphens
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	return name
}

// cleanDirectoryPath removes numeric prefixes from directory paths
func (r *Router) cleanDirectoryPath(path string) string {
	// Split path into segments and clean each one
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = r.cleanID(segment)
	}
	return strings.Join(segments, "/")
}

// sortNavItems sorts navigation items by numeric prefix
func (r *Router) sortNavItems(items []NavItem) {
	// Sort by extracting numeric prefix from original names
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			// Extract numeric prefixes for comparison using original IDs
			num1 := r.extractNumericPrefix(items[i].OriginalID)
			num2 := r.extractNumericPrefix(items[j].OriginalID)

			if num1 > num2 {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// extractNumericPrefix extracts the numeric prefix from a name (e.g., "1.start-guide" -> 1)
func (r *Router) extractNumericPrefix(name string) int {
	// Parse numeric prefix from original directory/file names
	parts := strings.Split(name, ".")
	if len(parts) >= 2 {
		// Try to parse first part as number
		num := 0
		for _, char := range parts[0] {
			if char >= '0' && char <= '9' {
				num = num*10 + int(char-'0')
			} else {
				break
			}
		}
		if num > 0 {
			return num
		}
	}

	// Fallback: assign high number for non-numbered items
	return 9999
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
	var toc []TOCEntry
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
