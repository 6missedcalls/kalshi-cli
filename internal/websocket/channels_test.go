package websocket

import (
	"encoding/json"
	"testing"
)

func TestNewSubscriptionManager(t *testing.T) {
	sm := NewSubscriptionManager()

	if sm == nil {
		t.Fatal("NewSubscriptionManager returned nil")
	}

	if sm.subscriptions == nil {
		t.Error("subscriptions map should be initialized")
	}

	if sm.nextID != 2 {
		t.Errorf("nextID should start at 2 (1 reserved for auth), got %d", sm.nextID)
	}
}

func TestSubscriptionManager_Subscribe(t *testing.T) {
	tests := []struct {
		name    string
		channel Channel
		params  map[string]string
	}{
		{
			name:    "subscribe to market_ticker",
			channel: ChannelMarketTicker,
			params:  map[string]string{"market_ticker": "BTC-100K"},
		},
		{
			name:    "subscribe to orderbook",
			channel: ChannelOrderbook,
			params:  map[string]string{"orderbook": "BTC-100K"},
		},
		{
			name:    "subscribe to public_trades",
			channel: ChannelPublicTrades,
			params:  map[string]string{"public_trades": "BTC-100K"},
		},
		{
			name:    "subscribe to user_orders",
			channel: ChannelUserOrders,
			params:  nil,
		},
		{
			name:    "subscribe to user_fills",
			channel: ChannelUserFills,
			params:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewSubscriptionManager()

			cmd, err := sm.Subscribe(tt.channel, tt.params)
			if err != nil {
				t.Fatalf("Subscribe failed: %v", err)
			}

			if cmd.Cmd != CmdSubscribe {
				t.Errorf("expected cmd '%s', got '%s'", CmdSubscribe, cmd.Cmd)
			}

			if cmd.ID < 2 {
				t.Errorf("command ID should be >= 2, got %d", cmd.ID)
			}

			// Verify channel is in params
			channels, ok := cmd.Params["channels"].([]Channel)
			if !ok {
				t.Fatal("channels not found in params")
			}

			if len(channels) != 1 || channels[0] != tt.channel {
				t.Errorf("expected channel %s, got %v", tt.channel, channels)
			}

			// Verify subscription is tracked
			if !sm.IsSubscribed(tt.channel) {
				t.Error("channel should be marked as subscribed")
			}
		})
	}
}

func TestSubscriptionManager_Unsubscribe(t *testing.T) {
	sm := NewSubscriptionManager()

	// Subscribe first
	_, err := sm.Subscribe(ChannelMarketTicker, map[string]string{"market_ticker": "BTC-100K"})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Then unsubscribe
	cmd, err := sm.Unsubscribe(ChannelMarketTicker)
	if err != nil {
		t.Fatalf("Unsubscribe failed: %v", err)
	}

	if cmd.Cmd != CmdUnsubscribe {
		t.Errorf("expected cmd '%s', got '%s'", CmdUnsubscribe, cmd.Cmd)
	}

	// Verify unsubscribed
	if sm.IsSubscribed(ChannelMarketTicker) {
		t.Error("channel should not be subscribed after unsubscribe")
	}
}

func TestSubscriptionManager_UnsubscribeNotSubscribed(t *testing.T) {
	sm := NewSubscriptionManager()

	_, err := sm.Unsubscribe(ChannelMarketTicker)
	if err == nil {
		t.Error("expected error when unsubscribing from non-subscribed channel")
	}
}

func TestSubscriptionManager_GetSubscriptions(t *testing.T) {
	sm := NewSubscriptionManager()

	sm.Subscribe(ChannelMarketTicker, map[string]string{"market_ticker": "BTC-100K"})
	sm.Subscribe(ChannelOrderbook, map[string]string{"orderbook": "ETH-5K"})

	subs := sm.GetSubscriptions()

	if len(subs) != 2 {
		t.Errorf("expected 2 subscriptions, got %d", len(subs))
	}
}

func TestCommand_MarshalJSON(t *testing.T) {
	cmd := Command{
		ID:  1,
		Cmd: CmdSubscribe,
		Params: map[string]interface{}{
			"channels":      []Channel{ChannelMarketTicker},
			"market_ticker": "BTC-100K",
		},
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Parse back to verify structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if parsed["id"].(float64) != 1 {
		t.Errorf("expected id 1, got %v", parsed["id"])
	}

	if parsed["cmd"].(string) != string(CmdSubscribe) {
		t.Errorf("expected cmd '%s', got %s", CmdSubscribe, parsed["cmd"])
	}

	params := parsed["params"].(map[string]interface{})
	if params["market_ticker"].(string) != "BTC-100K" {
		t.Errorf("expected market_ticker 'BTC-100K', got %v", params["market_ticker"])
	}
}

func TestBuildAuthCommand(t *testing.T) {
	apiKeyID := "test-key-id"
	signature := "test-signature"
	timestamp := "2024-01-15T12:00:00Z"

	cmd := BuildAuthCommand(apiKeyID, signature, timestamp)

	if cmd.ID != AuthCommandID {
		t.Errorf("auth command ID should be %d, got %d", AuthCommandID, cmd.ID)
	}

	if cmd.Cmd != CmdAuth {
		t.Errorf("expected cmd '%s', got '%s'", CmdAuth, cmd.Cmd)
	}

	params := cmd.Params
	if params["api_key"].(string) != apiKeyID {
		t.Errorf("expected api_key '%s', got '%v'", apiKeyID, params["api_key"])
	}

	if params["signature"].(string) != signature {
		t.Errorf("expected signature '%s', got '%v'", signature, params["signature"])
	}

	if params["timestamp"].(string) != timestamp {
		t.Errorf("expected timestamp '%s', got '%v'", timestamp, params["timestamp"])
	}
}

func TestSubscription_Restore(t *testing.T) {
	sm := NewSubscriptionManager()

	// Subscribe to multiple channels
	sm.Subscribe(ChannelMarketTicker, map[string]string{"market_ticker": "BTC-100K"})
	sm.Subscribe(ChannelOrderbook, map[string]string{"orderbook": "BTC-100K"})

	// Get subscriptions for restore
	subs := sm.GetSubscriptions()

	// Create new manager and restore
	newSM := NewSubscriptionManager()
	commands := newSM.RestoreSubscriptions(subs)

	if len(commands) != 2 {
		t.Errorf("expected 2 restore commands, got %d", len(commands))
	}

	for _, cmd := range commands {
		if cmd.Cmd != CmdSubscribe {
			t.Errorf("restore command should be subscribe, got %s", cmd.Cmd)
		}
	}
}
