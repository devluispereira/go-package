package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// defaultForwardHeaders defines the default list of headers to be forwarded by the middleware.
// These headers are commonly used for tracing, user identification, and platform information.
var defaultForwardHeaders = []string{
	"x-request-id",
	"x-canonical-uri",
	"x-client-user-agent",
	"x-platform-id",
	"x-device-id",
	"x-client-version",
	"x-tenant-id",
	"x-profile-type",
	"x-country-code",
	"x-origin-ip",
	"x-hsid",
	"x-glb-exp-id",
}

type ForwardedHeadersKeyType struct{}

// ForwardHeadersMiddleware collects specified headers from the incoming request and stores them in Fiber's Locals.
//
// Parameters:
//
//	appName: Name of the origin application. Will be added to the forwarded headers as "x-origin-app".
//	forwardHeaders: List of header names to forward. If empty, uses defaultForwardHeaders.
//
// Behavior:
//   - For each header in the list, if present in the request, adds it to a map.
//   - Adds "x-origin-app" with the value of appName to the map.
//   - Stores the map in c.Locals("forwardedHeaders") for use in subsequent handlers.
//
// Usage:
//
//	app.Use(ForwardHeadersMiddleware("my-app", []string{"x-request-id", "x-client-user-agent"}))
//
//	// To access forwarded headers in a handler:
//	headers := c.Locals("forwardedHeaders").(map[string]string)

func ForwardHeadersMiddleware(appName string, forwardHeaders []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headersMap := make(map[string]string)

		if len(forwardHeaders) == 0 {
			forwardHeaders = append(forwardHeaders, defaultForwardHeaders...)
		}

		for _, h := range forwardHeaders {
			val := c.Get(h)
			if val != "" {
				headersMap[h] = val
			}
		}

		ctx := context.WithValue(c.UserContext(), "forwardedHeaders", headersMap)

		c.SetUserContext(ctx)
		return c.Next()
	}
}
