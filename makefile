SHELL := /bin/bash

# Variables
ROOT_DIR := $(shell pwd)
VERSION := 1.0
BUILD_DIR := ${ROOT_DIR}/bin
CONFIG_DIR := ${ROOT_DIR}/config
SCRIPTS_DIR := ${ROOT_DIR}/scripts
DOCKER_COMPOSE_DEV := ${CONFIG_DIR}/docker-compose.dev.yml
DOCKER_COMPOSE_STG := ${CONFIG_DIR}/docker-compose.stg.yml

# Phony targets declaration
.PHONY: all build run test clean lint \
	docker-dev-build docker-dev-up docker-dev-down docker-dev-logs \
	docker-stg-build docker-stg-up docker-stg-down docker-stg-logs \
	tech-house-dev-build tech-house-dev-up tech-house-dev-down tech-house-dev-logs \
	tech-house-stg-build tech-house-stg-up tech-house-stg-down tech-house-stg-logs \
	build-lambda

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


build-lambda:
	# Compile for Linux (Lambda runtime)
	GOOS=linux GOARCH=amd64 go build -o ./tmp/bootstrap ./projects/customers-manager/cmd/lambda
	# Create the ZIP file for Lambda
	zip ./zip-lambda/lambda.zip ./tmp/bootstrap
	# Clean binary
	rm ./tmp/bootstrap