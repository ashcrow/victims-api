# Passed into cmd/main.go at build time
VERSION := $(shell cat ./VERSION)
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date +%s)

# Used in tagging images
IMAGE_VERSION_TAG := victims-api:$(VERSION)
IMAGE_DATE_TAG := victims-api:$(BUILD_TIME)

# Used during all builds
LDFLAGS := -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildTime=${BUILD_TIME}

.PHONY: help deps victims-api static-victims-api clean image gofmt golint lint

default: help

help:
	@echo "Targets:"
	@echo " deps: Install dependencies with govendor"
	@echo "	victims-api: Builds a victims-api binary"
	@echo "	clean: cleans up and removes built files"
	@echo "	image: builds a container image"

deps:
	go get github.com/kardianos/govendor
	govendor sync

victims-api:
	govendor build -ldflags '${LDFLAGS}' -o victims-api cmd/main.go

static-victims-api:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 govendor build --ldflags '-extldflags "-static" ${LDFLAGS}' -a -o victims-api cmd/main.go

clean:
	go clean
	rm -f victims-api

image: clean deps static-victims-api
	sudo docker build -t $(IMAGE_VERSION_TAG) -t $(IMAGE_DATE_TAG) .

gofmt:
	gofmt -l api/ cmd/ types/

golint:
	go get github.com/golang/lint/golint
	golint api/ cmd/ types/

lint: gofmt golint
