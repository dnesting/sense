package realtime

import (
	"fmt"
	"log"
	"strings"
)

var debugLogger *log.Logger

func debugf(format string, args ...interface{}) {
	if debugLogger != nil {
		debugLogger.Output(2, strings.TrimRight(fmt.Sprintf(format, args...), "\n"))
	}
}

/*
func debug(args ...interface{}) {
	if debugLogger != nil {
		debugLogger.Output(2, strings.TrimRight(fmt.Sprint(args...), "\n"))
	}
}
*/

// SetDebug enables debug logging using the given logger. Set to nil to disable.
func SetDebug(l *log.Logger) {
	debugLogger = l
}
