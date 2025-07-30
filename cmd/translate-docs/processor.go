package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// PlaceholderType represents different types of content that need preservation
type PlaceholderType string

const (
	PlaceholderCodeBlock   PlaceholderType = "CB"
	PlaceholderInlineCode  PlaceholderType = "CODE"
	PlaceholderURL         PlaceholderType = "URL"
	PlaceholderEmail       PlaceholderType = "EMAIL"
	PlaceholderCallout     PlaceholderType = "CALLOUT"
	PlaceholderVideo       PlaceholderType = "VIDEO"
	PlaceholderLink        PlaceholderType = "LINK"
)

// Placeholder represents a masked piece of content
type Placeholder struct {
	ID           string
	Type         PlaceholderType
	Content      string
	Position     int
	LineNumber   int
	BeforeText   string // Context for recovery
	AfterText    string // Context for recovery
}

// DocumentProcessor handles the parsing and processing of documentation
type DocumentProcessor struct {
	parser goldmark.Markdown
}

// NewDocumentProcessor creates a new document processor
func NewDocumentProcessor() *DocumentProcessor {
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	
	return &DocumentProcessor{
		parser: md,
	}
}

// ProcessDocument extracts and masks technical content, returning masked content and placeholder map
func (p *DocumentProcessor) ProcessDocument(content string) (string, map[string]Placeholder, error) {
	placeholderMap := make(map[string]Placeholder)
	
	// Extract frontmatter first to preserve it separately
	frontmatter, bodyContent := p.extractFrontmatter(content)
	
	// Pre-process to handle special markdown extensions (callouts)
	bodyContent, calloutMap := p.extractCallouts(bodyContent, placeholderMap)
	
	// Extract video tags
	bodyContent, videoMap := p.extractVideos(bodyContent, placeholderMap)
	
	// Extract code blocks first (they're the largest structures)
	bodyContent, codeBlockMap := p.extractCodeBlocks(bodyContent, placeholderMap)
	
	// Extract inline code
	bodyContent, inlineCodeMap := p.extractInlineCode(bodyContent, placeholderMap)
	
	// Extract internal documentation links
	bodyContent, linkMap := p.extractDocLinks(bodyContent, placeholderMap)
	
	// Extract URLs and emails
	bodyContent = p.extractURLsAndEmails(bodyContent, placeholderMap)
	
	// Merge all placeholder maps
	for k, v := range calloutMap {
		placeholderMap[k] = v
	}
	for k, v := range videoMap {
		placeholderMap[k] = v
	}
	for k, v := range codeBlockMap {
		placeholderMap[k] = v
	}
	for k, v := range inlineCodeMap {
		placeholderMap[k] = v
	}
	for k, v := range linkMap {
		placeholderMap[k] = v
	}
	
	// Combine frontmatter with processed body
	finalContent := frontmatter + bodyContent
	
	return finalContent, placeholderMap, nil
}

// extractFrontmatter separates frontmatter from body content
func (p *DocumentProcessor) extractFrontmatter(content string) (string, string) {
	// Check if content starts with frontmatter
	if !strings.HasPrefix(content, "---") {
		return "", content
	}
	
	// Find the closing frontmatter delimiter
	lines := strings.Split(content, "\n")
	closingIndex := -1
	
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closingIndex = i
			break
		}
	}
	
	if closingIndex == -1 {
		// No closing delimiter found
		return "", content
	}
	
	// Extract frontmatter (including delimiters) and body
	frontmatterLines := lines[:closingIndex+1]
	bodyLines := lines[closingIndex+1:]
	
	// Process frontmatter to ensure it's translatable
	processedFrontmatter := p.processFrontmatterForTranslation(frontmatterLines)
	
	return processedFrontmatter, strings.Join(bodyLines, "\n")
}

// processFrontmatterForTranslation preserves frontmatter structure exactly
func (p *DocumentProcessor) processFrontmatterForTranslation(lines []string) string {
	// Create a version with only translatable parts exposed
	var translatableVersion []string
	translatableVersion = append(translatableVersion, "---")
	
	for i := 1; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "title:") {
			// Extract the title value for translation
			titleValue := strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			translatableVersion = append(translatableVersion, fmt.Sprintf("title: %s", titleValue))
		} else if strings.HasPrefix(line, "description:") {
			// Extract the description value for translation
			descValue := strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			translatableVersion = append(translatableVersion, fmt.Sprintf("description: %s", descValue))
		}
	}
	
	translatableVersion = append(translatableVersion, "---")
	
	// Return the translatable version
	return strings.Join(translatableVersion, "\n") + "\n"
}

// extractCallouts handles special ::callout blocks
func (p *DocumentProcessor) extractCallouts(content string, placeholderMap map[string]Placeholder) (string, map[string]Placeholder) {
	// Match callouts and process their content
	calloutPattern := regexp.MustCompile(`(?s)(::callout\s*)((?:---\s*[^-]+---\s*)?)(.*?)(::)`)
	localMap := make(map[string]Placeholder)
	
	content = calloutPattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the parts of the callout
		matches := calloutPattern.FindStringSubmatch(match)
		if len(matches) != 5 {
			return match
		}
		
		calloutStart := matches[1]  // ::callout
		frontmatter := matches[2]   // ---\nicon: ...\n---
		calloutContent := matches[3] // The actual content to translate
		calloutEnd := matches[4]     // ::
		
		// Create placeholders for the callout structure
		startPlaceholder := p.generatePlaceholder(PlaceholderCallout)
		startUUID := p.extractPlaceholderID(startPlaceholder)
		localMap[startUUID] = Placeholder{
			ID:      startUUID,
			Type:    PlaceholderCallout,
			Content: calloutStart,
		}
		
		endPlaceholder := p.generatePlaceholder(PlaceholderCallout)
		endUUID := p.extractPlaceholderID(endPlaceholder)
		localMap[endUUID] = Placeholder{
			ID:      endUUID,
			Type:    PlaceholderCallout,
			Content: calloutEnd,
		}
		
		// Handle frontmatter if present
		result := startPlaceholder
		if frontmatter != "" {
			fmPlaceholder := p.generatePlaceholder(PlaceholderCallout)
			fmUUID := p.extractPlaceholderID(fmPlaceholder)
			localMap[fmUUID] = Placeholder{
				ID:      fmUUID,
				Type:    PlaceholderCallout,
				Content: frontmatter,
			}
			result += fmPlaceholder
		}
		
		// Leave the content for translation but process any nested technical content
		result += calloutContent + endPlaceholder
		
		return result
	})
	
	return content, localMap
}

// extractVideos handles HTML video tags
func (p *DocumentProcessor) extractVideos(content string, placeholderMap map[string]Placeholder) (string, map[string]Placeholder) {
	// Match video blocks
	videoPattern := regexp.MustCompile(`(?s)<video[^>]*>.*?</video>`)
	localMap := make(map[string]Placeholder)
	
	content = videoPattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := p.generatePlaceholder(PlaceholderVideo)
		uuid := p.extractPlaceholderID(placeholder)
		localMap[uuid] = Placeholder{
			ID:      uuid,
			Type:    PlaceholderVideo,
			Content: match,
		}
		return placeholder
	})
	
	return content, localMap
}

// extractCodeBlocks extracts all code blocks and replaces them with placeholders
func (p *DocumentProcessor) extractCodeBlocks(content string, placeholderMap map[string]Placeholder) (string, map[string]Placeholder) {
	// Match code blocks with optional language and labels
	codeBlockPattern := regexp.MustCompile("(?s)```[^`]*```")
	localMap := make(map[string]Placeholder)
	
	content = codeBlockPattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := p.generatePlaceholder(PlaceholderCodeBlock)
		uuid := p.extractPlaceholderID(placeholder)
		localMap[uuid] = Placeholder{
			ID:      uuid,
			Type:    PlaceholderCodeBlock,
			Content: match,
		}
		return placeholder
	})
	
	return content, localMap
}

// extractInlineCode extracts inline code snippets
func (p *DocumentProcessor) extractInlineCode(content string, placeholderMap map[string]Placeholder) (string, map[string]Placeholder) {
	inlineCodePattern := regexp.MustCompile("`([^`]+)`")
	localMap := make(map[string]Placeholder)
	
	content = inlineCodePattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := p.generatePlaceholder(PlaceholderInlineCode)
		uuid := p.extractPlaceholderID(placeholder)
		localMap[uuid] = Placeholder{
			ID:      uuid,
			Type:    PlaceholderInlineCode,
			Content: match,
		}
		return placeholder
	})
	
	return content, localMap
}

// extractDocLinks extracts internal documentation links
func (p *DocumentProcessor) extractDocLinks(content string, placeholderMap map[string]Placeholder) (string, map[string]Placeholder) {
	// Match markdown links that point to /docs/...
	docLinkPattern := regexp.MustCompile(`\[([^\]]+)\]\((/docs/[^)]+)\)`)
	localMap := make(map[string]Placeholder)
	
	content = docLinkPattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract link text and URL
		matches := docLinkPattern.FindStringSubmatch(match)
		if len(matches) != 3 {
			return match
		}
		
		linkText := matches[1]
		linkURL := matches[2]
		
		// Create placeholder for the URL part only
		urlPlaceholder := p.generatePlaceholder(PlaceholderLink)
		uuid := p.extractPlaceholderID(urlPlaceholder)
		localMap[uuid] = Placeholder{
			ID:      uuid,
			Type:    PlaceholderLink,
			Content: linkURL,
		}
		
		// Return the link with translatable text but preserved URL
		return fmt.Sprintf("[%s](%s)", linkText, urlPlaceholder)
	})
	
	return content, localMap
}

// extractURLsAndEmails finds and masks URLs and email addresses
func (p *DocumentProcessor) extractURLsAndEmails(content string, placeholderMap map[string]Placeholder) string {
	// Email pattern
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	content = emailPattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := p.generatePlaceholder(PlaceholderEmail)
		uuid := p.extractPlaceholderID(placeholder)
		placeholderMap[uuid] = Placeholder{
			ID:      uuid,
			Type:    PlaceholderEmail,
			Content: match,
		}
		return placeholder
	})
	
	// URL pattern (http/https)
	urlPattern := regexp.MustCompile(`https?://[^\s\)]+`)
	content = urlPattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := p.generatePlaceholder(PlaceholderURL)
		uuid := p.extractPlaceholderID(placeholder)
		placeholderMap[uuid] = Placeholder{
			ID:      uuid,
			Type:    PlaceholderURL,
			Content: match,
		}
		return placeholder
	})
	
	return content
}

// generatePlaceholder creates a unique placeholder
func (p *DocumentProcessor) generatePlaceholder(placeholderType PlaceholderType) string {
	id := uuid.New().String()
	return fmt.Sprintf("@@%s##%s##%s@@", placeholderType, id, placeholderType)
}

// extractPlaceholderID extracts just the UUID from a full placeholder
func (p *DocumentProcessor) extractPlaceholderID(placeholder string) string {
	// Extract UUID from format @@TYPE##UUID##TYPE@@
	parts := strings.Split(strings.Trim(placeholder, "@"), "##")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// ValidateTranslation checks if all placeholders are preserved correctly
func (p *DocumentProcessor) ValidateTranslation(original, translated string) error {
	originalPlaceholders := p.extractPlaceholders(original)
	translatedPlaceholders := p.extractPlaceholders(translated)
	
	// Check count
	if len(originalPlaceholders) != len(translatedPlaceholders) {
		return fmt.Errorf("placeholder count mismatch: original=%d, translated=%d", 
			len(originalPlaceholders), len(translatedPlaceholders))
	}
	
	// Check exact matches
	for _, placeholder := range originalPlaceholders {
		found := false
		for _, transPlaceholder := range translatedPlaceholders {
			if placeholder == transPlaceholder {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing or modified placeholder: %s", placeholder)
		}
	}
	
	return nil
}

// extractPlaceholders finds all placeholders in content
func (p *DocumentProcessor) extractPlaceholders(content string) []string {
	placeholderPattern := regexp.MustCompile(`@@[A-Z]+##[a-f0-9-]+##[A-Z]+@@`)
	return placeholderPattern.FindAllString(content, -1)
}

// RestoreContent replaces all placeholders with their original content
func (p *DocumentProcessor) RestoreContent(content string, placeholderMap map[string]Placeholder) string {
	result := content
	
	// Restore placeholders by reconstructing the full placeholder format
	for uuid, data := range placeholderMap {
		fullPlaceholder := fmt.Sprintf("@@%s##%s##%s@@", data.Type, uuid, data.Type)
		result = strings.ReplaceAll(result, fullPlaceholder, data.Content)
	}
	
	return result
}

// RecoverPlaceholders attempts to recover corrupted placeholders
func (p *DocumentProcessor) RecoverPlaceholders(original, corrupted string) string {
	// Extract all original placeholders
	originalPlaceholders := p.extractPlaceholders(original)
	
	// Try to find and fix corrupted placeholders
	result := corrupted
	
	for _, placeholder := range originalPlaceholders {
		if !strings.Contains(result, placeholder) {
			// Try to find partial matches or translations
			parts := strings.Split(placeholder, "##")
			if len(parts) == 3 {
				id := parts[1]
				// Look for the ID in various forms
				patterns := []string{
					fmt.Sprintf("@@%s##%s##%s@@", ".*", id, ".*"),  // Any type with same ID
					fmt.Sprintf("@@.*##%s##.*@@", id),              // More flexible
					id,                                               // Just the ID
				}
				
				for _, pattern := range patterns {
					re := regexp.MustCompile(pattern)
					if matches := re.FindAllString(result, -1); len(matches) > 0 {
						// Replace the corrupted version with the correct one
						result = strings.Replace(result, matches[0], placeholder, 1)
						break
					}
				}
			}
		}
	}
	
	return result
}