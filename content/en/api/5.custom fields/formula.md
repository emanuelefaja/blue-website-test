---
title: Formula Custom Field
description: Create calculated fields that automatically compute values based on other data
category: Custom Fields
---

Formula custom fields automatically calculate values based on other custom fields, todo properties, and data from your project. They update automatically when source data changes and support various display formats including numbers, currency, and percentages.

## Basic Example

Create a simple formula field:

```graphql
mutation CreateFormulaField {
  createCustomField(input: {
    name: "Budget Total"
    type: FORMULA
    projectId: "proj_123"
    formula: {
      logic: {
        text: "SUM(Budget)"
        html: "<span>SUM(Budget)</span>"
      }
      display: {
        type: NUMBER
        precision: 2
      }
    }
  }) {
    id
    name
    type
    formula
  }
}
```

## Advanced Example

Create a currency formula with complex calculations:

```graphql
mutation CreateCurrencyFormula {
  createCustomField(input: {
    name: "Profit Margin"
    type: FORMULA
    projectId: "proj_123"
    formula: {
      logic: {
        text: "SUM(Revenue) - SUM(Costs)"
        html: "<span>SUM(Revenue) - SUM(Costs)</span>"
      }
      display: {
        type: CURRENCY
        currency: {
          code: "USD"
          name: "US Dollar"
        }
        precision: 2
      }
    }
    description: "Automatically calculates profit by subtracting costs from revenue"
  }) {
    id
    name
    type
    formula
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the formula field |
| `type` | CustomFieldType! | ✅ Yes | Must be `FORMULA` |
| `formula` | JSON! | ✅ Yes | Formula definition (see below) |
| `description` | String | No | Help text shown to users |

### Formula Structure

```json
{
  "logic": {
    "text": "Plain text formula",
    "html": "HTML formatted formula"
  },
  "display": {
    "type": "NUMBER|CURRENCY|PERCENTAGE",
    "currency": {
      "code": "USD",
      "name": "US Dollar"
    },
    "precision": 2,
    "function": "SUM|AVERAGE|COUNT|MAX|MIN"
  }
}
```

## Supported Functions

### Aggregation Functions

| Function | Description | Example |
|----------|-------------|---------|
| `SUM` | Sum of all values | `SUM(Budget)` |
| `AVERAGE` | Average of numeric values | `AVERAGE(Score)` |
| `AVERAGEA` | Average including text/empty values | `AVERAGEA(Rating)` |
| `COUNT` | Count of non-empty values | `COUNT(Tasks)` |
| `COUNTA` | Count of all values | `COUNTA(Items)` |
| `MAX` | Maximum value | `MAX(Priority)` |
| `MIN` | Minimum value | `MIN(StartDate)` |

### Mathematical Operations

- Addition: `SUM(Budget) + SUM(Bonus)`
- Subtraction: `SUM(Revenue) - SUM(Costs)`
- Multiplication: `COUNT(Hours) * 25`
- Division: `SUM(Total) / COUNT(Items)`
- Parentheses: `(SUM(A) + SUM(B)) / 2`

## Display Types

### Number Display

```json
{
  "display": {
    "type": "NUMBER",
    "precision": 2
  }
}
```

Result: `1250.75`

### Currency Display

```json
{
  "display": {
    "type": "CURRENCY",
    "currency": {
      "code": "USD",
      "name": "US Dollar"
    },
    "precision": 2
  }
}
```

Result: `$1,250.75`

### Percentage Display

```json
{
  "display": {
    "type": "PERCENTAGE",
    "precision": 1
  }
}
```

Result: `87.5%`

## Editing Formula Fields

Update existing formula fields:

```graphql
mutation EditFormulaField {
  editCustomField(input: {
    customFieldId: "field_456"
    formula: {
      logic: {
        text: "AVERAGE(Score)"
        html: "<span>AVERAGE(Score)</span>"
      }
      display: {
        type: PERCENTAGE
        precision: 1
      }
    }
  }) {
    id
    formula
  }
}
```

## Formula Calculation

### Automatic Updates

Formula fields automatically recalculate when:
- Source custom field values change
- Referenced todo properties are updated
- Dependencies in other formulas change
- Records are added or removed from the project

### Asynchronous Processing

- Calculations are processed in the background
- Results are cached for fast retrieval
- Updates are published via GraphQL subscriptions
- No blocking of user interface during calculations

### Dependency Management

- Formulas can reference other formulas
- Circular dependency detection prevents infinite loops
- Nested dependencies are properly resolved
- Source data changes trigger cascading updates

## Formula Result Storage

Results are stored in the `formulaResult` field:

```json
{
  "number": 1250.75,
  "formulaResult": {
    "number": 1250.75,
    "display": {
      "type": "CURRENCY",
      "currency": {
        "code": "USD",
        "name": "US Dollar"
      },
      "precision": 2
    }
  }
}
```

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The formula field definition |
| `number` | Float | Calculated numeric result |
| `formulaResult` | JSON | Full result with display formatting |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last calculated |

## Data Sources

### Custom Fields

Reference other custom fields in the same project:

```javascript
SUM(Budget)           // Sum budget custom field
AVERAGE(Rating)       // Average rating scores
COUNT(Status)         // Count non-empty status values
MAX(Priority)         // Maximum priority value
```

### Todo Properties

Reference built-in todo properties:

```javascript
COUNT(Assignees)      // Number of assignees
COUNT(Tags)           // Number of tags
COUNT(Comments)       // Number of comments
SUM(TimeSpent)        // Total time tracked
```

### Cross-Project References

Use REFERENCE custom fields to access data from other projects:

```javascript
SUM(Project.Budget)   // Sum budget from referenced project
COUNT(Tasks.Done)     // Count completed tasks from referenced project
```

## Common Formula Examples

### Budget Calculations

```json
{
  "logic": {
    "text": "SUM(Approved Budget) - SUM(Spent)",
    "html": "<span>SUM(Approved Budget) - SUM(Spent)</span>"
  },
  "display": {
    "type": "CURRENCY",
    "currency": { "code": "USD", "name": "US Dollar" },
    "precision": 2
  }
}
```

### Completion Rate

```json
{
  "logic": {
    "text": "(COUNT(Completed) / COUNT(Total)) * 100",
    "html": "<span>(COUNT(Completed) / COUNT(Total)) * 100</span>"
  },
  "display": {
    "type": "PERCENTAGE",
    "precision": 1
  }
}
```

### Average Score

```json
{
  "logic": {
    "text": "AVERAGE(Quality Score)",
    "html": "<span>AVERAGE(Quality Score)</span>"
  },
  "display": {
    "type": "NUMBER",
    "precision": 1
  }
}
```

### Resource Utilization

```json
{
  "logic": {
    "text": "SUM(Hours Used) / SUM(Hours Available)",
    "html": "<span>SUM(Hours Used) / SUM(Hours Available)</span>"
  },
  "display": {
    "type": "PERCENTAGE",
    "precision": 0
  }
}
```

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create formula field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update formula field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| View formula results | Standard record view permissions |
| Edit formula logic | Same as update permissions |

## Error Handling

### Invalid Formula Syntax
```json
{
  "errors": [{
    "message": "Invalid formula syntax",
    "extensions": {
      "code": "INVALID_FORMULA"
    }
  }]
}
```

### Circular Dependency
```json
{
  "errors": [{
    "message": "Circular dependency detected in formula",
    "extensions": {
      "code": "CIRCULAR_DEPENDENCY"
    }
  }]
}
```

### Missing Reference
```json
{
  "errors": [{
    "message": "Referenced field not found",
    "extensions": {
      "code": "FIELD_NOT_FOUND"
    }
  }]
}
```

## Best Practices

### Formula Design
- Use clear, descriptive names for formula fields
- Add descriptions explaining the calculation logic
- Test formulas with sample data before deployment
- Keep formulas simple and readable

### Performance Optimization
- Avoid deeply nested formula dependencies
- Use specific field references rather than wildcards
- Consider caching strategies for complex calculations
- Monitor formula performance in large projects

### Data Quality
- Validate source data before using in formulas
- Handle empty or null values appropriately
- Use appropriate precision for display types
- Consider edge cases in calculations

## Common Use Cases

1. **Financial Tracking**
   - Budget calculations
   - Profit/loss statements
   - Cost analysis
   - Revenue projections

2. **Project Management**
   - Completion percentages
   - Resource utilization
   - Timeline calculations
   - Performance metrics

3. **Quality Control**
   - Average scores
   - Pass/fail rates
   - Quality metrics
   - Compliance tracking

4. **Business Intelligence**
   - KPI calculations
   - Trend analysis
   - Comparative metrics
   - Dashboard values

## Limitations

- Cannot execute arbitrary code (security restriction)
- Limited to supported functions and operations
- No direct database queries
- Cannot reference external APIs or services
- Circular dependencies are blocked
- Complex calculations may have performance impact

## Related Resources

- [Number Fields](/api/custom-fields/number) - For static numeric values
- [Currency Fields](/api/custom-fields/currency) - For monetary values
- [Reference Fields](/api/custom-fields/reference) - For cross-project data
- [Lookup Fields](/api/custom-fields/lookup) - For aggregated data
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts