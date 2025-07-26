---
title: Erstellung der benutzerdefinierten Berechtigungsengine von Blue
category: "Engineering"
description: Werfen Sie einen Blick hinter die Kulissen des Blue-Engineering-Teams, während es erklärt, wie eine KI-gestützte Auto-Kategorisierungs- und Tagging-Funktion entwickelt wurde.
date: 2024-07-25
---

Effektives Projekt- und Prozessmanagement ist entscheidend für Organisationen jeder Größe.

Bei Blue haben wir es uns zur Aufgabe gemacht, die Arbeit der Welt zu organisieren, indem wir die beste Projektmanagement-Plattform auf dem Planeten entwickeln – einfach, leistungsstark, flexibel und für alle erschwinglich.

Das bedeutet, dass unsere Plattform sich an die einzigartigen Bedürfnisse jedes Teams anpassen muss. Heute freuen wir uns, den Vorhang für eines unserer leistungsstärksten Features zu lüften: Benutzerdefinierte Berechtigungen.

Projektmanagement-Tools sind das Rückgrat moderner Arbeitsabläufe, in denen sensible Daten, wichtige Kommunikationen und strategische Pläne untergebracht sind. Daher ist die Fähigkeit, den Zugriff auf diese Informationen präzise zu steuern, nicht nur ein Luxus – sie ist eine Notwendigkeit.

<video autoplay loop muted playsinline>
  <source src="/videos/user-roles.mp4" type="video/mp4">
</video>

Benutzerdefinierte Berechtigungen spielen eine entscheidende Rolle in B2B-SaaS-Plattformen, insbesondere in Projektmanagement-Tools, wo das Gleichgewicht zwischen Zusammenarbeit und Sicherheit den Erfolg eines Projekts ausmachen kann.

Aber hier geht Blue einen anderen Weg: **Wir glauben, dass Funktionen auf Unternehmensniveau nicht nur für Budgets in Unternehmensgröße reserviert sein sollten.**

In einer Ära, in der KI kleinen Teams ermöglicht, in nie dagewesenem Maßstab zu operieren, warum sollten robuste Sicherheit und Anpassungsmöglichkeiten unerreichbar sein?

In diesem Blick hinter die Kulissen werden wir erkunden, wie wir unser Feature für benutzerdefinierte Berechtigungen entwickelt haben, die den Status quo der SaaS-Preisklassen in Frage stellen und leistungsstarke, flexible Sicherheitsoptionen für Unternehmen jeder Größe bereitstellen.

Egal, ob Sie ein Startup mit großen Träumen oder ein etabliertes Unternehmen sind, das seine Prozesse optimieren möchte, benutzerdefinierte Berechtigungen können neue Anwendungsfälle ermöglichen, von denen Sie nicht einmal wussten, dass sie möglich sind.

## Verständnis der benutzerdefinierten Benutzerberechtigungen

Bevor wir in unsere Reise zur Entwicklung benutzerdefinierter Berechtigungen für Blue eintauchen, lassen Sie uns einen Moment Zeit nehmen, um zu verstehen, was benutzerdefinierte Benutzerberechtigungen sind und warum sie im Projektmanagement-Software so wichtig sind.

Benutzerdefinierte Benutzerberechtigungen beziehen sich auf die Fähigkeit, Zugriffsrechte für einzelne Benutzer oder Gruppen innerhalb eines Softwaresystems anzupassen. Anstatt sich auf vordefinierte Rollen mit festen Berechtigungen zu verlassen, ermöglichen benutzerdefinierte Berechtigungen Administratoren, hochspezifische Zugriffsprofile zu erstellen, die perfekt mit der Struktur und den Arbeitsablaufbedürfnissen ihrer Organisation übereinstimmen.

Im Kontext von Projektmanagement-Software wie Blue umfassen benutzerdefinierte Berechtigungen:

* **Granulare Zugriffskontrolle**: Bestimmen, wer bestimmte Arten von Projektdaten anzeigen, bearbeiten oder löschen kann.
* **Funktionsbasierte Einschränkungen**: Aktivieren oder Deaktivieren bestimmter Funktionen für bestimmte Benutzer oder Teams.
* **Datensensitivitätsstufen**: Festlegen unterschiedlicher Zugriffslevels auf sensible Informationen innerhalb von Projekten.
* **Workflow-spezifische Berechtigungen**: Abstimmung der Benutzerfähigkeiten mit spezifischen Phasen oder Aspekten Ihres Projektablaufs.

Die Bedeutung benutzerdefinierter Berechtigungen im Projektmanagement kann nicht genug betont werden:

* **Erhöhte Sicherheit**: Durch die Bereitstellung von nur dem Zugriff, den Benutzer benötigen, verringern Sie das Risiko von Datenverletzungen oder unbefugten Änderungen.
* **Verbesserte Compliance**: Benutzerdefinierte Berechtigungen helfen Organisationen, branchenspezifische regulatorische Anforderungen zu erfüllen, indem sie den Datenzugriff steuern.
* **Optimierte Zusammenarbeit**: Teams können effizienter arbeiten, wenn jedes Mitglied das richtige Maß an Zugriff hat, um seine Rolle ohne unnötige Einschränkungen oder überwältigende Berechtigungen auszuführen.
* **Flexibilität für komplexe Organisationen**: Wenn Unternehmen wachsen und sich weiterentwickeln, ermöglichen benutzerdefinierte Berechtigungen der Software, sich an sich ändernde Organisationsstrukturen und Prozesse anzupassen.

## Zum JA kommen

[Wir haben zuvor geschrieben](/insights/value-proposition-blue), dass jede Funktion in Blue ein **hartes** JA sein muss, bevor wir entscheiden, sie zu entwickeln. Wir haben nicht die Luxus von Hunderten von Ingenieuren und können keine Zeit und kein Geld mit dem Bau von Dingen verschwenden, die die Kunden nicht benötigen.

Und so war der Weg zur Implementierung benutzerdefinierter Berechtigungen in Blue keine gerade Linie. Wie viele leistungsstarke Funktionen begann es mit einem klaren Bedarf unserer Benutzer und entwickelte sich durch sorgfältige Überlegungen und Planung.

Jahrelang hatten unsere Kunden nach mehr granularer Kontrolle über Benutzerberechtigungen gefragt. Als Organisationen jeder Größe begannen, zunehmend komplexe und sensible Projekte zu bearbeiten, wurden die Einschränkungen unserer standardmäßigen rollenbasierten Zugriffskontrolle offensichtlich.

Kleine Startups, die mit externen Kunden arbeiten, mittelständische Unternehmen mit komplexen Genehmigungsprozessen und große Unternehmen mit strengen Compliance-Anforderungen äußerten alle denselben Bedarf:

Mehr Flexibilität bei der Verwaltung des Benutzerzugriffs.

Trotz der klaren Nachfrage zögerten wir zunächst, in die Entwicklung benutzerdefinierter Berechtigungen einzutauchen.

Warum?

Wir verstanden die Komplexität, die damit verbunden ist!

Benutzerdefinierte Berechtigungen berühren jeden Teil eines Projektmanagementsystems, von der Benutzeroberfläche bis hin zur Datenbankstruktur. Wir wussten, dass die Implementierung dieser Funktion erhebliche Änderungen an unserer Kernarchitektur erfordern würde und sorgfältige Überlegungen zu den Auswirkungen auf die Leistung notwendig wären.

Als wir die Landschaft sondierten, stellten wir fest, dass nur sehr wenige unserer Wettbewerber versucht hatten, eine leistungsstarke benutzerdefinierte Berechtigungsengine wie die zu implementieren, die unsere Kunden anforderten. Diejenigen, die es taten, reservierten sie oft für ihre höchsten Unternehmenspläne.

Es wurde klar, warum: Der Entwicklungsaufwand ist erheblich, und die Einsätze sind hoch.

Die falsche Implementierung von benutzerdefinierten Berechtigungen könnte kritische Fehler oder Sicherheitsanfälligkeiten einführen, die das gesamte System gefährden könnten. Diese Erkenntnis verdeutlichte das Ausmaß der Herausforderung, die wir in Betracht zogen.

### Den Status Quo herausfordern

Als wir jedoch weiterhin wuchsen und uns weiterentwickelten, erreichten wir eine entscheidende Erkenntnis:

**Das traditionelle SaaS-Modell, leistungsstarke Funktionen für Unternehmenskunden zu reservieren, macht im heutigen Geschäftsumfeld keinen Sinn mehr.**

Im Jahr 2024, mit der Kraft von KI und fortschrittlichen Tools, können kleine Teams in einem Maßstab und einer Komplexität operieren, die mit viel größeren Organisationen konkurrieren. Ein Startup könnte sensible Kundendaten in mehreren Ländern bearbeiten. Eine kleine Marketingagentur könnte Dutzende von Kundenprojekten mit unterschiedlichen Vertraulichkeitsanforderungen jonglieren. Diese Unternehmen benötigen das gleiche Maß an Sicherheit und Anpassung wie *jede* große Unternehmung.

Wir fragten uns: Warum sollte die Größe der Belegschaft oder des Budgets eines Unternehmens darüber entscheiden, ob es seine Daten sicher und seine Prozesse effizient halten kann?

### Unternehmensniveau für alle

Diese Erkenntnis führte uns zu einer Kernphilosophie, die nun einen Großteil unserer Entwicklung bei Blue antreibt: Funktionen auf Unternehmensniveau sollten für Unternehmen jeder Größe zugänglich sein.

Wir glauben, dass:

- **Sicherheit kein Luxus sein sollte.** Jedes Unternehmen, unabhängig von seiner Größe, verdient die Werkzeuge, um seine Daten und Prozesse zu schützen.
- **Flexibilität Innovation antreibt.** Indem wir allen unseren Benutzern leistungsstarke Werkzeuge geben, ermöglichen wir es ihnen, Arbeitsabläufe und Systeme zu schaffen, die ihre Branchen voranbringen.
- **Wachstum keine Plattformänderungen erfordern sollte.** Wenn unsere Kunden wachsen, sollten ihre Werkzeuge nahtlos mit ihnen wachsen.

Mit dieser Denkweise beschlossen wir, die Herausforderung der benutzerdefinierten Berechtigungen direkt anzugehen, und verpflichteten uns, sie allen unseren Benutzern zur Verfügung zu stellen, nicht nur denen in höheren Plänen.

Diese Entscheidung setzte uns auf einen Weg sorgfältigen Designs, iterativer Entwicklung und kontinuierlichem Benutzerfeedback, das letztendlich zu dem Feature für benutzerdefinierte Berechtigungen führte, auf das wir heute stolz sind.

Im nächsten Abschnitt werden wir darauf eingehen, wie wir den Design- und Entwicklungsprozess angegangen sind, um dieses komplexe Feature zum Leben zu erwecken.

### Design und Entwicklung

Als wir uns entschieden, benutzerdefinierte Berechtigungen anzugehen, wurde uns schnell klar, dass wir es mit einer riesigen Aufgabe zu tun hatten.

Auf den ersten Blick mag "benutzerdefinierte Berechtigungen" einfach klingen, aber es ist ein täuschend komplexes Feature, das jeden Aspekt unseres Systems berührt.

Die Herausforderung war gewaltig: Wir mussten kaskadierende Berechtigungen implementieren, Bearbeitungen in Echtzeit ermöglichen, erhebliche Änderungen am Datenbankschema vornehmen und eine nahtlose Funktionalität in unserem gesamten Ökosystem – Web, Mac, Windows, iOS und Android-Apps sowie unserer API und Webhooks – sicherstellen.

Die Komplexität war genug, um selbst die erfahrensten Entwickler innehalten zu lassen.

Unser Ansatz konzentrierte sich auf zwei Grundprinzipien:

1. Die Funktion in handhabbare Versionen aufzuteilen
2. Inkrementelles Versenden zu akzeptieren.

Angesichts der Komplexität der benutzerdefinierten Berechtigungen stellten wir uns eine entscheidende Frage:

> Was wäre die einfachste mögliche erste Version dieses Features?

Dieser Ansatz entspricht dem agilen Prinzip, ein Minimum Viable Product (MVP) zu liefern und basierend auf Feedback iterativ zu arbeiten.

Unsere Antwort war erfrischend einfach:

1. Einführung eines Schalters zum Ausblenden des Projektaktivitätstags
2. Hinzufügen eines weiteren Schalters zum Ausblenden des Formularstags

**Das war's.**

Keine großen Effekte, keine komplexen Berechtigungsmatrizen – nur zwei einfache Ein/Aus-Schalter.

Obwohl es auf den ersten Blick unterwhelming erscheinen mag, bot dieser Ansatz mehrere bedeutende Vorteile:

* **Schnelle Implementierung**: Diese einfachen Schalter konnten schnell entwickelt und getestet werden, sodass wir eine grundlegende Version der benutzerdefinierten Berechtigungen schnell in die Hände der Benutzer bekommen konnten.
* **Klarer Benutzerwert**: Selbst mit nur diesen beiden Optionen boten wir greifbaren Wert. Einige Teams möchten möglicherweise den Aktivitätsfeed von Kunden ausblenden, während andere den Zugriff auf Formulare für bestimmte Benutzergruppen einschränken müssen.
* **Grundlage für Wachstum**: Dieser einfache Start legte den Grundstein für komplexere Berechtigungen. Er ermöglichte uns, die grundlegende Infrastruktur für benutzerdefinierte Berechtigungen einzurichten, ohne uns von Anfang an in der Komplexität zu verlieren.
* **Benutzerfeedback**: Durch die Veröffentlichung dieser einfachen Version konnten wir echtes Feedback darüber sammeln, wie Benutzer mit benutzerdefinierten Berechtigungen interagierten, was unsere zukünftige Entwicklung informierte.
* **Technisches Lernen**: Diese erste Implementierung gab unserem Entwicklungsteam praktische Erfahrungen in der Modifizierung von Berechtigungen über unsere Plattform hinweg und bereitete uns auf komplexere Iterationen vor.

Und wissen Sie, es ist tatsächlich ziemlich demütigend, eine große Vision für etwas zu haben und dann etwas zu liefern, das nur einen so kleinen Prozentsatz dieser Vision ausmacht.

Nachdem wir diese ersten beiden Schalter veröffentlicht hatten, beschlossen wir, etwas Sophistizierteres anzugehen. Wir entschieden uns für zwei neue benutzerdefinierte Benutzerrollenberechtigungen.

Die erste war die Möglichkeit, Benutzer auf die Ansicht von Datensätzen zu beschränken, die speziell ihnen zugewiesen wurden. Dies ist sehr nützlich, wenn Sie einen Kunden in einem Projekt haben und nur möchten, dass er Datensätze sieht, die speziell ihm zugewiesen sind, anstatt alles, woran Sie für ihn arbeiten.

Die zweite war eine Option für Projektadministratoren, Benutzergruppen daran zu hindern, andere Benutzer einzuladen. Dies ist gut, wenn Sie ein sensibles Projekt haben, das Sie sicherstellen möchten, dass es auf einer "Need to See"-Basis bleibt.

Nachdem wir dies veröffentlicht hatten, gewannen wir mehr Vertrauen und in unserer dritten Version nahmen wir uns die Spaltenberechtigungen vor, was bedeutet, dass wir entscheiden konnten, welche benutzerdefinierten Felder eine bestimmte Benutzergruppe anzeigen oder bearbeiten kann.

Dies ist äußerst leistungsstark. Stellen Sie sich vor, Sie haben ein CRM-Projekt, und Sie haben Daten, die nicht nur mit den Beträgen, die der Kunde zahlen wird, sondern auch mit Ihren Kosten und Gewinnmargen verbunden sind. Sie möchten möglicherweise nicht, dass Ihre Kostenfelder und das Projektmargenformelfeld für Junior-Mitarbeiter sichtbar sind, und benutzerdefinierte Berechtigungen ermöglichen es Ihnen, diese Felder zu sperren, sodass sie nicht angezeigt werden.

Als Nächstes gingen wir dazu über, listenbasierte Berechtigungen zu erstellen, bei denen Projektadministratoren entscheiden können, ob eine Benutzergruppe eine bestimmte Liste anzeigen, bearbeiten und löschen kann. Wenn sie eine Liste ausblenden, werden auch alle Datensätze in dieser Liste ausgeblendet, was großartig ist, da es bedeutet, dass Sie bestimmte Teile Ihres Prozesses vor Ihren Teammitgliedern oder Kunden verbergen können.

Das ist das Endergebnis:

<video autoplay loop muted playsinline>
  <source src="/videos/custom-user-roles.mp4" type="video/mp4">
</video>

## Technische Überlegungen

Im Herzen der technischen Architektur von Blue liegt GraphQL, eine entscheidende Wahl, die unsere Fähigkeit, komplexe Funktionen wie benutzerdefinierte Berechtigungen zu implementieren, erheblich beeinflusst hat. Aber bevor wir in die Einzelheiten eintauchen, lassen Sie uns einen Schritt zurücktreten und verstehen, was GraphQL ist und wie es sich von dem traditionelleren REST-API-Ansatz unterscheidet.
GraphQL vs REST API: Eine zugängliche Erklärung

Stellen Sie sich vor, Sie sind in einem Restaurant. Mit einer REST-API ist es wie das Bestellen von einer festen Speisekarte. Sie bitten um ein bestimmtes Gericht (Endpunkt), und Sie erhalten alles, was dazugehört, egal ob Sie alles wollen oder nicht. Wenn Sie Ihre Mahlzeit anpassen möchten, müssen Sie möglicherweise mehrere Bestellungen (API-Aufrufe) aufgeben oder um ein speziell zubereitetes Gericht (benutzerdefinierter Endpunkt) bitten.

GraphQL hingegen ist wie ein Gespräch mit einem Koch, der alles zubereiten kann. Sie sagen dem Koch genau, welche Zutaten Sie möchten (Datenfelder) und in welchen Mengen. Der Koch bereitet dann ein Gericht zu, das genau das ist, was Sie angefordert haben – nicht mehr, nicht weniger. Genau das macht GraphQL – es ermöglicht dem Client, genau die Daten anzufordern, die er benötigt, und der Server liefert nur das.

### Ein wichtiges Mittagessen

Etwa sechs Wochen nach der ursprünglichen Entwicklung von Blue gingen unser leitender Ingenieur und CEO zum Mittagessen.

Das Thema der Diskussion?

Ob wir von REST-APIs zu GraphQL wechseln sollten. Dies war keine Entscheidung, die leichtfertig getroffen werden sollte – die Annahme von GraphQL würde bedeuten, sechs Wochen ursprünglicher Arbeit abzulehnen.

Auf dem Rückweg ins Büro stellte der CEO dem leitenden Ingenieur eine entscheidende Frage: "Würden wir es in fünf Jahren bereuen, dies nicht getan zu haben?"

Die Antwort wurde klar: GraphQL war der Weg nach vorne.

Wir erkannten das Potenzial dieser Technologie frühzeitig und sahen, wie sie unsere Vision für eine flexible, leistungsstarke Projektmanagement-Plattform unterstützen könnte.

Unser Weitblick bei der Annahme von GraphQL zahlte sich aus, als es darum ging, benutzerdefinierte Berechtigungen zu implementieren. Mit einer REST-API hätten wir für jede mögliche Konfiguration von benutzerdefinierten Berechtigungen einen anderen Endpunkt benötigt – ein Ansatz, der schnell unhandlich und schwer zu warten geworden wäre.

GraphQL hingegen ermöglicht es uns, benutzerdefinierte Berechtigungen dynamisch zu handhaben. So funktioniert es:

- **Echtzeit-Berechtigungsprüfungen**: Wenn ein Client eine Anfrage stellt, kann unser GraphQL-Server die Berechtigungen des Benutzers direkt aus unserer Datenbank überprüfen.
- **Präzise Datenabfrage**: Basierend auf diesen Berechtigungen gibt GraphQL nur die angeforderten Daten zurück, die innerhalb der Zugriffsrechte des Benutzers liegen.
- **Flexible Abfragen**: Wenn sich die Berechtigungen ändern, müssen wir keine neuen Endpunkte erstellen oder bestehende ändern. Die gleiche GraphQL-Abfrage kann sich an unterschiedliche Berechtigungseinstellungen anpassen.
- **Effizientes Datenabrufen**: GraphQL ermöglicht es Clients, genau das anzufordern, was sie benötigen. Das bedeutet, dass wir keine Daten überholen, was potenziell Informationen offenlegen könnte, auf die der Benutzer keinen Zugriff haben sollte.

Diese Flexibilität ist entscheidend für ein Feature so komplex wie benutzerdefinierte Berechtigungen. Sie ermöglicht es uns, granularen Zugriff *ohne* Einbußen bei der Leistung oder Wartbarkeit anzubieten.

## Herausforderungen

Die Implementierung benutzerdefinierter Berechtigungen in Blue brachte ihre eigenen Herausforderungen mit sich, die uns dazu drängten, innovativ zu sein und unseren Ansatz zu verfeinern. Die Leistungsoptimierung wurde schnell zu einem kritischen Anliegen. Als wir mehr granulare Berechtigungsprüfungen hinzufügten, riskieren wir, unser System zu verlangsamen, insbesondere bei großen Projekten mit vielen Benutzern und komplexen Berechtigungseinstellungen. Um dem entgegenzuwirken, implementierten wir eine mehrstufige Caching-Strategie, optimierten unsere Datenbankabfragen und nutzten die Fähigkeit von GraphQL, nur die notwendigen Daten anzufordern. Dieser Ansatz ermöglichte es uns, schnelle Antwortzeiten aufrechtzuerhalten, selbst wenn Projekte skalierten und die Berechtigungskomplexität zunahm.

Die Benutzeroberfläche für benutzerdefinierte Berechtigungen stellte eine weitere bedeutende Hürde dar. Wir mussten die Benutzeroberfläche intuitiv und verwaltbar für Administratoren gestalten, selbst als wir mehr Optionen hinzufügten und die Komplexität des Systems erhöhten.

Unsere Lösung umfasste mehrere Runden von Benutzertests und iterativem Design.

Wir führten eine visuelle Berechtigungsmatrix ein, die es Administratoren ermöglichte, Berechtigungen schnell über verschiedene Rollen und Projektbereiche hinweg anzuzeigen und zu ändern.

Die Gewährleistung der plattformübergreifenden Konsistenz stellte ihre eigenen Herausforderungen dar. Wir mussten benutzerdefinierte Berechtigungen einheitlich über unsere Web-, Desktop- und mobilen Anwendungen implementieren, jede mit ihrer eigenen Benutzeroberfläche und Benutzererfahrung. Dies war besonders knifflig für unsere mobilen Apps, die dynamisch Funktionen basierend auf den Berechtigungen des Benutzers ausblenden und anzeigen mussten. Wir adressierten dies, indem wir unsere Berechtigungslogik in der API-Schicht zentralisierten, um sicherzustellen, dass alle Plattformen konsistente Berechtigungsdaten erhielten.

Dann entwickelten wir ein flexibles UI-Framework, das sich in Echtzeit an diese Berechtigungsänderungen anpassen konnte und ein nahtloses Erlebnis unabhängig von der verwendeten Plattform bot.

Benutzerschulung und -akzeptanz stellten die letzte Hürde in unserer Reise zu benutzerdefinierten Berechtigungen dar. Die Einführung eines so leistungsstarken Features bedeutete, dass wir unseren Benutzern helfen mussten, benutzerdefinierte Berechtigungen zu verstehen und effektiv zu nutzen.

Wir führten benutzerdefinierte Berechtigungen zunächst für einen Teil unserer Benutzerbasis ein und überwachten sorgfältig deren Erfahrungen und sammelten Erkenntnisse. Dieser Ansatz ermöglichte es uns, das Feature und unsere Schulungsmaterialien basierend auf der realen Nutzung zu verfeinern, bevor wir es unserer gesamten Benutzerbasis zur Verfügung stellten.

Der schrittweise Rollout erwies sich als unschätzbar, da er uns half, kleinere Probleme und Punkte der Benutzerverwirrung zu identifizieren und zu beheben, die wir nicht vorhergesehen hatten, was letztendlich zu einem polierteren und benutzerfreundlicheren Feature für alle unsere Benutzer führte.

Dieser Ansatz, zunächst an eine Teilmenge von Benutzern zu starten, sowie unser typischer 2-3-wöchiger "Beta"-Zeitraum in unserer öffentlichen Beta, hilft uns, nachts ruhig zu schlafen. :)

## Ausblick

Wie bei allen Funktionen ist nichts jemals *"fertig"*.

Unsere langfristige Vision für das Feature der benutzerdefinierten Berechtigungen erstreckt sich über Tags, benutzerdefinierte Feldfilter, anpassbare Projektnavigation und Kommentarsteuerungen.

Lassen Sie uns in jeden Aspekt eintauchen.

### Tag-Berechtigungen

Wir denken, es wäre großartig, Berechtigungen basierend darauf zu erstellen, ob ein Datensatz ein oder mehrere Tags hat. Der offensichtlichste Anwendungsfall wäre, dass Sie eine benutzerdefinierte Benutzerrolle namens "Kunden" erstellen und nur Benutzern in dieser Rolle erlauben, Datensätze zu sehen, die das Tag "Kunden" haben.

Dies gibt Ihnen einen Überblick darüber, ob ein Datensatz von Ihren Kunden gesehen werden kann oder nicht.

Dies könnte noch mächtiger werden mit UND/ODER-Kombinatoren, bei denen Sie komplexere Regeln festlegen können. Zum Beispiel könnten Sie eine Regel einrichten, die den Zugriff auf Datensätze erlaubt, die sowohl mit "Kunden" als auch mit "Öffentlich" gekennzeichnet sind, oder auf Datensätze, die entweder mit "Intern" oder "Vertraulich" gekennzeichnet sind. Dieses Maß an Flexibilität würde unglaublich nuancierte Berechtigungseinstellungen ermöglichen, die selbst den komplexesten organisatorischen Strukturen und Arbeitsabläufen gerecht werden.

Die potenziellen Anwendungen sind vielfältig. Projektmanager könnten sensible Informationen leicht segregieren, Vertriebsteams könnten automatisch auf relevante Kundendaten zugreifen, und externe Mitarbeiter könnten nahtlos in bestimmte Teile eines Projekts integriert werden, ohne das Risiko, sensible interne Informationen offenzulegen.

### Benutzerdefinierte Feldfilter

Unsere Vision für benutzerdefinierte Feldfilter stellt einen bedeutenden Fortschritt in der granularen Zugriffskontrolle dar. Diese Funktion wird Projektadministratoren ermächtigen, festzulegen, welche Datensätze spezifische Benutzergruppen basierend auf den Werten benutzerdefinierter Felder sehen können. Es geht darum, dynamische, datengestützte Grenzen für den Informationszugang zu schaffen.

Stellen Sie sich vor, Sie könnten Berechtigungen wie folgt einrichten:

- Nur Datensätze anzeigen, bei denen das Dropdown "Projektstatus" auf "Öffentlich" eingestellt ist
- Sichtbarkeit auf Elemente beschränken, bei denen das Mehrfachauswahlfeld "Abteilung" "Marketing" enthält
- Zugriff auf Aufgaben erlauben, bei denen das Kontrollkästchen "Priorität" aktiviert ist
- Projekte anzeigen, bei denen das Zahlenfeld "Budget" über einem bestimmten Schwellenwert liegt

### Anpassbare Projektnavigation

Dies ist einfach eine Erweiterung der Schalter, die wir bereits haben. Anstatt nur Schalter für "Aktivität" und "Formulare" zu haben, möchten wir dies auf jeden einzelnen Teil der Projektnavigation ausweiten. Auf diese Weise können Projektadministratoren fokussierte Schnittstellen erstellen und Werkzeuge entfernen, die sie nicht benötigen.

### Kommentarsteuerungen

In Zukunft möchten wir kreativ sein, wie wir unseren Kunden erlauben, zu entscheiden, wer Kommentare sehen kann und wer nicht. Möglicherweise erlauben wir mehrere tabbed Kommentarbereiche unter einem Datensatz, und jeder kann für verschiedene Benutzergruppen sichtbar oder nicht sichtbar sein.

Darüber hinaus könnten wir auch eine Funktion zulassen, bei der nur Kommentare, in denen ein Benutzer *speziell* erwähnt wird, sichtbar sind und nichts anderes. Dies würde es Teams, die Kunden in Projekten haben, ermöglichen, sicherzustellen, dass nur Kommentare, die sie möchten, dass Kunden sie sehen, sichtbar sind.

## Fazit

Das ist also unser Ansatz, wie wir eines der interessantesten und leistungsstärksten Features entwickelt haben! [Wie Sie in unserem Projektmanagement-Vergleichstool sehen können](/compare), haben nur sehr wenige Projektmanagement-Systeme eine so leistungsstarke Berechtigungsmatrix, und die, die es tun, reservieren sie für ihre teuersten Unternehmenspläne, wodurch sie für ein typisches kleines oder mittelständisches Unternehmen unzugänglich werden.

Mit Blue haben Sie *alle* Funktionen in unserem Plan verfügbar – wir glauben nicht, dass Funktionen auf Unternehmensniveau nur für Unternehmenskunden reserviert sein sollten!