---
title: How does Blue offer unlimited file storage?
slug: unlimited-file-storage-faq
category: "FAQ"
description: How does Blue offer unlimited file storage for project management? 
image: /patterns/letter.png
date: 2024-08-01
showdate: false
---

This is perhaps one of the most frequently asked questions that we receive here at Blue, typically from brand new or prospective customers. 

At Blue, we've taken a bold step: offering [unlimited file storage](/platform/file-management) to all our customers.

It's a feature that often raises eyebrows and prompts the question:
 
> How is this possible?

We understand the skepticism. In an era where every gigabyte seems to come at a price, unlimited storage sounds too good to be true. You might wonder about the sustainability of such an offer or worry that it's a temporary perk that will disappear once costs become too high.

These concerns are valid, and we're here to address them head-on. 

In this article, we'll pull back the curtain and explain exactly how Blue manages to offer unlimited file storage without compromising our service or financial stability. We'll dive into the technology, economics, and strategies that make this possible, giving you a clear understanding of why you can trust Blue with your data storage needs, now *and* in the future.

## The Short Answer

The short answer is that we can do this for several reasons:

1.  **Limited actual usage**: Most customers don't actually upload that much data!
While we offer unlimited storage, the majority of our users only use a fraction of what's available to them. This means we can confidently provide this feature without worrying about excessive usage.
2. **Affordable cloud storage**: Data storage is actually quite cheap.
We have 33TB of storage (that's 33 *thousand* GB) and it costs just a couple of hundred dollars a month. The economies of scale in cloud storage mean we can offer vast amounts of storage at a relatively low cost.
3. **Smart storage strategies**: We're intelligent about how we store the data.
We use intelligent tiering with multiple storage classes that move less frequently accessed files to cheaper, but slightly slower file storage. This allows us to optimize costs while still providing quick access to frequently used files.


Let's get into the long answer!

## Most customers don't actually upload that much data.

Right now, we store 959,386 files. Most of these files are relatively small. Think PDF reports, screenshots, excel files, and so on. Some customers, of course, *do* upload large files such as videos for review, and this is fine. 

But on average, most files are just a few MB.

This graph shows our total storage growth over the years:


![](/resources/totalfilestorageyearly.png)

## Data Storage is actually quite cheap

This is often surprising to many people, and for good reason. Most of us are used to thinking about storage costs in terms of our personal devices - laptops, smartphones, and tablets. When we look at upgrading these devices, the cost for additional storage can be substantial. 

For instance:

- Upgrading from a 256GB to a 512GB laptop might cost an extra $200
- Jumping from a 128GB to a 256GB smartphone could set you back $100 or more
- Even external hard drives for personal use can seem pricey, with 1TB drives costing $50-$100

Given these personal experiences, it's natural to assume that cloud storage must be equally, if not more, expensive. After all, aren't cloud providers just using a lot of hard drives?

Right now, all of your files are stored with Amazon Web Services on a service called [Simple Storage Services (S3)Amazon Web Services (S3)](https://aws.amazon.com/s3/) 

To simplify, there are two main costs related to file storage when running an application like Blue:

1. The monthly cost of storing 1GB of data. 
2. The cost of transferring the data back to customers when they need to download a file, charged per GB.


Well, let's take a look at the numbers! 

| S3 Service | Cost per GB stored |
|:----------:|:------------------:|
| S3 Standard (First 50 TB/month) | $0.023 |
| S3 Standard (Next 450 TB/month) | $0.022 |
| S3 Standard (Over 500 TB/month) | $0.021 |
| S3 Standard - Infrequent Access | $0.0125 |
| S3 One Zone - Infrequent Access | $0.01 |
| S3 Glacier Instant Retrieval | $0.004 |
| S3 Glacier Flexible Retrieval | $0.004 |
| S3 Glacier Deep Archive | $0.00099 |


However, the reality of cloud storage is quite different:

- **Economy of Scale**: Cloud providers like AWS, Google Cloud, and Microsoft Azure operate at an enormous scale. They buy storage in bulk and can negotiate prices that are far lower than what individual consumers pay.
- **Optimized Infrastructure**: These providers have highly optimized data centers designed specifically for efficient storage and retrieval. This allows them to maximize the use of their hardware and minimize costs.
- **Advanced Storage Technologies**: Cloud providers use a mix of storage technologies, including traditional hard drives, SSDs, and even tape storage for rarely accessed data. This allows them to balance performance and cost effectively.
- **Tiered Storage Options**: As mentioned earlier, we use intelligent tiering. Cloud providers offer various storage classes with different access speeds and costs. Frequently accessed data might be stored on faster, slightly more expensive storage, while rarely accessed data can be moved to much cheaper, slower storage.

To put this into perspective, let's look at some numbers:

To illustrate just how much our storage costs are lower than those of commercial storage consumers, let's compare our storage costs to what you might pay when upgrading your smartphone storage.

Our current storage costs are approximately $0.006 to $0.01 per GB per month. 

Of course, there are other costs involved beyond just the raw storage - such as data transfer, redundancy for reliability, and the systems to manage all this data. But even accounting for these, the economics of cloud storage allow us to provide unlimited storage without putting our business at risk.

Then, we have to consider data transfer when users download data. AWS gives 100GB of data transfer out to the internet free each month, but after that we have to pay. 

Interestingly, transferring 1GB out (i.e. when a customer of Blue downloads a file) is **more expensive** than storing 1Gb for a month:

| Data Transfer OUT Volume | Cost per GB |
|:------------------------:|:-----------:|
| First 10 TB/month | $0.09 |
| Next 40 TB/month | $0.085 |
| Next 100 TB/month | $0.07 |
| Exceeding 150 TB/month | $0.05 |


### Considering Cloudflare R2


Most of our cost is from data transfers, *not* data hosting. 

![](/resources/datahostingvstransfer.png)

We have been customers of [Cloudflare](https://cloudflare.com) for quite a while, and our marketing website runs on [CloudFlare Pages](https://pages.cloudflard.com), and our DNS is managed via Cloudflare. 

Honestly, it is really a great service, it comes with lots of great features right out of the box.

They have an [S3-like file hosting plan called R2](https://www.cloudflare.com/developer-platform/r2/) that has one fantastic feature: 

Data transfers are **completely** free. 

And this is why we are considering moving across in the future, because that would completely eliminate one large chunk of our bill and allow us to continue to offer unlimited file storage forever. 

## Intelligent Tiering

One of the key strategies that allows us to offer unlimited storage while keeping costs manageable is our use of intelligent tiering. This approach takes advantage of how files are typically used in [project management](/solutions/project-management) scenarios, allowing us to optimize storage costs without sacrificing performance or accessibility.

Intelligent tiering is an automated storage management system that moves data between different storage classes based on access patterns. Here's how it works:

- When a file is first uploaded, it's stored in a high-performance, readily accessible storage tier.
- The system continuously monitors how often each file is accessed.
- If a file hasn't been accessed for a certain period, it's automatically moved to a less expensive, slightly slower storage tier.
- If the file is accessed again, it's promptly moved back to the high-performance tier.

This process happens seamlessly in the background, without any noticeable impact on the user experience.

### Blue's File Access Patterns

Our analysis of file usage across projects has revealed an interesting pattern:

- **Initial Intense Usage:** When a file is first uploaded to a project, it tends to be accessed frequently. This might be due to team members reviewing a new document, collaborating on a freshly uploaded design, or discussing a recent report.
- **Rapid Decline:** After this initial period of high activity, which typically lasts a few days, there's often a sharp drop-off in access frequency.
- **Long-Term Low Access:** Many files then enter a phase where they're accessed only occasionally, if at all. They remain important for reference or documentation purposes but aren't part of the day-to-day project work.

![](/resources/fileaccesspattern.png)

This pattern aligns perfectly with the concept of intelligent tiering. 

Here's how we leverage it:

- **High-Performance Storage for New Files**: When you upload a file, it's immediately available in our fastest storage tier. This ensures that during those crucial first few days of collaboration, your team experiences no lag in accessing and working with the file.
- **Cost-Effective Storage for Older Files**: As access to a file decreases, our system automatically moves it to a more cost-effective storage tier. This might be something like Amazon S3's Infrequent Access tier or even Glacier for very rarely accessed files.
- **Seamless Retrieval**: If an older file suddenly needs to be accessed again, our system quickly retrieves it and moves it back to the high-performance tier. There might be a slight delay the first time an old file is accessed, but subsequent accesses will be at full speed.


 
## Our Fair Use Policy: A Summary

We do have some [fair use policies](https://blue.cc/legal/terms) in place to ensure that all our users receive the best possible service:

At Blue, we're proud to offer unlimited file storage, but it's important to understand the intended use of this feature. Our unlimited storage is specifically designed for project-related files and collaborative work. It's not meant to be used as a general backup solution or a public file-sharing platform. 

Blue is a project management tool, intended for sharing files with clients and vendors within the context of your projects.

It's not designed for public distribution of files on forums, social media, or other public platforms.

To ensure fair usage and maintain a high-quality service for all our users, we have a few guidelines in place. We may review accounts that fall into the top 5% for data use or bandwidth, especially those with high public link sharing. 

For cost-efficiency, files not accessed for over 30 days are moved to slower storage. We trust our users to upload only files they have the right to share, and we have a zero-tolerance policy for illegal content. 

If usage patterns risk affecting our system's stability, we may need to take action, such as temporary account suspension or service limitations. 

However, we always aim to communicate with users before taking any such steps. Rest assured, the vast majority of our users never approach these limits in their day-to-day use of Blue.

## Conclusion

Offering unlimited file storage is not just a marketing gimmick for us at Blue—it's a sustainable feature backed by smart technology and economics. Let's recap why we're confident in our ability to maintain this offer:

- **Realistic usage patterns**: While we offer unlimited storage, most of our customers use only a fraction of what's available, allowing us to confidently provide this feature.
- **Cost-effective cloud storage**: Thanks to the economies of scale in cloud storage, we can provide vast amounts of storage at a surprisingly low cost.
- **Intelligent tiering**: Our smart storage strategies, including automated movement of files between storage tiers based on access patterns, help us optimize costs while maintaining performance.
- **Continuous innovation**: We're always exploring new technologies, like Cloudflare R2, that could further reduce our costs and improve our service.

By leveraging these factors, we're able to offer you unlimited storage *without* compromising our financial stability or the quality of our service. This means you can focus on your projects without worrying about storage limits or unexpected costs.

At Blue, we're committed to providing you with the tools you need to work effectively, and unlimited storage is a key part of that commitment. We believe in transparent communication about how we operate, which is why we've shared these details with you. As we continue to grow and evolve, our promise remains the same: we'll keep providing you with the storage you need, when you need it, without limitations. 

Thank you for trusting us with your projects and your data. 

If you have any more questions about our unlimited storage or any other aspect of our service, don't hesitate to reach out. 