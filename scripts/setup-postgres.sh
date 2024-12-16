#!/bin/bash

# NOTE: no funciona el script del todo, pudo crear la base de datos (con los pasos de la documentacion de pgAdmin)
# Función para cargar las variables de entorno desde el archivo ./config/.env
load_env() {
  export $(grep -v '^#' .././config/.env | xargs)
}

# Llamar a la función para cargar las variables
load_env

# Variables de entorno cargadas desde el ./config/.env
POSTGRES_HOST=${POSTGRES_HOST:-postgres}
POSTGRES_PORT=${POSTGRES_PORT:-5432}
POSTGRES_NAME=${POSTGRES_DATABASE:-dev_events_db}
POSTGRES_USER=${POSTGRES_USERNAME:-postgres}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-root}
POSTGRES_ROOT_PASSWORD=${POSTGRES_ROOT_PASSWORD:-rootpassword}
POSTGRES_TABLE=${POSTGRES_TABLE:-events}
COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME:-my_project}

# Obtener el nombre de la red de Docker Compose
NETWORK_NAME=$(docker network ls --filter name=${COMPOSE_PROJECT_NAME}_default --format "{{.Name}}")

if [ -z "$NETWORK_NAME" ]; then
  echo "Docker network not found. Please make sure Docker Compose is up and running."
  exit 1
fi

# Esperar a que el contenedor de PostgreSQL esté listo para aceptar conexiones
until docker run --rm --network=${NETWORK_NAM} -e PGPASSWORD=$POSTGRES_ROOT_PASSWORD postgres:16 psql -U postgres -h $POSTGRES_HOST -p $POSTGRES_PORT -c "\q"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

# Crear la base de datos
docker run --rm --network=$NETWORK_NAME -e PGPASSWORD=$POSTGRES_ROOT_PASSWORD postgres:16 psql -U postgres -h $POSTGRES_HOST -p $POSTGRES_PORT -c "CREATE DATABASE $POSTGRES_NAME;"

# Crear el usuario y otorgar privilegios
docker run --rm --network=$NETWORK_NAME -e PGPASSWORD=$POSTGRES_ROOT_PASSWORD postgres:16 psql -U postgres -h $POSTGRES_HOST -p $POSTGRES_PORT -c "CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';"
docker run --rm --network=$NETWORK_NAME -e PGPASSWORD=$POSTGRES_ROOT_PASSWORD postgres:16 psql -U postgres -h $POSTGRES_HOST -p $POSTGRES_PORT -c "GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_NAME TO $POSTGRES_USER;"

# Crear la tabla
docker run --rm --network=$NETWORK_NAME -e PGPASSWORD=$POSTGRES_PASSWORD postgres:16 psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_NAME -c "CREATE TABLE $POSTGRES_TABLE (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  location VARCHAR(255),
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP,
  category VARCHAR(50),
  creator_id UUID NOT NULL,
  is_public BOOLEAN NOT NULL DEFAULT true,
  is_recurring BOOLEAN NOT NULL DEFAULT false,
  series_id UUID,
  status VARCHAR(50) NOT NULL
);"

echo "Database and table created successfully"
