export GO111MODULE := on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

TARGET := golang-base-service
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

IMAGE_BASE_NAME := lyckety/$(TARGET)
IMAGE_TAG ?= manual
IMAGE_NAME := $(IMAGE_BASE_NAME):$(IMAGE_TAG)


VERSION ?= manual
COMMIT ?= $(shell git rev-list -1 HEAD --abbrev-commit)
BRANCH ?= manual_branch
BUILD_TIME := $(shell date +%s)
BUILDER ?= $(shell hostname)

run:
	go run cmd/$(TARGET)/main.go


################################################################################
### Builds
###

build:
	go build -ldflags -o bin/$(TARGET) ./cmd/$(TARGET)/

image:
	docker build -t $(IMAGE_NAME) \
	-f ./Dockerfile .

################################################################################
### Generators
###

buf-build:
	buf generate --output pkg/grpc

################################################################################
### Uploads & pushes
###

imagetar: image
	docker save lyckety/$(IMAGE_NAME) -o $(IMAGE_NAME).tar

################################################################################
### Tests
###

test: build
	go test -v ./... -bench=.

################################################################################
### Linters
###
linters: golangci-lint

golangci-lint: build
	test $(TARGET) != email-sender && grep -r --exclude-dir='.git' 'email-sender' | grep -v 'should be absent in normal project' && echo "all email-sender should be absent in normal project" && exit 1 || :
	find -type f -name "*.go" | grep -v '.*\.pb\.go' | grep -v '\/[a-z_]*.go' && echo "Files should be named in snake case" && exit 1 || echo "All files named in snake case"
	golangci-lint version
	golangci-lint run


buf-lint:
	buf lint

################################################################################
### Golang helpers
###

tidy:
	go mod tidy

clean:
	rm -rf bin/$(TARGET)

download:
	go mod download

strip: build
	strip bin/$(TARGET)


modup: tidy
	go get -u ./...
	go mod tidy
