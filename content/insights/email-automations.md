--- 
title: How to create custom email automations
slug: email-automations
category: "Product Updates"
description: Custom email notifications are an incredibly powerful feature in Blue that can help keep work moving forwards and ensure communication is on auto-pilot.
image: /insights/email-background.png
---

Email automations in Blue are a [powerful project management automation](/platform/features/automations) for streamlining communication, ensuring [great teamwork](/insights/great-teamwork) and keeping projects moving forward. By leveraging data stored within your records, you can automatically send personalized emails when certain triggers occur, such as a new record being created or a task becoming overdue. 

In this article, we’ll explore how to set up and use email automations in Blue.

## Setting Up Email Automations.

Creating an email automation in Blue is a straightforward process. First, select the trigger that will initiate the automated email. Some common triggers include:

- A new record is created
- A tag is added to a record
- A record is moved to another list


Next, configure the email details, including:

- From name and reply-to address
- To address (can be static or dynamically pulled from an email custom field)
- CC or BCC addresses (optional)

![](/insights/email-automations-image.webp)

One of the key benefits of email automations in Blue is the ability to personalize the content using merge tags. When customizing the email subject and body, you can insert merge tags that reference specific record data, such as the record name or custom field values. Simply use the {curly bracket} syntax to insert merge tags.

You can also include file attachments by dragging and dropping them into the email or using the attachment icon. Files from File custom fields may automatically attach if they are under 10MB.

Before finalizing your email automation, it’s recommended to send a test email to yourself or a colleague to ensure everything is functioning as intended.

## Use Cases and Examples

Email automations in Blue can be used for a variety of purposes. Here are a few examples:

1. Send a confirmation email when a client submits a request via an intake form. Set the trigger to send an email when a new record is created through the form, and make sure to include an email field in the form to capture the client’s address.
2. Notify an assignee when a new high-priority task is created. Set the trigger to send an email when a “Priority” tag is added to a record, and use the {Assignee} merge tag to automatically send the email to the assigned user.
3. Send a survey to a customer after a support ticket is marked as resolved. Set the trigger to send an email when a record is marked as completed and moved to the “Done” list. Include the customer’s email in a custom field and provide detailed information about the resolved issue in the email body.
4. Automate a recruitment program by sending confirmation emails to applicants. Set the trigger to send an email when an application is submitted through a form and added to the “Received” list. Capture the applicant’s email in the form, and use it to send a thank-you response.

## Benefits of Email Automations

Email automations in Blue offer several key benefits:

- Personalized communication through the use of merge tags and custom field data
- Automatic notifications that reduce manual work and ensure timely updates
- Structured, data-driven workflows that move projects forward based on record data

## Conclusion 

Email automations in Blue are a valuable tool for streamlining communication and keeping projects on track. By leveraging triggers, merge tags, and custom field data, you can create personalized, automated emails that enhance your team’s productivity and ensure important updates are never missed. With a wide range of use cases and easy setup, email automations are a must-have feature for any Blue user looking to optimize their workflow.

