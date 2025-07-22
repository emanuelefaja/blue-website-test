/**
 * Auth Cookie Utility
 * Reads and parses the blueAuthState cookie set by the Blue app
 * This allows the marketing website to show personalized content
 */
window.AuthCookie = {
    /**
     * Read and parse the auth state cookie
     * @returns {Object|null} Auth state object with isLoggedIn and firstName, or null if not found
     */
    read() {
        const cookieValue = this.getCookieValue('blueAuthState');
        if (!cookieValue) {
            return null;
        }

        try {
            // The cookie value is URL encoded, so decode it first
            const decodedValue = decodeURIComponent(cookieValue);
            const authState = JSON.parse(decodedValue);
            
            // Validate the structure
            if (typeof authState === 'object' && 
                authState !== null && 
                typeof authState.isLoggedIn === 'boolean') {
                return {
                    isLoggedIn: authState.isLoggedIn,
                    firstName: authState.firstName || ''
                };
            }
            
            return null;
        } catch (error) {
            console.warn('Failed to parse auth cookie:', error);
            return null;
        }
    },

    /**
     * Get a cookie value by name
     * @param {string} name - Cookie name
     * @returns {string|null} Cookie value or null if not found
     */
    getCookieValue(name) {
        const cookies = document.cookie.split('; ');
        
        for (const cookie of cookies) {
            const [cookieName, cookieValue] = cookie.split('=');
            if (cookieName === name) {
                return cookieValue;
            }
        }
        
        return null;
    },

    /**
     * Check if user is logged in
     * @returns {boolean} True if logged in, false otherwise
     */
    isLoggedIn() {
        const authState = this.read();
        return authState?.isLoggedIn === true;
    },

    /**
     * Get user's first name
     * @returns {string} User's first name or empty string
     */
    getFirstName() {
        const authState = this.read();
        return authState?.firstName || '';
    }
};