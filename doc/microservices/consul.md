## Consul

**Consul** es una herramienta de código abierto desarrollada por HashiCorp que proporciona una solución completa para el descubrimiento de servicios, la gestión de configuraciones y la segmentación de red en un entorno distribuido. Consul está diseñado para ser altamente disponible, escalable y fácil de integrar con diversas aplicaciones y entornos de microservicios.

### Características Clave de Consul

1. **Descubrimiento de Servicios**:
   - Consul permite que los servicios se registren y se descubran entre sí. Los servicios pueden encontrar otros servicios sin necesidad de conocer sus ubicaciones exactas (IP y puertos).

2. **Gestión de Configuraciones**:
   - Consul proporciona una base de datos KV (key-value) distribuida que se puede usar para almacenar configuraciones dinámicas que los servicios pueden leer y reaccionar en tiempo real.

3. **Segmentación de Red (Service Segmentation)**:
   - Consul puede gestionar políticas de red y aplicar malla de servicios (service mesh) para controlar y asegurar la comunicación entre servicios.

4. **Supervisión y Salud de Servicios**:
   - Consul realiza comprobaciones de salud para asegurarse de que los servicios estén funcionando correctamente. Los servicios pueden registrarse con comprobaciones de salud que Consul ejecutará periódicamente.

5. **Multi-Datacenter**:
   - Consul es capaz de operar en múltiples centros de datos, proporcionando descubrimiento de servicios y configuración de clave-valor a través de diferentes ubicaciones geográficas.

### Componentes Principales de Consul

1. **Agentes**:
   - **Agentes de Cliente**: Se ejecutan en cada nodo donde los servicios están registrados. Se encargan de registrar servicios locales y ejecutar comprobaciones de salud.
   - **Agentes de Servidor**: Mantienen el estado del clúster, gestionan la base de datos KV y coordinan el descubrimiento de servicios y las operaciones de segmentación de red.

2. **Catálogo de Servicios**:
   - Un registro de todos los servicios registrados y sus instancias, junto con el estado de sus comprobaciones de salud.

3. **Base de Datos KV**:
   - Almacena configuraciones y otros datos que los servicios pueden necesitar para su funcionamiento.

4. **Comprobaciones de Salud**:
   - Verifican periódicamente el estado de los servicios para asegurarse de que están operativos.

5. **Interfaz de Usuario y API**:
   - Consul proporciona una interfaz web para visualizar el estado del clúster y una API HTTP para interactuar programáticamente con el sistema.

### Ejemplo de Uso de Consul en una Arquitectura de Microservicios

Supongamos que tenemos una aplicación de comercio electrónico con varios microservicios:

- **User Service**
- **Product Service**
- **Order Service**
- **Payment Service**

Estos servicios necesitan descubrirse entre sí y gestionar sus configuraciones dinámicamente.

#### 1. Despliegue de Consul en Kubernetes

Podemos usar Helm para desplegar Consul en un clúster de Kubernetes.

```sh
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install consul hashicorp/consul --set global.name=consul
```

#### 2. Configuración de Consul para Descubrimiento de Servicios

##### Archivo `consul-agent-config.json`:

```json
{
  "datacenter": "dc1",
  "data_dir": "/opt/consul",
  "log_level": "INFO",
  "node_name": "consul-server",
  "server": true,
  "bootstrap_expect": 1,
  "ui": true,
  "bind_addr": "0.0.0.0",
  "client_addr": "0.0.0.0",
  "advertise_addr": "<node-ip>",
  "retry_join": ["<ip-of-another-consul-server>"]
}
```

#### 3. Configuración de un Servicio para Registrar en Consul

Supongamos que tenemos un servicio de usuarios escrito en Go que necesita registrarse en Consul.

**Archivo `main.go` del User Service**:

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hashicorp/consul/api"
)

func main() {
    // Configurar Consul
    config := api.DefaultConfig()
    config.Address = "consul-server.consul:8500"
    client, err := api.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create Consul client: %v", err)
    }

    // Registrar el servicio en Consul
    registration := &api.AgentServiceRegistration{
        ID:      "user-service",
        Name:    "user-service",
        Address: "localhost",
        Port:    8081,
        Check: &api.AgentServiceCheck{
            HTTP:     "http://localhost:8081/health",
            Interval: "10s",
        },
    }

    err = client.Agent().ServiceRegister(registration)
    if err != nil {
        log.Fatalf("Failed to register service with Consul: %v", err)
    }

    // Configurar el router Gin
    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    r.GET("/users", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get Users"})
    })

    r.Run(":8081")
}
```

#### 4. Despliegue del User Service en Kubernetes

##### Archivo `k8s/user-service-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: your-docker-repo/user-service:latest
        ports:
        - containerPort: 8081
        env:
        - name: CONSUL_HTTP_ADDR
          value: "consul-server.consul:8500" # Dirección del servidor Consul
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
    - protocol: TCP
      port: 8081
      targetPort: 8081
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

### Desplegar los Recursos en Kubernetes

Aplica los manifiestos de Kubernetes para desplegar tu aplicación:

```sh
kubectl apply -f k8s/user-service-deployment.yaml
```

### Construir y Publicar la Imagen Docker

Construye y publica la imagen Docker de tu aplicación:

```sh
docker build -t your-docker-repo/user-service:latest .
docker push your-docker-repo/user-service:latest
```

### Registrar y Descubrir Servicios

Necesitas tener un servidor de Consul en funcionamiento para que los servicios puedan registrarse y descubrirse. Consul puede ejecutarse localmente en tu máquina de desarrollo o en un servidor dedicado en tu entorno de producción. Aquí hay una guía para configurar y levantar Consul:

### Paso 1: Instalar Consul

#### En macOS:

Puedes instalar Consul usando Homebrew:

```sh
brew install consul
```

#### En Linux:

Descarga el binario desde el sitio oficial de Consul y sigue las instrucciones de instalación:

```sh
wget https://releases.hashicorp.com/consul/1.10.4/consul_1.10.4_linux_amd64.zip
unzip consul_1.10.4_linux_amd64.zip
sudo mv consul /usr/local/bin/
```

##### Debian/Ubuntu/etc
```sh
sudo apt install consul 
```

#### En Windows:

Descarga el binario desde el sitio oficial de Consul y sigue las instrucciones de instalación.


#### docker-compose:

Si esta definido en docker-compose, no hace falta instalarlo localmente.

### Paso 2: Iniciar Consul

Puedes iniciar un servidor de Consul en modo desarrollo con el siguiente comando:

```sh
consul agent -dev
```

Este comando inicia un agente de Consul en modo desarrollo en tu máquina local, escuchando en `127.0.0.1:8500`.

### Paso 3: Verificar Consul

Una vez que Consul esté en funcionamiento, puedes verificar que está corriendo accediendo a la interfaz web en tu navegador:

```
http://localhost:8500
```

Esto te llevará a la UI de Consul, donde puedes ver los servicios registrados, los nodos y otros detalles.

### Paso 4: Configurar y Ejecutar los Servicios

Con Consul en funcionamiento, ahora puedes ejecutar los servicios `EventService` y `NotificationService` para que se registren en Consul y puedan ser descubiertos.

#### EventService

##### Archivo `main.go`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/consul"
	"github.com/micro/go-micro/v2/web"
	"path/to/your/project/handler"
	"path/to/your/project/internal/core"
)

func main() {
	// Configurar el registro de servicios con Consul
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
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
	r.POST("/event", restHandler.CreateEvent)

	// Registrar el handler REST con go-micro
	service.Handle("/", r)

	// Ejecutar el servicio
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

#### NotificationService

##### Archivo `main.go`

```go
package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/consul"
	"path/to/your/project/proto"
)

type NotificationService struct{}

func (s *NotificationService) SendNotification(ctx context.Context, req *proto.NotificationRequest, res *proto.NotificationResponse) error {
	// Implementar la lógica para enviar una notificación
	fmt.Printf("Sending notification: %s\n", req.Message)
	res.Status = "Notification sent successfully"
	return nil
}

func main() {
	// Configurar el registro de servicios con Consul
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
	})

	// Crear un nuevo servicio en un puerto específico
	service := micro.NewService(
		micro.Name("notification.service"),
		micro.Version("latest"),
		micro.Registry(reg),
		micro.Address(":8082"), // Especificar el puerto en el que este servicio escuchará
	)

	// Inicializar el servicio
	service.Init()

	// Registrar el handler del servicio
	proto.RegisterNotificationServiceHandler(service.Server(), new(NotificationService))

	// Ejecutar el servicio
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
```

### Paso 5: Ejecutar el Cliente

Finalmente, puedes ejecutar el cliente que consumirá ambos servicios:

##### Archivo `client/main.go`

```go
package main

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/consul"
	"path/to/your/project/proto"
)

func main() {
	// Configurar el registro de servicios con Consul
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{"127.0.0.1:8500"}
	})

	// Crear un nuevo servicio cliente
	service := micro.NewService(
		micro.Name("client.service"),
		micro.Registry(reg),
	)
	service.Init()

	// Crear clientes para los servicios Event y Notification
	eventClient := proto.NewEventService("event.service", service.Client())
	notificationClient := proto.NewNotificationService("notification.service", service.Client())

	// Crear un evento
	eventReq := &proto.EventRequest{
		Id:          "1",
		Title:       "Sample Event",
		Description: "This is a sample event",
	}

	eventRes, err := eventClient.CreateEvent(context.Background(), eventReq)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("Event Response:", eventRes.Message)

	// Enviar una notificación
	notificationReq := &proto.NotificationRequest{
		Message: "New event created: " + eventReq.Title,
	}

	notificationRes, err := notificationClient.SendNotification(context.Background(), notificationReq)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println("Notification Response:", notificationRes.Status)
}
```

### Conclusión

Al configurar y levantar un servidor de Consul, tus microservicios pueden registrarse y ser descubiertos dinámicamente, facilitando la comunicación y escalabilidad en una arquitectura de microservicios. Cada servicio escucha en su propio puerto, y Consul se encarga de mantener la información sobre estos servicios y sus ubicaciones, permitiendo a los clientes encontrar y comunicarse con los servicios necesarios sin configuraciones estáticas.


Sure, I'll translate the provided text into Spanish:

---

### Guía práctica de HashiCorp Consul — Parte 1

Velotio Technologies
Velotio Perspectives

Publicado en Velotio Perspectives
18 min de lectura
16 de abril de 2019

Esta es la parte 1 de una serie de 2 partes sobre una guía práctica de HashiCorp Consul. Esta parte se centra principalmente en entender los problemas que Consul resuelve y cómo los resuelve. La segunda parte se centrará más en una aplicación práctica de Consul en un ejemplo de la vida real y se publicará la próxima semana. Empecemos.

¿Qué tal configurar una malla de servicios descubierta, configurable y segura utilizando una sola herramienta?

¿Qué pasaría si te decimos que esta herramienta es independiente de la plataforma y está lista para la nube?

Y viene como una descarga binaria única.

Todo esto es cierto. La herramienta de la que hablamos es HashiCorp Consul.

Consul proporciona descubrimiento de servicios, comprobaciones de salud, balanceo de carga, gráfico de servicios, aplicación de identidad a través de TLS y gestión de configuración de servicios distribuidos.

Vamos a conocer Consul en detalle a continuación y ver cómo resuelve estos complejos desafíos y facilita la vida de un operador de sistemas distribuidos.

### Introducción

Los microservicios y otros sistemas distribuidos pueden permitir un desarrollo de software más rápido y sencillo. Pero hay un inconveniente que resulta en una mayor complejidad operativa en torno a la comunicación entre servicios, la gestión de la configuración y la segmentación de la red.

HashiCorp Consul es una herramienta de código abierto que resuelve estas nuevas complejidades al proporcionar descubrimiento de servicios, comprobaciones de salud, balanceo de carga, un gráfico de servicios, aplicación de identidad mutua TLS y un almacén de configuración de clave-valor. Estas características hacen de Consul un plano de control ideal para una malla de servicios.

HashiCorp anunció Consul en abril de 2014 y desde entonces ha tenido una buena aceptación en la comunidad.

Esta guía está dirigida a discutir algunos de estos problemas cruciales y explorar las diversas soluciones proporcionadas por HashiCorp Consul para abordar estos problemas.

Vamos a repasar los temas que vamos a cubrir en esta guía. Los temas están escritos para ser autocontenidos. Puedes saltar directamente a un tema específico si lo deseas.

### Breve historia sobre arquitecturas monolíticas vs. orientadas a servicios (SOA)

Al observar las arquitecturas tradicionales de entrega de aplicaciones, lo que encontramos es un monolito clásico. Cuando hablamos de monolito, tenemos un despliegue de aplicación única.

Incluso si es una aplicación única, generalmente tiene múltiples subcomponentes diferentes.

Uno de los ejemplos que dio el CTO de HashiCorp, Armon Dadgar, durante su video introductorio para Consul fue sobre la entrega de una aplicación bancaria de escritorio. Tiene un conjunto discreto de subcomponentes: por ejemplo, autenticación (digamos el subsistema A), gestión de cuentas (subsistema B), transferencia de fondos (subsistema C) y cambio de divisas (subsistema D).

Ahora, aunque estas son funciones independientes —autenticación del sistema A vs. transferencia de fondos del sistema C— lo desplegamos como una sola aplicación monolítica.

En los últimos años, hemos visto una tendencia alejarse de este tipo de arquitectura. Hay varias razones para este cambio.

El desafío con un monolito es: supongamos que hay un error en uno de los subsistemas, el sistema A, relacionado con la autenticación.

No podemos simplemente solucionarlo en el sistema A y actualizarlo en producción.

Tenemos que actualizar el sistema A y volver a implementar toda la aplicación, lo que significa que también necesitamos desplegar los subsistemas B, C y D.

Esta implementación completa no es ideal. En su lugar, nos gustaría hacer una implementación de servicios individuales.

La misma aplicación monolítica entregada como un conjunto de servicios individuales y discretos.

Por lo tanto, si hay una corrección de error en uno de nuestros servicios:

Y corregimos ese error:

Podemos hacer la reimplementación de ese servicio sin coordinar la implementación con otros servicios. De lo que estamos hablando esencialmente es de una forma de microservicios.

La corrección de errores resultará en la reimplementación solo del Servicio A dentro de toda nuestra aplicación.

Esto da un gran impulso a nuestra agilidad de desarrollo. No necesitamos coordinar nuestros esfuerzos de desarrollo entre diferentes equipos de desarrollo o incluso sistemas. Tendremos la libertad de desarrollar e implementar de forma independiente. Un servicio semanalmente y otro trimestralmente. Esto será una gran ventaja para los equipos de desarrollo.

Pero no hay tal cosa como un almuerzo gratis.

La eficiencia de desarrollo que hemos ganado introduce su propio conjunto de desafíos operativos. Veamos algunos de esos.

### Descubrimiento de servicios en un monolito, sus desafíos en un sistema distribuido y la solución de Consul

**Aplicaciones monolíticas**

Suponiendo que dos servicios en una sola aplicación quieren hablar entre sí. Una forma es exponer un método, hacerlo público y permitir que otros servicios lo llamen. En una aplicación monolítica, es una aplicación única y los servicios expondrían funciones públicas y simplemente significaría llamadas a funciones entre servicios.

Dado que esta es una llamada a función dentro de un proceso, ocurre en la memoria. Por lo tanto, es rápida y no necesitamos preocuparnos por cómo se movieron nuestros datos y si era seguro o no.

**Sistemas distribuidos**

En el mundo distribuido, el servicio A ya no se entrega como la misma aplicación que el servicio B. Entonces, ¿cómo encuentra el servicio A al servicio B si quiere hablar con él?

El servicio A puede que ni siquiera esté en la misma máquina que el servicio B. Por lo tanto, hay una red en juego. Y no es tan rápida y hay una latencia que podemos medir en líneas de milisegundos, en comparación con nanosegundos de una simple llamada a función.

**Desafíos en sistemas distribuidos**

Como ya sabemos, dos servicios en un sistema distribuido tienen que descubrirse para interactuar. Una de las formas tradicionales de resolver esto es mediante el uso de balanceadores de carga.

Los balanceadores de carga se situarían frente a cada servicio con una IP estática conocida por todos los demás servicios.

Esto da la capacidad de agregar múltiples instancias del mismo servicio detrás del balanceador de carga y dirigiría el tráfico en consecuencia. Pero esta IP del balanceador de carga es estática y está codificada dentro de todos los demás servicios, por lo que los servicios pueden omitir el descubrimiento.

El desafío ahora es mantener un conjunto de balanceadores de carga para cada servicio individual. Y podemos asumir con seguridad que originalmente había un balanceador de carga para toda la aplicación también. El costo y el esfuerzo para mantener estos balanceadores de carga han aumentado.

Con balanceadores de carga frente a los servicios, son un punto único de fallos. Incluso cuando tenemos múltiples instancias de servicio detrás del balanceador de carga, si está inactivo, nuestro servicio está inactivo. No importa cuántas instancias de ese servicio estén funcionando.

Los balanceadores de carga también aumentan la latencia de la comunicación entre servicios. Si el servicio A desea hablar con el servicio B, la solicitud de A tendrá que hablar primero con el balanceador de carga del servicio B y luego llegar a B. La respuesta de B también tendrá que pasar por el mismo proceso.

Por naturaleza, los balanceadores de carga se gestionan manualmente en la mayoría de los casos. Si agregamos otra instancia de servicio, no estará disponible de inmediato. Necesitaremos registrar ese servicio en el balanceador de carga para hacerlo accesible al mundo. Esto significaría esfuerzo manual y tiempo.

**Soluciones de Consul**

La solución de Consul al problema de descubrimiento de servicios en sistemas distribuidos es un registro de servicios central.

Consul mantiene un registro central que contiene la entrada para todos los servicios aguas arriba. Cuando se inicia una instancia de servicio, se registra en el registro central. El registro se llena con todas las instancias aguas arriba del servicio.

Cuando un servicio A quiere hablar con el servicio B, lo descubrirá y se comunicará con B consultando el registro sobre las instancias aguas arriba del servicio B. Por lo tanto, en lugar de hablar con un balanceador de carga, el servicio puede hablar directamente con la instancia de servicio de destino deseada.

Consul también proporciona comprobaciones de salud en estas instancias de servicio. Si una de las instancias de servicio o el servicio en sí no está saludable o falla en su comprobación de salud, el registro sabría sobre este escenario y evitaría devolver la dirección del servicio. El trabajo que haría el balanceador de carga es manejado por el registro en este caso.

Además, si hay múltiples instancias del mismo servicio, Consul enviaría el tráfico aleatoriamente a diferentes instancias. Por lo tanto, nivela la carga entre diferentes instancias.

Consul ha manejado nuestros desafíos de detección de fallos y distribución de carga en múltiples instancias de servicios sin la necesidad de desplegar un balanceador de carga centralizado.

El problema tradicional de balanceadores de carga lentos y gestionados manualmente se soluciona aquí. Consul gestiona el registro de forma programática, que se actualiza cuando cualquier nuevo servicio se registra y se vuelve disponible para recibir tráfico.

Esto ayuda con la escalabilidad de los servicios con facilidad.

### Gestión de la configuración en un monolito, sus desafíos en un entorno distribuido y la solución de Consul

**Aplicaciones monolíticas**

Cuando observamos la configuración para una aplicación monolítica, tienden a estar en archivos YAML, XML o JSON gigantes. Esa configuración está destinada a configurar toda la aplicación.

Dado un solo archivo, todos nuestros subsistemas

 en nuestra aplicación monolítica ahora consumirían la configuración del mismo archivo. Así, creando una vista consistente de todos nuestros subsistemas o servicios.

Si deseamos cambiar el estado de la aplicación mediante una actualización de configuración, estaría fácilmente disponible para todos los subsistemas. La nueva configuración es consumida simultáneamente por todos los componentes de nuestra aplicación.

**Sistemas distribuidos**

A diferencia del monolito, los servicios distribuidos no tendrían una vista común sobre la configuración. La configuración ahora está distribuida y cada servicio individual necesitaría ser configurado por separado.

**Desafíos en sistemas distribuidos**

La configuración debe estar repartida entre diferentes servicios. Mantener la consistencia entre la configuración en diferentes servicios después de cada actualización es un desafío.

Además, el desafío crece cuando esperamos que la configuración se actualice dinámicamente.

**Soluciones de Consul**

La solución de Consul para la gestión de configuración en un entorno distribuido es el almacén central de clave-valor.

Consul resuelve este desafío de una manera única. En lugar de distribuir la configuración entre diferentes servicios distribuidos como piezas de configuración, envía toda la configuración a todos los servicios y los configura dinámicamente en el sistema distribuido.

Tomemos un ejemplo de cambio de estado en la configuración. El estado cambiado se envía a todos los servicios en tiempo real. La configuración está presente consistentemente con todos los servicios.

### Segmentación de red en un monolito, sus desafíos en sistemas distribuidos y las soluciones de Consul

**Aplicaciones monolíticas**

Cuando observamos nuestra arquitectura monolítica clásica, la red generalmente se divide en tres zonas diferentes.

La primera zona en nuestra red es de acceso público. El tráfico que llega a nuestra aplicación a través de Internet y llega a nuestros balanceadores de carga.

La segunda zona es el tráfico desde nuestros balanceadores de carga hacia nuestra aplicación. Mayormente una zona de red interna sin acceso público directo.

La tercera zona es la zona de red cerrada, designada principalmente para datos. Esta se considera una zona aislada.

Solo la zona de balanceadores de carga puede llegar a la zona de aplicación y solo la zona de aplicación puede llegar a la zona de datos. Es un sistema de zonificación directo, simple de implementar y gestionar.

**Sistemas distribuidos**

El patrón cambia drásticamente para los servicios distribuidos.

Hay múltiples servicios dentro de nuestra zona de red de aplicación en sí. Cada uno de estos servicios habla con otros dentro de esta red, creando un patrón de tráfico complicado.

**Desafíos en sistemas distribuidos**

El principal desafío es que el tráfico no está en un flujo secuencial. A diferencia de la arquitectura monolítica, donde el flujo estaba definido desde los balanceadores de carga hasta la aplicación y la aplicación hacia los datos.

Dependiendo del patrón de acceso que deseamos soportar, el tráfico podría provenir de diferentes puntos finales y llegar a diferentes servicios.

El cliente esencialmente habla con cada servicio dentro de la aplicación directa o indirectamente.

Dado múltiples servicios y la capacidad de soportar múltiples puntos finales nos permite desplegar múltiples consumidores y proveedores de servicios.

Debido a la naturaleza del sistema, la seguridad es nuestro próximo desafío. Los servicios deberían ser capaces de identificar que el tráfico que reciben es de una entidad verificada y confiable en la red.

SOA exige control sobre las fuentes de tráfico confiables y no confiables.

Controlar el flujo de tráfico y segmentar la red en grupos o segmentos se convertirá en un problema mayor. Además, asegurarnos de que tenemos reglas estrictas que nos guíen con la partición de la red basada en quién debería poder hablar con quién y viceversa es también vital.

**Soluciones de Consul**

La solución de Consul al desafío general de segmentación de red en sistemas distribuidos es mediante la implementación de gráficos de servicios y TLS mutuo.

Consul resuelve el problema de la segmentación de red gestionando centralmente la definición sobre quién puede hablar con quién. Consul tiene una característica dedicada para esto llamada Consul Connect.

Consul Connect inscribe estas políticas de comunicación entre servicios que deseamos y las implementa como parte del gráfico de servicios. Por lo tanto, una política podría decir que el servicio A puede hablar con el servicio B, pero B no puede hablar con C, por ejemplo.

El mayor beneficio de esto es que no está restringido por IP. Más bien es a nivel de servicio. Esto lo hace escalable. La política se aplicará a todas las instancias del servicio y no habrá una regla de firewall específica de IP de un servicio. Haciéndonos independientes de la escala de nuestra red de distribución.

Consul Connect también maneja la identidad del servicio utilizando el popular protocolo TLS. Distribuye el certificado TLS asociado a un servicio.

Estos certificados ayudan a otros servicios a identificarse de manera segura entre sí. TLS también ayuda con la comunicación segura entre los servicios. Esto hace una implementación de red confiable.

Consul aplica TLS utilizando un proxy basado en agentes adjunto a cada instancia de servicio. Este proxy actúa como un sidecar. El uso de proxy, en este caso, nos previene de hacer cualquier cambio en el código del servicio original.

Esto permite el beneficio de nivel superior de aplicar cifrados en datos en reposo y en tránsito. Además, ayudará a cumplir con las leyes de privacidad e identidad del usuario.

### Arquitectura básica de Consul

Consul es un sistema distribuido y altamente disponible.

Consul se distribuye como una descarga binaria única para todas las plataformas populares. El ejecutable puede funcionar como un cliente o como un servidor.

Cada nodo que proporciona servicios a Consul ejecuta un agente de Consul. Cada uno de estos agentes habla con uno o más servidores de Consul.

El agente de Consul es responsable de verificar la salud de los servicios en el nodo, así como la verificación de la salud del nodo en sí. No es responsable del descubrimiento de servicios ni de mantener datos de clave/valor.

Los servidores de Consul son donde se almacenan y replican los datos.

Consul puede funcionar con un solo servidor, pero HashiCorp recomienda ejecutar un conjunto de 3 a 5 servidores para evitar fallos. Como todos los datos se almacenan en el lado del servidor de Consul, con un solo servidor, la falla podría causar una pérdida de datos.

Con un clúster de múltiples servidores, eligen un líder entre ellos mismos. También se recomienda por HashiCorp tener un clúster de servidores por centro de datos.

Durante el proceso de descubrimiento, cualquier servicio en busca de otros servicios puede consultar a los servidores de Consul o incluso a los agentes de Consul. Los agentes de Consul reenvían automáticamente las consultas a los servidores de Consul.

Si la consulta es intercentro de datos, las consultas son reenviadas por el servidor de Consul a los servidores de Consul remotos. Los resultados de los servidores de Consul remotos son devueltos al servidor de Consul original.

### Empezando con Consul

Esta sección está dedicada a observar de cerca a Consul como herramienta, con alguna experiencia práctica.

**Descarga e instalación**

Como se discutió anteriormente, Consul se distribuye como un binario único descargado desde el sitio web de HashiCorp o desde la sección de lanzamientos del repositorio de GitHub de Consul.

Un solo binario puede funcionar como servidor de Consul o incluso como agente cliente de Consul.

Puedes descargar Consul desde aquí — [Página de descarga de Consul](https://www.consul.io/downloads).

**Uso de Consul**

Una vez que descomprimas el archivo comprimido y pongas el binario en tu PATH, puedes ejecutarlo así.

Esto iniciará el agente en modo de desarrollo.

**Miembros de Consul**

Mientras el comando anterior está ejecutándose, puedes verificar todos los miembros en la red de Consul.

Dado que solo tenemos un nodo en ejecución, se trata como un servidor por defecto. Puedes designar un agente como servidor proporcionando el parámetro del servidor en la línea de comandos o el servidor como un parámetro de configuración en la configuración de Consul.

La salida del comando anterior se basa en el protocolo de gossip y es consistentemente eventual.

**API HTTP de Consul**

Para una vista consistentemente fuerte de la red de agentes de Consul, podemos usar la API HTTP proporcionada por Consul.

**Interfaz DNS de Consul**

Consul también proporciona una interfaz DNS para consultar nodos. Sirve DNS en el puerto 8600 por defecto. Ese puerto es configurable.

Registrar un servicio en Consul se puede lograr ya sea escribiendo una definición de servicio o enviando una solicitud a través de una API HTTP adecuada.

**Definición de servicio de Consul**

La definición de servicio es una de las formas populares de registrar un servicio. Veamos un ejemplo de definición de servicio.

Para alojar nuestras definiciones de servicio, agregaremos un directorio de configuración, convencionalmente nombrado como consul.d — '.d' representa que hay un conjunto de archivos de configuración bajo este directorio, en lugar de una única configuración bajo el nombre consul.

Escribe la definición de servicio para una aplicación web ficticia de Django ejecutándose en el puerto 80 en localhost.

Para hacer que nuestro agente de Consul sea consciente de esta definición de servicio, podemos proporcionar el directorio de configuración a él.

La información relevante en el log aquí son las declaraciones de sincronización relacionadas con el servicio "web". El agente de Consul ha aceptado nuestra configuración y la ha sincronizado en todos los nodos. En este caso, un nodo.

**Consulta de servicio DNS de Consul**

Podemos consultar el servicio con DNS, como lo hicimos con el nodo. Así:

También podemos consultar DNS para registros de servicio que nos den más información sobre los detalles del servicio como el puerto y el nodo.

También puedes usar la ETIQUETA que proporcionamos en la

 definición de servicio para consultar una etiqueta específica:

**Catálogo de servicios de Consul a través de API HTTP**

El servicio también podría ser consultado usando la API HTTP:

Podemos filtrar los servicios basándonos en las comprobaciones de salud en la API HTTP:

**Actualizar la definición de servicio de Consul**

Si deseas actualizar la definición de servicio en un agente de Consul en funcionamiento, es muy sencillo.

Hay tres formas de lograr esto. Puedes enviar una señal SIGHUP al proceso, recargar Consul que internamente envía SIGHUP en el nodo o puedes llamar a la API HTTP dedicada a las actualizaciones de definición de servicio que internamente recargará la configuración del agente.

Envía SIGHUP a 21289

O recarga Consul

Recarga de configuración activada

Deberías ver esto en tu log de Consul.

**Interfaz web de Consul**

Consul proporciona una hermosa interfaz de usuario web lista para usar. Puedes acceder a ella en el puerto 8500.

En este caso en http://localhost:8500. Veamos algunas de las pantallas.

La página de inicio de la interfaz de usuario de Consul con toda la información relevante relacionada con un agente de Consul y la comprobación del servicio web.

Explorando los servicios definidos en la interfaz web de Consul

Al profundizar en los detalles de un servicio dado, obtenemos un panel de servicio con todos los nodos y su salud para ese servicio.

Explorando la información a nivel de nodo para cada servicio en la interfaz web de Consul

En cada nodo individual, podemos ver las comprobaciones de salud, la información de los servicios y las sesiones.

Explorando la información de comprobación de salud específica del nodo, información de servicios e información de sesiones en la interfaz web de Consul

En general, la interfaz web de Consul es realmente impresionante y un gran complemento para las herramientas de línea de comandos que proporciona Consul.

### ¿Cómo es Consul diferente de Zookeeper, doozerd y etcd?

Consul tiene soporte de primera clase para descubrimiento de servicios, comprobación de salud, almacenamiento de clave-valor, centros de datos múltiples.

Zookeeper, doozerd y etcd están principalmente basados en el mecanismo de almacenamiento de clave-valor. Para lograr algo más allá de dicho clave-valor, el almacén necesita herramientas adicionales, bibliotecas y desarrollo personalizado alrededor de ellos.

Todas estas herramientas, incluyendo Consul, utilizan nodos de servidor que requieren un quórum de nodos para operar y son consistentemente fuertes.

Más o menos, todos tienen semánticas similares para la gestión del almacenamiento de clave/valor.

Estas semánticas son atractivas para construir sistemas de descubrimiento de servicios. Consul tiene soporte listo para usar para descubrimiento de servicios, que otros sistemas carecen.

Un sistema de descubrimiento de servicios también requiere una forma de realizar comprobaciones de salud. Ya que es importante verificar la salud del servicio antes de permitir que otros lo descubran. Algunos sistemas utilizan latidos con actualizaciones periódicas y TTL. El trabajo para estas comprobaciones de salud crece con la escala y requiere infraestructura fija. La ventana de detección de fallos es al menos tan larga como el TTL.

A diferencia de Zookeeper, Consul tiene agentes clientes en cada nodo del clúster, hablando entre sí en el grupo de gossip. Esto permite que los clientes sean ligeros, da mejor capacidad de comprobación de salud, reduce la complejidad del lado del cliente y resuelve desafíos de depuración.

Además, Consul proporciona soporte nativo para interfaces HTTP o DNS para realizar operaciones a nivel de sistema, nodo o servicio. Otros sistemas necesitan que estos sean desarrollados alrededor de las primitivas expuestas.

El sitio web de Consul da un buen comentario sobre las comparaciones entre Consul y otras herramientas.

### Herramientas de código abierto alrededor de HashiCorp Consul

HashiCorp y la comunidad han creado varias herramientas alrededor de Consul.

Estas herramientas de Consul son creadas y gestionadas por los ingenieros dedicados de HashiCorp:

- **Consul Template (3.3k estrellas)** — Renderización de plantillas genéricas y notificaciones con Consul. Renderización de plantillas, notificador y supervisor para los datos de Consul y Vault de @hashicorp. Proporciona una manera conveniente de poblar valores de Consul en el sistema de archivos utilizando el demonio de plantilla de consul.

- **Envconsul (1.2k estrellas)** — Leer y establecer variables ambientales para procesos desde Consul. Envconsul proporciona una manera conveniente de lanzar un subproceso con variables de entorno pobladas desde Consul y Vault de HashiCorp.

- **Consul Replicate (360 estrellas)** — Demonio de replicación de KV inter-centro de datos de Consul. Este proyecto proporciona una manera conveniente de replicar valores de un centro de datos de Consul a otro utilizando el demonio de replicación de consul.

- **Consul Migrate** — Herramienta de migración de datos para manejar actualizaciones de Consul a 0.5.1+.

La comunidad alrededor de Consul también ha creado varias herramientas para ayudar con el registro de servicios y la gestión de la configuración del servicio, me gustaría mencionar algunas de las más populares y bien mantenidas:

- **Confd (5.9k estrellas)** — Gestiona archivos de configuración de aplicaciones locales utilizando plantillas y datos de etcd o consul.

- **Fabio (5.4k estrellas)** — Fabio es un equilibrador de carga HTTP(S) y enrutador TCP rápido, moderno y sin configuración para el despliegue de aplicaciones gestionadas por consul. Registra tus servicios en consul, proporciona una comprobación de salud y Fabio comenzará a enrutar tráfico a ellos. No se requiere configuración.

- **Registrator (3.9k estrellas)** — Puente de registro de servicios para Docker con adaptadores enchufables. Registrator registra automáticamente y elimina el registro de servicios para cualquier contenedor Docker inspeccionando los contenedores a medida que se activan.

- **Hashi-UI (871 estrellas)** — Una interfaz de usuario moderna para Consul y Nomad de HashiCorp.

- **Git2consul (594 estrellas)** — Refleja el contenido de un repositorio git en KVs de Consul. Git2consul toma uno o muchos repositorios git y los refleja en KVs de Consul. El objetivo es que las organizaciones de cualquier tamaño usen git como el almacén de respaldo, la pista de auditoría y el mecanismo de control de acceso para cambios de configuración y Consul como el mecanismo de entrega.

- **Spring-cloud-consul (503 estrellas)** — Este proyecto proporciona integraciones de Consul para aplicaciones Spring Boot a través de la configuración automática y el enlace al entorno Spring y otros modelos de programación Spring. Con algunas anotaciones simples, puedes habilitar y configurar rápidamente los patrones comunes dentro de tu aplicación y construir grandes sistemas distribuidos con componentes basados en Consul.

- **Crypt (453 estrellas)** — Almacenar y recuperar configuraciones cifradas desde etcd o consul.

- **Mesos-Consul (344 estrellas)** — Puente de Mesos a Consul para descubrimiento de servicios. Mesos-consul registra automáticamente y elimina el registro de servicios ejecutados como tareas de Mesos.

- **Consul-cli (228 estrellas)** — Interfaz de línea de comandos para la API HTTP de Consul.

### Conclusión

Los sistemas distribuidos no son fáciles de construir y configurar. Mantenerlos y mantenerlos funcionando es otro trabajo aparte. HashiCorp Consul facilita la vida de los ingenieros que enfrentan tales desafíos.

A medida que pasamos por diferentes aspectos de Consul, aprendimos cuán directo se volvería para nosotros desarrollar e implementar una aplicación con arquitectura distribuida o de microservicios.

La facilidad de uso, la excelente documentación, el código robusto listo para producción y el respaldo de la comunidad, permite adaptar e introducir Consul de HashiCorp en nuestra pila tecnológica con bastante facilidad.

Esperamos que haya sido un viaje informativo en el camino de Consul. Nuestro viaje aún no ha terminado, esta fue solo la primera mitad. Nos volveremos a encontrar con la segunda parte de este artículo que nos llevará a través de ejemplos prácticos cercanos a aplicaciones de la vida real.

Háganos saber qué le gustaría escuchar más de nosotros o si tiene alguna pregunta sobre el tema, estaremos más que felices de responder a esas.

---

Espero que esta traducción sea de utilidad. Si necesitas ayuda con algo más, ¡avísame!