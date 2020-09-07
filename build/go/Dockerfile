FROM golang:latest

WORKDIR /go

RUN apt-get update && apt-get install -y \
        git \
        vim \
        net-tools

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
ENV POSTGRESQL_URL 'postgres://chatadmin:admin@store:5432/chatadmin?sslmode=disable'
WORKDIR /go/src/chat/build