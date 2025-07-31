package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// TranslationEntry represents a single text to translate
type TranslationEntry struct {
	Key   string   // JSON path like "stats.customers_label" or "page_title"
	Value string   // The text to translate
	Path  []string // JSON path as array for reconstruction
}

// TranslationTask represents a single translation job
type TranslationTask struct {
	Entry      TranslationEntry
	TargetLang string
}

// TranslationResult represents the result of a translation
type TranslationResult struct {
	Task         TranslationTask
	Translation  string
	InputTokens  int
	OutputTokens int
	Error        error
}

// TranslationStats tracks statistics
type TranslationStats struct {
	TotalTasks        int
	CompletedTasks    int32
	FailedTasks       int32
	TotalInputTokens  int64
	TotalOutputTokens int64
	StartTime         time.Time
}

// Language names mapping
var languageNames = map[string]string{
	"en":    "English",
	"es":    "Spanish",
	"fr":    "French",
	"de":    "German",
	"it":    "Italian",
	"pt":    "Portuguese",
	"ja":    "Japanese",
	"ko":    "Korean",
	"zh":    "Chinese (Simplified)",
	"zh-TW": "Chinese (Traditional)",
	"ru":    "Russian",
	"nl":    "Dutch",
	"pl":    "Polish",
	"sv":    "Swedish",
	"id":    "Indonesian",
	"km":    "Khmer",
}

const (
	inputTokenCost  = 0.15 // $0.15 per 1M input tokens for gpt-4o-mini
	outputTokenCost = 0.60 // $0.60 per 1M output tokens for gpt-4o-mini
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY not found in environment variables")
		os.Exit(1)
	}

	// Initialize OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	// Get target languages from web/languages.go
	targetLanguages, err := getTargetLanguages()
	if err != nil {
		fmt.Printf("Error reading target languages: %v\n", err)
		os.Exit(1)
	}

	// Check for help flags
	if len(os.Args) >= 2 && (os.Args[1] == "-help" || os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/translate-json/main.go [json-file]")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  go run cmd/translate-json/main.go about.json    # Translate specific file")
		fmt.Println("  go run cmd/translate-json/main.go              # Scan and translate all files needing work")
		fmt.Println("")
		fmt.Println("When no file is specified, the tool will automatically scan the translations/")
		fmt.Println("directory and process all files that have missing translations.")
		return
	}

	// Check if specific file provided or batch mode
	if len(os.Args) >= 2 {
		// Single file mode (existing behavior)
		jsonFile := os.Args[1]
		if !strings.HasSuffix(jsonFile, ".json") {
			jsonFile += ".json"
		}

		err := processFile(jsonFile, targetLanguages, &client)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", jsonFile, err)
			os.Exit(1)
		}
	} else {
		// Batch mode - scan all files and process those needing translations
		fmt.Println("No specific file provided. Scanning translations directory for files needing translation...")

		filesToProcess, err := findFilesNeedingTranslation(targetLanguages)
		if err != nil {
			fmt.Printf("Error scanning translations directory: %v\n", err)
			os.Exit(1)
		}

		if len(filesToProcess) == 0 {
			fmt.Println("🎉 All translation files are complete! No work needed.")
			return
		}

		fmt.Printf("Found %d files needing translation: %v\n\n", len(filesToProcess), filesToProcess)

		// Process each file
		for i, jsonFile := range filesToProcess {
			fmt.Printf("=== Processing file %d/%d: %s ===\n", i+1, len(filesToProcess), jsonFile)

			err := processFile(jsonFile, targetLanguages, &client)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", jsonFile, err)
				// Continue with next file instead of exiting
				continue
			}

			if i < len(filesToProcess)-1 {
				fmt.Printf("\n✅ Completed %s. Moving to next file...\n\n", jsonFile)
			}
		}

		fmt.Println("\n🎉 Batch processing complete! All files have been processed.")
	}
}

// getTargetLanguages reads supported languages from web/languages.go
func getTargetLanguages() ([]string, error) {
	// Read the languages.go file
	data, err := os.ReadFile("web/languages.go")
	if err != nil {
		return nil, fmt.Errorf("failed to read web/languages.go: %v", err)
	}

	content := string(data)

	// Find the SupportedLanguages slice
	start := strings.Index(content, "var SupportedLanguages = []string{")
	if start == -1 {
		return nil, fmt.Errorf("SupportedLanguages not found in web/languages.go")
	}

	end := strings.Index(content[start:], "}")
	if end == -1 {
		return nil, fmt.Errorf("end of SupportedLanguages not found")
	}

	// Extract the slice content
	sliceContent := content[start : start+end]

	var languages []string
	// Parse each line to extract language codes
	lines := strings.Split(sliceContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"") {
			// Extract language code between quotes
			if endQuote := strings.Index(line[1:], "\""); endQuote != -1 {
				lang := line[1 : endQuote+1]
				if lang != "en" { // Exclude English since it's the source
					languages = append(languages, lang)
				}
			}
		}
	}

	return languages, nil
}

func loadJSON(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var jsonObj map[string]interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, err
	}

	return jsonObj, nil
}

func extractEnglishEntries(jsonData map[string]interface{}) ([]TranslationEntry, error) {
	var entries []TranslationEntry

	englishData, ok := jsonData["en"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no 'en' section found in JSON")
	}

	// Recursively extract all string values
	extractStrings(englishData, []string{}, &entries)

	return entries, nil
}

func extractStrings(obj interface{}, path []string, entries *[]TranslationEntry) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPath := append(path, key)
			extractStrings(value, newPath, entries)
		}
	case string:
		// Create the dot-separated key
		key := strings.Join(path, ".")
		*entries = append(*entries, TranslationEntry{
			Key:   key,
			Value: v,
			Path:  append([]string{}, path...), // Copy the path
		})
	case []interface{}:
		// Handle arrays - for now, skip them as they're complex to translate
		// Could be extended to handle arrays of strings if needed
	}
}

func buildTranslationTasks(entries []TranslationEntry, languages []string, jsonData map[string]interface{}) []TranslationTask {
	var tasks []TranslationTask

	for _, lang := range languages {
		// Get existing translations for this language
		existingKeys := make(map[string]bool)
		if langData, ok := jsonData[lang].(map[string]interface{}); ok {
			// Recursively get all existing keys
			getAllKeys(langData, []string{}, existingKeys)
		}

		for _, entry := range entries {
			// Check if already translated
			if existingKeys[entry.Key] {
				continue
			}

			tasks = append(tasks, TranslationTask{
				Entry:      entry,
				TargetLang: lang,
			})
		}
	}

	return tasks
}

func getAllKeys(obj interface{}, path []string, keys map[string]bool) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newPath := append(path, key)
			getAllKeys(value, newPath, keys)
		}
	case string:
		key := strings.Join(path, ".")
		keys[key] = true
	}
}

func translationWorker(id int, client *openai.Client, taskChan <-chan TranslationTask, resultChan chan<- TranslationResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		result := translateEntry(client, task)
		resultChan <- result
	}
}

func translateEntry(client *openai.Client, task TranslationTask) TranslationResult {
	ctx := context.Background()

	// Prepare the prompt
	systemPrompt := `You are a professional translator for Blue, a B2B SaaS process management platform.
Translate the text accurately and completely.
Maintain professional tone appropriate for enterprise software.
Keep product names like "Blue" unchanged.
Keep URLs, technical terms, and acronyms unchanged.
Preserve any HTML tags if present.
For UI text, keep translations concise and clear.`

	userPrompt := fmt.Sprintf("Translate the following text from English to %s:\n\n%s",
		languageNames[task.TargetLang], task.Entry.Value)

	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Temperature: openai.Float(0.3),
		MaxTokens:   openai.Int(500),
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return TranslationResult{
			Task:  task,
			Error: err,
		}
	}

	if len(completion.Choices) == 0 {
		return TranslationResult{
			Task:  task,
			Error: fmt.Errorf("no response from OpenAI"),
		}
	}

	return TranslationResult{
		Task:         task,
		Translation:  completion.Choices[0].Message.Content,
		InputTokens:  int(completion.Usage.PromptTokens),
		OutputTokens: int(completion.Usage.CompletionTokens),
		Error:        nil,
	}
}

func processResults(resultChan <-chan TranslationResult, jsonData map[string]interface{},
	jsonPath string, stats *TranslationStats,
	fileMutex *sync.Mutex, doneChan chan<- bool) {

	for result := range resultChan {
		if result.Error != nil {
			atomic.AddInt32(&stats.FailedTasks, 1)
			fmt.Printf("\nERROR [%s -> %s]: %v\n", result.Task.Entry.Key, result.Task.TargetLang, result.Error)
			continue
		}

		// Update statistics
		atomic.AddInt32(&stats.CompletedTasks, 1)
		atomic.AddInt64(&stats.TotalInputTokens, int64(result.InputTokens))
		atomic.AddInt64(&stats.TotalOutputTokens, int64(result.OutputTokens))

		// Update JSON data
		fileMutex.Lock()

		// Ensure language exists in JSON
		if _, ok := jsonData[result.Task.TargetLang]; !ok {
			jsonData[result.Task.TargetLang] = make(map[string]interface{})
		}

		// Set the nested value
		setNestedValue(jsonData[result.Task.TargetLang], result.Task.Entry.Path, result.Translation)

		// Save to file
		saveJSON(jsonPath, jsonData)

		fileMutex.Unlock()

		// Display progress
		completed := atomic.LoadInt32(&stats.CompletedTasks)
		failed := atomic.LoadInt32(&stats.FailedTasks)
		total := stats.TotalTasks
		percentage := float64(completed+failed) / float64(total) * 100

		fmt.Printf("\r[%d/%d] %.1f%% | Completed: %d | Failed: %d | %s: %s -> %s",
			completed+failed, total, percentage, completed, failed,
			result.Task.Entry.Key, "en", result.Task.TargetLang)
	}

	fmt.Println() // New line after progress
	doneChan <- true
}

func setNestedValue(obj interface{}, path []string, value string) {
	if len(path) == 0 {
		return
	}

	current := obj
	for _, key := range path[:len(path)-1] {
		m, ok := current.(map[string]interface{})
		if !ok {
			return
		}

		if _, exists := m[key]; !exists {
			m[key] = make(map[string]interface{})
		}

		current = m[key]
	}

	// Set the final value
	if m, ok := current.(map[string]interface{}); ok {
		m[path[len(path)-1]] = value
	}
}

func saveJSON(path string, jsonData map[string]interface{}) error {
	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func printFinalSummary(stats *TranslationStats) {
	duration := time.Since(stats.StartTime)
	completed := atomic.LoadInt32(&stats.CompletedTasks)
	failed := atomic.LoadInt32(&stats.FailedTasks)
	inputTokens := atomic.LoadInt64(&stats.TotalInputTokens)
	outputTokens := atomic.LoadInt64(&stats.TotalOutputTokens)

	totalCost := float64(inputTokens)/1_000_000*inputTokenCost +
		float64(outputTokens)/1_000_000*outputTokenCost

	fmt.Printf("\n=== Translation Complete ===\n")
	fmt.Printf("Total tasks: %d\n", stats.TotalTasks)
	fmt.Printf("Completed: %d\n", completed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Duration: %s\n", duration.Round(time.Second))
	fmt.Printf("Tokens used: %d input, %d output\n", inputTokens, outputTokens)
	fmt.Printf("Estimated cost: $%.4f\n", totalCost)

	if failed > 0 {
		fmt.Printf("\nNote: %d translations failed. Run the script again to retry.\n", failed)
	}
}

// findFilesNeedingTranslation scans the translations directory and returns files that need translation
func findFilesNeedingTranslation(targetLanguages []string) ([]string, error) {
	translationsDir := "translations"

	// Read directory contents
	entries, err := os.ReadDir(translationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read translations directory: %v", err)
	}

	var filesToProcess []string

	// Check each JSON file
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			needsTranslation, err := fileNeedsTranslation(entry.Name(), targetLanguages)
			if err != nil {
				fmt.Printf("Warning: Error checking %s: %v\n", entry.Name(), err)
				continue
			}

			if needsTranslation {
				filesToProcess = append(filesToProcess, entry.Name())
			}
		}
	}

	return filesToProcess, nil
}

// fileNeedsTranslation checks if a file has missing translations
func fileNeedsTranslation(jsonFile string, targetLanguages []string) (bool, error) {
	jsonPath := "translations/" + jsonFile
	jsonData, err := loadJSON(jsonPath)
	if err != nil {
		return false, err
	}

	// Check if English section exists
	if _, ok := jsonData["en"].(map[string]interface{}); !ok {
		// No English section means no translations needed
		return false, nil
	}

	// Extract English entries
	englishEntries, err := extractEnglishEntries(jsonData)
	if err != nil {
		return false, err
	}

	if len(englishEntries) == 0 {
		// No English content to translate
		return false, nil
	}

	// Check if any language is missing translations
	for _, lang := range targetLanguages {
		langData, exists := jsonData[lang].(map[string]interface{})
		if !exists {
			// Language section doesn't exist - needs translation
			return true, nil
		}

		// Check if all English keys are translated in this language
		existingKeys := make(map[string]bool)
		getAllKeys(langData, []string{}, existingKeys)

		for _, entry := range englishEntries {
			if !existingKeys[entry.Key] {
				// Missing translation found
				return true, nil
			}
		}
	}

	// All languages have all translations
	return false, nil
}

// processFile handles the translation of a single file (extracted from main logic)
func processFile(jsonFile string, targetLanguages []string, client *openai.Client) error {
	// Load existing JSON file
	jsonPath := "translations/" + jsonFile
	jsonData, err := loadJSON(jsonPath)
	if err != nil {
		return fmt.Errorf("error loading %s: %v", jsonPath, err)
	}

	// Extract English entries
	englishEntries, err := extractEnglishEntries(jsonData)
	if err != nil {
		return fmt.Errorf("error extracting English entries: %v", err)
	}
	fmt.Printf("Found %d English entries to translate in %s\n", len(englishEntries), jsonFile)

	// Build translation tasks
	tasks := buildTranslationTasks(englishEntries, targetLanguages, jsonData)

	if len(tasks) == 0 {
		fmt.Printf("✅ All translations are complete for %s!\n", jsonFile)
		return nil
	}

	fmt.Printf("\n=== Translation Summary ===\n")
	fmt.Printf("File: %s\n", jsonFile)
	fmt.Printf("Total translations needed: %d\n", len(tasks))
	fmt.Printf("Languages: %v\n", targetLanguages)

	// Initialize statistics
	stats := &TranslationStats{
		TotalTasks: len(tasks),
		StartTime:  time.Now(),
	}

	// Number of concurrent workers
	const workerCount = 50

	// Create channels
	taskChan := make(chan TranslationTask, len(tasks))
	resultChan := make(chan TranslationResult, len(tasks))

	// Mutex for file writing
	var fileMutex sync.Mutex

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go translationWorker(i, client, taskChan, resultChan, &wg)
	}

	// Send tasks to channel
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// Start result processor
	doneChan := make(chan bool)
	go processResults(resultChan, jsonData, jsonPath, stats, &fileMutex, doneChan)

	// Wait for workers
	wg.Wait()
	close(resultChan)

	// Wait for result processor
	<-doneChan

	// Final summary
	printFinalSummary(stats)

	return nil
}
