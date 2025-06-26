---
title:  Scaling CSV Imports & Export to 250,000+ Records
slug: scaling-csv-import-export
tags: ["engineering-blog"]
description: Discover how Blue scaled CSV imports and exports 10x using Rust and scalable architecture and strategic tech choices in B2B SaaS.
image: /resources/import-export-background.png
date: 2024-07-18
showdate: true
sitemap:
  loc: /resources/scaling-csv-import-export
  lastmod: 2024-07-18
  changefreq: monthly
---
At Blue, we're [constantly pushing the boundaries](/platform/roadmap) of what's possible in project management software. Over the years, we've [released hundreds of features](/platform/changelog)

Our latest engineering feat? 

A complete overhaul of our [CSV import](https://documentation.blue.cc/integrations/csv-import) and [export](https://documentation.blue.cc/integrations/csv-export) system, dramatically improving performance and scalability. 

This post takes you behind the scenes of how we tackled this challenge, the technologies we employed, and the impressive results we achieved.

The most interesting thing here is that we had to step outside of our typical [technology stack](https://sop.blue.cc/product/technology-stack) to achieve the results that we wanted. This is decision that has to be made thoughtfully, because the long-term repercussion can be severe in terms of technology debt and long-term maintenance overhead. 

<video autoplay loop muted playsinline>
  <source src="/resources/import-export-video.mp4" type="video/mp4">
</video>

## Scaling for Enterprise Needs

Our journey began with a request from an enterprise customer in the events industry. This client uses Blue as their central hub for managing vast lists of events, venues, and speakers, integrating it seamlessly with their website. 

For them, Blue isn't just a tool — it's the single source of truth for their entire operation.

While we are always proud to hear that customers use us for such mission-critical needs, there is also a large amount of responsibility on our side to ensure a fast, reliable system.

As this customer scaled their operations, they faced a significant hurdle: **importing and exporting large CSV files contaning 100,000 to 200,000+ records.**

This was beyond the capabilty of our system at the time. In fact, our previous import/export system was already struggling with imports and exports containing more than 10,000 to 20,000 records! So 200,000+ records was out of the question. 

 Users experienced frustratingly long wait times, and in some cases, imports or exports would *fail to complete altogether.* This significantly affected their operations as they relied on daily imports and exports to manage certain aspects of their operations. 

> Multi-tenancy is an architecture where a single instance of software serves multiple customers (tenants). While efficient, it requires careful resource management to ensure one tenant's actions don't negatively impact others.

And this limitation wasn't just affecting this particular client. 

Due to our multi-tenant architecture—where multiple customers share the same infrastructure—a single resource-intensive import or export could potentially slow down operations for other users, which in practice often happened. 

As usual, we did a build vs buy analysis, to understand if we should spend the time to upgrade our own system or buy a system from someone esle. We looked at various possibilities.

The vendor that did stand out was a SaaS provider called [Flatfile](https://flatfile.com/). Their system and capabilities looked like exactly what we needed. 

But, after reviewing their [pricing](https://flatfile.com/pricing/), we decided that this would end up being an extremely expensive solution for application of our scale — *$2/file starts to add up really quickly!* —and it was better to extend our built-in CSV import/export engine. 

To tackle this challenge, we made a bold decision: introduce Rust into our primarly Javascript tech stack. This systems programming language, known for its performance and safety, was the perfect tool for our performance-critical CSV parsing and data mapping needs.

Here's how we approached the solution.

### Introducing Background Services

The foundation of our solution was the introduction of background services to handle resource-intensive tasks. This approach allowed us to offload heavy processing from our main server, significantly improving overall system performance.
Our background services architecture is designed with scalability in mind. Like all components of our infrastructure, these services auto-scale based on demand. 

This means that during peak times, when multiple large imports or exports are being processed simultaneously, the system automatically allocates more resources to handle the increased load. Conversely, during quieter periods, it scales down to optimize resource usage.

This scalable background service architecture has benefitted Blue not only for CSV imports & exports. Over time, we've moved a substantial number of features into background services to take load off our main servers:

- **[Formula Calculations](https://documentation.blue.cc/custom-fields/formula)**: Offloads complex mathematical operations to ensure rapid updates of derived fields without impacting main server performance.
- **[Dashboard/Charts](/platform/dashboards)**: Processes large datasets in the background to generate up-to-date visualizations without slowing down the user interface.
- **[Search Index](https://documentation.blue.cc/projects/search)**: Continuously updates the search index in the background, ensuring fast and accurate search results without impacting system performance.
- **[Copying Projects](https://documentation.blue.cc/projects/copying-projects)**: Handles the replication of large, complex projects in the background, allowing users to continue working while the copy is being created.
- **[Project Management Automations](/platform/project-management-automation)**: Executes user-defined automated workflows in the background, ensuring timely actions without blocking other operations.
- **[Repeating Records](https://documentation.blue.cc/records/repeat)**: Generates recurring tasks or events in the background, maintaining schedule accuracy without burdening the main application.
- **[Time Duration Custom Fields](https://documentation.blue.cc/custom-fields/duration)**: Continuously calculates and updates the time difference between two events in Blue, providing real-time duration data without impacting system responsiveness.

## New Rust Module for Data Parsing

The heart of our CSV processing solution is a custom Rust module. While this marked our first venture outside our core tech stack of Javascript, the decision to use Rust was driven by its exceptional performance in concurrent operations and file processing tasks.

Rust's strong suits align perfectly with the demands of CSV parsing and data mapping. Its zero-cost abstractions allow for high-level programming without sacrificing performance, while its ownership model ensures memory safety without the need for garbage collection. These features make Rust particularly adept at handling large datasets efficiently and safely.

For CSV parsing, we leveraged Rust's csv crate, which offers high-performance reading and writing of CSV data. We combined this with custom data mapping logic to ensure seamless integration with Blue's data structures.

The learning curve for Rust was steep but manageable. Our team dedicated about two weeks to intensive learning for this.

The improvements were impressive:

![](/resources/import-export.png)


Our new system can process the same amount of records that our old system could process in 15 minutes in around 30 seconds. 

## Web Server and Database Interaction

For the web server component of our Rust implementation, we chose Rocket as our framework. Rocket stood out for its combination of performance and developer-friendly features. Its static typing and compile-time checking align well with Rust's safety principles, helping us catch potential issues early in the development process.
On the database front, we opted for SQLx. This async SQL library for Rust offers several advantages that made it ideal for our needs:

- Type-safe SQL: SQLx allows us to write raw SQL with compile-time checked queries, ensuring type safety without sacrificing performance.
- Async support: This aligns well with Rocket and our need for efficient, non-blocking database operations.
- Database agnostic: While we primarily use [AWS Aurora](https://aws.amazon.com/rds/aurora/), which is MySQL compatible, SQLx's support for multiple databases gives us flexibility for the future in case we ever decide to change. 

## Optimization of Batching

Our journey to the optimal batching configuration was one of rigorous testing and careful analysis. We ran extensive benchmarks with various combinations of concurrent transactions and chunk sizes, measuring not just raw speed but also resource utilization and system stability.

The process involved creating test datasets of varying sizes and complexity, simulating real-world usage patterns. We then ran these datasets through our system, adjusting the number of concurrent transactions and the chunk size for each run.

After analyzing the results, we found that processing 5 concurrent transactions with a chunk size of 500 records provided the best balance of speed and resource utilization. This configuration allows us to maintain high throughput without overwhelming our database or consuming excessive memory.

Interestingly, we found that increasing concurrency beyond 5 transactions didn't yield significant performance gains and sometimes led to increased database contention. Similarly, larger chunk sizes improved raw speed but at the cost of higher memory usage and longer response times for small to medium-sized imports/exports.

## CSV Exports via Email Links

The final piece of our solution addresses the challenge of delivering large exported files to users. Instead of providing a direct download from our web app, which could lead to timeout issues and increased server load, we implemented a system of emailed download links.

When a user initiates a large export, our system processes the request in the background. Once complete, rather than holding the connection open or storing the file on our web servers, we upload the file to a secure, temporary storage location. We then generate a unique, secure download link and email it to the user.

These download links are valid for 2 hours, striking a balance between user convenience and information security. This timeframe gives users ample opportunity to retrieve their data while ensuring that sensitive information isn't left accessible indefinitely.

The security of these download links was a top priority in our design. Each link is:

- Unique and randomly generated, making it practically impossible to guess
- Valid for only 2 hours
- Encrypted in transit, ensuring the safety of data as it's downloaded

This approach offers several benefits:

- It reduces the load on our web servers, as they don't need to handle large file downloads directly
- It improves the user experience, especially for users with slower internet connections who might face browser timeout issues with direct downloads
- It provides a more reliable solution for very large exports that might exceed typical web timeout limits

User feedback on this feature has been overwhelmingly positive, with many appreciating the flexibility it offers in managing large data exports.

## Exporting Filtered Data

The other obvious improvement was to allow users to only export data that was already filtered in their project view. This means if there is active tag "priority", then only records that have this tag would end up in the CSV export. This means less time manipulating data in Excel to filter things out that are not important, and also helps us reduce the number of rows to process.

## Looking Ahead

While we don't have immediate plans to expand our use of Rust, this project has shown us the potential of this technology for performance-critical operations. It's an exciting option we now have in our toolkit for future optimization needs. This CSV import and export overhaul aligns perfectly with Blue's commitment to scalability. 

We're dedicated to providing a platform that grows with our customers, handling their expanding data needs without compromising on performance.

The decision to introduce Rust into our technology stack wasn't taken lightly. It raised an important question that many engineering teams face: When is it appropriate to venture outside your core tech stack, and when should you stick with familiar tools?

There's no one-size-fits-all answer, but at Blue, we've developed a framework for making these crucial decisions:

- **Problem-First Approach:** We always start by clearly defining the problem we're trying to solve. In this case, we needed to dramatically improve the performance of CSV imports and exports for large datasets.
- **Exhausting Existing Solutions:** Before looking outside our core stack, we thoroughly explore what can be achieved with our existing technologies. This often involves profiling, optimization, and rethinking our approach within familiar constraints.
- **Quantifying the Potential Gain:** If we're considering a new technology, we need to be able to clearly articulate and, ideally, quantify the benefits. For our CSV project, we projected order-of-magnitude improvements in processing speed.
- **Assessing the Costs:** Introducing a new technology isn't just about the immediate project. We consider the long-term costs:
  - Learning curve for the team
  - Ongoing maintenance and support
  - Potential complications in deployment and operations
  - Impact on hiring and team composition
- **Containment and Integration:** If we do introduce a new technology, we aim to contain it to a specific, well-defined part of our system. We also ensure we have a clear plan for how it will integrate with our existing stack.
- **Future-Proofing:** We consider whether this technology choice opens up future opportunities or if it might paint us into a corner.

One of the primary risks of frequently adopting new technologies is ending up with what we call a *"technology zoo"* - a fragmented ecosystem where different parts of your application are written in different languages or frameworks, requiring a wide range of specialized skills to maintain.


## Conclusion

This project exemplifies Blue's approach to engineering: *we're not afraid to step outside our comfort zone and adopt new technologies when it means delivering a significantly better experience for our users.* 

By reimagining our CSV import and export process, we've not only solved a pressing need for one enterprise client but improved the experience for all our users dealing with large datasets.

As we continue to push the boundaries of what's possible in [project management software](/solutions/project-management), we're excited to tackle more challenges like this. 

Stay tuned for more [deep dives into the engineering that powers Blue!](/resources/engineering-blog)
