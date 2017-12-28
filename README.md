# ezmq library (go)

protocol-ezmq-java is a library (jar) which provides a standard messaging interface over various data streaming
and serialization / deserialization middlewares along with some added functionalities.</br>
  - Currently supports streaming using 0mq and serialization / deserialization using protobuf.
  - Publisher -> Multiple Subscribers broadcasting.
  - Topic based subscription and data routing at source (read publisher).
  - High speed serialization and deserialization.

## Prerequisites ##
 - You must install basic prerequisites for build 
   - Install build-essential
   ```
   $ sudo apt-get install build-essential
   ```
- Python
  - Version : 2.4 to 3.0
  - [How to install](https://wiki.python.org/moin/BeginnersGuide/Download)

- SCons
  - Version : 2.3.0 or above
  - [How to install](http://scons.org/doc/2.3.0/HTML/scons-user/c95.html)
  
- Go compiler
  - Version : 1.9
  - [How to install](https://golang.org/doc/install)


## How to build ##
1. Goto: ~/protocol-ezmq-go/</br>
2. 
```
./build.sh <options>
```
**Notes:** </br>
(a) For getting help about script option: **$ ./build.sh --help** </br>
(b) Currently, Script needs sudo permission for installing zeroMQ and protobuf libraries. In future need for sudo will be removed by installing those libraries in ezmq library.
(c) While doing cross compilation, permission denied error may come.</br>
      **For example:** The below error while building for armhf:</br>
      `go install runtime/internal/sys: mkdir /usr/local/go/pkg/linux_arm: permission denied`</br>
       Do the following:</br>
       - sudo mkdir /usr/local/go/pkg/linux_arm/</br>
       - sudo chmod 777 /usr/local/go/pkg/linux_arm/</br>

## How to run ezmq samples ##

### Prerequisites ###
 Built ezmq library
 
### Subscriber sample ###
1. Goto: ~/${GOPATH}/src/go/samples/
2. 
```
./subscriber
```
3.  On successful running it will show following logs
```
2017-12-21T01:22:15.754+0530    DEBUG   ezmq/ezmqapi.go:38      EZMQ initialized
[Initialize] Error code is: 0
Enter 1 for General Event testing
Enter 2 for Topic Based delivery
```

### Publisher sample ###

1. Goto: ~/${GOPATH}/src/go/samples/
2. 
```
./subscriber
```
3.  On successful running it will show following logs
```
2017-12-21T01:23:25.754+0530    DEBUG   ezmq/ezmqapi.go:38      EZMQ initialized
[Initialize] Error code is: 0
Enter 1 for General Event testing
Enter 2 for Topic Based delivery
```

## Usage guide for ezmq library (For micro-services) ##
1. The micro-service which wants to use ezmq GO library has to import ezmq package:
    `import go/ezmq`

## Running static analyzer for ezmq library ##
1. Goto: **~/${GOPATH}/src/go/**
2. Run the below command:</br>
```
$ go tool vet -all .
```
</br>
## Future Work ##
  - High speed parallel ordered serialization / deserialization based on streaming load.
  - Threadpool for multi-subscriber handling.
  - Router pattern. For number of subscribers to single publisher use case.
  - Clustering Support.
</br></br>

