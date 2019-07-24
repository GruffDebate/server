#!/bin/sh
go run `ls *.go | grep -v _test.go`
