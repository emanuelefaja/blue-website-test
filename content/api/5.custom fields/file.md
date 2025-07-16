---
title: File Custom Field
description: Create file fields to attach documents, images, and other files to records
category: Custom Fields
---

File custom fields allow you to attach multiple files to records. Files are stored securely in AWS S3 with comprehensive metadata tracking, file type validation, and proper access controls.

## Basic Example

Create a simple file field:

```graphql
mutation CreateFileField {
  createCustomField(input: {
    name: "Attachments"
    type: FILE
  }) {
    id
    name
    type
  }
}
```

## Advanced Example

Create a file field with description:

```graphql
mutation CreateDetailedFileField {
  createCustomField(input: {
    name: "Project Documents"
    type: FILE
    description: "Upload project-related documents, images, and files"
  }) {
    id
    name
    type
    description
  }
}
```

## Input Parameters

### CreateCustomFieldInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | String! | ✅ Yes | Display name of the file field |
| `type` | CustomFieldType! | ✅ Yes | Must be `FILE` |
| `description` | String | No | Help text shown to users |

**Note**: Custom fields are automatically associated with the project based on the user's current project context. No `projectId` parameter is required.

## File Upload Process

### Step 1: Upload File

First, upload the file to get a file UID:

```graphql
mutation UploadFile {
  uploadFile(input: {
    file: $file  # File upload variable
    companyId: "company_123"
    projectId: "proj_123"
  }) {
    id
    uid
    name
    size
    type
    extension
    status
  }
}
```

### Step 2: Attach File to Record

Then attach the uploaded file to a record:

```graphql
mutation AttachFileToRecord {
  createTodoCustomFieldFile(input: {
    todoId: "todo_123"
    customFieldId: "file_field_456"
    fileUid: "file_uid_from_upload"
  }) {
    id
    file {
      uid
      name
      size
      type
    }
  }
}
```

## Managing File Attachments

### Adding Single Files

```graphql
mutation AddFileToField {
  createTodoCustomFieldFile(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    fileUid: "file_uid_789"
  }) {
    id
    position
    file {
      uid
      name
      size
      type
      extension
    }
  }
}
```

### Removing Files

```graphql
mutation RemoveFileFromField {
  deleteTodoCustomFieldFile(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    fileUid: "file_uid_789"
  })
}
```

### Bulk File Operations

Update multiple files at once using customFieldOptionIds:

```graphql
mutation SetMultipleFiles {
  setTodoCustomField(input: {
    todoId: "todo_123"
    customFieldId: "field_456"
    customFieldOptionIds: ["file_uid_1", "file_uid_2", "file_uid_3"]
  })
}
```

## File Upload Input Parameters

### UploadFileInput

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | Upload! | ✅ Yes | File to upload |
| `companyId` | String! | ✅ Yes | Company ID for file storage |
| `projectId` | String | No | Project ID for project-specific files |

### File Management Input Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `todoId` | String! | ✅ Yes | ID of the record |
| `customFieldId` | String! | ✅ Yes | ID of the file custom field |
| `fileUid` | String! | ✅ Yes | Unique identifier of the uploaded file |

## File Storage and Limits

### File Size Limits

| Limit Type | Size |
|------------|------|
| Maximum file size | 256MB per file |
| Batch upload limit | 10 files max, 1GB total |
| GraphQL upload limit | 256MB |

### Supported File Types

#### Images
- `jpg`, `jpeg`, `png`, `gif`, `bmp`, `webp`, `svg`, `ico`, `tiff`, `tif`

#### Videos
- `mp4`, `avi`, `mov`, `wmv`, `flv`, `webm`, `mkv`, `3gp`

#### Audio
- `mp3`, `wav`, `flac`, `aac`, `ogg`, `wma`

#### Documents
- `pdf`, `doc`, `docx`, `xls`, `xlsx`, `ppt`, `pptx`, `txt`, `rtf`

#### Archives
- `zip`, `rar`, `7z`, `tar`, `gz`

#### Code/Text
- `json`, `xml`, `csv`, `md`, `yaml`, `yml`

### Storage Architecture

- **Storage**: AWS S3 with organized folder structure
- **Path Format**: `companies/{companySlug}/projects/{projectSlug}/uploads/{year}/{month}/{username}/{fileUid}_{filename}`
- **Security**: Signed URLs for secure access
- **Backup**: Automatic S3 redundancy

## Response Fields

### File Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Database ID |
| `uid` | String! | Unique file identifier |
| `name` | String! | Original filename |
| `size` | Int! | File size in bytes |
| `type` | String! | MIME type |
| `extension` | String! | File extension |
| `status` | FileStatus! | PENDING or CONFIRMED |
| `shared` | Boolean! | Whether file is shared |
| `createdAt` | DateTime! | Upload timestamp |

### TodoCustomFieldFile Response

| Field | Type | Description |
|-------|------|-------------|
| `id` | String! | Junction record ID |
| `uid` | String! | Unique identifier |
| `position` | Float! | Display order |
| `file` | File! | Associated file object |
| `todoCustomField` | TodoCustomField! | Parent custom field |
| `createdAt` | DateTime! | When file was attached |

## Creating Records with Files

When creating records, you can attach files using their UIDs:

```graphql
mutation CreateRecordWithFiles {
  createTodo(input: {
    title: "Project deliverables"
    todoListId: "list_123"
    customFields: [{
      customFieldId: "file_field_id"
      customFieldOptionIds: ["file_uid_1", "file_uid_2"]
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
      todoCustomFieldFiles {
        id
        position
        file {
          uid
          name
          size
          type
        }
      }
    }
  }
}
```

## File Validation and Security

### Upload Validation

- **MIME Type Checking**: Validates against allowed types
- **File Extension Validation**: Fallback for `application/octet-stream`
- **Size Limits**: Enforced at upload time
- **Filename Sanitization**: Removes special characters

### Access Control

- **Upload Permissions**: Project/company membership required
- **File Association**: ADMIN, OWNER, MEMBER, CLIENT roles
- **File Access**: Inherited from project/company permissions
- **Secure URLs**: Time-limited signed URLs for file access

## Required Permissions

| Action | Required Permission |
|--------|-------------------|
| Create file field | `CUSTOM_FIELDS_CREATE` at company or project level |
| Update file field | `CUSTOM_FIELDS_UPDATE` at company or project level |
| Upload files | Project or company membership |
| Attach files | ADMIN, OWNER, MEMBER, or CLIENT role |
| View files | Standard record view permissions |
| Delete files | Same as attach permissions |

## Error Responses

### File Too Large
```json
{
  "errors": [{
    "message": "File \"filename.pdf\": Size exceeds maximum limit of 256MB",
    "extensions": {
      "code": "BAD_USER_INPUT"
    }
  }]
}
```

### File Not Found
```json
{
  "errors": [{
    "message": "File not found",
    "extensions": {
      "code": "FILE_NOT_FOUND"
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

### File Management
- Upload files before attaching to records
- Use descriptive filenames
- Organize files by project/purpose
- Clean up unused files periodically

### Performance
- Upload files in batches when possible
- Use appropriate file formats for content type
- Compress large files before upload
- Consider file preview requirements

### Security
- Validate file contents, not just extensions
- Use virus scanning for uploaded files
- Implement proper access controls
- Monitor file upload patterns

## Common Use Cases

1. **Document Management**
   - Project specifications
   - Contracts and agreements
   - Meeting notes and presentations
   - Technical documentation

2. **Asset Management**
   - Design files and mockups
   - Brand assets and logos
   - Marketing materials
   - Product images

3. **Compliance and Records**
   - Legal documents
   - Audit trails
   - Certificates and licenses
   - Financial records

4. **Collaboration**
   - Shared resources
   - Version-controlled documents
   - Feedback and annotations
   - Reference materials

## Integration Features

### With Automations
- Trigger actions when files are added/removed
- Process files based on type or metadata
- Send notifications for file changes
- Archive files based on conditions

### With Cover Images
- Use file fields as cover image sources
- Automatic image processing and thumbnails
- Dynamic cover updates when files change

### With Lookups
- Reference files from other records
- Aggregate file counts and sizes
- Find records by file metadata
- Cross-reference file attachments

## Limitations

- Maximum 256MB per file
- Dependent on S3 availability
- No built-in file versioning
- No automatic file conversion
- Limited file preview capabilities
- No real-time collaborative editing

## Related Resources

- [Upload Files API](/api/files/upload) - File upload endpoints
- [Custom Fields Overview](/custom-fields/list-custom-fields) - General concepts
- [Automations API](/api/automations/index) - File-based automations
- [AWS S3 Documentation](https://docs.aws.amazon.com/s3/) - Storage backend