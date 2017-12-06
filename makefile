up:
	go run `ls *.go | grep -v _test.go`

test:
	go test ./gruff ./api -v

up2:
	up start --address :8080