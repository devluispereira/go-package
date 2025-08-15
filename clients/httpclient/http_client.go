package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

type HTTPResponse struct {
	Body       any
	StatusCode int
	Headers    http.Header
}

// NewHTTPClient creates a new HTTPClient instance.
//
// Params:
//
//   - baseUrl: Base URL for requests (used if path is relative).
//
//   - timeout: Timeout for HTTP requests.
//
//   - middlewares: Optional RoundTripper middlewares.
//
//     Recommended order:
//
//     1. NewLoggingMiddleware;
//     (Should be outermost to log all requests and responses, including cache hits and circuit breaker events)
//
//     2. NewHeaderMiddleware;
//     (Sets custom headers before cache and circuit logic, ensuring cache keys and backend requests use the correct headers)
//
//     3. CacheMiddleware;
//     (Checks/sets cache after headers are set, and before circuit breaker, for maximum cache efficiency)
//
//     4. CircuitBreakerMiddleware.
//     (Protects backend only for requests that reach it, after cache and header logic)
//
// Returns: Configured HTTP client.
func NewHTTPClient(
	baseUrl string,
	timeout time.Duration,
	middlewares ...RoundTripperMiddleware) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout:   timeout,
			Transport: configMiddlewares(middlewares),
		},
		baseURL: baseUrl,
	}
}

// Get sends an HTTP GET request to the specified path.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - path: Request path or full URL.
//
// Returns:
//   - *HTTPResponse: The response object.
//   - error: Any error encountered.
func (c *HTTPClient) Get(ctx context.Context, path string) (*HTTPResponse, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// Post sends an HTTP POST request to the specified path with a request body.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - path: Request path or full URL.
//   - body: Request body as io.Reader.
//
// Returns:
//   - *HTTPResponse: The response object.
//   - error: Any error encountered.
func (c *HTTPClient) Post(ctx context.Context, path string, body io.Reader) (*HTTPResponse, error) {
	return c.doRequest(ctx, "POST", path, body)
}

// Put sends an HTTP PUT request to the specified path with a request body.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - path: Request path or full URL.
//   - body: Request body as io.Reader.
//
// Returns:
//   - *HTTPResponse: The response object.
//   - error: Any error encountered.
func (c *HTTPClient) Put(ctx context.Context, path string, body io.Reader) (*HTTPResponse, error) {
	return c.doRequest(ctx, "PUT", path, body)
}

// Patch sends an HTTP PATCH request to the specified path with a request body.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - path: Request path or full URL.
//   - body: Request body as io.Reader.
//
// Returns:
//   - *HTTPResponse: The response object.
//   - error: Any error encountered.
func (c *HTTPClient) Patch(ctx context.Context, path string, body io.Reader) (*HTTPResponse, error) {
	return c.doRequest(ctx, "PATCH", path, body)
}

// Delete sends an HTTP DELETE request to the specified path.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - path: Request path or full URL.
//
// Returns:
//   - *HTTPResponse: The response object.
//   - error: Any error encountered.
func (c *HTTPClient) Delete(ctx context.Context, path string) (*HTTPResponse, error) {
	return c.doRequest(ctx, "DELETE", path, nil)
}

// Head sends an HTTP HEAD request to the specified path.
//
// Parameters:
//   - ctx: Context for cancellation and timeout.
//   - path: Request path or full URL.
//
// Returns:
//   - *HTTPResponse: The response object.
//   - error: Any error encountered.
func (c *HTTPClient) Head(ctx context.Context, path string) (*HTTPResponse, error) {
	return c.doRequest(ctx, "HEAD", path, nil)
}

func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body io.Reader) (*HTTPResponse, error) {
	url := path
	if !strings.HasPrefix(path, "http") {
		url = strings.TrimSuffix(c.baseURL, "/") + "/" + strings.TrimPrefix(path, "/")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	forwardedHeaders := getForwardedHeaders(ctx)

	for k, value := range forwardedHeaders {
		req.Header.Set(k, value)
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	defer resp.Body.Close()
	var jsonBody any
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &HTTPResponse{
		Body:       jsonBody,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}

func getForwardedHeaders(ctx context.Context) map[string]string {
	headers, _ := ctx.Value("forwardedHeaders").(map[string]string)
	return headers
}
