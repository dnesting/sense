package sense

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// Implement the Options pattern for New and Connect.  Minimize the amount of
// state kept in this struct, since it will be re-used when needed to generate
// new clients.
type newOptions struct {
	httpClient     *http.Client
	rateLimit      *rate.Limit
	apiUrl         string
	realtimeApiUrl string
	realtimeOrigin string

	internalClient         internalClient
	internalRealtimeClient internalRealtimeClient
	deviceID               string
}

// Option is a function that can be passed to New or Connect to configure the
// resulting Client.
type Option func(*newOptions)

// WithHttpClient sets the HTTP client used to make requests to the Sense API.
// If this option is not provided, http.DefaultClient will be used.
//
// This option can be useful if you need special handling for proxies or TLS
// certificates, or if you have your own approach to authenticating with Sense.
func WithHttpClient(httpClient *http.Client) Option {
	return func(o *newOptions) {
		o.httpClient = httpClient
	}
}

// WithRateLimit applies a rate limit to requests made to the Sense API.
// Without this option, a default rate limit of 10 requests/second will be applied.
// You can disable rate limiting by providing a limit of 0.
func WithRateLimit(limit rate.Limit) Option {
	return func(o *newOptions) {
		o.rateLimit = &limit
	}
}

// WithApiUrl sets the base URLs for the Sense API.
// If this option is not provided, the standard production API URLs will be used
// (https://api.sense.com/apiservice/api/v1/).
func WithApiUrl(apiUrl, realtimeApiUrl string) Option {
	return func(o *newOptions) {
		o.apiUrl = apiUrl
		o.realtimeApiUrl = realtimeApiUrl
	}
}

// WithInternalClient is for internal use. All other options will be ignored.
//
// Deprecated: For internal use.
func WithInternalClient(cl internalClient, acl internalRealtimeClient) Option {
	return func(o *newOptions) {
		o.internalClient = cl
		o.internalRealtimeClient = acl
	}
}

// WithDeviceID sets the X-Sense-Device-Id header on requests to the Sense API.
// This appears to be intended to be a unique identifier for the client installation.
// If this option is not provided, a random value will be generated and used for all clients for the life of this process.
// Set this value to "" explicitly to disable this header.
func WithDeviceID(id string) Option {
	return func(o *newOptions) {
		o.deviceID = id
	}
}

func getOptions(build newOptions, opts ...Option) *newOptions {
	for _, o := range opts {
		o(&build)
	}
	// Be a little paranoid about nils
	if build.httpClient == nil {
		build.httpClient = defaultOptions.httpClient
	}
	if build.apiUrl == "" {
		build.apiUrl = defaultOptions.apiUrl
	}
	return &build
}

var defaultOptions = newOptions{
	httpClient:     http.DefaultClient,
	apiUrl:         defaultApiRoot,
	realtimeApiUrl: defaultRealtimeApiRoot,
	deviceID:       uuid.New().String(),
}
