#!/bin/bash
if [ $# != 1 ];then
   echo "USAGE: $0 v.x.x"
   exit 1
fi
touch go.mod
echo "module golang_gin
go 1.14
require (
        github.com/tencentyun/tcaplusdb-go-sdk/tdr $1
        github.com/tencentyun/tsf4g/tdrcom v0.0.0-20210518063825-e81b9a582a72
)" > go.mod
echo $?
if [ $? -ne 0 ]
then
    echo "go.mod导入文件失败"
    exit 2
fi
echo "export GO111MODULE=on
all: build
build:
	go test -c -o test
clean:
	rm -rf test" > Makefile
if [ $? -ne 0 ]
then
    echo "导入Makefile失败"
    exit 3
fi
go mod tidy
sleep 10
make
