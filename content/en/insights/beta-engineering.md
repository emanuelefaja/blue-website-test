---
title:  Why Blue has an Open Beta
category: "Engineering"
description: Learn why our project management system has an ongoing open beta. 
date: 2024-08-03
---

Many B2B SaaS startups launch in Beta, and for good reasons. It’s part of the traditional Silicon Valley motto of *“move fast and break things”*.

Putting a “beta” sticker on a product reduces expectations. 

Something is broken? Oh well, it is just a beta.

The system is slow? Oh well, it is just a beta. 

[The documentation](https://blue.cc/docs) is non-existent? Oh well…you get the point.

And this is *actually* a good thing.  Reid Hoffman, the founder of LinkedIN, famously said:

> If you are not embarrassed by the first version of your product, you've launched too late.

And the beta sticker is also good for customers. It helps them to self-select.

The customers who try beta products are the ones at the at the early stages of the Technology Adoption Lifecycle, also known as the Product Adoption Curve.

The Technology Adoption Lifecycle is typically divided into five main segments:

1. Innovators
2. Early Adopters
3. Early Majority
4. Late Majority
5. Laggards

![](/insights/technology-adoption-lifecycle-graph.png)


However, eventually the product has to mature, and customers expect a stable, working product. They don’t want access to a “beta” environment where things break.

Or do they?

*This* is the question we asked ourselves.

We believe we asked ourselves this question due to the nature of how Blue was initially built. [Blue started as an off-shoot of a busy design agency](/insights/agency-success-playbook), and so we worked *inside* the office of a business that was actively using Blue to run all their projects.

This means that for years, we were able to observe how *real* human beings — sitting right next to us! — used Blue in their daily lives. 

And because they used Blue from the early days, this team always used Blue Beta! 

And so it was natural for us to allow all our other customers to use it as well. 

**And this is why we do not have a dedicated testing team.**

That’s right. 

Nobody at Blue has the *sole* responsibility for ensuring that our platform is running well and stable. 

This is for a several reasons.

The first is a lower cost base. 

Not having a full-time testing team significantly reduces our costs, and we are able to pass these savings to our customers with the lowest prices in the industry. 

To put this into perspective, we offer enterprise-level feature sets that our competition charges $30-$55/user/month for just $7/month. 

This doesn’t happen by accident, it’s *by design*.

However, it’s not a good strategy to sell a cheaper product if it doesn’t work.

So the *real question is*, how do we manage to create a stable platform that thousands of customers can use without a dedicated testing team?

Of course, our approach to having an open Beta is crucial to this, but before we dive into this, we want to touch on developer responsibility.

We made the early decision at Blue that we would never split responsibilities for front-end and back-end technologies. We would only hire or train full stack developers. 

The reason we made this decision was to ensure that a developer would fully own the feature that they were working on. So there would be none of the *“throw the problem over the garden fence”* mentality that you sometimes get when there are joint responsibilities for features. 

And this extends to the testing of the feature, of understanding  the customer use cases and request, and of reading and commenting on the specifications. 

In other words, each developer builds a deep and intuitive understanding of the feature that they are building. 

Okay, let’s now talk about our open beta.

When we say it is “open” — we mean it. Any customer can try it simply by adding “beta” in front of our web application url.

So “app.blue.cc” becomes “beta.app.blue.cc” 

When they do this, they are able to see their usual data, as both the Beta and Production environments share the same database, but they will also be able to see new features. 

Customers can easily work even if they have some team members on Production and some curious ones on Beta. 

We typically have a few hundred customers using Beta at any one the, and we post feature previews on our community forums so they can checkout what’s new and try it out. 

And this is the point: we have *several hundred* testers! 

All of these customers will try out features in their wor kflows, and be quite vocal if something is not quite right, because they  are *already* implementing the feature inside of their business! 

The most common feedback are small but very useful changes that address edge cases that we did not consider. 

We leave new features on Beta between 2-4 week. Whenever we feel that they are stable, then we release to production.

We also have the ability to by-pass Beta if required, using a fast-track flag. This is typically done for bug-fixes that we do not want to hold for 2-4 weeks before shipping to production.

The result? 

Pushing to production feels…well boring! Like nothing — it’s simply not a big deal for us. 

And it means that this smooths out our release schedule, which is what has enabled us to [ship features monthly like clockwork for the last six years.](/changelog).

However, like any choice, there are some tradeoffs.

Customer supports is slightly more complex, as we have to support customers across two versions of our platform. Sometimes this can cause confusion to customers that have team members using two different versions.

Another pain point is that this approach can sometimes slow down the overall release schedule to production. This is especially true for larger features that can get '"stuck” in Beta if there is another related feature that is having problems and need some further work.

But on balance, we think that these tradeoffs are worth the benefits of a lower cost base and greater customer engagement.

We're one of the few software companies to embrace this approach, but it is now a fundamental part of our product development approach.

