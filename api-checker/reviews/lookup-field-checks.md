# Verification for: Lookup Custom Field
Path: /content/en/api/5.custom fields/lookup.md
Status: [✅] Completed

## 1. GraphQL Schema Verification

### Custom Field Type
- [✅] Verify `LOOKUP` exists in CustomFieldType enum
  - Location in schema: `/bloo-api/src/schema.graphql:4521`
  - Actual vs Documented: **MATCHES**

### Mutations
- [✅] Verify `createCustomField` mutation supports LOOKUP type
  - Input parameters for lookup: **lookupOption: CustomFieldLookupOptionInput**
  - Response fields: **Standard CustomField response**

### Lookup Configuration
- [✅] Verify `lookupOption` field exists in CreateCustomFieldInput - **YES** line 2766
- [✅] Verify `CustomFieldLookupOptionInput` type exists - **YES** line 2772
- [✅] Verify `CustomFieldLookupType` enum exists - **YES** line 4752

## 2. Input Parameter Verification

### CustomFieldLookupOptionInput Fields
- [✅] `referenceId` field - **String!** (required)
- [✅] `lookupId` field - **String** (optional, used for TODO_CUSTOM_FIELD)
- [✅] `lookupType` field - **CustomFieldLookupType!** (required)
- [❌] Check for additional/missing fields - **MANY HALLUCINATED FIELDS**:
  - targetField (doesn't exist)
  - function (doesn't exist)
  - filter (doesn't exist)
  - display (doesn't exist)

## 3. Lookup Type Enum Verification

### CustomFieldLookupType Values
- [✅] `TODO_DUE_DATE` - **EXISTS**
- [✅] `TODO_CREATED_AT` - **EXISTS**
- [✅] `TODO_UPDATED_AT` - **EXISTS**
- [✅] `TODO_TAG` - **EXISTS**
- [✅] `TODO_ASSIGNEE` - **EXISTS**
- [✅] `TODO_DESCRIPTION` - **EXISTS**
- [✅] `TODO_LIST` - **EXISTS**
- [✅] `TODO_CUSTOM_FIELD` - **EXISTS**

All 8 enum values documented are correct and exist in the schema.

## 4. Response Type Verification

### CustomFieldLookupOption Response
- [✅] `lookupType` field exists - **CustomFieldLookupType**
- [✅] `lookupResult` field exists and type - **JSON**
- [✅] `reference` field exists - **CustomField**
- [✅] `lookup` field exists - **CustomField**
- [✅] `parentCustomField` field exists - **CustomField**
- [✅] Additional field found: `parentLookup` - **CustomField**
- [✅] Deprecated field: `lookupValues` - **JSON** (deprecated)

## 5. Business Logic Verification

### Lookup Behavior
- [✅] Read-only nature of lookup fields - **CONFIRMED**
- [✅] Automatic calculation and updates - **CONFIRMED**
- [❌] Maximum records limit (1000) - **NOT FOUND** in implementation
- [✅] Circular dependency prevention - **CONFIRMED** with ValidationError

### Supported Custom Field Types for TODO_CUSTOM_FIELD
- [✅] Verify which custom field types can be looked up - **ALL TYPES SUPPORTED**
- [❌] Check actual implementation vs documented types - **NO AGGREGATION FUNCTIONS**:
  - No SUM, COUNT, AVERAGE, MAX, MIN
  - No filtering capabilities
  - No display formatting options
  - Simple data extraction only

## 6. Permission Verification

### Create/Update Permissions
- [✅] OWNER can create/update lookup fields - **CONFIRMED**
- [✅] ADMIN can create/update lookup fields - **CONFIRMED**
- [✅] View permissions required for referenced project - **CONFIRMED**

## 7. Error Code Verification

### Documented Error Codes
- [✅] `CUSTOM_FIELD_NOT_FOUND` - **EXISTS** and used correctly
- [⚠️] `VALIDATION_ERROR` for circular lookup - **Uses BAD_USER_INPUT** instead
- [✅] `PROJECT_NOT_FOUND` - **EXISTS** and used correctly

## 8. Documentation Issues

### Suspicious Content
- [❌] Complex lookup examples with filters and functions - **COMPLETELY HALLUCINATED**
  - No SUM, COUNT, AVERAGE, MAX, MIN functions exist
  - No filter parameter exists
  - No targetField parameter exists
- [❌] Display options and formatting - **DOESN'T EXIST**
- [❌] Target field paths like "customFields.budget" - **WRONG PATTERN**
- [❌] Filter objects in lookup configuration - **NOT SUPPORTED**

## 9. Link Verification

### Internal API Links
- [✅] `/api/custom-fields/reference` - **EXISTS**
- [✅] `/api/custom-fields/formula` - **EXISTS**
- [✅] `/api/custom-fields/number` - **EXISTS**
- [❌] `/api/custom-fields` - **WRONG PATH** - should be `/api/custom-fields/list-custom-fields`

## Summary

### Critical Issues (Must Fix)
1. **Massive hallucination of features**: Documentation describes aggregation functions (SUM, COUNT, AVERAGE, etc.) that don't exist
2. **Non-existent parameters**: targetField, function, filter, display options are all made up
3. **Wrong implementation model**: Docs suggest complex calculations, but lookups only extract data
4. **Incorrect examples**: Most code examples use non-existent features

### Minor Issues (Should Fix)
1. **Error code**: Uses BAD_USER_INPUT instead of VALIDATION_ERROR for circular dependencies
2. **Missing field**: Documentation doesn't mention `parentLookup` field
3. **Wrong link**: `/api/custom-fields` should be `/api/custom-fields/list-custom-fields`
4. **No mention of deprecated field**: `lookupValues` is deprecated

### Suggestions
1. **Complete rewrite needed**: Remove all aggregation function content
2. **Simplify examples**: Show only data extraction, not calculations
3. **Correct the mental model**: Lookups extract data, they don't calculate
4. **Fix all code examples**: Remove non-existent parameters

### Overall Assessment
This documentation is **severely inaccurate** (only ~30% correct). It describes a completely different feature than what's implemented. The actual lookup fields are much simpler - they only extract and display data from referenced records without any aggregation, filtering, or calculation capabilities.