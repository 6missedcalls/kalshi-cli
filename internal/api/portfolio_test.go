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

func TestGetBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/balance" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.BalanceResponse{Balance: 100000}
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
}

func TestGetPositions(t *testing.T) {
	expectedPositions := []models.Position{
		{
			Ticker:   "BTC-100K",
			Position: 10,
		},
		{
			Ticker:   "ETH-10K",
			Position: -5,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/positions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.PositionsResponse{
			Positions: expectedPositions,
			Cursor:    "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetPositions(context.Background(), PositionsOptions{})
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}

	if len(result.Positions) != 2 {
		t.Errorf("expected 2 positions, got %d", len(result.Positions))
	}
	if result.Positions[0].Ticker != "BTC-100K" {
		t.Errorf("expected ticker 'BTC-100K', got '%s'", result.Positions[0].Ticker)
	}
}

func TestGetPositionsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("ticker") != "BTC-100K" {
			t.Errorf("expected ticker 'BTC-100K', got '%s'", query.Get("ticker"))
		}
		if query.Get("event_ticker") != "BTC-2024" {
			t.Errorf("expected event_ticker 'BTC-2024', got '%s'", query.Get("event_ticker"))
		}
		if query.Get("limit") != "25" {
			t.Errorf("expected limit '25', got '%s'", query.Get("limit"))
		}

		resp := models.PositionsResponse{Positions: []models.Position{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetPositions(context.Background(), PositionsOptions{
		Ticker:      "BTC-100K",
		EventTicker: "BTC-2024",
		Limit:       25,
	})
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}
}

func TestGetFills(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedFills := []models.Fill{
		{
			TradeID:     "trade-1",
			OrderID:     "order-1",
			Ticker:      "BTC-100K",
			Side:        "yes",
			Count:       5,
			CreatedTime: now,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/fills" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.FillsResponse{
			Fills:  expectedFills,
			Cursor: "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetFills(context.Background(), FillsOptions{})
	if err != nil {
		t.Fatalf("GetFills failed: %v", err)
	}

	if len(result.Fills) != 1 {
		t.Errorf("expected 1 fill, got %d", len(result.Fills))
	}
	if result.Fills[0].TradeID != "trade-1" {
		t.Errorf("expected trade ID 'trade-1', got '%s'", result.Fills[0].TradeID)
	}
}

func TestGetFillsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("ticker") != "BTC-100K" {
			t.Errorf("expected ticker 'BTC-100K', got '%s'", query.Get("ticker"))
		}
		if query.Get("order_id") != "order-123" {
			t.Errorf("expected order_id 'order-123', got '%s'", query.Get("order_id"))
		}

		resp := models.FillsResponse{Fills: []models.Fill{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetFills(context.Background(), FillsOptions{
		Ticker:  "BTC-100K",
		OrderID: "order-123",
	})
	if err != nil {
		t.Fatalf("GetFills failed: %v", err)
	}
}

func TestGetSettlements(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedSettlements := []models.Settlement{
		{
			Ticker:       "BTC-100K",
			MarketResult: "yes",
			Revenue:      500,
			SettledTime:  now,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/settlements" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.SettlementsResponse{
			Settlements: expectedSettlements,
			Cursor:      "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetSettlements(context.Background(), SettlementsOptions{})
	if err != nil {
		t.Fatalf("GetSettlements failed: %v", err)
	}

	if len(result.Settlements) != 1 {
		t.Errorf("expected 1 settlement, got %d", len(result.Settlements))
	}
	if result.Settlements[0].Ticker != "BTC-100K" {
		t.Errorf("expected ticker 'BTC-100K', got '%s'", result.Settlements[0].Ticker)
	}
}

func TestGetSettlementsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("limit") != "100" {
			t.Errorf("expected limit '100', got '%s'", query.Get("limit"))
		}

		resp := models.SettlementsResponse{Settlements: []models.Settlement{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetSettlements(context.Background(), SettlementsOptions{
		Limit: 100,
	})
	if err != nil {
		t.Fatalf("GetSettlements failed: %v", err)
	}
}

func TestGetSubaccounts(t *testing.T) {
	expectedSubaccounts := []models.Subaccount{
		{SubaccountID: 1, Balance: 50000, AvailableBalance: 45000},
		{SubaccountID: 2, Balance: 25000, AvailableBalance: 25000},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/subaccounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.SubaccountsResponse{Subaccounts: expectedSubaccounts}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetSubaccounts(context.Background())
	if err != nil {
		t.Fatalf("GetSubaccounts failed: %v", err)
	}

	if len(result.Subaccounts) != 2 {
		t.Errorf("expected 2 subaccounts, got %d", len(result.Subaccounts))
	}
}

func TestCreateSubaccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/subaccounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.Subaccount{
			SubaccountID:     3,
			Balance:          0,
			AvailableBalance: 0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateSubaccount(context.Background())
	if err != nil {
		t.Fatalf("CreateSubaccount failed: %v", err)
	}

	if result.SubaccountID != 3 {
		t.Errorf("expected subaccount ID 3, got %d", result.SubaccountID)
	}
}

func TestGetTransfers(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// BUG FIX: Correct path per Kalshi API spec
		if r.URL.Path != "/trade-api/v2/portfolio/subaccounts/transfers" {
			t.Errorf("unexpected path: %s, expected /trade-api/v2/portfolio/subaccounts/transfers", r.URL.Path)
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
	if result.Transfers[0].TransferID != "transfer-1" {
		t.Errorf("expected transfer ID 'transfer-1', got '%s'", result.Transfers[0].TransferID)
	}
}

func TestTransfer(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		// BUG FIX: Correct path per Kalshi API spec
		if r.URL.Path != "/trade-api/v2/portfolio/subaccounts/transfers" {
			t.Errorf("unexpected path: %s, expected /trade-api/v2/portfolio/subaccounts/transfers", r.URL.Path)
		}

		var req models.TransferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.FromSubaccount != 1 {
			t.Errorf("expected from_subaccount 1, got %d", req.FromSubaccount)
		}
		if req.ToSubaccount != 2 {
			t.Errorf("expected to_subaccount 2, got %d", req.ToSubaccount)
		}
		if req.Amount != 5000 {
			t.Errorf("expected amount 5000, got %d", req.Amount)
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
	if result.Amount != 5000 {
		t.Errorf("expected amount 5000, got %d", result.Amount)
	}
}

func TestGetSubaccountBalances(t *testing.T) {
	expectedBalances := []models.SubaccountBalance{
		{SubaccountID: 1, Balance: 50000, AvailableBalance: 45000},
		{SubaccountID: 2, Balance: 25000, AvailableBalance: 25000},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/subaccounts/balances" {
			t.Errorf("unexpected path: %s", r.URL.Path)
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
}

func TestGetRestingOrderValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/resting-order-value" {
			t.Errorf("unexpected path: %s", r.URL.Path)
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

func TestGetBalanceFullResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/balance" {
			t.Errorf("unexpected path: %s", r.URL.Path)
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

func TestPortfolioAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
			"code":  "UNAUTHORIZED",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetBalance(context.Background())
	if err == nil {
		t.Fatal("expected error for unauthorized request")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
}
