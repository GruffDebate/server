up:
	go run `ls *.go | grep -v _test.go`

test:
	go generate && go build -v && go test ./api -v && go vet

up2:
	up start --address :8080