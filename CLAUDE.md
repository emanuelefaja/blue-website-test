- This is a new website for Blue, a b2b saas process management system 
- The idea is to make a site with a high degree of polish to be on a world class level like Stripe, OpenAI, and Linear. 
- I am a single developer, so this has to be simple and easy maintanble. 

## Rules

- Write and setup idiomatic golang


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
