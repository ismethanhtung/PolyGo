package models

import "time"

// Position represents a user's position in a market
type Position struct {
	Asset           string  `json:"asset"`
	ConditionID     string  `json:"conditionId"`
	Size            string  `json:"size"`
	AverageCost     string  `json:"avgCost"`
	CurrentPrice    string  `json:"currentPrice,omitempty"`
	PercentChange   float64 `json:"percentChange,omitempty"`
	RealizedPnL     string  `json:"realizedPnl,omitempty"`
	UnrealizedPnL   string  `json:"unrealizedPnl,omitempty"`
	CurVal          string  `json:"curVal,omitempty"`
	TotalBought     string  `json:"totalBought,omitempty"`
	TotalSold       string  `json:"totalSold,omitempty"`
	Outcome         string  `json:"outcome,omitempty"`
	OutcomeIndex    int     `json:"outcomeIndex,omitempty"`
	ProxyWalletAddr string  `json:"proxyWalletAddress,omitempty"`
}

// PositionsResponse represents positions list response
type PositionsResponse struct {
	Data       []Position `json:"data"`
	NextCursor string     `json:"next_cursor,omitempty"`
}

// Activity represents user activity entry
type Activity struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Action      string    `json:"action"`
	Description string    `json:"description,omitempty"`
	Market      string    `json:"market,omitempty"`
	Asset       string    `json:"asset,omitempty"`
	Amount      string    `json:"amount,omitempty"`
	Price       string    `json:"price,omitempty"`
	Side        string    `json:"side,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	TxHash      string    `json:"transactionHash,omitempty"`
}

// ActivityResponse represents activity list response
type ActivityResponse struct {
	Data       []Activity `json:"data"`
	NextCursor string     `json:"next_cursor,omitempty"`
}

// UserBalance represents user balance info
type UserBalance struct {
	Balance           string `json:"balance"`
	AvailableBalance  string `json:"availableBalance"`
	LockedBalance     string `json:"lockedBalance"`
	WithdrawableBalance string `json:"withdrawableBalance"`
}

// APICredentials represents user API credentials for trading
type APICredentials struct {
	APIKey       string `json:"api_key"`
	APISecret    string `json:"api_secret"`
	Passphrase   string `json:"passphrase"`
	PrivateKey   string `json:"private_key,omitempty"`
	FunderAddr   string `json:"funder_address,omitempty"`
}
