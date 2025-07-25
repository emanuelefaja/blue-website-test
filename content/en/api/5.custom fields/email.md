---
title: Email Custom Field
description: Create email fields to store and validate email addresses
category: Custom Fields
---

Email custom fields allow you to store email addresses in records with built-in validation. They're ideal for tracking contact information, assignee emails, or any email-related data in your projects.

## Basic Example

Create a simple email field:

```graphql
mutation CreateEmailField {
  createCustomField(input: {
    name: "Contact Email"
    type: EMAIL
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create an email field with description:

```graphql
mutation CreateDetailedEmailField {
  createCustomField(input: {
    name: "Client Email"
    type: EMAIL
    description: "Primary email address for client communications"
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
| `name` | String! | ✅ Yes | Display name of the email field |
| `type` | CustomFieldType! | ✅ Yes | Must be `EMAIL` |
| `description` | String | No | Help text shown to users |

## Setting Email Values

To set or update an email value on a record:

```graphql
mutation SetEmailValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    text: "john.doe@example.com"
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the email custom field |
| `text` | String | No | Email address to store |

## Creating Records with Email Values

When creating a new record with email values:

```graphql
mutation CreateRecordWithEmail {
  createTodo(input: {
    title: "Follow up with client"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "email_field_id"
      value: "client@company.com"
    }]
  }) {
    id
    title
    customFields {
      id
      name
      type
      text
    }
  }
}
```

## Response Fields

### CustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the custom field |
| `name` | String! | Display name of the email field |
| `type` | CustomFieldType! | The field type (EMAIL) |
| `description` | String | Help text for the field |
| `text` | String | The stored email address value |
| `createdAt` | DateTime! | When the field was created |
| `updatedAt` | DateTime! | When the field was last modified |

## Email Validation

### Form Validation
When email fields are used in forms, they automatically validate the email format:
- Uses standard email validation rules
- Trims whitespace from input
- Rejects invalid email formats

### Validation Rules
- Must contain an `@` symbol
- Must have a valid domain format
- Leading/trailing whitespace is automatically removed
- Common email formats are accepted

### Valid Email Examples
```
user@example.com
john.doe@company.co.uk
test+tag@domain.org
first.last@sub.domain.com
```

### Invalid Email Examples
```
plainaddress          # Missing @ symbol
@domain.com          # Missing local part
user@                # Missing domain
user@domain          # Missing TLD
user name@domain.com # Spaces not allowed
```

## Important Notes

### Direct API vs Forms
- **Forms**: Automatic email validation is applied
- **Direct API**: No validation - any text can be stored
- **Recommendation**: Use forms for user input to ensure validation

### Storage Format
- Email addresses are stored as plain text
- No special formatting or parsing
- Case sensitivity: EMAIL custom fields are stored case-sensitively (unlike user authentication emails which are normalized to lowercase)
- No maximum length limitations beyond database constraints (16MB limit)

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create email field | `OWNER` or `ADMIN` project-level role |
| Update email field | `OWNER` or `ADMIN` project-level role |
| Delete email field | `OWNER` or `ADMIN` project-level role |
| Set email value | Any role except `VIEW_ONLY` and `COMMENT_ONLY` |
| View email value | Any project role with field access |

## Error Responses

### Invalid Email Format (Forms Only)
```json
{
  "errors": [{
    "message": "ValidationError",
    "extensions": {
      "code": "BAD_USER_INPUT",
      "data": {
        "errors": [{
          "field": "email",
          "message": "Email format is invalid"
        }]
      }
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

### Data Entry
- Always validate email addresses in your application
- Use email fields only for actual email addresses
- Consider using forms for user input to get automatic validation

### Data Quality
- Trim whitespace before storing
- Consider case normalization (typically lowercase)
- Validate email format before important operations

### Privacy Considerations
- Email addresses are stored as plain text
- Consider data privacy regulations (GDPR, CCPA)
- Implement appropriate access controls

## Common Use Cases

1. **Contact Management**
   - Client email addresses
   - Vendor contact information
   - Team member emails
   - Support contact details

2. **Project Management**
   - Stakeholder emails
   - Approval contact emails
   - Notification recipients
   - External collaborator emails

3. **Customer Support**
   - Customer email addresses
   - Support ticket contacts
   - Escalation contacts
   - Feedback email addresses

4. **Sales & Marketing**
   - Lead email addresses
   - Campaign contact lists
   - Partner contact information
   - Referral source emails

## Integration Features

### With Automations
- Trigger actions when email fields are updated
- Send notifications to stored email addresses
- Create follow-up tasks based on email changes

### With Lookups
- Reference email data from other records
- Aggregate email lists from multiple sources
- Find records by email address

### With Forms
- Automatic email validation
- Email format checking
- Whitespace trimming

## Limitations

- No built-in email verification or validation beyond format checking
- No email-specific UI features (like clickable email links)
- Stored as plain text without encryption
- No email composition or sending capabilities
- No email metadata storage (display name, etc.)
- Direct API calls bypass validation (only forms validate)

## Related Resources

- [Text Fields](/api/custom-fields/text-single) - For non-email text data
- [URL Fields](/api/custom-fields/url) - For website addresses
- [Phone Fields](/api/custom-fields/phone) - For phone numbers
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts
- [Forms API](/api/forms) - For validated email input