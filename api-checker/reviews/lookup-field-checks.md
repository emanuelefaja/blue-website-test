# Verification for: Lookup Custom Field
Path: /content/en/api/5.custom fields/lookup.md
Status: [🔄] In Progress

## 1. GraphQL Schema Verification

### Custom Field Type
- [🔄] Verify `LOOKUP` exists in CustomFieldType enum
  - Location in schema: [searching...]
  - Actual vs Documented: [pending verification]

### Mutations
- [🔄] Verify `createCustomField` mutation supports LOOKUP type
  - Input parameters for lookup: [checking...]
  - Response fields: [checking...]

### Lookup Configuration
- [🔄] Verify `lookupOption` field exists in CreateCustomFieldInput
- [🔄] Verify `CustomFieldLookupOptionInput` type exists
- [🔄] Verify `CustomFieldLookupType` enum exists

## 2. Input Parameter Verification

### CustomFieldLookupOptionInput Fields
- [🔄] `referenceId` field - type and requirement
- [🔄] `lookupId` field - type and requirement  
- [🔄] `lookupType` field - type and requirement
- [🔄] Check for additional/missing fields mentioned in docs

## 3. Lookup Type Enum Verification

### CustomFieldLookupType Values
- [🔄] `TODO_DUE_DATE` - verify exists
- [🔄] `TODO_CREATED_AT` - verify exists
- [🔄] `TODO_UPDATED_AT` - verify exists
- [🔄] `TODO_TAG` - verify exists
- [🔄] `TODO_ASSIGNEE` - verify exists
- [🔄] `TODO_DESCRIPTION` - verify exists
- [🔄] `TODO_LIST` - verify exists
- [🔄] `TODO_CUSTOM_FIELD` - verify exists

## 4. Response Type Verification

### CustomFieldLookupOption Response
- [🔄] `lookupType` field exists
- [🔄] `lookupResult` field exists and type
- [🔄] `reference` field exists
- [🔄] `lookup` field exists
- [🔄] `parentCustomField` field exists

## 5. Business Logic Verification

### Lookup Behavior
- [🔄] Read-only nature of lookup fields
- [🔄] Automatic calculation and updates
- [🔄] Maximum records limit (1000)
- [🔄] Circular dependency prevention

### Supported Custom Field Types for TODO_CUSTOM_FIELD
- [🔄] Verify which custom field types can be looked up
- [🔄] Check actual implementation vs documented types

## 6. Permission Verification

### Create/Update Permissions
- [🔄] OWNER can create/update lookup fields
- [🔄] ADMIN can create/update lookup fields
- [🔄] View permissions required for referenced project

## 7. Error Code Verification

### Documented Error Codes
- [🔄] `CUSTOM_FIELD_NOT_FOUND` - verify usage
- [🔄] `VALIDATION_ERROR` for circular lookup - verify
- [🔄] `PROJECT_NOT_FOUND` - verify usage

## 8. Documentation Issues

### Suspicious Content
- [🔄] Complex lookup examples with filters and functions (SUM, COUNT, etc.)
- [🔄] Display options and formatting
- [🔄] Target field paths like "customFields.budget"
- [🔄] Filter objects in lookup configuration

## 9. Link Verification

### Internal API Links
- [🔄] `/api/custom-fields/reference` - exists?
- [🔄] `/api/custom-fields/formula` - exists?
- [🔄] `/api/custom-fields/number` - exists?
- [🔄] `/api/custom-fields` - exists?

## Summary

### Critical Issues (Must Fix)
[To be populated after verification]

### Minor Issues (Should Fix)
[To be populated after verification]

### Suggestions
[To be populated after verification]