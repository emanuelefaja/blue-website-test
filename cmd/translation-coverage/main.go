package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirectoryStats struct {
	Name        string
	TotalFiles  int
	PresentFiles int
	MissingFiles []string
	Coverage    float64
}

type LanguageStats struct {
	Language      string
	Emoji         string
	TotalFiles    int
	PresentFiles  int
	MissingFiles  []string
	Coverage      float64
	Directories   map[string]DirectoryStats
}

type TranslationFileStats struct {
	FileName          string
	TotalLanguages    int
	PresentLanguages  []string
	MissingLanguages  []string
	Coverage          float64
}

type PageTranslationMapping struct {
	PagePath          string
	TranslationFile   string
	HasTranslation    bool
}

type ComprehensiveCoverage struct {
	ContentStats      []LanguageStats
	TranslationStats  []TranslationFileStats
	PageMappings      []PageTranslationMapping
}

// SupportedLanguages must match the list in web/languages.go
var SupportedLanguages = []string{
	"en", "zh", "es", "fr", "de", "ja", "pt", "ru", 
	"ko", "it", "id", "nl", "pl", "zh-TW", "sv", "km",
}

var languageEmojis = map[string]string{
	"de":    "🇩🇪",
	"es":    "🇪🇸", 
	"fr":    "🇫🇷",
	"it":    "🇮🇹",
	"ja":    "🇯🇵",
	"ko":    "🇰🇷",
	"zh":    "🇨🇳",
	"zh-TW": "🇹🇼",
	"ru":    "🇷🇺",
	"pl":    "🇵🇱",
	"pt":    "🇵🇹",
	"nl":    "🇳🇱",
	"sv":    "🇸🇪",
	"id":    "🇮🇩",
	"km":    "🇰🇭",
}

var languageNames = map[string]string{
	"de":    "German",
	"es":    "Spanish",
	"fr":    "French", 
	"it":    "Italian",
	"ja":    "Japanese",
	"ko":    "Korean",
	"zh":    "Chinese",
	"zh-TW": "Chinese (Traditional)",
	"ru":    "Russian",
	"pl":    "Polish",
	"pt":    "Portuguese",
	"nl":    "Dutch",
	"sv":    "Swedish",
	"id":    "Indonesian",
	"km":    "Khmer",
}

func main() {
	fmt.Println("📊 Blue Translation Coverage Report")
	fmt.Println("===================================")
	
	// Part 1: Content Coverage (existing)
	fmt.Println("\n🌐 PART 1: CONTENT COVERAGE")
	fmt.Println(strings.Repeat("=", 50))
	
	// Scan baseline English content
	baselineFiles, err := scanDirectory("content/en")
	if err != nil {
		fmt.Printf("❌ Error scanning baseline directory: %v\n", err)
		return
	}
	
	if len(baselineFiles) == 0 {
		fmt.Println("❌ No files found in /content/en directory")
		return
	}
	
	fmt.Printf("\n📝 Baseline: English (/content/en) - %d files\n\n", len(baselineFiles))
	
	// Get all language directories
	languages, err := getLanguageDirectories()
	if err != nil {
		fmt.Printf("❌ Error getting language directories: %v\n", err)
		return
	}
	
	// Analyze each language
	var allStats []LanguageStats
	for _, lang := range languages {
		stats := analyzeLanguage(lang, baselineFiles)
		allStats = append(allStats, stats)
	}
	
	// Sort by coverage (highest first)
	sort.Slice(allStats, func(i, j int) bool {
		return allStats[i].Coverage > allStats[j].Coverage
	})
	
	// Generate content report
	generateReport(allStats, len(baselineFiles))
	
	// Part 2: Translation JSON Coverage
	fmt.Println("\n\n💬 PART 2: TRANSLATION JSON COVERAGE")
	fmt.Println(strings.Repeat("=", 50))
	
	translationStats, err := scanTranslationFiles()
	if err != nil {
		fmt.Printf("❌ Error scanning translation files: %v\n", err)
	} else {
		generateTranslationReport(translationStats)
	}
	
	// Part 3: Page-Translation Mapping
	fmt.Println("\n\n🔗 PART 3: PAGE-TRANSLATION MAPPING")
	fmt.Println(strings.Repeat("=", 50))
	
	mappings := mapPagesToTranslations()
	generateMappingReport(mappings)
	
	// Final Summary
	generateFinalSummary(allStats, translationStats, mappings)
}

func scanDirectory(basePath string) (map[string]bool, error) {
	files := make(map[string]bool)
	
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Only include .md files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// Get relative path from base directory
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}
			files[relPath] = true
		}
		
		return nil
	})
	
	return files, err
}

func getLanguageDirectories() ([]string, error) {
	contentDir := "content"
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}
	
	var languages []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "en" {
			languages = append(languages, entry.Name())
		}
	}
	
	sort.Strings(languages)
	return languages, nil
}

func analyzeLanguage(language string, baselineFiles map[string]bool) LanguageStats {
	langPath := filepath.Join("content", language)
	langFiles, err := scanDirectory(langPath)
	if err != nil {
		// Language directory doesn't exist or is empty
		langFiles = make(map[string]bool)
	}
	
	// Calculate overall stats
	presentCount := 0
	var missingFiles []string
	
	for file := range baselineFiles {
		if langFiles[file] {
			presentCount++
		} else {
			missingFiles = append(missingFiles, file)
		}
	}
	
	coverage := float64(presentCount) / float64(len(baselineFiles)) * 100
	
	// Analyze by directory
	directories := analyzeDirectories(baselineFiles, langFiles)
	
	emoji := languageEmojis[language]
	if emoji == "" {
		emoji = "🌐"
	}
	
	return LanguageStats{
		Language:      language,
		Emoji:         emoji,
		TotalFiles:    len(baselineFiles),
		PresentFiles:  presentCount,
		MissingFiles:  missingFiles,
		Coverage:      coverage,
		Directories:   directories,
	}
}

func analyzeDirectories(baselineFiles, langFiles map[string]bool) map[string]DirectoryStats {
	dirMap := make(map[string]map[string]bool)
	
	// Group files by directory
	for file := range baselineFiles {
		dir := filepath.Dir(file)
		if dir == "." {
			dir = "root"
		}
		
		if dirMap[dir] == nil {
			dirMap[dir] = make(map[string]bool)
		}
		dirMap[dir][file] = true
	}
	
	// Calculate stats for each directory
	directories := make(map[string]DirectoryStats)
	for dir, files := range dirMap {
		presentCount := 0
		var missingFiles []string
		
		for file := range files {
			if langFiles[file] {
				presentCount++
			} else {
				missingFiles = append(missingFiles, filepath.Base(file))
			}
		}
		
		coverage := float64(presentCount) / float64(len(files)) * 100
		
		directories[dir] = DirectoryStats{
			Name:         dir,
			TotalFiles:   len(files),
			PresentFiles: presentCount,
			MissingFiles: missingFiles,
			Coverage:     coverage,
		}
	}
	
	return directories
}

func generateReport(allStats []LanguageStats, totalFiles int) {
	// Generate summary table first
	generateSummaryTable(allStats, totalFiles)
	
	// Then detailed breakdown
	fmt.Println("\n📋 Detailed Breakdown:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	
	for _, stats := range allStats {
		generateLanguageReport(stats)
	}
	
	// Generate final summary
	fmt.Println("📈 Summary:")
	if len(allStats) > 0 {
		fmt.Printf("  🥇 Best: %s (%.1f%%)\n", getLanguageDisplayName(allStats[0].Language), allStats[0].Coverage)
	}
	if len(allStats) > 1 {
		fmt.Printf("  🥈 Second: %s (%.1f%%)\n", getLanguageDisplayName(allStats[1].Language), allStats[1].Coverage)
	}
	if len(allStats) > 2 {
		fmt.Printf("  🥉 Third: %s (%.1f%%)\n", getLanguageDisplayName(allStats[2].Language), allStats[2].Coverage)
	}
	
	// Find languages that need attention
	var needsAttention []string
	for _, stats := range allStats {
		if stats.Coverage < 10 {
			needsAttention = append(needsAttention, getLanguageDisplayName(stats.Language))
		}
	}
	
	if len(needsAttention) > 0 {
		fmt.Printf("  ⚠️  Needs attention: %s (< 10%%)\n", strings.Join(needsAttention, ", "))
	}
	
	fmt.Printf("\n💾 Total baseline files: %d\n", totalFiles)
	fmt.Printf("🌐 Languages analyzed: %d\n", len(allStats))
}

func generateLanguageReport(stats LanguageStats) {
	langName := getLanguageDisplayName(stats.Language)
	
	fmt.Printf("%s %s (%s) - %.1f%% coverage (%d/%d files)\n", 
		stats.Emoji, langName, stats.Language, stats.Coverage, 
		stats.PresentFiles, stats.TotalFiles)
	
	// Sort directories by name for consistent output
	var dirNames []string
	for dirName := range stats.Directories {
		dirNames = append(dirNames, dirName)
	}
	sort.Strings(dirNames)
	
	// Show directory breakdown
	for _, dirName := range dirNames {
		dirStats := stats.Directories[dirName]
		
		if dirStats.Coverage == 100 {
			fmt.Printf("  📁 %s/ - %.0f%% (%d/%d) ✅\n", 
				dirName, dirStats.Coverage, dirStats.PresentFiles, dirStats.TotalFiles)
		} else if dirStats.Coverage == 0 {
			fmt.Printf("  📁 %s/ - %.0f%% (%d/%d) ❌ Missing entire directory\n", 
				dirName, dirStats.Coverage, dirStats.PresentFiles, dirStats.TotalFiles)
		} else {
			fmt.Printf("  📁 %s/ - %.0f%% (%d/%d) ❌ Missing: %s\n", 
				dirName, dirStats.Coverage, dirStats.PresentFiles, dirStats.TotalFiles,
				formatMissingFiles(dirStats.MissingFiles))
		}
	}
	
	fmt.Println()
}

func formatMissingFiles(missingFiles []string) string {
	if len(missingFiles) == 0 {
		return "none"
	}
	
	// Sort for consistent output
	sort.Strings(missingFiles)
	
	if len(missingFiles) <= 3 {
		return strings.Join(missingFiles, ", ")
	}
	
	// Show first 3 and count
	return fmt.Sprintf("%s, and %d more", 
		strings.Join(missingFiles[:3], ", "), len(missingFiles)-3)
}

func generateSummaryTable(allStats []LanguageStats, totalFiles int) {
	fmt.Println("📊 Coverage Summary:")
	fmt.Println(strings.Repeat("=", 65))
	fmt.Printf("%-4s %-20s %-15s %-10s\n", "", "Language", "Files", "Coverage")
	fmt.Println(strings.Repeat("-", 65))
	
	for i, stats := range allStats {
		rankEmoji := ""
		if i == 0 {
			rankEmoji = "🥇"
		} else if i == 1 {
			rankEmoji = "🥈"
		} else if i == 2 {
			rankEmoji = "🥉"
		}
		
		langName := getLanguageDisplayName(stats.Language)
		filesRatio := fmt.Sprintf("%d/%d", stats.PresentFiles, stats.TotalFiles)
		coverage := fmt.Sprintf("%.1f%%", stats.Coverage)
		
		fmt.Printf("%-4s %s %-18s %-15s %s\n", 
			rankEmoji, stats.Emoji, langName, filesRatio, coverage)
	}
	
	fmt.Println(strings.Repeat("=", 65))
}

func getLanguageDisplayName(langCode string) string {
	if name, exists := languageNames[langCode]; exists {
		return name
	}
	return strings.ToUpper(langCode)
}

func scanTranslationFiles() ([]TranslationFileStats, error) {
	translationsDir := "translations"
	entries, err := os.ReadDir(translationsDir)
	if err != nil {
		return nil, err
	}
	
	var stats []TranslationFileStats
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			filePath := filepath.Join(translationsDir, entry.Name())
			fileStats, err := analyzeTranslationFile(filePath)
			if err != nil {
				fmt.Printf("  ⚠️  Error analyzing %s: %v\n", entry.Name(), err)
				continue
			}
			stats = append(stats, fileStats)
		}
	}
	
	// Sort by coverage (highest first)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Coverage > stats[j].Coverage
	})
	
	return stats, nil
}

func analyzeTranslationFile(filePath string) (TranslationFileStats, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TranslationFileStats{}, err
	}
	
	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return TranslationFileStats{}, err
	}
	
	fileName := filepath.Base(filePath)
	stats := TranslationFileStats{
		FileName:       fileName,
		TotalLanguages: len(SupportedLanguages),
	}
	
	// Check which languages are present
	for _, lang := range SupportedLanguages {
		if _, exists := translations[lang]; exists {
			stats.PresentLanguages = append(stats.PresentLanguages, lang)
		} else {
			stats.MissingLanguages = append(stats.MissingLanguages, lang)
		}
	}
	
	stats.Coverage = float64(len(stats.PresentLanguages)) / float64(stats.TotalLanguages) * 100
	
	return stats, nil
}

func mapPagesToTranslations() []PageTranslationMapping {
	var mappings []PageTranslationMapping
	
	// Known utility translation files that don't need corresponding pages
	utilityTranslations := map[string]bool{
		"common.json": true,
		"search.json": true,
		"schema.json": true,
		"chat.json":   true, // Also appears to be a utility file
	}
	
	// Get all HTML pages
	pages := []string{}
	err := filepath.Walk("pages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") && !strings.HasSuffix(info.Name(), ".bak") {
			pages = append(pages, path)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("  ⚠️  Error scanning pages directory: %v\n", err)
		return mappings
	}
	
	// Get all translation files
	translationFiles := make(map[string]bool)
	entries, err := os.ReadDir("translations")
	if err != nil {
		fmt.Printf("  ⚠️  Error reading translations directory: %v\n", err)
		return mappings
	}
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			translationFiles[entry.Name()] = true
		}
	}
	
	// Map pages to translation files
	for _, page := range pages {
		mapping := PageTranslationMapping{
			PagePath: page,
		}
		
		// Determine expected translation file
		translationFile := getExpectedTranslationFile(page)
		mapping.TranslationFile = translationFile
		
		// Check if translation file exists
		if translationFiles[translationFile] {
			mapping.HasTranslation = true
			delete(translationFiles, translationFile) // Remove from map to track orphans
		}
		
		mappings = append(mappings, mapping)
	}
	
	// Add orphaned translation files (no corresponding page, excluding utility files)
	for file := range translationFiles {
		if !utilityTranslations[file] {
			mappings = append(mappings, PageTranslationMapping{
				PagePath:        "",
				TranslationFile: file,
				HasTranslation:  false,
			})
		}
	}
	
	return mappings
}

func getExpectedTranslationFile(pagePath string) string {
	// Remove pages/ prefix and .html suffix
	path := strings.TrimPrefix(pagePath, "pages/")
	path = strings.TrimSuffix(path, ".html")
	
	// Special case mappings
	if path == "index" {
		return "home.json"
	}
	
	// Special case for solutions/index.html -> solutions.json
	if path == "solutions/index" {
		return "solutions.json"
	}
	
	// For nested paths, use the last segment
	segments := strings.Split(path, "/")
	baseName := segments[len(segments)-1]
	
	// Handle index.html in subdirectories
	if baseName == "index" {
		if len(segments) > 1 {
			// Use parent directory name
			baseName = segments[len(segments)-2]
		}
	}
	
	return baseName + ".json"
}

func generateTranslationReport(stats []TranslationFileStats) {
	if len(stats) == 0 {
		fmt.Println("❌ No translation files found")
		return
	}
	
	fmt.Printf("\n📂 Found %d translation files\n", len(stats))
	fmt.Printf("🌍 Expected %d languages per file\n\n", len(SupportedLanguages))
	
	// Summary table
	fmt.Println("📊 Translation File Coverage:")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-30s %-20s %-20s\n", "File", "Languages", "Coverage")
	fmt.Println(strings.Repeat("-", 70))
	
	var fullyCovered, partiallyCovered, empty int
	
	for _, stat := range stats {
		langRatio := fmt.Sprintf("%d/%d", len(stat.PresentLanguages), stat.TotalLanguages)
		coverage := fmt.Sprintf("%.1f%%", stat.Coverage)
		
		statusIcon := ""
		if stat.Coverage == 100 {
			statusIcon = " ✅"
			fullyCovered++
		} else if stat.Coverage == 0 {
			statusIcon = " ❌"
			empty++
		} else {
			statusIcon = " ⚠️"
			partiallyCovered++
		}
		
		fmt.Printf("%-30s %-20s %-20s%s\n", stat.FileName, langRatio, coverage, statusIcon)
	}
	
	fmt.Println(strings.Repeat("=", 70))
	
	// Summary stats
	fmt.Printf("\n📈 Summary:\n")
	fmt.Printf("  ✅ Fully translated: %d files (%.1f%%)\n", 
		fullyCovered, float64(fullyCovered)/float64(len(stats))*100)
	fmt.Printf("  ⚠️  Partially translated: %d files (%.1f%%)\n", 
		partiallyCovered, float64(partiallyCovered)/float64(len(stats))*100)
	fmt.Printf("  ❌ Empty/Missing: %d files (%.1f%%)\n", 
		empty, float64(empty)/float64(len(stats))*100)
	
	// Show files needing attention
	fmt.Println("\n🔍 Files Needing Attention:")
	needsAttention := false
	for _, stat := range stats {
		if stat.Coverage < 100 && stat.Coverage > 0 {
			needsAttention = true
			fmt.Printf("\n  📄 %s (%.1f%% complete)\n", stat.FileName, stat.Coverage)
			fmt.Printf("     Missing: %s\n", formatLanguageList(stat.MissingLanguages))
		}
	}
	
	if !needsAttention {
		fmt.Println("  ✅ All translation files are either complete or empty!")
	}
}

func generateMappingReport(mappings []PageTranslationMapping) {
	if len(mappings) == 0 {
		fmt.Println("❌ No page mappings found")
		return
	}
	
	var pagesWithTranslations, pagesWithoutTranslations, orphanedTranslations int
	var missingTranslations []PageTranslationMapping
	var orphans []PageTranslationMapping
	
	for _, mapping := range mappings {
		if mapping.PagePath == "" {
			// Orphaned translation file
			orphanedTranslations++
			orphans = append(orphans, mapping)
		} else if mapping.HasTranslation {
			pagesWithTranslations++
		} else {
			pagesWithoutTranslations++
			missingTranslations = append(missingTranslations, mapping)
		}
	}
	
	totalPages := pagesWithTranslations + pagesWithoutTranslations
	coverage := float64(pagesWithTranslations) / float64(totalPages) * 100
	
	fmt.Printf("\n📄 Page Coverage:\n")
	fmt.Printf("  Total pages: %d\n", totalPages)
	fmt.Printf("  Pages with translations: %d (%.1f%%)\n", pagesWithTranslations, coverage)
	fmt.Printf("  Pages without translations: %d\n", pagesWithoutTranslations)
	fmt.Printf("  Orphaned translation files: %d\n", orphanedTranslations)
	
	// Show pages missing translations
	if len(missingTranslations) > 0 {
		fmt.Println("\n❌ Pages Missing Translation Files:")
		// Sort by path for better readability
		sort.Slice(missingTranslations, func(i, j int) bool {
			return missingTranslations[i].PagePath < missingTranslations[j].PagePath
		})
		
		for i, mapping := range missingTranslations {
			if i >= 10 {
				fmt.Printf("  ... and %d more\n", len(missingTranslations)-10)
				break
			}
			fmt.Printf("  - %s → %s (missing)\n", mapping.PagePath, mapping.TranslationFile)
		}
	}
	
	// Show orphaned translation files
	if len(orphans) > 0 {
		fmt.Println("\n🗑️  Orphaned Translation Files (no corresponding page):")
		for _, orphan := range orphans {
			fmt.Printf("  - %s\n", orphan.TranslationFile)
		}
	}
}

func generateFinalSummary(contentStats []LanguageStats, translationStats []TranslationFileStats, mappings []PageTranslationMapping) {
	fmt.Println("\n\n📊 FINAL SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	
	// Content coverage summary
	if len(contentStats) > 0 {
		avgContentCoverage := 0.0
		for _, stat := range contentStats {
			avgContentCoverage += stat.Coverage
		}
		avgContentCoverage /= float64(len(contentStats))
		
		fmt.Printf("\n📝 Content Coverage:\n")
		fmt.Printf("  Average coverage across languages: %.1f%%\n", avgContentCoverage)
		fmt.Printf("  Best coverage: %s (%.1f%%)\n", 
			getLanguageDisplayName(contentStats[0].Language), contentStats[0].Coverage)
	}
	
	// Translation file coverage summary
	if len(translationStats) > 0 {
		fullyCovered := 0
		for _, stat := range translationStats {
			if stat.Coverage == 100 {
				fullyCovered++
			}
		}
		fmt.Printf("\n💬 Translation Files:\n")
		fmt.Printf("  Total files: %d\n", len(translationStats))
		fmt.Printf("  Fully translated: %d (%.1f%%)\n", 
			fullyCovered, float64(fullyCovered)/float64(len(translationStats))*100)
	}
	
	// Page mapping summary
	pagesWithTranslations := 0
	totalPages := 0
	for _, mapping := range mappings {
		if mapping.PagePath != "" {
			totalPages++
			if mapping.HasTranslation {
				pagesWithTranslations++
			}
		}
	}
	
	if totalPages > 0 {
		fmt.Printf("\n🔗 Page-Translation Mapping:\n")
		fmt.Printf("  Pages with translation support: %d/%d (%.1f%%)\n",
			pagesWithTranslations, totalPages, 
			float64(pagesWithTranslations)/float64(totalPages)*100)
	}
	
	fmt.Println("\n✨ Report generation complete!")
}

func formatLanguageList(languages []string) string {
	if len(languages) == 0 {
		return "none"
	}
	
	// Add emoji flags to languages
	var formatted []string
	for _, lang := range languages {
		emoji := languageEmojis[lang]
		if emoji == "" {
			emoji = "🌐"
		}
		formatted = append(formatted, fmt.Sprintf("%s %s", emoji, lang))
	}
	
	if len(formatted) <= 4 {
		return strings.Join(formatted, ", ")
	}
	
	// Show first 4 and count
	return fmt.Sprintf("%s, and %d more", 
		strings.Join(formatted[:4], ", "), len(formatted)-4)
}