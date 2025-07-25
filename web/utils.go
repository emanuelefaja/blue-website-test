package web

import (
	"strings"
)

// ContentType represents different content sections
type ContentType struct {
	Name       string
	BaseDir    string
	URLPrefix  string
}

// ContentTypes registry for all content types
var ContentTypes = map[string]ContentType{
	"docs": {
		Name:      "docs",
		BaseDir:   "content/en/docs",
		URLPrefix: "/docs",
	},
	"api": {
		Name:      "api",
		BaseDir:   "content/en/api",
		URLPrefix: "/api",
	},
	"legal": {
		Name:      "legal",
		BaseDir:   "content/en/legal",
		URLPrefix: "/legal",
	},
	"insights": {
		Name:      "insights",
		BaseDir:   "content/en/insights",
		URLPrefix: "/insights",
	},
}

// CleanTitle removes numeric prefixes and cleans up titles
func CleanTitle(name string) string {
	// Remove numeric prefix (e.g., "1.start-guide" -> "start-guide")
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		name = parts[1]
	}

	// Replace hyphens/underscores with spaces and title case
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Simple title case
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// CleanID creates clean IDs for navigation
func CleanID(name string) string {
	// Remove numeric prefix
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		name = parts[1]
	}

	// Convert to lowercase and replace spaces/special chars with hyphens
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	return name
}

// CleanDirectoryPath removes numeric prefixes from directory paths
func CleanDirectoryPath(path string) string {
	// Split path into segments and clean each one
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = CleanID(segment)
	}
	return strings.Join(segments, "/")
}

// ExtractNumericPrefix extracts the numeric prefix from a name (e.g., "1.start-guide" -> 1)
func ExtractNumericPrefix(name string) int {
	// Parse numeric prefix from original directory/file names
	parts := strings.Split(name, ".")
	if len(parts) >= 2 {
		// Try to parse first part as number
		num := 0
		for _, char := range parts[0] {
			if char >= '0' && char <= '9' {
				num = num*10 + int(char-'0')
			} else {
				break
			}
		}
		if num > 0 {
			return num
		}
	}

	// Fallback: assign high number for non-numbered items
	return 9999
}

// GenerateFilePatterns generates variations of a file pattern for matching
func GenerateFilePatterns(segment string, extension string) []string {
	patterns := []string{
		"*" + segment + extension,                               // e.g., *download-apps.md
		"*" + strings.ReplaceAll(segment, "-", " ") + extension, // e.g., *download apps.md
		"*" + strings.ReplaceAll(segment, "-", "_") + extension, // e.g., *download_apps.md
	}
	return patterns
}

// GetContentTypeFromPath determines content type from URL path
func GetContentTypeFromPath(path string) (ContentType, bool) {
	cleanPath := strings.Trim(path, "/")
	parts := strings.Split(cleanPath, "/")
	
	if len(parts) == 0 {
		return ContentType{}, false
	}

	// Check first part of path
	switch parts[0] {
	case "docs":
		return ContentTypes["docs"], true
	case "api":
		return ContentTypes["api"], true
	case "legal":
		return ContentTypes["legal"], true
	case "insights":
		return ContentTypes["insights"], true
	default:
		return ContentType{}, false
	}
}