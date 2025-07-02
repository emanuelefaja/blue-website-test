--- 
title: Project management automation — emails to stakeholders.
slug: project-management-automations-checkbox-email
category: "Product Updates"
description: Often, you want to be in control of your project management automations 
image: /resources/email-background.png
date: 2024-07-08
showdate: true
---

We’ve covered how to [create email automations before.](/resources/email-automations)

However, often there are stakeholders in projects that only need to be alerted when there is something *really* important. 

Wouldn't it be nice if there was a project management automation where you, as a project manager, could be in control of *exactly* when to notify a key stakeholder via the press of a button?

Well, turns out that with Blue, you can do precisely this! 

Today we are going to learn how to create a really useful project management automation: 

A checkbox that automatically notifies one or more key stakeholders, giving them all the key context of what you’re notifying them about. As a bonus point, we will also learn how to lock down this ability so only certain members of your project can trigger this email notification.

This will look something like this once you’re done:

![](/public/resources/checkbox-email-automation.png)

And just by checking this checkbox, you will be able to trigger a project management automation to send a custom notification email to stakeholders. 

Let’s go step by step.

## 1. Create your checkbox custom field

This is very easy, you can check out our [detailed documentation](https://documentation.blue.cc/custom-fields/introduction#creating-custom-fields) on creating custom fields.

Make sure that you name this field something obvious that you’ll remember such as “notify management” or “notify stakeholders”. 

## 2. Create your project management automation trigger.

On the records view in your project, click on the small robot on the top right to open the automation settings:

<video autoplay loop muted playsinline>
  <source src="/public/videos/notify-stakeholders-automation-setup.mp4" type="video/mp4">
</video>

## 3. Create your project management automation action.

In this case, our action will be to send a custom email notification to one or more email addresses. It's good to note here that these people do **not** have to be in Blue to receive these emails, you can send emails to *any* email address.  

You can learn more in our [detailed documentation guide on how to setup email automations](https://documentation.blue.cc/automations/actions/email-automations)

Your final result should look something like this:

![](/public/resources/email-automation-example.png)

## 4. Bonus: Restrict access to the checkbox.

You can use [custom user roles in Blue](/platform/user-permissions) to restrict access to the checkbox custom fields, ensuring that only authorized team members can trigger email notifications.

Blue allows Project Administrators to define roles and assign permissions to user groups. This system is crucial for maintaining control over who can interact with specific elements of your project, including custom fields like the notification checkbox.

1. Navigate to the User Management section in Blue and select "Custom User Roles."
2. Create a new role by providing a descriptive name and an optional description.
3. Within the role permissions, locate the section for Custom Fields Access.
4. Specify whether the role can view or edit the checkbox custom field. For example, restrict editing access to roles like "Project Administrator" while allowing a newly created custom role to manage this field.
5. Assign the newly created role to the appropriate users or user groups. This ensures that only the designated individuals have the capability to interact with the notification checkbox.

[Read more at our official documentation site.](https://documentation.blue.cc/user-management/roles/custom-user-roles)

By implementing these custom roles, you enhance the security and integrity of your project management processes. Only authorized team members can trigger critical email notifications, ensuring that stakeholders receive important updates without unnecessary alerts. 

## Conclusion

By implementing the project management automation outlined above, you gain precise control over when and how to notify key stakeholders. This approach ensures that important updates are communicated effectively, without overwhelming your stakeholders with unnecessary information. Utilizing Blue’s custom field and automation features, you can streamline your project management process, enhance communication, and maintain a high level of efficiency.

With just a simple checkbox, you can trigger custom email notifications tailored to your project's needs, ensuring that the right people are informed at the right time. Moreover, the ability to restrict this functionality to specific team members adds an extra layer of control and security.

Start leveraging this powerful feature in Blue today to keep your stakeholders informed and your projects running smoothly. For more detailed steps and additional customization options, refer to the provided documentation links. Happy automating!