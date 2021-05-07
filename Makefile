VERSION ?= $(shell git describe --tags)
BUILD ?= $(shell git rev-parse --short HEAD)
PROJECTNAME := goslice
TARGET := .target
GOFILES := ./cmd/goslice
PREFIX := /usr/local
DESTDIR :=
BIN := goslice

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
	GOFLAGS=$(GOFLAGS) \
	go build $(LDFLAGS) -o $(TARGET) $(GOFILES)

## clean the build folder
clean:
	@rm -Rf .target
	@rm -f .test_stl/*.gcode

test:
	@go test ./...

.PHONY: install
install: build
	install -Dm755 $(TARGET)/$(BIN) $(DESTDIR)$(PREFIX)/bin/${BIN}

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/${BIN}

all: build