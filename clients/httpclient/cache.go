package httpclient

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// IRedisClient defines the interface for a Redis client used by the cache middleware.
// It must implement Get and Set methods for string keys and values.
type IRedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

// cacheKeyHeaders is a list of HTTP header names used to compose the cache key.
type cacheKeyHeaders []string

// CacheConfig holds the configuration for the cache middleware, including Redis client, TTL, and headers for cache key.
type CacheConfig struct {
	RedisClient IRedisClient
	TTL         time.Duration
	OverrideTTL bool
	Headers     cacheKeyHeaders
}

// SerializableCache represents the structure of a cached HTTP response, ready for (de)serialization.
type SerializableCache struct {
	Status            string              `json:"status"`
	StatusCode        int                 `json:"status_code"`
	Proto             string              `json:"proto"`
	ResponseHeaders   map[string][]string `json:"header"`
	Body              string              `json:"body"`
	CacheControlValue int                 `json:"cacheControlValue"`
	Policy            CachePolicy         `json:"policy"`
}

// CachePolicy defines cache control policy for a cached response, including max-age and headers used.
type CachePolicy struct {
	MaxAge  int      `json:"maxAge"`
	Headers []string `json:"headers"`
}

// NewCacheMiddleware is an HTTP middleware that provides transparent caching for GET requests using a Redis backend.
//
// It checks if the cache is enabled and a Redis client is configured. For each GET request, it attempts to retrieve
// a cached response from Redis using a generated cache key. If a valid cached response is found, it is deserialized
// and returned immediately, setting the "X-Cache" header to "HIT". If not found, the request proceeds to the next
// RoundTripper, and the response is cached asynchronously if the status code is 2xx. The cache TTL can be overridden
// by configuration, and the middleware also updates the "Cache-Control" header accordingly.
//
// Parameters:
//
//	cfg *CacheConfig: Cache configuration struct.
//	  - RedisClient: Redis client used to store and retrieve cached data.
//	  - TTL: Default expiration time (Time To Live) for cache entries.
//	  - OverrideTTL: If true, overrides the TTL from the Cache-Control header with the configured TTL.
//	  - Headers: HTTP headers that will be considered when generating the cache key.
//
// Returns:
//
//	A function that wraps an http.RoundTripper with caching logic.
func NewCacheMiddleware(cfg *CacheConfig) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if cfg.RedisClient == nil {
				return next.RoundTrip(req)
			}

			if req.Method != "GET" {
				return next.RoundTrip(req)
			}

			cacheKey := getCacheKey(req, cfg.Headers)

			value, err := cfg.RedisClient.Get(req.Context(), cacheKey)

			if err == nil && value != "" {
				responseSerialized, err := parseCachedResponseFromString(value)

				if err != nil {
					logger.Error().Msg("Error deserializing cached response")
					return next.RoundTrip(req)
				}

				resp := &http.Response{
					StatusCode:    responseSerialized.StatusCode,
					Status:        responseSerialized.Status,
					Proto:         responseSerialized.Proto,
					ProtoMajor:    1,
					ProtoMinor:    1,
					Body:          io.NopCloser(strings.NewReader(responseSerialized.Body)),
					Header:        make(http.Header),
					ContentLength: int64(len(responseSerialized.Body)),
					Request:       req,
				}

				for k, v := range responseSerialized.ResponseHeaders {
					for _, vv := range v {
						resp.Header.Add(k, vv)
					}
				}

				newCacheControl := fmt.Sprintf("max-age=%v, public", responseSerialized.CacheControlValue)
				resp.Header.Set("Cache-Control", newCacheControl)
				resp.Header.Set("X-Cache", "HIT")

				return resp, nil
			}

			resp, err := next.RoundTrip(req)

			if err != nil {
				return resp, fmt.Errorf("error executing request: %w", err)
			}

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {

				responseCacheControl := getCacheControlHeaderValue(resp)

				var ttl time.Duration = time.Second * time.Duration(responseCacheControl)

				if cfg.OverrideTTL {
					ttl = cfg.TTL
				}

				newCacheControl := fmt.Sprintf("max-age=%v, public", ttl.Seconds())
				resp.Header.Set("Cache-Control", newCacheControl)

				policy := CachePolicy{
					MaxAge:  responseCacheControl,
					Headers: cfg.Headers,
				}

				cachedValue, err := responseToJSON(resp, policy)

				resp.Header.Set("X-Cache", "MISS")

				if err != nil {
					logger.Err(err).Msg("Error serializing response for cache")
					return resp, fmt.Errorf("error serializing response for cache: %w", err)
				}

				go func() {
					setErr := cfg.RedisClient.Set(req.Context(), cacheKey, cachedValue, ttl)

					if setErr != nil {
						logger.Error().Err(setErr).Msg("Error saving to cache")
					}
				}()

			}

			return resp, nil
		})
	}
}

func getCacheKey(req *http.Request, headers cacheKeyHeaders) string {
	keyParts := []string{
		buildURLPart(req),
		buildQueryPart(req),
		buildVaryHeadersPart(req, headers),
	}

	base := strings.Join(keyParts, "|")
	hash := sha256.Sum256([]byte(base))
	return hex.EncodeToString(hash[:])
}

func buildURLPart(req *http.Request) string {
	return req.URL.String()
}

func buildQueryPart(req *http.Request) string {
	query := req.URL.Query()
	var queryParts []string

	for k, v := range query {
		sort.Strings(v)
		queryParts = append(queryParts, k+"="+strings.Join(v, ","))
	}

	sort.Strings(queryParts)
	return strings.Join(queryParts, "&")
}

func buildVaryHeadersPart(req *http.Request, headers cacheKeyHeaders) string {
	var headersParts []string

	for _, key := range headers {
		if req.Header.Get(key) != "" {
			headersParts = append(headersParts, key+":"+req.Header.Get(key))
		}
	}

	fmt.Println("Vary Headers:", headersParts)
	return strings.Join(headersParts, "|")
}

func responseToJSON(resp *http.Response, policy CachePolicy) ([]byte, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	sr := SerializableCache{
		Status:            resp.Status,
		StatusCode:        resp.StatusCode,
		Proto:             resp.Proto,
		ResponseHeaders:   resp.Header,
		Policy:            policy,
		CacheControlValue: getCacheControlHeaderValue(resp),
		Body:              string(bodyBytes),
	}

	return json.Marshal(sr)
}

func parseCachedResponseFromString(jsonStr string) (*SerializableCache, error) {
	var sc SerializableCache

	err := json.Unmarshal([]byte(jsonStr), &sc)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached response: %w", err)
	}

	return &sc, nil
}

func getCacheControlHeaderValue(res *http.Response) int {
	cacheControlValue := res.Header.Get("Cache-Control")
	re := regexp.MustCompile(`max-age=(\d+)`)
	matches := re.FindStringSubmatch(cacheControlValue)

	if len(matches) > 1 {
		age, err := strconv.Atoi(matches[1])

		if err != nil {
			fmt.Println("error on convert to int", err)
			return 0
		}

		return age
	}

	return 0
}
