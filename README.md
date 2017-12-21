# ezmq GO Client Library

protocol-ezmq-go library provides a standard messaging interface over various data streaming 
and serialization / deserialization middlewares along with some added functionalities.</br>
Following is the architecture of ezmq client library: </br> </br>
![ezmq Architecture](doc/images/ezMQ_architecture_0.1.png?raw=true "ezmq Arch")

## Features:
* Currently supports streaming using 0mq and serialization / deserialization using protobuf.
* Publisher -> Multiple Subscribers broadcasting.
* Topic based subscription and data routing at source (read publisher).
* High speed serialization and deserialization.

## Future Work:
* High speed parallel ordered serialization / deserialization based on streaming load.
* Threadpool for multi-subscriber handling.
* Router pattern. For number of subscribers to single publisher use case.
* Clustering Support.
</br></br>

## How to build ezmq library and samples
### pre-requisites
1. go 1.9 should be installed on linux machine. </br>
   **Refer:** https://golang.org/doc/install

### Build Instructions
1. Goto: ~/protocol-ezmq-go/</br>
2. ./build.sh <options></br></br>
**Notes:** </br>
(i) For getting help about script: **$ ./build.sh --help**</br>
(ii)Currently, Script needs sudo permission for installing zeroMQ and protobuf libraries. In future need for sudo will be removed by installing those libraries in ezmq library.</br>
(iii)While doing cross compilation, permission denied error may come.</br>
      **For example:** The below error while building for armhf:</br>
      `go install runtime/internal/sys: mkdir /usr/local/go/pkg/linux_arm: permission denied`</br>
       Do the following:</br>
       - sudo mkdir /usr/local/go/pkg/linux_arm/</br>
       - sudo chmod 777 /usr/local/go/pkg/linux_arm/</br>

## How to run ezmq samples

### pre-requisites
Built ezmq
### Run the subscriber sample application

1. Goto: ~/${GOPATH}/src/go/samples/
2. ./subscriber
3.  On successful running it will show following logs:

```
2017-12-21T01:22:15.754+0530    DEBUG   ezmq/ezmqapi.go:38      EZMQ initialized
[Initialize] Error code is: 0
Enter 1 for General Event testing
Enter 2 for Topic Based delivery
```
**Follow the instructions on the screen.**

###  Run the publisher sample application

1. Goto: ~/${GOPATH}/src/go/samples/
2. ./subscriber
3.  On successful running it will show following logs:

```
2017-12-21T01:23:25.754+0530    DEBUG   ezmq/ezmqapi.go:38      EZMQ initialized
[Initialize] Error code is: 0
Enter 1 for General Event testing
Enter 2 for Topic Based delivery
```
**Follow the instructions on the screen.**

## ezmq Usage guide [For micro-services]
1. The micro-service which wants to use ezmq GO library has to import ezmq package:
    `import go/ezmq`
2. Refer ezmq sample applications to use ezmq library APIs. **[~/EZMQ/go/samples]**

## Generating godoc for ezmq library 
1. Goto: **~/${GOPATH}/src/go/**
2.  Use following command to generate godoc: </br>
   $ godoc -html go/ezmq  > ezmq.html</br>
     **Note:** godoc can be generated only after building and installing ezmq.
3. Open the ezmq.html in web browser. </br>
    **Note:** Refer [guide]( https://godoc.org/golang.org/x/tools/cmd/godoc) for trying more options.

## Running static analyzer for ezmq library
1. Goto: **~/${GOPATH}/src/go/**
2. Run the below command:</br>
   $ go tool vet -all .
