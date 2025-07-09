package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"blue-website/web"

	"github.com/joho/godotenv"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	startTime := time.Now()
	fmt.Printf("🚀 Starting Blue Website server at %s...\n", startTime.Format("15:04:05.000"))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Validate required environment variables
	requiredEnvVars := []string{
		"CLOUDFLARE_ACCOUNT_ID",
		"CLOUDFLARE_DATABASE_ID",
		"CLOUDFLARE_API_KEY",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Printf("Warning: %s environment variable not set, status monitoring disabled", envVar)
		}
	}

	// Parallelize independent startup tasks
	var wg sync.WaitGroup

	// Generate search index concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		searchStart := time.Now()
		fmt.Println("🔍 Generating search index...")
		if err := web.GenerateSearchIndex(); err != nil {
			log.Printf("⚠️  Warning: Failed to generate search index: %v", err)
		} else {
			fmt.Printf("✅ Search index generated successfully (took %v)\n", time.Since(searchStart))
		}
	}()

	// Generate sitemap concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		sitemapStart := time.Now()
		fmt.Println("🗺️  Generating sitemap...")
		seoService := web.NewSEOService()
		if err := seoService.LoadData(); err != nil {
			log.Printf("⚠️  Warning: Failed to load SEO data for sitemap: %v", err)
			return
		}
		if err := seoService.GenerateSitemap("https://blue.cc"); err != nil {
			log.Printf("⚠️  Warning: Failed to generate sitemap: %v", err)
		} else {
			fmt.Printf("✅ Sitemap generated successfully (took %v)\n", time.Since(sitemapStart))
		}
	}()

	// Wait for both tasks to complete before proceeding
	wg.Wait()

	// File-based routing handler
	routerStart := time.Now()
	fmt.Println("🛣️  Initializing router...")
	router := web.NewRouter("pages")
	fmt.Printf("✅ Router initialized (took %v)\n", time.Since(routerStart))

	// Run link checker in background after router is ready
	go func() {
		linkCheckerStart := time.Now()
		fmt.Println("🔗 Starting link checker...")

		// Create fresh services for link checker (they cache content internally)
		markdownService := web.NewMarkdownService()
		contentService := web.NewContentService("content")
		linkSeoService := web.NewSEOService()
		if err := linkSeoService.LoadData(); err != nil {
			log.Printf("⚠️  Warning: Failed to load SEO data for link checker: %v", err)
			return
		}
		htmlService := web.NewHTMLService("pages", "layouts", "components", markdownService)

		// Pre-render content for link checker
		if err := markdownService.PreRenderAllMarkdown(contentService, linkSeoService); err != nil {
			log.Printf("⚠️  Warning: Failed to pre-render markdown for link checker: %v", err)
			return
		}
		if err := htmlService.PreRenderAllHTMLPages(web.NewNavigationService(linkSeoService), linkSeoService); err != nil {
			log.Printf("⚠️  Warning: Failed to pre-render HTML for link checker: %v", err)
			return
		}

		// Run the link checker
		if err := web.RunLinkChecker(markdownService, htmlService, linkSeoService); err != nil {
			log.Printf("⚠️  Warning: link checker failed: %v", err)
		} else {
			fmt.Printf("✅ Link checker completed (took %v)\n", time.Since(linkCheckerStart))
		}
	}()

	// Initialize status monitoring in background if environment variables are set
	if os.Getenv("CLOUDFLARE_API_KEY") != "" {
		go func() {
			statusStart := time.Now()
			fmt.Println("🏥 Initializing status monitoring...")

			// Create D1 client
			d1Client := web.NewD1Client()

			// Create health checker
			healthChecker := web.NewHealthChecker(d1Client)

			// Initialize database and load historical data
			if err := healthChecker.Initialize(); err != nil {
				log.Printf("⚠️  Warning: Failed to initialize status monitoring: %v", err)
			} else {
				fmt.Printf("✅ Status monitoring initialized (took %v)\n", time.Since(statusStart))

				// Set the health checker in the router
				router.SetStatusChecker(healthChecker)

				// Start background health checks
				go func() {
					ticker := time.NewTicker(5 * time.Minute)
					defer ticker.Stop()

					// Run initial check only if needed
					log.Println("🔍 Checking if initial health checks are needed...")
					healthChecker.CheckAllServicesIfNeeded()

					// Run periodic checks
					for range ticker.C {
						log.Println("⏰ Running scheduled health checks...")
						healthChecker.CheckAllServices()
					}
				}()
			}
		}()
	} else {
		log.Println("⚠️  Status monitoring disabled (missing environment variables)")
	}

	// Create a handler that checks router first (for HTML pages), then falls back to static files
	cacheFS := web.NewCacheFileServer("public/")
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// For static asset paths, serve directly
		if strings.HasPrefix(path, "/css/") || strings.HasPrefix(path, "/js/") || 
		   strings.HasPrefix(path, "/font/") || strings.HasPrefix(path, "/images/") ||
		   strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js") ||
		   strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") ||
		   strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".ico") ||
		   strings.HasSuffix(path, ".woff") || strings.HasSuffix(path, ".woff2") ||
		   strings.HasSuffix(path, ".ttf") || strings.HasSuffix(path, ".webp") {
			cacheFS.ServeHTTP(w, r)
			return
		}
		
		// For everything else, try router first (handles HTML pages and content)
		// The router will return quickly if it finds a cached page
		routerWriter := &responseWriter{ResponseWriter: w, statusCode: 200}
		router.ServeHTTP(routerWriter, r)
		
		// If router returned 404, try static files as fallback
		if routerWriter.statusCode == 404 {
			// Reset the response
			w.Header().Del("Content-Type")
			
			// Check if static file exists
			fullPath := filepath.Join("public", path)
			if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
				// It's a file, serve it
				cacheFS.ServeHTTP(w, r)
			}
		}
	})

	http.Handle("/", mainHandler)

	// Get port from environment variable, default to 8080 for local development
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	totalStartupTime := time.Since(startTime)
	fmt.Printf("🚀 Blue Website server ready on :%s (total startup: %v)\n", port, totalStartupTime)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
