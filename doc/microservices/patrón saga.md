## Patrón Saga

El patrón Saga en el contexto de microservicios es una estrategia de manejo de transacciones que involucra múltiples servicios, donde cada paso afecta el estado en un sistema distribuido. Se utiliza principalmente cuando necesitas garantizar la consistencia de los datos a través de varios servicios que necesitan realizar una serie de operaciones relacionadas. Aquí te explico más sobre cómo funciona:

### Concepto
En sistemas tradicionales, las transacciones se manejan de forma atómica usando transacciones de bases de datos, generalmente con el soporte de transacciones ACID (Atomicidad, Consistencia, Aislamiento, Durabilidad). Sin embargo, en un sistema basado en microservicios, cada servicio podría tener su propia base de datos y es difícil mantener transacciones globales de manera eficiente y escalable.

### Cómo funciona el patrón Saga
El patrón Saga divide una transacción global en varias transacciones locales, cada una gestionada por un microservicio diferente. Cada transacción local puede tener su propio estado de confirmación y, en caso de que una transacción falle, el patrón Saga define compensaciones para cada una de las transacciones previas para revertir los cambios y mantener la consistencia del sistema.

### Tipos de Saga
Hay dos enfoques principales para implementar sagas:

1. **Sagas basadas en orquestación**: Un servicio coordinador (o orquestador) es responsable de iniciar cada transacción local y decidir los siguientes pasos, incluyendo compensaciones en caso de fallos. El orquestador centraliza la lógica de control y dirige el proceso global.

2. **Sagas basadas en coreografía**: Cada servicio sabe cuándo y cómo iniciar su transacción y qué hacer después, incluidas las compensaciones necesarias. Los servicios emiten eventos que otros servicios escuchan y reaccionan, coordinándose sin una autoridad central.

### Ventajas y desafíos
- **Ventajas**: Mayor resistencia y flexibilidad, ya que los servicios pueden manejar sus propias transacciones de forma autónoma. Facilita la escalabilidad porque no depende de una gestión centralizada de transacciones.

- **Desafíos**: La implementación puede ser más compleja, especialmente en el seguimiento del estado global y la gestión de compensaciones. La depuración y el monitoreo también pueden ser más difíciles debido a la naturaleza distribuida de la transacción.

El patrón Saga es especialmente útil en sistemas donde la disponibilidad y la escalabilidad son más críticas que la consistencia inmediata, adaptándose bien al principio de eventual consistencia en sistemas distribuidos.