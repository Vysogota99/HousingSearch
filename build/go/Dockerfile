FROM golang:latest

WORKDIR /go/src/school/web
WORKDIR /go


RUN apt-get update && apt-get install -y \
        git \
        vim \
        net-tools 

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
ENV POSTGRESQL_URL_AUTH 'postgres://auth_user:admin@store:5432/auth_user?sslmode=disable'
ENV POSTGRESQL_URL_MAIN 'postgres://server_user:admin@store:5432/server?sslmode=disable'

        
RUN apt install -y protobuf-compiler        
RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/school/build