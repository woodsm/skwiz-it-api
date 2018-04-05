#!/bin/bash

NAME=skwiz-it-api

CONTAINER_ID=$(docker ps -q -f name=^/${NAME}$);

if [[ -z "${CONTAINER_ID}" ]]; then
    echo "No existing container running"
else
    echo "Removing existing container ${CONTAINER_ID} : $(docker rm -f ${NAME})"
fi

IMAGE_ID=$(docker images | grep ${NAME});

if [[ -z "${IMAGE_ID}" ]]; then
    echo "Building new docker image ${NAME}"
    docker build -t ${NAME} .
    else
    echo "Using existing docker image."
    echo "If you want to rebuild it, use \"docker rm -f ${NAME}\" to remove it and run this script again"
fi

docker run -d -p 8081:3000 --name ${NAME} --restart always ${NAME}