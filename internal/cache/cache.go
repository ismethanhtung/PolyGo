package cache

import (
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dgraph-io/ristretto"
	"github.com/polygo/internal/config"
)

// Cache wraps ristretto cache with typed methods
type Cache struct {
	store  *ristretto.Cache
	config *config.CacheConfig
	pool   sync.Pool // Pool for byte slices
}

// CacheEntry represents a cached entry with metadata
type CacheEntry struct {
	Data      []byte
	CreatedAt time.Time
	TTL       time.Duration
}

// New creates a new cache instance
func New(cfg *config.CacheConfig) (*Cache, error) {
	store, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.NumCounters,
		MaxCost:     cfg.MaxCost,
		BufferItems: cfg.BufferItems,
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}

	return &Cache{
		store:  store,
		config: cfg,
		pool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate 4KB buffers
				return make([]byte, 0, 4096)
			},
		},
	}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) ([]byte, bool) {
	val, found := c.store.Get(key)
	if !found {
		return nil, false
	}
	
	data, ok := val.([]byte)
	if !ok {
		return nil, false
	}
	
	return data, true
}

// GetJSON retrieves and unmarshals a value from cache
func (c *Cache) GetJSON(key string, dest interface{}) bool {
	data, found := c.Get(key)
	if !found {
		return false
	}
	
	if err := sonic.Unmarshal(data, dest); err != nil {
		return false
	}
	
	return true
}

// Set stores a value in cache with TTL
func (c *Cache) Set(key string, value []byte, ttl time.Duration) bool {
	// Make a copy to avoid data races
	data := make([]byte, len(value))
	copy(data, value)
	
	return c.store.SetWithTTL(key, data, int64(len(data)), ttl)
}

// SetJSON marshals and stores a value in cache
func (c *Cache) SetJSON(key string, value interface{}, ttl time.Duration) bool {
	data, err := sonic.Marshal(value)
	if err != nil {
		return false
	}
	
	return c.Set(key, data, ttl)
}

// SetWithDefaultTTL stores a value with default TTL
func (c *Cache) SetWithDefaultTTL(key string, value []byte) bool {
	return c.Set(key, value, c.config.DefaultTTL)
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) {
	c.store.Del(key)
}

// Clear removes all values from cache
func (c *Cache) Clear() {
	c.store.Clear()
}

// Wait waits for all pending sets to complete
func (c *Cache) Wait() {
	c.store.Wait()
}

// Close closes the cache
func (c *Cache) Close() {
	c.store.Close()
}

// Metrics returns cache metrics
func (c *Cache) Metrics() *ristretto.Metrics {
	return c.store.Metrics
}

// HitRatio returns the cache hit ratio
func (c *Cache) HitRatio() float64 {
	metrics := c.store.Metrics
	if metrics == nil {
		return 0
	}
	return metrics.Ratio()
}

// GetConfig returns the cache configuration (for accessing TTL values)
func (c *Cache) GetConfig() *config.CacheConfig {
	return c.config
}

// CacheKey helpers for consistent key generation
const (
	PrefixMarkets   = "markets:"
	PrefixEvents    = "events:"
	PrefixPrice     = "price:"
	PrefixOrderBook = "book:"
	PrefixSpread    = "spread:"
	PrefixTrades    = "trades:"
	PrefixPositions = "positions:"
)

// MarketKey generates a cache key for market
func MarketKey(id string) string {
	return PrefixMarkets + id
}

// MarketsListKey generates a cache key for markets list
func MarketsListKey(params string) string {
	return PrefixMarkets + "list:" + params
}

// EventKey generates a cache key for event
func EventKey(id string) string {
	return PrefixEvents + id
}

// EventsListKey generates a cache key for events list
func EventsListKey(params string) string {
	return PrefixEvents + "list:" + params
}

// PriceKey generates a cache key for price
func PriceKey(tokenID string) string {
	return PrefixPrice + tokenID
}

// OrderBookKey generates a cache key for order book
func OrderBookKey(tokenID string) string {
	return PrefixOrderBook + tokenID
}

// SpreadKey generates a cache key for spread
func SpreadKey(tokenID string) string {
	return PrefixSpread + tokenID
}
