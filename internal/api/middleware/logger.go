package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logger returns a middleware that logs requests with latency
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Calculate latency
		latency := time.Since(start)
		
		// Get status code
		status := c.Response().StatusCode()
		
		// Log format: METHOD PATH STATUS LATENCY
		log.Printf("%s %s %d %v",
			c.Method(),
			c.Path(),
			status,
			latency,
		)
		
		// Set latency header for clients
		c.Set("X-Response-Time", latency.String())
		
		return err
	}
}

// LoggerWithConfig returns a configurable logger middleware
type LoggerConfig struct {
	// Skip defines a function to skip logging for certain paths
	Skip func(c *fiber.Ctx) bool
	// Format defines log format (not implemented, using default)
	Format string
	// TimeFormat defines time format
	TimeFormat string
}

// LoggerWithConfig returns a middleware with custom config
func LoggerWithConfig(config LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if we should skip logging
		if config.Skip != nil && config.Skip(c) {
			return c.Next()
		}
		
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Calculate latency
		latency := time.Since(start)
		
		// Get response info
		status := c.Response().StatusCode()
		
		// Log with timestamp
		timeFormat := config.TimeFormat
		if timeFormat == "" {
			timeFormat = "2006-01-02 15:04:05"
		}
		
		log.Printf("[%s] %s %s %d %v %s",
			time.Now().Format(timeFormat),
			c.Method(),
			c.Path(),
			status,
			latency,
			c.IP(),
		)
		
		c.Set("X-Response-Time", latency.String())
		
		return err
	}
}
