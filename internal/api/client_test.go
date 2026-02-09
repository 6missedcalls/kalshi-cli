package api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/6missedcalls/kalshi-cli/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			Production: false,
			Timeout:    30 * time.Second,
		},
	}

	signer := createTestSigner(t)

	client := NewClient(cfg, signer)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClient_NilConfig(t *testing.T) {
	signer := createTestSigner(t)

	client := NewClient(nil, signer)

	if client == nil {
		t.Fatal("NewClient with nil config should still create client with defaults")
	}
}

func TestNewClient_NilSigner(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			Production: false,
			Timeout:    30 * time.Second,
		},
	}

	client := NewClient(cfg, nil)

	if client == nil {
		t.Fatal("NewClient with nil signer should still create client")
	}
}

func TestClient_BaseURL_Demo(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			Production: false,
		},
	}

	client := NewClient(cfg, nil)

	if client.BaseURL() != config.DemoBaseURL {
		t.Errorf("expected demo base URL %s, got %s", config.DemoBaseURL, client.BaseURL())
	}
}

func TestClient_BaseURL_Production(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{
			Production: true,
		},
	}

	client := NewClient(cfg, nil)

	if client.BaseURL() != config.ProdBaseURL {
		t.Errorf("expected production base URL %s, got %s", config.ProdBaseURL, client.BaseURL())
	}
}

func TestClient_SetDebug(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg, nil)

	// Should not panic
	client.SetDebug(true)
	client.SetDebug(false)
}

func TestClient_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}
}

func TestClient_Get_WithContext_Cancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.Get(ctx, "/test")

	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}

func TestClient_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["ticker"] != "TEST-MARKET" {
			t.Errorf("expected ticker TEST-MARKET, got %v", body["ticker"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "order-123"})
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	body := map[string]interface{}{"ticker": "TEST-MARKET", "count": 10}
	resp, err := client.Post(ctx, "/orders", body)

	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode())
	}
}

func TestClient_Put_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"updated": "true"})
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	body := map[string]interface{}{"count": 20}
	resp, err := client.Put(ctx, "/orders/123", body)

	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}
}

func TestClient_Delete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	resp, err := client.DeleteRaw(ctx, "/orders/123")

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode())
	}
}

func TestClient_RequestSigning(t *testing.T) {
	var receivedAuth string
	var receivedTimestamp string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedTimestamp = r.Header.Get("KALSHI-ACCESS-TIMESTAMP")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	signer := createTestSigner(t)
	client := createTestClientWithURLAndSigner(t, server.URL, signer)

	ctx := context.Background()
	_, err := client.Get(ctx, "/test")

	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if receivedAuth == "" {
		t.Error("Authorization header not set")
	}

	if !strings.HasPrefix(receivedAuth, "KALSHI-API-KEY ") {
		t.Errorf("Authorization header has wrong format: %s", receivedAuth)
	}

	if receivedTimestamp == "" {
		t.Error("KALSHI-ACCESS-TIMESTAMP header not set")
	}

	// Timestamp should be in ISO 8601 format
	_, err = time.Parse("2006-01-02T15:04:05Z", receivedTimestamp)
	if err != nil {
		t.Errorf("KALSHI-ACCESS-TIMESTAMP has wrong format: %s", receivedTimestamp)
	}
}

func TestClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "INVALID_PARAM",
			"message": "Invalid ticker format",
		})
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")

	// The client should return the response even on error status codes
	if err != nil {
		t.Fatalf("Get should not return error for API errors: %v", err)
	}

	if resp.StatusCode() != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode())
	}

	// Parse the error
	apiErr := ParseAPIError(resp)
	if apiErr == nil {
		t.Fatal("ParseAPIError returned nil for error response")
	}

	if apiErr.Code != "INVALID_PARAM" {
		t.Errorf("expected code INVALID_PARAM, got %s", apiErr.Code)
	}

	if apiErr.Message != "Invalid ticker format" {
		t.Errorf("expected message 'Invalid ticker format', got '%s'", apiErr.Message)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		Code:       "RATE_LIMITED",
		Message:    "Too many requests",
		StatusCode: 429,
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "RATE_LIMITED") {
		t.Errorf("error string should contain code: %s", errStr)
	}
	if !strings.Contains(errStr, "Too many requests") {
		t.Errorf("error string should contain message: %s", errStr)
	}
}

func TestClient_RateLimiting_Retry(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := atomic.AddInt32(&attempts, 1)

		if attempt < 3 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    "RATE_LIMITED",
				"message": "Too many requests",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")

	if err != nil {
		t.Fatalf("Get failed after retries: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("expected status 200 after retries, got %d", resp.StatusCode())
	}

	if atomic.LoadInt32(&attempts) < 3 {
		t.Errorf("expected at least 3 attempts, got %d", attempts)
	}
}

func TestClient_RateLimiting_MaxRetries(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    "RATE_LIMITED",
			"message": "Too many requests",
		})
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	resp, _ := client.Get(ctx, "/test")

	// After max retries, should return the rate limit response
	if resp.StatusCode() != http.StatusTooManyRequests {
		t.Errorf("expected status 429 after max retries, got %d", resp.StatusCode())
	}

	// Should have made multiple attempts
	if atomic.LoadInt32(&attempts) < 2 {
		t.Errorf("expected multiple retry attempts, got %d", attempts)
	}
}

func TestClient_ExponentialBackoff(t *testing.T) {
	var requestTimes []time.Time

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTimes = append(requestTimes, time.Now())

		if len(requestTimes) < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := createTestClientWithURL(t, server.URL)

	ctx := context.Background()
	client.Get(ctx, "/test")

	if len(requestTimes) < 3 {
		t.Skip("not enough retries to verify backoff")
	}

	// Second retry should have longer delay than first
	delay1 := requestTimes[1].Sub(requestTimes[0])
	delay2 := requestTimes[2].Sub(requestTimes[1])

	if delay2 < delay1 {
		t.Errorf("expected exponential backoff: delay1=%v, delay2=%v", delay1, delay2)
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"429 is rate limit", 429, true},
		{"200 is not rate limit", 200, false},
		{"400 is not rate limit", 400, false},
		{"500 is not rate limit", 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRateLimitError(tt.statusCode)
			if result != tt.expected {
				t.Errorf("IsRateLimitError(%d) = %v, want %v", tt.statusCode, result, tt.expected)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"500 is server error", 500, true},
		{"502 is server error", 502, true},
		{"503 is server error", 503, true},
		{"504 is server error", 504, true},
		{"400 is not server error", 400, false},
		{"429 is not server error", 429, false},
		{"200 is not server error", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsServerError(tt.statusCode)
			if result != tt.expected {
				t.Errorf("IsServerError(%d) = %v, want %v", tt.statusCode, result, tt.expected)
			}
		})
	}
}

// Helper functions

func createTestSigner(t *testing.T) *Signer {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("test-api-key", privateKey)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	return signer
}

func createTestClientWithURL(t *testing.T, baseURL string) *Client {
	t.Helper()

	cfg := &config.Config{
		API: config.APIConfig{
			Production: false,
			Timeout:    5 * time.Second,
		},
	}

	client := NewClient(cfg, nil)
	client.SetBaseURL(baseURL)

	return client
}

func createTestClientWithURLAndSigner(t *testing.T, baseURL string, signer *Signer) *Client {
	t.Helper()

	cfg := &config.Config{
		API: config.APIConfig{
			Production: false,
			Timeout:    5 * time.Second,
		},
	}

	client := NewClient(cfg, signer)
	client.SetBaseURL(baseURL)

	return client
}
