package polymarket

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/models"
)

// GammaClient handles Gamma API requests (markets, events)
type GammaClient struct {
	client *Client
}

// NewGammaClient creates a new Gamma client
func NewGammaClient(client *Client) *GammaClient {
	return &GammaClient{client: client}
}

// GetEvents retrieves events from Gamma API
func (g *GammaClient) GetEvents(params *models.EventQueryParams) ([]byte, bool, error) {
	query := buildEventQuery(params)
	cacheKey := cache.EventsListKey(query)
	url := g.client.Gamma("/events" + query)

	return g.client.GetWithCache(url, cacheKey, g.client.config.ReadTimeout)
}

// GetEvent retrieves a single event by ID
func (g *GammaClient) GetEvent(id string) ([]byte, bool, error) {
	cacheKey := cache.EventKey(id)
	url := g.client.Gamma("/events/" + id)

	ttl := g.client.cache.GetConfig().EventsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// GetEventBySlug retrieves an event by slug
func (g *GammaClient) GetEventBySlug(slug string) ([]byte, bool, error) {
	cacheKey := cache.EventKey("slug:" + slug)
	url := g.client.Gamma("/events?slug=" + url.QueryEscape(slug))

	ttl := g.client.cache.GetConfig().EventsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// GetMarkets retrieves markets from Gamma API
func (g *GammaClient) GetMarkets(params *models.MarketQueryParams) ([]byte, bool, error) {
	query := buildMarketQuery(params)
	cacheKey := cache.MarketsListKey(query)
	url := g.client.Gamma("/markets" + query)

	ttl := g.client.cache.GetConfig().MarketsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// GetMarket retrieves a single market by ID
func (g *GammaClient) GetMarket(id string) ([]byte, bool, error) {
	cacheKey := cache.MarketKey(id)
	url := g.client.Gamma("/markets/" + id)

	ttl := g.client.cache.GetConfig().MarketsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// GetMarketBySlug retrieves a market by slug
func (g *GammaClient) GetMarketBySlug(slug string) ([]byte, bool, error) {
	cacheKey := cache.MarketKey("slug:" + slug)
	url := g.client.Gamma("/markets?slug=" + url.QueryEscape(slug))

	ttl := g.client.cache.GetConfig().MarketsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// GetMarketByConditionID retrieves a market by condition ID
func (g *GammaClient) GetMarketByConditionID(conditionID string) ([]byte, bool, error) {
	cacheKey := cache.MarketKey("condition:" + conditionID)
	url := g.client.Gamma("/markets?condition_id=" + conditionID)

	ttl := g.client.cache.GetConfig().MarketsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// GetMarketByClobTokenID retrieves a market by CLOB token ID
func (g *GammaClient) GetMarketByClobTokenID(tokenID string) ([]byte, bool, error) {
	cacheKey := cache.MarketKey("token:" + tokenID)
	url := g.client.Gamma("/markets?clob_token_id=" + tokenID)

	ttl := g.client.cache.GetConfig().MarketsTTL
	return g.client.GetWithCache(url, cacheKey, ttl)
}

// SearchEvents searches events by query
func (g *GammaClient) SearchEvents(query string, limit int) ([]byte, bool, error) {
	cacheKey := cache.EventsListKey("search:" + query + ":" + strconv.Itoa(limit))
	u := g.client.Gamma(fmt.Sprintf("/events?_q=%s&_limit=%d", url.QueryEscape(query), limit))

	ttl := g.client.cache.GetConfig().EventsTTL
	return g.client.GetWithCache(u, cacheKey, ttl)
}

// buildEventQuery builds query string for events
func buildEventQuery(params *models.EventQueryParams) string {
	if params == nil {
		return ""
	}

	v := url.Values{}

	if params.Limit > 0 {
		v.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		v.Set("next_cursor", params.Cursor)
	}
	if params.Active != nil {
		v.Set("active", strconv.FormatBool(*params.Active))
	}
	if params.Closed != nil {
		v.Set("closed", strconv.FormatBool(*params.Closed))
	}
	if params.Archived != nil {
		v.Set("archived", strconv.FormatBool(*params.Archived))
	}
	if params.Slug != "" {
		v.Set("slug", params.Slug)
	}
	if params.Tag != "" {
		v.Set("tag", params.Tag)
	}

	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

// buildMarketQuery builds query string for markets
func buildMarketQuery(params *models.MarketQueryParams) string {
	if params == nil {
		return ""
	}

	v := url.Values{}

	if params.Limit > 0 {
		v.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		v.Set("next_cursor", params.Cursor)
	}
	if params.Active != nil {
		v.Set("active", strconv.FormatBool(*params.Active))
	}
	if params.Closed != nil {
		v.Set("closed", strconv.FormatBool(*params.Closed))
	}
	if params.Slug != "" {
		v.Set("slug", params.Slug)
	}
	if params.EventSlug != "" {
		v.Set("event_slug", params.EventSlug)
	}
	if params.ClobTokenID != "" {
		v.Set("clob_token_id", params.ClobTokenID)
	}

	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}
