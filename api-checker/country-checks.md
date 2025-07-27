# Verification for: country.md
Path: /content/en/api/5.custom fields/country.md
Status: [❌] Has Issues - Multiple critical discrepancies between documentation and implementation

## 1. Country Validation Behavior

### setTodoCustomField Mutation
- [❌] **Documentation claims**: System validates country inputs
- [❌] **Reality**: NO validation occurs - accepts any string without checking
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/setTodoCustomField.ts (line 59)
  - Code: `countryCodes: countryCodes?.join(',')` - just joins array to string

### createTodo Mutation  
- [✓] Country validation ONLY happens here
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/createTodo.ts (lines 240-252)
  - Uses `i18n-iso-countries` library to validate and convert
  - Converts country names to ISO codes

## 2. Storage Format Issues

### Database Storage
- [❌] **Documentation claims**: `countryCodes` stores as array
- [❌] **Reality**: Stored as comma-separated STRING in database
  - Joined on save: `countryCodes?.join(',')`
  - Split on read: `countryCodes?.split(',')` in GraphQL resolver

### Field Independence
- [❌] **Documentation claims**: Can use `countryCodes` and `text` independently
- [❌] **Reality in createTodo**: 
  - `text` is automatically set to original input value
  - `countryCodes` is set to validated/converted code
  - Cannot control them independently

## 3. Multiple Country Support

### setTodoCustomField
- [✓] Supports multiple countries via `countryCodes` array parameter
- [✓] Stores as comma-separated string

### createTodo
- [❌] **Documentation claims**: Accepts comma-separated countries in value field
- [❌] **Reality**: Treats entire input as single country
  - No splitting by commas
  - No multiple country support in value field

## 4. Behavioral Differences Between Mutations

### Critical Undocumented Difference
- [❌] Documentation doesn't explain the major behavioral difference:
  - **createTodo**: Validates, converts names to codes, auto-populates both fields
  - **setTodoCustomField**: No validation, stores raw values as provided

## 5. Error Handling

### Error Messages
- [❌] **Documentation shows**: `"Invalid country value: 'Atlantis'"`
- [❌] **Reality**: `"Invalid country value."` (no specific value shown)
  - Location: createTodo.ts line 248

### Error Code
- [✓] `CUSTOM_FIELD_VALUE_PARSE_ERROR` is correct

## 6. Permission System

### Permission Constants
- [❌] **Documentation claims**: `CUSTOM_FIELDS_CREATE` permission required
- [❌] **Reality**: Uses standard project roles (OWNER, ADMIN, MEMBER, CLIENT)
  - No specific `CUSTOM_FIELDS_CREATE` constant found in codebase

## 7. Input Parameters

### createTodo Value Field
- [❌] Documentation doesn't explain that country values use generic `value` field
- [❌] Missing crucial detail about customFields array structure

## 8. Country Standards

### ISO Standards
- [✓] Correctly documents ISO 3166-1 Alpha-2 usage
- [✓] i18n-iso-countries library confirmed in implementation

## Summary

### Critical Issues (Must Fix)
1. **Validation behavior is completely different between mutations** - Not documented
2. **setTodoCustomField has NO validation** - Documentation incorrectly claims validation
3. **Storage format is string, not array** - Fundamental misrepresentation
4. **createTodo doesn't support multiple countries** - False claim about comma-separated input
5. **Permission system uses different constants** - Wrong permission names

### Moderate Issues
1. Error messages don't include the invalid value
2. `text` and `countryCodes` cannot be set independently in createTodo
3. Missing explanation of `value` field usage in createTodo

### Accurate Elements
- Country code standards (ISO 3166-1 Alpha-2)
- Basic field creation syntax
- General use cases and best practices

### Overall Assessment
The documentation presents an idealized, unified behavior that doesn't exist. In reality, the two mutations (`createTodo` and `setTodoCustomField`) have completely different implementations for country fields. This is a significant documentation failure that could lead to data quality issues and developer confusion.

**Recommendation**: Major rewrite needed to accurately represent the two different behaviors and clarify when validation occurs.