package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// These tests verify correct API paths according to Kalshi documentation.
// RFQs and Quotes should be under /communications namespace.

func TestCommunicationsPathsAudit(t *testing.T) {
	// Define expected paths according to Kalshi API documentation
	expectedPaths := map[string]string{
		"GetRFQs":       "/trade-api/v2/communications/rfqs",
		"GetRFQ":        "/trade-api/v2/communications/rfqs/rfq-123",
		"CreateRFQ":     "/trade-api/v2/communications/rfqs",
		"DeleteRFQ":     "/trade-api/v2/communications/rfqs/rfq-123",
		"GetQuotes":     "/trade-api/v2/communications/quotes",
		"GetQuote":      "/trade-api/v2/communications/quotes/quote-456",
		"CreateQuote":   "/trade-api/v2/communications/quotes",
		"AcceptQuote":   "/trade-api/v2/communications/quotes/quote-456/accept",
		"ConfirmQuote":  "/trade-api/v2/communications/quotes/quote-456/confirm",
		"DeleteQuote":   "/trade-api/v2/communications/quotes/quote-456",
		"GetCommID":     "/trade-api/v2/communications/id",
	}

	t.Run("GetRFQs uses correct path", func(t *testing.T) {
		var actualPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.RFQsResponse{RFQs: []models.RFQ{}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.GetRFQs(context.Background(), RFQsOptions{})

		if actualPath != expectedPaths["GetRFQs"] {
			t.Errorf("GetRFQs path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["GetRFQs"], actualPath)
		}
	})

	t.Run("GetRFQ uses correct path", func(t *testing.T) {
		var actualPath string
		now := time.Now().UTC()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.RFQResponse{RFQ: models.RFQ{
				RFQID:       "rfq-123",
				CreatedTime: now,
				ExpiresTime: now.Add(time.Hour),
			}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.GetRFQ(context.Background(), "rfq-123")

		if actualPath != expectedPaths["GetRFQ"] {
			t.Errorf("GetRFQ path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["GetRFQ"], actualPath)
		}
	})

	t.Run("CreateRFQ uses correct path", func(t *testing.T) {
		var actualPath string
		now := time.Now().UTC()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.RFQResponse{RFQ: models.RFQ{
				RFQID:       "new-rfq",
				CreatedTime: now,
				ExpiresTime: now.Add(time.Hour),
			}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.CreateRFQ(context.Background(), models.CreateRFQRequest{
			Ticker:   "TEST",
			Side:     "yes",
			Quantity: 100,
		})

		if actualPath != expectedPaths["CreateRFQ"] {
			t.Errorf("CreateRFQ path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["CreateRFQ"], actualPath)
		}
	})

	t.Run("CancelRFQ uses correct path", func(t *testing.T) {
		var actualPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.CancelRFQ(context.Background(), "rfq-123")

		if actualPath != expectedPaths["DeleteRFQ"] {
			t.Errorf("CancelRFQ path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["DeleteRFQ"], actualPath)
		}
	})

	t.Run("GetQuotes uses correct path", func(t *testing.T) {
		var actualPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.QuotesResponse{Quotes: []models.Quote{}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.GetQuotes(context.Background(), QuotesOptions{})

		if actualPath != expectedPaths["GetQuotes"] {
			t.Errorf("GetQuotes path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["GetQuotes"], actualPath)
		}
	})

	t.Run("GetQuote uses correct path", func(t *testing.T) {
		var actualPath string
		now := time.Now().UTC()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.QuoteResponse{Quote: models.Quote{
				QuoteID:     "quote-456",
				CreatedTime: now,
				ExpiresTime: now.Add(time.Minute * 5),
			}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.GetQuote(context.Background(), "quote-456")

		if actualPath != expectedPaths["GetQuote"] {
			t.Errorf("GetQuote path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["GetQuote"], actualPath)
		}
	})

	t.Run("CreateQuote uses correct path", func(t *testing.T) {
		var actualPath string
		now := time.Now().UTC()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.QuoteResponse{Quote: models.Quote{
				QuoteID:     "new-quote",
				CreatedTime: now,
				ExpiresTime: now.Add(time.Minute * 5),
			}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.CreateQuote(context.Background(), models.CreateQuoteRequest{
			RFQID: "rfq-123",
			Price: 50,
		})

		if actualPath != expectedPaths["CreateQuote"] {
			t.Errorf("CreateQuote path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["CreateQuote"], actualPath)
		}
	})

	t.Run("AcceptQuote uses correct path", func(t *testing.T) {
		var actualPath string
		now := time.Now().UTC()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.QuoteResponse{Quote: models.Quote{
				QuoteID:     "quote-456",
				Status:      "accepted",
				CreatedTime: now,
			}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.AcceptQuote(context.Background(), "quote-456")

		if actualPath != expectedPaths["AcceptQuote"] {
			t.Errorf("AcceptQuote path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["AcceptQuote"], actualPath)
		}
	})

	t.Run("ConfirmQuote uses correct path", func(t *testing.T) {
		var actualPath string
		now := time.Now().UTC()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.QuoteResponse{Quote: models.Quote{
				QuoteID:     "quote-456",
				Status:      "confirmed",
				CreatedTime: now,
			}})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.ConfirmQuote(context.Background(), "quote-456")

		if actualPath != expectedPaths["ConfirmQuote"] {
			t.Errorf("ConfirmQuote path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["ConfirmQuote"], actualPath)
		}
	})

	t.Run("CancelQuote uses correct path", func(t *testing.T) {
		var actualPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.CancelQuote(context.Background(), "quote-456")

		if actualPath != expectedPaths["DeleteQuote"] {
			t.Errorf("CancelQuote path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["DeleteQuote"], actualPath)
		}
	})

	t.Run("GetCommunicationsID uses correct path", func(t *testing.T) {
		var actualPath string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualPath = r.URL.Path
			json.NewEncoder(w).Encode(models.CommunicationsIDResponse{CommunicationsID: "comm-123"})
		}))
		defer server.Close()

		client := createPathTestClient(t, server.URL)
		client.GetCommunicationsID(context.Background())

		if actualPath != expectedPaths["GetCommID"] {
			t.Errorf("GetCommunicationsID path mismatch:\n  expected: %s\n  actual:   %s", expectedPaths["GetCommID"], actualPath)
		}
	})
}

func TestConfirmQuoteExists(t *testing.T) {
	now := time.Now().UTC()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.QuoteResponse{Quote: models.Quote{
			QuoteID:     "quote-789",
			Status:      "confirmed",
			CreatedTime: now,
		}})
	}))
	defer server.Close()

	client := createPathTestClient(t, server.URL)
	result, err := client.ConfirmQuote(context.Background(), "quote-789")
	if err != nil {
		t.Fatalf("ConfirmQuote failed: %v", err)
	}
	if result.Quote.Status != "confirmed" {
		t.Errorf("expected status 'confirmed', got '%s'", result.Quote.Status)
	}
}

func TestEmptyIDValidation(t *testing.T) {
	client := &Client{}

	t.Run("GetRFQ with empty ID", func(t *testing.T) {
		_, err := client.GetRFQ(context.Background(), "")
		if err == nil || err.Error() != "RFQ ID is required" {
			t.Errorf("expected 'RFQ ID is required' error, got: %v", err)
		}
	})

	t.Run("CancelRFQ with empty ID", func(t *testing.T) {
		err := client.CancelRFQ(context.Background(), "")
		if err == nil || err.Error() != "RFQ ID is required" {
			t.Errorf("expected 'RFQ ID is required' error, got: %v", err)
		}
	})

	t.Run("GetQuote with empty ID", func(t *testing.T) {
		_, err := client.GetQuote(context.Background(), "")
		if err == nil || err.Error() != "quote ID is required" {
			t.Errorf("expected 'quote ID is required' error, got: %v", err)
		}
	})

	t.Run("AcceptQuote with empty ID", func(t *testing.T) {
		_, err := client.AcceptQuote(context.Background(), "")
		if err == nil || err.Error() != "quote ID is required" {
			t.Errorf("expected 'quote ID is required' error, got: %v", err)
		}
	})

	t.Run("ConfirmQuote with empty ID", func(t *testing.T) {
		_, err := client.ConfirmQuote(context.Background(), "")
		if err == nil || err.Error() != "quote ID is required" {
			t.Errorf("expected 'quote ID is required' error, got: %v", err)
		}
	})

	t.Run("CancelQuote with empty ID", func(t *testing.T) {
		err := client.CancelQuote(context.Background(), "")
		if err == nil || err.Error() != "quote ID is required" {
			t.Errorf("expected 'quote ID is required' error, got: %v", err)
		}
	})
}

// Helper to create a test client with custom base URL
func createPathTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("test-api-key", privateKey)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	return NewClientLegacy(signer, WithBaseURL(baseURL))
}
