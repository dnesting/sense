// This is a dummy command used to do a basic test that the package is working correctly.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dnesting/sense"
	"github.com/dnesting/sense/realtime"
	"github.com/dnesting/sense/sensecli"
)

var (
	flagDebug = flag.Bool("debug", false, "enable debugging")
	// note: other flags set by sensecli.SetupStandardFlags()
)

func main() {
	configFile, flagCreds := sensecli.SetupStandardFlags()
	flag.Parse()

	ctx := context.Background()
	httpClient := http.DefaultClient
	if *flagDebug {
		// enable HTTP client logging
		httpClient = sense.SetDebug(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds), httpClient)
	}

	clients, err := sensecli.CreateClients(ctx,
		configFile, flagCreds,
		sense.WithHttpClient(httpClient))
	if err != nil {
		log.Fatal(err)
	}

	for _, client := range clients {
		fmt.Println("account:", client.GetAccountID())
		fmt.Println("monitors:")
		for _, monitor := range client.GetMonitors() {
			fmt.Println("- monitor: ", monitor.ID, monitor.SerialNumber)
			var count = 3
			fn := func(_ context.Context, msg realtime.Message) error {
				if rt, ok := msg.(*realtime.RealtimeUpdate); ok {
					fmt.Printf("  W: %.1f\n", rt.W)
					count--
					if count == 0 {
						return realtime.Stop
					}
				}
				return nil
			}
			if err := client.Stream(ctx, monitor.ID, fn); err != nil {
				log.Println(err)
			}
		}
	}
}
