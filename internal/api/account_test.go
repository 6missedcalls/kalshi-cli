package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListAPIKeys(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		serverResponse apiKeysResponse
		serverStatus   int
		wantErr        bool
		wantCount      int
	}{
		{
			name: "returns API keys successfully",
			serverResponse: apiKeysResponse{
				APIKeys: []APIKey{
					{
						ID:          "key-1",
						Name:        "Trading Bot",
						CreatedTime: JSONTime{Time: now.Add(-24 * time.Hour)},
						Scopes:      []string{"read", "trade"},
					},
					{
						ID:          "key-2",
						Name:        "Read Only",
						CreatedTime: JSONTime{Time: now.Add(-48 * time.Hour)},
						Scopes:      []string{"read"},
					},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name: "returns empty API keys list",
			serverResponse: apiKeysResponse{
				APIKeys: []APIKey{},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
			wantCount:    0,
		},
		{
			name:           "handles server error",
			serverResponse: apiKeysResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "handles unauthorized",
			serverResponse: apiKeysResponse{},
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

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			keys, err := client.ListAPIKeys(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(keys) != tt.wantCount {
				t.Errorf("expected %d API keys, got %d", tt.wantCount, len(keys))
			}
		})
	}
}

func TestCreateAPIKey(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		request        CreateAPIKeyRequest
		serverResponse CreateAPIKeyResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:    "creates API key successfully",
			request: CreateAPIKeyRequest{Name: "New Trading Bot"},
			serverResponse: CreateAPIKeyResponse{
				APIKey: APIKey{
					ID:          "new-key-id",
					Name:        "New Trading Bot",
					CreatedTime: JSONTime{Time: now},
					Scopes:      []string{"read", "trade"},
				},
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:    "creates API key with empty name",
			request: CreateAPIKeyRequest{Name: ""},
			serverResponse: CreateAPIKeyResponse{
				APIKey: APIKey{
					ID:          "new-key-id",
					Name:        "",
					CreatedTime: JSONTime{Time: now},
					Scopes:      []string{"read"},
				},
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles server error",
			request:        CreateAPIKeyRequest{Name: "Test Key"},
			serverResponse: CreateAPIKeyResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "handles unauthorized",
			request:        CreateAPIKeyRequest{Name: "Test Key"},
			serverResponse: CreateAPIKeyResponse{},
			serverStatus:   http.StatusUnauthorized,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				var req CreateAPIKeyRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				if req.Name != tt.request.Name {
					t.Errorf("expected name %q, got %q", tt.request.Name, req.Name)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			resp, err := client.CreateAPIKey(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.APIKey.Name != tt.request.Name {
				t.Errorf("expected name %q, got %q", tt.request.Name, resp.APIKey.Name)
			}

			if resp.PrivateKey == "" {
				t.Error("expected private key to be returned")
			}
		})
	}
}

func TestDeleteAPIKey(t *testing.T) {
	tests := []struct {
		name         string
		keyID        string
		serverStatus int
		wantErr      bool
	}{
		{
			name:         "deletes API key successfully",
			keyID:        "key-to-delete",
			serverStatus: http.StatusNoContent,
			wantErr:      false,
		},
		{
			name:         "deletes API key successfully with 200",
			keyID:        "key-to-delete",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "handles not found",
			keyID:        "nonexistent-key",
			serverStatus: http.StatusNotFound,
			wantErr:      true,
		},
		{
			name:         "handles server error",
			keyID:        "key-id",
			serverStatus: http.StatusInternalServerError,
			wantErr:      true,
		},
		{
			name:         "handles unauthorized",
			keyID:        "key-id",
			serverStatus: http.StatusUnauthorized,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}

				expectedPath := "/api-keys/" + tt.keyID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := newTestClient(t, server.URL)
			err := client.DeleteAPIKey(context.Background(), tt.keyID)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestGetAPILimits removed - APILimitsResponse and GetAPILimits not implemented
