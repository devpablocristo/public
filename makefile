SHELL := /bin/bash

# Variables
ROOT_DIR := $(shell pwd)
VERSION := 1.0
BUILD_DIR := ${ROOT_DIR}/bin
CONFIG_DIR := ${ROOT_DIR}/config
SCRIPTS_DIR := ${ROOT_DIR}/scripts
DOCKER_COMPOSE_DEV := ${CONFIG_DIR}/docker-compose.dev.yml
DOCKER_COMPOSE_STG := ${CONFIG_DIR}/docker-compose.stg.yml
DOCKER_COMPOSE_PROD := ${CONFIG_DIR}/docker-compose.prod.yml

# Phony targets declaration
.PHONY: all build run test clean lint \
	docker-dev-build docker-dev-up docker-dev-down docker-dev-logs \
	docker-stg-build docker-stg-up docker-stg-down docker-stg-logs \
	docker-prod-build docker-prod-up docker-prod-down docker-prod-logs \
	sg_backend-dev-build sg_backend-dev-up sg_backend-dev-down sg_backend-dev-logs \
	sg_backend-stg-build sg_backend-stg-up sg_backend-stg-down sg_backend-stg-logs \
	sg_backend-prod-build sg_backend-prod-up sg_backend-prod-down sg_backend-prod-logs \
	sg_requests-dev-build sg_requests-dev-up sg_requests-dev-down sg_requests-dev-logs \
	sg_requests-stg-build sg_requests-stg-up sg_requests-stg-down sg_requests-stg-logs \
	sg_requests-prod-build sg_requests-prod-up sg_requests-prod-down sg_requests-prod-logs \
	tech-house-dev-build tech-house-dev-up tech-house-dev-down tech-house-dev-logs \
	tech-house-stg-build tech-house-stg-up tech-house-stg-down tech-house-stg-logs \
	tech-house-prod-build tech-house-prod-up tech-house-prod-down tech-house-prod-logs

# Default target
all: build run

# Core commands
build:
	@echo "Building the project..."
	@mkdir -p ${BUILD_DIR}
	APP_NAME=main go build -gcflags "all=-N -l" -o ${BUILD_DIR}/${APP_NAME} -ldflags "-X main.Version=${VERSION}" ${ROOT_DIR}/cmd/main.go

run:
	@echo "Running the project..."
	@go run ${ROOT_DIR}/cmd/main.go

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -f ${BUILD_DIR}/main

lint:
	@echo "Linting the project..."
	@golangci-lint run --config .golangci.yml --verbose

# Development Docker commands
docker-dev-build:
	@echo "Building $(PROFILE) services in dev mode..."
	@chmod +x ${SCRIPTS_DIR}/entrypoint.sh
	docker compose -f ${DOCKER_COMPOSE_DEV} --profile $(PROFILE) up --build -d
	@$(MAKE) docker-dev-logs PROFILE=$(PROFILE)

docker-dev-up:
	@echo "Starting $(PROFILE) services in dev mode..."
	@chmod +x ${SCRIPTS_DIR}/entrypoint.sh
	docker compose -f ${DOCKER_COMPOSE_DEV} --profile $(PROFILE) up -d
	@$(MAKE) docker-dev-logs PROFILE=$(PROFILE)

docker-dev-down:
	@echo "Stopping $(PROFILE) services in dev mode..."
	docker compose -f ${DOCKER_COMPOSE_DEV} --profile $(PROFILE) down --remove-orphans

docker-dev-logs:
	@echo "Fetching logs for $(PROFILE) services in dev..."
	docker compose -f ${DOCKER_COMPOSE_DEV} --profile $(PROFILE) logs -f

# Staging Docker commands
docker-stg-build:
	@echo "Building $(PROFILE) services in staging mode..."
	docker compose -f ${DOCKER_COMPOSE_STG} --profile $(PROFILE) up --build -d
	@$(MAKE) docker-stg-logs PROFILE=$(PROFILE)

docker-stg-up:
	@echo "Starting $(PROFILE) services in staging mode..."
	docker compose -f ${DOCKER_COMPOSE_STG} --profile $(PROFILE) up -d
	@$(MAKE) docker-stg-logs PROFILE=$(PROFILE)

docker-stg-down:
	@echo "Stopping $(PROFILE) services in staging mode..."
	docker compose -f ${DOCKER_COMPOSE_STG} --profile $(PROFILE) down --remove-orphans

docker-stg-logs:
	@echo "Fetching logs for $(PROFILE) services in staging..."
	docker compose -f ${DOCKER_COMPOSE_STG} --profile $(PROFILE) logs -f

# Production Docker commands
docker-prod-build:
	@echo "Building $(PROFILE) services in production mode..."
	docker compose -f ${DOCKER_COMPOSE_PROD} --profile $(PROFILE) up --build -d
	@$(MAKE) docker-prod-logs PROFILE=$(PROFILE)

docker-prod-up:
	@echo "Starting $(PROFILE) services in production mode..."
	docker compose -f ${DOCKER_COMPOSE_PROD} --profile $(PROFILE) up -d
	@$(MAKE) docker-prod-logs PROFILE=$(PROFILE)

docker-prod-down:
	@echo "Stopping $(PROFILE) services in production mode..."
	docker compose -f ${DOCKER_COMPOSE_PROD} --profile $(PROFILE) down --remove-orphans

docker-prod-logs:
	@echo "Fetching logs for $(PROFILE) services in production..."
	docker compose -f ${DOCKER_COMPOSE_PROD} --profile $(PROFILE) logs -f

# Requests profile targets for dev
sg_requests-dev-build:
	@$(MAKE) docker-dev-build PROFILE=sg_requests

sg_requests-dev-up:
	@$(MAKE) docker-dev-up PROFILE=sg_requests

sg_requests-dev-down:
	@$(MAKE) docker-dev-down PROFILE=sg_requests

sg_requests-dev-logs:
	@$(MAKE) docker-dev-logs PROFILE=sg_requests

# Requests profile targets for staging
sg_requests-stg-build:
	@$(MAKE) docker-stg-build PROFILE=sg_requests

sg_requests-stg-up:
	@$(MAKE) docker-stg-up PROFILE=sg_requests

sg_requests-stg-down:
	@$(MAKE) docker-stg-down PROFILE=sg_requests

sg_requests-stg-logs:
	@$(MAKE) docker-stg-logs PROFILE=sg_requests

# Requests profile targets for production
sg_requests-prod-build:
	@$(MAKE) docker-prod-build PROFILE=sg_requests

sg_requests-prod-up:
	@$(MAKE) docker-prod-up PROFILE=sg_requests

sg_requests-prod-down:
	@$(MAKE) docker-prod-down PROFILE=sg_requests

sg_requests-prod-logs:
	@$(MAKE) docker-prod-logs PROFILE=sg_requests

# Backend profile targets for dev
sg_backend-dev-build:
	@$(MAKE) docker-dev-build PROFILE=sg_backend

sg_backend-dev-up:
	@$(MAKE) docker-dev-up PROFILE=sg_backend

sg_backend-dev-down:
	@$(MAKE) docker-dev-down PROFILE=sg_backend

sg_backend-dev-logs:
	@$(MAKE) docker-dev-logs PROFILE=sg_backend

# Backend profile targets for staging
sg_backend-stg-build:
	@$(MAKE) docker-stg-build PROFILE=sg_backend

sg_backend-stg-up:
	@$(MAKE) docker-stg-up PROFILE=sg_backend

sg_backend-stg-down:
	@$(MAKE) docker-stg-down PROFILE=sg_backend

sg_backend-stg-logs:
	@$(MAKE) docker-stg-logs PROFILE=sg_backend

# Backend profile targets for production
sg_backend-prod-build:
	@$(MAKE) docker-prod-build PROFILE=sg_backend

sg_backend-prod-up:
	@$(MAKE) docker-prod-up PROFILE=sg_backend

sg_backend-prod-down:
	@$(MAKE) docker-prod-down PROFILE=sg_backend

sg_backend-prod-logs:
	@$(MAKE) docker-prod-logs PROFILE=sg_backend

# Tech House profile targets for dev
tech-house-dev-build:
	@$(MAKE) docker-dev-build PROFILE=tech-house

tech-house-dev-up:
	@$(MAKE) docker-dev-up PROFILE=tech-house

tech-house-dev-down:
	@$(MAKE) docker-dev-down PROFILE=tech-house

tech-house-dev-logs:
	@$(MAKE) docker-dev-logs PROFILE=tech-house

# Tech House profile targets for staging
tech-house-stg-build:
	@$(MAKE) docker-stg-build PROFILE=tech-house

tech-house-stg-up:
	@$(MAKE) docker-stg-up PROFILE=tech-house

tech-house-stg-down:
	@$(MAKE) docker-stg-down PROFILE=tech-house

tech-house-stg-logs:
	@$(MAKE) docker-stg-logs PROFILE=tech-house

# Tech House profile targets for production
tech-house-prod-build:
	@$(MAKE) docker-prod-build PROFILE=tech-house

tech-house-prod-up:
	@$(MAKE) docker-prod-up PROFILE=tech-house

tech-house-prod-down:
	@$(MAKE) docker-prod-down PROFILE=tech-house

tech-house-prod-logs:
	@$(MAKE) docker-prod-logs PROFILE=tech-house