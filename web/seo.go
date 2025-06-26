package web

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"

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
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Language    string `json:"language"`
	Author      string `json:"author"`
}

// MetadataDefaults represents default metadata values
type MetadataDefaults struct {
	TitleSuffix string   `json:"title_suffix"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

// Metadata holds the complete metadata structure
type Metadata struct {
	Site     SiteMetadata                `json:"site"`
	Pages    map[string]PageMetadata     `json:"pages"`
	Defaults MetadataDefaults            `json:"defaults"`
}

// Frontmatter represents markdown file frontmatter
type Frontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
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

// PageData holds data for template rendering
type PageData struct {
	Title       string
	Content     template.HTML
	Navigation  *Navigation
	PageMeta    *PageMetadata
	SiteMeta    *SiteMetadata
	Description string
	Keywords    []string
	IsMarkdown  bool
	Frontmatter *Frontmatter
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

// PreparePageData creates PageData with metadata for the given path
func (s *SEOService) PreparePageData(path string, content template.HTML, isMarkdown bool, frontmatter *Frontmatter, navigation *Navigation) PageData {
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
			// Check if specific page metadata exists
			if meta, exists := s.metadata.Pages[pageKey]; exists {
				pageMeta = &meta
				if title == "" {
					title = meta.Title
				}
				if description == "" {
					description = meta.Description
				}
				keywords = meta.Keywords
			} else {
				// Use defaults
				if title == "" {
					title = s.getFallbackTitle(path) + s.metadata.Defaults.TitleSuffix
				}
				if description == "" {
					description = s.metadata.Defaults.Description
				}
				keywords = s.metadata.Defaults.Keywords
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
	
	var siteMeta *SiteMetadata
	if s.metadata != nil {
		siteMeta = &s.metadata.Site
	}
	
	return PageData{
		Title:       title,
		Content:     content,
		Navigation:  navigation,
		PageMeta:    pageMeta,
		SiteMeta:    siteMeta,
		Description: description,
		Keywords:    keywords,
		IsMarkdown:  isMarkdown,
		Frontmatter: frontmatter,
	}
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