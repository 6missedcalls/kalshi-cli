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

func TestGetRFQs(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedRFQs := []models.RFQ{
		{
			RFQID:       "rfq-1",
			Ticker:      "BTC-100K",
			Side:        "yes",
			Quantity:    100,
			Status:      "active",
			CreatedTime: now,
			ExpiresTime: now.Add(time.Hour),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/rfqs" {
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
	if result.RFQs[0].RFQID != "rfq-1" {
		t.Errorf("expected RFQ ID 'rfq-1', got '%s'", result.RFQs[0].RFQID)
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
	now := time.Now().UTC().Truncate(time.Second)
	expectedRFQ := models.RFQ{
		RFQID:       "rfq-123",
		Ticker:      "BTC-100K",
		Side:        "yes",
		Quantity:    50,
		Status:      "active",
		CreatedTime: now,
		ExpiresTime: now.Add(time.Hour),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/rfqs/rfq-123" {
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

	if result.RFQ.RFQID != "rfq-123" {
		t.Errorf("expected RFQ ID 'rfq-123', got '%s'", result.RFQ.RFQID)
	}
}

func TestCreateRFQ(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/rfqs" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.CreateRFQRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Ticker != "BTC-100K" {
			t.Errorf("expected ticker 'BTC-100K', got '%s'", req.Ticker)
		}
		if req.Side != "yes" {
			t.Errorf("expected side 'yes', got '%s'", req.Side)
		}
		if req.Quantity != 100 {
			t.Errorf("expected quantity 100, got %d", req.Quantity)
		}

		resp := models.RFQResponse{
			RFQ: models.RFQ{
				RFQID:       "new-rfq-id",
				Ticker:      req.Ticker,
				Side:        req.Side,
				Quantity:    req.Quantity,
				Status:      "active",
				CreatedTime: now,
				ExpiresTime: now.Add(time.Hour),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateRFQ(context.Background(), models.CreateRFQRequest{
		Ticker:   "BTC-100K",
		Side:     "yes",
		Quantity: 100,
	})
	if err != nil {
		t.Fatalf("CreateRFQ failed: %v", err)
	}

	if result.RFQ.RFQID != "new-rfq-id" {
		t.Errorf("expected RFQ ID 'new-rfq-id', got '%s'", result.RFQ.RFQID)
	}
}

func TestCancelRFQ(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/rfqs/rfq-to-cancel" {
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
	now := time.Now().UTC().Truncate(time.Second)
	expectedQuotes := []models.Quote{
		{
			QuoteID:     "quote-1",
			RFQID:       "rfq-1",
			Ticker:      "BTC-100K",
			Side:        "yes",
			Price:       55,
			Quantity:    100,
			Status:      "active",
			CreatedTime: now,
			ExpiresTime: now.Add(time.Minute * 5),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/quotes" {
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
	if result.Quotes[0].QuoteID != "quote-1" {
		t.Errorf("expected quote ID 'quote-1', got '%s'", result.Quotes[0].QuoteID)
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
	now := time.Now().UTC().Truncate(time.Second)
	expectedQuote := models.Quote{
		QuoteID:     "quote-123",
		RFQID:       "rfq-1",
		Ticker:      "BTC-100K",
		Price:       60,
		Quantity:    50,
		Status:      "active",
		CreatedTime: now,
		ExpiresTime: now.Add(time.Minute * 5),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/quotes/quote-123" {
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

	if result.Quote.QuoteID != "quote-123" {
		t.Errorf("expected quote ID 'quote-123', got '%s'", result.Quote.QuoteID)
	}
}

func TestCreateQuote(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/quotes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.CreateQuoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.RFQID != "rfq-123" {
			t.Errorf("expected rfq_id 'rfq-123', got '%s'", req.RFQID)
		}
		if req.Price != 55 {
			t.Errorf("expected price 55, got %d", req.Price)
		}

		resp := models.QuoteResponse{
			Quote: models.Quote{
				QuoteID:     "new-quote-id",
				RFQID:       req.RFQID,
				Price:       req.Price,
				Quantity:    req.Quantity,
				Status:      "active",
				CreatedTime: now,
				ExpiresTime: now.Add(time.Minute * 5),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateQuote(context.Background(), models.CreateQuoteRequest{
		RFQID:    "rfq-123",
		Price:    55,
		Quantity: 100,
	})
	if err != nil {
		t.Fatalf("CreateQuote failed: %v", err)
	}

	if result.Quote.QuoteID != "new-quote-id" {
		t.Errorf("expected quote ID 'new-quote-id', got '%s'", result.Quote.QuoteID)
	}
}

func TestAcceptQuote(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/quotes/quote-123/accept" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.QuoteResponse{
			Quote: models.Quote{
				QuoteID:     "quote-123",
				Status:      "accepted",
				CreatedTime: now,
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
		if r.URL.Path != "/trade-api/v2/quotes/quote-to-cancel" {
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
