---
title: Búsqueda en tiempo real
category: "Product Updates"
description: Blue presenta un nuevo motor de búsqueda ultrarrápido que devuelve resultados en todos tus proyectos en milisegundos, permitiéndote cambiar de contexto en un abrir y cerrar de ojos.
date: 2024-03-01
---

Estamos emocionados de anunciar el lanzamiento de nuestro nuevo motor de búsqueda, diseñado para revolucionar la forma en que encuentras información dentro de Blue. La funcionalidad de búsqueda eficiente es crucial para una gestión de proyectos fluida, y nuestra última actualización garantiza que puedas acceder a tus datos más rápido que nunca.

Nuestro nuevo motor de búsqueda te permite buscar todos los comentarios, archivos, registros, campos personalizados, descripciones y listas de verificación. Ya sea que necesites encontrar un comentario específico hecho en un proyecto, localizar rápidamente un archivo o buscar un registro o campo particular, nuestro motor de búsqueda proporciona resultados ultrarrápidos.

A medida que las herramientas se acercan a una capacidad de respuesta de 50-100 ms, tienden a desvanecerse y mezclarse con el fondo, proporcionando una experiencia de usuario fluida y casi invisible. Para dar contexto, un parpadeo humano toma aproximadamente 60-120 ms, ¡así que 50 ms es en realidad más rápido que un parpadeo! Este nivel de capacidad de respuesta te permite interactuar con Blue sin siquiera darte cuenta de que está ahí, liberándote para concentrarte en el trabajo real en cuestión. Al aprovechar este nivel de rendimiento, nuestro nuevo motor de búsqueda asegura que puedas acceder rápidamente a la información que necesitas, sin que nunca interfiera con tu flujo de trabajo.

Para lograr nuestro objetivo de búsqueda ultrarrápida, utilizamos las últimas tecnologías de código abierto. Nuestro motor de búsqueda está construido sobre MeiliSearch, un popular servicio de búsqueda como servicio de código abierto que utiliza procesamiento de lenguaje natural y búsqueda vectorial para encontrar rápidamente resultados relevantes. Además, implementamos almacenamiento en memoria, lo que nos permite almacenar datos de acceso frecuente en RAM, reduciendo el tiempo que tarda en devolver los resultados de búsqueda. Esta combinación de MeiliSearch y almacenamiento en memoria permite que nuestro motor de búsqueda entregue resultados en milisegundos, haciéndote posible encontrar rápidamente lo que necesitas sin tener que pensar en la tecnología subyacente.

La nueva barra de búsqueda está convenientemente ubicada en la barra de navegación, lo que te permite comenzar a buscar de inmediato. Para una experiencia de búsqueda más detallada, simplemente presiona la tecla Tab mientras buscas para acceder a la página de búsqueda completa. Además, puedes activar rápidamente la función de búsqueda desde cualquier lugar utilizando el atajo CMD/Ctrl+K, lo que facilita aún más encontrar lo que necesitas.

<video autoplay loop muted playsinline>
  <source src="/videos/search-demo.mp4" type="video/mp4">
</video>


## Desarrollos Futuros

Esto es solo el comienzo. Ahora que tenemos una infraestructura de búsqueda de próxima generación, podemos hacer cosas realmente interesantes en el futuro.

Lo siguiente será la búsqueda semántica, que es una mejora significativa respecto a la búsqueda típica por palabras clave. Permítenos explicarlo.

Esta función permitirá que el motor de búsqueda entienda el contexto de tus consultas. Por ejemplo, buscar "mar" recuperará documentos relevantes incluso si no se utiliza la frase exacta. Podrías estar pensando "¡pero escribí 'océano' en su lugar!" - y tienes razón. El motor de búsqueda también entenderá la similitud entre "mar" y "océano", y devolverá documentos relevantes incluso si no se utiliza la frase exacta. Esta función es particularmente útil al buscar documentos que contienen términos técnicos, acrónimos o simplemente palabras comunes que tienen múltiples variaciones o errores tipográficos.

Otra función próxima es la capacidad de buscar imágenes por su contenido. Para lograr esto, procesaremos cada imagen en tu proyecto, creando una representación para cada una. En términos generales, una representación es un conjunto matemático de coordenadas que corresponde al significado de una imagen. Esto significa que todas las imágenes pueden ser buscadas en función de lo que contienen, independientemente de su nombre de archivo o metadatos. Imagina buscar "diagrama de flujo" y encontrar todas las imágenes relacionadas con diagramas de flujo, *independientemente de sus nombres de archivo.*