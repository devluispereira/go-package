package httpclient

import (
	"net/http"
)

// NewHeaderMiddleware returns an HTTP middleware that adds custom headers to all outgoing requests.
//
// Parameters:
//
//	headers: A map of header keys and values to be set on each outgoing HTTP request.
//
// Returns:
//
//	A function that wraps an http.RoundTripper and sets the provided headers on every request before forwarding it.
//	Existing headers with the same key will be overwritten.
func NewHeaderMiddleware(headers map[string]string) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if len(headers) != 0 {
				for k, v := range headers {
					req.Header.Set(k, v)
				}
			}

			resp, err := next.RoundTrip(req)
			return resp, err
		})
	}
}
