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
	proto "github.com/golang/protobuf/proto"
	ezmq "go/ezmq"
	utils "go/unittests/utils"

	"testing"
)

func TestSerialization(t *testing.T) {
	var event ezmq.Event = utils.GetEvent()
	var readings []*ezmq.Reading = event.GetReading()

	//marshal [serialize] the protobuf event
	byteEvent, err := proto.Marshal(&event)
	if nil != err {
		t.Errorf("Error occured while marshalling event")
	}

	//Unmarshal [de-serialize] the protobuf event
	var convertedevent ezmq.Event
	er := proto.Unmarshal(byteEvent, &convertedevent)
	if nil != er {
		t.Errorf("Error in unmarshalling event")
	}

	if *(event.Id) != convertedevent.GetId() {
		t.Errorf("Id mismatch")
	}

	if *(event.Created) != convertedevent.GetCreated() {
		t.Errorf("Id mismatch")
	}

	if *(event.Modified) != convertedevent.GetModified() {
		t.Errorf("Modified mismatch")
	}

	if *(event.Origin) != convertedevent.GetOrigin() {
		t.Errorf("Origin mismatch")
	}

	if *(event.Pushed) != convertedevent.GetPushed() {
		t.Errorf("Pushed mismatch")
	}

	if *(event.Device) != convertedevent.GetDevice() {
		t.Errorf("Device mismatch")
	}

	if event.String() != convertedevent.String() {
		t.Errorf("Event string mismatch")
	}

	var covertedReadings []*ezmq.Reading = convertedevent.GetReading()
	if *(readings[0].Id) != covertedReadings[0].GetId() {
		t.Errorf("Reading: Id mismatch")
	}

	if *(readings[0].Created) != covertedReadings[0].GetCreated() {
		t.Errorf("Reading: Id mismatch")
	}

	if *(readings[0].Modified) != covertedReadings[0].GetModified() {
		t.Errorf("Reading: Modified mismatch")
	}

	if *(readings[0].Origin) != covertedReadings[0].GetOrigin() {
		t.Errorf("Reading: Origin mismatch")
	}

	if *(readings[0].Pushed) != covertedReadings[0].GetPushed() {
		t.Errorf("Reading: Pushed mismatch")
	}

	if *(readings[0].Name) != covertedReadings[0].GetName() {
		t.Errorf("Reading: Name mismatch")
	}

	if *(readings[0].Value) != covertedReadings[0].GetValue() {
		t.Errorf("Reading: Value mismatch")
	}

	if *(readings[0].Device) != covertedReadings[0].GetDevice() {
		t.Errorf("Reading: Device mismatch")
	}

	if (readings[0].String()) != covertedReadings[0].String() {
		t.Errorf("Reading: String mismatch")
	}
	//reset the reading
	covertedReadings[0].Reset()
}

func TestEventContentType(t *testing.T) {
	event := utils.GetEvent()
	if 0 != event.GetContentType() {
		t.Errorf("\nAssertion failed")
	}
}

func TestEmptyEvent(t *testing.T) {
	var event ezmq.Event
	var reading1 *ezmq.Reading = &ezmq.Reading{}
	event.Reading = make([]*ezmq.Reading, 1)
	event.Reading[0] = reading1
	stringValue := event.GetDevice()
	stringValue = event.GetId()
	intValue := event.GetCreated()
	intValue = event.GetModified()
	intValue = event.GetOrigin()
	intValue = event.GetPushed()
	var readings []*ezmq.Reading = event.GetReading()
	stringValue = readings[0].GetDevice()
	stringValue = readings[0].GetId()
	stringValue = readings[0].GetValue()
	stringValue = readings[0].GetName()
	intValue = readings[0].GetOrigin()
	intValue = readings[0].GetPushed()
	intValue = readings[0].GetModified()
	intValue = readings[0].GetCreated()
	if stringValue != "" {
		t.Errorf("\nInvalid value")
	}
	if intValue != 0 {
		t.Errorf("\nInvalid value")
	}
}
