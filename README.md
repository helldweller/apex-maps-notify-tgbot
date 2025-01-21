
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

## Helm

    kubectl create ns apex-maps-notify-tgbot
    kubectl create secret generic apex-maps-notify-tgbot -n apex-maps-notify-tgbot \
        --from-literal=TGBOT_API_KEY=xxx \
        --from-literal=APEX_API_KEY=yyyy \
        --from-literal=TGBOT_CHAT_ID=-11111111
    # OR
    kubectl create secret generic apex-maps-notify-tgbot \
        --from-env-file=secret.env --dry-run=client --output=yaml > chart/templates/secret.yaml

    helm upgrade -i apex-maps-notify-tgbot ./chart \
        --namespace apex-maps-notify-tgbot \
        --wait \
        --atomic \
        --cleanup-on-fail

## todo

* Map image upload
    * Change bot to https://github.com/go-telegram/bot
* Notification about the start of the map in a user-defined period of time
    * Need some db
* Unit tests
