--- 
title: Creating reusable checklists using automations
slug: reusable-checklists-automation
category: "Best Practices"
description: Learn how to create project management automations for reusable checklists.
image: /insights/check-background.png
date: 2024-07-08
showdate: true
---

In many projects and processes, you may need to use the same checklist across multiple records or tasks. 

However, it is not very efficient to manually retype the checklist each time you want to add it to a record. This is where you can leverage [powerful project management automations](/platform/features/automations) to automatically do this for you! 

As a reminder, automations in Blue require to key things:

1. A Trigger — What should happen to start the automation. This can be when a record is givena  specific tag, when it moves to a specific 
2. One or more Actions — In this case, it would be the automatic creation of one or more checklists.

Let's start with the action first, then discuss the possible triggers that you can use.

## Checklist Automation Action

You can create a new automation, and you can setup one or more checklists to be created, as per the example below:

![](/insights/checklist-automation.png)

These would be the checklist(s) that you want to be created each time you take the action.

## Checklist Automation Triggers

There are several ways you can trigger the creation of your reusable checklists. Here are some popular options:

- **Adding a Specific Tag:** You can set up the automation to trigger when a particular tag is added to a record. For example, when the tag "New Project" is added, it could automatically create your project initiation checklist.
- **Record Assignment:** The checklist creation can be triggered when a record is assigned to a specific individual or to anyone. This is useful for onboarding checklists or task-specific procedures.
- **Moving to a Specific List:** When a record is moved to a particular list in your project board, it can trigger the creation of a relevant checklist. For instance, moving an item to a "Quality Assurance" list could trigger a QA checklist.
- **Custom Checkbox Field:** Create a custom checkbox field and set the automation to trigger when this box is ticked. This gives you manual control over when to add the checklist.
- **Single or Multi-Select Custom Fields:** You can create a single or multi-select custom field with various options. Each option can be linked to a specific checklist template through separate automations. This allows for more granular control and the ability to have multiple checklist templates ready for different scenarios.

To enhance control over who can trigger these automations, you can hide these custom fields from certain users using custom user roles. This ensures that only project admins or other authorized personnel can trigger these options.

Remember, the key to effective use of reusable checklists with automations is to design your triggers thoughtfully. Consider your team's workflow, the types of projects you handle, and who should have the ability to initiate different processes. With well-planned automations, you can significantly streamline your project management and ensure consistency across your operations.

## Useful Resources

- [Project Management Automation Documentation](https://documentation.blue.cc/automations)
- [Custom User Roles Documentation](https://documentation.blue.cc/user-management/roles/custom-user-roles)
- [Custom Field Documentation](https://documentation.blue.cc/custom-fields)