package sense_test

import (
	"context"
	"log"
	"testing"

	"github.com/dnesting/sense"
)

func TestClient(t *testing.T) {
	// TODO
}

func doSomethingWith(c *sense.Client) {}

func ExampleConnect() {
	client, err := sense.Connect(
		context.Background(),
		sense.PasswordCredentials{
			Email:    "you@example.com",
			Password: "secret",
		})
	if err != nil {
		log.Fatal(err)
	}
	doSomethingWith(client)
}

func ExampleConnect_withMFA() {
	mfaFunc := func(_ context.Context) (string, error) {
		// obtain your MFA code somehow, we'll just use a fixed value to demonstrate
		return "12345", nil
	}
	client, err := sense.Connect(
		context.Background(),
		sense.PasswordCredentials{
			Email:    "you@example.com",
			Password: "secret",
			MfaFn:    mfaFunc, // <--
		})
	if err != nil {
		log.Fatal(err)
	}
	doSomethingWith(client)
}
