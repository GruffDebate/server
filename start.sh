#!/bin/sh

GRUFF_DB="dbname=gruff sslmode=disable" go run `ls *.go | grep -v _test.go`
