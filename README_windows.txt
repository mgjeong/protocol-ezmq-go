####### Building ezMQ GO for windows platform [64-bit : amd64]:

Pre-requisites:
1) Install Go 1.9.4. ( https://golang.org/doc/install )
2) Install Git. ( https://git-scm.com/download/win )
3) Install minGW ( https://sourceforge.net/projects/mingw-w64/files/?source=navbar ) 
	> Select posix thread model when asked in intallation setup.
4) Install libsodium.
	> Download libsodium-1.0.16-mingw.tar.gz from https://download.libsodium.org/libsodium/releases/
	> Copy libsodium-win64/include to <mingw_installation>\x86_64-w64-mingw32\include
	> Copy libsodium-win64/lib/ to <mingw_installation>\x86_64-w64-mingw32\lib
	#Above copying is required as for now we can not modify zmq makefile to point to appropriate path because of license issue
5) Update your path variable to point to above 3 installations, Ex in your command prompt: 
	set PATH=%PATH%;C:\Program Files\mingw-w64\mingw64\bin

Building ezMQ dependencies:
   (a) libzmq:
       $ cd ~/protocol-ezmq-go/dependencies
       $ git clone https://github.com/zeromq/libzmq.git
       $ cd libzmq
       $ git checkout v4.2.2
       $ cd builds\mingw32
       $ mingw32-make all -f Makefile.mingw32
	   
       On succesful build it will create libzmq.dll and libzmq.dll.a in current directory.
       $ cd ~/protocol-ezmq-go/
       $ set CGO_CPPFLAGS=-I%CD%\dependencies\libzmq\include 
       $ set CGO_LDFLAGS=-L%CD%\dependencies\libzmq\builds\mingw32
       $ go get github.com/pebbe/zmq4
	   
   (b) protobuf:
       $ cd ~/protocol-ezmq-go/
       $ go get -u github.com/golang/protobuf/protoc-gen-go

   (b) zap logger:
       $ cd ~/protocol-ezmq-go/
       $ go get -u go.uber.org/zap

Building ezMQ:
1) Check your GOPATH with command: 
       $ go env
2) Copy protocol-ezmq-go/* to <YOUR_GOPATH>/src/go/
3)     $ cd <YOUR_GOPATH>/src/go/ezmq/
       $ go build -tags=debug
       $ go install
4) Build samples:
       $ cd ../samples
       $ go build -a -tags=debug subscriber.go
       $ go build -a -tags=debug publisher.go
