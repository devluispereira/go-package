package server

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
}

// NewServer creates and configures a Fiber server instance.
//
// Parameters:
//
//	name: The name of the origin application. Used for the X-Origin-App header.
//	forwardHeaders: List of headers to be forwarded. If empty, uses defaults.
//
// Behavior:
//   - Removes default server identification headers.
//   - Sets the X-Origin-App header in the request.
//   - Applies ForwardHeadersMiddleware to collect and forward headers.
//   - Adds a /healthcheck endpoint for health monitoring.
//
// Usage:
//
//	server := NewServer("my-app", []string{"x-request-id", "x-client-user-agent"})
//	log.Fatal(server.App.Listen(":8080"))
func NewServer(name string, forwardHeaders []string) *Server {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Response().Header.Del("Server")
		c.Response().Header.Del("X-Powered-By")
		c.Set("X-Origin-App", name)

		return c.Next()
	})

	app.Use(ForwardHeadersMiddleware(name, forwardHeaders))

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("OK")
	})

	return &Server{
		App: app,
	}
}
