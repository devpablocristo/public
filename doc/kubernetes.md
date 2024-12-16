## Kubernetes

Para reemplazar Consul por Kubernetes en el ejemplo proporcionado, utilizaremos la integración de Kubernetes con `go-micro`. Esto implica cambiar la configuración del registro de servicios para usar Kubernetes en lugar de Consul. Aquí tienes el ejemplo modificado:

### Paso 1: Instalar las Dependencias Necesarias

Asegúrate de tener las dependencias necesarias para utilizar Kubernetes con `go-micro`. Si no las tienes, puedes instalarlas con `go get`:

```sh
go get github.com/micro/go-micro/v2/registry/kubernetes
```

### Paso 2: Configurar el Registro del Servicio con Kubernetes

Modificaremos el archivo `main.go` para utilizar Kubernetes en lugar de Consul:

#### Archivo `main.go`

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/kubernetes"
	"github.com/micro/go-micro/v2/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"path/to/your/project/handler"
	"path/to/your/project/internal/core"
	"path/to/your/project/proto"
)

// Middleware de logging
func loggingWrapper(fn micro.HandlerFunc) micro.HandlerFunc {
	return func(ctx context.Context, req micro.Request, rsp interface{}) error {
		logger.Info("Request received")
		return fn(ctx, req, rsp)
	}
}

// Handler para el EventService
type EventServiceHandler struct {
	notificationClient proto.NotificationService
}

func (h *EventServiceHandler) CreateEvent(ctx context.Context, req *proto.EventRequest, res *proto.EventResponse) error {
	// Implementar la lógica para crear un evento
	fmt.Printf("Creating event: %s\n", req.Title)

	// Simular la creación del evento
	res.Message = "Event created successfully"

	// Enviar notificación después de crear el evento
	notificationReq := &proto.NotificationRequest{
		Message: "New event created: " + req.Title,
	}

	notificationRes, err := h.notificationClient.SendNotification(ctx, notificationReq)
	if err != nil {
		res.Err = err.Error()
		return err
	}
	fmt.Println("Notification Response:", notificationRes.Status)

	return nil
}

func main() {
	// Configurar el registro de servicios con Kubernetes
	reg := kubernetes.NewRegistry()

	// Crear un nuevo servicio con go-micro
	service := micro.NewService(
		micro.Name("example.service"),
		micro.Registry(reg),
		micro.WrapHandler(loggingWrapper), // Añadir un middleware de logging
	)

	// Inicializar el servicio
	service.Init(
		micro.AfterStart(func() error {
			fmt.Println("Service started")
			return nil
		}),
		micro.BeforeStop(func() error {
			fmt.Println("Service stopping")
			return nil
		}),
	)

	// Crear el cliente para NotificationService
	notificationClient := proto.NewNotificationService("notification.service", service.Client())

	// Crear el handler del EventService
	eventHandler := &EventServiceHandler{
		notificationClient: notificationClient,
	}

	// Registrar el handler del EventService
	proto.RegisterEventServiceHandler(service.Server(), eventHandler)

	// Crear el servicio web para exponer el REST API
	webService := web.NewService(
		web.Name("event.web"),
		web.Registry(reg),
		web.Address(":8081"),
	)

	// Inicializar el servicio web
	webService.Init()

	// Crear el UseCasePort
	useCasePort := core.NewUseCasePort()

	// Crear el handler REST
	restHandler := handler.NewRestHandler(useCasePort)

	// Configurar el router Gin
	r := gin.Default()
	r.POST("/event", restHandler.CreateEvent)
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // Endpoint para métricas

	// Registrar el handler REST con go-micro
	webService.Handle("/", r)

	// Ejecutar ambos servicios
	go func() {
		if err := service.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	if err := webService.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

### Explicación del Cambio

1. **Configuración del Registro de Servicios con Kubernetes**:
   ```go
   reg := kubernetes.NewRegistry()
   ```
   - En lugar de utilizar Consul, configuramos el registro de servicios para usar Kubernetes. `kubernetes.NewRegistry()` crea un nuevo registro que utiliza la API de Kubernetes para registrar y descubrir servicios.

2. **Resto del Código**:
   - El resto del código permanece igual, ya que `go-micro` maneja el registro y descubrimiento de servicios de manera similar independientemente del backend de registro utilizado.

### Configurar Kubernetes

Para que el servicio y el cliente funcionen correctamente en Kubernetes, necesitas definir los recursos de Kubernetes (Deployment y Service) para ambos.

#### Archivo `k8s/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: event-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: event-service
  template:
    metadata:
      labels:
        app: event-service
    spec:
      containers:
      - name: event-service
        image: your-docker-repo/event-service:latest
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: event-service
spec:
  selector:
    app: event-service
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080
  type: ClusterIP
```

### Conclusión

Este ejemplo demuestra cómo reemplazar Consul por Kubernetes para el registro y descubrimiento de servicios en `go-micro`. Al utilizar `kubernetes.NewRegistry()`, podemos aprovechar las capacidades nativas de Kubernetes para gestionar nuestros servicios. Este enfoque permite que los servicios se registren y descubran automáticamente utilizando la infraestructura de Kubernetes, lo que simplifica la gestión de servicios en un entorno de contenedores.


### Consul y Kubernetes

Consul y Kubernetes tienen algunas funcionalidades superpuestas, especialmente en lo que respecta al descubrimiento de servicios y la gestión de configuraciones, pero son herramientas diseñadas con propósitos y enfoques diferentes. Aquí hay una comparación de las dos:

### Kubernetes

Kubernetes es una plataforma de orquestación de contenedores que automatiza el despliegue, la escalabilidad y la gestión de aplicaciones en contenedores. Algunas de sus características principales incluyen:

1. **Orquestación de Contenedores**:
   - Gestiona el ciclo de vida de los contenedores, incluyendo despliegues, actualizaciones y reinicios automáticos.
   - Facilita el escalado automático de aplicaciones basado en la carga.

2. **Descubrimiento de Servicios**:
   - Proporciona mecanismos internos para el descubrimiento de servicios utilizando DNS y configuraciones de servicio.
   - Cada servicio en Kubernetes puede ser descubierto y accedido a través de un nombre DNS.

3. **Balanceo de Carga**:
   - Distribuye automáticamente el tráfico de red entrante a las instancias de los contenedores para equilibrar la carga.

4. **Gestión de Configuraciones y Secretos**:
   - Facilita la gestión de configuraciones y secretos de manera segura usando ConfigMaps y Secrets.

5. **Autoscaling**:
   - Ajusta automáticamente el número de réplicas de las aplicaciones basadas en métricas de rendimiento.

6. **Resiliencia y Recuperación**:
   - Reemplaza y reinicia contenedores que fallan, programa y reschedule contenedores cuando los nodos fallan.

### Consul

Consul es una solución de HashiCorp diseñada para facilitar el descubrimiento y la configuración de servicios en entornos dinámicos y distribuidos. Sus características principales incluyen:

1. **Descubrimiento de Servicios**:
   - Proporciona un registro de servicios y un mecanismo de descubrimiento de servicios mediante APIs HTTP y DNS.
   - Permite a los servicios registrarse y descubrir otros servicios de manera dinámica.

2. **Balanceo de Carga**:
   - Puede trabajar con balanceadores de carga externos para distribuir el tráfico entre las instancias del servicio.
   - Integración con proxies como Envoy para realizar balanceo de carga inteligente.

3. **Gestión de Configuraciones**:
   - Almacena y gestiona configuraciones dinámicas mediante su Key-Value Store.

4. **Health Checks**:
   - Realiza verificaciones de salud configurables para determinar el estado de los servicios y los nodos.
   - Desregistrar servicios automáticamente si fallan los checks de salud.

5. **Redes de Servicios (Service Mesh)**:
   - Con Consul Connect, ofrece capacidades de malla de servicios para gestionar la comunicación segura entre servicios.
   - Proporciona mTLS (Mutual TLS) para la comunicación segura entre servicios y políticas de autorización.

6. **Integración con Herramientas de Orquestación**:
   - Aunque Consul no orquesta contenedores, se integra bien con herramientas como Kubernetes para proporcionar descubrimiento de servicios, configuración y gestión de redes de servicios.

### Diferencias Clave

1. **Enfoque**:
   - **Kubernetes** se centra en la orquestación de contenedores, gestionando el despliegue, escalabilidad y administración de aplicaciones contenedorizadas.
   - **Consul** se centra en el descubrimiento de servicios, la configuración y la gestión de redes de servicios.

2. **Orquestación de Contenedores**:
   - Kubernetes maneja directamente la orquestación de contenedores, programando, escalando y gestionando el ciclo de vida de los contenedores.
   - Consul no orquesta contenedores, pero puede integrarse con herramientas de orquestación como Kubernetes o Nomad.

3. **Gestión de Redes de Servicios**:
   - Kubernetes ofrece descubrimiento de servicios y balanceo de carga dentro del clúster utilizando DNS y otros mecanismos internos.
   - Consul ofrece capacidades de descubrimiento de servicios más flexibles y avanzadas, incluyendo malla de servicios con Consul Connect.

### Usos Combinados

Muchas veces, Consul y Kubernetes se usan juntos para aprovechar lo mejor de ambos mundos:

- **Kubernetes** se utiliza para la orquestación de contenedores y gestión de aplicaciones.
- **Consul** se utiliza para el descubrimiento de servicios, la configuración dinámica y la malla de servicios.

## Ejemplo

Cómo configurar y desplegar una aplicación Go que usa `go-micro` y `gin`, con Consul para el descubrimiento de servicios y Kubernetes para la orquestación. Además, hemos añadido un middleware y un endpoint de métricas con Prometheus.

### Archivo `main.go`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/consul"
	"github.com/micro/go-micro/v2/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"path/to/your/project/handler"
	"path/to/your/project/internal/core"
)

// Middleware de ejemplo para logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Registrar la solicitud entrante
		startTime := time.Now()
		logger.Infof("Incoming request: %s %s", c.Request.Method, c.Request.URL)

		// Procesar la solicitud
		c.Next()

		// Registrar la respuesta saliente
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()
		logger.Infof("Response: %d, Latency: %v", statusCode, latency)
	}
}

func main() {
	// Configurar el registro de servicios con Consul
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"consul-server.consul:8500"} // Dirección del servidor Consul en Kubernetes
	})

	// Crear un nuevo servicio web con go-micro
	service := web.NewService(
		web.Name("event.service"),
		web.Version("latest"),
		web.Registry(reg),
		web.Address(":8081"), // Especificar el puerto en el que este servicio escuchará
	)

	// Inicializar el servicio
	if err := service.Init(); err != nil {
		logger.Fatal(err)
	}

	// Crear el UseCasePort
	useCasePort := core.NewUseCasePort()

	// Crear el handler REST
	restHandler := handler.NewRestHandler(useCasePort)

	// Configurar el router Gin
	r := gin.Default()
	r.Use(loggingMiddleware()) // Añadir middleware de logging
	r.POST("/event", restHandler.CreateEvent)

	// Añadir endpoint para métricas
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Registrar el handler REST con go-micro
	service.Handle("/", r)

	// Ejecutar el servicio
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

### Archivo `Dockerfile`

```dockerfile
# Dockerfile
FROM golang:1.16

WORKDIR /app

COPY . .

RUN go build -o event-service ./cmd

EXPOSE 8081

CMD ["./event-service"]
```

### Manifiestos de Kubernetes

#### Archivo `k8s/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: event-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: event-service
  template:
    metadata:
      labels:
        app: event-service
    spec:
      containers:
      - name: event-service
        image: your-docker-repo/event-service:latest
        env:
        - name: CONSUL_HTTP_ADDR
          value: "consul-server.consul:8500" # Dirección del servidor Consul
        ports:
        - containerPort: 8081
---
apiVersion: v1
kind: Service
metadata:
  name: event-service
spec:
  selector:
    app: event-service
  ports:
    - protocol: TCP
      port: 8081
      targetPort: 8081
```

### Despliegue de Consul en Kubernetes

Si aún no has desplegado Consul en Kubernetes, puedes hacerlo con Helm:

```sh
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install consul hashicorp/consul --set global.name=consul
```

### Desplegar los Recursos en Kubernetes

Aplica los manifiestos de Kubernetes para desplegar tu aplicación:

```sh
kubectl apply -f k8s/deployment.yaml
```

### Construir y Publicar la Imagen Docker

Construye y publica la imagen Docker de tu aplicación:

```sh
docker build -t your-docker-repo/event-service:latest .
docker push your-docker-repo/event-service:latest
```

### Conclusión

En este ejemplo, hemos configurado un servicio en Go utilizando `go-micro` y `gin`, con Consul para el descubrimiento de servicios y Kubernetes para la orquestación. Hemos añadido un middleware de logging a la aplicación y un endpoint de métricas utilizando Prometheus. Además, hemos preparado los archivos de configuración necesarios para desplegar tanto la aplicación como Consul en un clúster de Kubernetes. Esta configuración permite aprovechar lo mejor de Kubernetes y Consul sin necesidad de modificar directamente el código de la aplicación para la integración básica.