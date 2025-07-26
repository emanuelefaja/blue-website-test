--- 
title: Automatisierung im Projektmanagement — E-Mails an Stakeholder.
category: "Product Updates"
description: Oft möchten Sie die Kontrolle über Ihre Automatisierungen im Projektmanagement haben.
date: 2024-07-08
---

Wir haben bereits behandelt, wie man [E-Mail-Automatisierungen erstellt.](/insights/email-automations)

Oft gibt es jedoch Stakeholder in Projekten, die nur dann benachrichtigt werden müssen, wenn es etwas *wirklich* Wichtiges gibt.

Wäre es nicht schön, wenn es eine Automatisierung im Projektmanagement gäbe, bei der Sie als Projektmanager *genau* kontrollieren können, wann Sie einen wichtigen Stakeholder mit einem Knopfdruck benachrichtigen?

Es stellt sich heraus, dass Sie mit Blue genau das tun können!

Heute werden wir lernen, wie man eine wirklich nützliche Automatisierung im Projektmanagement erstellt:

Ein Kontrollkästchen, das automatisch einen oder mehrere wichtige Stakeholder benachrichtigt und ihnen den gesamten wichtigen Kontext gibt, worüber Sie sie benachrichtigen. Als Bonuspunkt werden wir auch lernen, wie man diese Möglichkeit einschränkt, sodass nur bestimmte Mitglieder Ihres Projekts diese E-Mail-Benachrichtigung auslösen können.

So wird es aussehen, wenn Sie fertig sind:

![](/insights/checkbox-email-automation.png)

Und allein durch das Aktivieren dieses Kontrollkästchens können Sie eine Automatisierung im Projektmanagement auslösen, die eine benutzerdefinierte Benachrichtigungs-E-Mail an Stakeholder sendet.

Lassen Sie uns Schritt für Schritt vorgehen.

## 1. Erstellen Sie Ihr benutzerdefiniertes Kontrollkästchenfeld

Das ist sehr einfach, Sie können unsere [detaillierte Dokumentation](https://documentation.blue.cc/custom-fields/introduction#creating-custom-fields) zum Erstellen benutzerdefinierter Felder einsehen.

Stellen Sie sicher, dass Sie dieses Feld mit einem einprägsamen Namen benennen, den Sie sich merken können, wie „Management benachrichtigen“ oder „Stakeholder benachrichtigen“.

## 2. Erstellen Sie Ihren Auslöser für die Automatisierung im Projektmanagement.

Klicken Sie in der Datensatzansicht Ihres Projekts auf den kleinen Roboter oben rechts, um die Automatisierungseinstellungen zu öffnen:

<video autoplay loop muted playsinline>
  <source src="/videos/notify-stakeholders-automation-setup.mp4" type="video/mp4">
</video>

## 3. Erstellen Sie Ihre Aktion für die Automatisierung im Projektmanagement.

In diesem Fall wird unsere Aktion darin bestehen, eine benutzerdefinierte E-Mail-Benachrichtigung an eine oder mehrere E-Mail-Adressen zu senden. Es ist wichtig zu beachten, dass diese Personen **nicht** in Blue sein müssen, um diese E-Mails zu erhalten; Sie können E-Mails an *jede* E-Mail-Adresse senden.

Sie können mehr in unserem [detaillierten Dokumentationsleitfaden zur Einrichtung von E-Mail-Automatisierungen](https://documentation.blue.cc/automations/actions/email-automations) erfahren.

Ihr Endergebnis sollte ungefähr so aussehen:

![](/insights/email-automation-example.png)

## 4. Bonus: Zugriff auf das Kontrollkästchen einschränken.

Sie können [benutzerdefinierte Rollen in Blue](/platform/features/user-permissions) verwenden, um den Zugriff auf die benutzerdefinierten Kontrollkästchenfelder einzuschränken, sodass nur autorisierte Teammitglieder E-Mail-Benachrichtigungen auslösen können.

Blue ermöglicht es Projektadministratoren, Rollen zu definieren und Berechtigungen für Benutzergruppen zuzuweisen. Dieses System ist entscheidend, um die Kontrolle darüber zu behalten, wer mit bestimmten Elementen Ihres Projekts, einschließlich benutzerdefinierter Felder wie dem Benachrichtigungs-Kontrollkästchen, interagieren kann.

1. Navigieren Sie zum Abschnitt Benutzerverwaltung in Blue und wählen Sie „Benutzerdefinierte Rollen“.
2. Erstellen Sie eine neue Rolle, indem Sie einen beschreibenden Namen und eine optionale Beschreibung angeben.
3. Suchen Sie innerhalb der Rollenberechtigungen den Abschnitt für den Zugriff auf benutzerdefinierte Felder.
4. Geben Sie an, ob die Rolle das Kontrollkästchen benutzerdefiniertes Feld anzeigen oder bearbeiten kann. Beispielsweise können Sie den Bearbeitungszugriff auf Rollen wie „Projektadministrator“ einschränken, während Sie einer neu erstellten benutzerdefinierten Rolle erlauben, dieses Feld zu verwalten.
5. Weisen Sie die neu erstellte Rolle den entsprechenden Benutzern oder Benutzergruppen zu. So wird sichergestellt, dass nur die vorgesehenen Personen die Möglichkeit haben, mit dem Benachrichtigungs-Kontrollkästchen zu interagieren.

[Lesen Sie mehr auf unserer offiziellen Dokumentationsseite.](https://documentation.blue.cc/user-management/roles/custom-user-roles)

Durch die Implementierung dieser benutzerdefinierten Rollen verbessern Sie die Sicherheit und Integrität Ihrer Prozesse im Projektmanagement. Nur autorisierte Teammitglieder können kritische E-Mail-Benachrichtigungen auslösen, sodass Stakeholder wichtige Updates erhalten, ohne unnötige Benachrichtigungen.

## Fazit

Durch die Implementierung der oben skizzierten Automatisierung im Projektmanagement erhalten Sie eine präzise Kontrolle darüber, wann und wie Sie wichtige Stakeholder benachrichtigen. Dieser Ansatz stellt sicher, dass wichtige Updates effektiv kommuniziert werden, ohne Ihre Stakeholder mit unnötigen Informationen zu überfluten. Mit den benutzerdefinierten Feldern und Automatisierungsfunktionen von Blue können Sie Ihren Prozess im Projektmanagement optimieren, die Kommunikation verbessern und ein hohes Maß an Effizienz aufrechterhalten.

Mit nur einem einfachen Kontrollkästchen können Sie benutzerdefinierte E-Mail-Benachrichtigungen auslösen, die auf die Bedürfnisse Ihres Projekts zugeschnitten sind, und sicherstellen, dass die richtigen Personen zur richtigen Zeit informiert werden. Darüber hinaus fügt die Möglichkeit, diese Funktionalität auf bestimmte Teammitglieder einzuschränken, eine zusätzliche Ebene der Kontrolle und Sicherheit hinzu.

Beginnen Sie noch heute, diese leistungsstarke Funktion in Blue zu nutzen, um Ihre Stakeholder informiert zu halten und Ihre Projekte reibungslos zu gestalten. Für detailliertere Schritte und zusätzliche Anpassungsoptionen verweisen Sie auf die bereitgestellten Dokumentationslinks. Viel Spaß beim Automatisieren!