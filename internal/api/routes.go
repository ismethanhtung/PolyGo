package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	
	"github.com/polygo/internal/api/handlers"
	"github.com/polygo/internal/api/middleware"
	"github.com/polygo/internal/cache"
	"github.com/polygo/internal/config"
	"github.com/polygo/internal/polymarket"
)

// Server holds all dependencies for the API server
type Server struct {
	app       *fiber.App
	config    *config.Config
	cache     *cache.Cache
	client    *polymarket.Client
	gamma     *polymarket.GammaClient
	clob      *polymarket.ClobClient
	data      *polymarket.DataClient
	wsManager *polymarket.WSManager
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, c *cache.Cache) (*Server, error) {
	// Create Polymarket client
	client := polymarket.NewClient(&cfg.Polymarket, c)
	
	// Create sub-clients
	gamma := polymarket.NewGammaClient(client)
	clob := polymarket.NewClobClient(client)
	data := polymarket.NewDataClient(client)
	
	// Create WebSocket manager
	wsManager := polymarket.NewWSManager(&cfg.Polymarket)
	
	// Create Fiber app with optimized settings
	app := fiber.New(fiber.Config{
		Prefork:               cfg.Server.Prefork,
		ServerHeader:          "PolyGo",
		DisableStartupMessage: !cfg.Server.Debug,
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           cfg.Server.IdleTimeout,
		// Performance optimizations
		DisableDefaultDate:         true,
		DisableHeaderNormalizing:   true,
		DisablePreParseMultipartForm: true,
		StreamRequestBody:          true,
	})
	
	server := &Server{
		app:       app,
		config:    cfg,
		cache:     c,
		client:    client,
		gamma:     gamma,
		clob:      clob,
		data:      data,
		wsManager: wsManager,
	}
	
	// Setup routes
	server.setupMiddleware()
	server.setupRoutes()
	
	return server, nil
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// CORS
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,POLY-API-KEY,POLY-API-SECRET,POLY-PASSPHRASE,POLY-SIGNATURE,POLY-TIMESTAMP",
	}))
	
	// Recovery
	s.app.Use(middleware.Recovery())
	
	// Logger (skip health checks)
	s.app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skip: func(c *fiber.Ctx) bool {
			path := c.Path()
			return path == "/health" || path == "/ready"
		},
	}))
	
	// Rate limiting
	s.app.Use(middleware.RateLimit(middleware.RateLimitConfig{
		Max:    1000,
		Window: 10 * 1000 * 1000 * 1000, // 10 seconds in nanoseconds
		Skip: func(c *fiber.Ctx) bool {
			return c.Path() == "/health" || c.Path() == "/ready"
		},
	}))
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Create handlers
	healthHandler := handlers.NewHealthHandler(s.cache, s.wsManager)
	marketsHandler := handlers.NewMarketsHandler(s.gamma)
	eventsHandler := handlers.NewEventsHandler(s.gamma)
	pricesHandler := handlers.NewPricesHandler(s.clob)
	ordersHandler := handlers.NewOrdersHandler(s.clob, &s.config.Auth)
	dataHandler := handlers.NewDataHandler(s.data)
	wsHandler := handlers.NewWebSocketHandler(s.wsManager)
	
	// Health endpoints
	s.app.Get("/health", healthHandler.Health)
	s.app.Get("/ready", healthHandler.Ready)
	s.app.Get("/stats", healthHandler.Stats)
	
	// Swagger
	s.app.Get("/swagger/*", swagger.HandlerDefault)
	
	// API v1 routes
	v1 := s.app.Group("/api/v1")
	
	// Markets (public)
	markets := v1.Group("/markets")
	markets.Get("/", marketsHandler.GetMarkets)
	markets.Get("/:id", marketsHandler.GetMarket)
	markets.Get("/slug/:slug", marketsHandler.GetMarketBySlug)
	markets.Get("/token/:token_id", marketsHandler.GetMarketByToken)
	
	// Events (public)
	events := v1.Group("/events")
	events.Get("/", eventsHandler.GetEvents)
	events.Get("/search", eventsHandler.SearchEvents)
	events.Get("/:id", eventsHandler.GetEvent)
	events.Get("/slug/:slug", eventsHandler.GetEventBySlug)
	
	// Prices (public)
	v1.Get("/price/:token_id", pricesHandler.GetPrice)
	v1.Get("/prices", pricesHandler.GetPrices)
	v1.Get("/book/:token_id", pricesHandler.GetOrderBook)
	v1.Get("/books", pricesHandler.GetOrderBooks)
	v1.Get("/spread/:token_id", pricesHandler.GetSpread)
	v1.Get("/midpoint/:token_id", pricesHandler.GetMidpoint)
	v1.Get("/midpoints", pricesHandler.GetMidpoints)
	v1.Get("/last-trade/:token_id", pricesHandler.GetLastTradePrice)
	
	// Trades (public)
	v1.Get("/trades/:token_id", ordersHandler.GetTrades)
	v1.Get("/market-trades", dataHandler.GetMarketTrades)
	
	// Price history (public)
	v1.Get("/price-history/:token_id", dataHandler.GetPriceHistory)
	v1.Get("/timeseries", dataHandler.GetTimeseries)
	
	// Top movers & leaderboard (public)
	v1.Get("/top-movers", dataHandler.GetTopMovers)
	v1.Get("/leaderboard", dataHandler.GetLeaderboard)
	
	// User data (public, address-based)
	v1.Get("/positions", dataHandler.GetPositions)
	v1.Get("/positions/market", dataHandler.GetPositionsByMarket)
	v1.Get("/user/trades", dataHandler.GetUserTrades)
	v1.Get("/user/trades/market", dataHandler.GetUserTradesByMarket)
	v1.Get("/activity", dataHandler.GetActivity)
	
	// Orders (authenticated)
	orders := v1.Group("/orders")
	orders.Use(middleware.OptionalAuth(&s.config.Auth))
	
	orders.Get("/", ordersHandler.GetOrders)
	orders.Get("/open", ordersHandler.GetOpenOrders)
	orders.Get("/:id", ordersHandler.GetOrder)
	orders.Post("/", middleware.Auth(&s.config.Auth), ordersHandler.CreateOrder)
	orders.Delete("/:id", middleware.Auth(&s.config.Auth), ordersHandler.CancelOrder)
	orders.Delete("/cancel-all", middleware.Auth(&s.config.Auth), ordersHandler.CancelAllOrders)
	orders.Post("/batch-cancel", middleware.Auth(&s.config.Auth), ordersHandler.CancelOrders)
	
	// WebSocket endpoints
	ws := s.app.Group("/ws")
	ws.Use(handlers.WSMiddleware())
	
	ws.Get("/market/:market_id", websocket.New(wsHandler.HandleMarketWS))
	ws.Get("/markets", websocket.New(wsHandler.HandleAllMarketsWS))
}

// Start starts the server
func (s *Server) Start() error {
	// Connect WebSocket to Polymarket
	go func() {
		if err := s.wsManager.Connect(); err != nil {
			// Log but don't fail - WebSocket is optional
			println("Warning: Failed to connect WebSocket:", err.Error())
		}
	}()
	
	addr := s.config.Server.Host + ":" + itoa(s.config.Server.Port)
	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.wsManager.Close()
	s.client.Close()
	s.cache.Close()
	return s.app.Shutdown()
}

// GetApp returns the Fiber app (for testing)
func (s *Server) GetApp() *fiber.App {
	return s.app
}

// itoa converts int to string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
