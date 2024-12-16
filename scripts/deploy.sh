#!/bin/bash

# Variables de entorno
CONFIG_DIR="./config"

# Desplegar RabbitMQ Config
RABBITMQ_CONFIG="$CONFIG_DIR/rabbitmq/rabbitmq.conf"
RABBITMQ_DEFINITIONS="$CONFIG_DIR/rabbitmq/definitions.json"
cp $RABBITMQ_CONFIG /etc/rabbitmq/
rabbitmqadmin import $RABBITMQ_DEFINITIONS

# Desplegar Prometheus Config
PROMETHEUS_CONFIG="$CONFIG_DIR/prometheus/prometheus.yml"
cp $PROMETHEUS_CONFIG /etc/prometheus/

# Desplegar Docker Compose
cd $CONFIG_DIR/docker
docker-compose up -d

# Desplegar Ngrok Config
NGROK_CONFIG="$CONFIG_DIR/ngrok/ngrok.yml"
ngrok start --config=$NGROK_CONFIG
