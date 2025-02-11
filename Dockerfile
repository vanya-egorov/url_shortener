FROM golang:1.22-alpine AS builder

WORKDIR /app

# Environment variables which CompileDaemon requires to run
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

#RUN apk --no-cache add bash git make gcc gettext

COPY  ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . ./

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build cmd/url-shortener/main.go" -command="./main"
