# Meta info
NAME := dementor
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
		-X 'main.revision=$(REVISON)'

## Setup
setup:
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/laher/goxc
	go get github.com/Songmu/make2help/cmd/make2help

## Install dependencies
deps: setup
	dep ensure

## Update dependencies
update: setup
	dep ensure -update

## Format source codes
fmt: setup
	goimports -w $$(find . -type f -name "*.go" -not -path "./vendor/*")

## Test
test: deps fmt
	go test -v ./... -tags=unit
	go test -v ./... -tags=integration

## Lint
lint: setup fmt
	go vet ./...
	golint $$(find . -type f -name "*.go" -not -path "./vendor/*")

## Run
run: deps fmt
	go run *.go

## Build binaries
build: deps fmt
	goxc -tasks=xc,archive -d=./bin -bc="linux,windows,darwin" -pv=$(VERSION)

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps update test lint run help build