package sense_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dnesting/sense"
	"github.com/dnesting/sense/internal/senseutil"
	"github.com/dnesting/sense/realtime"
)

func TestSense(t *testing.T) {

}

func mockForExample() []sense.Option {
	type mockMsg = senseutil.RTMsg

	ch := make(chan mockMsg)
	go func() {
		ch <- mockMsg{M: &realtime.Hello{}, E: nil}
		ch <- mockMsg{M: &realtime.RealtimeUpdate{W: 590.4}, E: nil}
		ch <- mockMsg{M: &realtime.RealtimeUpdate{W: 591.4}, E: nil}
		ch <- mockMsg{M: &realtime.RealtimeUpdate{W: 592.4}, E: nil}
		ch <- mockMsg{M: &realtime.RealtimeUpdate{W: 593.4}, E: nil} // shouldn't see
		close(ch)
	}()
	return []sense.Option{
		sense.WithInternalClient(nil, &senseutil.MockRTClient{Ch: ch}),
		sense.WithHttpClient(&http.Client{
			Transport: &senseutil.MockTransport{
				RT: func(req *http.Request) (*http.Response, error) {
					fmt.Fprintln(os.Stderr, req.Method, req.URL)
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-token"}`)),
					}, nil
				},
			},
		})}
}

func ExampleClient_Stream() {
	// instantiate a Sense client (error checking omitted)
	client, _ := sense.Connect(
		context.Background(),
		sense.PasswordCredentials{
			Email:    "test@example.com",
			Password: "pass",
		},
		mockForExample()...)

	// start a stream and collect 3 data points
	stopAfter := 3
	client.Stream(
		context.Background(),
		123, // monitor ID to stream events from

		func(_ context.Context, msg realtime.Message) error {
			switch msg := msg.(type) {

			case *realtime.Hello:
				fmt.Println("We're online!")

			case *realtime.RealtimeUpdate:
				fmt.Printf("Power consumption is now: %.1f W\n", msg.W)
				stopAfter--
				if stopAfter == 0 {
					return realtime.Stop
				}

			}
			return nil
		})

	// Output:
	// We're online!
	// Power consumption is now: 590.4 W
	// Power consumption is now: 591.4 W
	// Power consumption is now: 592.4 W
}
