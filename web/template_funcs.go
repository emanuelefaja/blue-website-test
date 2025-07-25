package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"time"
)

// templateFuncs defines template functions used across all templates
var templateFuncs = template.FuncMap{
	"toJSON": func(v any) template.JS {
		data, _ := json.Marshal(v)
		return template.JS(data)
	},
	"dict": dict,
	"slice": slice,
	"html": func(s string) template.HTML {
		return template.HTML(s)
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