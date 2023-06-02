// Package senseauth implements the api.sense.com OAuth flow.
//
// Sense uses its own sort of OAuth flow here, which is not easily supported
// by the standard golang.org/x/oauth2 package.  Where possible, we make use of
// types and patterns from the standard package.
//
// To use, first authenticate with Sense:
//
//	conf := senseauth.DefaultConfig
//	tok, err := conf.PasswordCredentialsToken(ctx, username, password)
//	if err != nil {
//		return err
//	}
//
// Then use the token to create an HTTP client:
//
//	client := conf.Client(tok)
//
// There is also a TokenSource implementation that should be otherwise
// compatible with the types in the oauth2 package.
package senseauth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dnesting/sense/internal/client"
	"golang.org/x/oauth2"
)

type internalClient interface {
	client.ClientInterface
	client.ClientWithResponsesInterface
}

// Config is the configuration for the Sense OAuth flow.
type Config struct {
	// HttpClient is the client used for HTTP requests used by this flow.
	// If nil, http.DefaultClient is used.
	HttpClient *http.Client

	// BaseURL is the base URL for the Sense API.  If empty, "https://api.sense.com/apiservice/api/v1" is used.
	BaseURL string

	// InternalSenseClient is used internally.
	InternalSenseClient internalClient
}

var defaultApiUrl = "https://api.sense.com/apiservice/api/v1"

// DefaultConfig is the default configuration for the Sense OAuth flow.
var DefaultConfig = Config{
	BaseURL: defaultApiUrl,
}

type MfaFunc func(ctx context.Context) (string, error)

func (c Config) http() *http.Client {
	if c.HttpClient != nil {
		return c.HttpClient
	}
	return http.DefaultClient
}

// ErrMFANeeded indicates that MFA is needed to authenticate but no MFA func was provided.
var ErrMFANeeded = errors.New("senseauth: MFA needed")

// PasswordCredentials holds the credentials used to authenticate to the Sense API.
//
// MfaFn is an optional function that returns the MFA code.
// This may be used during authentication if the Sense account requires MFA.
// It should return the MFA code or an error to indicate no MFA code is available
// and to abort the authentication.
// If implementations need to block on something, consider respecting ctx.Done().
type PasswordCredentials struct {
	Email    string
	Password string
	MfaFn    func(ctx context.Context) (string, error)
}

func deref[T any](v *T) (result T) {
	if v == nil {
		return
	}
	return *v
}

func (c *Config) getClient() internalClient {
	cl := c.InternalSenseClient
	if cl == nil {
		var err error
		// Use our own unauthenticated Sense client to perform the authentication requests.
		cl, err = client.NewClientWithResponses(c.BaseURL, client.WithHTTPClient(c.http()))
		if err != nil {
			// shouldn't happen
			panic("senseauth: failed to create internal client: " + err.Error())
		}
		c.InternalSenseClient = cl
	}
	return cl
}

// PasswordCredentialsToken authenticates with the given email and password, and
// returns the token.
// Since Sense provides additional information in the authentication response,
// the response is returned as well so that the caller can do something with it.
// The caller does not have to call Close on the response Body.
// If multi-factor authentication is enabled for the account, and creds.MfaFn is
// not nil, it will be called to obtain the MFA code.
// If this returns an error, the authentication will be aborted.
func (c Config) PasswordCredentialsToken(ctx context.Context, creds PasswordCredentials) (tok *oauth2.Token, httpResponse *http.Response, err error) {
	if debugging() {
		defer func() {
			if err != nil {
				debug(err)
			} else {
				debug("senseauth: auth successful:", tok)
			}
		}()
	}

	cl := c.getClient()
	request := client.AuthenticateFormdataRequestBody{
		Email:    &creds.Email,
		Password: &creds.Password,
	}
	httpResponse, err = cl.AuthenticateWithFormdataBody(ctx, request)
	if err != nil {
		return nil, httpResponse, fmt.Errorf("senseauth: authenticate: %w", err)
	}
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, httpResponse, fmt.Errorf("senseauth: read response body: %w", err)
	}
	httpResponse.Body.Close()
	httpResponse.Body = io.NopCloser(strings.NewReader(string(body)))
	res, err := client.ParseAuthenticateResponse(httpResponse)
	if err != nil {
		return nil, nil, fmt.Errorf("senseauth: parse authenticate response: %w", err)
	}
	if res.StatusCode() == http.StatusUnauthorized && res.JSON401.MfaToken != nil {
		// We need MFA
		if creds.MfaFn == nil {
			return nil, nil, fmt.Errorf("%w: %s", ErrMFANeeded, creds.Email)
		}
		debug("senseauth: calling MFA function")
		mfaValue, err := creds.MfaFn(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("senseauth: mfa: %w", err)
		}
		request.MfaToken = res.JSON401.MfaToken
		request.Totp = &mfaValue
		res, err = cl.AuthenticateWithFormdataBodyWithResponse(ctx, request)
		if err != nil {
			return nil, nil, fmt.Errorf("senseauth: authenticate with mfa: %w", err)
		}
	}

	if err := client.Ensure(nil, "", res, 200); err != nil {
		return nil, nil, fmt.Errorf("senseauth: %w", err)
	}
	data := res.JSON200
	httpResponse.Body = io.NopCloser(strings.NewReader(string(body)))
	if deref(data.AccessToken) == "" {
		return nil, httpResponse, fmt.Errorf("senseauth: %s: not authorized", creds.Email)
	}
	tok = &oauth2.Token{
		AccessToken:  deref(data.AccessToken),
		RefreshToken: deref(data.RefreshToken),
		Expiry:       guessExpiry(deref(data.AccessToken)),
	}
	return withExtras(tok, deref(data.UserId)), httpResponse, nil
}

func withExtras(tok *oauth2.Token, userID int) *oauth2.Token {
	return tok.WithExtra(map[string]interface{}{
		"sense_user_id": userID,
	})
}

func fromExtras(tok *oauth2.Token) (userID *int) {
	if tok == nil {
		return nil
	}
	userID, _ = tok.Extra("sense_user_id").(*int)
	return
}

// TokenSource returns a token source that will renew the token when needed.
// The provided token must have been generated by the PasswordCredentialsToken method.
// Renewals will use the HTTP client configured in the Config in a background context.
func (c Config) TokenSource(tok *oauth2.Token) *TokenSource {
	return &TokenSource{
		tok:        tok,
		baseUrl:    c.BaseURL,
		httpClient: c.http(),
		client:     c.getClient(),
	}
}

// Client returns an HTTP client that will renew the token when needed.
// The provided token must have been generated by the PasswordCredentialsToken method.
// Renewals will use the HTTP client configured in the Config in a background context.
// The returned client will otherwise be unassociated with the HTTP client in the Config.
func (c Config) Client(tok *oauth2.Token) *http.Client {
	return c.ClientFrom(nil, tok)
}

// ClientFrom returns an HTTP client, derived from the provided client, that will
// renew the token when needed.
// The provided token must have been generated by the PasswordCredentialsToken method.
// Renewals will use the HTTP client configured in the Config in a background context.
func (c Config) ClientFrom(cli *http.Client, tok *oauth2.Token) *http.Client {
	return NewClientFrom(cli, c.TokenSource(tok))
}

// NewClientFrom returns an HTTP client, derived from the provided client, that
// will renew the token when needed. Renewals will use the provided token source
// in a background context.
func NewClientFrom(cli *http.Client, ts oauth2.TokenSource) *http.Client {
	ctx := context.Background()
	if cli != nil {
		// This is how oauth2.NewClient gets the underlying Transport wrapped.
		ctx = context.WithValue(ctx, oauth2.HTTPClient, cli)
	}
	return oauth2.NewClient(ctx, ts)
}

// TokenSource is a token source that will renew the token when needed.
// It can be used with standard golang.org/x/oauth2 types.
type TokenSource struct {
	tok        *oauth2.Token
	baseUrl    string
	httpClient *http.Client
	client     internalClient
}

var _ oauth2.TokenSource = (*TokenSource)(nil)

// Token returns a token that will renew the token when needed using the
// background context.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	return t.TokenContext(context.Background())
}

func ptr[T any](v T) *T {
	return &v
}

// TokenContext returns a token, renewing it if needed using
// the provided context.
func (t *TokenSource) TokenContext(ctx context.Context) (*oauth2.Token, error) {
	// Note: The standard oauth2 types will simply call Token(), without passing
	// a context.  This capabililty is nevertheless here for completeness and
	// future hopes.  We could implement our own HTTP transport that passes
	// the request context here, but that's unnecessary work right now.
	//
	// See also: https://github.com/golang/oauth2/issues/262
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	if t.tok != nil && t.tok.Valid() {
		debugf("senseauth: re-using existing valid token for %s:%d", file, line)
		return t.tok, nil
	}
	userID := fromExtras(t.tok)
	if userID == nil || *userID == 0 {
		// The Sense OAuth implementation seems to require a user ID to accompany
		// token renewals.  Since we don't want to connect this package to the main
		// sense package (where this data value would live), we smuggle it in the
		// token's extra data.  If it's not there, we can't renew the token.
		panic("senseauth: token was not generated correctly")
	}
	req := client.RenewAuthTokenFormdataRequestBody{
		UserId:        userID,
		RefreshToken:  &t.tok.RefreshToken,
		IsAccessToken: ptr(true),
	}
	res, err1 := t.client.RenewAuthTokenWithFormdataBodyWithResponse(ctx, req)
	if err := client.Ensure(err1, "RenewAuthToken", res, 200); err != nil {
		return nil, fmt.Errorf("senseauth: %w", err)
	}
	data := res.JSON200

	debugf("senseauth: successfully renewed token for %s:%d", file, line)
	tok := &oauth2.Token{
		AccessToken:  deref(data.AccessToken),
		RefreshToken: deref(data.RefreshToken),
		Expiry:       deref(data.Expires),
	}
	return withExtras(tok, *userID), nil
}

// This seems to be roughly when tokens seem to expire
const assumeExpire = 8 * time.Hour

// The /authenticate endpoint doesn't return a token expiry, so we have to guess.
// Fortunately it looks like this appears to be encoded in the token itself.
func guessExpiry(tok string) time.Time {
	jwtExpFieldRe := regexp.MustCompile(`"exp":(\d+),`)
	if strings.HasPrefix(tok, "t1.v2.") {
		parts := strings.SplitN(tok, ".", 5)
		if len(parts) == 5 {
			// This field appears to be a truncated base64-encoded JSON (JWT) object,
			// so we'll likely get an ErrUnexpectedEOF. Fortunately the bit we're
			// interested in is in the first part, so we can ignore the error.
			dr := base64.NewDecoder(base64.StdEncoding, strings.NewReader(parts[3]))
			data, err := io.ReadAll(dr)
			if err == nil || err == io.ErrUnexpectedEOF {
				if m := jwtExpFieldRe.FindStringSubmatch(string(data)); m != nil {
					secs, err := strconv.ParseInt(m[1], 10, 64)
					if err == nil {
						t := time.Unix(secs, 0)
						debug("senseauth: using expiry from token:", t)
						return t
					}
				}
			}
		}
	}
	debug("senseauth: assuming expiry: ", assumeExpire)
	return time.Now().Add(assumeExpire)
}
