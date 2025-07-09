---
title: FAQ on Blue Security
slug: faq-security
category: "FAQ"
description: This is a list of the most frequently asked questions on the security protocols and practices at Blue. 
image: /insights/security-background.png
date: 2024-07-19
showdate: false
---

Our mission is to organize the world's work by building the best project management platform on the planet.

Central to achieving this mission is ensuring that our platform is secure, reliable, and trustworthy. We understand that to be your single source of truth, Blue must safeguard your sensitive business data against outside threats, data loss, and downtime.

This means that we take security seriously at Blue. 

When we think about security, we consider a holistic approach that focusses on three key areas:

1.  **Infrastructure & Network Security**: Ensures that our physical and virtual systems are protected from external threats and unauthorized access.
2.  **Software Security**: Focuses on the security of the code itself, including secure coding practices, regular code reviews, and vulnerability management.
3.  **Platform Security**: Includes the features within Blue, such as [sophisticated access controls](/platform/features/user-permissions), ensuring that projects are private by default, and other measures to protect user data and privacy.


## How scalable is Blue?

This is an important question, as you want a system that can *grow* with you. You don't want to have to switch your project and process management platform in six or twelve months.

We choose platform providers with care, to ensure that they can handle the demanding workloads of our customers. We use cloud services from some of the world's top cloud providers that power companies such as [Spotify](https://spotify.com) and [Netflix](https://netflix.com), that have several orders of magnitude more traffic that we do. 

The main cloud providers we use are:

- **[Cloudflare](https://cloudflare.com)**: We manage of DNS (Domain Name Service) via Cloudflare as well as our marketing website which runs on [Cloudflare Pages](https://pages.cloudflare.com/). 
- **[Amazon Web Services](https://aws.amazon.com/)**: We use AWS for our database, which is [Aurora](https://aws.amazon.com/rds/aurora/), for file storage via [Simple Storage Service (S3)](https://aws.amazon.com/s3/), and also for sending emails via [Simple Email Service (SES)](https://aws.amazon.com/ses/)
- **[Render](https://render.com)**: We use Render for our front-end servers, application/API servers, our background services, queuing sytem, and Redis database. Interestingly, Render is actually built *on top* of AWS! 


## How secure are files in Blue? 

Let's start with data storage. Our files are hosted on [AWS S3](https://aws.amazon.com/s3/), which is the world's most popular cloud object storage with industry-leading scalability, data availability, security, and performance.

We have 99.99% file availability and 99.999999999% high durability. 

Let's break down what this means.

Availability refers to the amount of time that the data is operational and accessible. The 99.99% file availability means that we can expect files to be unavailable for no more than approximately 8.76 hours per year.

Durability refers to the likelihood that data remains intact and uncorrupted over time. This level of durability means we can expect to lose no more than one file out of 10 billion files uploaded, thanks to extensive redundancy and data replication across multiple data centers.

We use [S3 Intelligent-Tiering](https://aws.amazon.com/s3/storage-classes/intelligent-tiering/) to automatically move files to different storage classes based on the frequency of access. Based on the activity patterns of hundreds of thousands of projects, we notice that most files are access in a pattern that resembles an exponential backoff curve. This means that most files are accessed very frequently in the first few days, and then are quickly access less and less frequently. This allows us to move older files to slower, but significantly cheaper, storage without impacting the user experience in a meaningful way. This is how we are able to offer [unlimited file storage for all accounts.](/platform/features/file-management)

The cost savings for this are significant. S3 Standard-Infrequent Access (S3 Standard-IA) is approximately 1.84 times cheaper than S3 Standard. This means that for every dollar we would have spent on S3 Standard, we only spend about 54 cents on S3 Standard-IA for the same amount of data stored.

| Feature                  | S3 Standard             | S3 Standard-IA       |
|--------------------------|-------------------------|-----------------------|
| Storage Cost             | $0.023 - $0.021 per GB  | $0.0125 per GB        |
| Request Cost (PUT, etc.) | $0.005 per 1,000 requests | $0.01 per 1,000 requests |
| Request Cost (GET)       | $0.0004 per 1,000 requests | $0.001 per 1,000 requests |
| Data Retrieval Cost      | $0.00 per GB            | $0.01 per GB          |


The files you upload through Blue are encrypted both in transit and at rest. Data transferred to and from Amazon S3 is secured using [Transport Layer Security (TLS)](https://www.internetsociety.org/deploy360/tls/basics), protecting against [eavesdropping](https://en.wikipedia.org/wiki/Network_eavesdropping) and [man-in-the-middle attacks](https://en.wikipedia.org/wiki/Man-in-the-middle_attack). For at-rest encryption, Amazon S3 uses Server-Side Encryption (SSE-S3), which automatically encrypts all new uploads with AES-256 encryption, with Amazon managing the encryption keys. This ensures your data remains secure throughout its entire lifecycle.

## What about non-file data? 

Our database is powered by [AWS Aurora](https://aws.amazon.com/rds/aurora/), a modern relational database service that ensures high performance, availability, and security for your data.

Data in Aurora is encrypted both in transit and at rest. We use SSL (AES-256) to secure connections between your database instance and your application, protecting data during transfer. For at-rest encryption, Aurora uses keys managed through AWS Key Management Service (KMS), ensuring that all stored data, including automated backups, snapshots, and replicas, is encrypted and protected.

Aurora features a distributed, fault-tolerant, and self-healing storage system. This system is decoupled from compute resources and can auto-scale up to 128 TiB per database instance. Data is replicated across three [Availability Zones](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/) (AZs), providing resilience against data loss and ensuring high availability. In case of a database crash, Aurora reduces recovery times to less than 60 seconds, ensuring minimal disruption.

Blue continuously backs up our database to Amazon S3, enabling point-in-time recovery. This means we can restore the blue master database to any specific time within the last five minutes, ensuring that your data is always recoverable. We also take regular snapshots of the database for longer backup retention periods. 

As a fully managed service, Aurora automates time-consuming administration tasks such as hardware provisioning, database setup, patching, and backups. This reduces the operational overhead and ensures that our database is always up-to-date with the latest security patches and performance improvements. 

If we are more efficient, we can pass our cost savings to our customers with our [industry-leading pricing](/pricing). 

Aurora is compliant with various industry standards such as HIPAA, GDPR, and SOC 2, ensuring that your data management practices meet stringent regulatory requirements. Regular security audits and integration with [Amazon GuardDuty](https://aws.amazon.com/guardduty/) help detect and mitigate potential security threats.

## How does Blue ensure login security?

Blue uses [magic links via email](https://documentation.blue.cc/user-management/magic-links) to provide secure and convenient access to your account, eliminating the need for traditional passwords.

This approach significantly enhances security by mitigating common threats associated with password-based logins. By eliminating passwords, magic links protect against phishing attacks and password theft, *as there is no password to steal or exploit.* 

Each magic link is valid for only one login session, reducing the risk of unauthorized access. Additionally, these links expire after 15 minutes, ensuring that any unused links cannot be exploited, further enhancing security.

The convenience offered by magic links is also noteworthy. Magic links provide a hassle-free login experience, allowing you to access your account *without* the need to remember complex passwords.

This not only simplifies the login process but also prevents security breaches that occur when passwords are reused across multiple services. Many users tend to use the same password across various platforms, which means a security breach on one service could compromise their accounts on other services, including Blue. By using magic links, Blue's security is not dependent on the security practices of other services, providing a more robust and independent layer of protection for our users.

When you request to log in to your Blue account, a unique login URL is sent to your email. Clicking on this link will instantly log you into your account. The link is designed to expire after a single use or after 15 minutes, whichever comes first, adding an extra layer of security. By using magic links, Blue ensures that your login process is both secure and user-friendly, providing peace of mind and convenience.

## How can I check the reliability and uptime of Blue?

At Blue, we are committed to maintaining a high level of reliability and transparency for our users. To provide visibility into our platform's performance, we offer a [dedicated system status page](https://status.blue.cc) which is also linked from our footer on every page of our website. 

![](/insights/status-page.png)

This page displays our historical uptime data, allowing you to see how consistently our services have been available over time. Additionally, the status page includes detailed incident reports, providing transparency about any past issues, their impact, and the steps we have taken to resolve them and prevent future occurrences.



