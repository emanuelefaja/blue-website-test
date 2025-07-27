# Verification for: 2.list-projects.md
Path: /content/en/api/2.projects/2.list-projects.md
Status: [✓] In Progress / [ ] Completed

## 1. GraphQL Schema Verification

### Query Name
- [ ] Verify the GraphQL operation name exists in schema
  - Operation: `projectList`
  - Location in schema: [file:line]
  - Actual vs Documented: [any differences]

### Input Type Verification
- [ ] Verify input type name is correct
  - Documented: `ProjectListFilter`
  - Actual in schema: [actual name or "NOT FOUND"]
  - Location: [file:line]

### Filter Parameters
For each parameter in ProjectListFilter:
- [ ] `companyIds`
  - Documented type: `[String!]!`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

- [ ] `ids`
  - Documented type: `[String!]`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

- [ ] `archived`
  - Documented type: `Boolean`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

- [ ] `isTemplate`
  - Documented type: `Boolean`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

- [ ] `search`
  - Documented type: `String`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

- [ ] `folderId`
  - Documented type: `String`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

- [ ] `inProject`
  - Documented type: `Boolean`
  - Actual type: [actual or "NOT FOUND"]
  - Required status matches: [Yes/No]

### Response Fields Verification
Check all documented project fields exist in Project type:
- [ ] `id` - Type: [actual type]
- [ ] `uid` - Type: [actual type]
- [ ] `slug` - Type: [actual type]
- [ ] `name` - Type: [actual type]
- [ ] `description` - Type: [actual type]
- [ ] `archived` - Type: [actual type]
- [ ] `color` - Type: [actual type]
- [ ] `icon` - Type: [actual type]
- [ ] `createdAt` - Type: [actual type]
- [ ] `updatedAt` - Type: [actual type]
- [ ] `allowNotification` - Type: [actual type]
- [ ] `position` - Type: [actual type]
- [ ] `unseenActivityCount` - Type: [actual type]
- [ ] `todoListsMaxPosition` - Type: [actual type]
- [ ] `lastAccessedAt` - Type: [actual type]
- [ ] `isTemplate` - Type: [actual type]
- [ ] `automationsCount` - Type: [actual type]
- [ ] `totalFileCount` - Type: [actual type]
- [ ] `totalFileSize` - Type: [actual type]
- [ ] `todoAlias` - Type: [actual type]

### PageInfo Fields Verification
- [ ] `totalPages` - Type: [actual type]
- [ ] `totalItems` - Type: [actual type]  
- [ ] `page` - Type: [actual type]
- [ ] `perPage` - Type: [actual type]
- [ ] `hasNextPage` - Type: [actual type]
- [ ] `hasPreviousPage` - Type: [actual type]

## 2. Enum Verification

### ProjectSort Values
Check all documented sort values exist:
- [ ] `id_ASC` - [exists/missing]
- [ ] `id_DESC` - [exists/missing]
- [ ] `name_ASC` - [exists/missing]
- [ ] `name_DESC` - [exists/missing]
- [ ] `createdAt_ASC` - [exists/missing]
- [ ] `createdAt_DESC` - [exists/missing]
- [ ] `updatedAt_ASC` - [exists/missing]
- [ ] `updatedAt_DESC` - [exists/missing]
- [ ] `position_ASC` - [exists/missing]
- [ ] `position_DESC` - [exists/missing]

## 3. Implementation Verification

### Resolver Check
- [ ] Resolver exists for this operation
  - Location: [file:line]
  - Handler function: `[functionName]`

### Business Logic Verification
- [ ] All documented filter parameters are used
- [ ] Default filtering behavior for `inProject: false`
- [ ] Folder filtering restrictions with `inProject: false`
- [ ] Position sorting restrictions with non-member projects

### Validation Rules
- [ ] Required companyIds field enforced
- [ ] Pagination limits (skip/take) enforced
- [ ] Folder + inProject:false restriction enforced

## 4. Permission Verification

### Access Requirements
- [ ] Permission checks exist in resolver
- [ ] User membership in company verified
- [ ] `inProject: false` requires company owner permission

## 5. Error Response Verification

### Expected Errors
- [ ] Company not found error
- [ ] Permission denied for `inProject: false`
- [ ] Invalid folder ID error

## 6. Advanced Example Verification

### Basic Example
- [ ] GraphQL syntax is valid
- [ ] All fields exist and have correct types
- [ ] Required parameters included

### Advanced Example  
- [ ] All filter parameters exist
- [ ] Sort values are valid
- [ ] `totalCount` field exists (shown in advanced example)
- [ ] `take` parameter works (documented as `take`, not `limit`)

## 7. Documentation Claims Verification

### Business Logic Claims
- [ ] Verify "case-insensitive" search claim
- [ ] Verify position sorting restriction with non-member projects  
- [ ] Verify folder filtering restrictions
- [ ] Verify archived/template exclusion for non-member projects
- [ ] Verify deprecated parameters are actually deprecated

### Default Values
- [ ] Default skip: 0
- [ ] Default take: 20
- [ ] Default behavior for undefined filters

## 8. Consistency Checks

### Parameter Naming
- [ ] Verify `take` vs `limit` consistency across docs
- [ ] Check if `totalCount` field exists (used in advanced example)
- [ ] Verify pagination field names match actual schema

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