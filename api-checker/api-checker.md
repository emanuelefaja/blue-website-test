# API Documentation Verification Prompt

## Overview
This prompt guides the systematic verification of API documentation against the actual Blue codebase. Since all API docs were AI-generated, we need to verify that every detail is accurate and nothing is hallucinated.

## Process

### 1. Setup
- Create `/api-checker` folder in the project root
- Create `api-tracker.md` to track progress through all API files
- Create `/api-checker/reviews` subfolder for individual verification files
- Process one API documentation file at a time from `/content/en/api/`
- **CRITICAL**: IMMEDIATELY mark the file as 🔄 In Progress in `api-tracker.md` BEFORE doing any other work
- **IMPORTANT**: Check `api-tracker.md` first and pick the next file that is NOT marked as 🔄 In Progress
- Only work on files with [ ] (not started) status
- **WORKFLOW**: 
  1. Choose a file with [ ] status
  2. IMMEDIATELY update api-tracker.md to mark it as [🔄] In Progress
  3. Only then proceed with verification
  4. Create the review checklist file and complete verification
  5. **IMMEDIATELY after completing the review file, update the actual API documentation file based on your findings**
  6. Update api-tracker.md to mark as [🔧] Fixed or [✅] Verified

### 2. For Each API Documentation File

Create a verification checklist file in `/api-checker/reviews/` named after the API doc (e.g., `create-project-checks.md` for `create-project.md`) with the following sections:

## Verification Checklist Template

```markdown
# Verification for: [API File Name]
Path: /content/en/api/[path/to/file.md]
Status: [ ] In Progress / [✓] Completed

## 1. GraphQL Schema Verification

### Mutation/Query Name
- [ ] Verify the GraphQL operation name exists in schema
  - Operation: `[operationName]`
  - Location in schema: [file:line]
  - Actual vs Documented: [any differences]

### Input Type Verification
- [ ] Verify input type name is correct
  - Documented: `[InputTypeName]`
  - Actual in schema: [actual name or "NOT FOUND"]
  - Location: [file:line]

### Input Parameters
For each parameter in the input:
- [ ] `parameterName`
  - Documented type: `[Type]`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]
  - Description accurate: [Yes/No/Needs update]
  - Default value (if any): [matches/different/not documented]

### Response Fields
For each field in the response:
- [ ] `fieldName`
  - Documented type: `[Type]`
  - Actual type: [actual or "NOT FOUND"]
  - Is field actually returned: [Yes/No]

## 2. Enum Verification

For each enum mentioned:
- [ ] `EnumName`
  - Exists in schema: [Yes/No]
  - Location: [file:line]
  - Values match:
    - [ ] `VALUE1` - [exists/missing/different]
    - [ ] `VALUE2` - [exists/missing/different]
  - Missing values in docs: [list any]
  - Extra values in docs: [list any]

## 3. Implementation Verification

### Resolver Check
- [ ] Resolver exists for this operation
  - Location: [file:line]
  - Handler function: `[functionName]`

### Business Logic Verification
- [ ] All documented parameters are actually used
  - [ ] `param1` - used in: [file:line or "NOT USED"]
  - [ ] `param2` - used in: [file:line or "NOT USED"]

### Validation Rules
- [ ] Required fields enforced in code
- [ ] Max length/size limits match documentation
- [ ] Format validation (URLs, hex colors, etc.) matches

## 4. Permission Verification

### Required Permissions
- [ ] Permission checks exist in resolver
  - Location: [file:line]
  - Documented roles match code: [Yes/No]
  
### Role-based Access
For each role mentioned:
- [ ] `OWNER` - can perform: [matches docs/different]
- [ ] `ADMIN` - can perform: [matches docs/different]
- [ ] `MEMBER` - can perform: [matches docs/different]

## 5. Error Response Verification

### Error Codes
For each error code documented:
- [ ] `ERROR_CODE`
  - Exists in codebase: [Yes/No]
  - Location: [file:line]
  - Message matches: [Yes/No]

## 6. Link Verification

### Internal API Links
- [ ] All links to other API pages are valid
  - [ ] `/api/path/to/page` - [exists/broken]

### Related Endpoints
- [ ] All mentioned related endpoints exist
  - [ ] `endpointName` - [file exists/missing]

## 7. Code Example Verification

### Basic Example
- [ ] GraphQL syntax is valid
- [ ] All fields in query/mutation exist
- [ ] Required fields are included
- [ ] Response structure matches actual response

### Advanced Example
- [ ] All optional parameters shown actually exist
- [ ] Nested objects structure is correct
- [ ] Complex types (JSON fields) have correct structure

## 8. Documentation Completeness

### Missing from Docs
- [ ] List any parameters found in code but not documented
- [ ] List any response fields found but not documented
- [ ] List any error cases found but not documented

### Extra in Docs (Hallucinated)
- [ ] List any parameters documented but not in code
- [ ] List any fields documented but not in code
- [ ] List any features documented but not implemented

## 9. Special Considerations

### Database/Prisma Verification
- [ ] If mentions database fields, verify against Prisma schema
  - Model: `[ModelName]`
  - Fields match: [Yes/No]

### Type Definitions
- [ ] All TypeScript/GraphQL types mentioned exist
  - [ ] `TypeName` - [location or "NOT FOUND"]

### Custom Field Types
- [ ] If mentions custom fields, verify types exist
  - [ ] All custom field types are real

## Summary

### Critical Issues (Must Fix)
1. [List any non-existent operations]
2. [List any hallucinated parameters]
3. [List any wrong types]

### Minor Issues (Should Fix)
1. [List any missing descriptions]
2. [List any formatting inconsistencies]

### Suggestions
1. [Any improvements for clarity]
2. [Missing helpful information]
```

## Verification Steps

### Step 1: GraphQL Schema Search
1. Navigate to `../bloo-api/` from the website root
2. Search for the GraphQL schema files:
   - Look in `src/graphql/schema/` or similar
   - Find type definitions (`.graphql` or `.gql` files)
   - Check resolver implementations in TypeScript

### Step 2: Input Verification
For each input parameter:
1. Find the input type definition in schema
2. Verify each field exists with correct type
3. Check if field is required (! notation)
4. Find where it's used in resolver code
5. Verify any validation rules

### Step 3: Enum Verification
1. Search for enum definitions in GraphQL schema
2. List all values and compare with documentation
3. Check if any values are deprecated

### Step 4: Permission Verification
1. Find the resolver implementation
2. Look for permission check code (likely using decorators or middleware)
3. Verify which roles are checked
4. Confirm the logic matches documentation

### Step 5: Response Verification
1. Find the return type in GraphQL schema
2. Check resolver to see what's actually returned
3. Verify all documented fields are real
4. Check for any computed fields

### Step 6: Error Verification
1. Search for error codes in the codebase
2. Find where errors are thrown
3. Verify error messages and codes match

## Search Strategies

### Finding GraphQL Definitions
```bash
# Search for mutation/query definitions
grep -r "mutation createProject" ../bloo-api/
grep -r "type Mutation" ../bloo-api/

# Search for input types
grep -r "input CreateProjectInput" ../bloo-api/

# Search for enums
grep -r "enum ProjectCategory" ../bloo-api/
```

### Finding Implementations
```bash
# Search for resolver functions
grep -r "createProject.*resolver" ../bloo-api/
grep -r "async createProject" ../bloo-api/

# Search for validation
grep -r "validate.*Project" ../bloo-api/
```

### Finding Permissions
```bash
# Search for permission decorators/checks
grep -r "@RequireRole" ../bloo-api/
grep -r "checkPermission" ../bloo-api/
grep -r "canCreate.*Project" ../bloo-api/
```

## Red Flags to Watch For

1. **Hallucinated Features**
   - Parameters that seem too convenient
   - Fields that don't appear in any search results
   - Enums with suspiciously perfect values

2. **Type Mismatches**
   - String vs ID types
   - Required vs optional discrepancies
   - Array types vs single values

3. **Missing Implementation**
   - Documented features with no code
   - Parameters accepted but never used
   - Response fields that aren't computed

4. **Permission Gaps**
   - Documented restrictions not enforced
   - Missing permission checks in code

## Tracking Progress

Update `api-tracker.md` with:
```markdown
# API Documentation Verification Tracker

## Status Legend
- 🔄 In Progress
- ✅ Verified 
- ❌ Has Issues
- 🔧 Fixed

## Files to Verify

### Start Guide
- [ ] 1.introduction.md
- [ ] 2.authentication.md
- [ ] 3.making-requests.md
- [ ] 4.GraphQL-playground.md
- [ ] 5.capabilities.md
- [ ] 7.rate-limits.md
- [ ] 8.upload-files.md

### Projects
- [ ] 1.index.md
- [ ] 2.create-project.md
- [ ] 2.delete-project.md
- [ ] 2.list-projects.md
[... continue for all files ...]

## Summary
- Total Files: [X]
- Verified: [X]
- Issues Found: [X]
- Fixed: [X]
```

## Final Notes

1. **Be Extremely Skeptical**: Assume everything might be hallucinated until proven otherwise
2. **Document Everything**: Even small discrepancies could indicate larger issues
3. **Check Related Code**: Don't just verify the direct resolver, check related services
4. **Test When Possible**: If you have access to GraphQL playground, test the actual queries
5. **Version Awareness**: Some features might be version-specific or feature-flagged

## Important: Complete Workflow

**DO NOT STOP** after creating the review file! The complete workflow is:
1. ✅ Create review file with all findings
2. ✅ **IMMEDIATELY edit the actual API documentation file** in `/content/en/api/` to fix all issues found
3. ✅ Update api-tracker.md with final status (🔧 Fixed if changes made, ✅ Verified if no changes needed)

Remember: The goal is to ensure the documentation is 100% accurate and helpful for developers using the Blue API.