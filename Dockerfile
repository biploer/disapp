FROM golang:1.23.1-alpine AS builder

WORKDIR /usr/local/src/disapp

RUN apk --no-cache add --update bash git make gcc musl-dev gettext npm

# dependencies
COPY ["Makefile", "go.mod", "go.sum", "./"]
RUN go mod download

# build
ADD cmd ./cmd
ADD internal ./internal
ADD config ./config
ADD web ./web
RUN cd web && npm install && cd ..
RUN go build -o ./bin/disapp cmd/disapp/main.go
