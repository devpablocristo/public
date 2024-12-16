#!/bin/sh

# shellcheck disable=SC2154  # Desactivar aviso de shellcheck para APP_NAME
# shellcheck source=./config/.env disable=SC1091 # Desactivar aviso de shellcheck de archivo no especificado

# INFO: este script funciona perfectamente, hace lo mismo que air, pero con inotifywatch.
# Load environment variables from the ./config/.env file in the parent directory
loadEnv() {
  if [ -f ./config/.env ]; then
    log "Loading environment variables from ./config/.env"
    # Use `set -a` to export all variables
    set -a
    # shellcheck source=./config/.env
    . ./config/.env
    set +a
  else
    echo "ERROR: ./config/.env file not found in the parent directory. Please create ./config/.env with the necessary environment variables."
    exit 1
  fi
}

# Function to log messages
log() {
  echo "ENTRYPOINT: $1"
}

# Validate essential environment variables
validateEnv() {
  if [ -z "${APP_NAME}" ]; then
    log "ERROR: APP_NAME is not set. Please check ./config/.env file."
    exit 1
  fi

  if [ -z "${DEBUG}" ]; then
    log "ERROR: DEBUG is not set. Please check ./config/.env file."
    exit 1
  fi

  log "Environment variables loaded successfully"
  log "App Name: ${APP_NAME}"
  log "Debug: ${DEBUG}"
}

# Function to build the server binary
buildServer() {
  log "Building server binary"
  go build -gcflags "all=-N -l" -buildvcs=false -o "/app/bin/${APP_NAME}" "${BUILDING_FILES}"
  # Verify if the binary file has been created and is executable
  if [ -f "/app/bin/${APP_NAME}" ]; then
    log "Binary file created successfully"
    chmod +x "/app/bin/${APP_NAME}"
  else
    log "Failed to create binary file"
    exit 1
  fi
}

# Function to run the server
runServer() {
  log "Run server"

  log "Killing old server"
  pkill -f dlv || true
  pkill -f "/app/bin/${APP_NAME}" || true

  if [ "${DEBUG}" = "true" ]; then
    log "Run in debug mode"
    dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec "/app/bin/${APP_NAME}" &
    liveReloading
  else
    log "Run in production mode"
    "/app/bin/${APP_NAME}"
  fi
}

# Function to rebuild and rerun the server
rerunServer() {
  log "Rerun server"
  buildServer
  runServer
}

# Function to monitor file changes and trigger server restart
liveReloading() {
  log "Run liveReloading"
  inotifywait -e modify,delete,move -m -r --format '%w%f' --exclude '.*(\.tmp|\.swp)$' /app | (
    while read -r file; do
      # Use [ ] instead of [[ ]] for POSIX compatibility
      if [ "${file##*.}" = "go" ]; then
        log "File ${file} changed. Reloading..."
        rerunServer
      fi
    done
  )
}

# Function to initialize the file change logger
initializeFileChangeLogger() {
  echo "" > /tmp/filechanges.log
  tail -f /tmp/filechanges.log &
}

# Main function to orchestrate the process
main() {
  log "Starting script"
  log "Current directory: $(pwd)"
  loadEnv
  validateEnv
  initializeFileChangeLogger
  buildServer
  runServer
}

# Call the main function to start the process
main