package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/polymarket"
	"github.com/polygo/pkg/response"
)

// DataHandler handles data-related endpoints (positions, trades, activity)
type DataHandler struct {
	data *polymarket.DataClient
}

// NewDataHandler creates a new data handler
func NewDataHandler(data *polymarket.DataClient) *DataHandler {
	return &DataHandler{data: data}
}

// GetPositions godoc
// @Summary Get user positions
// @Description Get all positions for a user address
// @Tags User Data
// @Accept json
// @Produce json
// @Param address query string true "User wallet address"
// @Param limit query int false "Limit results" default(100)
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} response.Response{data=[]models.Position}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/positions [get]
func (h *DataHandler) GetPositions(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return response.BadRequest(c, "Address is required")
	}
	
	limit := c.QueryInt("limit", 100)
	cursor := c.Query("cursor")
	
	data, err := h.data.GetPositions(address, limit, cursor)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetPositionsByMarket godoc
// @Summary Get positions for a specific market
// @Description Get user positions for a specific market
// @Tags User Data
// @Accept json
// @Produce json
// @Param address query string true "User wallet address"
// @Param market query string true "Market ID"
// @Success 200 {object} response.Response{data=[]models.Position}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/positions/market [get]
func (h *DataHandler) GetPositionsByMarket(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return response.BadRequest(c, "Address is required")
	}
	
	marketID := c.Query("market")
	if marketID == "" {
		return response.BadRequest(c, "Market ID is required")
	}
	
	data, err := h.data.GetPositionsByMarket(address, marketID)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetUserTrades godoc
// @Summary Get user trades
// @Description Get trade history for a user
// @Tags User Data
// @Accept json
// @Produce json
// @Param address query string true "User wallet address"
// @Param limit query int false "Limit results" default(100)
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} response.Response{data=[]models.Trade}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/trades [get]
func (h *DataHandler) GetUserTrades(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return response.BadRequest(c, "Address is required")
	}
	
	limit := c.QueryInt("limit", 100)
	cursor := c.Query("cursor")
	
	data, err := h.data.GetTrades(address, limit, cursor)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetUserTradesByMarket godoc
// @Summary Get user trades for a specific market
// @Description Get trade history for a user in a specific market
// @Tags User Data
// @Accept json
// @Produce json
// @Param address query string true "User wallet address"
// @Param market query string true "Market ID"
// @Param limit query int false "Limit results" default(100)
// @Success 200 {object} response.Response{data=[]models.Trade}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/trades/market [get]
func (h *DataHandler) GetUserTradesByMarket(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return response.BadRequest(c, "Address is required")
	}
	
	marketID := c.Query("market")
	if marketID == "" {
		return response.BadRequest(c, "Market ID is required")
	}
	
	limit := c.QueryInt("limit", 100)
	
	data, err := h.data.GetTradesByMarket(address, marketID, limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetActivity godoc
// @Summary Get user activity
// @Description Get activity log for a user
// @Tags User Data
// @Accept json
// @Produce json
// @Param address query string true "User wallet address"
// @Param limit query int false "Limit results" default(100)
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} response.Response{data=[]models.Activity}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/activity [get]
func (h *DataHandler) GetActivity(c *fiber.Ctx) error {
	address := c.Query("address")
	if address == "" {
		return response.BadRequest(c, "Address is required")
	}
	
	limit := c.QueryInt("limit", 100)
	cursor := c.Query("cursor")
	
	data, err := h.data.GetActivity(address, limit, cursor)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetMarketTrades godoc
// @Summary Get public market trades
// @Description Get trade history for a market (no auth required)
// @Tags Trades
// @Accept json
// @Produce json
// @Param market query string true "Market ID"
// @Param limit query int false "Limit results" default(100)
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} response.Response{data=[]models.Trade}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/market-trades [get]
func (h *DataHandler) GetMarketTrades(c *fiber.Ctx) error {
	marketID := c.Query("market")
	if marketID == "" {
		return response.BadRequest(c, "Market ID is required")
	}
	
	limit := c.QueryInt("limit", 100)
	cursor := c.Query("cursor")
	
	data, err := h.data.GetMarketTrades(marketID, limit, cursor)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetPriceHistory godoc
// @Summary Get price history
// @Description Get historical price data for a token
// @Tags Prices
// @Accept json
// @Produce json
// @Param token_id path string true "CLOB Token ID"
// @Param interval query string false "Time interval (1h, 1d, max)" default(1d)
// @Param fidelity query int false "Data fidelity/resolution"
// @Success 200 {object} response.Response{data=object}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/price-history/{token_id} [get]
func (h *DataHandler) GetPriceHistory(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	interval := c.Query("interval", "1d")
	fidelity := c.QueryInt("fidelity", 0)
	
	data, err := h.data.GetPriceHistory(tokenID, interval, fidelity)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetTimeseries godoc
// @Summary Get timeseries data
// @Description Get timeseries data for a market
// @Tags Prices
// @Accept json
// @Produce json
// @Param condition_id query string true "Condition ID"
// @Param start_ts query int false "Start timestamp (unix)"
// @Param end_ts query int false "End timestamp (unix)"
// @Success 200 {object} response.Response{data=object}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/timeseries [get]
func (h *DataHandler) GetTimeseries(c *fiber.Ctx) error {
	conditionID := c.Query("condition_id")
	if conditionID == "" {
		return response.BadRequest(c, "Condition ID is required")
	}
	
	startTs := int64(c.QueryInt("start_ts", 0))
	endTs := int64(c.QueryInt("end_ts", 0))
	
	data, err := h.data.GetTimeseriesData(conditionID, startTs, endTs)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetTopMovers godoc
// @Summary Get top moving markets
// @Description Get markets with the highest price changes
// @Tags Markets
// @Accept json
// @Produce json
// @Param limit query int false "Limit results" default(10)
// @Success 200 {object} response.Response{data=object}
// @Failure 500 {object} response.Response
// @Router /api/v1/top-movers [get]
func (h *DataHandler) GetTopMovers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	
	data, err := h.data.GetTopMovers(limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetLeaderboard godoc
// @Summary Get trading leaderboard
// @Description Get the top traders leaderboard
// @Tags User Data
// @Accept json
// @Produce json
// @Param limit query int false "Limit results" default(100)
// @Success 200 {object} response.Response{data=object}
// @Failure 500 {object} response.Response
// @Router /api/v1/leaderboard [get]
func (h *DataHandler) GetLeaderboard(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 100)
	
	data, err := h.data.GetLeaderboard(limit)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}
