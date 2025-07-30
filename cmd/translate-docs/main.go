package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// TranslationTask represents a single doc translation job
type TranslationTask struct {
	SourceFile string
	TargetFile string
	SourceLang string
	TargetLang string
	SourcePath string
	TargetPath string
}

// TranslationStats tracks statistics for the translation session
type TranslationStats struct {
	TotalTasks        int
	CompletedTasks    int32 // Use atomic for thread-safe updates
	FailedTasks       int32 // Track failed translations
	TotalInputTokens  int64 // Use atomic for thread-safe updates
	TotalOutputTokens int64 // Use atomic for thread-safe updates
	StartTime         time.Time
	mu                sync.Mutex
}

// TranslationResult represents the result of a translation task
type TranslationResult struct {
	Task         TranslationTask
	InputTokens  int
	OutputTokens int
	Error        error
	Duration     time.Duration
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
	// OpenAI pricing per 1M tokens
	inputTokenCost  = 0.4  // $0.4 per 1M input tokens
	outputTokenCost = 1.60 // $1.60 per 1M output tokens
)

func main() {
	// Load environment variables from root directory
	if err := godotenv.Load(".env"); err != nil {
		// Try without .env file (in case env vars are already set)
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY not found in environment variables")
		fmt.Println("Please ensure OPENAI_API_KEY is set in your environment or in the .env file in the root directory")
		os.Exit(1)
	}

	// Initialize OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	// Initialize the document processor
	processor := NewDocumentProcessor()

	// Get base content directory
	contentDir := filepath.Join(".", "content")

	// Discover available languages
	languages, err := discoverLanguages(contentDir)
	if err != nil {
		fmt.Printf("Error discovering languages: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d languages: %v\n", len(languages), languages)

	// Get English documentation files
	englishDocs, err := getDocs(filepath.Join(contentDir, "en", "docs"))
	if err != nil {
		fmt.Printf("Error reading English documentation: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d English documentation files\n", len(englishDocs))

	// Build translation queue
	tasks := buildTranslationQueue(contentDir, languages, englishDocs)

	if len(tasks) == 0 {
		fmt.Println("No translations needed. All documentation is up to date!")
		return
	}

	// Display translation summary
	fmt.Printf("\n=== Translation Summary ===\n")
	tasksByLang := make(map[string]int)
	for _, task := range tasks {
		tasksByLang[task.TargetLang]++
	}

	for lang, count := range tasksByLang {
		fmt.Printf("%s (%s): %d translations needed\n", languageNames[lang], lang, count)
	}
	fmt.Printf("\nTotal translations required: %d\n", len(tasks))

	// Initialize statistics
	stats := &TranslationStats{
		TotalTasks: len(tasks),
		StartTime:  time.Now(),
	}

	// Number of concurrent workers (moderate for docs due to complexity)
	const workerCount = 75

	// Start translation process
	fmt.Printf("\n=== Starting Translation Process ===\n")
	fmt.Printf("Running %d parallel workers...\n", workerCount)

	systemPrompt := getSystemPrompt()

	// Create channels
	taskChan := make(chan TranslationTask, len(tasks))
	resultChan := make(chan TranslationResult, len(tasks))

	// Create wait group for workers
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go translationWorker(i, &client, processor, systemPrompt, taskChan, resultChan, &wg)
	}

	// Send all tasks to the channel
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	// Start result collector in a separate goroutine
	doneChan := make(chan bool)
	go collectResults(stats, resultChan, len(tasks), doneChan)

	// Wait for all workers to complete
	wg.Wait()
	close(resultChan)

	// Wait for result collector to finish
	<-doneChan

	// Final summary
	totalTime := time.Since(stats.StartTime)
	completed := atomic.LoadInt32(&stats.CompletedTasks)
	failed := atomic.LoadInt32(&stats.FailedTasks)
	inputTokens := atomic.LoadInt64(&stats.TotalInputTokens)
	outputTokens := atomic.LoadInt64(&stats.TotalOutputTokens)
	totalCost := calculateCost(int(inputTokens), int(outputTokens))
	avgTime := totalTime.Seconds() / float64(completed)

	fmt.Printf("\n=== Translation Complete ===\n")
	fmt.Printf("Total tasks: %d\n", stats.TotalTasks)
	fmt.Printf("Completed: %d\n", completed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total time: %s\n", formatDuration(totalTime))
	fmt.Printf("Average time: %.1fs/translation\n", avgTime)
	fmt.Printf("Total tokens: %d input, %d output\n", inputTokens, outputTokens)
	fmt.Printf("Total cost: $%.4f\n", totalCost)
}

// translationWorker processes translation tasks from the channel
func translationWorker(id int, client *openai.Client, processor *DocumentProcessor, systemPrompt string, taskChan <-chan TranslationTask, resultChan chan<- TranslationResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		startTime := time.Now()
		inputTokens, outputTokens, err := translateDoc(client, processor, task, systemPrompt)
		duration := time.Since(startTime)

		result := TranslationResult{
			Task:         task,
			InputTokens:  inputTokens,
			OutputTokens: outputTokens,
			Error:        err,
			Duration:     duration,
		}

		resultChan <- result
	}
}

// collectResults collects results from workers and updates progress
func collectResults(stats *TranslationStats, resultChan <-chan TranslationResult, totalTasks int, doneChan chan<- bool) {
	for result := range resultChan {
		if result.Error == nil {
			completed := atomic.AddInt32(&stats.CompletedTasks, 1)
			atomic.AddInt64(&stats.TotalInputTokens, int64(result.InputTokens))
			atomic.AddInt64(&stats.TotalOutputTokens, int64(result.OutputTokens))

			// Calculate progress
			progress := float64(completed) / float64(totalTasks) * 100
			elapsed := time.Since(stats.StartTime)
			avgTime := elapsed.Seconds() / float64(completed)

			// Get current totals
			inputTokens := atomic.LoadInt64(&stats.TotalInputTokens)
			outputTokens := atomic.LoadInt64(&stats.TotalOutputTokens)
			currentCost := calculateCost(int(inputTokens), int(outputTokens))

			// Display progress
			fmt.Printf("\r[%d/%d] %.1f%% | Avg: %.1fs | Elapsed: %s | Tokens: %d in, %d out | Cost: $%.4f",
				completed, totalTasks, progress, avgTime, formatDuration(elapsed),
				inputTokens, outputTokens, currentCost)
		} else {
			atomic.AddInt32(&stats.FailedTasks, 1)
			// Show errors on new line
			fmt.Printf("\n  ERROR [%s -> %s]: %v\n",
				result.Task.SourceFile, result.Task.TargetLang, result.Error)
		}
	}

	// Clear the progress line and move to next line
	fmt.Println()
	doneChan <- true
}

func discoverLanguages(contentDir string) ([]string, error) {
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "en" {
			// Check if it has a docs directory
			docsPath := filepath.Join(contentDir, entry.Name(), "docs")
			if _, err := os.Stat(docsPath); err == nil {
				languages = append(languages, entry.Name())
			}
		}
	}

	return languages, nil
}

func getDocs(dir string) ([]string, error) {
	var docs []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			// Get relative path from docs directory
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			docs = append(docs, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return docs, nil
}

func buildTranslationQueue(contentDir string, languages []string, englishDocs []string) []TranslationTask {
	var tasks []TranslationTask

	for _, lang := range languages {
		targetDocsDir := filepath.Join(contentDir, lang, "docs")

		// Create docs directory if it doesn't exist
		os.MkdirAll(targetDocsDir, 0755)

		// Check which docs are missing
		for _, doc := range englishDocs {
			targetPath := filepath.Join(targetDocsDir, doc)

			// Check if translation already exists
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				task := TranslationTask{
					SourceFile: doc,
					TargetFile: doc,
					SourceLang: "en",
					TargetLang: lang,
					SourcePath: filepath.Join(contentDir, "en", "docs", doc),
					TargetPath: targetPath,
				}
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}

func getSystemPrompt() string {
	return `You are a professional translator for Blue's documentation. You must translate the ENTIRE markdown file while preserving all technical content.

CRITICAL PRESERVATION RULES:
1. NEVER translate any text matching these patterns:
   - @@[A-Z]+##[a-f0-9-]+##[A-Z]+@@
   - These are technical markers that MUST remain EXACTLY as written
   
2. Examples that MUST NOT change:
   - @@CB##550e8400-e29b-41d4-a716-446655440000##CB@@
   - @@CODE##abc123##CODE@@
   - @@LINK##def456##LINK@@
   - @@URL##ghi789##URL@@
   - @@CALLOUT##xyz789##CALLOUT@@
   - @@VIDEO##vid123##VIDEO@@

3. If you see these markers, copy them EXACTLY character-for-character

WHAT TO TRANSLATE:
1. The "title" field value in the frontmatter
2. The "description" field value in the frontmatter  
3. ALL BODY CONTENT after the closing --- of the frontmatter
4. All headings, paragraphs, lists, and descriptive text
5. Comments in code examples (preceded by # or //)
6. ALL TEXT CONTENT between callout markers, including:
   - Callout headings (e.g., "**Warning**")
   - Callout body text and lists
   - Any warnings, notes, or tips inside callouts

NEVER TRANSLATE:
- Technical markers (@@...@@)
- "Blue" (product name)
- Technical terms (API, GraphQL, HTTP, JSON, OAuth, CSV, etc.)
- Email addresses
- URLs and internal documentation links (preserve /docs/... paths)
- Parameter names, field names, type names
- File paths and filenames
- Code blocks (except comments)
- The callout structure markers themselves (::callout, ::)
- Callout frontmatter keys (icon:, type:) - but DO translate their values if they are descriptive
- HTML tags like <video>, <source>, etc.

IMPORTANT CALLOUT RULES:
- Callout markers (@@CALLOUT##...##CALLOUT@@) preserve the structure
- ALWAYS translate the content BETWEEN these markers
- This includes bold text, lists, paragraphs, etc. inside callouts

FORMAT:
- Return the complete translated markdown file
- Preserve all markdown formatting
- Keep the same structure as the original
- Maintain all placeholders exactly as they appear`
}

func translateDoc(client *openai.Client, processor *DocumentProcessor, task TranslationTask, systemPrompt string) (int, int, error) {
	// Read source file
	content, err := os.ReadFile(task.SourcePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read source file: %w", err)
	}

	// Process the document (extract and mask technical content)
	maskedContent, placeholderMap, err := processor.ProcessDocument(string(content))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to process document: %w", err)
	}

	// Prepare user prompt with masked content
	userPrompt := fmt.Sprintf("Translate the following documentation from English to %s. Remember to preserve all @@...@@ markers exactly as they appear:\n\n%s",
		languageNames[task.TargetLang], maskedContent)

	// Try translation with retries
	var translatedContent string
	var inputTokens, outputTokens int
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create chat completion request
		ctx := context.Background()
		params := openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT4oMini,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt),
				openai.UserMessage(userPrompt),
			},
			Temperature: openai.Float(0.3), // Lower temperature for more consistent translations
		}

		completion, err := client.Chat.Completions.New(ctx, params)
		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(time.Second * time.Duration(attempt+1)) // Backoff
				continue
			}
			return 0, 0, fmt.Errorf("OpenAI API error after %d attempts: %w", maxRetries, err)
		}

		if len(completion.Choices) == 0 {
			return 0, 0, fmt.Errorf("no response from OpenAI")
		}

		translatedContent = completion.Choices[0].Message.Content
		inputTokens = int(completion.Usage.PromptTokens)
		outputTokens = int(completion.Usage.CompletionTokens)

		// Validate translation
		if err := processor.ValidateTranslation(maskedContent, translatedContent); err == nil {
			break // Translation is valid
		} else if attempt < maxRetries-1 {
			// Enhance prompt for next attempt
			systemPrompt = enhancePromptForRetry(systemPrompt, err)
		} else {
			// Last attempt failed, try to recover
			translatedContent = processor.RecoverPlaceholders(maskedContent, translatedContent)
		}
	}

	// Restore technical content
	finalContent := processor.RestoreContent(translatedContent, placeholderMap)

	// Ensure target directory exists (including subdirectories)
	targetDir := filepath.Dir(task.TargetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write translated file
	if err := os.WriteFile(task.TargetPath, []byte(finalContent), 0644); err != nil {
		return 0, 0, fmt.Errorf("failed to write translated file: %w", err)
	}

	return inputTokens, outputTokens, nil
}

func enhancePromptForRetry(prompt string, err error) string {
	addition := fmt.Sprintf("\n\nIMPORTANT: The previous translation failed validation: %v\nPlease ensure all @@...@@ markers are preserved EXACTLY as they appear in the source.", err)
	return prompt + addition
}

func calculateCost(inputTokens, outputTokens int) float64 {
	inputCost := float64(inputTokens) / 1_000_000 * inputTokenCost
	outputCost := float64(outputTokens) / 1_000_000 * outputTokenCost
	return inputCost + outputCost
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}