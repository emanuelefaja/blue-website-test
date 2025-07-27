# Verification for: 2.create-project.md
Path: /content/en/api/2.projects/2.create-project.md
Status: [ ] In Progress / [✓] Completed - FIXED

## 1. GraphQL Schema Verification

### Mutation/Query Name
- [✓] Verify the GraphQL operation name exists in schema
  - Operation: `createProject`
  - Location in schema: /Users/manny/Blue/bloo-api/src/schema.graphql:2295
  - Actual vs Documented: MATCHES - `createProject(input: CreateProjectInput!): Project`

### Input Type Verification
- [✓] Verify input type name is correct
  - Documented: `CreateProjectInput`
  - Actual in schema: CreateProjectInput (MATCHES)
  - Location: /Users/manny/Blue/bloo-api/src/schema.graphql (input section)

### Input Parameters
For each parameter in the input:
- [✓] `name`
  - Documented type: `String` (Required: ✅ Yes)
  - Actual type: `String!` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes - "The project name (max 50 characters). URLs will be stripped from the name."
  - Default value (if any): None

- [✓] `companyId`
  - Documented type: `String` (Required: ✅ Yes)
  - Actual type: `String!` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes - "The ID or slug of the company where the project will be created."
  - Default value (if any): None

- [✓] `description`
  - Documented type: `String` (Required: No)
  - Actual type: `String` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes
  - Default value (if any): None

- [✓] `color`
  - Documented type: `String` (Required: No)
  - Actual type: `String` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes - "Project color in hex format"
  - Default value (if any): None

- [✓] `icon`
  - Documented type: `String` (Required: No)
  - Actual type: `String` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes - "Icon identifier for the project"
  - Default value (if any): None

- [✓] `category`
  - Documented type: `ProjectCategory` (Required: No)
  - Actual type: `ProjectCategory` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes - "Defaults to GENERAL if not specified"
  - Default value (if any): GENERAL (correctly documented)

- [✓] `templateId`
  - Documented type: `String` (Required: No)
  - Actual type: `String` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes
  - Default value (if any): None

- [✓] `coverConfig` ⚠️ IMPLEMENTATION BUG FOUND
  - Documented type: `TodoCoverConfigInput` (Required: No)
  - Actual type: `TodoCoverConfigInput` (MATCHES)
  - Required status matches: Yes
  - Description accurate: Yes
  - Default value (if any): None
  - **BUG**: Parameter is accepted but NOT USED when creating non-template projects

### Response Fields
For each field in the response:
- [✓] `id`
  - Documented type: `String!`
  - Actual type: `ID!` (from Project type)
  - Is field actually returned: Yes

- [✓] `name`
  - Documented type: `String!`
  - Actual type: `String!` (MATCHES)
  - Is field actually returned: Yes

- [✓] `slug`
  - Documented type: `String!`
  - Actual type: `String!` (MATCHES)
  - Is field actually returned: Yes

- [✓] `description`
  - Documented type: `String`
  - Actual type: `String` (MATCHES)
  - Is field actually returned: Yes

- [✓] `color`
  - Documented type: `String`
  - Actual type: `String` (MATCHES)
  - Is field actually returned: Yes

- [✓] `icon`
  - Documented type: `String`
  - Actual type: `String` (MATCHES)
  - Is field actually returned: Yes

- [✓] `category`
  - Documented type: `String`
  - Actual type: `ProjectCategory` (ENUM TYPE - documentation should show this as ProjectCategory not String)
  - Is field actually returned: Yes

## 2. Enum Verification

For each enum mentioned:
- [✓] `ProjectCategory`
  - Exists in schema: Yes
  - Location: /Users/manny/Blue/bloo-api/src/generated/types.ts
  - Values match:
    - [✓] `CRM` - exists
    - [✓] `CROSS_FUNCTIONAL` - exists
    - [✓] `CUSTOMER_SUCCESS` - exists
    - [✓] `DESIGN` - exists
    - [✓] `ENGINEERING` - exists
    - [✓] `GENERAL` - exists
    - [✓] `HR` - exists
    - [✓] `IT` - exists
    - [✓] `MARKETING` - exists
    - [✓] `OPERATIONS` - exists
    - [✓] `PRODUCT` - exists
    - [✓] `SALES` - exists
  - Missing values in docs: None
  - Extra values in docs: None

- [✓] `ImageFit`
  - Exists in schema: Yes
  - Location: /Users/manny/Blue/bloo-api/src/schema.graphql
  - Values match:
    - [✓] `COVER` - exists
    - [✓] `CONTAIN` - exists
    - [✓] `FILL` - exists
    - [✓] `SCALE_DOWN` - exists
  - Missing values in docs: None
  - Extra values in docs: None

- [✓] `ImageSelectionType`
  - Exists in schema: Yes
  - Location: /Users/manny/Blue/bloo-api/src/schema.graphql
  - Values match:
    - [✓] `FIRST` - exists
    - [✓] `LAST` - exists
    - [❌] `SPECIFIC` - DOES NOT EXIST IN SCHEMA
  - Missing values in docs: None
  - Extra values in docs: `SPECIFIC` (this is hallucinated - not in actual schema)

- [✓] `ImageSource`
  - Exists in schema: Yes
  - Location: /Users/manny/Blue/bloo-api/src/schema.graphql
  - Values match:
    - [✓] `DESCRIPTION` - exists
    - [✓] `COMMENTS` - exists
    - [✓] `CUSTOM_FIELD` - exists
  - Missing values in docs: None
  - Extra values in docs: None

## 3. Implementation Verification

### Resolver Check
- [✓] Resolver exists for this operation
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/createProject.ts
  - Handler function: `createProject` (default export)

### Business Logic Verification
- [✓] All documented parameters are actually used
  - [✓] `name` - used in: line 56-75 (project creation)
  - [✓] `companyId` - used in: line 26 (company lookup), line 56-75 (project creation)
  - [✓] `description` - used in: line 56-75 (project creation)
  - [✓] `color` - used in: line 56-75 (project creation)
  - [✓] `icon` - used in: line 56-75 (project creation)
  - [✓] `category` - used in: line 56-75 (project creation with default to GENERAL)
  - [✓] `templateId` - used in: line 35 (template lookup), line 139 (copy queue)
  - [❌] `coverConfig` - **BUG: NOT USED for non-template projects** (extracted on line 19 but never passed to project creation)

### Validation Rules
- [✓] Required fields enforced in code (TypeScript/GraphQL enforces required fields)
- [❌] Max length/size limits match documentation - **NOT ENFORCED IN CODE** (docs say max 50 chars for name)
- [✓] Format validation - URLs are stripped from project name (line 61)

## 4. Permission Verification

### Required Permissions
- [✓] Permission checks exist in resolver
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/createProject.ts:32
  - Documented roles match code: Yes
  
### Role-based Access
For each role mentioned:
- [✓] `OWNER` - can perform: matches docs
- [✓] `ADMIN` - can perform: matches docs  
- [✓] `MEMBER` - can perform: matches docs

Code uses: `requireCompanyRole(companyId, userId, ['OWNER', 'ADMIN', 'MEMBER'])`

## 5. Error Response Verification

### Error Codes
For each error code documented:
- [✓] Template size limit error (250,000 todos)
  - Exists in codebase: Yes
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/createProject.ts:46-48
  - Message matches: Yes - "Template cannot have more than 250000 todos"

## 6. Link Verification

### Internal API Links
- [✓] All links to other API pages are valid
  - [✓] No internal links found in this doc

### Related Endpoints
- [✓] All mentioned related endpoints exist
  - [✓] `copyProjectStatus` query - EXISTS in /Users/manny/Blue/bloo-api/src/resolvers/Query/copyProjectStatus.ts

## 7. Code Example Verification

### Basic Example
- [✓] GraphQL syntax is valid
- [✓] All fields in query/mutation exist
- [✓] Required fields are included
- [✓] Response structure matches actual response

### Advanced Example
- [✓] All optional parameters shown actually exist
- [✓] Nested objects structure is correct (coverConfig)
- [✓] Complex types (TodoCoverConfigInput) have correct structure

## 8. Documentation Completeness

### Missing from Docs
- [✓] List any parameters found in code but not documented
  - None - all parameters are documented
- [✓] List any response fields found but not documented
  - Many Project type fields are available but not shown in examples (e.g., createdAt, updatedAt, companyId, etc.)
- [✓] List any error cases found but not documented
  - Company not found error
  - Template not found error
  - Permission denied errors

### Extra in Docs (Hallucinated)
- [✓] List any parameters documented but not in code
  - None
- [✓] List any fields documented but not in code
  - None
- [✓] List any features documented but not implemented
  - `SPECIFIC` value for ImageSelectionType enum (hallucinated)
  - coverConfig functionality for non-template projects (parameter accepted but not implemented)

## 9. Special Considerations

### Database/Prisma Verification
- [✓] If mentions database fields, verify against Prisma schema
  - Model: `Project`
  - Fields match: Yes (name, description, color, icon, category all exist in Prisma schema)

### Type Definitions
- [✓] All TypeScript/GraphQL types mentioned exist
  - [✓] `TodoCoverConfigInput` - Found in schema.graphql

### Custom Field Types
- [✓] If mentions custom fields, verify types exist
  - [✓] N/A for this endpoint

## Summary

### Critical Issues (Must Fix) - ALL FIXED ✅
1. **HALLUCINATED ENUM VALUE**: `SPECIFIC` value for `ImageSelectionType` does not exist in schema - ✅ FIXED: Removed from docs
2. **IMPLEMENTATION BUG**: `coverConfig` parameter is accepted but NOT USED when creating non-template projects - ✅ FIXED: Added warning and clarified it only works with templates
3. **TYPE ERROR**: Response field `category` is documented as `String` but is actually `ProjectCategory` enum - ✅ FIXED: Added response fields table showing correct type

### Minor Issues (Should Fix) - ALL FIXED ✅
1. **MISSING VALIDATION**: Name max length (50 chars) is documented but not enforced in code - ✅ FIXED: Removed the unenforceable limit from docs
2. **INCOMPLETE RESPONSE DOCS**: Many available Project fields not shown in response examples - ✅ FIXED: Added complete response fields table
3. **MISSING ERROR DOCS**: Common errors (company not found, permission denied) not documented - ✅ FIXED: Added error responses section

### Suggestions - ALL IMPLEMENTED ✅
1. Add note that `coverConfig` currently only works with template-based project creation - ✅ DONE
2. Document all available Project response fields or link to Project type documentation - ✅ DONE
3. Add common error responses section - ✅ DONE
4. Fix the `category` response type to show `ProjectCategory` instead of `String` - ✅ DONE (in response fields table)
5. Remove `SPECIFIC` from ImageSelectionType enum values - ✅ DONE