---
title: Country Custom Field
description: Create country selection fields with ISO country code validation
category: Custom Fields
---

Country custom fields allow you to store and manage country information for records. The field supports both country names and ISO Alpha-2 country codes, automatically validating and converting between formats.

## Basic Example

Create a simple country field:

```graphql
mutation CreateCountryField {
  createCustomField(input: {
    name: "Country of Origin"
    type: COUNTRY
    projectId: "proj_123"
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a country field with description:

```graphql
mutation CreateDetailedCountryField {
  createCustomField(input: {
    name: "Customer Location"
    type: COUNTRY
    projectId: "proj_123"
    description: "Primary country where the customer is located"
    isActive: true
  }) {
    id
    name
    type
    description
    isActive
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the country field |
| `type` | CustomFieldType! | ✅ Yes | Must be `COUNTRY` |
| `description` | String | No | Help text shown to users |

**Note**: The `projectId` is not passed in the input but is determined by the GraphQL context (typically from request headers or authentication).

## Setting Country Values

Country fields store data in two separate database fields that can be used independently or together:
- **`countryCodes`**: Stores validated ISO Alpha-2 country codes as an array
- **`text`**: Stores display text or country names as a string

### Understanding the Parameters

The `setTodoCustomField` mutation accepts two optional parameters for country fields:

| Parameter | Type | Required | Description | What it does |
|-----------|------|----------|-------------|--------------|
| `todoId` | String! | ✅ Yes | ID of the record to update | - |
| `customFieldId` | String! | ✅ Yes | ID of the country custom field | - |
| `countryCodes` | [String!] | No | Array of ISO Alpha-2 country codes | Stored in the `countryCodes` field |
| `text` | String | No | Display text or country names | Stored in the `text` field |

**Important**: Both `countryCodes` and `text` are optional and stored independently. You can use either one or both.

### Option 1: Using Only Country Codes

Store validated ISO codes without display text:

```graphql
mutation SetCountryByCode {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    countryCodes: ["US"]
  })
}
```

Result: `countryCodes` = `["US"]`, `text` = `null`

### Option 2: Using Only Text

Store display text without validated codes:

```graphql
mutation SetCountryByText {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    text: "United States"
  })
}
```

Result: `countryCodes` = `null`, `text` = `"United States"`

**Note**: When using only `text`, no validation or conversion to country codes occurs.

### Option 3: Using Both (Recommended)

Store both validated codes and display text:

```graphql
mutation SetCountryComplete {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    countryCodes: ["US"]
    text: "United States"
  })
}
```

Result: `countryCodes` = `["US"]`, `text` = `"United States"`

### Multiple Countries

Store multiple countries using arrays:

```graphql
mutation SetMultipleCountries {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    countryCodes: ["US", "CA", "MX"]
    text: "North American Markets"  # Can be any descriptive text
  })
}
```

## Creating Records with Country Values

When creating records, the `createTodo` mutation **automatically validates and converts** country values:

```graphql
mutation CreateRecordWithCountry {
  createTodo(input: {
    title: "International Client"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "country_field_id"
      value: "France"  # Can use country name or code
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
      text
      countryCodes
    }
  }
}
```

### Accepted Input Formats

| Input Type | Example | Result |
|------------|---------|---------|
| Country Name | `"United States"` | Stored as `US` |
| ISO Alpha-2 Code | `"GB"` | Stored as `GB` |
| Multiple (comma-separated) | `"US, CA"` | Stored as array `["US", "CA"]` |
| Mixed format | `"United States, CA"` | Converted to `["US", "CA"]` |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `text` | String | Display text (country names) |
| `countryCodes` | [String!] | Array of ISO Alpha-2 country codes |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## Country Standards

Blue uses the **ISO 3166-1 Alpha-2** standard for country codes:

- Two-letter country codes (e.g., US, GB, FR, DE)
- Validated using the `i18n-iso-countries` library
- Supports all officially recognized countries

### Example Country Codes

| Country | ISO Code |
|---------|----------|
| United States | `US` |
| United Kingdom | `GB` |
| Canada | `CA` |
| Germany | `DE` |
| France | `FR` |
| Japan | `JP` |
| Australia | `AU` |
| Brazil | `BR` |

For the complete official list of ISO 3166-1 alpha-2 country codes, visit the [ISO Online Browsing Platform](https://www.iso.org/obp/ui/#search/code/).

## Validation

The system validates country inputs:

1. **Valid ISO Code**: Accepts any valid ISO Alpha-2 code
2. **Country Name**: Automatically converts recognized country names to codes
3. **Invalid Input**: Throws `CustomFieldValueParseError` for unrecognized values

### Error Example

```json
{
  "errors": [{
    "message": "Invalid country value: 'Atlantis'",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

## Integration Features

### Lookup Fields
Country fields can be referenced by LOOKUP custom fields, allowing you to pull country data from related records.

### Automations
Use country values in automation conditions:
- Filter actions by specific countries
- Send notifications based on country
- Route tasks based on geographic regions

### Forms
Country fields in forms automatically validate user input and convert country names to codes.

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create country field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update country field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set country value | Standard record edit permissions |
| View country value | Standard record view permissions |

## Error Responses

### Invalid Country Value
```json
{
  "errors": [{
    "message": "Invalid country value provided",
    "extensions": {
      "code": "CUSTOM_FIELD_VALUE_PARSE_ERROR"
    }
  }]
}
```

### Field Type Mismatch
```json
{
  "errors": [{
    "message": "Field type mismatch: expected COUNTRY",
    "extensions": {
      "code": "INVALID_FIELD_TYPE"
    }
  }]
}
```

## Best Practices

### Input Handling
- Accept both country names and codes for user convenience
- Always store as ISO codes for consistency
- Display full country names in UI for clarity

### Data Quality
- Validate country inputs at entry point
- Use consistent formats across your system
- Consider regional groupings for reporting

### Multiple Countries
- Use array support for records that span multiple countries
- Separate multiple values with commas in text input
- Store all codes for complete data

## Common Use Cases

1. **Customer Management**
   - Customer headquarters location
   - Shipping destinations
   - Tax jurisdictions

2. **Project Tracking**
   - Project location
   - Team member locations
   - Market targets

3. **Compliance & Legal**
   - Regulatory jurisdictions
   - Data residency requirements
   - Export controls

4. **Sales & Marketing**
   - Territory assignments
   - Market segmentation
   - Campaign targeting

## Limitations

- Only supports ISO 3166-1 Alpha-2 codes (2-letter codes)
- No built-in support for country subdivisions (states/provinces)
- No automatic country flag icons (text-based only)
- Cannot validate historical country codes
- No built-in region or continent grouping

## Related Resources

- [Custom Fields Overview](/custom-fields/list-custom-fields) - General custom field concepts
- [Lookup Fields](/api/custom-fields/lookup) - Reference country data from other records
- [Forms API](/api/forms) - Include country fields in custom forms