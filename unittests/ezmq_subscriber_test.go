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

package unittests

import (
	"go/ezmq"
	"go/unittests/utils"

	List "container/list"
	"fmt"
	"testing"
	"time"
)

var subResult ezmq.EZMQErrorCode
var subApiInstance *ezmq.EZMQAPI
var subscriber *ezmq.EZMQSubscriber

var published = make(chan bool, 1)
var eventCount = 0

func subCB(ezmqMsg ezmq.EZMQMessage) {
	fmt.Printf("\nsubCB")
	eventCount++
}
func subTopicCB(topic string, ezmqMsg ezmq.EZMQMessage) {
	fmt.Printf("\nsubTopicCB")
	eventCount++
}

func TestGetSubInstance(t *testing.T) {
	subApiInstance = ezmq.GetInstance()
	subApiInstance.Initialize()
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subApiInstance.Terminate()
}

func TestGeSubInstanceNegative(t *testing.T) {
	subscriber = nil
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil != subscriber {
		t.Errorf("\nSubscriber instace is NULL")
	}
	subApiInstance = ezmq.GetInstance()
	subApiInstance.Initialize()
}

func TestSubscriberStart(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}
	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscribe1(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.Subscribe()
	if subResult != 0 {
		t.Errorf("\nError while subscribing\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscribe2(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.SubscribeForTopic(utils.Topic)
	if subResult != 0 {
		t.Errorf("\nError while subscribing on Utils.Topic\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func publish() {
	var pubApiInstance *ezmq.EZMQAPI
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	publisher.Start()
	var event ezmq.Event = utils.GetEvent()
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		publisher.PublishOnTopic(utils.Topic, event)
	}
	publisher.Stop()
	published <- true
}

func TestSubscribe3(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.SubscribeForTopic(utils.Topic)
	if subResult != 0 {
		t.Errorf("\nError while subscribing on Utils.Topic\n")
	}

	go publish()
	<-published // wait for publisher to stop

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscribe4(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	topicList := List.New()
	subResult = subscriber.SubscribeForTopicList(*topicList)
	if subResult != 2 {
		t.Errorf("\nWrong error code\n")
	}
	e1 := topicList.PushFront("topic1")
	_ = e1
	e2 := topicList.PushFront("topic2")
	_ = e2
	subResult = subscriber.SubscribeForTopicList(*topicList)
	if subResult != 0 {
		t.Errorf("\nError while subscribing on Utils.Topic list\n")
	}

	e3 := topicList.PushFront("")
	_ = e3
	subResult = subscriber.SubscribeForTopicList(*topicList)
	if subResult == 0 {
		t.Errorf("\nSubscribd on invalid Utils.Topic list\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscribeTopic(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	var testingTopic string = ""

	// Empty Utils.Topic test
	if 2 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscribed on invalid Utils.Topic\n")
	}

	// Alphabet test
	testingTopic = "Utils.Topic"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Numeric test
	testingTopic = "123"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Alpha-Numeric test
	testingTopic = "1a2b3"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Alphabet forward slash test
	testingTopic = "Utils.Topic/"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Alphabet-Numeric, forward slash test
	testingTopic = "Utils.Topic/13/4jtjos/"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Alphabet-Numeric, forward slash test
	testingTopic = "123a/1this3/4jtjos"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Alphabet, backslash test
	testingTopic = "Utils.Topic\";"
	if 2 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscribed on invalid Utils.Topic\n")
	}

	// Alphabet-Numeric, forward slash and space test
	testingTopic = "Utils.Topic/13/4jtjos/ "
	if 2 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscribed on invalid Utils.Topic\n")
	}

	// Special character test
	testingTopic = "*123a"
	if 2 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscribed on invalid Utils.Topic\n")
	}

	// Sentence test
	testingTopic = "This is a Utils.Topic"
	if 2 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscribed on invalid Utils.Topic\n")
	}

	// Topic contain forward slash at last
	testingTopic = "Utils.Topic/122/livingroom/"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Topic contain -
	testingTopic = "Utils.Topic/122/livingroom/-"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Topic contain _
	testingTopic = "Utils.Topic/122/livingroom_"
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	// Topic contain .
	testingTopic = "Utils.Topic/122.livingroom."
	if 0 != (subscriber.SubscribeForTopic(testingTopic)) {
		t.Errorf("\nSubscription failed for valid Utils.Topic\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscribeSecured(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}

	clientPrivateKey := "ZB1@RS6Kv^zucova$kH(!o>tZCQ.<!Q)6-0aWFmW"
	clientPublicKey := "-QW?Ved(f:<::3d5tJ$[4Er&]6#9yr=vha/caBc("
	subResult = subscriber.SetClientKeys([]byte(clientPrivateKey), []byte(clientPublicKey))
	if subResult != ezmq.EZMQ_OK {
		t.Errorf("\nError while setting client keys\n")
	}
	serverPublicKey := "tXJx&1^QE2g7WCXbF.$$TVP.wCtxwNhR8?iLi&S<"
	subResult = subscriber.SetServerPublicKey([]byte(serverPublicKey))
	if subResult != ezmq.EZMQ_OK {
		t.Errorf("\nError while setting server key\n")
	}

	//Negative case
	subResult = subscriber.SetClientKeys([]byte(""), []byte(""))
	if subResult != ezmq.EZMQ_ERROR {
		t.Errorf("\nError while setting client keys\n")
	}
	subResult = subscriber.SetServerPublicKey([]byte(""))
	if subResult != ezmq.EZMQ_ERROR {
		t.Errorf("\nError while setting server key\n")
	}

	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.Subscribe()
	if subResult != 0 {
		t.Errorf("\nError while subscribing\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscribeNegative(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, -1, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult == 0 {
		t.Errorf("\nStarted subscriber on invalid utils.Port\n")
	}

	subResult = subscriber.Subscribe()
	if subResult != 0 {
		t.Errorf("\nSubscribed on invalid subscriber\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestUnSubscribe(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.Subscribe()
	if subResult != 0 {
		t.Errorf("\nError while subscribing\n")
	}
	subResult = subscriber.UnSubscribe()
	if subResult != 0 {
		t.Errorf("\nError while unsubscribing\n")
	}

	subResult = subscriber.SubscribeForTopic(utils.Topic)
	if subResult != 0 {
		t.Errorf("\nError while subscribing for Utils.Topic\n")
	}
	subResult = subscriber.UnSubscribeForTopic(utils.Topic)
	if subResult != 0 {
		t.Errorf("\nError while unsubscribing for Utils.Topic\n")
	}

	topicList := List.New()
	e1 := topicList.PushFront("topic1")
	_ = e1
	e2 := topicList.PushFront("topic2")
	_ = e2

	subResult = subscriber.SubscribeForTopicList(*topicList)
	if subResult != 0 {
		t.Errorf("\nError while subscribing for Utils.Topic list\n")
	}
	subResult = subscriber.UnSubscribeForTopicList(*topicList)
	if subResult != 0 {
		t.Errorf("\nError while unsubscribing for Utils.Topic list\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestUnSubscribeNegative(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.SubscribeForTopic("")
	if subResult == 0 {
		t.Errorf("\nSubscribed for invalid Utils.Topic\n")
	}
	subResult = subscriber.UnSubscribeForTopic("")
	if subResult == 0 {
		t.Errorf("\nUnSubscribed for invalid Utils.Topic\n")
	}

	topicList := List.New()
	e1 := topicList.PushFront("topic1")
	_ = e1
	e2 := topicList.PushFront("!$topic2")
	_ = e2

	subResult = subscriber.SubscribeForTopicList(*topicList)
	if subResult == 0 {
		t.Errorf("\nSubscribed for invalid Utils.Topic list\n")
	}
	subResult = subscriber.UnSubscribeForTopicList(*topicList)
	if subResult == 0 {
		t.Errorf("\nUnSubscribed for invalid Utils.Topic list\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubscriberIPPort(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}
	subResult = subscriber.SubscribeWithIPPort("192.168.1.1", 5562, utils.Topic)
	if subResult != 0 {
		t.Errorf("\nError while subscribing\n")
	}
	subResult = subscriber.SubscribeWithIPPort("192.168.1.1", -1, "")
	if subResult != 1 {
		t.Errorf("\nReturned wrong value\n")
	}
	subResult = subscriber.SubscribeWithIPPort("192.168.1.1", 5562, "")
	if subResult != 2 {
		t.Errorf("\nReturned wrong value\n")
	}
}

func TestSubscribeIPPortSecured(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}

	clientPrivateKey := "ZB1@RS6Kv^zucova$kH(!o>tZCQ.<!Q)6-0aWFmW"
	clientPublicKey := "-QW?Ved(f:<::3d5tJ$[4Er&]6#9yr=vha/caBc("
	subResult = subscriber.SetClientKeys([]byte(clientPrivateKey), []byte(clientPublicKey))
	if subResult != ezmq.EZMQ_OK {
		t.Errorf("\nError while setting client keys\n")
	}
	serverPublicKey := "tXJx&1^QE2g7WCXbF.$$TVP.wCtxwNhR8?iLi&S<"
	subResult = subscriber.SetServerPublicKey([]byte(serverPublicKey))
	if subResult != ezmq.EZMQ_OK {
		t.Errorf("\nError while setting server key\n")
	}

	subResult = subscriber.Start()
	if subResult != 0 {
		t.Errorf("\nError while starting subscriber\n")
	}

	subResult = subscriber.SubscribeWithIPPort("192.168.1.1", 5562, utils.Topic)
	if subResult != 0 {
		t.Errorf("\nError while subscribing\n")
	}

	subResult = subscriber.Stop()
	if subResult != 0 {
		t.Errorf("\nError while Stopping subscriber")
	}
}

func TestSubStartStop(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}

	for i := 0; i < 15; i++ {
		subResult = subscriber.Start()
		if subResult != 0 {
			t.Errorf("\nError while starting subscriber\n")
		}
		subResult = subscriber.Stop()
		if subResult != 0 {
			t.Errorf("\nError while Stopping subscriber")
		}
	}
}

func TestSubGetIp(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	ip := subscriber.GetIP()
	if ip != utils.Ip {
		t.Errorf("\nAssertion failed")
	}
}

func TestSubGetPort(t *testing.T) {
	subscriber = ezmq.GetEZMQSubscriber(utils.Ip, utils.Port, subCB, subTopicCB)
	if nil == subscriber {
		t.Errorf("\nSubscriber instance is NULL")
	}
	port := subscriber.GetPort()
	if port != utils.Port {
		t.Errorf("\nAssertion failed")
	}
}
