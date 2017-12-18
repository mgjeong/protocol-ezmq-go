#!/bin/bash
./install_dependencies_armhf.sh
export GOPATH=$PWD

#get and install required go dependencies
GOOS=linux GOARCH=arm CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm'  CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 CGO_ENABLED=1 go get github.com/pebbe/zmq4
GOARCH=arm go get -u github.com/golang/protobuf/protoc-gen-go
GOARCH=arm go get -u go.uber.org/zap

#build ezmq SDK
mkdir ./src/go
cp -r  ezmq ./src/go
cd ./src/go/ezmq
CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build
CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go install

#build samples
cd ./../../../
cp -r samples ./src/go
cd ./src/go/samples
CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -a subscriber.go
CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -a publisher.go
