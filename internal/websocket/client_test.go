package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func TestNewClient(t *testing.T) {
	opts := ClientOptions{
		URL:       "wss://demo-api.kalshi.co/trade-api/ws/v2",
		APIKeyID:  "test-key",
		Signature: "test-sig",
		Timestamp: "2024-01-15T12:00:00Z",
	}

	client := NewClient(opts)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.url != opts.URL {
		t.Errorf("expected URL '%s', got '%s'", opts.URL, client.url)
	}

	if client.apiKeyID != opts.APIKeyID {
		t.Errorf("expected apiKeyID '%s', got '%s'", opts.APIKeyID, client.apiKeyID)
	}

	if client.subscriptions == nil {
		t.Error("subscriptions should be initialized")
	}

	if client.router == nil {
		t.Error("router should be initialized")
	}
}

func TestClientOptions_Validate(t *testing.T) {
	tests := []struct {
		name        string
		opts        ClientOptions
		expectError bool
	}{
		{
			name: "valid options",
			opts: ClientOptions{
				URL:       "wss://demo-api.kalshi.co/trade-api/ws/v2",
				APIKeyID:  "test-key",
				Signature: "test-sig",
				Timestamp: "2024-01-15T12:00:00Z",
			},
			expectError: false,
		},
		{
			name: "missing URL",
			opts: ClientOptions{
				APIKeyID:  "test-key",
				Signature: "test-sig",
				Timestamp: "2024-01-15T12:00:00Z",
			},
			expectError: true,
		},
		{
			name: "missing APIKeyID",
			opts: ClientOptions{
				URL:       "wss://demo-api.kalshi.co/trade-api/ws/v2",
				Signature: "test-sig",
				Timestamp: "2024-01-15T12:00:00Z",
			},
			expectError: true,
		},
		{
			name: "missing Signature",
			opts: ClientOptions{
				URL:       "wss://demo-api.kalshi.co/trade-api/ws/v2",
				APIKeyID:  "test-key",
				Timestamp: "2024-01-15T12:00:00Z",
			},
			expectError: true,
		},
		{
			name: "missing Timestamp",
			opts: ClientOptions{
				URL:       "wss://demo-api.kalshi.co/trade-api/ws/v2",
				APIKeyID:  "test-key",
				Signature: "test-sig",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestClient_Connect(t *testing.T) {
	// Create test WebSocket server
	server := newTestWSServer(t, func(conn *websocket.Conn) {
		ctx := context.Background()

		// Expect auth command
		_, data, err := conn.Read(ctx)
		if err != nil {
			t.Logf("server read error: %v", err)
			return
		}

		var cmd Command
		if err := json.Unmarshal(data, &cmd); err != nil {
			t.Logf("unmarshal error: %v", err)
			return
		}

		if cmd.Cmd != CmdAuth {
			t.Errorf("expected auth command, got %s", cmd.Cmd)
		}

		// Send auth response
		response := Message{
			ID:   cmd.ID,
			Type: "response",
		}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:       wsURL,
		APIKeyID:  "test-key",
		Signature: "test-sig",
		Timestamp: "2024-01-15T12:00:00Z",
	}

	client := NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !client.IsConnected() {
		t.Error("client should be connected")
	}

	client.Close()
}

func TestClient_Reconnect(t *testing.T) {
	connectionCount := 0
	var mu sync.Mutex

	server := newTestWSServer(t, func(conn *websocket.Conn) {
		mu.Lock()
		connectionCount++
		count := connectionCount
		mu.Unlock()

		ctx := context.Background()

		// Read auth command
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}

		var cmd Command
		json.Unmarshal(data, &cmd)

		// Send auth response
		response := Message{ID: cmd.ID, Type: "response"}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		// First connection: close immediately to trigger reconnect
		if count == 1 {
			conn.Close(websocket.StatusGoingAway, "test disconnect")
			return
		}

		// Second connection: stay open
		<-ctx.Done()
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:                wsURL,
		APIKeyID:           "test-key",
		Signature:          "test-sig",
		Timestamp:          "2024-01-15T12:00:00Z",
		ReconnectBaseDelay: 10 * time.Millisecond,
		ReconnectMaxDelay:  50 * time.Millisecond,
	}

	client := NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Wait for reconnect
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	finalCount := connectionCount
	mu.Unlock()

	if finalCount < 2 {
		t.Errorf("expected at least 2 connections (initial + reconnect), got %d", finalCount)
	}

	client.Close()
}

func TestClient_Subscribe(t *testing.T) {
	server := newTestWSServer(t, func(conn *websocket.Conn) {
		ctx := context.Background()

		// Handle auth
		_, data, _ := conn.Read(ctx)
		var cmd Command
		json.Unmarshal(data, &cmd)
		response := Message{ID: cmd.ID, Type: "response"}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		// Handle subscribe
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}

		json.Unmarshal(data, &cmd)
		if cmd.Cmd != CmdSubscribe {
			t.Errorf("expected subscribe command, got %s", cmd.Cmd)
		}

		// Send subscribe response
		response = Message{ID: cmd.ID, Type: "response"}
		respData, _ = json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:       wsURL,
		APIKeyID:  "test-key",
		Signature: "test-sig",
		Timestamp: "2024-01-15T12:00:00Z",
	}

	client := NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	err := client.Subscribe(ctx, ChannelMarketTicker, map[string]string{"market_ticker": "BTC-100K"})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	client.Close()
}

func TestClient_Unsubscribe(t *testing.T) {
	server := newTestWSServer(t, func(conn *websocket.Conn) {
		ctx := context.Background()

		// Handle auth
		_, data, _ := conn.Read(ctx)
		var cmd Command
		json.Unmarshal(data, &cmd)
		response := Message{ID: cmd.ID, Type: "response"}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		// Handle subscribe
		_, data, _ = conn.Read(ctx)
		json.Unmarshal(data, &cmd)
		response = Message{ID: cmd.ID, Type: "response"}
		respData, _ = json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		// Handle unsubscribe
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}

		json.Unmarshal(data, &cmd)
		if cmd.Cmd != CmdUnsubscribe {
			t.Errorf("expected unsubscribe command, got %s", cmd.Cmd)
		}

		response = Message{ID: cmd.ID, Type: "response"}
		respData, _ = json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:       wsURL,
		APIKeyID:  "test-key",
		Signature: "test-sig",
		Timestamp: "2024-01-15T12:00:00Z",
	}

	client := NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Subscribe first
	if err := client.Subscribe(ctx, ChannelMarketTicker, map[string]string{"market_ticker": "BTC-100K"}); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Then unsubscribe
	err := client.Unsubscribe(ctx, ChannelMarketTicker)
	if err != nil {
		t.Fatalf("Unsubscribe failed: %v", err)
	}

	client.Close()
}

func TestClient_RegisterHandler(t *testing.T) {
	called := make(chan bool, 1)

	server := newTestWSServer(t, func(conn *websocket.Conn) {
		ctx := context.Background()

		// Handle auth
		_, data, _ := conn.Read(ctx)
		var cmd Command
		json.Unmarshal(data, &cmd)
		response := Message{ID: cmd.ID, Type: "response"}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		// Send a ticker message
		msg := Message{
			Type:    "ticker",
			Channel: ChannelMarketTicker,
			Data:    json.RawMessage(`{"ticker":"BTC-100K","yes_price":55}`),
		}
		msgData, _ := json.Marshal(msg)
		conn.Write(ctx, websocket.MessageText, msgData)

		// Keep connection open
		<-ctx.Done()
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:       wsURL,
		APIKeyID:  "test-key",
		Signature: "test-sig",
		Timestamp: "2024-01-15T12:00:00Z",
	}

	client := NewClient(opts)

	// Register handler before connecting
	client.RegisterHandler(ChannelMarketTicker, &MockHandler{
		onMessage: func(msg Message) error {
			called <- true
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	select {
	case <-called:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("handler was not called within timeout")
	}

	client.Close()
}

func TestClient_Ping(t *testing.T) {
	pingReceived := make(chan bool, 1)

	server := newTestWSServer(t, func(conn *websocket.Conn) {
		ctx := context.Background()

		// Handle auth
		_, data, _ := conn.Read(ctx)
		var cmd Command
		json.Unmarshal(data, &cmd)
		response := Message{ID: cmd.ID, Type: "response"}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		// Wait for ping
		for {
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}

			var cmd Command
			if err := json.Unmarshal(data, &cmd); err != nil {
				continue
			}

			if cmd.Cmd == CmdPing {
				pingReceived <- true
				// Send pong response
				response := Message{ID: cmd.ID, Type: "pong"}
				respData, _ := json.Marshal(response)
				conn.Write(ctx, websocket.MessageText, respData)
			}
		}
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:          wsURL,
		APIKeyID:     "test-key",
		Signature:    "test-sig",
		Timestamp:    "2024-01-15T12:00:00Z",
		PingInterval: 100 * time.Millisecond,
	}

	client := NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	select {
	case <-pingReceived:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("ping was not received within timeout")
	}

	client.Close()
}

func TestClient_Close(t *testing.T) {
	server := newTestWSServer(t, func(conn *websocket.Conn) {
		ctx := context.Background()

		// Handle auth
		_, data, _ := conn.Read(ctx)
		var cmd Command
		json.Unmarshal(data, &cmd)
		response := Message{ID: cmd.ID, Type: "response"}
		respData, _ := json.Marshal(response)
		conn.Write(ctx, websocket.MessageText, respData)

		<-ctx.Done()
	})
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	opts := ClientOptions{
		URL:       wsURL,
		APIKeyID:  "test-key",
		Signature: "test-sig",
		Timestamp: "2024-01-15T12:00:00Z",
	}

	client := NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	client.Close()

	if client.IsConnected() {
		t.Error("client should not be connected after Close")
	}
}

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		name     string
		attempt  int
		base     time.Duration
		max      time.Duration
		expected time.Duration
	}{
		{
			name:     "first attempt",
			attempt:  0,
			base:     100 * time.Millisecond,
			max:      10 * time.Second,
			expected: 100 * time.Millisecond,
		},
		{
			name:     "second attempt",
			attempt:  1,
			base:     100 * time.Millisecond,
			max:      10 * time.Second,
			expected: 200 * time.Millisecond,
		},
		{
			name:     "third attempt",
			attempt:  2,
			base:     100 * time.Millisecond,
			max:      10 * time.Second,
			expected: 400 * time.Millisecond,
		},
		{
			name:     "capped at max",
			attempt:  10,
			base:     100 * time.Millisecond,
			max:      1 * time.Second,
			expected: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateBackoff(tt.attempt, tt.base, tt.max)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function to create test WebSocket server
func newTestWSServer(t *testing.T, handler func(conn *websocket.Conn)) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Logf("websocket accept error: %v", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "")

		handler(conn)
	}))
}
