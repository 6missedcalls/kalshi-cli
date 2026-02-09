package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

func TestGetOrders(t *testing.T) {
	expectedOrders := []models.Order{
		{
			OrderID: "order-1",
			Ticker:  "BTC-100K",
			Status:  models.OrderStatusResting,
			Side:    models.OrderSideYes,
		},
		{
			OrderID: "order-2",
			Ticker:  "ETH-10K",
			Status:  models.OrderStatusExecuted,
			Side:    models.OrderSideNo,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.OrdersResponse{
			Orders: expectedOrders,
			Cursor: "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetOrders(context.Background(), OrdersOptions{})
	if err != nil {
		t.Fatalf("GetOrders failed: %v", err)
	}

	if len(result.Orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result.Orders))
	}
	if result.Orders[0].OrderID != "order-1" {
		t.Errorf("expected order ID 'order-1', got '%s'", result.Orders[0].OrderID)
	}
	if result.Cursor != "next-cursor" {
		t.Errorf("expected cursor 'next-cursor', got '%s'", result.Cursor)
	}
}

func TestGetOrdersWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("ticker") != "BTC-100K" {
			t.Errorf("expected ticker 'BTC-100K', got '%s'", query.Get("ticker"))
		}
		if query.Get("status") != "resting" {
			t.Errorf("expected status 'resting', got '%s'", query.Get("status"))
		}
		if query.Get("limit") != "50" {
			t.Errorf("expected limit '50', got '%s'", query.Get("limit"))
		}

		resp := models.OrdersResponse{Orders: []models.Order{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetOrders(context.Background(), OrdersOptions{
		Ticker: "BTC-100K",
		Status: "resting",
		Limit:  50,
	})
	if err != nil {
		t.Fatalf("GetOrders failed: %v", err)
	}
}

func TestGetOrder(t *testing.T) {
	expectedOrder := models.Order{
		OrderID: "order-123",
		Ticker:  "BTC-100K",
		Status:  models.OrderStatusResting,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/order-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.OrderResponse{Order: expectedOrder}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetOrder(context.Background(), "order-123")
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}

	if result.Order.OrderID != "order-123" {
		t.Errorf("expected order ID 'order-123', got '%s'", result.Order.OrderID)
	}
}

func TestCreateOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Ticker != "BTC-100K" {
			t.Errorf("expected ticker 'BTC-100K', got '%s'", req.Ticker)
		}
		if req.Side != models.OrderSideYes {
			t.Errorf("expected side 'yes', got '%s'", req.Side)
		}
		if req.Count != 10 {
			t.Errorf("expected count 10, got %d", req.Count)
		}

		resp := models.CreateOrderResponse{
			Order: models.Order{
				OrderID: "new-order-id",
				Ticker:  req.Ticker,
				Side:    req.Side,
				Status:  models.OrderStatusResting,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateOrder(context.Background(), models.CreateOrderRequest{
		Ticker:   "BTC-100K",
		Side:     models.OrderSideYes,
		Action:   models.OrderActionBuy,
		Type:     models.OrderTypeLimit,
		Count:    10,
		YesPrice: 50,
	})
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	if result.Order.OrderID != "new-order-id" {
		t.Errorf("expected order ID 'new-order-id', got '%s'", result.Order.OrderID)
	}
}

func TestCancelOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/order-to-cancel" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.OrderResponse{
			Order: models.Order{
				OrderID: "order-to-cancel",
				Status:  models.OrderStatusCanceled,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CancelOrder(context.Background(), "order-to-cancel")
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}

	if result.Order.Status != models.OrderStatusCanceled {
		t.Errorf("expected status 'canceled', got '%s'", result.Order.Status)
	}
}

func TestAmendOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/order-to-amend" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.AmendOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Price != 55 {
			t.Errorf("expected price 55, got %d", req.Price)
		}

		resp := models.OrderResponse{
			Order: models.Order{
				OrderID:  "order-to-amend",
				YesPrice: 55,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.AmendOrder(context.Background(), "order-to-amend", models.AmendOrderRequest{
		Price: 55,
	})
	if err != nil {
		t.Fatalf("AmendOrder failed: %v", err)
	}

	if result.Order.YesPrice != 55 {
		t.Errorf("expected price 55, got %d", result.Order.YesPrice)
	}
}

func TestDecreaseOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/order-to-decrease/decrease" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.DecreaseOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.ReduceBy != 5 {
			t.Errorf("expected reduce_by 5, got %d", req.ReduceBy)
		}

		resp := models.OrderResponse{
			Order: models.Order{
				OrderID:           "order-to-decrease",
				RemainingQuantity: 5,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.DecreaseOrder(context.Background(), "order-to-decrease", 5)
	if err != nil {
		t.Fatalf("DecreaseOrder failed: %v", err)
	}

	if result.Order.RemainingQuantity != 5 {
		t.Errorf("expected remaining quantity 5, got %d", result.Order.RemainingQuantity)
	}
}

func TestBatchCreateOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/batched" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.BatchCreateOrdersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Orders) != 2 {
			t.Errorf("expected 2 orders, got %d", len(req.Orders))
		}

		resp := models.BatchCreateOrdersResponse{
			Orders: []models.Order{
				{OrderID: "batch-order-1"},
				{OrderID: "batch-order-2"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.BatchCreateOrders(context.Background(), []models.CreateOrderRequest{
		{Ticker: "BTC-100K", Side: models.OrderSideYes, Count: 10},
		{Ticker: "ETH-10K", Side: models.OrderSideNo, Count: 20},
	})
	if err != nil {
		t.Fatalf("BatchCreateOrders failed: %v", err)
	}

	if len(result.Orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result.Orders))
	}
}

func TestBatchCancelOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.BatchCancelOrdersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		resp := models.BatchCancelOrdersResponse{
			Orders: []models.Order{
				{OrderID: "cancel-1", Status: models.OrderStatusCanceled},
				{OrderID: "cancel-2", Status: models.OrderStatusCanceled},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.BatchCancelOrders(context.Background(), models.BatchCancelOrdersRequest{
		OrderIDs: []string{"cancel-1", "cancel-2"},
	})
	if err != nil {
		t.Fatalf("BatchCancelOrders failed: %v", err)
	}

	if len(result.Orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result.Orders))
	}
}

func TestGetQueuePosition(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/order-123/position" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.QueuePosition{
			OrderID:       "order-123",
			QueuePosition: 5,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetQueuePosition(context.Background(), "order-123")
	if err != nil {
		t.Fatalf("GetQueuePosition failed: %v", err)
	}

	if result.QueuePosition != 5 {
		t.Errorf("expected queue position 5, got %d", result.QueuePosition)
	}
}

func TestGetAllQueuePositions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/positions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.QueuePositionsResponse{
			Positions: []models.QueuePosition{
				{OrderID: "order-1", QueuePosition: 3},
				{OrderID: "order-2", QueuePosition: 7},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetAllQueuePositions(context.Background())
	if err != nil {
		t.Fatalf("GetAllQueuePositions failed: %v", err)
	}

	if len(result.Positions) != 2 {
		t.Errorf("expected 2 positions, got %d", len(result.Positions))
	}
}

func TestOrdersAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Order not found",
			"code":  "ORDER_NOT_FOUND",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetOrder(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent order")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}

func createTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("test-api-key", privateKey)
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	return NewClientLegacy(signer, WithBaseURL(baseURL))
}
