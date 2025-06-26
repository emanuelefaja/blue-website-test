---
title: Reference & lookup custom fields
slug: reference-lookups-custom-fields
tags: ["product-updates"]
description: Effortlessly create interconnected projects in Blue, transforming it into a single source of truth for your business with the new Reference and Lookup Fields.
image: /resources/link-background.png
date: 1/11/2023
showdate: true
sitemap:
  loc: /resources/reference-lookups-custom-fields
  lastmod: 2023/11/01
  changefreq: monthly
---

Projects in Blue are already a powerful way to manage your business data and move work forwards.

Today, we're taking the next logical step and allowing you to interconnect your data *between* projects for the ultimate flexibility and power. 

Interconnecting projects within Blue transforms it into a single source of truth for your business. This capability allows for the creation of a comprehensive and interconnected dataset, enabling seamless data flow and enhanced visibility across projects. By linking projects, teams can achieve a unified view of operations, enhancing decision-making and operational efficiency.

## An example

Consider ACME Company, which uses Blue's Reference and Lookup custom fields to create an interconnected data ecosystem across its Customer, Sales, and Inventory projects. Customer records in the Customers project are linked via Reference fields to sales transactions in the Sales project. This linkage allows Lookup fields to pull associated customer details, such as phone numbers and account statuses, directly into each sales record. Additionally, inventory items sold are displayed in the sales record through a Lookup field referencing the Quantity Sold data from the Inventory project. Finally, inventory withdrawals are connected to relevant sales via a Reference field in Inventory, pointing back to the Sales records. This setup provides full visibility into which sale triggered the inventory removal, creating an integrated 360-degree view across projects.

## How Reference Fields Work

Reference custom fields enable you to create relationships between records across different projects in Blue. When creating a Reference field, the Project Administrator selects the specific project that will provide the list of reference records. Configuration options include:

* **Single Select**: Allows choosing one reference record.
* **Multi-select**: Allows choosing multiple reference records.
* **Filtering**: Set filters to allow users to select only records that match the filter criteria.

Once set up, users can select specific records from the dropdown menu within the Reference field, establishing a link between projects.

## Extending reference fields using lookups

Lookup custom fields are used to import data from records in other projects, creating one-way visibility. They are always read-only and are connected to a specific Reference custom field. When a user selects one or more records using a Reference custom field, the Lookup custom field will show data from those records. Lookups can display data such as:

* Created at
* Updated at
* Due Date
* Description
* List
* Tag
* Assignee
* Any supported custom field from the referenced record — including other lookup fields!


For example, imagine a scenario where you have three projects: **Project A** is a sales project, **Project B** is an inventory management project, and **Project C** is a customer relationship project. In Project A, you have a Reference custom field that links sales records to the corresponding customer records in Project C. In Project B, you have a Lookup custom field that imports information from Project A, such as the quantity sold. This way, when a sales record is created in Project A, the customer information associated with that sale is automatically pulled in from Project C, and the quantity sold is automatically pulled in from Project B. This allows you to keep all relevant information in one place and view without having to create duplicate data or manually update records across projects.

A real-life example of this is an e-commerce company that uses Blue to manage their sales, inventory, and customer relationships. In their **Sales** project, they have a Reference custom field that links each sales record to the corresponding **Customer** record in their **Customers** project. In their **Inventory** project, they have a Lookup custom field that imports information from the Sales project, such as the quantity sold, and displays it in the inventory item record. This allows them to easily see which sales are driving inventory removals and keep their inventory levels up to date without having to manually update records across projects.

## Conclusion

Imagine a world where your project data isn't siloed but flows freely between projects, providing comprehensive insights and driving efficiency. That's the power of Blue's Reference and Lookup fields. By enabling seamless data connections and providing real-time visibility across projects, these features transform how teams collaborate and make decisions. Whether you're managing customer relationships, tracking sales, or overseeing inventory, Reference and Lookup fields in Blue empower your team to work smarter, faster, and more effectively. Dive into the interconnected world of Blue and watch your productivity soar.

[Check out the documentation](https://documentation.blue.cc/custom-fields/reference) or [sign up and try it for yourself.](https://app.blue.cc)