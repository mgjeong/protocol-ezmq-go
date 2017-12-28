/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
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
	"syscall"
	"time"
)

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

func main() {
	var port int = 5562
	var result ezmq.EZMQErrorCode
	var publisher *ezmq.EZMQPublisher = nil
	var instance *ezmq.EZMQAPI = nil
	startCB := func(code ezmq.EZMQErrorCode) { fmt.Printf("startCB") }
	stopCB := func(code ezmq.EZMQErrorCode) { fmt.Printf("stopCB") }
	errorCB := func(code ezmq.EZMQErrorCode) { fmt.Printf("errorCB") }

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

	//User choice
	var choice int
	var topic string
	fmt.Printf("\nEnter 1 for General Event testing")
	fmt.Printf("\nEnter 2 for Topic Based delivery\n")
	fmt.Scanf("%d", &choice)

	switch choice {
	case 1:
		publisher = ezmq.GetEZMQPublisher(port, startCB, stopCB, errorCB)
	case 2:
		publisher = ezmq.GetEZMQPublisher(port, startCB, stopCB, errorCB)
		fmt.Printf("\nEnter the topic: ")
		fmt.Scanf("%s", &topic)
		fmt.Printf("Topic is: %s\n", topic)
	default:
		fmt.Printf("\nInvalid choice..[Re-run application]\n")
		os.Exit(-1)
	}

	//start publisher
	result = publisher.Start()
	if result != 0 {
		fmt.Printf("\nError while starting publisher\n")
		os.Exit(-1)
	}
	fmt.Printf("\n[Start] Error code is: %d", result)

	var event ezmq.Event = getEvent()
	fmt.Printf("\n--------- Will Publish 15 events at interval of 1 seconds ---------\n")
	for i := 0; i < 15; i++ {
		if topic == "" {
			result = publisher.Publish(event)
		} else {
			result = publisher.PublishOnTopic(topic, event)
		}
		if result != 0 {
			fmt.Printf("\nError while publishing")
		}
		fmt.Printf("\nPublished event result: %d\n", result)
		time.Sleep(1000 * time.Millisecond)
	}

	//stop publisher
	result = publisher.Stop()
	if result != 0 {
		fmt.Printf("Error while Stopping publisher")
	}
	fmt.Printf("\n[Stop] Error code is: %d\n", result)
}
