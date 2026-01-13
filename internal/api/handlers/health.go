package handlers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/polymarket"
	"github.com/polygo/pkg/response"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	cache     *cache.Cache
	wsManager *polymarket.WSManager
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(c *cache.Cache, ws *polymarket.WSManager) *HealthHandler {
	return &HealthHandler{
		cache:     c,
		wsManager: ws,
		startTime: time.Now(),
	}
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Uptime    string            `json:"uptime"`
	Timestamp int64             `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// Health godoc
// @Summary Health check
// @Description Check if the server is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	services := map[string]string{
		"cache": "healthy",
	}
	
	if h.wsManager.IsConnected() {
		services["websocket"] = "connected"
	} else {
		services["websocket"] = "disconnected"
	}
	
	resp := HealthResponse{
		Status:    "healthy",
		Uptime:    time.Since(h.startTime).String(),
		Timestamp: time.Now().UnixMilli(),
		Services:  services,
	}
	
	return response.Success(c, resp)
}

// ReadyResponse represents readiness check response
type ReadyResponse struct {
	Ready     bool   `json:"ready"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Ready godoc
// @Summary Readiness check
// @Description Check if the server is ready to accept requests
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} ReadyResponse
// @Failure 503 {object} ReadyResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Check if cache is working
	testKey := "__ready_check__"
	h.cache.Set(testKey, []byte("ok"), time.Second)
	_, found := h.cache.Get(testKey)
	h.cache.Delete(testKey)
	
	if !found {
		return c.Status(fiber.StatusServiceUnavailable).JSON(ReadyResponse{
			Ready:     false,
			Message:   "Cache not ready",
			Timestamp: time.Now().UnixMilli(),
		})
	}
	
	return response.Success(c, ReadyResponse{
		Ready:     true,
		Timestamp: time.Now().UnixMilli(),
	})
}

// StatsResponse represents server statistics
type StatsResponse struct {
	Uptime       string  `json:"uptime"`
	GoVersion    string  `json:"go_version"`
	NumGoroutine int     `json:"num_goroutine"`
	NumCPU       int     `json:"num_cpu"`
	MemAlloc     uint64  `json:"mem_alloc_bytes"`
	MemTotal     uint64  `json:"mem_total_bytes"`
	MemSys       uint64  `json:"mem_sys_bytes"`
	CacheHitRate float64 `json:"cache_hit_rate"`
	Timestamp    int64   `json:"timestamp"`
}

// Stats godoc
// @Summary Server statistics
// @Description Get server runtime statistics
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} StatsResponse
// @Router /stats [get]
func (h *HealthHandler) Stats(c *fiber.Ctx) error {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	resp := StatsResponse{
		Uptime:       time.Since(h.startTime).String(),
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		MemAlloc:     mem.Alloc,
		MemTotal:     mem.TotalAlloc,
		MemSys:       mem.Sys,
		CacheHitRate: h.cache.HitRatio(),
		Timestamp:    time.Now().UnixMilli(),
	}
	
	return response.Success(c, resp)
}
