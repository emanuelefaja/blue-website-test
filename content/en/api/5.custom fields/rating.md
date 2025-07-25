---
title: Rating Custom Field
description: Create rating fields to store numeric ratings with configurable scales and validation
category: Custom Fields
---

Rating custom fields allow you to store numeric ratings in records with configurable minimum and maximum values. They're ideal for performance ratings, satisfaction scores, priority levels, or any numeric scale-based data in your projects.

## Basic Example

Create a simple rating field with default 0-5 scale:

```graphql
mutation CreateRatingField {
  createCustomField(input: {
    name: "Performance Rating"
    type: RATING
    projectId: "proj_123"
    max: 5
  }) {
    id
    name
    type
    min
    max
  }
}
```

## Advanced Example

Create a rating field with custom scale and description:

```graphql
mutation CreateDetailedRatingField {
  createCustomField(input: {
    name: "Customer Satisfaction"
    type: RATING
    projectId: "proj_123"
    description: "Rate customer satisfaction from 1-10"
    min: 1
    max: 10
  }) {
    id
    name
    type
    description
    min
    max
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the rating field |
| `type` | CustomFieldType! | ✅ Yes | Must be `RATING` |
| `description` | String | No | Help text shown to users |
| `min` | Float | No | Minimum rating value (defaults to 0) |
| `max` | Float | No | Maximum rating value |

## Setting Rating Values

To set or update a rating value on a record:

```graphql
mutation SetRatingValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    number: 4.5
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the rating custom field |
| `number` | Float! | ✅ Yes | Rating value within the configured range |

## Creating Records with Rating Values

When creating a new record with rating values:

```graphql
mutation CreateRecordWithRating {
  createTodo(input: {
    title: "Review customer feedback"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "rating_field_id"
      value: "4.5"
    }]
  }) {
    id
    title
    customFields {
      id
      customField {
        name
        type
        min
        max
      }
      value
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
| `value` | Float | The stored rating value |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

### CustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field |
| `name` | String! | Display name of the rating field |
| `type` | CustomFieldType! | Always `RATING` |
| `min` | Float | Minimum allowed rating value |
| `max` | Float | Maximum allowed rating value |
| `description` | String | Help text for the field |

## Rating Validation

### Value Constraints
- Rating values must be numeric (Float type)
- Values must be within the configured min/max range
- If no minimum is specified, defaults to 0
- Maximum value is optional but recommended

### Validation Rules
- Input is parsed as a float number
- Must be greater than or equal to the minimum value
- Must be less than or equal to the maximum value (if specified)
- Invalid numeric input throws validation error

### Valid Rating Examples
For a field with min=1, max=5:
```
1       # Minimum value
5       # Maximum value
3.5     # Decimal values allowed
2.75    # Precise decimal ratings
```

### Invalid Rating Examples
For a field with min=1, max=5:
```
0       # Below minimum
6       # Above maximum
-1      # Negative value (below min)
abc     # Non-numeric value
```

## Configuration Options

### Rating Scale Setup
```graphql
# 1-5 star rating
mutation CreateStarRating {
  createCustomField(input: {
    name: "Star Rating"
    type: RATING
    projectId: "proj_123"
    min: 1
    max: 5
  }) {
    id
    min
    max
  }
}

# 0-100 percentage rating
mutation CreatePercentageRating {
  createCustomField(input: {
    name: "Completion Percentage"
    type: RATING
    projectId: "proj_123"
    min: 0
    max: 100
  }) {
    id
    min
    max
  }
}
```

### Common Rating Scales
- **1-5 Stars**: `min: 1, max: 5`
- **0-10 NPS**: `min: 0, max: 10`
- **1-10 Performance**: `min: 1, max: 10`
- **0-100 Percentage**: `min: 0, max: 100`
- **Custom Scale**: Any numeric range

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create rating field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update rating field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set rating value | Standard record edit permissions |
| View rating value | Standard record view permissions |

## Error Responses

### Value Below Minimum
```json
{
  "errors": [{
    "message": "Rating must be greater than or equal to 1.",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

### Value Above Maximum
```json
{
  "errors": [{
    "message": "Rating must be less than or equal to 5.",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

### Invalid Numeric Value
```json
{
  "errors": [{
    "message": "Invalid rating value.",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
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

### Scale Design
- Use consistent rating scales across similar fields
- Consider user familiarity (1-5 stars, 0-10 NPS)
- Set appropriate minimum values (0 vs 1)
- Define clear meaning for each rating level

### Data Quality
- Validate rating values before storing
- Use decimal precision appropriately
- Consider rounding for display purposes
- Provide clear guidance on rating meanings

### User Experience
- Display rating scales visually (stars, progress bars)
- Show current value and scale limits
- Provide context for rating meanings
- Consider default values for new records

## Common Use Cases

1. **Performance Management**
   - Employee performance ratings
   - Project quality scores
   - Task completion ratings
   - Skill level assessments

2. **Customer Feedback**
   - Satisfaction ratings
   - Product quality scores
   - Service experience ratings
   - Net Promoter Score (NPS)

3. **Priority and Importance**
   - Task priority levels
   - Urgency ratings
   - Risk assessment scores
   - Impact ratings

4. **Quality Assurance**
   - Code review ratings
   - Testing quality scores
   - Documentation quality
   - Process adherence ratings

## Integration Features

### With Automations
- Trigger actions based on rating thresholds
- Send notifications for low ratings
- Create follow-up tasks for high ratings
- Route work based on rating values

### With Lookups
- Calculate average ratings across records
- Find records by rating ranges
- Reference rating data from other records
- Aggregate rating statistics

### With Forms
- Automatic range validation
- Visual rating input controls
- Real-time validation feedback
- Star or slider input options

## Activity Tracking

Rating field changes are automatically tracked:
- Old and new rating values are logged
- Activity shows numeric changes
- Timestamps for all rating updates
- User attribution for changes

## Limitations

- Only numeric values are supported
- No built-in visual rating display (stars, etc.)
- Decimal precision depends on database configuration
- No rating metadata storage (comments, context)
- No automatic rating aggregation or statistics
- No built-in rating conversion between scales

## Related Resources

- [Number Fields](/api/custom-fields/number) - For general numeric data
- [Percent Fields](/api/custom-fields/percent) - For percentage values
- [Select Fields](/api/custom-fields/select-single) - For discrete choice ratings
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts
- [Forms API](/api/forms) - For visual rating input controls