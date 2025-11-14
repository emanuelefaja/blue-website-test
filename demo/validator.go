package demo

import (
	"errors"
	"html"
	"regexp"
	"strings"
)

var (
	// ErrBotDetected is returned when a bot is detected via honeypot
	ErrBotDetected = errors.New("bot detected")
	// ErrInvalidEmail is returned when email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrMissingRequired is returned when required fields are missing
	ErrMissingRequired = errors.New("missing required fields")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateRequest validates and sanitizes the demo request
func ValidateRequest(req *DemoRequest) error {
	// Check honeypot field - if filled, it's likely a bot
	if req.URL != "" {
		return ErrBotDetected
	}
	
	// Check required fields
	if req.FullName == "" || req.Email == "" || req.Company == "" ||
		req.JobTitle == "" || req.CompanySize == "" || req.UseCase == "" {
		return ErrMissingRequired
	}
	
	// Validate email format
	if !emailRegex.MatchString(req.Email) {
		return ErrInvalidEmail
	}
	
	// Sanitize all fields
	req.FullName = sanitizeInput(req.FullName, 100)
	req.Email = sanitizeInput(req.Email, 255)
	req.Company = sanitizeInput(req.Company, 100)
	req.JobTitle = sanitizeInput(req.JobTitle, 100)
	req.Phone = sanitizeInput(req.Phone, 50)
	req.CompanySize = sanitizeInput(req.CompanySize, 20)
	req.UseCase = sanitizeInput(req.UseCase, 50)
	req.Message = sanitizeInput(req.Message, 1000)
	
	// Validate company size options
	validSizes := map[string]bool{
		"1-25":      true,
		"25-50":     true,
		"50-250":    true,
		"250-1000":  true,
		"1000-5000": true,
		"5000+":     true,
	}
	if !validSizes[req.CompanySize] {
		return errors.New("invalid company size")
	}
	
	// Validate use case options
	validUseCases := map[string]bool{
		"process-automation": true,
		"project-management": true,
		"client-portal":      true,
		"service-tickets":    true,
		"sales-crm":          true,
		"it-operations":      true,
		"other":              true,
	}
	if !validUseCases[req.UseCase] {
		return errors.New("invalid use case")
	}
	
	return nil
}

// sanitizeInput removes HTML tags and limits length
func sanitizeInput(input string, maxLength int) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Escape HTML entities to prevent XSS
	input = html.EscapeString(input)
	
	// Limit length
	if len(input) > maxLength {
		input = input[:maxLength]
	}
	
	return input
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r string) string {
	// Check X-Forwarded-For header (for proxies/load balancers)
	if r != "" {
		// Take the first IP if there are multiple
		parts := strings.Split(r, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	
	// Default to remote address
	return "unknown"
}