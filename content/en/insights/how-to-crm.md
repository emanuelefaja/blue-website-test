---
title: How to setup Blue as a CRM
category: "Best Practices"
description: Learn how to setup Blue to track your customers and deals in an easy way. 
date: 2024-08-11
---

## Introduction

One of the key advantages of using Blue is not using it for a *specific* use case, but using *across* use cases. This means that you do not have to pay for multiple tools, and you also have one place where you can easily switch between your various projects and processes such as hiring, sales, marketing, and more. 

In helping thousands of customers to get setup in Blue over the years, we’ve noticed that the tough part if *not* setting up Blue itself, but thinking through the processes and making the most out of out platform.

The key parts are thinking of the step-by-step workflow for each of your business processes that you want to track, and also the specifics of the data that you want to capture, and how this translates to the custom fields that you setup. 

Today, we are going to walk you through creating [an easy to use, yet powerful sales CRM system](/solutions/use-case/sales-crm) with a customer database that is linked to a pipeline of opportunities. All this data will flow into a dashboard where you can see realtime data on your total sales, forecasted sales, and more. 

## Customer Database

The first thing to do is to setup a new project to store your customer data. This data is then going to be cross-referenced in another project where you track specific sales opportunities. 

The reason we split your customer information from the opportunities, is that they do not map one-to-one.

One customer may have multiple opportunities or projects. 

For instance, if you are a marketing and design agency, you may initially engage with a customer for their branding, and then do a separate project for their website, and then another for their social media management.

All of these would be separate sales opportunities that require their own tracking and proposals, but they are all linked to that one customer. 

The advantage of splitting out your customer database into a separate project is that if you update any details in your customer database, all your opportunities will automatically have the new data, which means that you now have one source of truth in your business! You don’t have to go back and edit everything manually! 

So, the first thing to decide is if you are going to be company-centric or person-centric.

This decision really depends on what you are selling, and to who you sell. If you sell primarily to businesses, then you will likely want the record name to be the company name. However, if you sell mostly to individuals (i.e. you’re a personal health coach, or a personal branding consultant), then you would most likely take a person-centric approach.

So the record name field is going either be the company name or the person name, depending on your choice. The reason for this, is that it means you can easily identify a customer at a glance, just be looking at your board or database. 

Next, you need to consider what information you want to capture as part of your customer database. These are going to become your custom fields.

The usual suspects here are:

- Email
- Phone Number
- Website
- Address
- Source (i.e. where did this customer first come from?)
- Category 

In Blue, you can also remove any default fields that you do not need. For this customer database, we typically recommend you remove due date, assignee, dependencies, and checklists. You may want to keep our default description field available in case you have general notes about that customer that are not specific to any sale opportunity. 

We recommended that you keep the "Reference by" field, as this will be useful later on. Once we setup our opportunity database, we will be able to see every sales record that is linked to this particular customer here. 

In terms of lists, we typically see our customers just keep it simple and have one list called "Customers" and leave it at that. It is better to use tags or custom fields for categorization.

What is great here is that once you have this setup, you can easily import your data from other systems or Excel sheets into Blue via our CSV import function, and you can also create a form for new potential customers to submit their details so you can **automatically** capture them into your database. 

## Opportunities Database

Now that we have our customer database, we need to create another project to capture our actual sales opportunities. You can call this project "Sales CRM" or "Opportunities".

### Lists as Process Steps

To setup your sales process, you need to think about what are the usual steps that an opportunity goes through from the moment you receive a request from a customer all the way to getting a signed contract. 

Each list in your project will be a step in your process.

Regardless of your specific process, there will be a few common lists that ALL Sales CRMs should have:

- Unqualified — All incoming requests, where you have not yet qualified a customer. 
- Closed Won — All of the opportunities that you won and turned into sales! 
- Closed Lost — All of the opportunities where you quoted a customer, and they did not accept. 
- N/A — This is where you place all the opportunities that you did not win, but also were not "lost". It could be the ones that you turned down, the ones where the customer, for whatever reason, ghosted you, and so on. 

In terms of thinking through your sales crm business process, you should consider the level of granularity that you want. We do not recommend having 20 or 30 columns, this typically gets confusing and stops you from being able to see the bigger picture. 

However, it is also important not to make each process to broad, as otherwise deals will get "stuck" at a specific stage for weeks or months, even when they are in fact moving forwards. Here is a typical recommended approach:

- **Unqualified**: All incoming requests, where you have not yet qualified a customer.
- **Qualification**: This is where you take the opportunity and start the process of understanding if this is a good fit for your firm.
- **Writing Proposal**: This is where you start to turn the opportunity into a pitch for your firm. This is a document that you would send to the client.
- **Proposal Sent**: This is where you've sent the proposal to the client and are waiting for a response.
- **Negotiations**: This is where you are in the process of finalizing the deal.
- **Contract Out for Signature**: This is where you are just waiting for the client to sign the contract.
- **Closed Won**: This is where you've won the deal and are now working on the project.
- **Closed Lost**: This is where you've quoted the client, but they have not accepted the terms.
- **N/A**: This is where you place all the opportunities that you did not win, but also were not "lost". It could be the ones that you turned down, the ones where the customer, for whatever reason, ghosted you, and so on.

### Tags as Service Categories 
Let's now talk about tags. 

We recommended that you use tags for the different types of services that you offer. So, going back to our marketing and design agency example, you have have tags for "branding","website","SEO","Facebook Management", and so on.

The advantages here is that you can easily filter down by service in one click, which can give you a brief overview of which services are more popular, and this can also inform future hiring, as typically different services require different team members. 

### Sales CRM Custom Fields

Net up, we need to consider what custom fields we want to have. 

Typical ones that we see used are:

- **Amount**: This is a currency field for the amount of the project
- **Cost**: Your expected cost to fulfill the sale, also a currency field
- **Profit**: A formula field to calculate the profit based on the amount and cost fields. 
- **Proposal URL**: This can include a link to an online Google Doc or Word document of your proposal, so you can easily click and review it. 
- **Received files**: This can be a file custom field where you can drop any files received from the client such as research materials, NDAs, and so on.
- **Contracts**: Another file custom field where you can add signed contracts for safekeeping. 
- **Confidence Level**: A star custom field with 5 stars, indicating how confident you are of winning this particular opportunity. This can be used later on in the dashboard for forecasting! 
- **Expected Close Date**: A date field to estimate when the deal is likely to close.
- **Customer**: A reference field linking to the primary contact person in the customer database.
- **Customer Name**: A lookup field that pull the customer name from the particular linked record in the customer database.
- **Customer Email**: A lookup field that pull the customer email from the particular linked record in the customer database.
- **Deal Source**: A dropdown field to track where the opportunity originated (e.g., referral, website, cold call, trade show).
- **Reason Lost**: A dropdown field (for closed lost deals) to categorize why the opportunity was lost.
- **Customer Size**: A dropdown field to categorize customers by size (e.g., small, medium, large enterprise).

Again, it is really **up to you** to decide precisely which fields you want to have. One word of warning: it is easy when setting up to add lots and lots of fields to your Sales CRM of data you would like to capture. However, you must be realistic in terms of the discipline and time commitment. There's no point having 30 fields in your Sales CRM if 90% of records will not have any data in them. 

The great thing about custom fields is that they integrate well into [Custom Permissions](/platform/features/user-permissions). This means you can decide exactly which fields team members in your team can view or edit. For instance, you may want to hide cost and profit information from junior staff. 

### Automations 

[Sales CRM Automations](/platform/features/automations) are a powerful feature in Blue that can streamline your sales process, ensure consistency, and save time on repetitive tasks. By setting up intelligent automations, you can enhance your sales CRM's effectiveness and allow your team to focus on what matters most - closing deals. Here are some key automations to consider for your sales CRM:

- **New Lead Assignment**: Automatically assign new leads to sales representatives based on predefined criteria such as location, deal size, or industry. This ensures quick follow-up and balanced workload distribution.
- **Follow-up Reminders**: Set up automated reminders for sales reps to follow up with prospects after a certain period of inactivity. This helps prevent leads from falling through the cracks.
- **Stage Progression Notifications**: Notify relevant team members when a deal moves to a new stage in the pipeline. This keeps everyone informed of progress and allows for timely interventions if needed.
- **Deal Aging Alerts**: Create alerts for deals that have been in a particular stage for longer than expected. This helps identify stalled deals that may need extra attention.


## Linking Customers and Deals

One of the most powerful features of Blue for creating an effective CRM system is the ability to link your customer database with your sales opportunities. This connection allows you to maintain a single source of truth for customer information while tracking multiple deals associated with each customer. Let's explore how to set this up using Reference and Lookup custom fields.

### Setting Up the Reference Field


1. In your Opportunities (or Sales CRM) project, create a new custom field.
2. Choose the "Reference" field type.
3. Select your Customer Database project as the source for the reference.
4. Configure the field to allow single select (as each opportunity is typically associated with one customer).
5. Name this field something like "Customer" or "Associated Company".

Now, when creating or editing an opportunity, you'll be able to select the associated customer from a dropdown menu populated with records from your Customer Database.

### Enhancing with Lookup Fields

Once you've established the reference connection, you can use Lookup fields to bring relevant customer information directly into your opportunities view. Here's how:

1. In your Opportunities project, create a new custom field.
2. Choose the "Lookup" field type.
3. Select the Reference field you just created ("Customer" or "Associated Company") as the source.
4. Choose which customer information you want to display. You might consider fields like: Email, Phone Number, Customer Category, Account Manager

Repeat this process for each piece of customer information you want to display in your opportunities view.

The benefits of this are:

- **Single Source of Truth**: Update customer information once in the Customer Database, and it automatically reflects in all linked opportunities.
- **Efficiency**: Quickly access relevant customer details while working on opportunities without switching between projects.
- **Data Integrity**: Reduce errors from manual data entry by automatically pulling in customer information.
- **Holistic View**: Easily see all opportunities associated with a customer by using the "Referenced By" field in your Customer Database.

### Advanced Tip: Lookup a Lookup

Blue offers an advanced feature called "Lookup a Lookup" that can be incredibly useful for complex CRM setups. This feature allows you to create connections across multiple projects, enabling you to access information from both your Customer Database and Opportunities project in a third project.

For instance, let's say you have a "Projects" workspace where you manage the actual work for your clients. You want this workspace to have access to both customer details and opportunity information. Here's how you can set this up:

First, create a Reference field in your Projects workspace that links to the Opportunities project. This establishes the initial connection. Next, create Lookup fields based on this Reference to pull in specific details from the opportunities, such as deal value or expected close date.

The real power comes in the next step: you can create additional Lookup fields that reach through the opportunity's Reference to the Customer Database. This allows you to pull in customer information like contact details or account status directly into your Projects workspace.

This chain of connections gives you a comprehensive view in your Projects workspace, combining data from both your opportunities and customer database. It's a powerful way to ensure that your project teams have all the relevant information at their fingertips without needing to switch between different projects.

### Best Practices for Linked CRM Systems

Maintain your Customer Database as the single source of truth for all customer information. Whenever you need to update customer details, always do it in the Customer Database first. This ensures that the information remains consistent across all linked projects.

When creating Reference and Lookup fields, use clear and meaningful names. This helps maintain clarity, especially as your system grows more complex. 

Regularly review your setup to ensure you're pulling in the most relevant information. As your business needs evolve, you might need to add new Lookup fields or remove ones that are no longer useful. Periodic reviews help keep your system streamlined and effective.

Consider leveraging Blue's automation features to keep your data synchronized and up-to-date across projects. For example, you could set up an automation to notify relevant team members when key customer information is updated in the Customer Database.

By effectively implementing these strategies and making full use of Reference and Lookup fields, you can create a powerful, interconnected CRM system in Blue. This system will provide you with a comprehensive 360-degree view of your customer relationships and sales pipeline, enabling more informed decision-making and smoother operations across your organization.

## Dashboards

Dashboards are a crucial component of any effective CRM system, providing at-a-glance insights into your sales performance and customer relationships. Blue's dashboard feature is particularly powerful because it allows you to combine real-time data from multiple projects simultaneously, giving you a comprehensive view of your sales operations.

When setting up your CRM dashboard in Blue, consider including several key metrics. Pipeline generated per month shows the total value of new opportunities added to your pipeline, helping you track your team's ability to generate new business. Sales per month displays your actual closed deals, allowing you to monitor your team's performance in converting opportunities into sales.

Introducing the concept of pipeline discounts can lead to more accurate forecasting. For example, you might count 90% of the value of deals in the "Contract Out for Signature" stage, but only 50% of deals in the "Proposal Sent" stage. This weighted approach provides a more realistic sales forecast.

Tracking new opportunities per month helps you monitor the number of new potential deals entering your pipeline, which is a good indicator of your sales team's prospecting efforts. Breaking down sales by type can help you identify your most successful offerings. If you set up an invoice tracking project linked to your opportunities, you can also track actual revenue on your dashboard, providing a complete picture from opportunity to cash.

Blue offers several powerful features to help you create an informative and interactive CRM dashboard. The platform provides three main types of charts: stat cards, pie charts, and bar graphs. Stat cards are ideal for displaying key metrics like total pipeline value or number of active opportunities. Pie charts are perfect for showing the composition of your sales by type or the distribution of deals across different stages. Bar graphs excel at comparing metrics over time, such as monthly sales or new opportunities.

Blue's sophisticated filtering capabilities allow you to segment your data by project, list, tag, and timeframe. This is particularly useful for drilling down into specific aspects of your sales data or comparing performance across different teams or products. The platform smartly consolidates lists and tags with the same name across projects, enabling seamless cross-project analysis. This is invaluable for a CRM setup where you might have separate projects for customers, opportunities, and invoices.

Customization is a key strength of Blue's dashboard feature. The drag-and-drop functionality and display flexibility allow you to create a dashboard that perfectly suits your needs. You can easily rearrange charts and choose the most appropriate visualization for each metric.
While dashboards are currently for internal use only, you can easily share them with team members, granting either view or edit permissions. This ensures that everyone in your sales team has access to the insights they need.

By leveraging these features and including the key metrics we've discussed, you can create a comprehensive CRM dashboard in Blue that provides real-time insights into your sales performance, pipeline health, and overall business growth. This dashboard will become an invaluable tool for making data-driven decisions and keeping your entire team aligned on your sales goals and progress.

## Conclusion

Setting up a comprehensive sales CRM in Blue is a powerful way to streamline your sales process and gain valuable insights into your customer relationships and business performance. By following the steps outlined in this guide, you've created a robust system that integrates customer information, sales opportunities, and performance metrics into a single, cohesive platform.

We started by creating a customer database, establishing a single source of truth for all your customer information. This foundation allows you to maintain accurate and up-to-date records for all your clients and prospects. We then built upon this with an opportunities database, enabling you to track and manage your sales pipeline effectively.

One of the key strengths of using Blue for your CRM is the ability to link these databases using reference and lookup fields. This integration creates a dynamic system where updates to customer information are instantly reflected across all related opportunities, ensuring data consistency and saving time on manual updates.
We explored how to leverage Blue's powerful automation features to streamline your workflow, from assigning new leads to sending follow-up reminders. These automations help ensure that no opportunities fall through the cracks and that your team can focus on high-value activities rather than administrative tasks.

Finally, we delved into creating dashboards that provide at-a-glance insights into your sales performance. By combining data from your customer and opportunity databases, these dashboards offer a comprehensive view of your sales pipeline, closed deals, and overall business health.


Remember, the key to getting the most out of your CRM is consistent use and regular refinement. Encourage your team to fully adopt the system, regularly review your processes and automations, and continue to explore new ways to leverage Blue's features to support your sales efforts.

With this sales CRM setup in Blue, you're well-equipped to nurture customer relationships, close more deals, and drive your business forward. 