#!/bin/bash
./install_dependencies_arm.sh
export GOPATH=$PWD

GOARCH=arm go get github.com/pebbe/zmq4
GOARCH=arm go get -u github.com/golang/protobuf/protoc-gen-go
GOARCH=arm go get -u go.uber.org/zap

mkdir ./src
mkdir ./src/go
cp -r  ezmq ./src/go
cd ./src/go/ezmq
CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build
CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go install
