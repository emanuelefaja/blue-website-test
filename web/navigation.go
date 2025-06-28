package web

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

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

// NavigationService handles all navigation-related operations
type NavigationService struct {
	navigation      *Navigation
	docsNavigation  *Navigation
	apiNavigation   *Navigation
	legalNavigation *Navigation
	seoService      *SEOService
}

// NewNavigationService creates a new navigation service
func NewNavigationService(seoService *SEOService) *NavigationService {
	service := &NavigationService{
		seoService: seoService,
	}

	// Load static navigation
	if err := service.LoadNavigation(); err != nil {
		log.Printf("Error loading navigation: %v", err)
	}

	// Generate dynamic navigations
	service.GenerateDynamicNavigations()

	return service
}

// LoadNavigation loads navigation data from JSON file
func (ns *NavigationService) LoadNavigation() error {
	data, err := os.ReadFile("data/nav.json")
	if err != nil {
		return err
	}

	ns.navigation = &Navigation{}
	return json.Unmarshal(data, ns.navigation)
}

// GenerateDynamicNavigations generates all dynamic navigation structures
func (ns *NavigationService) GenerateDynamicNavigations() {
	// Generate for each content type
	for key, contentType := range ContentTypes {
		nav, err := ns.GenerateContentNavigation(contentType.BaseDir, contentType.URLPrefix)
		if err != nil {
			log.Printf("Error generating %s navigation: %v", key, err)
			continue
		}

		switch key {
		case "docs":
			ns.docsNavigation = nav
		case "api":
			ns.apiNavigation = nav
		case "legal":
			ns.legalNavigation = nav
		}
	}
}

// GenerateContentNavigation creates navigation tree from content directory
func (ns *NavigationService) GenerateContentNavigation(contentDir, baseURL string) (*Navigation, error) {
	var sections []NavItem

	// Read the content directory
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}

	// Process each directory/file
	for _, entry := range entries {
		if entry.IsDir() {
			navItem, err := ns.processDirectory(filepath.Join(contentDir, entry.Name()), entry.Name(), baseURL, contentDir)
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
			fileTitle := CleanTitle(fileName)

			// Try to get title from frontmatter
			if filePath := filepath.Join(contentDir, entry.Name()); filePath != "" {
				if data, err := os.ReadFile(filePath); err == nil {
					if frontmatter, _, err := ns.seoService.ParseFrontmatter(data); err == nil && frontmatter != nil && frontmatter.Title != "" {
						fileTitle = frontmatter.Title
					}
				}
			}

			href := baseURL + "/" + CleanID(fileName)

			sections = append(sections, NavItem{
				ID:         CleanID(fileName),
				Name:       fileTitle,
				Href:       href,
				OriginalID: fileName,
			})
		}
	}

	// Sort sections by numeric prefix
	ns.sortNavItems(sections)

	return &Navigation{Sections: sections}, nil
}

// processDirectory processes a content directory and creates NavItem
func (ns *NavigationService) processDirectory(dirPath, dirName, baseURL, rootContentDir string) (*NavItem, error) {
	// Read directory metadata
	title := CleanTitle(dirName)
	dirMetaPath := filepath.Join(dirPath, "_dir.yml")
	if data, err := os.ReadFile(dirMetaPath); err == nil {
		var dirMeta DirMetadata
		if err := yaml.Unmarshal(data, &dirMeta); err == nil && dirMeta.Title != "" {
			title = dirMeta.Title
		}
	}

	// Create nav item
	navItem := &NavItem{
		ID:         CleanID(dirName),
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
			childNav, err := ns.processDirectory(filepath.Join(dirPath, entry.Name()), entry.Name(), baseURL, rootContentDir)
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
			fileTitle := CleanTitle(fileName)

			// Try to get title from frontmatter
			if filePath := filepath.Join(dirPath, entry.Name()); filePath != "" {
				if data, err := os.ReadFile(filePath); err == nil {
					if frontmatter, _, err := ns.seoService.ParseFrontmatter(data); err == nil && frontmatter != nil && frontmatter.Title != "" {
						fileTitle = frontmatter.Title
					}
				}
			}

			// Create relative path for href
			relDir := strings.TrimPrefix(dirPath, rootContentDir+"/")

			// Get content type from root directory
			contentTypeName := strings.Split(rootContentDir, "/")[1] // e.g., "content/docs" -> "docs"
			
			// Remove the content type prefix
			if strings.HasPrefix(relDir, contentTypeName+"/") {
				relDir = strings.TrimPrefix(relDir, contentTypeName+"/")
			} else if contentTypeName == "api-docs" && strings.HasPrefix(relDir, "api-docs/") {
				relDir = strings.TrimPrefix(relDir, "api-docs/")
			} else if contentTypeName == "legal" && strings.HasPrefix(relDir, "legal/") {
				relDir = strings.TrimPrefix(relDir, "legal/")
			}

			// Clean numeric prefixes from directory path
			relDir = CleanDirectoryPath(relDir)

			href := baseURL + "/" + relDir + "/" + CleanID(fileName)

			children = append(children, NavItem{
				ID:         CleanID(fileName),
				Name:       fileTitle,
				Href:       href,
				OriginalID: fileName,
			})
		}
	}

	// Sort children by numeric prefix
	ns.sortNavItems(children)

	if len(children) > 0 {
		navItem.Children = children
	}

	return navItem, nil
}

// sortNavItems sorts navigation items by numeric prefix
func (ns *NavigationService) sortNavItems(items []NavItem) {
	// Sort by extracting numeric prefix from original names
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			// Extract numeric prefixes for comparison using original IDs
			num1 := ExtractNumericPrefix(items[i].OriginalID)
			num2 := ExtractNumericPrefix(items[j].OriginalID)

			if num1 > num2 {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// GetNavigationForPath returns the appropriate navigation based on the URL path
func (ns *NavigationService) GetNavigationForPath(path string) *Navigation {
	// Always start with static navigation
	if ns.navigation == nil {
		return &Navigation{}
	}

	// Make a copy of the static navigation
	nav := &Navigation{
		Sections: make([]NavItem, len(ns.navigation.Sections)),
	}
	copy(nav.Sections, ns.navigation.Sections)

	// Always add Documentation section if available
	if ns.docsNavigation != nil {
		docSection := NavItem{
			ID:       "documentation",
			Name:     "Documentation",
			Expanded: strings.HasPrefix(path, "/docs"), // Only expand when on docs pages
			Children: ns.docsNavigation.Sections,
		}
		nav.Sections = append(nav.Sections, docSection)
	}

	// Always add API Reference section if available
	if ns.apiNavigation != nil {
		apiSection := NavItem{
			ID:       "api-reference",
			Name:     "API Reference",
			Expanded: strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/api-docs"), // Only expand when on API pages
			Children: ns.apiNavigation.Sections,
		}
		nav.Sections = append(nav.Sections, apiSection)
	}

	// Always add Legal section if available (placed at end for bottom positioning)
	if ns.legalNavigation != nil {
		legalSection := NavItem{
			ID:       "legal",
			Name:     "Legal",
			Expanded: strings.HasPrefix(path, "/legal"), // Only expand when on legal pages
			Children: ns.legalNavigation.Sections,
		}
		nav.Sections = append(nav.Sections, legalSection)
	}

	return nav
}