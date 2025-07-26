---
title:  Por qué Blue tiene una Beta Abierta
category: "Engineering"
description: Descubre por qué nuestro sistema de gestión de proyectos tiene una beta abierta en curso. 
date: 2024-08-03
---

Muchas startups de B2B SaaS lanzan en Beta, y por buenas razones. Es parte del tradicional lema de Silicon Valley *“moverse rápido y romper cosas”*.

Poner una etiqueta de “beta” en un producto reduce las expectativas.

¿Algo está roto? Bueno, es solo una beta.

¿El sistema es lento? Bueno, es solo una beta.

¿La [documentación](https://blue.cc/docs) no existe? Bueno… ya entiendes el punto.

Y esto es *realmente* algo bueno. Reid Hoffman, el fundador de LinkedIn, dijo famosamente:

> Si no te sientes avergonzado por la primera versión de tu producto, has lanzado demasiado tarde.

Y la etiqueta beta también es buena para los clientes. Les ayuda a auto-seleccionarse.

Los clientes que prueban productos beta son aquellos que están en las primeras etapas del Ciclo de Adopción de Tecnología, también conocido como la Curva de Adopción de Productos.

El Ciclo de Adopción de Tecnología se divide típicamente en cinco segmentos principales:

1. Innovadores
2. Primeros Adoptantes
3. Mayoría Temprana
4. Mayoría Tardía
5. Rezagados

![](/insights/technology-adoption-lifecycle-graph.png)

Sin embargo, eventualmente el producto tiene que madurar, y los clientes esperan un producto estable y funcional. No quieren acceso a un entorno “beta” donde las cosas se rompen.

¿O sí?

*Esta* es la pregunta que nos hicimos.

Creemos que nos hicimos esta pregunta debido a la naturaleza de cómo se construyó inicialmente Blue. [Blue comenzó como una ramificación de una agencia de diseño ocupada](/insights/agency-success-playbook), y así trabajamos *dentro* de la oficina de un negocio que estaba utilizando activamente Blue para gestionar todos sus proyectos.

Esto significa que durante años, pudimos observar cómo *seres humanos reales* — ¡sentados justo al lado de nosotros! — usaban Blue en su vida diaria.

Y como usaron Blue desde los primeros días, ¡este equipo siempre usó Blue Beta!

Por lo tanto, fue natural para nosotros permitir que todos nuestros otros clientes también lo usaran.

**Y por eso no tenemos un equipo de pruebas dedicado.**

Así es.

Nadie en Blue tiene la *única* responsabilidad de garantizar que nuestra plataforma funcione bien y de manera estable.

Esto es por varias razones.

La primera es una base de costos más baja.

No tener un equipo de pruebas a tiempo completo reduce significativamente nuestros costos, y podemos transferir estos ahorros a nuestros clientes con los precios más bajos de la industria.

Para poner esto en perspectiva, ofrecemos conjuntos de características de nivel empresarial que nuestra competencia cobra entre $30 y $55/usuario/mes por solo $7/mes.

Esto no sucede por accidente, es *intencional*.

Sin embargo, no es una buena estrategia vender un producto más barato si no funciona.

Así que la *verdadera pregunta es*, ¿cómo logramos crear una plataforma estable que miles de clientes puedan usar sin un equipo de pruebas dedicado?

Por supuesto, nuestro enfoque de tener una Beta abierta es crucial para esto, pero antes de profundizar en esto, queremos tocar la responsabilidad del desarrollador.

Tomamos la decisión temprana en Blue de que nunca dividiríamos las responsabilidades para tecnologías de front-end y back-end. Solo contrataríamos o capacitaríamos desarrolladores de pila completa.

La razón por la que tomamos esta decisión fue para asegurar que un desarrollador tuviera plena propiedad de la función en la que estaba trabajando. Así que no habría la mentalidad de *“lanzar el problema por encima de la cerca del jardín”* que a veces se obtiene cuando hay responsabilidades compartidas para las funciones.

Y esto se extiende a las pruebas de la función, a entender los casos de uso y solicitudes del cliente, y a leer y comentar sobre las especificaciones.

En otras palabras, cada desarrollador construye una comprensión profunda e intuitiva de la función que está construyendo.

Bien, hablemos ahora de nuestra beta abierta.

Cuando decimos que es “abierta”, lo decimos en serio. Cualquier cliente puede probarla simplemente añadiendo “beta” delante de la URL de nuestra aplicación web.

Así que “app.blue.cc” se convierte en “beta.app.blue.cc”.

Cuando hacen esto, pueden ver sus datos habituales, ya que tanto los entornos Beta como de Producción comparten la misma base de datos, pero también podrán ver nuevas características.

Los clientes pueden trabajar fácilmente incluso si tienen algunos miembros del equipo en Producción y algunos curiosos en Beta.

Normalmente tenemos unos pocos cientos de clientes usando Beta en cualquier momento, y publicamos avances de características en nuestros foros comunitarios para que puedan ver qué hay de nuevo y probarlo.

Y este es el punto: ¡tenemos *varios cientos* de testers!

Todos estos clientes probarán características en sus flujos de trabajo y serán bastante vocales si algo no está del todo bien, porque ya están *implementando* la función dentro de su negocio.

Los comentarios más comunes son pequeños pero muy útiles cambios que abordan casos extremos que no consideramos.

Dejamos nuevas características en Beta entre 2 y 4 semanas. Siempre que sintamos que son estables, las lanzamos a producción.

También tenemos la capacidad de omitir Beta si es necesario, utilizando una bandera de vía rápida. Esto se hace típicamente para correcciones de errores que no queremos retener durante 2-4 semanas antes de enviarlas a producción.

¿El resultado?

Enviar a producción se siente… bueno, aburrido. Como nada, simplemente no es un gran problema para nosotros.

Y significa que esto suaviza nuestro calendario de lanzamientos, lo que nos ha permitido [enviar características mensualmente como un reloj durante los últimos seis años.](/changelog).

Sin embargo, como cualquier elección, hay algunos compromisos.

El soporte al cliente es ligeramente más complejo, ya que tenemos que apoyar a los clientes en dos versiones de nuestra plataforma. A veces esto puede causar confusión a los clientes que tienen miembros del equipo usando dos versiones diferentes.

Otro punto problemático es que este enfoque puede ralentizar el calendario general de lanzamientos a producción. Esto es especialmente cierto para características más grandes que pueden quedar “atascadas” en Beta si hay otra característica relacionada que está teniendo problemas y necesita más trabajo.

Pero en general, creemos que estos compromisos valen la pena por los beneficios de una base de costos más baja y un mayor compromiso del cliente.

Somos una de las pocas empresas de software que adopta este enfoque, pero ahora es una parte fundamental de nuestro enfoque de desarrollo de productos.