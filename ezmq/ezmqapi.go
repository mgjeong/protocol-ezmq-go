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

// EZMQ package which provides simplified APIs for Publisher and Subscriber.
package ezmq

import (
	zmq "github.com/pebbe/zmq4"

	"math/rand"
	"time"
)

// Structure represents EZMQAPI
type EZMQAPI struct {
	context *zmq.Context
	status  EZMQStatusCode
}

var instance *EZMQAPI

// Get EZMQAPI instance.
func GetInstance() *EZMQAPI {
	if nil == instance {
		instance = &EZMQAPI{}
		instance.status = EZMQ_Constructed
		InitLogger()
		rand.Seed(time.Now().UnixNano())
	}
	return instance
}

// Initialize required EZMQ components. This API should be called first,
// before using any EZMQ APIs.
func (ezmqInstance *EZMQAPI) Initialize() EZMQErrorCode {
	if nil == ezmqInstance.context {
		var err error
		ezmqInstance.context, err = zmq.NewContext()
		if err != nil {
			logger.Error("EZMQ initialization failed")
			return EZMQ_ERROR
		}
		zmq.SetIoThreads(1)
	}
	logger.Debug("EZMQ initialized")

	ezmqInstance.status = EZMQ_Initialized
	return EZMQ_OK
}

// Perform cleanup of EZMQ components.
func (ezmqInstance *EZMQAPI) Terminate() EZMQErrorCode {
	if ezmqInstance.context != nil {
		err := ezmqInstance.context.Term()
		if nil != err {
			logger.Error("EZMQ termination failed")
			return EZMQ_ERROR
		}
		ezmqInstance.context = nil
		logger.Debug("EZMQ terminated")
	}
	ezmqInstance.status = EZMQ_Terminated
	return EZMQ_OK
}

func (ezmqInstance *EZMQAPI) GetStatus() EZMQStatusCode {
	return ezmqInstance.status
}

func (ezmqInstance *EZMQAPI) GetContext() *zmq.Context {
	return ezmqInstance.context
}
