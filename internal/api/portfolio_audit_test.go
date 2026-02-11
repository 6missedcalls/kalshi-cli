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

// TestTransferPathFix verifies the transfer endpoints use the correct path
func TestTransferPathFix(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedTransfers := []models.Transfer{
		{
			TransferID:     "transfer-1",
			FromSubaccount: 1,
			ToSubaccount:   2,
			Amount:         10000,
			CreatedTime:    now,
		},
	}

	t.Run("GetTransfers uses correct subaccounts path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			expectedPath := "/trade-api/v2/portfolio/subaccounts/transfers"
			if r.URL.Path != expectedPath {
				t.Errorf("incorrect path: got %s, want %s", r.URL.Path, expectedPath)
			}

			resp := models.TransfersResponse{Transfers: expectedTransfers}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		result, err := client.GetTransfers(context.Background())
		if err != nil {
			t.Fatalf("GetTransfers failed: %v", err)
		}

		if len(result.Transfers) != 1 {
			t.Errorf("expected 1 transfer, got %d", len(result.Transfers))
		}
	})

	t.Run("Transfer (POST) uses correct subaccounts path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			expectedPath := "/trade-api/v2/portfolio/subaccounts/transfers"
			if r.URL.Path != expectedPath {
				t.Errorf("incorrect path: got %s, want %s", r.URL.Path, expectedPath)
			}

			var req models.TransferRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			resp := models.Transfer{
				TransferID:     "new-transfer",
				FromSubaccount: req.FromSubaccount,
				ToSubaccount:   req.ToSubaccount,
				Amount:         req.Amount,
				CreatedTime:    now,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		result, err := client.Transfer(context.Background(), models.TransferRequest{
			FromSubaccount: 1,
			ToSubaccount:   2,
			Amount:         5000,
		})
		if err != nil {
			t.Fatalf("Transfer failed: %v", err)
		}

		if result.TransferID != "new-transfer" {
			t.Errorf("expected transfer ID 'new-transfer', got '%s'", result.TransferID)
		}
	})
}

// TestGetSubaccountBalances verifies the new endpoint
func TestGetSubaccountBalancesEndpoint(t *testing.T) {
	expectedBalances := []models.SubaccountBalance{
		{SubaccountID: 1, Balance: 50000, AvailableBalance: 45000},
		{SubaccountID: 2, Balance: 25000, AvailableBalance: 25000},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/trade-api/v2/portfolio/subaccounts/balances"
		if r.URL.Path != expectedPath {
			t.Errorf("incorrect path: got %s, want %s", r.URL.Path, expectedPath)
		}

		resp := models.SubaccountBalancesResponse{Balances: expectedBalances}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetSubaccountBalances(context.Background())
	if err != nil {
		t.Fatalf("GetSubaccountBalances failed: %v", err)
	}

	if len(result.Balances) != 2 {
		t.Errorf("expected 2 balances, got %d", len(result.Balances))
	}
	if result.Balances[0].SubaccountID != 1 {
		t.Errorf("expected subaccount ID 1, got %d", result.Balances[0].SubaccountID)
	}
	if result.Balances[0].Balance != 50000 {
		t.Errorf("expected balance 50000, got %d", result.Balances[0].Balance)
	}
}

// TestGetRestingOrderValue verifies the new FCM endpoint
func TestGetRestingOrderValueEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/trade-api/v2/portfolio/resting-order-value"
		if r.URL.Path != expectedPath {
			t.Errorf("incorrect path: got %s, want %s", r.URL.Path, expectedPath)
		}

		resp := models.RestingOrderValueResponse{RestingOrderValue: 150000}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetRestingOrderValue(context.Background())
	if err != nil {
		t.Fatalf("GetRestingOrderValue failed: %v", err)
	}

	if result.RestingOrderValue != 150000 {
		t.Errorf("expected resting order value 150000, got %d", result.RestingOrderValue)
	}
}

// TestBalanceResponseFullFields verifies all balance fields are parsed
func TestBalanceResponseFullFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/trade-api/v2/portfolio/balance"
		if r.URL.Path != expectedPath {
			t.Errorf("incorrect path: got %s, want %s", r.URL.Path, expectedPath)
		}

		resp := models.BalanceResponse{
			Balance:        100000,
			PortfolioValue: 50000,
			UpdatedTs:      1700000000,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetBalance(context.Background())
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}

	if result.Balance != 100000 {
		t.Errorf("expected balance 100000, got %d", result.Balance)
	}
	if result.PortfolioValue != 50000 {
		t.Errorf("expected portfolio_value 50000, got %d", result.PortfolioValue)
	}
	if result.UpdatedTs != 1700000000 {
		t.Errorf("expected updated_ts 1700000000, got %d", result.UpdatedTs)
	}
}

// TestEdgeCases tests edge cases and error handling
func TestPortfolioEdgeCases(t *testing.T) {
	t.Run("empty transfers list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := models.TransfersResponse{Transfers: []models.Transfer{}}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		result, err := client.GetTransfers(context.Background())
		if err != nil {
			t.Fatalf("GetTransfers failed: %v", err)
		}

		if result.Transfers == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(result.Transfers) != 0 {
			t.Errorf("expected 0 transfers, got %d", len(result.Transfers))
		}
	})

	t.Run("empty subaccount balances", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := models.SubaccountBalancesResponse{Balances: []models.SubaccountBalance{}}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		result, err := client.GetSubaccountBalances(context.Background())
		if err != nil {
			t.Fatalf("GetSubaccountBalances failed: %v", err)
		}

		if result.Balances == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(result.Balances) != 0 {
			t.Errorf("expected 0 balances, got %d", len(result.Balances))
		}
	})

	t.Run("zero resting order value", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := models.RestingOrderValueResponse{RestingOrderValue: 0}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		result, err := client.GetRestingOrderValue(context.Background())
		if err != nil {
			t.Fatalf("GetRestingOrderValue failed: %v", err)
		}

		if result.RestingOrderValue != 0 {
			t.Errorf("expected 0, got %d", result.RestingOrderValue)
		}
	})

	t.Run("API error on transfer", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "insufficient_balance",
				"code":  "INSUFFICIENT_BALANCE",
			})
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		_, err := client.Transfer(context.Background(), models.TransferRequest{
			FromSubaccount: 1,
			ToSubaccount:   2,
			Amount:         9999999999,
		})
		if err == nil {
			t.Fatal("expected error for insufficient balance")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != 400 {
			t.Errorf("expected status 400, got %d", apiErr.StatusCode)
		}
	})

	t.Run("API error on resting order value (non-FCM)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "fcm_only",
				"code":  "FCM_ONLY",
			})
		}))
		defer server.Close()

		client := createTestClient(t, server.URL)
		_, err := client.GetRestingOrderValue(context.Background())
		if err == nil {
			t.Fatal("expected error for non-FCM account")
		}

		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != 403 {
			t.Errorf("expected status 403, got %d", apiErr.StatusCode)
		}
	})
}
