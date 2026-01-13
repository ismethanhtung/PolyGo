package unit

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/polygo/pkg/response"
)

func TestResponse_Success(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Success(c, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result response.Response
	err = sonic.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.Nil(t, result.Error)
	assert.NotZero(t, result.Timestamp)
}

func TestResponse_Error(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Error(c, 400, "BAD_REQUEST", "Invalid input", "Details here")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result response.Response
	err = sonic.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.False(t, result.Success)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "BAD_REQUEST", result.Error.Code)
	assert.Equal(t, "Invalid input", result.Error.Message)
	assert.Equal(t, "Details here", result.Error.Details)
}

func TestResponse_BadRequest(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.BadRequest(c, "Missing parameter")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestResponse_NotFound(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.NotFound(c, "Resource not found")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestResponse_Unauthorized(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Unauthorized(c, "Invalid credentials")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestResponse_TooManyRequests(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.TooManyRequests(c)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)
}

func TestResponse_Raw(t *testing.T) {
	app := fiber.New()

	rawData := []byte(`{"raw": "data"}`)

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Raw(c, rawData)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, rawData, body)
}

func TestResponse_RawWithCacheHeader(t *testing.T) {
	app := fiber.New()

	rawData := []byte(`{"cached": "data"}`)

	app.Get("/hit", func(c *fiber.Ctx) error {
		return response.RawWithCacheHeader(c, rawData, true)
	})

	app.Get("/miss", func(c *fiber.Ctx) error {
		return response.RawWithCacheHeader(c, rawData, false)
	})

	// Test cache hit
	req := httptest.NewRequest("GET", "/hit", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, "HIT", resp.Header.Get("X-Cache"))

	// Test cache miss
	req = httptest.NewRequest("GET", "/miss", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, "MISS", resp.Header.Get("X-Cache"))
}

func TestResponse_SuccessWithMeta(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		meta := &response.Meta{
			NextCursor: "abc123",
			Limit:      100,
			Total:      500,
			CacheHit:   true,
			LatencyMs:  15,
		}
		return response.SuccessWithMeta(c, []string{"item1", "item2"}, meta)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result response.Response
	err = sonic.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.NotNil(t, result.Meta)
	assert.Equal(t, "abc123", result.Meta.NextCursor)
	assert.Equal(t, 100, result.Meta.Limit)
	assert.Equal(t, 500, result.Meta.Total)
	assert.True(t, result.Meta.CacheHit)
}

func BenchmarkResponse_Success(b *testing.B) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Success(c, map[string]string{"message": "hello"})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Test(req)
	}
}

func BenchmarkResponse_Raw(b *testing.B) {
	app := fiber.New()

	rawData := []byte(`{"raw": "data", "number": 123, "array": [1,2,3]}`)

	app.Get("/test", func(c *fiber.Ctx) error {
		return response.Raw(c, rawData)
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Test(req)
	}
}
