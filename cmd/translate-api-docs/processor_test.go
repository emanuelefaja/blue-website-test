package main

import (
	"strings"
	"testing"
)

func TestDocumentProcessor(t *testing.T) {
	processor := NewDocumentProcessor()

	t.Run("Extract and restore code blocks", func(t *testing.T) {
		input := "# API Example\n\nHere's how to use it:\n\n```graphql\nmutation CreateProject {\n  createProject(input: { name: \"Test\" }) {\n    id\n  }\n}\n```\n\nAnd another example:\n\n```json\n{\n  \"data\": {\n    \"id\": \"123\"\n  }\n}\n```"

		masked, placeholders, err := processor.ProcessDocument(input)
		if err != nil {
			t.Fatalf("ProcessDocument failed: %v", err)
		}

		// Check that code blocks were replaced
		if strings.Contains(masked, "```") {
			t.Error("Code blocks were not properly masked")
		}

		// Check that placeholders were created
		if len(placeholders) < 2 {
			t.Errorf("Expected at least 2 placeholders, got %d", len(placeholders))
		}

		// Validate placeholder format
		for _, placeholder := range processor.extractPlaceholders(masked) {
			if !strings.HasPrefix(placeholder, "@@") || !strings.HasSuffix(placeholder, "@@") {
				t.Errorf("Invalid placeholder format: %s", placeholder)
			}
		}

		// Test restoration
		restored := processor.RestoreContent(masked, placeholders)
		if !strings.Contains(restored, "mutation CreateProject") {
			t.Error("Code block content was not properly restored")
		}
	})

	t.Run("Table intelligence", func(t *testing.T) {
		input := "## Parameters\n\n| Parameter | Type | Required | Description |\n|-----------|------|----------|-------------|\n| `todoId` | String! | ✅ Yes | The ID of the todo to update |\n| `title` | String | No | New title for the todo |"

		masked, _, err := processor.ProcessDocument(input)
		if err != nil {
			t.Fatalf("ProcessDocument failed: %v", err)
		}

		// Check that parameter names and types are masked
		if strings.Contains(masked, "`todoId`") || strings.Contains(masked, "String!") {
			t.Error("Technical identifiers in table were not properly masked")
		}

		// Check that descriptions are preserved for translation
		if !strings.Contains(masked, "The ID of the todo to update") {
			t.Error("Description text should be preserved for translation")
		}

		// Check that Yes/No values are preserved
		if !strings.Contains(masked, "Yes") {
			t.Error("Yes/No values should be preserved for translation")
		}
	})

	t.Run("Inline code preservation", func(t *testing.T) {
		input := "Use the `createTodo` mutation with `todoListId` parameter."

		masked, placeholders, err := processor.ProcessDocument(input)
		if err != nil {
			t.Fatalf("ProcessDocument failed: %v", err)
		}

		// Check that inline code is masked
		if strings.Contains(masked, "`createTodo`") || strings.Contains(masked, "`todoListId`") {
			t.Error("Inline code was not properly masked")
		}

		// Check placeholders were created
		inlineCodeCount := 0
		for _, p := range placeholders {
			if p.Type == PlaceholderInlineCode {
				inlineCodeCount++
			}
		}
		if inlineCodeCount != 2 {
			t.Errorf("Expected 2 inline code placeholders, got %d", inlineCodeCount)
		}
	})

	t.Run("Email preservation", func(t *testing.T) {
		input := "Contact us at support@blue.cc or sales@blue.cc for assistance."

		masked, placeholders, err := processor.ProcessDocument(input)
		if err != nil {
			t.Fatalf("ProcessDocument failed: %v", err)
		}

		// Check that emails are masked
		if strings.Contains(masked, "@blue.cc") {
			t.Error("Email addresses were not properly masked")
		}

		// Verify restoration
		restored := processor.RestoreContent(masked, placeholders)
		if !strings.Contains(restored, "support@blue.cc") || !strings.Contains(restored, "sales@blue.cc") {
			t.Error("Email addresses were not properly restored")
		}
	})

	t.Run("Callout blocks", func(t *testing.T) {
		input := "Some text\n::callout\n---\nicon: lightbulb\n---\nThis is important information.\n::"

		masked, placeholders, err := processor.ProcessDocument(input)
		if err != nil {
			t.Fatalf("ProcessDocument failed: %v", err)
		}

		// Check that callout structure is preserved
		calloutFound := false
		for _, p := range placeholders {
			if p.Type == PlaceholderCallout {
				calloutFound = true
			}
		}
		if !calloutFound {
			t.Error("Callout markers were not properly extracted")
		}

		// Check that the content "This is important information." is exposed for translation
		if !strings.Contains(masked, "This is important information.") {
			t.Error("Callout content should be exposed for translation")
		}
		
		// Check that callout markers are replaced with placeholders
		if strings.Contains(masked, "::callout") || strings.Contains(masked, "::") {
			t.Error("Raw callout markers should be replaced with placeholders")
		}
	})

	t.Run("Validation", func(t *testing.T) {
		original := "Test @@CB##123-456##CB@@ content @@CODE##789##CODE@@ here"
		
		// Test valid translation
		valid := "Prueba @@CB##123-456##CB@@ contenido @@CODE##789##CODE@@ aquí"
		if err := processor.ValidateTranslation(original, valid); err != nil {
			t.Errorf("Valid translation failed validation: %v", err)
		}

		// Test missing placeholder
		invalid1 := "Prueba @@CB##123-456##CB@@ contenido aquí"
		if err := processor.ValidateTranslation(original, invalid1); err == nil {
			t.Error("Missing placeholder should fail validation")
		}

		// Test modified placeholder
		invalid2 := "Prueba @@CB##123-456##CB@@ contenido @@CODIGO##789##CODIGO@@ aquí"
		if err := processor.ValidateTranslation(original, invalid2); err == nil {
			t.Error("Modified placeholder should fail validation")
		}
	})

	t.Run("Complex document", func(t *testing.T) {
		input := "---\ntitle: Create a Project\ndescription: How to create a new project using the API\n---\n\n## Create a Project\n\nUse the `createProject` mutation:\n\n```graphql\nmutation {\n  createProject(input: { name: \"New Project\" }) {\n    id\n    name\n  }\n}\n```\n\n### Parameters\n\n| Parameter | Type | Required | Description |\n|-----------|------|----------|-------------|\n| `name` | String! | ✅ Yes | The project name |\n| `companyId` | String! | ✅ Yes | Company ID |\n\n::callout\n---\nicon: warning\n---\nRemember to include authentication headers.\n::\n\nFor support, contact support@blue.cc."

		masked, placeholders, err := processor.ProcessDocument(input)
		if err != nil {
			t.Fatalf("ProcessDocument failed: %v", err)
		}

		// Verify all technical content is masked
		technicalTerms := []string{"```", "`name`", "`companyId`", "String!", "@blue.cc", "::callout"}
		for _, term := range technicalTerms {
			if strings.Contains(masked, term) {
				t.Errorf("Technical term '%s' was not properly masked", term)
			}
		}

		// Verify translatable content is preserved
		translatableContent := []string{"How to create a new project", "The project name", "Company ID"}
		for _, content := range translatableContent {
			if !strings.Contains(masked, content) {
				t.Errorf("Translatable content '%s' was not preserved", content)
			}
		}

		// Test full restoration
		restored := processor.RestoreContent(masked, placeholders)
		if restored != input {
			t.Errorf("Restored content does not match original.\nOriginal:\n%s\n\nRestored:\n%s", input, restored)
		}
	})
}

func TestPlaceholderRecovery(t *testing.T) {
	processor := NewDocumentProcessor()

	t.Run("Recover corrupted placeholders", func(t *testing.T) {
		original := "Test @@CB##550e8400-e29b-41d4-a716-446655440000##CB@@ content"
		corrupted := "Test @@BLOQUE_CODIGO##550e8400-e29b-41d4-a716-446655440000##BLOQUE_CODIGO@@ content"

		recovered := processor.RecoverPlaceholders(original, corrupted)
		if !strings.Contains(recovered, "@@CB##550e8400-e29b-41d4-a716-446655440000##CB@@") {
			t.Error("Failed to recover corrupted placeholder")
		}
	})
}