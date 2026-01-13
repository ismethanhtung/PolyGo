package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/models"
	"github.com/polygo/internal/polymarket"
	"github.com/polygo/pkg/response"
)

// MarketsHandler handles market-related endpoints
type MarketsHandler struct {
	gamma *polymarket.GammaClient
}

// NewMarketsHandler creates a new markets handler
func NewMarketsHandler(gamma *polymarket.GammaClient) *MarketsHandler {
	return &MarketsHandler{gamma: gamma}
}

// GetMarkets godoc
// @Summary List all markets
// @Description Get a list of markets with optional filtering
// @Tags Markets
// @Accept json
// @Produce json
// @Param limit query int false "Limit results" default(100)
// @Param cursor query string false "Pagination cursor"
// @Param active query bool false "Filter by active status"
// @Param closed query bool false "Filter by closed status"
// @Param slug query string false "Filter by slug"
// @Param event_slug query string false "Filter by event slug"
// @Param clob_token_id query string false "Filter by CLOB token ID"
// @Success 200 {object} response.Response{data=[]models.Market}
// @Failure 500 {object} response.Response
// @Router /api/v1/markets [get]
func (h *MarketsHandler) GetMarkets(c *fiber.Ctx) error {
	params := &models.MarketQueryParams{
		Limit:       c.QueryInt("limit", 100),
		Cursor:      c.Query("cursor"),
		Slug:        c.Query("slug"),
		EventSlug:   c.Query("event_slug"),
		ClobTokenID: c.Query("clob_token_id"),
	}
	
	// Handle bool pointers
	if c.Query("active") != "" {
		active := c.QueryBool("active")
		params.Active = &active
	}
	if c.Query("closed") != "" {
		closed := c.QueryBool("closed")
		params.Closed = &closed
	}
	
	data, cacheHit, err := h.gamma.GetMarkets(params)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetMarket godoc
// @Summary Get market by ID
// @Description Get detailed information about a specific market
// @Tags Markets
// @Accept json
// @Produce json
// @Param id path string true "Market ID"
// @Success 200 {object} response.Response{data=models.Market}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/markets/{id} [get]
func (h *MarketsHandler) GetMarket(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Market ID is required")
	}
	
	data, cacheHit, err := h.gamma.GetMarket(id)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	if len(data) == 0 || string(data) == "null" {
		return response.NotFound(c, "Market not found")
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetMarketBySlug godoc
// @Summary Get market by slug
// @Description Get market by its URL slug
// @Tags Markets
// @Accept json
// @Produce json
// @Param slug path string true "Market slug"
// @Success 200 {object} response.Response{data=models.Market}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/markets/slug/{slug} [get]
func (h *MarketsHandler) GetMarketBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.BadRequest(c, "Slug is required")
	}
	
	data, cacheHit, err := h.gamma.GetMarketBySlug(slug)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetMarketByToken godoc
// @Summary Get market by CLOB token ID
// @Description Get market associated with a specific CLOB token
// @Tags Markets
// @Accept json
// @Produce json
// @Param token_id path string true "CLOB Token ID"
// @Success 200 {object} response.Response{data=models.Market}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/markets/token/{token_id} [get]
func (h *MarketsHandler) GetMarketByToken(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	data, cacheHit, err := h.gamma.GetMarketByClobTokenID(tokenID)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}
