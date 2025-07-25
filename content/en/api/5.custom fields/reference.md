---
title: Reference Custom Field
description: Create reference fields that link to records in other projects for cross-project relationships
category: Custom Fields
---

Reference custom fields allow you to create links between records in different projects, enabling cross-project relationships and data sharing. They provide a powerful way to connect related work across your organization's project structure.

## Basic Example

Create a simple reference field:

```graphql
mutation CreateReferenceField {
  createCustomField(input: {
    name: "Related Project"
    type: REFERENCE
    referenceProjectId: "proj_456"
    description: "Link to related project records"
  }) {
    id
    name
    type
    referenceProjectId
  }
}
```

## Advanced Example

Create a reference field with filtering and multiple selection:

```graphql
mutation CreateFilteredReferenceField {
  createCustomField(input: {
    name: "Dependencies"
    type: REFERENCE
    referenceProjectId: "proj_456"
    referenceMultiple: true
    referenceFilter: {
      status: ACTIVE
      tags: ["dependency"]
    }
    description: "Select multiple dependency records from the project"
  }) {
    id
    name
    type
    referenceProjectId
    referenceMultiple
    referenceFilter
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the reference field |
| `type` | CustomFieldType! | ✅ Yes | Must be `REFERENCE` |
| `referenceProjectId` | String! | ✅ Yes | ID of the project to reference |
| `referenceMultiple` | Boolean | No | Allow multiple record selection (default: false) |
| `referenceFilter` | TodoFilterInput | No | Filter criteria for referenced records |
| `description` | String | No | Help text shown to users |

**Note**: Custom fields are automatically associated with the project based on the user's current project context.

## Reference Configuration

### Single vs Multiple References

**Single Reference (default):**
```graphql
{
  referenceMultiple: false  # or omit this field
}
```
- Users can select one record from the referenced project
- Returns a single Todo object

**Multiple References:**
```graphql
{
  referenceMultiple: true
}
```
- Users can select multiple records from the referenced project
- Returns an array of Todo objects

### Reference Filtering

Use `referenceFilter` to limit which records can be selected:

```graphql
{
  referenceFilter: {
    status: ACTIVE
    assigneeIds: ["user_123"]
    tags: ["important"]
    dueDateFrom: "2024-01-01"
    dueDateTo: "2024-12-31"
  }
}
```

## Setting Reference Values

### Single Reference

```graphql
mutation SetSingleReference {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    customFieldReferenceTodoIds: ["referenced_todo_789"]
  })
}
```

### Multiple References

```graphql
mutation SetMultipleReferences {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    customFieldReferenceTodoIds: [
      "referenced_todo_789",
      "referenced_todo_012",
      "referenced_todo_345"
    ]
  })
}
```

### SetTodoCustomFieldInput Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record to update |
| `customFieldId` | String! | ✅ Yes | ID of the reference custom field |
| `customFieldReferenceTodoIds` | [String!] | ✅ Yes | Array of referenced record IDs |

## Creating Records with References

```graphql
mutation CreateRecordWithReference {
  createTodo(input: {
    title: "Implementation Task"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "reference_field_id"
      value: "referenced_todo_789"
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
      selectedTodos {
        id
        title
        status
      }
    }
  }
}
```

## Response Fields

### TodoCustomField Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier for the field value |
| `customField` | CustomField! | The reference field definition |
| `selectedTodos` | [Todo!] | Array of referenced records |
| `todo` | Todo! | The record this value belongs to |
| `createdAt` | DateTime! | When the value was created |
| `updatedAt` | DateTime! | When the value was last modified |

### Referenced Todo Fields

Each referenced Todo includes:

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Unique identifier of the referenced record |
| `title` | String! | Title of the referenced record |
| `status` | TodoStatus! | Current status (ACTIVE, COMPLETED, etc.) |
| `description` | String | Description of the referenced record |
| `dueDate` | DateTime | Due date if set |
| `assignees` | [User!] | Assigned users |
| `tags` | [Tag!] | Associated tags |
| `project` | Project! | Project containing the referenced record |

## Querying Reference Data

### Basic Query

```graphql
query GetRecordsWithReferences {
  todos(projectId: "project_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
      }
      selectedTodos {
        id
        title
        status
        project {
          id
          name
        }
      }
    }
  }
}
```

### Advanced Query with Nested Data

```graphql
query GetDetailedReferences {
  todos(projectId: "project_123") {
    id
    title
    customFields {
      id
      customField {
        name
        type
        referenceProjectId
        referenceMultiple
      }
      selectedTodos {
        id
        title
        description
        status
        dueDate
        assignees {
          id
          name
          email
        }
        tags {
          id
          name
          color
        }
        project {
          id
          name
        }
      }
    }
  }
}
```

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create reference field | `OWNER` or `ADMIN` role at project level |
| Update reference field | `OWNER` or `ADMIN` role at project level |
| Set reference value | Standard record edit permissions |
| View reference value | Standard record view permissions |
| Access referenced records | View permissions on referenced project |

**Important**: Users must have view permissions on the referenced project to see the linked records.

## Cross-Project Access

### Project Visibility

- Users can only reference records from projects they have access to
- Referenced records respect the original project's permissions
- Changes to referenced records appear in real-time
- Deleting referenced records removes them from reference fields

### Permission Inheritance

- Reference fields inherit permissions from both projects
- Users need view access to the referenced project
- Edit permissions are based on the current project's rules
- Referenced data is read-only in the context of the reference field

## Error Responses

### Invalid Reference Project

```json
{
  "errors": [{
    "message": "Project not found",
    "extensions": {
      "code": "PROJECT_NOT_FOUND"
    }
  }]
}
```

### Referenced Record Not Found

```json
{
  "errors": [{
    "message": "Custom field not found",
    "extensions": {
      "code": "CUSTOM_FIELD_NOT_FOUND"
    }
  }]
}
```

### Permission Denied

```json
{
  "errors": [{
    "message": "Forbidden",
    "extensions": {
      "code": "FORBIDDEN"
    }
  }]
}
```

## Best Practices

### Field Design

1. **Clear naming** - Use descriptive names that indicate the relationship
2. **Appropriate filtering** - Set filters to show only relevant records
3. **Consider permissions** - Ensure users have access to referenced projects
4. **Document relationships** - Provide clear descriptions of the connection

### Performance Considerations

1. **Limit reference scope** - Use filters to reduce the number of selectable records
2. **Avoid deep nesting** - Don't create complex chains of references
3. **Consider caching** - Referenced data is cached for performance
4. **Monitor usage** - Track how references are being used across projects

### Data Integrity

1. **Handle deletions** - Plan for when referenced records are deleted
2. **Validate permissions** - Ensure users can access referenced projects
3. **Update dependencies** - Consider impact when changing referenced records
4. **Audit trails** - Track reference relationships for compliance

## Common Use Cases

### Project Dependencies

```graphql
# Link to prerequisite tasks in other projects
{
  name: "Prerequisites"
  type: REFERENCE
  referenceProjectId: "infrastructure_project"
  referenceMultiple: true
  referenceFilter: {
    status: COMPLETED
    tags: ["prerequisite"]
  }
}
```

### Client Requirements

```graphql
# Reference client requirements from a requirements project
{
  name: "Client Requirements"
  type: REFERENCE
  referenceProjectId: "requirements_project"
  referenceFilter: {
    assigneeIds: ["client_user_id"]
    status: ACTIVE
  }
}
```

### Resource Allocation

```graphql
# Link to resource records in a resource management project
{
  name: "Assigned Resources"
  type: REFERENCE
  referenceProjectId: "resources_project"
  referenceMultiple: true
  referenceFilter: {
    tags: ["available"]
  }
}
```

### Quality Assurance

```graphql
# Reference QA test cases from a testing project
{
  name: "Test Cases"
  type: REFERENCE
  referenceProjectId: "qa_project"
  referenceMultiple: true
  referenceFilter: {
    status: ACTIVE
    tags: ["test-case"]
  }
}
```

## Integration with Lookups

Reference fields work seamlessly with [Lookup fields](/api/custom-fields/lookup) to pull data from referenced records:

```graphql
# Reference field links to records
{
  name: "Related Tasks"
  type: REFERENCE
  referenceProjectId: "other_project"
}

# Lookup field pulls data from referenced records
{
  name: "Total Budget"
  type: LOOKUP
  lookupOption: {
    customFieldId: "related_tasks_field_id"
    targetField: "budget"
    function: SUM
  }
}
```

## Limitations

- Maximum 100 records can be referenced in a single field
- Referenced projects must be accessible to the user
- Circular references between projects are not allowed
- Reference fields cannot point to records in the same project
- Changes to referenced project permissions affect reference field access
- Deep nesting of references may impact performance

## Related Resources

- [Lookup Fields](/api/custom-fields/lookup) - Pull data from referenced records
- [Projects API](/api/projects) - Managing projects that contain references
- [Records API](/api/records) - Working with records that have references
- [Custom Fields Overview](/api/custom-fields) - General concepts