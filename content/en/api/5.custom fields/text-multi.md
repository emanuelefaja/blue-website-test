---
title: Multi-Line Text Custom Field
description: Create multi-line text fields for longer content like descriptions, notes, and comments
category: Custom Fields
---

Multi-line text custom fields allow you to store longer text content with line breaks and formatting. They're ideal for descriptions, notes, comments, or any text data that needs multiple lines.

## Basic Example

Create a simple multi-line text field:

```graphql
mutation CreateTextMultiField {
  createCustomField(input: {
    name: "Description"
    type: TEXT_MULTI
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a multi-line text field with description:

```graphql
mutation CreateDetailedTextMultiField {
  createCustomField(input: {
    name: "Project Notes"
    type: TEXT_MULTI
    description: "Detailed notes and observations about the project"
  }) {
    id
    name
    type
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the text field |
| `type` | CustomFieldType! | ✅ Yes | Must be `TEXT_MULTI` |
| `description` | String | No | Help text shown to users |

**Note:** The project context is determined automatically from the `X-Bloo-Project-ID` header in your GraphQL request. Custom fields are always created within the project specified in the request header.

## Setting Text Values

To set or update a multi-line text value on a record:

```graphql
mutation SetTextMultiValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    text: "This is a multi-line text value.\n\nIt can contain line breaks and longer content."
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the text custom field |
| `text` | String | No | Multi-line text content to store |

## Creating Records with Text Values

When creating a new record with multi-line text values:

```graphql
mutation CreateRecordWithTextMulti {
  createTodo(input: {
    title: "Project Planning"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "text_multi_field_id"
      value: "Project Overview:\n\n1. Research phase\n2. Design phase\n3. Implementation phase\n\nKey considerations:\n- Budget constraints\n- Timeline requirements\n- Resource allocation"
    }]
  }) {
    id
    title
    customFields {
      id
      customField {
        name
        type
      }
      text
    }
  }
}
```

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `text` | String | The stored multi-line text content |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## Text Validation

### Form Validation
When multi-line text fields are used in forms:
- Leading and trailing whitespace is automatically trimmed
- Required validation is applied if the field is marked as required
- No specific format validation is applied

### Validation Rules
- Accepts any string content including line breaks
- No character length limits (up to database limits)
- Supports Unicode characters and special symbols
- Line breaks are preserved in storage

### Valid Text Examples
```
Single line text

Multi-line text with
line breaks

Text with special characters:
- Bullets
- Numbers: 123
- Symbols: @#$%
- Unicode: 🚀 ✅ ⭐

Code snippets:
function example() {
  return "hello world";
}
```

## Important Notes

### Storage Capacity
- Stored using MySQL `MediumText` type
- Supports up to 16MB of text content
- Line breaks and formatting are preserved
- UTF-8 encoding for international characters

### Direct API vs Forms
- **Forms**: Automatic whitespace trimming and required validation
- **Direct API**: Text is stored exactly as provided
- **Recommendation**: Use forms for user input to ensure consistent formatting

### TEXT_MULTI vs TEXT_SINGLE
- **TEXT_MULTI**: Multi-line textarea input, ideal for longer content
- **TEXT_SINGLE**: Single-line text input, ideal for short values
- **Backend**: Both use identical storage and validation
- **Frontend**: Different UI components for data entry

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create text field | `OWNER` or `ADMIN` project-level role |
| Update text field | `OWNER` or `ADMIN` project-level role |
| Set text value | Any role except `VIEW_ONLY` or `COMMENT_ONLY` |
| View text value | Any project-level role |

## Error Responses

### Required Field Validation (Forms Only)
```json
{
  "errors": [{
    "message": "This field is required",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

### Field Not Found
```json
{
  "errors": [{
    "message": "Custom field not found",
    "extensions": {
      "code": "CUSTOM_FIELD_NOT_FOUND"
    }
  }]
}
```

## Best Practices

### Content Organization
- Use consistent formatting for structured content
- Consider using markdown-like syntax for readability
- Break long content into logical sections
- Use line breaks to improve readability

### Data Entry
- Provide clear field descriptions to guide users
- Use forms for user input to ensure validation
- Consider character limits based on your use case
- Validate content format in your application if needed

### Performance Considerations
- Very long text content may affect query performance
- Consider pagination for displaying large text fields
- Index considerations for search functionality
- Monitor storage usage for fields with large content

## Filtering and Search

### Contains Search
Multi-line text fields support substring searching through the todos query filtering system:

```graphql
query SearchTextMulti {
  todos(
    # Custom field filtering is available through the query system
    # Exact parameter structure should be verified from current API schema
    where: {
      customFields: {
        some: {
          customFieldId: "text_multi_field_id"
          text: {
            contains: "project"
          }
        }
      }
    }
  ) {
    id
    title
    customFields {
      text
    }
  }
}
```

### Search Capabilities
- Substring matching within text fields
- Searches across all lines of text
- Supports partial word matching
- Search behavior depends on database collation settings

## Common Use Cases

1. **Project Management**
   - Task descriptions
   - Project requirements
   - Meeting notes
   - Status updates

2. **Customer Support**
   - Issue descriptions
   - Resolution notes
   - Customer feedback
   - Communication logs

3. **Content Management**
   - Article content
   - Product descriptions
   - User comments
   - Review details

4. **Documentation**
   - Process descriptions
   - Instructions
   - Guidelines
   - Reference materials

## Integration Features

### With Automations
- Trigger actions when text content changes
- Extract keywords from text content
- Create summaries or notifications
- Process text content with external services

### With Lookups
- Reference text data from other records
- Aggregate text content from multiple sources
- Find records by text content
- Display related text information

### With Forms
- Automatic whitespace trimming
- Required field validation
- Multi-line textarea UI
- Character count display (if configured)

## Limitations

- No built-in text formatting or rich text editing
- No automatic link detection or conversion
- No spell checking or grammar validation
- No built-in text analysis or processing
- No versioning or change tracking
- Limited search capabilities (no full-text search)
- No content compression for very large text

## Related Resources

- [Single-Line Text Fields](/api/custom-fields/text-single) - For short text values
- [Email Fields](/api/custom-fields/email) - For email addresses
- [URL Fields](/api/custom-fields/url) - For website addresses
- [Custom Fields Overview](/api/custom-fields/list-custom-fields) - General concepts
- [Forms API](/api/forms) - For validated text input