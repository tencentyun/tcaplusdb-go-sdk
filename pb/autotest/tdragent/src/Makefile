# Go parameters
export GO111MODULE=off
export GOPATH:=${GOPATH}:$(shell pwd)/../../../../
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=tdrAgent

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) main.go
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
