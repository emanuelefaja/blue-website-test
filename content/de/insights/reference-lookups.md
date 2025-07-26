---
title: Referenz- und Lookup-Benutzerdefinierte Felder
category: "Product Updates"
description: Erstellen Sie mühelos miteinander verbundene Projekte in Blue und verwandeln Sie es in eine einzige Informationsquelle für Ihr Unternehmen mit den neuen Referenz- und Lookup-Feldern.
date: 2023-11-01
---

Projekte in Blue sind bereits eine leistungsstarke Möglichkeit, Ihre Geschäftsdaten zu verwalten und die Arbeit voranzutreiben.

Heute machen wir den nächsten logischen Schritt und ermöglichen es Ihnen, Ihre Daten *zwischen* Projekten zu verknüpfen, um maximale Flexibilität und Leistung zu erzielen.

Die Verknüpfung von Projekten innerhalb von Blue verwandelt es in eine einzige Informationsquelle für Ihr Unternehmen. Diese Funktion ermöglicht die Erstellung eines umfassenden und miteinander verbundenen Datensatzes, der einen nahtlosen Datenfluss und eine verbesserte Sichtbarkeit über Projekte hinweg ermöglicht. Durch die Verknüpfung von Projekten können Teams eine einheitliche Sicht auf die Abläufe erreichen, was die Entscheidungsfindung und die betriebliche Effizienz verbessert.

## Ein Beispiel

Betrachten Sie die ACME Company, die die Referenz- und Lookup-Benutzerdefinierten Felder von Blue verwendet, um ein miteinander verbundenes Datenökosystem über ihre Projekte Kunden, Vertrieb und Inventar zu schaffen. Kundenakten im Projekt Kunden sind über Referenzfelder mit Verkaufstransaktionen im Projekt Vertrieb verknüpft. Diese Verknüpfung ermöglicht es Lookup-Feldern, zugehörige Kundendetails, wie Telefonnummern und Kontostände, direkt in jeden Verkaufsdatensatz zu ziehen. Darüber hinaus werden verkaufte Inventarartikel im Verkaufsdatensatz über ein Lookup-Feld angezeigt, das auf die Verkaufsmenge aus dem Projekt Inventar verweist. Schließlich sind Inventarentnahmen über ein Referenzfeld im Inventar mit den relevanten Verkäufen verbunden, das auf die Verkaufsdatensätze verweist. Dieses Setup bietet vollständige Transparenz darüber, welcher Verkauf die Inventarentnahme ausgelöst hat, und schafft eine integrierte 360-Grad-Sicht über die Projekte hinweg.

## Wie Referenzfelder funktionieren

Referenzbenutzerdefinierte Felder ermöglichen es Ihnen, Beziehungen zwischen Datensätzen in verschiedenen Projekten in Blue zu erstellen. Bei der Erstellung eines Referenzfeldes wählt der Projektadministrator das spezifische Projekt aus, das die Liste der Referenzdatensätze bereitstellt. Die Konfigurationsoptionen umfassen:

* **Einzelauswahl**: Ermöglicht die Auswahl eines Referenzdatensatzes.
* **Mehrfachauswahl**: Ermöglicht die Auswahl mehrerer Referenzdatensätze.
* **Filterung**: Setzen Sie Filter, um Benutzern die Auswahl nur der Datensätze zu ermöglichen, die den Filterkriterien entsprechen.

Sobald dies eingerichtet ist, können Benutzer spezifische Datensätze aus dem Dropdown-Menü innerhalb des Referenzfeldes auswählen und eine Verknüpfung zwischen den Projekten herstellen.

## Erweiterung von Referenzfeldern mit Lookups

Lookup-Benutzerdefinierte Felder werden verwendet, um Daten aus Datensätzen in anderen Projekten zu importieren und eine einseitige Sichtbarkeit zu schaffen. Sie sind immer schreibgeschützt und mit einem bestimmten Referenzbenutzerdefinierten Feld verbunden. Wenn ein Benutzer einen oder mehrere Datensätze mithilfe eines Referenzbenutzerdefinierten Feldes auswählt, zeigt das Lookup-Benutzerdefinierte Feld Daten aus diesen Datensätzen an. Lookups können Daten wie Folgendes anzeigen:

* Erstellt am
* Aktualisiert am
* Fälligkeitsdatum
* Beschreibung
* Liste
* Tag
* Zuweiser
* Jedes unterstützte benutzerdefinierte Feld aus dem referenzierten Datensatz — einschließlich anderer Lookup-Felder!

Stellen Sie sich beispielsweise ein Szenario vor, in dem Sie drei Projekte haben: **Projekt A** ist ein Verkaufsprojekt, **Projekt B** ist ein Projekt zur Verwaltung von Beständen und **Projekt C** ist ein Projekt für Kundenbeziehungen. In Projekt A haben Sie ein Referenzbenutzerdefiniertes Feld, das Verkaufsdatensätze mit den entsprechenden Kundendatensätzen in Projekt C verknüpft. In Projekt B haben Sie ein Lookup-Benutzerdefiniertes Feld, das Informationen aus Projekt A importiert, wie die verkaufte Menge. Auf diese Weise wird, wenn ein Verkaufsdatensatz in Projekt A erstellt wird, die mit diesem Verkauf verbundene Kundeninformation automatisch aus Projekt C übernommen, und die verkaufte Menge wird automatisch aus Projekt B übernommen. Dies ermöglicht es Ihnen, alle relevanten Informationen an einem Ort zu halten und anzuzeigen, ohne doppelte Daten zu erstellen oder Datensätze manuell über Projekte hinweg zu aktualisieren.

Ein reales Beispiel dafür ist ein E-Commerce-Unternehmen, das Blue zur Verwaltung seiner Verkäufe, Bestände und Kundenbeziehungen nutzt. In ihrem **Verkaufs**-Projekt haben sie ein Referenzbenutzerdefiniertes Feld, das jeden Verkaufsdatensatz mit dem entsprechenden **Kunden**-Datensatz in ihrem **Kunden**-Projekt verknüpft. In ihrem **Inventar**-Projekt haben sie ein Lookup-Benutzerdefiniertes Feld, das Informationen aus dem Verkaufsprojekt importiert, wie die verkaufte Menge, und diese im Inventarartikel-Datensatz anzeigt. Dies ermöglicht es ihnen, leicht zu sehen, welche Verkäufe die Inventarentnahmen vorantreiben, und ihre Bestände aktuell zu halten, ohne Datensätze manuell über Projekte hinweg aktualisieren zu müssen.

## Fazit

Stellen Sie sich eine Welt vor, in der Ihre Projektdaten nicht isoliert sind, sondern frei zwischen Projekten fließen, umfassende Einblicke bieten und die Effizienz steigern. Das ist die Kraft von Blues Referenz- und Lookup-Feldern. Durch die Ermöglichung nahtloser Datenverbindungen und die Bereitstellung von Echtzeit-Sichtbarkeit über Projekte hinweg verändern diese Funktionen, wie Teams zusammenarbeiten und Entscheidungen treffen. Egal, ob Sie Kundenbeziehungen verwalten, Verkäufe verfolgen oder Bestände überwachen, Referenz- und Lookup-Felder in Blue ermöglichen es Ihrem Team, intelligenter, schneller und effektiver zu arbeiten. Tauchen Sie ein in die vernetzte Welt von Blue und sehen Sie zu, wie Ihre Produktivität in die Höhe schnellt.

[Überprüfen Sie die Dokumentation](https://documentation.blue.cc/custom-fields/reference) oder [melden Sie sich an und probieren Sie es selbst aus.](https://app.blue.cc)