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
     * Get current language from URL
     * @returns {string} Language code or empty string
     */
    getCurrentLanguage() {
        const path = window.location.pathname;
        const langMatch = path.match(/^\/([a-z]{2}|[a-z]{2}-[A-Z]{2})(\/|$)/);
        return langMatch ? langMatch[1] : '';
    },

    /**
     * Add language prefix to path if needed
     * @param {string} path - The path to process
     * @returns {string} Path with language prefix
     */
    addLanguagePrefix(path) {
        const currentLang = this.getCurrentLanguage();
        if (!currentLang) return path;
        
        // Check if path already has a language prefix
        const hasLangPrefix = /^\/([a-z]{2}|[a-z]{2}-[A-Z]{2})(\/|$)/.test(path);
        if (hasLangPrefix) return path;
        
        // Add current language prefix
        return path === '/' ? `/${currentLang}` : `/${currentLang}${path}`;
    },

    /**
     * Set up client-side routing for internal links
     */
    setupClientRouting() {
        const navLinks = document.querySelectorAll('a[href^="/"]:not([data-spa-handled])');
        navLinks.forEach(link => {
            // Skip external links, special links, download links, and pages with dynamic scripts
            if (link.target === '_blank' || link.href.includes('#') || link.hasAttribute('download')) return;
            
            // Skip SPA routing for pages that need full page loads
            const url = new URL(link.href);
            const skipSPA = [].includes(url.pathname); // No pages skip SPA
            if (skipSPA) return;
            
            // Mark as handled to prevent duplicate listeners
            link.setAttribute('data-spa-handled', 'true');
            
            // Prefetch on hover
            link.addEventListener('mouseenter', () => {
                const targetUrl = new URL(link.href);
                const languageAwarePath = this.addLanguagePrefix(targetUrl.pathname);
                targetUrl.pathname = languageAwarePath;
                
                const prefetchLink = document.createElement('link');
                prefetchLink.rel = 'prefetch';
                prefetchLink.href = targetUrl.href;
                document.head.appendChild(prefetchLink);
            });
            
            // Handle click
            link.addEventListener('click', async (e) => {
                // Allow modifier keys to work normally (Cmd/Ctrl + Click for new tab)
                if (e.metaKey || e.ctrlKey || e.shiftKey || e.button === 1) {
                    return; // Let browser handle it normally
                }
                
                e.preventDefault();
                await this.handleNavigation(link);
            });
        });
    },

    /**
     * Handle navigation to a new page
     * @param {HTMLElement} link - The clicked link element
     * @param {boolean} isLanguageSwitch - Whether this is a language switch
     */
    async handleNavigation(link, isLanguageSwitch = false) {
        try {
            // Add language prefix to the URL if needed
            const targetUrl = new URL(link.href);
            const languageAwarePath = this.addLanguagePrefix(targetUrl.pathname);
            targetUrl.pathname = languageAwarePath;
            
            const response = await fetch(targetUrl.href, {
                headers: {
                    'X-Requested-With': 'XMLHttpRequest'
                }
            });
            if (!response.ok) throw new Error('Network response was not ok');
            
            // Check if the page needs code highlighting
            const needsCodeHighlight = response.headers.get('X-Needs-Code-Highlight') === 'true';
            
            const html = await response.text();
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            
            if (isLanguageSwitch) {
                // For language switches, update more of the DOM
                await this.updateFullPageContent(doc, targetUrl, needsCodeHighlight);
            } else {
                // For regular navigation, just update main content
                const newMain = doc.querySelector('main');
                const currentMain = document.querySelector('main');
                
                if (newMain && currentMain) {
                    await this.updatePageContent(newMain, currentMain, targetUrl, doc, needsCodeHighlight);
                } else {
                    window.location.href = targetUrl.href;
                }
            }
        } catch (error) {
            console.error('Client routing failed:', error);
            // Add language prefix to fallback URL too
            const targetUrl = new URL(link.href);
            targetUrl.pathname = this.addLanguagePrefix(targetUrl.pathname);
            window.location.href = targetUrl.href;
        }
    },

    /**
     * Update full page content for language switching
     * @param {Document} doc - Parsed document from response
     * @param {URL} targetUrl - The target URL object
     * @param {boolean} needsCodeHighlight - Whether the page needs code highlighting
     */
    async updateFullPageContent(doc, targetUrl, needsCodeHighlight) {
        // Preserve current state
        const preservedState = {
            darkMode: localStorage.getItem('theme') === 'dark',
            customLogo: localStorage.getItem('customLogo'),
            sidebarOpen: document.querySelector('body')._x_dataStack?.[0]?.sidebarsOpen ?? true,
            searchQuery: document.getElementById('searchInput')?.value || '',
            authState: window.AuthStateManager?.getState() || null
        };

        // Get all components that need updating
        const componentsToUpdate = [
            { selector: '.fixed.top-0.left-0.right-0.z-\\[60\\]', name: 'top-news-bar' },
            { selector: 'header', name: 'topbar' },
            { selector: 'nav[x-data*="currentPath"]', name: 'left-sidebar' },
            { selector: 'main', name: 'main' },
            { selector: '.image-lightbox', name: 'lightbox' },
            { selector: '.fixed.right-0.top-0.bottom-0.bg-white', name: 'ai-sidebar' }
        ];

        // Fade out all components
        const fadePromises = componentsToUpdate.map(component => {
            const element = document.querySelector(component.selector);
            if (element) {
                element.style.transition = 'opacity 50ms ease-out';
                element.style.opacity = '0';
            }
        });

        // Wait for fade out
        await new Promise(resolve => setTimeout(resolve, 50));

        // Update each component
        componentsToUpdate.forEach(component => {
            const currentElement = document.querySelector(component.selector);
            const newElement = doc.querySelector(component.selector);
            
            if (currentElement && newElement) {
                // Replace the HTML content
                currentElement.outerHTML = newElement.outerHTML;
            }
        });

        // Update title and meta tags
        const newTitle = doc.querySelector('title');
        if (newTitle) document.title = newTitle.textContent;
        this.updateMetaTags(doc);

        // Update body data-language attribute
        const newBodyLang = doc.body.getAttribute('data-language');
        if (newBodyLang) {
            document.body.setAttribute('data-language', newBodyLang);
        }

        // Update URL
        window.history.pushState(null, '', targetUrl.href);

        // Restore preserved state
        if (preservedState.darkMode) {
            document.documentElement.classList.add('dark');
        }
        if (preservedState.customLogo) {
            localStorage.setItem('customLogo', preservedState.customLogo);
        }

        // Re-initialize Alpine.js components with preserved state
        setTimeout(async () => {
            // Re-initialize Alpine on updated elements
            componentsToUpdate.forEach(component => {
                const element = document.querySelector(component.selector);
                if (element && window.Alpine) {
                    window.Alpine.initTree(element);
                }
            });

            // Restore states after Alpine initialization
            const bodyData = document.querySelector('body')._x_dataStack?.[0];
            if (bodyData) {
                bodyData.sidebarsOpen = preservedState.sidebarOpen;
            }

            // Restore search query
            const searchInput = document.getElementById('searchInput');
            if (searchInput && preservedState.searchQuery) {
                searchInput.value = preservedState.searchQuery;
            }

            // Fade in all components
            componentsToUpdate.forEach(component => {
                const element = document.querySelector(component.selector);
                if (element) {
                    element.style.opacity = '1';
                }
            });

            // Load Highlight.js if needed before reinitializing
            if (needsCodeHighlight && typeof HighlightLoader !== 'undefined') {
                await HighlightLoader.load();
            }
            
            // Re-initialize content
            this.reinitializeContent();
            
            // Re-setup routing for new elements
            this.setupClientRouting();

            // Update sidebar state
            const pathname = targetUrl.pathname;
            const pathWithoutLang = pathname.replace(/^\/([a-z]{2}|[a-z]{2}-[A-Z]{2})(\/|$)/, '/');
            this.updateSidebarState(pathWithoutLang);

        }, 100);
    },

    /**
     * Update page content with smooth transition
     * @param {HTMLElement} newMain - New main content element
     * @param {HTMLElement} currentMain - Current main content element
     * @param {URL} targetUrl - The target URL object
     * @param {Document} doc - Parsed document from response
     * @param {boolean} needsCodeHighlight - Whether the page needs code highlighting
     */
    async updatePageContent(newMain, currentMain, targetUrl, doc, needsCodeHighlight) {
        // Much faster transition - 50ms total
        currentMain.style.opacity = '0';
        currentMain.style.transition = 'opacity 50ms ease-out';
        
        // Update content after very brief fade
        setTimeout(async () => {
            currentMain.innerHTML = newMain.innerHTML;
            // Instant fade in
            currentMain.style.opacity = '1';
            
            // Update title
            const newTitle = doc.querySelector('title');
            if (newTitle) document.title = newTitle.textContent;
            
            // Update meta tags for SEO
            this.updateMetaTags(doc);
            
            // Update body data-language attribute
            const newBodyLang = doc.body.getAttribute('data-language');
            if (newBodyLang) {
                document.body.setAttribute('data-language', newBodyLang);
            }
            
            // Update URL
            window.history.pushState(null, '', targetUrl.href);
            
            // Update sidebar active state - remove language prefix for sidebar matching
            const pathname = targetUrl.pathname;
            const pathWithoutLang = pathname.replace(/^\/([a-z]{2}|[a-z]{2}-[A-Z]{2})(\/|$)/, '/');
            this.updateSidebarState(pathWithoutLang);
            
            // Collapse sidebar menus when navigating to home
            if (pathWithoutLang === '/') {
                window.dispatchEvent(new CustomEvent('collapse-sidebar-menus'));
            }
            
            // Load Highlight.js if needed before reinitializing
            if (needsCodeHighlight && typeof HighlightLoader !== 'undefined') {
                await HighlightLoader.load();
            }
            
            // Re-initialize content
            this.reinitializeContent();
            
            // For pages with Alpine.js content, re-setup after Alpine renders
            if (pathWithoutLang === '/insights' || pathWithoutLang === '/platform/integrations' || pathWithoutLang === '/platform/api' || pathWithoutLang === '/platform/status') {
                setTimeout(() => {
                    this.setupClientRouting();
                }, 100);
            }
            
            // Scroll to top
            window.scrollTo({ top: 0, behavior: 'smooth' });
        }, 50);
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
        // Use centralized initialization
        if (typeof BlueInit !== 'undefined') {
            BlueInit.reinitAfterNavigation();
        } else {
            console.error('BlueInit not available - initialization may fail');
        }
    },

    /**
     * Set up browser back/forward button handling
     */
    setupPopStateHandler() {
        window.addEventListener('popstate', () => {
            // Update sidebar active state for browser navigation - remove language prefix
            const pathname = window.location.pathname;
            const pathWithoutLang = pathname.replace(/^\/([a-z]{2}|[a-z]{2}-[A-Z]{2})(\/|$)/, '/');
            this.updateSidebarState(pathWithoutLang);
            window.location.reload();
        });
    }
};