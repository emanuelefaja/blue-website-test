/**
 * White Label Logo Upload Functionality
 * Handles logo upload, preview, and localStorage management
 */
window.logoUploader = function() {
    return {
        isDragging: false,
        previewUrl: null,
        message: '',
        messageType: '',
        
        init() {
            // Load existing logo from localStorage
            const savedLogo = localStorage.getItem('customLogo');
            if (savedLogo) {
                this.previewUrl = savedLogo;
                this.showMessage('Custom logo loaded', 'success');
            }
        },
        
        handleDrop(event) {
            this.isDragging = false;
            const files = event.dataTransfer.files;
            if (files.length > 0) {
                this.processFile(files[0]);
            }
        },
        
        handleFileSelect(event) {
            const files = event.target.files;
            if (files.length > 0) {
                this.processFile(files[0]);
            }
        },
        
        processFile(file) {
            // Validate file type
            const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/svg+xml'];
            if (!validTypes.includes(file.type)) {
                this.showMessage('Please upload a PNG, JPG, or SVG file', 'error');
                return;
            }
            
            // Validate file size (500KB)
            if (file.size > 500 * 1024) {
                this.showMessage('File size must be less than 500KB', 'error');
                return;
            }
            
            // Read and display the file
            const reader = new FileReader();
            reader.onload = (e) => {
                this.previewUrl = e.target.result;
                localStorage.setItem('customLogo', e.target.result);
                this.showMessage('Logo uploaded successfully!', 'success');
                
                // Dispatch event to update topbar
                window.dispatchEvent(new CustomEvent('logoChanged', { 
                    detail: { logoUrl: e.target.result } 
                }));
            };
            reader.readAsDataURL(file);
        },
        
        resetLogo() {
            this.previewUrl = null;
            localStorage.removeItem('customLogo');
            this.showMessage('Reset to default Blue logo', 'success');
            
            // Dispatch event to update topbar
            window.dispatchEvent(new CustomEvent('logoChanged', { 
                detail: { logoUrl: null } 
            }));
        },
        
        showMessage(text, type) {
            this.message = text;
            this.messageType = type;
            setTimeout(() => {
                this.message = '';
            }, 3000);
        }
    }
};