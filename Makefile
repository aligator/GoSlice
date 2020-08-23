VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)" | awk '{print tolower($$0)}')

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/.target

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

## build: Compile the binary.
build:
	@echo compile
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOARM=$(GOARM) CGO_CPPFLAGS=$(CGO_CPPFLAGS) CGO_CFLAGS=$(CGO_CFLAGS) CGO_CXXFLAGS=$(CGO_CXXFLAGS) CGO_LDFLAGS=$(CGO_LDFLAGS) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

## clear the build folder
clear:
	@rm -R .target
	@rm .test_stl/*.gcode

test:
	@go test ./...

all: build