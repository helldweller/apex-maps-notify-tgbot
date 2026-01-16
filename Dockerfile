# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Speed up builds by separating mod download and sources
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src ./

# Static, trimmed binary
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app/main.go

FROM --platform=$TARGETPLATFORM alpine:3.19
WORKDIR /app

# Add CA certs if the app makes HTTPS calls
RUN apk add --no-cache ca-certificates

COPY --from=build /out/app /app/main
ENTRYPOINT ["/app/main"]
