###############################################################################
# Copyright 2017 Samsung Electronics All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################

#!/bin/bash
set +e
#Colors
RED="\033[0;31m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
NO_COLOUR="\033[0m"

PROJECT_ROOT=$(pwd)
export GOPATH=$PWD
DEP_ROOT=$(pwd)/dependencies
EZMQ_TARGET_ARCH="$(uname -m)"
EZMQ_WITH_DEP=false
EZMQ_BUILD_MODE="release"
EZMQ_WITH_SECURITY=true

ZMQ_LIBSODIUM="yes"
IS_SECURED="secure"

install_dependencies() {
    # download required tool chain for cross compilation [arm/arm64/armhf]
    if [ "arm" = ${EZMQ_TARGET_ARCH} ]; then
        if [ -x "/usr/bin/arm-linux-gnueabi-g++" ] && [ -x "/usr/bin/arm-linux-gnueabi-gcc" ]; then
            echo -e "${BLUE}Cross compile tool-chain found for arm${NO_COLOUR}"
        else
            echo -e "${RED}No cross compile tool-chain found for arm, install using below commands:${NO_COLOUR}"
            echo -e "${BLUE} $ sudo apt-get update"
            echo -e " $ sudo apt-get install gcc-arm-linux-gnueabi"
            echo -e " $ sudo apt-get install g++-arm-linux-gnueabi${NO_COLOUR}"
            exit 0
        fi
    elif [ "arm64" = ${EZMQ_TARGET_ARCH} ]; then
        if [ -x "/usr/bin/aarch64-linux-gnu-g++-4.8" ] && [ -x "/usr/bin/aarch64-linux-gnu-gcc-4.8" ] && [ -x "/usr/bin/aarch64-linux-gnu-g++" ] && [ -x "/usr/bin/aarch64-linux-gnu-gcc" ]; then
            echo -e "${BLUE}Cross compile tool-chain found for arm64${NO_COLOUR}"
        else
            echo -e "${RED}No cross compile tool-chain found for arm64, install using below commands:${NO_COLOUR}"
            echo -e "${BLUE} $ sudo apt-get update"
            echo -e " $ sudo apt-get install gcc-4.8-aarch64-linux-gnu"
            echo -e " $ sudo apt-get install g++-4.8-aarch64-linux-gnu"
            echo -e " $ cd /usr/bin"
            echo -e " $ sudo ln -s aarch64-linux-gnu-g++-4.8 aarch64-linux-gnu-g++"
            echo -e " $ sudo ln -s aarch64-linux-gnu-gcc-4.8 aarch64-linux-gnu-gcc${NO_COLOUR}"
            exit 0
        fi
    elif [ "armhf" = ${EZMQ_TARGET_ARCH} ]; then
        if [ -x "/usr/bin/arm-linux-gnueabihf-g++-4.8" ] && [ -x "/usr/bin/arm-linux-gnueabihf-gcc-4.8" ] && [ -x "/usr/bin/arm-linux-gnueabihf-g++" ] && [ -x "/usr/bin/arm-linux-gnueabihf-gcc" ]; then
            echo -e "${BLUE}Cross compile tool-chain found for armhf${NO_COLOUR}"
        else
            echo -e "${RED}No cross compile tool-chain found for armhf, install using below commands:${NO_COLOUR}"
            echo -e "${BLUE} $ sudo apt-get update"
            echo -e " $ sudo apt-get install gcc-4.8-arm-linux-gnueabihf"
            echo -e " $ sudo apt-get install g++-4.8-arm-linux-gnueabihf"
            echo -e " $ cd /usr/bin"
            echo -e " $ sudo ln -s arm-linux-gnueabihf-g++-4.8 arm-linux-gnueabihf-g++"
            echo -e " $ sudo ln -s arm-linux-gnueabihf-gcc-4.8 arm-linux-gnueabihf-gcc${NO_COLOUR}"
            exit 0
        fi
    fi

    # Check, clone, build and install zeroMQ and protobuf libraries
    cd ./dependencies
    DEP_ROOT=$(pwd)

    # Clone, build and install ZeroMQ library
    if [ -d "./libzmq" ] ; then
        echo "ZMQ library folder exist"
    else
        git clone https://github.com/zeromq/libzmq.git
    fi
    cd libzmq
    git checkout v4.2.2
    chmod +x version.sh
    ./version.sh
    chmod +x autogen.sh
    ./autogen.sh
    if [ "arm" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}ZeroMQ configuring for arm${NO_COLOUR}"
        ./configure --host=arm-none-linux-gnueabi CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ --with-libsodium=${ZMQ_LIBSODIUM}
    elif [ "arm64" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}ZeroMQ configuring for arm64${NO_COLOUR}"
        ./configure --host=aarch64-unknown-linux-gnu CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 --with-libsodium=${ZMQ_LIBSODIUM}
    elif [ "armhf" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}ZeroMQ configuring for armhf${NO_COLOUR}"
        ./configure --host=arm-linux-gnueabihf CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 --with-libsodium=${ZMQ_LIBSODIUM}
    else
        ./configure --with-libsodium=${ZMQ_LIBSODIUM}
    fi
    make -j 4
    sudo make install
    sudo ldconfig

    # build, install protobuf library
    FILENAME="protobuf-cpp-3.4.0.tar.gz"
    cd $DEP_ROOT
    if [ -e "$FILENAME" ]; then
        echo "Protobuf tar exist"
    else
        wget https://github.com/google/protobuf/releases/download/v3.4.0/protobuf-cpp-3.4.0.tar.gz
    fi

    if [ -d "./protobuf-3.4.0" ]; then
        echo "Protobuf library folder exist"
    else
        tar -xvf protobuf-cpp-3.4.0.tar.gz
    fi

    cd protobuf-3.4.0/
    chmod +x autogen.sh
    ./autogen.sh
    if [ "arm" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Protobuf configuring for arm${NO_COLOUR}"
        ./configure --host=arm-linux CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++
        make -j 4
    elif [ "arm64" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Protobuf configuring for arm64${NO_COLOUR}"
        ./configure --host=aarch64-unknown-linux-gnu CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8
        make -j 4
    elif [ "armhf" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Protobuf configuring for armhf${NO_COLOUR}"
        ./configure --host=arm-linux-gnueabihf CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8
        make -j 4
    else
        ./configure
        make -j 4
        sudo make install
    fi

    #handle protobuf error for cross compilation
    if [ "arm" = ${EZMQ_TARGET_ARCH} ] || [ "arm64" = ${EZMQ_TARGET_ARCH} ] || [ "armhf" = ${EZMQ_TARGET_ARCH} ] || [ "armhf-qemu" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Copying libs from protobuf built directory to /usr/local/lib${NO_COLOUR}"
        sudo cp src/.libs/protoc /usr/bin
        sudo cp src/.libs/* /usr/local/lib
    fi
    sudo ldconfig

    # install ezMQ go pkg dependencies
    if [ "arm" = ${EZMQ_TARGET_ARCH} ] || [ "armhf-qemu" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Installing zmq4, protoc and zap for arm/armhf-qemu${NO_COLOUR}"
        GOARCH=arm go get github.com/pebbe/zmq4
        GOARCH=arm go get -u github.com/golang/protobuf/protoc-gen-go
        GOARCH=arm go get -u go.uber.org/zap
    elif [ "arm64" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Installing zmq4, protoc and zap for arm64${NO_COLOUR}"
        GOARCH=arm64 go get github.com/pebbe/zmq4
        GOARCH=arm64 go get -u github.com/golang/protobuf/protoc-gen-go
        GOARCH=arm64 go get -u go.uber.org/zap
        make -j 4
    elif [ "armhf" = ${EZMQ_TARGET_ARCH} ]; then
        echo -e "${BLUE}Installing zmq4, protoc and zap for armhf${NO_COLOUR}"
        GOOS=linux GOARCH=arm CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm'  CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 CGO_ENABLED=1 go get github.com/pebbe/zmq4
        GOARCH=arm go get -u github.com/golang/protobuf/protoc-gen-go
        GOARCH=arm go get -u go.uber.org/zap
        make -j 4
    else
        echo -e "${BLUE}Installing zmq4, protoc and zap for x86/x86_64/armhf-native${NO_COLOUR}"
        go get github.com/pebbe/zmq4
        go get -u github.com/golang/protobuf/protoc-gen-go
        go get -u go.uber.org/zap
    fi
    echo -e "${GREEN}Install dependencies done${NO_COLOUR}"
}

build_x86_and_64() {
    cd $PROJECT_ROOT/src/go/
    #build ezmq SDK
    cd ./ezmq
    go build -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" 
    go install
    
    #build samples
    cd ../samples
    if [ ${EZMQ_WITH_SECURITY} = true ]; then
        go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber_secured.go
        go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher_secured.go  
    else
        go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber.go
        go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher.go  
    fi
}

build_arm() {
    cd $PROJECT_ROOT/src/go/
    #build ezmq SDK
    cd ./ezmq
    CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" 
    CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go install
    #build samples
    cd ../samples
    
    if [ ${EZMQ_WITH_SECURITY} = true ]; then
        CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber_secured.go
        CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher_secured.go 
    else
        CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber.go
        CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher.go
    fi
    
}

build_arm64() {
    cd $PROJECT_ROOT/src/go/
    #build ezmq SDK
    cd ./ezmq
    CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go build -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" 
    CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go install
    #build samples
    cd ../samples
    if [ ${EZMQ_WITH_SECURITY} = true ]; then
        CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}"subscriber_secured.go
        CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher_secured.go
    else
        CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber.go
        CGO_ENABLED=1 CC=/usr/bin/aarch64-linux-gnu-gcc-4.8 CXX=/usr/bin/aarch64-linux-gnu-g++-4.8 GOOS=linux GOARCH=arm64 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher.go
    fi
}

build_armhf() {
    cd $PROJECT_ROOT/src/go/
    #build ezmq SDK
    cd ./ezmq
    CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}"
    CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go install
    #build samples
    cd ../samples
    if [ ${EZMQ_WITH_SECURITY} = true ]; then
        CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber_secured.go
        CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher_secured.go 
    else
        CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber.go
        CGO_LDFLAGS+='-Bstatic -lzmq -lprotobuf -Bdynamic -lstdc++ -lm' GOOS=linux GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc-4.8 CXX=arm-linux-gnueabihf-g++-4.8 go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher.go
    fi  
}

build_armhf_native() {
    cd $PROJECT_ROOT/src/go/
    #build ezmq SDK
    cd ./ezmq
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}"
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go install
    #build samples
    cd ../samples
    if [ ${EZMQ_WITH_SECURITY} = true ]; then
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber_secured.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher_secured.go 
    else
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher.go
    fi
}

build_armhf_qemu() {
    cd $PROJECT_ROOT/src/go/
    #build ezmq SDK
    cd ./ezmq
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" 
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go install
    #build samples
    cd ../samples
    if [ ${EZMQ_WITH_SECURITY} = true ]; then
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber_secured.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher_secured.go 
    else
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" subscriber.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQ_BUILD_MODE} ${IS_SECURED}" publisher.go
    fi
}

clean_ezmq() {
    echo -e "Cleaning ${BLUE}${PROJECT_ROOT}${NO_COLOUR}"
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/src/${NO_COLOUR}"
    rm -rf ./src
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/pkg/${NO_COLOUR}"
    rm -rf ./pkg
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/bin/${NO_COLOUR}"
    rm -rf ./bin
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/dependencies/${NO_COLOUR}"
    rm -rf ./dependencies/libzmq
    rm -rf ./dependencies/protobuf-3.4.0
    echo -e "Finished Cleaning ${BLUE}${EZMQ}${NO_COLOUR}"
}

usage() {
    echo -e "${BLUE}Usage:${NO_COLOUR} ./build_auto.sh <option>"
    echo -e "${GREEN}Options:${NO_COLOUR}"
    echo "  --target_arch=[x86|x86_64|arm|arm64|armhf|armhf-qemu|armhf-native] :  Choose Target Architecture"
    echo "  --with_dependencies=[true|false](default: false)                               :  Build ezmq along with dependencies [zmq and protobuf]"
    echo "  --build_mode=[release|debug](default: release)                     :  Build ezmq library and samples in release or debug mode"
    echo "  --with_security=[true|false](default: true)                        :  Build ezmq library with or without Security feature"
    echo "  -c                                                                 :  Clean ezmq Repository and its dependencies"
    echo "  -h / --help                                                        :  Display help and exit [Be careful it will also remove GOPATH:src, pkg and bin]"
    echo -e "${GREEN}Notes: ${NO_COLOUR}"
    echo "  - While building newly for any architecture use -with_dependencies=true option."
    echo "  - Make sure libsodium is installed, it can be installed using:"
    echo "    $ sudo apt-get install libsodium-dev"
}

build_ezmq() {
    echo -e "${GREEN}Target Arch is: $EZMQ_TARGET_ARCH${NO_COLOUR}" 
    echo -e "${GREEN}Build mode is: $EZMQ_BUILD_MODE${NO_COLOUR}"
    echo -e "${GREEN}Is security enabled: $EZMQ_WITH_SECURITY${NO_COLOUR}"
    echo -e "${GREEN}Build with depedencies: ${EZMQ_WITH_DEP}${NO_COLOUR}"
    #dependencies
    if [ ${EZMQ_WITH_DEP} = true ]; then
        install_dependencies
    fi
    if [ ${EZMQ_WITH_SECURITY} = false ]; then
        ZMQ_LIBSODIUM="no"
        IS_SECURED="unsecure"
    fi
    cd $PROJECT_ROOT
    if [ -d "./src/go" ] ; then
        echo "src/go folder exits"
    else
        mkdir ./src/go
    fi
    #copy ezmq SDK files
    cp -r ezmq ./src/go
    #copy ezmq samples
    cp -r samples ./src/go
    # Copy unit test cases
    cp -r unittests ./src/go

    if [ "x86" = ${EZMQ_TARGET_ARCH} ]; then
        build_x86_and_64; exit 0;
    elif [ "x86_64" = ${EZMQ_TARGET_ARCH} ]; then
         build_x86_and_64; exit 0;
    elif [ "arm" = ${EZMQ_TARGET_ARCH} ]; then
         build_arm; exit 0;
    elif [ "arm64" = ${EZMQ_TARGET_ARCH} ]; then
         build_arm64; exit 0;
    elif [ "armhf" = ${EZMQ_TARGET_ARCH} ]; then
         build_armhf; exit 0;
    elif [ "armhf-native" = ${EZMQ_TARGET_ARCH} ]; then
         build_armhf_native; exit 0;
    elif [ "armhf-qemu" = ${EZMQ_TARGET_ARCH} ]; then
         build_armhf_qemu; exit 0;
    else
         echo -e "${RED}Not a supported architecture${NO_COLOUR}"
         usage; exit 1;
    fi
    if [ $? -ne 0 ]; then 
        echo -e "${RED}Build failed${NO_COLOUR}" 
        exit 1 
    fi
}

process_cmd_args() {
    if [ "$#" -eq 0  ]; then
        echo -e "No argument.."
        usage; exit 1
    fi

    while [ "$#" -gt 0  ]; do
        case "$1" in
            --with_dependencies=*)
                EZMQ_WITH_DEP="${1#*=}";
                if [ ${EZMQ_WITH_DEP} != true ] && [ ${EZMQ_WITH_DEP} != false ]; then
                    echo -e "${RED}Unknown option for --with_dependencies${NO_COLOUR}"
                    shift 1; exit 0
                fi
                shift 1;
                ;;
            --target_arch=*)
                EZMQ_TARGET_ARCH="${1#*=}";
                shift 1
                ;;
            --build_mode=*)
                EZMQ_BUILD_MODE="${1#*=}";
                shift 1;
                ;;
            --with_security=*)
                EZMQ_WITH_SECURITY="${1#*=}";
                if [ ${EZMQ_WITH_SECURITY} != true ] && [ ${EZMQ_WITH_SECURITY} != false ]; then
                    echo -e "${RED}Unknown option for --with_security${NO_COLOUR}"
                    shift 1; exit 0
                fi              
                shift 1;
                ;; 
            -c)
                clean_ezmq
                shift 1; exit 0
                ;;
            -h)
                usage; exit 0
                ;;
            --help)
                usage; exit 0
                ;;
            -*)
                echo -e "${RED}"
                echo "unknown option: $1" >&2;
                echo -e "${NO_COLOUR}"
                usage; exit 1
                ;;
             *)
                echo -e "${RED}"
                echo "unknown option: $1" >&2;
                echo -e "${NO_COLOUR}"
                usage; exit 1
                ;;
        esac
    done
}

process_cmd_args "$@"
echo -e "Building ezMQ library("${EZMQ_TARGET_ARCH}").."
build_ezmq
echo -e "Done building ezMQ library("${EZMQ_TARGET_ARCH}")"
