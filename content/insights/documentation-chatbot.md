---
title: Why We Built Our Own AI Documentation Chatbot
slug: documentation-chatbot
category: "Product Updates"
description: We built our own documentation AI Chatbot that is trained on the Blue platform documentation.
image: /insights/chat-background.png
date: 2024-07-09
showdate: true
---

At Blue, we're always looking to find ways to make life easier for our customers. We have [in-depth documentation of every feature](https://documentation.blue.cc), [YouTube Videos](https://www.youtube.com/@HeyBlueTeam), [Tips & Tricks](/insights/tips-tricks), and [various support channels](/support). 

We've been keeping a close eye on the development of AI (Artificial Intelligence) as we are very much into [project management automations](/platform/features/automations). We also released features such as [AI Auto Categorization](/insights/ai-auto-categorization) and [AI Summaries](/insights/ai-content-summarization) to make work easier for our thousands of customers. 

One thing that is clear is that AI is here to stay, and it is going to have an incredible effect across most industries, and project management is no exception. So we asked ourselves how we could further leverage AI to help the full lifecycle of a customer, from discovery, pre-sales, onboarding, and also with ongoing questions.

The answer was quite clear: **We needed an AI chatbot trained on our documentation.**

Let's face it: *every* organization should probably have a chatbot. They are great ways for customers to get instant answers to typical questions, without having to dig through pages of dense documentation or your website. The importance of chatbots in modern marketing websites cannot be overstated. 

![](/insights/ai-chatbot-regular.png)

For software companies specifically, one should not consider the marketing website as a separate "thing" — it *is* part of your product. This is because it fits in the typical customer life:

- **Awareness** (Discovery): This is where potential customers first stumble upon your awesome product. Your chatbot can be their friendly guide, pointing them to key features and benefits right off the bat.
- **Consideration** (Education): Now they're curious and want to learn more. Your chatbot becomes their personal tutor, dishing out information tailored to their specific needs and questions.
- **Purchase/Conversion**: This is the moment of truth - when a prospect decides to take the plunge and become a customer. Your chatbot can smooth out any last-minute hiccups, answer those "just before I buy" questions, and maybe even throw in a sweet deal to seal the deal.
- **Onboarding**: They've bought in, now what? Your chatbot transforms into a helpful sidekick, walking new users through setup, showing them the ropes, and making sure they don't feel lost in your product's wonderland.
- **Retention**: Keeping customers happy is the name of the game. Your chatbot is on call 24/7, ready to troubleshoot issues, offer tips and tricks, and make sure your customers feel the love.
- **Expansion**: Time to level up! Your chatbot can subtly suggest new features, upsells, or cross-sells that align with how the customer is already using your product. It's like having a really smart, non-pushy salesperson always on standby.
- **Advocacy**: Happy customers become your biggest cheerleaders. Your chatbot can encourage satisfied users to spread the word, leave reviews, or participate in referral programs. It's like having a hype machine built right into your product!

## Build vs Buy Decision

Once we decided to implement an AI chatbot, the next big question was: build or buy? As a small team laser-focused on our core product, we generally prefer "as-a-service" solutions or popular open-source platforms. We're not in the business of reinventing the wheel for every part of our tech stack, after all.
So, we rolled up our sleeves and dove into the market, hunting for both paid and open-source AI chatbot solutions. 

Our requirements were straightforward, but non-negotiable:

- **Unbranded Experience**: This chatbot isn't just a nice-to-have widget; it's going on our marketing website and eventually in our product. We're not keen on advertising someone else's brand in our own digital real estate.
- **Great UX**: For many potential customers, this chatbot might be their first point of contact with Blue. It sets the tone for their perception of our company. Let's face it: if we can't nail a proper chatbot on our website, how can we expect customers to trust us with their mission-critical projects and processes?
- **Reasonable Cost**: With a large user base and plans to integrate the chatbot into our core product, we needed a solution that wouldn't break the bank as usage scales up. Ideally, we wanted a **BYOK (Bring Your Own Key) option**. This would allow us to use our own OpenAI or other AI service key, paying for direct variable costs instead of a markup to a third-party vendor that doesn't actually run the models.
- **OpenAI Assistants API Compatible**: If we were going to go with an open-source software, we did not want to have the hassle of having to manage a pipeline for document ingestion, indexing, vector databases, and all of that. We wanted to use the [OpenAI Assistants API](https://platform.openai.com/docs/assistants/overview) that would abstract away all the complexity behind an API. Honestly — this is really well done. 
- **Scalability**: We want to have this chatbot in multiple places, with potentially tens of thousands of users a year. We expect significant usage, and we don't want to be locked into a solution that can't scale with our needs.

## Commercial AI Chatbots

The ones we reviewed tended to have a better UX than open-source solutions — as unfortunately is often the case. There is probably a separate discussion to be had one day about *why* many open-source solutions ignore or underplay the importance of UX. 

We'll provide a list here, in case you're looking for some solid commercial offerings:

- **[Chatbase](https://chatbase.co):** Chatbase allows you to build a custom AI chatbot trained on your knowledge base and add it to your website or interact with it through their API. It offers features like trustworthy answers, lead generation, advanced analytics, and the ability to connect to multiple data sources. To us, this felt like one of the most polished commercial offerings out there. 
- **[DocsBot AI](https://docsbot.ai/):** DocsBot AI creates custom ChatGPT bots trained on your documentation and content for support, presales, research, and more. It provides embeddable widgets to easily add the chatbot to your website, the ability to reply to support tickets automatically, and a powerful API for integration.
- **[CustomGPT.ai](https://customgpt.ai):** CustomGPT.ai creates a personal chatbot experience by ingesting your business data, including website content, helpdesk, knowledge bases, documents, and more. It allows leads to ask questions and get instant answers based on your content, without needing to search. Interestingly, they also [claim to beat OpenAI at RAG (Retrieval Augemented Generation)!](https://customgpt.ai/customgpt-beats-open-ai-in-rag-benchmark/)
- **[OpenAssistantGPT](https://www.openassistantgpt.io/)**: This is an interesting commercial offering, because it *also* happens to be open-source software. It seems a little early stage, and the pricing did not feel realistic ($27/month for unlimited messages will never work commercially for them).

We also looked at [InterCom Fin](https://www.intercom.com/fin) which is part of their customer support software. This would have meant switching away from [HelpScout](https://wwww.helpscout.com) which we have used since we started Blue. This could have been possible, but InterCom Fin has some crazy pricing that simply excluded it out of consideration.

And this is actually the problem with many of the commercial offerings. InterCom Fin charges $0.99 per customer support request handled, and ChatBase charges $399/month for 40,000 messages. That's almost $5k a year for a simple chat widget. 

Considering that the prices for AI inference are dropping like crazy. OpenAI reduced their prices quite dramatically:

- The original GPT-4 (8k context) was priced at $0.03 per 1K prompt tokens.
- The GPT-4 Turbo (128k context) was priced at $0.01 per 1K prompt tokens, a 50% reduction from the original GPT-4.
- The GPT-4o model is priced at $0.005 per 1K tokens, which is a further 50% reduction from the GPT-4 Turbo pricing.

That's an 83% reduction in costs, and we do not expect that to stay stagnant. 

Considering that we were looking for a scalable solution that would be used by tens of thousands of users a year with a significant amount of messages, it makes sense to go directly to the source and pay for the API costs directly, not use a commercial version which marksup the costs.

## Open Source AI Chatbots

As mentioned, the open source options we reviewed were mostly dissapointing with regards to the "Great UX" requirement. 

We looked at:

- **[Deepchat](https://deepchat.dev/)**: This is a framework-agnostic chat component for AI services, which connects to various AI APIs including OpenAI. It also has the ability for users to download an AI model that runs directly in the browser. We played around with this and got a version working, but the OpenAI Assistants API implemented felt quite buggy with several issues. However, this is a very promising project, and their playground is really nicely done. 
- **[OpenAssistantGPT](https://www.openassistantgpt.io/)**: Looking at this again from an open-source perspective, this would need us to spin up quite a bit of infrastructure, something that we did not want to do, because we wanted to rely as much as possible on OpenAI's Assistants API. 


## Building Our Own ChatBot

And so, without being able to find something that matched all of our requirements, we decided to build our own AI chatbot that could interface with the OpenAI Assistants API. This, in the end, turned out to be relatively painless! 

Our website uses [Nuxt3](https://nuxt.com), [Vue3](https://vuejs.org/) (which is the same framework as the Blue Platform), and [TailwindUI](https://tailwindui.com/).

The first step was to create an API (Application Programming Interface) in Nuxt3 that can "speak" to the OpenAI Assistants API. This was necessary as we did not want to do everything on the front-end, as this would expose our OpenAI API key to the world, with the potential for abuse. 

Our backend API acts as a secure middleman between the user's browser and OpenAI. Here's what it does:

- **Conversation Management:** It creates and manages "threads" for each conversation. Think of a thread as a unique chat session that remembers everything you've said.
- **Message Handling:** When you send a message, our API adds it to the right thread and asks OpenAI's assistant to craft a response.
- **Smart Waiting:** Instead of making you stare at a loading screen, our API checks in with OpenAI every second to see if your response is ready. It's like having a waiter who keeps an eye on your order without bothering the chef every two seconds.
- **Security First:** By handling all this on the server, we keep your data and our API keys safe and sound.

Then, there was the front-end and user experience. As discussed earlier, this was *critically* important, because we do not get a second chance at making a first impression! 

In designing our chatbot, we paid meticulous attention to user experience, ensuring that every interaction is smooth, intuitive, and reflective of Blue's commitment to quality. The chatbot interface begins with a simple, elegant Blue circle,, using [HeroIcons for our icons](https://heroicons.com/) (which we use throughout the Blue website) to act as our chatbot opening widget.  This design choice ensures visual consistency and immediate brand recognition.


![](/insights/ai-chatbot-circle.png)

We understand that sometimes users might need additional support or more in-depth information. That's why we've included convenient links within the chatbot interface. An email link for support is readily available, allowing users to reach out to our team directly if they need more personalized assistance. Additionally, we've incorporated a documentation link, providing easy access to more comprehensive resources for those who want to dive deeper into Blue's offerings.

The user experience is further enhanced by tasteful fade-in and fade-up animations when opening the chatbot window. These subtle animations add a touch of sophistication to the interface, making the interaction feel more dynamic and engaging. We've also implemented a typing indicator, a small but crucial feature that lets users know the chatbot is processing their query and crafting a response. This visual cue helps manage user expectations and maintains a sense of active communication.

<video autoplay loop muted playsinline>
  <source src="/videos/ai-chatbot-animation.mp4" type="video/mp4">
</video>


Recognizing that some conversations might require more screen real estate, we've added the ability to open the conversation in a larger window. This feature is particularly useful for longer exchanges or when reviewing detailed information, giving users the flexibility to adapt the chatbot to their needs.

Behind the scenes, we've implemented some intelligent processing to optimize the chatbot's responses. Our system automatically parses the AI's replies to remove references to our internal documents, ensuring that the information presented is clean, relevant, and focused solely on addressing the user's query.
To enhance readability and allow for more nuanced communication, we've incorporated markdown support using the 'marked' library. This feature enables our AI to provide richly formatted text, including bold and italic emphasis, structured lists, and even code snippets when necessary. It's akin to receiving a well-formatted, tailored mini-document in response to your questions.

Last but certainly not least, we've prioritized security in our implementation. Using the DOMPurify library, we sanitize the HTML generated from markdown parsing. This crucial step ensures that any potentially harmful scripts or code are stripped out before the content is displayed to you. It's our way of guaranteeing that the helpful information you receive is not only informative but also safe to consume.


## Future Developments

So this is just the start, we've got some exciting things on the roadmap for this feature. 

One of our upcoming features is the ability to stream responses in real-time. Soon, you'll see the chatbot's replies appear character by character, making conversations feel more natural and dynamic. It's like watching the AI think, creating a more engaging and interactive experience that keeps you in the loop every step of the way.

For our valued Blue users, we're working on personalization. The chatbot will recognize when you're logged in, tailoring its responses based on your account information, usage history, and preferences. Imagine a chatbot that not only answers your questions but understands your specific context within the Blue ecosystem, providing more relevant and personalized assistance.

We understand that you might be working on multiple projects or have various queries. That's why we're developing the ability to maintain several distinct conversation threads with our chatbot. This feature will allow you to switch between different topics seamlessly, without losing context – just like having multiple tabs open in your browser.

To make your interactions even more productive, we're creating a feature that will offer suggested follow-up questions based on your current conversation. This will help you explore topics more deeply and discover related information you might not have thought to ask about, making each chat session more comprehensive and valuable.

We're also excited about creating a suite of specialized AI assistants, each tailored for specific needs. Whether you're looking to answer pre-sales questions, set up a new project, or troubleshoot advanced features, you'll be able to choose the assistant that best fits your current needs. It's like having a team of Blue experts at your fingertips, each specializing in different aspects of our platform.

Lastly, we're working on allowing you to upload screenshots directly to the chat. The AI will analyze the image and provide explanations or troubleshooting steps based on what it sees. This feature will make it easier than ever to get help with specific issues you encounter while using Blue, bridging the gap between visual information and textual assistance.

## Conclusion

We hope this deep dive into our AI chatbot development process has provided some valuable insights into our product development thinking at Blue. Our journey from identifying the need for a chatbot to building our own solution showcases how we approach decision-making and innovation.

![](/insights/ai-chatbot-modal.png)

At Blue, we carefully weigh the options of building versus buying, always with an eye on what will best serve our users and align with our long-term goals. In this case, we identified a significant gap in the market for a cost-effective yet visually appealing chatbot that could meet our specific needs. While we generally advocate for leveraging existing solutions rather than reinventing the wheel, sometimes the best path forward is to create something tailored to your unique requirements.

Our decision to build our own chatbot wasn't taken lightly. It was the result of thorough market research, a clear understanding of our needs, and a commitment to providing the best possible experience for our users. By developing in-house, we were able to create a solution that not only meets our current needs but also lays the groundwork for future enhancements and integrations.

This project exemplifies our approach at Blue: we're not afraid to roll up our sleeves and build something from scratch when it's the right choice for our product and our users. It's this willingness to go the extra mile that allows us to deliver innovative solutions that truly meet the needs of our customers.
We're excited about the future of our AI chatbot and the value it will bring to both potential and existing Blue users. As we continue to refine and expand its capabilities, we remain committed to pushing the boundaries of what's possible in project management and customer interaction.

Thank you for joining us on this journey through our development process. We hope it's given you a glimpse into the thoughtful, user-centric approach we take with every aspect of Blue. Stay tuned for more updates as we continue to evolve and enhance our platform to better serve you.

If you're interested, you can find the link to source code for this project here:

- **[ChatWidget](https://gitlab.com/bloohq/blue-website/-/blob/main/components/ChatWidget.vue)**: This is a Vue Component that powers the chat widget itself. 
- **[Chat API](https://gitlab.com/bloohq/blue-website/-/blob/main/server/api/chat.post.ts)**: This is the middleware that works in-between the chat component and the OpenAI Assistants API