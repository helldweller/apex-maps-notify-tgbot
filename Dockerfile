FROM golang:1.23-alpine AS build
WORKDIR /app
COPY ./src/go.sum ./src/go.mod ./
RUN time go mod download
COPY ./src ./
RUN time go build ./cmd/app/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/main /app/main
# COPY ./config.yaml ./
ENTRYPOINT /app/main
