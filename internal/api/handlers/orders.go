package handlers

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/api/middleware"
	"github.com/polygo/internal/config"
	"github.com/polygo/internal/models"
	"github.com/polygo/internal/polymarket"
	"github.com/polygo/pkg/response"
)

// OrdersHandler handles order-related endpoints
type OrdersHandler struct {
	clob       *polymarket.ClobClient
	authConfig *config.AuthConfig
}

// NewOrdersHandler creates a new orders handler
func NewOrdersHandler(clob *polymarket.ClobClient, authConfig *config.AuthConfig) *OrdersHandler {
	return &OrdersHandler{
		clob:       clob,
		authConfig: authConfig,
	}
}

// getAuthHeaders extracts auth headers from context
func (h *OrdersHandler) getAuthHeaders(c *fiber.Ctx) map[string]string {
	creds := middleware.GetAuthCredentials(c)
	if creds == nil {
		return nil
	}
	return middleware.GetAuthHeaders(creds, h.authConfig)
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Place a new order on the market
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order details"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=models.Order}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders [post]
func (h *OrdersHandler) CreateOrder(c *fiber.Ctx) error {
	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	
	// Validate required fields
	if req.TokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	if req.Price == "" {
		return response.BadRequest(c, "Price is required")
	}
	if req.Size == "" {
		return response.BadRequest(c, "Size is required")
	}
	if req.Side != models.SideBuy && req.Side != models.SideSell {
		return response.BadRequest(c, "Side must be BUY or SELL")
	}
	
	// Default order type
	if req.Type == "" {
		req.Type = models.OrderTypeGTC
	}
	
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	data, err := h.clob.CreateOrder(&req, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetOrders godoc
// @Summary Get user orders
// @Description Get orders for the authenticated user
// @Tags Orders
// @Accept json
// @Produce json
// @Param market query string false "Filter by market"
// @Param status query string false "Filter by status"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]models.Order}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders [get]
func (h *OrdersHandler) GetOrders(c *fiber.Ctx) error {
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	params := make(map[string]string)
	if market := c.Query("market"); market != "" {
		params["market"] = market
	}
	if status := c.Query("status"); status != "" {
		params["status"] = status
	}
	
	data, err := h.clob.GetOrders(params, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Get details of a specific order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=models.Order}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders/{id} [get]
func (h *OrdersHandler) GetOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}
	
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	data, err := h.clob.GetOrder(orderID, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetOpenOrders godoc
// @Summary Get open orders
// @Description Get all open orders for the authenticated user
// @Tags Orders
// @Accept json
// @Produce json
// @Param market query string false "Filter by market"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]models.Order}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders/open [get]
func (h *OrdersHandler) GetOpenOrders(c *fiber.Ctx) error {
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	market := c.Query("market")
	
	data, err := h.clob.GetOpenOrders(market, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an existing order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders/{id} [delete]
func (h *OrdersHandler) CancelOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return response.BadRequest(c, "Order ID is required")
	}
	
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	data, err := h.clob.CancelOrder(orderID, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// CancelAllOrders godoc
// @Summary Cancel all orders
// @Description Cancel all orders for a specific market
// @Tags Orders
// @Accept json
// @Produce json
// @Param market query string true "Market ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders/cancel-all [delete]
func (h *OrdersHandler) CancelAllOrders(c *fiber.Ctx) error {
	market := c.Query("market")
	if market == "" {
		return response.BadRequest(c, "Market is required")
	}
	
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	data, err := h.clob.CancelAll(market, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// GetTrades godoc
// @Summary Get public trades
// @Description Get trade history for a token
// @Tags Trades
// @Accept json
// @Produce json
// @Param token_id path string true "Token ID"
// @Param limit query int false "Limit results" default(100)
// @Param before query string false "Cursor for pagination (before)"
// @Param after query string false "Cursor for pagination (after)"
// @Success 200 {object} response.Response{data=[]models.Trade}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/trades/{token_id} [get]
func (h *OrdersHandler) GetTrades(c *fiber.Ctx) error {
	tokenID := c.Params("token_id")
	if tokenID == "" {
		return response.BadRequest(c, "Token ID is required")
	}
	
	limit := c.QueryInt("limit", 100)
	before := c.Query("before")
	after := c.Query("after")
	
	data, err := h.clob.GetTradesHistory(tokenID, limit, before, after)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}

// BatchCancelRequest represents batch cancel request
type BatchCancelRequest struct {
	OrderIDs []string `json:"orderIds"`
}

// CancelOrders godoc
// @Summary Cancel multiple orders
// @Description Cancel multiple orders by IDs
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body BatchCancelRequest true "Order IDs to cancel"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/orders/batch-cancel [post]
func (h *OrdersHandler) CancelOrders(c *fiber.Ctx) error {
	var req BatchCancelRequest
	if err := sonic.Unmarshal(c.Body(), &req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	
	if len(req.OrderIDs) == 0 {
		return response.BadRequest(c, "At least one order ID is required")
	}
	
	authHeaders := h.getAuthHeaders(c)
	if authHeaders == nil {
		return response.Unauthorized(c, "Authentication required")
	}
	
	data, err := h.clob.CancelOrders(req.OrderIDs, authHeaders)
	if err != nil {
		return response.InternalError(c, err)
	}
	
	return response.Raw(c, data)
}
