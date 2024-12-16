#!/bin/bash

log() {
  echo "$(date +"%Y-%m-%d %H:%M:%S") [INFO] $1"
}

buildServer() {
  log "Building server binary"
  go build -gcflags "all=-N -l" -buildvcs=false -o "/app/bin/${APP_NAME}" "${BUILDING_FILES}"
  if [ -f "/app/bin/${APP_NAME}" ]; then
    log "Binary file created successfully"
    chmod +x "/app/bin/${APP_NAME}"
  else
    log "Failed to create binary file"
    exit 1
  fi
}

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

rerunServer() {
  log "Rerun server"
  buildServer
  runServer
}

liveReloading() {
  log "Run liveReloading"
  inotifywait -e modify,delete,move -m -r --format '%w%f' --exclude '.*(\.tmp|\.swp)$' /app | (
    while read -r file; do
      if [ "${file##*.}" = "go" ]; then
        log "File ${file} changed. Reloading..."
        rerunServer
      fi
    done
  )
}

# Ejecuta la función basada en el argumento pasado
case "$1" in
  "rerunServer")
    rerunServer
    ;;
  *)
    echo "Usage: $0 {rerunServer}"
    exit 1
    ;;
esac
