package websocket

import (
	"encoding/json"
	"sync"
)

// Channel represents a WebSocket channel type
type Channel string

const (
	// Public channels (no auth required)
	ChannelMarketTicker        Channel = "ticker"
	ChannelMarketTickerV2      Channel = "ticker_v2"
	ChannelPublicTrades        Channel = "trade"
	ChannelMarketLifecycle     Channel = "market_lifecycle_v2"
	ChannelMultivariateLookups Channel = "multivariate"

	// Authenticated channels
	ChannelOrderbook         Channel = "orderbook_delta"
	ChannelUserOrders        Channel = "user_orders"
	ChannelUserFills         Channel = "fill"
	ChannelMarketPositions   Channel = "market_positions"
	ChannelOrderGroupUpdates Channel = "order_group_updates"
	ChannelCommunications    Channel = "communications"
)

// ChannelRequiresAuth returns true if the channel requires authentication
func ChannelRequiresAuth(channel Channel) bool {
	switch channel {
	case ChannelOrderbook,
		ChannelUserOrders,
		ChannelUserFills,
		ChannelMarketPositions,
		ChannelOrderGroupUpdates,
		ChannelCommunications:
		return true
	default:
		return false
	}
}

// Message represents a WebSocket message from the server
type Message struct {
	ID      int             `json:"id,omitempty"`
	Type    string          `json:"type,omitempty"`
	Channel Channel         `json:"channel,omitempty"`
	Msg     json.RawMessage `json:"msg,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Handler defines the interface for handling WebSocket messages
type Handler interface {
	HandleMessage(msg Message) error
}

// HandlerFunc is a function adapter for Handler interface
type HandlerFunc func(msg Message) error

// HandleMessage implements Handler interface
func (f HandlerFunc) HandleMessage(msg Message) error {
	return f(msg)
}

// MessageRouter routes messages to appropriate handlers by channel
type MessageRouter struct {
	handlers map[Channel]Handler
	mu       sync.RWMutex
}

// NewMessageRouter creates a new message router
func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		handlers: make(map[Channel]Handler),
	}
}

// Register registers a handler for a specific channel
func (r *MessageRouter) Register(channel Channel, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[channel] = handler
}

// Unregister removes a handler for a specific channel
func (r *MessageRouter) Unregister(channel Channel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, channel)
}

// Route routes a message to the appropriate handler
func (r *MessageRouter) Route(msg Message) error {
	r.mu.RLock()
	handler, ok := r.handlers[msg.Channel]
	r.mu.RUnlock()

	if !ok {
		// No handler registered for this channel, ignore silently
		return nil
	}

	return handler.HandleMessage(msg)
}

// HasHandler returns true if a handler is registered for the channel
func (r *MessageRouter) HasHandler(channel Channel) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.handlers[channel]
	return ok
}

// ParseMessage parses a raw JSON message into a Message struct
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
