---
title: Button Custom Field
description: Create interactive button fields that trigger automations when clicked
category: Custom Fields
---

Button custom fields provide interactive UI elements that trigger automations when clicked. Unlike other custom field types that store data, button fields serve as action triggers to execute configured workflows.

## Basic Example

Create a simple button field that triggers an automation:

```graphql
mutation CreateButtonField {
  createCustomField(input: {
    name: "Send Invoice"
    type: BUTTON
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a button with confirmation requirements:

```graphql
mutation CreateButtonWithConfirmation {
  createCustomField(input: {
    name: "Delete All Attachments"
    type: BUTTON
    projectId: "proj_123"
    buttonType: "hardConfirmation"
    buttonConfirmText: "DELETE"
    description: "Permanently removes all attachments from this task"
  }) {
    id
    name
    type
    buttonType
    buttonConfirmText
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the button |
| `type` | CustomFieldType! | ✅ Yes | Must be `BUTTON` |
| `projectId` | String! | ✅ Yes | Project ID where the field will be created |
| `buttonType` | String | No | Confirmation behavior (see Button Types below) |
| `buttonConfirmText` | String | No | Text users must type for hard confirmation |
| `description` | String | No | Help text shown to users |
| `required` | Boolean | No | Whether the field is required (defaults to false) |
| `isActive` | Boolean | No | Whether the field is active (defaults to true) |

### Button Type Field

The `buttonType` field is a free-form string that can be used by UI clients to determine confirmation behavior. Common values include:

- `""` (empty) - No confirmation
- `"soft"` - Simple confirmation dialog
- `"hard"` - Require typing confirmation text

**Note**: These are UI hints only. The API does not validate or enforce specific values.

## Triggering Button Clicks

To trigger a button click and execute associated automations:

```graphql
mutation ClickButton {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
  })
}
```

### Click Input Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the task containing the button |
| `customFieldId` | String! | ✅ Yes | ID of the button custom field |

### Important: API Behavior

**All button clicks through the API execute immediately** regardless of any `buttonType` or `buttonConfirmText` settings. These fields are stored for UI clients to implement confirmation dialogs, but the API itself:

- Does not validate confirmation text
- Does not enforce any confirmation requirements
- Executes the button action immediately when called

Confirmation is purely a client-side UI safety feature.

### Example: Clicking Different Button Types

```graphql
# Button with no confirmation
mutation ClickSimpleButton {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "simple_button_id"
  })
}

# Button with soft confirmation (API call is the same!)
mutation ClickSoftConfirmButton {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "soft_confirm_button_id"
  })
}

# Button with hard confirmation (API call is still the same!)
mutation ClickHardConfirmButton {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "hard_confirm_button_id"
  })
}
```

All three mutations above will execute the button action immediately when called through the API, bypassing any confirmation requirements.

## Response Fields

### Custom Field Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the custom field |
| `name` | String! | Display name of the button |
| `type` | CustomFieldType! | Always `BUTTON` for button fields |
| `buttonType` | String | Confirmation behavior setting |
| `buttonConfirmText` | String | Required confirmation text (if using hard confirmation) |
| `description` | String | Help text for users |
| `required` | Boolean! | Whether the field is required |
| `isActive` | Boolean! | Whether the field is currently active |
| `projectId` | String! | ID of the project this field belongs to |
| `createdAt` | DateTime! | When the field was created |
| `updatedAt` | DateTime! | When the field was last modified |

## How Button Fields Work

### Automation Integration

Button fields are designed to work with Blue's automation system:

1. **Create the button field** using the mutation above
2. **Configure automations** that listen for `CUSTOM_FIELD_BUTTON_CLICKED` events
3. **Users click the button** in the UI
4. **Automations execute** the configured actions

### Event Flow

When a button is clicked:

```
User Click → setTodoCustomField mutation → CUSTOM_FIELD_BUTTON_CLICKED event → Automation execution
```

### No Data Storage

Important: Button fields don't store any value data. They purely serve as action triggers. Each click:
- Generates an event
- Triggers associated automations
- Records an action in the task history
- Does not modify any field value

## Required Permissions

Users need appropriate project roles to create and use button fields:

| Action | Required Role |
|--------|-------------------|
| Create button field | `OWNER` or `ADMIN` at project level |
| Update button field | `OWNER` or `ADMIN` at project level |
| Click button | `OWNER`, `ADMIN`, `MEMBER`, or `CLIENT` (based on field permissions) |
| Configure automations | `OWNER` or `ADMIN` at project level |

## Error Responses

### Permission Denied
```json
{
  "errors": [{
    "message": "You don't have permission to edit this custom field",
    "extensions": {
      "code": "FORBIDDEN"
    }
  }]
}
```

### Custom Field Not Found
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

**Note**: The API does not return specific errors for missing automations or confirmation mismatches.

## Best Practices

### Naming Conventions
- Use action-oriented names: "Send Invoice", "Create Report", "Notify Team"
- Be specific about what the button does
- Avoid generic names like "Button 1" or "Click Here"

### Confirmation Settings
- Leave `buttonType` empty for safe, reversible actions
- Set `buttonType` to suggest confirmation behavior to UI clients
- Use `buttonConfirmText` to specify what users should type in UI confirmations
- Remember: These are UI hints only - API calls always execute immediately

### Automation Design
- Keep button actions focused on a single workflow
- Provide clear feedback about what happened after clicking
- Consider adding description text to explain the button's purpose

## Common Use Cases

1. **Workflow Transitions**
   - "Mark as Complete"
   - "Send for Approval"
   - "Archive Task"

2. **External Integrations**
   - "Sync to CRM"
   - "Generate Invoice"
   - "Send Email Update"

3. **Batch Operations**
   - "Update All Subtasks"
   - "Copy to Projects"
   - "Apply Template"

4. **Reporting Actions**
   - "Generate Report"
   - "Export Data"
   - "Create Summary"

## Limitations

- Buttons cannot store or display data values
- Each button can only trigger automations, not direct API calls (however, automations can include HTTP request actions to call external APIs or Blue's own APIs)
- Button visibility cannot be conditionally controlled
- Maximum of one automation execution per click (though that automation can trigger multiple actions)

## Related Resources

- [Automations API](/api/automations/index) - Configure actions triggered by buttons
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General custom field concepts