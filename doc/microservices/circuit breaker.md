## Circuit Breaker

El patrón Circuit Breaker (Interruptor de Circuito) es un patrón de diseño utilizado en sistemas distribuidos para prevenir que fallos en un servicio o dependencia se propaguen a otros servicios y evitar cascadas de errores. Funciona de manera similar a un interruptor eléctrico, abriéndose para interrumpir el flujo cuando detecta fallos y cerrándose de nuevo cuando el sistema se recupera.

### Estados del Circuit Breaker

1. **Closed (Cerrado)**:
   - **Normal Operation**: Todas las solicitudes pasan al servicio de destino.
   - **Error Tracking**: Se monitorean los errores y tiempos de respuesta.
   - **Threshold**: Si los errores superan un umbral predefinido (por ejemplo, un porcentaje de fallos en un período de tiempo), el circuito se abre.

2. **Open (Abierto)**:
   - **Blocking**: Las solicitudes se bloquean inmediatamente y se devuelve una respuesta de error o se ejecuta un método de fallback.
   - **Cool Down Period**: Después de un período de enfriamiento, el circuito entra en un estado de prueba.

3. **Half-Open (Medio Abierto)**:
   - **Trial Requests**: Permite pasar un número limitado de solicitudes para probar si el servicio de destino se ha recuperado.
   - **Recovery**: Si las solicitudes son exitosas, el circuito se cierra de nuevo.
   - **Failure**: Si las solicitudes fallan, el circuito se abre nuevamente.

### Funcionamiento del Circuit Breaker

1. **Monitoreo de Solicitudes**: Mientras el circuito está cerrado, el sistema monitorea las solicitudes para detectar fallos y tiempos de respuesta.
2. **Apertura del Circuito**: Si los fallos superan el umbral definido, el circuito se abre. Esto significa que las solicitudes posteriores no se enviarán al servicio fallido, sino que recibirán una respuesta de error o se ejecutará un método de fallback.
3. **Período de Enfriamiento**: Después de que el circuito se abre, hay un período de tiempo en el que no se envían solicitudes al servicio fallido.
4. **Reintento de Solicitudes**: Después del período de enfriamiento, el circuito pasa a un estado medio abierto y permite que un número limitado de solicitudes pasen. Esto prueba si el servicio se ha recuperado.
5. **Cierre del Circuito**: Si las solicitudes de prueba son exitosas, el circuito se cierra y el flujo normal de solicitudes se reanuda. Si fallan, el circuito se abre nuevamente.

### Ejemplo de Implementación en Go con Hystrix

Aquí tienes un ejemplo básico de cómo implementar el patrón Circuit Breaker usando Netflix Hystrix en una aplicación Go:

1. **Instala Hystrix**:
   ```sh
   go get github.com/afex/hystrix-go/hystrix
   ```

2. **Configura y Usa Hystrix**:
   ```go
   package main

   import (
       "context"
       "log"

       "github.com/afex/hystrix-go/hystrix"
       "github.com/micro/go-micro/v2"
       "github.com/micro/go-micro/v2/client"
   )

   func main() {
       // Configura Hystrix
       hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
           Timeout:               1000, // tiempo de espera en milisegundos
           MaxConcurrentRequests: 100,  // número máximo de solicitudes concurrentes
           ErrorPercentThreshold: 25,   // umbral de porcentaje de errores
       })

       // Crea un nuevo servicio
       service := micro.NewService(
           micro.Name("example.service"),
       )

       // Inicializa el servicio
       service.Init()

       // Crea un cliente Hystrix
       hystrixClient := client.NewClient(
           client.Wrap(hystrixWrapper),
       )

       // Llama a un servicio usando el cliente Hystrix
       req := hystrixClient.NewRequest("example.service", "Example.Endpoint", &Request{})
       rsp := &Response{}

       // Ejecuta la llamada
       err := hystrix.Do("my_command", func() error {
           return hystrixClient.Call(context.Background(), req, rsp)
       }, nil)

       if err != nil {
           log.Fatalf("Error calling service: %v", err)
       }

       log.Printf("Response: %v", rsp)
   }

   func hystrixWrapper(c client.Client) client.Client {
       return &hystrixClient{c}
   }

   type hystrixClient struct {
       client.Client
   }

   type Request struct {
       // define tus campos aquí
   }

   type Response struct {
       // define tus campos aquí
   }
   ```

### Ejemplo en Go con `goresilience`

Vamos a usar la biblioteca `goresilience` para implementar un Circuit Breaker en Go.

1. **Instalación de la biblioteca**:
   ```sh
   go get github.com/slok/goresilience/circuitbreaker
   ```

2. **Código de Ejemplo**:
   ```go
   package main

   import (
       "errors"
       "fmt"
       "time"

       "github.com/slok/goresilience"
       "github.com/slok/goresilience/circuitbreaker"
       "github.com/slok/goresilience/retry"
   )

   func main() {
       // Configuración del Circuit Breaker
       cbConfig := circuitbreaker.Config{
           FailureRatio:    0.5,             // Permitir 50% de fallos antes de abrir
           MinimumRequests: 10,              // Número mínimo de solicitudes antes de evaluar el estado
           OpenTimeout:     5 * time.Second, // Tiempo de espera en estado abierto
           HalfOpenTimeout: 2 * time.Second, // Tiempo de espera en estado medio abierto
       }
       cb := circuitbreaker.NewMiddleware(cbConfig)

       // Middleware de reintento para hacer la demostración más robusta
       retryConfig := retry.Config{
           Times: 3,                          // Intentar 3 veces
           WaitBase: 500 * time.Millisecond,  // Esperar 500ms entre intentos
       }
       rt := retry.NewMiddleware(retryConfig)

       // Resiliencia compuesta con reintento y circuit breaker
       runner := goresilience.RunnerChain(rt, cb)

       // Simular una función que puede fallar
       err := runner.Run(func() error {
           // Aquí iría la lógica de la solicitud a un servicio externo
           // Simularemos un fallo
           fmt.Println("Intentando realizar la solicitud...")
           return errors.New("servicio no disponible")
       })

       if err != nil {
           fmt.Println("La solicitud falló después de varios intentos:", err)
       } else {
           fmt.Println("La solicitud fue exitosa.")
       }
   }
   ```

### Comportamiento del Circuit Breaker

1. **Configuración del Circuit Breaker**: Se configura con un `FailureRatio` de 0.5, lo que significa que el Circuit Breaker se abrirá si más del 50% de las solicitudes fallan. `MinimumRequests` indica que se necesitan al menos 10 solicitudes para evaluar el estado. `OpenTimeout` y `HalfOpenTimeout` definen los tiempos de espera para cambiar entre estados.

2. **Middleware de Reintento**: Se configura un middleware de reintento para intentar la operación hasta 3 veces con un tiempo de espera de 500ms entre intentos. Esto no es parte del Circuit Breaker, pero se incluye para mostrar cómo manejar fallos de manera más robusta.

3. **Runner Compuesto**: Se compone el runner con el middleware de reintento y el Circuit Breaker, asegurando que ambas capas de resiliencia se apliquen.

4. **Función Simulada**: La función dentro de `runner.Run` simula una solicitud a un servicio que siempre falla. En un escenario real, esta función contendría la lógica para realizar una solicitud HTTP u otro tipo de operación remota.

5. **Ejecución y Manejo de Errores**: Si la solicitud falla después de varios intentos, se imprime un mensaje de error. Si es exitosa, se imprime un mensaje de éxito.

### Detonantes para Abrir el Circuito

El Circuit Breaker se abre cuando detecta un umbral de fallos que indica que el servicio dependiente no está funcionando correctamente. Los detonantes específicos pueden variar según la configuración del Circuit Breaker, pero típicamente incluyen:

1. **Ratio de Fallos Excedido**: Si un porcentaje específico de solicitudes fallan en un período de tiempo, el Circuit Breaker se abrirá. Por ejemplo, si el `FailureRatio` está configurado en 0.5 (50%), y más del 50% de las solicitudes fallan, el Circuit Breaker se abrirá.

2. **Número Mínimo de Solicitudes**: El Circuit Breaker solo comenzará a evaluar el ratio de fallos después de un número mínimo de solicitudes (configurado en `MinimumRequests`). Esto evita que se abra el circuito por unas pocas solicitudes fallidas.

3. **Excepciones Específicas**: En algunos casos, el Circuit Breaker puede configurarse para abrirse cuando se detectan tipos específicos de errores, como errores de red o de tiempo de espera.

### Recuperación del Circuito

Cuando un Circuit Breaker pasa al estado "abierto", evita que las solicitudes lleguen al servicio problemático durante un período de tiempo definido (`OpenTimeout`). Esto da al servicio problemático tiempo para recuperarse sin ser sobrecargado con nuevas solicitudes. Una vez que ha transcurrido el tiempo definido

, el Circuit Breaker pasa al estado "medio abierto" (`Half-Open`).

### Recuperación del Sistema en Estado "Medio Abierto"

En el estado "medio abierto", el Circuit Breaker permite que un número limitado de solicitudes pasen al servicio problemático para verificar si ha vuelto a funcionar correctamente. Si estas solicitudes tienen éxito, el Circuit Breaker pasa al estado "cerrado" (`Closed`), permitiendo que todas las solicitudes fluyan nuevamente al servicio. Si las solicitudes fallan, el Circuit Breaker vuelve al estado "abierto", y el proceso se repite.

### Por Qué se Recupera

El propósito fundamental de abrir un circuit breaker es permitir que el microservicio problemático tenga tiempo para resolver sus problemas sin la carga adicional de nuevas peticiones. Esto podría implicar diversas acciones automatizadas o manuales, como:

1. **Gestión de Recursos**: Liberar recursos que podrían estar agotados, como memoria o conexiones a bases de datos.
2. **Reinicio de Componentes**: Algunos sistemas pueden configurarse para reiniciar servicios o componentes automáticamente.
3. **Reducción de la Carga**: Al no recibir nuevas peticiones, el microservicio puede procesar lo que ya tiene en cola más eficientemente.
4. **Resolución de Problemas Subyacentes**: El tiempo sin recibir nuevas peticiones permite a los operadores del sistema investigar y corregir la causa raíz del problema.


### Ejemplo Detallado con Explicación de Recuperación

Aquí te muestro cómo puedes implementar y observar la recuperación de un Circuit Breaker usando Go con la biblioteca `goresilience`:

1. **Instalación de la biblioteca**:
   ```sh
   go get github.com/slok/goresilience/circuitbreaker
   ```

2. **Código de Ejemplo**:
   ```go
   package main

   import (
       "errors"
       "fmt"
       "time"

       "github.com/slok/goresilience"
       "github.com/slok/goresilience/circuitbreaker"
       "github.com/slok/goresilience/retry"
   )

   var attempt int

   func main() {
       // Configuración del Circuit Breaker
       cbConfig := circuitbreaker.Config{
           FailureRatio:    0.5,             // Permitir 50% de fallos antes de abrir
           MinimumRequests: 2,               // Número mínimo de solicitudes antes de evaluar el estado
           OpenTimeout:     5 * time.Second, // Tiempo de espera en estado abierto
           HalfOpenTimeout: 2 * time.Second, // Tiempo de espera en estado medio abierto
       }
       cb := circuitbreaker.NewMiddleware(cbConfig)

       // Middleware de reintento para hacer la demostración más robusta
       retryConfig := retry.Config{
           Times: 3,                          // Intentar 3 veces
           WaitBase: 500 * time.Millisecond,  // Esperar 500ms entre intentos
       }
       rt := retry.NewMiddleware(retryConfig)

       // Resiliencia compuesta con reintento y circuit breaker
       runner := goresilience.RunnerChain(rt, cb)

       // Realizar varias solicitudes para demostrar la recuperación
       for i := 0; i < 10; i++ {
           err := runner.Run(func() error {
               return simulatedRequest()
           })

           if err != nil {
               fmt.Printf("Intento %d: La solicitud falló: %v\n", i+1, err)
           } else {
               fmt.Printf("Intento %d: La solicitud fue exitosa.\n", i+1)
           }

           time.Sleep(1 * time.Second) // Esperar entre intentos
       }
   }

   // Función simulada que falla inicialmente y luego se recupera
   func simulatedRequest() error {
       attempt++
       if attempt < 5 {
           fmt.Println("Simulando un fallo...")
           return errors.New("servicio no disponible")
       }
       fmt.Println("Simulando una solicitud exitosa...")
       return nil
   }
   ```

### Explicación del Código

1. **Configuración del Circuit Breaker**: Similar al ejemplo anterior, pero con `MinimumRequests` reducido a 2 para demostrar rápidamente el cambio de estado.

2. **Función `simulatedRequest`**: Simula fallos en las primeras cuatro solicitudes y luego simula éxito. Esto permite observar cómo el Circuit Breaker cambia de estado y eventualmente se recupera.

3. **Bucle de Solicitudes**: Se realizan 10 solicitudes, con un segundo de espera entre cada una. Esto permite observar cómo el Circuit Breaker reacciona a los fallos y se recupera cuando la solicitud empieza a tener éxito.

### Comportamiento Esperado

- **Intentos 1-4**: El Circuit Breaker estará en estado "cerrado" y las solicitudes fallarán.
- **Después del 4to fallo**: El Circuit Breaker cambiará a estado "abierto" y bloqueará las solicitudes durante `OpenTimeout`.
- **Después de `OpenTimeout`**: El Circuit Breaker cambiará a estado "medio abierto" y permitirá algunas solicitudes.
- **Si las solicitudes son exitosas en estado "medio abierto"**: El Circuit Breaker cambiará a estado "cerrado" y permitirá todas las solicitudes nuevamente.

### Beneficios del Circuit Breaker

1. **Resiliencia**: Mejora la resiliencia del sistema al aislar los fallos y evitar la propagación de problemas.
2. **Estabilidad**: Mantiene la estabilidad de la aplicación incluso cuando algunas partes del sistema están experimentando problemas.
3. **Rendimiento**: Previene que los fallos de un servicio degraden el rendimiento general del sistema.
4. **Visibilidad**: Proporciona métricas y visibilidad sobre el comportamiento de los servicios y sus fallos.

Este patrón es crucial en arquitecturas de microservicios y sistemas distribuidos donde las dependencias pueden fallar, y es importante que el sistema general continúe funcionando de manera eficiente y estable.