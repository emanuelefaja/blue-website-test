---
title:  Cómo usamos Blue para construir Blue. 
category: "CEO Blog"
description: ¡Descubre cómo utilizamos nuestra propia plataforma de gestión de proyectos para construir nuestra plataforma de gestión de proyectos! 
date: 2024-08-07
---

Estás a punto de recibir un recorrido interno sobre cómo Blue construye Blue.

En Blue, consumimos nuestro propio producto.

Esto significa que usamos Blue para *construir* Blue.

Este término que suena extraño, a menudo llamado "dogfooding", se atribuye a Paul Maritz, un gerente de Microsoft en la década de 1980. Se dice que envió un correo electrónico con el asunto *"Comiendo nuestra propia comida para perros"* para alentar a los empleados de Microsoft a usar los productos de la empresa.

La idea de usar tus propias herramientas para construir tus herramientas es que conduce a un ciclo de retroalimentación positivo.

La idea de usar tus propias herramientas para construir tus herramientas lleva a un ciclo de retroalimentación positivo, creando numerosos beneficios:

- **Nos ayuda a identificar rápidamente problemas de usabilidad en el mundo real.** Al usar Blue a diario, encontramos los mismos desafíos que nuestros usuarios pueden enfrentar, lo que nos permite abordarlos de manera proactiva.
- **Acelera el descubrimiento de errores.** El uso interno a menudo revela errores antes de que lleguen a nuestros clientes, mejorando la calidad general del producto.
- **Aumenta nuestra empatía hacia los usuarios finales.** Nuestro equipo obtiene experiencia de primera mano sobre las fortalezas y debilidades de Blue, lo que nos ayuda a tomar decisiones más centradas en el usuario.
- **Impulsa una cultura de calidad dentro de nuestra organización.** Cuando todos usan el producto, hay un interés compartido en su excelencia.
- **Fomenta la innovación.** El uso regular a menudo genera ideas para nuevas características o mejoras, manteniendo a Blue a la vanguardia.

[Hemos hablado antes sobre por qué no tenemos un equipo de pruebas dedicado](/insights/open-beta) y esta es otra razón más.

Si hay errores en nuestro sistema, casi siempre los encontramos en nuestro uso diario constante de la plataforma. Y esto también crea una función de presión para solucionarlos, ya que obviamente los encontraremos muy molestos, ¡ya que probablemente somos uno de los principales usuarios de Blue!

Este enfoque demuestra nuestro compromiso con el producto. Al confiar en Blue nosotros mismos, mostramos a nuestros clientes que realmente creemos en lo que estamos construyendo. No es solo un producto que vendemos, es una herramienta de la que dependemos todos los días.

## Proceso Principal

Tenemos un proyecto en Blue, apropiadamente llamado "Producto".

**Todo** lo relacionado con el desarrollo de nuestro producto se rastrea aquí. Comentarios de clientes, errores, ideas de características, trabajo en curso, etc. La idea de tener un proyecto donde rastreamos todo es que [promueve un mejor trabajo en equipo.](/insights/great-teamwork)

Cada registro es una característica o parte de una característica. Así es como pasamos de "¿no sería genial si..." a "¡mira esta increíble nueva característica!"

El proyecto tiene las siguientes listas:

- **Ideas/Comentarios**: Esta es una lista de ideas del equipo o comentarios de clientes basados en llamadas o intercambios de correos electrónicos. ¡Siéntete libre de agregar cualquier idea aquí! En esta lista, aún no hemos decidido que construiremos alguna de estas características, pero revisamos regularmente esto en busca de ideas que queremos explorar más a fondo.
- **Backlog (Largo Plazo)**: Aquí es donde van las características de la lista de Ideas/Comentarios si decidimos que serían una buena adición a Blue.
- **{Trimestre Actual}**: Esto se estructura típicamente como "Qx AAAA" y muestra nuestras prioridades trimestrales.
- **Errores**: Esta es una lista de errores conocidos reportados por el equipo o clientes. Los errores añadidos aquí tendrán automáticamente la etiqueta "Error".
- **Especificaciones**: Estas características están actualmente siendo especificadas. No todas las características requieren una especificación o diseño; depende del tamaño esperado de la característica y del nivel de confianza que tengamos con respecto a los casos límite y la complejidad.
- **Backlog de Diseño**: Este es el backlog para los diseñadores; cada vez que terminan algo que está en progreso, pueden elegir cualquier elemento de esta lista.
- **Diseño en Progreso**: Estas son las características actuales que los diseñadores están diseñando.
- **Revisión de Diseño**: Aquí es donde están las características cuyos diseños están siendo revisados actualmente.
- **Backlog (Corto Plazo)**: Esta es una lista de características en las que probablemente comenzaremos a trabajar en las próximas semanas. Aquí es donde se realizan las asignaciones. El CEO y el Jefe de Ingeniería deciden qué características se asignan a qué ingeniero según la experiencia previa y la carga de trabajo. [Los miembros del equipo pueden luego trasladarlas a En Progreso](/insights/push-vs-pull-kanban) una vez que hayan completado su trabajo actual.
- **En Progreso**: Estas son características que se están desarrollando actualmente.
- **Revisión de Código**: Una vez que una característica ha terminado su desarrollo, pasa por una revisión de código. Luego, se moverá de nuevo a "En Progreso" si se necesitan ajustes o se desplegará en el entorno de Desarrollo.
- **Desarrollo**: Estas son todas las características actualmente en el entorno de Desarrollo. Otros miembros del equipo y ciertos clientes pueden revisarlas.
- **Beta**: Estas son todas las características actualmente en el [entorno Beta](https://beta.app.blue.cc). Muchos clientes utilizan esto como su plataforma diaria de Blue y también proporcionarán comentarios.
- **Producción**: Cuando una característica llega a producción, se considera finalizada.

A veces, al desarrollar una característica, nos damos cuenta de que ciertos subcomponentes son más difíciles de implementar de lo que inicialmente esperábamos, y podemos optar por no hacerlo en la versión inicial que desplegamos a los clientes. En este caso, podemos crear un nuevo registro con un nombre que siga el formato "{NombreDeLaCaracterística} V2" e incluir todos los subcomponentes como elementos de lista de verificación.

## Etiquetas

- **Móvil**: Esto significa que la característica es específica para nuestras aplicaciones de iOS, Android o iPad.
- **{NombreDelClienteEmpresarial}**: Se está construyendo una característica específicamente para un cliente empresarial. El seguimiento es importante, ya que generalmente hay acuerdos comerciales adicionales para cada característica.
- **Error**: Esto significa que es un error que requiere ser corregido.
- **Fast-Track**: Esto significa que es un Cambio de Fast-Track que no tiene que pasar por el ciclo completo de lanzamiento como se describe arriba.
- **Principal**: Esto es un desarrollo de característica importante. Típicamente se reserva para trabajos de infraestructura importantes, grandes actualizaciones de dependencias y módulos nuevos significativos dentro de Blue.
- **IA**: Esta característica contiene un componente de inteligencia artificial.
- **Seguridad**: Esto significa que se debe revisar una implicación de seguridad o se requiere un parche.

La etiqueta de fast-track es particularmente interesante. Esto se reserva para actualizaciones más pequeñas y menos complejas que no requieren nuestro ciclo completo de lanzamiento y que queremos enviar a los clientes dentro de 24-48 horas.

Los cambios de fast-track son típicamente ajustes menores que pueden mejorar significativamente la experiencia del usuario sin alterar la funcionalidad central. Piensa en corregir errores tipográficos en la interfaz de usuario, ajustar el relleno de los botones o agregar nuevos íconos para una mejor guía visual. Estos son el tipo de cambios que, aunque pequeños, pueden marcar una gran diferencia en cómo los usuarios perciben e interactúan con nuestro producto. ¡También son molestos si tardan mucho en enviarse!

Nuestro proceso de fast-track es sencillo.

Comienza creando una nueva rama desde la principal, implementando los cambios y luego creando solicitudes de fusión para cada rama objetivo: Desarrollo, Beta y Producción. Generamos un enlace de vista previa para revisión, asegurándonos de que incluso estos pequeños cambios cumplan con nuestros estándares de calidad. Una vez aprobados, los cambios se fusionan simultáneamente en todas las ramas, manteniendo nuestros entornos sincronizados.

## Campos Personalizados

No tenemos muchos campos personalizados en nuestro proyecto de Producto.

- **Especificaciones**: Esto enlaza a un documento de Blue que tiene la especificación para esa característica en particular. Esto no siempre se hace, ya que depende de la complejidad de la característica.
- **MR**: Este es el enlace a la Solicitud de Fusión en [Gitlab](https://gitlab.com) donde alojamos nuestro código.
- **Enlace de Vista Previa**: Para características que cambian principalmente el front-end, podemos crear una URL única que tenga esos cambios para cada commit, para que podamos revisar fácilmente los cambios.
- **Líder**: Este campo nos dice qué ingeniero senior está a cargo de la revisión de código. Asegura que cada característica reciba la atención experta que merece y siempre hay una persona de referencia clara para preguntas o inquietudes.

## Listas de Verificación

Durante nuestras demostraciones semanales, colocaremos los comentarios discutidos en una lista de verificación llamada "comentarios" y también habrá otra lista de verificación que contiene el [WBS (Estructura de Desglose del Trabajo)](/insights/simple-work-breakdown-structure) principal de la característica, para que podamos decir fácilmente qué está hecho y qué queda por hacer.

## Conclusión

¡Y eso es todo!

Creemos que a veces la gente se sorprende de lo sencillo que es nuestro proceso, pero creemos que los procesos simples son a menudo muy superiores a los procesos excesivamente complejos que no se pueden entender fácilmente.

Esta simplicidad es intencional. Nos permite mantenernos ágiles, responder rápidamente a las necesidades de los clientes y mantener a todo nuestro equipo alineado.

Al usar Blue para construir Blue, no solo estamos desarrollando un producto, ¡lo estamos viviendo!

Así que la próxima vez que uses Blue, recuerda: no solo estás usando un producto que hemos construido. Estás usando un producto del que dependemos personalmente todos los días.

Y eso marca toda la diferencia.