---
title: Campos personalizados de referencia y búsqueda
category: "Product Updates"
description: Crea proyectos interconectados en Blue sin esfuerzo, transformándolo en una única fuente de verdad para tu negocio con los nuevos Campos de Referencia y Búsqueda.
date: 2023-11-01
---

Los proyectos en Blue ya son una forma poderosa de gestionar los datos de tu negocio y avanzar en el trabajo.

Hoy, estamos dando el siguiente paso lógico y permitiéndote interconectar tus datos *entre* proyectos para obtener la máxima flexibilidad y poder.

Interconectar proyectos dentro de Blue lo transforma en una única fuente de verdad para tu negocio. Esta capacidad permite la creación de un conjunto de datos integral e interconectado, lo que facilita un flujo de datos sin problemas y una mayor visibilidad a través de los proyectos. Al vincular proyectos, los equipos pueden lograr una visión unificada de las operaciones, mejorando la toma de decisiones y la eficiencia operativa.

## Un ejemplo

Considera la empresa ACME, que utiliza los campos personalizados de Referencia y Búsqueda de Blue para crear un ecosistema de datos interconectado a través de sus proyectos de Clientes, Ventas e Inventario. Los registros de clientes en el proyecto de Clientes están vinculados a las transacciones de ventas en el proyecto de Ventas a través de campos de Referencia. Esta vinculación permite que los campos de Búsqueda extraigan detalles asociados del cliente, como números de teléfono y estados de cuenta, directamente en cada registro de venta. Además, los artículos de inventario vendidos se muestran en el registro de ventas a través de un campo de Búsqueda que hace referencia a los datos de Cantidad Vendida del proyecto de Inventario. Finalmente, los retiros de inventario están conectados a las ventas relevantes a través de un campo de Referencia en Inventario, apuntando de vuelta a los registros de Ventas. Esta configuración proporciona una visibilidad completa sobre qué venta desencadenó la eliminación de inventario, creando una vista integrada de 360 grados a través de los proyectos.

## Cómo funcionan los campos de referencia

Los campos personalizados de Referencia te permiten crear relaciones entre registros en diferentes proyectos en Blue. Al crear un campo de Referencia, el Administrador del Proyecto selecciona el proyecto específico que proporcionará la lista de registros de referencia. Las opciones de configuración incluyen:

* **Selección única**: Permite elegir un registro de referencia.
* **Selección múltiple**: Permite elegir múltiples registros de referencia.
* **Filtrado**: Establece filtros para permitir a los usuarios seleccionar solo registros que coincidan con los criterios de filtro.

Una vez configurado, los usuarios pueden seleccionar registros específicos del menú desplegable dentro del campo de Referencia, estableciendo un vínculo entre proyectos.

## Ampliando los campos de referencia utilizando búsquedas

Los campos personalizados de Búsqueda se utilizan para importar datos de registros en otros proyectos, creando visibilidad unidireccional. Siempre son de solo lectura y están conectados a un campo personalizado de Referencia específico. Cuando un usuario selecciona uno o más registros utilizando un campo personalizado de Referencia, el campo personalizado de Búsqueda mostrará datos de esos registros. Las búsquedas pueden mostrar datos como:

* Creado en
* Actualizado en
* Fecha de vencimiento
* Descripción
* Lista
* Etiqueta
* Asignado a
* Cualquier campo personalizado admitido del registro referenciado, ¡incluidos otros campos de búsqueda!

Por ejemplo, imagina un escenario donde tienes tres proyectos: **Proyecto A** es un proyecto de ventas, **Proyecto B** es un proyecto de gestión de inventario, y **Proyecto C** es un proyecto de relaciones con clientes. En el Proyecto A, tienes un campo personalizado de Referencia que vincula los registros de ventas a los registros de clientes correspondientes en el Proyecto C. En el Proyecto B, tienes un campo personalizado de Búsqueda que importa información del Proyecto A, como la cantidad vendida. De esta manera, cuando se crea un registro de venta en el Proyecto A, la información del cliente asociada con esa venta se extrae automáticamente del Proyecto C, y la cantidad vendida se extrae automáticamente del Proyecto B. Esto te permite mantener toda la información relevante en un solo lugar y visualizarla sin tener que crear datos duplicados o actualizar manualmente registros a través de proyectos.

Un ejemplo de la vida real de esto es una empresa de comercio electrónico que utiliza Blue para gestionar sus ventas, inventario y relaciones con clientes. En su proyecto de **Ventas**, tienen un campo personalizado de Referencia que vincula cada registro de venta al correspondiente registro de **Cliente** en su proyecto de **Clientes**. En su proyecto de **Inventario**, tienen un campo personalizado de Búsqueda que importa información del proyecto de Ventas, como la cantidad vendida, y la muestra en el registro del artículo de inventario. Esto les permite ver fácilmente qué ventas están impulsando las eliminaciones de inventario y mantener sus niveles de inventario actualizados sin tener que actualizar manualmente registros a través de proyectos.

## Conclusión

Imagina un mundo donde los datos de tus proyectos no están aislados, sino que fluyen libremente entre proyectos, proporcionando información integral y fomentando la eficiencia. Ese es el poder de los campos de Referencia y Búsqueda de Blue. Al permitir conexiones de datos sin problemas y proporcionar visibilidad en tiempo real a través de proyectos, estas características transforman la forma en que los equipos colaboran y toman decisiones. Ya sea que estés gestionando relaciones con clientes, rastreando ventas o supervisando inventario, los campos de Referencia y Búsqueda en Blue empoderan a tu equipo para trabajar de manera más inteligente, rápida y efectiva. Sumérgete en el mundo interconectado de Blue y observa cómo se eleva tu productividad.

[Consulta la documentación](https://documentation.blue.cc/custom-fields/reference) o [regístrate y pruébalo por ti mismo.](https://app.blue.cc)