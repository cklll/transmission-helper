FROM golang:1.16.2-alpine

WORKDIR /usr/src/transmission-helper

RUN apk add build-base

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
