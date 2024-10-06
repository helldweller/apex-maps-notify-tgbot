#!/bin/bash

docker build -t ghcr.io/helldweller/apex-maps-notify-tgbot:${IMAGE_TAG} .
docker push  ghcr.io/helldweller/apex-maps-notify-tgbot:${IMAGE_TAG}
