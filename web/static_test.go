package web

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCacheFileServer(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"test.css":  "body { color: red; }",
		"test.js":   "console.log('test');",
		"test.png":  "fake-png-data",
		"test.html": "<html></html>",
	}

	for filename, content := range testFiles {
		filePath := tempDir + "/" + filename
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Create cache file server
	cacheFS := NewCacheFileServer(tempDir)

	// Test cases
	testCases := []struct {
		path           string
		expectedCache  bool
		expectedMaxAge string
		description    string
	}{
		{
			path:           "/test.css",
			expectedCache:  true,
			expectedMaxAge: "max-age=86400",
			description:    "CSS should have 1 day cache",
		},
		{
			path:           "/test.js",
			expectedCache:  true,
			expectedMaxAge: "max-age=604800",
			description:    "JS should have 1 week cache",
		},
		{
			path:           "/test.png",
			expectedCache:  true,
			expectedMaxAge: "max-age=31536000",
			description:    "PNG should have 1 year cache",
		},
		{
			path:           "/test.html",
			expectedCache:  false,
			expectedMaxAge: "",
			description:    "HTML should have no cache",
		},
		{
			path:           "/nonexistent.txt",
			expectedCache:  false,
			expectedMaxAge: "",
			description:    "Unknown file type should have no cache",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			// Serve the request
			cacheFS.ServeHTTP(w, req)

			// Check cache headers
			cacheControl := w.Header().Get("Cache-Control")

			if tc.expectedCache {
				if cacheControl == "" {
					t.Errorf("Expected cache header for %s, got none", tc.path)
				}
				if !strings.Contains(cacheControl, tc.expectedMaxAge) {
					t.Errorf("Expected %s in cache header for %s, got: %s", tc.expectedMaxAge, tc.path, cacheControl)
				}
			} else {
				// For non-existent files, we might not get any cache headers
				// since http.FileServer handles 404s differently
				if cacheControl != "" && !strings.Contains(cacheControl, "no-cache") {
					t.Errorf("Expected no-cache header or no header for %s, got: %s", tc.path, cacheControl)
				}
			}
		})
	}
}

func TestCacheFileServerDevelopment(t *testing.T) {
	// Set development environment
	os.Setenv("ENV", "development")
	defer os.Unsetenv("ENV")

	tempDir := t.TempDir()

	// Create test CSS file
	cssContent := "body { color: red; }"
	cssPath := tempDir + "/test.css"
	if err := os.WriteFile(cssPath, []byte(cssContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSS file: %v", err)
	}

	// Create cache file server
	cacheFS := NewCacheFileServer(tempDir)

	// Test that CSS has shorter cache in development
	req := httptest.NewRequest("GET", "/test.css", nil)
	w := httptest.NewRecorder()

	cacheFS.ServeHTTP(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	if !strings.Contains(cacheControl, "max-age=300") {
		t.Errorf("Expected 5-minute cache in development, got: %s", cacheControl)
	}
}

func TestCachePolicyMethods(t *testing.T) {
	cacheFS := NewCacheFileServer("public/")

	// Test getting existing policy
	policy, exists := cacheFS.GetCachePolicy(".css")
	if !exists {
		t.Error("Expected CSS policy to exist")
	}
	if policy.MaxAge != 86400 {
		t.Errorf("Expected CSS max-age to be 86400, got %d", policy.MaxAge)
	}

	// Test getting non-existent policy
	_, exists = cacheFS.GetCachePolicy(".xyz")
	if exists {
		t.Error("Expected non-existent policy to return false")
	}

	// Test setting custom policy
	customPolicy := CachePolicy{MaxAge: 3600, Public: true}
	cacheFS.SetCachePolicy(".custom", customPolicy)

	policy, exists = cacheFS.GetCachePolicy(".custom")
	if !exists {
		t.Error("Expected custom policy to exist after setting")
	}
	if policy.MaxAge != 3600 {
		t.Errorf("Expected custom max-age to be 3600, got %d", policy.MaxAge)
	}
}
