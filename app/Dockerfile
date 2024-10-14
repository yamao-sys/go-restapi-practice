# syntax=docker/dockerfile:1

FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY . .

RUN go install golang.org/x/tools/cmd/goimports@latest
