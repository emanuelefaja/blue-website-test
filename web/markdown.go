package web

import (
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownService handles markdown processing
type MarkdownService struct {
	markdown goldmark.Markdown
}

// NewMarkdownService creates a new markdown service
func NewMarkdownService() *MarkdownService {
	// Configure Goldmark with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			NewYouTubeExtension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	return &MarkdownService{
		markdown: md,
	}
}

// Convert converts markdown content to HTML
func (ms *MarkdownService) Convert(markdownContent []byte) (string, error) {
	var htmlBuffer strings.Builder
	if err := ms.markdown.Convert(markdownContent, &htmlBuffer); err != nil {
		return "", err
	}
	return htmlBuffer.String(), nil
}

// ProcessMarkdownFile reads a markdown file, parses frontmatter, and converts to HTML
func (ms *MarkdownService) ProcessMarkdownFile(filePath string, seoService *SEOService) (string, *Frontmatter, error) {
	// Read file
	mdBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	// Parse frontmatter
	var markdownContent []byte
	frontmatter, markdownContent, err := seoService.ParseFrontmatter(mdBytes)
	if err != nil {
		// Continue without frontmatter
		markdownContent = mdBytes
	}

	// Convert to HTML
	html, err := ms.Convert(markdownContent)
	if err != nil {
		return "", nil, err
	}

	return html, frontmatter, nil
}