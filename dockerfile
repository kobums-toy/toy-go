FROM golang:1.21 AS builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/gofiber.linux main.go


FROM        alpine

COPY ./bin/gofiber.linux /usr/local/main/main
COPY ./config/config.json /usr/local/main/config/config.json
CMD mkdir -p /usr/local/main/webdata
#ADD ./assets /usr/local/main/assets
#ADD ./views /usr/local/main/views

WORKDIR /usr/local/main
CMD ./main