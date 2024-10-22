# syntax=docker/dockerfile:1

FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY . .

RUN go mod tidy && \
	go install honnef.co/go/tools/cmd/staticcheck@latest
