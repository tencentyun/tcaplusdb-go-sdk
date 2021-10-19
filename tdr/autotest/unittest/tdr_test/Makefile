# Go parameters
export GOPATH:=$(shell pwd)/../../../../
GOCMD=go
GOBUILD=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=apiTest

all: test
test:
	$(GOBUILD) -c -o $(BINARY_NAME) 
