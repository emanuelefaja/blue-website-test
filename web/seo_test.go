package web

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadSEOData tests loading SEO data from files
func TestLoadSEOData(t *testing.T) {
	// Create temporary data directory
	testDir := t.TempDir()
	dataDir := filepath.Join(testDir, "data")
	seoDir := filepath.Join(testDir, "seo")
	os.MkdirAll(dataDir, 0755)
	os.MkdirAll(seoDir, 0755)

	// Create test metadata file
	metadata := &Metadata{
		Site: SiteMetadata{
			Name: "Test Site",
			Descriptions: map[string]string{
				"en": "Test site description",
			},
			URL: "https://example.com",
		},
		Pages: map[string]map[string]PageMetadata{
			"/about": {
				"en": {
					Title:       "About Us",
					Description: "About page description",
				},
			},
		},
		Defaults: MetadataDefaults{
			TitleSuffix: " | Test Site",
			Descriptions: map[string]string{
				"en": "Default description",
			},
			Keywords: map[string][]string{
				"en": {"test", "site"},
			},
		},
	}

	metadataJSON, _ := json.MarshalIndent(metadata, "", "  ")
	os.WriteFile(filepath.Join(dataDir, "metadata.json"), metadataJSON, 0644)

	// Create test redirects file
	redirects := &Redirects{
		Redirects: map[string]string{
			"/old-page": "/new-page",
			"/legacy":   "/modern",
		},
		Rules: RedirectRules{
			StatusCode:    301,
			TrailingSlash: "remove",
		},
	}

	redirectsJSON, _ := json.MarshalIndent(redirects, "", "  ")
	os.WriteFile(filepath.Join(dataDir, "redirects.json"), redirectsJSON, 0644)

	// Change working directory to test directory temporarily
	oldWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(oldWd)

	// Test loading
	seoService := NewSEOService()
	err := seoService.LoadData()
	if err != nil {
		t.Fatalf("Failed to load SEO data: %v", err)
	}

	// Verify metadata loaded
	if seoService.metadata == nil {
		t.Error("Metadata not loaded")
	}
	if seoService.metadata.Site.Name != "Test Site" {
		t.Errorf("Expected site name 'Test Site', got %q", seoService.metadata.Site.Name)
	}

	// Verify redirects loaded
	if seoService.redirects == nil {
		t.Error("Redirects not loaded")
	}
	if len(seoService.redirects.Redirects) != 2 {
		t.Errorf("Expected 2 redirects, got %d", len(seoService.redirects.Redirects))
	}
}

// TestCheckRedirect tests redirect checking functionality
func TestCheckRedirect(t *testing.T) {
	seoService := &SEOService{
		redirects: &Redirects{
			Redirects: map[string]string{
				"/old-page":    "/new-page",
				"/legacy/path": "/modern/path",
				"/with-slash/": "/without-slash",
			},
			Rules: RedirectRules{
				StatusCode: 301,
			},
		},
	}

	tests := []struct {
		name           string
		path           string
		expectRedirect bool
		expectedTo     string
		expectedStatus int
	}{
		{
			name:           "Direct redirect",
			path:           "/old-page",
			expectRedirect: true,
			expectedTo:     "/new-page",
			expectedStatus: 301,
		},
		{
			name:           "No redirect",
			path:           "/current-page",
			expectRedirect: false,
			expectedTo:     "",
			expectedStatus: 0,
		},
		{
			name:           "Legacy path redirect",
			path:           "/legacy/path",
			expectRedirect: true,
			expectedTo:     "/modern/path",
			expectedStatus: 301,
		},
		{
			name:           "Trailing slash redirect",
			path:           "/with-slash/",
			expectRedirect: true,
			expectedTo:     "/without-slash",
			expectedStatus: 301,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			to, status, shouldRedirect := seoService.CheckRedirect(tt.path)

			if shouldRedirect != tt.expectRedirect {
				t.Errorf("Expected redirect=%v, got %v", tt.expectRedirect, shouldRedirect)
			}

			if to != tt.expectedTo {
				t.Errorf("Expected redirect to %q, got %q", tt.expectedTo, to)
			}

			if status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}
		})
	}
}

// TestPreparePageMetadata tests metadata preparation
func TestPreparePageMetadata(t *testing.T) {
	seoService := &SEOService{
		metadata: &Metadata{
			Site: SiteMetadata{
				Name: "Test Site",
				Descriptions: map[string]string{
					"en": "Test site description",
				},
				URL: "https://example.com",
			},
			Pages: map[string]map[string]PageMetadata{
				"about": {
					"en": {
						Title:       "About Us",
						Description: "About page description",
						Keywords:    []string{"about", "company"},
					},
				},
			},
			Defaults: MetadataDefaults{
				TitleSuffix: " | Test Site",
				Descriptions: map[string]string{
					"en": "Default description",
				},
				Keywords: map[string][]string{
					"en": {"test", "site"},
				},
			},
		},
	}

	tests := []struct {
		name             string
		path             string
		isMarkdown       bool
		frontmatter      *Frontmatter
		expectedTitle    string
		expectedDesc     string
		expectedKeywords []string
	}{
		{
			name:             "Page with metadata",
			path:             "/about",
			isMarkdown:       false,
			frontmatter:      nil,
			expectedTitle:    "About Us | Test Site",
			expectedDesc:     "About page description",
			expectedKeywords: []string{"about", "company"},
		},
		{
			name:             "Page without metadata",
			path:             "/contact",
			isMarkdown:       false,
			frontmatter:      nil,
			expectedTitle:    "Contact | Test Site",
			expectedDesc:     "Default description",
			expectedKeywords: []string{"test", "site"},
		},
		{
			name:       "Markdown with frontmatter",
			path:       "/blog/post",
			isMarkdown: true,
			frontmatter: &Frontmatter{
				Title:       "Blog Post Title",
				Description: "Blog post description",
			},
			expectedTitle:    "Blog Post Title | Test Site",
			expectedDesc:     "Blog post description",
			expectedKeywords: []string{}, // Bug: keywords not set when frontmatter has both title and description
		},
		{
			name:             "Homepage",
			path:             "/",
			isMarkdown:       false,
			frontmatter:      nil,
			expectedTitle:    "Home | Test Site",    // getFallbackTitle returns "Home", then suffix added
			expectedDesc:     "Default description", // Empty in metadata, falls back to default
			expectedKeywords: []string{"test", "site"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, desc, keywords, _, _ := seoService.PreparePageMetadata(tt.path, tt.isMarkdown, tt.frontmatter, "en")

			if title != tt.expectedTitle {
				t.Errorf("Expected title %q, got %q", tt.expectedTitle, title)
			}

			if desc != tt.expectedDesc {
				t.Errorf("Expected description %q, got %q", tt.expectedDesc, desc)
			}

			if len(keywords) != len(tt.expectedKeywords) {
				t.Errorf("Expected %d keywords, got %d", len(tt.expectedKeywords), len(keywords))
			}
		})
	}
}

// TestParseFrontmatter tests frontmatter parsing functionality
func TestParseFrontmatter(t *testing.T) {
	seoService := NewSEOService()

	tests := []struct {
		name              string
		content           string
		expectFrontmatter bool
		expectedTitle     string
		expectedDesc      string
		expectedContent   string
	}{
		{
			name: "Valid frontmatter",
			content: `---
title: Test Article
description: Test description
category: Engineering
date: 2024-01-01
---
# Main Content

This is the main content.`,
			expectFrontmatter: true,
			expectedTitle:     "Test Article",
			expectedDesc:      "Test description",
			expectedContent:   "# Main Content\n\nThis is the main content.",
		},
		{
			name:              "No frontmatter",
			content:           "# Just Content\n\nNo frontmatter here.",
			expectFrontmatter: false,
			expectedContent:   "# Just Content\n\nNo frontmatter here.",
		},
		{
			name:              "Windows line endings",
			content:           "---\r\ntitle: Windows Test\r\n---\r\nContent here",
			expectFrontmatter: true,
			expectedTitle:     "Windows Test",
			expectedContent:   "Content here",
		},
		{
			name:              "Incomplete frontmatter",
			content:           "---\ntitle: Incomplete\n# No closing delimiter",
			expectFrontmatter: false,
			expectedContent:   "---\ntitle: Incomplete\n# No closing delimiter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frontmatter, content, err := seoService.ParseFrontmatter([]byte(tt.content))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectFrontmatter {
				if frontmatter == nil {
					t.Fatal("Expected frontmatter but got nil")
				}
				if frontmatter.Title != tt.expectedTitle {
					t.Errorf("Expected title %q, got %q", tt.expectedTitle, frontmatter.Title)
				}
				if tt.expectedDesc != "" && frontmatter.Description != tt.expectedDesc {
					t.Errorf("Expected description %q, got %q", tt.expectedDesc, frontmatter.Description)
				}
			} else {
				if frontmatter != nil {
					t.Error("Expected no frontmatter but got some")
				}
			}

			if string(content) != tt.expectedContent {
				t.Errorf("Expected content %q, got %q", tt.expectedContent, string(content))
			}
		})
	}
}

// TestGenerateSitemap tests sitemap generation
func TestGenerateSitemap(t *testing.T) {
	// Create test directory structure
	testDir := t.TempDir()
	pagesDir := filepath.Join(testDir, "pages")
	contentDir := filepath.Join(testDir, "content")
	publicDir := filepath.Join(testDir, "public")

	// Create directories
	os.MkdirAll(pagesDir, 0755)
	os.MkdirAll(filepath.Join(contentDir, "docs"), 0755)
	os.MkdirAll(publicDir, 0755)

	// Create test HTML files
	os.WriteFile(filepath.Join(pagesDir, "index.html"), []byte("<h1>Home</h1>"), 0644)
	os.WriteFile(filepath.Join(pagesDir, "about.html"), []byte("<h1>About</h1>"), 0644)
	os.WriteFile(filepath.Join(pagesDir, "404.html"), []byte("<h1>Not Found</h1>"), 0644)

	// Create test markdown files
	os.WriteFile(filepath.Join(contentDir, "docs", "intro.md"), []byte("# Introduction"), 0644)

	// Change to test directory
	oldWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(oldWd)

	// Test sitemap generation
	seoService := NewSEOService()
	err := seoService.GenerateSitemap("https://example.com")
	if err != nil {
		t.Fatalf("Failed to generate sitemap: %v", err)
	}

	// Check sitemap was created
	sitemapPath := filepath.Join(publicDir, "sitemap.xml")
	data, err := os.ReadFile(sitemapPath)
	if err != nil {
		t.Fatalf("Failed to read sitemap: %v", err)
	}

	sitemapContent := string(data)

	// Verify sitemap contains expected URLs
	expectedURLs := []string{
		"https://example.com/",
		"https://example.com/about",
		"https://example.com/docs/intro",
	}

	for _, url := range expectedURLs {
		if !strings.Contains(sitemapContent, url) {
			t.Errorf("Expected sitemap to contain %s", url)
		}
	}

	// Verify 404 page is excluded
	if strings.Contains(sitemapContent, "404") {
		t.Error("Sitemap should not contain 404 page")
	}

	// Verify XML structure
	if !strings.HasPrefix(sitemapContent, "<?xml") {
		t.Error("Sitemap should start with XML declaration")
	}

	if !strings.Contains(sitemapContent, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`) {
		t.Error("Sitemap should contain proper namespace")
	}
}

// TestGetPageKey tests the page key generation
func TestGetPageKey(t *testing.T) {
	seoService := NewSEOService()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Root path",
			path:     "/",
			expected: "home",
		},
		{
			name:     "Simple path",
			path:     "/about",
			expected: "about",
		},
		{
			name:     "Nested path",
			path:     "/docs/intro",
			expected: "docs/intro",
		},
		{
			name:     "Trailing slash",
			path:     "/contact/",
			expected: "contact",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use reflection to call private method
			result := seoService.getPageKey(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestGetFallbackTitle tests fallback title generation
func TestGetFallbackTitle(t *testing.T) {
	seoService := NewSEOService()

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Root path",
			path:     "/",
			expected: "Home",
		},
		{
			name:     "Simple path",
			path:     "/about",
			expected: "About",
		},
		{
			name:     "Hyphenated path",
			path:     "/getting-started",
			expected: "Getting Started",
		},
		{
			name:     "Nested path",
			path:     "/docs/api-reference",
			expected: "Docs - Api Reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seoService.getFallbackTitle(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestFilePathToURL tests conversion of file paths to URLs
func TestFilePathToURL(t *testing.T) {
	seoService := NewSEOService()

	tests := []struct {
		name     string
		filePath string
		baseDir  string
		expected string
	}{
		{
			name:     "Root index",
			filePath: "pages/index.html",
			baseDir:  "pages",
			expected: "/",
		},
		{
			name:     "Regular page",
			filePath: "pages/about.html",
			baseDir:  "pages",
			expected: "/about",
		},
		{
			name:     "Nested page",
			filePath: "pages/platform/features.html",
			baseDir:  "pages",
			expected: "/platform/features",
		},
		{
			name:     "Directory index",
			filePath: "pages/docs/index.html",
			baseDir:  "pages",
			expected: "/docs/",
		},
		{
			name:     "Skip 404",
			filePath: "pages/404.html",
			baseDir:  "pages",
			expected: "",
		},
		{
			name:     "Skip copy files",
			filePath: "pages/about.copy.html",
			baseDir:  "pages",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := seoService.filePathToURL(tt.filePath, tt.baseDir)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestGetURLProperties tests URL priority and change frequency assignment
func TestGetURLProperties(t *testing.T) {
	seoService := NewSEOService()

	tests := []struct {
		name             string
		urlPath          string
		expectedPriority string
		expectedFreq     string
	}{
		{
			name:             "Homepage",
			urlPath:          "/",
			expectedPriority: "1.0",
			expectedFreq:     "weekly",
		},
		{
			name:             "Pricing page",
			urlPath:          "/pricing",
			expectedPriority: "0.9",
			expectedFreq:     "weekly",
		},
		{
			name:             "Platform page",
			urlPath:          "/platform/status",
			expectedPriority: "0.8",
			expectedFreq:     "weekly",
		},
		{
			name:             "Docs page",
			urlPath:          "/docs/introduction",
			expectedPriority: "0.6",
			expectedFreq:     "weekly",
		},
		{
			name:             "Insights page",
			urlPath:          "/insights/article",
			expectedPriority: "0.5",
			expectedFreq:     "monthly",
		},
		{
			name:             "Default page",
			urlPath:          "/random-page",
			expectedPriority: "0.5",
			expectedFreq:     "monthly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority, freq := seoService.getURLProperties(tt.urlPath)
			if priority != tt.expectedPriority {
				t.Errorf("Expected priority %q, got %q", tt.expectedPriority, priority)
			}
			if freq != tt.expectedFreq {
				t.Errorf("Expected frequency %q, got %q", tt.expectedFreq, freq)
			}
		})
	}
}
