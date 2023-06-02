package sense

import (
	"context"

	"github.com/dnesting/sense/realtime"
	"golang.org/x/oauth2"
)

const (
	defaultRealtimeApiRoot = "wss://clientrt.sense.com/"
	defaultRealtimeOrigin  = "https://home.sense.com"
)

// Stream begins streaming real-time data via callback.  If the callback returns
// realtime.Stop, the stream will be closed and this function will return without error.
// Otherwise, if any other error occurs, it will be returned.
func (s *Client) Stream(ctx context.Context, monitor int, callback realtime.Callback) error {
	return s.realtimeClient.Stream(ctx, monitor, callback)
}

// Deprecated: For use with testing.
type internalRealtimeClient interface {
	Stream(ctx context.Context, monitor int, callback realtime.Callback) error
}

func newRealtimeClient(opts *newOptions, src oauth2.TokenSource) internalRealtimeClient {
	if opts.internalRealtimeClient != nil {
		// for testing
		return opts.internalRealtimeClient
	}
	c := &realtime.Client{
		HttpClient: opts.httpClient,
		BaseUrl:    opts.realtimeApiUrl,
		DeviceID:   opts.deviceID,
		Origin:     opts.realtimeOrigin,
		TokenSrc:   src,
	}
	if c.Origin == "" {
		c.Origin = defaultRealtimeOrigin
	}
	return c
}
