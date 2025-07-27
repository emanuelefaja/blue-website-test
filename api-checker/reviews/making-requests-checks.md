# Verification for: making-requests.md
Path: /content/en/api/1.start-guide/3.making-requests.md
Status: [🔧] Fixed - Replaced hallucinated subscription with real one, added error examples

## 1. GraphQL Schema Verification

### Queries Verification
- [✓] `projectList` query exists
  - Location: /Users/manny/Blue/bloo-api/src/schema.graphql:301
  - Parameters match: `filter: ProjectListFilter!`
  - Returns `ProjectPagination!` with `items` array
  - Fields `name`, `id`, `updatedAt` all exist in Project type

### Mutations Verification
- [✓] `createTodo` mutation exists
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/createTodo.ts
  - Input type: `CreateTodoInput`
  - Parameters verified:
    - [✓] `todoListId` - exists (optional, not required as shown)
    - [✓] `title` - exists (required)
    - [✓] `position` - exists (optional)
  - Returns `Todo!` with fields `id`, `title`, `position`

- [✓] `deleteTodo` mutation exists
  - Location: /Users/manny/Blue/bloo-api/src/resolvers/Mutation/deleteTodo.ts
  - Input type: `DeleteTodoInput`
  - Parameter verified:
    - [✓] `todoId` - exists (required)
  - Returns `MutationResult!` with field `success`

### Subscription Verification
- [❌] `projectUpdated` subscription
  - Documented subscription does NOT exist
  - Location searched: entire schema and resolvers
  - Actual project subscriptions available:
    - `onAddedToProject`
    - `onRemovedFromProject`
    - `onArchiveProject`
    - `onUnarchiveProject`

## 2. Implementation Verification

### API Endpoint
- [✓] Production endpoint `https://api.blue.cc/graphql` is correct
  - Development: `http://localhost:4000/graphql`

### Headers
- [✓] All required headers shown correctly:
  - `Content-Type: application/json`
  - `X-Bloo-Token-ID`
  - `X-Bloo-Token-Secret`
  - `X-Bloo-Company-ID`

### WebSocket Support
- [✓] WebSocket endpoint exists
  - Would be `wss://api.blue.cc/graphql` in production
  - Uses `graphql-ws` library
  - Location: /Users/manny/Blue/bloo-api/src/lib/server.ts

## 3. Code Example Verification

### Reading Data Examples
- [✓] curl example syntax is valid
- [✓] Python example syntax is valid
- [✓] Node.js example syntax is valid
- [✓] Query structure is correct
- [✓] Response format matches actual schema

### Writing Data Examples
- [✓] createTodo mutation examples are valid
- [✓] deleteTodo mutation examples are valid
- [✓] All programming language examples have correct syntax

### Subscription Example
- [❌] WebSocket example shows non-existent subscription
- [❌] Message format may not match actual implementation
  - Shown: `{ type: 'start', payload: { query: '...' } }`
  - Need to verify against `graphql-ws` protocol

## 4. Response Format Verification

### projectList Response
- [✓] JSON structure is valid
- [✓] Fields match schema:
  - `data.projectList.items` array structure correct
  - Each item can have `name`, `id`, `updatedAt`

### Mutation Responses
- [✓] createTodo response structure would be valid
- [✓] deleteTodo response structure matches (`{ success: Boolean }`)

## 5. Documentation Completeness

### Missing Information
- [ ] No error response examples
- [ ] No mention of pagination parameters for projectList
- [ ] No mention of rate limiting
- [ ] No authentication error handling

### Extra/Hallucinated Information
- [❌] `projectUpdated` subscription does not exist
- [❌] WebSocket message format not verified

## Summary

### Critical Issues (Must Fix)
1. **Hallucinated subscription**: `projectUpdated` subscription does not exist in the API
2. **WebSocket example**: Shows a subscription that doesn't exist

### Minor Issues (Should Fix)
1. WebSocket message format needs verification against actual protocol
2. No error response examples provided
3. Missing information about available subscriptions

### Suggestions
1. Replace `projectUpdated` with an actual subscription like `subscribeToActivity`
2. Add error response examples
3. List available subscriptions
4. Verify WebSocket protocol message format with graphql-ws documentation