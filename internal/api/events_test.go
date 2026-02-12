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

func TestListEvents(t *testing.T) {
	tests := []struct {
		name           string
		params         ListEventsParams
		serverResponse models.EventsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
		wantCursor     string
	}{
		{
			name:   "returns events successfully",
			params: ListEventsParams{},
			serverResponse: models.EventsResponse{
				Events: []models.Event{
					{EventTicker: "ELECTION-2024", Title: "2024 Presidential Election"},
					{EventTicker: "FED-MAR-2024", Title: "March 2024 Fed Decision"},
				},
				Cursor: "next-cursor-123",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
			wantCursor:   "next-cursor-123",
		},
		{
			name:   "returns events with status filter",
			params: ListEventsParams{Status: "open"},
			serverResponse: models.EventsResponse{
				Events: []models.Event{
					{EventTicker: "ELECTION-2024", Title: "2024 Presidential Election"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
			wantCursor:   "",
		},
		{
			name:   "returns events with pagination",
			params: ListEventsParams{Cursor: "prev-cursor", Limit: 10},
			serverResponse: models.EventsResponse{
				Events: []models.Event{
					{EventTicker: "EVENT-1", Title: "Event 1"},
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
			params:         ListEventsParams{},
			serverResponse: models.EventsResponse{},
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

			client := newTestClient(t, server.URL)
			events, cursor, err := client.ListEvents(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(events) != tt.wantCount {
				t.Errorf("expected %d events, got %d", tt.wantCount, len(events))
			}

			if cursor != tt.wantCursor {
				t.Errorf("expected cursor %q, got %q", tt.wantCursor, cursor)
			}
		})
	}
}

func TestGetEvent(t *testing.T) {
	tests := []struct {
		name           string
		eventTicker    string
		serverResponse models.EventResponse
		serverStatus   int
		wantErr        bool
		wantTicker     string
	}{
		{
			name:        "returns single event successfully",
			eventTicker: "ELECTION-2024",
			serverResponse: models.EventResponse{
				Event: models.Event{
					EventTicker: "ELECTION-2024",
					Title:       "2024 Presidential Election",
					Category:    "politics",
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantTicker:   "ELECTION-2024",
		},
		{
			name:           "handles not found",
			eventTicker:    "INVALID-EVENT",
			serverResponse: models.EventResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:        "returns event with markets list",
			eventTicker: "FED-MAR-2024",
			serverResponse: models.EventResponse{
				Event: models.Event{
					EventTicker: "FED-MAR-2024",
					Title:       "March 2024 Fed Decision",
					Markets:     []string{"FED-RATE-25", "FED-RATE-50", "FED-RATE-HOLD"},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantTicker:   "FED-MAR-2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := TradeAPIPrefix + "/events/" + tt.eventTicker
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
			resp, err := client.GetEvent(context.Background(), tt.eventTicker)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.EventTicker != tt.wantTicker {
				t.Errorf("expected event ticker %q, got %q", tt.wantTicker, resp.EventTicker)
			}
		})
	}
}

func TestListMultivariateEvents(t *testing.T) {
	tests := []struct {
		name           string
		params         ListMultivariateParams
		serverResponse models.MultivariateEventsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns multivariate events successfully",
			params: ListMultivariateParams{},
			serverResponse: models.MultivariateEventsResponse{
				Events: []models.MultivariateEvent{
					{Ticker: "MV-EVENT-1", Title: "Multivariate Event 1"},
					{Ticker: "MV-EVENT-2", Title: "Multivariate Event 2"},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:   "returns multivariate events with pagination",
			params: ListMultivariateParams{Cursor: "cursor-123", Limit: 10},
			serverResponse: models.MultivariateEventsResponse{
				Events: []models.MultivariateEvent{
					{Ticker: "MV-EVENT-3", Title: "Multivariate Event 3"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:   "returns multivariate events with status filter",
			params: ListMultivariateParams{Status: "open"},
			serverResponse: models.MultivariateEventsResponse{
				Events: []models.MultivariateEvent{
					{Ticker: "MV-EVENT-4", Title: "Multivariate Event 4", Status: "open"},
				},
				Cursor: "",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles server error",
			params:         ListMultivariateParams{},
			serverResponse: models.MultivariateEventsResponse{},
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

				if tt.params.Status != "" {
					if got := r.URL.Query().Get("status"); got != tt.params.Status {
						t.Errorf("expected status=%s, got %s", tt.params.Status, got)
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
			events, _, err := client.ListMultivariateEvents(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(events) != tt.wantCount {
				t.Errorf("expected %d multivariate events, got %d", tt.wantCount, len(events))
			}
		})
	}
}

func TestGetMultivariateEvent(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		serverResponse MultivariateEventResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "returns multivariate event successfully",
			ticker: "MV-EVENT-1",
			serverResponse: MultivariateEventResponse{
				Event: models.MultivariateEvent{
					Ticker: "MV-EVENT-1",
					Title:  "Multivariate Event 1",
					Status: "open",
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			ticker:         "INVALID-MV-EVENT",
			serverResponse: MultivariateEventResponse{},
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

				expectedPath := TradeAPIPrefix + "/events/multivariate/" + tt.ticker
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
			resp, err := client.GetMultivariateEvent(context.Background(), tt.ticker)

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

func TestGetEventCandlesticks(t *testing.T) {
	tests := []struct {
		name         string
		params       CandlesticksParams
		rawResponse  string // raw JSON in Kalshi v2 event candlestick format
		serverStatus int
		wantErr      bool
		wantCount    int
	}{
		{
			name: "returns candlesticks successfully",
			params: CandlesticksParams{
				SeriesTicker: "ELECTION-SERIES",
				Ticker:       "ELECTION-2024",
				Period:       "1h",
			},
			rawResponse: `{
				"market_tickers": ["ELECTION-2024"],
				"market_candlesticks": [[
					{"end_period_ts": 1704067200, "price": {"open": 50, "high": 55, "low": 48, "close": 52}, "volume": 1000, "open_interest": 500},
					{"end_period_ts": 1704070800, "price": {"open": 52, "high": 58, "low": 51, "close": 56}, "volume": 1200, "open_interest": 550}
				]]
			}`,
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns candlesticks with time range",
			params: CandlesticksParams{
				SeriesTicker: "FED-SERIES",
				Ticker:       "FED-MAR-2024",
				Period:       "15m",
				StartTime:    func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
				EndTime:      func() *time.Time { t := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC); return &t }(),
			},
			rawResponse: `{
				"market_tickers": ["FED-MAR-2024"],
				"market_candlesticks": [[
					{"end_period_ts": 1704067200, "price": {"open": 45, "high": 50, "low": 44, "close": 49}, "volume": 800, "open_interest": 400}
				]]
			}`,
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name: "returns empty candlesticks for no data",
			params: CandlesticksParams{
				SeriesTicker: "NO-DATA-SERIES",
				Ticker:       "NO-DATA-EVENT",
				Period:       "1d",
			},
			rawResponse:  `{"market_tickers": [], "market_candlesticks": []}`,
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name: "handles not found error",
			params: CandlesticksParams{
				SeriesTicker: "INVALID-SERIES",
				Ticker:       "INVALID-EVENT",
				Period:       "1h",
			},
			rawResponse:  `{}`,
			serverStatus: http.StatusNotFound,
			wantErr:      true,
		},
		{
			name: "handles server error",
			params: CandlesticksParams{
				SeriesTicker: "ERROR-SERIES",
				Ticker:       "ERROR-EVENT",
				Period:       "1h",
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

				expectedPath := TradeAPIPrefix + "/series/" + tt.params.SeriesTicker + "/events/" + tt.params.Ticker + "/candlesticks"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.params.Period != "" {
					expectedInterval := periodToInterval(tt.params.Period)
					if got := r.URL.Query().Get("period_interval"); got != expectedInterval {
						t.Errorf("expected period_interval=%s, got %s", expectedInterval, got)
					}
				}

				if tt.params.StartTime != nil {
					if got := r.URL.Query().Get("start_ts"); got == "" {
						t.Error("expected start_ts query param")
					}
				}

				if tt.params.EndTime != nil {
					if got := r.URL.Query().Get("end_ts"); got == "" {
						t.Error("expected end_ts query param")
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
			candlesticks, err := client.GetEventCandlesticks(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(candlesticks) != tt.wantCount {
				t.Errorf("expected %d candlesticks, got %d", tt.wantCount, len(candlesticks))
			}
		})
	}
}

func TestGetEventMetadata(t *testing.T) {
	tests := []struct {
		name           string
		ticker         string
		serverResponse models.EventMetadataResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "returns event metadata successfully",
			ticker: "ELECTION-2024",
			serverResponse: models.EventMetadataResponse{
				EventMetadata: models.EventMetadata{
					EventTicker: "ELECTION-2024",
					Metadata: map[string]string{
						"source":      "AP News",
						"resolution":  "Official election results",
						"last_update": "2024-11-05T00:00:00Z",
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			ticker:         "INVALID-EVENT",
			serverResponse: models.EventMetadataResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "handles server error",
			ticker:         "ERROR-EVENT",
			serverResponse: models.EventMetadataResponse{},
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

				expectedPath := TradeAPIPrefix + "/events/" + tt.ticker + "/metadata"
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
			resp, err := client.GetEventMetadata(context.Background(), tt.ticker)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.EventTicker != tt.ticker {
				t.Errorf("expected event ticker %q, got %q", tt.ticker, resp.EventTicker)
			}
		})
	}
}

func TestGetForecastPercentileHistory(t *testing.T) {
	tests := []struct {
		name           string
		params         ForecastPercentileHistoryParams
		serverResponse models.ForecastPercentileHistoryResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns forecast history successfully",
			params: ForecastPercentileHistoryParams{
				Ticker: "ELECTION-2024",
			},
			serverResponse: models.ForecastPercentileHistoryResponse{
				History: []models.ForecastPercentilePoint{
					{Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), P10: 40, P25: 45, P50: 50, P75: 55, P90: 60},
					{Timestamp: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), P10: 42, P25: 47, P50: 52, P75: 57, P90: 62},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns forecast history with time range",
			params: ForecastPercentileHistoryParams{
				Ticker:    "FED-MAR-2024",
				StartTime: func() *time.Time { t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
				EndTime:   func() *time.Time { t := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC); return &t }(),
			},
			serverResponse: models.ForecastPercentileHistoryResponse{
				History: []models.ForecastPercentilePoint{
					{Timestamp: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), P10: 35, P25: 40, P50: 45, P75: 50, P90: 55},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name: "handles not found",
			params: ForecastPercentileHistoryParams{
				Ticker: "INVALID-EVENT",
			},
			serverResponse: models.ForecastPercentileHistoryResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "handles server error",
			params: ForecastPercentileHistoryParams{
				Ticker: "ERROR-EVENT",
			},
			serverResponse: models.ForecastPercentileHistoryResponse{},
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

				expectedPath := TradeAPIPrefix + "/events/" + tt.params.Ticker + "/forecast-percentile-history"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.params.StartTime != nil {
					if got := r.URL.Query().Get("start_ts"); got == "" {
						t.Error("expected start_ts query param")
					}
				}

				if tt.params.EndTime != nil {
					if got := r.URL.Query().Get("end_ts"); got == "" {
						t.Error("expected end_ts query param")
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
			history, err := client.GetForecastPercentileHistory(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(history) != tt.wantCount {
				t.Errorf("expected %d history points, got %d", tt.wantCount, len(history))
			}
		})
	}
}
