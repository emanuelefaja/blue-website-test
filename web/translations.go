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
	translations map[string]map[string]interface{} // [language][namespaced_key]value
	mu           sync.RWMutex
}

// Global translation manager instance
var translationManager *TranslationManager

// InitTranslations loads all translation files at startup
func InitTranslations() error {
	translationManager = &TranslationManager{
		translations: make(map[string]map[string]interface{}),
	}
	
	// Initialize language maps
	for _, lang := range SupportedLanguages {
		translationManager.translations[lang] = make(map[string]interface{})
	}
	
	// Find all translation files in the translations directory
	translationFiles, err := filepath.Glob("translations/*.json")
	if err != nil {
		return fmt.Errorf("failed to glob translation files: %w", err)
	}
	
	// Load each translation file
	for _, file := range translationFiles {
		if err := translationManager.loadFeatureFile(file); err != nil {
			return fmt.Errorf("failed to load translation file %s: %w", file, err)
		}
	}
	
	return nil
}

// loadFeatureFile loads a feature-based translation file (e.g., search.json, about.json)
func (tm *TranslationManager) loadFeatureFile(filePath string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	// Extract feature namespace from filename (e.g., "search" from "search.json")
	filename := filepath.Base(filePath)
	namespace := strings.TrimSuffix(filename, ".json")
	
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	// Parse JSON - expected format: { "en": {...}, "es": {...}, ... }
	var featureTranslations map[string]map[string]interface{}
	if err := json.Unmarshal(data, &featureTranslations); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", filePath, err)
	}
	
	// Merge into main translations with namespace prefix
	for lang, translations := range featureTranslations {
		if _, exists := tm.translations[lang]; !exists {
			tm.translations[lang] = make(map[string]interface{})
		}
		
		// Add namespace prefix to all keys and merge
		tm.addNamespacedTranslations(tm.translations[lang], namespace, translations)
	}
	
	return nil
}

// addNamespacedTranslations recursively adds translations with namespace prefix
func (tm *TranslationManager) addNamespacedTranslations(target map[string]interface{}, namespace string, source map[string]interface{}) {
	for key, value := range source {
		namespacedKey := namespace + "." + key
		
		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively flatten nested objects
			tm.addNamespacedTranslations(target, namespacedKey, v)
		default:
			target[namespacedKey] = value
		}
	}
}

// Translate retrieves a translation for the given language and key
func Translate(lang, key string, args ...interface{}) string {
	if translationManager == nil {
		return key // Return key if translations not initialized
	}
	
	translationManager.mu.RLock()
	defer translationManager.mu.RUnlock()
	
	// Get translation for requested language (key is already namespaced)
	value := translationManager.getValue(lang, key)
	
	// Fallback to default language if not found
	if value == "" && lang != DefaultLanguage {
		value = translationManager.getValue(DefaultLanguage, key)
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

// getValue retrieves a value from the flat translation map
func (tm *TranslationManager) getValue(lang, key string) string {
	langTranslations, exists := tm.translations[lang]
	if !exists {
		return ""
	}
	
	// Direct lookup since we store keys in flattened format
	if value, exists := langTranslations[key]; exists {
		// Convert value to string
		switch v := value.(type) {
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
	
	return ""
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