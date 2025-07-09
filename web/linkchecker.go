package web

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// LinkChecker validates internal links across the site
type LinkChecker struct {
	markdownService *MarkdownService
	htmlService     *HTMLService
	seoService      *SEOService
	validRoutes     map[string]bool
	publicFiles     map[string]bool
	brokenLinks     map[string][]BrokenLink
	totalLinks      int
	validLinks      int
}

// BrokenLink represents a broken link found in the site
type BrokenLink struct {
	URL         string
	SourcePath  string
	LineContext string
}

// NewLinkChecker creates a new link checker instance
func NewLinkChecker(markdownService *MarkdownService, htmlService *HTMLService, seoService *SEOService) *LinkChecker {
	return &LinkChecker{
		markdownService: markdownService,
		htmlService:     htmlService,
		seoService:      seoService,
		validRoutes:     make(map[string]bool),
		publicFiles:     make(map[string]bool),
		brokenLinks:     make(map[string][]BrokenLink),
	}
}

// CheckAllLinks checks all internal links in the site
func (lc *LinkChecker) CheckAllLinks() error {
	// Build valid routes from cached content
	lc.buildValidRoutes()

	// Check links in all cached markdown content
	for path, content := range lc.markdownService.GetAllCachedContent() {
		links := lc.extractLinksFromHTML(content.HTML)
		links = append(links, lc.extractLinksFromMarkdown(content.HTML)...)
		lc.validateLinks(path, links)
	}

	// Check links in all cached HTML content
	for path, content := range lc.htmlService.GetAllCachedContent() {
		links := lc.extractLinksFromHTML(content.HTML)
		lc.validateLinks(path, links)
	}

	// Print results
	lc.printResults()

	return nil
}

// buildValidRoutes builds a map of all valid routes
func (lc *LinkChecker) buildValidRoutes() {
	// Add all cached markdown routes
	for path := range lc.markdownService.GetAllCachedContent() {
		lc.validRoutes[path] = true
		// Also add with trailing slash
		if !strings.HasSuffix(path, "/") {
			lc.validRoutes[path+"/"] = true
		}
	}

	// Add all cached HTML routes
	for path := range lc.htmlService.GetAllCachedContent() {
		lc.validRoutes[path] = true
		// Also add with trailing slash
		if !strings.HasSuffix(path, "/") {
			lc.validRoutes[path+"/"] = true
		}
	}

	// Explicitly add common routes that might be handled specially by router
	lc.validRoutes["/"] = true

	// Add common API and special routes
	lc.validRoutes["/health"] = true
	lc.validRoutes["/api/status/current"] = true
	lc.validRoutes["/api/status/history"] = true
	lc.validRoutes["/platform/status"] = true

	// Scan public directory for static files
	lc.scanPublicFiles()
}

// scanPublicFiles scans the public directory and adds all files as valid routes
func (lc *LinkChecker) scanPublicFiles() {
	err := filepath.WalkDir("public", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		if !d.IsDir() {
			// Convert file path to URL path
			urlPath := "/" + path
			lc.publicFiles[urlPath] = true
		}

		return nil
	})

	if err != nil {
		// Fallback to common files if scanning fails
		lc.publicFiles["/favicon.ico"] = true
		lc.publicFiles["/css/style.css"] = true
	}
}

// extractLinksFromHTML extracts links from HTML content
func (lc *LinkChecker) extractLinksFromHTML(htmlContent string) []string {
	var links []string
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return links
	}

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			case "img":
				for _, attr := range n.Attr {
					if attr.Key == "src" && !strings.HasPrefix(attr.Val, "data:") {
						links = append(links, attr.Val)
					}
				}
			case "video":
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						links = append(links, attr.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)

	return links
}

// extractLinksFromMarkdown extracts markdown-style links
func (lc *LinkChecker) extractLinksFromMarkdown(content string) []string {
	var links []string

	// Match markdown links [text](url)
	mdLinkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	matches := mdLinkRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			links = append(links, match[2])
		}
	}

	// Match markdown images ![alt](url)
	mdImageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	imageMatches := mdImageRegex.FindAllStringSubmatch(content, -1)
	for _, match := range imageMatches {
		if len(match) > 2 {
			links = append(links, match[2])
		}
	}

	return links
}

// validateLinks validates a list of links for a given source path
func (lc *LinkChecker) validateLinks(sourcePath string, links []string) {
	for _, link := range links {
		lc.totalLinks++

		// Skip external links
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") || strings.HasPrefix(link, "//") {
			lc.validLinks++
			continue
		}

		// Skip mailto and other protocols
		if strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") || strings.HasPrefix(link, "javascript:") {
			lc.validLinks++
			continue
		}

		// Skip plain email addresses (without mailto: prefix)
		if strings.Contains(link, "@") && !strings.Contains(link, "/") {
			lc.validLinks++
			continue
		}

		// Skip empty links and anchors on same page
		if link == "" || link == "#" {
			lc.validLinks++
			continue
		}

		// Parse the link
		u, err := url.Parse(link)
		if err != nil {
			lc.addBrokenLink(sourcePath, link, "Invalid URL format")
			continue
		}

		// Handle anchor-only links
		if strings.HasPrefix(link, "#") {
			// TODO: Validate anchor exists in current page
			lc.validLinks++
			continue
		}

		// Get the path without fragment and query
		checkPath := u.Path
		if checkPath == "" {
			checkPath = "/"
		}

		// Check if it's a valid route
		if lc.validRoutes[checkPath] {
			lc.validLinks++
			continue
		}

		// Check if it's a public file
		if strings.HasPrefix(checkPath, "/") || lc.publicFiles[checkPath] {
			lc.validLinks++
			continue
		}

		// Check for redirects
		if redirectTo, _, shouldRedirect := lc.seoService.CheckRedirect(checkPath); shouldRedirect {
			// Link redirects to valid location
			if lc.validRoutes[redirectTo] {
				lc.validLinks++
				continue
			}
		}

		// Link is broken
		lc.addBrokenLink(sourcePath, link, "Page not found")
	}
}

// addBrokenLink adds a broken link to the results
func (lc *LinkChecker) addBrokenLink(sourcePath, link, context string) {
	brokenLink := BrokenLink{
		URL:         link,
		SourcePath:  sourcePath,
		LineContext: context,
	}

	if _, exists := lc.brokenLinks[link]; !exists {
		lc.brokenLinks[link] = []BrokenLink{}
	}
	lc.brokenLinks[link] = append(lc.brokenLinks[link], brokenLink)
}

// printResults prints the link checking results and generates CSV
func (lc *LinkChecker) printResults() {
	brokenCount := lc.totalLinks - lc.validLinks

	if brokenCount == 0 {
		log.Printf("✅ %d/%d links valid", lc.validLinks, lc.totalLinks)
	} else {
		log.Printf("🔗 %d/%d links valid", lc.validLinks, lc.totalLinks)
		log.Printf("❌ %d broken links found", brokenCount)

		// Generate CSV file
		lc.generateCSV()
	}
}

// generateCSV creates a CSV file with broken links and their sources
func (lc *LinkChecker) generateCSV() {
	file, err := os.Create("404.csv")
	if err != nil {
		log.Printf("Failed to create 404.csv: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Broken URL", "Source Page", "References"}
	if err := writer.Write(header); err != nil {
		log.Printf("Failed to write CSV header: %v", err)
		return
	}

	// Sort broken links by URL for consistent output
	var brokenURLs []string
	for url := range lc.brokenLinks {
		brokenURLs = append(brokenURLs, url)
	}
	sort.Strings(brokenURLs)

	// Write each broken link with its sources
	for _, url := range brokenURLs {
		sources := lc.brokenLinks[url]

		// Group sources and count references
		sourceCount := make(map[string]int)
		for _, source := range sources {
			sourceCount[source.SourcePath]++
		}

		// Write one row per unique source page
		for sourcePage, count := range sourceCount {
			record := []string{
				url,
				sourcePage,
				fmt.Sprintf("%d", count),
			}
			if err := writer.Write(record); err != nil {
				log.Printf("Failed to write CSV record: %v", err)
			}
		}
	}

	log.Printf("📝 Broken links exported to 404.csv")
}

// RunLinkChecker runs the link checker with the provided services
func RunLinkChecker(markdownService *MarkdownService, htmlService *HTMLService, seoService *SEOService) error {
	log.Printf("🔗 Checking internal links...")

	checker := NewLinkChecker(markdownService, htmlService, seoService)
	return checker.CheckAllLinks()
}
