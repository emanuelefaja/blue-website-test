---
title: Como configurar o Blue como um CRM
category: "Melhores Práticas"
description: Aprenda como configurar o Blue para rastrear seus clientes e negócios de forma fácil.
date: 2024-08-11
---

## Introdução

Uma das principais vantagens de usar o Blue é não usá-lo para um caso de uso *específico*, mas usá-lo *entre* casos de uso. Isso significa que você não precisa pagar por várias ferramentas, e também tem um local onde pode facilmente alternar entre seus diversos projetos e processos, como contratação, vendas, marketing e muito mais.

Ao ajudar milhares de clientes a se configurarem no Blue ao longo dos anos, notamos que a parte difícil *não* é configurar o próprio Blue, mas pensar através dos processos e aproveitar ao máximo nossa plataforma.

As partes principais são pensar no fluxo de trabalho passo a passo para cada um dos processos empresariais que você quer rastrear, e também os detalhes específicos dos dados que você quer capturar, e como isso se traduz nos campos personalizados que você configura.

Hoje, vamos guiá-lo através da criação de [um sistema CRM de vendas fácil de usar, mas poderoso](/solutions/use-case/sales-crm) com uma base de dados de clientes que está vinculada a um pipeline de oportunidades. Todos esses dados fluirão para um painel onde você pode ver dados em tempo real sobre suas vendas totais, vendas previstas e muito mais.

## Base de Dados de Clientes

A primeira coisa a fazer é configurar um novo projeto para armazenar seus dados de clientes. Esses dados serão então referenciados cruzadamente em outro projeto onde você rastreia oportunidades específicas de vendas.

A razão pela qual separamos suas informações de clientes das oportunidades é que elas não mapeiam uma para uma.

Um cliente pode ter múltiplas oportunidades ou projetos.

Por exemplo, se você é uma agência de marketing e design, pode inicialmente se envolver com um cliente para sua marca, e depois fazer um projeto separado para seu site, e depois outro para o gerenciamento de suas redes sociais.

Todas essas seriam oportunidades de vendas separadas que requerem seu próprio rastreamento e propostas, mas todas estão vinculadas a esse um cliente.

A vantagem de separar sua base de dados de clientes em um projeto separado é que se você atualizar qualquer detalhe em sua base de dados de clientes, todas suas oportunidades automaticamente terão os novos dados, o que significa que agora você tem uma fonte única de verdade em seu negócio! Você não precisa voltar e editar tudo manualmente!

Então, a primeira coisa a decidir é se você vai ser centrado na empresa ou centrado na pessoa.

Esta decisão realmente depende do que você está vendendo e para quem você vende. Se você vende principalmente para empresas, então provavelmente vai querer que o nome do registro seja o nome da empresa. No entanto, se você vende principalmente para indivíduos (ou seja, você é um coach de saúde pessoal ou um consultor de marca pessoal), então você provavelmente tomaria uma abordagem centrada na pessoa.

Então o campo nome do registro vai ser ou o nome da empresa ou o nome da pessoa, dependendo da sua escolha. A razão para isso é que significa que você pode facilmente identificar um cliente de relance, apenas olhando para seu quadro ou base de dados.

Em seguida, você precisa considerar quais informações você quer capturar como parte de sua base de dados de clientes. Essas vão se tornar seus campos personalizados.

Os suspeitos usuais aqui são:

- Email
- Número de Telefone
- Website
- Endereço
- Fonte (ou seja, de onde esse cliente veio pela primeira vez?)
- Categoria

No Blue, você também pode remover quaisquer campos padrão de que não precisa. Para esta base de dados de clientes, normalmente recomendamos que você remova data de vencimento, responsável, dependências e listas de verificação. Você pode querer manter nosso campo de descrição padrão disponível caso tenha notas gerais sobre esse cliente que não são específicas para qualquer oportunidade de venda.

Recomendamos que você mantenha o campo "Referenciado por", pois isso será útil mais tarde. Uma vez que configurarmos nossa base de dados de oportunidades, poderemos ver todos os registros de vendas que estão vinculados a este cliente específico aqui.

Em termos de listas, normalmente vemos nossos clientes manterem simples e terem uma lista chamada "Clientes" e deixarem assim. É melhor usar tags ou campos personalizados para categorização.

O que é ótimo aqui é que uma vez que você tenha isso configurado, pode facilmente importar seus dados de outros sistemas ou planilhas Excel para o Blue via nossa função de importação CSV, e também pode criar um formulário para novos clientes potenciais submeterem seus detalhes para que você possa **automaticamente** capturá-los em sua base de dados.

## Base de Dados de Oportunidades

Agora que temos nossa base de dados de clientes, precisamos criar outro projeto para capturar nossas oportunidades de vendas reais. Você pode chamar este projeto de "CRM de Vendas" ou "Oportunidades".

### Listas como Etapas do Processo

Para configurar seu processo de vendas, você precisa pensar sobre quais são as etapas usuais pelas quais uma oportunidade passa desde o momento em que você recebe uma solicitação de um cliente até conseguir um contrato assinado.

Cada lista em seu projeto será uma etapa em seu processo.

Independentemente do seu processo específico, haverá algumas listas comuns que TODOS os CRMs de Vendas devem ter:

- Não Qualificado — Todas as solicitações recebidas, onde você ainda não qualificou um cliente.
- Fechado Ganho — Todas as oportunidades que você ganhou e transformou em vendas!
- Fechado Perdido — Todas as oportunidades onde você fez uma cotação para um cliente, e eles não aceitaram.
- N/A — Aqui é onde você coloca todas as oportunidades que você não ganhou, mas também não foram "perdidas". Podem ser as que você recusou, as onde o cliente, por qualquer razão, desapareceu, e assim por diante.

Em termos de pensar através do seu processo empresarial de CRM de vendas, você deve considerar o nível de granularidade que você quer. Não recomendamos ter 20 ou 30 colunas, isso normalmente fica confuso e impede você de conseguir ver o quadro geral.

No entanto, também é importante não tornar cada processo muito amplo, caso contrário as negociações ficarão "presas" em uma etapa específica por semanas ou meses, mesmo quando estão de fato avançando. Aqui está uma abordagem recomendada típica:

- **Não Qualificado**: Todas as solicitações recebidas, onde você ainda não qualificou um cliente.
- **Qualificação**: Aqui é onde você pega a oportunidade e inicia o processo de entender se isso é um bom ajuste para sua empresa.
- **Escrevendo Proposta**: Aqui é onde você começa a transformar a oportunidade em uma apresentação para sua empresa. Este é um documento que você enviaria para o cliente.
- **Proposta Enviada**: Aqui é onde você enviou a proposta para o cliente e está esperando uma resposta.
- **Negociações**: Aqui é onde você está no processo de finalizar o negócio.
- **Contrato Para Assinatura**: Aqui é onde você está apenas esperando o cliente assinar o contrato.
- **Fechado Ganho**: Aqui é onde você ganhou o negócio e agora está trabalhando no projeto.
- **Fechado Perdido**: Aqui é onde você fez uma cotação para o cliente, mas eles não aceitaram os termos.
- **N/A**: Aqui é onde você coloca todas as oportunidades que você não ganhou, mas também não foram "perdidas". Podem ser as que você recusou, as onde o cliente, por qualquer razão, desapareceu, e assim por diante.

### Tags como Categorias de Serviços
Vamos agora falar sobre tags.

Recomendamos que você use tags para os diferentes tipos de serviços que você oferece. Então, voltando ao nosso exemplo de agência de marketing e design, você pode ter tags para "branding", "website", "SEO", "Gerenciamento Facebook", e assim por diante.

As vantagens aqui são que você pode facilmente filtrar por serviço em um clique, o que pode dar uma visão geral de quais serviços são mais populares, e isso também pode informar futuras contratações, já que normalmente diferentes serviços requerem diferentes membros da equipe.

### Campos Personalizados do CRM de Vendas

Em seguida, precisamos considerar quais campos personalizados queremos ter.

Os típicos que vemos sendo usados são:

- **Valor**: Este é um campo de moeda para o valor do projeto
- **Custo**: Seu custo esperado para cumprir a venda, também um campo de moeda
- **Lucro**: Um campo de fórmula para calcular o lucro baseado nos campos de valor e custo.
- **URL da Proposta**: Pode incluir um link para um Google Doc ou documento Word online da sua proposta, para que você possa facilmente clicar e revisá-lo.
- **Arquivos Recebidos**: Pode ser um campo de arquivo personalizado onde você pode colocar quaisquer arquivos recebidos do cliente, como materiais de pesquisa, NDAs, e assim por diante.
- **Contratos**: Outro campo de arquivo personalizado onde você pode adicionar contratos assinados para guarda.
- **Nível de Confiança**: Um campo de estrelas personalizado com 5 estrelas, indicando quão confiante você está de ganhar esta oportunidade específica. Isso pode ser usado mais tarde no painel para previsão!
- **Data de Fechamento Esperada**: Um campo de data para estimar quando o negócio provavelmente fechará.
- **Cliente**: Um campo de referência vinculando à pessoa de contato principal na base de dados de clientes.
- **Nome do Cliente**: Um campo de consulta que puxa o nome do cliente do registro vinculado específico na base de dados de clientes.
- **Email do Cliente**: Um campo de consulta que puxa o email do cliente do registro vinculado específico na base de dados de clientes.
- **Fonte do Negócio**: Um campo dropdown para rastrear de onde a oportunidade se originou (por exemplo, indicação, website, chamada fria, feira comercial).
- **Razão da Perda**: Um campo dropdown (para negócios fechados perdidos) para categorizar por que a oportunidade foi perdida.
- **Tamanho do Cliente**: Um campo dropdown para categorizar clientes por tamanho (por exemplo, pequeno, médio, grande empresa).

Novamente, realmente **cabe a você** decidir precisamente quais campos você quer ter. Uma palavra de advertência: é fácil ao configurar adicionar muitos e muitos campos ao seu CRM de Vendas de dados que você gostaria de capturar. No entanto, você deve ser realista em termos da disciplina e comprometimento de tempo. Não adianta ter 30 campos em seu CRM de Vendas se 90% dos registros não terão nenhum dado neles.

A coisa ótima sobre campos personalizados é que eles se integram bem com [Permissões Personalizadas](/platform/features/user-permissions). Isso significa que você pode decidir exatamente quais campos os membros da equipe em sua equipe podem visualizar ou editar. Por exemplo, você pode querer esconder informações de custo e lucro de funcionários juniores.

### Automações

[Automações de CRM de Vendas](/platform/features/automations) são um recurso poderoso no Blue que pode simplificar seu processo de vendas, garantir consistência e economizar tempo em tarefas repetitivas. Ao configurar automações inteligentes, você pode aumentar a eficácia do seu CRM de vendas e permitir que sua equipe se concentre no que mais importa - fechar negócios. Aqui estão algumas automações principais a considerar para seu CRM de vendas:

- **Atribuição de Novos Leads**: Atribuir automaticamente novos leads a representantes de vendas baseado em critérios predefinidos como localização, tamanho do negócio ou indústria. Isso garante acompanhamento rápido e distribuição equilibrada da carga de trabalho.
- **Lembretes de Acompanhamento**: Configure lembretes automatizados para representantes de vendas acompanharem prospects após um certo período de inatividade. Isso ajuda a prevenir que leads caiam através das rachaduras.
- **Notificações de Progressão de Etapa**: Notificar membros relevantes da equipe quando um negócio se move para uma nova etapa no pipeline. Isso mantém todos informados do progresso e permite intervenções oportunas se necessário.
- **Alertas de Envelhecimento de Negócios**: Criar alertas para negócios que estiveram em uma etapa específica por mais tempo que o esperado. Isso ajuda a identificar negócios parados que podem precisar de atenção extra.

## Vinculando Clientes e Negócios

Um dos recursos mais poderosos do Blue para criar um sistema CRM eficaz é a capacidade de vincular sua base de dados de clientes com suas oportunidades de vendas. Esta conexão permite manter uma fonte única de verdade para informações de clientes enquanto rastreia múltiplos negócios associados a cada cliente. Vamos explorar como configurar isso usando campos personalizados de Referência e Consulta.

### Configurando o Campo de Referência

1. Em seu projeto de Oportunidades (ou CRM de Vendas), crie um novo campo personalizado.
2. Escolha o tipo de campo "Referência".
3. Selecione seu projeto de Base de Dados de Clientes como a fonte para a referência.
4. Configure o campo para permitir seleção única (já que cada oportunidade está tipicamente associada a um cliente).
5. Nomeie este campo algo como "Cliente" ou "Empresa Associada".

Agora, ao criar ou editar uma oportunidade, você poderá selecionar o cliente associado de um menu dropdown populado com registros de sua Base de Dados de Clientes.

### Aprimorando com Campos de Consulta

Uma vez que você estabeleceu a conexão de referência, pode usar campos de Consulta para trazer informações relevantes de clientes diretamente para sua visualização de oportunidades. Veja como:

1. Em seu projeto de Oportunidades, crie um novo campo personalizado.
2. Escolha o tipo de campo "Consulta".
3. Selecione o campo de Referência que você acabou de criar ("Cliente" ou "Empresa Associada") como a fonte.
4. Escolha quais informações de cliente você quer exibir. Você pode considerar campos como: Email, Número de Telefone, Categoria do Cliente, Gerente de Conta

Repita este processo para cada pedaço de informação de cliente que você quer exibir em sua visualização de oportunidades.

Os benefícios disso são:

- **Fonte Única de Verdade**: Atualize informações de cliente uma vez na Base de Dados de Clientes, e isso automaticamente reflete em todas as oportunidades vinculadas.
- **Eficiência**: Acesse rapidamente detalhes relevantes de clientes enquanto trabalha em oportunidades sem alternar entre projetos.
- **Integridade de Dados**: Reduza erros de entrada manual de dados puxando automaticamente informações de clientes.
- **Visão Holística**: Veja facilmente todas as oportunidades associadas a um cliente usando o campo "Referenciado Por" em sua Base de Dados de Clientes.

### Dica Avançada: Consulta de uma Consulta

O Blue oferece um recurso avançado chamado "Consulta de uma Consulta" que pode ser incrivelmente útil para configurações CRM complexas. Este recurso permite criar conexões através de múltiplos projetos, permitindo acessar informações tanto de sua Base de Dados de Clientes quanto do projeto de Oportunidades em um terceiro projeto.

Por exemplo, digamos que você tenha um espaço de trabalho "Projetos" onde gerencia o trabalho real para seus clientes. Você quer que este espaço de trabalho tenha acesso tanto a detalhes de clientes quanto informações de oportunidades. Veja como você pode configurar isso:

Primeiro, crie um campo de Referência em seu espaço de trabalho de Projetos que vincula ao projeto de Oportunidades. Isso estabelece a conexão inicial. Em seguida, crie campos de Consulta baseados nesta Referência para puxar detalhes específicos das oportunidades, como valor do negócio ou data de fechamento esperada.

O poder real vem no próximo passo: você pode criar campos de Consulta adicionais que alcançam através da Referência da oportunidade para a Base de Dados de Clientes. Isso permite puxar informações de clientes como detalhes de contato ou status da conta diretamente para seu espaço de trabalho de Projetos.

Esta cadeia de conexões dá uma visão abrangente em seu espaço de trabalho de Projetos, combinando dados tanto de suas oportunidades quanto base de dados de clientes. É uma maneira poderosa de garantir que suas equipes de projeto tenham todas as informações relevantes na ponta dos dedos sem precisar alternar entre diferentes projetos.

### Melhores Práticas para Sistemas CRM Vinculados

Mantenha sua Base de Dados de Clientes como a fonte única de verdade para todas as informações de clientes. Sempre que precisar atualizar detalhes de clientes, sempre faça isso na Base de Dados de Clientes primeiro. Isso garante que a informação permaneça consistente através de todos os projetos vinculados.

Ao criar campos de Referência e Consulta, use nomes claros e significativos. Isso ajuda a manter clareza, especialmente conforme seu sistema cresce mais complexo.

Revise regularmente sua configuração para garantir que está puxando as informações mais relevantes. Conforme suas necessidades empresariais evoluem, você pode precisar adicionar novos campos de Consulta ou remover os que não são mais úteis. Revisões periódicas ajudam a manter seu sistema simplificado e eficaz.

Considere aproveitar os recursos de automação do Blue para manter seus dados sincronizados e atualizados através dos projetos. Por exemplo, você poderia configurar uma automação para notificar membros relevantes da equipe quando informações importantes de clientes forem atualizadas na Base de Dados de Clientes.

Ao implementar eficazmente essas estratégias e fazer uso completo dos campos de Referência e Consulta, você pode criar um sistema CRM poderoso e interconectado no Blue. Este sistema fornecerá uma visão abrangente de 360 graus de seus relacionamentos com clientes e pipeline de vendas, permitindo tomada de decisões mais informada e operações mais suaves através de sua organização.

## Painéis

Painéis são um componente crucial de qualquer sistema CRM eficaz, fornecendo insights rápidos sobre seu desempenho de vendas e relacionamentos com clientes. O recurso de painel do Blue é particularmente poderoso porque permite combinar dados em tempo real de múltiplos projetos simultaneamente, dando uma visão abrangente de suas operações de vendas.

Ao configurar seu painel CRM no Blue, considere incluir várias métricas principais. Pipeline gerado por mês mostra o valor total de novas oportunidades adicionadas ao seu pipeline, ajudando você a rastrear a capacidade da sua equipe de gerar novos negócios. Vendas por mês exibe seus negócios realmente fechados, permitindo monitorar o desempenho da sua equipe em converter oportunidades em vendas.

Introduzir o conceito de descontos de pipeline pode levar a previsões mais precisas. Por exemplo, você pode contar 90% do valor de negócios na etapa "Contrato Para Assinatura", mas apenas 50% de negócios na etapa "Proposta Enviada". Esta abordagem ponderada fornece uma previsão de vendas mais realista.

Rastrear novas oportunidades por mês ajuda você a monitorar o número de novos negócios potenciais entrando em seu pipeline, que é um bom indicador dos esforços de prospecção da sua equipe de vendas. Dividir vendas por tipo pode ajudar você a identificar suas ofertas mais bem-sucedidas. Se você configurar um projeto de rastreamento de faturas vinculado às suas oportunidades, também pode rastrear receita real em seu painel, fornecendo um quadro completo de oportunidade até dinheiro.

O Blue oferece vários recursos poderosos para ajudá-lo a criar um painel CRM informativo e interativo. A plataforma fornece três tipos principais de gráficos: cartões de estatística, gráficos de pizza e gráficos de barras. Cartões de estatística são ideais para exibir métricas principais como valor total do pipeline ou número de oportunidades ativas. Gráficos de pizza são perfeitos para mostrar a composição de suas vendas por tipo ou a distribuição de negócios através de diferentes etapas. Gráficos de barras se destacam em comparar métricas ao longo do tempo, como vendas mensais ou novas oportunidades.

As capacidades sofisticadas de filtragem do Blue permitem segmentar seus dados por projeto, lista, tag e período de tempo. Isso é particularmente útil para aprofundar aspectos específicos dos seus dados de vendas ou comparar desempenho através de diferentes equipes ou produtos. A plataforma consolida inteligentemente listas e tags com o mesmo nome através de projetos, permitindo análise cruzada de projetos sem problemas. Isso é inestimável para uma configuração CRM onde você pode ter projetos separados para clientes, oportunidades e faturas.

Personalização é uma força principal do recurso de painel do Blue. A funcionalidade de arrastar e soltar e flexibilidade de exibição permitem criar um painel que se adequa perfeitamente às suas necessidades. Você pode facilmente reorganizar gráficos e escolher a visualização mais apropriada para cada métrica.

Embora os painéis sejam atualmente apenas para uso interno, você pode facilmente compartilhá-los com membros da equipe, concedendo permissões de visualização ou edição. Isso garante que todos em sua equipe de vendas tenham acesso aos insights de que precisam.

Ao aproveitar esses recursos e incluir as métricas principais que discutimos, você pode criar um painel CRM abrangente no Blue que fornece insights em tempo real sobre seu desempenho de vendas, saúde do pipeline e crescimento geral do negócio. Este painel se tornará uma ferramenta inestimável para tomar decisões baseadas em dados e manter toda sua equipe alinhada em seus objetivos e progresso de vendas.

## Conclusão

Configurar um CRM de vendas abrangente no Blue é uma maneira poderosa de simplificar seu processo de vendas e obter insights valiosos sobre seus relacionamentos com clientes e desempenho empresarial. Seguindo os passos descritos neste guia, você criou um sistema robusto que integra informações de clientes, oportunidades de vendas e métricas de desempenho em uma plataforma única e coesa.

Começamos criando uma base de dados de clientes, estabelecendo uma fonte única de verdade para todas as suas informações de clientes. Esta fundação permite manter registros precisos e atualizados para todos os seus clientes e prospects. Depois construímos sobre isso com uma base de dados de oportunidades, permitindo rastrear e gerenciar seu pipeline de vendas eficazmente.

Uma das principais forças de usar o Blue para seu CRM é a capacidade de vincular essas bases de dados usando campos de referência e consulta. Esta integração cria um sistema dinâmico onde atualizações nas informações de clientes são instantaneamente refletidas através de todas as oportunidades relacionadas, garantindo consistência de dados e economizando tempo em atualizações manuais.

Exploramos como aproveitar os recursos poderosos de automação do Blue para simplificar seu fluxo de trabalho, desde atribuir novos leads até enviar lembretes de acompanhamento. Essas automações ajudam a garantir que nenhuma oportunidade caia através das rachaduras e que sua equipe possa se concentrar em atividades de alto valor em vez de tarefas administrativas.

Finalmente, nos aprofundamos em criar painéis que fornecem insights rápidos sobre seu desempenho de vendas. Combinando dados de suas bases de dados de clientes e oportunidades, esses painéis oferecem uma visão abrangente do seu pipeline de vendas, negócios fechados e saúde geral do negócio.

Lembre-se, a chave para obter o máximo do seu CRM é uso consistente e refinamento regular. Encoraje sua equipe a adotar completamente o sistema, revise regularmente seus processos e automações, e continue a explorar novas maneiras de aproveitar os recursos do Blue para apoiar seus esforços de vendas.

Com esta configuração de CRM de vendas no Blue, você está bem equipado para nutrir relacionamentos com clientes, fechar mais negócios e impulsionar seu negócio para frente.