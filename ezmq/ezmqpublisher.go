// +build unsecure

package ezmq

import (
	proto "github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"go.uber.org/zap"

	List "container/list"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Address Prefix for Publisher.
const PUB_TCP_PREFIX = "tcp://*:"

// Regex Pattern for Topic validation.
const TOPIC_PATTERN = "^[a-zA-Z0-9-_./]+$"

// Callback to get error code for start of EZMQ publisher.
// [As of now, Not being used]
type EZMQStartCB func(code EZMQErrorCode)

// Callback to get error code for stop of EZMQ publisher.
// [As of now, Not being used]
type EZMQStopCB func(code EZMQErrorCode)

// Error Callback for start/stop of EZMQ publisher.
// [As of now, Not being used]
type EZMQErrorCB func(code EZMQErrorCode)

//Structure represents EZMQPublisher.
type EZMQPublisher struct {
	port          int
	startCallback EZMQStartCB
	stopCallback  EZMQStopCB
	errorCallback EZMQErrorCB

	publisher *zmq.Socket
	context   *zmq.Context
	mutex     *sync.Mutex
}

// Constructs EZMQPublisher.
func GetEZMQPublisher(port int, startCallback EZMQStartCB, stopCallback EZMQStopCB, errorCallback EZMQErrorCB) *EZMQPublisher {
	var instance *EZMQPublisher
	instance = &EZMQPublisher{}
	instance.port = port
	instance.startCallback = startCallback
	instance.stopCallback = stopCallback
	instance.errorCallback = errorCallback
	instance.context = GetInstance().GetContext()

	if nil == instance.context {
		logger.Error("Context is null")
		return nil
	}
	instance.publisher = nil
	instance.mutex = &sync.Mutex{}
	InitLogger()
	return instance
}

// Starts PUB instance.
func (pubInstance *EZMQPublisher) Start() EZMQErrorCode {
	if nil == pubInstance.context {
		logger.Error("Context is null")
		return EZMQ_ERROR
	}

	pubInstance.mutex.Lock()
	defer pubInstance.mutex.Unlock()
	if nil == pubInstance.publisher {
		var err error
		pubInstance.publisher, err = instance.context.NewSocket(zmq.PUB)
		if nil != err {
			logger.Error("Publisher Socket creation failed")
		}
		var address string = getPubSocketAddress(pubInstance.port)
		err = pubInstance.publisher.Bind(address)
		if nil != err {
			logger.Error("Error while starting publisher")
			pubInstance.publisher = nil
			return EZMQ_ERROR
		}
		logger.Debug("Publisher started", zap.String("address", address))
	}
	return EZMQ_OK
}

func getHeader(content EZMQContentType) []byte {

	var ezmqHeader byte = 0x00
	var version byte = 1
	var contentType byte = (byte)(content)

	version = (byte)(version << 2)
	ezmqHeader = (byte)(ezmqHeader | version)
	contentType = (byte)(contentType << 5)
	ezmqHeader = (byte)(ezmqHeader | contentType)

	header := [1]byte{ezmqHeader}
	return header[:]
}

func (pubInstance *EZMQPublisher) publishInternal(topic string, ezmqMsg EZMQMessage) EZMQErrorCode {
	if nil == ezmqMsg {
		return EZMQ_ERROR
	}
	// form the EZMQ header
	contentType := ezmqMsg.GetContentType()
	header := getHeader(contentType)

	// form the EZMQ data
	var byteEvent []byte = nil
	var err error
	if contentType == EZMQ_CONTENT_TYPE_PROTOBUF {
		event := ezmqMsg.(Event)
		byteEvent, err = proto.Marshal(&event)
		if nil != err {
			logger.Error("Error occured while marshalling proto event")
			return EZMQ_ERROR
		}
	} else if contentType == EZMQ_CONTENT_TYPE_BYTEDATA {
		ezmqByteData := ezmqMsg.(EZMQByteData)
		byteEvent = ezmqByteData.GetByteData()
	} else {
		logger.Error("Not a supported messsage type")
		return EZMQ_INVALID_CONTENT_TYPE
	}

	if nil == byteEvent {
		logger.Error("nil byte event")
		return EZMQ_ERROR
	}

	pubInstance.mutex.Lock()
	defer pubInstance.mutex.Unlock()
	if nil == pubInstance.publisher {
		logger.Error("Publisher is nil")
		return EZMQ_ERROR
	}
	// send topic [if any]
	if topic != "" {
		result, err := pubInstance.publisher.Send(topic, zmq.SNDMORE)
		if nil != err {
			logger.Error("Error while sending topic", zap.Int("Sent bytes", result))
			return EZMQ_ERROR
		}
	}

	// send header
	result, err := pubInstance.publisher.SendBytes(header, zmq.SNDMORE)
	if nil != err {
		logger.Error("Error while sending header", zap.Int("Sent bytes", result))
		return EZMQ_ERROR
	}

	// send data
	result, err = pubInstance.publisher.SendBytes(byteEvent, 0)
	if nil != err {
		logger.Error("Error while publishing data", zap.Int("Sent bytes", result))
		return EZMQ_ERROR
	}
	logger.Debug("Published data")
	return EZMQ_OK
}

// Publish events on the socket for subscribers.
func (pubInstance *EZMQPublisher) Publish(ezmqMsg EZMQMessage) EZMQErrorCode {
	if nil == pubInstance.publisher {
		logger.Error("Publisher is null")
		return EZMQ_ERROR
	}
	return pubInstance.publishInternal("", ezmqMsg)
}

// Publish events on a specific topic on socket for subscribers.
//
// Note:
// (1) Topic name should be as path format. For example:home/livingroom/
//
// (2) Topic name can have letters [a-z, A-z], numerics [0-9] and special characters _ - / and .
func (pubInstance *EZMQPublisher) PublishOnTopic(topic string, ezmqMsg EZMQMessage) EZMQErrorCode {
	if nil == pubInstance.publisher {
		return EZMQ_ERROR
	}

	//validate the topic
	validTopic := sanitizeTopic(topic)
	if validTopic == "" {
		return EZMQ_INVALID_TOPIC
	}
	return pubInstance.publishInternal(validTopic, ezmqMsg)
}

// Publish an events on list of topics on socket for subscribers. On any of
// the topic in list, if it failed to publish event it will return
// EZMQ_ERROR/EZMQ_INVALID_TOPIC.
//
// Note:
// (1) Topic name should be as path format. For example:home/livingroom/
//
// (2) Topic name can have letters [a-z, A-z], numerics [0-9] and special characters _ - / and .
func (pubInstance *EZMQPublisher) PublishOnTopicList(topicList List.List, ezmqMsg EZMQMessage) EZMQErrorCode {
	if topicList.Len() == 0 {
		return EZMQ_INVALID_TOPIC
	}
	for topic := topicList.Front(); topic != nil; topic = topic.Next() {
		result := pubInstance.PublishOnTopic(topic.Value.(string), ezmqMsg)
		if result != EZMQ_OK {
			return EZMQ_ERROR
		}
	}
	return EZMQ_OK
}

// Stops PUB instance.
func (pubInstance *EZMQPublisher) Stop() EZMQErrorCode {
	pubInstance.mutex.Lock()
	defer pubInstance.mutex.Unlock()

	if nil == pubInstance.publisher {
		logger.Error("Publisher is null")
		return EZMQ_ERROR
	}
	// Sync close
	result := pubInstance.syncClose()
	if result == EZMQ_OK {
		pubInstance.publisher = nil
		logger.Debug("Publisher Stopped")
	}
	return result
}

// Get publisher port.
func (pubInstance *EZMQPublisher) GetPort() int {
	return pubInstance.port
}

func getPubSocketAddress(port int) string {
	return string(PUB_TCP_PREFIX) + strconv.Itoa(port)
}

func sanitizeTopic(topic string) string {
	if topic == "" {
		return topic
	}

	result, _ := regexp.MatchString(TOPIC_PATTERN, topic)
	if false == result {
		return ""
	}

	index := strings.LastIndex(topic, "/")
	if index == -1 {
		return topic + "/"
	}

	if index != (len(topic) - 1) {
		return topic + "/"
	}
	return topic
}

func getMonitorAddress() string {
	var address string = "inproc://monitor-"
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	return address + strconv.Itoa(random.Intn(10000))
}

func (pubInstance *EZMQPublisher) syncClose() EZMQErrorCode {
	var errMonitor error
	var address string = getMonitorAddress()
	errMonitor = pubInstance.publisher.Monitor(address, zmq.EVENT_CLOSED)
	if errMonitor != nil {
		logger.Info("Error in monitor")
	}
	socket, errMonitor := pubInstance.context.NewSocket(zmq.PAIR)
	if errMonitor == nil {
		errMonitor = socket.Connect(address)
		if errMonitor != nil {
			logger.Info("Pair socket connection failed")
		}
	} else {
		logger.Info("Pair socket creation failed")
	}

	//close the publisher socket
	err := pubInstance.publisher.Close()
	if nil != err {
		return EZMQ_ERROR
	}
	logger.Debug("Closed publisher socket")

	if nil == socket || errMonitor != nil {
		return EZMQ_OK
	}
	// set receive timeout
	socket.SetRcvtimeo(time.Second)
	for {
		event, addr, value, err := socket.RecvEvent(0)
		if err != nil {
			logger.Error("Error while receiving Event")
			socket.Close()
			break
		}
		logger.Debug("Event received", zap.Int("eventType", int(event)), zap.String("address", addr),
			zap.Int("Value", value))
		if event == zmq.EVENT_CLOSED {
			logger.Debug("Received EVENT_CLOSED from socket")
			socket.Close()
			break
		}
	}
	return EZMQ_OK
}
