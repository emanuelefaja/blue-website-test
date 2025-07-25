---
title: Single-Select Custom Field
description: Create single-select fields to allow users to choose one option from a predefined list
category: Custom Fields
---

Single-select custom fields allow users to choose exactly one option from a predefined list. They're ideal for status fields, categories, priorities, or any scenario where only one choice should be made from a controlled set of options.

## Basic Example

Create a simple single-select field:

```graphql
mutation CreateSingleSelectField {
  createCustomField(input: {
    name: "Project Status"
    type: SELECT_SINGLE
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a single-select field with predefined options:

```graphql
mutation CreateDetailedSingleSelectField {
  createCustomField(input: {
    name: "Priority Level"
    type: SELECT_SINGLE
    projectId: "proj_123"
    description: "Set the priority level for this task"
    customFieldOptions: [
      { title: "Low", color: "#28a745" }
      { title: "Medium", color: "#ffc107" }
      { title: "High", color: "#fd7e14" }
      { title: "Critical", color: "#dc3545" }
    ]
  }) {
    id
    name
    type
    description
    customFieldOptions {
      id
      title
      color
      position
    }
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the single-select field |
| `type` | CustomFieldType! | ✅ Yes | Must be `SELECT_SINGLE` |
| `description` | String | No | Help text shown to users |
| `customFieldOptions` | [CreateCustomFieldOptionInput!] | No | Initial options for the field |

### CreateCustomFieldOptionInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `title` | String! | ✅ Yes | Display text for the option |
| `color` | String | No | Hex color code for the option |

## Adding Options to Existing Fields

Add new options to an existing single-select field:

```graphql
mutation AddSingleSelectOption {
  createCustomFieldOption(input: {
    customFieldId: "field_123"
    title: "Urgent"
    color: "#6f42c1"
  }) {
    id
    title
    color
    position
  }
}
```

## Setting Single-Select Values

To set the selected option on a record:

```graphql
mutation SetSingleSelectValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    customFieldOptionId: "option_789"
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the single-select custom field |
| `customFieldOptionId` | String! | ✅ Yes | ID of the option to select |

## Creating Records with Single-Select Values

When creating a new record with single-select values:

```graphql
mutation CreateRecordWithSingleSelect {
  createTodo(input: {
    title: "Review user feedback"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "priority_field_id"
      value: "option_high_priority"
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
      selectedOption {
        id
        title
        color
      }
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
| `selectedOption` | CustomFieldOption | The currently selected option (null if none) |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

### CustomFieldOption Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the option |
| `title` | String! | Display text for the option |
| `color` | String | Hex color code for visual representation |
| `position` | Float | Sort order for the option |
| `customField` | CustomField! | The custom field this option belongs to |

### CustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field |
| `name` | String! | Display name of the single-select field |
| `type` | CustomFieldType! | Always `SELECT_SINGLE` |
| `description` | String | Help text for the field |
| `customFieldOptions` | [CustomFieldOption!] | All available options |
| `selectedOption` | CustomFieldOption | Current selection for this record |

## Value Format

### Input Format
- **API Parameter**: Single option ID (`"option_123"`)
- **String Format**: Single option ID (`"option_123"`)
- **Multiple IDs**: If multiple IDs provided, only the first is used

### Output Format
- **GraphQL Response**: Single CustomFieldOption object
- **Activity Log**: Option title as string
- **Automation Data**: Option title as string

## Selection Behavior

### Exclusive Selection
- Setting a new option automatically removes the previous selection
- Only one option can be selected at a time
- Setting `null` or empty value clears the selection

### Fallback Logic
- If `customFieldOptionIds` array is provided, only the first option is used
- This ensures compatibility with multi-select input formats
- Empty arrays or null values clear the selection

## Managing Options

### Update Option Properties
```graphql
mutation UpdateOption {
  editCustomFieldOption(input: {
    id: "option_123"
    title: "Updated Priority"
    color: "#ff6b6b"
  }) {
    id
    title
    color
  }
}
```

### Delete Option
```graphql
mutation DeleteOption {
  deleteCustomFieldOption(id: "option_123")
}
```

**Note**: Deleting an option will clear it from all records where it was selected.

### Reorder Options
```graphql
mutation ReorderOptions {
  reorderCustomFieldOptions(input: {
    customFieldId: "field_123"
    optionIds: ["option_1", "option_3", "option_2"]
  }) {
    id
    position
  }
}
```

## Validation Rules

### Option Validation
- The provided option ID must exist
- Option must belong to the specified custom field
- Only one option can be selected (enforced automatically)
- Null/empty values are valid (no selection)

### Field Validation
- Must have at least one option defined to be usable
- Option titles must be unique within the field
- Color codes must be valid hex format (if provided)

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create single-select field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update single-select field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Add/edit options | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set selected value | Standard record edit permissions |
| View selected value | Standard record view permissions |

## Error Responses

### Invalid Option ID
```json
{
  "errors": [{
    "message": "Custom field option not found",
    "extensions": {
      "code": "NOT_FOUND"
    }
  }]
}
```

### Option Doesn't Belong to Field
```json
{
  "errors": [{
    "message": "Option does not belong to this custom field",
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

### Option Design
- Use clear, descriptive option titles
- Apply meaningful color coding
- Keep option lists focused and relevant
- Order options logically (by priority, frequency, etc.)

### Status Field Patterns
- Use consistent status workflows across projects
- Consider the natural progression of options
- Include clear final states (Done, Canceled, etc.)
- Use colors that reflect option meaning

### Data Management
- Review and clean up unused options periodically
- Use consistent naming conventions
- Consider the impact of option deletion on existing records
- Plan for option updates and migrations

## Common Use Cases

1. **Status and Workflow**
   - Task status (To Do, In Progress, Done)
   - Approval status (Pending, Approved, Rejected)
   - Project phase (Planning, Development, Testing, Released)
   - Issue resolution status

2. **Classification and Categorization**
   - Priority levels (Low, Medium, High, Critical)
   - Task types (Bug, Feature, Enhancement, Documentation)
   - Project categories (Internal, Client, Research)
   - Department assignments

3. **Quality and Assessment**
   - Review status (Not Started, In Review, Approved)
   - Quality ratings (Poor, Fair, Good, Excellent)
   - Risk levels (Low, Medium, High)
   - Confidence levels

4. **Assignment and Ownership**
   - Team assignments
   - Department ownership
   - Role-based assignments
   - Regional assignments

## Integration Features

### With Automations
- Trigger actions when specific options are selected
- Route work based on selected categories
- Send notifications for status changes
- Create conditional workflows based on selections

### With Lookups
- Filter records by selected options
- Reference option data from other records
- Create reports based on option selections
- Group records by selected values

### With Forms
- Dropdown input controls
- Radio button interfaces
- Option validation and filtering
- Conditional field display based on selections

## Activity Tracking

Single-select field changes are automatically tracked:
- Shows old and new option selections
- Displays option titles in activity log
- Timestamps for all selection changes
- User attribution for modifications

## Differences from Multi-Select

| Feature | Single-Select | Multi-Select |
|---------|---------------|--------------|
| **Selection Limit** | Exactly 1 option | Multiple options |
| **Input Parameter** | `customFieldOptionId` | `customFieldOptionIds` |
| **Response Field** | `selectedOption` | `selectedOptions` |
| **Storage Behavior** | Replaces existing selection | Adds to existing selections |
| **Common Use Cases** | Status, category, priority | Tags, skills, categories |

## Limitations

- Only one option can be selected at a time
- No hierarchical or nested option structure
- Options are shared across all records using the field
- No built-in option analytics or usage tracking
- Color codes are for display only, no functional impact
- Cannot set different permissions per option

## Related Resources

- [Multi-Select Fields](/api/custom-fields/select-multi) - For multiple-choice selections
- [Checkbox Fields](/api/custom-fields/checkbox) - For simple boolean choices
- [Text Fields](/api/custom-fields/text-single) - For free-form text input
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts
- [Custom Field Options](/api/custom-fields/options) - Managing field options