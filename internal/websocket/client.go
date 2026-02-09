package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
)

const (
	// Kalshi spec: "Ping frames every 10 seconds, respond with Pong"
	defaultPingInterval       = 10 * time.Second
	defaultReconnectBaseDelay = 1 * time.Second
	defaultReconnectMaxDelay  = 60 * time.Second
	defaultWriteTimeout       = 10 * time.Second
	defaultReadTimeout        = 60 * time.Second
)

// ClientOptions contains configuration for the WebSocket client
type ClientOptions struct {
	URL                string
	APIKeyID           string
	Signature          string
	Timestamp          string
	PingInterval       time.Duration
	ReconnectBaseDelay time.Duration
	ReconnectMaxDelay  time.Duration
	WriteTimeout       time.Duration
	ReadTimeout        time.Duration
}

// Validate checks that required options are set
func (o *ClientOptions) Validate() error {
	if o.URL == "" {
		return errors.New("URL is required")
	}
	if o.APIKeyID == "" {
		return errors.New("APIKeyID is required")
	}
	if o.Signature == "" {
		return errors.New("Signature is required")
	}
	if o.Timestamp == "" {
		return errors.New("Timestamp is required")
	}
	return nil
}

// Client manages a WebSocket connection to Kalshi
type Client struct {
	url       string
	apiKeyID  string
	signature string
	timestamp string

	conn          *websocket.Conn
	connected     atomic.Bool
	subscriptions *SubscriptionManager
	router        *MessageRouter

	pingInterval       time.Duration
	reconnectBaseDelay time.Duration
	reconnectMaxDelay  time.Duration
	writeTimeout       time.Duration
	readTimeout        time.Duration

	pendingResponses map[int]chan *Message
	pendingMu        sync.RWMutex
	nextPingID       int

	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	closeMu    sync.Mutex
	closed     bool

	onReconnect func()
	onError     func(error)
}

// NewClient creates a new WebSocket client
func NewClient(opts ClientOptions) *Client {
	pingInterval := opts.PingInterval
	if pingInterval == 0 {
		pingInterval = defaultPingInterval
	}

	reconnectBaseDelay := opts.ReconnectBaseDelay
	if reconnectBaseDelay == 0 {
		reconnectBaseDelay = defaultReconnectBaseDelay
	}

	reconnectMaxDelay := opts.ReconnectMaxDelay
	if reconnectMaxDelay == 0 {
		reconnectMaxDelay = defaultReconnectMaxDelay
	}

	writeTimeout := opts.WriteTimeout
	if writeTimeout == 0 {
		writeTimeout = defaultWriteTimeout
	}

	readTimeout := opts.ReadTimeout
	if readTimeout == 0 {
		readTimeout = defaultReadTimeout
	}

	return &Client{
		url:                opts.URL,
		apiKeyID:           opts.APIKeyID,
		signature:          opts.Signature,
		timestamp:          opts.Timestamp,
		subscriptions:      NewSubscriptionManager(),
		router:             NewMessageRouter(),
		pingInterval:       pingInterval,
		reconnectBaseDelay: reconnectBaseDelay,
		reconnectMaxDelay:  reconnectMaxDelay,
		writeTimeout:       writeTimeout,
		readTimeout:        readTimeout,
		pendingResponses:   make(map[int]chan *Message),
		nextPingID:         1000, // Start ping IDs at 1000 to avoid conflicts
	}
}

// Connect establishes a WebSocket connection and authenticates
func (c *Client) Connect(ctx context.Context) error {
	if err := c.connect(ctx); err != nil {
		return err
	}

	// Create a context for the background goroutines
	bgCtx, cancel := context.WithCancel(context.Background())
	c.cancelFunc = cancel

	// Start read loop
	c.wg.Add(1)
	go c.readLoop(bgCtx)

	// Start ping loop
	c.wg.Add(1)
	go c.pingLoop(bgCtx)

	return nil
}

// connect performs the actual connection and authentication
func (c *Client) connect(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to dial WebSocket: %w", err)
	}

	c.conn = conn

	// Authenticate
	authCmd := BuildAuthCommand(c.apiKeyID, c.signature, c.timestamp)
	if err := c.sendCommand(ctx, authCmd); err != nil {
		conn.Close(websocket.StatusNormalClosure, "auth failed")
		return fmt.Errorf("failed to send auth command: %w", err)
	}

	// Wait for auth response
	respChan := c.registerPendingResponse(authCmd.ID)
	defer c.unregisterPendingResponse(authCmd.ID)

	// Read auth response
	_, data, err := conn.Read(ctx)
	if err != nil {
		conn.Close(websocket.StatusNormalClosure, "read failed")
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	msg, err := ParseMessage(data)
	if err != nil {
		conn.Close(websocket.StatusNormalClosure, "parse failed")
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	if msg.Type == "error" {
		conn.Close(websocket.StatusNormalClosure, "auth error")
		return errors.New("authentication failed")
	}

	// Send to pending response channel if waiting
	select {
	case respChan <- msg:
	default:
	}

	c.connected.Store(true)
	return nil
}

// readLoop continuously reads messages from the WebSocket
func (c *Client) readLoop(ctx context.Context) {
	defer c.wg.Done()

	reconnectAttempt := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !c.IsConnected() {
			// Attempt reconnect
			delay := calculateBackoff(reconnectAttempt, c.reconnectBaseDelay, c.reconnectMaxDelay)
			reconnectAttempt++

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}

			if err := c.reconnect(ctx); err != nil {
				if c.onError != nil {
					c.onError(fmt.Errorf("reconnect failed: %w", err))
				}
				continue
			}

			reconnectAttempt = 0
			if c.onReconnect != nil {
				c.onReconnect()
			}
		}

		_, data, err := c.conn.Read(ctx)
		if err != nil {
			c.connected.Store(false)
			if c.onError != nil {
				c.onError(fmt.Errorf("read error: %w", err))
			}
			continue
		}

		msg, err := ParseMessage(data)
		if err != nil {
			if c.onError != nil {
				c.onError(fmt.Errorf("parse error: %w", err))
			}
			continue
		}

		c.handleMessage(msg)
	}
}

// handleMessage processes an incoming message
func (c *Client) handleMessage(msg *Message) {
	// Check if this is a response to a pending command
	if msg.ID > 0 {
		c.pendingMu.RLock()
		respChan, ok := c.pendingResponses[msg.ID]
		c.pendingMu.RUnlock()

		if ok {
			select {
			case respChan <- msg:
			default:
			}
			return
		}
	}

	// Route to channel handler
	if msg.Channel != "" {
		if err := c.router.Route(*msg); err != nil {
			if c.onError != nil {
				c.onError(fmt.Errorf("handler error: %w", err))
			}
		}
	}
}

// pingLoop sends periodic ping messages
func (c *Client) pingLoop(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !c.IsConnected() {
				continue
			}

			c.nextPingID++
			pingCmd := BuildPingCommand(c.nextPingID)

			writeCtx, cancel := context.WithTimeout(ctx, c.writeTimeout)
			err := c.sendCommand(writeCtx, pingCmd)
			cancel()

			if err != nil {
				if c.onError != nil {
					c.onError(fmt.Errorf("ping error: %w", err))
				}
			}
		}
	}
}

// reconnect attempts to re-establish the connection
func (c *Client) reconnect(ctx context.Context) error {
	// Store current subscriptions
	subs := c.subscriptions.GetSubscriptions()

	// Clear subscription tracking (will be restored)
	c.subscriptions.Clear()

	// Reconnect
	if err := c.connect(ctx); err != nil {
		return err
	}

	// Restore subscriptions
	commands := c.subscriptions.RestoreSubscriptions(subs)
	for _, cmd := range commands {
		if err := c.sendCommand(ctx, cmd); err != nil {
			return fmt.Errorf("failed to restore subscription: %w", err)
		}
	}

	return nil
}

// sendCommand sends a command to the WebSocket server
func (c *Client) sendCommand(ctx context.Context, cmd *Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	writeCtx, cancel := context.WithTimeout(ctx, c.writeTimeout)
	defer cancel()

	return c.conn.Write(writeCtx, websocket.MessageText, data)
}

// Subscribe subscribes to a channel
func (c *Client) Subscribe(ctx context.Context, channel Channel, params map[string]string) error {
	cmd, err := c.subscriptions.Subscribe(channel, params)
	if err != nil {
		return err
	}

	return c.sendCommand(ctx, cmd)
}

// Unsubscribe unsubscribes from a channel
func (c *Client) Unsubscribe(ctx context.Context, channel Channel) error {
	cmd, err := c.subscriptions.Unsubscribe(channel)
	if err != nil {
		return err
	}

	return c.sendCommand(ctx, cmd)
}

// RegisterHandler registers a handler for a channel
func (c *Client) RegisterHandler(channel Channel, handler Handler) {
	c.router.Register(channel, handler)
}

// UnregisterHandler removes a handler for a channel
func (c *Client) UnregisterHandler(channel Channel) {
	c.router.Unregister(channel)
}

// IsConnected returns true if the client is connected
func (c *Client) IsConnected() bool {
	return c.connected.Load()
}

// Close closes the WebSocket connection
func (c *Client) Close() {
	c.closeMu.Lock()
	if c.closed {
		c.closeMu.Unlock()
		return
	}
	c.closed = true
	c.closeMu.Unlock()

	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	c.connected.Store(false)

	if c.conn != nil {
		c.conn.Close(websocket.StatusNormalClosure, "client closed")
	}

	c.wg.Wait()
}

// OnReconnect sets a callback for when the client reconnects
func (c *Client) OnReconnect(fn func()) {
	c.onReconnect = fn
}

// OnError sets a callback for error handling
func (c *Client) OnError(fn func(error)) {
	c.onError = fn
}

// registerPendingResponse creates a channel to receive a response for a command
func (c *Client) registerPendingResponse(id int) chan *Message {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	ch := make(chan *Message, 1)
	c.pendingResponses[id] = ch
	return ch
}

// unregisterPendingResponse removes a pending response channel
func (c *Client) unregisterPendingResponse(id int) {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	delete(c.pendingResponses, id)
}

// calculateBackoff calculates exponential backoff delay
func calculateBackoff(attempt int, base, max time.Duration) time.Duration {
	delay := base
	for i := 0; i < attempt; i++ {
		delay *= 2
		if delay > max {
			return max
		}
	}
	return delay
}

// TickerData represents market ticker data
type TickerData struct {
	Ticker       string `json:"ticker"`
	YesPrice     int    `json:"yes_price"`
	NoPrice      int    `json:"no_price"`
	YesBid       int    `json:"yes_bid"`
	YesAsk       int    `json:"yes_ask"`
	Volume       int    `json:"volume"`
	OpenInterest int    `json:"open_interest"`
}

// OrderbookData represents orderbook data
type OrderbookData struct {
	Ticker  string           `json:"ticker"`
	YesBids []OrderbookLevel `json:"yes_bids"`
	YesAsks []OrderbookLevel `json:"yes_asks"`
	NoBids  []OrderbookLevel `json:"no_bids"`
	NoAsks  []OrderbookLevel `json:"no_asks"`
}

// OrderbookLevel represents a price level
type OrderbookLevel struct {
	Price    int `json:"price"`
	Quantity int `json:"quantity"`
}

// TradeData represents trade data
type TradeData struct {
	TradeID   string `json:"trade_id"`
	Ticker    string `json:"ticker"`
	Price     int    `json:"price"`
	Count     int    `json:"count"`
	TakerSide string `json:"taker_side"`
	Timestamp string `json:"ts"`
}

// OrderUpdateData represents order update data
type OrderUpdateData struct {
	OrderID           string `json:"order_id"`
	Ticker            string `json:"ticker"`
	Status            string `json:"status"`
	Side              string `json:"side"`
	Action            string `json:"action"`
	YesPrice          int    `json:"yes_price"`
	NoPrice           int    `json:"no_price"`
	InitialQuantity   int    `json:"initial_quantity"`
	RemainingQuantity int    `json:"remaining_quantity"`
	FilledQuantity    int    `json:"filled_quantity"`
}

// FillData represents fill data
type FillData struct {
	FillID    string `json:"fill_id"`
	TradeID   string `json:"trade_id"`
	OrderID   string `json:"order_id"`
	Ticker    string `json:"ticker"`
	Side      string `json:"side"`
	Action    string `json:"action"`
	YesPrice  int    `json:"yes_price"`
	NoPrice   int    `json:"no_price"`
	Count     int    `json:"count"`
	IsTaker   bool   `json:"is_taker"`
	Timestamp string `json:"ts"`
}

// PositionData represents position data
type PositionData struct {
	Ticker      string `json:"ticker"`
	Position    int    `json:"position"`
	TotalCost   int    `json:"total_cost"`
	RealizedPnl int    `json:"realized_pnl"`
	Exposure    int    `json:"exposure"`
}

// TickerV2Data represents incremental delta ticker updates (market_ticker_v2)
type TickerV2Data struct {
	Ticker    string `json:"ticker"`
	DeltaType string `json:"delta_type"`
	YesPrice  int    `json:"yes_price"`
	NoPrice   int    `json:"no_price"`
	Delta     int    `json:"delta"`
	Volume    int    `json:"volume"`
}

// OrderGroupUpdateData represents order group lifecycle updates
type OrderGroupUpdateData struct {
	OrderGroupID string `json:"order_group_id"`
	Status       string `json:"status"`
	TotalOrders  int    `json:"total_orders"`
	FilledOrders int    `json:"filled_orders"`
}

// MarketLifecycleData represents market state changes
type MarketLifecycleData struct {
	Ticker    string `json:"ticker"`
	Status    string `json:"status"`
	OldStatus string `json:"old_status"`
}

// CommunicationData represents RFQ/quote notifications
type CommunicationData struct {
	Type     string `json:"type"`
	Ticker   string `json:"ticker"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Side     string `json:"side"`
}

// MultivariateLookupData represents collection lookup notifications
type MultivariateLookupData struct {
	SeriesID    string `json:"series_id"`
	LookupValue string `json:"lookup_value"`
}
