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
    name: "Total Budget"
    type: LOOKUP
    lookupOption: {
      customFieldId: "related_projects_field_id"
      targetField: "budget"
      function: SUM
    }
    description: "Sum of budget from related projects"
  }) {
    id
    name
    type
    lookupOption
  }
}
```

## Advanced Example

Create a lookup field with filtering and formatting:

```graphql
mutation CreateAdvancedLookupField {
  createCustomField(input: {
    name: "Completed Tasks Count"
    type: LOOKUP
    lookupOption: {
      customFieldId: "dependencies_field_id"
      targetField: "status"
      function: COUNT
      filter: {
        status: COMPLETED
      }
      display: {
        type: NUMBER
        precision: 0
      }
    }
    description: "Count of completed dependency tasks"
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
| `customFieldId` | String! | ✅ Yes | ID of the reference field to pull data from |
| `targetField` | String! | ✅ Yes | Field name to aggregate from referenced records |
| `function` | LookupFunction! | ✅ Yes | Aggregation function to apply |
| `filter` | TodoFilterInput | No | Filter criteria for referenced records |
| `display` | LookupDisplayInput | No | Display formatting options |

## Aggregation Functions

### LookupFunction Types

| Function | Description | Use Cases | Returns |
|----------|-------------|-----------|---------|
| `SUM` | Sum of numeric values | Budgets, quantities, scores | Number |
| `AVERAGE` | Average of numeric values | Ratings, performance metrics | Number |
| `COUNT` | Count of records | Task counts, completion rates | Number |
| `MIN` | Minimum value | Earliest dates, lowest scores | Number/Date |
| `MAX` | Maximum value | Latest dates, highest scores | Number/Date |
| `CONCAT` | Concatenate text values | Names, descriptions | String |
| `FIRST` | First value found | Latest status, primary contact | Mixed |
| `LAST` | Last value found | Most recent update | Mixed |

### Function Examples

```graphql
# Sum all budget values from referenced records
{
  function: SUM
  targetField: "budget"
}

# Count completed tasks
{
  function: COUNT
  filter: { status: COMPLETED }
}

# Get average rating
{
  function: AVERAGE
  targetField: "rating"
}

# Find maximum due date
{
  function: MAX
  targetField: "dueDate"
}

# Concatenate all assignee names
{
  function: CONCAT
  targetField: "assignees.name"
}
```

## Target Field Options

### Built-in Todo Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | String | Record title |
| `description` | String | Record description |
| `status` | TodoStatus | Record status |
| `dueDate` | DateTime | Due date |
| `createdAt` | DateTime | Creation date |
| `updatedAt` | DateTime | Last update date |
| `assignees.name` | String | Assignee names |
| `assignees.email` | String | Assignee emails |
| `tags.name` | String | Tag names |

### Custom Field Values

Reference custom fields from referenced records:

```graphql
{
  targetField: "customFields.budget"  # Custom field named "budget"
  targetField: "customFields.priority"  # Custom field named "priority"
  targetField: "customFields.score"  # Custom field named "score"
}
```

## Display Configuration

### LookupDisplayInput

| Parameter | Type | Description |
|-----------|------|-------------|
| `type` | DisplayType | NUMBER, CURRENCY, PERCENTAGE, TEXT |
| `precision` | Int | Decimal places for numbers |
| `currency` | CurrencyInput | Currency configuration |
| `format` | String | Custom format string |

### Display Examples

```graphql
# Currency display
{
  display: {
    type: CURRENCY
    currency: {
      code: "USD"
      name: "US Dollar"
    }
    precision: 2
  }
}

# Percentage display
{
  display: {
    type: PERCENTAGE
    precision: 1
  }
}

# Number display
{
  display: {
    type: NUMBER
    precision: 0
  }
}
```

## Filtering Referenced Data

Use the `filter` parameter to limit which referenced records are included:

```graphql
{
  lookupOption: {
    customFieldId: "tasks_field_id"
    targetField: "customFields.effort"
    function: SUM
    filter: {
      status: ACTIVE
      assigneeIds: ["user_123"]
      tags: ["high-priority"]
      dueDateFrom: "2024-01-01"
      dueDateTo: "2024-12-31"
    }
  }
}
```

## Lookup Field Values

Lookup fields are read-only and automatically calculated. Values are stored in the `lookupResult` field:

```json
{
  "number": 15750.50,
  "lookupResult": {
    "value": 15750.50,
    "display": {
      "type": "CURRENCY",
      "currency": {
        "code": "USD",
        "name": "US Dollar"
      },
      "formatted": "$15,750.50"
    },
    "sourceCount": 12,
    "lastCalculated": "2024-03-15T10:30:00Z"
  }
}
```

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The lookup field definition |
| `number` | Float | Calculated numeric result |
| `text` | String | Calculated text result |
| `lookupResult` | JSON | Full result with metadata |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last calculated |

### LookupResult Structure

| Field | Type | Description |
|-------|------|-------------|
| `value` | Mixed | The calculated value |
| `display` | DisplayInfo | Formatted display information |
| `sourceCount` | Int | Number of records used in calculation |
| `lastCalculated` | DateTime | When the calculation was performed |
| `error` | String | Error message if calculation failed |

## Querying Lookup Data

### Basic Query

```graphql
query GetRecordsWithLookups {
  todos(projectId: "project_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
      }
      number
      text
      lookupResult
    }
  }
}
```

### Advanced Query with Metadata

```graphql
query GetDetailedLookups {
  todos(projectId: "project_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
        lookupOption {
          customFieldId
          targetField
          function
          filter
        }
      }
      lookupResult {
        value
        display {
          formatted
          type
        }
        sourceCount
        lastCalculated
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
    "message": "Reference field not found",
    "extensions": {
      "code": "REFERENCE_FIELD_NOT_FOUND"
    }
  }]
}
```

### Invalid Target Field

```json
{
  "errors": [{
    "message": "Target field not found in referenced records",
    "extensions": {
      "code": "TARGET_FIELD_NOT_FOUND"
    }
  }]
}
```

### Calculation Error

```json
{
  "errors": [{
    "message": "Lookup calculation failed",
    "extensions": {
      "code": "LOOKUP_CALCULATION_ERROR"
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