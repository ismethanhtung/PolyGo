package integration

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/polygo/internal/api"
	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/config"
)

func setupTestServer(t *testing.T) *fiber.App {
	cfg := config.DefaultConfig()
	cfg.Server.Debug = true

	c, err := cache.New(&cfg.Cache)
	require.NoError(t, err)

	server, err := api.NewServer(cfg, c)
	require.NoError(t, err)

	return server.GetApp()
}

func TestHealthEndpoint(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = sonic.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
}

func TestReadyEndpoint(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestStatsEndpoint(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/stats", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = sonic.Unmarshal(body, &result)
	require.NoError(t, err)

	data := result["data"].(map[string]interface{})
	assert.NotEmpty(t, data["go_version"])
	assert.NotZero(t, data["num_cpu"])
}

func TestMarketsEndpoint_RequiresNoAuth(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/markets?limit=10", nil)
	resp, err := app.Test(req, 5*int(time.Second))
	require.NoError(t, err)

	// Should not return 401
	assert.NotEqual(t, 401, resp.StatusCode)
}

func TestEventsEndpoint_RequiresNoAuth(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/events?limit=10", nil)
	resp, err := app.Test(req, 5*int(time.Second))
	require.NoError(t, err)

	// Should not return 401
	assert.NotEqual(t, 401, resp.StatusCode)
}

func TestPriceEndpoint_InvalidTokenID(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/price/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	// Should return 404 (route not found) or redirect
	assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 301)
}

func TestOrdersEndpoint_RequiresAuth(t *testing.T) {
	app := setupTestServer(t)

	// POST without auth should fail
	req := httptest.NewRequest("POST", "/api/v1/orders", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestCORS_Headers(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("OPTIONS", "/api/v1/markets", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestRateLimitHeaders(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/markets", nil)
	resp, err := app.Test(req, 5*int(time.Second))
	require.NoError(t, err)

	// Rate limit headers should be present
	assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Limit"))
}

func TestResponseTime_Header(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.NotEmpty(t, resp.Header.Get("X-Response-Time"))
}

func TestPositions_RequiresAddress(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/positions", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	assert.Equal(t, 400, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	sonic.Unmarshal(body, &result)

	assert.False(t, result["success"].(bool))
}

func TestTopMovers_Endpoint(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/top-movers?limit=5", nil)
	resp, err := app.Test(req, 5*int(time.Second))
	require.NoError(t, err)

	// Should not require auth
	assert.NotEqual(t, 401, resp.StatusCode)
}

func TestLeaderboard_Endpoint(t *testing.T) {
	app := setupTestServer(t)

	req := httptest.NewRequest("GET", "/api/v1/leaderboard", nil)
	resp, err := app.Test(req, 5*int(time.Second))
	require.NoError(t, err)

	// Should not require auth
	assert.NotEqual(t, 401, resp.StatusCode)
}

func BenchmarkHealthEndpoint(b *testing.B) {
	cfg := config.DefaultConfig()
	c, _ := cache.New(&cfg.Cache)
	server, _ := api.NewServer(cfg, c)
	app := server.GetApp()

	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Test(req, -1)
	}
}

func BenchmarkMarketsEndpoint(b *testing.B) {
	cfg := config.DefaultConfig()
	c, _ := cache.New(&cfg.Cache)
	server, _ := api.NewServer(cfg, c)
	app := server.GetApp()

	req := httptest.NewRequest("GET", "/api/v1/markets?limit=10", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.Test(req, 5*int(time.Second))
	}
}
