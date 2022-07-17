
## Develop

### In minikube

    cp skaffold/app-secret-example.yaml skaffold/app-secret.yaml
    # and edit secrets
    minikube start
    eval $(minikube docker-env)
    kubectl config use-context minikube
    skaffold run -p init
    skaffold dev -p app

### Local

    cd src
    go build ./cmd/app/main.go
    export APEX_API_KEY=xxx
    export LOG_LEVEL=info
    export TGBOT_API_KEY='yyy'
    export MONGODB_URI=mongodb://localhost:27017/tgbot
    ./main

### Docker

    docker run -d --rm -v "data:/var/lib/mongodb" --name mongodb docker.io/mongo:5.0
    docker build -t apex-maps-tgbot .
    docker run --rm -it \
        --env "APEX_API_KEY=xxx" \
        --env "LOG_LEVEL=info" \
        --env "TGBOT_API_KEY=yyy" \
        --env "MONGODB_URI=mongodb://mongodb:27017/tgbot"
        apex-maps-tgbot


## todo

Фичреквест: Добавть уведомление в указанный чат(или лично) о начале олимпуса (с ограничением отправки в определенные промежутки времени)

    Helm chart
    GitHub Actions:
        pipeline тестов кач-ва кода для каждого коммита master/pr
        gitops pipeline для деплоя при мерже мастера или создании релиза
        secrets можно хранить в GitHub
    