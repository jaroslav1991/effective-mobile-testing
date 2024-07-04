FROM golang:1.21-alpine as builder

RUN apk update && apk upgrade && apk add --no-cache git openssh

RUN mkdir /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app

COPY . .

COPY ./migrations /go/src/app/migrations

RUN go build -o service ./cmd/

ENTRYPOINT "/go/src/app/service"