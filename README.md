# sense

This package is an incomplete implementation of an entirely UNOFFICIAL
and UNSUPPORTED API to access data for a [Sense](https://sense.com/)
Energy Monitor account.  This repository has no affiliation with Sense.

Because this is unsupported, this package may stop working at any time.

The API in this package is not stable and I may change it at any time.

## Usage

```go
import (
	"github.com/dnesting/sense"
	"github.com/dnesting/sense/realtime"
)

func main() {
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
```

### MFA

If your account requires multi-factor authentication, you can accommodate that like:

```go
mfaFunc := func() (string, error) {
    // obtain your MFA code somehow, we'll just use a fixed value to demonstrate
    return "12345", nil
}
client, err := sense.Connect(ctx, sense.PasswordCredentials{
    Email: "you@example.com",
    Password: "secret",
    MfaFn: mfaFunc,
})
```

Your `mfaFunc` will be called when needed.

## Notes

This implementation is incomplete, and what's there is incompletely tested.
If you wish to contribute, here's how the project is laid out:

```
|-- internal
|   |-- client         contains an (incomplete) OpenAPI spec and
|   |                  auto-generated code that does the heavy lifting
|   |-- ratelimited    implements some HTTP rate limiting
|   `-- senseutil      helper functions, mocks for testing, etc.
|-- realtime           contains a complete-ish AsyncAPI spec but
|                      hand-generated code implementing the real-time
|                      WebSockets API
|-- senseauth          implements the Sense artisinal OAuth
`-- sensecli           helpers that CLI tools might find useful
```

### Debugging

If you need the gory internals to figure something out:

```go
httpClient := sense.SetDebug(log.Default(), nil)
client, err := sense.Connect(ctx, credentials, sense.WithHTTPClient(httpClient))
```

## Development

### CI/CD

This repository uses GitHub Actions for continuous integration and release validation:

- **Pull Request Checks**: All PRs are automatically tested on multiple Go versions (1.22, 1.23) and platforms (Linux, Windows, macOS). The CI runs tests, vet, format checks, and security vulnerability scanning.

- **Release Validation**: When you create a GitHub release with a tag (e.g., `v1.0.0`), the release workflow validates the code by running the full test suite and security checks to ensure the tagged version is ready for use.

To create a release:
1. Create a new tag: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. Create a GitHub release from the tag
4. The release workflow will validate the code quality

Users can then import the library using: `go get github.com/dnesting/sense@v1.0.0`
