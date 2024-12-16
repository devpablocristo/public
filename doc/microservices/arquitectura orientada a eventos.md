Una arquitectura orientada a eventos es un enfoque de diseño de software en el que los componentes del sistema se comunican entre sí mediante la producción y el consumo de eventos. En lugar de realizar llamadas directas entre servicios, los eventos se generan y se publican en un sistema de mensajería, y otros servicios los consumen para realizar acciones en respuesta.

**Características**:
- **Desacoplamiento**: Los servicios no necesitan conocer la ubicación ni la implementación de otros servicios. Solo necesitan conocer el formato del evento que deben producir o consumir.
- **Asincronía**: La comunicación se realiza de forma asíncrona, permitiendo que los servicios continúen su trabajo sin esperar a que otros servicios respondan.
- **Escalabilidad**: Facilita la escalabilidad horizontal, ya que los servicios pueden ser escalados independientemente según la carga de trabajo.
- **Resiliencia**: Mejora la resiliencia del sistema, ya que los servicios pueden seguir funcionando incluso si otros servicios están temporalmente inactivos.
- **Flexibilidad**: Permite añadir o modificar funcionalidades sin afectar significativamente a otros componentes del sistema.
- **Tolerancia a Fallos**: La arquitectura puede manejar fallos de manera más robusta, ya que los mensajes se pueden almacenar y reintentar en caso de fallos temporales.

### Patrones Comunes en Arquitecturas Orientadas a Eventos

1. **Pub/Sub (Publicación/Suscripción)**:
   - Los productores publican eventos en un topic y los consumidores se suscriben a esos topics para recibir los eventos.

2. **Event Sourcing**:
   - En lugar de almacenar el estado actual del sistema, se almacenan todos los eventos que cambiaron el estado. El estado actual se puede reconstruir reproduciendo estos eventos.

3. **CQRS (Command Query Responsibility Segregation)**:
   - Separa las operaciones de lectura y escritura en diferentes modelos, utilizando eventos para mantener el modelo de consulta actualizado.

### Ejemplo de Arquitectura Orientada a Eventos

Supongamos que estamos construyendo un sistema de comercio electrónico con una arquitectura orientada a eventos. Los servicios clave pueden incluir:

- **Order Service**: Gestiona los pedidos.
- **Inventory Service**: Gestiona el inventario.
- **Notification Service**: Envía notificaciones a los usuarios.
- **Payment Service**: Gestiona los pagos.

#### Diagrama de Arquitectura

```
Cliente -> Order Service -> Message Broker (Kafka, RabbitMQ, etc.)
                           -> Inventory Service (suscriptor)
                           -> Notification Service (suscriptor)
                           -> Payment Service (suscriptor)
```

#### Implementación en Go con RabbitMQ

##### Productor de Eventos (Order Service)

**Archivo `main.go` del Order Service**:

```go
package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/streadway/amqp"
)

func main() {
    // Conectar a RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-service:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    // Configurar el router Gin
    r := gin.Default()

    r.POST("/order", func(c *gin.Context) {
        // Lógica para crear un pedido
        order := map[string]string{"id": "1234", "status": "created"}

        // Publicar un mensaje en RabbitMQ
        body := "New order created with ID 1234"
        err = ch.Publish(
            "",            // exchange
            "orderQueue",  // routing key
            false,         // mandatory
            false,         // immediate
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(body),
            })
        if err != nil {
            log.Fatalf("Failed to publish a message: %v", err)
        }

        c.JSON(http.StatusCreated, order)
    })

    r.Run(":8080")
}
```

##### Consumidor de Eventos (Inventory Service)

**Archivo `main.go` del Inventory Service**:

```go
package main

import (
    "log"
    "github.com/streadway/amqp"
)

func main() {
    // Conectar a RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-service:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    msgs, err := ch.Consume(
        "orderQueue", // queue
        "",           // consumer
        true,         // auto-ack
        false,        // exclusive
        false,        // no-local
        false,        // no-wait
        nil,          // args
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    forever := make(chan bool)

    go func() {
        for d := range msgs {
            log.Printf("Received a message: %s", d.Body)
            // Lógica para actualizar el inventario
        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}
```

### Despliegue en Kubernetes

#### Archivo `k8s/deployment.yaml` para RabbitMQ y los servicios

```yaml
apiVersion: v1
kind: Service
metadata:
  name: rabbitmq-service
spec:
  ports:
  - port: 5672
    targetPort: 5672
  selector:
    app: rabbitmq
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rabbitmq
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rabbitmq
  template:
    metadata:
      labels:
        app: rabbitmq
    spec:
      containers:
      - name: rabbitmq
        image: rabbitmq:3-management
        ports:
        - containerPort: 5672
        - containerPort: 15672
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
      - name: order-service
        image: your-docker-repo/order-service:latest
        ports:
        - containerPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inventory-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: inventory-service
  template:
    metadata:
      labels:
        app: inventory-service
    spec:
      containers:
      - name: inventory-service
       

 image: your-docker-repo/inventory-service:latest
```

### Construir y Publicar la Imagen Docker

Construye y publica las imágenes Docker de tus servicios:

```sh
docker build -t your-docker-repo/order-service:latest .
docker push your-docker-repo/order-service:latest

docker build -t your-docker-repo/inventory-service:latest .
docker push your-docker-repo/inventory-service:latest
```

### Desplegar los Recursos en Kubernetes

Aplica los manifiestos de Kubernetes para desplegar tus servicios:

```sh
kubectl apply -f k8s/deployment.yaml
```

### Despliegue de Prometheus, Grafana e Istio

#### Despliegue de Prometheus y Grafana

Usa Helm para desplegar Prometheus y Grafana:

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Desplegar Prometheus
helm install prometheus prometheus-community/prometheus

# Desplegar Grafana
helm install grafana grafana/grafana --set adminPassword='YourPassword' --set service.type=NodePort
```

#### Despliegue de Istio

Instala Istio CLI y despliega Istio en tu clúster de Kubernetes:

```sh
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.9.0
export PATH=$PWD/bin:$PATH
istioctl install --set profile=demo -y
```

Etiquetar el namespace para la inyección automática de Istio:

```sh
kubectl label namespace default istio-injection=enabled
```

### Despliegue de Consul en Kubernetes

Despliega Consul en Kubernetes usando Helm:

```sh
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install consul hashicorp/consul --set global.name=consul
```

### Construir y Publicar la Imagen Docker

Construye y publica la imagen Docker de tu aplicación:

```sh
docker build -t your-docker-repo/user-service:latest .
docker push your-docker-repo/user-service:latest
```