package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"
)

// Global variable to store the sprite version
var spriteVersion string

// initSpriteVersion calculates a version based on sprite.svg modification time
func initSpriteVersion() {
	file, err := os.Open("public/icons/sprite.svg")
	if err != nil {
		// If file doesn't exist, use timestamp
		spriteVersion = fmt.Sprintf("%d", time.Now().Unix())
		return
	}
	defer file.Close()

	// Get file modification time
	info, err := file.Stat()
	if err != nil {
		spriteVersion = fmt.Sprintf("%d", time.Now().Unix())
		return
	}

	// Use modification time as version
	spriteVersion = fmt.Sprintf("%d", info.ModTime().Unix())
}

// templateFuncs defines template functions used across all templates
var templateFuncs = template.FuncMap{
	"toJSON": func(v any) template.JS {
		data, _ := json.Marshal(v)
		return template.JS(data)
	},
	"dict": dict,
	"slice": slice,
	"spriteURL": func(iconName string) string {
		// Return sprite URL with cache-busting query parameter
		if spriteVersion == "" {
			initSpriteVersion()
		}
		return fmt.Sprintf("/icons/sprite.svg?v=%s#%s", spriteVersion, iconName)
	},
	"html": func(s string) template.HTML {
		// Sanitize HTML to only allow specific safe tags
		return template.HTML(sanitizeHTML(s))
	},
	"parseJSON": func(s string) (interface{}, error) {
		var data interface{}
		err := json.Unmarshal([]byte(s), &data)
		return data, err
	},
	"jsonEscape": func(s string) string {
		data, _ := json.Marshal(s)
		// Remove the surrounding quotes from the JSON string
		escaped := string(data)
		if len(escaped) >= 2 && escaped[0] == '"' && escaped[len(escaped)-1] == '"' {
			escaped = escaped[1 : len(escaped)-1]
		}
		return escaped
	},
	"buildJSON": func(templateStr string, args ...interface{}) (interface{}, error) {
		// Escape all string arguments for JSON safety
		escapedArgs := make([]interface{}, len(args))
		for i, arg := range args {
			if str, ok := arg.(string); ok {
				data, _ := json.Marshal(str)
				// Remove surrounding quotes
				escaped := string(data)
				if len(escaped) >= 2 && escaped[0] == '"' && escaped[len(escaped)-1] == '"' {
					escaped = escaped[1 : len(escaped)-1]
				}
				escapedArgs[i] = escaped
			} else {
				escapedArgs[i] = arg
			}
		}
		
		// Build the JSON string with escaped arguments
		jsonStr := fmt.Sprintf(templateStr, escapedArgs...)
		
		// Parse and return
		var data interface{}
		err := json.Unmarshal([]byte(jsonStr), &data)
		return data, err
	},
	"safeURL": func(s string) template.URL {
		return template.URL(s)
	},
	"formatDate": func(dateStr string) string {
		// Handle empty or invalid dates
		if dateStr == "" {
			return ""
		}

		// Parse ISO date format (YYYY-MM-DD)
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// If parsing fails, return original string
			return dateStr
		}

		// Format as "Month Day, Year" (e.g., "July 12, 2024")
		return parsedDate.Format("January 2, 2006")
	},
	// Translation function - will be overridden with language-specific version
	"t": func(key string, args ...interface{}) string {
		// Default fallback - just return the key
		return key
	},
	// normalizeCategory converts a category name to its translation key format
	"normalizeCategory": func(category string) string {
		// Convert to lowercase and replace spaces with hyphens
		normalized := strings.ToLower(category)
		normalized = strings.ReplaceAll(normalized, " ", "-")
		// Remove any non-alphanumeric characters except hyphens
		var result strings.Builder
		for _, r := range normalized {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				result.WriteRune(r)
			}
		}
		return result.String()
	},
	// lower converts a string to lowercase
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
}

// getTemplateFuncs returns template functions with language-specific translation function
func getTemplateFuncs(lang string) template.FuncMap {
	// Create a copy of the base template functions
	funcs := make(template.FuncMap)
	for k, v := range templateFuncs {
		funcs[k] = v
	}
	
	// Override the translation function with language-specific version
	funcs["t"] = func(key string, args ...interface{}) string {
		return Translate(lang, key, args...)
	}
	
	// tjs is like t but escapes the result for safe use in JavaScript strings
	funcs["tjs"] = func(key string, args ...interface{}) string {
		result := Translate(lang, key, args...)
		// Escape backslashes first, then single quotes for JavaScript
		result = strings.ReplaceAll(result, "\\", "\\\\")
		result = strings.ReplaceAll(result, "'", "\\'")
		return result
	}
	
	return funcs
}

// dict creates a map from key-value pairs for use in templates
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("dict requires even number of arguments")
	}
	m := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m, nil
}

// slice creates a slice from the given values for use in templates
func slice(values ...interface{}) []interface{} {
	return values
}

// sanitizeHTML removes all HTML tags except for a whitelist of safe tags
func sanitizeHTML(input string) string {
	// First, escape all HTML to prevent any XSS
	escaped := html.EscapeString(input)
	
	// Define allowed tags - only self-closing tags and simple formatting tags
	allowedTags := map[string]bool{
		"br":     true,
		"strong": true,
		"em":     true,
		"b":      true,
		"i":      true,
	}
	
	// Replace allowed tags back to their unescaped form
	// This regex matches escaped HTML tags like &lt;br&gt; or &lt;br/&gt; or &lt;br /&gt;
	tagPattern := regexp.MustCompile(`&lt;(/?)(\w+)(\s*/?)&gt;`)
	
	result := tagPattern.ReplaceAllStringFunc(escaped, func(match string) string {
		// Extract the tag name from the escaped HTML
		parts := tagPattern.FindStringSubmatch(match)
		if len(parts) >= 3 {
			tagName := strings.ToLower(parts[2])
			if allowedTags[tagName] {
				// Convert back to actual HTML tag
				return strings.ReplaceAll(strings.ReplaceAll(match, "&lt;", "<"), "&gt;", ">")
			}
		}
		// Keep the escaped version for non-allowed tags
		return match
	})
	
	return result
}