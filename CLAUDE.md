- This is a new website for Blue, a b2b saas process management system 
- The idea is to make a site with a high degree of polish to be on a world class level like Stripe, OpenAI, and Linear. 
- I am a single developer, so this has to be simple and easy maintanble. 

## Rules

- Write and setup idiomatic golang
- Tailwind CLI is alwasy in watch mode, no need for you to rebuild.


## Plan

1. Setup initial folder structure ✅
2. Setup Golang ✅
3. Install Goldmark, Air, godotev, https://github.com/go-yaml/yaml ✅
4. Setup file based route navigation for /pages with a dedicated router inside a web package ✅
5. Setup main.html with left/right sidebar and navbar components ✅
6. Add favicon.ico and logo.svg ✅
7. Added data/nav.json for dynamic rendering of navigation acros components. ✅
8. Made topbar buttons hotkeys work ✅
9. Make left/right sidebar toggle open/close ✅
10. Setup TailwindCSS Properly with CLI and Watcher ✅
11. Dynamic left sidebar ✅
12. Metadata.json for dynmaic titles and so on ✅
13. Markdown generation ✅
14. Dynamic sidebar for docs and api docs ✅
15. Centralized redirets in redirects.json ✅
16. Prose styling
17. Right sidebar dynamic both for html pages and markdown files with up and down arrows
18. Replace font-awesome with svgs for arrows and light/dark mode?
19. dark mode implementation
20. Make legal navigation dynamic and page rendering working ✅
21. Sales modal?
22. Changelog Implementation ✅
23. Pricing Page
24. Roadmap 
25. Review what components to "borrow" from other sites.
26. Realtime search ✅
27. YouTube widget for markdown ✅
28. Improve search





## Tech Stack

- Golang
- AlpineJS
- @alpinejs/collapse@3 plugin
- https://github.com/markmead/alpinejs-component (for components)
- TailwindCSS v4


## Folder Strucuture

- Layouts
- Pages (main pages and subfolder with pages, we will have route based navigation)
- Public (all static assets, images, etc)
- SEO (files for SEO like metadata and redirects, schema, robots, etc)
