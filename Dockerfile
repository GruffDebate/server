# FROM scratch
# ADD gruff-server /gruff-server
# ENTRYPOINT ["/gruff-server"]

FROM golang:1.9.2-alpine

LABEL brunoksato <bruno.sato@live.com> and timothy.high <timothy.high@gmail.com>

RUN apk add --no-cache g++ glide

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir -p /go/src/github.com/bigokro/gruff-server

ADD . /go/src/github.com/bigokro/gruff-server

WORKDIR /go/src/github.com/bigokro/gruff-server
RUN glide install
    
RUN go install github.com/bigokro/gruff-server

ENTRYPOINT ["/go/bin/gruff-server"]