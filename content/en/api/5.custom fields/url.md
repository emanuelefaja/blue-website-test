---
title: URL Custom Field
description: Create URL fields to store website addresses and links
category: Custom Fields
---

URL custom fields allow you to store website addresses and links in your records. They're ideal for tracking project websites, reference links, documentation URLs, or any web-based resources related to your work.

## Basic Example

Create a simple URL field:

```graphql
mutation CreateUrlField {
  createCustomField(input: {
    name: "Project Website"
    type: URL
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a URL field with description:

```graphql
mutation CreateDetailedUrlField {
  createCustomField(input: {
    name: "Reference Link"
    type: URL
    projectId: "proj_123"
    description: "Link to external documentation or resources"
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
| `name` | String! | ✅ Yes | Display name of the URL field |
| `type` | CustomFieldType! | ✅ Yes | Must be `URL` |
| `description` | String | No | Help text shown to users |

## Setting URL Values

To set or update a URL value on a record:

```graphql
mutation SetUrlValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    text: "https://example.com/documentation"
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the URL custom field |
| `text` | String! | ✅ Yes | URL address to store |

## Creating Records with URL Values

When creating a new record with URL values:

```graphql
mutation CreateRecordWithUrl {
  createTodo(input: {
    title: "Review documentation"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "url_field_id"
      value: "https://docs.example.com/api"
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
| `text` | String | The stored URL address |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## URL Validation

### Current Implementation
- **Direct API**: No URL format validation is currently enforced
- **Forms**: URL validation is planned but not currently active
- **Storage**: Any string value can be stored in URL fields

### Planned Validation
Future versions will include:
- HTTP/HTTPS protocol validation
- Valid URL format checking
- Domain name validation
- Automatic protocol prefix addition

### Recommended URL Formats
While not currently enforced, use these standard formats:

```
https://example.com
https://www.example.com
https://subdomain.example.com
https://example.com/path
https://example.com/path?param=value
http://localhost:3000
https://docs.example.com/api/v1
```

## Important Notes

### Storage Format
- URLs are stored as plain text without modification
- No automatic protocol addition (http://, https://)
- Case sensitivity preserved as entered
- No URL encoding/decoding performed

### Direct API vs Forms
- **Forms**: Planned URL validation (not currently active)
- **Direct API**: No validation - any text can be stored
- **Recommendation**: Validate URLs in your application before storing

### URL vs Text Fields
- **URL**: Semantically intended for web addresses
- **TEXT_SINGLE**: General single-line text
- **Backend**: Currently identical storage and validation
- **Frontend**: Different UI components for data entry

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create URL field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update URL field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set URL value | Standard record edit permissions |
| View URL value | Standard record view permissions |

## Error Responses

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

## Best Practices

### URL Format Standards
- Always include protocol (http:// or https://)
- Use HTTPS when possible for security
- Test URLs before storing to ensure they're accessible
- Consider using shortened URLs for display purposes

### Data Quality
- Validate URLs in your application before storing
- Check for common typos (missing protocols, incorrect domains)
- Standardize URL formats across your organization
- Consider URL accessibility and availability

### Security Considerations
- Be cautious with user-provided URLs
- Validate domains if restricting to specific sites
- Consider URL scanning for malicious content
- Use HTTPS URLs when handling sensitive data

## Filtering and Search

### Contains Search
URL fields support substring searching:

```graphql
query SearchUrls {
  todos(
    customFieldFilters: [{
      customFieldId: "url_field_id"
      operation: CONTAINS
      value: "docs.example.com"
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
- Partial domain matching
- Path and parameter searching
- No protocol-specific filtering

## Common Use Cases

1. **Project Management**
   - Project websites
   - Documentation links
   - Repository URLs
   - Demo sites

2. **Content Management**
   - Reference materials
   - Source links
   - Media resources
   - External articles

3. **Customer Support**
   - Customer websites
   - Support documentation
   - Knowledge base articles
   - Video tutorials

4. **Sales & Marketing**
   - Company websites
   - Product pages
   - Marketing materials
   - Social media profiles

## Integration Features

### With Lookups
- Reference URLs from other records
- Find records by domain or URL pattern
- Display related web resources
- Aggregate links from multiple sources

### With Forms
- URL-specific input components
- Planned validation for proper URL format
- Link preview capabilities (frontend)
- Clickable URL display

### With Reporting
- Track URL usage and patterns
- Monitor broken or inaccessible links
- Categorize by domain or protocol
- Export URL lists for analysis

## Limitations

### Current Limitations
- No active URL format validation
- No automatic protocol addition
- No link verification or accessibility checking
- No URL shortening or expansion
- No favicon or preview generation

### Automation Restrictions
- Not available as automation trigger fields
- Cannot be used in automation field updates
- Can be referenced in automation conditions
- Available in email templates and webhooks

### General Constraints
- No built-in link preview functionality
- No automatic URL shortening
- No click tracking or analytics
- No URL expiration checking
- No malicious URL scanning

## Future Enhancements

### Planned Features
- HTTP/HTTPS protocol validation
- Custom regex validation patterns
- Automatic protocol prefix addition
- URL accessibility checking

### Potential Improvements
- Link preview generation
- Favicon display
- URL shortening integration
- Click tracking capabilities
- Broken link detection

## Related Resources

- [Text Fields](/api/custom-fields/text-single) - For non-URL text data
- [Email Fields](/api/custom-fields/email) - For email addresses
- [Custom Fields Overview](/api/custom-fields/list-custom-fields) - General concepts
- [Forms API](/api/forms) - For URL input validation

## Migration from Text Fields

If you're migrating from text fields to URL fields:

1. **Create URL field** with the same name and configuration
2. **Export existing text values** to verify they're valid URLs
3. **Update records** to use the new URL field
4. **Delete old text field** after successful migration
5. **Update applications** to use URL-specific UI components

### Migration Example
```graphql
# Step 1: Create URL field
mutation CreateUrlField {
  createCustomField(input: {
    name: "Website Link"
    type: URL
    projectId: "proj_123"
  }) {
    id
  }
}

# Step 2: Update records (repeat for each record)
mutation MigrateToUrlField {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "new_url_field_id"
    text: "https://example.com"  # Value from old text field
  })
}
```