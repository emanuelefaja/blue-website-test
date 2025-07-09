/**
 * Copy Code Button Utilities
 * Handles adding copy buttons to code blocks and managing clipboard operations
 */
window.CopyCodeUtils = {
    /**
     * Initialize copy buttons on all code blocks
     */
    init() {
        this.addCopyButtons();
    },

    /**
     * Add copy buttons to all code blocks that don't already have them
     */
    addCopyButtons() {
        const codeBlocks = document.querySelectorAll('.prose pre:not([data-copy-added])');
        codeBlocks.forEach(pre => {
            pre.setAttribute('data-copy-added', 'true');
            
            const copyBtn = document.createElement('button');
            copyBtn.className = 'code-copy-btn';
            copyBtn.innerHTML = '<svg class="w-4 h-4" fill="currentColor"><use href="/icons/sprite.svg#copy"></use></svg>';
            copyBtn.title = 'Copy code';
            
            copyBtn.addEventListener('click', async () => {
                const code = pre.querySelector('code');
                const text = code ? code.textContent : pre.textContent;
                
                try {
                    await navigator.clipboard.writeText(text);
                    copyBtn.innerHTML = '<svg class="w-4 h-4" fill="currentColor"><use href="/icons/sprite.svg#check"></use></svg>';
                    copyBtn.classList.add('copied');
                    copyBtn.title = 'Copied!';
                    
                    setTimeout(() => {
                        copyBtn.innerHTML = '<svg class="w-4 h-4" fill="currentColor"><use href="/icons/sprite.svg#copy"></use></svg>';
                        copyBtn.classList.remove('copied');
                        copyBtn.title = 'Copy code';
                    }, 2000);
                } catch (err) {
                    console.error('Failed to copy:', err);
                }
            });
            
            pre.appendChild(copyBtn);
        });
    }
};