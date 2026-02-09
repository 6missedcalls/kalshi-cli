package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
					{EventTicker: "ELECTION-2024", Title: "2024 Presidential Election", Status: "open"},
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
					Status:      "open",
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
