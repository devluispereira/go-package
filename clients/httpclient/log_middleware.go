package httpclient

import (
	"net/http"
	"time"
)

// NewLoggingMiddleware returns an HTTP middleware that logs all outgoing requests and responses.
//
// Parameters:
//   name: The name of the service or component making the HTTP request (used for log context).
//
// Returns:
//   A function that wraps an http.RoundTripper and logs request and response details, including method, URL, status, duration, cache status, and errors.
//   Logs at INFO level for successful requests and ERROR level for failed requests.

func NewLoggingMiddleware(name string) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			start := time.Now()
			resp, err := next.RoundTrip(req)
			duration := time.Since(start)

			if err != nil {
				logger.Error().
					Str("service", name).
					Str("method", req.Method).
					Str("url", req.URL.String()).
					Int("status", 500).
					Int64("duration_ms", duration.Milliseconds()).
					Msg(err.Error())

				return resp, err
			}

			logger.Info().
				Str("service", name).
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Int("status", resp.StatusCode).
				Int64("duration_ms", duration.Milliseconds()).
				Str("cache", resp.Header.Get("X-Cache")).
				Msg(resp.Status)

			return resp, err
		})
	}
}
