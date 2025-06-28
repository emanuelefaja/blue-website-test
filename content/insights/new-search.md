---
title: Real-time search
slug: new-search
category: "Product Updates"
description: Blue unveils a new blazing-fast search engine that returns results across all your projects in milliseconds, empowering you to switch context in a blink of an eye. 
image: /resources/search-background.png
date: 2024-03-01
showdate: true
---

We are thrilled to announce the launch of our new search engine, designed to revolutionize how you find information within Blue. Efficient search functionality is crucial for seamless project management, and our latest update ensures that you can access your data faster than ever.

Our new search engine allows you to search for all comments, files, records, custom fields, descriptions, and checklists. Whether you need to find a specific comment made on a project, quickly locate a file, or search for a particular record or field, our search engine provides lightning-fast results.

As tools approach the 50-100ms responsiveness, they tend to fade away and blend into the background, providing a seamless and almost invisible user experience. For context, a human blink takes approximately 60-120ms, so 50ms is actually faster than a blink of an eye! This level of responsiveness allows you to interact with Blue without even realizing it's there, freeing you up to focus on the actual work at hand. By leveraging this level of performance, our new search engine ensures that you can quickly access the information you need, without it ever getting in the way of your workflow.

To achieve our goal of lightning-fast search, we leveraged the latest open-source technologies. Our search engine is built on top of MeiliSearch, a popular open-source search-as-a-service that uses natural language processing and vector search to quickly find relevant results. Additionally, we implemented in-memory storage, which allows us to store frequently accessed data in RAM, reducing the time it takes to return search results. This combination of MeiliSearch and in-memory storage enables our search engine to deliver results in milliseconds, making it possible for you to quickly find what you need without ever having to think about the underlying technology.

The new search bar is conveniently located on the navigation bar, allowing you to start searching right away. For a more detailed search experience, simply press the Tab key while searching to access the full search page. Additionally, you can quickly activate the search function from anywhere using the CMD/Ctrl+K shortcut, making it even easier to find what you need.

<video autoplay loop muted playsinline>
  <source src="/resources/search-demo.mp4" type="video/mp4">
</video>


## Future Developments

This is just the start. Now that we have a next-generation search infrastructure, we can do some really interesting things in the future.

Next up is going to be semantic search, which is a significant improvement to the typical keyword search. Allow us to explain. 

This feature will allow the search engine to understand the context of your queries. For example, searching for "sea" will retrieve relevant documents even if the exact phrase isn't used. You might be thinking "but I typed 'ocean' instead!" - and you're right. The search engine will also understand the similarity between "sea" and "ocean", and return relevant documents even if the exact phrase isn't used. This feature is particularly useful when searching for documents containing technical terms, acronyms, or just common words that have multiple variations or typos.
 
Another upcoming feature is the ability to search for images by their content. To achieve this, we will be processing every image in your project, creating an embedding for each one. In high-level terms, an embedding is a mathematical set of coordinates that corresponds to the meaning of an image. This means that all images can be searched based on what they contain, regardless of their filename or metadata. Imagine searching for "flowchart" and finding all images related to flowcharts, *regardless of their filenames.*