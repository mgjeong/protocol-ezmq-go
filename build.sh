#!/bin/bash
./install_dependencies.sh
export GOPATH=$PWD

#get and install required go dependencies
go get github.com/pebbe/zmq4
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u go.uber.org/zap

#build ezmq SDK
mkdir ./src/go
cp -r ezmq ./src/go
cp -r samples ./src/go
cd ./src/go/ezmq
go build
go install

#build samples
cd ./../../../
cp -r samples ./src/go
cd ./src/go/samples
go build -a subscriber.go
go build -a publisher.go