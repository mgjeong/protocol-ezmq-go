#!/bin/bash
./install_dependencies_armhf_qemu.sh
export GOPATH=$PWD

#get and install required go dependencies
GOARCH=arm go get github.com/pebbe/zmq4
GOARCH=arm go get -u github.com/golang/protobuf/protoc-gen-go
GOARCH=arm go get -u go.uber.org/zap

#build ezmq SDK
mkdir ./src/go
cp -r  ezmq ./src/go
cd ./src/go/ezmq
CGO_ENABLED=1 GOOS=linux GOARCH=arm go build
CGO_ENABLED=1 GOOS=linux GOARCH=arm go install

#build samples
cd ./../../../
cp -r samples ./src/go
cd ./src/go/samples
CGO_ENABLED=1 GOOS=linux GOARCH=arm go build subscriber.go
CGO_ENABLED=1 GOOS=linux GOARCH=arm go build publisher.go
