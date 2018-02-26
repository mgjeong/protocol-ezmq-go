// +build !debug

package ezmq

import (
	"go.uber.org/zap"

	"fmt"
)

var logger *zap.Logger

func InitLogger() {
	var err error
	logger = zap.NewNop()
	if nil != err {
		_ = fmt.Errorf("\nlogger creation failed")
	}
}
