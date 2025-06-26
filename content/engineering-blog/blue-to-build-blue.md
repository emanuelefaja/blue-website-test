---
title:  How we use Blue to build Blue. 
slug: blue-to-build-blue
tags: ["engineering-blog"]
description: Learn how we use our own project management platform to build our project management platform! 
image: /patterns/bluewhitecircle.png
date: 07/08/2024
showdate: true
sitemap:
  loc: /resources/blue-to-build-blue
  lastmod: 2024/08/07
  changefreq: monthly
---

You're about to get an insider's tour of how Blue builds Blue.

At Blue, we eat our own dog food. 

This means that we use Blue to *build* Blue. 

This strange-sounding term, often referred to as "dogfooding" is often attributed to Paul Maritz, a manager at Microsoft in the 1980s. He reportedly sent an email with the subject line *"Eating our own dog food"* to encourage Microsoft employees to use the company's products.

The idea of using your own tools to build your tools is that it leads to positive feedback cycle.

The idea of using your own tools to build your tools leads to a positive feedback cycle, creating numerous benefits:

- **It helps us identify real-world usability issues quickly.** As we use Blue daily, we encounter the same challenges our users might face, allowing us to address them proactively.
- **It accelerates bug discovery.** Internal use often reveals bugs before they reach our customers, improving overall product quality.
- **It enhances our empathy for end-users.** Our team gains firsthand experience of Blue's strengths and weaknesses, helping us make more user-centric decisions.
- **It drives a culture of quality within our organization.** When everyone uses the product, there's a shared stake in its excellence.
- **It fosters innovation.** Regular use often sparks ideas for new features or improvements, keeping Blue at the cutting edge.


[We've spoken before about why we have no dedicated testing team](/resources/open-beta) and this is yet another reason. 

If there are bugs in our system, we almost always find them in our constant daily use of the platform. And this also creates a forcing function to fix them, as we will obviously find them very annoying as we are probably one of the top users of Blue! 

This approach demonstrates our commitment to the product. By relying on Blue ourselves, we show our customers that we truly believe in what we're building. It's not just a product we sell – it's a tool we depend on every day. 

## Main Process

We have one project in Blue, aptly named "Product". 

**Everything** related to our product development is tracked here. Customer feedback, bugs, feature ideas, ongoing work, and so on. The idea of having one project where we track everything is that it [promotes better teamwork.](/resources/great-teamwork) 

Each records is a feature or part of a features.   This is how  move from "wouldn't it be cool if..." to "check out this awesome new feature!" 

The project has the following lists:

- **Ideas/Feedback**: This is a list of team ideas or customer feedback based on calls or email exchanges. Feel free to add any ideas here! In this list, we have not decided that we will yet build any of these features, but we regularly review this for ideas that we want to explore further
- **Backlog (Long Term)**: This is where features from the Ideas/Feedback list go if we decide they would be a good addition to Blue.
- **{Current Quarter}**: This is typically structured as "Qx YYYY" and shows our quarter priorities. 
- **Bugs**: This is a list of known bugs reported by the team or customers. Bugs added here will automatically have the "Bug" tag added.
- **Specifications**: These features are currently being specified. Not every feature requires a specification or design; it depends on the expected size of the feature and the confidence level that we have with regards to edge cases and complexity. 
- **Design Backlog**: This is the backlog for the designers, any time they have finished something that is in progress they can pick any item from this list. 
- **In Progress Design**: This is the current features that the designers are designing.
- **Design Review**: This is where the features whose. designs are currently being reviewed. 
- **Backlog (Short Term)**: This is a list of features we will likely start working on in the next few weeks. This is where assignments take place. The CEO and Head of Engineering decide which features are assigned to which engineer based on previous experience and workload.  [Team members can then pull these into the In Progress](/resources/push-vs-pull-kanban) once they have completed their current work. 
- **In Progress**: These are features that are currently being developed.
- **Code Review**: Once a feature has finished development, it undergoes a code review. Then it will either be moved back to "In Progress" if adjustments are needed or deployed to the Development environment.
- **Dev**: These are all the features currently in the Development environment.  Other team members and certain customers can review these. 
- **Beta**: These are all the features currently in the [Beta environment](https://beta.app.blue.cc). Many customers use this as their daily Blue platform and will also provide feedback. 
- **Production**: When a feature reaches production, it is then considered done.

Sometimes, as we develop a feature, we realize that certain sub-features are more difficult to implement than initially expected, and we may choose not to do them in the initial version that we deploy to customers. In this case, we can spin up a new record with a name following the format "{FeatureName} V2" and include all the sub-features as checklist items.

## Tags

- **Mobile**: This means that the feature is specific to either our iOS, Android, or iPad apps.
- **{EnterpriseCustomerName}**: A feature is specifically being built for an enterprise customer. Tracking is important as there are typically additional commercial agreements for each feature.
- **Bug**: This means that this is a bug that requires fixing.
- **Fast-Track**: This means that this is a Fast-Track Change that does not have to go through the full release cycle as described above.
- **Main**: This is a major feature development. It is typically reserved for major infrastructure work, big dependency upgrades, and significant new modules within Blue.
- **AI**: This feature contains an artificial intelligence component.
- **Security**: This means a security implication must be reviewed or a patch is required.


The fast-track tag is particularly interesting. This is reserved for smaller, less complex updates that don't require our full release cycle, and that we want to ship to customers within 24-48 hours. 

Fast-track changes are typically minor adjustments that can significantly improve user experience without altering core functionality. Think fixing typos in the UI, tweaking button padding, or adding new icons for better visual guidance. These are the kind of changes that, while small, can make a big difference in how users perceive and interact with our product. They are also annoying if they take ages to ship! 

Our fast-track process is straightforward.

It starts with creating a new branch from main, implementing the changes, and then creating merge requests for each target branch - Dev, Beta, and Production. We generate a preview link for review, ensuring that even these small changes meet our quality standards. Once approved, the changes are merged simultaneously into all branches, keeping our environments in sync.

## Custom Fields

We don't have many custom fields in our Product project. 

- **Specifications**: This links to a Blue doc that has the specification for that particular feature. This is not always done, as it depends on the complexity of the feature.
- **MR**: This is the link to the Merge Request in [Gitlab](https://gitlab.com) where we host our code. 
- **Preview Link**: For features that primarily change the front-end, we can create a unique URL that has those changes for each commit, so we can easily review the changes. 
- **Lead**: This field tells us which senior engineer is taking point on the code review. It ensures that every feature gets the expert attention it deserves, and there's always a clear go-to person for questions or concerns.

## Checklists

During our weekly demos, we will drop the discussed feedback in a checklist called "feedback" and there will also be another checklist that contains the main [WBS (Work Breakdown Scope)](/resources/simple-work-breakdown-structure) of the feature, so we can easily tell what's done and what's yet to do. 

## Conclusion

And that's it!

We think that sometimes people are surprised and how straight-forward our process is, but we believe that simple processes are often far superior than overly complex processes that you cannot easily understand. 

This simplicity is intentional. It allows us to stay agile, respond quickly to customer needs, and keep our entire team aligned. 

By using Blue to build Blue, we're not just developing a product – we're living it.

So the next time you're using Blue, remember: you're not just using a product we've built. You're using a product we personally rely on every single day. 

And that makes all the difference.