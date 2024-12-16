# Presentación sobre Circuit Breaker

## Introducción al Circuit Breaker

- **Definición**: El patrón Circuit Breaker (Interruptor de Circuito) es un diseño utilizado en sistemas distribuidos para prevenir que fallos en un servicio o dependencia se propaguen a otros servicios y evitar cascadas de errores.
- **Funcionamiento**: Similar a un interruptor eléctrico, se abre para interrumpir el flujo cuando detecta fallos y se cierra de nuevo cuando el sistema se recupera.

---

## Estados del Circuit Breaker

### 1. Closed (Cerrado)

- **Operación Normal**: Todas las solicitudes pasan al servicio de destino.
- **Monitoreo de Errores**: Se rastrean errores y tiempos de respuesta.
- **Umbral**: Si los errores superan un umbral predefinido, el circuito se abre.

### 2. Open (Abierto)

- **Bloqueo**: Las solicitudes se bloquean inmediatamente y se devuelve una respuesta de error o se ejecuta un método de fallback.
- **Período de Enfriamiento**: Después de un tiempo, el circuito entra en un estado de prueba.

### 3. Half-Open (Medio Abierto)

- **Solicitudes de Prueba**: Permite pasar un número limitado de solicitudes para probar si el servicio se ha recuperado.
- **Recuperación**: Si las solicitudes son exitosas, el circuito se cierra.
- **Fallo**: Si las solicitudes fallan, el circuito se abre nuevamente.

---

## Funcionamiento del Circuit Breaker

1. **Monitoreo de Solicitudes**: Mientras el circuito está cerrado, se monitorean las solicitudes para detectar fallos y tiempos de respuesta.
2. **Apertura del Circuito**: Si los fallos superan el umbral definido, el circuito se abre. Las solicitudes posteriores reciben una respuesta de error o se ejecuta un método de fallback.
3. **Período de Enfriamiento**: Durante este tiempo, no se envían solicitudes al servicio fallido.
4. **Reintento de Solicitudes**: Después del período de enfriamiento, el circuito pasa a un estado medio abierto y permite que un número limitado de solicitudes pasen.
5. **Cierre del Circuito**: Si las solicitudes de prueba son exitosas, el circuito se cierra y el flujo normal de solicitudes se reanuda. Si fallan, el circuito se abre nuevamente.

---

## Ejemplo de Implementación en Go con Hystrix

1. **Instalación de Hystrix**:
   ```sh
   go get github.com/afex/hystrix-go/hystrix
   ```

2. **Código de Ejemplo**:
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
           Timeout:               1000,
           MaxConcurrentRequests: 100,
           ErrorPercentThreshold: 25,
       })

       // Crea un nuevo servicio
       service := micro.NewService(
           micro.Name("example.service"),
       )
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

---

## Ejemplo en Go con `goresilience`

1. **Instalación de la Biblioteca**:
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
           FailureRatio:    0.5,
           MinimumRequests: 10,
           OpenTimeout:     5 * time.Second,
           HalfOpenTimeout: 2 * time.Second,
       }
       cb := circuitbreaker.NewMiddleware(cbConfig)

       // Middleware de reintento
       retryConfig := retry.Config{
           Times: 3,
           WaitBase: 500 * time.Millisecond,
       }
       rt := retry.NewMiddleware(retryConfig)

       // Resiliencia compuesta
       runner := goresilience.RunnerChain(rt, cb)

       // Simular una función que puede fallar
       err := runner.Run(func() error {
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

---

## Comportamiento del Circuit Breaker

1. **Configuración**: Se configura con un `FailureRatio` de 0.5, `MinimumRequests` de 10, `OpenTimeout` de 5 segundos y `HalfOpenTimeout` de 2 segundos.
2. **Middleware de Reintento**: Se configura para intentar la operación hasta 3 veces con 500ms de espera entre intentos.
3. **Runner Compuesto**: Combina el middleware de reintento y el Circuit Breaker para mayor resiliencia.
4. **Función Simulada**: Simula una solicitud a un servicio que siempre falla inicialmente.
5. **Ejecución y Manejo de Errores**: Muestra mensajes según el resultado de las solicitudes.

---

## Detonantes para Abrir el Circuito

1. **Ratio de Fallos Excedido**: Si un porcentaje específico de solicitudes fallan, el Circuit Breaker se abrirá. Por ejemplo, si `FailureRatio` es 0.5 (50%), y más del 50% de las solicitudes fallan, el Circuit Breaker se abrirá.
2. **Número Mínimo de Solicitudes**: El Circuit Breaker comenzará a evaluar el ratio de fallos después de un número mínimo de solicitudes (`MinimumRequests`).
3. **Excepciones Específicas**: El Circuit Breaker puede configurarse para abrirse cuando se detectan tipos específicos de errores, como errores de red o de tiempo de espera.

---

## Recuperación del Circuito

- **Estado "Abierto"**: El Circuit Breaker evita que las solicitudes lleguen al servicio problemático durante un período de tiempo definido (`OpenTimeout`), permitiendo la recuperación del servicio sin sobrecargas adicionales.
- **Estado "Medio Abierto"**: Permite un número limitado de solicitudes para verificar si el servicio ha vuelto a funcionar correctamente. Si son exitosas, el circuito se cierra. Si fallan, el circuito vuelve al estado "abierto".

---

## Ejemplo de Recuperación en Go con `goresilience`

1. **Código de Ejemplo**:
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
           FailureRatio:    0.5,
           MinimumRequests: 2,
           OpenTimeout:     5 * time.Second,
           HalfOpenTimeout: 2 * time.Second,
       }
       cb := circuitbreaker.NewMiddleware(cbConfig)

       // Middleware de reintento
       retryConfig := retry.Config{
           Times: 3,
           WaitBase: 500 * time.Millisecond,
       }
       rt := retry.NewMiddleware(retryConfig)

       // Resiliencia compuesta
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
---

## Beneficios del Circuit Breaker

1. **Resiliencia**: Mejora la resiliencia del sistema al aislar los fallos y evitar la propagación de problemas.
2. **Estabilidad**: Mantiene la estabilidad de la aplicación incluso cuando algunas partes del sistema están experimentando problemas.
3. **Rendimiento**: Previene que los fallos de un servicio degraden el rendimiento general del sistema.
4. **Visibilidad**: Proporciona métricas y visibilidad sobre el comportamiento de los servicios y sus fallos.

---

## Recuperación del Circuito

El propósito fundamental de abrir un circuit breaker es permitir que el microservicio problemático tenga tiempo para resolver sus problemas sin la carga adicional de nuevas peticiones. Esto podría implicar diversas acciones automatizadas o manuales, tales como:

- **Gestión de Recursos**: Liberación de recursos que podrían estar agotados, como memoria o conexiones a bases de datos.
- **Reinicio de Componentes**: Algunos sistemas pueden configurarse para reiniciar servicios o componentes automáticamente.
- **Reducción de la Carga**: Al no recibir nuevas peticiones, el microservicio puede procesar lo que ya tiene en cola más eficientemente.
- **Resolución de Problemas Subyacentes**: El tiempo sin recibir nuevas peticiones permite a los operadores del sistema investigar y corregir la causa raíz del problema.

### Estado "Abierto"

- **Prevención de Sobrecarga**: El Circuit Breaker evita que las solicitudes lleguen al servicio problemático durante un período de tiempo definido (`OpenTimeout`), permitiendo la recuperación del servicio sin sobrecargas adicionales.

### Estado "Medio Abierto"

- **Prueba de Recuperación**: Permite un número limitado de solicitudes para verificar si el servicio ha vuelto a funcionar correctamente. Si estas solicitudes son exitosas, el circuito se cierra. Si fallan, el circuito vuelve al estado "abierto".

---

## Beneficios del Circuit Breaker

1. **Resiliencia**: Mejora la resiliencia del sistema al aislar los fallos y evitar la propagación de problemas.
2. **Estabilidad**: Mantiene la estabilidad de la aplicación incluso cuando algunas partes del sistema están experimentando problemas.
3. **Rendimiento**: Previene que los fallos de un servicio degraden el rendimiento general del sistema.
4. **Visibilidad**: Proporciona métricas y visibilidad sobre el comportamiento de los servicios y sus fallos.

---

## Conclusión

El patrón Circuit Breaker es esencial en arquitecturas de microservicios y sistemas distribuidos para mantener la eficiencia y estabilidad del sistema. Implementar un Circuit Breaker mejora la resiliencia, estabilidad y rendimiento de la aplicación al manejar de manera efectiva los fallos de los servicios dependientes.

---

