// Package realtime implements the unofficial and unsupported Sense real-time API.
//
// The Sense real-time API is a WebSocket API that provides real-time updates
// about the state of your Sense monitor.  It is used by the Sense mobile app
// and the Sense web app.
//
// WARNING: Sense does not provide a supported API. This package may stop
// working without notice.
//
// The current implementation of this appears to be reasonably complete, however
// there may be a few fields with interface{} types that could be better
// investsigated.
//
// The first incoming messages appear to follow this pattern:
//
// 1. Hello
// 2. MonitorInfo
// 3. DataChange
// 4. DeviceStates
// 5. RealtimeUpdate
//
// After this, you'll get RealtimeUpdate messages once a second, with an occasional
// DataChange message.
package realtime

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/oauth2"
	"nhooyr.io/websocket"
)

const traceName = "github.com/dnesting/sense/realtime"

// hello

type helloMsg struct {
	Type    string `json:"type"`
	Payload Hello  `json:"payload"`
}

// Hello is the first message sent by the server after connecting.
type Hello struct {
	Online bool `json:"online"`
}

func (h Hello) GetType() string { return "hello" }

// monitor_info

type monitorInfoMsg struct {
	Type    string      `json:"type"`
	Payload MonitorInfo `json:"payload"`
}

// MonitorInfo is sent by the server after connecting and returns a CSV string.
type MonitorInfo struct {
	Features string `json:"features"`
}

func (m MonitorInfo) GetType() string {
	return "monitor_info"
}

type deviceStatesMsg struct {
	Type    string       `json:"type"`
	Payload DeviceStates `json:"payload"`
}

// DeviceStates is sent by the server shortly after connecting and provides
// the current "online" state of all devices.
type DeviceStates struct {
	States     []DeviceState `json:"states"`
	UpdateType string        `json:"update_type"`
}

type DeviceState struct {
	DeviceID string `json:"device_id"`
	Mode     string `json:"mode"`
	State    string `json:"state"`
}

func (d DeviceStates) GetType() string {
	return "device_states"
}

// realtime_update

type realtimeUpdateMsg struct {
	Type    string         `json:"type"`
	Payload RealtimeUpdate `json:"payload"`
}

// RealtimeUpdate is the message periodically sent by the server with the
// current state of the monitor and all known devices.
type RealtimeUpdate struct {
	C int `json:"c"`

	// Appears to be the wattage reading for each of the monitor's sensors.
	Channels []float32 `json:"channels"`
	// This appears to be the same as W but as an integer.
	DW          int     `json:"d_w"`
	DefaultCost float32 `json:"default_cost"`
	// Deltas are usually missing, but when they are present they appear to
	// contain the difference in the W value between this update and the previous.
	Deltas []Delta `json:"deltas"`
	// Devices contains details about the current consumption of known devices.
	Devices []Device `json:"devices"`
	// Epoch appears to be the Unix time_t of the start of the stream.
	Epoch int `json:"epoch"`
	// Frame appears to be a counter that increases between updates, seemingly 30 per update.
	Frame int     `json:"frame"`
	GridW float32 `json:"grid_w"`
	// This appears to be the AC frequency in Hz.
	Hz        float32 `json:"hz"`
	PowerFlow struct {
		Grid []string `json:"grid"`
	} `json:"power_flow"`
	// Appears to be the AC voltage reading on each of the monitor's sensors.
	Voltage []float32 `json:"voltage"`
	// W appears to be the total wattage observed being consumed by the monitor.
	W float32 `json:"w"`
	// These appear to be Unix time timestamps, with subsecond precision.  Possibly
	// used to gauge the latency of the stream servers.
	Stats struct {
		Brcv float32
		Mrcv float32
		Msnd float32
	}
}

func (r *RealtimeUpdate) GetType() string {
	return "realtime_update"
}

type Delta struct {
	Frame      int     `json:"frame"`
	Channel    int     `json:"channel"`
	StartFrame int     `json:"start_frame"`
	W          float32 `json:"w"`
}

type Device struct {
	Attrs interface{}            `json:"attrs"`
	Icon  string                 `json:"icon"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Tags  map[string]interface{} `json:"tags"`
	W     float32                `json:"w"`
}

// data_change

type dataChangeMsg struct {
	Type    string     `json:"type"`
	Payload DataChange `json:"payload"`
}

// DataChange is a message sent by the server when some specific pieces of
// data have changed.  These are likely used to signal web or mobile clients
// of the need to refresh their data.
type DataChange struct {
	DeviceDataChecksum      string `json:"device_data_checksum"`
	MonitorOverviewChecksum string `json:"monitor_overview_checksum"`
	PartnerChecksum         string `json:"partner_checksum"`
	PendingEvents           struct {
		Type string `json:"type"`
		Goal struct {
			Guid           string      `json:"guid"`
			NotificationID interface{} `json:"notification_id"`
			Timestamp      interface{} `json:"timestamp"`
		} `json:"goal"`
		NewDeviceFound struct {
			DeviceID  interface{} `json:"device_id"`
			Guid      string      `json:"guid"`
			Timestamp interface{} `json:"timestamp"`
		} `json:"new_device_found"`
	}
	SettingsVersion int `json:"settings_version"`
	UserVersion     int `json:"user_version"`
}

func (d DataChange) GetType() string {
	return "data_change"
}

// new_timeline_event

type newTimelineEventMsg struct {
	Type    string           `json:"type"`
	Payload NewTimelineEvent `json:"payload"`
}

type NewTimelineEvent struct {
	ItemsAdded   []TimelineEvent `json:"items_added"`
	ItemsRemoved []TimelineEvent `json:"items_removed"`
	ItemsUpdated []TimelineEvent `json:"items_updated"`
	UserID       int             `json:"user_id"`
}

func (n NewTimelineEvent) GetType() string {
	return "new_timeline_event"
}

type TimelineEvent struct {
	AllowSticky               bool                     `json:"allow_sticky"`
	Body                      string                   `json:"body"`
	BodyArgs                  []map[string]interface{} `json:"body_args"`
	BodyKey                   string                   `json:"body_key"`
	Destination               string                   `json:"destination"`
	DeviceID                  string                   `json:"device_id"`
	DeviceState               string                   `json:"device_state"`
	DeviceTransitionFromState string                   `json:"device_transition_from_state"`
	GUID                      string                   `json:"guid"`
	Icon                      string                   `json:"icon"`
	MonitorID                 int                      `json:"monitor_id"`
	ShowAction                bool                     `json:"show_action"`
	Time                      *time.Time               `json:"time"`
	Type                      string                   `json:"type"`
	UserDeviceType            string                   `json:"user_device_type"`
}

var unexpectedMessageTypes = map[string]bool{}

// Given bytes from a websocket message, parse out the Type and then
// the ultimate payload message.
func parseMessage(buf []byte) (msg Message, err error) {
	var discriminator struct {
		Type string `json:"type"`
	}
	if err = json.Unmarshal(buf, &discriminator); err != nil {
		if !unexpectedMessageTypes[""] {
			unexpectedMessageTypes[""] = true
			log.Print("unable to parse message type (future messages suppressed): ", err, "\n", hex.Dump(buf))
		}
		return nil, err
	}

	switch discriminator.Type {
	case "hello":
		var m helloMsg
		if err = json.Unmarshal(buf, &m); err == nil {
			msg = &m.Payload
		}
	case "monitor_info":
		var m monitorInfoMsg
		if err = json.Unmarshal(buf, &m); err == nil {
			msg = &m.Payload
		}
	case "device_states":
		var m deviceStatesMsg
		if err = json.Unmarshal(buf, &m); err == nil {
			msg = &m.Payload
		}
	case "realtime_update":
		var m realtimeUpdateMsg
		if err = json.Unmarshal(buf, &m); err == nil {
			msg = &m.Payload
		}
	case "data_change":
		var m dataChangeMsg
		if err = json.Unmarshal(buf, &m); err == nil {
			msg = &m.Payload
		}
	case "new_timeline_event":
		var m newTimelineEventMsg
		if err = json.Unmarshal(buf, &m); err == nil {
			msg = &m.Payload
		}
	default:
		err = fmt.Errorf("unknown message type: %s", discriminator.Type)
	}
	if err != nil {
		if !unexpectedMessageTypes[discriminator.Type] {
			unexpectedMessageTypes[discriminator.Type] = true
			log.Print("unable to parse message (future messages suppressed): ", err, "\n", hex.Dump(buf))
		}
		return nil, err
	}
	return msg, nil
}

// Message is a generic interface for all messages sent by the server.
// The GetType method will return a string indicating which type of message
// it is.
type Message interface {
	GetType() string
}

// Conn is the method set we use to interact with a websocket.
// It is used for testing.
type Conn interface {
	Read(ctx context.Context) (websocket.MessageType, []byte, error)
	Close(websocket.StatusCode, string) error
}

// Dialer is the method set we use to connect to the websockets API.
// It is used for testing.
type Dialer interface {
	Dial(ctx context.Context, url string, opts *websocket.DialOptions) (Conn, *http.Response, error)
}

type realDialer struct{}

func (d *realDialer) Dial(ctx context.Context, url string, opts *websocket.DialOptions) (Conn, *http.Response, error) {
	return websocket.Dial(ctx, url, opts)
}

// Client is a real-time data client for the Sense API.  All fields are
// optional, and will be populated with default values if not provided.
//
// The [sense.Client] has a [sense.Client.Stream] method that uses
// this client behind the scenes; you shouldn't need to instantiate one
// directly.
//
// If you do want to use it directly, you'll need a [oauth2.TokenSource],
// such as the one generated by the [github.com/dnesting/sense/senseauth] package.
type Client struct {
	BaseUrl    string
	Origin     string
	HttpClient *http.Client
	DeviceID   string
	TokenSrc   oauth2.TokenSource
	Dialer     Dialer
}

// Stop is a sentinel error that can be returned from a callback to stop the
// stream.
//
//lint:ignore ST1012 sentinel value
var Stop = errors.New("stop sentinel")

// Callback is called for each message received during a Stream call.  Use
// [Message.GetType] or a type assertion to determine the type of the message.
type Callback func(context.Context, Message) error

func (c *Client) buildRequest(monitorID int) (string, websocket.DialOptions, error) {
	opts := websocket.DialOptions{
		HTTPClient: c.HttpClient,
		HTTPHeader: http.Header{
			"Origin": []string{c.Origin},
		},
	}

	u, err := url.Parse(c.BaseUrl)
	if err != nil {
		return "", opts, err
	}
	u = u.ResolveReference(&url.URL{Path: "monitors/" + strconv.Itoa(monitorID) + "/realtimefeed"})

	params := url.Values{
		"client_type": []string{"web"},
		"ui_language": []string{"en-US"},
	}
	if c.DeviceID != "" {
		params["device_id"] = []string{c.DeviceID}
	}

	var tok *oauth2.Token
	if c.TokenSrc != nil {
		tok, err = c.TokenSrc.Token()
		if err != nil {
			return "", opts, err
		}
		params["access_token"] = []string{tok.AccessToken}
	}
	u.RawQuery = params.Encode()
	return u.String(), opts, nil
}

// Reads incoming messages from the websocket and relays them back to the messageLoop.
func readLoop(ctx context.Context, ws Conn, ch chan<- Message) error {
	for {
		mtype, buf, err := ws.Read(ctx)
		if err != nil {
			return err
		}
		if mtype != websocket.MessageText {
			continue
		}
		msg, err := parseMessage(buf)
		if err == nil {
			ch <- msg
		}
	}
}

// Spawns readLoop and calls the callback for each message received.
func messageLoop(ctx context.Context, ws Conn, callback Callback) error {
	ch := make(chan Message)
	var readErr error
	go func() {
		readErr = readLoop(ctx, ws, ch)
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return ws.Close(websocket.StatusNormalClosure, "")
		case msg, ok := <-ch:
			if !ok {
				return readErr
			}
			err := func() error {
				ctx, span := otel.Tracer(traceName).Start(ctx, fmt.Sprintf("Handle %T", msg))
				defer span.End()
				span.SetAttributes(attribute.String("message.type", msg.GetType()))
				debugf("running callback for %T", msg)
				if err := callback(ctx, msg); err != nil {
					if err == Stop {
						ws.Close(websocket.StatusNormalClosure, "")
						return nil
					}
					ws.Close(websocket.StatusInternalError, "")
					return err
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}
	}
}

// Stream opens a websocket connection to the given monitor and calls the
// callback for each message received.  The callback can return the sentinel
// error [Stop] to stop the stream.
func (c *Client) Stream(ctx context.Context, monitorID int, callback Callback) error {
	ctx, span := otel.Tracer(traceName).Start(ctx, "Stream")
	defer span.End()

	uri, opts, err := c.buildRequest(monitorID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	dialer := c.Dialer
	if dialer == nil {
		dialer = &realDialer{}
	}

	debugf("dialing %q", uri)
	ws, _, err := dialer.Dial(ctx, uri, &opts)
	if err != nil {
		err = fmt.Errorf("dial %q: %w", uri, err)
		span.RecordError(err)
		return err
	}

	if err = messageLoop(ctx, ws, callback); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			span.RecordError(err)
			return err
		}
	}
	return nil
}
