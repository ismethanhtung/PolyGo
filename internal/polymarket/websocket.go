package polymarket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/polygo/internal/config"
)

// WSMessageType represents WebSocket message types
type WSMessageType string

const (
	WSMessageTypeSubscribe   WSMessageType = "subscribe"
	WSMessageTypeUnsubscribe WSMessageType = "unsubscribe"
	WSMessageTypePing        WSMessageType = "ping"
	WSMessageTypePong        WSMessageType = "pong"
)

// WSChannel represents WebSocket channel types
type WSChannel string

const (
	WSChannelMarket WSChannel = "market"
	WSChannelUser   WSChannel = "user"
	WSChannelPrice  WSChannel = "price"
	WSChannelTrade  WSChannel = "trade"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      WSMessageType   `json:"type"`
	Channel   WSChannel       `json:"channel,omitempty"`
	Markets   []string        `json:"markets,omitempty"`
	Assets    []string        `json:"assets,omitempty"`
	Auth      *WSAuth         `json:"auth,omitempty"`
	Data      sonic.RawMessage `json:"data,omitempty"`
	Timestamp int64           `json:"timestamp,omitempty"`
}

// WSAuth represents WebSocket authentication
type WSAuth struct {
	APIKey    string `json:"apiKey"`
	Secret    string `json:"secret"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Signature  string `json:"signature"`
}

// WSManager manages WebSocket connections to Polymarket
type WSManager struct {
	config     *config.PolymarketConfig
	clobConn   *websocket.Conn
	liveConn   *websocket.Conn
	mu         sync.RWMutex
	
	// Subscriptions
	marketSubs map[string][]chan []byte
	userSubs   map[string]chan []byte
	
	// Callbacks
	onMessage  func(channel WSChannel, data []byte)
	onError    func(err error)
	onConnect  func()
	onDisconnect func()
	
	// State
	connected  bool
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// NewWSManager creates a new WebSocket manager
func NewWSManager(cfg *config.PolymarketConfig) *WSManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WSManager{
		config:     cfg,
		marketSubs: make(map[string][]chan []byte),
		userSubs:   make(map[string]chan []byte),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// SetCallbacks sets WebSocket callbacks
func (w *WSManager) SetCallbacks(onMessage func(WSChannel, []byte), onError func(error), onConnect, onDisconnect func()) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	w.onMessage = onMessage
	w.onError = onError
	w.onConnect = onConnect
	w.onDisconnect = onDisconnect
}

// Connect establishes WebSocket connections
func (w *WSManager) Connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	// Connect to CLOB WebSocket
	clobConn, _, err := websocket.DefaultDialer.DialContext(w.ctx, w.config.WsClobURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to CLOB WebSocket: %w", err)
	}
	w.clobConn = clobConn
	
	// Connect to Live Data WebSocket
	liveConn, _, err := websocket.DefaultDialer.DialContext(w.ctx, w.config.WsLiveDataURL, nil)
	if err != nil {
		w.clobConn.Close()
		return fmt.Errorf("failed to connect to Live Data WebSocket: %w", err)
	}
	w.liveConn = liveConn
	
	w.connected = true
	
	// Start message handlers
	w.wg.Add(2)
	go w.handleClobMessages()
	go w.handleLiveMessages()
	
	// Start ping routine
	w.wg.Add(1)
	go w.pingRoutine()
	
	if w.onConnect != nil {
		w.onConnect()
	}
	
	return nil
}

// handleClobMessages handles messages from CLOB WebSocket
func (w *WSManager) handleClobMessages() {
	defer w.wg.Done()
	
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			_, message, err := w.clobConn.ReadMessage()
			if err != nil {
				if w.onError != nil {
					w.onError(err)
				}
				w.reconnect()
				return
			}
			
			w.processMessage(WSChannelMarket, message)
		}
	}
}

// handleLiveMessages handles messages from Live Data WebSocket
func (w *WSManager) handleLiveMessages() {
	defer w.wg.Done()
	
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			_, message, err := w.liveConn.ReadMessage()
			if err != nil {
				if w.onError != nil {
					w.onError(err)
				}
				return
			}
			
			w.processMessage(WSChannelPrice, message)
		}
	}
}

// processMessage processes incoming WebSocket messages
func (w *WSManager) processMessage(channel WSChannel, data []byte) {
	if w.onMessage != nil {
		w.onMessage(channel, data)
	}
	
	// Parse message to route to subscribers
	var msg WSMessage
	if err := sonic.Unmarshal(data, &msg); err != nil {
		return
	}
	
	w.mu.RLock()
	defer w.mu.RUnlock()
	
	// Route to market subscribers
	if len(msg.Markets) > 0 {
		for _, market := range msg.Markets {
			if subs, ok := w.marketSubs[market]; ok {
				for _, ch := range subs {
					select {
					case ch <- data:
					default:
						// Channel full, skip
					}
				}
			}
		}
	}
}

// pingRoutine sends periodic pings to keep connection alive
func (w *WSManager) pingRoutine() {
	defer w.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.mu.RLock()
			if w.clobConn != nil {
				ping := WSMessage{Type: WSMessageTypePing, Timestamp: time.Now().UnixMilli()}
				data, _ := sonic.Marshal(ping)
				w.clobConn.WriteMessage(websocket.TextMessage, data)
			}
			w.mu.RUnlock()
		}
	}
}

// reconnect attempts to reconnect WebSocket
func (w *WSManager) reconnect() {
	w.mu.Lock()
	w.connected = false
	w.mu.Unlock()
	
	if w.onDisconnect != nil {
		w.onDisconnect()
	}
	
	// Attempt reconnection with exponential backoff
	backoff := time.Second
	maxBackoff := 30 * time.Second
	
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-time.After(backoff):
			if err := w.Connect(); err != nil {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			}
			return
		}
	}
}

// SubscribeMarket subscribes to market updates
func (w *WSManager) SubscribeMarket(marketID string) (chan []byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	ch := make(chan []byte, 100)
	w.marketSubs[marketID] = append(w.marketSubs[marketID], ch)
	
	// Send subscribe message
	msg := WSMessage{
		Type:    WSMessageTypeSubscribe,
		Channel: WSChannelMarket,
		Markets: []string{marketID},
	}
	
	data, err := sonic.Marshal(msg)
	if err != nil {
		return nil, err
	}
	
	if w.clobConn != nil {
		if err := w.clobConn.WriteMessage(websocket.TextMessage, data); err != nil {
			return nil, err
		}
	}
	
	return ch, nil
}

// UnsubscribeMarket unsubscribes from market updates
func (w *WSManager) UnsubscribeMarket(marketID string, ch chan []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	if subs, ok := w.marketSubs[marketID]; ok {
		for i, sub := range subs {
			if sub == ch {
				w.marketSubs[marketID] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
		
		// If no more subscribers, unsubscribe from server
		if len(w.marketSubs[marketID]) == 0 {
			delete(w.marketSubs, marketID)
			
			msg := WSMessage{
				Type:    WSMessageTypeUnsubscribe,
				Channel: WSChannelMarket,
				Markets: []string{marketID},
			}
			
			data, _ := sonic.Marshal(msg)
			if w.clobConn != nil {
				w.clobConn.WriteMessage(websocket.TextMessage, data)
			}
		}
	}
}

// SubscribeUser subscribes to user updates (requires authentication)
func (w *WSManager) SubscribeUser(userID string, auth *WSAuth) (chan []byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	ch := make(chan []byte, 100)
	w.userSubs[userID] = ch
	
	msg := WSMessage{
		Type:    WSMessageTypeSubscribe,
		Channel: WSChannelUser,
		Auth:    auth,
	}
	
	data, err := sonic.Marshal(msg)
	if err != nil {
		return nil, err
	}
	
	if w.clobConn != nil {
		if err := w.clobConn.WriteMessage(websocket.TextMessage, data); err != nil {
			return nil, err
		}
	}
	
	return ch, nil
}

// Close closes all WebSocket connections
func (w *WSManager) Close() {
	w.cancel()
	
	w.mu.Lock()
	if w.clobConn != nil {
		w.clobConn.Close()
	}
	if w.liveConn != nil {
		w.liveConn.Close()
	}
	w.connected = false
	w.mu.Unlock()
	
	w.wg.Wait()
	
	// Close all subscriber channels
	for _, subs := range w.marketSubs {
		for _, ch := range subs {
			close(ch)
		}
	}
	for _, ch := range w.userSubs {
		close(ch)
	}
}

// IsConnected returns connection status
func (w *WSManager) IsConnected() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.connected
}
