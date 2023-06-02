package sense

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dnesting/sense/internal/senseutil"
	"github.com/dnesting/sense/realtime"
	"github.com/dnesting/sense/senseauth"
)

var debugLogger *log.Logger

type loggingTransport struct {
	transport http.RoundTripper
}

func debug(args ...interface{}) {
	if debugLogger != nil {
		debugLogger.Output(2, strings.TrimRight(fmt.Sprint(args...), "\n"))
	}
}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	senseutil.DumpRequest(debugLogger, r)

	tr := s.transport
	if tr == nil {
		tr = http.DefaultTransport
	}
	resp, err := tr.RoundTrip(r)
	if err != nil {
		debug("error:", err)
		return nil, err
	}
	senseutil.DumpResponse(debugLogger, resp)
	return resp, err
}

// SetDebug enables debug logging using the given logger and returns an
// [*http.Client] that wraps baseClient, logging requests and responses
// to the same logger.  Passing nil will disable debug logging.
func SetDebug(l *log.Logger, baseClient *http.Client) *http.Client {
	if baseClient == nil {
		baseClient = http.DefaultClient
	}
	debugLogger = l
	senseauth.SetDebug(l)
	realtime.SetDebug(l)
	return &http.Client{
		Transport: &loggingTransport{baseClient.Transport},
	}
}
