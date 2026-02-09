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

func TestListMarkets(t *testing.T) {
	tests := []struct {
		name           string
		params         ListMarketsParams
		serverResponse models.MarketsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
		wantCursor     string
	}{
		{
			name:   "returns markets successfully",
			params: ListMarketsParams{},
			serverResponse: models.MarketsResponse{
				Markets: []models.Market{
					{Ticker: "BTC-100K", Title: "Bitcoin to $100K"},
					{Ticker: "ETH-10K", Title: "Ethereum to $10K"},
				},
				Cursor: "next-cursor-123",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
			wantCursor:   "next-cursor-123",
		},
		{
			name:   "returns markets with status filter",
			params: ListMarketsParams{Status: "open"},
			serverResponse: models.MarketsResponse{
				Markets: []models.Market{
					{Ticker: "BTC-100K", Title: "Bitcoin to $100K", Status: "open"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
			wantCursor:   "",
		},
		{
			name:   "returns markets with pagination",
			params: ListMarketsParams{Cursor: "prev-cursor", Limit: 10},
			serverResponse: models.MarketsResponse{
				Markets: []models.Market{
					{Ticker: "MARKET-1", Title: "Market 1"},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
			wantCursor:   "next-cursor",
		},
		{
			name:           "handles server error",
			params:         ListMarketsParams{},
			serverResponse: models.MarketsResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:   "returns markets with series ticker filter",
			params: ListMarketsParams{SeriesTicker: "FED-RATES"},
			serverResponse: models.MarketsResponse{
				Markets: []models.Market{
					{Ticker: "FED-RATE-MAR", Title: "March Rate Decision"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if tt.params.Status != "" {
					if got := r.URL.Query().Get("status"); got != tt.params.Status {
						t.Errorf("expected status=%s, got %s", tt.params.Status, got)
					}
				}

				if tt.params.SeriesTicker != "" {
					if got := r.URL.Query().Get("series_ticker"); got != tt.params.SeriesTicker {
						t.Errorf("expected series_ticker=%s, got %s", tt.params.SeriesTicker, got)
					}
				}

				if tt.params.Cursor != "" {
					if got := r.URL.Query().Get("cursor"); got != tt.params.Cursor {
						t.Errorf("expected cursor=%s, got %s", tt.params.Cursor, got)
					}
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.ListMarkets(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Markets) != tt.wantCount {
				t.Errorf("expected %d markets, got %d", tt.wantCount, len(resp.Markets))
			}

			if resp.Cursor != tt.wantCursor {
				t.Errorf("expected cursor %q, got %q", tt.wantCursor, resp.Cursor)
			}
		})
	}
}

func TestGetMarket(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		serverResponse models.MarketResponse
		serverStatus   int
		wantErr        bool
		wantTicker     string
	}{
		{
			name:   "returns single market successfully",
			ticker: "BTC-100K",
			serverResponse: models.MarketResponse{
				Market: models.Market{
					Ticker: "BTC-100K",
					Title:  "Bitcoin to $100K",
					Status: "open",
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantTicker:   "BTC-100K",
		},
		{
			name:           "handles not found",
			ticker:         "INVALID-TICKER",
			serverResponse: models.MarketResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:   "returns market with all fields",
			ticker: "ETH-5K",
			serverResponse: models.MarketResponse{
				Market: models.Market{
					Ticker:       "ETH-5K",
					EventTicker:  "ETH-PRICE",
					Title:        "Ethereum to $5K",
					Status:       "open",
					YesBid:       45,
					YesAsk:       47,
					Volume:       10000,
					OpenInterest: 5000,
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantTicker:   "ETH-5K",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/markets/" + tt.ticker
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetMarket(context.Background(), tt.ticker)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Ticker != tt.wantTicker {
				t.Errorf("expected ticker %q, got %q", tt.wantTicker, resp.Ticker)
			}
		})
	}
}

func TestGetOrderbook(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		serverResponse models.OrderbookResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "returns orderbook successfully",
			ticker: "BTC-100K",
			serverResponse: models.OrderbookResponse{
				Orderbook: models.Orderbook{
					Ticker: "BTC-100K",
					YesBids: []models.OrderbookLevel{
						{Price: 45, Quantity: 100},
						{Price: 44, Quantity: 200},
					},
					YesAsks: []models.OrderbookLevel{
						{Price: 47, Quantity: 150},
						{Price: 48, Quantity: 250},
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles invalid ticker",
			ticker:         "INVALID",
			serverResponse: models.OrderbookResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:   "returns empty orderbook",
			ticker: "ETH-5K",
			serverResponse: models.OrderbookResponse{
				Orderbook: models.Orderbook{
					Ticker:  "ETH-5K",
					YesBids: []models.OrderbookLevel{},
					YesAsks: []models.OrderbookLevel{},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/markets/" + tt.ticker + "/orderbook"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetOrderbook(context.Background(), tt.ticker)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Ticker != tt.ticker {
				t.Errorf("expected ticker %q, got %q", tt.ticker, resp.Ticker)
			}
		})
	}
}

func TestGetTrades(t *testing.T) {
	tests := []struct {
		name           string
		params         GetTradesParams
		serverResponse models.TradesResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns trades successfully",
			params: GetTradesParams{Ticker: "BTC-100K"},
			serverResponse: models.TradesResponse{
				Trades: []models.Trade{
					{TradeID: "trade-1", Ticker: "BTC-100K", Price: 45, Count: 10},
					{TradeID: "trade-2", Ticker: "BTC-100K", Price: 46, Count: 5},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:   "returns trades with pagination",
			params: GetTradesParams{Cursor: "cursor-123", Limit: 50},
			serverResponse: models.TradesResponse{
				Trades: []models.Trade{
					{TradeID: "trade-3", Ticker: "ETH-5K", Price: 30, Count: 20},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles server error",
			params:         GetTradesParams{},
			serverResponse: models.TradesResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if tt.params.Ticker != "" {
					if got := r.URL.Query().Get("ticker"); got != tt.params.Ticker {
						t.Errorf("expected ticker=%s, got %s", tt.params.Ticker, got)
					}
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetTrades(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Trades) != tt.wantCount {
				t.Errorf("expected %d trades, got %d", tt.wantCount, len(resp.Trades))
			}
		})
	}
}

func TestGetCandlesticks(t *testing.T) {
	tests := []struct {
		name           string
		params         GetCandlesticksParams
		serverResponse models.CandlesticksResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns candlesticks successfully",
			params: GetCandlesticksParams{Ticker: "BTC-100K", Period: "1h"},
			serverResponse: models.CandlesticksResponse{
				Candlesticks: []models.Candlestick{
					{Ticker: "BTC-100K", Open: 45, High: 48, Low: 44, Close: 47, Volume: 1000},
					{Ticker: "BTC-100K", Open: 47, High: 50, Low: 46, Close: 49, Volume: 1200},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:           "handles invalid ticker",
			params:         GetCandlesticksParams{Ticker: "INVALID", Period: "1h"},
			serverResponse: models.CandlesticksResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "returns candlesticks with time range",
			params: GetCandlesticksParams{
				Ticker:    "ETH-5K",
				Period:    "1d",
				StartTime: time.Now().Add(-24 * time.Hour).Unix(),
				EndTime:   time.Now().Unix(),
			},
			serverResponse: models.CandlesticksResponse{
				Candlesticks: []models.Candlestick{
					{Ticker: "ETH-5K", Open: 30, High: 32, Low: 29, Close: 31},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/markets/" + tt.params.Ticker + "/candlesticks"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetCandlesticks(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Candlesticks) != tt.wantCount {
				t.Errorf("expected %d candlesticks, got %d", tt.wantCount, len(resp.Candlesticks))
			}
		})
	}
}

func TestListSeries(t *testing.T) {
	tests := []struct {
		name           string
		params         ListSeriesParams
		serverResponse models.SeriesResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns series successfully",
			params: ListSeriesParams{},
			serverResponse: models.SeriesResponse{
				Series: []models.Series{
					{Ticker: "FED-RATES", Title: "Federal Reserve Rate Decisions"},
					{Ticker: "BTC-PRICE", Title: "Bitcoin Price Milestones"},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:   "returns series with cursor pagination",
			params: ListSeriesParams{Cursor: "cursor-123"},
			serverResponse: models.SeriesResponse{
				Series: []models.Series{
					{Ticker: "ETH-PRICE", Title: "Ethereum Price Milestones"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:   "returns series with category filter",
			params: ListSeriesParams{Category: "crypto"},
			serverResponse: models.SeriesResponse{
				Series: []models.Series{
					{Ticker: "BTC-PRICE", Title: "Bitcoin Price Milestones", Category: "crypto"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles server error",
			params:         ListSeriesParams{},
			serverResponse: models.SeriesResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if tt.params.Category != "" {
					if got := r.URL.Query().Get("category"); got != tt.params.Category {
						t.Errorf("expected category=%s, got %s", tt.params.Category, got)
					}
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.ListSeries(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Series) != tt.wantCount {
				t.Errorf("expected %d series, got %d", tt.wantCount, len(resp.Series))
			}
		})
	}
}

func TestGetSeries(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		serverResponse GetSeriesResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "returns series by ticker successfully",
			ticker: "FED-RATES",
			serverResponse: GetSeriesResponse{
				Series: models.Series{
					Ticker:   "FED-RATES",
					Title:    "Federal Reserve Rate Decisions",
					Category: "economics",
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			ticker:         "INVALID-SERIES",
			serverResponse: GetSeriesResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/series/" + tt.ticker
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetSeries(context.Background(), tt.ticker)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Ticker != tt.ticker {
				t.Errorf("expected ticker %q, got %q", tt.ticker, resp.Ticker)
			}
		})
	}
}

// newTestClient creates a test client with a mock server URL
func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	signer := newTestSigner(t)
	return NewClientLegacy(signer, WithBaseURL(serverURL))
}

// newTestSigner creates a test signer with a generated key
func newTestSigner(t *testing.T) *Signer {
	t.Helper()
	key, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	signer, err := NewSigner("test-api-key-id", key)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}
	return signer
}
