---
title:  AI Auto Categorization (Engineering Deep Dive)
slug: ai-auto-categorization-engineering
tags: ["engineering-blog"]
description: Go behind the scenes with the Blue engineering team as they explain how to built an AI-powered auto-categorization and tagging feature.
image: /resources/ai-tag-background.png
date: 12/07/2024
showdate: true
sitemap:
  loc: /resources/ai-auto-categorization-engineering
  lastmod: 2024/07/12
  changefreq: monthly
---

We recently released [AI Auto Categorization](/resources/ai-auto-categorization) to all Blue users. This is an AI feature that is bundled in to the core subscription of Blue, without additional costs. In this post, we dig into the engineering behind making this feature happen.

---
At Blue, our approach to feature development is rooted in a deep understanding of user needs and market trends, coupled with a commitment to maintaining the simplicity and ease of use that defines our platform. This is what drives our [roadmap](/platform/roadmap), and what has [allowed us to consistently ship features each month for years](/platform/changelog). 

The introduction of AI-powered auto-tagging to Blue is a perfect example of this philosophy in action. Before we dive into the technical details of how we built this feature, it's crucial to understand the problem we were solving and the careful consideration that went into its development.

The project management landscape is rapidly evolving, with AI capabilities becoming increasingly central to user expectations. Our customers, particularly those managing large-scale [projects](/platform) with millions of [records](/platform/records), had been vocal about their desire for smarter, more efficient ways to organize and categorize their data. 

However, at Blue, we don't simply add features because they're trendy or requested. Our philosophy is that every new addition must prove its worth, with the default answer being a firm *"no"* until a feature demonstrates strong demand and clear utility.

To truly understand the depth of the problem and the potential of AI auto-tagging, we conducted extensive customer interviews, focusing on long-time users who manage complex, data-rich projects across multiple domains. 

These conversations revealed a common thread: *while tagging was invaluable for organization and searchability, the manual nature of the process was becoming a bottleneck, especially for teams dealing with high volumes of records.*

But we saw beyond just solving the immediate pain point of manual tagging. 

We envisioned a future where AI-powered tagging could become the foundation for more intelligent, automated workflows. 

The real power of this feature, we realized, lay in its potential integration with our [project management automation system](/platform/project-management-automation). Imagine a project management tool that not only categorizes information intelligently but also uses those categories to route tasks, trigger actions, and adapt workflows in real-time. 

This vision aligned perfectly with our goal of keeping Blue simple yet powerful.

Furthermore, we recognized the potential to extend this capability beyond the confines of our platform. By developing a robust AI tagging system, we were laying the groundwork for a "categorization API" that could work out of the box, potentially opening up new avenues for how our users interact with and leverage Blue in their broader tech ecosystems.

This feature, therefore, wasn't just about adding an AI checkbox to our feature list. 

It was about taking a significant step towards a more intelligent, adaptive project management platform while staying true to our core philosophy of simplicity and user-centricity.

In the following sections, we'll dive into the technical challenges we faced in bringing this vision to life, the architecture we designed to support it, and the solutions we implemented. We'll also explore the future possibilities this feature opens up, showcasing how a carefully considered addition can pave the way for transformative changes in project management.

---
## The Problem

As discussed above, manual tagging of project records can be time-consuming and inconsistent. 

We set out to solve this by leveraging AI to automatically suggest tags based on record content. 

The main challenges were:

1. Choosing an appropriate AI model
2. Efficiently processing large volumes of records
3. Ensuring data privacy and security
4. Integrating the feature seamlessly into our existing architecture

## Selecting the AI Model

We evaluated several AI platforms, including [OpenAI](https://openai.com), open-source models on [HuggingFace](https://huggingface.co/), and [Replicate](https://replicate.com). 

Our criteria included:

- Cost-effectiveness
- Accuracy in understanding context
- Ability to adhere to specific output formats
- Data privacy guarantees

After thorough testing, we chose OpenAI's [GPT-3.5 Turbo](https://platform.openai.com/docs/models/gpt-3-5-turbo). While [GPT-4](https://softgist.com/the-ultimate-guide-to-prompt-engineering) might offer marginal improvements in accuracy, our tests showed that GPT-3.5's performance was more than adequate for our auto-tagging needs. The balance of cost-effectiveness and strong categorization capabilities made GPT-3.5 the ideal choice for this feature.


The higher cost of GPT-4 would have forced us to offer the feature as a paid add-on, conflicting with our goal of **bundling AI within our main product at no additional cost to end users.** 

As of our implementation, the pricing for GPT-3.5 Turbo is:

- $0.0005 per 1K input tokens (or $0.50 per 1M input tokens)
- $0.0015 per 1K output tokens (or $1.50 per 1M output tokens)

Let's make some assumptions about an average record in Blue:

- **Title**: ~10 tokens
- **Description**: ~50 tokens
- **2 comments**: ~30 tokens each
- **5 custom fields**: ~10 tokens each
- **List name, due date, and other metadata**: ~20 tokens
- **System prompt and available tags**: ~50 tokens

Total input tokens per record: 10 + 50 + (30 * 2) + (10 * 5) + 20 + 50 ≈ 240 tokens

For the output, let's assume an average of 3 tags suggested per record, which might total around 20 output tokens including the JSON formatting.

For 1 million records:

- Input cost: (240 * 1,000,000 / 1,000,000) * $0.50 = $120
- Output cost: (20 * 1,000,000 / 1,000,000) * $1.50 = $30

**Total cost for auto-tagging 1 million records: $120 + $30 = $150**

## GPT3.5 Turbo Performance

Categorization is a task that large language models (LLMs) like  GPT-3.5 Turbo excel at, making them particularly well-suited for our auto-tagging feature. LLMs are trained on vast amounts of text data, allowing them to understand context, semantics, and relationships between concepts. This broad knowledge base enables them to perform categorization tasks with high accuracy across a wide range of domains.

For our specific use case of project management tagging,  GPT-3.5 Turbo demonstrates several key strengths:

- **Contextual Understanding:** Can grasp the overall context of a project record, considering not just individual words but the meaning conveyed by the entire description, comments, and other fields.
- **Flexibility:** Can adapt to various project types and industries without requiring extensive reprogramming.
- **Handling Ambiguity:** Can weigh multiple factors to make nuanced decisions.
- **Learning from Examples:** Can quickly understand and apply new categorization schemes without additional training.
- **Multi-label Classification:** Can suggest multiple relevant tags for a single record, which was crucial for our requirements. 

GPT-3.5 Turbo also stood out for its reliability in adhering to our required JSON output format, which was *crucial* for seamless integration with our existing systems. Open-source models, while promising, often added extra comments or deviated from the expected format, which would have required additional post-processing. This consistency in output format was a key factor in our decision, as it significantly simplified our implementation and reduced potential points of failure.

Opting for GPT-3.5 Turbo with its consistent JSON output allowed us to implement a more straightforward, reliable, and maintainable solution. 

Had we chosen a model with less reliable formatting, we would have faced a cascade of complications: the need for robust parsing logic to handle various output formats, extensive error handling for inconsistent outputs, potential performance impacts from additional processing, increased testing complexity to cover all output variations, and a greater long-term maintenance burden. 

Parsing errors could lead to incorrect tagging, negatively impacting user experience. By avoiding these pitfalls, we were able to focus our engineering efforts on critical aspects like performance optimization and user interface design, rather than wrestling with unpredictable AI outputs.

## System Architecture

Our AI auto-tagging feature is built on a robust, scalable architecture designed to handle high volumes of requests efficiently while providing a seamless user experience. As with all our systems, we've architected this feature to support one order of magnitude more traffic than we currently experience. This approach, while seemingly overengineered for current needs, is a best practice that allows us to seamlessly handle sudden spikes in usage and gives us ample runway for growth without major architectural overhauls. Otherwise, we would have to reengineer all of our systems every 18 months — something that we have learned the hard way in the past! 

Let's break down the components and flow of our system:


- **User Interaction:** The process begins when a user presses the "Autotag" button in the Blue interface. This action triggers the auto-tagging workflow.
- **Blue API Call:** The user's action is translated into an API call to our Blue backend. This API endpoint is designed to handle auto-tagging requests.
- **Queue Management:** Instead of processing the request immediately, which could lead to performance issues under high load, we add the tagging request to a queue. We use Redis for this queuing mechanism, which allows us to manage load effectively and ensure system scalability.
- **Background Service:** We implemented a background service that continuously monitors the queue for new requests. This service is responsible for processing queued requests.
- **OpenAI API Integration:** The background service prepares the necessary data and makes API calls to OpenAI's GPT-3.5 model. This is where the actual AI-powered tagging occurs. We send relevant project data and receive suggested tags in return.
- **Result Processing:** The background service processes the results received from OpenAI. This involves parsing the AI's response and preparing the data for application to the project.
- **Tag Application:** The processed results are used to apply the new tags to the relevant items in the project. This step updates our database with the AI-suggested tags.
- **User Interface Reflection:** Finally, the new tags appear in the user's project view, completing the auto-tagging process from the user's perspective.

This architecture offers several key benefits that enhance both system performance and user experience. By utilizing a queue and background processing, we've achieved impressive scalability, allowing us to handle numerous requests simultaneously without overwhelming our system or hitting the rate limits of the OpenAI API. Implementing this architecture required careful consideration of various factors to ensure optimal performance and reliability. For queue management, we chose Redis, leveraging its speed and reliability in handling distributed queues. 

This approach also contributes to the overall responsiveness of the feature. Users receive immediate feedback that their request is being processed, even if the actual tagging takes some time, creating a sense of real-time interaction. The architecture's fault tolerance is another crucial advantage. If any part of the process encounters issues, such as temporary OpenAI API disruptions, we can gracefully retry or handle the failure without impacting the entire system. 

This robustness, combined with the real-time appearance of tags, enhances the user experience, giving the impression of AI "magic" at work.

## Data & Prompts

A crucial step in our auto-tagging process is preparing the data to be sent to the GPT-3.5 model. This step required careful consideration to balance providing enough context for accurate tagging while maintaining efficiency and protecting user privacy. Here's a detailed look at our data preparation process.

For each record, we compile the following information:

- **List Name**: Provides context about the broader category or phase of the project.
- **Record Title**: Often contains key information about the record's purpose or content.
- **Custom Fields**: We include text and number-based [custom fields](/platform/custom-fields), which often contain crucial project-specific information.
- **Description**: Typically contains the most detailed information about the record.
- **Comments**: Can provide additional context or updates that might be relevant for tagging.
- **Due Date**: Temporal information that might influence tag selection.

Interestingly, we do not send existing tag data to GPT-3.5, and we do this to avoid biasing the model.

The core of our auto-tagging feature lies in how we interact with the GPT-3.5 model and process its responses. This section of our pipeline required careful design to ensure accurate, consistent, and efficient tagging.

We use a carefully crafted system prompt to instruct the AI on its task. Here's a breakdown of our prompt and the rationale behind each component:

```
You will be provided with record data, and your task is to choose the tags that are relevant to the record.
You can respond with an empty array if you are unsure.
Available tags: ${tags}.
Today: ${currentDate}.
Please respond in JSON using the following format:
{ "tags": ["tag-1", "tag-2"] }
```

- **Task Definition:** We clearly state the AI's task to ensure focused responses.
- **Uncertainty Handling:** We explicitly allow for empty responses, preventing forced tagging when the AI is unsure.
- **Available Tags:** We provide a list of valid tags (${tags}) to constrain the AI's choices to existing project tags.
- **Current Date:** Including ${currentDate} helps the AI understand the temporal context, which can be crucial for certain types of projects.
- **Response Format:** We specify a JSON format for easy parsing and error checking.

This prompt is the result of extensive testing and iteration. We found that being explicit about the task, available options, and desired output format significantly improved the accuracy and consistency of the AI's responses — simplicity is key! 

 The list of available tags is generated server-side and validated before inclusion in the prompt. We implement strict character limits on tag names to prevent oversized prompts.

As mentioned above, we had no issue with GPT-3.5 Turbo in getting back the pure JSON response in the correct format 100% of the time. 

So in summary,

- We combine the system prompt with the prepared record data.
- This combined prompt is sent to the GPT-3.5 model via OpenAI's API.
- We use a temperature setting of 0.3 to balance creativity and consistency in the AI's responses.
- Our API call includes a max_tokens parameter to limit response size and control costs.

Once we receive the AI's response, we go through several steps to process and apply the suggested tags:

* **JSON Parsing**: We attempt to parse the response as JSON. If parsing fails, we log the error and skip tagging for that record.
* **Schema Validation**: We validate the parsed JSON against our expected schema (an object with a "tags" array). This catches any structural deviations in the AI's response.
* **Tag Validation**: We cross-reference the suggested tags against our list of valid project tags. This step filters out any tags that don't exist in the project, which could occur if the AI misunderstood or if project tags changed between prompt creation and response processing.
* **Deduplication**: We remove any duplicate tags from the AI's suggestion to avoid redundant tagging.
* **Application**: The validated and deduplicated tags are then applied to the record in our database.
* **Logging and Analytics**: We log the the final applied tags. This data is valuable for monitoring the system's performance and improving it over time.

## Challenges

Implementing AI-powered auto-tagging in Blue presented several unique challenges, each requiring innovative solutions to ensure a robust, efficient, and user-friendly feature.

### Undo Bulk Operation

The AI Tagging feature can be done both on individual records as well as in bulk. The problem with the bulk operation is that if the user does not like the result, they would have to manually go through thousands of records and undo the work of the AI. Clearly, that's unacceptable. 

To solve this, we implemented an innovative tagging session system. Each bulk tagging operation is assigned a unique session ID, which is associated with all tags applied during that session. This allows us to efficiently manage undo operations by simply deleting all tags associated with a particular session ID. We also remove related audit trails, ensuring that undone operations leave no trace in the system. This approach gives users the confidence to experiment with AI tagging, knowing they can easily revert changes if needed.

### Data Privacy

Data privacy was another critical challenge we faced. 

Our users trust us with their project data, and it was paramount to ensure this information wasn't retained or used for model training by OpenAI. We tackled this on multiple fronts. 

First, we formed an agreement with OpenAI that explicitly prohibits the use of our data for model training. Additionally, OpenAI deletes the data after processing, providing an extra layer of privacy protection. 

On our end, we took the precaution of excluding sensitive information, such as assignee details, from the data sent to the AI so this ensures that specific individuals names are not sent to third parties along with other data. 

This comprehensive approach allows us to leverage AI capabilities while maintaining the highest standards of data privacy and security.

### Rate Limits and Catching Errors

One of our primary concerns was scalability and rate limiting. Direct API calls to OpenAI for each record would have been inefficient and could quickly hit rate limits, especially for large projects or during peak usage times. To address this, we developed a background service architecture that allows us to batch requests and implement our own queuing system. This approach helps us manage API call frequency and enables more efficient processing of large volumes of records, ensuring smooth performance even under heavy load.

The nature of AI interactions meant we also had to prepare for occasional errors or unexpected outputs. There were instances where the AI might produce invalid JSON or outputs that didn't match our expected format. To handle this, we implemented robust error handling and parsing logic throughout our system. If the AI response isn't valid JSON or doesn't contain the expected "tags" key, our system is designed to treat it as if no tags were suggested, rather than attempting to process potentially corrupt data. This ensures that even in the face of AI unpredictability, our system remains stable and reliable.

## Future Developments

We believe that features, and the Blue product as a whole, is never "done" — there is always room for improvement. 

There were some features that we considered in the initial build that did not pass the scoping phase, but are interesting to note as we will likely implement some version of them in the future. 

The first is adding tag description. This would allow end users to not only give tags a name and a colour, but also an optional description. This would be also passed to the AI to help provide further context and potentially improve accuracy.

While additional context could be valuable, we're mindful of the potential complexity it might introduce. There's a delicate balance to strike between providing useful information and overwhelming users with too much detail. As we develop this feature, we'll be focusing on finding that sweet spot where added context enhances rather than complicates the user experience.

Perhaps the most exciting enhancement on our horizon is the integration of AI auto-tagging with our [project management automation system](/platform/project-management-automation). 

This means that the AI tagging feature could either be a trigger, or an action from an automation. This could be huge as it could turn this rather simple AI categorization feature into an AI-based routing system for work.

Imagine an automation that states:

When AI tags a record as "Critical" -> Assigned to "Manager" and Send a Custom Email

This means that when you AI-tag a record, if the AI decides it is a critical issue, then it can automatically assign the project manager and send them a custom email. This extends the [benefits of our project management automation system](/platform/project-management-automation) from purely a rule-based system into a true flexible AI system. 

By continuously exploring the frontiers of AI in project management, we aim to provide our users with tools that not only meet their current needs but anticipate and shape the future of work. As always, we'll be developing these features in close collaboration with our user community, ensuring that every enhancement adds real, practical value to the project management process.

## Conclusion

So that's it!  

This was a fun feature to implement, and one our first steps into AI, alongside the [AI Content Summarization](/resources/ai-content-summarization) we've previously launched. We know that AI is going to be a bigger and bigger role in project management in the future, and we can't wait to roll out more innovative features leveraging advanced LLMs (Large Language Models). 

There was quite a bit to think about while implementing this, and we're especially excited about how we can leverage this feature in the future with Blue's existing [project management automation engine](/resources/benefits-project-management-automation). 

We also hope it's been an interesting read, and that it gives you a glimpse into how we think about engineering the features that you use everyday. 
