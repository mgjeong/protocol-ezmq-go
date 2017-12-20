# ezMQ GO Client SDK

protocol-ezmq-go library provides a standard messaging interface over various data streaming 
and serialization / deserialization middlewares along with some added functionalities.</br>
Following is the architecture of ezMQ client library: </br> </br>
![ezMQ Architecture](doc/images/ezMQ_architecture_0.1.png?raw=true "ezMQ Arch")

## Features:
* Currently supports streaming using 0mq and serialization / deserialization using protobuf.
* Publisher -> Multiple Subscribers broadcasting.
* Topic based subscription and data routing at source (read publisher).
* High speed serialization and deserialization.

## Future Work:
* High speed parallel ordered serialization / deserialization based on streaming load.
* Threadpool for multi-subscriber handling.
* Router pattern.
* Clustering Support.
</br></br>

## How to build ezMQ SDK and samples
### pre-requisites
1. go 1.9 should be installed on linux machine. </br>
   **Refer:** https://golang.org/doc/install

### Build Instructions
1. Goto: ~/protocol-ezmq-go/
2. ./build.sh </br>
   **Note:** For getting help about script: **$ ./build.sh --help**

## How to run ezMQ samples

### pre-requisites
Built ezMQ
### Run the subscriber sample application

1. Goto: ~/${GOPATH}/src/go/linux/samples/
2. ./subscriber
3.  On successful running it will show following logs:

```
Initialize API [result]: 0
```
**Follow the instructions on the screen.**

###  Run the publisher sample application

1. Goto: ~/${GOPATH}/src/go/linux/samples/
2. ./subscriber
3.  On successful running it will show following logs:

```
Initialize API [result]: 0
```
**Follow the instructions on the screen.**

