---
title: Lookup Custom Field
description: Create lookup fields that automatically pull and aggregate data from referenced records
category: Custom Fields
---

Lookup custom fields automatically pull data from records referenced by [Reference fields](/api/custom-fields/reference), enabling powerful data aggregation and computed values across projects. They update automatically when referenced data changes.

## Basic Example

Create a simple lookup field:

```graphql
mutation CreateLookupField {
  createCustomField(input: {
    name: "Related Todo Tags"
    type: LOOKUP
    lookupOption: {
      referenceId: "reference_field_id"
      lookupType: TODO_TAG
    }
    description: "Tags from related todos"
  }) {
    id
    name
    type
    lookupOption
  }
}
```

## Advanced Example

Create a lookup field for custom field data:

```graphql
mutation CreateAdvancedLookupField {
  createCustomField(input: {
    name: "Referenced Budget Data"
    type: LOOKUP
    lookupOption: {
      referenceId: "project_reference_field_id"
      lookupId: "budget_custom_field_id"
      lookupType: TODO_CUSTOM_FIELD
    }
    description: "Budget data from referenced project todos"
  }) {
    id
    name
    type
    lookupOption
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the lookup field |
| `type` | CustomFieldType! | ✅ Yes | Must be `LOOKUP` |
| `lookupOption` | CustomFieldLookupOptionInput! | ✅ Yes | Lookup configuration |
| `description` | String | No | Help text shown to users |

**Note**: Custom fields are automatically associated with the project based on the user's current project context.

## Lookup Configuration

### CustomFieldLookupOptionInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `referenceId` | String! | ✅ Yes | ID of the reference field to pull data from |
| `lookupId` | String | No | ID of the specific custom field to lookup (for TODO_CUSTOM_FIELD) |
| `lookupType` | CustomFieldLookupType! | ✅ Yes | Type of data to extract from referenced records |

## Lookup Types

### CustomFieldLookupType Values

| Type | Description | Returns | Use Cases |
|------|-------------|---------|-----------|
| `TODO_DUE_DATE` | Due dates from referenced todos | Date range | Project timelines, deadline tracking |
| `TODO_CREATED_AT` | Creation dates from referenced todos | Date range | Creation time analysis |
| `TODO_UPDATED_AT` | Last updated dates from referenced todos | Date range | Activity tracking |
| `TODO_TAG` | Tags from referenced todos | Array of tags | Tag aggregation, categorization |
| `TODO_ASSIGNEE` | Assignees from referenced todos | Array of users | Team member tracking |
| `TODO_DESCRIPTION` | Descriptions from referenced todos | Array of text | Content aggregation |
| `TODO_LIST` | Todo list names from referenced todos | Array of list names | List organization |
| `TODO_CUSTOM_FIELD` | Custom field values from referenced todos | Varies by field type | Cross-project data aggregation |

### Lookup Examples

```graphql
# Get tags from referenced todos
{
  lookupType: TODO_TAG
  referenceId: "reference_field_id"
}

# Get assignees from referenced todos
{
  lookupType: TODO_ASSIGNEE
  referenceId: "reference_field_id"
}

# Get due dates from referenced todos
{
  lookupType: TODO_DUE_DATE
  referenceId: "reference_field_id"
}

# Get custom field data from referenced todos
{
  lookupType: TODO_CUSTOM_FIELD
  referenceId: "reference_field_id"
  lookupId: "budget_custom_field_id"
}
```

## TODO_CUSTOM_FIELD Lookup

When using `TODO_CUSTOM_FIELD` type, the lookup extracts data from a specific custom field in the referenced todos. The `lookupId` parameter specifies which custom field to look up.

### Supported Custom Field Types

| Field Type | Returns | Example |
|------------|---------|---------|
| `CURRENCY` | `{number, currency}` objects | Budget amounts with currency |
| `NUMBER` | Numeric values | Scores, quantities |
| `UNIQUE_ID` | Formatted ID strings | Ticket numbers |
| `SELECT_SINGLE` | Option objects | Status values |
| `SELECT_MULTI` | Arrays of option objects | Category selections |
| `LOCATION` | Location data with grouped todos | Geographic data |
| `LOOKUP` | Nested lookup results | Chained lookups |

### Example Usage

```graphql
# Get budget amounts from referenced todos
{
  lookupType: TODO_CUSTOM_FIELD
  referenceId: "project_reference_field"
  lookupId: "budget_custom_field_id"
}
```

## Lookup Results

Lookup fields automatically calculate their values based on the referenced data. The results are accessible through the `CustomFieldLookupOption` type.

## Lookup Field Values

Lookup fields are read-only and automatically calculated. Values are accessed through the `CustomFieldLookupOption` type:

```json
{
  "customFieldLookupOption": {
    "lookupType": "TODO_TAG",
    "lookupResult": [
      {
        "id": "tag_123",
        "name": "urgent",
        "color": "#ff0000"
      },
      {
        "id": "tag_456", 
        "name": "development",
        "color": "#00ff00"
      }
    ]
  }
}
```

## Response Fields

### CustomField Response (for lookup fields)

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field |
| `name` | String! | Display name of the lookup field |
| `type` | CustomFieldType! | Will be `LOOKUP` |
| `customFieldLookupOption` | CustomFieldLookupOption | Lookup configuration and results |
| `createdAt` | DateTime! | When the field was created |
| `updatedAt` | DateTime! | When the field was last updated |

### CustomFieldLookupOption Structure

| Field | Type | Description |
|-------|------|-------------|
| `lookupType` | CustomFieldLookupType! | Type of lookup (TODO_TAG, TODO_ASSIGNEE, etc.) |
| `lookupResult` | JSON | The calculated lookup results |
| `reference` | CustomField | The reference field being looked up |
| `lookup` | CustomField | The specific field being looked up (for TODO_CUSTOM_FIELD) |
| `parentCustomField` | CustomField | The parent lookup field |

## Querying Lookup Data

### Basic Query

```graphql
query GetLookupFields {
  customFields(projectId: "project_123") {
    id
    name
    type
    customFieldLookupOption {
      lookupType
      lookupResult
      reference {
        id
        name
      }
      lookup {
        id
        name
      }
    }
  }
}
```

### Advanced Query with Results

```graphql
query GetDetailedLookups {
  customFields(projectId: "project_123") {
    id
    name
    type
    customFieldLookupOption {
      lookupType
      lookupResult
      reference {
        id
        name
        referenceProjectId
      }
      lookup {
        id
        name
        type
      }
      parentCustomField {
        id
        name
      }
    }
  }
}
```

## Automatic Updates

### Calculation Triggers

Lookup fields automatically recalculate when:
- Referenced records are modified
- Reference field values change
- New records are added to referenced projects
- Referenced records are deleted
- Filters match different records

### Update Process

1. **Change Detection** - System monitors referenced data
2. **Calculation Queue** - Updates are queued for processing
3. **Batch Processing** - Multiple updates are batched for efficiency
4. **Result Storage** - New values are stored and cached
5. **Notification** - Subscriptions notify of changes

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create lookup field | `OWNER` or `ADMIN` role at project level |
| Update lookup field | `OWNER` or `ADMIN` role at project level |
| View lookup results | Standard record view permissions |
| Access source data | View permissions on referenced project |

**Important**: Users must have view permissions on the referenced project to see lookup results.

## Error Responses

### Invalid Reference Field

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

### Invalid Lookup Configuration

```json
{
  "errors": [{
    "message": "Circular lookup detected",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

### Project Access Error

```json
{
  "errors": [{
    "message": "Project not found",
    "extensions": {
      "code": "PROJECT_NOT_FOUND"
    }
  }]
}
```

## Best Practices

### Field Design

1. **Clear naming** - Use descriptive names that indicate the calculation
2. **Appropriate functions** - Choose functions that match your data type
3. **Useful filters** - Apply filters to get meaningful results
4. **Proper formatting** - Use display options for user-friendly output

### Performance Optimization

1. **Limit scope** - Use filters to reduce calculation complexity
2. **Choose efficient functions** - COUNT is faster than SUM for large datasets
3. **Monitor dependencies** - Avoid complex chains of lookups
4. **Cache considerations** - Results are cached but recalculation takes time

### Data Quality

1. **Validate source data** - Ensure referenced fields contain expected data
2. **Handle null values** - Functions handle nulls differently
3. **Consider data types** - Match functions to appropriate data types
4. **Test calculations** - Verify results with sample data

## Common Use Cases

### Budget Tracking

```graphql
# Sum budget from all related project tasks
{
  name: "Total Project Budget"
  type: LOOKUP
  lookupOption: {
    customFieldId: "related_tasks_field"
    targetField: "customFields.budget"
    function: SUM
    display: {
      type: CURRENCY
      currency: { code: "USD" }
    }
  }
}
```

### Progress Monitoring

```graphql
# Count completed dependencies
{
  name: "Completed Dependencies"
  type: LOOKUP
  lookupOption: {
    customFieldId: "dependencies_field"
    targetField: "status"
    function: COUNT
    filter: { status: COMPLETED }
  }
}
```

### Quality Metrics

```graphql
# Average quality score from reviews
{
  name: "Average Quality Score"
  type: LOOKUP
  lookupOption: {
    customFieldId: "reviews_field"
    targetField: "customFields.quality_score"
    function: AVERAGE
    display: {
      type: NUMBER
      precision: 1
    }
  }
}
```

### Resource Utilization

```graphql
# Sum allocated hours from team members
{
  name: "Total Allocated Hours"
  type: LOOKUP
  lookupOption: {
    customFieldId: "team_members_field"
    targetField: "customFields.allocated_hours"
    function: SUM
    filter: {
      status: ACTIVE
      assigneeIds: ["team_lead_id"]
    }
  }
}
```

## Integration with References

Lookup fields require [Reference fields](/api/custom-fields/reference) to work:

```graphql
# Step 1: Create reference field
mutation CreateReference {
  createCustomField(input: {
    name: "Project Dependencies"
    type: REFERENCE
    referenceProjectId: "dependencies_project"
    referenceMultiple: true
  }) {
    id
  }
}

# Step 2: Create lookup field using the reference
mutation CreateLookup {
  createCustomField(input: {
    name: "Dependencies Budget"
    type: LOOKUP
    lookupOption: {
      customFieldId: "project_dependencies_field_id"
      targetField: "customFields.budget"
      function: SUM
    }
  }) {
    id
  }
}
```

## Limitations

- Lookup fields are read-only and cannot be directly edited
- Maximum 1000 referenced records can be processed per lookup
- Complex calculations may have performance impact
- Circular lookup dependencies are not allowed
- Results are cached and may not reflect real-time changes
- Some functions (like CONCAT) have string length limits

## Related Resources

- [Reference Fields](/api/custom-fields/reference) - Link to records for lookup source
- [Formula Fields](/api/custom-fields/formula) - Calculated fields within the same project
- [Number Fields](/api/custom-fields/number) - For static numeric values
- [Custom Fields Overview](/api/custom-fields) - General concepts