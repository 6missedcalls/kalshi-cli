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

			client := accountTestClient(t, server.URL)
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

			client := accountTestClient(t, server.URL)
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

			client := accountTestClient(t, server.URL)
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

func TestGetAPILimits(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse APILimitsResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "returns API limits successfully",
			serverResponse: APILimitsResponse{
				RateLimit:        100,
				MaxOrdersPerCall: 50,
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles server error",
			serverResponse: APILimitsResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "handles unauthorized",
			serverResponse: APILimitsResponse{},
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

				expectedPath := "/account/api-limits"
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

			client := accountTestClient(t, server.URL)
			limits, err := client.GetAPILimits(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if limits.RateLimit != tt.serverResponse.RateLimit {
				t.Errorf("expected rate limit %d, got %d", tt.serverResponse.RateLimit, limits.RateLimit)
			}

			if limits.MaxOrdersPerCall != tt.serverResponse.MaxOrdersPerCall {
				t.Errorf("expected max orders per call %d, got %d", tt.serverResponse.MaxOrdersPerCall, limits.MaxOrdersPerCall)
			}
		})
	}
}

func TestGenerateAPIKey(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name           string
		request        GenerateAPIKeyRequest
		serverResponse GenerateAPIKeyResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:    "generates API key pair successfully",
			request: GenerateAPIKeyRequest{Name: "Generated Key"},
			serverResponse: GenerateAPIKeyResponse{
				APIKey: APIKey{
					ID:          "generated-key-id",
					Name:        "Generated Key",
					CreatedTime: JSONTime{Time: now},
					Scopes:      []string{"read", "trade"},
				},
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "handles server error",
			request:        GenerateAPIKeyRequest{Name: "Test Key"},
			serverResponse: GenerateAPIKeyResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "handles unauthorized",
			request:        GenerateAPIKeyRequest{Name: "Test Key"},
			serverResponse: GenerateAPIKeyResponse{},
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

				expectedPath := "/api-keys/generate"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				var req GenerateAPIKeyRequest
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

			client := accountTestClient(t, server.URL)
			resp, err := client.GenerateAPIKey(context.Background(), tt.request)

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

func TestCreateAPIKeyWithPublicKey(t *testing.T) {
	now := time.Now().UTC()
	testPublicKey := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBg...\n-----END PUBLIC KEY-----"

	tests := []struct {
		name           string
		request        CreateAPIKeyWithPublicKeyRequest
		serverResponse CreateAPIKeyWithPublicKeyResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "creates API key with public key successfully",
			request: CreateAPIKeyWithPublicKeyRequest{
				Name:      "My Key",
				PublicKey: testPublicKey,
			},
			serverResponse: CreateAPIKeyWithPublicKeyResponse{
				APIKey: APIKey{
					ID:          "custom-key-id",
					Name:        "My Key",
					CreatedTime: JSONTime{Time: now},
					Scopes:      []string{"read", "trade"},
				},
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "handles missing public key error",
			request: CreateAPIKeyWithPublicKeyRequest{
				Name:      "My Key",
				PublicKey: "",
			},
			serverResponse: CreateAPIKeyWithPublicKeyResponse{},
			serverStatus:   http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "handles server error",
			request: CreateAPIKeyWithPublicKeyRequest{
				Name:      "Test Key",
				PublicKey: testPublicKey,
			},
			serverResponse: CreateAPIKeyWithPublicKeyResponse{},
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "handles unauthorized",
			request: CreateAPIKeyWithPublicKeyRequest{
				Name:      "Test Key",
				PublicKey: testPublicKey,
			},
			serverResponse: CreateAPIKeyWithPublicKeyResponse{},
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

				expectedPath := "/api-keys"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				var req CreateAPIKeyWithPublicKeyRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				if req.Name != tt.request.Name {
					t.Errorf("expected name %q, got %q", tt.request.Name, req.Name)
				}

				if req.PublicKey != tt.request.PublicKey {
					t.Errorf("expected public key %q, got %q", tt.request.PublicKey, req.PublicKey)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			client := accountTestClient(t, server.URL)
			resp, err := client.CreateAPIKeyWithPublicKey(context.Background(), tt.request)

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
		})
	}
}

// accountTestClient creates a test client with the given base URL for account tests
func accountTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client := NewClient(nil, nil)
	client.SetBaseURL(serverURL)
	return client
}
