export GO111MODULE := on
export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=arm64

TARGET := golang-service

run:
	go run cmd/$(TARGET)/main.go

buf-generate:
	buf generate --output pkg/grpc

test: build
	go test -v ./... -bench=.

linters: golangci-lint buf-lint

golangci-lint:
	golangci-lint version
	golangci-lint run

buf-lint:
	buf lint

tidy:
	go mod tidy

clean:
	rm -rf bin/$(TARGET)

download:
	go mod download

modup: tidy
	go get -u ./...
	go mod tidy
