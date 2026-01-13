package response

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

// Response represents a standardized API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta contains metadata for paginated responses
type Meta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Total      int    `json:"total,omitempty"`
	CacheHit   bool   `json:"cache_hit,omitempty"`
	LatencyMs  int64  `json:"latency_ms,omitempty"`
}

// Pre-allocated byte slices for common responses
var (
	successPrefix = []byte(`{"success":true,"data":`)
	errorPrefix   = []byte(`{"success":false,"error":`)
	timestampKey  = []byte(`,"timestamp":`)
	closeBrace    = []byte(`}`)
)

// Success sends a successful response with data
func Success(c *fiber.Ctx, data interface{}) error {
	return SuccessWithMeta(c, data, nil)
}

// SuccessWithMeta sends a successful response with data and metadata
func SuccessWithMeta(c *fiber.Ctx, data interface{}, meta *Meta) error {
	resp := Response{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().UnixMilli(),
	}
	
	// Use sonic for faster JSON encoding
	body, err := sonic.Marshal(resp)
	if err != nil {
		return Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Failed to encode response", err.Error())
	}
	
	c.Set("Content-Type", "application/json")
	return c.Send(body)
}

// Error sends an error response
func Error(c *fiber.Ctx, status int, code, message, details string) error {
	resp := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UnixMilli(),
	}
	
	body, _ := sonic.Marshal(resp)
	c.Set("Content-Type", "application/json")
	return c.Status(status).Send(body)
}

// BadRequest sends a 400 error response
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, "BAD_REQUEST", message, "")
}

// NotFound sends a 404 error response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", message, "")
}

// Unauthorized sends a 401 error response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, "")
}

// InternalError sends a 500 error response
func InternalError(c *fiber.Ctx, err error) error {
	return Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", err.Error())
}

// TooManyRequests sends a 429 error response
func TooManyRequests(c *fiber.Ctx) error {
	return Error(c, fiber.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", "Please slow down")
}

// Raw sends raw JSON bytes directly (zero-copy for cached responses)
func Raw(c *fiber.Ctx, body []byte) error {
	c.Set("Content-Type", "application/json")
	return c.Send(body)
}

// RawWithCacheHeader sends raw JSON with cache indicator
func RawWithCacheHeader(c *fiber.Ctx, body []byte, cacheHit bool) error {
	c.Set("Content-Type", "application/json")
	if cacheHit {
		c.Set("X-Cache", "HIT")
	} else {
		c.Set("X-Cache", "MISS")
	}
	return c.Send(body)
}
