# Meta info
NAME := dementor
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
		-X 'main.revision=$(REVISON)'

## Setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/laher/goxc
	go get github.com/Songmu/make2help/cmd/make2help

## Install dependencies
deps: setup
	glide install

## Update dependencies
update: setup
	glide update

## Format source codes
fmt: setup
	goimports -w $$(glide nv -x)

## Test
test: deps fmt
	go test -v $$(glide novendor) -tags=unit
	go test -v $$(glide novendor) -tags=integration

## Lint
lint: setup fmt
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

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