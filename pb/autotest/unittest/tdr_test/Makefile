# Go parameters
#export GO111MODULE=off
#export GOPATH:=$(shell pwd)/../../../../
GOCMD=go
GOBUILD=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=tdrTest

all: test
test:
	$(GOBUILD) -c -coverpkg=../../../../tcaplus-go-api/... -covermode=atomic -o $(BINARY_NAME)
