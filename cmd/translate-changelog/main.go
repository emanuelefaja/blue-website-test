package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// ChangelogEntry represents a single changelog entry to translate
type ChangelogEntry struct {
	Key     string // e.g., "2025_june_zapier_integration_title"
	Value   string // The text to translate
	IsTitle bool   // true for _title, false for _description
	Year    string
	Month   string
	BaseKey string // e.g., "2025_june_zapier_integration"
}

// TranslationTask represents a single translation job
type TranslationTask struct {
	Entry      ChangelogEntry
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

	// Load existing changelog
	changelogPath := "translations/changelog.json"
	changelogData, err := loadChangelog(changelogPath)
	if err != nil {
		fmt.Printf("Error loading changelog: %v\n", err)
		os.Exit(1)
	}

	// Get target languages from web/languages.go
	targetLanguages, err := getTargetLanguages()
	if err != nil {
		fmt.Printf("Error reading target languages: %v\n", err)
		os.Exit(1)
	}

	// Extract English entries (skip UI translations)
	englishEntries := extractEnglishEntries(changelogData)
	fmt.Printf("Found %d English changelog entries to translate\n", len(englishEntries))

	// Build translation tasks
	tasks := buildTranslationTasks(englishEntries, targetLanguages, changelogData)

	if len(tasks) == 0 {
		fmt.Println("All translations are complete!")
		return
	}

	fmt.Printf("\n=== Translation Summary ===\n")
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
		go translationWorker(i, &client, taskChan, resultChan, &wg)
	}

	// Send tasks to channel
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// Start result processor
	doneChan := make(chan bool)
	go processResults(resultChan, changelogData, changelogPath, stats, &fileMutex, doneChan)

	// Wait for workers
	wg.Wait()
	close(resultChan)

	// Wait for result processor
	<-doneChan

	// Final summary
	printFinalSummary(stats)
}

func loadChangelog(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var changelog map[string]interface{}
	if err := json.Unmarshal(data, &changelog); err != nil {
		return nil, err
	}

	return changelog, nil
}

func extractEnglishEntries(changelog map[string]interface{}) []ChangelogEntry {
	var entries []ChangelogEntry

	englishData, ok := changelog["en"].(map[string]interface{})
	if !ok {
		return entries
	}

	// First, add UI translations
	uiKeys := []string{"page_title", "page_subtitle", "no_entries"}
	for _, key := range uiKeys {
		if value, ok := englishData[key].(string); ok {
			entries = append(entries, ChangelogEntry{
				Key:     key,
				Value:   value,
				IsTitle: false,  // UI strings are more like descriptions
				Year:    "0000", // Sort UI translations first
				Month:   "ui",
				BaseKey: key,
			})
		}
	}

	// Handle months separately
	if monthsData, ok := englishData["months"].(map[string]interface{}); ok {
		for monthKey, monthValue := range monthsData {
			if strValue, ok := monthValue.(string); ok {
				fullKey := "months." + monthKey
				entries = append(entries, ChangelogEntry{
					Key:     fullKey,
					Value:   strValue,
					IsTitle: false,
					Year:    "0000",
					Month:   "ui",
					BaseKey: fullKey,
				})
			}
		}
	}

	// Then add changelog entries
	for key, value := range englishData {
		strValue, ok := value.(string)
		if !ok {
			continue
		}

		// Skip UI translations we already handled
		if key == "page_title" || key == "page_subtitle" || key == "no_entries" || key == "months" {
			continue
		}

		// Parse the key to extract year, month, and type
		parts := strings.Split(key, "_")
		if len(parts) < 4 {
			continue
		}

		year := parts[0]
		month := parts[1]
		isTitle := strings.HasSuffix(key, "_title")
		baseKey := strings.TrimSuffix(key, "_title")
		baseKey = strings.TrimSuffix(baseKey, "_description")

		entries = append(entries, ChangelogEntry{
			Key:     key,
			Value:   strValue,
			IsTitle: isTitle,
			Year:    year,
			Month:   month,
			BaseKey: baseKey,
		})
	}

	// Sort entries: UI first, then chronologically
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Year != entries[j].Year {
			return entries[i].Year < entries[j].Year
		}
		return entries[i].Key < entries[j].Key
	})

	return entries
}

func buildTranslationTasks(entries []ChangelogEntry, languages []string, changelog map[string]interface{}) []TranslationTask {
	var tasks []TranslationTask

	for _, lang := range languages {
		// Get existing translations for this language
		var existingKeys map[string]bool
		if langData, ok := changelog[lang].(map[string]interface{}); ok {
			existingKeys = make(map[string]bool)

			// Check regular keys
			for key := range langData {
				if key != "months" {
					existingKeys[key] = true
				}
			}

			// Check months keys
			if monthsData, ok := langData["months"].(map[string]interface{}); ok {
				for monthKey := range monthsData {
					existingKeys["months."+monthKey] = true
				}
			}
		} else {
			existingKeys = make(map[string]bool)
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

func translationWorker(id int, client *openai.Client, taskChan <-chan TranslationTask, resultChan chan<- TranslationResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		result := translateEntry(client, task)
		resultChan <- result
	}
}

func translateEntry(client *openai.Client, task TranslationTask) TranslationResult {
	ctx := context.Background()

	// Prepare the prompt based on entry type
	systemPrompt := getSystemPrompt(task.Entry.IsTitle)
	userPrompt := fmt.Sprintf("Translate the following changelog entry from English to %s:\n\n%s",
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

func getSystemPrompt(isTitle bool) string {
	if isTitle {
		return `You are a professional translator for Blue, a B2B SaaS process management platform. 
Translate the changelog entry title to be concise and clear.
Maintain technical accuracy and professional tone.
Keep product names like "Blue", "Zapier", "Make.com", etc. unchanged.
Keep technical terms like API, CSV, PDF, OAuth, etc. unchanged.`
	}

	return `You are a professional translator for Blue, a B2B SaaS process management platform.
Translate the text accurately and completely.
Maintain professional tone appropriate for enterprise software.
Keep product names like "Blue" unchanged.
Keep URLs, technical terms, and acronyms unchanged.
For UI text (page titles, navigation), keep translations concise and clear.
For changelog descriptions, preserve any markdown formatting if present.`
}

func processResults(resultChan <-chan TranslationResult, changelog map[string]interface{},
	changelogPath string, stats *TranslationStats,
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

		// Update changelog data
		fileMutex.Lock()

		// Ensure language exists in changelog
		if _, ok := changelog[result.Task.TargetLang]; !ok {
			changelog[result.Task.TargetLang] = make(map[string]interface{})
		}

		langData := changelog[result.Task.TargetLang].(map[string]interface{})

		// Handle nested months structure
		if strings.HasPrefix(result.Task.Entry.Key, "months.") {
			// Ensure months object exists
			if _, ok := langData["months"]; !ok {
				langData["months"] = make(map[string]interface{})
			}
			monthsData := langData["months"].(map[string]interface{})
			monthKey := strings.TrimPrefix(result.Task.Entry.Key, "months.")
			monthsData[monthKey] = result.Translation
		} else {
			// Regular key
			langData[result.Task.Entry.Key] = result.Translation
		}

		// Save to file
		saveChangelog(changelogPath, changelog)

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

func saveChangelog(path string, changelog map[string]interface{}) error {
	data, err := json.MarshalIndent(changelog, "", "  ")
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
