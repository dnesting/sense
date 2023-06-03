// Package ratelimited provides a rate-limited HTTP client.
package ratelimited

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
)

const traceName = "github.com/dnesting/sense"

// Limiter is a function that is expected to wait until rate limiting
// conditions have been met.  The implementation should return early
// if the provided context expires.
type Limiter func(context.Context) error

// Transport is an http.RoundTripper that calls L to apply a rate limit
// to requests.  If Transport is nil, http.DefaultTransport is used.
// If L is nil, no rate limiting is applied.
type Transport struct {
	L    Limiter
	Base interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Apply the rate limit for this Transport.
	if t.L != nil {
		ctx, span := otel.Tracer(traceName).Start(r.Context(), "rate limit")
		if err := t.L(ctx); err != nil {
			span.RecordError(err)
			span.End()
			return nil, err
		}
		span.End()
	}
	base := t.Base
	if base == nil {
		base = http.DefaultClient
	}
	return base.Do(r)
}

// NewClient creates a new HTTP client that calls limiter to apply a
// rate limit to requests.  If client is nil, http.DefaultClient is used.
// If limiter is nil, no rate limiting is applied.
func NewClient(client *http.Client, limiter Limiter) *http.Client {
	return &http.Client{
		Transport: &Transport{
			L:    limiter,
			Base: client,
		},
	}
}
