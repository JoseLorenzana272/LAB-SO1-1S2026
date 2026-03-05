#!/bin/bash

IMAGES=("hello-world" "alpine")
RANDOM_IMAGE=${IMAGES[$((RANDOM % 2))]}

CONTAINER_NAME="container_$(date +%s)_$$"

echo "$(date): Creating container $CONTAINER_NAME with image $RANDOM_IMAGE"

if [ "$RANDOM_IMAGE" = "alpine" ]; then
    docker run -d --name "$CONTAINER_NAME" "$RANDOM_IMAGE" sleep 300
else
    docker run -d --name "$CONTAINER_NAME" "$RANDOM_IMAGE"
fi
