# Review: Button Custom Field Documentation
Path: /content/en/api/5.custom fields/button.md
Date: 2025-01-27
Status: ❌ Has Significant Issues

## Summary
The button custom field documentation contains several inaccuracies and describes functionality that doesn't exist in the current implementation. The API is simpler than documented.

## Major Issues Found

### 1. ❌ Button Types Not Enforced
**Documentation Claims:**
- Specific buttonType values: `noConfirmation`, `softConfirmation`, `hardConfirmation`
- These control confirmation behavior

**Reality:**
- These specific values don't exist in the code
- `buttonType` accepts any string value
- No validation or enforcement of button types
- Confirmation is purely a UI feature, not API-enforced

### 2. ❌ Non-Existent Error Codes
**Documentation Lists:**
- `NO_AUTOMATIONS` error
- `INVALID_CONFIRMATION` error

**Reality:**
- Neither error code exists in the codebase
- No special error handling for buttons

### 3. ❌ Incorrect Permission Names
**Documentation Claims:**
- `CUSTOM_FIELDS_CREATE` permission
- `CUSTOM_FIELDS_UPDATE` permission
- `AUTOMATIONS_CREATE` permission

**Reality:**
- These permission constants don't exist
- Actual permissions use role-based system (OWNER, ADMIN, MEMBER, CLIENT)

### 4. ⚠️ Misleading Confirmation Behavior
**Documentation Claims:**
- API handles confirmation differently based on buttonType
- Mentions confirmation validation

**Reality:**
- API doesn't enforce any confirmation
- All button clicks execute immediately via API
- Confirmation is 100% client-side UI behavior

## What IS Correct

### ✅ Accurate Information:
1. BUTTON is a valid CustomFieldType
2. Buttons trigger `CUSTOM_FIELD_BUTTON_CLICKED` automation events
3. Buttons don't store data values
4. Uses `setTodoCustomField` mutation to click
5. Basic field creation structure is correct
6. Automation integration concept is accurate

## Recommendations

### 1. Simplify Button Types Section
Remove the specific enum values and explain that `buttonType` is a free-form string for UI hints only.

### 2. Remove Non-Existent Errors
Delete the entire error responses section or replace with actual errors that can occur.

### 3. Fix Permissions
Replace specific permission constants with role-based permissions:
- Creating fields: OWNER/ADMIN roles
- Clicking buttons: Based on custom field edit permissions

### 4. Clarify API vs UI Behavior
Explicitly state that ALL confirmation logic is client-side and the API always executes immediately.

### 5. Update Examples
Remove examples that suggest different API behavior for different button types.

## Suggested Rewrite Priority
**High** - This documentation is significantly misleading and could cause developer confusion.

## Overall Accuracy Score
**40%** - While the core concept is correct, most specific details are wrong or misleading.