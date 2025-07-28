package web

import (
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// KeywordExtractor handles NLP-based keyword extraction from text
type KeywordExtractor struct {
	stopWords      map[string]bool
	wordRegex      *regexp.Regexp
	minWordLength  int
	maxKeywords    int
	phrasesEnabled bool
}

// NewKeywordExtractor creates a new keyword extractor with default settings
func NewKeywordExtractor() *KeywordExtractor {
	return &KeywordExtractor{
		stopWords:      initializeStopWords(),
		wordRegex:      regexp.MustCompile(`\b[\w'-]+\b`),
		minWordLength:  3,
		maxKeywords:    100,
		phrasesEnabled: true,
	}
}

// ExtractKeywords extracts the most relevant keywords from text
func (ke *KeywordExtractor) ExtractKeywords(text string, limit int) []string {
	if limit <= 0 {
		limit = ke.maxKeywords
	}

	// Limit text processing to first 5000 words for performance
	words := ke.tokenize(text)
	if len(words) > 5000 {
		words = words[:5000]
	}

	// Calculate term frequencies
	termFreq := ke.calculateTermFrequencies(words)

	// Extract phrases if enabled
	var phraseFreq map[string]int
	if ke.phrasesEnabled {
		phraseFreq = ke.extractPhrases(words)
	}

	// Combine and score all terms
	allTerms := ke.combineAndScore(termFreq, phraseFreq, len(words))

	// Return top N terms
	return ke.selectTopTerms(allTerms, limit)
}

// tokenize converts text to lowercase tokens
func (ke *KeywordExtractor) tokenize(text string) []string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Find all word matches
	matches := ke.wordRegex.FindAllString(text, -1)

	// Filter tokens
	var tokens []string
	for _, match := range matches {
		// Skip if too short
		if len(match) < ke.minWordLength {
			continue
		}

		// Skip if it's a number
		if isNumber(match) {
			continue
		}

		// Skip stop words
		if ke.stopWords[match] {
			continue
		}

		tokens = append(tokens, match)
	}

	return tokens
}

// calculateTermFrequencies counts occurrences of each term
func (ke *KeywordExtractor) calculateTermFrequencies(words []string) map[string]int {
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}
	return freq
}

// extractPhrases finds common 2-word and 3-word phrases
func (ke *KeywordExtractor) extractPhrases(words []string) map[string]int {
	phrases := make(map[string]int)

	// Extract 2-word phrases
	for i := 0; i < len(words)-1; i++ {
		// Skip if either word would be a stop word in isolation
		if ke.isCommonWord(words[i]) || ke.isCommonWord(words[i+1]) {
			continue
		}
		phrase := words[i] + " " + words[i+1]
		phrases[phrase]++
	}

	// Extract 3-word phrases (less common, so higher threshold)
	for i := 0; i < len(words)-2; i++ {
		// Skip if any word is too common
		if ke.isCommonWord(words[i]) || ke.isCommonWord(words[i+1]) || ke.isCommonWord(words[i+2]) {
			continue
		}
		phrase := words[i] + " " + words[i+1] + " " + words[i+2]
		phrases[phrase]++
	}

	// Filter phrases that appear only once
	filtered := make(map[string]int)
	for phrase, count := range phrases {
		if count > 1 {
			filtered[phrase] = count
		}
	}

	return filtered
}

// combineAndScore combines single words and phrases with scoring
func (ke *KeywordExtractor) combineAndScore(termFreq, phraseFreq map[string]int, totalWords int) map[string]float64 {
	scores := make(map[string]float64)

	// Score single terms
	for term, freq := range termFreq {
		// Simple TF scoring with slight length bonus for technical terms
		score := float64(freq)
		if len(term) > 8 {
			score *= 1.2 // Slight bonus for longer technical terms
		}
		scores[term] = score
	}

	// Score phrases (with bonus for being phrases)
	for phrase, freq := range phraseFreq {
		// Phrases get a multiplier since they're more specific
		scores[phrase] = float64(freq) * 2.0
	}

	return scores
}

// selectTopTerms selects the top N terms by score
func (ke *KeywordExtractor) selectTopTerms(scores map[string]float64, limit int) []string {
	// Create slice for sorting
	type termScore struct {
		term  string
		score float64
	}

	var sortedTerms []termScore
	for term, score := range scores {
		sortedTerms = append(sortedTerms, termScore{term, score})
	}

	// Sort by score descending
	sort.Slice(sortedTerms, func(i, j int) bool {
		return sortedTerms[i].score > sortedTerms[j].score
	})

	// Extract top N terms
	result := make([]string, 0, limit)
	for i := 0; i < limit && i < len(sortedTerms); i++ {
		result = append(result, sortedTerms[i].term)
	}

	return result
}

// isNumber checks if a string is purely numeric
func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// isCommonWord checks if a word is too common for phrases
func (ke *KeywordExtractor) isCommonWord(word string) bool {
	// Very common words that make phrases less meaningful
	veryCommon := map[string]bool{
		"get": true, "set": true, "use": true, "make": true,
		"can": true, "will": true, "may": true, "must": true,
		"new": true, "old": true, "good": true, "bad": true,
		"see": true, "go": true, "come": true, "take": true,
	}
	return veryCommon[word]
}

// initializeStopWords creates the stop words map
func initializeStopWords() map[string]bool {
	// Common English stop words
	stopWordsList := []string{
		// Articles
		"a", "an", "the",
		// Pronouns
		"i", "me", "my", "myself", "we", "our", "ours", "ourselves",
		"you", "your", "yours", "yourself", "yourselves",
		"he", "him", "his", "himself", "she", "her", "hers", "herself",
		"it", "its", "itself", "they", "them", "their", "theirs", "themselves",
		"what", "which", "who", "whom", "this", "that", "these", "those",
		// Prepositions
		"in", "on", "at", "by", "for", "with", "about", "against", "between",
		"into", "through", "during", "before", "after", "above", "below",
		"to", "from", "up", "down", "out", "off", "over", "under", "again",
		// Conjunctions
		"and", "or", "but", "if", "because", "as", "until", "while",
		"of", "at", "by", "for", "with", "about", "against", "between",
		// Common verbs
		"am", "is", "are", "was", "were", "be", "been", "being",
		"have", "has", "had", "having", "do", "does", "did", "doing",
		"could", "should", "would", "might", "shall", "can", "will", "may", "must",
		// Other common words
		"here", "there", "when", "where", "why", "how", "all", "both", "each",
		"few", "more", "most", "other", "some", "such", "no", "nor", "not",
		"only", "own", "same", "so", "than", "too", "very", "just", "also",
	}

	stopWords := make(map[string]bool, len(stopWordsList))
	for _, word := range stopWordsList {
		stopWords[word] = true
	}

	return stopWords
}

// ParallelExtractor handles concurrent keyword extraction
type ParallelExtractor struct {
	extractor *KeywordExtractor
	workers   int
}

// NewParallelExtractor creates a new parallel keyword extractor
func NewParallelExtractor(workers int) *ParallelExtractor {
	if workers <= 0 {
		workers = 4 // Default number of workers
	}
	return &ParallelExtractor{
		extractor: NewKeywordExtractor(),
		workers:   workers,
	}
}

// ExtractFromDocuments processes multiple documents in parallel
func (pe *ParallelExtractor) ExtractFromDocuments(docs []Document, keywordLimit int) []DocumentKeywords {
	// Create channels
	jobs := make(chan extractJob, len(docs))
	results := make(chan DocumentKeywords, len(docs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < pe.workers; i++ {
		wg.Add(1)
		go pe.worker(jobs, results, keywordLimit, &wg)
	}

	// Send jobs
	for i, doc := range docs {
		jobs <- extractJob{index: i, doc: doc}
	}
	close(jobs)

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	keywordResults := make([]DocumentKeywords, len(docs))
	for result := range results {
		keywordResults[result.Index] = result
	}

	return keywordResults
}

// worker processes extraction jobs
func (pe *ParallelExtractor) worker(jobs <-chan extractJob, results chan<- DocumentKeywords, limit int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		keywords := pe.extractor.ExtractKeywords(job.doc.Content, limit)
		results <- DocumentKeywords{
			Index:    job.index,
			URL:      job.doc.URL,
			Keywords: keywords,
		}
	}
}

// Document represents a document to process
type Document struct {
	URL     string
	Content string
}

// DocumentKeywords represents extracted keywords for a document
type DocumentKeywords struct {
	Index    int
	URL      string
	Keywords []string
}

// extractJob represents a job for parallel processing
type extractJob struct {
	index int
	doc   Document
}