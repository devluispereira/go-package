package httpclient

import (
	"net/http"
)

// RoundTripperMiddleware defines a function that wraps an http.RoundTripper with additional behavior (middleware).
type RoundTripperMiddleware func(http.RoundTripper) http.RoundTripper

// RoundTripperFunc allows using ordinary functions as http.RoundTripper implementations.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls the underlying function, allowing RoundTripperFunc to satisfy http.RoundTripper.
func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// configMiddlewares composes a slice of RoundTripperMiddleware into a single http.RoundTripper chain.
// The first middleware in the slice will be the outermost (executed first).
func configMiddlewares(middlewares []RoundTripperMiddleware) http.RoundTripper {
	composed := http.DefaultTransport

	for i := len(middlewares) - 1; i >= 0; i-- {
		composed = middlewares[i](composed)
	}

	return composed
}
