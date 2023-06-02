package client

import (
	"errors"
	"fmt"
	"reflect"
)

// The bulk of the client implementation is generated from the OpenAPI spec
// using oapi-codegen.  https://github.com/deepmap/oapi-codegen

//go:generate oapi-codegen --config=oapi-codegen.config openapi.yaml

// AsErr is a helper that returns nil only if parent is nil and withErrorReason
// does not have a value in its ErrorReason field.  Otherwise, it returns an
// error derived from both of these values as appropriate.
//
// This is useful since JSON APIs have multiple ways they can return errors,
// and this function helps normalize them all behind an error type.
func AsErr(parent error, withErrorReason interface{}) (err error) {
	err = parent
	if withErrorReason == nil {
		return err
	}

	// Try to extract an error message from the provided value.  Multiple
	// structs may have an ErrorReason field, so we use reflect to find it.
	v := reflect.ValueOf(withErrorReason)
	// Dereference if needed
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName("ErrorReason")
	if field.IsValid() && !field.IsNil() {
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if val, ok := field.Interface().(string); ok {
			err = errors.New(val)
			if parent != nil {
				err = fmt.Errorf("%w: %w", parent, err)
			}
		}
	}

	return err
}

type statusCoder interface {
	Status() string
	StatusCode() int
}

// ErrAuthenticationNeeded is wrapped by any errors resulting from a 401 response.
// You can test it using errors.Is(err, client.ErrAuthenticationNeeded).
var ErrAuthenticationNeeded = errors.New("authentication needed")

// For reference, a typical Response struct in the generated API client for a call
// that specifies 200 and 401 are valid responses normally looks like this:
//
//	type AuthenticateResponse struct {
//		Body         []byte
//		HTTPResponse *http.Response
//		JSON200      *Hello
//		JSON401      *struct {
//			ErrorReason *string `json:"error_reason,omitempty"`
//			MfaToken    *string `json:"mfa_token,omitempty"`
//			Status      *string `json:"status,omitempty"`
//		}
//		JSONDefault *Error
//	}
func Ensure(parent error, prepend string, response statusCoder, code int) error {
	if parent != nil && prepend != "" {
		parent = fmt.Errorf("%s: %w", prepend, parent)
	}
	if response == nil {
		// we can't do any better than parent, so just return it
		return AsErr(parent, nil)
	}
	v := reflect.ValueOf(response)
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// response pointed to nil, so, like above, we can't do any better
			return AsErr(parent, nil)
		}
		v = v.Elem()
	}

	if response.StatusCode() == code {
		// this means success, but don't ignore the error we were passed
		return parent
	}

	// Special treatment for 401 responses, which seem to also use an *Error body.
	if response.StatusCode() == 401 {
		if parent == nil {
			parent = ErrAuthenticationNeeded
		} else {
			parent = fmt.Errorf("%w: %w", parent, ErrAuthenticationNeeded)
		}
		field := v.FieldByName("JSON401")
		if field.IsValid() && !field.IsNil() {
			if val, ok := field.Interface().(*Error); ok {
				return AsErr(parent, val)
			}
		}
	}

	if parent == nil {
		parent = errors.New(response.Status())
	} else {
		parent = fmt.Errorf("%w: %s", parent, response.Status())
	}

	// Code isn't what we expected, it's not a 401, and there's no JSON{code} field.
	// Return an error derived from the JSONDefault field if it is available.
	field := v.FieldByName("JSONDefault")
	if field.IsValid() && !field.IsNil() {
		if val, ok := field.Interface().(*Error); ok {
			return AsErr(parent, val)
		}
	}
	return AsErr(parent, nil)
}

// AuthPayloadKey is something we use to smuggle the auth response payload out via the auth token.
var AuthPayloadKey = struct{}{}
