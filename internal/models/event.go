package models

import "time"

// Event represents a Polymarket event
type Event struct {
	ID                  string    `json:"id"`
	Ticker              string    `json:"ticker"`
	Slug                string    `json:"slug"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	StartDate           time.Time `json:"startDate,omitempty"`
	EndDate             time.Time `json:"endDate,omitempty"`
	Volume              string    `json:"volume"`
	Liquidity           string    `json:"liquidity"`
	Active              bool      `json:"active"`
	Closed              bool      `json:"closed"`
	Archived            bool      `json:"archived"`
	New                 bool      `json:"new"`
	Featured            bool      `json:"featured"`
	Restricted          bool      `json:"restricted"`
	LiquidityClaimable  bool      `json:"liquidityClaimable"`
	RewardsMinSize      float64   `json:"rewardsMinSize,omitempty"`
	RewardsMaxSpread    float64   `json:"rewardsMaxSpread,omitempty"`
	SpreadMultiplierMin float64   `json:"spreadMultiplierMin,omitempty"`
	SpreadMultiplierMax float64   `json:"spreadMultiplierMax,omitempty"`
	Icon                string    `json:"icon,omitempty"`
	Image               string    `json:"image,omitempty"`
	CoverImage          string    `json:"coverImage,omitempty"`
	Markets             []Market  `json:"markets,omitempty"`
	Tags                []Tag     `json:"tags,omitempty"`
}

// Tag represents an event tag
type Tag struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// EventsResponse represents the API response for events list
type EventsResponse struct {
	Data       []Event `json:"data"`
	NextCursor string  `json:"next_cursor,omitempty"`
	Limit      int     `json:"limit"`
}

// EventQueryParams represents query parameters for event filtering
type EventQueryParams struct {
	Limit    int    `query:"limit"`
	Cursor   string `query:"cursor"`
	Active   *bool  `query:"active"`
	Closed   *bool  `query:"closed"`
	Archived *bool  `query:"archived"`
	Slug     string `query:"slug"`
	Tag      string `query:"tag"`
}
