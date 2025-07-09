/**
 * Single Page Application (SPA) Routing Utilities
 * Handles client-side navigation with smooth transitions
 */
window.SPAUtils = {
    /**
     * Initialize SPA routing system
     */
    init() {
        this.setupClientRouting();
        this.setupPopStateHandler();
        
        // Re-setup routing after Alpine.js renders (for dynamically shown/hidden content)
        setTimeout(() => {
            this.setupClientRouting();
        }, 100);
    },

    /**
     * Set up client-side routing for internal links
     */
    setupClientRouting() {
        const navLinks = document.querySelectorAll('a[href^="/"]:not([data-spa-handled])');
        navLinks.forEach(link => {
            // Skip external links, special links, and pages with dynamic scripts
            if (link.target === '_blank' || link.href.includes('#')) return;
            
            // Skip SPA routing for pages that need full page loads
            const url = new URL(link.href);
            const skipSPA = [].includes(url.pathname); // No pages skip SPA
            if (skipSPA) return;
            
            // Mark as handled to prevent duplicate listeners
            link.setAttribute('data-spa-handled', 'true');
            
            // Prefetch on hover
            link.addEventListener('mouseenter', () => {
                const prefetchLink = document.createElement('link');
                prefetchLink.rel = 'prefetch';
                prefetchLink.href = link.href;
                document.head.appendChild(prefetchLink);
            });
            
            // Handle click
            link.addEventListener('click', async (e) => {
                e.preventDefault();
                await this.handleNavigation(link);
            });
        });
    },

    /**
     * Handle navigation to a new page
     * @param {HTMLElement} link - The clicked link element
     */
    async handleNavigation(link) {
        try {
            const response = await fetch(link.href);
            if (!response.ok) throw new Error('Network response was not ok');
            
            const html = await response.text();
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            const newMain = doc.querySelector('main');
            const currentMain = document.querySelector('main');
            
            if (newMain && currentMain) {
                await this.updatePageContent(newMain, currentMain, link, doc);
            } else {
                window.location.href = link.href;
            }
        } catch (error) {
            console.error('Client routing failed:', error);
            window.location.href = link.href;
        }
    },

    /**
     * Update page content with smooth transition
     * @param {HTMLElement} newMain - New main content element
     * @param {HTMLElement} currentMain - Current main content element
     * @param {HTMLElement} link - The clicked link element
     * @param {Document} doc - Parsed document from response
     */
    async updatePageContent(newMain, currentMain, link, doc) {
        // Fade out current content
        currentMain.style.opacity = '0';
        currentMain.style.transition = 'opacity 150ms ease-out';
        
        // Update content after fade out
        setTimeout(() => {
            currentMain.innerHTML = newMain.innerHTML;
            // Fade in new content
            currentMain.style.opacity = '1';
            
            // Update title
            const newTitle = doc.querySelector('title');
            if (newTitle) document.title = newTitle.textContent;
            
            // Update meta tags for SEO
            this.updateMetaTags(doc);
            
            // Update URL
            window.history.pushState(null, '', link.href);
            
            // Update sidebar active state
            this.updateSidebarState(new URL(link.href).pathname);
            
            // Collapse sidebar menus when navigating to home
            if (new URL(link.href).pathname === '/') {
                window.dispatchEvent(new CustomEvent('collapse-sidebar-menus'));
            }
            
            // Re-initialize content
            this.reinitializeContent();
            
            // For pages with Alpine.js content, re-setup after Alpine renders
            const pathname = new URL(link.href).pathname;
            if (pathname === '/insights' || pathname === '/platform/integrations' || pathname === '/platform/api' || pathname === '/platform/status') {
                setTimeout(() => {
                    this.setupClientRouting();
                }, 100);
            }
            
            // Scroll to top
            window.scrollTo({ top: 0, behavior: 'smooth' });
        }, 150);
    },

    /**
     * Update meta tags for SEO
     * @param {Document} doc - Parsed document from response
     */
    updateMetaTags(doc) {
        ['description', 'og:title', 'og:description', 'og:url', 'twitter:title', 'twitter:description'].forEach(name => {
            const oldMeta = document.querySelector(`meta[name="${name}"], meta[property="${name}"]`);
            const newMeta = doc.querySelector(`meta[name="${name}"], meta[property="${name}"]`);
            if (oldMeta && newMeta) {
                if (oldMeta.hasAttribute('content')) {
                    oldMeta.content = newMeta.content;
                } else if (oldMeta.hasAttribute('property')) {
                    oldMeta.setAttribute('content', newMeta.content);
                }
            }
        });
    },

    /**
     * Update sidebar active state
     * @param {string} pathname - The new pathname
     */
    updateSidebarState(pathname) {
        const sidebar = document.querySelector('nav[x-data*="currentPath"]');
        if (sidebar && sidebar._x_dataStack && sidebar._x_dataStack[0].updateCurrentPath) {
            sidebar._x_dataStack[0].updateCurrentPath(pathname);
        }
    },

    /**
     * Re-initialize content after navigation
     */
    reinitializeContent() {
        // Re-initialize syntax highlighting with retry logic
        function initializeHighlightJS() {
            if (typeof hljs !== 'undefined') {
                hljs.highlightAll();
            } else {
                // Retry if hljs is still loading
                setTimeout(initializeHighlightJS, 100);
            }
        }
        initializeHighlightJS();
        
        // Re-initialize copy buttons with error handling
        if (typeof CopyCodeUtils !== 'undefined') {
            CopyCodeUtils.init();
        } else {
            console.warn('CopyCodeUtils not available during navigation - copy buttons may not work');
        }
        
        // Re-initialize image zoom with error handling
        if (typeof ImageZoomUtils !== 'undefined') {
            ImageZoomUtils.init();
        } else {
            console.warn('ImageZoomUtils not available during navigation - image zoom may not work');
        }
        
        // Re-setup client routing
        this.setupClientRouting();
    },

    /**
     * Set up browser back/forward button handling
     */
    setupPopStateHandler() {
        window.addEventListener('popstate', () => {
            // Update sidebar active state for browser navigation
            this.updateSidebarState(window.location.pathname);
            window.location.reload();
        });
    }
};