FROM golang:1.12.4-alpine as build_base

RUN apk update && apk upgrade && \
  apk add --no-cache bash git openssh g++ glide ca-certificates

RUN adduser -D -g '' appuser

RUN mkdir -p /go/src/github.com/GruffDebate/server
ADD . /go/src/github.com/GruffDebate/server
WORKDIR /go/src/github.com/GruffDebate/server
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS server_builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server .

FROM alpine:latest as gruff
COPY --from=server_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server_builder /etc/passwd /etc/passwd
COPY --from=server_builder /go/src/github.com/GruffDebate/server/server /app/
WORKDIR /app
EXPOSE 8080
CMD ["./server"]
