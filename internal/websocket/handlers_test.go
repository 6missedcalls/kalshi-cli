package websocket

import (
	"encoding/json"
	"testing"
)

func TestMessageHandler_HandleMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     Message
		expectCall  bool
		expectError bool
	}{
		{
			name: "market_ticker message",
			message: Message{
				Type:    "ticker",
				Channel: ChannelMarketTicker,
				Data:    json.RawMessage(`{"ticker":"BTC-100K","yes_price":50}`),
			},
			expectCall:  true,
			expectError: false,
		},
		{
			name: "orderbook message",
			message: Message{
				Type:    "orderbook_snapshot",
				Channel: ChannelOrderbook,
				Data:    json.RawMessage(`{"ticker":"BTC-100K"}`),
			},
			expectCall:  true,
			expectError: false,
		},
		{
			name: "public_trades message",
			message: Message{
				Type:    "trade",
				Channel: ChannelPublicTrades,
				Data:    json.RawMessage(`{"ticker":"BTC-100K","price":50}`),
			},
			expectCall:  true,
			expectError: false,
		},
		{
			name: "user_orders message",
			message: Message{
				Type:    "order_update",
				Channel: ChannelUserOrders,
				Data:    json.RawMessage(`{"order_id":"abc123"}`),
			},
			expectCall:  true,
			expectError: false,
		},
		{
			name: "user_fills message",
			message: Message{
				Type:    "fill",
				Channel: ChannelUserFills,
				Data:    json.RawMessage(`{"fill_id":"xyz789"}`),
			},
			expectCall:  true,
			expectError: false,
		},
		{
			name: "unknown channel message",
			message: Message{
				Type:    "unknown",
				Channel: "unknown_channel",
				Data:    json.RawMessage(`{}`),
			},
			expectCall:  false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			handler := &MockHandler{
				onMessage: func(msg Message) error {
					called = true
					return nil
				},
			}

			router := NewMessageRouter()
			if tt.expectCall {
				router.Register(tt.message.Channel, handler)
			}

			err := router.Route(tt.message)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectCall && !called {
				t.Error("expected handler to be called")
			}
			if !tt.expectCall && called {
				t.Error("handler should not have been called")
			}
		})
	}
}

func TestMessageRouter_Register(t *testing.T) {
	router := NewMessageRouter()

	handler := &MockHandler{}
	router.Register(ChannelMarketTicker, handler)

	if len(router.handlers) != 1 {
		t.Errorf("expected 1 handler, got %d", len(router.handlers))
	}

	if router.handlers[ChannelMarketTicker] != handler {
		t.Error("handler not registered correctly")
	}
}

func TestMessageRouter_Unregister(t *testing.T) {
	router := NewMessageRouter()

	handler := &MockHandler{}
	router.Register(ChannelMarketTicker, handler)
	router.Unregister(ChannelMarketTicker)

	if len(router.handlers) != 0 {
		t.Errorf("expected 0 handlers after unregister, got %d", len(router.handlers))
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
		validate    func(t *testing.T, msg *Message)
	}{
		{
			name: "valid ticker message",
			data: []byte(`{"type":"ticker","channel":"ticker","data":{"ticker":"BTC-100K"}}`),
			validate: func(t *testing.T, msg *Message) {
				if msg.Type != "ticker" {
					t.Errorf("expected type 'ticker', got '%s'", msg.Type)
				}
				if msg.Channel != ChannelMarketTicker {
					t.Errorf("expected channel '%s', got '%s'", ChannelMarketTicker, msg.Channel)
				}
			},
		},
		{
			name: "valid command response",
			data: []byte(`{"id":1,"type":"response","msg":{"channels":["ticker"]}}`),
			validate: func(t *testing.T, msg *Message) {
				if msg.ID != 1 {
					t.Errorf("expected id 1, got %d", msg.ID)
				}
				if msg.Type != "response" {
					t.Errorf("expected type 'response', got '%s'", msg.Type)
				}
			},
		},
		{
			name:        "invalid json",
			data:        []byte(`{invalid`),
			expectError: true,
		},
		{
			name: "error response",
			data: []byte(`{"id":1,"type":"error","msg":{"error":"authentication failed"}}`),
			validate: func(t *testing.T, msg *Message) {
				if msg.Type != "error" {
					t.Errorf("expected type 'error', got '%s'", msg.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseMessage(tt.data)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && tt.validate != nil {
				tt.validate(t, msg)
			}
		})
	}
}

// MockHandler implements Handler interface for testing
type MockHandler struct {
	onMessage func(msg Message) error
}

func (m *MockHandler) HandleMessage(msg Message) error {
	if m.onMessage != nil {
		return m.onMessage(msg)
	}
	return nil
}
