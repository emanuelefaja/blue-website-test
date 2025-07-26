---
title: Escalando Importaciones y Exportaciones CSV a 250,000+ Registros
category: "Engineering"
description: Descubre cómo Blue escaló las importaciones y exportaciones CSV 10x utilizando Rust y una arquitectura escalable junto con elecciones tecnológicas estratégicas en B2B SaaS.
date: 2024-07-18
---
En Blue, estamos [constantemente empujando los límites](/platform/roadmap) de lo que es posible en el software de gestión de proyectos. A lo largo de los años, hemos [lanzado cientos de funciones](/platform/changelog).

¿Nuestra última hazaña de ingeniería?

Una revisión completa de nuestro sistema de [importación](https://documentation.blue.cc/integrations/csv-import) y [exportación](https://documentation.blue.cc/integrations/csv-export) CSV, mejorando drásticamente el rendimiento y la escalabilidad.

Esta publicación te lleva detrás de escena sobre cómo abordamos este desafío, las tecnologías que empleamos y los impresionantes resultados que logramos.

Lo más interesante aquí es que tuvimos que salir de nuestro [stack tecnológico](https://sop.blue.cc/product/technology-stack) típico para lograr los resultados que deseábamos. Esta es una decisión que debe tomarse con cuidado, ya que la repercusión a largo plazo puede ser severa en términos de deuda tecnológica y costos de mantenimiento a largo plazo.

<video autoplay loop muted playsinline>
  <source src="/videos/import-export-video.mp4" type="video/mp4">
</video>

## Escalando para Necesidades Empresariales

Nuestro viaje comenzó con una solicitud de un cliente empresarial en la industria de eventos. Este cliente utiliza Blue como su centro central para gestionar vastas listas de eventos, lugares y oradores, integrándolo sin problemas con su sitio web.

Para ellos, Blue no es solo una herramienta: es la única fuente de verdad para toda su operación.

Si bien siempre estamos orgullosos de escuchar que los clientes nos utilizan para necesidades tan críticas, también hay una gran responsabilidad de nuestra parte para garantizar un sistema rápido y confiable.

A medida que este cliente escaló sus operaciones, se enfrentó a un obstáculo significativo: **importar y exportar grandes archivos CSV que contenían de 100,000 a 200,000+ registros.**

Esto estaba más allá de la capacidad de nuestro sistema en ese momento. De hecho, nuestro sistema anterior de importación/exportación ya estaba luchando con importaciones y exportaciones que contenían más de 10,000 a 20,000 registros. ¡Así que 200,000+ registros estaba fuera de cuestión!

Los usuarios experimentaron tiempos de espera frustrantemente largos y, en algunos casos, las importaciones o exportaciones *no se completaban en absoluto.* Esto afectó significativamente sus operaciones, ya que dependían de importaciones y exportaciones diarias para gestionar ciertos aspectos de sus operaciones.

> La multi-tenencia es una arquitectura donde una sola instancia de software sirve a múltiples clientes (inquilinos). Si bien es eficiente, requiere una gestión cuidadosa de los recursos para garantizar que las acciones de un inquilino no impacten negativamente a otros.

Y esta limitación no solo afectaba a este cliente en particular.

Debido a nuestra arquitectura multi-inquilino, donde múltiples clientes comparten la misma infraestructura, una única importación o exportación intensiva en recursos podría ralentizar potencialmente las operaciones para otros usuarios, lo que en la práctica a menudo sucedía.

Como de costumbre, realizamos un análisis de construir frente a comprar, para entender si debíamos dedicar tiempo a actualizar nuestro propio sistema o comprar un sistema de alguien más. Consideramos varias posibilidades.

El proveedor que se destacó fue un proveedor de SaaS llamado [Flatfile](https://flatfile.com/). Su sistema y capacidades parecían ser exactamente lo que necesitábamos.

Pero, después de revisar su [precio](https://flatfile.com/pricing/), decidimos que esto terminaría siendo una solución extremadamente costosa para una aplicación de nuestra escala — *¡$2/archivo se suma rápidamente!* — y era mejor extender nuestro motor de importación/exportación CSV integrado.

Para abordar este desafío, tomamos una decisión audaz: introducir Rust en nuestro stack tecnológico principalmente de Javascript. Este lenguaje de programación de sistemas, conocido por su rendimiento y seguridad, era la herramienta perfecta para nuestras necesidades críticas de análisis de CSV y mapeo de datos.

Así es como abordamos la solución.

### Introduciendo Servicios en Segundo Plano

La base de nuestra solución fue la introducción de servicios en segundo plano para manejar tareas intensivas en recursos. Este enfoque nos permitió descargar el procesamiento pesado de nuestro servidor principal, mejorando significativamente el rendimiento general del sistema. Nuestra arquitectura de servicios en segundo plano está diseñada con la escalabilidad en mente. Como todos los componentes de nuestra infraestructura, estos servicios se escalan automáticamente según la demanda.

Esto significa que durante los momentos pico, cuando se procesan simultáneamente múltiples importaciones o exportaciones grandes, el sistema asigna automáticamente más recursos para manejar la carga aumentada. Por el contrario, durante períodos más tranquilos, se reduce para optimizar el uso de recursos.

Esta arquitectura de servicios en segundo plano escalables ha beneficiado a Blue no solo para importaciones y exportaciones CSV. Con el tiempo, hemos trasladado un número sustancial de funciones a servicios en segundo plano para aliviar la carga de nuestros servidores principales:

- **[Cálculos de Fórmulas](https://documentation.blue.cc/custom-fields/formula)**: Descarga operaciones matemáticas complejas para asegurar actualizaciones rápidas de campos derivados sin afectar el rendimiento del servidor principal.
- **[Tableros/Gráficos](/platform/features/dashboards)**: Procesa grandes conjuntos de datos en segundo plano para generar visualizaciones actualizadas sin ralentizar la interfaz de usuario.
- **[Índice de Búsqueda](https://documentation.blue.cc/projects/search)**: Actualiza continuamente el índice de búsqueda en segundo plano, asegurando resultados de búsqueda rápidos y precisos sin afectar el rendimiento del sistema.
- **[Copiar Proyectos](https://documentation.blue.cc/projects/copying-projects)**: Maneja la replicación de proyectos grandes y complejos en segundo plano, permitiendo a los usuarios continuar trabajando mientras se crea la copia.
- **[Automatizaciones de Gestión de Proyectos](/platform/features/automations)**: Ejecuta flujos de trabajo automatizados definidos por el usuario en segundo plano, asegurando acciones oportunas sin bloquear otras operaciones.
- **[Registros Repetidos](https://documentation.blue.cc/records/repeat)**: Genera tareas o eventos recurrentes en segundo plano, manteniendo la precisión del horario sin sobrecargar la aplicación principal.
- **[Campos Personalizados de Duración de Tiempo](https://documentation.blue.cc/custom-fields/duration)**: Calcula y actualiza continuamente la diferencia de tiempo entre dos eventos en Blue, proporcionando datos de duración en tiempo real sin afectar la capacidad de respuesta del sistema.

## Nuevo Módulo Rust para Análisis de Datos

El corazón de nuestra solución de procesamiento CSV es un módulo Rust personalizado. Si bien esta marcó nuestra primera incursión fuera de nuestro stack tecnológico central de Javascript, la decisión de usar Rust fue impulsada por su excepcional rendimiento en operaciones concurrentes y tareas de procesamiento de archivos.

Las fortalezas de Rust se alinean perfectamente con las demandas del análisis de CSV y el mapeo de datos. Sus abstracciones de costo cero permiten una programación de alto nivel sin sacrificar rendimiento, mientras que su modelo de propiedad garantiza la seguridad de la memoria sin la necesidad de recolección de basura. Estas características hacen que Rust sea particularmente adecuado para manejar grandes conjuntos de datos de manera eficiente y segura.

Para el análisis de CSV, aprovechamos la crate csv de Rust, que ofrece lectura y escritura de datos CSV de alto rendimiento. Combinamos esto con lógica de mapeo de datos personalizada para asegurar una integración fluida con las estructuras de datos de Blue.

La curva de aprendizaje para Rust fue empinada pero manejable. Nuestro equipo dedicó aproximadamente dos semanas a un aprendizaje intensivo sobre esto.

Las mejoras fueron impresionantes:

![](/insights/import-export.png)

Nuestro nuevo sistema puede procesar la misma cantidad de registros que nuestro antiguo sistema podía procesar en 15 minutos en alrededor de 30 segundos.

## Interacción entre el Servidor Web y la Base de Datos

Para el componente del servidor web de nuestra implementación de Rust, elegimos Rocket como nuestro marco. Rocket se destacó por su combinación de rendimiento y características amigables para desarrolladores. Su tipado estático y verificación en tiempo de compilación se alinean bien con los principios de seguridad de Rust, ayudándonos a detectar problemas potenciales temprano en el proceso de desarrollo. En el frente de la base de datos, optamos por SQLx. Esta biblioteca SQL asíncrona para Rust ofrece varias ventajas que la hicieron ideal para nuestras necesidades:

- SQL seguro por tipo: SQLx nos permite escribir SQL en bruto con consultas verificadas en tiempo de compilación, asegurando la seguridad de tipo sin sacrificar rendimiento.
- Soporte asíncrono: Esto se alinea bien con Rocket y nuestra necesidad de operaciones de base de datos eficientes y no bloqueantes.
- Agnóstico a la base de datos: Si bien utilizamos principalmente [AWS Aurora](https://aws.amazon.com/rds/aurora/), que es compatible con MySQL, el soporte de SQLx para múltiples bases de datos nos brinda flexibilidad para el futuro en caso de que decidamos cambiar.

## Optimización del Procesamiento por Lotes

Nuestro viaje hacia la configuración óptima de procesamiento por lotes fue uno de pruebas rigurosas y análisis cuidadosos. Realizamos extensos benchmarks con varias combinaciones de transacciones concurrentes y tamaños de bloques, midiendo no solo la velocidad bruta sino también la utilización de recursos y la estabilidad del sistema.

El proceso involucró la creación de conjuntos de datos de prueba de diferentes tamaños y complejidades, simulando patrones de uso del mundo real. Luego, ejecutamos estos conjuntos de datos a través de nuestro sistema, ajustando el número de transacciones concurrentes y el tamaño del bloque para cada ejecución.

Después de analizar los resultados, encontramos que procesar 5 transacciones concurrentes con un tamaño de bloque de 500 registros proporcionaba el mejor equilibrio entre velocidad y utilización de recursos. Esta configuración nos permite mantener un alto rendimiento sin abrumar nuestra base de datos o consumir memoria excesiva.

Curiosamente, descubrimos que aumentar la concurrencia más allá de 5 transacciones no producía ganancias significativas en el rendimiento y, a veces, conducía a una mayor contención en la base de datos. De manera similar, tamaños de bloque más grandes mejoraban la velocidad bruta, pero a costa de un mayor uso de memoria y tiempos de respuesta más largos para importaciones/exportaciones pequeñas a medianas.

## Exportaciones CSV a través de Enlaces por Correo Electrónico

La pieza final de nuestra solución aborda el desafío de entregar grandes archivos exportados a los usuarios. En lugar de proporcionar una descarga directa desde nuestra aplicación web, lo que podría llevar a problemas de tiempo de espera y aumentar la carga del servidor, implementamos un sistema de enlaces de descarga por correo electrónico.

Cuando un usuario inicia una exportación grande, nuestro sistema procesa la solicitud en segundo plano. Una vez completada, en lugar de mantener la conexión abierta o almacenar el archivo en nuestros servidores web, subimos el archivo a una ubicación de almacenamiento temporal segura. Luego generamos un enlace de descarga único y seguro y se lo enviamos por correo electrónico al usuario.

Estos enlaces de descarga son válidos por 2 horas, logrando un equilibrio entre la conveniencia del usuario y la seguridad de la información. Este plazo brinda a los usuarios la oportunidad de recuperar sus datos mientras se asegura que la información sensible no quede accesible indefinidamente.

La seguridad de estos enlaces de descarga fue una prioridad en nuestro diseño. Cada enlace es:

- Único y generado aleatoriamente, lo que hace prácticamente imposible adivinarlo.
- Válido solo por 2 horas.
- Encriptado en tránsito, asegurando la seguridad de los datos mientras se descargan.

Este enfoque ofrece varios beneficios:

- Reduce la carga en nuestros servidores web, ya que no necesitan manejar descargas de archivos grandes directamente.
- Mejora la experiencia del usuario, especialmente para usuarios con conexiones a internet más lentas que podrían enfrentar problemas de tiempo de espera del navegador con descargas directas.
- Proporciona una solución más confiable para exportaciones muy grandes que podrían exceder los límites típicos de tiempo de espera web.

Los comentarios de los usuarios sobre esta función han sido abrumadoramente positivos, con muchos apreciando la flexibilidad que ofrece para gestionar grandes exportaciones de datos.

## Exportando Datos Filtrados

La otra mejora obvia fue permitir a los usuarios exportar solo los datos que ya estaban filtrados en su vista de proyecto. Esto significa que si hay una etiqueta activa "prioridad", entonces solo los registros que tienen esta etiqueta terminarían en la exportación CSV. Esto significa menos tiempo manipulando datos en Excel para filtrar cosas que no son importantes, y también nos ayuda a reducir el número de filas a procesar.

## Mirando Hacia Adelante

Si bien no tenemos planes inmediatos para expandir nuestro uso de Rust, este proyecto nos ha mostrado el potencial de esta tecnología para operaciones críticas de rendimiento. Es una opción emocionante que ahora tenemos en nuestra caja de herramientas para futuras necesidades de optimización. Esta revisión de importación y exportación CSV se alinea perfectamente con el compromiso de Blue con la escalabilidad.

Estamos dedicados a proporcionar una plataforma que crezca con nuestros clientes, manejando sus crecientes necesidades de datos sin comprometer el rendimiento.

La decisión de introducir Rust en nuestro stack tecnológico no se tomó a la ligera. Planteó una pregunta importante que muchos equipos de ingeniería enfrentan: ¿Cuándo es apropiado aventurarse fuera de tu stack tecnológico central y cuándo deberías apegarte a herramientas familiares?

No hay una respuesta única, pero en Blue hemos desarrollado un marco para tomar estas decisiones cruciales:

- **Enfoque Primero en el Problema:** Siempre comenzamos definiendo claramente el problema que estamos tratando de resolver. En este caso, necesitábamos mejorar drásticamente el rendimiento de las importaciones y exportaciones CSV para grandes conjuntos de datos.
- **Agotar Soluciones Existentes:** Antes de mirar fuera de nuestro stack central, exploramos a fondo lo que se puede lograr con nuestras tecnologías existentes. Esto a menudo implica perfilado, optimización y repensar nuestro enfoque dentro de las limitaciones familiares.
- **Cuantificar la Ganancia Potencial:** Si estamos considerando una nueva tecnología, necesitamos ser capaces de articular claramente y, idealmente, cuantificar los beneficios. Para nuestro proyecto CSV, proyectamos mejoras de orden de magnitud en la velocidad de procesamiento.
- **Evaluar los Costos:** Introducir una nueva tecnología no se trata solo del proyecto inmediato. Consideramos los costos a largo plazo:
  - Curva de aprendizaje para el equipo.
  - Mantenimiento y soporte continuos.
  - Complicaciones potenciales en la implementación y operaciones.
  - Impacto en la contratación y composición del equipo.
- **Contención e Integración:** Si introducimos una nueva tecnología, buscamos contenerla en una parte específica y bien definida de nuestro sistema. También aseguramos tener un plan claro sobre cómo se integrará con nuestro stack existente.
- **Preparación para el Futuro:** Consideramos si esta elección tecnológica abre oportunidades futuras o si podría encerrarnos en un rincón.

Uno de los principales riesgos de adoptar frecuentemente nuevas tecnologías es terminar con lo que llamamos un *"zoológico tecnológico"* - un ecosistema fragmentado donde diferentes partes de tu aplicación están escritas en diferentes lenguajes o marcos, requiriendo una amplia gama de habilidades especializadas para mantener.

## Conclusión

Este proyecto ejemplifica el enfoque de Blue hacia la ingeniería: *no tenemos miedo de salir de nuestra zona de confort y adoptar nuevas tecnologías cuando significa ofrecer una experiencia significativamente mejor para nuestros usuarios.*

Al reimaginar nuestro proceso de importación y exportación CSV, no solo hemos resuelto una necesidad urgente para un cliente empresarial, sino que hemos mejorado la experiencia para todos nuestros usuarios que manejan grandes conjuntos de datos.

A medida que continuamos empujando los límites de lo que es posible en [software de gestión de proyectos](/solutions/use-case/project-management), estamos emocionados de abordar más desafíos como este.

¡Mantente atento para más [profundizaciones en la ingeniería que impulsa a Blue!](/insights/engineering-blog)