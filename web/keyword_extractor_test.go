package web

import (
	"strings"
	"testing"
	"time"
)

func TestKeywordExtractor(t *testing.T) {
	extractor := NewKeywordExtractor()

	// Test basic keyword extraction
	text := `Blue is a powerful task management platform that helps teams collaborate efficiently. 
	With features like real-time collaboration, project tracking, and automated workflows, 
	Blue makes it easy to manage complex projects. The platform includes advanced security features
	and integrates with popular tools like Slack, GitHub, and Jira. Teams can create custom workflows,
	track progress, and generate detailed reports. Blue's AI-powered insights help identify bottlenecks
	and optimize team performance.`

	keywords := extractor.ExtractKeywords(text, 20)

	// Check that we got keywords
	if len(keywords) == 0 {
		t.Error("No keywords extracted")
	}

	// Check that common words are filtered out
	for _, keyword := range keywords {
		if keyword == "the" || keyword == "and" || keyword == "is" {
			t.Errorf("Stop word not filtered: %s", keyword)
		}
	}

	// Print keywords for manual verification
	t.Logf("Extracted keywords: %v", keywords)
}

func TestPhrasesExtraction(t *testing.T) {
	extractor := NewKeywordExtractor()

	text := `Machine learning algorithms are important. Machine learning models help predict outcomes.
	Deep learning networks process data efficiently. Deep learning systems improve accuracy.
	Natural language processing is advancing rapidly. Natural language understanding helps computers.
	Computer vision enables automation. Computer vision systems detect objects.`

	keywords := extractor.ExtractKeywords(text, 15)

	// Should include multi-word phrases
	found := false
	for _, keyword := range keywords {
		if strings.Contains(keyword, " ") {
			found = true
			t.Logf("Found phrase: %s", keyword)
		}
	}

	if !found {
		t.Error("No phrases extracted")
	}
}

func BenchmarkKeywordExtraction(b *testing.B) {
	extractor := NewKeywordExtractor()

	// Sample document (about 500 words)
	text := strings.Repeat(`Blue is a comprehensive project management platform designed for modern teams. 
	It offers powerful features for task tracking, collaboration, and workflow automation. 
	Teams can create projects, assign tasks, set deadlines, and monitor progress in real-time. 
	The platform includes built-in communication tools, file sharing capabilities, and advanced reporting. 
	With integrations for popular tools and a flexible API, Blue adapts to any workflow. 
	Security is a top priority with enterprise-grade encryption and compliance certifications. `, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractor.ExtractKeywords(text, 50)
	}
}

func TestPerformanceComparison(t *testing.T) {
	extractor := NewKeywordExtractor()

	// Create a sample document
	text := strings.Repeat(`Blue platform features include task management, project tracking, team collaboration,
	workflow automation, real-time updates, custom fields, advanced reporting, security controls,
	API integrations, mobile apps, desktop applications, web interface, user permissions, data export,
	backup systems, audit logs, compliance tools, and enterprise support. `, 50)

	// Measure extraction time
	start := time.Now()
	keywords := extractor.ExtractKeywords(text, 80)
	extractTime := time.Since(start)

	// Calculate size reduction
	originalSize := len(text)
	keywordsSize := 0
	for _, k := range keywords {
		keywordsSize += len(k) + 1 // +1 for separator
	}

	reduction := float64(originalSize-keywordsSize) / float64(originalSize) * 100

	t.Logf("Original text size: %d bytes", originalSize)
	t.Logf("Keywords size: %d bytes", keywordsSize)
	t.Logf("Size reduction: %.1f%%", reduction)
	t.Logf("Extraction time: %v", extractTime)
	t.Logf("Keywords extracted: %d", len(keywords))
}