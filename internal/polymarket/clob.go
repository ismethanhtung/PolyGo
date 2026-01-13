package polymarket

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/models"
)

// ClobClient handles CLOB API requests (prices, order books, orders)
type ClobClient struct {
	client *Client
}

// NewClobClient creates a new CLOB client
func NewClobClient(client *Client) *ClobClient {
	return &ClobClient{client: client}
}

// GetPrice retrieves the current price for a token
func (c *ClobClient) GetPrice(tokenID string, side models.Side) ([]byte, bool, error) {
	cacheKey := cache.PriceKey(tokenID + ":" + string(side))
	url := c.client.CLOB(fmt.Sprintf("/price?token_id=%s&side=%s", tokenID, side))

	ttl := c.client.cache.config.PricesTTL
	return c.client.GetWithCache(url, cacheKey, ttl)
}

// GetPrices retrieves prices for multiple tokens
func (c *ClobClient) GetPrices(tokenIDs []string, side models.Side) ([]byte, error) {
	// Build comma-separated token IDs
	var tokens string
	for i, id := range tokenIDs {
		if i > 0 {
			tokens += ","
		}
		tokens += id
	}

	url := c.client.CLOB(fmt.Sprintf("/prices?token_ids=%s&side=%s", url.QueryEscape(tokens), side))
	return c.client.Get(url, nil)
}

// GetOrderBook retrieves the order book for a token
func (c *ClobClient) GetOrderBook(tokenID string) ([]byte, bool, error) {
	cacheKey := cache.OrderBookKey(tokenID)
	url := c.client.CLOB("/book?token_id=" + tokenID)

	ttl := c.client.cache.config.OrderBookTTL
	return c.client.GetWithCache(url, cacheKey, ttl)
}

// GetOrderBooks retrieves order books for multiple tokens
func (c *ClobClient) GetOrderBooks(tokenIDs []string) ([]byte, error) {
	var tokens string
	for i, id := range tokenIDs {
		if i > 0 {
			tokens += ","
		}
		tokens += id
	}

	url := c.client.CLOB("/books?token_ids=" + url.QueryEscape(tokens))
	return c.client.Get(url, nil)
}

// GetSpread retrieves the spread for a token
func (c *ClobClient) GetSpread(tokenID string) ([]byte, bool, error) {
	cacheKey := cache.SpreadKey(tokenID)
	url := c.client.CLOB("/spread?token_id=" + tokenID)

	ttl := c.client.cache.config.PricesTTL
	return c.client.GetWithCache(url, cacheKey, ttl)
}

// GetMidpoint retrieves the midpoint price for a token
func (c *ClobClient) GetMidpoint(tokenID string) ([]byte, bool, error) {
	cacheKey := cache.PriceKey("mid:" + tokenID)
	url := c.client.CLOB("/midpoint?token_id=" + tokenID)

	ttl := c.client.cache.config.PricesTTL
	return c.client.GetWithCache(url, cacheKey, ttl)
}

// GetMidpoints retrieves midpoints for multiple tokens
func (c *ClobClient) GetMidpoints(tokenIDs []string) ([]byte, error) {
	var tokens string
	for i, id := range tokenIDs {
		if i > 0 {
			tokens += ","
		}
		tokens += id
	}

	url := c.client.CLOB("/midpoints?token_ids=" + url.QueryEscape(tokens))
	return c.client.Get(url, nil)
}

// GetLastTradePrice retrieves the last trade price for a token
func (c *ClobClient) GetLastTradePrice(tokenID string) ([]byte, bool, error) {
	cacheKey := cache.PriceKey("last:" + tokenID)
	url := c.client.CLOB("/last-trade-price?token_id=" + tokenID)

	ttl := c.client.cache.config.PricesTTL
	return c.client.GetWithCache(url, cacheKey, ttl)
}

// OrderRequest represents an order request body
type OrderRequest struct {
	Order         interface{} `json:"order"`
	Owner         string      `json:"owner,omitempty"`
	OrderType     string      `json:"orderType,omitempty"`
}

// CreateOrder creates a new order (requires authentication)
func (c *ClobClient) CreateOrder(order *models.CreateOrderRequest, authHeaders map[string]string) ([]byte, error) {
	url := c.client.CLOB("/order")
	
	body, err := sonic.Marshal(order)
	if err != nil {
		return nil, err
	}

	return c.client.Post(url, body, &RequestOptions{Headers: authHeaders})
}

// CancelOrder cancels an existing order (requires authentication)
func (c *ClobClient) CancelOrder(orderID string, authHeaders map[string]string) ([]byte, error) {
	url := c.client.CLOB("/order/" + orderID)
	return c.client.Delete(url, &RequestOptions{Headers: authHeaders})
}

// CancelOrders cancels multiple orders (requires authentication)
func (c *ClobClient) CancelOrders(orderIDs []string, authHeaders map[string]string) ([]byte, error) {
	url := c.client.CLOB("/orders")
	
	body, err := sonic.Marshal(map[string][]string{"orderIds": orderIDs})
	if err != nil {
		return nil, err
	}

	return c.client.Delete(url, &RequestOptions{Headers: authHeaders})
	// Note: fasthttp doesn't support DELETE with body directly
	// May need to use POST with _method override or custom implementation
	_ = body
	return nil, fmt.Errorf("batch cancel not implemented")
}

// CancelAll cancels all orders for a market (requires authentication)
func (c *ClobClient) CancelAll(marketID string, authHeaders map[string]string) ([]byte, error) {
	url := c.client.CLOB("/cancel-all?market=" + marketID)
	return c.client.Delete(url, &RequestOptions{Headers: authHeaders})
}

// GetOrders retrieves orders for the authenticated user
func (c *ClobClient) GetOrders(params map[string]string, authHeaders map[string]string) ([]byte, error) {
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	u := c.client.CLOB("/orders")
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	return c.client.Get(u, &RequestOptions{Headers: authHeaders})
}

// GetOrder retrieves a specific order
func (c *ClobClient) GetOrder(orderID string, authHeaders map[string]string) ([]byte, error) {
	url := c.client.CLOB("/order/" + orderID)
	return c.client.Get(url, &RequestOptions{Headers: authHeaders})
}

// GetOpenOrders retrieves open orders for the authenticated user
func (c *ClobClient) GetOpenOrders(market string, authHeaders map[string]string) ([]byte, error) {
	url := c.client.CLOB("/orders/open")
	if market != "" {
		url += "?market=" + market
	}
	return c.client.Get(url, &RequestOptions{Headers: authHeaders})
}

// GetTradesHistory retrieves trade history
func (c *ClobClient) GetTradesHistory(tokenID string, limit int, before, after string) ([]byte, error) {
	query := url.Values{}
	query.Set("token_id", tokenID)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if before != "" {
		query.Set("before", before)
	}
	if after != "" {
		query.Set("after", after)
	}

	u := c.client.CLOB("/trades?" + query.Encode())
	return c.client.Get(u, nil)
}

// GetMarketTradesHistory retrieves trade history for a market
func (c *ClobClient) GetMarketTradesHistory(conditionID string, limit int) ([]byte, error) {
	query := url.Values{}
	query.Set("condition_id", conditionID)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	u := c.client.CLOB("/trades?" + query.Encode())
	return c.client.Get(u, nil)
}

// GetTickSize retrieves tick size for a token
func (c *ClobClient) GetTickSize(tokenID string) ([]byte, error) {
	url := c.client.CLOB("/tick-size?token_id=" + tokenID)
	return c.client.Get(url, nil)
}

// GetNegRisk retrieves neg risk info for a token
func (c *ClobClient) GetNegRisk(tokenID string) ([]byte, error) {
	url := c.client.CLOB("/neg-risk?token_id=" + tokenID)
	return c.client.Get(url, nil)
}
