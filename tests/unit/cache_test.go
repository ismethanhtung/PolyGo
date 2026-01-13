package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/config"
)

func TestCache_SetAndGet(t *testing.T) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 20, // 1MB
		NumCounters: 1e6,
		BufferItems: 64,
		DefaultTTL:  time.Second,
	}

	c, err := cache.New(cfg)
	require.NoError(t, err)
	defer c.Close()

	// Test Set and Get
	key := "test-key"
	value := []byte("test-value")

	ok := c.Set(key, value, time.Minute)
	assert.True(t, ok)

	// Wait for set to complete
	c.Wait()

	// Get value
	result, found := c.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, result)
}

func TestCache_GetMiss(t *testing.T) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 20,
		NumCounters: 1e6,
		BufferItems: 64,
		DefaultTTL:  time.Second,
	}

	c, err := cache.New(cfg)
	require.NoError(t, err)
	defer c.Close()

	// Get non-existent key
	result, found := c.Get("non-existent")
	assert.False(t, found)
	assert.Nil(t, result)
}

func TestCache_TTLExpiration(t *testing.T) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 20,
		NumCounters: 1e6,
		BufferItems: 64,
		DefaultTTL:  50 * time.Millisecond,
	}

	c, err := cache.New(cfg)
	require.NoError(t, err)
	defer c.Close()

	key := "expiring-key"
	value := []byte("expiring-value")

	c.Set(key, value, 50*time.Millisecond)
	c.Wait()

	// Should exist immediately
	_, found := c.Get(key)
	assert.True(t, found)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	_, found = c.Get(key)
	assert.False(t, found)
}

func TestCache_Delete(t *testing.T) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 20,
		NumCounters: 1e6,
		BufferItems: 64,
		DefaultTTL:  time.Second,
	}

	c, err := cache.New(cfg)
	require.NoError(t, err)
	defer c.Close()

	key := "delete-key"
	value := []byte("delete-value")

	c.Set(key, value, time.Minute)
	c.Wait()

	// Delete
	c.Delete(key)

	// Should not exist
	_, found := c.Get(key)
	assert.False(t, found)
}

func TestCache_SetJSON(t *testing.T) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 20,
		NumCounters: 1e6,
		BufferItems: 64,
		DefaultTTL:  time.Second,
	}

	c, err := cache.New(cfg)
	require.NoError(t, err)
	defer c.Close()

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	key := "json-key"
	original := TestStruct{Name: "test", Value: 42}

	ok := c.SetJSON(key, original, time.Minute)
	assert.True(t, ok)
	c.Wait()

	var result TestStruct
	found := c.GetJSON(key, &result)
	assert.True(t, found)
	assert.Equal(t, original, result)
}

func TestCacheKey_Generation(t *testing.T) {
	assert.Equal(t, "markets:123", cache.MarketKey("123"))
	assert.Equal(t, "markets:list:active=true", cache.MarketsListKey("active=true"))
	assert.Equal(t, "events:456", cache.EventKey("456"))
	assert.Equal(t, "price:token123", cache.PriceKey("token123"))
	assert.Equal(t, "book:token456", cache.OrderBookKey("token456"))
	assert.Equal(t, "spread:token789", cache.SpreadKey("token789"))
}

func BenchmarkCache_Set(b *testing.B) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 30,
		NumCounters: 1e7,
		BufferItems: 64,
		DefaultTTL:  time.Second,
	}

	c, _ := cache.New(cfg)
	defer c.Close()

	value := []byte("benchmark-value-with-some-data-to-simulate-real-payload")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("bench-key", value, time.Minute)
	}
}

func BenchmarkCache_Get(b *testing.B) {
	cfg := &config.CacheConfig{
		MaxCost:     1 << 30,
		NumCounters: 1e7,
		BufferItems: 64,
		DefaultTTL:  time.Second,
	}

	c, _ := cache.New(cfg)
	defer c.Close()

	value := []byte("benchmark-value-with-some-data-to-simulate-real-payload")
	c.Set("bench-key", value, time.Minute)
	c.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("bench-key")
	}
}
