package senseutil

import (
	"context"
	"io"
	"net/http"

	"github.com/coder/websocket"
	"github.com/dnesting/sense/realtime"
)

type MockTransport struct {
	RT func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RT(req)
}

type RTMsg struct {
	M realtime.Message
	E error
}

type MockRTClient struct {
	Ch <-chan RTMsg
}

func (c *MockRTClient) Stream(ctx context.Context, deviceID int, f realtime.Callback) error {
	for msg := range c.Ch {
		if msg.E != nil {
			return msg.E
		}
		if err := f(ctx, msg.M); err != nil {
			if err == realtime.Stop {
				return nil
			}
			return err
		}
	}
	return nil
}

type WSMsg struct {
	T websocket.MessageType
	D string
	E error
}

type MockWSConn struct {
	Ch <-chan WSMsg
}

func (mc *MockWSConn) Close(_ websocket.StatusCode, _ string) error {
	return nil
}

func (mc *MockWSConn) Read(_ context.Context) (websocket.MessageType, []byte, error) {
	msg, ok := <-mc.Ch
	if ok {
		return msg.T, []byte(msg.D), msg.E
	}
	return 0, nil, io.EOF
}

var _ realtime.Conn = &MockWSConn{}

type MockWSDialer struct {
	Dialed string
	Ch     <-chan WSMsg
}

var _ realtime.Dialer = &MockWSDialer{}

func (d *MockWSDialer) Dial(ctx context.Context, urlStr string, opts *websocket.DialOptions) (realtime.Conn, *http.Response, error) {
	d.Dialed = urlStr
	return &MockWSConn{
		Ch: d.Ch,
	}, nil, nil
}
