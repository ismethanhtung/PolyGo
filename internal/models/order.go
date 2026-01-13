package models

import "time"

// Side represents order side (BUY/SELL)
type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

// OrderType represents order type
type OrderType string

const (
	OrderTypeGTC OrderType = "GTC" // Good Till Cancelled
	OrderTypeFOK OrderType = "FOK" // Fill Or Kill
	OrderTypeGTD OrderType = "GTD" // Good Till Date
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusLive      OrderStatus = "LIVE"
	OrderStatusMatched   OrderStatus = "MATCHED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order represents a trading order
type Order struct {
	ID              string      `json:"id"`
	MarketID        string      `json:"market"`
	Asset           string      `json:"asset_id"`
	Side            Side        `json:"side"`
	Price           string      `json:"price"`
	OriginalSize    string      `json:"original_size"`
	SizeMatched     string      `json:"size_matched"`
	Status          OrderStatus `json:"status"`
	Type            OrderType   `json:"type"`
	Owner           string      `json:"owner"`
	Expiration      int64       `json:"expiration,omitempty"`
	AssociateTradeID string     `json:"associate_trade_id,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	Outcome         string      `json:"outcome,omitempty"`
}

// OrderBook represents the order book for a token
type OrderBook struct {
	TokenID   string      `json:"token_id"`
	Bids      []PriceLevel `json:"bids"`
	Asks      []PriceLevel `json:"asks"`
	Hash      string      `json:"hash"`
	Timestamp int64       `json:"timestamp"`
}

// PriceLevel represents a price level in the order book
type PriceLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// Price represents current price info
type Price struct {
	TokenID string `json:"token_id"`
	Price   string `json:"price"`
	Side    Side   `json:"side"`
}

// Spread represents bid-ask spread
type Spread struct {
	TokenID   string `json:"token_id"`
	BidPrice  string `json:"bid"`
	AskPrice  string `json:"ask"`
	SpreadAbs string `json:"spread_abs"`
	SpreadPct string `json:"spread_pct"`
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	TokenID    string    `json:"tokenID" validate:"required"`
	Side       Side      `json:"side" validate:"required"`
	Price      string    `json:"price" validate:"required"`
	Size       string    `json:"size" validate:"required"`
	Type       OrderType `json:"type"`
	Expiration int64     `json:"expiration,omitempty"`
}

// OrdersResponse represents orders list response
type OrdersResponse struct {
	Data       []Order `json:"data"`
	NextCursor string  `json:"next_cursor,omitempty"`
}

// Trade represents a completed trade
type Trade struct {
	ID            string    `json:"id"`
	TakerOrderID  string    `json:"taker_order_id"`
	Market        string    `json:"market"`
	Asset         string    `json:"asset_id"`
	Side          Side      `json:"side"`
	Price         string    `json:"price"`
	Size          string    `json:"size"`
	Fee           string    `json:"fee,omitempty"`
	TradeOwner    string    `json:"trader_side,omitempty"`
	Bucket        int       `json:"bucket_index,omitempty"`
	TransactionHash string  `json:"transaction_hash,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	MatchTime     time.Time `json:"match_time,omitempty"`
}

// TradesResponse represents trades list response
type TradesResponse struct {
	Data       []Trade `json:"data"`
	NextCursor string  `json:"next_cursor,omitempty"`
}
