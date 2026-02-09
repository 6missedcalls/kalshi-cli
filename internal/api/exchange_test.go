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
	now := time.Now().UTC()
	tests := []struct {
		name           string
		serverResponse models.ExchangeScheduleResponse
		serverStatus   int
		wantErr        bool
		wantEntries    int
	}{
		{
			name: "returns schedule successfully",
			serverResponse: models.ExchangeScheduleResponse{
				Schedule: models.ExchangeSchedule{
					ScheduleEntries: []models.ScheduleEntry{
						{
							StartTime:   now,
							EndTime:     now.Add(8 * time.Hour),
							Maintenance: false,
						},
						{
							StartTime:   now.Add(8 * time.Hour),
							EndTime:     now.Add(10 * time.Hour),
							Maintenance: true,
						},
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantEntries:  2,
		},
		{
			name: "returns empty schedule",
			serverResponse: models.ExchangeScheduleResponse{
				Schedule: models.ExchangeSchedule{
					ScheduleEntries: []models.ScheduleEntry{},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantEntries:  0,
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
			var resp models.ExchangeScheduleResponse
			err := client.GetExchangeSchedule(context.Background(), &resp)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Schedule.ScheduleEntries) != tt.wantEntries {
				t.Errorf("expected %d schedule entries, got %d", tt.wantEntries, len(resp.Schedule.ScheduleEntries))
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
			var resp models.AnnouncementsResponse
			err := client.GetAnnouncements(context.Background(), &resp)

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

func TestGetFeeChanges(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		serverResponse models.FeeChangesResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns fee changes successfully",
			serverResponse: models.FeeChangesResponse{
				FeeChanges: []models.FeeChange{
					{
						Ticker:      "BTC-100K",
						OldFee:      5,
						NewFee:      3,
						EffectiveAt: now.Add(24 * time.Hour),
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name: "returns empty fee changes",
			serverResponse: models.FeeChangesResponse{
				FeeChanges: []models.FeeChange{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: models.FeeChangesResponse{},
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

				expectedPath := TradeAPIPrefix + "/exchange/fee-changes"
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
			resp, err := client.GetFeeChanges(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.FeeChanges) != tt.wantCount {
				t.Errorf("expected %d fee changes, got %d", tt.wantCount, len(resp.FeeChanges))
			}
		})
	}
}
