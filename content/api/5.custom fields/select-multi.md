---
title: Multi-Select Custom Field
description: Create multi-select fields to allow users to choose multiple options from a predefined list
category: Custom Fields
---

Multi-select custom fields allow users to choose multiple options from a predefined list. They're ideal for categories, tags, skills, features, or any scenario where multiple selections are needed from a controlled set of options.

## Basic Example

Create a simple multi-select field:

```graphql
mutation CreateMultiSelectField {
  createCustomField(input: {
    name: "Project Categories"
    type: SELECT_MULTI
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a multi-select field with predefined options:

```graphql
mutation CreateDetailedMultiSelectField {
  createCustomField(input: {
    name: "Required Skills"
    type: SELECT_MULTI
    projectId: "proj_123"
    description: "Select all skills required for this task"
    customFieldOptions: [
      { title: "JavaScript", color: "#f7df1e" }
      { title: "React", color: "#61dafb" }
      { title: "Node.js", color: "#339933" }
      { title: "GraphQL", color: "#e10098" }
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
| `name` | String! | ✅ Yes | Display name of the multi-select field |
| `type` | CustomFieldType! | ✅ Yes | Must be `SELECT_MULTI` |
| `description` | String | No | Help text shown to users |
| `customFieldOptions` | [CreateCustomFieldOptionInput!] | No | Initial options for the field |

### CreateCustomFieldOptionInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `title` | String! | ✅ Yes | Display text for the option |
| `color` | String | No | Hex color code for the option |

## Adding Options to Existing Fields

Add new options to an existing multi-select field:

```graphql
mutation AddMultiSelectOption {
  createCustomFieldOption(input: {
    customFieldId: "field_123"
    title: "Python"
    color: "#3776ab"
  }) {
    id
    title
    color
    position
  }
}
```

## Setting Multi-Select Values

To set multiple selected options on a record:

```graphql
mutation SetMultiSelectValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    customFieldOptionIds: ["option_1", "option_2", "option_3"]
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the multi-select custom field |
| `customFieldOptionIds` | [String!] | ✅ Yes | Array of option IDs to select |

## Creating Records with Multi-Select Values

When creating a new record with multi-select values:

```graphql
mutation CreateRecordWithMultiSelect {
  createTodo(input: {
    title: "Develop new feature"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "skills_field_id"
      value: "option1,option2,option3"
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
      selectedOptions {
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
| `selectedOptions` | [CustomFieldOption!] | Array of selected options |
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
| `name` | String! | Display name of the multi-select field |
| `type` | CustomFieldType! | Always `SELECT_MULTI` |
| `description` | String | Help text for the field |
| `customFieldOptions` | [CustomFieldOption!] | All available options |

## Value Format

### Input Format
- **API Parameter**: Array of option IDs (`["option1", "option2", "option3"]`)
- **String Format**: Comma-separated option IDs (`"option1,option2,option3"`)

### Output Format
- **GraphQL Response**: Array of CustomFieldOption objects
- **Activity Log**: Comma-separated option titles
- **Automation Data**: Array of option titles

## Managing Options

### Update Option Properties
```graphql
mutation UpdateOption {
  editCustomFieldOption(input: {
    id: "option_123"
    title: "Updated Title"
    color: "#ff0000"
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
- All provided option IDs must exist
- Options must belong to the specified custom field
- Only SELECT_MULTI fields can have multiple options selected
- Empty array is valid (no selections)

### Field Validation
- Must have at least one option defined to be usable
- Option titles must be unique within the field
- Color codes must be valid hex format (if provided)

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create multi-select field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update multi-select field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Add/edit options | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set selected values | Standard record edit permissions |
| View selected values | Standard record view permissions |

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

### Multiple Options on Non-Multi Field
```json
{
  "errors": [{
    "message": "custom fields can only have one option",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

## Best Practices

### Option Design
- Use descriptive, concise option titles
- Apply consistent color coding schemes
- Keep option lists manageable (typically 3-20 options)
- Order options logically (alphabetically, by frequency, etc.)

### Data Management
- Review and clean up unused options periodically
- Use consistent naming conventions across projects
- Consider option reusability when creating fields
- Plan for option updates and migrations

### User Experience
- Provide clear field descriptions
- Use colors to improve visual distinction
- Group related options together
- Consider default selections for common cases

## Common Use Cases

1. **Project Management**
   - Task categories and tags
   - Priority levels and types
   - Team member assignments
   - Status indicators

2. **Content Management**
   - Article categories and topics
   - Content types and formats
   - Publication channels
   - Approval workflows

3. **Customer Support**
   - Issue categories and types
   - Affected products or services
   - Resolution methods
   - Customer segments

4. **Product Development**
   - Feature categories
   - Technical requirements
   - Testing environments
   - Release channels

## Integration Features

### With Automations
- Trigger actions when specific options are selected
- Route work based on selected categories
- Send notifications for high-priority selections
- Create follow-up tasks based on option combinations

### With Lookups
- Filter records by selected options
- Aggregate data across option selections
- Reference option data from other records
- Create reports based on option combinations

### With Forms
- Multi-select input controls
- Option validation and filtering
- Dynamic option loading
- Conditional field display

## Activity Tracking

Multi-select field changes are automatically tracked:
- Shows added and removed options
- Displays option titles in activity log
- Timestamps for all selection changes
- User attribution for modifications

## Limitations

- Maximum practical limit of options depends on UI performance
- No hierarchical or nested option structure
- Options are shared across all records using the field
- No built-in option analytics or usage tracking
- Color codes are for display only, no functional impact
- Cannot set different permissions per option

## Related Resources

- [Single-Select Fields](/api/custom-fields/select-single) - For single-choice selections
- [Checkbox Fields](/api/custom-fields/checkbox) - For simple boolean choices
- [Text Fields](/api/custom-fields/text-single) - For free-form text input
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts
- [Custom Field Options](/api/custom-fields/options) - Managing field options