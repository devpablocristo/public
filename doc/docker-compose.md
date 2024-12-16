Si un servicio no tiene un perfil asignado en su configuración de `docker-compose.yml`, Docker Compose lo considera parte de la configuración predeterminada y lo levantará junto con cualquier otro servicio cuando ejecutas el comando `docker compose up` o incluso cuando especificas un perfil específico.

### Solución para Controlar los Servicios que se Levantan

Para evitar que Docker Compose levante servicios que no quieres, debes asegurarte de asignar un perfil a **todos los servicios** en tu archivo `docker-compose.yml`. De esta manera, solo se levantarán los servicios que coincidan con los perfiles que especifiques en el comando de ejecución.

### Ajuste en el Archivo `docker-compose.yml`

Modifica tu archivo `docker-compose.yml` para asegurarte de que cada servicio tenga un perfil asignado. Aquí te muestro cómo puedes hacerlo:

```yaml
version: "3.8"

services:
  golang-sdk:
    container_name: "${APP_NAME:-golang-sdk}"
    build:
      context: ..
      dockerfile: config/Dockerfile.dev
    image: "${APP_NAME:-golang-sdk}:${APP_VERSION:-1.0}"
    ports:
      - "${ROUTER_PORT:-8080}:${ROUTER_TARGET_PORT:-8080}"
      - "${DELVE_PORT:-2345}:${DELVE_TARGET_PORT:-2345}"
    volumes:
      - type: bind
        source: ..
        target: /app
      - type: bind
        source: ../migrations
        target: /app/migrations
    environment:
      - APP_NAME=${APP_NAME:-golang-sdk}
      # ... otras variables de entorno
    depends_on:
      - postgres
      - consul
      - cassandra
      - mysql
      - redis
      - prometheus
      - grafana
      - pyroscope
      - mongodb
      - rabbitmq
      - dynamodb
    networks:
      - app-network
    restart: on-failure
    profiles:
      - golang-sdk  # Asignar perfil específico

  monitoring-api:
    container_name: "monitoring-api"
    build:
      context: ..
      dockerfile: config/Dockerfile.dev
    image: "monitoring-api:${APP_VERSION:-1.0}"
    working_dir: /app/cmd/examples/monitoring
    command: go run main.go
    ports:
      - "${ROUTER_PORT:-8080}:${ROUTER_TARGET_PORT:-8080}"
      - "${DELVE_PORT:-2345}:${DELVE_TARGET_PORT:-2345}"
    volumes:
      - type: bind
        source: ..
        target: /app
    environment:
      - APP_NAME=monitoring-api
      - APP_VERSION=${APP_VERSION:-1.0}
      - DEBUG=${DEBUG:-true}
      - PROMETHEUS_URL=http://prometheus:9090
    depends_on:
      - mysql
      - prometheus
    networks:
      - app-network
    restart: on-failure
    profiles:
      - monitoring  # Asignar perfil específico

  mysql:
    image: mysql:8.0
    container_name: mysql
    ports:
      - "${MYSQL_PORT:-3306}:${MYSQL_PORT:-3306}"
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-root}
      - MYSQL_DATABASE=${MYSQL_DATABASE:-dev_events_db}
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - app-network
    restart: on-failure
    profiles:
      - monitoring  # Asignar perfil específico

  prometheus:
    image: prom/prometheus:v2.45.6
    container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "${PROMETHEUS_PORT:-9090}:${PROMETHEUS_PORT:-9090}"
    networks:
      - app-network
    restart: on-failure
    profiles:
      - monitoring  # Asignar perfil específico

  grafana:
    image: grafana/grafana:10.2.8
    container_name: grafana
    ports:
      - "${GRAFANA_PORT:-3000}:${GRAFANA_PORT:-3000}"
    depends_on:
      - prometheus
    networks:
      - app-network
    restart: on-failure
    profiles:
      - monitoring  # Asignar perfil específico

  # Asegúrate de asignar un perfil a TODOS los demás servicios también.
  # Por ejemplo:
  consul:
    image: consul:1.15.4
    container_name: consul
    ports:
      - "${CONSUL_PORT:-8500}:${CONSUL_TARGET_PORT:-8500}"
    networks:
      - app-network
    restart: on-failure
    profiles:
      - consul  # Asignar perfil específico

  # Continúa con el resto de tus servicios...
```

### Ejecución con Perfiles

Con esta configuración, puedes ejecutar solo los servicios que pertenecen al perfil `monitoring` sin que Docker Compose levante otros servicios predeterminados:

```bash
docker compose -f ./config/docker-compose.dev.yml --profile monitoring up
```

Solo levantará los servicios que tienen `profiles: [monitoring]` asignado.

### Conclusión

Asignando perfiles explícitos a todos los servicios en tu archivo `docker-compose.yml`, puedes controlar exactamente qué servicios se levantan cuando usas perfiles específicos. Esto evitará que Docker Compose inicie cualquier servicio de forma predeterminada que no desees ejecutar.