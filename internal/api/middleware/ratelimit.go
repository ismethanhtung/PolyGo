package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/polygo/pkg/response"
)

// RateLimitConfig holds rate limiter configuration
type RateLimitConfig struct {
	// Max requests per window
	Max int
	// Window duration
	Window time.Duration
	// Key generator function
	KeyGenerator func(c *fiber.Ctx) string
	// Skip function
	Skip func(c *fiber.Ctx) bool
}

// rateLimitEntry holds rate limit state for a key
type rateLimitEntry struct {
	count     int
	resetAt   time.Time
	mu        sync.Mutex
}

// rateLimiter holds all rate limit entries
type rateLimiter struct {
	entries map[string]*rateLimitEntry
	mu      sync.RWMutex
	config  RateLimitConfig
}

// newRateLimiter creates a new rate limiter
func newRateLimiter(config RateLimitConfig) *rateLimiter {
	rl := &rateLimiter{
		entries: make(map[string]*rateLimitEntry),
		config:  config,
	}
	
	// Start cleanup goroutine
	go rl.cleanup()
	
	return rl
}

// cleanup removes expired entries
func (r *rateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for key, entry := range r.entries {
			if now.After(entry.resetAt) {
				delete(r.entries, key)
			}
		}
		r.mu.Unlock()
	}
}

// check checks if request is allowed
func (r *rateLimiter) check(key string) (bool, int, time.Time) {
	r.mu.RLock()
	entry, exists := r.entries[key]
	r.mu.RUnlock()
	
	now := time.Now()
	
	if !exists {
		r.mu.Lock()
		entry = &rateLimitEntry{
			count:   1,
			resetAt: now.Add(r.config.Window),
		}
		r.entries[key] = entry
		r.mu.Unlock()
		return true, r.config.Max - 1, entry.resetAt
	}
	
	entry.mu.Lock()
	defer entry.mu.Unlock()
	
	// Reset if window expired
	if now.After(entry.resetAt) {
		entry.count = 1
		entry.resetAt = now.Add(r.config.Window)
		return true, r.config.Max - 1, entry.resetAt
	}
	
	// Check limit
	if entry.count >= r.config.Max {
		return false, 0, entry.resetAt
	}
	
	entry.count++
	return true, r.config.Max - entry.count, entry.resetAt
}

// RateLimit returns a rate limiting middleware
func RateLimit(config RateLimitConfig) fiber.Handler {
	if config.Max == 0 {
		config.Max = 100
	}
	if config.Window == 0 {
		config.Window = time.Minute
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}
	
	limiter := newRateLimiter(config)
	
	return func(c *fiber.Ctx) error {
		// Check skip
		if config.Skip != nil && config.Skip(c) {
			return c.Next()
		}
		
		key := config.KeyGenerator(c)
		allowed, remaining, resetAt := limiter.check(key)
		
		// Set headers
		c.Set("X-RateLimit-Limit", string(rune(config.Max)))
		c.Set("X-RateLimit-Remaining", string(rune(remaining)))
		c.Set("X-RateLimit-Reset", resetAt.Format(time.RFC3339))
		
		if !allowed {
			c.Set("Retry-After", resetAt.Sub(time.Now()).String())
			return response.TooManyRequests(c)
		}
		
		return c.Next()
	}
}

// DefaultRateLimit returns a rate limiter with default settings
func DefaultRateLimit() fiber.Handler {
	return RateLimit(RateLimitConfig{
		Max:    1000,
		Window: 10 * time.Second,
	})
}
