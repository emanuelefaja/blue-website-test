package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	// Keys to remove from roadmap translations
	keysToRemove := []string{
		"multi_homing_records",
		"relative_date_actions",
		"relative_date_triggers",
		"copy_automation",
		"custom_date_automations",
		"owner_project_access",
		"n8n_integration",
		"scheduled_automations",
		"reporting_module",
	}

	// Read the roadmap.json file
	filePath := "translations/roadmap.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse JSON into a map structure
	var translations map[string]any
	if err := json.Unmarshal(data, &translations); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Track what we remove
	removedCount := 0

	// Iterate through each language
	for lang, content := range translations {
		// Get the items map
		if langMap, ok := content.(map[string]any); ok {
			if items, ok := langMap["items"].(map[string]any); ok {
				// Remove each unwanted key
				for _, key := range keysToRemove {
					if _, exists := items[key]; exists {
						delete(items, key)
						removedCount++
						fmt.Printf("Removed '%s' from language '%s'\n", key, lang)
					}
				}
			}
		}
	}

	// Marshal back to JSON with proper formatting
	output, err := json.MarshalIndent(translations, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Write back to file
	if err := os.WriteFile(filePath, output, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("\n✅ Successfully removed %d translation keys from roadmap.json\n", removedCount)
	fmt.Printf("📝 Removed keys: %v\n", keysToRemove)
}
