package web

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestSearchItemStructure tests the SearchItem struct
func TestSearchItemStructure(t *testing.T) {
	item := SearchItem{
		Title:       "Test Page",
		Description: "Test description",
		Content:     "Test content",
		URL:         "/test",
		Type:        "page",
		Section:     "docs",
		Category:    "Feature",
	}

	if item.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got %q", item.Title)
	}
}

// TestSearchFrontmatterStructure tests the SearchFrontmatter struct
func TestSearchFrontmatterStructure(t *testing.T) {
	fm := SearchFrontmatter{
		Title:       "Test Title",
		Description: "Test description",
		Category:    "Test category",
	}

	if fm.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %q", fm.Title)
	}
}

// TestSplitFrontmatter tests the splitFrontmatter function
func TestSplitFrontmatter(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expectedFM      string
		expectedContent string
	}{
		{
			name: "Valid frontmatter with content",
			content: `---
title: Test Page
description: A test page
---
# Main Content

This is the body.`,
			expectedFM: `title: Test Page
description: A test page`,
			expectedContent: `# Main Content

This is the body.`,
		},
		{
			name:            "No frontmatter",
			content:         "# Just Content\n\nNo frontmatter here.",
			expectedFM:      "",
			expectedContent: "# Just Content\n\nNo frontmatter here.",
		},
		{
			name:            "Frontmatter with Windows line endings",
			content:         "---\r\ntitle: Test\r\n---\r\nContent",
			expectedFM:      "title: Test",
			expectedContent: "Content",
		},
		{
			name:            "Invalid frontmatter (no closing delimiter)",
			content:         "---\ntitle: Test\nContent without closing",
			expectedFM:      "",
			expectedContent: "---\ntitle: Test\nContent without closing",
		},
		{
			name:            "Empty content",
			content:         "",
			expectedFM:      "",
			expectedContent: "",
		},
		{
			name: "Frontmatter with triple dash in content",
			content: `---
title: Test
---
# Content

Some text --- more text`,
			expectedFM: "title: Test",
			expectedContent: `# Content

Some text --- more text`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, content := splitFrontmatter(tt.content)
			if fm != tt.expectedFM {
				t.Errorf("Expected frontmatter %q, got %q", tt.expectedFM, fm)
			}
			if content != tt.expectedContent {
				t.Errorf("Expected content %q, got %q", tt.expectedContent, content)
			}
		})
	}
}

// TestExtractHTMLTitle tests the extractHTMLTitle function
func TestExtractHTMLTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Valid title tag",
			html:     `<html><head><title>Test Page Title</title></head><body>Content</body></html>`,
			expected: "Test Page Title",
		},
		{
			name:     "No title tag",
			html:     `<html><head></head><body>Content</body></html>`,
			expected: "",
		},
		{
			name:     "Empty title tag",
			html:     `<html><head><title></title></head><body>Content</body></html>`,
			expected: "",
		},
		{
			name:     "Title with special characters",
			html:     `<title>Test & Page | Blue™</title>`,
			expected: "Test & Page | Blue™",
		},
		{
			name:     "Multiple title tags (use first)",
			html:     `<title>First Title</title><title>Second Title</title>`,
			expected: "First Title",
		},
		{
			name:     "Unclosed title tag",
			html:     `<title>Unclosed Title`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHTMLTitle(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestExtractH1Title tests the extractH1Title function
func TestExtractH1Title(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple H1",
			html:     `<h1>Main Title</h1>`,
			expected: "Main Title",
		},
		{
			name:     "H1 with attributes",
			html:     `<h1 class="title" id="main">Main Title</h1>`,
			expected: "Main Title",
		},
		{
			name:     "H1 with nested elements",
			html:     `<h1>Main <span>Title</span> Text</h1>`,
			expected: "Main Title Text",
		},
		{
			name:     "Multiple H1 tags (use first)",
			html:     `<h1>First H1</h1><h1>Second H1</h1>`,
			expected: "First H1",
		},
		{
			name:     "No H1 tag",
			html:     `<h2>Not H1</h2><p>Content</p>`,
			expected: "",
		},
		{
			name:     "Case insensitive H1",
			html:     `<H1>UPPERCASE H1</H1>`,
			expected: "UPPERCASE H1",
		},
		{
			name: "H1 with line breaks",
			html: `<h1>
				Multi
				Line
				Title
			</h1>`,
			expected: "Multi Line Title",
		},
		{
			name:     "Unclosed H1 tag",
			html:     `<h1>Unclosed Title`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractH1Title(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestGenerateTitleFromFilename tests the generateTitleFromFilename function
func TestGenerateTitleFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Simple filename",
			filePath: "about.html",
			expected: "About",
		},
		{
			name:     "Filename with hyphens",
			filePath: "getting-started.html",
			expected: "Getting Started",
		},
		{
			name:     "Filename with underscores",
			filePath: "user_guide_intro.html",
			expected: "User Guide Intro",
		},
		{
			name:     "Filename with path",
			filePath: "/pages/docs/api-reference.html",
			expected: "Api Reference",
		},
		{
			name:     "Filename with multiple extensions",
			filePath: "template.min.html",
			expected: "Template.min",
		},
		{
			name:     "All lowercase filename",
			filePath: "pricing.html",
			expected: "Pricing",
		},
		{
			name:     "Mixed case filename",
			filePath: "TeamManagement.html",
			expected: "Teammanagement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTitleFromFilename(tt.filePath)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestGenerateTitleFromURL tests the generateTitleFromURL function
func TestGenerateTitleFromURL(t *testing.T) {
	tests := []struct {
		name     string
		urlPath  string
		expected string
	}{
		{
			name:     "Simple URL",
			urlPath:  "/about",
			expected: "About",
		},
		{
			name:     "Nested URL",
			urlPath:  "/docs/getting-started",
			expected: "Getting Started",
		},
		{
			name:     "URL with trailing slash",
			urlPath:  "/platform/features/",
			expected: "Features",
		},
		{
			name:     "Root URL",
			urlPath:  "/",
			expected: "Home",
		},
		{
			name:     "Empty URL",
			urlPath:  "",
			expected: "Home",
		},
		{
			name:     "URL with underscores",
			urlPath:  "/api/user_management",
			expected: "User Management",
		},
		{
			name:     "Deep nested URL",
			urlPath:  "/docs/guides/advanced/configuration",
			expected: "Configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTitleFromURL(tt.urlPath)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestGetPageKey is already defined in seo_test.go, so we skip it here

// TestExtractTextFromHTML tests the extractTextFromHTML function
func TestExtractTextFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "HTML with H1 and content",
			html: `<html>
				<head><title>Test</title></head>
				<body>
					<nav>Navigation</nav>
					<h1>Main Title</h1>
					<p>This is a paragraph.</p>
					<p>Another paragraph.</p>
				</body>
			</html>`,
			expected: "Main Title This is a paragraph. Another paragraph.",
		},
		{
			name: "HTML with scripts (should be removed)",
			html: `<h1>Title</h1>
				<script>console.log('test');</script>
				<p>Content</p>
				<script type="text/javascript">
					function test() { return true; }
				</script>`,
			expected: "Title Content",
		},
		{
			name: "HTML with styles (should be removed)",
			html: `<h1>Title</h1>
				<style>.test { color: red; }</style>
				<p>Content</p>
				<style type="text/css">
					body { margin: 0; }
				</style>`,
			expected: "Title Content",
		},
		{
			name: "HTML with AlpineJS directives",
			html: `<h1>Title</h1>
				<div x-data="{ open: false }" x-show="open" @click="toggle()">
					<p>Interactive content</p>
				</div>`,
			expected: "Title Interactive content",
		},
		{
			name: "HTML without H1",
			html: `<html>
				<body>
					<h2>Subtitle</h2>
					<p>Content without H1</p>
				</body>
			</html>`,
			expected: "",
		},
		{
			name: "HTML with nested tags",
			html: `<h1>Title</h1>
				<div>
					<section>
						<article>
							<p>Deeply <strong>nested</strong> content</p>
						</article>
					</section>
				</div>`,
			expected: "Title Deeply nested content",
		},
		{
			name: "HTML with special characters and entities",
			html: `<h1>Title &amp; More</h1>
				<p>Content with &copy; and &trade; symbols</p>`,
			expected: "Title &amp; More Content with &copy; and &trade; symbols",
		},
		{
			name: "HTML with JavaScript patterns",
			html: `<h1>Title</h1>
				<p>Content with function() { return true; } and const x = 5;</p>`,
			expected: "Title Content with true; and x 5;",
		},
		{
			name: "HTML with unclosed script tag",
			html: `<h1>Title</h1>
				<script>console.log('test'
				<p>This should not appear</p>`,
			expected: "Title",
		},
		{
			name: "HTML with SVG elements",
			html: `<h1>Title</h1>
				<svg><path d="M10 10"/></svg>
				<p>Content after SVG</p>`,
			expected: "Title Content after SVG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTextFromHTML(tt.html)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestExtractPageTitle tests the extractPageTitle function
func TestExtractPageTitle(t *testing.T) {
	// Create test metadata
	metadata := &Metadata{
		Pages: map[string]map[string]PageMetadata{
			"about": {
				"en": {
					Title:       "About Blue",
					Description: "Learn about Blue",
				},
			},
			"platform/features": {
				"en": {
					Title:       "Platform Features",
					Description: "Feature overview",
				},
			},
		},
	}

	tests := []struct {
		name     string
		html     string
		url      string
		filePath string
		metadata *Metadata
		expected string
	}{
		{
			name:     "Title from metadata",
			html:     `<h1>Different Title</h1>`,
			url:      "/about",
			filePath: "pages/about.html",
			metadata: metadata,
			expected: "About Blue",
		},
		{
			name:     "Title from H1 when no metadata",
			html:     `<h1>H1 Title</h1>`,
			url:      "/contact",
			filePath: "pages/contact.html",
			metadata: metadata,
			expected: "H1 Title",
		},
		{
			name:     "Title from HTML title tag",
			html:     `<title>HTML Title</title><body>No H1</body>`,
			url:      "/test",
			filePath: "pages/test.html",
			metadata: nil,
			expected: "HTML Title",
		},
		{
			name:     "Title from filename as fallback",
			html:     `<body>No title anywhere</body>`,
			url:      "/no-title",
			filePath: "pages/no-title.html",
			metadata: nil,
			expected: "No Title",
		},
		{
			name:     "Nested path with metadata",
			html:     `<h1>Features</h1>`,
			url:      "/platform/features",
			filePath: "pages/platform/features.html",
			metadata: metadata,
			expected: "Platform Features",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPageTitle(tt.html, tt.url, tt.filePath, tt.metadata)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestLoadMetadata tests the loadMetadata function
func TestLoadMetadata(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create data directory
	os.MkdirAll("data", 0755)

	t.Run("Valid metadata file", func(t *testing.T) {
		// Create test metadata.json
		metadataContent := `{
			"pages": {
				"home": {
					"en": {
						"title": "Home",
						"description": "Welcome to Blue"
					}
				},
				"about": {
					"en": {
						"title": "About",
						"description": "About Blue"
					}
				}
			}
		}`
		os.WriteFile("data/metadata.json", []byte(metadataContent), 0644)

		metadata, err := loadMetadata()
		if err != nil {
			t.Fatalf("Failed to load metadata: %v", err)
		}

		if len(metadata.Pages) != 2 {
			t.Errorf("Expected 2 pages, got %d", len(metadata.Pages))
		}

		if metadata.Pages["home"]["en"].Title != "Home" {
			t.Errorf("Expected home title 'Home', got %q", metadata.Pages["home"]["en"].Title)
		}
	})

	t.Run("Missing metadata file", func(t *testing.T) {
		os.Remove("data/metadata.json")

		_, err := loadMetadata()
		if err == nil {
			t.Error("Expected error for missing metadata file")
		}
	})

	t.Run("Invalid JSON in metadata file", func(t *testing.T) {
		os.WriteFile("data/metadata.json", []byte(`{invalid json}`), 0644)

		_, err := loadMetadata()
		if err == nil {
			t.Error("Expected error for invalid JSON")
		}
	})
}

// TestWriteSearchIndex tests the writeSearchIndex function
func TestWriteSearchIndex(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create public directory
	os.MkdirAll("public", 0755)

	items := []SearchItem{
		{
			Title:       "Test Page 1",
			Description: "Description 1",
			Content:     "Content 1",
			URL:         "/test1",
			Type:        "page",
		},
		{
			Title:       "Test Page 2",
			Description: "Description 2",
			Content:     "Content 2",
			URL:         "/test2",
			Type:        "content",
			Section:     "docs",
		},
	}

	err := writeSearchIndex(items)
	if err != nil {
		t.Fatalf("Failed to write search index: %v", err)
	}

	// Read and verify the written file
	data, err := os.ReadFile("public/searchIndex.json")
	if err != nil {
		t.Fatalf("Failed to read search index: %v", err)
	}

	var readItems []SearchItem
	err = json.Unmarshal(data, &readItems)
	if err != nil {
		t.Fatalf("Failed to parse search index: %v", err)
	}

	if len(readItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(readItems))
	}

	if readItems[0].Title != "Test Page 1" {
		t.Errorf("Expected first item title 'Test Page 1', got %q", readItems[0].Title)
	}
}

// TestIndexHTMLPages tests the indexHTMLPages function
func TestIndexHTMLPages(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create test directory structure
	os.MkdirAll("pages/platform", 0755)

	// Create test HTML files
	files := map[string]string{
		"pages/index.html": `<html>
			<head><title>Home</title></head>
			<body><h1>Welcome to Blue</h1><p>Task management platform.</p></body>
		</html>`,
		"pages/about.html": `<html>
			<body><h1>About Us</h1><p>Learn about our company.</p></body>
		</html>`,
		"pages/platform/features.html": `<html>
			<body><h1>Platform Features</h1><p>Feature overview.</p></body>
		</html>`,
		"pages/not-html.txt": "This should be ignored",
	}

	for path, content := range files {
		os.WriteFile(path, []byte(content), 0644)
	}

	// Test without metadata
	var items []SearchItem
	err := indexHTMLPages(&items, nil, nil)
	if err != nil {
		t.Fatalf("Failed to index HTML pages: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Verify URLs
	expectedURLs := map[string]bool{
		"/":                  false,
		"/about":             false,
		"/platform/features": false,
	}

	for _, item := range items {
		if item.Type != "page" {
			t.Errorf("Expected type 'page', got %q", item.Type)
		}
		expectedURLs[item.URL] = true
	}

	for url, found := range expectedURLs {
		if !found {
			t.Errorf("Expected URL %q not found", url)
		}
	}

	// Test with indexedURLs filter
	indexedURLs := map[string]bool{
		"/about": true,
	}

	var filteredItems []SearchItem
	err = indexHTMLPages(&filteredItems, nil, indexedURLs)
	if err != nil {
		t.Fatalf("Failed to index HTML pages with filter: %v", err)
	}

	// Should skip /about since it's already indexed
	if len(filteredItems) != 2 {
		t.Errorf("Expected 2 items after filtering, got %d", len(filteredItems))
	}
}

// TestIndexMarkdownContent tests the indexMarkdownContent function
func TestIndexMarkdownContent(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create test directory structure
	os.MkdirAll("content/docs/guides", 0755)
	os.MkdirAll("content/api", 0755)

	// Create test markdown files
	files := map[string]string{
		"content/docs/1. introduction.md": `---
title: Getting Started
description: Learn the basics
category: Documentation
---
# Getting Started

This is the introduction.`,
		"content/docs/guides/2. advanced-usage.md": `---
title: Advanced Usage Guide
---
# Advanced Usage

Advanced features explained here.`,
		"content/api/users.md": `# User API

API documentation for users.`,
		"content/not-markdown.txt": "This should be ignored",
	}

	for path, content := range files {
		os.WriteFile(path, []byte(content), 0644)
	}

	var items []SearchItem
	err := indexMarkdownContent(&items)
	if err != nil {
		t.Fatalf("Failed to index markdown content: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Check specific items
	for _, item := range items {
		if item.Type != "content" {
			t.Errorf("Expected type 'content', got %q", item.Type)
		}

		switch item.URL {
		case "/docs/introduction":
			if item.Title != "Getting Started" {
				t.Errorf("Expected title 'Getting Started', got %q", item.Title)
			}
			if item.Description != "Learn the basics" {
				t.Errorf("Expected description 'Learn the basics', got %q", item.Description)
			}
			if item.Section != "docs" {
				t.Errorf("Expected section 'docs', got %q", item.Section)
			}
		case "/docs/guides/advanced-usage":
			if item.Title != "Advanced Usage Guide" {
				t.Errorf("Expected title 'Advanced Usage Guide', got %q", item.Title)
			}
		case "/api/users":
			if item.Title != "User API" {
				t.Errorf("Expected title 'User API', got %q", item.Title)
			}
			if item.Section != "api" {
				t.Errorf("Expected section 'api', got %q", item.Section)
			}
		default:
			t.Errorf("Unexpected URL: %q", item.URL)
		}
	}
}

// TestIndexCachedMarkdownContentForLanguage tests the indexCachedMarkdownContentForLanguage function
func TestIndexCachedMarkdownContentForLanguage(t *testing.T) {
	// Create a mock MarkdownService with cached content
	markdownService := &MarkdownService{
		cache: &MarkdownCache{
			cache: make(map[string]*CachedContent),
		},
	}

	// Add test content to cache with language prefix
	testContent := map[string]*CachedContent{
		"en:/docs/intro": {
			HTML: `<h1>Introduction</h1><p>Welcome to the docs.</p>`,
			Frontmatter: &Frontmatter{
				Title:       "Introduction to Blue",
				Description: "Get started with Blue",
			},
		},
		"en:/api/projects": {
			HTML: `<h1>Projects API</h1><p>Manage projects via API.</p>`,
			Frontmatter: &Frontmatter{
				Title: "Projects API Reference",
			},
		},
		"en:/insights/best-practices": {
			HTML: `<h1>Best Practices</h1><p>Learn best practices.</p>`,
			Frontmatter: &Frontmatter{
				Title:       "Project Management Best Practices",
				Description: "Expert tips and tricks",
			},
		},
	}

	for path, content := range testContent {
		markdownService.cache.cache[path] = content
	}

	var items []SearchItem
	err := indexCachedMarkdownContentForLanguage(&items, markdownService, "en")
	if err != nil {
		t.Fatalf("Failed to index cached markdown content: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Verify items
	for _, item := range items {
		switch item.URL {
		case "/docs/intro":
			if item.Title != "Introduction to Blue" {
				t.Errorf("Expected title 'Introduction to Blue', got %q", item.Title)
			}
			if item.Category != "Docs" {
				t.Errorf("Expected category 'Docs', got %q", item.Category)
			}
			if item.Section != "docs" {
				t.Errorf("Expected section 'docs', got %q", item.Section)
			}
		case "/api/projects":
			if item.Title != "Projects API Reference" {
				t.Errorf("Expected title 'Projects API Reference', got %q", item.Title)
			}
			if item.Category != "API" {
				t.Errorf("Expected category 'API', got %q", item.Category)
			}
		case "/insights/best-practices":
			if item.Title != "Project Management Best Practices" {
				t.Errorf("Expected title 'Project Management Best Practices', got %q", item.Title)
			}
			if item.Category != "Insights" {
				t.Errorf("Expected category 'Insights', got %q", item.Category)
			}
		}

		// Content should be extracted text without HTML tags
		if strings.Contains(strings.ToLower(item.Content), "<h1") || strings.Contains(strings.ToLower(item.Content), "</h1>") {
			t.Errorf("Content still contains HTML tags: %q", item.Content)
		}
		// Content should have some text
		if item.Content == "" {
			t.Errorf("Content is empty for %q", item.URL)
		}
	}
}

// TestIndexCachedHTMLPagesForLanguage tests the indexCachedHTMLPagesForLanguage function
func TestIndexCachedHTMLPagesForLanguage(t *testing.T) {
	// Create mock services
	htmlService := &HTMLService{
		cache: &HTMLCache{
			cache: make(map[string]*CachedContent),
		},
	}

	metadata := &Metadata{
		Pages: map[string]map[string]PageMetadata{
			"pricing": {
				"en": {
					Title:       "Pricing Plans",
					Description: "Choose the right plan for your team",
				},
			},
		},
	}

	// Add test content to cache with language prefix
	testContent := map[string]*CachedContent{
		"en:/": {
			HTML:     `<h1>Welcome to Blue</h1><p>Task management made simple.</p>`,
			FilePath: "pages/index.html",
		},
		"en:/pricing": {
			HTML:     `<h1>Pricing</h1><p>Our pricing plans.</p>`,
			FilePath: "pages/pricing.html",
		},
		"en:/platform/features/views": {
			HTML:     `<h1>Multiple Views</h1><p>View your tasks your way.</p>`,
			FilePath: "pages/platform/features/views.html",
		},
	}

	for path, content := range testContent {
		htmlService.cache.cache[path] = content
	}

	var items []SearchItem
	err := indexCachedHTMLPagesForLanguage(&items, htmlService, metadata, "en")
	if err != nil {
		t.Fatalf("Failed to index cached HTML pages: %v", err)
	}

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	// Verify items
	for _, item := range items {
		if item.Type != "page" {
			t.Errorf("Expected type 'page', got %q", item.Type)
		}

		switch item.URL {
		case "/":
			if item.Title != "Welcome to Blue" {
				t.Errorf("Expected title 'Welcome to Blue', got %q", item.Title)
			}
		case "/pricing":
			// Should use metadata title
			if item.Title != "Pricing Plans" {
				t.Errorf("Expected title 'Pricing Plans', got %q", item.Title)
			}
			if item.Description != "Choose the right plan for your team" {
				t.Errorf("Expected description from metadata, got %q", item.Description)
			}
		case "/platform/features/views":
			if item.Title != "Multiple Views" {
				t.Errorf("Expected title 'Multiple Views', got %q", item.Title)
			}
			if item.Category != "Feature" {
				t.Errorf("Expected category 'Feature', got %q", item.Category)
			}
			if item.Section != "platform" {
				t.Errorf("Expected section 'platform', got %q", item.Section)
			}
		}
	}
}

// TestGenerateSearchIndex tests the main GenerateSearchIndex function
func TestGenerateSearchIndex(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create directory structure
	os.MkdirAll("pages", 0755)
	os.MkdirAll("content/docs", 0755)
	os.MkdirAll("public", 0755)
	os.MkdirAll("data", 0755)

	// Create test files
	os.WriteFile("pages/test.html", []byte(`<h1>Test Page</h1><p>Content</p>`), 0644)
	os.WriteFile("content/docs/guide.md", []byte(`---
title: User Guide
---
# User Guide

Guide content here.`), 0644)

	// Create metadata.json
	metadataContent := `{
		"pages": {
			"test": {
				"title": "Test Page Title",
				"description": "Test description"
			}
		}
	}`
	os.WriteFile("data/metadata.json", []byte(metadataContent), 0644)

	err := GenerateSearchIndex()
	if err != nil {
		t.Fatalf("Failed to generate search index: %v", err)
	}

	// Verify searchIndex.json was created
	if _, err := os.Stat("public/searchIndex.json"); os.IsNotExist(err) {
		t.Error("searchIndex.json was not created")
	}

	// Read and verify content
	data, _ := os.ReadFile("public/searchIndex.json")
	var items []SearchItem
	json.Unmarshal(data, &items)

	if len(items) != 2 {
		t.Errorf("Expected 2 items in search index, got %d", len(items))
	}
}

// TestGenerateSearchIndexWithCache tests the GenerateSearchIndexWithCache function
func TestGenerateSearchIndexWithCache(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create directory structure
	os.MkdirAll("pages", 0755)
	os.MkdirAll("public", 0755)
	os.MkdirAll("data", 0755)

	// Create test HTML file
	os.WriteFile("pages/test.html", []byte(`<h1>Test Page</h1>`), 0644)

	// Create mock MarkdownService with cached content
	markdownService := &MarkdownService{
		cache: &MarkdownCache{
			cache: make(map[string]*CachedContent),
		},
	}

	// Add cached markdown content
	markdownService.cache.cache["/docs/cached"] = &CachedContent{
		HTML: `<h1>Cached Doc</h1><p>Cached content.</p>`,
		Frontmatter: &Frontmatter{
			Title: "Cached Documentation",
		},
	}

	err := GenerateSearchIndexWithCache(markdownService)
	if err != nil {
		t.Fatalf("Failed to generate search index with cache: %v", err)
	}

	// Verify searchIndex.json was created
	if _, err := os.Stat("public/searchIndex.json"); os.IsNotExist(err) {
		t.Error("searchIndex.json was not created")
	}

	// Read and verify content
	data, _ := os.ReadFile("public/searchIndex.json")
	var items []SearchItem
	json.Unmarshal(data, &items)

	if len(items) < 2 {
		t.Errorf("Expected at least 2 items in search index, got %d", len(items))
	}

	// Check for cached content
	foundCached := false
	for _, item := range items {
		if item.URL == "/docs/cached" {
			foundCached = true
			if item.Title != "Cached Documentation" {
				t.Errorf("Expected cached title 'Cached Documentation', got %q", item.Title)
			}
		}
	}

	if !foundCached {
		t.Error("Cached markdown content not found in search index")
	}
}

// TestGenerateSearchIndexWithCaches tests the GenerateSearchIndexWithCaches function
func TestGenerateSearchIndexWithCaches(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create directory structure
	os.MkdirAll("pages", 0755)
	os.MkdirAll("public", 0755)
	os.MkdirAll("data", 0755)

	// Create a non-cached HTML file (like status page)
	os.WriteFile("pages/status.html", []byte(`<h1>System Status</h1><p>All systems operational.</p>`), 0644)

	// Create mock services
	markdownService := &MarkdownService{
		cache: &MarkdownCache{
			cache: make(map[string]*CachedContent),
		},
	}

	htmlService := &HTMLService{
		cache: &HTMLCache{
			cache: make(map[string]*CachedContent),
		},
	}

	// Add cached content
	htmlService.cache.cache["/pricing"] = &CachedContent{
		HTML:     `<h1>Pricing</h1><p>Our plans.</p>`,
		FilePath: "pages/pricing.html",
	}

	markdownService.cache.cache["/docs/intro"] = &CachedContent{
		HTML: `<h1>Introduction</h1><p>Get started.</p>`,
		Frontmatter: &Frontmatter{
			Title: "Introduction",
		},
	}

	err := GenerateSearchIndexWithCaches(markdownService, htmlService)
	if err != nil {
		t.Fatalf("Failed to generate search index with caches: %v", err)
	}

	// Verify searchIndex.json was created
	if _, err := os.Stat("public/searchIndex.json"); os.IsNotExist(err) {
		t.Error("searchIndex.json was not created")
	}

	// Read and verify content
	data, _ := os.ReadFile("public/searchIndex.json")
	var items []SearchItem
	json.Unmarshal(data, &items)

	if len(items) < 3 {
		t.Errorf("Expected at least 3 items in search index, got %d", len(items))
	}

	// Verify we have all types of content
	hasStatus := false
	hasCachedHTML := false
	hasCachedMarkdown := false

	for _, item := range items {
		switch item.URL {
		case "/status":
			hasStatus = true
		case "/pricing":
			hasCachedHTML = true
		case "/docs/intro":
			hasCachedMarkdown = true
		}
	}

	if !hasStatus {
		t.Error("Non-cached status page not found in search index")
	}
	if !hasCachedHTML {
		t.Error("Cached HTML page not found in search index")
	}
	if !hasCachedMarkdown {
		t.Error("Cached markdown content not found in search index")
	}
}

// TestMarkdownURLCleaning tests the URL cleaning logic for markdown files
func TestMarkdownURLCleaning(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create test directory structure
	os.MkdirAll("content/docs/guides", 0755)

	// Create test files with various naming patterns
	files := map[string]string{
		"content/docs/1. introduction.md":              "# Introduction",
		"content/docs/2.getting-started.md":            "# Getting Started",
		"content/docs/guides/10. advanced features.md": "# Advanced",
		"content/api/11.users.md":                      "# Users API",
	}

	for path, content := range files {
		os.WriteFile(path, []byte(content), 0644)
	}

	var items []SearchItem
	err := indexMarkdownContent(&items)
	if err != nil {
		t.Fatalf("Failed to index markdown content: %v", err)
	}

	// Check URL cleaning
	expectedURLs := map[string]string{
		"/docs/introduction":             "Introduction",
		"/docs/getting-started":          "getting started",
		"/docs/guides/advanced-features": "advanced features",
		"/api/users":                     "users",
	}

	for _, item := range items {
		expectedTitle, exists := expectedURLs[item.URL]
		if !exists {
			t.Errorf("Unexpected URL: %q", item.URL)
			continue
		}

		if item.Title != expectedTitle {
			t.Errorf("For URL %q, expected title %q, got %q", item.URL, expectedTitle, item.Title)
		}
	}
}

// TestConcurrentSearchIndexGeneration tests thread safety of search index generation
func TestConcurrentSearchIndexGeneration(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Change to test directory
	originalWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(originalWd)

	// Create directory structure
	os.MkdirAll("public", 0755)

	// Run multiple concurrent writes
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func(n int) {
			items := []SearchItem{
				{
					Title:   "Concurrent Test",
					Content: "Test content",
					URL:     "/test",
					Type:    "page",
				},
			}
			writeSearchIndex(items)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify the file exists and is valid JSON
	data, err := os.ReadFile("public/searchIndex.json")
	if err != nil {
		t.Fatalf("Failed to read search index after concurrent writes: %v", err)
	}

	var items []SearchItem
	err = json.Unmarshal(data, &items)
	if err != nil {
		t.Fatalf("Search index contains invalid JSON after concurrent writes: %v", err)
	}
}

// TestEdgeCasesAndErrorHandling tests various edge cases
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	t.Run("Empty frontmatter values", func(t *testing.T) {
		fm, _ := splitFrontmatter(`---
title: 
description: 
---
Content`)

		if !strings.Contains(fm, "title:") {
			t.Error("Expected frontmatter to contain empty title field")
		}
	})

	t.Run("Very long content", func(t *testing.T) {
		longContent := strings.Repeat("This is a very long content. ", 1000)
		html := "<h1>Title</h1><p>" + longContent + "</p>"

		result := extractTextFromHTML(html)
		if !strings.Contains(result, "This is a very long content") {
			t.Error("Failed to extract long content")
		}
	})

	t.Run("Unicode in titles and content", func(t *testing.T) {
		unicodeTitle := extractH1Title(`<h1>Unicode Test 你好 🌟</h1>`)
		if unicodeTitle != "Unicode Test 你好 🌟" {
			t.Errorf("Failed to extract Unicode title, got %q", unicodeTitle)
		}
	})

	t.Run("Malformed HTML", func(t *testing.T) {
		malformed := `<h1>Title<h1><p>Content</h1>`
		text := extractTextFromHTML(malformed)
		// Should still extract something despite malformed HTML
		if text == "" {
			t.Error("Should extract some text even from malformed HTML")
		}
	})

	t.Run("Nested script tags", func(t *testing.T) {
		html := `<h1>Title</h1>
			<script>
				var code = '<script>nested</script>';
			</script>
			<p>Content</p>`

		result := extractTextFromHTML(html)
		if strings.Contains(result, "nested") || strings.Contains(result, "var code") {
			t.Error("Script content should be completely removed")
		}
	})
}
