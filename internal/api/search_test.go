package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// =============================================================================
// TDD Step 1: Write FAILING tests FIRST (RED)
// =============================================================================

func TestGetSportsFilters(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse models.SportsFiltersResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns sports filters successfully",
			serverResponse: models.SportsFiltersResponse{
				Filters: []models.SportsFilter{
					{
						ID:       "nfl-teams",
						Name:     "NFL Teams",
						Sport:    "football",
						League:   "NFL",
						Category: "team",
					},
					{
						ID:       "nba-teams",
						Name:     "NBA Teams",
						Sport:    "basketball",
						League:   "NBA",
						Category: "team",
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns empty filters",
			serverResponse: models.SportsFiltersResponse{
				Filters: []models.SportsFilter{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: models.SportsFiltersResponse{},
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

				expectedPath := TradeAPIPrefix + "/search/sports/filters"
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

			client := searchTestClient(t, server.URL)
			resp, err := client.GetSportsFilters(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Filters) != tt.wantCount {
				t.Errorf("expected %d filters, got %d", tt.wantCount, len(resp.Filters))
			}
		})
	}
}

func TestGetSearchTags(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse models.TagsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns tags successfully",
			serverResponse: models.TagsResponse{
				Mappings: []models.TagMapping{
					{
						Category: "politics",
						Tags:     []string{"election", "congress", "senate"},
					},
					{
						Category: "economics",
						Tags:     []string{"gdp", "inflation", "employment"},
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns empty tags",
			serverResponse: models.TagsResponse{
				Mappings: []models.TagMapping{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: models.TagsResponse{},
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

				expectedPath := TradeAPIPrefix + "/search/tags"
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

			client := searchTestClient(t, server.URL)
			resp, err := client.GetSearchTags(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Mappings) != tt.wantCount {
				t.Errorf("expected %d mappings, got %d", tt.wantCount, len(resp.Mappings))
			}
		})
	}
}

func TestGetStructuredTarget(t *testing.T) {
	tests := []struct {
		name           string
		targetID       string
		serverResponse models.StructuredTargetResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:     "returns target successfully",
			targetID: "target-123",
			serverResponse: models.StructuredTargetResponse{
				Target: models.StructuredTarget{
					ID:          "target-123",
					Name:        "Bitcoin Price Target",
					Description: "Bitcoin reaches $100,000",
					Type:        "price_target",
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles not found",
			targetID:       "invalid-id",
			serverResponse: models.StructuredTargetResponse{},
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

				expectedPath := TradeAPIPrefix + "/structured-targets/" + tt.targetID
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

			client := searchTestClient(t, server.URL)
			resp, err := client.GetStructuredTarget(context.Background(), tt.targetID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.ID != tt.targetID {
				t.Errorf("expected target ID %q, got %q", tt.targetID, resp.ID)
			}
		})
	}
}

func TestListStructuredTargets(t *testing.T) {
	tests := []struct {
		name           string
		params         ListStructuredTargetsParams
		serverResponse models.StructuredTargetsResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name:   "returns targets successfully",
			params: ListStructuredTargetsParams{Limit: 100},
			serverResponse: models.StructuredTargetsResponse{
				Targets: []models.StructuredTarget{
					{ID: "target-1", Name: "Target 1"},
					{ID: "target-2", Name: "Target 2"},
				},
				Cursor: "next-cursor",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:   "returns targets with pagination",
			params: ListStructuredTargetsParams{Cursor: "page-2", Limit: 50},
			serverResponse: models.StructuredTargetsResponse{
				Targets: []models.StructuredTarget{
					{ID: "target-3", Name: "Target 3"},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    1,
		},
		{
			name:           "handles server error",
			params:         ListStructuredTargetsParams{},
			serverResponse: models.StructuredTargetsResponse{},
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

				expectedPath := TradeAPIPrefix + "/structured-targets"
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

			client := searchTestClient(t, server.URL)
			resp, err := client.ListStructuredTargets(context.Background(), tt.params)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Targets) != tt.wantCount {
				t.Errorf("expected %d targets, got %d", tt.wantCount, len(resp.Targets))
			}
		})
	}
}

func TestGetIncentives(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse models.IncentivesResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns incentives successfully",
			serverResponse: models.IncentivesResponse{
				Incentives: []models.Incentive{
					{
						ID:          "incentive-1",
						Name:        "Referral Bonus",
						Description: "Get $20 for referring a friend",
						Type:        "referral",
						Value:       20.0,
						Status:      "active",
					},
					{
						ID:          "incentive-2",
						Name:        "Welcome Bonus",
						Description: "New user welcome bonus",
						Type:        "signup",
						Value:       10.0,
						Status:      "active",
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns empty incentives",
			serverResponse: models.IncentivesResponse{
				Incentives: []models.Incentive{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: models.IncentivesResponse{},
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

				expectedPath := TradeAPIPrefix + "/incentives"
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

			client := searchTestClient(t, server.URL)
			resp, err := client.GetIncentives(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(resp.Incentives) != tt.wantCount {
				t.Errorf("expected %d incentives, got %d", tt.wantCount, len(resp.Incentives))
			}
		})
	}
}

// searchTestClient creates a test client for search tests
func searchTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client := NewClient(nil, nil)
	client.SetBaseURL(serverURL)
	return client
}
