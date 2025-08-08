package web

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v3"
)

// Global variable to cache sprite version for callouts
var calloutSpriteVersion string

// getCalloutSpriteVersion gets or initializes the sprite version for callouts
func getCalloutSpriteVersion() string {
	if calloutSpriteVersion == "" {
		file, err := os.Open("public/icons/sprite.svg")
		if err != nil {
			calloutSpriteVersion = fmt.Sprintf("%d", time.Now().Unix())
			return calloutSpriteVersion
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			calloutSpriteVersion = fmt.Sprintf("%d", time.Now().Unix())
			return calloutSpriteVersion
		}

		calloutSpriteVersion = fmt.Sprintf("%d", info.ModTime().Unix())
	}
	return calloutSpriteVersion
}

// CalloutNode represents a callout node in the AST
type CalloutNode struct {
	ast.BaseBlock
	Icon    string
	Target  string
	To      string
	Content string
}

// Dump implements ast.Node.Dump
func (n *CalloutNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// Kind implements ast.Node.Kind
func (n *CalloutNode) Kind() ast.NodeKind {
	return CalloutNodeKind
}

// CalloutNodeKind is the node kind for callouts
var CalloutNodeKind = ast.NewNodeKind("CalloutNode")

// CalloutParser parses callout blocks
type CalloutParser struct{}

// Trigger returns trigger characters for the parser
func (s *CalloutParser) Trigger() []byte {
	return []byte{':'}
}

// Open checks if the current line contains a callout start
func (s *CalloutParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()

	// Look for ::callout pattern on its own line
	if !bytes.HasPrefix(bytes.TrimSpace(line), []byte("::callout")) {
		return nil, parser.NoChildren
	}

	// Advance past the ::callout line
	reader.Advance(segment.Len())

	// Collect all lines until we find the closing ::
	var allLines []string

	for {
		line, segment := reader.PeekLine()
		if segment.Len() == 0 {
			break
		}

		lineStr := strings.TrimSpace(string(line))

		// Check for end of callout
		if lineStr == "::" {
			reader.Advance(segment.Len())
			break
		}

		allLines = append(allLines, string(line))
		reader.Advance(segment.Len())
	}

	// Parse the collected content
	node := &CalloutNode{}

	if len(allLines) > 0 {
		content := strings.Join(allLines, "")

		// Check if there's frontmatter
		if strings.Contains(content, "---") {
			parts := strings.Split(content, "---")
			if len(parts) >= 3 {
				// Extract frontmatter (between first and second ---)
				frontmatterStr := strings.TrimSpace(parts[1])

				// Remove empty lines from frontmatter
				lines := strings.Split(frontmatterStr, "\n")
				var cleanLines []string
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						cleanLines = append(cleanLines, line)
					}
				}
				frontmatterStr = strings.Join(cleanLines, "\n")

				var frontmatter map[string]string
				if err := yaml.Unmarshal([]byte(frontmatterStr), &frontmatter); err == nil {
					node.Icon = frontmatter["icon"]
					node.Target = frontmatter["target"]
					node.To = frontmatter["to"]
				}

				// Content comes after the second ---
				if len(parts) > 2 {
					node.Content = strings.TrimSpace(strings.Join(parts[2:], "---"))
				}
			}
		} else {
			// No frontmatter, treat entire content as markdown
			node.Content = strings.TrimSpace(content)
		}
	}

	return node, parser.NoChildren
}

// Continue is not used for this parser
func (s *CalloutParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

// Close finalizes the node
func (s *CalloutParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	// Nothing to do
}

// CanInterruptParagraph returns true if this parser can interrupt paragraphs
func (s *CalloutParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine returns false
func (s *CalloutParser) CanAcceptIndentedLine() bool {
	return false
}

// CalloutRenderer renders callout nodes to HTML
type CalloutRenderer struct{}

// RegisterFuncs registers rendering functions
func (r *CalloutRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(CalloutNodeKind, r.renderCallout)
}

// renderCallout renders a callout to HTML
func (r *CalloutRenderer) renderCallout(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	callout, ok := node.(*CalloutNode)
	if !ok {
		return ast.WalkContinue, nil
	}

	// Use icon name directly
	iconName := callout.Icon

	// Convert markdown content to HTML using basic regex replacements
	contentHTML := callout.Content

	// Handle bold text **text**
	contentHTML = strings.ReplaceAll(contentHTML, "**", "</strong>")
	contentHTML = strings.ReplaceAll(contentHTML, "</strong>", "<strong>")
	// Fix the pattern for proper bold formatting
	contentHTML = strings.ReplaceAll(contentHTML, "<strong>", "**")

	// Use a more precise approach for bold formatting
	parts := strings.Split(contentHTML, "**")
	if len(parts) > 1 {
		contentHTML = ""
		for i, part := range parts {
			if i%2 == 0 {
				contentHTML += part
			} else {
				contentHTML += "<strong>" + part + "</strong>"
			}
		}
	}

	// Handle line breaks
	contentHTML = strings.ReplaceAll(contentHTML, "<br/>", "<br>")
	contentHTML = strings.ReplaceAll(contentHTML, "\n", "<br>")

	if callout.To != "" {
		// Make entire callout clickable
		target := ""
		if callout.Target == "_blank" {
			target = ` target="_blank" rel="noopener noreferrer"`
		}

		// Clickable callout container with hover effects
		w.WriteString(fmt.Sprintf(`<a href="%s"%s class="block no-underline group">`, callout.To, target))
		w.WriteString(`<div class="callout border border-gray-200 rounded-lg p-4 bg-gray-50 dark:bg-gray-800 dark:border-gray-700 my-4 hover:border-brand-blue hover:bg-blue-50 dark:hover:border-blue-400 dark:hover:bg-blue-900/20 transition-all duration-200">`)
		w.WriteString(`<div class="flex items-start gap-3">`)

		// Icon container
		w.WriteString(`<div class="w-10 h-10 bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg flex items-center justify-center flex-shrink-0">`)
		w.WriteString(fmt.Sprintf(`<svg class="w-5 h-5 text-brand-blue dark:text-brand-blue" fill="currentColor">
			<use href="/icons/sprite.svg?v=%s#%s"></use>
		</svg>`, getCalloutSpriteVersion(), iconName))
		w.WriteString(`</div>`)

		// Content with hover color change
		w.WriteString(`<div class="flex-1 text-gray-600 dark:text-gray-400 group-hover:text-brand-blue dark:group-hover:text-brand-blue transition-colors">`)
		w.WriteString(contentHTML)
		w.WriteString(`</div>`)

		// Close containers
		w.WriteString(`</div></div></a>`)
	} else {
		// Non-clickable callout container
		w.WriteString(`<div class="callout border border-gray-200 rounded-lg p-4 bg-gray-50 dark:bg-gray-800 dark:border-gray-700 my-4">`)
		w.WriteString(`<div class="flex items-start gap-3">`)

		// Icon container
		w.WriteString(`<div class="w-10 h-10 bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg flex items-center justify-center flex-shrink-0">`)
		w.WriteString(fmt.Sprintf(`<svg class="w-5 h-5 text-brand-blue dark:text-brand-blue" fill="currentColor">
			<use href="/icons/sprite.svg?v=%s#%s"></use>
		</svg>`, getCalloutSpriteVersion(), iconName))
		w.WriteString(`</div>`)

		// Content without hover effects
		w.WriteString(`<div class="flex-1 text-gray-600 dark:text-gray-400">`)
		w.WriteString(contentHTML)
		w.WriteString(`</div>`)

		// Close containers
		w.WriteString(`</div></div>`)
	}

	return ast.WalkContinue, nil
}

// CalloutExtension is the Goldmark extension for callouts
type CalloutExtension struct{}

// Extend extends the Goldmark parser with callout support
func (e *CalloutExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&CalloutParser{}, 200),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&CalloutRenderer{}, 200),
		),
	)
}

// NewCalloutExtension creates a new callout extension
func NewCalloutExtension() goldmark.Extender {
	return &CalloutExtension{}
}
