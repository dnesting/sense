package sense_test

import (
	"context"
	"log"
	"testing"

	"github.com/dnesting/sense"
)

func TestClient(t *testing.T) {
	// Test that an unauthenticated client returns zero values
	client := sense.New()

	if userID := client.GetUserID(); userID != 0 {
		t.Errorf("Expected GetUserID() to return 0 for unauthenticated client, got %d", userID)
	}

	if accountID := client.GetAccountID(); accountID != 0 {
		t.Errorf("Expected GetAccountID() to return 0 for unauthenticated client, got %d", accountID)
	}

	if monitors := client.GetMonitors(); monitors != nil {
		t.Errorf("Expected GetMonitors() to return nil for unauthenticated client, got %v", monitors)
	}
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
