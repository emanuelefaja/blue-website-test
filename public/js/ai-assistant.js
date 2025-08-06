// AI Assistant enhancement functionality
console.log('[AI Assistant] Script loading...');

(function() {
    'use strict';
    
    console.log('[AI Assistant] Script executing inside IIFE...');
    
    // Enhanced markdown processing with syntax highlighting
    window.formatAIMessage = function(content) {
        if (!content) return '';
        
        // First, escape any HTML in the content to prevent XSS
        content = escapeHtml(content);
        
        // Process code blocks with syntax highlighting
        content = content.replace(/```(\w+)?\n([\s\S]*?)```/g, function(match, lang, code) {
            const trimmedCode = code.trim();
            // We'll apply highlighting after the content is rendered
            const langClass = lang ? `language-${lang}` : 'language-plaintext';
            
            return `<div class="relative group my-3">
                <div class="absolute top-2 right-2 flex items-center gap-2">
                    ${lang ? `<span class="text-xs text-gray-400 px-2 py-1 bg-gray-800 rounded">${lang}</span>` : ''}
                    <button onclick="copyCodeToClipboard(this)" class="opacity-0 group-hover:opacity-100 transition-opacity p-1.5 bg-gray-700 hover:bg-gray-600 rounded text-gray-300 hover:text-white" title="Copy code">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                        </svg>
                    </button>
                </div>
                <pre class="bg-gray-900 dark:bg-gray-950 rounded-lg p-4 pr-20 overflow-x-auto"><code class="${langClass} text-sm" data-raw="${escapeHtml(trimmedCode)}">${trimmedCode}</code></pre>
            </div>`;
        });
        
        // Process inline code
        content = content.replace(/`([^`]+)`/g, '<code class="px-1.5 py-0.5 bg-gray-200 dark:bg-gray-700 rounded text-sm font-mono">$1</code>');
        
        // Process bold text (must come before italic to handle **text** correctly)
        content = content.replace(/\*\*([^*]+)\*\*/g, '<strong class="font-semibold">$1</strong>');
        
        // Process italic text
        content = content.replace(/\*([^*]+)\*/g, '<em>$1</em>');
        
        // Process links
        content = content.replace(/\[([^\]]+)\]\(([^\)]+)\)/g, '<a href="$2" class="text-brand-blue hover:underline" target="_blank" rel="noopener">$1</a>');
        
        // Process headings
        content = content.replace(/^### (.+)$/gm, '<h3 class="text-base font-semibold mt-3 mb-2">$1</h3>');
        content = content.replace(/^## (.+)$/gm, '<h2 class="text-lg font-semibold mt-4 mb-2">$1</h2>');
        content = content.replace(/^# (.+)$/gm, '<h1 class="text-xl font-bold mt-4 mb-2">$1</h1>');
        
        // Process lists - handle multi-line list items
        let lines = content.split('\n');
        let inList = false;
        let listType = null;
        let listItems = [];
        let processedLines = [];
        
        for (let i = 0; i < lines.length; i++) {
            let line = lines[i];
            let orderedMatch = line.match(/^(\d+)\.\s+(.+)$/);
            let unorderedMatch = line.match(/^[-*]\s+(.+)$/);
            
            if (orderedMatch) {
                if (!inList || listType !== 'ordered') {
                    if (inList) {
                        processedLines.push(wrapList(listItems, listType));
                        listItems = [];
                    }
                    inList = true;
                    listType = 'ordered';
                }
                listItems.push(orderedMatch[2]);
            } else if (unorderedMatch) {
                if (!inList || listType !== 'unordered') {
                    if (inList) {
                        processedLines.push(wrapList(listItems, listType));
                        listItems = [];
                    }
                    inList = true;
                    listType = 'unordered';
                }
                listItems.push(unorderedMatch[1]);
            } else {
                if (inList) {
                    processedLines.push(wrapList(listItems, listType));
                    listItems = [];
                    inList = false;
                    listType = null;
                }
                processedLines.push(line);
            }
        }
        
        if (inList) {
            processedLines.push(wrapList(listItems, listType));
        }
        
        content = processedLines.join('\n');
        
        // Process line breaks and paragraphs
        content = content.replace(/\n\n+/g, '</p><p class="mb-3">');
        content = content.replace(/\n/g, '<br>');
        
        // Wrap in paragraph tags if not already wrapped and not starting with a block element
        if (!content.match(/^<[hpuold]/)) {
            content = `<p class="mb-3">${content}</p>`;
        }
        
        return content;
    };
    
    // Helper function to wrap list items
    function wrapList(items, type) {
        const tag = type === 'ordered' ? 'ol' : 'ul';
        const listClass = type === 'ordered' ? 'list-decimal' : 'list-disc';
        const itemsHtml = items.map(item => `<li>${item}</li>`).join('');
        return `<${tag} class="${listClass} ml-5 my-2 space-y-1">${itemsHtml}</${tag}>`;
    }
    
    // Helper function to escape HTML
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    // Copy code to clipboard function
    window.copyCodeToClipboard = function(button) {
        const codeBlock = button.closest('.relative').querySelector('code');
        // Use the raw data if available, otherwise use textContent
        const text = codeBlock.dataset.raw || codeBlock.textContent || codeBlock.innerText;
        
        navigator.clipboard.writeText(text).then(() => {
            // Show success feedback
            const originalHTML = button.innerHTML;
            button.innerHTML = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>';
            button.classList.add('!text-green-400');
            
            setTimeout(() => {
                button.innerHTML = originalHTML;
                button.classList.remove('!text-green-400');
            }, 2000);
        }).catch(err => {
            console.error('Failed to copy code:', err);
            // Fallback to older method
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed';
            textArea.style.opacity = '0';
            document.body.appendChild(textArea);
            textArea.select();
            try {
                document.execCommand('copy');
                // Show success feedback
                const originalHTML = button.innerHTML;
                button.innerHTML = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>';
                button.classList.add('!text-green-400');
                
                setTimeout(() => {
                    button.innerHTML = originalHTML;
                    button.classList.remove('!text-green-400');
                }, 2000);
            } catch (err) {
                console.error('Fallback copy failed:', err);
            }
            document.body.removeChild(textArea);
        });
    };
    
    // Initialize syntax highlighting for AI responses
    window.highlightAICode = function(container) {
        if (window.hljs) {
            const codeBlocks = container.querySelectorAll('pre code:not(.hljs)');
            codeBlocks.forEach((block) => {
                // Store the raw content before highlighting
                if (!block.dataset.raw) {
                    block.dataset.raw = block.textContent;
                }
                hljs.highlightElement(block);
            });
        }
    };
    
    // Add keyboard shortcut to open AI sidebar - following the same pattern as topbar.html
    document.addEventListener('DOMContentLoaded', function() {
        console.log('[AI Assistant] Setting up keyboard shortcuts...');
        
        function handleAIKeyPress(event) {
            // Debug: Log all Cmd/Ctrl key combinations
            if (event.metaKey || event.ctrlKey) {
                console.log('[AI Assistant] Key pressed:', event.key, 'Meta:', event.metaKey, 'Ctrl:', event.ctrlKey);
            }
            
            // Check if the active element is an input, textarea, or any element within a form
            const activeElement = document.activeElement;
            const isFormElement = activeElement.tagName === 'INPUT' ||
                                activeElement.tagName === 'TEXTAREA' ||
                                activeElement.tagName === 'SELECT' ||
                                activeElement.isContentEditable ||
                                activeElement.closest('form') !== null;

            // Handle Cmd+I / Ctrl+I for AI sidebar (works globally except in form elements)
            if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 'i') {
                console.log('[AI Assistant] Cmd/Ctrl+I detected! isFormElement:', isFormElement);
                
                // Don't trigger if we're in a form element (to allow italic formatting)
                if (!isFormElement) {
                    event.preventDefault();
                    console.log('[AI Assistant] Attempting to toggle sidebar...');
                    
                    // Check if Alpine is available
                    if (window.Alpine) {
                        console.log('[AI Assistant] Alpine found');
                        if (window.Alpine.store('aiSidebar')) {
                            console.log('[AI Assistant] AI Sidebar store found, toggling...');
                            window.Alpine.store('aiSidebar').toggle();
                            console.log('[AI Assistant] Toggle complete, new state:', window.Alpine.store('aiSidebar').open);
                        } else {
                            console.error('[AI Assistant] AI Sidebar store NOT found!');
                            // Try to check if Alpine stores are initialized
                            console.log('[AI Assistant] Available Alpine stores:', Object.keys(window.Alpine.store()));
                        }
                    } else {
                        console.error('[AI Assistant] Alpine NOT found!');
                    }
                }
            }
        }

        document.addEventListener('keydown', handleAIKeyPress);
        console.log('[AI Assistant] Keyboard listener attached');
        
        // Also check Alpine status on load
        setTimeout(() => {
            console.log('[AI Assistant] Alpine check after timeout:');
            console.log('  - Alpine exists:', !!window.Alpine);
            if (window.Alpine) {
                console.log('  - Alpine started:', window.Alpine.version);
                try {
                    console.log('  - aiSidebar store exists:', !!window.Alpine.store('aiSidebar'));
                    if (window.Alpine.store('aiSidebar')) {
                        console.log('  - aiSidebar store state:', window.Alpine.store('aiSidebar').open);
                    }
                } catch (e) {
                    console.error('  - Error accessing store:', e);
                }
            }
        }, 1000);
    });
    
    // Auto-resize textarea
    window.autoResizeTextarea = function(textarea) {
        textarea.style.height = 'auto';
        textarea.style.height = Math.min(textarea.scrollHeight, 200) + 'px';
    };
})();