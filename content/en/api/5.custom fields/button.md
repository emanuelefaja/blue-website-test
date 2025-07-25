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

### Button Types

| Value | Description |
|-------|-------------|
| `noConfirmation` | Button clicks immediately without any confirmation (default) |
| `softConfirmation` | Shows a simple confirmation dialog before executing |
| `hardConfirmation` | Requires users to type specific text before executing |

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

### Important: Confirmation Handling

**The API call is identical for all button types** - whether the button has no confirmation, soft confirmation, or hard confirmation. The confirmation logic is handled entirely by the Blue web interface:

- **No confirmation**: Button clicks execute immediately
- **Soft confirmation**: Web UI shows a confirmation dialog
- **Hard confirmation**: Web UI requires typing specific text

When using the API directly, **all button clicks execute immediately** regardless of the button's confirmation settings. The API does not validate or enforce confirmations - this is purely a UI safety feature.

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

Users need appropriate permissions to create and use button fields:

| Action | Required Permission |
|--------|-------------------|
| Create button field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update button field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Click button | Standard task permissions for the containing task |
| Configure automations | `AUTOMATIONS_CREATE` at company or project level |

## Error Responses

### No Automations Configured
```json
{
  "errors": [{
    "message": "Button has no automations configured",
    "extensions": {
      "code": "NO_AUTOMATIONS"
    }
  }]
}
```

### Invalid Confirmation Text
```json
{
  "errors": [{
    "message": "Confirmation text does not match",
    "extensions": {
      "code": "INVALID_CONFIRMATION"
    }
  }]
}
```

### Permission Denied
```json
{
  "errors": [{
    "message": "You don't have permission to use this button",
    "extensions": {
      "code": "FORBIDDEN"
    }
  }]
}
```

## Best Practices

### Naming Conventions
- Use action-oriented names: "Send Invoice", "Create Report", "Notify Team"
- Be specific about what the button does
- Avoid generic names like "Button 1" or "Click Here"

### Confirmation Settings
- Use `noConfirmation` for safe, reversible actions
- Use `softConfirmation` for important but recoverable actions
- Use `hardConfirmation` for destructive or irreversible actions

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