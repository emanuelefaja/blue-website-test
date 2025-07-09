/**
 * Image Zoom/Lightbox Utilities
 * Handles adding zoom functionality to images in markdown content
 * Desktop only - disabled on mobile devices
 */
window.ImageZoomUtils = {
    /**
     * Initialize image zoom functionality
     */
    init() {
        // Skip on mobile devices (viewport width < 768px)
        if (window.innerWidth < 768) {
            return;
        }
        
        this.addImageZoom();
    },

    /**
     * Add zoom functionality to all prose images
     */
    addImageZoom() {
        // Select all images in prose content that aren't already handled or inside links
        const images = document.querySelectorAll('.prose img:not([data-zoom-added]):not(a img), .prose-lg img:not([data-zoom-added]):not(a img)');
        
        images.forEach(img => {
            // Mark as handled to prevent duplicate handlers
            img.setAttribute('data-zoom-added', 'true');
            
            // Add click handler
            img.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                
                // Dispatch custom event with image details
                window.dispatchEvent(new CustomEvent('image-zoom', {
                    detail: {
                        src: img.src,
                        alt: img.alt || ''
                    }
                }));
            });
        });
    }
};