package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Polymarket PolymarketConfig `mapstructure:"polymarket"`
	Cache      CacheConfig      `mapstructure:"cache"`
	Auth       AuthConfig       `mapstructure:"auth"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	Prefork      bool          `mapstructure:"prefork"`
	Debug        bool          `mapstructure:"debug"`
}

// PolymarketConfig holds Polymarket API configuration
type PolymarketConfig struct {
	ClobBaseURL      string        `mapstructure:"clob_base_url"`
	GammaBaseURL     string        `mapstructure:"gamma_base_url"`
	DataBaseURL      string        `mapstructure:"data_base_url"`
	WsClobURL        string        `mapstructure:"ws_clob_url"`
	WsLiveDataURL    string        `mapstructure:"ws_live_data_url"`
	MaxConnsPerHost  int           `mapstructure:"max_conns_per_host"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	MaxIdleConnDur   time.Duration `mapstructure:"max_idle_conn_dur"`
	RetryCount       int           `mapstructure:"retry_count"`
	RetryWaitTime    time.Duration `mapstructure:"retry_wait_time"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	MaxCost        int64         `mapstructure:"max_cost"`
	NumCounters    int64         `mapstructure:"num_counters"`
	BufferItems    int64         `mapstructure:"buffer_items"`
	MarketsTTL     time.Duration `mapstructure:"markets_ttl"`
	EventsTTL      time.Duration `mapstructure:"events_ttl"`
	PricesTTL      time.Duration `mapstructure:"prices_ttl"`
	OrderBookTTL   time.Duration `mapstructure:"order_book_ttl"`
	DefaultTTL     time.Duration `mapstructure:"default_ttl"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKeyHeader     string `mapstructure:"api_key_header"`
	APISecretHeader  string `mapstructure:"api_secret_header"`
	PassphraseHeader string `mapstructure:"passphrase_header"`
	SignatureHeader  string `mapstructure:"signature_header"`
	TimestampHeader  string `mapstructure:"timestamp_header"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  30 * time.Second,
			Prefork:      false,
			Debug:        false,
		},
		Polymarket: PolymarketConfig{
			ClobBaseURL:     "https://clob.polymarket.com",
			GammaBaseURL:    "https://gamma-api.polymarket.com",
			DataBaseURL:     "https://data-api.polymarket.com",
			WsClobURL:       "wss://ws-subscriptions-clob.polymarket.com/ws/",
			WsLiveDataURL:   "wss://ws-live-data.polymarket.com",
			MaxConnsPerHost: 1000,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			MaxIdleConnDur:  30 * time.Second,
			RetryCount:      3,
			RetryWaitTime:   100 * time.Millisecond,
		},
		Cache: CacheConfig{
			MaxCost:      1 << 30,      // 1GB
			NumCounters:  1e7,          // 10M counters
			BufferItems:  64,           // 64 buffer items
			MarketsTTL:   30 * time.Second,
			EventsTTL:    30 * time.Second,
			PricesTTL:    100 * time.Millisecond,
			OrderBookTTL: 50 * time.Millisecond,
			DefaultTTL:   5 * time.Second,
		},
		Auth: AuthConfig{
			APIKeyHeader:     "POLY-API-KEY",
			APISecretHeader:  "POLY-API-SECRET",
			PassphraseHeader: "POLY-PASSPHRASE",
			SignatureHeader:  "POLY-SIGNATURE",
			TimestampHeader:  "POLY-TIMESTAMP",
		},
	}
}

// Load loads configuration from environment and config file
func Load() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/polygo")

	// Environment variables
	viper.SetEnvPrefix("POLYGO")
	viper.AutomaticEnv()

	// Bind environment variables
	bindEnvVars()

	// Try to read config file (not required)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, using defaults + env vars
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindEnvVars() {
	// Server
	viper.BindEnv("server.host", "POLYGO_HOST")
	viper.BindEnv("server.port", "POLYGO_PORT")
	viper.BindEnv("server.debug", "POLYGO_DEBUG")
	viper.BindEnv("server.prefork", "POLYGO_PREFORK")

	// Polymarket URLs
	viper.BindEnv("polymarket.clob_base_url", "POLYGO_CLOB_URL")
	viper.BindEnv("polymarket.gamma_base_url", "POLYGO_GAMMA_URL")
	viper.BindEnv("polymarket.data_base_url", "POLYGO_DATA_URL")

	// Cache
	viper.BindEnv("cache.max_cost", "POLYGO_CACHE_MAX_COST")
	viper.BindEnv("cache.markets_ttl", "POLYGO_CACHE_MARKETS_TTL")
	viper.BindEnv("cache.prices_ttl", "POLYGO_CACHE_PRICES_TTL")
}

// GetAddress returns the full address string
func (c *ServerConfig) GetAddress() string {
	return c.Host + ":" + string(rune(c.Port))
}
