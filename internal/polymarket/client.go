package polymarket

import (
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/config"
	"github.com/valyala/fasthttp"
)

// Client is the main Polymarket API client
type Client struct {
	httpClient *fasthttp.Client
	cache      *cache.Cache
	config     *config.PolymarketConfig

	// Base URLs
	clobURL  string
	gammaURL string
	dataURL  string

	// Request/Response pools for zero-allocation
	reqPool  sync.Pool
	respPool sync.Pool
}

// NewClient creates a new Polymarket client with optimized settings
func NewClient(cfg *config.PolymarketConfig, c *cache.Cache) *Client {
	client := &Client{
		httpClient: &fasthttp.Client{
			Name:                     "PolyGo/1.0",
			MaxConnsPerHost:          cfg.MaxConnsPerHost,
			MaxIdleConnDuration:      cfg.MaxIdleConnDur,
			ReadTimeout:              cfg.ReadTimeout,
			WriteTimeout:             cfg.WriteTimeout,
			NoDefaultUserAgentHeader: true,
			DisableHeaderNamesNormalizing: true,
			DisablePathNormalizing:   true,
		},
		cache:    c,
		config:   cfg,
		clobURL:  cfg.ClobBaseURL,
		gammaURL: cfg.GammaBaseURL,
		dataURL:  cfg.DataBaseURL,
	}

	// Initialize pools
	client.reqPool = sync.Pool{
		New: func() interface{} {
			return fasthttp.AcquireRequest()
		},
	}
	client.respPool = sync.Pool{
		New: func() interface{} {
			return fasthttp.AcquireResponse()
		},
	}

	return client
}

// acquireRequest gets a request from pool
func (c *Client) acquireRequest() *fasthttp.Request {
	return fasthttp.AcquireRequest()
}

// releaseRequest returns a request to pool
func (c *Client) releaseRequest(req *fasthttp.Request) {
	fasthttp.ReleaseRequest(req)
}

// acquireResponse gets a response from pool
func (c *Client) acquireResponse() *fasthttp.Response {
	return fasthttp.AcquireResponse()
}

// releaseResponse returns a response to pool
func (c *Client) releaseResponse(resp *fasthttp.Response) {
	fasthttp.ReleaseResponse(resp)
}

// RequestOptions holds options for HTTP requests
type RequestOptions struct {
	Headers map[string]string
	Timeout time.Duration
}

// doRequest performs an HTTP request with retry logic
func (c *Client) doRequest(method, url string, body []byte, opts *RequestOptions) ([]byte, error) {
	req := c.acquireRequest()
	resp := c.acquireResponse()
	defer c.releaseRequest(req)
	defer c.releaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	if body != nil {
		req.SetBody(body)
	}

	timeout := c.config.ReadTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	var lastErr error
	for i := 0; i <= c.config.RetryCount; i++ {
		if i > 0 {
			time.Sleep(c.config.RetryWaitTime * time.Duration(i))
		}

		err := c.httpClient.DoTimeout(req, resp, timeout)
		if err != nil {
			lastErr = err
			continue
		}

		statusCode := resp.StatusCode()
		if statusCode >= 200 && statusCode < 300 {
			// Make a copy of the body
			result := make([]byte, len(resp.Body()))
			copy(result, resp.Body())
			return result, nil
		}

		if statusCode >= 500 {
			lastErr = fmt.Errorf("server error: %d", statusCode)
			continue
		}

		// Client error, don't retry
		return nil, fmt.Errorf("request failed with status %d: %s", statusCode, resp.Body())
	}

	return nil, fmt.Errorf("request failed after %d retries: %v", c.config.RetryCount, lastErr)
}

// Get performs a GET request
func (c *Client) Get(url string, opts *RequestOptions) ([]byte, error) {
	return c.doRequest("GET", url, nil, opts)
}

// GetWithCache performs a GET request with caching
func (c *Client) GetWithCache(url, cacheKey string, ttl time.Duration) ([]byte, bool, error) {
	// Check cache first
	if data, found := c.cache.Get(cacheKey); found {
		return data, true, nil
	}

	// Fetch from API
	data, err := c.Get(url, nil)
	if err != nil {
		return nil, false, err
	}

	// Store in cache
	c.cache.Set(cacheKey, data, ttl)

	return data, false, nil
}

// Post performs a POST request
func (c *Client) Post(url string, body []byte, opts *RequestOptions) ([]byte, error) {
	return c.doRequest("POST", url, body, opts)
}

// Delete performs a DELETE request
func (c *Client) Delete(url string, opts *RequestOptions) ([]byte, error) {
	return c.doRequest("DELETE", url, nil, opts)
}

// GetJSON performs a GET request and unmarshals the response
func (c *Client) GetJSON(url string, dest interface{}, opts *RequestOptions) error {
	data, err := c.Get(url, opts)
	if err != nil {
		return err
	}
	return sonic.Unmarshal(data, dest)
}

// PostJSON performs a POST request with JSON body and unmarshals the response
func (c *Client) PostJSON(url string, body interface{}, dest interface{}, opts *RequestOptions) error {
	bodyBytes, err := sonic.Marshal(body)
	if err != nil {
		return err
	}

	data, err := c.Post(url, bodyBytes, opts)
	if err != nil {
		return err
	}

	if dest != nil {
		return sonic.Unmarshal(data, dest)
	}
	return nil
}

// CLOB returns the CLOB API URL
func (c *Client) CLOB(path string) string {
	return c.clobURL + path
}

// Gamma returns the Gamma API URL
func (c *Client) Gamma(path string) string {
	return c.gammaURL + path
}

// Data returns the Data API URL
func (c *Client) Data(path string) string {
	return c.dataURL + path
}

// Close cleans up client resources
func (c *Client) Close() {
	c.httpClient.CloseIdleConnections()
}
