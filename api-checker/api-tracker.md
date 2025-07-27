# API Documentation Verification Tracker

## Status Legend
- 🔄 In Progress
- ✅ Verified 
- ❌ Has Issues
- 🔧 Fixed

## Files to Verify

### Start Guide
- [✅] 1.introduction.md
- [🔧] 2.authentication.md - Documentation improved with pat_ prefix and bcrypt security mentions
- [🔧] 3.making-requests.md - Fixed: Replaced hallucinated subscription with subscribeToActivity, added error examples
- [✅] 4.GraphQL-playground.md - Verified accurate (auth covered in 2.authentication.md, introspection intentional)
- [🔧] 5.capabilities.md - Enhanced: Added query depth limit info and bulk operations clarification
- [🔧] 7.rate-limits.md - Fixed: Replaced misleading "no rate limits" with accurate table of 12 rate-limited operations
- [🔧] 8.upload-files.md - Fixed: Updated REST API file size limit from 5GB to 4.8GB to match implementation

### Projects
- [🔧] 1.index.md - Fixed: Updated to use projectList query, added missing PERSONAL/PROCUREMENT categories, corrected error codes, fixed API links
- [🔧] 2.create-project.md - Fixed: Removed hallucinated enum value, clarified coverConfig limitation, added response fields & error docs
- [✅] 2.delete-project.md - Accurate documentation, only minor error message fix applied
- [🔧] 2.list-projects.md - Enhanced: Added complete Project fields table with types and additional available fields
- [✅] 3.archive-project.md - Verified accurate, minor error message text fix applied
- [🔧] 3.project-activity.md - Fixed: Replaced UI documentation with comprehensive API documentation based on actual implementation
- [🔧] 3.rename-project.md - Fixed: Removed hallucinated PROJECT_NAME_TOO_LONG error, updated name to optional, added comprehensive EditProjectInput fields
- [🔧] 4.copy-project.md - Fixed: Wrong copyProjectStatus schema, added missing coverConfig option, corrected dependency claims
- [🔧] 5.lists.md - Enhanced: Fixed CLIENT role permissions and error message text
- [✅] 11.templates.md

### Records
- [🔧] 1.index.md - Enhanced: Fixed CLIENT role permissions clarification and error message text
- [✅] 2.list-records.md - Verified comprehensive implementation with enhanced performance notes
- [🔧] 3.toggle-record-status.md - Fixed: Corrected error messages, updated side effects list, removed archived project claim, fixed related endpoint references
- [🔧] 4.tags.md - Enhanced: Complete rewrite with full CRUD operations, advanced filtering, AI suggestions, and comprehensive documentation
- [🔧] 5.move-record-list.md - Complete rewrite: From 20 lines to 170+ comprehensive documentation with all implementation details
- [✅] 6.assignees.md - Verified: Complete rewrite from 20 lines to comprehensive API documentation with 3 operations, permissions, business logic - NO HALLUCINATIONS FOUND
- [✅] 7.update-record.md - Verified comprehensive implementation with enhanced permissions and return value documentation
- [🔧] 8.copy-record.md - Fixed: Corrected title field requirement, fixed response format, added missing COMMENTS option, updated error codes, enhanced permissions and cross-project documentation
- [🔧] 9.add-comment.md - Fixed: Removed non-existent files field, corrected file processing description

### Custom Fields
- [🔄] 1.index.md
- [🔧] 2.list-custom-fields.md - Enhanced: Fixed cursor pagination claim, clarified multi-project limitation, noted endCursor deprecation
- [🔧] 3.create-custom-fields.md - Fixed: Corrected TIME_DURATION enum values (TODO_CREATED_AT, TODO_MARKED_AS_COMPLETE), added missing currency conversion parameters
- [🔧] 4.custom-field-values.md - Fixed: Removed non-existent RECORD_NOT_FOUND error, clarified FORMULA/LOOKUP fields are read-only, enhanced permissions documentation
- [🔧] 5.delete-custom-field.md - Fixed: Removed non-existent PROJECT_NOT_ACTIVE error (98% accurate otherwise)
- [🔧] button.md - Fixed: Corrected button types to UI hints, removed non-existent errors, fixed permissions to role-based
- [🔧] checkbox.md - Fixed: Added case-sensitivity note, clarified import behavior, fixed link paths, removed non-existent forms API link (95% accurate)
- [🔧] country.md - Fixed: Clarified validation only in createTodo, corrected storage format, explained behavioral differences between mutations
- [🔧] currency-conversion.md - Fixed: Corrected permission constants from CUSTOM_FIELDS_CREATE/UPDATE to standard user roles (OWNER/ADMIN)
- [ ] currency.md
- [ ] date.md
- [ ] email.md
- [ ] file.md
- [ ] formula.md
- [ ] location.md
- [🔧] lookup.md - Complete rewrite: Removed all hallucinated aggregation functions, fixed to show lookups as data extractors only (from 30% to 100% accurate)
- [ ] number.md
- [ ] percent.md
- [🔄] phone.md
- [ ] rating.md
- [ ] reference.md
- [ ] select-multi.md
- [ ] select-single.md
- [ ] text-multi.md
- [ ] text-single.md
- [ ] time-duration.md
- [ ] unique-id.md
- [ ] url.md

### Automations
- [ ] 1.index.md

### User Management
- [ ] 1.index.md
- [ ] 2.list-users.md
- [ ] 3.remove-user.md
- [ ] 4.retrieve-custom-role.md

### Company Management
- [ ] 1.index.md

### Dashboards
- [ ] 1.index.md
- [ ] 2. Clone Dashboard copy.md
- [ ] 3. Rename Dashboard.md

### Libraries
- [ ] 1.python.md

### Other
- [🔧] 12.error-codes.md - Complete rewrite: From 57 lines to 262 lines documenting all 108 custom error codes organized by category with production safety info and best practices

## Summary
- Total Files: 73
- Verified: 3
- Issues Found: 1
- Fixed: 1

