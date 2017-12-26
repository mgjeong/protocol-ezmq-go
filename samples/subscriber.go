package main

import (
	ezmq "go/ezmq"

	"fmt"
	"os"
	"os/signal"
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

func main() {
	var ip string = "localhost"
	var port int = 5562
	var result ezmq.EZMQErrorCode
	var instance *ezmq.EZMQAPI
	var subscriber *ezmq.EZMQSubscriber
	var isSubscribed bool = false

	//Handler for ctrl+c
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

	//callbacks
	subCB := func(event ezmq.Event) { printEvent(event) }
	subTopicCB := func(topic string, event ezmq.Event) {
		fmt.Printf("\nTopic: %s", topic)
		printEvent(event)
	}

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
		subscriber = ezmq.GetEZMQSubscriber(ip, port, subCB, subTopicCB)
	case 2:
		subscriber = ezmq.GetEZMQSubscriber(ip, port, subCB, subTopicCB)
		fmt.Printf("\nEnter the topic: ")
		fmt.Scanf("%s", &topic)
		fmt.Printf("Topic is: %s\n", topic)
	default:
		fmt.Printf("\nInvalid choice..[Re-run application]\n")
		os.Exit(-1)
	}

	//start subscriber
	result = subscriber.Start()
	if result != 0 {
		fmt.Printf("Error while starting subscriber\n")
		os.Exit(-1)
	}
	fmt.Printf("\n[Start] Error code is: %d\n", result)

	if topic == "" {
		result = subscriber.Subscribe()
	} else {
		result = subscriber.SubscribeForTopic(topic)
	}

	if result != 0 {
		fmt.Printf("Error while Subscribing\n")
		os.Exit(-1)
	}
	isSubscribed = true
	fmt.Printf("\nSuscribed to publisher.. -- Waiting for Events --\n")

	<-exit
	fmt.Println("exiting")
}
