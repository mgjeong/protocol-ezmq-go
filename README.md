# ezmq library (go)

protocol-ezmq-go is a go package which provides a standard messaging interface over various data streaming
and serialization / deserialization middlewares along with some added functionalities.</br>
  - Currently supports streaming using 0mq and serialization / deserialization using protobuf.
  - Publisher -> Multiple Subscribers broadcasting.
  - Topic based subscription and data routing at source (read publisher).
  - High speed serialization and deserialization.

## Prerequisites ##
 - You must install basic prerequisites for build
   - Make sure that libtool, pkg-config, build-essential, autoconf, and automake are installed.
      ```
      $ sudo apt-get install libtool pkg-config build-essential autoconf automake
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

- You must install **libsodium**: [It is required for secured mode]
   ```
   $ sudo apt-get install libsodium-dev 
   ```

## How to build ##
1. Goto: ~/protocol-ezmq-go/</br>
2. Run the script:
   ```
   ./build_auto.sh <options>
   ```
**Notes:** </br>
(a) For getting help about script option: **$ ./build_auto.sh --help** </br>
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

### Subscriber sample [Secured] ###
1. Goto: ~/${GOPATH}/src/go/samples/
2. Run the sample:
    ```
     ./subscriber_secured
    ```
    - **It will give list of options for running the sample.** </br>
    - **Update ip, port and topic as per requirement.** </br>
    - **With secured sample unsecured features can be tested** </br>

### Publisher sample [Secured] ###

1. Goto: ~/${GOPATH}/src/go/samples/
2. Run the sample:
   ```
   ./publisher_secured
   ```
   - **It will give list of options for running the sample.** </br>
   - **Update port and topic as per requirement.** </br>
   - **With secured sample unsecured features can be tested** </br>

### Subscriber sample  ###
1. Goto: ~/${GOPATH}/src/go/samples/
2. Run the sample:
    ```
     ./subscriber
    ```
    - **It will give list of options for running the sample.** </br>
    - **Update ip, port and topic as per requirement.** </br>
    - **This sample will be built, only if ezmq package is built in unsecured mode.** </br>
    
### Publisher sample ###

1. Goto: ~/${GOPATH}/src/go/samples/
2. Run the sample:
   ```
   ./publisher
   ```
   - **It will give list of options for running the sample.** </br>
   - **Update port and topic as per requirement.** </br>    
   - **This sample will be built, only if ezmq package is built in unsecured mode.** </br>

## Unit test and code coverage report

### Pre-requisite
Built ezmq package.

### Run the unit test cases
1. Goto:  `~/protocol-ezmq-go/src/go/unittests`
2. Run the script </br>
   ` $ build.sh`

### Generate the code coverage report
1. Goto `~/protocol-ezmq-go/src/go/unittests` </br>
2. Run the script </br>
    ` $ build.sh`
3. Run the below command to open coverage report in web browser: </br>
     `$ go tool cover -html=coverage.out`

## Usage guide for ezmq library (For micro-services) ##
1. The microservice which wants to use ezmq GO library has to import ezmq package:
    `import go/ezmq`
2. Reference ezmq library APIs : [doc/godoc/ezmq.html](doc/godoc/ezmq.html)

## Running static analyzer for ezmq library ##
1. Goto: **~/${GOPATH}/src/go/**
2. Run the below command:</br>
   ```
   $ go tool vet -all .
   ```
## Future Work ##
  - High speed parallel ordered serialization / deserialization based on streaming load.
  - Threadpool for multi-subscriber handling.
  - Router pattern. For number of subscribers to single publisher use case.
  - Clustering Support.
</br></br>
