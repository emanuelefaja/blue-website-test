package web

import (
	"html/template"
	"sync"
	"time"
)

// CachedContent represents pre-rendered markdown content
type CachedContent struct {
	HTML        string
	Frontmatter *Frontmatter
	ModTime     time.Time
	FilePath    string
}

// MarkdownCache provides thread-safe caching for pre-rendered markdown content
type MarkdownCache struct {
	mu    sync.RWMutex
	cache map[string]*CachedContent
}

// HTMLCache provides thread-safe caching for pre-rendered HTML pages
type HTMLCache struct {
	mu    sync.RWMutex
	cache map[string]*CachedContent
}

// NewMarkdownCache creates a new markdown cache
func NewMarkdownCache() *MarkdownCache {
	return &MarkdownCache{
		cache: make(map[string]*CachedContent),
	}
}

// Get retrieves cached content by URL path
func (mc *MarkdownCache) Get(urlPath string) (*CachedContent, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	content, exists := mc.cache[urlPath]
	return content, exists
}

// Set stores pre-rendered content in cache
func (mc *MarkdownCache) Set(urlPath string, content *CachedContent) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.cache[urlPath] = content
}

// GetAll returns all cached content (for search indexing)
func (mc *MarkdownCache) GetAll() map[string]*CachedContent {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	result := make(map[string]*CachedContent)
	for k, v := range mc.cache {
		result[k] = v
	}
	return result
}

// Size returns the number of cached items
func (mc *MarkdownCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	return len(mc.cache)
}

// Clear removes all cached content
func (mc *MarkdownCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.cache = make(map[string]*CachedContent)
}

// GetHTML returns the cached HTML as template.HTML
func (cc *CachedContent) GetHTML() template.HTML {
	return template.HTML(cc.HTML)
}

// Delete removes a specific item from cache
func (mc *MarkdownCache) Delete(urlPath string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	delete(mc.cache, urlPath)
}

// GetCacheStats returns cache statistics
func (mc *MarkdownCache) GetCacheStats() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	totalSize := 0
	for _, content := range mc.cache {
		totalSize += len(content.HTML)
	}
	
	return map[string]interface{}{
		"count":     len(mc.cache),
		"totalSize": totalSize,
		"avgSize":   totalSize / max(len(mc.cache), 1),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewHTMLCache creates a new HTML cache
func NewHTMLCache() *HTMLCache {
	return &HTMLCache{
		cache: make(map[string]*CachedContent),
	}
}

// Get retrieves cached HTML content by URL path
func (hc *HTMLCache) Get(urlPath string) (*CachedContent, bool) {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	content, exists := hc.cache[urlPath]
	return content, exists
}

// Set stores pre-rendered HTML content in cache
func (hc *HTMLCache) Set(urlPath string, content *CachedContent) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	hc.cache[urlPath] = content
}

// GetAll returns all cached HTML content
func (hc *HTMLCache) GetAll() map[string]*CachedContent {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	result := make(map[string]*CachedContent)
	for k, v := range hc.cache {
		result[k] = v
	}
	return result
}

// Size returns the number of cached HTML items
func (hc *HTMLCache) Size() int {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	return len(hc.cache)
}

// Clear removes all cached HTML content
func (hc *HTMLCache) Clear() {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	hc.cache = make(map[string]*CachedContent)
}

// Delete removes a specific HTML item from cache
func (hc *HTMLCache) Delete(urlPath string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	delete(hc.cache, urlPath)
}

// GetCacheStats returns HTML cache statistics
func (hc *HTMLCache) GetCacheStats() map[string]interface{} {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	totalSize := 0
	for _, content := range hc.cache {
		totalSize += len(content.HTML)
	}
	
	return map[string]interface{}{
		"count":     len(hc.cache),
		"totalSize": totalSize,
		"avgSize":   totalSize / max(len(hc.cache), 1),
	}
}