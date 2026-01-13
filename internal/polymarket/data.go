package polymarket

import (
	"net/url"
	"strconv"
)

// DataClient handles Data API requests (positions, trades, activity)
type DataClient struct {
	client *Client
}

// NewDataClient creates a new Data client
func NewDataClient(client *Client) *DataClient {
	return &DataClient{client: client}
}

// GetPositions retrieves user positions
func (d *DataClient) GetPositions(address string, limit int, cursor string) ([]byte, error) {
	query := url.Values{}
	query.Set("user", address)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		query.Set("next_cursor", cursor)
	}

	u := d.client.Data("/positions?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetPositionsByMarket retrieves positions for a specific market
func (d *DataClient) GetPositionsByMarket(address, marketID string) ([]byte, error) {
	query := url.Values{}
	query.Set("user", address)
	query.Set("market", marketID)

	u := d.client.Data("/positions?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetTrades retrieves user trades
func (d *DataClient) GetTrades(address string, limit int, cursor string) ([]byte, error) {
	query := url.Values{}
	query.Set("user", address)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		query.Set("next_cursor", cursor)
	}

	u := d.client.Data("/trades?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetTradesByMarket retrieves trades for a specific market
func (d *DataClient) GetTradesByMarket(address, marketID string, limit int) ([]byte, error) {
	query := url.Values{}
	query.Set("user", address)
	query.Set("market", marketID)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	u := d.client.Data("/trades?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetActivity retrieves user activity
func (d *DataClient) GetActivity(address string, limit int, cursor string) ([]byte, error) {
	query := url.Values{}
	query.Set("user", address)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		query.Set("next_cursor", cursor)
	}

	u := d.client.Data("/activity?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetMarketTrades retrieves public trades for a market (no auth required)
func (d *DataClient) GetMarketTrades(marketID string, limit int, cursor string) ([]byte, error) {
	query := url.Values{}
	query.Set("market", marketID)
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		query.Set("next_cursor", cursor)
	}

	u := d.client.Data("/trades?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetPriceHistory retrieves price history for a market
func (d *DataClient) GetPriceHistory(tokenID string, interval string, fidelity int) ([]byte, error) {
	query := url.Values{}
	query.Set("clob_token_id", tokenID)
	if interval != "" {
		query.Set("interval", interval) // e.g., "1d", "1h", "max"
	}
	if fidelity > 0 {
		query.Set("fidelity", strconv.Itoa(fidelity))
	}

	u := d.client.Data("/prices-history?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetTimeseriesData retrieves timeseries data for a market
func (d *DataClient) GetTimeseriesData(conditionID string, startTs, endTs int64) ([]byte, error) {
	query := url.Values{}
	query.Set("condition_id", conditionID)
	if startTs > 0 {
		query.Set("start_ts", strconv.FormatInt(startTs, 10))
	}
	if endTs > 0 {
		query.Set("end_ts", strconv.FormatInt(endTs, 10))
	}

	u := d.client.Data("/timeseries?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetTopMovers retrieves top moving markets
func (d *DataClient) GetTopMovers(limit int) ([]byte, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	u := d.client.Data("/top-movers?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetVolume retrieves volume data
func (d *DataClient) GetVolume(conditionID string) ([]byte, error) {
	query := url.Values{}
	query.Set("condition_id", conditionID)

	u := d.client.Data("/volume?" + query.Encode())
	return d.client.Get(u, nil)
}

// GetLeaderboard retrieves the trading leaderboard
func (d *DataClient) GetLeaderboard(limit int) ([]byte, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	u := d.client.Data("/leaderboard?" + query.Encode())
	return d.client.Get(u, nil)
}
