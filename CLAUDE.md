# CLAUDE.md

## General
- This is a new website for Blue, a b2b saas process management system 
- This is going to look enteprise-ready
- The idea is to make a site with a high degree of polish to be on a world class level like Stripe, OpenAI, and Linear, Figma, and so on. 
- People should really be amazed by the website. This should blow Asana, Clickup, Monday, and Trello out of the water.
- I am a single developer, so this has to be simple and easy maintanble long-term.

## Key Information

### Performance & Architecture
- **All pages are pre-rendered** at startup for maximum performance:
  - **HTML pages** (e.g., `/pages/index.html`) cached via `HTMLService` with full template processing
  - **Markdown content** (e.g., `/content/docs/`, `/content/insights/`) processed by Goldmark and cached via `MarkdownService`
  - **Automatic routing** - `/content/docs/introduction.md` becomes `/docs/introduction`, `/pages/pricing.html` becomes `/pricing`
  - **YAML frontmatter** in markdown files provides metadata (title, description, category, tags)
- **SPA-like experience** implemented in `layouts/main.html` using client-side routing:
  - Intercepts internal link clicks and fetches content via AJAX
  - Smoothly transitions between pages with fade effects (150ms)
  - Updates browser history and sidebar state
  - Re-initializes JavaScript (syntax highlighting, copy buttons, Alpine.js)
  - Prefetches pages on hover for instant navigation
  - Falls back to full page loads for dynamic pages like `/platform/status`
- **Hybrid content system**:
  - **Static HTML pages** for marketing/product pages with Go template variables
  - **Markdown pages** for documentation, blog posts, and content-heavy pages
  - **Automatic table of contents** generation for markdown content
  - **Component reuse** across both HTML and markdown pages
- **Component-based architecture** with reusable HTML components
- **In-memory caching** for all pre-rendered content for instant serving
- **Static asset optimization** with proper cache headers

## Rules

- Write and setup idiomatic golang
- Tailwind CLI is alwasy in watch mode, no need for you to rebuild.
- Air is always running, so no need to build the server or run main.go
- We are using tailwind4, so there is no tailwind.config file
- When making components/sections, provide an ASCII mockup if appropriate before designing.

## Tech Stack

### Backend (Go)
- **Go 1.24.4** - Main server language
- **github.com/yuin/goldmark** - Markdown processing with CommonMark compliance
- **gopkg.in/yaml.v3** - YAML parsing for frontmatter and configuration
- **github.com/joho/godotenv** - Environment variable loading from .env files
- **golang.org/x/net** - Extended networking libraries
- **net/http** (stdlib) - HTTP server and routing
- **html/template** (stdlib) - Template engine for server-side rendering
- **encoding/json** (stdlib) - JSON processing for APIs and data
- **path/filepath** (stdlib) - File path manipulation
- **Cloudflare D1** - Direct REST API integration for status monitoring (no SDK)

### Frontend
- **AlpineJS 3.x** - Lightweight JavaScript framework for interactivity
- **@alpinejs/collapse@3** - Collapse/expand animations
- **@alpinejs/intersect@3** - Intersection observer utilities
- **TailwindCSS v4** - Utility-first CSS framework (no config file)
- **Highlight.js 11.9.0** - Syntax highlighting for code blocks
- **Fuse.js 7.0.0** - Fuzzy search for site search functionality

### Development Tools
- **Air** - Hot reload for Go development
- **Tailwind CLI** - CSS compilation and watch mode


## Full Website Structure

```
.
├── components/                 # Reusable HTML Golang components
│   ├── head.html              # HTML head with meta tags and assets
│   ├── left-sidebar.html      # Documentation navigation
│   ├── right-sidebar.html     # Table of contents
│   ├── topbar.html           # Main navigation bar
│   ├── page-heading.html     # Page title and subtitle component
│   └── ...                   # Other UI components
├── content/                   # Markdown content (auto-rendered)
│   ├── api/                  # API documentation
│   ├── docs/                 # User documentation
│   ├── insights/             # Blog posts and insights (80+ articles)
│   ├── legal/                # Legal pages (terms, privacy, etc.)
│   └── project-management-dictionary/  # Glossary terms
├── data/                     # Configuration and metadata
│   ├── nav.json             # Navigation structure for all menus
│   ├── metadata.json        # Page titles, descriptions, SEO data
│   └── redirects.json       # URL redirect mappings
├── layouts/
│   └── main.html            # Main layout with SPA routing logic
├── pages/                    # Static HTML pages (file-based routing)
│   ├── index.html           # Homepage
│   ├── platform/            # Platform pages (features, changelog, etc.)
│   ├── solutions/           # Use case and industry solutions
│   ├── company/             # About, values, charter
│   └── ...                  # Other static pages
├── public/                   # Static assets served directly
│   ├── css/                # Stylesheets and CSS files
│   │   ├── style.css       # Compiled Tailwind CSS
│   │   ├── input.css       # Tailwind input file
│   │   └── highlight-github-dark.min.css # Syntax highlighting theme
│   ├── js/                 # JavaScript files
│   │   ├── vendor/         # Third-party JavaScript libraries
│   │   └── ...             # Custom JavaScript files
│   ├── og.png              # Open Graph social media image
│   ├── sitemap.xml         # Generated sitemap
│   ├── robots.txt          # Search engine directives
│   ├── font/               # Inter font files (variable weight)
│   ├── logo/               # Brand assets
│   ├── product/            # Product screenshots and videos
│   └── ...                 # Other static assets
├── seo/
│   └── schema.json         # Structured data for search engines
├── web/                     # Go backend services (30+ files)
│   ├── router.go           # Main HTTP router setup
│   ├── handlers.go         # HTTP request handlers
│   ├── static.go           # Static file serving
│   ├── security.go         # Security headers and middleware
│   ├── cache.go            # In-memory caching system
│   ├── html.go             # HTML page pre-rendering service
│   ├── markdown.go         # Markdown processing and caching
│   ├── content.go          # Content management utilities
│   ├── page_data.go        # Page data context building
│   ├── template_funcs.go   # Custom template functions
│   ├── seo.go              # SEO metadata and sitemap generation
│   ├── navigation.go       # Navigation structure management
│   ├── search.go           # Site search index generation
│   ├── status.go           # System status monitoring
│   ├── health.go           # Health check endpoint
│   ├── types.go            # Common type definitions
│   ├── utils.go            # Utility functions
│   ├── toc.go              # Table of contents generation
│   ├── callout.go          # Markdown callout processing
│   ├── youtube.go          # YouTube embed processing
│   ├── linkchecker.go      # Link validation
│   ├── png.go              # PNG asset handling
│   ├── schema.go           # Schema.org structured data
│   └── *_test.go           # Test files for each service
├── .air.toml               # Hot reload configuration
├── go.mod                  # Go module dependencies
├── main.go                 # Server entry point
├── tailwindcss             # Tailwind CLI binary
└── CLAUDE.md              # Project documentation (this file)
```

## Development Workflow

### Starting Development
```bash
# Start the Go server with hot reload
air

# Tailwind CSS is already in watch mode
# Access site at http://localhost:8080
```

### Common Commands
```bash
# Run Go server
go run main.go

# Tailwind Watch
./tailwindcss -i public/css/input.css -o public/css/style.css --watch --minify

# Build for production
go build -o blue-website

# Format Go code
go fmt ./...

# Run tests
go test ./...
```

## Key Architecture Decisions

### Routing System
- File-based routing from `/pages` directory
- Router setup in `web/router.go` with handlers in `web/handlers.go`
- Static file serving handled by `web/static.go`
- Security middleware in `web/security.go`
- HTML pages map directly to URLs (e.g., `/pages/pricing.html` → `/pricing`)

### Component System
- Server-side HTML components in `/components`
- AlpineJS for client-side interactivity
- Components are included via Go templates

### Data Management
- `data/nav.json` - Navigation structure for all menus
- `data/metadata.json` - Page titles, descriptions, and SEO data
- `seo/redirects.json` - URL redirect mapping

### Content Structure
- Markdown files in `/content` are rendered with Goldmark
- YAML frontmatter support for metadata
- Automatic table of contents generation for docs

### Multi-Language System
- **Language detection**: 
  - URL-based with language prefixes (e.g., `/en/`, `/zh/`, `/es/`)
  - Cookie preference (`lang` cookie) takes precedence
  - Falls back to browser Accept-Language header
  - Default to English if no preference detected
- **Translation structure**:
  - Modular JSON files in `/translations/` organized by section:
    - `common.json` - Shared UI elements (navigation, buttons, etc.)
    - `home.json` - Homepage specific translations
    - `features.json` - Features page translations
    - `about.json`, `values.json`, `charter.json` - Company pages
    - `search.json` - Search functionality
  - Each file contains language objects with nested keys for organization
- **Template functions**:
  - `{{t "section.key"}}` - For simple text translations
  - `{{t "section.key" "Default Text"}}` - With fallback default
  - `parseJSON` with `printf` for complex data structures with translations
- **Supported languages** (16 total, configured in `web/languages.go`):
  - en (English), zh (简体中文), es (Español), fr (Français)
  - de (Deutsch), ja (日本語), pt (Português), ru (Русский)
  - ko (한국어), it (Italiano), id (Indonesian), nl (Nederlands)
  - pl (Polski), zh-TW (繁體中文), sv (Svenska), km (ភាសាខ្មែរ)
- **Content structure**: 
  - HTML pages: Use translation keys with `{{t}}` function
  - Components: Receive translated data via template context
  - Markdown: Future support via language-specific directories
- **Language switching**:
  - Language picker in navbar shows native language names
  - Sets cookie and redirects to language-specific URL
  - Preserves current page when switching languages
- **SEO optimization**:
  - Proper `og:locale` meta tags for each language
  - Language-specific URLs for better indexing
  - Alternate language links in HTML head
- **Fallback behavior**: 
  - Missing translations show the key itself as fallback
  - All languages fall back to showing the translation key
  - Default text can be provided as second parameter to `{{t}}`

## Component Patterns

### Creating New Components
1. Add HTML file to `/components`
2. Use AlpineJS directives for interactivity
3. Include in layouts via `{{template "component-name" .}}`

### Dynamic Data in Components
- Pass data through Go template context
- Use `x-data` for AlpineJS state
- Global data available via `window.blueData`
- **IMPORTANT**: Use `parseJSON` with JSON strings for complex data structures (arrays, nested objects) - DO NOT use Go template `slice` or `dict` functions for arrays
- Simple single-level data can use `dict` function
- Example: `{{$logoData := parseJSON \`{"logos": [{"name": "Company", "src": "/logo.png"}]}\`}}`

### Component Data Access Rules
- **Components receive parsed data directly** - Access with Go template syntax (`.customers`, `.title`, etc.)
- **Do NOT use jsonify or other custom filters** - These don't exist in Go templates
- **Use {{range}} for iterations** - Not JavaScript loops in x-data
- **Mix Go templates with AlpineJS carefully**:
  - Go templates render server-side first
  - AlpineJS handles client-side interactivity after
  - Example: `{{range $index, $item := .items}}<div x-show="showAll || {{$index}} < limit">{{$item.name}}</div>{{end}}`
- **For conditional rendering based on data length**: Use `{{$total := len .items}}{{if gt $total 20}}...{{end}}`

## Important Files

### Core Configuration
- `main.go` - Entry point and server setup
- `web/router.go` - Main router setup and middleware
- `web/handlers.go` - HTTP request handlers
- `web/static.go` - Static file serving logic
- `web/security.go` - Security headers and CORS
- `web/cache.go` - In-memory caching system
- `web/markdown.go` - Markdown processing
- `web/page_data.go` - Page context building

### Key Components
- `layouts/main.html` - Main layout wrapper
- `components/topbar.html` - Navigation bar
- `components/left-sidebar.html` - Documentation navigation
- `components/right-sidebar.html` - Page table of contents

## Debugging Tips

### Common Issues
1. **404 Errors** - Check file exists in `/pages` or `/content`
2. **Styling Issues** - Ensure Tailwind CLI is running
3. **Component Not Rendering** - Verify template name matches
4. **Markdown Not Parsing** - Check frontmatter format

### Debug Mode
- Air provides hot reload and error messages
- Check browser console for AlpineJS errors
- Go template errors appear in terminal


## API Documentation Standards

Each API documentation page should follow this consistent structure to ensure clarity and completeness:

### Required Sections for API Pages

1. **Title and Description** (frontmatter)
   - Clear, action-oriented title (e.g., "Create a Project", "List all Projects")
   - Concise description of what the endpoint does

2. **Overview/Introduction**
   - Brief explanation of the endpoint's purpose
   - When and why to use it

3. **Basic Example**
   - Simple, minimal GraphQL query/mutation
   - Shows only required parameters
   - Easy to copy and test

4. **Advanced Example** (if applicable)
   - Complex example with all optional parameters
   - Real-world use case
   - Shows nested configurations

5. **Parameter Tables**
   - **Input Parameters**: All available input fields with:
     - Parameter name
     - Type (with GraphQL notation)
     - Required (✅ Yes / No)
     - Description
   - **Nested Types**: Separate tables for complex input types
   - **Enum Values**: Tables listing all possible values with descriptions

6. **Response Fields**
   - Table showing all fields in the response
   - Field types and descriptions
   - Note which fields are always present vs optional

7. **Permissions/Authorization**
   - Required company-level roles
   - Required project-level roles (if applicable)
   - Clear table showing which roles can/cannot perform the action

8. **Error Responses**
   - Common error scenarios with example JSON
   - Error codes and their meanings
   - How to handle each error type

9. **Important Notes/Considerations**
   - Performance implications
   - Limitations (e.g., max records, rate limits)
   - Best practices
   - Related endpoints

### Example Structure Template

```markdown
---
title: [Action] [Resource]
description: [What this endpoint does]
---

## [Action] a [Resource]

[Brief explanation of what this endpoint does and when to use it]

### Basic Example

```graphql
[Minimal working example]
```

### Advanced Example

```graphql
[Complex example with all options]
```

## Input Parameters

### [InputType]

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `field1` | String! | ✅ Yes | Description of field1 |
| `field2` | Boolean | No | Description of field2 |

### [EnumType] Values

| Value | Description |
|-------|-------------|
| `VALUE1` | What VALUE1 means |
| `VALUE2` | What VALUE2 means |

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier |
| `success` | Boolean! | Operation success status |

## Required Permissions

[Explain permission requirements]

| Role | Can Perform Action |
|------|-------------------|
| `OWNER` | ✅ Yes |
| `ADMIN` | ✅ Yes |
| `MEMBER` | ❌ No |

## Error Responses

### [Error Type]
```json
{
  "errors": [{
    "message": "Error message",
    "extensions": {
      "code": "ERROR_CODE"
    }
  }]
}
```
