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

func TestGetMilestone(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		milestoneID    string
		serverResponse models.MilestoneResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:        "returns milestone successfully",
			milestoneID: "milestone-123",
			serverResponse: models.MilestoneResponse{
				Milestone: models.Milestone{
					ID:          "milestone-123",
					Title:       "Q1 2024 GDP Report",
					Description: "Quarterly GDP growth rate",
					Status:      "active",
					Category:    "economics",
					TargetDate:  now.Add(30 * 24 * time.Hour),
					CreatedTime: now.Add(-7 * 24 * time.Hour),
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			milestoneID:    "invalid-id",
			serverResponse: models.MilestoneResponse{},
			serverStatus:   http.StatusNotFound,
			wantErr:        true,
		},
		{
			name:           "handles server error",
			milestoneID:    "milestone-error",
			serverResponse: models.MilestoneResponse{},
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

				expectedPath := TradeAPIPrefix + "/milestones/" + tt.milestoneID
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

			client := milestonesTestClient(t, server.URL)
			resp, err := client.GetMilestone(context.Background(), tt.milestoneID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.ID != tt.milestoneID {
				t.Errorf("expected milestone ID %q, got %q", tt.milestoneID, resp.ID)
			}
		})
	}
}

func TestListMilestones(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		params         ListMilestonesParams
		serverResponse models.MilestonesResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns milestones successfully",
			params: ListMilestonesParams{},
			serverResponse: models.MilestonesResponse{
				Milestones: []models.Milestone{
					{
						ID:    "milestone-1",
						Title: "GDP Report",
					},
					{
						ID:    "milestone-2",
						Title: "Jobs Report",
					},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns milestones with date filter",
			params: ListMilestonesParams{
				MinDate: now.Add(-30 * 24 * time.Hour),
				MaxDate: now,
			},
			serverResponse: models.MilestonesResponse{
				Milestones: []models.Milestone{
					{ID: "milestone-recent", Title: "Recent Milestone"},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles server error",
			params:         ListMilestonesParams{},
			serverResponse: models.MilestonesResponse{},
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

				expectedPath := TradeAPIPrefix + "/milestones"
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

			client := milestonesTestClient(t, server.URL)
			resp, err := client.ListMilestones(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Milestones) != tt.wantCount {
				t.Errorf("expected %d milestones, got %d", tt.wantCount, len(resp.Milestones))
			}
		})
	}
}

func TestGetLiveData(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		milestoneID    string
		serverResponse models.LiveDataResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:        "returns live data successfully",
			milestoneID: "milestone-123",
			serverResponse: models.LiveDataResponse{
				Data: models.LiveData{
					MilestoneID: "milestone-123",
					Value:       3.2,
					Unit:        "percent",
					Source:      "BEA",
					Timestamp:   now,
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			milestoneID:    "invalid-id",
			serverResponse: models.LiveDataResponse{},
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

				expectedPath := TradeAPIPrefix + "/live-data/" + tt.milestoneID
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

			client := milestonesTestClient(t, server.URL)
			resp, err := client.GetLiveData(context.Background(), tt.milestoneID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.MilestoneID != tt.milestoneID {
				t.Errorf("expected milestone ID %q, got %q", tt.milestoneID, resp.MilestoneID)
			}
		})
	}
}

func TestGetBatchLiveData(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		milestoneIDs   []string
		serverResponse models.BatchLiveDataResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:         "returns batch live data successfully",
			milestoneIDs: []string{"milestone-1", "milestone-2"},
			serverResponse: models.BatchLiveDataResponse{
				Data: []models.LiveData{
					{MilestoneID: "milestone-1", Value: 3.2, Timestamp: now},
					{MilestoneID: "milestone-2", Value: 4.5, Timestamp: now},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:           "handles empty request",
			milestoneIDs:   []string{},
			serverResponse: models.BatchLiveDataResponse{Data: []models.LiveData{}},
			serverStatus:   http.StatusOK,
			wantErr:        false,
			wantCount:      0,
		},
		{
			name:           "handles server error",
			milestoneIDs:   []string{"error-id"},
			serverResponse: models.BatchLiveDataResponse{},
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

				expectedPath := TradeAPIPrefix + "/live-data"
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

			client := milestonesTestClient(t, server.URL)
			resp, err := client.GetBatchLiveData(context.Background(), tt.milestoneIDs)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Data) != tt.wantCount {
				t.Errorf("expected %d data items, got %d", tt.wantCount, len(resp.Data))
			}
		})
	}
}

// milestonesTestClient creates a test client for milestones tests
func milestonesTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client := NewClient(nil, nil)
	client.SetBaseURL(serverURL)
	return client
}
