package models

import "time"

// Market represents a Polymarket market
type Market struct {
	ID                  string    `json:"id"`
	Question            string    `json:"question"`
	Description         string    `json:"description"`
	ConditionID         string    `json:"conditionId"`
	Slug                string    `json:"slug"`
	EndDate             time.Time `json:"endDate"`
	Liquidity           string    `json:"liquidity"`
	Volume              string    `json:"volume"`
	Volume24hr          string    `json:"volume24hr"`
	Active              bool      `json:"active"`
	Closed              bool      `json:"closed"`
	MarketType          string    `json:"marketType"`
	OutcomePrices       []string  `json:"outcomePrices"`
	Outcomes            []string  `json:"outcomes"`
	ClobTokenIDs        []string  `json:"clobTokenIds"`
	AcceptingOrders     bool      `json:"acceptingOrders"`
	AcceptingOrdersTS   time.Time `json:"acceptingOrdersTimestamp"`
	EnableOrderBook     bool      `json:"enableOrderBook"`
	NegRisk             bool      `json:"negRisk"`
	NegRiskMarketID     string    `json:"negRiskMarketId,omitempty"`
	NegRiskRequestID    string    `json:"negRiskRequestId,omitempty"`
	Icon                string    `json:"icon,omitempty"`
	Image               string    `json:"image,omitempty"`
	RewardsMinSize      float64   `json:"rewardsMinSize,omitempty"`
	RewardsMaxSpread    float64   `json:"rewardsMaxSpread,omitempty"`
	SpreadMultiplierMin float64   `json:"spreadMultiplierMin,omitempty"`
	SpreadMultiplierMax float64   `json:"spreadMultiplierMax,omitempty"`
}

// MarketsResponse represents the API response for markets list
type MarketsResponse struct {
	Data       []Market `json:"data"`
	NextCursor string   `json:"next_cursor,omitempty"`
	Limit      int      `json:"limit"`
}

// MarketQueryParams represents query parameters for market filtering
type MarketQueryParams struct {
	Limit      int    `query:"limit"`
	Cursor     string `query:"cursor"`
	Active     *bool  `query:"active"`
	Closed     *bool  `query:"closed"`
	Slug       string `query:"slug"`
	EventSlug  string `query:"event_slug"`
	ClobTokenID string `query:"clob_token_id"`
}
