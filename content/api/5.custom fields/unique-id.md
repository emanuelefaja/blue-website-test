---
title: Unique ID Custom Field
description: Create auto-generated unique identifier fields with sequential numbering and custom formatting
category: Custom Fields
---

Unique ID custom fields automatically generate sequential, unique identifiers for your records. They're perfect for creating ticket numbers, order IDs, invoice numbers, or any sequential identifier system in your workflow.

## Basic Example

Create a simple unique ID field with auto-sequencing:

```graphql
mutation CreateUniqueIdField {
  createCustomField(input: {
    name: "Ticket Number"
    type: UNIQUE_ID
    projectId: "proj_123"
    useSequenceUniqueId: true
  }) {
    id
    name
    type
    useSequenceUniqueId
  }
}
```

## Advanced Example

Create a formatted unique ID field with prefix and zero-padding:

```graphql
mutation CreateFormattedUniqueIdField {
  createCustomField(input: {
    name: "Order ID"
    type: UNIQUE_ID
    projectId: "proj_123"
    description: "Auto-generated order identifier"
    useSequenceUniqueId: true
    prefix: "ORD-"
    sequenceDigits: 4
    sequenceStartingNumber: 1000
  }) {
    id
    name
    type
    description
    useSequenceUniqueId
    prefix
    sequenceDigits
    sequenceStartingNumber
  }
}
```

## Input Parameters

### CreateCustomFieldInput (UNIQUE_ID)

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the unique ID field |
| `type` | CustomFieldType! | ✅ Yes | Must be `UNIQUE_ID` |
| `description` | String | No | Help text shown to users |
| `useSequenceUniqueId` | Boolean | No | Enable auto-sequencing (default: false) |
| `prefix` | String | No | Text prefix for generated IDs (e.g., "TASK-") |
| `sequenceDigits` | Int | No | Number of digits for zero-padding |
| `sequenceStartingNumber` | Int | No | Starting number for the sequence |

## Configuration Options

### Auto-Sequencing (`useSequenceUniqueId`)
- **true**: Automatically generates sequential IDs when records are created
- **false**: Manual entry required (functions like a text field)

### Prefix (`prefix`)
- Optional text prefix added to all generated IDs
- Examples: "TASK-", "ORD-", "BUG-", "REQ-"
- No length limit, but keep reasonable for display

### Sequence Digits (`sequenceDigits`)
- Number of digits for zero-padding the sequence number
- Example: `sequenceDigits: 3` produces `001`, `002`, `003`
- If not specified, no padding is applied

### Starting Number (`sequenceStartingNumber`)
- The first number in the sequence
- Example: `sequenceStartingNumber: 1000` starts at 1000, 1001, 1002...
- If not specified, starts at 1

## Generated ID Format

The final ID format combines all configuration options:

```
{prefix}{paddedSequenceNumber}
```

### Format Examples

| Configuration | Generated IDs |
|---------------|---------------|
| No options | `1`, `2`, `3` |
| `prefix: "TASK-"` | `TASK-1`, `TASK-2`, `TASK-3` |
| `sequenceDigits: 3` | `001`, `002`, `003` |
| `prefix: "ORD-", sequenceDigits: 4` | `ORD-0001`, `ORD-0002`, `ORD-0003` |
| `prefix: "BUG-", sequenceStartingNumber: 500` | `BUG-500`, `BUG-501`, `BUG-502` |
| All options combined | `TASK-1001`, `TASK-1002`, `TASK-1003` |

## Reading Unique ID Values

### Query Records with Unique IDs
```graphql
query GetRecordsWithUniqueIds {
  todos(projectId: "proj_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
        prefix
        sequenceDigits
      }
      sequenceId    # The generated sequence number
      value        # Same as sequenceId for UNIQUE_ID fields
    }
  }
}
```

### Response Format
```json
{
  "data": {
    "todos": [
      {
        "id": "todo_123",
        "title": "Fix login issue",
        "customFields": [
          {
            "id": "field_value_456",
            "customField": {
              "name": "Ticket Number",
              "type": "UNIQUE_ID",
              "prefix": "TASK-",
              "sequenceDigits": 3
            },
            "sequenceId": 42,
            "value": 42
          }
        ]
      }
    ]
  }
}
```

## Automatic ID Generation

### When IDs Are Generated
- **Record Creation**: IDs are automatically assigned when new records are created
- **Field Addition**: When adding a UNIQUE_ID field to existing records, IDs are generated for all records
- **Background Processing**: ID generation happens asynchronously via job queues

### Generation Process
1. **Trigger**: New record is created or UNIQUE_ID field is added
2. **Sequence Lookup**: System finds the next available sequence number
3. **ID Assignment**: Sequence number is assigned to the record
4. **Counter Update**: Sequence counter is incremented for future records
5. **Formatting**: ID is formatted with prefix and padding when displayed

### Uniqueness Guarantees
- **Database Constraints**: Unique constraint on sequence IDs within each field
- **Atomic Operations**: Sequence generation uses database locks to prevent duplicates
- **Project Scoping**: Sequences are independent per project
- **Race Condition Protection**: Concurrent requests are handled safely

## Manual vs Automatic Mode

### Automatic Mode (`useSequenceUniqueId: true`)
- IDs are automatically generated
- Users cannot manually edit values
- Sequential numbering is guaranteed
- Background processing handles ID assignment

### Manual Mode (`useSequenceUniqueId: false`)
- Functions like a regular text field
- Users can input custom values
- No automatic generation
- No uniqueness enforcement beyond database constraints

## Setting Manual Values (Manual Mode Only)

When `useSequenceUniqueId` is false, you can set values manually:

```graphql
mutation SetUniqueIdValue {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    text: "CUSTOM-ID-001"
  })
}
```

## Response Fields

### TodoCustomField Response (UNIQUE_ID)

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `sequenceId` | Int | The generated sequence number (auto mode) |
| `text` | String | The stored text value (manual mode) |
| `value` | Any | Returns sequenceId for auto mode, text for manual mode |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last updated |

### CustomField Response (UNIQUE_ID)

| Field | Type | Description |
|-------|------|-------------|
| `useSequenceUniqueId` | Boolean | Whether auto-sequencing is enabled |
| `prefix` | String | Text prefix for generated IDs |
| `sequenceDigits` | Int | Number of digits for zero-padding |
| `sequenceStartingNumber` | Int | Starting number for the sequence |

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create unique ID field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update unique ID field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Set manual value | Standard record edit permissions |
| View unique ID value | Standard record view permissions |

## Error Responses

### Field Configuration Error
```json
{
  "errors": [{
    "message": "Invalid sequence configuration",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

### Manual Edit in Auto Mode
```json
{
  "errors": [{
    "message": "Cannot manually set value for auto-sequence field",
    "extensions": {
      "code": "FORBIDDEN"
    }
  }]
}
```

## Important Notes

### Auto-Generated IDs
- **Read-Only**: Auto-generated IDs cannot be manually edited
- **Permanent**: Once assigned, sequence IDs don't change
- **Chronological**: IDs reflect creation order
- **Scoped**: Sequences are independent per project

### Performance Considerations
- ID generation is asynchronous and may take a moment
- Large batches of records are processed in background jobs
- Sequence generation uses database locks (minimal performance impact)
- Consider sequence starting numbers for high-volume projects

### Migration and Updates
- Adding auto-sequencing to existing records triggers background processing
- Changing sequence settings affects only future records
- Existing IDs remain unchanged when configuration updates
- Sequence counters continue from current maximum

## Best Practices

### Configuration Design
- Choose descriptive prefixes that won't conflict with other systems
- Use appropriate digit padding for your expected volume
- Set reasonable starting numbers to avoid conflicts
- Test configuration with sample data before deployment

### Prefix Guidelines
- Keep prefixes short and memorable (2-5 characters)
- Use uppercase for consistency
- Include separators (hyphens, underscores) for readability
- Avoid special characters that might cause issues in URLs or systems

### Sequence Planning
- Estimate your record volume to choose appropriate digit padding
- Consider future growth when setting starting numbers
- Plan for different sequence ranges for different record types
- Document your ID schemes for team reference

## Common Use Cases

1. **Support Systems**
   - Ticket numbers: `TICK-001`, `TICK-002`
   - Case IDs: `CASE-2024-001`
   - Support requests: `SUP-001`

2. **Project Management**
   - Task IDs: `TASK-001`, `TASK-002`
   - Sprint items: `SPRINT-001`
   - Deliverable numbers: `DEL-001`

3. **Business Operations**
   - Order numbers: `ORD-2024-001`
   - Invoice IDs: `INV-001`
   - Purchase orders: `PO-001`

4. **Quality Management**
   - Bug reports: `BUG-001`
   - Test case IDs: `TEST-001`
   - Review numbers: `REV-001`

## Integration Features

### With Automations
- Trigger actions when unique IDs are assigned
- Use ID patterns in automation rules
- Reference IDs in email templates and notifications

### With Lookups
- Reference unique IDs from other records
- Find records by unique ID
- Display related record identifiers

### With Reporting
- Group and filter by ID patterns
- Track ID assignment trends
- Monitor sequence usage and gaps

## Limitations

- **Sequential Only**: IDs are assigned in chronological order
- **No Gaps**: Deleted records leave gaps in sequences
- **No Reuse**: Sequence numbers are never reused
- **Project Scoped**: Cannot share sequences across projects
- **Format Constraints**: Limited formatting options
- **No Bulk Updates**: Cannot bulk update existing sequence IDs
- **No Custom Logic**: Cannot implement custom ID generation rules

## Related Resources

- [Text Fields](/api/custom-fields/text-single) - For manual text identifiers
- [Number Fields](/api/custom-fields/number) - For numeric sequences
- [Custom Fields Overview](/api/custom-fields/list-custom-fields) - General concepts
- [Automations](/api/automations) - For ID-based automation rules