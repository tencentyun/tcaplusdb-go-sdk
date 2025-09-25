#export GO111MODULE=off
#export GOPATH:=$(shell pwd)../../../../../

all: build
build:
	go test -c -coverpkg=../../../../tcaplus-go-api/... -covermode=atomic -o pbTest
clean:
	rm -rf pbTest
