
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

Фичреквест: Добавть уведомление в указанный чат(или лично) о начале олимпуса (с ограничением отправки в определенные промежутки времени)

    Helm chart
    GitHub Actions:
        pipeline тестов кач-ва кода для каждого коммита master/pr
        gitops pipeline для деплоя при мерже мастера или создании релиза
        secrets можно хранить в GitHub
    