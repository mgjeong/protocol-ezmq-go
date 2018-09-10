/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

package main

import (
	ezmq "go/ezmq"

	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const ServerSecretKey = "[:X%Q3UfY+kv2A^.wv:(qy2E=bk0L][cm=mS3Hcx"

func getEvent() ezmq.Event {
	var event ezmq.Event

	var id string = "id1"
	event.Id = &id
	var created int64 = 1
	event.Created = &created
	var modified int64 = 2
	event.Modified = &modified
	var origin int64 = 3
	event.Origin = &origin
	var pushed int64 = 4
	event.Pushed = &pushed
	var device string = "device1"
	event.Device = &device

	//form the reading
	var reading1 *ezmq.Reading = &ezmq.Reading{}
	var rId string = "id1"
	reading1.Id = &rId
	var rCreated int64 = 1
	reading1.Created = &rCreated
	var rModified int64 = 2
	reading1.Modified = &rModified
	var rOrigin int64 = 3
	reading1.Origin = &rOrigin
	var rPushed int64 = 4
	reading1.Pushed = &rPushed
	var rDevice string = "device1"
	reading1.Device = &rDevice
	var rName string = "temperature"
	reading1.Name = &rName
	var rValue string = "20"
	reading1.Value = &rValue

	event.Reading = make([]*ezmq.Reading, 1)
	event.Reading[0] = reading1
	return event
}

func getByteDataEvent() ezmq.EZMQByteData {
	var bytes ezmq.EZMQByteData
	byteArray := [5]byte{0x40, 0x05, 0x10, 0x11, 0x12}
	bytes.ByteData = byteArray[:]
	return bytes
}

func printError() {
	fmt.Printf("\nRe-run the application as shown in below example: \n")
	fmt.Printf("\n  (1) For publishing without topic: ")
	fmt.Printf("\n      ./publisher_secured -port 5562\n")
	fmt.Printf("\n  (2) For publishing without topic: [Secured] ")
	fmt.Printf("\n      ./publisher_secured -port 5562 -secured 1\n")
	fmt.Printf("\n  (3) For publishing with topic: ")
	fmt.Printf("\n      ./publisher_secured -port 5562 -t topic1 \n")
	fmt.Printf("\n  (4) For publishing with topic: [Secured] ")
	fmt.Printf("\n      ./publisher_secured -port 5562 -t topic1 -secured 1\n")
	os.Exit(-1)
}

func main() {
	var port int
	var topic string
	var isSecured int = 0
	var result ezmq.EZMQErrorCode
	var publisher *ezmq.EZMQPublisher = nil
	var instance *ezmq.EZMQAPI = nil
	startCB := func(code ezmq.EZMQErrorCode) { fmt.Printf("startCB") }
	stopCB := func(code ezmq.EZMQErrorCode) { fmt.Printf("stopCB") }
	errorCB := func(code ezmq.EZMQErrorCode) { fmt.Printf("errorCB") }

	// get port from command line arguments
	if len(os.Args) != 3 && len(os.Args) != 5 && len(os.Args) != 7 {
		printError()
	}
	for n := 1; n < len(os.Args); n++ {
		if 0 == strings.Compare(os.Args[n], "-port") {
			port, _ = strconv.Atoi(os.Args[n+1])
			fmt.Printf("\nGiven Port %d: ", port)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-t") {
			topic = os.Args[n+1]
			fmt.Printf("Topic is : %s", topic)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-secured") {
			isSecured, _ = strconv.Atoi(os.Args[n+1])
			fmt.Printf("\nSecured %d: ", isSecured)
			n = n + 1
		} else {
			printError()
		}
	}

	//Handler for ctrl+c
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignal
		fmt.Println(sig)
		if nil != publisher {
			publisher.Stop()
		}
		if nil != instance {
			instance.Terminate()
		}
		os.Exit(0)
	}()

	//get singleton instance
	instance = ezmq.GetInstance()

	//Initilize the EZMQ SDK
	result = instance.Initialize()
	fmt.Printf("\n[Initialize] Error code is: %d", result)
	if result != ezmq.EZMQ_OK {
		fmt.Printf("Error while initializing\n")
		os.Exit(-1)
	}
	publisher = ezmq.GetEZMQPublisher(port, startCB, stopCB, errorCB)

	//set key
	if 1 == isSecured {
		result = publisher.SetServerPrivateKey([]byte(ServerSecretKey))
		if result != ezmq.EZMQ_OK {
			fmt.Printf("\nError while setting key\n")
			os.Exit(-1)
		}
	}

	//start publisher
	result = publisher.Start()
	if result != ezmq.EZMQ_OK {
		fmt.Printf("\nError while starting publisher\n")
		os.Exit(-1)
	}
	fmt.Printf("\n[Start] Error code is: %d", result)

	var event ezmq.Event = getEvent()
	//var byteData ezmq.EZMQByteData = getByteDataEvent()

	// This delay is added to prevent ZeroMQ first packet drop during
	// initial connection of publisher and subscriber.
	time.Sleep(1000 * time.Millisecond)

	fmt.Printf("\n--------- Will Publish 15 events at interval of 1 seconds ---------\n")
	for i := 0; i < 15; i++ {
		if topic == "" {
			result = publisher.Publish(event)
			// result = publisher.Publish(byteData)
		} else {
			result = publisher.PublishOnTopic(topic, event)
			//result = publisher.PublishOnTopic(topic, byteData)
		}
		if result != ezmq.EZMQ_OK {
			fmt.Printf("\nError while publishing")
		}
		fmt.Printf("\nPublished event result: %d\n", result)
		time.Sleep(1000 * time.Millisecond)
	}

	//stop publisher
	result = publisher.Stop()
	if result != ezmq.EZMQ_OK {
		fmt.Printf("Error while Stopping publisher")
	}
	fmt.Printf("\n[Stop] Error code is: %d\n", result)
}
