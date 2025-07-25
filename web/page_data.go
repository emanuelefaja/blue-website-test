package web

import (
	"html/template"
	"log"
	"sort"
	"strings"
)

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
func (r *Router) preparePageData(path string, content template.HTML, isMarkdown bool, frontmatter *Frontmatter, navigation *Navigation, lang string) PageData {
	// Get metadata from SEO service
	title, description, keywords, pageMeta, siteMeta := r.seoService.PreparePageMetadata(path, isMarkdown, frontmatter)

	// Prepare insights data - only include if on insights page
	var insights []InsightData
	if path == "/insights" {
		// Get all cached insights from MarkdownService
		cachedInsights := r.markdownService.GetAllCachedContent()

		// Initialize PNG generator
		pngGen := NewPNGGenerator()

		for urlPath, content := range cachedInsights {
			if strings.HasPrefix(urlPath, "/insights/") && content.Frontmatter != nil {
				// Extract category from tags or category field
				category := content.Frontmatter.Category
				if category == "" && len(content.Frontmatter.Tags) > 0 {
					// Fallback to first tag if category is not set
					category = content.Frontmatter.Tags[0]
				}

				// Generate unique PNG for this insight based on title
				pngPath, err := pngGen.GenerateOrGetPNG(content.Frontmatter.Title, content.Frontmatter.Slug)
				if err != nil {
					log.Printf("Error generating PNG for insight %s: %v", content.Frontmatter.Title, err)
					continue
				}

				insights = append(insights, InsightData{
					Title:       content.Frontmatter.Title,
					Description: content.Frontmatter.Description,
					Category:    category,
					Slug:        content.Frontmatter.Slug,
					PNGPath:     pngPath,
					Date:        content.Frontmatter.Date,
					URL:         urlPath,
				})
			}
		}
		
		// Sort insights by date in reverse chronological order (latest first)
		sort.Slice(insights, func(i, j int) bool {
			// Handle empty dates by putting them at the end
			if insights[i].Date == "" {
				return false
			}
			if insights[j].Date == "" {
				return true
			}
			// Compare dates as strings (assumes YYYY-MM-DD format)
			return insights[i].Date > insights[j].Date
		})
	}

	// Extract table of contents (skip if path is excluded)
	toc := make([]TOCEntry, 0)
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
			// TOC extracted
		}
	} else if r.isTOCExcluded(path) {
		// TOC generation skipped
	}

	// Generate schema data
	schemaData := template.JS("[]")
	if r.schemaService != nil {
		pageType := r.schemaService.GetPageType(path)
		schemaData = r.schemaService.GenerateSchema(pageType, path, frontmatter)
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
		TOC:            toc,
		CustomerNumber: 17000,
		Insights:       insights,
		Path:           path,
		SchemaData:     schemaData,
		Language:       lang,
	}
}