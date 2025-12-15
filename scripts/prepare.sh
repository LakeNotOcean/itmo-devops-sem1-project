#!/bin/bash

set -e

PORT_VALUE=8080
DOCKER_IMAGE_NAME=backend

# Переменные сервиса
source ./configs/.env

echo "Starting to build the ${DOCKER_IMAGE_NAME} Docker image..."

docker build --build-arg PORT=$PORT_VALUE -t $DOCKER_IMAGE_NAME .

if [ $? -eq 0 ]; then
    echo "Docker image ${DOCKER_IMAGE_NAME} successfully built!"
else
    echo $?
    exit 1
fi