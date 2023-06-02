package senseauth

import (
	"fmt"
	"log"
)

var debugLogger *log.Logger

func debugging() bool {
	return debugLogger != nil
}

func debugf(format string, args ...interface{}) {
	if debugging() {
		debugLogger.Output(2, fmt.Sprintf(format, args...))
	}
}

func debug(args ...interface{}) {
	if debugging() {
		debugLogger.Output(2, fmt.Sprint(args...))
	}
}

func SetDebug(l *log.Logger) {
	debugLogger = l
}
