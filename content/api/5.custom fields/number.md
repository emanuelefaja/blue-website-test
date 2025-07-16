---
title: Number Custom Field
description: Create number fields to store numeric values with optional min/max constraints and prefix formatting
category: Custom Fields
---

Number custom fields allow you to store numeric values for records. They support validation constraints, decimal precision, and can be used for quantities, scores, measurements, or any numeric data that doesn't require special formatting.

## Basic Example

Create a simple number field:

```graphql
mutation CreateNumberField {
  createCustomField(input: {
    name: "Priority Score"
    type: NUMBER
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a number field with constraints and prefix:

```graphql
mutation CreateConstrainedNumberField {
  createCustomField(input: {
    name: "Team Size"
    type: NUMBER
    min: 1
    max: 100
    prefix: "#"
    description: "Number of team members assigned to this project"
  }) {
    id
    name
    type
    min
    max
    prefix
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the number field |
| `type` | CustomFieldType! | ✅ Yes | Must be `NUMBER` |
| `min` | Float | No | Minimum allowed value |
| `max` | Float | No | Maximum allowed value |
| `prefix` | String | No | Display prefix (e.g., "#", "~", "$") |
| `description` | String | No | Help text shown to users |

## Setting Number Values

Number fields store decimal values with optional validation:

### Simple Number Value

```graphql
mutation SetNumberValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    number: 42.5
  })
}
```

### Integer Value

```graphql
mutation SetIntegerValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    number: 100
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the number custom field |
| `number` | Float | No | Numeric value to store |

## Value Validation

### Min/Max Constraints

When creating a field with constraints:

```graphql
mutation CreateValidatedField {
  createCustomField(input: {
    name: "Rating"
    type: NUMBER
    min: 1
    max: 10
    description: "Rating from 1 to 10"
  }) {
    id
    name
    min
    max
  }
}
```

### Supported Value Types

| Type | Example | Description |
|------|---------|-------------|
| Integer | `42` | Whole numbers |
| Decimal | `42.5` | Numbers with decimal places |
| Negative | `-10` | Negative values (if no min constraint) |
| Zero | `0` | Zero value |

**Note**: Values are validated against `min` and `max` constraints if they are set on the field.

## Creating Records with Number Values

When creating a new record with number values:

```graphql
mutation CreateRecordWithNumber {
  createTodo(input: {
    title: "Performance Review"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "score_field_id"
      value: "85.5"
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
        prefix
      }
      number
      value
    }
  }
}
```

### Supported Input Formats

| Format | Example | Result |
|--------|---------|---------|
| String number | `"42.5"` | Parsed to Float: 42.5 |
| Integer string | `"100"` | Parsed to Float: 100.0 |
| Decimal string | `"3.14159"` | Parsed to Float: 3.14159 |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `number` | Float | The numeric value |
| `value` | Object | Combined value object (see below) |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

### CustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field definition |
| `name` | String! | Display name of the field |
| `type` | CustomFieldType! | Always `NUMBER` |
| `min` | Float | Minimum allowed value |
| `max` | Float | Maximum allowed value |
| `prefix` | String | Display prefix |
| `description` | String | Help text |

### Value Object Structure

```json
{
  "value": {
    "number": 42.5
  }
}
```

**Note**: If the number value is not set, the value will be `null`.

## Filtering and Querying

Number fields support comprehensive numeric filtering:

```graphql
query FilterByNumberRange {
  todos(filter: {
    customFields: [{
      customFieldId: "score_field_id"
      numberRange: {
        min: 80
        max: 100
      }
      operator: BETWEEN
    }]
  }) {
    id
    title
    customFields {
      number
    }
  }
}
```

### Supported Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `EQ` | Equal to | `number = 42` |
| `NE` | Not equal to | `number ≠ 42` |
| `GT` | Greater than | `number > 42` |
| `GTE` | Greater than or equal | `number ≥ 42` |
| `LT` | Less than | `number < 42` |
| `LTE` | Less than or equal | `number ≤ 42` |
| `BETWEEN` | Within range | `42 ≤ number ≤ 100` |
| `NULL` | No value set | `number is null` |
| `NOT_NULL` | Has any value | `number is not null` |

### Range Filtering

```graphql
query FilterByRange {
  todos(filter: {
    customFields: [{
      customFieldId: "priority_field_id"
      operator: GTE
      number: 5
    }]
  }) {
    id
    title
  }
}
```

## Aggregation and Analysis

Number fields support aggregation functions:

```graphql
query AggregateNumbers {
  todos(filter: {
    customFields: [{
      customFieldId: "score_field_id"
      operator: NOT_NULL
    }]
  }) {
    aggregates {
      sum
      average
      min
      max
      count
    }
  }
}
```

### Supported Aggregations

| Function | Description |
|----------|-------------|
| `SUM` | Total of all values |
| `AVERAGE` | Mean of all values |
| `MIN` | Smallest value |
| `MAX` | Largest value |
| `COUNT` | Number of records with values |

## Display Formatting

### With Prefix

If a prefix is set, it will be displayed:

| Value | Prefix | Display |
|-------|--------|---------|
| `42` | `"#"` | `#42` |
| `100` | `"~"` | `~100` |
| `3.14` | `"π"` | `π3.14` |

### Decimal Precision

Numbers maintain their decimal precision:

| Input | Stored | Displayed |
|-------|--------|-----------|
| `42` | `42.0` | `42` |
| `42.5` | `42.5` | `42.5` |
| `42.123` | `42.123` | `42.123` |

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create number field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update number field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set number value | Standard record edit permissions |
| View number value | Standard record view permissions |
| Use aggregation functions | Standard record query permissions |

## Error Responses

### Invalid Number Format
```json
{
  "errors": [{
    "message": "Invalid number format",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

### Value Below Minimum
```json
{
  "errors": [{
    "message": "Number must be greater than or equal to 1",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

### Value Above Maximum
```json
{
  "errors": [{
    "message": "Number must be less than or equal to 100",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

### Not a Number
```json
{
  "errors": [{
    "message": "Value is not a valid number",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

## Best Practices

### Constraint Design
- Set realistic min/max values based on your use case
- Use constraints to prevent data entry errors
- Consider if negative values are valid

### Value Precision
- Use appropriate decimal precision for your needs
- Consider rounding for display purposes
- Be consistent with precision across related fields

### Display Enhancement
- Use meaningful prefixes for context
- Consider units in field names (e.g., "Weight (kg)")
- Provide clear descriptions for validation rules

## Common Use Cases

1. **Scoring Systems**
   - Performance ratings
   - Quality scores
   - Priority levels
   - Customer satisfaction ratings

2. **Measurements**
   - Quantities and amounts
   - Dimensions and sizes
   - Durations (in numeric format)
   - Capacities and limits

3. **Business Metrics**
   - Revenue figures
   - Conversion rates
   - Budget allocations
   - Target numbers

4. **Technical Data**
   - Version numbers
   - Configuration values
   - Performance metrics
   - Threshold settings

## Integration Features

### With Formulas
- Reference NUMBER fields in calculations
- Perform mathematical operations
- Create derived metrics

### With Automations
- Trigger actions based on number thresholds
- Update related fields based on number changes
- Send notifications for specific values

### With Lookups
- Aggregate numbers from related records
- Calculate totals and averages
- Find min/max values across relationships

### With Charts
- Create numerical visualizations
- Track trends over time
- Compare values across records

## Limitations

- No built-in currency formatting (use CURRENCY type instead)
- No automatic percentage symbol (use PERCENT type instead)
- No unit conversion capabilities
- Decimal precision limited by Float type
- No mathematical formula evaluation in the field itself

## Related Resources

- [Custom Fields Overview](/custom-fields/list-custom-fields) - General custom field concepts
- [Currency Custom Field](/custom-fields/currency) - For monetary values
- [Percent Custom Field](/custom-fields/percent) - For percentage values
- [Formula Custom Field](/custom-fields/formula) - For calculated values
- [Automations API](/api/automations/index) - Create number-based automations
- [Filtering API](/api/filtering) - Advanced numeric filtering techniques