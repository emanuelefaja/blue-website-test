---
title: Currency Custom Field
description: Create currency fields to track monetary values with proper formatting and validation
category: Custom Fields
---

Currency custom fields allow you to store and manage monetary values with associated currency codes. The field supports 72 different currencies including major fiat currencies and cryptocurrencies, with automatic formatting and optional min/max constraints.

## Basic Example

Create a simple currency field:

```graphql
mutation CreateCurrencyField {
  createCustomField(input: {
    name: "Budget"
    type: CURRENCY
    projectId: "proj_123"
    currency: "USD"
  }) {
    id
    name
    type
    currency
  }
}
```

## Advanced Example

Create a currency field with validation constraints:

```graphql
mutation CreateConstrainedCurrencyField {
  createCustomField(input: {
    name: "Deal Value"
    type: CURRENCY
    projectId: "proj_123"
    currency: "EUR"
    min: 0
    max: 1000000
    description: "Estimated deal value in euros"
    isActive: true
  }) {
    id
    name
    type
    currency
    min
    max
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the currency field |
| `type` | CustomFieldType! | ✅ Yes | Must be `CURRENCY` |
| `projectId` | String! | ✅ Yes | Project ID where the field will be created |
| `currency` | String | No | Default currency code (3-letter ISO code) |
| `min` | Float | No | Minimum allowed value |
| `max` | Float | No | Maximum allowed value |
| `description` | String | No | Help text shown to users |
| `isActive` | Boolean | No | Whether the field is active (defaults to true) |

## Setting Currency Values

To set or update a currency value on a record:

```graphql
mutation SetCurrencyValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    number: 1500.50
    currency: "USD"
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the currency custom field |
| `number` | Float! | ✅ Yes | The monetary amount |
| `currency` | String! | ✅ Yes | 3-letter currency code |

## Creating Records with Currency Values

When creating a new record with currency values:

```graphql
mutation CreateRecordWithCurrency {
  createTodo(input: {
    title: "Q4 Marketing Campaign"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "currency_field_id"
      value: "25000.00"
      currency: "GBP"
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
      number
      currency
    }
  }
}
```

### Input Format for Create

When creating records, currency values are passed differently:

| Parameter | Type | Description |
|-----------|------|-------------|
| `customFieldId` | String! | ID of the currency field |
| `value` | String! | Amount as a string (e.g., "1500.50") |
| `currency` | String! | 3-letter currency code |

## Supported Currencies

Blue supports 72 currencies including 70 fiat currencies and 2 cryptocurrencies:

### Fiat Currencies

#### Americas
| Currency | Code | Name |
|----------|------|------|
| US Dollar | `USD` | US Dollar |
| Canadian Dollar | `CAD` | Canadian Dollar |
| Mexican Peso | `MXN` | Mexican Peso |
| Brazilian Real | `BRL` | Brazilian Real |
| Argentine Peso | `ARS` | Argentine Peso |
| Chilean Peso | `CLP` | Chilean Peso |
| Colombian Peso | `COP` | Colombian Peso |
| Peruvian Sol | `PEN` | Peruvian Sol |
| Uruguayan Peso | `UYU` | Uruguayan Peso |
| Venezuelan Bolívar | `VES` | Venezuelan Bolívar Soberano |
| Bolivian Boliviano | `BOB` | Bolivian Boliviano |
| Costa Rican Colón | `CRC` | Costa Rican Colón |
| Dominican Peso | `DOP` | Dominican Peso |
| Guatemalan Quetzal | `GTQ` | Guatemalan Quetzal |
| Jamaican Dollar | `JMD` | Jamaican Dollar |

#### Europe
| Currency | Code | Name |
|----------|------|------|
| Euro | `EUR` | Euro |
| British Pound | `GBP` | Pound Sterling |
| Swiss Franc | `CHF` | Swiss Franc |
| Swedish Krona | `SEK` | Swedish Krona |
| Norwegian Krone | `NOK` | Norwegian Krone |
| Danish Krone | `DKK` | Danish Krone |
| Polish Złoty | `PLN` | Polish Złoty |
| Czech Koruna | `CZK` | Czech Koruna |
| Hungarian Forint | `HUF` | Hungarian Forint |
| Romanian Leu | `RON` | Romanian Leu |
| Bulgarian Lev | `BGN` | Bulgarian Lev |
| Turkish Lira | `TRY` | Turkish Lira |
| Ukrainian Hryvnia | `UAH` | Ukrainian Hryvnia |
| Russian Ruble | `RUB` | Russian Ruble |
| Georgian Lari | `GEL` | Georgian Lari |
| Icelandic króna | `ISK` | Icelandic króna |
| Bosnia-Herzegovina Mark | `BAM` | Bosnia-Herzegovina Convertible Mark |

#### Asia-Pacific
| Currency | Code | Name |
|----------|------|------|
| Japanese Yen | `JPY` | Yen |
| Chinese Yuan | `CNY` | Yuan |
| Hong Kong Dollar | `HKD` | Hong Kong Dollar |
| Singapore Dollar | `SGD` | Singapore Dollar |
| Australian Dollar | `AUD` | Australian Dollar |
| New Zealand Dollar | `NZD` | New Zealand Dollar |
| South Korean Won | `KRW` | South Korean Won |
| Indian Rupee | `INR` | Indian Rupee |
| Indonesian Rupiah | `IDR` | Indonesian Rupiah |
| Thai Baht | `THB` | Thai Baht |
| Malaysian Ringgit | `MYR` | Malaysian Ringgit |
| Philippine Peso | `PHP` | Philippine Peso |
| Vietnamese Dong | `VND` | Vietnamese Dong |
| Taiwanese Dollar | `TWD` | New Taiwan Dollar |
| Pakistani Rupee | `PKR` | Pakistani Rupee |
| Sri Lankan Rupee | `LKR` | Sri Lankan Rupee |
| Cambodian Riel | `KHR` | Cambodian Riel |
| Kazakhstani Tenge | `KZT` | Kazakhstani Tenge |

#### Middle East & Africa
| Currency | Code | Name |
|----------|------|------|
| UAE Dirham | `AED` | UAE Dirham |
| Saudi Riyal | `SAR` | Saudi Riyal |
| Kuwaiti Dinar | `KWD` | Kuwaiti Dinar |
| Bahraini Dinar | `BHD` | Bahraini Dinar |
| Qatari Riyal | `QAR` | Qatari Riyal |
| Israeli Shekel | `ILS` | Israeli New Shekel |
| Egyptian Pound | `EGP` | Egyptian Pound |
| Moroccan Dirham | `MAD` | Moroccan Dirham |
| Tunisian Dinar | `TND` | Tunisian Dinar |
| South African Rand | `ZAR` | South African Rand |
| Kenyan Shilling | `KES` | Kenyan Shilling |
| Nigerian Naira | `NGN` | Nigerian Naira |
| Ghanaian Cedi | `GHS` | Ghanaian Cedi |
| Zambian Kwacha | `ZMW` | Zambian Kwacha |
| Malagasy Ariary | `MGA` | Malagasy Ariary |

### Cryptocurrencies
| Currency | Code |
|----------|------|
| Bitcoin | `BTC` |
| Ethereum | `ETH` |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `number` | Float | The monetary amount |
| `currency` | String | The 3-letter currency code |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

## Currency Formatting

The system automatically formats currency values based on locale:

- **Symbol placement**: Correctly positions currency symbols (before/after)
- **Decimal separators**: Uses locale-specific separators (. or ,)
- **Thousand separators**: Applies appropriate grouping
- **Decimal places**: Shows 0-2 decimal places based on the amount
- **Special handling**: USD/CAD show currency code prefix for clarity

### Formatting Examples

| Value | Currency | Display |
|-------|----------|---------|
| 1500.50 | USD | USD $1,500.50 |
| 1500.50 | EUR | €1.500,50 |
| 1500 | JPY | ¥1,500 |
| 1500.99 | GBP | £1,500.99 |

## Validation

### Amount Validation
- Must be a valid number
- Respects `min` and `max` constraints if set
- Supports up to 2 decimal places

### Currency Code Validation
- Must be one of the 72 supported currency codes
- Case-sensitive (use uppercase)
- Invalid codes return an error

## Integration Features

### Formulas
Currency fields can be used in FORMULA custom fields for calculations:
- Sum multiple currency fields
- Calculate percentages
- Perform arithmetic operations

### Currency Conversion
Use CURRENCY_CONVERSION fields to automatically convert between currencies (see [Currency Conversion Fields](/api/custom-fields/currency-conversion))

### Automations
Currency values can trigger automations based on:
- Amount thresholds
- Currency type
- Value changes

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create currency field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update currency field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set currency value | Standard record edit permissions |
| View currency value | Standard record view permissions |

## Error Responses

### Invalid Currency Code
```json
{
  "errors": [{
    "message": "Invalid currency code: 'XYZ'",
    "extensions": {
      "code": "INVALID_CURRENCY"
    }
  }]
}
```

### Value Out of Range
```json
{
  "errors": [{
    "message": "Value must be between 0 and 1000000",
    "extensions": {
      "code": "VALUE_OUT_OF_RANGE"
    }
  }]
}
```

### Invalid Number Format
```json
{
  "errors": [{
    "message": "Invalid number format",
    "extensions": {
      "code": "INVALID_NUMBER"
    }
  }]
}
```

## Best Practices

### Currency Selection
- Set a default currency that matches your primary market
- Use ISO 4217 currency codes consistently
- Consider user location when choosing defaults

### Value Constraints
- Set reasonable min/max values to prevent data entry errors
- Use 0 as minimum for fields that shouldn't be negative
- Consider your use case when setting maximums

### Multi-Currency Projects
- Use consistent base currency for reporting
- Implement CURRENCY_CONVERSION fields for automatic conversion
- Document which currency should be used for each field

## Common Use Cases

1. **Project Budgeting**
   - Project budget tracking
   - Cost estimates
   - Expense tracking

2. **Sales & Deals**
   - Deal values
   - Contract amounts
   - Revenue tracking

3. **Financial Planning**
   - Investment amounts
   - Funding rounds
   - Financial targets

4. **International Business**
   - Multi-currency pricing
   - Foreign exchange tracking
   - Cross-border transactions

## Limitations

- Maximum 2 decimal places for display (though more precision is stored)
- No built-in currency conversion in standard CURRENCY fields
- Cannot mix currencies in a single field value
- No automatic exchange rate updates (use CURRENCY_CONVERSION for this)
- Currency symbols are not customizable

## Related Resources

- [Currency Conversion Fields](/api/custom-fields/currency-conversion) - Automatic currency conversion
- [Number Fields](/api/custom-fields/number) - For non-monetary numeric values
- [Formula Fields](/api/custom-fields/formula) - Calculate with currency values
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General custom field concepts