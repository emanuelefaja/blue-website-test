package web

import (
	"net/http"
	"strings"
)

// SupportedLanguages is the ONLY place to add new languages
// Just add a new language code here and everything else works automatically
var SupportedLanguages = []string{
	"en", // English
	"zh", // 简体中文 (Simplified Chinese)
	"es", // Español (Spanish)
	"fr", // Français (French)
	"de", // Deutsch (German)
	"ja", // 日本語 (Japanese)
	"pt", // Português (Portuguese)
	"ru", // Русский (Russian)
	"ko", // 한국어 (Korean)
	"it", // Italiano (Italian)
	"id", // Indonesian
	"nl", // Nederlands (Dutch)
	"pl", // Polski (Polish)
	"zh-TW", // 繁體中文 (Traditional Chinese)
	"sv", // Svenska (Swedish)
	"km", // ភាសាខ្មែរ (Khmer)
}

// DefaultLanguage is the fallback language
const DefaultLanguage = "en"

// detectPreferredLanguage detects user's preferred language from cookie or Accept-Language header
func detectPreferredLanguage(r *http.Request) string {
	// 1. Check cookie first (user's explicit choice)
	if cookie, err := r.Cookie("lang"); err == nil {
		// Validate the cookie value is a supported language
		for _, lang := range SupportedLanguages {
			if cookie.Value == lang {
				return lang
			}
		}
	}
	
	// 2. Check Accept-Language header
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		// Parse Accept-Language header (e.g., "es-ES,es;q=0.9,en;q=0.8")
		// Simple approach: check if any supported language appears in the header
		acceptLower := strings.ToLower(acceptLang)
		for _, lang := range SupportedLanguages {
			if strings.Contains(acceptLower, lang) {
				return lang
			}
		}
	}
	
	// 3. Default to English
	return DefaultLanguage
}

// isLanguagePrefix checks if a path segment is a valid language code
func isLanguagePrefix(segment string) bool {
	for _, lang := range SupportedLanguages {
		if segment == lang {
			return true
		}
	}
	return false
}

// extractLanguageFromPath extracts language code and clean path from URL
// Examples:
// "/en/about" -> "en", "/about"
// "/es/pricing" -> "es", "/pricing"
// "/about" -> "", "/about"
func extractLanguageFromPath(path string) (string, string) {
	// Remove leading slash for easier parsing
	trimmed := strings.TrimPrefix(path, "/")
	
	// Split path into segments
	segments := strings.SplitN(trimmed, "/", 2)
	
	// Check if first segment is a language code
	if len(segments) > 0 && isLanguagePrefix(segments[0]) {
		lang := segments[0]
		remainingPath := "/"
		if len(segments) > 1 {
			remainingPath = "/" + segments[1]
		}
		return lang, remainingPath
	}
	
	// No language prefix found
	return "", path
}

// setLanguageCookie sets the language preference cookie
func setLanguageCookie(w http.ResponseWriter, lang string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "lang",
		Value:    lang,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}