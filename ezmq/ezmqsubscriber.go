package ezmq

import (
	proto "github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"go.uber.org/zap"

	List "container/list"
	"strconv"
)

// Address prefix to bind subscriber.
const SUB_TCP_PREFIX = "tcp://"

// Callback to get all the subscribed events.
type EZMQSubCB func(event Event)

// Callback to get all the subscribed events for a specific topic.
type EZMQSubTopicCB func(topic string, event Event)

// Structure represents EZMQSubscriber.
type EZMQSubscriber struct {
	ip               string
	port             int
	subCallback      EZMQSubCB
	subTopicCallback EZMQSubTopicCB

	subscriber *zmq.Socket
	context    *zmq.Context

	isReceiverStarted bool
}

// Contructs EZMQSubscriber.
func GetEZMQSubscriber(ip string, port int, subCallback EZMQSubCB, subTopicCallback EZMQSubTopicCB) *EZMQSubscriber {
	var instance *EZMQSubscriber
	instance = &EZMQSubscriber{}
	instance.ip = ip
	instance.port = port
	instance.subCallback = subCallback
	instance.subTopicCallback = subTopicCallback
	instance.context = GetInstance().getContext()
	InitLogger()
	if nil == instance.context {
		logger.Error("Context is null")
	}
	instance.subscriber = nil
	instance.isReceiverStarted = false
	return instance
}

func receive(subInstance *EZMQSubscriber) {
	var data []byte
	var event Event
	var err error
	var more bool
	var topic string

	for subInstance.isReceiverStarted {
		if nil == subInstance.subscriber {
			logger.Error("subscriber or poller is null")
			break
		}
		data, err = subInstance.subscriber.RecvBytes(0)
		if err != nil {
			break
		}
		more, err = subInstance.subscriber.GetRcvmore()
		if err != nil {
			break
		}
		if more {
			topic = string(data[:])
			data, err = subInstance.subscriber.RecvBytes(0)
		}

		//change byte array to Event
		err := proto.Unmarshal(data, &event)
		if nil != err {
			logger.Error("Error in unmarshalling data")
		}

		if more {
			subInstance.subTopicCallback(topic, event)
		} else {
			subInstance.subCallback(event)
		}
	}
	logger.Debug("Received the shut down request")
}

// Starts SUB instance.
func (subInstance *EZMQSubscriber) Start() EZMQErrorCode {
	if nil == subInstance.context {
		logger.Error("Context is null")
		return EZMQ_ERROR
	}

	if nil == subInstance.subscriber {
		var err error
		subInstance.subscriber, err = zmq.NewSocket(zmq.SUB)
		if nil != err {
			logger.Error("Subscriber Socket creation failed")
			return EZMQ_ERROR
		}
		var address string = getSubSocketAddress(subInstance.ip, subInstance.port)
		err = subInstance.subscriber.Connect(address)
		if nil != err {
			logger.Error("Subscriber Socket connect failed")
			return EZMQ_ERROR
		}
		logger.Debug("Starting subscriber", zap.String("Address", address))
	}

	//call a go routine [new thread] for receiver
	if false == subInstance.isReceiverStarted {
		subInstance.isReceiverStarted = true
		go receive(subInstance)
	}
	return EZMQ_OK
}

func (subInstance *EZMQSubscriber) subscribeInternal(topic string) EZMQErrorCode {
	if nil != subInstance.subscriber {
		err := subInstance.subscriber.SetSubscribe(topic)
		if nil != err {
			logger.Error("subscribeInternal error occured")
			return EZMQ_ERROR
		}
	} else {
		logger.Error("subscriber is null")
		return EZMQ_ERROR
	}
	logger.Debug("subscribed for events")
	return EZMQ_OK
}

// Subscribe for event/messages.
func (subInstance *EZMQSubscriber) Subscribe() EZMQErrorCode {
	return subInstance.subscribeInternal("")
}

// Subscribe for event/messages on a particular topic.
func (subInstance *EZMQSubscriber) SubscribeForTopic(topic string) EZMQErrorCode {
	//validate the topic
	validTopic := sanitizeTopic(topic)
	if validTopic == "" {
		return EZMQ_INVALID_TOPIC
	}
	logger.Debug("subscribing for events", zap.String("Topic", validTopic))
	return subInstance.subscribeInternal(validTopic)
}

// Subscribe for event/messages on given list of topics. On any of the topic
// in list, if it failed to subscribe events it will return
// EZMQ_ERROR/EZMQ_INVALID_TOPIC.
//
// Note:
// (1) Topic name should be as path format. For example:home/livingroom/
// (2) Topic name can have letters [a-z, A-z], numerics [0-9] and special characters _ - / and .
func (subInstance *EZMQSubscriber) SubscribeForTopicList(topicList List.List) EZMQErrorCode {
	if topicList.Len() == 0 {
		return EZMQ_INVALID_TOPIC
	}
	for topic := topicList.Front(); topic != nil; topic = topic.Next() {
		result := subInstance.SubscribeForTopic(topic.Value.(string))
		if result != EZMQ_OK {
			return result
		}
	}
	return EZMQ_OK
}

func (subInstance *EZMQSubscriber) unSubscribeInternal(topic string) EZMQErrorCode {
	if nil != subInstance.subscriber {
		err := subInstance.subscriber.SetUnsubscribe(topic)
		if nil != err {
			logger.Error("subscriber is null")
			return EZMQ_ERROR
		}
	} else {
		return EZMQ_ERROR
	}
	return EZMQ_OK
}

// Un-subscribe all the events from publisher.
func (subInstance *EZMQSubscriber) UnSubscribe() EZMQErrorCode {
	return subInstance.unSubscribeInternal("")
}

// Un-subscribe specific topic events.
//
// Note:
// (1) Topic name should be as path format. For example:home/livingroom/
// (2) Topic name can have letters [a-z, A-z], numerics [0-9] and special characters _ - / and .
func (subInstance *EZMQSubscriber) UnSubscribeForTopic(topic string) EZMQErrorCode {
	//validate the topic
	validTopic := sanitizeTopic(topic)
	if validTopic == "" {
		return EZMQ_INVALID_TOPIC
	}
	logger.Debug("Unsubscribe for events", zap.String("Topic", validTopic))
	return subInstance.unSubscribeInternal(validTopic)
}

// Un-subscribe event/messages on given list of topics. On any of the topic
// in list, if it failed to unsubscribe events it will return
// EZMQ_ERROR/EZMQ_INVALID_TOPIC.
//
// Note:
// (1) Topic name should be as path format. For example:home/livingroom/ .
// (2) Topic name can have letters [a-z, A-z], numerics [0-9] and special characters _ - / and .
func (subInstance *EZMQSubscriber) UnSubscribeForTopicList(topicList List.List) EZMQErrorCode {
	if topicList.Len() == 0 {
		return EZMQ_INVALID_TOPIC
	}
	for topic := topicList.Front(); topic != nil; topic = topic.Next() {
		result := subInstance.UnSubscribeForTopic(topic.Value.(string))
		if result != EZMQ_OK {
			return result
		}
	}
	return EZMQ_OK
}

// Stops SUB instance.
func (subInstance *EZMQSubscriber) Stop() EZMQErrorCode {
	if nil != subInstance.subscriber {
		err := subInstance.subscriber.Close()
		if nil != err {
			logger.Error("Error while stopping subscriber")
			return EZMQ_ERROR
		}
	}
	subInstance.subscriber = nil
	subInstance.isReceiverStarted = false
	logger.Debug("Subscriber stopped")
	return EZMQ_OK
}

// Get Ip of publisher to which subscribed.
func (subInstance *EZMQSubscriber) GetIP() string {
	return subInstance.ip
}

// Get Port of publisher to which subscribed.
func (subInstance *EZMQSubscriber) GetPort() int {
	return subInstance.port
}

func getSubSocketAddress(ip string, port int) string {
	return string(SUB_TCP_PREFIX) + ip + ":" + strconv.Itoa(port)
}
