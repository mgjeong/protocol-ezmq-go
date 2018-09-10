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
)

var pubResult ezmq.EZMQErrorCode
var pubApiInstance *ezmq.EZMQAPI
var publisher *ezmq.EZMQPublisher

func startCB(code ezmq.EZMQErrorCode) { fmt.Printf("startCB") }
func stopCB(code ezmq.EZMQErrorCode)  { fmt.Printf("stopCB") }
func errorCB(code ezmq.EZMQErrorCode) { fmt.Printf("errorCB") }

func TestGetPubInstance(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubApiInstance.Terminate()
}

func TestGePubInstanceNegative(t *testing.T) {
	publisher = nil
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil != publisher {
		t.Errorf("\nPublisher instance is NUvLL")
	}
}

func TestPublisherStart(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}
	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("\nError while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPublish1(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}

	var event ezmq.Event = utils.GetEvent()
	pubResult = publisher.Publish(event)
	if pubResult != 0 {
		t.Errorf("\nError while publishing event\n")
	}

	byteData := utils.GetByteDataEvent()
	pubResult = publisher.Publish(byteData)
	if pubResult != 0 {
		t.Errorf("\nError while publishing event\n")
	}

	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("\nError while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPublish2(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}

	var event ezmq.Event = utils.GetEvent()
	pubResult = publisher.PublishOnTopic(utils.Topic, event)
	if pubResult != 0 {
		t.Errorf("\nError while publishing event on utils.Topic\n")
	}

	byteData := utils.GetByteDataEvent()
	pubResult = publisher.PublishOnTopic(utils.Topic, byteData)
	if pubResult != 0 {
		t.Errorf("\nError while publishing byte data on utils.Topic\n")
	}
	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("\nError while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPublish3(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}

	var event ezmq.Event = utils.GetEvent()
	topicList := List.New()
	pubResult = publisher.PublishOnTopicList(*topicList, event)
	if pubResult != 2 {
		t.Errorf("\nWrong error code\n")
	}
	e1 := topicList.PushFront("topic1")
	_ = e1
	e2 := topicList.PushFront("topic2")
	_ = e2
	pubResult = publisher.PublishOnTopicList(*topicList, event)
	if pubResult != 0 {
		t.Errorf("\nError while publishing event on utils.Topic list\n")
	}
	byteData := utils.GetByteDataEvent()
	pubResult = publisher.PublishOnTopicList(*topicList, byteData)
	if pubResult != 0 {
		t.Errorf("\nError while publishing byte data on utils.Topic list\n")
	}

	e3 := topicList.PushFront("")
	_ = e3
	pubResult = publisher.PublishOnTopicList(*topicList, event)
	if pubResult == 0 {
		t.Errorf("\nPublished on wrong utils.Topic list\n")
	}
	pubResult = publisher.PublishOnTopicList(*topicList, byteData)
	if pubResult == 0 {
		t.Errorf("\nError while publishing byte data on utils.Topic list\n")
	}

	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("Error while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPublishSecured(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}

	serverSecretKey := "[:X%Q3UfY+kv2A^.wv:(qy2E=bk0L][cm=mS3Hcx"
	pubResult = publisher.SetServerPrivateKey([]byte(serverSecretKey))
	if pubResult != ezmq.EZMQ_OK {
		t.Errorf("\nError while setting server private key\n")
	}

	//negative case
	pubResult = publisher.SetServerPrivateKey([]byte(""))
	if pubResult != ezmq.EZMQ_ERROR {
		t.Errorf("\nWrong error code\n")
	}

	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}

	var event ezmq.Event = utils.GetEvent()
	pubResult = publisher.Publish(event)
	if pubResult != 0 {
		t.Errorf("\nError while publishing event\n")
	}

	byteData := utils.GetByteDataEvent()
	pubResult = publisher.Publish(byteData)
	if pubResult != 0 {
		t.Errorf("\nError while publishing event\n")
	}

	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("\nError while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPublicTopic(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}

	var event ezmq.Event = utils.GetEvent()
	byteData := utils.GetByteDataEvent()

	var testingTopic string = ""

	// Empty utils.Topic test
	if 2 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}
	if 2 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}

	// Alphabet test
	testingTopic = "utils.Topic"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Numeric test
	testingTopic = "123"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Alpha-Numeric test
	testingTopic = "1a2b3"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Alphabet forward slash test
	testingTopic = "utils.Topic/"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Alphabet-Numeric, forward slash test
	testingTopic = "utils.Topic/13/4jtjos/"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Alphabet-Numeric, forward slash test
	testingTopic = "123a/1this3/4jtjos"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Alphabet, backslash test
	testingTopic = "utils.Topic\";"
	if 2 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}
	if 2 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}

	// Alphabet-Numeric, forward slash and space test
	testingTopic = "utils.Topic/13/4jtjos/ "
	if 2 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}
	if 2 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}

	// Special character test
	testingTopic = "*123a"
	if 2 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}
	if 2 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}

	// Sentence test
	testingTopic = "This is a utils.Topic"
	if 2 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}
	if 2 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished on invalid utils.Topic\n")
	}

	// Topic contain forward slash at last
	testingTopic = "utils.Topic/122/livingroom/"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Topic contain -
	testingTopic = "utils.Topic/122/livingroom/-"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Topic contain _
	testingTopic = "utils.Topic/122/livingroom_"
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	// Topic contain .
	testingTopic = "utils.Topic/122.livingroom."
	if 0 != (publisher.PublishOnTopic(testingTopic, event)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}
	if 0 != (publisher.PublishOnTopic(testingTopic, byteData)) {
		t.Errorf("\nPublished failed for valid utils.Topic\n")
	}

	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("\nError while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPublishNegative(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	pubResult = publisher.Start()
	if pubResult != 0 {
		t.Errorf("\nError while starting publisher\n")
	}

	var event ezmq.Event = utils.GetWrongEvent()
	pubResult = publisher.Publish(event)
	if pubResult == 0 {
		t.Errorf("\nPublished wrong event\n")
	}
	pubResult = publisher.Publish(nil)
	if pubResult == 0 {
		t.Errorf("\nPublished nil event\n")
	}
	pubResult = publisher.PublishOnTopic(utils.Topic, event)
	if pubResult == 0 {
		t.Errorf("\nPublished wrong event\n")
	}
	pubResult = publisher.PublishOnTopic(utils.Topic, nil)
	if pubResult == 0 {
		t.Errorf("\nPublished wrong event\n")
	}

	topicList := List.New()
	e1 := topicList.PushFront("topic1")
	_ = e1
	e2 := topicList.PushFront("")
	_ = e2
	pubResult = publisher.PublishOnTopicList(*topicList, event)
	if pubResult == 0 {
		t.Errorf("\nPublished invalid topiclist\n")
	}
	byteData := utils.GetByteDataEvent()
	pubResult = publisher.PublishOnTopicList(*topicList, byteData)
	if pubResult == 0 {
		t.Errorf("\nPublished invalid topiclist\n")
	}

	pubResult = publisher.Stop()
	if pubResult != 0 {
		t.Errorf("\nError while Stopping publisher")
	}
	pubApiInstance.Terminate()
}

func TestPubStartStop(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}

	for i := 0; i < 15; i++ {
		pubResult = publisher.Start()
		if pubResult != 0 {
			t.Errorf("\nError while starting publisher\n")
		}
		pubResult = publisher.Stop()
		if pubResult != 0 {
			t.Errorf("\nError while stopping publisher\n")
		}
	}
	pubApiInstance.Terminate()
}

func TestPubGetPort(t *testing.T) {
	pubApiInstance = ezmq.GetInstance()
	pubApiInstance.Initialize()
	publisher = ezmq.GetEZMQPublisher(utils.Port, startCB, stopCB, errorCB)
	if nil == publisher {
		t.Errorf("\nPublisher instance is NULL")
	}
	port := publisher.GetPort()
	if port != utils.Port {
		t.Errorf("\nAssertion failed")
	}
}
