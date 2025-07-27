---
title: Currency Conversion Custom Field
description: Create fields that automatically convert currency values using real-time exchange rates
category: Custom Fields
---

Currency Conversion custom fields automatically convert values from a source CURRENCY field to different target currencies using real-time exchange rates. These fields update automatically whenever the source currency value changes.

The conversion rates are provided by the [Frankfurter API](https://github.com/hakanensari/frankfurter), an open-source service that tracks reference exchange rates published by the [European Central Bank](https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html). This ensures accurate, reliable, and up-to-date currency conversions for your international business needs.

## Basic Example

Create a simple currency conversion field:

```graphql
mutation CreateCurrencyConversionField {
  createCustomField(input: {
    name: "Price in EUR"
    type: CURRENCY_CONVERSION
    currencyFieldId: "source_currency_field_id"
    conversionDateType: "currentDate"
  }) {
    id
    name
    type
    currencyFieldId
    conversionDateType
  }
}
```

## Advanced Example

Create a conversion field with a specific date for historical rates:

```graphql
mutation CreateHistoricalConversionField {
  createCustomField(input: {
    name: "Q1 Budget in Local Currency"
    type: CURRENCY_CONVERSION
    currencyFieldId: "budget_field_id"
    conversionDateType: "specificDate"
    conversionDate: "2024-01-01T00:00:00Z"
    description: "Budget converted at Q1 exchange rates"
  }) {
    id
    name
    type
    currencyFieldId
    conversionDateType
    conversionDate
  }
}
```

## Complete Setup Process

Setting up a currency conversion field requires three steps:

### Step 1: Create a Source CURRENCY Field

```graphql
mutation CreateSourceCurrencyField {
  createCustomField(input: {
    name: "Contract Value"
    type: CURRENCY
    currency: "USD"
  }) {
    id  # Save this ID for Step 2
    name
    type
  }
}
```

### Step 2: Create the CURRENCY_CONVERSION Field

```graphql
mutation CreateConversionField {
  createCustomField(input: {
    name: "Contract Value (Local Currency)"
    type: CURRENCY_CONVERSION
    currencyFieldId: "source_field_id_from_step_1"
    conversionDateType: "currentDate"
  }) {
    id  # Save this ID for Step 3
    name
    type
  }
}
```

### Step 3: Create Conversion Options

```graphql
mutation CreateConversionOptions {
  createCustomFieldOptions(input: {
    customFieldId: "conversion_field_id_from_step_2"
    customFieldOptions: [
      {
        title: "USD to EUR"
        currencyConversionFrom: "USD"
        currencyConversionTo: "EUR"
      },
      {
        title: "USD to GBP"
        currencyConversionFrom: "USD"
        currencyConversionTo: "GBP"
      },
      {
        title: "Any to JPY"
        currencyConversionFrom: "Any"
        currencyConversionTo: "JPY"
      }
    ]
  }) {
    id
    title
    currencyConversionFrom
    currencyConversionTo
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the conversion field |
| `type` | CustomFieldType! | ✅ Yes | Must be `CURRENCY_CONVERSION` |
| `currencyFieldId` | String | No | ID of the source CURRENCY field to convert from |
| `conversionDateType` | String | No | Date strategy for exchange rates (see below) |
| `conversionDate` | String | No | Date string for conversion (based on conversionDateType) |
| `description` | String | No | Help text shown to users |

**Note**: Custom fields are automatically associated with the project based on the user's current project context. No `projectId` parameter is required.

### Conversion Date Types

| Type | Description | conversionDate Parameter |
|------|-------------|-------------------------|
| `currentDate` | Uses real-time exchange rates | Not required |
| `specificDate` | Uses rates from a fixed date | ISO date string (e.g., "2024-01-01T00:00:00Z") |
| `fromDateField` | Uses date from another field | "todoDueDate" or DATE field ID |

## Creating Conversion Options

Conversion options define which currency pairs can be converted:

### CreateCustomFieldOptionInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `customFieldId` | String! | ✅ Yes | ID of the CURRENCY_CONVERSION field |
| `title` | String! | ✅ Yes | Display name for this conversion option |
| `currencyConversionFrom` | String! | ✅ Yes | Source currency code or "Any" |
| `currencyConversionTo` | String! | ✅ Yes | Target currency code |

### Using "Any" as Source

The special value "Any" as `currencyConversionFrom` creates a fallback option:

```graphql
mutation CreateUniversalConversion {
  createCustomFieldOption(input: {
    customFieldId: "conversion_field_id"
    title: "Any currency to EUR"
    currencyConversionFrom: "Any"
    currencyConversionTo: "EUR"
  }) {
    id
  }
}
```

This option will be used when no specific currency pair match is found.

## How Automatic Conversion Works

1. **Value Update**: When a value is set in the source CURRENCY field
2. **Option Matching**: System finds matching conversion option based on source currency
3. **Rate Fetching**: Retrieves exchange rate from Frankfurter API
4. **Calculation**: Multiplies source amount by exchange rate
5. **Storage**: Saves converted value with target currency code

### Example Flow

```graphql
# 1. Set value in source CURRENCY field
mutation SetSourceValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "source_currency_field_id"
    number: 1000
    currency: "USD"
  })
}

# 2. CURRENCY_CONVERSION fields automatically update
# If you have USD→EUR and USD→GBP options configured,
# both conversion fields will calculate and store their values
```

## Date-Based Conversions

### Using Current Date

```graphql
mutation CreateRealtimeConversion {
  createCustomField(input: {
    name: "Current EUR Value"
    type: CURRENCY_CONVERSION
    currencyFieldId: "source_field_id"
    conversionDateType: "currentDate"
  })
}
```

Conversions update with current exchange rates each time the source value changes.

### Using Specific Date

```graphql
mutation CreateFixedDateConversion {
  createCustomField(input: {
    name: "Year-End 2023 Value"
    type: CURRENCY_CONVERSION
    currencyFieldId: "source_field_id"
    conversionDateType: "specificDate"
    conversionDate: "2023-12-31T00:00:00Z"
  })
}
```

Always uses exchange rates from the specified date.

### Using Date from Field

```graphql
mutation CreateDateFieldConversion {
  createCustomField(input: {
    name: "Value at Contract Date"
    type: CURRENCY_CONVERSION
    currencyFieldId: "source_field_id"
    conversionDateType: "fromDateField"
    conversionDate: "contract_date_field_id"  # ID of a DATE custom field
  })
}
```

Uses the date from another field (either todo due date or a DATE custom field).

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The conversion field definition |
| `number` | Float | The converted amount |
| `currency` | String | The target currency code |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last updated |

## Exchange Rate Source

Blue uses the **Frankfurter API** for exchange rates:
- Open-source API hosted by the European Central Bank
- Updates daily with official exchange rates
- Supports historical rates back to 1999
- Free and reliable for business use

## Error Handling

### Conversion Failures

When conversion fails (API error, invalid currency, etc.):
- The converted value is set to `0`
- The target currency is still stored
- No error is thrown to the user

### Common Scenarios

| Scenario | Result |
|----------|---------|
| Same currency (USD→USD) | Value copied without API call |
| Invalid currency code | Conversion returns 0 |
| API unavailable | Conversion returns 0 |
| No matching option | Uses "Any" option if available |
| Missing source value | No conversion performed |

## Required Permissions

Custom field management requires project-level access:

| Role | Can Create/Update Fields |
|------|-------------------------|
| `OWNER` | ✅ Yes |
| `ADMIN` | ✅ Yes |
| `MEMBER` | ❌ No |
| `CLIENT` | ❌ No |

View permissions for converted values follow standard record access rules.

## Best Practices

### Option Configuration
- Create specific currency pairs for common conversions
- Add an "Any" fallback option for flexibility
- Use descriptive titles for options

### Date Strategy Selection
- Use `currentDate` for live financial tracking
- Use `specificDate` for historical reporting
- Use `fromDateField` for transaction-specific rates

### Performance Considerations
- Multiple conversion fields update in parallel
- API calls are made only when source value changes
- Same-currency conversions skip API calls

## Common Use Cases

1. **Multi-Currency Projects**
   - Track project costs in local currencies
   - Report total budget in company currency
   - Compare values across regions

2. **International Sales**
   - Convert deal values to reporting currency
   - Track revenue in multiple currencies
   - Historical conversion for closed deals

3. **Financial Reporting**
   - Period-end currency conversions
   - Consolidated financial statements
   - Budget vs. actual in local currency

4. **Contract Management**
   - Convert contract values at signing date
   - Track payment schedules in multiple currencies
   - Currency risk assessment

## Limitations

- No support for cryptocurrency conversions
- Cannot set converted values manually (always calculated)
- Fixed 2 decimal places precision for all converted amounts
- No support for custom exchange rates
- No caching of exchange rates (fresh API call for each conversion)
- Depends on Frankfurter API availability

## Related Resources

- [Currency Fields](/api/custom-fields/currency) - Source fields for conversions
- [Date Fields](/api/custom-fields/date) - For date-based conversions
- [Formula Fields](/api/custom-fields/formula) - Alternative calculations
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts