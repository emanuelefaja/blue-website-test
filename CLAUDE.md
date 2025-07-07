# CLAUDE.md

## General
- This is a new website for Blue, a b2b saas process management system 
- The idea is to make a site with a high degree of polish to be on a world class level like Stripe, OpenAI, and Linear. People should really be amazed by the website. This should blow Asana, Clickup, Monday, and Trello out of the water.
- I am a single developer, so this has to be simple and easy maintanble. 

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

## Rules

- Write and setup idiomatic golang
- Tailwind CLI is alwasy in watch mode, no need for you to rebuild.
- Air is always running, so no need to build the server or run main.go
- We are using tailwind4, so there is no tailwind.config file
- When making components/sections, provide an ASCII mockup if appropriate before designing.


## Pages to work on

- Home page
- Compare Platforms
- All Feature pages
- platform/white-label
- platform/api
- platform/overview
- solutions/industry/government
- solutions/industry/non-profits
- solutions/industry/education
- solutions/industry/logistics
- solutions/industry/real-estate
- solutions/industry/professional-services
- solutions/industry/production
- solutions/industry/agencies
- solutions/industry/solopreneurs
- solutions/use-case/process-management
- solutions/use-case/project-management  
- solutions/use-case/sales-crm
- solutions/use-case/service-tickets
- solutions/use-case/it
- solutions/use-case/asset-management
- solutions/use-case/content-calendar
- solutions/use-case/personal-productivity

## To Go Live
- Also images/videos on documentation page require completely redoing, they are quite amatuer. 
- I need some more components for features I think
- Consider alignign titles tot he stat card components for other components, looks very elegant.
- Make all components mobile responive
- Clear all 404 issues.
- handle /docs/automations/actions/send-email in left sidebr etc
- Sort out card inconsistencies in brand page
- Align CTA in brand guidelines to backgroud color of FAQ backgrounds
- Preview image for all pages of the website (similar to linear)
- Make the video in markdown pages also have the same curved corners.
- Make testimonial components from the testimonials parts in the brand page.
- Check if these is a more specific javascript focus for the layout switch I am doing on main.
- Consider thin grey line (like in hero video) for FAQ and CTAs that are light grey?



##  Main Plan

- **Refactor web/router.go** - Extract template functions, page data logic, and service orchestration into separate services (currently 586 lines doing too much)
- Add 404 component to brand guidelines, looks quite ncie actually! 
- Improve experts page (later on)
- redirect bug: /docs redirects but /docs/ does not, this should work seamlessly
- lazy load components on scroll (or at least logo images) because excessive DOM size
- Polyfills and transforms enable legacy browsers to use new JavaScript features. However, many aren't necessary for modern browsers. Consider modifying your JavaScript build process to not transpile Baseline features, unless you know you must support legacy browser (related to fusejs)
- Check on curve of videos corners overlapping the border
- update policy on file storage and also FAQ
- Consider adding https://billing.blue.cc/p/login/14k7w17SddUW0yk288  link to the sidebar somewhere? 
- Switch logo grid to use tailwind styles later on
- Figure out a scalable way to handle the copy button for code blocks
- On dark mode, system status page load flashes light mode
- Make an awards component with crozdesk awards.
- Figure out a scalable way to hande the mute button on the hero video
- for /features ensure that there is no extra spacing on mobile for the column
- I saw that for the dynamic api claude recommende **File:** `/public/js/api-credentials.js` (new file), I wonder if I should centralize some other JS logic as well?
- Review what components to "borrow" from other sites.
- Add savings calculator to blue page (or not mentioned competitors at all?)
- Make the website multi linguagal?  Really like that to be honest.
- For api docs, put company or project id and all the mutations update in the examples. Like Stripe.
- Find all "big" brand customers via AI and add them to logos on customers page.
- Add SOP to website?
- status page link causes flash on sidebar (how to make part of SPA)
- Cehck security of APIs for status uptime, and perhaps expose them later on?
- Add terms about bypassing community etc. (David ho)
- Add Git commit history to website?
- Evetnaulyl split out professional services into multiple sub ones (law firms, accountants, doctors, etc)
- Add ability to edit logo size in grid manually pass through to component
- Remove unused images across various folders.
- Fix white label add on button on pricing page
- Customer stories (eventually)
- Changelog page from roadmap flahses the sidebar.
- Pull some latest reviews from Appsumo?
- Add GDOR and HIPPA complaince logos to the security page
- Topbar notice: https://devdojo.com/pines/docs/banner#
- Create a api endpoint on blue to count companies, and then use that to power website customer count.
- dual button CTA has far too much padding, but there seems to be a bug with p-16 that makes it looks super squashed.
- Improve insight SVG pattern generation (similar to openai patterns) or consider PNG generation instead because of guassian blur support?
- Review buffer about page for inspiration
- Add AI Chatbot 
- Confirm the paragraph text styles in the brand page.
- Consider back to top button like old site, but may get too busy?
- Add FAQ dropdowns to Brand Guidelines
- Add videos to FAQ possibily? 
- when pressing c to go to contact the sidebar flashes.
- image zoom effect found on https://linear.app/changelog
- switch page: https://linear.app/switch eventually here we should have details about buying out the entire contract.
- Blue University? Could be quite good huh? 
- Start/Setup guide: https://ghost.org/resources/
- Consider full width blue section that break out of the main content area, but when the sidebar goes over it, it turns white text instead of dark grey. Meh?
- Anthropic footer — looks great, perhaps in blue or dark blue?
- countdown to next update on status page?

# Ideas

- Sales modal


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
│   ├── style.css            # Compiled Tailwind CSS
│   ├── og.png              # Open Graph social media image
│   ├── sitemap.xml         # Generated sitemap
│   ├── robots.txt          # Search engine directives
│   ├── font/               # Inter font files (variable weight)
│   ├── logo/               # Brand assets
│   ├── product/            # Product screenshots and videos
│   └── ...                 # Other static assets
├── seo/
│   └── schema.json         # Structured data for search engines
├── web/                     # Go backend services (17 files)
│   ├── router.go           # HTTP routing and request handling
│   ├── html.go             # HTML page pre-rendering service
│   ├── markdown.go         # Markdown processing and caching
│   ├── seo.go              # SEO metadata and sitemap generation
│   ├── navigation.go       # Navigation structure management
│   ├── search.go           # Site search index generation
│   ├── status.go           # System status monitoring
│   └── ...                 # Other services
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
./tailwindcss -i public/input.css -o public/style.css --watch --minify

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
- Dynamic route handling in `web/router.go`
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

## Component Patterns

### Creating New Components
1. Add HTML file to `/components`
2. Use AlpineJS directives for interactivity
3. Include in layouts via `{{template "component-name" .}}`

### Dynamic Data in Components
- Pass data through Go template context
- Use `x-data` for AlpineJS state
- Global data available via `window.blueData`

## Important Files

### Core Configuration
- `main.go` - Entry point and server setup
- `web/router.go` - Route handling logic
- `web/markdown.go` - Markdown processing
- `tailwind.config.js` - Tailwind configuration

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


## self writing text

This is just a note on the implementation, we will use this in the future:

<div 
    x-data="{
        text: '',
        textArray : ['process', 'projects', 'customers', 'leads', 'records', 'purchase orders', 'tasks', 'teams', 'workflows', 'operations', 'logistics', 'marketing', 'milestones', 'deliverables', 'expectations', 'deadlines', 'work'],
        textIndex: 0,
        charIndex: 0,
        typeSpeed: 110,
        cursorSpeed: 550,
        pauseEnd: 1500,
        pauseStart: 20,
        direction: 'forward',
    }" 
    x-init="$nextTick(() => {
        let typingInterval = setInterval(startTyping, $data.typeSpeed);
    
        function startTyping(){
            let current = $data.textArray[ $data.textIndex ];
            
            // check to see if we hit the end of the string
            if($data.charIndex > current.length){
                    $data.direction = 'backward';
                    clearInterval(typingInterval);
                    
                    setTimeout(function(){
                        typingInterval = setInterval(startTyping, $data.typeSpeed);
                    }, $data.pauseEnd);
            }   
                
            $data.text = current.substring(0, $data.charIndex);
            
            if($data.direction == 'forward')
            {
                $data.charIndex += 1;
            } 
            else 
            {
                if($data.charIndex == 0)
                {
                    $data.direction = 'forward';
                    clearInterval(typingInterval);
                    setTimeout(function(){
                        $data.textIndex += 1;
                        if($data.textIndex >= $data.textArray.length)
                        {
                            $data.textIndex = 0;
                        }
                        typingInterval = setInterval(startTyping, $data.typeSpeed);
                    }, $data.pauseStart);
                }
                $data.charIndex -= 1;
            }
        }
                    
        setInterval(function(){
            if($refs.cursor.classList.contains('hidden'))
            {
                $refs.cursor.classList.remove('hidden');
            } 
            else 
            {
                $refs.cursor.classList.add('hidden');
            }
        }, $data.cursorSpeed);

    })"
    class="max-w-7xl mx-auto py-20">
    <div class="text-left">
        <h3 class="text-3xl md:text-4xl font-medium tracking-tight text-gray-900">The modern way<br>to manage
            <span class="relative inline-block">
                <span class="text-blue-600" x-text="text"></span>
                <span class="absolute right-1 top-0 h-full w-2 -mr-3 bg-brand-blue" x-ref="cursor"></span>
            </span>
        </h3>
    </div>
</div>