package web

import (
	"html/template"
	"testing"
	"time"
)

// TestIsTOCExcluded tests the TOC exclusion logic
func TestIsTOCExcluded(t *testing.T) {
	router := &Router{
		tocExcludedPaths: []string{
			"/changelog",
			"/roadmap",
			"/platform/status",
		},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Excluded changelog",
			path:     "/changelog",
			expected: true,
		},
		{
			name:     "Excluded roadmap",
			path:     "/roadmap",
			expected: true,
		},
		{
			name:     "Excluded status",
			path:     "/platform/status",
			expected: true,
		},
		{
			name:     "Not excluded docs",
			path:     "/docs/introduction",
			expected: false,
		},
		{
			name:     "Not excluded insights",
			path:     "/insights",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := router.isTOCExcluded(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestPreparePageData tests the page data preparation
func TestPreparePageData(t *testing.T) {
	// Create test router with services
	router := &Router{
		tocExcludedPaths:  []string{"/changelog"},
		seoService:        NewSEOService(),
		markdownService:   NewMarkdownService(),
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	// Test basic page data preparation
	content := template.HTML("<h2 id='section-1'>Section 1</h2><p>Content</p>")
	pageData := router.preparePageData("/test", content, false, nil, &Navigation{}, "en")

	// Verify basic fields
	if pageData.Path != "/test" {
		t.Errorf("Expected path /test, got %s", pageData.Path)
	}

	if pageData.Content != content {
		t.Error("Content mismatch")
	}

	if pageData.IsMarkdown {
		t.Error("Expected IsMarkdown to be false")
	}

	if pageData.CustomerNumber != 17000 {
		t.Errorf("Expected CustomerNumber 17000, got %d", pageData.CustomerNumber)
	}
}

// TestInsightsSorting tests that insights are sorted by date in reverse chronological order
func TestInsightsSorting(t *testing.T) {
	// Create test router with mocked markdown service
	markdownService := NewMarkdownService()

	// Pre-populate cache with test insights
	insights := map[string]*CachedContent{
		"/insights/old-article": {
			Frontmatter: &Frontmatter{
				Title:       "Old Article",
				Description: "An old article",
				Category:    "Engineering",
				Date:        "2024-01-01",
			},
		},
		"/insights/new-article": {
			Frontmatter: &Frontmatter{
				Title:       "New Article",
				Description: "A new article",
				Category:    "Product Updates",
				Date:        "2024-12-25",
			},
		},
		"/insights/medium-article": {
			Frontmatter: &Frontmatter{
				Title:       "Medium Article",
				Description: "A medium article",
				Category:    "News",
				Date:        "2024-06-15",
			},
		},
		"/insights/no-date-article": {
			Frontmatter: &Frontmatter{
				Title:       "No Date Article",
				Description: "Article without date",
				Category:    "FAQ",
				Date:        "", // Empty date
			},
		},
		"/insights/future-article": {
			Frontmatter: &Frontmatter{
				Title:       "Future Article",
				Description: "Article from the future",
				Category:    "Engineering",
				Date:        "2025-01-15",
			},
		},
		"/not-insight/other": {
			Frontmatter: &Frontmatter{
				Title: "Not an insight",
				Date:  "2024-01-01",
			},
		},
	}

	// Add insights to cache
	for path, content := range insights {
		markdownService.cache.Set(path, content)
	}

	router := &Router{
		tocExcludedPaths:  []string{},
		seoService:        NewSEOService(),
		markdownService:   markdownService,
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	// Prepare page data for insights page
	pageData := router.preparePageData("/insights", template.HTML(""), false, nil, &Navigation{}, "en")

	// Verify insights are present and sorted
	if len(pageData.Insights) != 5 {
		t.Errorf("Expected 5 insights, got %d", len(pageData.Insights))
	}

	// Check order - should be: Future (2025-01-15), New (2024-12-25), Medium (2024-06-15), Old (2024-01-01), No Date ("")
	expectedOrder := []string{
		"Future Article",
		"New Article",
		"Medium Article",
		"Old Article",
		"No Date Article",
	}

	for i, expectedTitle := range expectedOrder {
		if i < len(pageData.Insights) && pageData.Insights[i].Title != expectedTitle {
			t.Errorf("Expected insight %d to be %q, got %q", i, expectedTitle, pageData.Insights[i].Title)
		}
	}

	// Verify dates are in descending order (except empty dates at end)
	for i := 0; i < len(pageData.Insights)-1; i++ {
		current := pageData.Insights[i].Date
		next := pageData.Insights[i+1].Date

		// Skip comparison if we hit empty dates
		if current == "" || next == "" {
			continue
		}

		if current < next {
			t.Errorf("Insights not in reverse chronological order: %s < %s", current, next)
		}
	}
}

// TestTOCGeneration tests table of contents generation
func TestTOCGeneration(t *testing.T) {
	router := &Router{
		tocExcludedPaths:  []string{"/excluded"},
		seoService:        NewSEOService(),
		markdownService:   NewMarkdownService(),
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	tests := []struct {
		name       string
		path       string
		content    string
		isMarkdown bool
		expectTOC  bool
		tocLength  int
	}{
		{
			name:       "Markdown with H2s",
			path:       "/docs/test",
			content:    `<h2 id="section-1">Section 1</h2><p>Content</p><h2 id="section-2">Section 2</h2>`,
			isMarkdown: true,
			expectTOC:  true,
			tocLength:  2,
		},
		{
			name:       "HTML with sections",
			path:       "/features",
			content:    `<section id="feature-1"><h2>Feature 1</h2></section><section id="feature-2"><h2>Feature 2</h2></section>`,
			isMarkdown: false,
			expectTOC:  true,
			tocLength:  2,
		},
		{
			name:       "Excluded path",
			path:       "/excluded",
			content:    `<h2 id="section">Section</h2>`,
			isMarkdown: true,
			expectTOC:  false,
			tocLength:  0,
		},
		{
			name:       "Empty content",
			path:       "/empty",
			content:    "",
			isMarkdown: true,
			expectTOC:  false,
			tocLength:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageData := router.preparePageData(tt.path, template.HTML(tt.content), tt.isMarkdown, nil, &Navigation{}, "en")

			if tt.expectTOC && len(pageData.TOC) != tt.tocLength {
				t.Errorf("Expected %d TOC entries, got %d", tt.tocLength, len(pageData.TOC))
			}

			if !tt.expectTOC && len(pageData.TOC) > 0 {
				t.Errorf("Expected no TOC entries, got %d", len(pageData.TOC))
			}
		})
	}
}

// TestPNGGeneration tests PNG generation for insights
func TestPNGGeneration(t *testing.T) {
	// This test verifies that PNG paths are generated for insights
	markdownService := NewMarkdownService()

	// Add a test insight
	markdownService.cache.Set("/insights/test-article", &CachedContent{
		Frontmatter: &Frontmatter{
			Title:       "Test Article",
			Description: "Test description",
			Category:    "Testing",
			Slug:        "test-article",
			Date:        "2024-01-01",
		},
	})

	router := &Router{
		tocExcludedPaths:  []string{},
		seoService:        NewSEOService(),
		markdownService:   markdownService,
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	// Prepare page data for insights page
	pageData := router.preparePageData("/insights", template.HTML(""), false, nil, &Navigation{}, "en")

	// Verify PNG path is generated
	if len(pageData.Insights) != 1 {
		t.Fatalf("Expected 1 insight, got %d", len(pageData.Insights))
	}

	insight := pageData.Insights[0]
	if insight.PNGPath == "" {
		t.Error("Expected PNG path to be generated")
	}

	// PNG path should be based on the title converted to filename format
	// "Test Article" -> "test-article"
	expectedPath := "/insights/test-article.png"
	if insight.PNGPath != expectedPath {
		t.Errorf("Expected PNG path %s, got %s", expectedPath, insight.PNGPath)
	}
}

// TestInsightCategoryExtraction tests category extraction from tags
func TestInsightCategoryExtraction(t *testing.T) {
	markdownService := NewMarkdownService()

	// Test cases with different category/tag combinations
	testCases := map[string]*CachedContent{
		"/insights/with-category": {
			Frontmatter: &Frontmatter{
				Title:    "With Category",
				Category: "Engineering",
				Tags:     []string{"backend", "golang"},
			},
		},
		"/insights/no-category-with-tags": {
			Frontmatter: &Frontmatter{
				Title: "No Category With Tags",
				Tags:  []string{"Product Updates", "feature"},
			},
		},
		"/insights/no-category-no-tags": {
			Frontmatter: &Frontmatter{
				Title: "No Category No Tags",
			},
		},
	}

	for path, content := range testCases {
		markdownService.cache.Set(path, content)
	}

	router := &Router{
		tocExcludedPaths:  []string{},
		seoService:        NewSEOService(),
		markdownService:   markdownService,
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	pageData := router.preparePageData("/insights", template.HTML(""), false, nil, &Navigation{}, "en")

	// Create a map for easier lookup
	insightsByTitle := make(map[string]InsightData)
	for _, insight := range pageData.Insights {
		insightsByTitle[insight.Title] = insight
	}

	// Test category is used when available
	if insight, ok := insightsByTitle["With Category"]; ok {
		if insight.Category != "Engineering" {
			t.Errorf("Expected category 'Engineering', got %q", insight.Category)
		}
	} else {
		t.Error("Insight 'With Category' not found")
	}

	// Test first tag is used when no category
	if insight, ok := insightsByTitle["No Category With Tags"]; ok {
		if insight.Category != "Product Updates" {
			t.Errorf("Expected category 'Product Updates' from first tag, got %q", insight.Category)
		}
	} else {
		t.Error("Insight 'No Category With Tags' not found")
	}

	// Test empty category when no category or tags
	if insight, ok := insightsByTitle["No Category No Tags"]; ok {
		if insight.Category != "" {
			t.Errorf("Expected empty category, got %q", insight.Category)
		}
	} else {
		t.Error("Insight 'No Category No Tags' not found")
	}
}

// TestSchemaDataGeneration tests that schema data is generated
func TestSchemaDataGeneration(t *testing.T) {
	router := &Router{
		tocExcludedPaths:  []string{},
		seoService:        NewSEOService(),
		markdownService:   NewMarkdownService(),
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	// Test with frontmatter
	frontmatter := &Frontmatter{
		Title:       "Test Article",
		Description: "Test description",
		Date:        "2024-01-01",
	}

	pageData := router.preparePageData("/insights/test", template.HTML(""), true, frontmatter, &Navigation{}, "en")

	// Schema data should be generated
	if pageData.SchemaData == "" {
		t.Error("Expected schema data to be generated")
	}

	// Should not be the default empty array
	if string(pageData.SchemaData) == "[]" {
		t.Error("Expected non-empty schema data")
	}
}

// BenchmarkInsightsSorting benchmarks the insights sorting performance
func BenchmarkInsightsSorting(b *testing.B) {
	// Create a large number of insights
	markdownService := NewMarkdownService()

	// Generate 1000 insights with random dates
	for i := 0; i < 1000; i++ {
		date := ""
		if i%10 != 0 { // 90% have dates
			year := 2020 + (i % 5)
			month := 1 + (i % 12)
			day := 1 + (i % 28)
			date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		}

		markdownService.cache.Set(
			"/insights/article-"+string(rune(i)),
			&CachedContent{
				Frontmatter: &Frontmatter{
					Title:    "Article " + string(rune(i)),
					Category: "Test",
					Date:     date,
				},
			},
		)
	}

	router := &Router{
		tocExcludedPaths:  []string{},
		seoService:        NewSEOService(),
		markdownService:   markdownService,
		navigationService: NewNavigationService(NewSEOService()),
		schemaService:     NewSchemaService(nil, "https://example.com"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.preparePageData("/insights", template.HTML(""), false, nil, &Navigation{}, "en")
	}
}
