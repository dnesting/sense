package sense_test

import (
	"context"
	"fmt"
	"log"

	"github.com/dnesting/sense"
	"github.com/dnesting/sense/realtime"
)

func Example() {
	ctx := context.Background()
	client, err := sense.Connect(ctx, sense.PasswordCredentials{
		Email:    "you@example.com",
		Password: "secret",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Monitors configured under account", client.AccountID)
	for _, m := range client.Monitors {
		fmt.Println("-", m.ID)

		// Use the realtime Stream API to grab one data point.
		fn := func(_ context.Context, msg realtime.Message) error {
			if rt, ok := msg.(*realtime.RealtimeUpdate); ok {
				fmt.Println("  current power consumption", rt.W, "W")
				return realtime.Stop
			}
			return nil
		}
		if err := client.Stream(ctx, m.ID, fn); err != nil {
			log.Fatal(err)
		}
	}
}
