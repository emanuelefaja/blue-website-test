package web

import (
	"os"
	"path/filepath"
	"strings"
)

// ContentService handles content file discovery and management
type ContentService struct {
	contentDir string
}

// NewContentService creates a new content service
func NewContentService(contentDir string) *ContentService {
	return &ContentService{
		contentDir: contentDir,
	}
}

// FindMarkdownFile searches for a markdown file matching the given path
func (cs *ContentService) FindMarkdownFile(path string) (string, error) {
	// Default to English
	return cs.FindMarkdownFileForLang(path, DefaultLanguage)
}

// FindMarkdownFileForLang searches for a language-specific markdown file
func (cs *ContentService) FindMarkdownFileForLang(path, lang string) (string, error) {
	// Convert URL path to potential file paths
	cleanPath := strings.Trim(path, "/")

	// Check if this is a content type path
	contentType, isContentPath := GetContentTypeFromPath(cleanPath)
	if isContentPath {
		return cs.findNumberedMarkdownFileForLang(cleanPath, contentType, lang)
	}

	// Try simple patterns for non-content paths with language
	patterns := []string{
		filepath.Join(cs.contentDir, lang, cleanPath+".md"),
		filepath.Join(cs.contentDir, lang, cleanPath, "index.md"),
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern, nil
		}
	}

	// If not found and not English, fall back to English
	if lang != DefaultLanguage {
		return cs.FindMarkdownFileForLang(path, DefaultLanguage)
	}

	return "", os.ErrNotExist
}

// findNumberedMarkdownFileForLang handles finding language-specific files with numeric prefixes
func (cs *ContentService) findNumberedMarkdownFileForLang(cleanPath string, contentType ContentType, lang string) (string, error) {
	parts := strings.Split(cleanPath, "/")
	if len(parts) < 2 {
		return "", os.ErrNotExist
	}

	// Build path progressively, finding numbered directories/files
	// Update BaseDir to include language
	currentPath := filepath.Join(cs.contentDir, lang, strings.TrimPrefix(contentType.BaseDir, "content/"))
	for i := 1; i < len(parts); i++ {
		cleanSegment := parts[i]

		if i == len(parts)-1 {
			// Last segment - look for numbered file
			filePath, err := cs.findNumberedFile(currentPath, cleanSegment)
			if err == nil {
				return filePath, nil
			}

			// Also try as directory with index
			dirPath, err := cs.findNumberedDirectory(currentPath, cleanSegment)
			if err == nil {
				indexPath := filepath.Join(dirPath, "index.md")
				if _, err := os.Stat(indexPath); err == nil {
					return indexPath, nil
				}
			}
		} else {
			// Intermediate segment - look for numbered directory
			dirPath, err := cs.findNumberedDirectory(currentPath, cleanSegment)
			if err != nil {
				return "", os.ErrNotExist
			}
			currentPath = dirPath
		}
	}

	// If not found and not English, fall back to English
	if lang != DefaultLanguage {
		return cs.findNumberedMarkdownFileForLang(cleanPath, contentType, DefaultLanguage)
	}

	return "", os.ErrNotExist
}

// findNumberedFile finds a file with numeric prefix matching the segment
func (cs *ContentService) findNumberedFile(dir, segment string) (string, error) {
	// Generate file patterns
	patterns := GenerateFilePatterns(segment, ".md")

	// Try each pattern
	for _, pattern := range patterns {
		glob := filepath.Join(dir, pattern)
		matches, err := filepath.Glob(glob)
		if err == nil && len(matches) > 0 {
			return matches[0], nil
		}

		// If no matches, try case-insensitive search
		if match, err := cs.findFileIgnoreCase(dir, pattern); err == nil {
			return match, nil
		}
	}

	return "", os.ErrNotExist
}

// findNumberedDirectory finds a directory with numeric prefix matching the segment
func (cs *ContentService) findNumberedDirectory(dir, segment string) (string, error) {
	// Generate directory patterns
	patterns := GenerateFilePatterns(segment, "")

	for _, pattern := range patterns {
		glob := filepath.Join(dir, pattern)
		matches, err := filepath.Glob(glob)
		if err == nil && len(matches) > 0 {
			// Check if it's a directory
			if info, err := os.Stat(matches[0]); err == nil && info.IsDir() {
				return matches[0], nil
			}
		}
	}

	return "", os.ErrNotExist
}

// findFileIgnoreCase performs case-insensitive file matching
func (cs *ContentService) findFileIgnoreCase(dir, pattern string) (string, error) {
	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Extract the pattern without the wildcard
	patternName := strings.TrimPrefix(pattern, "*")
	patternLower := strings.ToLower(patternName)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// Check if file matches pattern (case-insensitive)
		if strings.HasSuffix(strings.ToLower(fileName), patternLower) {
			// Also check if it has a numeric prefix (to match our numbered file pattern)
			parts := strings.SplitN(fileName, ".", 2)
			if len(parts) == 2 {
				// Check if first part starts with a number
				if len(parts[0]) > 0 && parts[0][0] >= '0' && parts[0][0] <= '9' {
					return filepath.Join(dir, fileName), nil
				}
			}
		}
	}

	return "", os.ErrNotExist
}
