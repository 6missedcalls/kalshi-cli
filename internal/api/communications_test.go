package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

func TestGetRFQs(t *testing.T) {
	expectedRFQs := []models.RFQ{
		{
			ID:           "rfq-1",
			MarketTicker: "BTC-100K",
			Contracts:    100,
			Status:       "active",
			CreatedTs:    "2024-01-15T12:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/rfqs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.RFQsResponse{
			RFQs:   expectedRFQs,
			Cursor: "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetRFQs(context.Background(), RFQsOptions{})
	if err != nil {
		t.Fatalf("GetRFQs failed: %v", err)
	}

	if len(result.RFQs) != 1 {
		t.Errorf("expected 1 RFQ, got %d", len(result.RFQs))
	}
	if result.RFQs[0].ID != "rfq-1" {
		t.Errorf("expected RFQ ID 'rfq-1', got '%s'", result.RFQs[0].ID)
	}
}

func TestGetRFQsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("ticker") != "BTC-100K" {
			t.Errorf("expected ticker 'BTC-100K', got '%s'", query.Get("ticker"))
		}
		if query.Get("status") != "active" {
			t.Errorf("expected status 'active', got '%s'", query.Get("status"))
		}

		resp := models.RFQsResponse{RFQs: []models.RFQ{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetRFQs(context.Background(), RFQsOptions{
		Ticker: "BTC-100K",
		Status: "active",
	})
	if err != nil {
		t.Fatalf("GetRFQs failed: %v", err)
	}
}

func TestGetRFQ(t *testing.T) {
	expectedRFQ := models.RFQ{
		ID:           "rfq-123",
		MarketTicker: "BTC-100K",
		Contracts:    50,
		Status:       "active",
		CreatedTs:    "2024-01-15T12:00:00Z",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/rfqs/rfq-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.RFQResponse{RFQ: expectedRFQ}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetRFQ(context.Background(), "rfq-123")
	if err != nil {
		t.Fatalf("GetRFQ failed: %v", err)
	}

	if result.RFQ.ID != "rfq-123" {
		t.Errorf("expected RFQ ID 'rfq-123', got '%s'", result.RFQ.ID)
	}
}

func TestCreateRFQ(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/rfqs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.CreateRFQRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.MarketTicker != "BTC-100K" {
			t.Errorf("expected market_ticker 'BTC-100K', got '%s'", req.MarketTicker)
		}
		if req.Contracts != 100 {
			t.Errorf("expected contracts 100, got %d", req.Contracts)
		}

		resp := models.RFQResponse{
			RFQ: models.RFQ{
				ID:           "new-rfq-id",
				MarketTicker: req.MarketTicker,
				Contracts:    req.Contracts,
				Status:       "active",
				CreatedTs:    "2024-01-15T12:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateRFQ(context.Background(), models.CreateRFQRequest{
		MarketTicker: "BTC-100K",
		Contracts:    100,
	})
	if err != nil {
		t.Fatalf("CreateRFQ failed: %v", err)
	}

	if result.RFQ.ID != "new-rfq-id" {
		t.Errorf("expected RFQ ID 'new-rfq-id', got '%s'", result.RFQ.ID)
	}
}

func TestCancelRFQ(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/rfqs/rfq-to-cancel" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	err := client.CancelRFQ(context.Background(), "rfq-to-cancel")
	if err != nil {
		t.Fatalf("CancelRFQ failed: %v", err)
	}
}

func TestGetQuotes(t *testing.T) {
	expectedQuotes := []models.Quote{
		{
			ID:        "quote-1",
			RFQID:     "rfq-1",
			YesBid:    55,
			Contracts: 100,
			Status:    "active",
			CreatedTs: "2024-01-15T12:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/quotes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.QuotesResponse{
			Quotes: expectedQuotes,
			Cursor: "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetQuotes(context.Background(), QuotesOptions{})
	if err != nil {
		t.Fatalf("GetQuotes failed: %v", err)
	}

	if len(result.Quotes) != 1 {
		t.Errorf("expected 1 quote, got %d", len(result.Quotes))
	}
	if result.Quotes[0].ID != "quote-1" {
		t.Errorf("expected quote ID 'quote-1', got '%s'", result.Quotes[0].ID)
	}
}

func TestGetQuotesWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("rfq_id") != "rfq-123" {
			t.Errorf("expected rfq_id 'rfq-123', got '%s'", query.Get("rfq_id"))
		}
		if query.Get("status") != "active" {
			t.Errorf("expected status 'active', got '%s'", query.Get("status"))
		}

		resp := models.QuotesResponse{Quotes: []models.Quote{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetQuotes(context.Background(), QuotesOptions{
		RFQID:  "rfq-123",
		Status: "active",
	})
	if err != nil {
		t.Fatalf("GetQuotes failed: %v", err)
	}
}

func TestGetQuote(t *testing.T) {
	expectedQuote := models.Quote{
		ID:        "quote-123",
		RFQID:     "rfq-1",
		YesBid:    60,
		Contracts: 50,
		Status:    "active",
		CreatedTs: "2024-01-15T12:00:00Z",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/quotes/quote-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.QuoteResponse{Quote: expectedQuote}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetQuote(context.Background(), "quote-123")
	if err != nil {
		t.Fatalf("GetQuote failed: %v", err)
	}

	if result.Quote.ID != "quote-123" {
		t.Errorf("expected quote ID 'quote-123', got '%s'", result.Quote.ID)
	}
}

func TestCreateQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/quotes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.CreateQuoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.RFQID != "rfq-123" {
			t.Errorf("expected rfq_id 'rfq-123', got '%s'", req.RFQID)
		}
		if req.YesBid != 55 {
			t.Errorf("expected yes_bid 55, got %d", req.YesBid)
		}

		resp := models.QuoteResponse{
			Quote: models.Quote{
				ID:        "new-quote-id",
				RFQID:     req.RFQID,
				YesBid:    req.YesBid,
				Contracts: 100,
				Status:    "active",
				CreatedTs: "2024-01-15T12:00:00Z",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateQuote(context.Background(), models.CreateQuoteRequest{
		RFQID:  "rfq-123",
		YesBid: 55,
	})
	if err != nil {
		t.Fatalf("CreateQuote failed: %v", err)
	}

	if result.Quote.ID != "new-quote-id" {
		t.Errorf("expected quote ID 'new-quote-id', got '%s'", result.Quote.ID)
	}
}

func TestAcceptQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/quotes/quote-123/accept" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.QuoteResponse{
			Quote: models.Quote{
				ID:     "quote-123",
				Status: "accepted",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.AcceptQuote(context.Background(), "quote-123")
	if err != nil {
		t.Fatalf("AcceptQuote failed: %v", err)
	}

	if result.Quote.Status != "accepted" {
		t.Errorf("expected status 'accepted', got '%s'", result.Quote.Status)
	}
}

func TestCancelQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/quotes/quote-to-cancel" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	err := client.CancelQuote(context.Background(), "quote-to-cancel")
	if err != nil {
		t.Fatalf("CancelQuote failed: %v", err)
	}
}

func TestGetCommunicationsID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/communications/id" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.CommunicationsIDResponse{
			CommunicationsID: "comm-id-123",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetCommunicationsID(context.Background())
	if err != nil {
		t.Fatalf("GetCommunicationsID failed: %v", err)
	}

	if result.CommunicationsID != "comm-id-123" {
		t.Errorf("expected communications ID 'comm-id-123', got '%s'", result.CommunicationsID)
	}
}

func TestCommunicationsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request",
			"code":  "INVALID_REQUEST",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.CreateRFQ(context.Background(), models.CreateRFQRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}
