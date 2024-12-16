### ¿Qué es Kong?

Kong es una plataforma de API Gateway y Microservices Management. Se utiliza para gestionar, monitorear y asegurar el tráfico de las API de manera eficiente. Kong ofrece funcionalidades como autenticación, autorización, limitación de tasa, registros, monitoreo y más, sin necesidad de implementarlas directamente en las APIs.

### ¿Para qué se usa?

Kong se usa para:

1. **Gestión de API**: Permite centralizar la gestión de múltiples APIs en una sola plataforma.
2. **Seguridad**: Implementa autenticación, autorización y cifrado de tráfico.
3. **Monitoreo**: Proporciona herramientas para monitorear el rendimiento y el uso de las APIs.
4. **Escalabilidad**: Facilita el balanceo de carga y la escalabilidad de las APIs.
5. **Registro y Análisis**: Permite registrar y analizar las solicitudes a las APIs.

### Ejemplo de uso en una API Golang

Supongamos que tienes una API en Golang que proporciona información sobre eventos. Kong puede ayudar a gestionar esta API, proporcionando autenticación, limitación de tasa y monitoreo.

### Ventajas de usar Kong

1. **Facilidad de Configuración**: Kong es fácil de configurar y desplegar.
2. **Extensibilidad**: Ofrece una gran cantidad de plugins que pueden extender sus funcionalidades.
3. **Rendimiento**: Kong está diseñado para manejar grandes volúmenes de tráfico.
4. **Comunidad y Soporte**: Tiene una comunidad activa y soporte comercial disponible.
5. **Escalabilidad**: Facilita la escalabilidad horizontal de tus servicios.

### ¿Por qué usa PostgreSQL?

Kong utiliza PostgreSQL como su base de datos para almacenar configuraciones, datos de plugins y otros metadatos importantes. PostgreSQL es una base de datos robusta, escalable y ampliamente utilizada, lo que la convierte en una opción fiable para Kong.

### ¿Cómo se implementa Golang?

Para implementar una API en Golang, se puede usar el framework Gin como un ejemplo. Aquí hay un ejemplo básico de una API en Golang:

#### Ejemplo de API en Golang con Gin

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/events", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "List of events",
        })
    })
    r.Run(":8080") // La API estará disponible en el puerto 8080
}
```

### Configuración de Kong en Docker Compose

La configuración de Kong en Docker Compose incluye la inicialización de migraciones de la base de datos y la definición de variables de entorno para conectar con PostgreSQL y definir los puertos.

#### `docker-compose.yml` con Kong

```yaml
services:
  kong-migrations:
    image: kong:latest
    command: kong migrations bootstrap
    environment:
      - KONG_DATABASE=postgres
      - KONG_PG_HOST=${POSTGRES_HOST:-postgres}
      - KONG_PG_PORT=${POSTGRES_HOST_PORT:-5432}
      - KONG_PG_USER=${POSTGRES_USERNAME:-postgres}
      - KONG_PG_PASSWORD=${POSTGRES_USER_PASSWORD:-root}
      - KONG_PG_DATABASE=${POSTGRES_DATABASE:-dev_events_db}
    depends_on:
      - postgres
    networks:
      - app-network

  kong:
    image: kong:latest
    container_name: kong
    environment:
      - KONG_DATABASE=postgres
      - KONG_PG_HOST=${POSTGRES_HOST:-postgres}
      - KONG_PG_PORT=${POSTGRES_HOST_PORT:-5432}
      - KONG_PG_USER=${POSTGRES_USERNAME:-postgres}
      - KONG_PG_PASSWORD=${POSTGRES_USER_PASSWORD:-root}
      - KONG_PG_DATABASE=${POSTGRES_DATABASE:-dev_events_db}
      - KONG_PROXY_ACCESS_LOG=${KONG_PROXY_ACCESS_LOG:-/dev/stdout}
      - KONG_ADMIN_ACCESS_LOG=${KONG_ADMIN_ACCESS_LOG:-/dev/stdout}
      - KONG_PROXY_ERROR_LOG=${KONG_PROXY_ERROR_LOG:-/dev/stderr}
      - KONG_ADMIN_ERROR_LOG=${KONG_ADMIN_ERROR_LOG:-/dev/stderr}
      - KONG_ADMIN_LISTEN=${KONG_ADMIN_LISTEN:-0.0.0.0:8001}
    ports:
      - "${KONG_PROXY_PORT:-8000}:${KONG_PROXY_PORT:-8000}"
      - "${KONG_PROXY_SSL_PORT:-8443}:${KONG_PROXY_SSL_PORT:-8443}"
      - "${KONG_ADMIN_PORT:-8001}:${KONG_ADMIN_PORT:-8001}"
      - "${KONG_ADMIN_SSL_PORT:-8444}:${KONG_ADMIN_SSL_PORT:-8444}"
    depends_on:
      - kong-migrations
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

### Variables de entorno necesarias

Asegúrate de definir las siguientes variables de entorno:

```env
POSTGRES_HOST=postgres
POSTGRES_HOST_PORT=5432
POSTGRES_USERNAME=postgres
POSTGRES_USER_PASSWORD=root
POSTGRES_DATABASE=dev_events_db
KONG_PROXY_ACCESS_LOG=/dev/stdout
KONG_ADMIN_ACCESS_LOG=/dev/stdout
KONG_PROXY_ERROR_LOG=/dev/stderr
KONG_ADMIN_ERROR_LOG=/dev/stderr
KONG_ADMIN_LISTEN=0.0.0.0:8001
KONG_PROXY_PORT=8000
KONG_PROXY_SSL_PORT=8443
KONG_ADMIN_PORT=8001
KONG_ADMIN_SSL_PORT=8444
```

### Explicación de la configuración de Kong

- **kong-migrations**: Este servicio ejecuta las migraciones necesarias en la base de datos de PostgreSQL para preparar la configuración de Kong.
- **kong**: Este es el servicio principal de Kong, que se conecta a PostgreSQL para almacenar configuraciones y metadatos. También define los puertos para el acceso proxy y administrativo.
- **ports**: Define los puertos que Kong utilizará para el proxy de las APIs y para el acceso administrativo.

### Probar la Configuración

1. Inicia los contenedores con `docker-compose up -d`.
2. Configura Kong para enrutar solicitudes a tu API en Golang con los siguientes comandos:

```sh
# Añadir un Servicio
curl -i -X POST \
  --url http://localhost:8001/services/ \
  --data 'name=my-service' \
  --data 'url=http://rest:8080'

# Añadir una Ruta
curl -i -X POST \
  --url http://localhost:8001/services/my-service/routes \
  --data 'paths[]=/my-api'
```

3. Prueba la configuración haciendo una solicitud a través de Kong:

```sh
curl -i -X GET \
  --url http://localhost:8000/my-api
```

Kong puede usar diferentes bases de datos para almacenar su configuración y datos de plugins. Las opciones más comunes son PostgreSQL y Cassandra. Aquí hay un resumen de las opciones de base de datos que Kong soporta:

### Opciones de Base de Datos para Kong

1. **PostgreSQL**:
   - Es la opción más común y recomendada para la mayoría de las implementaciones de Kong.
   - Proporciona una base de datos relacional robusta y escalable.
   - Facilita la integración con muchas herramientas y servicios.

2. **Cassandra**:
   - Es una base de datos NoSQL distribuida.
   - Adecuada para implementaciones donde se necesita alta disponibilidad y escalabilidad horizontal.
   - Puede ser más compleja de configurar y mantener en comparación con PostgreSQL.

### Configuración de Kong con Cassandra

Si decides usar Cassandra en lugar de PostgreSQL, aquí tienes cómo ajustar la configuración:

#### Configuración de Cassandra en `docker-compose.yml`

```yaml
services:
  cassandra:
    image: cassandra:4.1
    container_name: cassandra
    ports:
      - "${CASSANDRA_PORT:-9042}:${CASSANDRA_TARGET_PORT:-9042}"
    environment:
      - CASSANDRA_CLUSTER_NAME=${CASSANDRA_CLUSTER_NAME:-CassandraCluster}
      - CASSANDRA_DC=${CASSANDRA_DC:-datacenter1}
      - CASSANDRA_RACK=${CASSANDRA_RACK:-rack1}
      - CASSANDRA_SEEDS=${CASSANDRA_HOST:-cassandra}
      - CASSANDRA_ENDPOINT_SNITCH=${CASSANDRA_ENDPOINT_SNITCH:-GossipingPropertyFileSnitch}
      - CASSANDRA_USERNAME=${CASSANDRA_USERNAME:-cassandra}
      - CASSANDRA_PASSWORD=${CASSANDRA_PASSWORD:-cassandra}
    volumes:
      - cassandra_data:/var/lib/cassandra
    networks:
      - app-network
    restart: on-failure

  kong-migrations:
    image: kong:latest
    command: kong migrations bootstrap
    environment:
      - KONG_DATABASE=cassandra
      - KONG_CASSANDRA_CONTACT_POINTS=${CASSANDRA_HOST:-cassandra}
      - KONG_CASSANDRA_PORT=${CASSANDRA_PORT:-9042}
      - KONG_CASSANDRA_KEYSPACE=${CASSANDRA_KEYSPACE:-kong}
      - KONG_CASSANDRA_USERNAME=${CASSANDRA_USERNAME:-cassandra}
      - KONG_CASSANDRA_PASSWORD=${CASSANDRA_PASSWORD:-cassandra}
    depends_on:
      - cassandra
    networks:
      - app-network

  kong:
    image: kong:latest
    container_name: kong
    environment:
      - KONG_DATABASE=cassandra
      - KONG_CASSANDRA_CONTACT_POINTS=${CASSANDRA_HOST:-cassandra}
      - KONG_CASSANDRA_PORT=${CASSANDRA_PORT:-9042}
      - KONG_CASSANDRA_KEYSPACE=${CASSANDRA_KEYSPACE:-kong}
      - KONG_CASSANDRA_USERNAME=${CASSANDRA_USERNAME:-cassandra}
      - KONG_CASSANDRA_PASSWORD=${CASSANDRA_PASSWORD:-cassandra}
      - KONG_PROXY_ACCESS_LOG=${KONG_PROXY_ACCESS_LOG:-/dev/stdout}
      - KONG_ADMIN_ACCESS_LOG=${KONG_ADMIN_ACCESS_LOG:-/dev/stdout}
      - KONG_PROXY_ERROR_LOG=${KONG_PROXY_ERROR_LOG:-/dev/stderr}
      - KONG_ADMIN_ERROR_LOG=${KONG_ADMIN_ERROR_LOG:-/dev/stderr}
      - KONG_ADMIN_LISTEN=${KONG_ADMIN_LISTEN:-0.0.0.0:8001}
    ports:
      - "${KONG_PROXY_PORT:-8000}:${KONG_PROXY_PORT:-8000}"
      - "${KONG_PROXY_SSL_PORT:-8443}:${KONG_PROXY_SSL_PORT:-8443}"
      - "${KONG_ADMIN_PORT:-8001}:${KONG_ADMIN_PORT:-8001}"
      - "${KONG_ADMIN_SSL_PORT:-8444}:${KONG_ADMIN_SSL_PORT:-8444}"
    depends_on:
      - kong-migrations
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

volumes:
  cassandra_data:
```

### Variables de entorno necesarias

Asegúrate de definir las siguientes variables de entorno para Cassandra:

```env
CASSANDRA_HOST=cassandra
CASSANDRA_PORT=9042
CASSANDRA_KEYSPACE=kong
CASSANDRA_USERNAME=cassandra
CASSANDRA_PASSWORD=cassandra
KONG_PROXY_ACCESS_LOG=/dev/stdout
KONG_ADMIN_ACCESS_LOG=/dev/stdout
KONG_PROXY_ERROR_LOG=/dev/stderr
KONG_ADMIN_ERROR_LOG=/dev/stderr
KONG_ADMIN_LISTEN=0.0.0.0:8001
KONG_PROXY_PORT=8000
KONG_PROXY_SSL_PORT=8443
KONG_ADMIN_PORT=8001
KONG_ADMIN_SSL_PORT=8444
```

### Explicación de la Configuración

1. **Cassandra**: Define el servicio Cassandra que Kong utilizará como su base de datos.
2. **kong-migrations**: Ejecuta las migraciones necesarias para configurar Cassandra para su uso con Kong.
3. **kong**: Configura el servicio principal de Kong para utilizar Cassandra como base de datos.

### Ventajas y Consideraciones

- **PostgreSQL**:
  - **Ventajas**: Fácil de configurar, soporta transacciones, familiar para la mayoría de los desarrolladores.
  - **Consideraciones**: Puede no ser tan eficiente en distribuciones masivas comparado con Cassandra.

- **Cassandra**:
  - **Ventajas**: Alta disponibilidad, escalabilidad horizontal, adecuada para grandes volúmenes de datos distribuidos.
  - **Consideraciones**: Más compleja de configurar y mantener, no soporta transacciones de la misma manera que PostgreSQL.

