package web

import (
	"regexp"
	"strings"
	"golang.org/x/net/html"
)

// TOCEntry represents a table of contents entry
type TOCEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Level int    `json:"level"`
}

// ExtractHTMLTOC extracts table of contents from HTML content
// Looks for H2 elements with ID attributes
func ExtractHTMLTOC(content string) ([]TOCEntry, error) {
	var toc []TOCEntry
	
	// Wrap content in a proper HTML structure for parsing
	htmlContent := "<html><body>" + content + "</body></html>"
	
	// Parse HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return toc, err
	}
	
	// Walk the HTML tree to find sections with IDs
	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "section" {
			// Extract ID attribute from section
			var id string
			for _, attr := range n.Attr {
				if attr.Key == "id" {
					id = attr.Val
					break
				}
			}
			
			// If section has ID, create TOC entry from ID
			if id != "" {
				// Generate title from ID (e.g., "key-features" → "Key Features")
				title := generateTitleFromID(id)
				toc = append(toc, TOCEntry{
					ID:    id,
					Title: title,
					Level: 2,
				})
			}
		}
		
		// Recursively walk children
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walker(child)
		}
	}
	
	walker(doc)
	return toc, nil
}

// ExtractH2TOC extracts table of contents from HTML content with H2 elements
// Looks for H2 elements with ID attributes (used for markdown-generated pages)
func ExtractH2TOC(content string) ([]TOCEntry, error) {
	var toc []TOCEntry
	
	// Wrap content in a proper HTML structure for parsing
	htmlContent := "<html><body>" + content + "</body></html>"
	
	// Parse HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return toc, err
	}
	
	// Walk the HTML tree to find H2 elements with IDs
	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h2" {
			// Extract ID attribute from h2
			var id string
			for _, attr := range n.Attr {
				if attr.Key == "id" {
					id = attr.Val
					break
				}
			}
			
			// If h2 has ID, extract the text content as title
			if id != "" {
				// Extract text content from h2 element
				var textBuilder strings.Builder
				var extractText func(*html.Node)
				extractText = func(node *html.Node) {
					if node.Type == html.TextNode {
						textBuilder.WriteString(node.Data)
					}
					for child := node.FirstChild; child != nil; child = child.NextSibling {
						extractText(child)
					}
				}
				extractText(n)
				title := strings.TrimSpace(textBuilder.String())
				
				if title != "" {
					toc = append(toc, TOCEntry{
						ID:    id,
						Title: title,
						Level: 2,
					})
				}
			}
		}
		
		// Recursively walk children
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walker(child)
		}
	}
	
	walker(doc)
	return toc, nil
}


// generateHeadingID generates a heading ID from text (similar to goldmark's algorithm)
func generateHeadingID(text string) string {
	// Convert to lowercase
	id := strings.ToLower(text)
	
	// Replace spaces with hyphens
	id = strings.ReplaceAll(id, " ", "-")
	
	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	id = reg.ReplaceAllString(id, "")
	
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	id = reg.ReplaceAllString(id, "-")
	
	// Trim leading/trailing hyphens
	id = strings.Trim(id, "-")
	
	return id
}

// generateTitleFromID converts section IDs to readable titles
// e.g., "key-features" → "Key Features", "dns-settings" → "DNS Settings"
func generateTitleFromID(id string) string {
	// Replace hyphens with spaces
	title := strings.ReplaceAll(id, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	
	// Apply smart title casing
	return formatTitle(title)
}

// formatTitle formats titles for display
// Converts "pricing" to "Pricing", preserves acronyms like "DNS Settings"
func formatTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return title
	}
	
	// Smart title case: preserve acronyms (all caps words)
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			// Check if word is all uppercase (likely an acronym)
			isAcronym := true
			for _, char := range word {
				if char >= 'a' && char <= 'z' {
					isAcronym = false
					break
				}
			}
			
			// Preserve acronyms, otherwise apply title case
			if isAcronym && len(word) > 1 {
				words[i] = word // Keep as-is for acronyms
			} else {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
	}
	
	return strings.Join(words, " ")
}