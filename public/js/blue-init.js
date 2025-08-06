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
        
        // Initialize auth cookie checking (for components that need it)
        this.initAuthCookie();
        
        // Initialize font loader
        this.initFontLoader();
        
        // Initialize AI Assistant
        this.initAIAssistant();
        
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
        
        // Re-initialize auth cookie checking
        this.initAuthCookie();
        
        // Check and maintain font loading state
        this.initFontLoader();
        
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
                
                // Manually trigger x-init for status monitoring components
                this.initStatusMonitoring(main);
            }
        }
    },
    
    /**
     * Initialize Highlight.js with retry logic
     */
    initHighlightJS() {
        if (typeof hljs !== 'undefined') {
            hljs.highlightAll();
        }
        // If hljs is not loaded, skip initialization
        // It will be loaded dynamically if needed during SPA navigation
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
    },
    
    /**
     * Initialize auth cookie functionality
     */
    initAuthCookie() {
        if (typeof AuthCookie !== 'undefined' && typeof AuthStateManager !== 'undefined') {
            // Initialize the centralized auth state manager
            AuthStateManager.init();
            console.log('✅ Auth state manager initialized');
        } else {
            console.warn('AuthCookie or AuthStateManager not loaded - auth state detection will not work');
        }
    },
    
    /**
     * Initialize status monitoring components after SPA navigation
     * @param {HTMLElement} container - Container element to search within
     */
    initStatusMonitoring(container) {
        // Find all status monitoring components
        const statusComponents = container.querySelectorAll('[x-data*="StatusUtils.createStatusPageData"], [x-data*="StatusUtils.createApiStatusData"]');
        
        statusComponents.forEach(component => {
            // Get the Alpine data object
            const alpineData = component._x_dataStack?.[0];
            
            if (alpineData && typeof alpineData.loadData === 'function') {
                console.log('📊 Loading status data for component...');
                // Trigger data loading
                alpineData.loadData();
            }
        });
    },
    
    /**
     * Initialize font loader
     */
    initFontLoader() {
        if (typeof FontLoader !== 'undefined') {
            FontLoader.init();
        } else {
            console.warn('FontLoader not loaded - dynamic fonts will not work');
        }
    },
    
    /**
     * Initialize AI Assistant functionality
     */
    initAIAssistant() {
        // The AI assistant script is self-initializing with its own DOMContentLoaded handler
        // We just need to ensure Alpine store is available
        if (typeof Alpine !== 'undefined') {
            // Check if the store is already initialized (it should be from main.html)
            setTimeout(() => {
                if (Alpine.store('aiSidebar')) {
                    console.log('✅ AI Assistant store verified');
                } else {
                    console.warn('AI Assistant store not found - keyboard shortcuts may not work');
                }
            }, 100);
        } else {
            console.warn('Alpine not loaded - AI Assistant will not work');
        }
    }
};