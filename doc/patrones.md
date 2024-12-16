Los patrones de diseño se clasifican generalmente en tres categorías principales, sumando un total de 23 patrones clásicos identificados por primera vez por Erich Gamma, Richard Helm, Ralph Johnson y John Vlissides, conocidos colectivamente como la "Banda de los Cuatro" (Gang of Four, GoF) en su libro "Design Patterns: Elements of Reusable Object-Oriented Software". Las categorías son:

### 1. Patrones Creacionales (5 patrones)
Estos patrones se centran en cómo se crean los objetos. Su objetivo principal es facilitar la creación de objetos, especialmente en situaciones complejas. Los patrones creacionales incluyen:

1 - **Singleton**: Asegura que una clase tenga una única instancia y proporciona un punto de acceso global a ella.
2 - **Abstract Factory**: Proporciona una interfaz para crear familias de objetos relacionados o dependientes sin especificar sus clases concretas.
3 - **Factory Method**: Define una interfaz para crear un objeto, pero deja que las subclases decidan qué clase instanciar. El Factory Method permite a una clase deferir la instanciación a subclases.
5 - **Builder**: Separa la construcción de un objeto complejo de su representación, de modo que el mismo proceso de construcción pueda crear diferentes representaciones.
- **Prototype**: Crea nuevos objetos clonándolos de un objeto existente.

### 2. Patrones Estructurales (7 patrones)
Estos patrones se ocupan de cómo los objetos y clases se componen para formar estructuras más grandes. Los patrones estructurales facilitan la composición de interfaces o la implementación de nuevas funcionalidades a sistemas existentes. Incluyen:

- **Adapter (o Wrapper)**: Permite que interfaces incompatibles colaboren.
- **Bridge**: Separa una abstracción de su implementación, de modo que ambas puedan variar de forma independiente.
- **Composite**: Compone objetos en estructuras de árbol para representar jerarquías de parte-todo.
- **Decorator**: Añade responsabilidades adicionales a un objeto de manera dinámica.
- **Facade**: Proporciona una interfaz unificada a un conjunto de interfaces en un subsistema.
- **Flyweight**: Utiliza el compartir para soportar eficientemente grandes cantidades de objetos de grano fino.
- **Proxy**: Proporciona un sustituto o marcador de posición para otro objeto para controlar el acceso a él.

### 3. Patrones de Comportamiento (11 patrones)
Estos patrones se centran en la comunicación efectiva y la asignación de responsabilidades entre objetos. Los patrones de comportamiento incluyen:

- **Chain of Responsibility**: Pasa la solicitud a lo largo de una cadena de potenciales manejadores hasta que uno de ellos maneja la solicitud.
- **Command**: Encapsula una solicitud como un objeto, permitiendo la parametrización de clientes con colas, solicitudes y operaciones.
- **Interpreter**: Implementa un interpretador para un lenguaje.
- **Iterator**: Proporciona una manera de acceder a los elementos de un objeto agregado secuencialmente sin exponer su representación subyacente.
- **Mediator**: Define un objeto que encapsula cómo un conjunto de objetos interactúa.
- **Memento**: Sin violar el encapsulamiento, captura y externaliza el estado interno de un objeto para que el objeto pueda ser restaurado a este estado más tarde.
- **Observer**: Define una dependencia uno-a-muchos entre objetos de manera que cuando uno cambia de estado, todos sus dependientes son notificados automáticamente.
- **State**: Permite a un objeto alterar su comportamiento cuando su estado interno cambia.
- **Strategy**: Define una familia de algoritmos, encapsula cada uno de ellos y los hace intercambiables. La estrategia permite que el algoritmo varíe independientemente de los clientes que lo utilizan.
- **Template Method**: Define el esqueleto de un algoritmo en la operación, postergando algunos pasos a las subclases.
- **Visitor**: Define una nueva operación a una clase sin cambiar la clase en la que se opera.

Cada patrón tiene un propósito específico y soluciona problemas comunes de diseño en el desarrollo de software, facilitando así el desarrollo de sistemas más limpios, más eficientes y más mantenibles.