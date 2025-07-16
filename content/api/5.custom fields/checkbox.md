---
title: Checkbox Custom Field
description: Create boolean checkbox fields for yes/no or true/false data
category: Custom Fields
---

Checkbox custom fields provide a simple boolean (true/false) input for tasks. They're perfect for binary choices, status indicators, or tracking whether something has been completed.

## Basic Example

Create a simple checkbox field:

```graphql
mutation CreateCheckboxField {
  createCustomField(input: {
    name: "Reviewed"
    type: CHECKBOX
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a checkbox field with description and validation:

```graphql
mutation CreateDetailedCheckbox {
  createCustomField(input: {
    name: "Customer Approved"
    type: CHECKBOX
    description: "Check this box when the customer has approved the work"
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
| `name` | String! | ✅ Yes | Display name of the checkbox |
| `type` | CustomFieldType! | ✅ Yes | Must be `CHECKBOX` |
| `description` | String | No | Help text shown to users |

## Setting Checkbox Values

To set or update a checkbox value on a task:

```graphql
mutation CheckTheBox {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    checked: true
  })
}
```

To uncheck a checkbox:

```graphql
mutation UncheckTheBox {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    checked: false
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the task to update |
| `customFieldId` | String! | ✅ Yes | ID of the checkbox custom field |
| `checked` | Boolean | No | True to check, false to uncheck |

## Creating Tasks with Checkbox Values

When creating a new task with checkbox values:

```graphql
mutation CreateTaskWithCheckbox {
  createTodo(input: {
    title: "Review contract"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "checkbox_field_id"
      value: "true"  # Pass as string
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
      checked
    }
  }
}
```

### Accepted String Values

When creating tasks, checkbox values must be passed as strings:

| String Value | Result |
|--------------|---------|
| `"true"` | ✅ Checked |
| `"1"` | ✅ Checked |
| `"checked"` | ✅ Checked |
| Any other value | ❌ Unchecked |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | ID! | Unique identifier for the field value |
| `uid` | String! | Alternative unique identifier |
| `customField` | CustomField! | The custom field definition |
| `checked` | Boolean | The checkbox state (true/false/null) |
| `todo` | Todo! | The task this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## Automation Integration

Checkbox fields trigger different automation events based on state changes:

| Action | Event Triggered | Description |
|--------|----------------|-------------|
| Check (false → true) | `CUSTOM_FIELD_ADDED` | Triggered when checkbox is checked |
| Uncheck (true → false) | `CUSTOM_FIELD_REMOVED` | Triggered when checkbox is unchecked |

This allows you to create automations that respond to checkbox state changes, such as:
- Sending notifications when items are approved
- Moving tasks when review checkboxes are checked
- Updating related fields based on checkbox states

## Data Import/Export

### Importing Checkbox Values

When importing data via CSV or other formats:
- `"true"`, `"yes"`, `"1"` → Checked
- `"false"`, `"no"`, `"0"`, empty → Unchecked

### Exporting Checkbox Values

When exporting data:
- Checked boxes export as `"X"`
- Unchecked boxes export as empty string `""`

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create checkbox field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update checkbox field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set checkbox value | Standard task edit permissions |
| View checkbox value | Standard task view permissions |

## Error Responses

### Invalid Value Type
```json
{
  "errors": [{
    "message": "Invalid value type for checkbox field",
    "extensions": {
      "code": "INVALID_FIELD_VALUE"
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

### Naming Conventions
- Use clear, action-oriented names: "Approved", "Reviewed", "Is Complete"
- Avoid negative names that confuse users: prefer "Is Active" over "Is Inactive"
- Be specific about what the checkbox represents

### When to Use Checkboxes
- **Binary choices**: Yes/No, True/False, Done/Not Done
- **Status indicators**: Approved, Reviewed, Published
- **Feature flags**: Has Priority Support, Requires Signature
- **Simple tracking**: Email Sent, Invoice Paid, Item Shipped

### When NOT to Use Checkboxes
- When you need more than two options (use SELECT_SINGLE instead)
- For numeric or text data (use NUMBER or TEXT fields)
- When you need to track who checked it or when (use audit logs)

## Common Use Cases

1. **Approval Workflows**
   - "Manager Approved"
   - "Client Sign-off"
   - "Legal Review Complete"

2. **Task Management**
   - "Is Blocked"
   - "Ready for Review"
   - "High Priority"

3. **Quality Control**
   - "QA Passed"
   - "Documentation Complete"
   - "Tests Written"

4. **Administrative Flags**
   - "Invoice Sent"
   - "Contract Signed"
   - "Follow-up Required"

## Limitations

- Checkbox fields can only store true/false values (no tri-state or null after initial set)
- No default value configuration (always starts as null until set)
- Cannot store additional metadata like who checked it or when
- No conditional visibility based on other field values

## Related Resources

- [Custom Fields Overview](/custom-fields/list-custom-fields) - General custom field concepts
- [Automations API](/api/automations/index) - Create automations triggered by checkbox changes
- [Forms API](/api/forms) - Include checkboxes in custom forms