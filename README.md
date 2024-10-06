
## Develop

### In minikube

    cp skaffold/app-secret-example.yaml skaffold/app-secret.yaml
    # and edit secrets
    minikube start
    kubectl apply -f skaffold/app-secret.yaml
    eval $(minikube docker-env)
    skaffold dev

### Local

    cd src
    go build ./cmd/app/main.go
    APEX_API_KEY=xxx LOG_LEVEL=info TGBOT_API_KEY='yyy' ./main

### Docker

    docker build -t apex-maps-tgbot .
    docker run --rm -it --env "APEX_API_KEY=xxx" --env "LOG_LEVEL=info" --env "TGBOT_API_KEY=yyy" apex-maps-tgbot


## todo

* Map image upload
    * Change bot to https://github.com/go-telegram/bot
* Notification about the start of the map in a user-defined period of time
    * Need some db
* Helm chart
    * External secrets from Github
* Unit tests
