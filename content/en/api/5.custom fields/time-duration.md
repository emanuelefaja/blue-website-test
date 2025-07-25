---
title: Time Duration Custom Field
description: Create calculated time duration fields that track time between events in your workflow
category: Custom Fields
---

Time Duration custom fields automatically calculate and display the duration between two events in your workflow. They're ideal for tracking processing times, response times, cycle times, or any time-based metrics in your projects.

## Basic Example

Create a simple time duration field that tracks how long tasks take to complete:

```graphql
mutation CreateTimeDurationField {
  createCustomField(input: {
    name: "Processing Time"
    type: TIME_DURATION
    projectId: "proj_123"
    timeDurationDisplay: FULL_DATE_SUBSTRING
    timeDurationStartInput: {
      type: TODO_CREATED_AT
      condition: FIRST
    }
    timeDurationEndInput: {
      type: TODO_MARKED_AS_COMPLETE
      condition: FIRST
    }
  }) {
    id
    name
    type
    timeDurationDisplay
    timeDurationStart {
      type
      condition
    }
    timeDurationEnd {
      type
      condition
    }
  }
}
```

## Advanced Example

Create a complex time duration field that tracks time between custom field changes:

```graphql
mutation CreateAdvancedTimeDurationField {
  createCustomField(input: {
    name: "Review Cycle Time"
    type: TIME_DURATION
    projectId: "proj_123"
    description: "Time from review request to approval"
    timeDurationDisplay: FULL_DATE_STRING
    timeDurationStartInput: {
      type: TODO_CUSTOM_FIELD
      condition: FIRST
      customFieldId: "status_field_id"
      customFieldOptionIds: ["review_requested_option_id"]
    }
    timeDurationEndInput: {
      type: TODO_CUSTOM_FIELD
      condition: FIRST
      customFieldId: "status_field_id"
      customFieldOptionIds: ["approved_option_id"]
    }
  }) {
    id
    name
    type
    description
    timeDurationDisplay
    timeDurationStart {
      type
      condition
      customField {
        name
      }
    }
    timeDurationEnd {
      type
      condition
      customField {
        name
      }
    }
  }
}
```

## Input Parameters

### CreateCustomFieldInput (TIME_DURATION)

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the duration field |
| `type` | CustomFieldType! | ✅ Yes | Must be `TIME_DURATION` |
| `description` | String | No | Help text shown to users |
| `timeDurationDisplay` | CustomFieldTimeDurationDisplayType! | ✅ Yes | How to display the duration |
| `timeDurationStartInput` | CustomFieldTimeDurationInput! | ✅ Yes | Start event configuration |
| `timeDurationEndInput` | CustomFieldTimeDurationInput! | ✅ Yes | End event configuration |

### CustomFieldTimeDurationInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | CustomFieldTimeDurationType! | ✅ Yes | Type of event to track |
| `condition` | CustomFieldTimeDurationCondition! | ✅ Yes | `FIRST` or `LAST` occurrence |
| `customFieldId` | String | Conditional | Required for `TODO_CUSTOM_FIELD` type |
| `customFieldOptionIds` | [String!] | Conditional | Required for select field changes |
| `todoListId` | String | Conditional | Required for `TODO_MOVED` type |
| `tagId` | String | Conditional | Required for `TODO_TAG_ADDED` type |
| `assigneeId` | String | Conditional | Required for `TODO_ASSIGNEE_ADDED` type |

### CustomFieldTimeDurationType Values

| Value | Description |
|-------|-------------|
| `TODO_CREATED_AT` | When the record was created |
| `TODO_CUSTOM_FIELD` | When a custom field value changed |
| `TODO_DUE_DATE` | When the due date was set |
| `TODO_MARKED_AS_COMPLETE` | When the record was marked complete |
| `TODO_MOVED` | When the record was moved to a different list |
| `TODO_TAG_ADDED` | When a tag was added to the record |
| `TODO_ASSIGNEE_ADDED` | When an assignee was added to the record |

### CustomFieldTimeDurationCondition Values

| Value | Description |
|-------|-------------|
| `FIRST` | Use the first occurrence of the event |
| `LAST` | Use the last occurrence of the event |

### CustomFieldTimeDurationDisplayType Values

| Value | Description | Example |
|-------|-------------|---------|
| `FULL_DATE` | Days:Hours:Minutes:Seconds format | `"01:02:03:04"` |
| `FULL_DATE_STRING` | Written out in full words | `"Two hours, two minutes, three seconds"` |
| `FULL_DATE_SUBSTRING` | Numeric with units | `"1 hour, 2 minutes, 3 seconds"` |

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The custom field definition |
| `number` | Float | The duration in seconds |
| `value` | Float | Alias for number (duration in seconds) |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last updated |

### CustomField Response (TIME_DURATION)

| Field | Type | Description |
|-------|------|-------------|
| `timeDurationDisplay` | CustomFieldTimeDurationDisplayType | Display format for the duration |
| `timeDurationStart` | CustomFieldTimeDuration | Start event configuration |
| `timeDurationEnd` | CustomFieldTimeDuration | End event configuration |

## Duration Calculation

### How It Works
1. **Start Event**: System monitors for the specified start event
2. **End Event**: System monitors for the specified end event
3. **Calculation**: Duration = End Time - Start Time
4. **Storage**: Duration stored in seconds as a number
5. **Display**: Formatted according to `timeDurationDisplay` setting

### Update Triggers
Duration values are automatically recalculated when:
- Records are created or updated
- Custom field values change
- Tags are added or removed
- Assignees are added or removed
- Records are moved between lists
- Records are marked complete/incomplete

## Reading Duration Values

### Query Duration Fields
```graphql
query GetTaskWithDuration {
  todo(id: "todo_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
        timeDurationDisplay
      }
      number    # Duration in seconds
      value     # Same as number
    }
  }
}
```

### Formatted Display Values
Duration values are automatically formatted based on the `timeDurationDisplay` setting:

```javascript
// FULL_DATE format
93784 seconds → "01:02:03:04" (1 day, 2 hours, 3 minutes, 4 seconds)

// FULL_DATE_STRING format
7323 seconds → "Two hours, two minutes, three seconds"

// FULL_DATE_SUBSTRING format
3723 seconds → "1 hour, 2 minutes, 3 seconds"
```

## Common Configuration Examples

### Task Completion Time
```graphql
timeDurationStartInput: {
  type: TODO_CREATED_AT
  condition: FIRST
}
timeDurationEndInput: {
  type: TODO_MARKED_AS_COMPLETE
  condition: FIRST
}
```

### Status Change Duration
```graphql
timeDurationStartInput: {
  type: TODO_CUSTOM_FIELD
  condition: FIRST
  customFieldId: "status_field_id"
  customFieldOptionIds: ["in_progress_option_id"]
}
timeDurationEndInput: {
  type: TODO_CUSTOM_FIELD
  condition: FIRST
  customFieldId: "status_field_id"
  customFieldOptionIds: ["completed_option_id"]
}
```

### Time in Specific List
```graphql
timeDurationStartInput: {
  type: TODO_MOVED
  condition: FIRST
  todoListId: "review_list_id"
}
timeDurationEndInput: {
  type: TODO_MOVED
  condition: FIRST
  todoListId: "approved_list_id"
}
```

### Assignment Response Time
```graphql
timeDurationStartInput: {
  type: TODO_ASSIGNEE_ADDED
  condition: FIRST
  assigneeId: "user_123"
}
timeDurationEndInput: {
  type: TODO_CUSTOM_FIELD
  condition: FIRST
  customFieldId: "status_field_id"
  customFieldOptionIds: ["started_option_id"]
}
```

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create duration field | Project-level `OWNER` or `ADMIN` role |
| Update duration field | Project-level `OWNER` or `ADMIN` role |
| View duration value | Any project member role |

## Error Responses

### Invalid Configuration
```json
{
  "errors": [{
    "message": "Custom field is required for TODO_CUSTOM_FIELD type",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

### Referenced Field Not Found
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

### Missing Required Options
```json
{
  "errors": [{
    "message": "Custom field options are required for select field changes",
    "extensions": {
      "code": "VALIDATION_ERROR"
    }
  }]
}
```

## Important Notes

### Automatic Calculation
- Duration fields are **read-only** - values are automatically calculated
- You cannot manually set duration values via API
- Calculations happen asynchronously via background jobs
- Values update automatically when trigger events occur

### Performance Considerations
- Duration calculations are queued and processed asynchronously
- Large numbers of duration fields may impact performance
- Consider the frequency of trigger events when designing duration fields
- Use specific conditions to avoid unnecessary recalculations

### Null Values
Duration fields will show `null` when:
- Start event hasn't occurred yet
- End event hasn't occurred yet
- Configuration references non-existent entities
- Calculation encounters an error

## Best Practices

### Configuration Design
- Use specific event types rather than generic ones when possible
- Choose appropriate `FIRST` vs `LAST` conditions based on your workflow
- Test duration calculations with sample data before deployment
- Document your duration field logic for team members

### Display Formatting
- Use `FULL_DATE_SUBSTRING` for most readable format
- Use `FULL_DATE` for compact, consistent width display
- Use `FULL_DATE_STRING` for formal reports and documents
- Consider your UI space constraints when choosing format

### Workflow Integration
- Design duration fields to match your actual business processes
- Use duration data for process improvement and optimization
- Monitor duration trends to identify workflow bottlenecks
- Set up alerts for duration thresholds if needed

## Common Use Cases

1. **Process Performance**
   - Task completion times
   - Review cycle times
   - Approval processing times
   - Response times

2. **SLA Monitoring**
   - Time to first response
   - Resolution times
   - Escalation timeframes
   - Service level compliance

3. **Workflow Analytics**
   - Bottleneck identification
   - Process optimization
   - Team performance metrics
   - Quality assurance timing

4. **Project Management**
   - Phase durations
   - Milestone timing
   - Resource allocation time
   - Delivery timeframes

## Limitations

- Duration fields are **read-only** and cannot be manually set
- Values are calculated asynchronously and may not be immediately available
- Requires proper event triggers to be set up in your workflow
- Cannot calculate durations for events that haven't occurred
- Limited to tracking time between discrete events (not continuous time tracking)
- No built-in SLA alerts or notifications
- Cannot aggregate multiple duration calculations into a single field

## Related Resources

- [Number Fields](/api/custom-fields/number) - For manual numeric values
- [Date Fields](/api/custom-fields/date) - For specific date tracking
- [Custom Fields Overview](/api/custom-fields/list-custom-fields) - General concepts
- [Automations](/api/automations) - For triggering actions based on duration thresholds