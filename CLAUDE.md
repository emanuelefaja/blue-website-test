- This is a new website for Blue, a b2b saas process management system 
- The idea is to make a site with a high degree of polish to be on a world class level like Stripe, OpenAI, and Linear. People should really be amazed by the website.
- I am a single developer, so this has to be simple and easy maintanble. 

## Rules

- Write and setup idiomatic golang
- Tailwind CLI is alwasy in watch mode, no need for you to rebuild.


## Plan


- Prose styling
- dark mode implementation
- Sales modal? ❌ (decided not to do it for now)
- Pricing Page
- Review what components to "borrow" from other sites.
- Align all elements in respective footers perfectly. 
- Individual insight 
- Brand page
- Home page structure / Messaging
- Make the website multi linguagal?  Really like that to be honest.
- For api docs, put company or project id and all the mutations update in the examples. Like Stripe.
- Find all "big" brand customers via AI and add them to logos on customers page.
- Add SOP to website?
- Add Git commit history to website?
- Fix white label add on button
- Status page with backend
- Customer stories (eventually)
- Create a api endpoint on blue to count companies, and then use that to power website customer count.
- Status page improvement, show live operational status of services with last checked time.
- handle /docs/automations/actions/send-email in left sidebr etc
- handle case where users go direct to docs/automations for instance
- Improve insight SVG pattern generation (similar to openai patterns) or consider PNG generation instead because of guassian blur support?
- Figure out insights ordering (latest first)
- Consider centralizing markdown content like the about page.
- Sort out card inconsistencies in brand page
- Review buffer about page for inspiration
- Highlighting text should be brand-blue 
- system status does not show on search
- insights on search I see code.
- Confirm the paragraph text styles in the brand page.
- cmd+C goes to contact page instead of copying.
- Consider back to top button like old site, but may get too busy?
- Align CTA in brand guidelines to backgroud color of FAQ backgrounds
- Add FAQ dropdowns to Brand Guidelines
- Figure out what I am doing with changelog router and if there is an easier way
- Add videos to FAQ possibily? 
- Mouse pointer when on FAQ question should be hand
- animation for heading title
- subtle animation for right sidebar?
- image zoom effect found on https://linear.app/changelog
- switch page: https://linear.app/switch
- Blue University? Could be quite good huh? 
- Start/Setup guide: https://ghost.org/resources/
- Consider full width blue section that break out of the main content area, but when the sidebar goes over it, it turns white text instead of dark grey. Meh?
- Make the video in markdown pages also have the same curved corners.
- Make testimonial components from the testimonials parts in the brand page.
- Reconsider the /platform dropdown and what to put on there, perhaps "how it works" and key features, not all features




## Tech Stack

- Golang
- AlpineJS
- @alpinejs/collapse@3 plugin
- https://github.com/markmead/alpinejs-component (for components)
- TailwindCSS v4


## Folder Structure

- Layouts
- Pages (main pages and subfolder with pages, we will have route based navigation)
- Public (all static assets, images, etc)
- SEO (files for SEO like metadata and redirects, schema, robots, etc)

## Full Website Structure

```
.
├── components/
│   ├── client-logos.html
│   ├── head.html
│   ├── left-sidebar.html
│   ├── page-heading.html
│   ├── right-sidebar.html
│   ├── testimonial-videos.html
│   └── topbar.html
├── content/
│   ├── agency-success-guide/
│   ├── alternatives/
│   ├── api-docs/
│   │   ├── 1.start-guide/
│   │   ├── 11.libraries/
│   │   ├── 2.projects/
│   │   ├── 3.records/
│   │   ├── 5.custom fields/
│   │   ├── 6.automations/
│   │   ├── 7.user management/
│   │   ├── 8.company-management/
│   │   └── 9.dashboards/
│   ├── company-news/
│   ├── docs/
│   │   ├── 1.start-guide/
│   │   ├── 10.use cases/
│   │   ├── 2.projects/
│   │   ├── 3.records/
│   │   ├── 4.views/
│   │   ├── 5.custom fields/
│   │   ├── 6.automations/
│   │   │   └── 4.actions/
│   │   ├── 7.user management/
│   │   │   └── 8.roles/
│   │   ├── 8.dashboards/
│   │   └── 9.integrations/
│   ├── engineering-blog/
│   ├── frequently-asked-questions/
│   ├── insights/
│   ├── legal/
│   ├── modern-work-practices/
│   ├── product-updates/
│   ├── project-management-dictionary/
│   └── tips-&-tricks/
├── data/
│   ├── nav.json
│   ├── metadata.json
│   └── (other data files)
├── layouts/
│   └── main.html
├── pages/
│   ├── company/
│   ├── platform/
│   └── (main page HTML files)
├── public/
│   ├── customers/
│   ├── logo/
│   ├── testimonials/
│   └── (static assets)
├── seo/
│   ├── redirects.json
│   └── (SEO-related files)
├── web/
│   └── (Go web package files)
├── go.mod
├── go.sum
├── main.go
├── tailwind.config.js
└── CLAUDE.md
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

## Future Enhancements Needed

### High Priority
- [ ] Internationalization (i18n) system
- [ ] API endpoint for live customer count
- [ ] Backend for status page
- [ ] Search improvements (fuzzy matching, filters)

### Medium Priority
- [ ] Dark mode implementation
- [ ] Component library documentation
- [ ] Performance optimization (caching, CDN)
- [ ] Analytics integration

### Nice to Have
- [ ] Git commit history display
- [ ] Multi-language support
- [ ] Advanced search with Algolia/Meilisearch
- [ ] A/B testing framework
