#!/bin/bash
./install_dependencies_arm64.sh
export GOPATH=$PWD

GOARCH=arm64 go get github.com/pebbe/zmq4
GOARCH=arm64 go get -u github.com/golang/protobuf/protoc-gen-go
GOARCH=arm64 go get -u go.uber.org/zap

mkdir ./src
mkdir ./src/go
cp -r  ezmq ./src/go
cd ./src/go/ezmq

CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go build
CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go install

