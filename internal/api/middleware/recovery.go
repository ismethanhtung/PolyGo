package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/polygo/pkg/response"
)

// Recovery returns a middleware that recovers from panics
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic and stack trace
				log.Printf("PANIC RECOVERED: %v\n%s", r, debug.Stack())
				
				// Return 500 error
				response.Error(c, fiber.StatusInternalServerError, 
					"INTERNAL_ERROR", 
					"An unexpected error occurred", 
					"")
			}
		}()
		
		return c.Next()
	}
}

// RecoveryWithConfig returns a recovery middleware with custom handler
type RecoveryConfig struct {
	// EnableStackTrace enables logging stack trace
	EnableStackTrace bool
	// StackTraceHandler handles stack trace
	StackTraceHandler func(c *fiber.Ctx, e interface{})
}

// RecoveryWithConfig returns a middleware with custom config
func RecoveryWithConfig(config RecoveryConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				if config.EnableStackTrace {
					log.Printf("PANIC RECOVERED: %v\n%s", r, debug.Stack())
				} else {
					log.Printf("PANIC RECOVERED: %v", r)
				}
				
				if config.StackTraceHandler != nil {
					config.StackTraceHandler(c, r)
				}
				
				response.Error(c, fiber.StatusInternalServerError,
					"INTERNAL_ERROR",
					"An unexpected error occurred",
					"")
			}
		}()
		
		return c.Next()
	}
}
