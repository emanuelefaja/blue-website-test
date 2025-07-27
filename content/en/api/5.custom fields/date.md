---
title: Date Custom Field
description: Create date fields to track single dates or date ranges with timezone support
category: Custom Fields
---

Date custom fields allow you to store single dates or date ranges for records. They support timezone handling, intelligent formatting, and can be used to track deadlines, event dates, or any time-based information.

## Basic Example

Create a simple date field:

```graphql
mutation CreateDateField {
  createCustomField(input: {
    name: "Deadline"
    type: DATE
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a due date field with description:

```graphql
mutation CreateDueDateField {
  createCustomField(input: {
    name: "Contract Expiration"
    type: DATE
    isDueDate: true
    description: "When the contract expires and needs renewal"
  }) {
    id
    name
    type
    isDueDate
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the date field |
| `type` | CustomFieldType! | ✅ Yes | Must be `DATE` |
| `isDueDate` | Boolean | No | Whether this field represents a due date |
| `description` | String | No | Help text shown to users |

**Note**: Custom fields are automatically associated with the project based on the user's current project context. No `projectId` parameter is required.

## Setting Date Values

Date fields can store either a single date or a date range:

### Single Date

```graphql
mutation SetSingleDate {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    startDate: "2025-01-15T10:00:00Z"
    endDate: "2025-01-15T10:00:00Z"
    timezone: "America/New_York"
  }) {
    id
    customField {
      value  # Contains { startDate, endDate, timezone }
    }
  }
}
```

### Date Range

```graphql
mutation SetDateRange {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    startDate: "2025-01-01T09:00:00Z"
    endDate: "2025-01-31T17:00:00Z"
    timezone: "Europe/London"
  }) {
    id
    customField {
      value  # Contains { startDate, endDate, timezone }
    }
  }
}
```

### All-Day Event

```graphql
mutation SetAllDayEvent {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    startDate: "2025-01-15T00:00:00Z"
    endDate: "2025-01-15T23:59:59Z"
    timezone: "Asia/Tokyo"
  }) {
    id
    customField {
      value  # Contains { startDate, endDate, timezone }
    }
  }
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the date custom field |
| `startDate` | DateTime | No | Start date/time in ISO 8601 format |
| `endDate` | DateTime | No | End date/time in ISO 8601 format |
| `timezone` | String | No | Timezone identifier (e.g., "America/New_York") |

**Note**: If only `startDate` is provided, `endDate` automatically defaults to the same value.

## Date Formats

### ISO 8601 Format
All dates must be provided in ISO 8601 format:
- `2025-01-15T14:30:00Z` - UTC time
- `2025-01-15T14:30:00+05:00` - With timezone offset
- `2025-01-15T14:30:00.123Z` - With milliseconds

### Timezone Identifiers
Use standard timezone identifiers:
- `America/New_York`
- `Europe/London`
- `Asia/Tokyo`
- `Australia/Sydney`

If no timezone is provided, the system defaults to the user's detected timezone.

## Creating Records with Date Values

When creating a new record with date values:

```graphql
mutation CreateRecordWithDate {
  createTodo(input: {
    title: "Project Milestone"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "date_field_id"
      value: "2025-02-15"  # Simple date format
    }]
  }) {
    id
    title
    customFields {
      id
      customField {
        name
        type
        value  # Date values are accessed here
      }
    }
  }
}
```

### Supported Input Formats

When creating records, dates can be provided in various formats:

| Format | Example | Result |
|--------|---------|---------|
| ISO Date | `"2025-01-15"` | Single date (start and end same) |
| ISO DateTime | `"2025-01-15T10:00:00Z"` | Single date/time |
| Date Range | `"2025-01-01,2025-01-31"` | Start and end dates |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | ID! | Unique identifier for the field value |
| `uid` | String! | Unique identifier string |
| `customField` | CustomField! | The custom field definition (contains the date values) |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

**Important**: Date values (`startDate`, `endDate`, `timezone`) are accessed through the `customField.value` field, not directly on TodoCustomField.

### Value Object Structure

Date values are returned through the `customField.value` field as a JSON object:

```json
{
  "customField": {
    "value": {
      "startDate": "2025-01-15T10:00:00.000Z",
      "endDate": "2025-01-15T17:00:00.000Z",
      "timezone": "America/New_York"
    }
  }
}
```

**Note**: The `value` field is on the `CustomField` type, not on `TodoCustomField`.

## Querying Date Values

When querying records with date custom fields, access the date values through the `customField.value` field:

```graphql
query GetRecordWithDateField {
  todo(id: "todo_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
        value  # For DATE type, contains { startDate, endDate, timezone }
      }
    }
  }
}
```

The response will include the date values in the `value` field:

```json
{
  "data": {
    "todo": {
      "customFields": [{
        "customField": {
          "name": "Deadline",
          "type": "DATE",
          "value": {
            "startDate": "2025-01-15T10:00:00.000Z",
            "endDate": "2025-01-15T10:00:00.000Z",
            "timezone": "America/New_York"
          }
        }
      }]
    }
  }
}
```

## Date Display Intelligence

The system automatically formats dates based on the range:

| Scenario | Display Format |
|----------|----------------|
| Single date | `Jan 15, 2025` |
| All-day event | `Jan 15, 2025` (no time shown) |
| Same day with times | `Jan 15, 2025 10:00 AM - 5:00 PM` |
| Multi-day range | `Jan 1 → Jan 31, 2025` |

**All-day detection**: Events from 00:00 to 23:59 are automatically detected as all-day events.

## Timezone Handling

### Storage
- All dates are stored in UTC in the database
- Timezone information is preserved separately
- Conversion happens on display

### Best Practices
- Always provide timezone for accuracy
- Use consistent timezones within a project
- Consider user locations for global teams

### Common Timezones

| Region | Timezone ID | UTC Offset |
|--------|-------------|------------|
| US Eastern | `America/New_York` | UTC-5/-4 |
| US Pacific | `America/Los_Angeles` | UTC-8/-7 |
| UK | `Europe/London` | UTC+0/+1 |
| EU Central | `Europe/Berlin` | UTC+1/+2 |
| Japan | `Asia/Tokyo` | UTC+9 |
| Australia Eastern | `Australia/Sydney` | UTC+10/+11 |

## Filtering and Querying

Date fields support complex filtering:

```graphql
query FilterByDateRange {
  todos(filter: {
    customFields: [{
      customFieldId: "date_field_id"
      dateRange: {
        startDate: "2025-01-01T00:00:00Z"
        endDate: "2025-12-31T23:59:59Z"
      }
      operator: EQ  # Returns todos whose dates overlap with this range
    }]
  }) {
    id
    title
  }
}
```

### Checking for Empty Date Fields

```graphql
query FilterEmptyDates {
  todos(filter: {
    customFields: [{
      customFieldId: "date_field_id"
      values: null
      operator: IS  # Returns todos with no date set
    }]
  }) {
    id
    title
  }
}
```

### Supported Operators

| Operator | Usage | Description |
|----------|-------|-------------|
| `EQ` | With dateRange | Date overlaps with specified range (any intersection) |
| `NE` | With dateRange | Date does not overlap with range |
| `IS` | With `values: null` | Date field is empty (startDate or endDate is null) |
| `NOT` | With `values: null` | Date field has a value (both dates are not null) |

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create date field | `OWNER` or `ADMIN` role at company or project level |
| Update date field | `OWNER` or `ADMIN` role at company or project level |
| Set date value | Standard record edit permissions |
| View date value | Standard record view permissions |

## Error Responses

### Invalid Date Format
```json
{
  "errors": [{
    "message": "Invalid date format. Use ISO 8601 format",
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


## Limitations

- No recurring date support (use automations for recurring events)
- Cannot set time without date
- No built-in working days calculation
- Date ranges don't validate end > start automatically
- Maximum precision is to the second (no millisecond storage)

## Related Resources

- [Custom Fields Overview](/api/custom-fields/list-custom-fields) - General custom field concepts
- [Automations API](/api/automations/index) - Create date-based automations