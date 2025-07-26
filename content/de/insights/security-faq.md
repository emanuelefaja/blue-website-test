---
title: FAQ zu Blue Sicherheit
category: "FAQ"
description: Dies ist eine Liste der am häufigsten gestellten Fragen zu den Sicherheitsprotokollen und -praktiken bei Blue.
date: 2024-07-19
---

Unsere Mission ist es, die Arbeit der Welt zu organisieren, indem wir die beste Projektmanagement-Plattform auf dem Planeten entwickeln.

Zentral für die Erreichung dieser Mission ist die Gewährleistung, dass unsere Plattform sicher, zuverlässig und vertrauenswürdig ist. Wir verstehen, dass Blue, um Ihre einzige Quelle der Wahrheit zu sein, Ihre sensiblen Geschäftsdaten vor externen Bedrohungen, Datenverlust und Ausfallzeiten schützen muss.

Das bedeutet, dass wir bei Blue Sicherheit ernst nehmen.

Wenn wir über Sicherheit nachdenken, verfolgen wir einen ganzheitlichen Ansatz, der sich auf drei Schlüsselbereiche konzentriert:

1.  **Infrastruktur- und Netzwerksicherheit**: Stellt sicher, dass unsere physischen und virtuellen Systeme vor externen Bedrohungen und unbefugtem Zugriff geschützt sind.
2.  **Software-Sicherheit**: Konzentriert sich auf die Sicherheit des Codes selbst, einschließlich sicherer Codierungspraktiken, regelmäßiger Codeüberprüfungen und Schwachstellenmanagement.
3.  **Plattform-Sicherheit**: Umfasst die Funktionen innerhalb von Blue, wie z.B. [ausgeklügelte Zugriffskontrollen](/platform/features/user-permissions), die sicherstellen, dass Projekte standardmäßig privat sind, sowie andere Maßnahmen zum Schutz von Benutzerdaten und -privatsphäre.


## Wie skalierbar ist Blue?

Dies ist eine wichtige Frage, da Sie ein System möchten, das mit Ihnen *wachsen* kann. Sie möchten nicht in sechs oder zwölf Monaten Ihre Projekt- und Prozessmanagement-Plattform wechseln müssen.

Wir wählen Plattformanbieter sorgfältig aus, um sicherzustellen, dass sie die anspruchsvollen Arbeitslasten unserer Kunden bewältigen können. Wir nutzen Cloud-Dienste von einigen der weltweit führenden Cloud-Anbieter, die Unternehmen wie [Spotify](https://spotify.com) und [Netflix](https://netflix.com) unterstützen, die mehrere Größenordnungen mehr Traffic haben als wir.

Die Haupt-Cloud-Anbieter, die wir nutzen, sind:

- **[Cloudflare](https://cloudflare.com)**: Wir verwalten DNS (Domain Name Service) über Cloudflare sowie unsere Marketing-Website, die auf [Cloudflare Pages](https://pages.cloudflare.com/) läuft.
- **[Amazon Web Services](https://aws.amazon.com/)**: Wir verwenden AWS für unsere Datenbank, die [Aurora](https://aws.amazon.com/rds/aurora/) ist, für die Dateispeicherung über [Simple Storage Service (S3)](https://aws.amazon.com/s3/) und auch für den Versand von E-Mails über [Simple Email Service (SES)](https://aws.amazon.com/ses/).
- **[Render](https://render.com)**: Wir nutzen Render für unsere Front-End-Server, Anwendungs-/API-Server, unsere Hintergrunddienste, das Warteschlangensystem und die Redis-Datenbank. Interessanterweise ist Render tatsächlich *auf* AWS aufgebaut!


## Wie sicher sind Dateien in Blue?

Beginnen wir mit der Datenspeicherung. Unsere Dateien werden auf [AWS S3](https://aws.amazon.com/s3/) gehostet, dem weltweit beliebtesten Cloud-Objektspeicher mit branchenführender Skalierbarkeit, Datenverfügbarkeit, Sicherheit und Leistung.

Wir haben eine Verfügbarkeit von 99,99 % und eine hohe Haltbarkeit von 99,999999999 %.

Lassen Sie uns aufschlüsseln, was das bedeutet.

Verfügbarkeit bezieht sich auf die Zeit, in der die Daten betriebsbereit und zugänglich sind. Die 99,99 % Verfügbarkeit von Dateien bedeutet, dass wir erwarten können, dass Dateien nicht länger als etwa 8,76 Stunden pro Jahr nicht verfügbar sind.

Haltbarkeit bezieht sich auf die Wahrscheinlichkeit, dass Daten im Laufe der Zeit intakt und unbeschädigt bleiben. Dieses Maß an Haltbarkeit bedeutet, dass wir erwarten können, nicht mehr als eine Datei von 10 Milliarden hochgeladenen Dateien zu verlieren, dank umfangreicher Redundanz und Datenreplikation über mehrere Rechenzentren hinweg.

Wir verwenden [S3 Intelligent-Tiering](https://aws.amazon.com/s3/storage-classes/intelligent-tiering/), um Dateien automatisch in verschiedene Speicherklassen basierend auf der Zugriffsfrequenz zu verschieben. Basierend auf den Aktivitätsmustern von Hunderttausenden von Projekten stellen wir fest, dass die meisten Dateien in einem Muster abgerufen werden, das einer exponentiellen Rückoffkurve ähnelt. Das bedeutet, dass die meisten Dateien in den ersten Tagen sehr häufig abgerufen werden und dann schnell immer seltener abgerufen werden. Dies ermöglicht es uns, ältere Dateien in langsameren, aber deutlich günstigeren Speicher zu verschieben, ohne die Benutzererfahrung wesentlich zu beeinträchtigen.

Die Kosteneinsparungen sind erheblich. S3 Standard-Infrequent Access (S3 Standard-IA) ist etwa 1,84-mal günstiger als S3 Standard. Das bedeutet, dass wir für jeden Dollar, den wir für S3 Standard ausgegeben hätten, nur etwa 54 Cent für S3 Standard-IA für die gleiche Menge gespeicherter Daten ausgeben.

| Funktion                  | S3 Standard             | S3 Standard-IA       |
|--------------------------|-------------------------|-----------------------|
| Speicherkosten           | $0.023 - $0.021 pro GB  | $0.0125 pro GB        |
| Anforderungskosten (PUT usw.) | $0.005 pro 1.000 Anforderungen | $0.01 pro 1.000 Anforderungen |
| Anforderungskosten (GET) | $0.0004 pro 1.000 Anforderungen | $0.001 pro 1.000 Anforderungen |
| Datenabrufkosten         | $0.00 pro GB            | $0.01 pro GB          |


Die Dateien, die Sie über Blue hochladen, sind sowohl während der Übertragung als auch im Ruhezustand verschlüsselt. Daten, die zu und von Amazon S3 übertragen werden, sind durch [Transport Layer Security (TLS)](https://www.internetsociety.org/deploy360/tls/basics) gesichert, um [Abhören](https://en.wikipedia.org/wiki/Network_eavesdropping) und [Man-in-the-Middle-Angriffe](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) zu verhindern. Für die Verschlüsselung im Ruhezustand verwendet Amazon S3 die serverseitige Verschlüsselung (SSE-S3), die automatisch alle neuen Uploads mit AES-256-Verschlüsselung verschlüsselt, wobei Amazon die Verschlüsselungsschlüssel verwaltet. Dies stellt sicher, dass Ihre Daten während ihres gesamten Lebenszyklus sicher bleiben.

## Was ist mit nicht-dateibasierten Daten?

Unsere Datenbank wird von [AWS Aurora](https://aws.amazon.com/rds/aurora/) betrieben, einem modernen relationalen Datenbankdienst, der hohe Leistung, Verfügbarkeit und Sicherheit für Ihre Daten gewährleistet.

Daten in Aurora sind sowohl während der Übertragung als auch im Ruhezustand verschlüsselt. Wir verwenden SSL (AES-256), um Verbindungen zwischen Ihrer Datenbankinstanz und Ihrer Anwendung zu sichern und Daten während der Übertragung zu schützen. Für die Verschlüsselung im Ruhezustand verwendet Aurora Schlüssel, die über den AWS Key Management Service (KMS) verwaltet werden, um sicherzustellen, dass alle gespeicherten Daten, einschließlich automatisierter Backups, Snapshots und Replikate, verschlüsselt und geschützt sind.

Aurora verfügt über ein verteiltes, fehlertolerantes und selbstheilendes Speichersystem. Dieses System ist von den Rechenressourcen entkoppelt und kann bis zu 128 TiB pro Datenbankinstanz automatisch skalieren. Daten werden über drei [Verfügbarkeitszonen](https://aws.amazon.com/about-aws/global-infrastructure/regions_az/) (AZs) repliziert, um Widerstandsfähigkeit gegen Datenverlust zu gewährleisten und eine hohe Verfügbarkeit sicherzustellen. Im Falle eines Datenbankausfalls reduziert Aurora die Wiederherstellungszeiten auf weniger als 60 Sekunden, um minimale Unterbrechungen zu gewährleisten.

Blue sichert unsere Datenbank kontinuierlich auf Amazon S3, um eine Wiederherstellung zu einem bestimmten Zeitpunkt zu ermöglichen. Das bedeutet, dass wir die Blue-Hauptdatenbank auf jeden bestimmten Zeitpunkt innerhalb der letzten fünf Minuten wiederherstellen können, um sicherzustellen, dass Ihre Daten immer wiederherstellbar sind. Wir machen auch regelmäßige Snapshots der Datenbank für längere Backup-Aufbewahrungsfristen.

Als vollständig verwalteter Dienst automatisiert Aurora zeitaufwändige Verwaltungsaufgaben wie Hardwarebereitstellung, Datenbankeinrichtung, Patchen und Backups. Dies reduziert den operativen Aufwand und stellt sicher, dass unsere Datenbank immer auf dem neuesten Stand mit den neuesten Sicherheitspatches und Leistungsverbesserungen ist.

Wenn wir effizienter sind, können wir unsere Kosteneinsparungen an unsere Kunden mit unseren [branchenführenden Preisen](/pricing) weitergeben.

Aurora erfüllt verschiedene Branchenstandards wie HIPAA, GDPR und SOC 2 und stellt sicher, dass Ihre Datenmanagementpraktiken strengen regulatorischen Anforderungen entsprechen. Regelmäßige Sicherheitsprüfungen und die Integration mit [Amazon GuardDuty](https://aws.amazon.com/guardduty/) helfen, potenzielle Sicherheitsbedrohungen zu erkennen und zu mindern.

## Wie stellt Blue die Sicherheit des Logins sicher?

Blue verwendet [magische Links per E-Mail](https://documentation.blue.cc/user-management/magic-links), um sicheren und bequemen Zugriff auf Ihr Konto zu ermöglichen, wodurch die Notwendigkeit traditioneller Passwörter entfällt.

Dieser Ansatz verbessert die Sicherheit erheblich, indem er häufige Bedrohungen im Zusammenhang mit passwortbasierten Logins mindert. Durch die Eliminierung von Passwörtern schützen magische Links vor Phishing-Angriffen und Passwortdiebstahl, *da es kein Passwort gibt, das gestohlen oder ausgenutzt werden kann.*

Jeder magische Link ist nur für eine Login-Sitzung gültig, wodurch das Risiko unbefugten Zugriffs verringert wird. Darüber hinaus verfallen diese Links nach 15 Minuten, sodass ungenutzte Links nicht ausgenutzt werden können, was die Sicherheit weiter erhöht.

Der Komfort, den magische Links bieten, ist ebenfalls bemerkenswert. Magische Links ermöglichen eine mühelose Anmeldung, sodass Sie auf Ihr Konto *ohne* die Notwendigkeit zugreifen können, sich komplexe Passwörter zu merken.

Dies vereinfacht nicht nur den Anmeldeprozess, sondern verhindert auch Sicherheitsverletzungen, die auftreten, wenn Passwörter über mehrere Dienste hinweg wiederverwendet werden. Viele Benutzer tendieren dazu, dasselbe Passwort über verschiedene Plattformen hinweg zu verwenden, was bedeutet, dass ein Sicherheitsvorfall bei einem Dienst ihre Konten bei anderen Diensten, einschließlich Blue, gefährden könnte. Durch die Verwendung von magischen Links ist die Sicherheit von Blue nicht von den Sicherheitspraktiken anderer Dienste abhängig, was unseren Benutzern eine robustere und unabhängige Schutzschicht bietet.

Wenn Sie sich in Ihr Blue-Konto einloggen möchten, wird eine einzigartige Anmeldungs-URL an Ihre E-Mail gesendet. Ein Klick auf diesen Link wird Sie sofort in Ihr Konto einloggen. Der Link ist so konzipiert, dass er nach einmaliger Verwendung oder nach 15 Minuten, je nachdem, was zuerst eintritt, abläuft, was eine zusätzliche Sicherheitsebene hinzufügt. Durch die Verwendung von magischen Links stellt Blue sicher, dass Ihr Anmeldeprozess sowohl sicher als auch benutzerfreundlich ist und Ihnen Sicherheit und Komfort bietet.

## Wie kann ich die Zuverlässigkeit und Verfügbarkeit von Blue überprüfen?

Bei Blue sind wir bestrebt, ein hohes Maß an Zuverlässigkeit und Transparenz für unsere Benutzer aufrechtzuerhalten. Um Einblick in die Leistung unserer Plattform zu geben, bieten wir eine [dedizierte Systemstatusseite](https://status.blue.cc) an, die auch in der Fußzeile jeder Seite unserer Website verlinkt ist.

![](/insights/status-page.png)

Diese Seite zeigt unsere historischen Verfügbarkeitsdaten an, sodass Sie sehen können, wie konsistent unsere Dienste im Laufe der Zeit verfügbar waren. Darüber hinaus enthält die Statusseite detaillierte Vorfallberichte, die Transparenz über vergangene Probleme, deren Auswirkungen und die Schritte, die wir unternommen haben, um sie zu lösen und zukünftige Vorkommen zu verhindern, bieten.