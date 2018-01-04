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
	"strconv"
	"strings"
	"syscall"
)

func printEvent(event ezmq.Event) {
	fmt.Printf("\n--------------------------------------")
	fmt.Printf("\nDevice: %s", event.GetDevice())
	fmt.Printf("\nReadings:")

	var readings []*ezmq.Reading = event.GetReading()
	for i := 0; i < len(readings); i++ {
		fmt.Printf("\nKey: %s", readings[i].GetName())
		fmt.Printf("\nValue: %s", readings[i].GetValue())
	}
	fmt.Printf("\n--------------------------------------\n")
}

func printError() {
	fmt.Printf("\nRe-run the application as shown in below examples: \n")
	fmt.Printf("\n  (1) For subscribing without topic: ")
	fmt.Printf("\n     ./subscriber -ip 107.108.81.116 -port 5562")
	fmt.Printf("\n     ./subscriber -ip localhost -port 5562\n")
	fmt.Printf("\n  (2) For subscribing with topic: ")
	fmt.Printf("\n     ./subscriber -ip 107.108.81.116 -port 5562 -t topic1")
	fmt.Printf("\n     ./subscriber -ip localhost -port 5562 -t topic1\n")
	os.Exit(-1)
}

func main() {
	var ip string
	var port int
	var topic string
	var result ezmq.EZMQErrorCode
	var instance *ezmq.EZMQAPI
	var subscriber *ezmq.EZMQSubscriber
	var isSubscribed bool = false

	// get ip and port from command line arguments
	if len(os.Args) != 5 && len(os.Args) != 7 {
		printError()
	}

	for n := 1; n < len(os.Args); n++ {
		if 0 == strings.Compare(os.Args[n], "-ip") {
			ip = os.Args[n+1]
			fmt.Printf("\nGiven Ip: %s", ip)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-port") {
			port, _ = strconv.Atoi(os.Args[n+1])
			fmt.Printf("\nGiven Port %d: ", port)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-t") {
			topic = os.Args[n+1]
			fmt.Printf("Topic is : %s", topic)
			n = n + 1
		} else {
			printError()
		}
	}

	// Handler for ctrl+c
	osSignal := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignal
		fmt.Println(sig)
		if false == isSubscribed {
			os.Exit(-1)
		}
		subscriber.Stop()
		instance.Terminate()
		exit <- true
	}()

	// callbacks
	subCB := func(event ezmq.Event) { printEvent(event) }
	subTopicCB := func(topic string, event ezmq.Event) {
		fmt.Printf("\nTopic: %s", topic)
		printEvent(event)
	}

	// get singleton instance
	instance = ezmq.GetInstance()

	//Initilize the EZMQ SDK
	result = instance.Initialize()
	fmt.Printf("\n[Initialize] Error code is: %d", result)
	if result != ezmq.EZMQ_OK {
		fmt.Printf("Error while initializing\n")
		os.Exit(-1)
	}

	subscriber = ezmq.GetEZMQSubscriber(ip, port, subCB, subTopicCB)

	// start subscriber
	result = subscriber.Start()
	if result != ezmq.EZMQ_OK {
		fmt.Printf("Error while starting subscriber\n")
		os.Exit(-1)
	}
	fmt.Printf("\n[Start] Error code is: %d\n", result)

	if topic == "" {
		result = subscriber.Subscribe()
	} else {
		result = subscriber.SubscribeForTopic(topic)
	}

	if result != ezmq.EZMQ_OK {
		fmt.Printf("Error while Subscribing\n")
		os.Exit(-1)
	}
	isSubscribed = true
	fmt.Printf("\nSuscribed to publisher.. -- Waiting for Events --\n")

	<-exit
	fmt.Println("exiting")
}
