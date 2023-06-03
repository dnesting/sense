// Package sense implements a high-level client for the UNSUPPORTED Sense Energy API.
//
// WARNING: Sense does not provide a supported API. This package may stop
// working without notice.
package sense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dnesting/sense/internal/client"
	"github.com/dnesting/sense/internal/ratelimited"
	"github.com/dnesting/sense/senseauth"
	"go.opentelemetry.io/otel"
	"golang.org/x/time/rate"
)

const (
	defaultApiRoot   = "https://api.sense.com/apiservice/api/v1/"
	defaultRateLimit = rate.Limit(time.Second / 10) // arbitrarily chosen
	userAgent        = "go-sense-library (github.com/dnesting/sense)"
	traceName        = "github.com/dnesting/sense"
)

// Client is the primary high-level object used to interact with the Sense API.
// It represents an "account", which can have some number of Monitors.
// Instantiate a Client using [New] or [Connect].
type Client struct {
	// Account fields are set after successful authentication.
	UserID    int
	AccountID int
	Monitors  []Monitor

	client         internalClient
	realtimeClient internalRealtimeClient
	opt            newOptions
}

// Monitor is a Sense monitor, which is a physical device that measures power usage.
// One account can have multiple Monitors.
type Monitor struct {
	ID           int
	SerialNumber string
}

// PasswordCredentials holds the credentials used to authenticate to the Sense API.
//
// MfaFn is an optional function that will be called (if it is provided)
// if the Sense API requests an MFA code for the account.  It should return
// the MFA code. If it returns an error, authentication will fail with
// that error.
type PasswordCredentials struct {
	Email    string
	Password string
	MfaFn    func(ctx context.Context) (string, error)
}

func (u PasswordCredentials) internalOnly() {}

// Credentials holds the credentials used to authenticate to the Sense API.
// The only implementation of this is [PasswordCredentials].
type Credentials interface {
	internalOnly()
}

var (
	_ senseauth.PasswordCredentials = senseauth.PasswordCredentials(PasswordCredentials{})
	_ Credentials                   = &PasswordCredentials{}
)

// Deprecated: For use with testing.
type internalClient interface {
	client.ClientInterface
	client.ClientWithResponsesInterface
}

// New creates a new unauthenticated Sense client, configured according to
// the provided options.
//
// Most callers will prefer to use [Connect] instead.
func New(opts ...Option) *Client {
	opt := getOptions(defaultOptions, opts...)
	return newClient(opt)
}

func newClient(opt *newOptions) (cl *Client) {
	cl = &Client{opt: *opt}
	cl.client = newInternalClient(opt)
	cl.realtimeClient = newRealtimeClient(opt, nil)
	return cl
}

// newInternalClient constructs a new [client.Client] from the provided options.
// Most of the handling for options occurs here.
func newInternalClient(opt *newOptions) (cl internalClient) {
	// Use caution when modifying fields in opt, since these fields will
	// be re-used if the client needs to be reset, such as by a call
	// to Authenticate().
	// Changes to opt.httpClient for instance may result in unnecessary
	// wrapping of the HTTP client.

	// If the user (tests) provided a client with WithInternalClient, use it and
	// stop processing options.
	cl = opt.internalClient
	if cl != nil {
		return cl
	}

	var httpClient = opt.httpClient

	// Apply rate-limiting, unless the user set a rate limit of 0.
	var rl rate.Limit
	if opt.rateLimit == nil {
		rl = defaultRateLimit
	} else if *opt.rateLimit != 0 {
		rl = *opt.rateLimit
	}
	burst := time.Second / time.Duration(rl)
	if burst < 1 {
		burst = 1
	}
	if rl != 0 {
		httpClient = ratelimited.NewClient(httpClient, rate.NewLimiter(rl, int(burst)).Wait)
	}
	headers := map[string]string{
		"User-Agent": userAgent,
	}
	if opt.deviceID != "" {
		headers["X-Sense-Device-ID"] = opt.deviceID
	}
	httpClient = addHeaders(httpClient, headers)

	cl, err := client.NewClientWithResponses(opt.apiUrl, client.WithHTTPClient(httpClient))
	if err != nil {
		// this should never happen
		panic(err)
	}
	return cl
}

func addHeaders(base *http.Client, headers map[string]string) *http.Client {
	return &http.Client{
		Transport: &addHeadersHttpTransport{
			headers: headers,
			base:    base,
		},
	}
}

type addHeadersHttpTransport struct {
	headers map[string]string
	base    *http.Client
}

func (t *addHeadersHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Set(k, v)
	}
	return t.base.Do(req)
}

// Connect instantiates a new Sense [Client], configured with the provided options.
//
// If credentials are provided, the client will be authenticated using those credentials.
// Otherwise it will be unauthenticated and will have limited abilities.
// This function is equivalent to calling [New] (opts...), possibly followed by [Client.Authenticate] (ctx).
func Connect(ctx context.Context, creds Credentials, opts ...Option) (*Client, error) {
	s := New(opts...)
	if creds != nil {
		return s, s.Authenticate(ctx, creds)
	}
	return s, nil
}

// ErrAuthenticationNeeded is wrapped by errors returned from many functions in this package
// whenever authentication is needed and the client is unauthenticated or its credentials
// are no longer valid.
//
// Test for this using errors.Is(err, sense.ErrAuthenticationNeeded).
var ErrAuthenticationNeeded = client.ErrAuthenticationNeeded

// Authenticate authenticates the client using the provided credentials.
// If the client was previously authenticated (including with Connect),
// those credentials will be replaced.
// If creds is nil, the client will be unauthenticated.
//
// See the [senseauth] package if you need more direct control over how
// the user is authenticated.  This package can generate an HTTP client
// that you can use here with [WithHttpClient].
func (s *Client) Authenticate(ctx context.Context, creds Credentials) error {
	ctx, span := otel.Tracer(traceName).Start(ctx, "Authenticate")
	defer span.End()

	// reset to unauthenticated state
	s.client = newInternalClient(&s.opt)
	s.realtimeClient = newRealtimeClient(&s.opt, nil)
	if creds == nil {
		return nil
	}

	if _, ok := creds.(*PasswordCredentials); ok {
		creds = *(creds.(*PasswordCredentials))
	}

	// The meat of authentication is handled by the senseauth package.
	screds := senseauth.PasswordCredentials(creds.(PasswordCredentials))
	config := senseauth.DefaultConfig
	config.InternalSenseClient = s.client // use our client in case it was customized
	tok, httpResponse, err := config.PasswordCredentialsToken(ctx, screds)
	if err != nil {
		return err
	}
	var hello client.Hello
	if err := json.NewDecoder(httpResponse.Body).Decode(&hello); err != nil {
		return fmt.Errorf("sense: authenticate: parse response: %w", err)
	}

	// We have an authentication token, so we can now build the HTTP client
	// that we want our Sense client to use.
	opt := s.opt // copy because we'll be overriding things we don't want to be persistent
	tokenSrc := config.TokenSource(tok)
	opt.httpClient = senseauth.NewClientFrom(opt.httpClient, tokenSrc)

	// Re-create the clients using this new authenticated HTTP client.
	s.client = newInternalClient(&opt)
	s.realtimeClient = newRealtimeClient(&opt, tokenSrc)

	s.UserID = deref(hello.UserId)
	s.AccountID = deref(hello.AccountId)
	for _, m := range deref(hello.Monitors) {
		s.Monitors = append(s.Monitors, Monitor{
			ID:           deref(m.Id),
			SerialNumber: deref(m.SerialNumber),
		})
	}

	return nil
}

// deref accepts a pointer type and returns the dereferenced value,
// or the underlying type's zero value.
func deref[T any](v *T) (out T) {
	if v == nil {
		return
	}
	return *v
}
