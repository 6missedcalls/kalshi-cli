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

// =============================================================================
// TDD Step 1: Write FAILING tests FIRST (RED)
// =============================================================================

func TestListMultivariateCollections(t *testing.T) {
	tests := []struct {
		name           string
		params         ListMultivariateCollectionsParams
		serverResponse models.MultivariateCollectionsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns collections successfully",
			params: ListMultivariateCollectionsParams{},
			serverResponse: models.MultivariateCollectionsResponse{
				Collections: []models.MultivariateCollection{
					{
						Ticker:      "PRES-2024",
						Title:       "2024 Presidential Election",
						Description: "Who will win the 2024 US Presidential Election?",
						Status:      "active",
					},
					{
						Ticker:      "SUPERBOWL-LIX",
						Title:       "Super Bowl LIX Winner",
						Description: "Which team will win Super Bowl LIX?",
						Status:      "active",
					},
				},
				Cursor: "next-cursor-123",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:   "returns collections with status filter",
			params: ListMultivariateCollectionsParams{Status: "active"},
			serverResponse: models.MultivariateCollectionsResponse{
				Collections: []models.MultivariateCollection{
					{
						Ticker: "ACTIVE-COLLECTION",
						Status: "active",
					},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:   "returns collections with pagination",
			params: ListMultivariateCollectionsParams{Cursor: "prev-cursor", Limit: 10},
			serverResponse: models.MultivariateCollectionsResponse{
				Collections: []models.MultivariateCollection{
					{Ticker: "COLLECTION-PAGE-2"},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles server error",
			params:         ListMultivariateCollectionsParams{},
			serverResponse: models.MultivariateCollectionsResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "handles unauthorized",
			params:         ListMultivariateCollectionsParams{},
			serverResponse: models.MultivariateCollectionsResponse{},
			serverStatus:   http.StatusUnauthorized,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/multivariate-collections"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.params.Status != "" {
					if got := r.URL.Query().Get("status"); got != tt.params.Status {
						t.Errorf("expected status=%s, got %s", tt.params.Status, got)
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

			client := multivariateTestClient(t, server.URL)
			resp, err := client.ListMultivariateCollections(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Collections) != tt.wantCount {
				t.Errorf("expected %d collections, got %d", tt.wantCount, len(resp.Collections))
			}
		})
	}
}

func TestGetMultivariateCollection(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		serverResponse models.MultivariateCollectionResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "returns collection successfully",
			ticker: "PRES-2024",
			serverResponse: models.MultivariateCollectionResponse{
				Collection: models.MultivariateCollection{
					Ticker:      "PRES-2024",
					Title:       "2024 Presidential Election",
					Description: "Who will win the 2024 US Presidential Election?",
					Status:      "active",
					LookupType:  "candidate",
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			ticker:         "INVALID-COLLECTION",
			serverResponse: models.MultivariateCollectionResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "handles server error",
			ticker:         "ERROR-COLLECTION",
			serverResponse: models.MultivariateCollectionResponse{},
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

				expectedPath := TradeAPIPrefix + "/multivariate-collections/" + tt.ticker
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

			client := multivariateTestClient(t, server.URL)
			resp, err := client.GetMultivariateCollection(context.Background(), tt.ticker)

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

func TestGetCollectionLookupHistory(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		ticker         string
		params         LookupHistoryParams
		serverResponse models.LookupHistoryResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns lookup history successfully",
			ticker: "PRES-2024",
			params: LookupHistoryParams{},
			serverResponse: models.LookupHistoryResponse{
				History: []models.LookupHistoryEntry{
					{
						Ticker:      "PRES-2024-TRUMP",
						LookupValue: "Donald Trump",
						CreatedTime: now.Add(-24 * time.Hour),
					},
					{
						Ticker:      "PRES-2024-HARRIS",
						LookupValue: "Kamala Harris",
						CreatedTime: now.Add(-12 * time.Hour),
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:   "returns lookup history with limit",
			ticker: "SUPERBOWL",
			params: LookupHistoryParams{Limit: 5},
			serverResponse: models.LookupHistoryResponse{
				History: []models.LookupHistoryEntry{
					{Ticker: "SB-CHIEFS", LookupValue: "Kansas City Chiefs"},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles not found",
			ticker:         "INVALID",
			params:         LookupHistoryParams{},
			serverResponse: models.LookupHistoryResponse{},
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

				expectedPath := TradeAPIPrefix + "/multivariate-collections/" + tt.ticker + "/lookup-history"
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

			client := multivariateTestClient(t, server.URL)
			resp, err := client.GetCollectionLookupHistory(context.Background(), tt.ticker, tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.History) != tt.wantCount {
				t.Errorf("expected %d history entries, got %d", tt.wantCount, len(resp.History))
			}
		})
	}
}

func TestCreateCollectionMarket(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		request        CreateCollectionMarketRequest
		serverResponse models.CreateCollectionMarketResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "creates market successfully",
			ticker: "PRES-2024",
			request: CreateCollectionMarketRequest{
				LookupValue: "Donald Trump",
			},
			serverResponse: models.CreateCollectionMarketResponse{
				MarketTicker: "PRES-2024-TRUMP",
				Created:      true,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:   "returns existing market",
			ticker: "PRES-2024",
			request: CreateCollectionMarketRequest{
				LookupValue: "Kamala Harris",
			},
			serverResponse: models.CreateCollectionMarketResponse{
				MarketTicker: "PRES-2024-HARRIS",
				Created:      false, // Market already existed
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:   "handles invalid lookup value",
			ticker: "PRES-2024",
			request: CreateCollectionMarketRequest{
				LookupValue: "",
			},
			serverResponse: models.CreateCollectionMarketResponse{},
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "handles collection not found",
			ticker:         "INVALID",
			request:        CreateCollectionMarketRequest{LookupValue: "Test"},
			serverResponse: models.CreateCollectionMarketResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/multivariate-collections/" + tt.ticker + "/markets"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				var req CreateCollectionMarketRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				if req.LookupValue != tt.request.LookupValue {
					t.Errorf("expected lookup_value %q, got %q", tt.request.LookupValue, req.LookupValue)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := multivariateTestClient(t, server.URL)
			resp, err := client.CreateCollectionMarket(context.Background(), tt.ticker, tt.request)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.MarketTicker == "" {
				t.Error("expected market ticker to be returned")
			}
		})
	}
}

func TestLookupCollectionMarket(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		lookupValue    string
		serverResponse models.LookupCollectionMarketResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:        "finds market successfully",
			ticker:      "PRES-2024",
			lookupValue: "Donald Trump",
			serverResponse: models.LookupCollectionMarketResponse{
				MarketTicker: "PRES-2024-TRUMP",
				Found:        true,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "market not found (never created)",
			ticker:         "PRES-2024",
			lookupValue:    "Unknown Candidate",
			serverResponse: models.LookupCollectionMarketResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "collection not found",
			ticker:         "INVALID",
			lookupValue:    "Test",
			serverResponse: models.LookupCollectionMarketResponse{},
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

				expectedPath := TradeAPIPrefix + "/multivariate-collections/" + tt.ticker + "/markets/lookup"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if got := r.URL.Query().Get("lookup_value"); got != tt.lookupValue {
					t.Errorf("expected lookup_value=%s, got %s", tt.lookupValue, got)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := multivariateTestClient(t, server.URL)
			resp, err := client.LookupCollectionMarket(context.Background(), tt.ticker, tt.lookupValue)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !resp.Found {
				t.Error("expected market to be found")
			}

			if resp.MarketTicker == "" {
				t.Error("expected market ticker to be returned")
			}
		})
	}
}

// multivariateTestClient creates a test client for multivariate collection tests
func multivariateTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client := NewClient(nil, nil)
	client.SetBaseURL(serverURL)
	return client
}
