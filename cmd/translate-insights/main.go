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


// TranslationTask represents a single translation job
type TranslationTask struct {
	SourceFile   string
	TargetFile   string
	SourceLang   string
	TargetLang   string
	SourcePath   string
	TargetPath   string
}

// TranslationStats tracks statistics for the translation session
type TranslationStats struct {
	TotalTasks        int
	CompletedTasks    int32  // Use atomic for thread-safe updates
	TotalInputTokens  int64  // Use atomic for thread-safe updates
	TotalOutputTokens int64  // Use atomic for thread-safe updates
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
	inputTokenCost  = 0.4    // $0.4 per 1M input tokens
	outputTokenCost = 1.60   // $1.60 per 1M output tokens
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY not found in environment variables")
		os.Exit(1)
	}

	// Initialize OpenAI client with timeout configuration
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithRequestTimeout(5 * time.Minute), // Increase timeout for large translations
	)

	// Get base content directory
	contentDir := filepath.Join(".", "content")

	// Discover available languages
	languages, err := discoverLanguages(contentDir)
	if err != nil {
		fmt.Printf("Error discovering languages: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d languages: %v\n", len(languages), languages)

	// Get English insights
	englishInsights, err := getInsights(filepath.Join(contentDir, "en", "insights"))
	if err != nil {
		fmt.Printf("Error reading English insights: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d English insights\n", len(englishInsights))

	// Build translation queue
	tasks := buildTranslationQueue(contentDir, languages, englishInsights)

	if len(tasks) == 0 {
		fmt.Println("No translations needed. All insights are up to date!")
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

	// Number of concurrent workers
	const workerCount = 100

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
		go translationWorker(i, &client, systemPrompt, taskChan, resultChan, &wg)
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
	inputTokens := atomic.LoadInt64(&stats.TotalInputTokens)
	outputTokens := atomic.LoadInt64(&stats.TotalOutputTokens)
	totalCost := calculateCost(int(inputTokens), int(outputTokens))
	avgTime := totalTime.Seconds() / float64(completed)

	fmt.Printf("\n=== Translation Complete ===\n")
	fmt.Printf("Total tasks: %d\n", stats.TotalTasks)
	fmt.Printf("Completed: %d\n", completed)
	fmt.Printf("Total time: %s\n", formatDuration(totalTime))
	fmt.Printf("Average time: %.1fs/translation\n", avgTime)
	fmt.Printf("Total tokens: %d input, %d output\n", inputTokens, outputTokens)
	fmt.Printf("Total cost: $%.4f\n", totalCost)
}

// translationWorker processes translation tasks from the channel
func translationWorker(id int, client *openai.Client, systemPrompt string, taskChan <-chan TranslationTask, resultChan chan<- TranslationResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChan {
		fmt.Printf("\n  WORKER %d: Starting %s -> %s\n", id, task.SourceFile, task.TargetLang)
		startTime := time.Now()
		inputTokens, outputTokens, err := translateInsight(client, task, systemPrompt)
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
		completed := atomic.AddInt32(&stats.CompletedTasks, 1)
		
		// Update token counts atomically
		if result.Error == nil {
			atomic.AddInt64(&stats.TotalInputTokens, int64(result.InputTokens))
			atomic.AddInt64(&stats.TotalOutputTokens, int64(result.OutputTokens))
		}

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

		// Show errors on new line
		if result.Error != nil {
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
			// Check if it has an insights directory
			insightsPath := filepath.Join(contentDir, entry.Name(), "insights")
			if _, err := os.Stat(insightsPath); err == nil {
				languages = append(languages, entry.Name())
			}
		}
	}

	return languages, nil
}

func getInsights(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var insights []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			insights = append(insights, entry.Name())
		}
	}

	return insights, nil
}

func buildTranslationQueue(contentDir string, languages []string, englishInsights []string) []TranslationTask {
	var tasks []TranslationTask

	for _, lang := range languages {
		targetInsightsDir := filepath.Join(contentDir, lang, "insights")
		
		// Get existing translations
		existingInsights, _ := getInsights(targetInsightsDir)
		existingMap := make(map[string]bool)
		for _, insight := range existingInsights {
			existingMap[insight] = true
		}

		// Find missing translations
		for _, insight := range englishInsights {
			if !existingMap[insight] {
				task := TranslationTask{
					SourceFile: insight,
					TargetFile: insight,
					SourceLang: "en",
					TargetLang: lang,
					SourcePath: filepath.Join(contentDir, "en", "insights", insight),
					TargetPath: filepath.Join(contentDir, lang, "insights", insight),
				}
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}

func getSystemPrompt() string {
	return `You are a professional translator for the Blue website, a B2B SaaS process management platform. You must translate the ENTIRE markdown file including all body content.

WHAT TO TRANSLATE:
1. The "title" field value in the frontmatter
2. The "description" field value in the frontmatter  
3. ALL BODY CONTENT after the closing --- of the frontmatter
4. All headings, paragraphs, lists, and text content

NEVER TRANSLATE:
- "Blue" (product name)
- Email addresses (support@blue.cc, sales@blue.cc, etc.)
- URLs and links
- Image filenames
- Code blocks
- Technical terms (API, GraphQL, webhook, OAuth, etc.)
- Person names
- Field names like "title:" and "description:" (only translate their values)

IMPORTANT: You MUST translate the ENTIRE document body, not just the frontmatter. Every paragraph, heading, and text element must be translated to the target language.

FORMAT:
- Return the complete translated markdown file
- Preserve all markdown formatting (headers, lists, links, etc.)
- Keep the same structure as the original
- Do not add any explanations or comments`
}

func translateInsight(client *openai.Client, task TranslationTask, systemPrompt string) (int, int, error) {
	// Read source file
	content, err := os.ReadFile(task.SourcePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read source file: %w", err)
	}

	// Check if file is large (more than 100 lines after stripping frontmatter)
	lines := strings.Split(string(content), "\n")
	if len(lines) > 100 {
		// Use chunked translation for large files
		return translateLargeInsight(client, task, systemPrompt, string(content))
	}

	// For smaller files, use the original single-request method
	return translateSmallInsight(client, task, systemPrompt, string(content))
}

// translateSmallInsight handles files that can be translated in a single request
func translateSmallInsight(client *openai.Client, task TranslationTask, systemPrompt string, content string) (int, int, error) {
	// Strip out category and date before sending to OpenAI
	contentForTranslation, originalLines := stripPreservableFields(content)

	// Prepare user prompt with stripped content
	userPrompt := fmt.Sprintf("Translate the following ENTIRE markdown document from English to %s. Translate ALL content including the body text, not just the frontmatter:\n\n%s", 
		languageNames[task.TargetLang], contentForTranslation)

	// Create chat completion request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
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
		return 0, 0, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(completion.Choices) == 0 {
		return 0, 0, fmt.Errorf("no response from OpenAI")
	}

	translatedContent := completion.Choices[0].Message.Content

	// Restore the original category and date fields
	translatedContent = restorePreservableFields(translatedContent, originalLines)

	// Ensure target directory exists
	targetDir := filepath.Dir(task.TargetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write translated file
	if err := os.WriteFile(task.TargetPath, []byte(translatedContent), 0644); err != nil {
		return 0, 0, fmt.Errorf("failed to write translated file: %w", err)
	}

	// Get usage information
	inputTokens := int(completion.Usage.PromptTokens)
	outputTokens := int(completion.Usage.CompletionTokens)

	return inputTokens, outputTokens, nil
}

// translateLargeInsight handles large files by splitting them into chunks
func translateLargeInsight(client *openai.Client, task TranslationTask, systemPrompt string, content string) (int, int, error) {
	fmt.Printf("\n  CHUNKING [%s -> %s]: File has %d lines, using chunked translation\n", 
		task.SourceFile, task.TargetLang, len(strings.Split(content, "\n")))
	
	// First, extract and preserve the frontmatter
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return 0, 0, fmt.Errorf("invalid markdown format")
	}

	frontmatter := parts[1]
	body := parts[2]

	// Strip category and date from frontmatter
	strippedFrontmatter, originalLines := stripPreservableFieldsFromFrontmatter(frontmatter)

	// Translate frontmatter separately
	fmt.Printf("  CHUNK [%s -> %s]: Translating frontmatter...\n", task.SourceFile, task.TargetLang)
	frontmatterTokensIn, frontmatterTokensOut, translatedFrontmatter, err := translateFrontmatter(client, task, systemPrompt, strippedFrontmatter)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to translate frontmatter: %w", err)
	}

	// Split body into chunks (by paragraphs to maintain context)
	chunks := splitBodyIntoChunks(body, 25) // 25 lines per chunk for smaller requests
	fmt.Printf("  CHUNK [%s -> %s]: Split into %d chunks\n", task.SourceFile, task.TargetLang, len(chunks))
	
	var translatedChunks []string
	totalInputTokens := frontmatterTokensIn
	totalOutputTokens := frontmatterTokensOut

	// Translate each chunk
	for i, chunk := range chunks {
		fmt.Printf("  CHUNK [%s -> %s]: Translating chunk %d/%d (%d chars)...\n", 
			task.SourceFile, task.TargetLang, i+1, len(chunks), len(chunk))
		
		chunkPrompt := fmt.Sprintf("Translate the following markdown content from English to %s. This is part %d of %d of a larger document. Maintain all markdown formatting:\n\n%s", 
			languageNames[task.TargetLang], i+1, len(chunks), chunk)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		params := openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT4oMini,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a professional translator. Translate the given markdown content while preserving all formatting, links, and structure."),
				openai.UserMessage(chunkPrompt),
			},
			Temperature: openai.Float(0.3),
		}

		completion, err := client.Chat.Completions.New(ctx, params)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to translate chunk %d: %w", i+1, err)
		}

		if len(completion.Choices) == 0 {
			return 0, 0, fmt.Errorf("no response for chunk %d", i+1)
		}

		translatedChunks = append(translatedChunks, completion.Choices[0].Message.Content)
		totalInputTokens += int(completion.Usage.PromptTokens)
		totalOutputTokens += int(completion.Usage.CompletionTokens)
	}

	// Reassemble the document
	restoredFrontmatter := restoreFrontmatterFields(translatedFrontmatter, originalLines)
	translatedBody := strings.Join(translatedChunks, "\n\n")
	finalContent := fmt.Sprintf("---\n%s\n---\n%s", restoredFrontmatter, translatedBody)

	// Ensure target directory exists
	targetDir := filepath.Dir(task.TargetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write translated file
	if err := os.WriteFile(task.TargetPath, []byte(finalContent), 0644); err != nil {
		return 0, 0, fmt.Errorf("failed to write translated file: %w", err)
	}

	return totalInputTokens, totalOutputTokens, nil
}

// stripPreservableFields removes category and date from content before translation
func stripPreservableFields(content string) (strippedContent string, originalLines []string) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return content, nil
	}

	frontmatterStr := parts[1]
	body := parts[2]

	// Store original lines that contain category and date
	originalLines = []string{}
	
	// Process frontmatter line by line to preserve exact formatting
	var newFrontmatterLines []string
	lines := strings.Split(strings.TrimSpace(frontmatterStr), "\n")
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check if this line contains category or date
		if strings.HasPrefix(trimmed, "category:") || strings.HasPrefix(trimmed, "date:") {
			// Store the exact original line
			originalLines = append(originalLines, line)
		} else if trimmed != "" {
			// Keep other lines
			newFrontmatterLines = append(newFrontmatterLines, line)
		}
	}

	// Reconstruct content without category and date
	if len(newFrontmatterLines) > 0 {
		strippedContent = fmt.Sprintf("---\n%s\n---\n%s", strings.Join(newFrontmatterLines, "\n"), body)
	} else {
		strippedContent = fmt.Sprintf("---\n---\n%s", body)
	}
	
	return strippedContent, originalLines
}

// restorePreservableFields adds back the original category and date fields
func restorePreservableFields(translatedContent string, originalLines []string) string {
	parts := strings.SplitN(translatedContent, "---", 3)
	if len(parts) < 3 || originalLines == nil {
		return translatedContent
	}

	frontmatterStr := parts[1]
	body := parts[2]

	// Process translated frontmatter line by line
	lines := strings.Split(strings.TrimSpace(frontmatterStr), "\n")
	
	// Fix any accidentally translated field names
	translations := map[string]string{
		"título": "title", "titre": "title", "titel": "title", "titolo": "title", 
		"tytuł": "title", "заголовок": "title", "タイトル": "title", "제목": "title", 
		"标题": "title", "標題": "title",
		"descripción": "description", "beschreibung": "description", "descrizione": "description", 
		"descrição": "description", "opis": "description", "описание": "description", 
		"説明": "description", "설명": "description", "描述": "description",
	}
	
	var fixedLines []string
	for _, line := range lines {
		fixedLine := line
		for wrong, correct := range translations {
			if strings.Contains(line, wrong+":") {
				fixedLine = strings.Replace(line, wrong+":", correct+":", 1)
				break
			}
		}
		if strings.TrimSpace(fixedLine) != "" {
			fixedLines = append(fixedLines, fixedLine)
		}
	}

	// Add back the original category and date lines
	fixedLines = append(fixedLines, originalLines...)

	// Reconstruct with restored fields
	return fmt.Sprintf("---\n%s\n---\n%s", strings.Join(fixedLines, "\n"), body)
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

// Helper functions for large file translation

// stripPreservableFieldsFromFrontmatter strips category and date from frontmatter only
func stripPreservableFieldsFromFrontmatter(frontmatter string) (stripped string, originalLines []string) {
	lines := strings.Split(strings.TrimSpace(frontmatter), "\n")
	var newLines []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "category:") || strings.HasPrefix(trimmed, "date:") {
			originalLines = append(originalLines, line)
		} else if trimmed != "" {
			newLines = append(newLines, line)
		}
	}
	
	return strings.Join(newLines, "\n"), originalLines
}

// translateFrontmatter translates only the frontmatter portion
func translateFrontmatter(client *openai.Client, task TranslationTask, systemPrompt string, frontmatter string) (int, int, string, error) {
	prompt := fmt.Sprintf("Translate the following YAML frontmatter fields from English to %s. Only translate the values of 'title' and 'description' fields, not the field names themselves:\n\n%s", 
		languageNames[task.TargetLang], frontmatter)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	params := openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a professional translator. Translate only the values of the title and description fields. Keep field names in English."),
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.3),
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return 0, 0, "", err
	}

	if len(completion.Choices) == 0 {
		return 0, 0, "", fmt.Errorf("no response from OpenAI")
	}

	translatedFrontmatter := completion.Choices[0].Message.Content
	inputTokens := int(completion.Usage.PromptTokens)
	outputTokens := int(completion.Usage.CompletionTokens)

	return inputTokens, outputTokens, translatedFrontmatter, nil
}

// splitBodyIntoChunks splits the markdown body into manageable chunks
func splitBodyIntoChunks(body string, linesPerChunk int) []string {
	// Split by double newlines to preserve paragraphs
	paragraphs := strings.Split(body, "\n\n")
	
	var chunks []string
	var currentChunk []string
	currentLines := 0
	
	for _, para := range paragraphs {
		paraLines := len(strings.Split(para, "\n"))
		
		// If adding this paragraph would exceed the limit, start a new chunk
		if currentLines > 0 && currentLines+paraLines > linesPerChunk {
			chunks = append(chunks, strings.Join(currentChunk, "\n\n"))
			currentChunk = []string{para}
			currentLines = paraLines
		} else {
			currentChunk = append(currentChunk, para)
			currentLines += paraLines
		}
	}
	
	// Add the last chunk
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, "\n\n"))
	}
	
	return chunks
}

// restoreFrontmatterFields restores the original category and date fields
func restoreFrontmatterFields(translatedFrontmatter string, originalLines []string) string {
	lines := strings.Split(strings.TrimSpace(translatedFrontmatter), "\n")
	
	// Fix any accidentally translated field names
	translations := map[string]string{
		"título": "title", "titre": "title", "titel": "title", "titolo": "title", 
		"tytuł": "title", "заголовок": "title", "タイトル": "title", "제목": "title", 
		"标题": "title", "標題": "title",
		"descripción": "description", "beschreibung": "description", "descrizione": "description", 
		"descrição": "description", "opis": "description", "описание": "description", 
		"説明": "description", "설명": "description", "描述": "description",
	}
	
	var fixedLines []string
	for _, line := range lines {
		fixedLine := line
		for wrong, correct := range translations {
			if strings.Contains(line, wrong+":") {
				fixedLine = strings.Replace(line, wrong+":", correct+":", 1)
				break
			}
		}
		if strings.TrimSpace(fixedLine) != "" {
			fixedLines = append(fixedLines, fixedLine)
		}
	}
	
	// Add back the original category and date lines
	fixedLines = append(fixedLines, originalLines...)
	
	return strings.Join(fixedLines, "\n")
}