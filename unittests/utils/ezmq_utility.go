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

package utils

import (
	ezmq "go/ezmq"
)

var Ip string = "localhost"
var Port int = 5562
var Topic string = "topic"
var Result ezmq.EZMQErrorCode
var StatusCode ezmq.EZMQStatusCode
var ApiInstance *ezmq.EZMQAPI

func GetEvent() ezmq.Event {
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

func GetWrongEvent() ezmq.Event {
	var event ezmq.Event

	var id string = "id1"
	event.Id = &id
	var device string = "device1"
	event.Device = &device

	//form the reading
	var reading1 *ezmq.Reading = &ezmq.Reading{}
	var rId string = "id1"
	reading1.Id = &rId
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

func GetByteDataEvent() ezmq.EZMQByteData {
	var bytes ezmq.EZMQByteData
	byteArray := [5]byte{0x40, 0x05, 0x10, 0x11, 0x12}
	bytes.ByteData = byteArray[:]
	return bytes
}
