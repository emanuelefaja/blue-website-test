---
title: Single-Line Text Custom Field
description: Create single-line text fields for short text values like names, titles, and labels
category: Custom Fields
---

Single-line text custom fields allow you to store short text values intended for single-line input. They're ideal for names, titles, labels, or any text data that should be displayed on a single line.

## Basic Example

Create a simple single-line text field:

```graphql
mutation CreateTextSingleField {
  createCustomField(input: {
    name: "Client Name"
    type: TEXT_SINGLE
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a single-line text field with description:

```graphql
mutation CreateDetailedTextSingleField {
  createCustomField(input: {
    name: "Product SKU"
    type: TEXT_SINGLE
    projectId: "proj_123"
    description: "Unique product identifier code"
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
| `type` | CustomFieldType! | ✅ Yes | Must be `TEXT_SINGLE` |
| `description` | String | No | Help text shown to users |

## Setting Text Values

To set or update a single-line text value on a record:

```graphql
mutation SetTextSingleValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    text: "ABC-123-XYZ"
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the text custom field |
| `text` | String! | ✅ Yes | Single-line text content to store |

## Creating Records with Text Values

When creating a new record with single-line text values:

```graphql
mutation CreateRecordWithTextSingle {
  createTodo(input: {
    title: "Process Order"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "text_single_field_id"
      value: "ORD-2024-001"
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
| `text` | String | The stored single-line text content |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## Text Validation

### Form Validation
When single-line text fields are used in forms:
- Leading and trailing whitespace is automatically trimmed
- Required validation is applied if the field is marked as required
- No specific format validation is applied

### Validation Rules
- Accepts any string content including line breaks (though not recommended)
- No character length limits (up to database limits)
- Supports Unicode characters and special symbols
- Line breaks are preserved but not intended for this field type

### Typical Text Examples
```
Product Name
SKU-123-ABC
Client Reference #2024-001
Version 1.2.3
Project Alpha
Status: Active
```

## Important Notes

### Storage Capacity
- Stored using MySQL `MediumText` type
- Supports up to 16MB of text content
- Identical storage to multi-line text fields
- UTF-8 encoding for international characters

### Direct API vs Forms
- **Forms**: Automatic whitespace trimming and required validation
- **Direct API**: Text is stored exactly as provided
- **Recommendation**: Use forms for user input to ensure consistent formatting

### TEXT_SINGLE vs TEXT_MULTI
- **TEXT_SINGLE**: Single-line text input, ideal for short values
- **TEXT_MULTI**: Multi-line textarea input, ideal for longer content
- **Backend**: Both use identical storage and validation
- **Frontend**: Different UI components for data entry
- **Intent**: TEXT_SINGLE is semantically meant for single-line values

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create text field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update text field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set text value | Standard record edit permissions |
| View text value | Standard record view permissions |

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
      "code": "NOT_FOUND"
    }
  }]
}
```

## Best Practices

### Content Guidelines
- Keep text concise and single-line appropriate
- Avoid line breaks for intended single-line display
- Use consistent formatting for similar data types
- Consider character limits based on your UI requirements

### Data Entry
- Provide clear field descriptions to guide users
- Use forms for user input to ensure validation
- Validate content format in your application if needed
- Consider using dropdowns for standardized values

### Performance Considerations
- Single-line text fields are lightweight and performant
- Consider indexing for frequently searched fields
- Use appropriate display widths in your UI
- Monitor content length for display purposes

## Filtering and Search

### Contains Search
Single-line text fields support substring searching:

```graphql
query SearchTextSingle {
  todos(
    customFieldFilters: [{
      customFieldId: "text_single_field_id"
      operation: CONTAINS
      value: "SKU"
    }]
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
- Case-insensitive substring matching
- Supports partial word matching
- Exact value matching
- No full-text search or ranking

## Common Use Cases

1. **Identifiers and Codes**
   - Product SKUs
   - Order numbers
   - Reference codes
   - Version numbers

2. **Names and Titles**
   - Client names
   - Project titles
   - Product names
   - Category labels

3. **Short Descriptions**
   - Brief summaries
   - Status labels
   - Priority indicators
   - Classification tags

4. **External References**
   - Ticket numbers
   - Invoice references
   - External system IDs
   - Document numbers

## Integration Features

### With Lookups
- Reference text data from other records
- Find records by text content
- Display related text information
- Aggregate text values from multiple sources

### With Forms
- Automatic whitespace trimming
- Required field validation
- Single-line text input UI
- Character limit display (if configured)

### With Imports/Exports
- Direct CSV column mapping
- Automatic text value assignment
- Bulk data import support
- Export to spreadsheet formats

## Limitations

### Automation Restrictions
- Not directly available as automation trigger fields
- Cannot be used in automation field updates
- Can be referenced in automation conditions
- Available in email templates and webhooks

### General Limitations
- No built-in text formatting or styling
- No automatic validation beyond required fields
- No built-in uniqueness enforcement
- No content compression for very large text
- No versioning or change tracking
- Limited search capabilities (no full-text search)

## Related Resources

- [Multi-Line Text Fields](/api/custom-fields/text-multi) - For longer text content
- [Email Fields](/api/custom-fields/email) - For email addresses
- [URL Fields](/api/custom-fields/url) - For website addresses
- [Unique ID Fields](/api/custom-fields/unique-id) - For auto-generated identifiers
- [Custom Fields Overview](/api/custom-fields/list-custom-fields) - General concepts
- [Forms API](/api/forms) - For validated text input