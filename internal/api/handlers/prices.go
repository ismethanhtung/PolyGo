package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/models"
	"github.com/polygo/internal/polymarket"
	"github.com/polygo/pkg/response"
)

// PricesHandler handles price-related endpoints
type PricesHandler struct {
	clob *polymarket.ClobClient
}

// NewPricesHandler creates a new prices handler
func NewPricesHandler(clob *polymarket.ClobClient) *PricesHandler {
	return &PricesHandler{clob: clob}
}

// GetPrice godoc
// @Summary Get current price
// @Description Get the current price for a token
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_id path string true "Token ID"
// @Param side query string false "Order side (BUY/SELL)" default(BUY)
// @Success 200 {object} response.Response{data=models.Price}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/price/{token_id} [get]
func (h *PricesHandler) GetPrice(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	side := models.Side(strings.ToUpper(c.Query("side", "BUY")))
	if side != models.SideBuy && side != models.SideSell {
		return response.BadRequest(c, "Side must be BUY or SELL")
	}
	
	data, cacheHit, err := h.clob.GetPrice(tokenID, side)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetPrices godoc
// @Summary Get prices for multiple tokens
// @Description Get current prices for multiple tokens at once
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_ids query string true "Comma-separated token IDs"
// @Param side query string false "Order side (BUY/SELL)" default(BUY)
// @Success 200 {object} response.Response{data=[]models.Price}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/prices [get]
func (h *PricesHandler) GetPrices(c *fiber.Ctx) error {
	tokenIDsStr := c.Query("token_ids")
	if tokenIDsStr == "" {
		return response.BadRequest(c, "Token IDs are required")
	}
	
	tokenIDs := strings.Split(tokenIDsStr, ",")
	if len(tokenIDs) == 0 {
		return response.BadRequest(c, "At least one token ID is required")
	}
	
	side := models.Side(strings.ToUpper(c.Query("side", "BUY")))
	if side != models.SideBuy && side != models.SideSell {
		return response.BadRequest(c, "Side must be BUY or SELL")
	}
	
	data, err := h.clob.GetPrices(tokenIDs, side)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetOrderBook godoc
// @Summary Get order book
// @Description Get the full order book for a token
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_id path string true "Token ID"
// @Success 200 {object} response.Response{data=models.OrderBook}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/book/{token_id} [get]
func (h *PricesHandler) GetOrderBook(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	data, cacheHit, err := h.clob.GetOrderBook(tokenID)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetOrderBooks godoc
// @Summary Get order books for multiple tokens
// @Description Get order books for multiple tokens at once
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_ids query string true "Comma-separated token IDs"
// @Success 200 {object} response.Response{data=[]models.OrderBook}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/books [get]
func (h *PricesHandler) GetOrderBooks(c *fiber.Ctx) error {
	tokenIDsStr := c.Query("token_ids")
	if tokenIDsStr == "" {
		return response.BadRequest(c, "Token IDs are required")
	}
	
	tokenIDs := strings.Split(tokenIDsStr, ",")
	if len(tokenIDs) == 0 {
		return response.BadRequest(c, "At least one token ID is required")
	}
	
	data, err := h.clob.GetOrderBooks(tokenIDs)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetSpread godoc
// @Summary Get spread
// @Description Get the bid-ask spread for a token
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_id path string true "Token ID"
// @Success 200 {object} response.Response{data=models.Spread}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/spread/{token_id} [get]
func (h *PricesHandler) GetSpread(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	data, cacheHit, err := h.clob.GetSpread(tokenID)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetMidpoint godoc
// @Summary Get midpoint price
// @Description Get the midpoint price for a token
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_id path string true "Token ID"
// @Success 200 {object} response.Response{data=object}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/midpoint/{token_id} [get]
func (h *PricesHandler) GetMidpoint(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	data, cacheHit, err := h.clob.GetMidpoint(tokenID)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}

// GetMidpoints godoc
// @Summary Get midpoints for multiple tokens
// @Description Get midpoint prices for multiple tokens at once
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_ids query string true "Comma-separated token IDs"
// @Success 200 {object} response.Response{data=object}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/midpoints [get]
func (h *PricesHandler) GetMidpoints(c *fiber.Ctx) error {
	tokenIDsStr := c.Query("token_ids")
	if tokenIDsStr == "" {
		return response.BadRequest(c, "Token IDs are required")
	}
	
	tokenIDs := strings.Split(tokenIDsStr, ",")
	
	data, err := h.clob.GetMidpoints(tokenIDs)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetLastTradePrice godoc
// @Summary Get last trade price
// @Description Get the last trade price for a token
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_id path string true "Token ID"
// @Success 200 {object} response.Response{data=object}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/last-trade/{token_id} [get]
func (h *PricesHandler) GetLastTradePrice(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	data, cacheHit, err := h.clob.GetLastTradePrice(tokenID)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.RawWithCacheHeader(c, data, cacheHit)
}
