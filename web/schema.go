package web

import (
	"encoding/json"
	"html/template"
	"log"
	"strings"
	"time"
)

// SchemaService handles structured data generation for SEO
type SchemaService struct {
	metadata *Metadata
	siteURL  string
	language string
}

// NewSchemaService creates a new schema service instance
func NewSchemaService(metadata *Metadata, siteURL string) *SchemaService {
	return &SchemaService{
		metadata: metadata,
		siteURL:  strings.TrimSuffix(siteURL, "/"),
		language: DefaultLanguage, // Default to English
	}
}

// SetLanguage updates the language for schema generation
func (s *SchemaService) SetLanguage(lang string) {
	s.language = lang
}

// GenerateSchema generates appropriate schema based on page type
func (s *SchemaService) GenerateSchema(pageType, path string, frontmatter *Frontmatter) template.JS {
	var schemaData []map[string]interface{}
	
	// Always include organization schema
	schemaData = append(schemaData, s.generateOrganizationSchema())
	
	// Add page-specific schemas
	switch pageType {
	case "platform":
		schemaData = append(schemaData, s.generateSoftwareApplicationSchema())
	case "pricing":
		schemaData = append(schemaData, s.generateProductSchema())
	case "insights":
		if frontmatter != nil {
			schemaData = append(schemaData, s.generateArticleSchema(frontmatter, path))
		}
	case "faq":
		schemaData = append(schemaData, s.generateFAQSchema())
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(schemaData)
	if err != nil {
		log.Printf("Error marshaling schema data: %v", err)
		return template.JS("[]")
	}
	
	return template.JS(jsonData)
}

// GetPageType determines the page type for schema generation
func (s *SchemaService) GetPageType(path string) string {
	switch {
	case strings.HasPrefix(path, "/platform"):
		return "platform"
	case path == "/pricing":
		return "pricing"
	case strings.HasPrefix(path, "/insights/"):
		return "insights"
	case strings.Contains(path, "/faq") || path == "/resources/faq":
		return "faq"
	default:
		return "general"
	}
}

// generateOrganizationSchema creates organization structured data
func (s *SchemaService) generateOrganizationSchema() map[string]interface{} {
	schema := map[string]interface{}{
		"@context": "https://schema.org",
		"@type": "Organization",
		"name": "Blue",
		"url": s.siteURL,
		"logo": s.siteURL + "/logo/blue-logo.png",
		"description": Translate(s.language, "schema.org_description"),
		"sameAs": []string{
			"https://twitter.com/bluehq",
			"https://www.linkedin.com/company/blue-teamwork",
			"https://www.youtube.com/@HeyBlueTeam",
		},
		"contactPoint": map[string]interface{}{
			"@type": "ContactPoint",
			"contactType": "customer support",
			"email": "support@blue.cc",
		},
	}
	
	// Description is now handled by Translate function above
	
	return schema
}

// generateSoftwareApplicationSchema creates software application structured data
func (s *SchemaService) generateSoftwareApplicationSchema() map[string]interface{} {
	// Translate feature list
	features := []string{
		Translate(s.language, "schema.feature_process_management"),
		Translate(s.language, "schema.feature_team_collaboration"),
		Translate(s.language, "schema.feature_workflow_automation"),
		Translate(s.language, "schema.feature_custom_fields"),
		Translate(s.language, "schema.feature_api_access"),
		Translate(s.language, "schema.feature_realtime_updates"),
		Translate(s.language, "schema.feature_file_attachments"),
		Translate(s.language, "schema.feature_custom_permissions"),
	}
	
	return map[string]interface{}{
		"@context": "https://schema.org",
		"@type": "SoftwareApplication",
		"name": "Blue",
		"applicationCategory": Translate(s.language, "schema.software_category"),
		"operatingSystem": Translate(s.language, "schema.software_operating_system"),
		"description": Translate(s.language, "schema.software_description"),
		"screenshot": s.siteURL + "/product/dashboard-screenshot.png",
		"featureList": features,
		"offers": map[string]interface{}{
			"@type": "Offer",
			"price": "7.00",
			"priceCurrency": "USD",
			"priceSpecification": map[string]interface{}{
				"@type": "UnitPriceSpecification",
				"price": "7.00",
				"priceCurrency": "USD",
				"unitText": Translate(s.language, "schema.price_unit_text"),
			},
		},
		"aggregateRating": map[string]interface{}{
			"@type": "AggregateRating",
			"ratingValue": "4.8",
			"reviewCount": "156",
			"bestRating": "5",
			"worstRating": "1",
		},
	}
}

// generateArticleSchema creates article structured data
func (s *SchemaService) generateArticleSchema(frontmatter *Frontmatter, path string) map[string]interface{} {
	schema := map[string]interface{}{
		"@context": "https://schema.org",
		"@type": "Article",
		"headline": frontmatter.Title,
		"description": frontmatter.Description,
		"url": s.siteURL + path,
		"author": map[string]interface{}{
			"@type": "Person",
			"name": Translate(s.language, "schema.author_name"),
		},
		"publisher": map[string]interface{}{
			"@type": "Organization",
			"name": Translate(s.language, "schema.publisher_name"),
			"logo": map[string]interface{}{
				"@type": "ImageObject",
				"url": s.siteURL + "/logo/blue-logo.png",
			},
		},
	}
	
	// Add dates if available
	if frontmatter.Date != "" {
		// Parse date and format as ISO 8601
		if t, err := time.Parse("2006-01-02", frontmatter.Date); err == nil {
			isoDate := t.Format(time.RFC3339)
			schema["datePublished"] = isoDate
			schema["dateModified"] = isoDate
		}
	}
	
	// Add article image if available
	if frontmatter.Image != "" {
		schema["image"] = s.siteURL + frontmatter.Image
	} else {
		schema["image"] = s.siteURL + "/og.png"
	}
	
	// Add category as articleSection
	if frontmatter.Category != "" {
		schema["articleSection"] = frontmatter.Category
	}
	
	// Add tags as keywords
	if len(frontmatter.Tags) > 0 {
		schema["keywords"] = strings.Join(frontmatter.Tags, ", ")
	}
	
	return schema
}

// generateProductSchema creates product/offer structured data for pricing
func (s *SchemaService) generateProductSchema() map[string]interface{} {
	return map[string]interface{}{
		"@context": "https://schema.org",
		"@type": "Product",
		"name": Translate(s.language, "schema.product_name"),
		"description": Translate(s.language, "schema.product_description"),
		"brand": map[string]interface{}{
			"@type": "Brand",
			"name": "Blue",
		},
		"offers": map[string]interface{}{
			"@type": "Offer",
			"url": s.siteURL + "/pricing",
			"priceCurrency": "USD",
			"price": "7.00",
			"priceValidUntil": "2025-12-31",
			"availability": "https://schema.org/InStock",
			"priceSpecification": map[string]interface{}{
				"@type": "UnitPriceSpecification",
				"price": "7.00",
				"priceCurrency": "USD",
				"unitText": Translate(s.language, "schema.price_unit_text"),
				"billingIncrement": 1,
				"billingDuration": "P1M",
			},
		},
		"aggregateRating": map[string]interface{}{
			"@type": "AggregateRating",
			"ratingValue": "4.8",
			"reviewCount": "156",
		},
	}
}

// generateFAQSchema creates FAQ structured data
func (s *SchemaService) generateFAQSchema() map[string]interface{} {
	// Build FAQ items using translations
	faqKeys := []struct {
		questionKey string
		answerKey   string
	}{
		{"schema.faq_what_is_blue_q", "schema.faq_what_is_blue_a"},
		{"schema.faq_how_different_q", "schema.faq_how_different_a"},
		{"schema.faq_cost_q", "schema.faq_cost_a"},
		{"schema.faq_trial_q", "schema.faq_trial_a"},
		{"schema.faq_integrations_q", "schema.faq_integrations_a"},
		{"schema.faq_mobile_q", "schema.faq_mobile_a"},
		{"schema.faq_gdpr_q", "schema.faq_gdpr_a"},
		{"schema.faq_encryption_q", "schema.faq_encryption_a"},
		{"schema.faq_import_q", "schema.faq_import_a"},
		{"schema.faq_support_q", "schema.faq_support_a"},
	}
	
	faqs := make([]map[string]interface{}, 0, len(faqKeys))
	for _, faq := range faqKeys {
		faqs = append(faqs, map[string]interface{}{
			"@type": "Question",
			"name": Translate(s.language, faq.questionKey),
			"acceptedAnswer": map[string]interface{}{
				"@type": "Answer",
				"text": Translate(s.language, faq.answerKey),
			},
		})
	}
	
	return map[string]interface{}{
		"@context": "https://schema.org",
		"@type": "FAQPage",
		"mainEntity": faqs,
	}
}