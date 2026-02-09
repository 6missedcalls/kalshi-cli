package websocket

import (
	"testing"
	"time"
)

// TDD RED PHASE: Tests for missing channels and spec compliance issues
// These tests document the expected API behavior from Kalshi's WebSocket spec

// Test 1: All Kalshi channels should be defined
func TestAllKalshiChannelsDefined(t *testing.T) {
	expectedChannels := []struct {
		name         string
		channel      Channel
		authRequired bool
	}{
		{"market_ticker", ChannelMarketTicker, false},
		{"market_ticker_v2", ChannelMarketTickerV2, false},
		{"public_trades", ChannelPublicTrades, false},
		{"orderbook", ChannelOrderbook, true},
		{"user_orders", ChannelUserOrders, true},
		{"user_fills", ChannelUserFills, true},
		{"market_positions", ChannelMarketPositions, true},
		{"order_group_updates", ChannelOrderGroupUpdates, true},
		{"communications", ChannelCommunications, true},
		{"market_lifecycle", ChannelMarketLifecycle, false},
		{"multivariate_lookups", ChannelMultivariateLookups, false},
	}

	for _, tc := range expectedChannels {
		t.Run(tc.name, func(t *testing.T) {
			if string(tc.channel) != tc.name {
				t.Errorf("expected channel constant for %s, got %s", tc.name, tc.channel)
			}
		})
	}
}

// Test 2: Ping interval should be 10 seconds per Kalshi spec
func TestPingIntervalMatchesKalshiSpec(t *testing.T) {
	// Kalshi docs specify: "Ping frames every 10 seconds, respond with Pong"
	expectedInterval := 10 * time.Second

	if defaultPingInterval != expectedInterval {
		t.Errorf("ping interval should be %v (per Kalshi spec), got %v", expectedInterval, defaultPingInterval)
	}
}

// Test 3: Channel requires authentication check should include orderbook
func TestChannelAuthRequirements(t *testing.T) {
	authRequiredChannels := map[Channel]bool{
		ChannelOrderbook:          true,
		ChannelUserOrders:         true,
		ChannelUserFills:          true,
		ChannelMarketPositions:    true,
		ChannelOrderGroupUpdates:  true,
		ChannelCommunications:     true,
		ChannelMarketTicker:       false,
		ChannelMarketTickerV2:     false,
		ChannelPublicTrades:       false,
		ChannelMarketLifecycle:    false,
		ChannelMultivariateLookups: false,
	}

	for channel, expected := range authRequiredChannels {
		t.Run(string(channel), func(t *testing.T) {
			result := ChannelRequiresAuth(channel)
			if result != expected {
				t.Errorf("channel %s auth requirement: expected %v, got %v", channel, expected, result)
			}
		})
	}
}

// Test 4: MarketPositions data structure should exist
func TestMarketPositionsDataStructure(t *testing.T) {
	// Per Kalshi spec, market_positions channel sends position updates
	data := PositionData{
		Ticker:      "BTC-100K",
		Position:    100,
		TotalCost:   5000,
		RealizedPnl: 250,
		Exposure:    1000,
	}

	if data.Ticker != "BTC-100K" {
		t.Errorf("expected ticker BTC-100K, got %s", data.Ticker)
	}
}

// Test 5: OrderGroupUpdate data structure should exist
func TestOrderGroupUpdateDataStructure(t *testing.T) {
	data := OrderGroupUpdateData{
		OrderGroupID: "og-123",
		Status:       "active",
		TotalOrders:  5,
		FilledOrders: 2,
	}

	if data.OrderGroupID != "og-123" {
		t.Errorf("expected order group ID og-123, got %s", data.OrderGroupID)
	}
}

// Test 6: MarketLifecycle data structure should exist
func TestMarketLifecycleDataStructure(t *testing.T) {
	data := MarketLifecycleData{
		Ticker:   "BTC-100K",
		Status:   "active",
		OldStatus: "inactive",
	}

	if data.Status != "active" {
		t.Errorf("expected status active, got %s", data.Status)
	}
}

// Test 7: Communication (RFQ/quote) data structure should exist
func TestCommunicationDataStructure(t *testing.T) {
	data := CommunicationData{
		Type:     "rfq",
		Ticker:   "BTC-100K",
		Quantity: 1000,
	}

	if data.Type != "rfq" {
		t.Errorf("expected type rfq, got %s", data.Type)
	}
}

// Test 8: market_ticker_v2 should send incremental delta updates
func TestMarketTickerV2DataStructure(t *testing.T) {
	data := TickerV2Data{
		Ticker:    "BTC-100K",
		DeltaType: "price_change",
		YesPrice:  55,
		NoPrice:   45,
		Delta:     5,
	}

	if data.DeltaType != "price_change" {
		t.Errorf("expected delta_type price_change, got %s", data.DeltaType)
	}
}

// Test 9: MultivariateLookups data structure should exist
func TestMultivariateLookupDataStructure(t *testing.T) {
	data := MultivariateLookupData{
		SeriesID:    "series-abc",
		LookupValue: "result",
	}

	if data.SeriesID != "series-abc" {
		t.Errorf("expected series ID series-abc, got %s", data.SeriesID)
	}
}

// Test 10: Subscription manager should track all channel types
func TestSubscriptionManagerAllChannels(t *testing.T) {
	allChannels := []Channel{
		ChannelMarketTicker,
		ChannelMarketTickerV2,
		ChannelPublicTrades,
		ChannelOrderbook,
		ChannelUserOrders,
		ChannelUserFills,
		ChannelMarketPositions,
		ChannelOrderGroupUpdates,
		ChannelCommunications,
		ChannelMarketLifecycle,
		ChannelMultivariateLookups,
	}

	sm := NewSubscriptionManager()

	for _, ch := range allChannels {
		t.Run(string(ch), func(t *testing.T) {
			_, err := sm.Subscribe(ch, nil)
			if err != nil {
				t.Errorf("failed to subscribe to %s: %v", ch, err)
			}
			if !sm.IsSubscribed(ch) {
				t.Errorf("channel %s should be subscribed", ch)
			}
		})
	}
}
