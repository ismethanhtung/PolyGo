package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/polygo/internal/polymarket"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	wsManager   *polymarket.WSManager
	clients     map[*websocket.Conn]map[string]bool // client -> subscribed markets
	clientsMu   sync.RWMutex
	broadcast   chan *WSBroadcast
}

// WSBroadcast represents a broadcast message
type WSBroadcast struct {
	MarketID string
	Data     []byte
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(wsManager *polymarket.WSManager) *WebSocketHandler {
	h := &WebSocketHandler{
		wsManager: wsManager,
		clients:   make(map[*websocket.Conn]map[string]bool),
		broadcast: make(chan *WSBroadcast, 1000),
	}
	
	// Setup callbacks from polymarket WebSocket
	wsManager.SetCallbacks(
		func(channel polymarket.WSChannel, data []byte) {
			// Forward to all subscribed clients
			h.handleUpstreamMessage(channel, data)
		},
		func(err error) {
			log.Printf("WebSocket error: %v", err)
		},
		func() {
			log.Println("WebSocket connected to Polymarket")
		},
		func() {
			log.Println("WebSocket disconnected from Polymarket")
		},
	)
	
	// Start broadcast handler
	go h.handleBroadcasts()
	
	return h
}

// handleUpstreamMessage handles messages from Polymarket WebSocket
func (h *WebSocketHandler) handleUpstreamMessage(channel polymarket.WSChannel, data []byte) {
	// Parse message to get market ID
	var msg struct {
		Markets []string `json:"markets"`
		Market  string   `json:"market"`
	}
	
	if err := sonic.Unmarshal(data, &msg); err != nil {
		return
	}
	
	// Broadcast to relevant clients
	markets := msg.Markets
	if msg.Market != "" {
		markets = append(markets, msg.Market)
	}
	
	for _, marketID := range markets {
		h.broadcast <- &WSBroadcast{
			MarketID: marketID,
			Data:     data,
		}
	}
}

// handleBroadcasts processes broadcast messages
func (h *WebSocketHandler) handleBroadcasts() {
	for msg := range h.broadcast {
		h.clientsMu.RLock()
		for conn, subs := range h.clients {
			if subs[msg.MarketID] || subs["*"] {
				go func(c *websocket.Conn, data []byte) {
					if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
						log.Printf("Failed to write to WebSocket: %v", err)
					}
				}(conn, msg.Data)
			}
		}
		h.clientsMu.RUnlock()
	}
}

// UpgradeCheck checks if the request can be upgraded to WebSocket
func (h *WebSocketHandler) UpgradeCheck(c *fiber.Ctx) bool {
	return websocket.IsWebSocketUpgrade(c)
}

// HandleMarketWS handles WebSocket connections for market updates
// @Summary Market WebSocket
// @Description WebSocket endpoint for real-time market updates
// @Tags WebSocket
// @Param market_id path string true "Market ID to subscribe"
// @Router /ws/market/{market_id} [get]
func (h *WebSocketHandler) HandleMarketWS(c *websocket.Conn) {
	marketID := c.Params("market_id")
	
	// Register client
	h.clientsMu.Lock()
	h.clients[c] = map[string]bool{marketID: true}
	h.clientsMu.Unlock()
	
	// Subscribe to market on upstream
	ch, err := h.wsManager.SubscribeMarket(marketID)
	if err != nil {
		log.Printf("Failed to subscribe to market %s: %v", marketID, err)
		c.Close()
		return
	}
	
	// Cleanup on disconnect
	defer func() {
		h.wsManager.UnsubscribeMarket(marketID, ch)
		h.clientsMu.Lock()
		delete(h.clients, c)
		h.clientsMu.Unlock()
		c.Close()
	}()
	
	// Forward messages from upstream
	go func() {
		for data := range ch {
			if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}()
	
	// Handle incoming messages from client
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}
		
		// Parse client message
		var clientMsg struct {
			Type    string   `json:"type"`
			Markets []string `json:"markets"`
		}
		
		if err := sonic.Unmarshal(msg, &clientMsg); err != nil {
			continue
		}
		
		switch clientMsg.Type {
		case "subscribe":
			for _, m := range clientMsg.Markets {
				h.clientsMu.Lock()
				h.clients[c][m] = true
				h.clientsMu.Unlock()
				h.wsManager.SubscribeMarket(m)
			}
		case "unsubscribe":
			for _, m := range clientMsg.Markets {
				h.clientsMu.Lock()
				delete(h.clients[c], m)
				h.clientsMu.Unlock()
			}
		case "ping":
			pong := map[string]interface{}{
				"type":      "pong",
				"timestamp": time.Now().UnixMilli(),
			}
			data, _ := sonic.Marshal(pong)
			c.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// HandleAllMarketsWS handles WebSocket for all market updates
// @Summary All Markets WebSocket
// @Description WebSocket endpoint for all real-time market updates
// @Tags WebSocket
// @Router /ws/markets [get]
func (h *WebSocketHandler) HandleAllMarketsWS(c *websocket.Conn) {
	// Register client for all markets
	h.clientsMu.Lock()
	h.clients[c] = map[string]bool{"*": true}
	h.clientsMu.Unlock()
	
	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, c)
		h.clientsMu.Unlock()
		c.Close()
	}()
	
	// Handle incoming messages
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		
		var clientMsg struct {
			Type string `json:"type"`
		}
		
		if err := sonic.Unmarshal(msg, &clientMsg); err != nil {
			continue
		}
		
		if clientMsg.Type == "ping" {
			pong := map[string]interface{}{
				"type":      "pong",
				"timestamp": time.Now().UnixMilli(),
			}
			data, _ := sonic.Marshal(pong)
			c.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// WSMiddleware returns middleware for WebSocket upgrade check
func WSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
