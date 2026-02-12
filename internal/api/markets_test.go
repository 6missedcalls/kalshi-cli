package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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

				w.Header().Set("Content-Type", "application/json")
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

				w.Header().Set("Content-Type", "application/json")
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

				w.Header().Set("Content-Type", "application/json")
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

				w.Header().Set("Content-Type", "application/json")
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
		name         string
		params       GetCandlesticksParams
		rawResponse  string // raw JSON in Kalshi v2 format
		serverStatus int
		wantErr      bool
		wantCount    int
	}{
		{
			name:   "returns candlesticks successfully",
			params: GetCandlesticksParams{SeriesTicker: "BTC-SERIES", Ticker: "BTC-100K", Period: "1h"},
			rawResponse: `{
				"ticker": "BTC-100K",
				"candlesticks": [
					{"end_period_ts": 1704067200, "price": {"open": 45, "high": 48, "low": 44, "close": 47}, "volume": 1000, "open_interest": 0},
					{"end_period_ts": 1704070800, "price": {"open": 47, "high": 50, "low": 46, "close": 49}, "volume": 1200, "open_interest": 0}
				]
			}`,
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:         "handles invalid ticker",
			params:       GetCandlesticksParams{SeriesTicker: "INVALID-SERIES", Ticker: "INVALID", Period: "1h"},
			rawResponse:  `{}`,
			serverStatus: http.StatusNotFound,
			wantErr:      true,
		},
		{
			name: "returns candlesticks with time range",
			params: GetCandlesticksParams{
				SeriesTicker: "ETH-SERIES",
				Ticker:       "ETH-5K",
				Period:       "1d",
				StartTime:    time.Now().Add(-24 * time.Hour).Unix(),
				EndTime:      time.Now().Unix(),
			},
			rawResponse: `{
				"ticker": "ETH-5K",
				"candlesticks": [
					{"end_period_ts": 1704067200, "price": {"open": 30, "high": 32, "low": 29, "close": 31}, "volume": 0, "open_interest": 0}
				]
			}`,
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

				expectedPath := TradeAPIPrefix + "/series/" + tt.params.SeriesTicker + "/markets/" + tt.params.Ticker + "/candlesticks"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					w.Write([]byte(tt.rawResponse))
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

				w.Header().Set("Content-Type", "application/json")
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

				w.Header().Set("Content-Type", "application/json")
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

// TDD RED: Tests for missing ListMarketsParams fields (min_close_ts, max_close_ts, event_ticker, tickers)
func TestListMarkets_WithTimeFilters(t *testing.T) {
	tests := []struct {
		name         string
		params       ListMarketsParams
		wantQueryKey string
		wantQueryVal string
	}{
		{
			name: "filters by min_close_ts",
			params: ListMarketsParams{
				MinCloseTs: 1704067200,
			},
			wantQueryKey: "min_close_ts",
			wantQueryVal: "1704067200",
		},
		{
			name: "filters by max_close_ts",
			params: ListMarketsParams{
				MaxCloseTs: 1704153600,
			},
			wantQueryKey: "max_close_ts",
			wantQueryVal: "1704153600",
		},
		{
			name: "filters by event_ticker",
			params: ListMarketsParams{
				EventTicker: "PRES-2024",
			},
			wantQueryKey: "event_ticker",
			wantQueryVal: "PRES-2024",
		},
		{
			name: "filters by multiple tickers",
			params: ListMarketsParams{
				Tickers: []string{"BTC-100K", "ETH-10K"},
			},
			wantQueryKey: "tickers",
			wantQueryVal: "BTC-100K,ETH-10K",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got := r.URL.Query().Get(tt.wantQueryKey)
				if got != tt.wantQueryVal {
					t.Errorf("expected %s=%s, got %s", tt.wantQueryKey, tt.wantQueryVal, got)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(models.MarketsResponse{
					Markets: []models.Market{},
					Cursor:  "",
				})
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			_, err := client.ListMarkets(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TDD RED: Test for orderbook depth parameter
func TestGetOrderbook_WithDepth(t *testing.T) {
	tests := []struct {
		name      string
		ticker    string
		depth     int
		wantDepth string
	}{
		{
			name:      "requests orderbook with depth 5",
			ticker:    "BTC-100K",
			depth:     5,
			wantDepth: "5",
		},
		{
			name:      "requests orderbook with depth 10",
			ticker:    "ETH-5K",
			depth:     10,
			wantDepth: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got := r.URL.Query().Get("depth")
				if got != tt.wantDepth {
					t.Errorf("expected depth=%s, got %s", tt.wantDepth, got)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(models.OrderbookResponse{
					Orderbook: models.Orderbook{
						Ticker:  tt.ticker,
						YesBids: []models.OrderbookLevel{},
						YesAsks: []models.OrderbookLevel{},
					},
				})
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			_, err := client.GetOrderbookWithDepth(context.Background(), tt.ticker, tt.depth)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TDD RED: Test for trades with timestamp filters
func TestGetTrades_WithTimeFilters(t *testing.T) {
	tests := []struct {
		name         string
		params       GetTradesParams
		wantQueryKey string
		wantQueryVal string
	}{
		{
			name: "filters by min_ts",
			params: GetTradesParams{
				Ticker: "BTC-100K",
				MinTs:  1704067200,
			},
			wantQueryKey: "min_ts",
			wantQueryVal: "1704067200",
		},
		{
			name: "filters by max_ts",
			params: GetTradesParams{
				Ticker: "BTC-100K",
				MaxTs:  1704153600,
			},
			wantQueryKey: "max_ts",
			wantQueryVal: "1704153600",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got := r.URL.Query().Get(tt.wantQueryKey)
				if got != tt.wantQueryVal {
					t.Errorf("expected %s=%s, got %s", tt.wantQueryKey, tt.wantQueryVal, got)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(models.TradesResponse{
					Trades: []models.Trade{},
					Cursor: "",
				})
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			_, err := client.GetTrades(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TDD RED: Test for batch candlesticks endpoint
func TestGetBatchCandlesticks(t *testing.T) {
	tests := []struct {
		name         string
		params       GetBatchCandlesticksParams
		rawResponse  string // raw JSON in Kalshi v2 format
		serverStatus int
		wantErr      bool
		wantCount    int
	}{
		{
			name: "returns batch candlesticks for multiple tickers",
			params: GetBatchCandlesticksParams{
				Tickers: []string{"BTC-100K", "ETH-5K", "SOL-500"},
				Period:  "1h",
			},
			rawResponse: `{
				"market_candlesticks": [
					{
						"ticker": "BTC-100K",
						"candlesticks": [
							{"end_period_ts": 1704067200, "price": {"open": 45, "high": 48, "low": 44, "close": 47}, "volume": 0, "open_interest": 0}
						]
					},
					{
						"ticker": "ETH-5K",
						"candlesticks": [
							{"end_period_ts": 1704067200, "price": {"open": 30, "high": 32, "low": 29, "close": 31}, "volume": 0, "open_interest": 0}
						]
					},
					{
						"ticker": "SOL-500",
						"candlesticks": [
							{"end_period_ts": 1704067200, "price": {"open": 55, "high": 58, "low": 54, "close": 57}, "volume": 0, "open_interest": 0}
						]
					}
				]
			}`,
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    3,
		},
		{
			name: "handles server error",
			params: GetBatchCandlesticksParams{
				Tickers: []string{"INVALID"},
				Period:  "1h",
			},
			rawResponse:  `{}`,
			serverStatus: http.StatusInternalServerError,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/markets/candlesticks/batch"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify tickers query param
				tickersParam := r.URL.Query().Get("tickers")
				if tt.params.Tickers != nil && len(tt.params.Tickers) > 0 && len(tt.params.Tickers) <= 100 {
					expectedTickers := strings.Join(tt.params.Tickers, ",")
					if tickersParam != expectedTickers {
						t.Errorf("expected tickers=%s, got %s", expectedTickers, tickersParam)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					w.Write([]byte(tt.rawResponse))
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetBatchCandlesticks(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.MarketCandlesticks) != tt.wantCount {
				t.Errorf("expected %d market candlesticks, got %d", tt.wantCount, len(resp.MarketCandlesticks))
			}
		})
	}
}

// TestGetBatchCandlesticks_MaxTickersLimit tests client-side validation of ticker limit
func TestGetBatchCandlesticks_MaxTickersLimit(t *testing.T) {
	// Generate 101 tickers to exceed limit
	tickers := make([]string, 101)
	for i := range tickers {
		tickers[i] = "TICKER-" + strconv.Itoa(i)
	}

	client := newTestClient(t, "http://localhost:9999") // URL doesn't matter, will fail before request
	_, err := client.GetBatchCandlesticks(context.Background(), GetBatchCandlesticksParams{
		Tickers: tickers,
		Period:  "1h",
	})

	if err == nil {
		t.Error("expected error for exceeding max tickers, got nil")
	}

	if !strings.Contains(err.Error(), "exceeds maximum") {
		t.Errorf("expected error message about exceeding maximum, got: %v", err)
	}
}

// TDD RED: Test periodToInterval function covers all documented periods
func TestPeriodToInterval(t *testing.T) {
	tests := []struct {
		period   string
		expected string
	}{
		{"1m", "1"},
		{"5m", "5m"},
		{"15m", "15m"},
		{"1h", "60"},
		{"4h", "4h"},
		{"1d", "1440"},
		{"unknown", "unknown"}, // passthrough for unknown periods
	}

	for _, tt := range tests {
		t.Run(tt.period, func(t *testing.T) {
			result := periodToInterval(tt.period)
			if result != tt.expected {
				t.Errorf("periodToInterval(%q) = %q, want %q", tt.period, result, tt.expected)
			}
		})
	}
}
