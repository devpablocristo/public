#!/bin/bash

# Variables importantes
PROJECT_ROOT=".." # Raíz del proyecto (donde está el script)
SWAG_CMD="swag"
DOCS_DIR="$PROJECT_ROOT/docs"
HANDLER_FILE="$PROJECT_ROOT/projects/customers-manager/internal/customer/adapters/inbound/gin-handler.go"
SWAG_INIT_FLAGS="-g \"$HANDLER_FILE\" --parseDependency -o \"$DOCS_DIR/swagger.yaml\""

# Mensajes informativos
echo "Generando documentación de Swagger..."

# Crear el directorio de documentación (si no existe)
mkdir -p "$DOCS_DIR"

# Verificar si swag está instalado
if ! command -v "$SWAG_CMD" &> /dev/null; then
    echo "Error: $SWAG_CMD no está instalado. Por favor, instala swag: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

# Verificar si el archivo handler existe
if [ ! -f "$HANDLER_FILE" ]; then
    echo "Error: Archivo handler no encontrado: $HANDLER_FILE"
    exit 1
fi

# Cambiar al directorio raíz del proyecto antes de ejecutar swag init
pushd "$PROJECT_ROOT" > /dev/null

# Ejecutar swag init
$SWAG_CMD init $SWAG_INIT_FLAGS

# Volver al directorio original
popd > /dev/null

# Verificar si hubo errores
if [ $? -eq 0 ]; then
    GENERATED_FILE="$DOCS_DIR/swagger.yaml"
    echo "Documentación generada exitosamente en: $(realpath "$GENERATED_FILE")"
else
    echo "Error al generar la documentación. Revisa los mensajes anteriores."
    exit 1
fi

echo "Proceso completado."