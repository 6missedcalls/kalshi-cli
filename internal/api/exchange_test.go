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

func TestGetExchangeStatus(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse ExchangeStatusResponse
		serverStatus   int
		wantErr        bool
		wantActive     bool
		wantTrading    bool
	}{
		{
			name: "returns exchange status - all active",
			serverResponse: ExchangeStatusResponse{
				ExchangeActive: true,
				TradingActive:  true,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantActive:   true,
			wantTrading:  true,
		},
		{
			name: "returns exchange status - maintenance",
			serverResponse: ExchangeStatusResponse{
				ExchangeActive: true,
				TradingActive:  false,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantActive:   true,
			wantTrading:  false,
		},
		{
			name: "returns exchange status - offline",
			serverResponse: ExchangeStatusResponse{
				ExchangeActive: false,
				TradingActive:  false,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantActive:   false,
			wantTrading:  false,
		},
		{
			name:           "handles server error",
			serverResponse: ExchangeStatusResponse{},
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

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.GetExchangeStatus(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.ExchangeActive != tt.wantActive {
				t.Errorf("expected exchange_active=%v, got %v", tt.wantActive, resp.ExchangeActive)
			}

			if resp.TradingActive != tt.wantTrading {
				t.Errorf("expected trading_active=%v, got %v", tt.wantTrading, resp.TradingActive)
			}
		})
	}
}

func TestGetExchangeSchedule(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse models.ExchangeScheduleResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "returns schedule successfully",
			serverResponse: models.ExchangeScheduleResponse{
				Schedule: models.ExchangeSchedule{
					StandardHours: []models.WeeklySchedule{
						{
							StartTime: "2024-01-15",
							EndTime:   "2024-06-15",
						},
					},
					MaintenanceWindows: []models.MaintenanceWindow{
						{
							StartDatetime: "2024-01-20T02:00:00Z",
							EndDatetime:   "2024-01-20T04:00:00Z",
						},
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "returns empty schedule",
			serverResponse: models.ExchangeScheduleResponse{
				Schedule: models.ExchangeSchedule{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles server error",
			serverResponse: models.ExchangeScheduleResponse{},
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

				expectedPath := TradeAPIPrefix + "/exchange/schedule"
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
			resp, err := client.GetExchangeSchedule(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp == nil {
				t.Fatal("expected non-nil response")
			}
		})
	}
}

func TestGetAnnouncements(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		serverResponse models.AnnouncementsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns announcements successfully",
			serverResponse: models.AnnouncementsResponse{
				Announcements: []models.Announcement{
					{
						ID:          "ann-1",
						Title:       "Scheduled Maintenance",
						Message:     "The exchange will be down for maintenance.",
						Status:      "active",
						Type:        "maintenance",
						CreatedTime: now,
					},
					{
						ID:          "ann-2",
						Title:       "New Markets Available",
						Message:     "We have added new crypto markets.",
						Status:      "active",
						Type:        "info",
						CreatedTime: now,
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns empty announcements",
			serverResponse: models.AnnouncementsResponse{
				Announcements: []models.Announcement{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: models.AnnouncementsResponse{},
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

				expectedPath := TradeAPIPrefix + "/exchange/announcements"
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
			resp, err := client.GetAnnouncements(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Announcements) != tt.wantCount {
				t.Errorf("expected %d announcements, got %d", tt.wantCount, len(resp.Announcements))
			}
		})
	}
}

func TestGetSeriesFeeChanges(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		serverResponse models.SeriesFeeChangesResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns series fee changes successfully",
			serverResponse: models.SeriesFeeChangesResponse{
				SeriesFeeChanges: []models.SeriesFeeChange{
					{
						SeriesTicker:  "PRES",
						OldFeeRate:    0.05,
						NewFeeRate:    0.03,
						EffectiveDate: now.Add(24 * time.Hour),
						AnnouncedDate: now,
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name: "returns empty series fee changes",
			serverResponse: models.SeriesFeeChangesResponse{
				SeriesFeeChanges: []models.SeriesFeeChange{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: models.SeriesFeeChangesResponse{},
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

				expectedPath := TradeAPIPrefix + "/exchange/series-fee-changes"
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
			resp, err := client.GetSeriesFeeChanges(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.SeriesFeeChanges) != tt.wantCount {
				t.Errorf("expected %d series fee changes, got %d", tt.wantCount, len(resp.SeriesFeeChanges))
			}
		})
	}
}

func TestGetUserDataTimestamp(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		serverResponse models.UserDataTimestampResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "returns user data timestamp successfully",
			serverResponse: models.UserDataTimestampResponse{
				Timestamp: now,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles server error",
			serverResponse: models.UserDataTimestampResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "handles unauthorized error",
			serverResponse: models.UserDataTimestampResponse{},
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

				expectedPath := TradeAPIPrefix + "/exchange/user-data-timestamp"
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
			resp, err := client.GetUserDataTimestamp(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.Timestamp.IsZero() {
				t.Error("expected non-zero timestamp")
			}
		})
	}
}

func TestGetExchangeStatusUsesCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := TradeAPIPrefix + "/exchange/status"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ExchangeStatusResponse{
			ExchangeActive: true,
			TradingActive:  true,
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	_, err := client.GetExchangeStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
