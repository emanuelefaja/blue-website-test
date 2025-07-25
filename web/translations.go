package web

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TranslationManager handles loading and caching of translations
type TranslationManager struct {
	translations map[string]map[string]interface{}
	mu           sync.RWMutex
}

// Global translation manager instance
var translationManager *TranslationManager

// InitTranslations loads all translation files at startup
func InitTranslations() error {
	translationManager = &TranslationManager{
		translations: make(map[string]map[string]interface{}),
	}
	
	// Load translations for all supported languages
	for _, lang := range SupportedLanguages {
		if err := translationManager.loadLanguage(lang); err != nil {
			return fmt.Errorf("failed to load language %s: %w", lang, err)
		}
	}
	
	return nil
}

// loadLanguage loads a single language file
func (tm *TranslationManager) loadLanguage(lang string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	// Construct path to translation file
	path := filepath.Join("translations", lang+".json")
	
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist and it's not English, that's okay (not translated yet)
		if os.IsNotExist(err) && lang != DefaultLanguage {
			tm.translations[lang] = make(map[string]interface{})
			return nil
		}
		return err
	}
	
	// Parse JSON
	var trans map[string]interface{}
	if err := json.Unmarshal(data, &trans); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	
	tm.translations[lang] = trans
	return nil
}

// Translate retrieves a translation for the given language and key
func Translate(lang, key string, args ...interface{}) string {
	if translationManager == nil {
		return key // Return key if translations not initialized
	}
	
	translationManager.mu.RLock()
	defer translationManager.mu.RUnlock()
	
	// Get translation for requested language
	value := translationManager.getNestedValue(lang, key)
	
	// Fallback to default language if not found
	if value == "" && lang != DefaultLanguage {
		value = translationManager.getNestedValue(DefaultLanguage, key)
	}
	
	// Return key if still not found
	if value == "" {
		return key
	}
	
	// Handle formatting if args provided
	if len(args) > 0 {
		return fmt.Sprintf(value, args...)
	}
	
	return value
}

// getNestedValue retrieves a nested value from translations using dot notation
func (tm *TranslationManager) getNestedValue(lang, key string) string {
	langTranslations, exists := tm.translations[lang]
	if !exists {
		return ""
	}
	
	// Split key by dots for nested access
	parts := strings.Split(key, ".")
	
	// Navigate through nested structure
	var current interface{} = langTranslations
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
			if current == nil {
				return ""
			}
		default:
			return ""
		}
	}
	
	// Convert final value to string
	switch v := current.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return ""
	}
}

// GetTranslations returns all translations for a language (used for client-side)
func GetTranslations(lang string) map[string]interface{} {
	if translationManager == nil {
		return nil
	}
	
	translationManager.mu.RLock()
	defer translationManager.mu.RUnlock()
	
	if trans, exists := translationManager.translations[lang]; exists {
		return trans
	}
	
	// Fallback to default language
	if trans, exists := translationManager.translations[DefaultLanguage]; exists {
		return trans
	}
	
	return nil
}

// ReloadTranslations reloads all translation files (useful for development)
func ReloadTranslations() error {
	return InitTranslations()
}