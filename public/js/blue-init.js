/**
 * Blue Initialization Module
 * Centralizes all initialization logic for the Blue website
 * Handles both initial page load and SPA navigation reinitialization
 */
window.BlueInit = {
    /**
     * Initialize all components on initial page load
     * Called from main.html when DOM is ready
     */
    initAll() {
        console.log('🚀 Initializing Blue components...');
        
        // Initialize Highlight.js
        this.initHighlightJS();
        
        // Initialize copy code buttons
        this.initCopyCode();
        
        // Initialize heading anchors
        this.initHeadingAnchors();
        
        // Initialize SPA routing
        this.initSPA();
        
        // Initialize image zoom
        this.initImageZoom();
        
        console.log('✅ Blue initialization complete');
    },
    
    /**
     * Reinitialize components after SPA navigation
     * Called from spa.js after content update
     */
    reinitAfterNavigation() {
        console.log('🔄 Reinitializing after navigation...');
        
        // Re-initialize syntax highlighting
        this.initHighlightJS();
        
        // Re-initialize copy buttons
        this.initCopyCode();
        
        // Re-initialize heading anchors
        this.initHeadingAnchors();
        
        // Re-initialize image zoom
        this.initImageZoom();
        
        // Re-initialize lazy loaded videos
        this.reinitLazyVideos();
        
        // Re-setup client routing (handled by SPA utils)
        if (typeof SPAUtils !== 'undefined' && SPAUtils.setupClientRouting) {
            SPAUtils.setupClientRouting();
        }
        
        // Re-initialize Alpine.js components for dynamic content
        if (typeof Alpine !== 'undefined') {
            // Initialize any new Alpine components in the main content area
            const main = document.querySelector('main');
            if (main) {
                Alpine.initTree(main);
            }
        }
    },
    
    /**
     * Initialize Highlight.js with retry logic
     */
    initHighlightJS() {
        if (typeof hljs !== 'undefined') {
            hljs.highlightAll();
        } else {
            // Retry if hljs is still loading
            setTimeout(() => this.initHighlightJS(), 100);
        }
    },
    
    /**
     * Initialize copy code buttons
     */
    initCopyCode() {
        if (typeof CopyCodeUtils !== 'undefined') {
            CopyCodeUtils.init();
        } else {
            console.warn('CopyCodeUtils not loaded - copy buttons will not work');
        }
    },
    
    /**
     * Initialize heading anchors
     */
    initHeadingAnchors() {
        if (typeof HeadingAnchors !== 'undefined') {
            HeadingAnchors.init();
        } else {
            console.warn('HeadingAnchors not loaded - heading anchor links will not work');
        }
    },
    
    /**
     * Initialize SPA routing
     */
    initSPA() {
        if (typeof SPAUtils !== 'undefined') {
            SPAUtils.init();
        } else {
            console.warn('SPAUtils not loaded - SPA navigation will not work');
        }
    },
    
    /**
     * Initialize image zoom functionality
     */
    initImageZoom() {
        if (typeof ImageZoomUtils !== 'undefined') {
            ImageZoomUtils.init();
        } else {
            console.warn('ImageZoomUtils not loaded - image zoom will not work');
        }
    },
    
    /**
     * Reinitialize lazy loaded videos after navigation
     */
    reinitLazyVideos() {
        const lazyVideos = document.querySelectorAll('[x-data*="loaded"]');
        lazyVideos.forEach(el => {
            if (el._x_dataStack) {
                el._x_dataStack[0].loaded = false;
            }
        });
    }
};