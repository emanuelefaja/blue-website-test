# Adding a New Language to Blue Website

This guide provides step-by-step instructions for adding a new language to the Blue website. The process involves updating configuration files, creating translation files, and running translation commands.

## Prerequisites

- OpenAI API key (for automated translations)
- Go 1.24.4 or higher installed
- Access to the codebase

## Step 1: Update Language Configuration

Edit `/web/languages.go` and add your new language code to the `SupportedLanguages` array:

```go
var SupportedLanguages = []string{
    "en", // English
    "zh", // 简体中文 (Simplified Chinese)
    // ... existing languages ...
    "ar", // العربية (Arabic) - NEW LANGUAGE
}
```

Also add the locale mapping in the same file:

```go
var LanguageLocales = map[string]string{
    "en":    "en_US",
    "zh":    "zh_CN",
    // ... existing mappings ...
    "ar":    "ar_SA", // NEW LANGUAGE
}
```

## Step 2: Create Content Directory Structure

Create the language directory structure in `/content`:

```bash
mkdir -p content/ar/api
mkdir -p content/ar/docs
mkdir -p content/ar/insights
mkdir -p content/ar/legal
```

## Step 3: Translate UI Text (JSON Files)

Set your OpenAI API key:

```bash
export OPENAI_API_KEY="your-api-key-here"
```

Run the JSON translation command to translate all UI text files:

```bash
cd /Users/manny/Blue/blue-new-website
go run cmd/translate-json/main.go ar
```

This will:
- Read all JSON files from `/translations` directory
- Translate each text entry to Arabic
- Create Arabic translations in the same JSON files
- Show progress and token usage

## Step 4: Translate Markdown Content

### 4.1 Translate Insights/Blog Posts

```bash
go run cmd/translate-insights/main.go ar
```

This will:
- Read all markdown files from `/content/en/insights/`
- Translate them to Arabic
- Save to `/content/ar/insights/`
- Preserve YAML frontmatter and markdown formatting

### 4.2 Translate API Documentation

```bash
go run cmd/translate-api-docs/main.go ar
```

This will:
- Translate all API documentation from `/content/en/api/`
- Save to `/content/ar/api/`

### 4.3 Translate General Documentation

```bash
go run cmd/translate-docs/main.go ar
```

This will:
- Translate all documentation from `/content/en/docs/`
- Save to `/content/ar/docs/`

### 4.4 Translate Changelog

```bash
go run cmd/translate-changelog/main.go ar
```

This will:
- Translate the changelog from `/content/en/docs/changelog.md`
- Save to `/content/ar/docs/changelog.md`

## Step 5: Update Language Names in Translation Commands

Each translation command has a `languageNames` map that should include your new language. Update these files:

1. `/cmd/translate-json/main.go`
2. `/cmd/translate-insights/main.go`
3. `/cmd/translate-api-docs/main.go`
4. `/cmd/translate-docs/main.go`
5. `/cmd/translate-changelog/main.go`

Add your language to the map:

```go
var languageNames = map[string]string{
    "en":    "English",
    "zh":    "Simplified Chinese",
    // ... existing languages ...
    "ar":    "Arabic", // NEW LANGUAGE
}
```

## Step 6: Verify Translations

After running all translation commands, verify:

1. **Check JSON translations**: Open any file in `/translations` and verify your language key exists
2. **Check content structure**: Ensure `/content/ar/` has the same structure as `/content/en/`
3. **Test the website**: Start the server and navigate to `/ar/` URLs

```bash
# Start the development server
air

# Visit http://localhost:8080/ar/
```

## Step 7: Check Translation Coverage

Run the translation coverage tool to see what's missing:

```bash
go run cmd/translation-coverage/main.go
```

This will show:
- Which JSON files have complete translations
- Which content files exist for each language
- Coverage percentage for each language

## Step 8: Manual Review and Adjustments

1. **Review critical pages**: Check important pages like homepage, pricing, and main features
2. **Fix RTL issues** (if applicable): For RTL languages like Arabic, you may need CSS adjustments
3. **Update language picker**: The language switcher should automatically show your new language
4. **Test language switching**: Ensure cookie-based language persistence works

## Additional Notes

### Translation Quality

The automated translations use OpenAI's GPT-4 model with specific prompts to:
- Maintain professional, business-appropriate tone
- Preserve technical terms and product names
- Keep markdown formatting intact
- Translate YAML frontmatter fields appropriately

### Partial Translations

If you only need to translate specific sections:
- You can manually create/edit individual files
- The system will fall back to showing translation keys for missing translations
- Use `{{t "key" "Default Text"}}` in templates to provide fallbacks

### Cost Estimation

Based on current token usage:
- UI translations (JSON files): ~50,000-100,000 tokens
- Insights/Blog posts: ~500,000-1,000,000 tokens
- Documentation: ~200,000-400,000 tokens
- Total: ~1-2 million tokens per language

### Maintenance

When adding new content:
1. Always add to English (`/content/en/`) first
2. Run the appropriate translation command for other languages
3. Or manually translate if preferred

### Troubleshooting

If translations fail:
- Check OpenAI API key is set correctly
- Verify rate limits aren't exceeded (commands use concurrent workers)
- Check console output for specific error messages
- Reduce worker count in the Go files if needed (default is 10)

## Summary Checklist

- [ ] Update `SupportedLanguages` in `/web/languages.go`
- [ ] Add locale mapping in `LanguageLocales`
- [ ] Create content directory structure
- [ ] Set OpenAI API key
- [ ] Run `translate-json` command
- [ ] Run `translate-insights` command
- [ ] Run `translate-api-docs` command
- [ ] Run `translate-docs` command
- [ ] Run `translate-changelog` command
- [ ] Update `languageNames` in all translation commands
- [ ] Test website with new language
- [ ] Review translations for accuracy
- [ ] Check translation coverage