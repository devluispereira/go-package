package server

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CacheType defines allowed values for the Cache-Control header type.
// Use the provided constants for safety and validation.
type CacheType string

const (
	CachePublic  CacheType = "public"
	CachePrivate CacheType = "private"
	CacheNoStore CacheType = "no-store"
	CacheNoCache CacheType = "no-cache"
)

// SetCacheControlMiddleware sets the Cache-Control header for a route or group in Fiber.
//
// Parameters:
//
//	cacheType: One of the allowed CacheType values (public, private, no-store, no-cache).
//	ttl: Time to live in seconds for the cache (max-age). If <= 0, max-age is omitted.
//
// Usage:
//
//	app.Get("/route", SetCacheControlMiddleware(CachePublic, 60), handler)
//
// If an invalid cacheType is provided, the middleware returns an error and does not set the header.
func SetCacheControlMiddleware(cacheType CacheType, ttl int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isValidCacheType(cacheType) {
			return fmt.Errorf("invalid cache type: %s", cacheType)
		}

		err := c.Next()

		if err != nil {
			return err
		}

		value := string(cacheType)
		if ttl > 0 {
			value += ", max-age=" + strconv.Itoa(ttl)
		}
		c.Response().Header.Set("Cache-Control", value)
		return nil
	}
}

func isValidCacheType(ct CacheType) bool {
	switch ct {
	case CachePublic, CachePrivate, CacheNoStore, CacheNoCache:
		return true
	default:
		return false
	}
}
