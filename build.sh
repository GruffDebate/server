#!/bin/sh

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o `ls *.go | grep -v _test.go`
#CGO_ENABLED=0 GOOS=linux go build -a --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo -o server .
GOOS=linux bash build