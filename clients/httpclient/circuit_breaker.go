package httpclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

// NewCircuitBreaker wraps an http.RoundTripper with a circuit breaker using gobreaker.
//
// The circuit breaker monitors HTTP requests and opens the circuit when the error rate
// reaches a threshold (default: 50% errors out of at least 20 requests, considering status >= 500 or 429 as errors).
// While open, requests will fail fast without calling the underlying transport. After a short interval,
// a limited number of requests are allowed to test recovery. If successful, the circuit closes again.
//
// Parameters:
//
//	cfg: Configuration for the circuit breaker.
//	     - cfg.Enabled: activates/deactivates the breaker.
//	     - cfg.Name: identifies the breaker instance (useful for logging/metrics).
//	next: The next http.RoundTripper to be wrapped. This is usually http.DefaultTransport or a custom transport.
//
// Returns:
//
//	An http.RoundTripper that applies circuit breaker logic to all requests.
func NewCircuitBreakerMiddleware(name string) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {

			breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:        name,
				MaxRequests: 10,
				Interval:    time.Second * 10,

				ReadyToTrip: func(counts gobreaker.Counts) bool {
					total := counts.Requests
					failures := counts.TotalFailures
					return total >= 20 && failures*100/total >= 50
				},

				IsSuccessful: func(err error) bool {
					if err == nil {
						return true
					}

					if httpErr, ok := err.(*HTTPStatusError); ok {
						return httpErr.Status < 500 && httpErr.Status != 429
					}

					return false
				},
			})

			logState(name, breaker, req)

			result, err := breaker.Execute(func() (any, error) {
				resp, err := next.RoundTrip(req)
				if err != nil {
					return nil, err
				}

				if resp.StatusCode >= 500 || resp.StatusCode == 429 {
					return nil, &HTTPStatusError{Status: resp.StatusCode, Err: fmt.Errorf("HTTP error")}
				}

				return resp, nil
			})

			if err != nil {
				return nil, err
			}

			return result.(*http.Response), nil
		})
	}
}

type HTTPStatusError struct {
	Status int
	Err    error
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("HTTP status %d: %v", e.Status, e.Err)
}

func logState(name string, breaker *gobreaker.CircuitBreaker, req *http.Request) {
	state := breaker.State()
	if state != gobreaker.StateClosed {
		var stateStr string
		switch state {
		case gobreaker.StateOpen:
			stateStr = "OPEN"
		case gobreaker.StateHalfOpen:
			stateStr = "HALF-OPEN"
		default:
			stateStr = "UNKNOWN"
		}

		logger.Info().
			Str("cb", name).
			Str("url", req.URL.String()).
			Str("state", stateStr).
			Msg("circuit-breaker:state change")
	}
}
