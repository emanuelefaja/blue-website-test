---
title: Warum wir unseren eigenen AI-Dokumentations-Chatbot entwickelt haben
category: "Product Updates"
description: Wir haben unseren eigenen Dokumentations-AI-Chatbot entwickelt, der auf der Dokumentation der Blue-Plattform trainiert ist.
date: 2024-07-09
---

Bei Blue sind wir ständig auf der Suche nach Möglichkeiten, das Leben unserer Kunden zu erleichtern. Wir haben [ausführliche Dokumentationen zu jeder Funktion](https://documentation.blue.cc), [YouTube-Videos](https://www.youtube.com/@HeyBlueTeam), [Tipps & Tricks](/insights/tips-tricks) und [verschiedene Support-Kanäle](/support).

Wir haben die Entwicklung von KI (Künstliche Intelligenz) genau im Auge behalten, da wir sehr an [Projektmanagement-Automatisierungen](/platform/features/automations) interessiert sind. Wir haben auch Funktionen wie [AI Auto-Kategorisierung](/insights/ai-auto-categorization) und [AI-Zusammenfassungen](/insights/ai-content-summarization) veröffentlicht, um die Arbeit für unsere tausenden von Kunden zu erleichtern.

Eines ist klar: KI ist hier, um zu bleiben, und sie wird einen unglaublichen Einfluss auf die meisten Branchen haben, wobei das Projektmanagement keine Ausnahme ist. Daher haben wir uns gefragt, wie wir KI weiter nutzen können, um den gesamten Lebenszyklus eines Kunden zu unterstützen, von der Entdeckung über den Vorverkauf, das Onboarding bis hin zu laufenden Fragen.

Die Antwort war ziemlich klar: **Wir benötigten einen KI-Chatbot, der auf unserer Dokumentation trainiert ist.**

Lassen Sie uns ehrlich sein: *Jede* Organisation sollte wahrscheinlich einen Chatbot haben. Sie sind großartige Möglichkeiten für Kunden, um sofort Antworten auf typische Fragen zu erhalten, ohne durch Seiten dichter Dokumentation oder Ihre Website blättern zu müssen. Die Bedeutung von Chatbots auf modernen Marketing-Websites kann nicht genug betont werden.

![](/insights/ai-chatbot-regular.png)

Für Softwareunternehmen sollte man die Marketing-Website nicht als separates "Ding" betrachten – sie *ist* Teil Ihres Produkts. Das liegt daran, dass sie in den typischen Lebenszyklus des Kunden passt:

- **Bewusstsein** (Entdeckung): Hier stoßen potenzielle Kunden zum ersten Mal auf Ihr großartiges Produkt. Ihr Chatbot kann ihr freundlicher Führer sein, der sie sofort auf wichtige Funktionen und Vorteile hinweist.
- **Überlegung** (Bildung): Jetzt sind sie neugierig und möchten mehr erfahren. Ihr Chatbot wird zu ihrem persönlichen Tutor und liefert Informationen, die auf ihre spezifischen Bedürfnisse und Fragen zugeschnitten sind.
- **Kauf/Umwandlung**: Dies ist der entscheidende Moment - wenn ein Interessent sich entscheidet, den Sprung zu wagen und Kunde zu werden. Ihr Chatbot kann letzte Hürden überwinden, Antworten auf die Fragen "gerade bevor ich kaufe" geben und vielleicht sogar ein tolles Angebot machen, um den Deal abzuschließen.
- **Onboarding**: Sie haben gekauft, und jetzt? Ihr Chatbot verwandelt sich in einen hilfreichen Begleiter, der neuen Benutzern beim Setup hilft, ihnen die Grundlagen zeigt und sicherstellt, dass sie sich in der Wunderwelt Ihres Produkts nicht verloren fühlen.
- **Bindung**: Kunden glücklich zu halten, ist das A und O. Ihr Chatbot ist rund um die Uhr verfügbar, bereit, Probleme zu beheben, Tipps und Tricks anzubieten und sicherzustellen, dass Ihre Kunden die Wertschätzung spüren.
- **Expansion**: Zeit, das Niveau zu erhöhen! Ihr Chatbot kann subtil neue Funktionen, Upsells oder Cross-Sells vorschlagen, die mit der Nutzung Ihres Produkts durch den Kunden übereinstimmen. Es ist, als hätte man einen wirklich intelligenten, nicht aufdringlichen Verkäufer immer in Bereitschaft.
- **Befürwortung**: Zufriedene Kunden werden zu Ihren größten Unterstützern. Ihr Chatbot kann zufriedene Benutzer ermutigen, das Wort zu verbreiten, Bewertungen abzugeben oder an Empfehlungsprogrammen teilzunehmen. Es ist, als hätte man eine Werbemaschine direkt in Ihr Produkt integriert!

## Build vs Buy Entscheidung

Nachdem wir uns entschieden hatten, einen KI-Chatbot zu implementieren, war die nächste große Frage: selbst entwickeln oder kaufen? Als kleines Team, das sich auf unser Kernprodukt konzentriert, ziehen wir im Allgemeinen "as-a-service"-Lösungen oder beliebte Open-Source-Plattformen vor. Schließlich sind wir nicht im Geschäft, das Rad für jeden Teil unseres Tech-Stacks neu zu erfinden.
Also krempelten wir die Ärmel hoch und tauchten in den Markt ein, auf der Suche nach sowohl kostenpflichtigen als auch Open-Source-KI-Chatbot-Lösungen.

Unsere Anforderungen waren einfach, aber nicht verhandelbar:

- **Unbranded Experience**: Dieser Chatbot ist nicht nur ein nettes Widget; er wird auf unserer Marketing-Website und schließlich in unserem Produkt eingesetzt. Wir sind nicht daran interessiert, die Marke eines anderen in unserem eigenen digitalen Raum zu bewerben.
- **Großartige UX**: Für viele potenzielle Kunden könnte dieser Chatbot ihr erster Kontakt mit Blue sein. Er bestimmt den Ton für ihre Wahrnehmung unseres Unternehmens. Lassen Sie uns ehrlich sein: Wenn wir keinen ordentlichen Chatbot auf unserer Website hinbekommen, wie können wir dann erwarten, dass Kunden uns mit ihren geschäftskritischen Projekten und Prozessen vertrauen?
- **Angemessene Kosten**: Mit einer großen Benutzerbasis und Plänen, den Chatbot in unser Kernprodukt zu integrieren, benötigten wir eine Lösung, die nicht die Bank sprengt, wenn die Nutzung zunimmt. Idealerweise wollten wir eine **BYOK (Bring Your Own Key) Option**. Dies würde es uns ermöglichen, unseren eigenen OpenAI- oder anderen KI-Dienstschlüssel zu verwenden und nur die direkten variablen Kosten zu zahlen, anstatt einen Aufschlag an einen Drittanbieter zu zahlen, der die Modelle nicht tatsächlich betreibt.
- **OpenAI Assistants API Kompatibel**: Wenn wir uns für eine Open-Source-Software entscheiden würden, wollten wir nicht die Mühe haben, eine Pipeline für die Dokumentenaufnahme, Indizierung, Vektordatenbanken und all das zu verwalten. Wir wollten die [OpenAI Assistants API](https://platform.openai.com/docs/assistants/overview) verwenden, die die gesamte Komplexität hinter einer API abstrahiert. Ehrlich gesagt – das ist wirklich gut gemacht.
- **Skalierbarkeit**: Wir möchten diesen Chatbot an mehreren Orten haben, mit potenziell zehntausenden von Benutzern pro Jahr. Wir erwarten eine erhebliche Nutzung und möchten nicht an eine Lösung gebunden sein, die nicht mit unseren Bedürfnissen skalieren kann.

## Kommerzielle KI-Chatbots

Die von uns überprüften hatten tendenziell eine bessere UX als Open-Source-Lösungen – wie es leider oft der Fall ist. Es gibt wahrscheinlich eines Tages eine separate Diskussion darüber, *warum* viele Open-Source-Lösungen die Bedeutung der UX ignorieren oder unterspielen.

Wir stellen hier eine Liste zur Verfügung, falls Sie nach soliden kommerziellen Angeboten suchen:

- **[Chatbase](https://chatbase.co):** Chatbase ermöglicht es Ihnen, einen benutzerdefinierten KI-Chatbot zu erstellen, der auf Ihrer Wissensdatenbank trainiert ist, und ihn zu Ihrer Website hinzuzufügen oder über ihre API zu interagieren. Es bietet Funktionen wie vertrauenswürdige Antworten, Lead-Generierung, erweiterte Analysen und die Möglichkeit, sich mit mehreren Datenquellen zu verbinden. Für uns fühlte sich dies wie eines der am besten ausgearbeiteten kommerziellen Angebote an.
- **[DocsBot AI](https://docsbot.ai/):** DocsBot AI erstellt benutzerdefinierte ChatGPT-Bots, die auf Ihrer Dokumentation und Ihren Inhalten für Support, Vorverkauf, Forschung und mehr trainiert sind. Es bietet einbettbare Widgets, um den Chatbot einfach zu Ihrer Website hinzuzufügen, die Möglichkeit, Support-Tickets automatisch zu beantworten, und eine leistungsstarke API für die Integration.
- **[CustomGPT.ai](https://customgpt.ai):** CustomGPT.ai schafft ein persönliches Chatbot-Erlebnis, indem es Ihre Geschäftsdaten, einschließlich Website-Inhalte, Helpdesk, Wissensdatenbanken, Dokumente und mehr, aufnimmt. Es ermöglicht Interessenten, Fragen zu stellen und sofortige Antworten basierend auf Ihren Inhalten zu erhalten, ohne suchen zu müssen. Interessanterweise [behaupten sie auch, OpenAI im RAG (Retrieval Augmented Generation) zu schlagen!](https://customgpt.ai/customgpt-beats-open-ai-in-rag-benchmark/)
- **[OpenAssistantGPT](https://www.openassistantgpt.io/)**: Dies ist ein interessantes kommerzielles Angebot, da es *auch* Open-Source-Software ist. Es scheint sich noch in einer frühen Phase zu befinden, und die Preisgestaltung fühlte sich nicht realistisch an (27 $/Monat für unbegrenzte Nachrichten wird kommerziell nie funktionieren).

Wir haben auch [InterCom Fin](https://www.intercom.com/fin) betrachtet, das Teil ihrer Kundenservice-Software ist. Dies hätte bedeutet, von [HelpScout](https://wwww.helpscout.com) wegzuwechseln, das wir seit der Gründung von Blue verwenden. Das hätte möglich sein können, aber InterCom Fin hat einige verrückte Preise, die es einfach aus der Überlegung ausschlossen.

Und das ist tatsächlich das Problem mit vielen der kommerziellen Angebote. InterCom Fin berechnet 0,99 $ pro bearbeitetem Kundenanfrage, und ChatBase berechnet 399 $/Monat für 40.000 Nachrichten. Das sind fast 5.000 $ pro Jahr für ein einfaches Chat-Widget.

Angesichts der Tatsache, dass die Preise für KI-Inferenz wie verrückt fallen. OpenAI hat seine Preise ziemlich drastisch gesenkt:

- Der ursprüngliche GPT-4 (8k Kontext) war mit 0,03 $ pro 1K Eingabetokens bepreist.
- Der GPT-4 Turbo (128k Kontext) war mit 0,01 $ pro 1K Eingabetokens bepreist, was einer 50%igen Reduzierung gegenüber dem ursprünglichen GPT-4 entspricht.
- Das GPT-4o-Modell ist mit 0,005 $ pro 1K Tokens bepreist, was eine weitere 50%ige Reduzierung gegenüber der Preisgestaltung des GPT-4 Turbo darstellt.

Das ist eine Reduzierung der Kosten um 83%, und wir erwarten nicht, dass dies stagnieren wird.

Angesichts der Tatsache, dass wir nach einer skalierbaren Lösung suchten, die von zehntausenden von Benutzern pro Jahr mit einer erheblichen Anzahl von Nachrichten genutzt wird, macht es Sinn, direkt zur Quelle zu gehen und die API-Kosten direkt zu bezahlen, anstatt eine kommerzielle Version zu verwenden, die die Kosten aufschlägt.

## Open Source KI-Chatbots

Wie bereits erwähnt, waren die Open-Source-Optionen, die wir überprüft haben, in Bezug auf die Anforderung "Großartige UX" größtenteils enttäuschend.

Wir haben uns angesehen:

- **[Deepchat](https://deepchat.dev/)**: Dies ist ein frameworkunabhängiges Chat-Komponente für KI-Dienste, die sich mit verschiedenen KI-APIs, einschließlich OpenAI, verbindet. Es hat auch die Fähigkeit für Benutzer, ein KI-Modell herunterzuladen, das direkt im Browser läuft. Wir haben damit experimentiert und eine funktionierende Version erhalten, aber die implementierte OpenAI Assistants API fühlte sich ziemlich fehlerhaft an, mit mehreren Problemen. Dennoch ist dies ein sehr vielversprechendes Projekt, und ihr Playground ist wirklich gut gemacht.
- **[OpenAssistantGPT](https://www.openassistantgpt.io/)**: Wenn wir dies erneut aus einer Open-Source-Perspektive betrachten, müssten wir eine Menge Infrastruktur aufbauen, was wir nicht tun wollten, da wir so viel wie möglich auf die OpenAI Assistants API angewiesen sein wollten.

## Erstellung unseres eigenen ChatBots

Und so, ohne etwas zu finden, das all unseren Anforderungen entsprach, haben wir beschlossen, unseren eigenen KI-Chatbot zu entwickeln, der mit der OpenAI Assistants API kommunizieren kann. Dies stellte sich letztendlich als relativ schmerzlos heraus!

Unsere Website verwendet [Nuxt3](https://nuxt.com), [Vue3](https://vuejs.org/) (was dasselbe Framework wie die Blue-Plattform ist) und [TailwindUI](https://tailwindui.com/).

Der erste Schritt bestand darin, eine API (Application Programming Interface) in Nuxt3 zu erstellen, die mit der OpenAI Assistants API "kommunizieren" kann. Dies war notwendig, da wir nicht alles im Frontend machen wollten, da dies unseren OpenAI API-Schlüssel der Welt aussetzen würde, mit dem Potenzial für Missbrauch.

Unsere Backend-API fungiert als sicherer Mittelsmann zwischen dem Browser des Benutzers und OpenAI. Hier ist, was sie tut:

- **Konversationsmanagement:** Sie erstellt und verwaltet "Threads" für jede Konversation. Denken Sie an einen Thread als eine einzigartige Chatsitzung, die sich an alles erinnert, was Sie gesagt haben.
- **Nachrichtenverarbeitung:** Wenn Sie eine Nachricht senden, fügt unsere API sie dem richtigen Thread hinzu und bittet den Assistenten von OpenAI, eine Antwort zu formulieren.
- **Intelligentes Warten:** Anstatt Sie auf einem Ladebildschirm starren zu lassen, fragt unsere API jede Sekunde bei OpenAI nach, ob Ihre Antwort bereit ist. Es ist, als hätte man einen Kellner, der Ihre Bestellung im Auge behält, ohne den Koch alle zwei Sekunden zu belästigen.
- **Sicherheit zuerst:** Indem wir all dies auf dem Server abwickeln, halten wir Ihre Daten und unsere API-Schlüssel sicher und geschützt.

Dann gab es das Frontend und die Benutzererfahrung. Wie bereits erwähnt, war dies *entscheidend* wichtig, da wir keine zweite Chance haben, einen ersten Eindruck zu hinterlassen!

Bei der Gestaltung unseres Chatbots haben wir akribisch auf die Benutzererfahrung geachtet und sichergestellt, dass jede Interaktion reibungslos, intuitiv und ein Spiegelbild von Blues Engagement für Qualität ist. Die Chatbot-Oberfläche beginnt mit einem einfachen, eleganten blauen Kreis, wobei wir [HeroIcons für unsere Icons](https://heroicons.com/) verwenden (die wir auf der gesamten Blue-Website verwenden), um als unser Chatbot-Öffnungs-Widget zu fungieren. Diese Designwahl sorgt für visuelle Konsistenz und sofortige Markenwiedererkennung.

![](/insights/ai-chatbot-circle.png)

Wir verstehen, dass Benutzer manchmal zusätzliche Unterstützung oder detailliertere Informationen benötigen. Deshalb haben wir praktische Links in die Chatbot-Oberfläche integriert. Ein E-Mail-Link für den Support ist leicht verfügbar, sodass Benutzer unser Team direkt kontaktieren können, wenn sie personalisierte Unterstützung benötigen. Darüber hinaus haben wir einen Dokumentationslink integriert, der einen einfachen Zugang zu umfassenderen Ressourcen bietet, für diejenigen, die tiefer in Blues Angebote eintauchen möchten.

Die Benutzererfahrung wird durch geschmackvolle Fade-In- und Fade-Up-Animationen beim Öffnen des Chatbot-Fensters weiter verbessert. Diese subtilen Animationen verleihen der Oberfläche einen Hauch von Raffinesse und machen die Interaktion dynamischer und ansprechender. Wir haben auch einen Tipp-Indikator implementiert, ein kleines, aber entscheidendes Feature, das den Benutzern mitteilt, dass der Chatbot ihre Anfrage verarbeitet und eine Antwort formuliert. Dieses visuelle Signal hilft, die Erwartungen der Benutzer zu steuern und ein Gefühl aktiver Kommunikation aufrechtzuerhalten.

<video autoplay loop muted playsinline>
  <source src="/videos/ai-chatbot-animation.mp4" type="video/mp4">
</video>

Da wir erkannt haben, dass einige Gespräche möglicherweise mehr Bildschirmfläche erfordern, haben wir die Möglichkeit hinzugefügt, die Konversation in einem größeren Fenster zu öffnen. Diese Funktion ist besonders nützlich für längere Gespräche oder beim Überprüfen detaillierter Informationen, da sie den Benutzern die Flexibilität gibt, den Chatbot an ihre Bedürfnisse anzupassen.

Hinter den Kulissen haben wir einige intelligente Verarbeitung implementiert, um die Antworten des Chatbots zu optimieren. Unser System analysiert automatisch die Antworten der KI, um Verweise auf unsere internen Dokumente zu entfernen, sodass die präsentierten Informationen sauber, relevant und ausschließlich darauf ausgerichtet sind, die Anfrage des Benutzers zu beantworten. Um die Lesbarkeit zu verbessern und nuanciertere Kommunikation zu ermöglichen, haben wir die Unterstützung von Markdown mit der 'marked'-Bibliothek integriert. Diese Funktion ermöglicht es unserer KI, reichhaltig formatierte Texte bereitzustellen, einschließlich fett und kursiv hervorgehobener Texte, strukturierter Listen und sogar Code-Snippets, wenn nötig. Es ist, als würde man eine gut formatierte, maßgeschneiderte Mini-Dokumentation als Antwort auf Ihre Fragen erhalten.

Zu guter Letzt haben wir die Sicherheit in unserer Implementierung priorisiert. Mit der DOMPurify-Bibliothek reinigen wir das HTML, das aus der Markdown-Analyse generiert wird. Dieser entscheidende Schritt stellt sicher, dass potenziell schädliche Skripte oder Codes entfernt werden, bevor der Inhalt Ihnen angezeigt wird. Es ist unsere Art zu garantieren, dass die hilfreichen Informationen, die Sie erhalten, nicht nur informativ, sondern auch sicher sind.

## Zukünftige Entwicklungen

Das ist also nur der Anfang, wir haben einige aufregende Dinge auf der Roadmap für dieses Feature.

Eine unserer kommenden Funktionen ist die Möglichkeit, Antworten in Echtzeit zu streamen. Bald werden Sie sehen, wie die Antworten des Chatbots Zeichen für Zeichen erscheinen, was Gespräche natürlicher und dynamischer macht. Es ist, als würde man der KI beim Denken zusehen und eine ansprechendere und interaktive Erfahrung schaffen, die Sie in jedem Schritt auf dem Laufenden hält.

Für unsere geschätzten Blue-Nutzer arbeiten wir an der Personalisierung. Der Chatbot wird erkennen, wenn Sie angemeldet sind, und seine Antworten basierend auf Ihren Kontoinformationen, Nutzungshistorie und Vorlieben anpassen. Stellen Sie sich einen Chatbot vor, der nicht nur Ihre Fragen beantwortet, sondern auch Ihren spezifischen Kontext innerhalb des Blue-Ökosystems versteht und relevantere und personalisierte Unterstützung bietet.

Wir verstehen, dass Sie möglicherweise an mehreren Projekten arbeiten oder verschiedene Anfragen haben. Deshalb entwickeln wir die Möglichkeit, mehrere unterschiedliche Konversationsstränge mit unserem Chatbot zu führen. Diese Funktion ermöglicht es Ihnen, nahtlos zwischen verschiedenen Themen zu wechseln, ohne den Kontext zu verlieren – genau wie beim Öffnen mehrerer Tabs in Ihrem Browser.

Um Ihre Interaktionen noch produktiver zu gestalten, erstellen wir eine Funktion, die vorgeschlagene Folgefragen basierend auf Ihrem aktuellen Gespräch anbietet. Dies wird Ihnen helfen, Themen tiefer zu erkunden und verwandte Informationen zu entdecken, an die Sie vielleicht nicht gedacht haben, was jede Chatsitzung umfassender und wertvoller macht.

Wir freuen uns auch darauf, eine Suite von spezialisierten KI-Assistenten zu schaffen, die jeweils auf spezifische Bedürfnisse zugeschnitten sind. Egal, ob Sie Antworten auf Vorverkaufsfragen suchen, ein neues Projekt einrichten oder fortgeschrittene Funktionen beheben möchten, Sie werden in der Lage sein, den Assistenten auszuwählen, der am besten zu Ihren aktuellen Bedürfnissen passt. Es ist, als hätte man ein Team von Blue-Experten jederzeit zur Hand, die sich auf verschiedene Aspekte unserer Plattform spezialisiert haben.

Zuletzt arbeiten wir daran, Ihnen die Möglichkeit zu geben, Screenshots direkt in den Chat hochzuladen. Die KI wird das Bild analysieren und Erklärungen oder Schritte zur Fehlerbehebung basierend auf dem, was sie sieht, bereitstellen. Diese Funktion wird es einfacher denn je machen, Hilfe bei spezifischen Problemen zu erhalten, auf die Sie beim Verwenden von Blue stoßen, und die Kluft zwischen visuellen Informationen und textlicher Unterstützung zu überbrücken.

## Fazit

Wir hoffen, dass dieser tiefere Einblick in unseren Entwicklungsprozess des KI-Chatbots einige wertvolle Einblicke in unser Produktentwicklungsdenken bei Blue gegeben hat. Unser Weg von der Identifizierung des Bedarfs an einem Chatbot bis hin zur Entwicklung unserer eigenen Lösung zeigt, wie wir Entscheidungsfindung und Innovation angehen.

![](/insights/ai-chatbot-modal.png)

Bei Blue wägen wir sorgfältig die Optionen zwischen Eigenentwicklung und Kauf ab, immer mit dem Blick darauf, was unseren Nutzern am besten dient und mit unseren langfristigen Zielen übereinstimmt. In diesem Fall haben wir eine signifikante Lücke im Markt für einen kostengünstigen, aber visuell ansprechenden Chatbot identifiziert, der unseren spezifischen Bedürfnissen gerecht werden kann. Während wir im Allgemeinen dafür plädieren, bestehende Lösungen zu nutzen, anstatt das Rad neu zu erfinden, ist manchmal der beste Weg nach vorne, etwas zu schaffen, das auf Ihre einzigartigen Anforderungen zugeschnitten ist.

Unsere Entscheidung, unseren eigenen Chatbot zu entwickeln, wurde nicht leichtfertig getroffen. Sie war das Ergebnis gründlicher Marktforschung, eines klaren Verständnisses unserer Bedürfnisse und eines Engagements, die bestmögliche Erfahrung für unsere Nutzer zu bieten. Durch die Entwicklung im eigenen Haus konnten wir eine Lösung schaffen, die nicht nur unseren aktuellen Bedürfnissen entspricht, sondern auch die Grundlage für zukünftige Verbesserungen und Integrationen legt.

Dieses Projekt exemplifiziert unseren Ansatz bei Blue: Wir scheuen uns nicht, die Ärmel hochzukrempeln und etwas von Grund auf neu zu bauen, wenn es die richtige Wahl für unser Produkt und unsere Nutzer ist. Diese Bereitschaft, einen Schritt weiter zu gehen, ermöglicht es uns, innovative Lösungen zu liefern, die wirklich den Bedürfnissen unserer Kunden entsprechen. 
Wir sind gespannt auf die Zukunft unseres KI-Chatbots und den Wert, den er sowohl potenziellen als auch bestehenden Blue-Nutzern bringen wird. Während wir weiterhin seine Fähigkeiten verfeinern und erweitern, bleiben wir dem Ziel verpflichtet, die Grenzen dessen, was im Projektmanagement und in der Kundeninteraktion möglich ist, zu verschieben.

Vielen Dank, dass Sie uns auf dieser Reise durch unseren Entwicklungsprozess begleitet haben. Wir hoffen, dass es Ihnen einen Einblick in den durchdachten, benutzerzentrierten Ansatz gegeben hat, den wir bei jedem Aspekt von Blue verfolgen. Bleiben Sie dran für weitere Updates, während wir weiterhin unsere Plattform weiterentwickeln und verbessern, um Ihnen besser zu dienen.

Wenn Sie interessiert sind, finden Sie hier den Link zum Quellcode für dieses Projekt:

- **[ChatWidget](https://gitlab.com/bloohq/blue-website/-/blob/main/components/ChatWidget.vue)**: Dies ist eine Vue-Komponente, die das Chat-Widget selbst antreibt.
- **[Chat API](https://gitlab.com/bloohq/blue-website/-/blob/main/server/api/chat.post.ts)**: Dies ist die Middleware, die zwischen der Chat-Komponente und der OpenAI Assistants API arbeitet.