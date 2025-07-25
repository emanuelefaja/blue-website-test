package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type TranslationJob struct {
	Language string
	Keys     map[string]string
	Context  map[string]string
	KeyOrder []string // Track order for numbered response conversion
}

type TranslationResult struct {
	Language string
	Success  bool
	Data     map[string]interface{}
	Error    error
	KeyCount int
	KeyOrder []string
}

type TranslationManager struct {
	totalLanguages int
	completed      int
	successful     []string
	failed         []string
	mutex          sync.Mutex
	fileMutex      sync.Mutex
}

func main() {
	fmt.Println("🌐 Blue Translation Tool")
	
	// Detect changes since last translation
	changedKeys, err := detectChanges()
	if err != nil {
		fmt.Printf("❌ Error detecting changes: %v\n", err)
		return
	}
	
	if len(changedKeys) == 0 {
		fmt.Println("✅ No changes detected since last translation")
		return
	}
	
	// Get target languages
	targetLanguages := getTargetLanguages()
	if len(targetLanguages) == 0 {
		fmt.Println("❌ No target languages found. Create language files (e.g. touch translations/es.json)")
		return
	}
	
	// Extract context keys
	contextMap := extractContextKeys(changedKeys)
	
	fmt.Printf("📝 Detected changes: %d keys\n", len(changedKeys))
	fmt.Printf("🎯 Target languages: %s (%d languages)\n", strings.Join(targetLanguages, ", "), len(targetLanguages))
	fmt.Printf("⚡ Worker pool: %d goroutines\n\n", min(10, len(targetLanguages)))
	
	// Show interactive menu
	choice := showMenu(changedKeys, targetLanguages)
	
	switch choice {
	case 1:
		previewChanges(changedKeys, contextMap)
	case 2:
		translateAllLanguages(changedKeys, contextMap, targetLanguages)
	case 3:
		translateSpecificLanguage(changedKeys, contextMap, targetLanguages)
	case 4:
		updateTag()
	case 5:
		fmt.Println("👋 Goodbye!")
		return
	}
}

func showMenu(changedKeys map[string]string, targetLanguages []string) int {
	fmt.Println("Options:")
	fmt.Println("[1] Preview changes (show what will be translated)")
	fmt.Println("[2] Translate all languages")
	fmt.Println("[3] Translate specific language")
	fmt.Println("[4] Update last-translation tag")
	fmt.Println("[5] Exit")
	fmt.Print("\nYour choice: ")
	
	var choice int
	fmt.Scanln(&choice)
	return choice
}

func detectChanges() (map[string]string, error) {
	// Check if last-translation tag exists
	cmd := exec.Command("git", "tag", "-l", "last-translation")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git command failed: %v", err)
	}
	
	var changedKeys map[string]string
	
	if strings.TrimSpace(string(output)) == "" {
		// No tag exists, translate everything
		fmt.Println("🏷️  No 'last-translation' tag found - will translate all keys")
		changedKeys, err = loadAllKeys()
	} else {
		// Compare with last translation
		fmt.Println("🔍 Comparing with last-translation tag...")
		changedKeys, err = getChangedKeysSinceTag()
	}
	
	return changedKeys, err
}

func loadAllKeys() (map[string]string, error) {
	data, err := os.ReadFile("translations/en.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read en.json: %v", err)
	}
	
	var enJson map[string]interface{}
	if err := json.Unmarshal(data, &enJson); err != nil {
		return nil, fmt.Errorf("invalid JSON in en.json: %v", err)
	}
	
	return flattenJSON(enJson, ""), nil
}

func getChangedKeysSinceTag() (map[string]string, error) {
	// Get en.json from last-translation tag
	cmd := exec.Command("git", "show", "last-translation:translations/en.json")
	oldData, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get old en.json: %v", err)
	}
	
	// Get current en.json
	newData, err := os.ReadFile("translations/en.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read current en.json: %v", err)
	}
	
	// Parse both versions
	var oldJson, newJson map[string]interface{}
	json.Unmarshal(oldData, &oldJson)
	json.Unmarshal(newData, &newJson)
	
	// Flatten and compare
	oldFlat := flattenJSON(oldJson, "")
	newFlat := flattenJSON(newJson, "")
	
	// Find changes
	changed := make(map[string]string)
	for key, newVal := range newFlat {
		if oldVal, exists := oldFlat[key]; !exists || oldVal != newVal {
			changed[key] = newVal
		}
	}
	
	return changed, nil
}

func flattenJSON(obj map[string]interface{}, prefix string) map[string]string {
	result := make(map[string]string)
	
	for key, value := range obj {
		// Skip context keys
		if strings.HasSuffix(key, "_context") {
			continue
		}
		
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			nested := flattenJSON(v, fullKey)
			for k, val := range nested {
				result[k] = val
			}
		}
	}
	
	return result
}

func extractContextKeys(changedKeys map[string]string) map[string]string {
	data, err := os.ReadFile("translations/en.json")
	if err != nil {
		return make(map[string]string)
	}
	
	var enJson map[string]interface{}
	json.Unmarshal(data, &enJson)
	
	contextMap := make(map[string]string)
	extractContextFromJSON(enJson, "", contextMap)
	
	return contextMap
}

func extractContextFromJSON(obj map[string]interface{}, prefix string, contextMap map[string]string) {
	for key, value := range obj {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		if strings.HasSuffix(key, "_context") {
			// This is a context key for the sibling
			siblingKey := strings.TrimSuffix(fullKey, "_context")
			if str, ok := value.(string); ok {
				contextMap[siblingKey] = str
			}
		} else if nested, ok := value.(map[string]interface{}); ok {
			extractContextFromJSON(nested, fullKey, contextMap)
		}
	}
}

func getTargetLanguages() []string {
	files, err := filepath.Glob("translations/*.json")
	if err != nil {
		return []string{}
	}
	
	var languages []string
	for _, file := range files {
		basename := filepath.Base(file)
		lang := strings.TrimSuffix(basename, ".json")
		
		// Skip source language and files with underscore prefix
		if lang != "en" && !strings.HasPrefix(lang, "_") {
			languages = append(languages, lang)
		}
	}
	
	sort.Strings(languages)
	return languages
}

func previewChanges(changedKeys map[string]string, contextMap map[string]string) {
	fmt.Println("\n📋 Preview of changes to be translated:")
	fmt.Println(strings.Repeat("=", 60))
	
	// Sort keys for consistent output
	var sortedKeys []string
	for key := range changedKeys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	
	for i, key := range sortedKeys {
		text := changedKeys[key]
		context := contextMap[key]
		
		fmt.Printf("\n%d. Key: %s\n", i+1, key)
		fmt.Printf("   Text: %q\n", text)
		if context != "" {
			fmt.Printf("   Context: %s\n", context)
		} else {
			fmt.Printf("   Context: (inferred from key structure)\n")
		}
	}
	
	fmt.Printf("\n📊 Total: %d keys to translate\n", len(changedKeys))
	
	// Ask to continue
	fmt.Print("\nContinue with translation? (y/n): ")
	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) == "y" {
		targetLanguages := getTargetLanguages()
		translateAllLanguages(changedKeys, contextMap, targetLanguages)
	}
}

func translateAllLanguages(changedKeys map[string]string, contextMap map[string]string, targetLanguages []string) {
	fmt.Printf("\n🚀 Starting translation with %d workers for %d languages...\n\n", 
		min(10, len(targetLanguages)), len(targetLanguages))
	
	// Initialize target language files with empty structure matching en.json
	for _, lang := range targetLanguages {
		initializeEmptyTranslationFile(lang)
	}
	
	// Create worker pool
	maxWorkers := min(10, len(targetLanguages))
	jobs := make(chan TranslationJob, len(targetLanguages))
	results := make(chan TranslationResult, len(targetLanguages))
	
	manager := &TranslationManager{
		totalLanguages: len(targetLanguages),
		completed:      0,
		successful:     []string{},
		failed:         []string{},
	}
	
	// Start workers
	for i := 0; i < maxWorkers; i++ {
		go translationWorker(jobs, results, manager)
	}
	
	// Create key order for consistent numbering
	var keyOrder []string
	for key := range changedKeys {
		keyOrder = append(keyOrder, key)
	}
	sort.Strings(keyOrder)
	
	// Send jobs
	for _, lang := range targetLanguages {
		jobs <- TranslationJob{
			Language: lang,
			Keys:     changedKeys,
			Context:  contextMap,
			KeyOrder: keyOrder,
		}
	}
	close(jobs)
	
	// Collect results
	for i := 0; i < len(targetLanguages); i++ {
		result := <-results
		manager.handleResult(result)
	}
	
	// Show final summary
	manager.showSummary()
	
	// Reorder successful language files to match source structure
	if len(manager.successful) > 0 {
		fmt.Print("\n🔄 Reordering translation files to match source structure...\n")
		for _, lang := range manager.successful {
			reorderTranslationFile(lang)
		}
		fmt.Println("✅ Translation files reordered")
	}
	
	// Ask to update tag
	if len(manager.successful) > 0 {
		fmt.Print("\nUpdate 'last-translation' tag? (y/n): ")
		var response string
		fmt.Scanln(&response)
		
		if strings.ToLower(response) == "y" {
			updateTag()
		}
	}
}

func translateSpecificLanguage(changedKeys map[string]string, contextMap map[string]string, targetLanguages []string) {
	fmt.Println("\nAvailable languages:")
	for i, lang := range targetLanguages {
		fmt.Printf("[%d] %s\n", i+1, lang)
	}
	
	fmt.Print("\nChoose language number: ")
	var choice int
	fmt.Scanln(&choice)
	
	if choice < 1 || choice > len(targetLanguages) {
		fmt.Println("❌ Invalid choice")
		return
	}
	
	selectedLang := targetLanguages[choice-1]
	fmt.Printf("\n🚀 Translating to %s...\n", selectedLang)
	
	// Create single job
	manager := &TranslationManager{
		totalLanguages: 1,
		completed:      0,
	}
	
	jobs := make(chan TranslationJob, 1)
	results := make(chan TranslationResult, 1)
	
	go translationWorker(jobs, results, manager)
	
	// Create key order for single language too
	var keyOrder []string
	for key := range changedKeys {
		keyOrder = append(keyOrder, key)
	}
	sort.Strings(keyOrder)
	
	jobs <- TranslationJob{
		Language: selectedLang,
		Keys:     changedKeys,
		Context:  contextMap,
		KeyOrder: keyOrder,
	}
	close(jobs)
	
	result := <-results
	manager.handleResult(result)
	manager.showSummary()
}

func translationWorker(jobs <-chan TranslationJob, results chan<- TranslationResult, manager *TranslationManager) {
	for job := range jobs {
		result := TranslationResult{
			Language: job.Language,
			KeyCount: len(job.Keys),
			KeyOrder: job.KeyOrder,
		}
		
		// Process keys in batches of 5
		allTranslations := make(map[string]interface{})
		success := true
		var lastError error
		
		// Split keys into batches of 5
		batchSize := 5
		keyBatches := make([][]string, 0)
		
		for i := 0; i < len(job.KeyOrder); i += batchSize {
			end := i + batchSize
			if end > len(job.KeyOrder) {
				end = len(job.KeyOrder)
			}
			keyBatches = append(keyBatches, job.KeyOrder[i:end])
		}
		
		fmt.Printf("🔄 %s: Processing %d batches of ~%d keys each...\n", job.Language, len(keyBatches), batchSize)
		
		// Process each batch
		for batchIndex, keyBatch := range keyBatches {
			batchKeys := make(map[string]string)
			batchContext := make(map[string]string)
			
			// Extract keys for this batch
			for _, key := range keyBatch {
				batchKeys[key] = job.Keys[key]
				if context, exists := job.Context[key]; exists {
					batchContext[key] = context
				}
			}
			
			// Try translation with retries for this batch
			batchSuccess := false
			for attempt := 1; attempt <= 5; attempt++ {
				data, err := callClaudeTranslate(batchKeys, batchContext, job.Language, keyBatch)
				if err == nil {
					// Convert numbered response to proper structure for this batch
					batchTranslations := convertNumberedResponse(data, keyBatch)
					
					// Immediately write this batch to the file
					manager.writeBatchToFile(job.Language, batchTranslations)
					
					// Also keep in memory for final result
					for k, v := range data {
						allTranslations[k] = v
					}
					batchSuccess = true
					fmt.Printf("   ✅ Batch %d/%d completed and saved (%d keys)\n", batchIndex+1, len(keyBatches), len(keyBatch))
					break
				}
				
				lastError = err
				
				if attempt < 5 {
					backoffDuration := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
					fmt.Printf("   ❌ Batch %d/%d retry %d/5 (%v) - waiting %v...\n", 
						batchIndex+1, len(keyBatches), attempt, err, backoffDuration)
					time.Sleep(backoffDuration)
				}
			}
			
			if !batchSuccess {
				success = false
				fmt.Printf("   ❌ Batch %d/%d failed after 5 retries\n", batchIndex+1, len(keyBatches))
				break
			}
		}
		
		result.Success = success
		result.Data = allTranslations
		result.Error = lastError
		
		results <- result
	}
}

func callClaudeTranslate(keys map[string]string, contextMap map[string]string, targetLang string, keyOrder []string) (map[string]interface{}, error) {
	// Build the prompt
	prompt := buildTranslationPrompt(keys, contextMap, targetLang, keyOrder)
	
	// Call Claude Code CLI
	cmd := exec.Command("claude", "-p", prompt)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("claude command failed: %v", err)
	}
	
	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response from Claude: %v", err)
	}
	
	return result, nil
}

func buildTranslationPrompt(keys map[string]string, contextMap map[string]string, targetLang string, keyOrder []string) string {
	// Language names for the prompt
	langNames := map[string]string{
		"es": "Spanish",
		"fr": "French", 
		"de": "German",
		"pt": "Portuguese",
		"it": "Italian",
		"ja": "Japanese",
		"ko": "Korean",
		"zh": "Chinese",
		"ru": "Russian",
		"ar": "Arabic",
		"hi": "Hindi",
		"nl": "Dutch",
		"sv": "Swedish",
		"pl": "Polish",
		"cs": "Czech",
		"da": "Danish",
	}
	
	langName := langNames[targetLang]
	if langName == "" {
		langName = targetLang
	}
	
	prompt := fmt.Sprintf(`You are translating UI text for Blue, a B2B SaaS process management platform. Blue helps teams manage workflows, tasks, and business processes efficiently.

CRITICAL RULES:
1. The word "Blue" is a product name and must NEVER be translated
2. Keep all placeholders like {name}, {{count}}, %%s unchanged
3. NEVER translate email addresses (keep support@blue.cc, sales@blue.cc, etc. as-is)
4. Maintain the same tone and formality as the English text
5. Return ONLY raw JSON - no markdown code blocks, no explanations, no backticks

Target Language: %s

Translations needed:

`, langName)
	
	// Use provided key order for consistent numbering
	for i, key := range keyOrder {
		text := keys[key]
		context := contextMap[key]
		
		prompt += fmt.Sprintf("%d. Text: %q\n", i+1, text)
		if context != "" {
			prompt += fmt.Sprintf("   Context: %s\n", context)
		} else {
			prompt += fmt.Sprintf("   Context: %s\n", inferContextFromKey(key))
		}
		prompt += "\n"
	}
	
	prompt += "Please respond with JSON format:\n{\n"
	for i := range keyOrder {
		prompt += fmt.Sprintf("  \"%d\": \"translation here\"", i+1)
		if i < len(keyOrder)-1 {
			prompt += ","
		}
		prompt += "\n"
	}
	prompt += "}"
	
	return prompt
}

func inferContextFromKey(key string) string {
	parts := strings.Split(key, ".")
	
	// Infer context from key structure
	if len(parts) >= 2 {
		section := parts[0]
		field := parts[len(parts)-1]
		
		contexts := map[string]string{
			"title":       "Page or section title",
			"subtitle":    "Descriptive subtitle text", 
			"description": "Longer explanatory text",
			"button":      "Button text - keep short and action-oriented",
			"cta":         "Call-to-action button",
			"error":       "Error message - be helpful, not technical",
			"success":     "Success message",
			"label":       "Form field label",
			"placeholder": "Form field placeholder text",
		}
		
		if context, exists := contexts[field]; exists {
			return fmt.Sprintf("%s in %s section", context, section)
		}
		
		return fmt.Sprintf("Text in %s section", section)
	}
	
	return "General UI text"
}

func (tm *TranslationManager) handleResult(result TranslationResult) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	tm.completed++
	
	if result.Success {
		// Update language file
		tm.updateLanguageFile(result.Language, result.Data, result.KeyOrder)
		tm.successful = append(tm.successful, result.Language)
		fmt.Printf("Progress: %d/%d ✅ %s - %d keys updated\n", 
			tm.completed, tm.totalLanguages, result.Language, result.KeyCount)
	} else {
		tm.failed = append(tm.failed, result.Language)
		fmt.Printf("Progress: %d/%d ❌ %s - failed after 5 retries: %v\n", 
			tm.completed, tm.totalLanguages, result.Language, result.Error)
	}
}

func (tm *TranslationManager) showRetry(lang string, attempt, maxAttempts int, err error, backoff time.Duration) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	fmt.Printf("Progress: %d/%d ❌ %s - retry %d/%d (%v) - waiting %v...\n", 
		tm.completed, tm.totalLanguages, lang, attempt, maxAttempts, err, backoff)
}

func (tm *TranslationManager) writeBatchToFile(lang string, batchTranslations map[string]interface{}) {
	tm.fileMutex.Lock()
	defer tm.fileMutex.Unlock()
	
	// Load existing language file
	filename := fmt.Sprintf("translations/%s.json", lang)
	existing := make(map[string]interface{})
	
	if data, err := os.ReadFile(filename); err == nil {
		json.Unmarshal(data, &existing)
	}
	
	// Ensure existing is not nil
	if existing == nil {
		existing = make(map[string]interface{})
	}
	
	// Deep merge batch translations with existing
	deepMergeMap(existing, batchTranslations)
	
	// Write back to file immediately
	data, _ := json.MarshalIndent(existing, "", "  ")
	os.WriteFile(filename, data, 0644)
}

func (tm *TranslationManager) updateLanguageFile(lang string, translations map[string]interface{}, keyOrder []string) {
	tm.fileMutex.Lock()
	defer tm.fileMutex.Unlock()
	
	// Load existing language file
	filename := fmt.Sprintf("translations/%s.json", lang)
	var existing map[string]interface{}
	
	if data, err := os.ReadFile(filename); err == nil {
		json.Unmarshal(data, &existing)
	} else {
		existing = make(map[string]interface{})
	}
	
	// Convert numbered response to proper structure
	updatedKeys := convertNumberedResponse(translations, keyOrder)
	
	// Merge with existing translations
	merged := mergeTranslations(existing, updatedKeys)
	
	// Write back to file
	data, _ := json.MarshalIndent(merged, "", "  ")
	os.WriteFile(filename, data, 0644)
}

func convertNumberedResponse(numbered map[string]interface{}, keyOrder []string) map[string]interface{} {
	// Convert {"1": "translation", "2": "translation"} back to proper nested structure
	// using the keyOrder to map numbered responses to actual keys
	
	result := make(map[string]interface{})
	
	for i, key := range keyOrder {
		numberKey := fmt.Sprintf("%d", i+1)
		if translation, exists := numbered[numberKey]; exists {
			if translationStr, ok := translation.(string); ok {
				setNestedValue(result, key, translationStr)
			}
		}
	}
	
	return result
}

// setNestedValue sets a value in a nested map structure using dot notation
// e.g., setNestedValue(map, "about.ceo_message.title", "value")
func setNestedValue(target map[string]interface{}, key string, value string) {
	parts := strings.Split(key, ".")
	current := target
	
	// Navigate/create nested structure
	for _, part := range parts[:len(parts)-1] {
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}
		
		if nested, ok := current[part].(map[string]interface{}); ok {
			current = nested
		} else {
			// If the key already exists but isn't a map, we need to create a new structure
			current[part] = make(map[string]interface{})
			current = current[part].(map[string]interface{})
		}
	}
	
	// Set the final value
	finalKey := parts[len(parts)-1]
	current[finalKey] = value
}

func mergeTranslations(existing, updates map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Deep copy existing structure
	result = deepCopyMap(existing)
	
	// Deep merge updates
	deepMergeMap(result, updates)
	
	return result
}

// deepCopyMap creates a deep copy of a map[string]interface{}
func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range original {
		switch v := value.(type) {
		case map[string]interface{}:
			result[key] = deepCopyMap(v)
		default:
			result[key] = v
		}
	}
	
	return result
}

// deepMergeMap merges updates into target, handling nested maps recursively
func deepMergeMap(target, updates map[string]interface{}) {
	if target == nil || updates == nil {
		return
	}
	
	for key, updateValue := range updates {
		if existingValue, exists := target[key]; exists {
			// If both values are maps, merge recursively
			if existingMap, ok := existingValue.(map[string]interface{}); ok {
				if updateMap, ok := updateValue.(map[string]interface{}); ok {
					deepMergeMap(existingMap, updateMap)
					continue
				}
			}
		}
		// Otherwise, overwrite or add the new value
		target[key] = updateValue
	}
}

func (tm *TranslationManager) showSummary() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 Translation Summary:")
	
	for _, lang := range tm.successful {
		fmt.Printf("✅ %s: translations/%s.json updated\n", lang, lang)
	}
	
	for _, lang := range tm.failed {
		fmt.Printf("❌ %s: translation failed\n", lang)
	}
	
	fmt.Printf("\n📝 %d/%d languages completed\n", len(tm.successful), tm.totalLanguages)
	
	if len(tm.failed) > 0 {
		fmt.Println("\nNext run will automatically retry failed languages.")
	}
}

func updateTag() {
	cmd := exec.Command("git", "tag", "-f", "last-translation")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to update tag: %v\n", err)
		return
	}
	
	fmt.Println("🏷️  Updated 'last-translation' tag to current commit")
}

func reorderTranslationFile(lang string) {
	// Load source file as raw text to preserve exact order
	sourceData, err := os.ReadFile("translations/en.json")
	if err != nil {
		fmt.Printf("❌ Failed to read en.json: %v\n", err)
		return
	}
	
	// Load target translations
	targetFile := fmt.Sprintf("translations/%s.json", lang)
	targetData, err := os.ReadFile(targetFile)
	if err != nil {
		fmt.Printf("❌ Failed to read %s: %v\n", targetFile, err)
		return
	}
	
	var targetTranslations map[string]interface{}
	if err := json.Unmarshal(targetData, &targetTranslations); err != nil {
		fmt.Printf("❌ Failed to parse %s: %v\n", targetFile, err)
		return
	}
	
	// Parse source to get flattened keys in order they appear
	sourceFlat := flattenJSON(parseJSONPreservingOrder(string(sourceData)), "")
	
	// Rebuild the target structure in source order
	rebuilt := make(map[string]interface{})
	targetFlat := flattenJSON(targetTranslations, "")
	
	// First, add keys in the order they appear in source
	for _, key := range getOrderedKeys(sourceFlat) {
		if value, exists := targetFlat[key]; exists {
			setNestedValue(rebuilt, key, value)
		}
	}
	
	// Write back to file
	data, _ := json.MarshalIndent(rebuilt, "", "  ")
	os.WriteFile(targetFile, data, 0644)
}

func parseJSONPreservingOrder(jsonStr string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}

func getOrderedKeys(flatMap map[string]string) []string {
	keys := make([]string, 0, len(flatMap))
	for key := range flatMap {
		keys = append(keys, key)
	}
	sort.Strings(keys) // This gives us a consistent order
	return keys
}

func reorderMapToMatchSource(source, target map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Iterate through source keys in order to preserve structure
	for key, sourceValue := range source {
		if targetValue, exists := target[key]; exists {
			// If both are maps, recursively reorder
			if sourceMap, sourceIsMap := sourceValue.(map[string]interface{}); sourceIsMap {
				if targetMap, targetIsMap := targetValue.(map[string]interface{}); targetIsMap {
					result[key] = reorderMapToMatchSource(sourceMap, targetMap)
				} else {
					result[key] = targetValue
				}
			} else {
				// For non-map values, just copy the target value
				result[key] = targetValue
			}
		}
	}
	
	return result
}

func initializeEmptyTranslationFile(lang string) error {
	// Load en.json as raw text to preserve order
	sourceData, err := os.ReadFile("translations/en.json")
	if err != nil {
		return fmt.Errorf("failed to read en.json: %v", err)
	}
	
	// Replace all string values with empty strings while preserving structure
	// This is a bit hacky but preserves exact JSON order
	emptyJSON := string(sourceData)
	
	// Use regex to replace all string values with empty strings
	// Match: "key": "any value here"
	// Replace with: "key": ""
	re := regexp.MustCompile(`(":\s*)"[^"]*"`)
	emptyJSON = re.ReplaceAllString(emptyJSON, `${1}""`)
	
	// Write to target file
	targetFile := fmt.Sprintf("translations/%s.json", lang)
	return os.WriteFile(targetFile, []byte(emptyJSON), 0644)
}

func createEmptyStructure(source map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Preserve exact key order by iterating through source
	for key, value := range source {
		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively create empty structure for nested maps
			result[key] = createEmptyStructure(v)
		case string:
			// For string values, use empty string as placeholder
			result[key] = ""
		default:
			// For other types, preserve the type but empty
			result[key] = ""
		}
	}
	
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}