package unittests

import (
	"bytes"
	ezmq "go/ezmq"
	test_utils "go/unittests/utils"

	"testing"
)

func TestGetByteData(t *testing.T) {
	var byteData ezmq.EZMQByteData
	byteArray := [5]byte{0x40, 0x05, 0x10, 0x11, 0x12}
	byteData.ByteData = byteArray[:]
	var array []byte = byteArray[:]
	var getArray []byte = byteData.GetByteData()
	if false == bytes.Equal(array, getArray) {
		t.Errorf("\nAssertion failed")
	}
}

func TestSetByteData(t *testing.T) {
	var byteData ezmq.EZMQByteData
	byteArray := [5]byte{0x40, 0x05, 0x10, 0x11, 0x12}
	byteData.ByteData = byteArray[:]

	byteArray = [5]byte{0x10, 0x15, 0x11, 0x14, 0x16}
	byteData.SetByteData(byteArray[:])

	var array []byte = byteArray[:]
	var getArray []byte = byteData.GetByteData()
	if false == bytes.Equal(array, getArray) {
		t.Errorf("\nAssertion failed")
	}
}

func TestSetNilByteData(t *testing.T) {
	var byteData ezmq.EZMQByteData
	if 0 == byteData.SetByteData(nil) {
		t.Errorf("\nAssertion failed")
	}
}

func TestDataContentType(t *testing.T) {
	byteData := test_utils.GetByteDataEvent()
	if 1 != byteData.GetContentType() {
		t.Errorf("\nAssertion failed")
	}
}
