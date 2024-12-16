Aquí tienes un listado de los temas tratados en el documento, en español:

1. **Introducción a la Concurrencia**:
   - Definición de concurrencia.
   - Importancia de la concurrencia.
   - Ejecución de funciones en entornos de múltiples núcleos.

2. **Entorno de Computación**:
   - Procesadores de múltiples núcleos.
   - Ejecución paralela de funciones.

3. **Concurrencia vs. Paralelismo**:
   - Diferencia entre concurrencia y paralelismo.
   - Ejemplos de procesadores y procesos de navegador web.

4. **Conceptos de Procesos y Hilos**:
   - Definición de procesos y hilos.
   - Estado de los procesos y hilos.
   - Problemas de conmutación de contexto y tamaño de pila fijo.

5. **Condiciones de Carrera y Atomicidad**:
   - Problemas de acceso concurrente a la memoria compartida.
   - Sincronización de acceso a la memoria.
   - Uso de locks y posibles bloqueos.

6. **Goroutines en Go**:
   - Introducción a goroutines.
   - Ventajas de goroutines sobre los hilos del sistema operativo.
   - Cambios de contexto y uso de la pila.

7. **Scheduler de Go**:
   - Funcionamiento del scheduler de Go.
   - Asynchronous Preemption en Go 1.14.
   - Trabajo de robo (work stealing) y gestión de colas de ejecución.

8. **Canales en Go**:
   - Definición y uso de canales.
   - Canales sin buffer y con buffer.
   - Dirección de los canales y propiedad de los canales.

9. **Select en Go**:
   - Uso de la instrucción select para operaciones no bloqueantes y timeouts.
   - Implementación de comunicación no bloqueante.

10. **sync Package**:
    - Uso de sync.Mutex y sync.RWMutex.
    - Uso de sync.WaitGroup para la sincronización de goroutines.
    - Operaciones atómicas y variables de condición (sync.Cond).

11. **Cancelación y Contexto en Go**:
    - Uso del paquete context para la cancelación de operaciones.
    - Contextos con tiempo límite (WithTimeout y WithDeadline).
    - Propagación de contexto y datos de alcance de la solicitud.

12. **Patrones de Concurrencia**:
    - Construcción de pipelines para procesamiento de datos.
    - Fan-out y fan-in para paralelizar etapas computacionales.
    - Gestión de cancelación en pipelines.

13. **Interfaces en Go**:
    - Definición y uso de interfaces.
    - Implementación implícita de interfaces.
    - Asignación de interfaces y uso de type assertions.
    - Interface vacía y sus usos.

14. **Timeouts en Servidores HTTP**:
    - Importancia de establecer timeouts en servidores HTTP.
    - Uso de http.TimeoutHandler y propagación de timeouts a través de contextos.

15. **Detector de Condiciones de Carrera en Go**:
    - Uso de la herramienta de detección de condiciones de carrera.
    - Ejecución de pruebas con detección de carreras habilitada.

16. **Patrones de Concurrencia en Aplicaciones Reales**:
    - Implementación de pipelines para procesamiento de imágenes.
    - Uso de canales para la comunicación entre goroutines.
    - Estrategias para evitar filtraciones de goroutines (goroutine leaks).

Este listado abarca los temas principales y subtemas discutidos en el documento.