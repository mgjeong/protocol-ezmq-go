#!/bin/bash
./install_dependencies.sh
export GOPATH=$PWD
go get github.com/pebbe/zmq4
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u go.uber.org/zap

mkdir ./src
mkdir ./src/go
cp -r ezmq ./src/go
cp -r samples ./src/go
cd ./src/go/ezmq
go build
go install
