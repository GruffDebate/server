run:
	bra run

install:
	gopm install && gopm build

up:
	go run `ls *.go | grep -v _test.go`

test:
	go generate && go build -v && go test ./api -v && go vet

perfomance1:
	GOMAXPROCS=1 go test -bench=NetHTTPServerGet -benchmem -benchtime=10s

perfomance4:
	GOMAXPROCS=4 go test -bench=NetHTTPServerGet -benchmem -benchtime=10s