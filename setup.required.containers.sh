#!/bin/bash
cd "$(dirname "$0")"

#RUN AS ROOT!
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

NGINX_DEFAULT_HOST=$1
NGINX_DIR=$HOME/docker/nginx

if [ -z $NGINX_DEFAULT_HOST ]; then
    echo "NGINX_DEFAULT_HOST variable not set, exiting..."
    exit 1
fi

NGINX_ID=$(docker ps -q -f name=^/nginx$)

if [ -z "$NGINX_ID" ]
    then
        echo "NGINX docker container is not running, starting it..."
        docker run -d -p 80:80 -p 443:443 \
            --restart always \
            --name nginx \
            -e DEFAULT_HOST=${NGINX_DEFAULT_HOST} \
            -v $NGINX_DIR/certs:/etc/nginx/certs:ro \
            -v $NGINX_DIR/vhost.d:/etc/nginx/vhost.d \
            -v $NGINX_DIR/conf.d:/etc/nginx/conf.d \
            -v $NGINX_DIR/share/nginx/html:/usr/share/nginx/html \
            -v /var/run/docker.sock:/tmp/docker.sock:ro \
            krashid/nginx-proxy:1.11.10
    else
        echo "NGINX docker container is already running..."
fi

LETS_ENCRYPT_ID=$(docker ps -q -f name=^/letsencrypt$)

if [ -z "$LETS_ENCRYPT_ID" ]
    then
        echo "LETSENCRYPT docker container is not running, starting it..."
        docker run -d \
            --restart always \
            --name letsencrypt \
            -v $NGINX_DIR/certs:/etc/nginx/certs:rw \
            --volumes-from nginx \
            -v /var/run/docker.sock:/var/run/docker.sock:ro \
	        -e ACME_TOS_HASH=cc88d8d9517f490191401e7b54e9ffd12a2b9082ec7a1d4cec6101f9f1647e7b \
            jrcs/letsencrypt-nginx-proxy-companion
    else
        echo "LETSENCRYPT docker container is already running..."
fi


echo "Setting up required autoheal container for recovering failed instances..."
HEAL_ID=$(docker ps -aqf "name=autoheal")

if [ -z "$HEAL_ID" ]
    then
        echo "autoheal docker container is not running, starting it..."
        docker run -d \
            --restart always \
            --name autoheal \
            -e AUTOHEAL_CONTAINER_LABEL=all \
            -v /var/run/docker.sock:/var/run/docker.sock \
            willfarrell/autoheal
    else
        echo "autoheal docker container is already running..."
fi
