package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// PageMetadata represents metadata for a specific page
type PageMetadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

// SiteMetadata represents global site metadata
type SiteMetadata struct {
	Name         string            `json:"name"`
	Descriptions map[string]string `json:"descriptions"`
	URL          string            `json:"url"`
	Author       string            `json:"author"`
}

// MetadataDefaults represents default metadata values
type MetadataDefaults struct {
	TitleSuffix  string              `json:"title_suffix"`
	Descriptions map[string]string   `json:"descriptions"`
	Keywords     map[string][]string `json:"keywords"`
}

// Metadata holds the complete metadata structure
type Metadata struct {
	Site     SiteMetadata                       `json:"site"`
	Pages    map[string]map[string]PageMetadata `json:"pages"`
	Defaults MetadataDefaults                   `json:"defaults"`
}

// Frontmatter represents markdown file frontmatter
type Frontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Slug        string   `yaml:"slug"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	Image       string   `yaml:"image"`
	Date        string   `yaml:"date"`
	ShowDate    bool     `yaml:"showdate"`
}

// RedirectRules represents redirect configuration rules
type RedirectRules struct {
	StatusCode    int    `json:"status_code"`
	TrailingSlash string `json:"trailing_slash"`
}

// Redirects holds the complete redirect configuration
type Redirects struct {
	Redirects map[string]string `json:"redirects"`
	Rules     RedirectRules     `json:"rules"`
}

// URLEntry represents a single URL in the sitemap
type URLEntry struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	LastMod    string   `xml:"lastmod,omitempty"`
	ChangeFreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

// URLSet represents the root sitemap element
type URLSet struct {
	XMLName xml.Name   `xml:"urlset"`
	Xmlns   string     `xml:"xmlns,attr"`
	URLs    []URLEntry `xml:"url"`
}

// SitemapEntry represents a single sitemap in the sitemap index
type SitemapEntry struct {
	XMLName xml.Name `xml:"sitemap"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod,omitempty"`
}

// SitemapIndex represents the root sitemap index element
type SitemapIndex struct {
	XMLName  xml.Name       `xml:"sitemapindex"`
	Xmlns    string         `xml:"xmlns,attr"`
	Sitemaps []SitemapEntry `xml:"sitemap"`
}

// SEOService handles all SEO-related functionality
type SEOService struct {
	metadata  *Metadata
	redirects *Redirects
}

// NewSEOService creates a new SEO service instance
func NewSEOService() *SEOService {
	return &SEOService{}
}

// LoadData loads metadata and redirects from files
func (s *SEOService) LoadData() error {
	if err := s.loadMetadata(); err != nil {
		log.Printf("Error loading metadata: %v", err)
	}

	if err := s.loadRedirects(); err != nil {
		log.Printf("Error loading redirects: %v", err)
	}

	return nil
}

// loadMetadata loads metadata from JSON file
func (s *SEOService) loadMetadata() error {
	data, err := os.ReadFile("data/metadata.json")
	if err != nil {
		return err
	}

	s.metadata = &Metadata{}
	return json.Unmarshal(data, s.metadata)
}

// loadRedirects loads redirect configuration from JSON file
func (s *SEOService) loadRedirects() error {
	data, err := os.ReadFile("data/redirects.json")
	if err != nil {
		return err
	}

	s.redirects = &Redirects{}
	return json.Unmarshal(data, s.redirects)
}

// CheckRedirect checks if a path should be redirected
func (s *SEOService) CheckRedirect(path string) (string, int, bool) {
	if s.redirects != nil {
		if redirectTo, exists := s.redirects.Redirects[path]; exists {
			statusCode := s.redirects.Rules.StatusCode
			if statusCode == 0 {
				statusCode = 301 // Default to permanent redirect
			}
			return redirectTo, statusCode, true
		}
	}
	return "", 0, false
}

// ParseFrontmatter extracts frontmatter from markdown content
func (s *SEOService) ParseFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	contentStr := string(content)

	// Normalize line endings to Unix style
	contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
	contentStr = strings.ReplaceAll(contentStr, "\r", "\n")

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(contentStr, "---\n") {
		return nil, content, nil
	}

	// Find the end of frontmatter - look for closing --- on its own line
	lines := strings.Split(contentStr, "\n")
	var frontmatterLines []string
	var contentStart int

	// Skip the opening ---
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			// Found closing delimiter
			contentStart = i + 1
			break
		}
		frontmatterLines = append(frontmatterLines, lines[i])
	}

	// If we didn't find a closing delimiter, treat as regular content
	if contentStart == 0 {
		return nil, content, nil
	}

	// Parse the frontmatter YAML
	frontmatterYAML := strings.Join(frontmatterLines, "\n")
	var frontmatter Frontmatter
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
		return nil, content, err
	}

	// Return frontmatter and content without frontmatter
	remainingLines := lines[contentStart:]
	markdownContent := []byte(strings.Join(remainingLines, "\n"))
	return &frontmatter, markdownContent, nil
}

// PreparePageMetadata creates page metadata for the given path
func (s *SEOService) PreparePageMetadata(path string, isMarkdown bool, frontmatter *Frontmatter, lang string) (string, string, []string, *PageMetadata, *SiteMetadata) {
	// Get page key for metadata lookup
	pageKey := s.getPageKey(path)

	// Get page metadata
	var pageMeta *PageMetadata
	var title, description string
	var keywords []string

	// For markdown files, prioritize frontmatter over metadata.json
	if isMarkdown && frontmatter != nil {
		if frontmatter.Title != "" {
			title = frontmatter.Title
		}
		if frontmatter.Description != "" {
			description = frontmatter.Description
		}
	}

	// If no frontmatter or missing fields, use metadata.json
	if title == "" || description == "" {
		if s.metadata != nil {
			// Check if specific page metadata exists for the language
			if pageData, exists := s.metadata.Pages[pageKey]; exists {
				if langMeta, langExists := pageData[lang]; langExists {
					pageMeta = &langMeta
					if title == "" {
						title = langMeta.Title
					}
					if description == "" {
						description = langMeta.Description
					}
					keywords = langMeta.Keywords
				} else if enMeta, enExists := pageData["en"]; enExists {
					// Fallback to English if language not found
					pageMeta = &enMeta
					if title == "" {
						title = enMeta.Title
					}
					if description == "" {
						description = enMeta.Description
					}
					keywords = enMeta.Keywords
				}
			} else {
				// Use defaults
				if title == "" {
					title = s.getFallbackTitle(path)
				}
				if description == "" {
					// Try language-specific default description
					if s.metadata.Defaults.Descriptions != nil {
						if desc, ok := s.metadata.Defaults.Descriptions[lang]; ok && desc != "" {
							description = desc
						} else if desc, ok := s.metadata.Defaults.Descriptions["en"]; ok && desc != "" {
							description = desc
						}
					}
				}
				// Try language-specific default keywords
				if s.metadata.Defaults.Keywords != nil {
					if kw, ok := s.metadata.Defaults.Keywords[lang]; ok && len(kw) > 0 {
						keywords = kw
					} else if kw, ok := s.metadata.Defaults.Keywords["en"]; ok && len(kw) > 0 {
						keywords = kw
					}
				}
			}
		} else {
			// Fallback if no metadata loaded
			if title == "" {
				title = s.getFallbackTitle(path)
			}
			if description == "" {
				description = "Blue - Powerful platform to create, manage, and scale processes for modern teams."
			}
			keywords = []string{"blue", "process management", "team collaboration"}
		}
	}

	// Apply title suffix if defined and not already present
	if s.metadata != nil && s.metadata.Defaults.TitleSuffix != "" {
		if !strings.HasSuffix(title, s.metadata.Defaults.TitleSuffix) {
			title = title + s.metadata.Defaults.TitleSuffix
		}
	}

	var siteMeta *SiteMetadata
	if s.metadata != nil {
		siteMeta = &s.metadata.Site
	}

	return title, description, keywords, pageMeta, siteMeta
}

// getPageKey converts URL path to metadata key
func (s *SEOService) getPageKey(path string) string {
	if path == "/" {
		return "home"
	}

	// Remove leading/trailing slashes
	cleanPath := strings.Trim(path, "/")
	return cleanPath
}

// getFallbackTitle creates a fallback title from URL path
func (s *SEOService) getFallbackTitle(path string) string {
	if path == "/" {
		return "Home"
	}

	// Remove leading slash and convert to title case
	cleanPath := strings.TrimPrefix(path, "/")
	cleanPath = strings.TrimSuffix(cleanPath, "/")

	// Replace slashes with spaces and title case
	parts := strings.Split(cleanPath, "/")
	for i, part := range parts {
		// Simple title case replacement for strings.Title
		words := strings.Split(strings.ReplaceAll(part, "-", " "), " ")
		for j, word := range words {
			if len(word) > 0 {
				words[j] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
		parts[i] = strings.Join(words, " ")
	}

	return strings.Join(parts, " - ")
}

// GenerateSitemap creates a sitemap.xml file in the public directory
func (s *SEOService) GenerateSitemap(baseURL string) error {
	// Generate main sitemap index
	if err := s.generateSitemapIndex(baseURL); err != nil {
		return fmt.Errorf("failed to generate sitemap index: %w", err)
	}

	// Generate language-specific sitemaps
	for _, lang := range SupportedLanguages {
		if err := s.generateLanguageSitemap(baseURL, lang); err != nil {
			return fmt.Errorf("failed to generate sitemap for %s: %w", lang, err)
		}
	}

	return nil
}

// generateSitemapIndex creates the main sitemap index file
func (s *SEOService) generateSitemapIndex(baseURL string) error {
	// Create sitemap index entries
	var sitemaps []SitemapEntry
	currentTime := time.Now().Format("2006-01-02")

	for _, lang := range SupportedLanguages {
		filename := "sitemap.xml"
		if lang != DefaultLanguage {
			filename = fmt.Sprintf("sitemap-%s.xml", lang)
		}

		sitemaps = append(sitemaps, SitemapEntry{
			Loc:     baseURL + "/" + filename,
			LastMod: currentTime,
		})
	}

	// Create sitemap index structure
	sitemapIndex := SitemapIndex{
		Xmlns:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: sitemaps,
	}

	// Generate XML
	xmlData, err := xml.MarshalIndent(sitemapIndex, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sitemap index XML: %w", err)
	}

	// Add XML declaration
	xmlContent := xml.Header + string(xmlData)

	// Write to public/sitemap_index.xml
	if err := os.WriteFile("public/sitemap_index.xml", []byte(xmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write sitemap_index.xml: %w", err)
	}

	// Also create a copy as sitemap.xml for backward compatibility
	if err := os.WriteFile("public/sitemap.xml", []byte(xmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write sitemap.xml: %w", err)
	}

	return nil
}

// generateLanguageSitemap creates a language-specific sitemap
func (s *SEOService) generateLanguageSitemap(baseURL, lang string) error {
	var urls []URLEntry
	currentTime := time.Now().Format("2006-01-02")

	// Add static HTML pages for this language
	htmlUrls, err := s.scanHTMLPagesForLanguage("pages", baseURL, currentTime, lang)
	if err != nil {
		log.Printf("Warning: failed to scan HTML pages for %s: %v", lang, err)
	} else {
		urls = append(urls, htmlUrls...)
	}

	// Add markdown content pages for this language
	markdownUrls, err := s.scanMarkdownContentForLanguage("content", baseURL, currentTime, lang)
	if err != nil {
		log.Printf("Warning: failed to scan markdown content for %s: %v", lang, err)
	} else {
		urls = append(urls, markdownUrls...)
	}

	// Create sitemap structure
	sitemap := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	// Generate XML
	xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sitemap XML: %w", err)
	}

	// Add XML declaration
	xmlContent := xml.Header + string(xmlData)

	// Determine filename
	filename := "public/sitemap.xml"
	if lang != DefaultLanguage {
		filename = fmt.Sprintf("public/sitemap-%s.xml", lang)
	}

	// Write to file
	if err := os.WriteFile(filename, []byte(xmlContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	return nil
}

// filePathToURL converts a file path to a URL path
func (s *SEOService) filePathToURL(filePath, baseDir string) string {
	// Remove base directory prefix
	urlPath := strings.TrimPrefix(filePath, baseDir)
	urlPath = strings.TrimPrefix(urlPath, "/")

	// Skip certain files
	if strings.Contains(filePath, "copy.html") {
		return "" // Skip backup files
	}
	if strings.Contains(filePath, "404.html") {
		return "" // Skip 404 page from sitemap
	}

	// Convert index.html to directory URLs
	if strings.HasSuffix(urlPath, "/index.html") {
		urlPath = strings.TrimSuffix(urlPath, "/index.html")
		if urlPath == "" {
			return "/" // Root page
		}
		return "/" + urlPath + "/"
	}

	// Convert regular HTML files
	if strings.HasSuffix(urlPath, ".html") {
		urlPath = strings.TrimSuffix(urlPath, ".html")
		// Handle special case for index.html at root
		if urlPath == "index" {
			return "/"
		}
		return "/" + urlPath
	}

	return ""
}

// markdownFilePathToURL converts a markdown file path to a clean URL path using the same logic as the router
func (s *SEOService) markdownFilePathToURL(filePath, baseDir string) string {
	// Remove base directory prefix and .md extension
	urlPath := strings.TrimPrefix(filePath, baseDir)
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, ".md")

	if urlPath == "" {
		return ""
	}

	// Handle index files (remove /index suffix)
	if strings.HasSuffix(urlPath, "/index") {
		urlPath = strings.TrimSuffix(urlPath, "/index")
		if urlPath == "" {
			urlPath = baseDir // For content section root
		}
	}

	// Apply content type URL mapping
	contentType, found := s.getContentTypeFromBaseDir(baseDir)
	if found {
		// Clean the path segments by removing numeric prefixes
		cleanPath := s.cleanPathSegments(urlPath)
		return contentType.URLPrefix + "/" + cleanPath
	}

	// Default fallback - clean the path and use as-is
	cleanPath := s.cleanPathSegments(urlPath)
	return "/" + cleanPath
}

// getContentTypeFromBaseDir maps base directories to their content types
func (s *SEOService) getContentTypeFromBaseDir(baseDir string) (ContentType, bool) {
	for _, contentType := range ContentTypes {
		if strings.HasSuffix(baseDir, contentType.BaseDir) {
			return contentType, true
		}
	}
	return ContentType{}, false
}

// cleanPathSegments removes numeric prefixes from URL path segments
func (s *SEOService) cleanPathSegments(urlPath string) string {
	if urlPath == "" {
		return ""
	}

	// Split path into segments and clean each one
	segments := strings.Split(urlPath, "/")
	cleanedSegments := make([]string, 0, len(segments))

	for _, segment := range segments {
		if segment != "" {
			// Use CleanID to remove numeric prefixes and normalize
			cleaned := CleanID(segment)
			if cleaned != "" {
				cleanedSegments = append(cleanedSegments, cleaned)
			}
		}
	}

	return strings.Join(cleanedSegments, "/")
}

// getURLProperties returns priority and change frequency based on URL path
func (s *SEOService) getURLProperties(urlPath string) (priority, changeFreq string) {
	// Default values
	priority = "0.5"
	changeFreq = "monthly"

	// High priority pages
	switch urlPath {
	case "/":
		priority = "1.0"
		changeFreq = "weekly"
	case "/pricing", "/platform", "/features":
		priority = "0.9"
		changeFreq = "weekly"
	case "/contact", "/about":
		priority = "0.8"
		changeFreq = "monthly"
	}

	// Category-based priorities
	if strings.HasPrefix(urlPath, "/platform/") {
		priority = "0.8"
		changeFreq = "weekly"
	} else if strings.HasPrefix(urlPath, "/solutions/") {
		priority = "0.7"
		changeFreq = "monthly"
	} else if strings.HasPrefix(urlPath, "/docs/") {
		priority = "0.6"
		changeFreq = "weekly"
	} else if strings.HasPrefix(urlPath, "/api/") {
		priority = "0.6"
		changeFreq = "weekly"
	} else if strings.HasPrefix(urlPath, "/insights/") {
		priority = "0.5"
		changeFreq = "monthly"
	} else if strings.HasPrefix(urlPath, "/company-news/") || strings.HasPrefix(urlPath, "/product-updates/") {
		priority = "0.4"
		changeFreq = "weekly"
	}

	return priority, changeFreq
}

// scanHTMLPagesForLanguage scans HTML pages for a specific language
func (s *SEOService) scanHTMLPagesForLanguage(pagesDir, baseURL, currentTime, lang string) ([]URLEntry, error) {
	var urls []URLEntry

	err := filepath.Walk(pagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// Convert file path to URL path
			urlPath := s.filePathToURL(path, pagesDir)
			if urlPath == "" {
				return nil // Skip invalid paths
			}

			// Build language-specific URL
			if lang != DefaultLanguage {
				urlPath = "/" + lang + urlPath
			}

			// Determine priority and change frequency based on path
			priority, changeFreq := s.getURLProperties(urlPath)

			urls = append(urls, URLEntry{
				Loc:        baseURL + urlPath,
				LastMod:    currentTime,
				ChangeFreq: changeFreq,
				Priority:   priority,
			})
		}

		return nil
	})

	return urls, err
}

// scanMarkdownContentForLanguage scans markdown content for a specific language
func (s *SEOService) scanMarkdownContentForLanguage(contentDir, baseURL, currentTime, lang string) ([]URLEntry, error) {
	var urls []URLEntry

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			// Convert file path to URL path
			urlPath := s.markdownFilePathToURL(path, contentDir)
			if urlPath == "" {
				return nil // Skip invalid paths
			}

			// Build language-specific URL
			if lang != DefaultLanguage {
				urlPath = "/" + lang + urlPath
			}

			// Determine priority and change frequency based on path
			priority, changeFreq := s.getURLProperties(urlPath)

			// Try to get actual modification time from file
			lastMod := currentTime
			if stat, err := os.Stat(path); err == nil {
				lastMod = stat.ModTime().Format("2006-01-02")
			}

			urls = append(urls, URLEntry{
				Loc:        baseURL + urlPath,
				LastMod:    lastMod,
				ChangeFreq: changeFreq,
				Priority:   priority,
			})
		}

		return nil
	})

	return urls, err
}
