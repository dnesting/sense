package realtime_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/dnesting/sense/internal/senseutil"
	"github.com/dnesting/sense/realtime"
	"github.com/dnesting/sense/senseauth"
	"nhooyr.io/websocket"
)

type msg = senseutil.WSMsg

func TestStream(t *testing.T) {
	ch := make(chan msg, 1)
	dialer := &senseutil.MockWSDialer{Ch: ch}
	client := &realtime.Client{
		Dialer:  dialer,
		BaseUrl: "wss://clientrt.example.test/path/",
	}

	ch <- msg{T: websocket.MessageText, D: `{"type":"hello","payload":{"online": true}}`, E: nil}

	var gotMsg realtime.Message
	err := client.Stream(context.Background(), 123, func(_ context.Context, msg realtime.Message) error {
		gotMsg = msg
		return realtime.Stop
	})

	expectedUrl := "wss://clientrt.example.test/path/monitors/123/realtimefeed?client_type=web&ui_language=en-US"
	if dialer.Dialed != expectedUrl {
		t.Errorf("expected dialer to dial\n%q, got\n%q", expectedUrl, dialer.Dialed)
	}

	if err != nil {
		t.Error("unexpected error:", err)
	}
	if gotMsg == nil {
		t.Error("expected message, got nil")
	}
	if gotMsg.GetType() != "hello" {
		t.Errorf("expected type=hello, got %q", gotMsg.GetType())
	}
	if hello, ok := gotMsg.(*realtime.Hello); ok {
		if !hello.Online {
			t.Error("expected online=true, got false")
		}
	} else {
		t.Errorf("expected *realtime.Hello, got %T", gotMsg)
	}
}

func TestMessageTypes(t *testing.T) {
	testCases := []struct {
		msgType string
		data    string
	}{
		{"hello", `{"type":"hello","payload":{"online": true}}`},
		{"monitor_info", `{"payload":{"features":"MONITOR_RESET,DEVICE_DELETION,DEVICE_STATE_UPDATE,NDI,UPDATE_CONFIG_ACK,HUE,NETEVENT_WM,TPL_CTRL,WEMO_CTRL,SMARTPLUGS,HS300,AO_SANS_NDI,CONNECTION_TEST,GENERATOR,SSI,SPLIT_PANEL,DEDICATED_CIRCUIT,RTS_HISTORICAL,WISER_HOME_DEVICES,RTU_WATTS_SIGNED,REGISTRATION_TOKEN,WISER_RELAY"},"type":"monitor_info"}`},
		{"data_change", `{
			"payload" : {
			   "device_data_checksum" : "B79DC787CB404D288169413E1C0B72A0B4476037",
			   "monitor_overview_checksum" : "2ED7614F7DD0909F83F43B70658244D10CCB7C57",
			   "partner_checksum" : "FE8BB93DCBA62E8CE8A062D724DBC8A69309BD57",
			   "pending_events" : {
				  "goal" : {
					 "guid" : "0",
					 "notification_id" : null,
					 "timestamp" : null
				  },
				  "monitor_id" : 216148,
				  "new_device_found" : {
					 "device_id" : null,
					 "guid" : "0",
					 "timestamp" : null
				  },
				  "type" : "GoalEvent"
			   },
			   "settings_version" : 1,
			   "user_version" : 2
			},
			"type" : "data_change"
		 }`},
		{"device_states", `{
			"payload" : {
			   "states" : [
				  {
					 "device_id" : "ssi-12345678",
					 "mode" : "active",
					 "state" : "online"
				  },
				  {
					 "device_id" : "ssi-23456789",
					 "mode" : "off",
					 "state" : "online"
				  }
			   ],
			   "update_type" : "full"
			},
			"type" : "device_states"
		 }`},
		{"realtime_update", `{
			"payload" : {
			   "_stats" : {
				  "brcv" : 1685500000.13076,
				  "mrcv" : 1685500000.157,
				  "msnd" : 1685000005.157
			   },
			   "c" : 4,
			   "channels" : [
				  327.123450703125,
				  264.234565625
			   ],
			   "d_w" : 590,
			   "defaultCost" : 7.00,
			   "deltas" : [],
			   "devices" : [
				  {
					 "attrs" : [],
					 "icon" : "alwayson",
					 "id" : "always_on",
					 "name" : "Always On",
					 "tags" : {
						"DefaultUserDeviceType" : "AlwaysOn",
						"DeviceListAllowed" : "true",
						"TimelineAllowed" : "false",
						"UserDeleted" : "false",
						"UserDeviceType" : "AlwaysOn",
						"UserDeviceTypeDisplayString" : "Always On",
						"UserEditable" : "false",
						"UserMergeable" : "false",
						"UserShowBubble" : "true",
						"UserShowInDeviceList" : "true"
					 },
					 "w" : 300
				  },
				  {
					 "attrs" : [],
					 "icon" : "home",
					 "id" : "unknown",
					 "name" : "Other",
					 "tags" : {
						"DefaultUserDeviceType" : "Unknown",
						"DeviceListAllowed" : "true",
						"TimelineAllowed" : "false",
						"UserDeleted" : "false",
						"UserDeviceType" : "Unknown",
						"UserDeviceTypeDisplayString" : "Unknown",
						"UserEditable" : "false",
						"UserMergeable" : "false",
						"UserShowBubble" : "true",
						"UserShowInDeviceList" : "true"
					 },
					 "w" : 123.96463
				  },
				  {
					 "attrs" : [],
					 "given_location" : "Kitchen",
					 "given_make" : "Signify Netherlands B.V.",
					 "icon" : "lightbulb",
					 "id" : "12345678",
					 "location" : "Kitchen",
					 "make" : "Signify Netherlands B.V.",
					 "name" : "Kitchen lights",
					 "sd" : {
						"extra" : {
						   "bri" : 254,
						   "ct" : 198,
						   "hue" : 40703,
						   "sat" : 42
						},
						"intensity" : 1,
						"w" : 7.83000272508612
					 },
					 "tags" : {
						"Alertable" : "true",
						"ControlCapabilities" : [
						   "OnOff",
						   "Brightness"
						],
						"DateCreated" : "2020-01-01T18:20:41.000Z",
						"DefaultLocation" : "Kitchen",
						"DefaultMake" : "Signify Netherlands B.V.",
						"DefaultUserDeviceType" : "Light",
						"DeviceListAllowed" : "true",
						"IntegrationType" : "Hue",
						"MergedDevices" : "ssi-12345678,ssi-23456789",
						"OriginalName" : "Kitchen lights",
						"PeerNames" : [],
						"Revoked" : "false",
						"SSIEnabled" : "true",
						"SSIModel" : "SelfReporting",
						"TimelineAllowed" : "true",
						"TimelineDefault" : "true",
						"UserDeletable" : "false",
						"UserDeviceType" : "Light",
						"UserDeviceTypeDisplayString" : "Light",
						"UserEditable" : "true",
						"UserEditableMeta" : "false",
						"UserMergeable" : "false",
						"UserShowBubble" : "true",
						"UserShowInDeviceList" : "true",
						"Virtual" : "true",
						"name_useredit" : "false"
					 },
					 "w" : 30.72001
				  }
			   ],
			   "epoch" : 1685567126,
			   "frame" : 13285440,
			   "grid_w" : 590,
			   "hz" : 59.9812185668945,
			   "power_flow" : {
				  "grid" : [
					 "home"
				  ]
			   },
			   "voltage" : [
				  123.177742004395,
				  123.033126831055
			   ],
			   "w" : 590.417236328125
			},
			"type" : "realtime_update"
		 }
		 `},
		{"new_timeline_event", `{
		   "payload" : {
		      "items_added" : [
		         {
		            "allow_sticky" : false,
		            "body" : "Device 3 turned off",
		            "body_args" : [
		               {
		                  "@type" : "String",
		                  "is_key" : false,
		                  "value" : "Device 3"
		               }
		            ],
		            "body_key" : "device_turned_off",
		            "destination" : "device:12345678",
		            "device_id" : "12345678",
		            "device_state" : "DeviceOff",
		            "device_transition_from_state" : "On",
		            "guid" : "12345678-1234-1234-1234-123456789012",
		            "icon" : "socket",
		            "monitor_id" : 12345,
		            "show_action" : false,
		            "time" : "2023-06-01T00:01:01.862Z",
		            "type" : "DeviceOff",
		            "user_device_type" : "MysteryDevice"
		         }
		      ],
		      "items_removed" : [],
		      "items_updated" : [],
		      "user_id" : 12345
		   },
		   "type" : "new_timeline_event"
		}
			`},
	}

	ch := make(chan msg)
	dialer := &senseutil.MockWSDialer{Ch: ch}
	client := &realtime.Client{
		Dialer: dialer,
	}

	go func() {
		for _, test := range testCases {
			ch <- msg{T: websocket.MessageText, D: test.data, E: nil}
		}
		close(ch)
	}()

	var gotMsgs = make(map[string]bool)
	err := client.Stream(context.Background(), 123, func(_ context.Context, msg realtime.Message) error {
		typ := msg.GetType()
		gotMsgs[typ] = true

		switch typ {
		case "hello":
			if _, ok := msg.(*realtime.Hello); !ok {
				t.Errorf("expected %s message, got %T", typ, msg)
			}
		case "monitor_info":
			if _, ok := msg.(*realtime.MonitorInfo); !ok {
				t.Errorf("expected %s message, got %T", typ, msg)
			}
		case "data_change":
			if _, ok := msg.(*realtime.DataChange); !ok {
				t.Errorf("expected %s message, got %T", typ, msg)
			}
		case "device_states":
			if _, ok := msg.(*realtime.DeviceStates); !ok {
				t.Errorf("expected %s message, got %T", typ, msg)
			}
		case "realtime_update":
			if _, ok := msg.(*realtime.RealtimeUpdate); !ok {
				t.Errorf("expected %s message, got %T", typ, msg)
			}
		case "new_timeline_event":
			if _, ok := msg.(*realtime.NewTimelineEvent); !ok {
				t.Errorf("expected %s message, got %T", typ, msg)
			}
		default:
			t.Errorf("got unexpected message type: %s", typ)
		}
		return nil
	})

	for _, test := range testCases {
		if !gotMsgs[test.msgType] {
			t.Errorf("expected to receive %s message", test.msgType)
		}
	}

	if err != nil {
		t.Error("unexpected error:", err)
	}
}

type mockTransport struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTrip(req)
}

func mockForExample(cl *realtime.Client, auth *senseauth.Config) {
	realtime.SetDebug(log.Default())
	ch := make(chan msg, 1)
	dialer := &senseutil.MockWSDialer{Ch: ch}
	cl.Dialer = dialer
	go func() {
		ch <- msg{T: websocket.MessageText, D: `{"type":"hello"}`, E: nil}
		ch <- msg{T: websocket.MessageText, D: `{"type":"realtime_update","payload":{"W": 590.4}}`, E: nil}
		ch <- msg{T: websocket.MessageText, D: `{"type":"realtime_update","payload":{"W": 591.4}}`, E: nil}
		ch <- msg{T: websocket.MessageText, D: `{"type":"realtime_update","payload":{"W": 592.4}}`, E: nil}
		ch <- msg{T: websocket.MessageText, D: `{"type":"realtime_update","payload":{"W": 593.4}}`, E: nil} // shouldn't see
		close(ch)
	}()

	auth.HttpClient = &http.Client{
		Transport: &mockTransport{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-token"}`)),
				}, nil
			},
		},
	}

}

func Example() {
	var client realtime.Client
	ctx := context.Background()

	// authenticate
	auth := senseauth.DefaultConfig
	mockForExample(&client, &auth)

	tok, _, err := auth.PasswordCredentialsToken(ctx,
		senseauth.PasswordCredentials{
			Email:    "user@example.com",
			Password: "pass",
		})
	if err != nil {
		log.Fatal(err)
	}
	client.TokenSrc = auth.TokenSource(tok)

	// start a stream and collect 3 data points
	stopAfter := 3
	err = client.Stream(ctx, 123, func(_ context.Context, msg realtime.Message) error {
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
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// We're online!
	// Power consumption is now: 590.4 W
	// Power consumption is now: 591.4 W
	// Power consumption is now: 592.4 W
}
