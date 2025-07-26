---
title: Scalare Importazioni ed Esportazioni CSV a 250.000+ Record
category: "Engineering"
description: Scopra come Blue ha scalato le importazioni ed esportazioni CSV di 10x utilizzando Rust e architettura scalabile e scelte tecnologiche strategiche nel B2B SaaS.
date: 2024-07-18
---
In Blue, stiamo [costantemente spingendo i confini](/platform/roadmap) di ciò che è possibile nel software di gestione progetti. Nel corso degli anni, abbiamo [rilasciato centinaia di funzionalità](/platform/changelog)

La nostra ultima impresa ingegneristica? 

Una revisione completa del nostro sistema di [importazione CSV](https://documentation.blue.cc/integrations/csv-import) ed [esportazione](https://documentation.blue.cc/integrations/csv-export), migliorando drasticamente le prestazioni e la scalabilità. 

Questo post La porta dietro le quinte di come abbiamo affrontato questa sfida, le tecnologie che abbiamo impiegato e i risultati impressionanti che abbiamo ottenuto.

La cosa più interessante qui è che abbiamo dovuto uscire dal nostro tipico [stack tecnologico](https://sop.blue.cc/product/technology-stack) per ottenere i risultati che volevamo. Questa è una decisione che deve essere presa con attenzione, perché le ripercussioni a lungo termine possono essere severe in termini di debito tecnologico e overhead di manutenzione a lungo termine. 

<video autoplay loop muted playsinline>
  <source src="/videos/import-export-video.mp4" type="video/mp4">
</video>

## Scalare per le Esigenze Enterprise

Il nostro viaggio è iniziato con una richiesta da un cliente enterprise nel settore eventi. Questo cliente utilizza Blue come hub centrale per gestire vaste liste di eventi, venue e relatori, integrandolo perfettamente con il loro sito web. 

Per loro, Blue non è solo uno strumento — è l'unica fonte di verità per tutta la loro operazione.

Mentre siamo sempre orgogliosi di sentire che i clienti ci utilizzano per esigenze così mission-critical, c'è anche una grande responsabilità dalla nostra parte per garantire un sistema veloce e affidabile.

Mentre questo cliente scalava le sue operazioni, ha affrontato un ostacolo significativo: **importare ed esportare grandi file CSV contenenti da 100.000 a 200.000+ record.**

Questo era oltre la capacità del nostro sistema al momento. Infatti, il nostro precedente sistema di import/export stava già lottando con importazioni ed esportazioni contenenti più di 10.000-20.000 record! Quindi 200.000+ record era fuori questione. 

Gli utenti sperimentavano tempi di attesa frustrantemente lunghi, e in alcuni casi, le importazioni o esportazioni *non riuscivano a completarsi del tutto.* Questo influenzava significativamente le loro operazioni poiché si affidavano a importazioni ed esportazioni quotidiane per gestire certi aspetti delle loro operazioni. 

> La multi-tenancy è un'architettura dove una singola istanza di software serve più clienti (tenant). Pur essendo efficiente, richiede una gestione attenta delle risorse per garantire che le azioni di un tenant non impattino negativamente gli altri.

E questa limitazione non stava influenzando solo questo particolare cliente. 

A causa della nostra architettura multi-tenant—dove più clienti condividono la stessa infrastruttura—una singola importazione o esportazione che richiede molte risorse poteva potenzialmente rallentare le operazioni per altri utenti, cosa che in pratica accadeva spesso. 

Come al solito, abbiamo fatto un'analisi build vs buy, per capire se dovevamo spendere il tempo per aggiornare il nostro sistema o comprare un sistema da qualcun altro. Abbiamo esaminato varie possibilità.

Il fornitore che si è distinto è stato un provider SaaS chiamato [Flatfile](https://flatfile.com/). Il loro sistema e le loro capacità sembravano esattamente quello di cui avevamo bisogno. 

Ma, dopo aver rivisto i loro [prezzi](https://flatfile.com/pricing/), abbiamo deciso che questo sarebbe finito per essere una soluzione estremamente costosa per un'applicazione della nostra scala — *$2/file inizia ad accumularsi molto velocemente!* —ed era meglio estendere il nostro motore integrato di import/export CSV. 

Per affrontare questa sfida, abbiamo preso una decisione audace: introdurre Rust nel nostro stack tecnologico principalmente Javascript. Questo linguaggio di programmazione di sistema, noto per le sue prestazioni e sicurezza, era lo strumento perfetto per le nostre esigenze critiche in termini di prestazioni di parsing CSV e mappatura dati.

Ecco come abbiamo affrontato la soluzione.

### Introduzione dei Servizi in Background

La base della nostra soluzione è stata l'introduzione di servizi in background per gestire compiti che richiedono molte risorse. Questo approccio ci ha permesso di scaricare l'elaborazione pesante dal nostro server principale, migliorando significativamente le prestazioni complessive del sistema.
La nostra architettura di servizi in background è progettata con la scalabilità in mente. Come tutti i componenti della nostra infrastruttura, questi servizi scalano automaticamente in base alla domanda. 

Questo significa che durante i picchi, quando più importazioni o esportazioni grandi vengono elaborate simultaneamente, il sistema alloca automaticamente più risorse per gestire il carico aumentato. Al contrario, durante i periodi più tranquilli, scala verso il basso per ottimizzare l'uso delle risorse.

Questa architettura scalabile di servizi in background ha beneficiato Blue non solo per le importazioni ed esportazioni CSV. Nel tempo, abbiamo spostato un numero sostanziale di funzionalità nei servizi in background per alleggerire il carico dei nostri server principali:

- **[Calcoli delle Formule](https://documentation.blue.cc/custom-fields/formula)**: Scarica operazioni matematiche complesse per garantire aggiornamenti rapidi dei campi derivati senza impattare le prestazioni del server principale.
- **[Dashboard/Grafici](/platform/features/dashboards)**: Elabora grandi dataset in background per generare visualizzazioni aggiornate senza rallentare l'interfaccia utente.
- **[Indice di Ricerca](https://documentation.blue.cc/projects/search)**: Aggiorna continuamente l'indice di ricerca in background, garantendo risultati di ricerca veloci e accurati senza impattare le prestazioni del sistema.
- **[Copia Progetti](https://documentation.blue.cc/projects/copying-projects)**: Gestisce la replicazione di progetti grandi e complessi in background, permettendo agli utenti di continuare a lavorare mentre la copia viene creata.
- **[Automazioni di Gestione Progetti](/platform/features/automations)**: Esegue flussi di lavoro automatizzati definiti dall'utente in background, garantendo azioni tempestive senza bloccare altre operazioni.
- **[Record Ripetuti](https://documentation.blue.cc/records/repeat)**: Genera compiti o eventi ricorrenti in background, mantenendo l'accuratezza della programmazione senza appesantire l'applicazione principale.
- **[Campi Personalizzati Durata Tempo](https://documentation.blue.cc/custom-fields/duration)**: Calcola e aggiorna continuamente la differenza di tempo tra due eventi in Blue, fornendo dati di durata in tempo reale senza impattare la reattività del sistema.

## Nuovo Modulo Rust per l'Elaborazione Dati

Il cuore della nostra soluzione di elaborazione CSV è un modulo Rust personalizzato. Mentre questo ha segnato la nostra prima avventura al di fuori del nostro stack tecnologico principale di Javascript, la decisione di utilizzare Rust è stata guidata dalle sue prestazioni eccezionali nelle operazioni concorrenti e nei compiti di elaborazione file.

I punti di forza di Rust si allineano perfettamente con le esigenze del parsing CSV e della mappatura dati. Le sue astrazioni a costo zero permettono programmazione ad alto livello senza sacrificare le prestazioni, mentre il suo modello di ownership garantisce la sicurezza della memoria senza la necessità di garbage collection. Queste caratteristiche rendono Rust particolarmente abile nel gestire grandi dataset in modo efficiente e sicuro.

Per il parsing CSV, abbiamo sfruttato il crate csv di Rust, che offre lettura e scrittura ad alte prestazioni di dati CSV. Abbiamo combinato questo con logica di mappatura dati personalizzata per garantire un'integrazione perfetta con le strutture dati di Blue.

La curva di apprendimento per Rust è stata ripida ma gestibile. Il nostro team ha dedicato circa due settimane di apprendimento intensivo per questo.

I miglioramenti sono stati impressionanti:

![](/insights/import-export.png)


Il nostro nuovo sistema può elaborare la stessa quantità di record che il nostro vecchio sistema poteva elaborare in 15 minuti in circa 30 secondi. 

## Interazione Web Server e Database

Per il componente web server della nostra implementazione Rust, abbiamo scelto Rocket come nostro framework. Rocket si è distinto per la sua combinazione di prestazioni e caratteristiche user-friendly per gli sviluppatori. La sua tipizzazione statica e il controllo in fase di compilazione si allineano bene con i principi di sicurezza di Rust, aiutandoci a catturare potenziali problemi presto nel processo di sviluppo.
Sul fronte database, abbiamo optato per SQLx. Questa libreria SQL asincrona per Rust offre diversi vantaggi che l'hanno resa ideale per le nostre esigenze:

- SQL type-safe: SQLx ci permette di scrivere SQL crudo con query controllate in fase di compilazione, garantendo la sicurezza dei tipi senza sacrificare le prestazioni.
- Supporto async: Questo si allinea bene con Rocket e la nostra necessità di operazioni database efficienti e non bloccanti.
- Database agnostico: Mentre utilizziamo principalmente [AWS Aurora](https://aws.amazon.com/rds/aurora/), che è compatibile con MySQL, il supporto di SQLx per più database ci dà flessibilità per il futuro nel caso decidessimo mai di cambiare. 

## Ottimizzazione del Batching

Il nostro viaggio verso la configurazione ottimale di batching è stato uno di test rigorosi e analisi attenta. Abbiamo condotto benchmark estensivi con varie combinazioni di transazioni concorrenti e dimensioni dei chunk, misurando non solo la velocità grezza ma anche l'utilizzo delle risorse e la stabilità del sistema.

Il processo ha coinvolto la creazione di dataset di test di dimensioni e complessità variabili, simulando modelli di utilizzo del mondo reale. Abbiamo poi fatto passare questi dataset attraverso il nostro sistema, regolando il numero di transazioni concorrenti e la dimensione del chunk per ogni esecuzione.

Dopo aver analizzato i risultati, abbiamo scoperto che l'elaborazione di 5 transazioni concorrenti con una dimensione del chunk di 500 record forniva il miglior equilibrio tra velocità e utilizzo delle risorse. Questa configurazione ci permette di mantenere un alto throughput senza sovraccaricare il nostro database o consumare memoria eccessiva.

Interessante, abbiamo scoperto che aumentare la concorrenza oltre 5 transazioni non produceva guadagni di prestazioni significativi e a volte portava a maggiore contesa del database. Allo stesso modo, dimensioni di chunk più grandi miglioravano la velocità grezza ma al costo di un maggiore utilizzo della memoria e tempi di risposta più lunghi per importazioni/esportazioni di piccole e medie dimensioni.

## Esportazioni CSV tramite Link Email

L'ultimo pezzo della nostra soluzione affronta la sfida di consegnare grandi file esportati agli utenti. Invece di fornire un download diretto dalla nostra app web, che potrebbe portare a problemi di timeout e aumento del carico del server, abbiamo implementato un sistema di link di download via email.

Quando un utente avvia una grande esportazione, il nostro sistema elabora la richiesta in background. Una volta completata, piuttosto che tenere la connessione aperta o memorizzare il file sui nostri server web, carichiamo il file in una posizione di archiviazione temporanea sicura. Generiamo poi un link di download unico e sicuro e lo inviamo via email all'utente.

Questi link di download sono validi per 2 ore, trovando un equilibrio tra convenienza dell'utente e sicurezza delle informazioni. Questo lasso di tempo dà agli utenti ampia opportunità di recuperare i loro dati garantendo al contempo che le informazioni sensibili non vengano lasciate accessibili indefinitamente.

La sicurezza di questi link di download è stata una priorità principale nel nostro design. Ogni link è:

- Unico e generato casualmente, rendendolo praticamente impossibile da indovinare
- Valido solo per 2 ore
- Crittografato in transito, garantendo la sicurezza dei dati mentre vengono scaricati

Questo approccio offre diversi benefici:

- Riduce il carico sui nostri server web, poiché non devono gestire download di file grandi direttamente
- Migliora l'esperienza utente, specialmente per utenti con connessioni internet più lente che potrebbero affrontare problemi di timeout del browser con download diretti
- Fornisce una soluzione più affidabile per esportazioni molto grandi che potrebbero superare i tipici limiti di timeout web

Il feedback degli utenti su questa funzionalità è stato estremamente positivo, con molti che apprezzano la flessibilità che offre nella gestione di grandi esportazioni di dati.

## Esportazione di Dati Filtrati

L'altro ovvio miglioramento è stato permettere agli utenti di esportare solo i dati che erano già filtrati nella loro vista del progetto. Questo significa che se c'è un tag attivo "priorità", allora solo i record che hanno questo tag finirebbero nell'esportazione CSV. Questo significa meno tempo a manipolare i dati in Excel per filtrare via le cose che non sono importanti, e ci aiuta anche a ridurre il numero di righe da elaborare.

## Guardando Avanti

Mentre non abbiamo piani immediati per espandere il nostro uso di Rust, questo progetto ci ha mostrato il potenziale di questa tecnologia per operazioni critiche in termini di prestazioni. È un'opzione entusiasmante che ora abbiamo nel nostro toolkit per future esigenze di ottimizzazione. Questa revisione dell'importazione ed esportazione CSV si allinea perfettamente con l'impegno di Blue verso la scalabilità. 

Siamo dedicati a fornire una piattaforma che cresce con i nostri clienti, gestendo le loro esigenze di dati in espansione senza compromettere le prestazioni.

La decisione di introdurre Rust nel nostro stack tecnologico non è stata presa alla leggera. Ha sollevato una domanda importante che molti team di ingegneria affrontano: Quando è appropriato avventurarsi al di fuori del proprio stack tecnologico principale, e quando si dovrebbe rimanere con strumenti familiari?

Non c'è una risposta valida per tutti, ma in Blue, abbiamo sviluppato un framework per prendere queste decisioni cruciali:

- **Approccio Problem-First:** Iniziamo sempre definendo chiaramente il problema che stiamo cercando di risolvere. In questo caso, dovevamo migliorare drasticamente le prestazioni delle importazioni ed esportazioni CSV per grandi dataset.
- **Esaurire le Soluzioni Esistenti:** Prima di guardare al di fuori del nostro stack principale, esploriamo a fondo quello che può essere ottenuto con le nostre tecnologie esistenti. Questo spesso coinvolge profiling, ottimizzazione e ripensare il nostro approccio entro vincoli familiari.
- **Quantificare il Guadagno Potenziale:** Se stiamo considerando una nuova tecnologia, dobbiamo essere in grado di articolare chiaramente e, idealmente, quantificare i benefici. Per il nostro progetto CSV, abbiamo proiettato miglioramenti di ordini di grandezza nella velocità di elaborazione.
- **Valutare i Costi:** Introdurre una nuova tecnologia non riguarda solo il progetto immediato. Consideriamo i costi a lungo termine:
  - Curva di apprendimento per il team
  - Manutenzione e supporto continui
  - Potenziali complicazioni nel deployment e nelle operazioni
  - Impatto sull'assunzione e composizione del team
- **Contenimento e Integrazione:** Se introduciamo una nuova tecnologia, miriamo a contenerla in una parte specifica e ben definita del nostro sistema. Ci assicuriamo anche di avere un piano chiaro per come si integrerà con il nostro stack esistente.
- **Future-Proofing:** Consideriamo se questa scelta tecnologica apre opportunità future o se potrebbe metterci in un angolo.

Uno dei rischi primari dell'adozione frequente di nuove tecnologie è finire con quello che chiamiamo uno *"zoo tecnologico"* - un ecosistema frammentato dove diverse parti della vostra applicazione sono scritte in linguaggi o framework diversi, richiedendo un'ampia gamma di competenze specializzate per la manutenzione.


## Conclusione

Questo progetto esemplifica l'approccio di Blue all'ingegneria: *non abbiamo paura di uscire dalla nostra zona di comfort e adottare nuove tecnologie quando significa fornire un'esperienza significativamente migliore per i nostri utenti.* 

Ripensando il nostro processo di importazione ed esportazione CSV, non abbiamo solo risolto un'esigenza pressante per un cliente enterprise ma migliorato l'esperienza per tutti i nostri utenti che lavorano con grandi dataset.

Mentre continuiamo a spingere i confini di ciò che è possibile nel [software di gestione progetti](/solutions/use-case/project-management), siamo entusiasti di affrontare più sfide come questa. 

Rimanga sintonizzato per più [approfondimenti sull'ingegneria che alimenta Blue!](/insights/engineering-blog)