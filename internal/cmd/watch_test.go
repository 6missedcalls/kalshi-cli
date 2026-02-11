package cmd

import (
	"testing"

	"github.com/6missedcalls/kalshi-cli/internal/websocket"
)

func TestRequiresAuth(t *testing.T) {
	tests := []struct {
		name     string
		channels []websocket.Channel
		expected bool
	}{
		{
			name:     "public channel market_ticker",
			channels: []websocket.Channel{websocket.ChannelMarketTicker},
			expected: false,
		},
		{
			name:     "public channel market_ticker_v2",
			channels: []websocket.Channel{websocket.ChannelMarketTickerV2},
			expected: false,
		},
		{
			name:     "public channel public_trades",
			channels: []websocket.Channel{websocket.ChannelPublicTrades},
			expected: false,
		},
		{
			name:     "public channel market_lifecycle",
			channels: []websocket.Channel{websocket.ChannelMarketLifecycle},
			expected: false,
		},
		{
			name:     "public channel multivariate_lookups",
			channels: []websocket.Channel{websocket.ChannelMultivariateLookups},
			expected: false,
		},
		{
			name:     "auth channel orderbook",
			channels: []websocket.Channel{websocket.ChannelOrderbook},
			expected: true,
		},
		{
			name:     "auth channel user_orders",
			channels: []websocket.Channel{websocket.ChannelUserOrders},
			expected: true,
		},
		{
			name:     "auth channel user_fills",
			channels: []websocket.Channel{websocket.ChannelUserFills},
			expected: true,
		},
		{
			name:     "auth channel market_positions",
			channels: []websocket.Channel{websocket.ChannelMarketPositions},
			expected: true,
		},
		{
			name:     "auth channel order_group_updates",
			channels: []websocket.Channel{websocket.ChannelOrderGroupUpdates},
			expected: true,
		},
		{
			name:     "auth channel communications",
			channels: []websocket.Channel{websocket.ChannelCommunications},
			expected: true,
		},
		{
			name:     "mixed channels with one auth",
			channels: []websocket.Channel{websocket.ChannelMarketTicker, websocket.ChannelUserOrders},
			expected: true,
		},
		{
			name:     "multiple public channels",
			channels: []websocket.Channel{websocket.ChannelMarketTicker, websocket.ChannelPublicTrades},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := requiresAuth(tt.channels)
			if result != tt.expected {
				t.Errorf("requiresAuth(%v) = %v, expected %v", tt.channels, result, tt.expected)
			}
		})
	}
}

// TestOrderbookRequiresAuth specifically tests the bug fix
// where orderbook was not being flagged as requiring auth
func TestOrderbookRequiresAuth(t *testing.T) {
	channels := []websocket.Channel{websocket.ChannelOrderbook}
	if !requiresAuth(channels) {
		t.Error("orderbook channel should require authentication per Kalshi API spec")
	}
}

// TestMarketPositionsChannelUsed verifies positions command uses correct channel
func TestMarketPositionsChannelUsed(t *testing.T) {
	channel := websocket.ChannelMarketPositions
	if channel != "market_positions" {
		t.Errorf("expected market_positions channel, got %s", channel)
	}
}

func TestChannelNamesMatchKalshiAPI(t *testing.T) {
	// Document correct Kalshi WebSocket v2 channel names
	tests := []struct {
		name     string
		channel  websocket.Channel
		expected string
	}{
		{"ticker", websocket.ChannelMarketTicker, "ticker"},
		{"ticker_v2", websocket.ChannelMarketTickerV2, "ticker_v2"},
		{"orderbook_delta", websocket.ChannelOrderbook, "orderbook_delta"},
		{"trade", websocket.ChannelPublicTrades, "trade"},
		{"user_orders", websocket.ChannelUserOrders, "user_orders"},
		{"fill", websocket.ChannelUserFills, "fill"},
		{"market_positions", websocket.ChannelMarketPositions, "market_positions"},
		{"order_group_updates", websocket.ChannelOrderGroupUpdates, "order_group_updates"},
		{"communications", websocket.ChannelCommunications, "communications"},
		{"market_lifecycle_v2", websocket.ChannelMarketLifecycle, "market_lifecycle_v2"},
		{"multivariate", websocket.ChannelMultivariateLookups, "multivariate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.channel) != tt.expected {
				t.Errorf("channel %s = %q, want %q", tt.name, tt.channel, tt.expected)
			}
		})
	}
}
