#!/bin/bash

# Nombre exacto del contenedor Docker
CONTAINER_NAME="sg_auth"

# Comando a ejecutar dentro del contenedor (ejecutado en segundo plano)
COMMAND="sh /app/scripts/entrypoint.sh &>/dev/null &"

# Verificar si el contenedor está corriendo utilizando coincidencia exacta
CONTAINER_ID=$(docker ps --filter "name=^${CONTAINER_NAME}$" --format "{{.ID}}")

# Si se encontró un contenedor con el nombre exacto
if [ -n "${CONTAINER_ID}" ]; then
    echo "Ejecutando script dentro del contenedor: ${CONTAINER_NAME} (ID: ${CONTAINER_ID})"
    docker exec ${CONTAINER_ID} sh -c "${COMMAND}"
else
    echo "Error: El contenedor ${CONTAINER_NAME} no está corriendo."
fi
