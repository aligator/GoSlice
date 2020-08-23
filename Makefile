VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := goslice
TARGET := .target
GOFILES := ./cmd/goslice

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

## build: Compile the binary.
build: clean
	@mkdir -p $(TARGET)
	@GOPATH=$(GOPATH) \
	GOBIN=$(GOBIN) \
	GOARM=$(GOARM) \
	CGO_CPPFLAGS=$(CGO_CPPFLAGS) \
	CGO_CFLAGS=$(CGO_CFLAGS) \
	CGO_CXXFLAGS=$(CGO_CXXFLAGS) \
	CGO_LDFLAGS=$(CGO_LDFLAGS) \
	go build $(LDFLAGS) -o $(TARGET) $(GOFILES)

## clean the build folder
clean:
	@rm -Rf .target
	@rm -f .test_stl/*.gcode

test:
	@go test ./...

all: build