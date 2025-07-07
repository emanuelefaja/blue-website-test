package web

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter with specified limit and window
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing requests for this IP
	requests, exists := rl.requests[ip]
	if !exists {
		requests = []time.Time{}
	}

	// Filter out old requests
	validRequests := []time.Time{}
	for _, req := range requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if we're over the limit
	if len(validRequests) >= rl.limit {
		rl.requests[ip] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	return true
}

// Cleanup removes old entries to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.window)
	for ip, requests := range rl.requests {
		validRequests := []time.Time{}
		for _, req := range requests {
			if req.After(cutoff) {
				validRequests = append(validRequests, req)
			}
		}
		if len(validRequests) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validRequests
		}
	}
}

// Global rate limiter instance
var apiRateLimiter = NewRateLimiter(60, time.Minute)

// Start cleanup goroutine
func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				apiRateLimiter.Cleanup()
			}
		}
	}()
}

// RateLimitMiddleware returns a middleware that enforces rate limiting
func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		
		if !apiRateLimiter.Allow(ip) {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			
			// Add rate limit headers
			w.Header().Set("X-RateLimit-Limit", "60")
			w.Header().Set("X-RateLimit-Window", "60")
			w.Header().Set("Retry-After", "60")
			
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}
		
		next(w, r)
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// CORS headers - restrictive for API security
		w.Header().Set("Access-Control-Allow-Origin", "https://blue.cc")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Cache headers for API responses
		w.Header().Set("Cache-Control", "public, max-age=60")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

// LoggingMiddleware logs API requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := getClientIP(r)
		
		log.Printf("API Request: %s %s from %s", r.Method, r.URL.Path, ip)
		
		next(w, r)
		
		duration := time.Since(start)
		log.Printf("API Response: %s %s completed in %v", r.Method, r.URL.Path, duration)
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (common with proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	
	// Check for X-Real-IP header (common with reverse proxies)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fallback to RemoteAddr, but strip the port
	remoteAddr := r.RemoteAddr
	if lastColon := strings.LastIndex(remoteAddr, ":"); lastColon != -1 {
		return remoteAddr[:lastColon]
	}
	return remoteAddr
}