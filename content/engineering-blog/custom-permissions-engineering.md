---
title:  Creating Blue's Custom Permissions Engine
slug: engineering-custom-permissions-project-management
tags: ["engineering-blog"]
description: Go behind the scenes with the Blue engineering team as they explain how to built an AI-powered auto-categorization and tagging feature.
image: /resources/custom-permissons-background.png
date: 2024-07-25
showdate: true
sitemap:
  loc: /resources/engineering-custom-permissions-project-management
  lastmod: 2024-07-25
  changefreq: monthly
---

Effective project and process management is crucial for organizations of all sizes. 

At Blue, [we've made it our mission](/about) to organize the world's work by building the best project management platform on the planet—simple, powerful, flexible, and affordable for all.

This means that our platform must adapt to the unique needs of each team. Today, we're excited to pull back the curtain on one of our most powerful features: Custom Permissions.

Project management tools are the backbone of modern workflows, housing sensitive data, crucial communications, and strategic plans. As such, the ability to finely control access to this information isn't just a luxury—it's a necessity. 

<video autoplay loop muted playsinline>
  <source src="/product/user-roles.mp4" type="video/mp4">
</video>


Custom permissions play a critical role in B2B SaaS platforms, especially in project management tools, where the balance between collaboration and security can make or break a project's success.

But here's where Blue takes a different approach: **we believe that enterprise-grade features shouldn't be reserved for enterprise-sized budgets.** 

In an era where AI empowers small teams to operate at unprecedented scales, why should robust security and customization be out of reach?

In this behind-the-scenes look, we'll explore how we developed our Custom Permissions feature, challenging the status quo of SaaS pricing tiers, and bringing powerful, flexible security options to businesses of all sizes. 

Whether you're a startup with big dreams or an established player looking to optimize your processes, custom permissions can enable new use cases you never even knew were possible. 

## Understanding Custom User Permissions

Before we dive into our journey of developing custom permissions for Blue, let's take a moment to understand what custom user permissions are and why they're so crucial in project management software.

Custom user permissions refer to the ability to tailor access rights for individual users or groups within a software system. Instead of relying on predefined roles with fixed sets of permissions, custom permissions allow administrators to create highly specific access profiles that align perfectly with their organization's structure and workflow needs.

In the context of project management software like Blue, custom permissions include:

* **Granular access control**: Determining who can view, edit, or delete specific types of project data.
* **Feature-based restrictions**: Enabling or disabling certain features for particular users or teams.
* **Data sensitivity levels**: Setting varying levels of access to sensitive information within projects.
* **Workflow-specific permissions**: Aligning user capabilities with specific stages or aspects of your project workflow.

The importance of custom permissions in project management cannot be overstated:

* **Enhanced security**: By providing users with only the access they need, you reduce the risk of data breaches or unauthorized changes.
* **Improved compliance**: Custom permissions help organizations meet industry-specific regulatory requirements by controlling data access.
* **Streamlined collaboration**: Teams can work more efficiently when each member has the right level of access to perform their role without unnecessary restrictions or overwhelming privileges.
* **Flexibility for complex organizations**: As companies grow and evolve, custom permissions allow the software to adapt to changing organizational structures and processes.

## Getting to YES

[We've written before](/resources/value-proposition-blue), that every feature in Blue has to be a **hard** YES before we decide to build it. We don't have the luxury of hundreds of engineers and to waste time and money building things that customers don't need.

And so, the path to implementing custom permissions in Blue wasn't a straight line. Like many powerful features, it began with a clear need from our users and evolved through careful consideration and planning.

For years, our customers had been requesting more granular control over user permissions. As organizations of all sizes began handling increasingly complex and sensitive projects, the limitations of our standard role-based access control became apparent. 

Small startups working with external clients, mid-size companies with intricate approval processes, and large enterprises with strict compliance requirements all voiced the same need: 

More flexibility in managing user access.

Despite the clear demand, we initially hesitated to dive into developing custom permissions. 

Why? 

We understood the complexity involved! 

Custom permissions touch every part of a project management system, from the user interface down to the database structure. We knew that implementing this feature would require significant changes to our core architecture and careful consideration of performance implications. 

As we surveyed the landscape, we noticed that very few of our competitors had attempted to implement a powerful custom permissions engine like the one our customers were requesting. Those who did often reserved it for their highest-tier enterprise plans. 

It became clear why: the development effort is substantial, and the stakes are high.

Implementing custom permissions incorrectly could introduce critical bugs or security vulnerabilities, potentially compromising the entire system. This realization underscored the magnitude of the challenge we were considering.

### Challenging the Status Quo

However, as we continued to grow and evolve, we reached a pivotal realization: 

**The traditional SaaS model of reserving powerful features for enterprise customers no longer made sense in today's business landscape.**

In 2024, with the power of AI and advanced tools, small teams can operate at a scale and complexity that rivals much larger organizations. A startup might be handling sensitive client data across multiple countries. A small marketing agency could be juggling dozens of client projects with varying confidentiality requirements. These businesses need the same level of security and customization as *any* large enterprise.

We asked ourselves: Why should the size of a company's workforce or budget determine their ability to keep their data safe and their processes efficient?

### Enterprise-Grade for All

This realization led us to a core philosophy that now drives much of our development at Blue: Enterprise-grade features should be accessible to businesses of all sizes.

We believe that:

- **Security shouldn't be a luxury.** Every company, regardless of size, deserves the tools to protect their data and processes.
- **Flexibility drives innovation.** By giving all our users powerful tools, we enable them to create workflows and systems that push their industries forward.
- **Growth shouldn't require platform changes.** As our customers grow, their tools should seamlessly grow with them.

With this mindset, we decided to tackle the challenge of custom permissions head-on, committed to making it available to all our users, not just those on higher-tier plans.

This decision set us on a path of careful design, iterative development, and continuous user feedback that ultimately led to the custom permissions feature we're proud to offer today. 

In the next section, we'll dive into how we approached the design and development process to bring this complex feature to life.

### Design and Development 

When we decided to tackle custom permissions, we quickly realized we were facing a behemoth of a task. 

At first glance, "custom permissions" might sound straightforward, but it's a deceptively complex feature that touches every aspect of our system.

The challenge was daunting: we needed to implement cascading permissions, allow on-the-fly edits, make significant database schema changes, and ensure seamless functionality across our entire ecosystem – web, Mac, Windows, iOS, and Android apps, as well as our API and webhooks. 

The complexity was enough to make even the most seasoned developers pause.

Our approach centered on two key principles: 

1. Breaking down the feature into manageable versions 
2. Embracing incremental shipping. 

Faced with the complexity of full-scale custom permissions, we asked ourselves a crucial question: 

> What would be the simplest possible first version of this feature? 

This approach aligns with the agile principle of delivering a Minimum Viable Product (MVP) and iterating based on feedback.

Our answer was refreshingly straightforward:

1. Introduce a toggle to hide the project activity tab
2. Add another toggle to hide the forms tab

**That was it.**

No bells and whistles, no complex permission matrices—just two simple on/off switches. 

While it might seem underwhelming at first glance, this approach offered several significant advantages:

* **Quick Implementation**: These simple toggles could be developed and tested rapidly, allowing us to get a basic version of custom permissions into users' hands quickly.
* **Clear User Value**: Even with just these two options, we were providing tangible value. Some teams might want to hide the activity feed from clients, while others might need to restrict access to forms for certain user groups.
* **Foundation for Growth**: This simple start laid the groundwork for more complex permissions. It allowed us to set up the basic infrastructure for custom permissions without getting bogged down in complexity from the outset.
* **User Feedback**: By releasing this simple version, we could gather real-world feedback on how users interacted with custom permissions, informing our future development.
* **Technical Learning**: This initial implementation gave our development team practical experience in modifying permissions across our platform, preparing us for more complex iterations.

And you know, it is actually quite humbling to have a huge vision for something, and then to ship something that is such a small percentage of that vision. 

After shipping these first two toggles, we decided to tackle something more sophisticated. We landed on two new custom user role permissions.

The first, was the ability to limit users to only view records that have been specifically assigned to them. This is very useful if you have a client in a project and you only want them to see records that are specfically assigned to them instead of everything that you are working on for them. 

The second, was an option for project administrators to block user groups from being able to invite other users. This is good if you have a sensitive project that you want to ensure stays on a "need to see" basis. 

Once we had shipped this, we gained more confidence and for our third version we tackled column-level permissions, which means being able to decide which custom fields a specific user group can view or edit.

This is extremely powerful. Imagine that you have a CRM project, and you have data in there that is not only related to the amounts that the customer will pay, but also your cost and profit margins. You may not want your cost fields and project margin formula field to be visible to junior staff, and custom permission allows you to lock down those fields so they are not shown. 

Next up, we moved onto creating list-based permissions, where proejct administrators can decide if a user group can view, edit, and delete a specific list. If they hide a list, all the records inside that list also become hidden, which is great because it means that you can hide certain parts of your process from your team members or clients. 

This is the end result:

<video autoplay loop muted playsinline>
  <source src="/resources/custom-user-roles.mp4" type="video/mp4">
</video>

## Technical Considerations

At the heart of Blue's technical architecture lies GraphQL, a pivotal choice that has significantly influenced our ability to implement complex features like custom permissions. But before we dive into the specifics, let's take a step back and understand what GraphQL is and how it differs from the more traditional REST API approach.
GraphQL vs REST API: An Accessible Explanation

Imagine you're at a restaurant. With a REST API, it's like ordering from a fixed menu. You ask for a specific dish (endpoint), and you get everything that comes with it, whether you want it all or not. If you want to customize your meal, you might need to make multiple orders (API calls) or ask for a specially prepared dish (custom endpoint).

GraphQL, on the other hand, is like having a conversation with a chef who can prepare anything. You tell the chef exactly what ingredients you want (data fields), and in what quantities. The chef then prepares a dish that's precisely what you asked for - no more, no less. This is essentially what GraphQL does - it allows the client to ask for exactly the data it needs, and the server provides just that.

### An Important Lunch

About six weeks into Blue's initial development, our lead engineer and CEO went out for lunch. 

The topic of discussion? 

Whether to switch from REST APIs to GraphQL. This wasn't a decision to be taken lightly - adopting GraphQL would mean discarding six weeks of initial work.

On the walk back to the office, the CEO posed a crucial question to the lead engineer: "Would we regret not doing this five years from now?" 

The answer became clear: GraphQL was the way forward. 

We recognized the potential of this technology early on, seeing how it could support our vision for a flexible, powerful project management platform.

Our foresight in adopting GraphQL paid dividends when it came to implementing custom permissions. With a REST API, we would have needed a different endpoint for every possible configuration of custom permissions - an approach that would quickly become unwieldy and difficult to maintain.

GraphQL, however, allows us to handle custom permissions dynamically. Here's how it works:

- **On-the-fly Permission Checks**: When a client makes a request, our GraphQL server can check the user's permissions right from our database.
- **Precise Data Retrieval**: Based on these permissions, GraphQL returns only the requested data that fits within the user's access rights.
- **Flexible Queries**: As permissions change, we don't need to create new endpoints or alter existing ones. The same GraphQL query can adapt to different permission setups.
- **Efficient Data Fetching**: GraphQL allows clients to request exactly what they need. This means we're not overfetching data, which could potentially expose information the user shouldn't access.

This flexibility is crucial for a feature as complex as custom permissions. It allows us to offer granular control *without* sacrificing performance or maintainability.

## Challenges

Implementing custom permissions in Blue brought its share of challenges, each pushing us to innovate and refine our approach. Performance optimization quickly emerged as a critical concern. As we added more granular permission checks, we risked slowing down our system, especially for large projects with many users and complex permission setups. To address this, we implemented a multi-tiered caching strategy, optimized our database queries, and leveraged GraphQL's ability to request only necessary data. This approach allowed us to maintain swift response times even as projects scaled and permission complexity grew.

The user interface for custom permissions presented another significant hurdle. We needed to make the interface intuitive and manageable for administrators, even as we added more options and increased the system's complexity. 

Our solution involved multiple rounds of user testing and iterative design. 

We introduced a visual permissions matrix that allowed admins to quickly view and modify permissions across different roles and project areas. 

Ensuring cross-platform consistency presented its own set of challenges. We needed to implement custom permissions uniformly across our web, desktop, and mobile applications, each with its unique interface and user experience considerations. This was particularly tricky for our mobile apps, which had to dynamically hide and show features based on the user's permissions. We addressed this by centralizing our permission logic in the API layer, ensuring that all platforms received consistent permission data. 

Then, we developed a flexible UI framework that could adapt to these permission changes in real-time, providing a seamless experience regardless of the platform used.

User education and adoption presented the final hurdle in our custom permissions journey. Introducing such a powerful feature meant we needed to help our users understand and effectively leverage custom permissions.

We initially launched custom permissions to a subset of our user base, carefully monitoring their experiences and collecting insights. This approach allowed us to refine the feature and our educational materials based on real-world usage before launching to our entire user base. 

The phased rollout proved invaluable, helping us identify and address minor issues and user confusion points that we hadn't anticipated, ultimately leading to a more polished and user-friendly feature for all of our users. 

This apporach of launching to a sub-set of users, as well as our typically 2-3 week "Beta" period on our public Beta helps us sleep at night. :)

## Looking Ahead

As with all features, nothing is ever *"done"*.

Our long-term vision for the custom permissions feature stretches across tags, custom-field filters, customizable project navigation, and comment controls.

Let's dive into each aspects.

### Tag Permissions

We think it would be amazing to be able to create permissions based on whether a record has one or more tags. The most obvious use case would be that you create a custom user role called "Customers" and only allow users in that role to see records that have the tag "Customers".

This gives you an at-a-glance view of whether a record can or cannot be seen by your customers. 

This could become even more powerful with AND/OR combinators, where you can specify more complex rules. For example, you could set up a rule that allows access to records tagged both "Customers" AND "Public", or records tagged either "Internal" OR "Confidential". This level of flexibility would allow for incredibly nuanced permission settings, catering to even the most complex organizational structures and workflows.

The potential applications are vast. Project managers could easily segregate sensitive information, sales teams could have automatic access to relevant client data, and external collaborators could be seamlessly integrated into specific parts of a project without risking exposure to sensitive internal information.

### Custom Field Filters

Our vision for Custom Field Filters represents a significant leap forward in granular access control. This feature will empower project administrators to define what records specific user groups can see based on the values of custom fields. It's about creating dynamic, data-driven boundaries for information access.

Imagine being able to set up permissions like this:

- Show only records where the "Project Status" dropdown is set to "Public"
- Restrict visibility to items where the "Department" multi-select field includes "Marketing"
- Allow access to tasks where the "Priority" checkbox is ticked
- Display projects where the "Budget" number field is above a certain threshold


### Customizable Project Navigation

This is simply an extension of the toggles that we already have. Instead of just having toggles for "activity" and "forms", we want to extend that to every single part of the project navigation.  This way, proejct administrators can create focussed interfaces and remove tools that they don't need. 

### Comment Controls

In the future, we want to be creative in how we allow our customers to decide who can and cannot see comments. We may allow multiple tabbed comment areas under one record, and each can be visible or not visible to different user groups.

Additionally, we may also allow a feature where only comments where a user is *specifically* mention are visible, and nothing else is. Thsi would allow teams that have clients on projects to ensure that only comments that they want clients to see are visible. 

## Conclusion

So there we have it, this is how we approached building one of most interesting and powerful features! [As you can see on our project management comparison tool](/compare), very few project management systems have such a powerful permission-matrix setup, and the ones that do reserve it for their most expensive enterprise plans, making it inaccessible to a typical small or medium company.

With Blue, you have *all* features available with our plan — we don't believe that enteprise-level features should be reserved for enterprise customers! 