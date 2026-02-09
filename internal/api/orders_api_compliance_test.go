package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// TestAmendOrderUsesCorrectHTTPMethod verifies that AmendOrder uses PATCH method
// per Kalshi API spec: PATCH /orders/{order_id}
func TestAmendOrderUsesCorrectHTTPMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API spec requires PATCH, not POST or PUT
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH method, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/trade-api/v2/portfolio/orders/order-to-amend" {
			t.Errorf("unexpected path: %s", r.URL.Path)
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
	_, err := client.AmendOrder(context.Background(), "order-to-amend", models.AmendOrderRequest{
		Price: 55,
	})
	if err != nil {
		t.Fatalf("AmendOrder failed: %v", err)
	}
}

// TestDecreaseOrderUsesCorrectHTTPMethod verifies that DecreaseOrder uses PATCH method
// per Kalshi API spec: PATCH /orders/{order_id}/decrease
func TestDecreaseOrderUsesCorrectHTTPMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API spec requires PATCH, not POST
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH method, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
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

// TestBatchCreateOrdersUsesCorrectPath verifies the correct endpoint path
// per Kalshi API spec: POST /orders/batch (not /orders/batched)
func TestBatchCreateOrdersUsesCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		// API spec uses /batch, not /batched
		expectedPath := "/trade-api/v2/portfolio/orders/batch"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req models.BatchCreateOrdersRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
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

// TestBatchCancelOrdersUsesCorrectPath verifies the correct endpoint path
// per Kalshi API spec: DELETE /orders/batch (not /orders)
func TestBatchCancelOrdersUsesCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		// API spec uses /batch for batch operations
		expectedPath := "/trade-api/v2/portfolio/orders/batch"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
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

// TestGetQueuePositionUsesCorrectPath verifies the correct endpoint path
// per Kalshi API spec: GET /orders/{order_id}/queue-position (not /position)
func TestGetQueuePositionUsesCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// API spec uses /queue-position, not /position
		expectedPath := "/trade-api/v2/portfolio/orders/order-123/queue-position"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
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

// TestGetAllQueuePositionsUsesCorrectPath verifies the correct endpoint path
// per Kalshi API spec: GET /orders/queue-positions (not /positions)
func TestGetAllQueuePositionsUsesCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// API spec uses /queue-positions, not /positions
		expectedPath := "/trade-api/v2/portfolio/orders/queue-positions"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
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

// TestBatchCreateOrdersLimitValidation verifies the 20 order limit for batch creation
// per Kalshi API spec: max 20 orders per batch
func TestBatchCreateOrdersLimitValidation(t *testing.T) {
	client := createTestClient(t, "http://localhost")

	// Create 21 orders (exceeds limit)
	orders := make([]models.CreateOrderRequest, 21)
	for i := range orders {
		orders[i] = models.CreateOrderRequest{
			Ticker: "TEST-TICKER",
			Side:   models.OrderSideYes,
			Count:  1,
		}
	}

	_, err := client.BatchCreateOrders(context.Background(), orders)
	if err == nil {
		t.Error("expected error for exceeding batch limit of 20 orders")
	}
}

// TestBatchCancelOrdersLimitValidation verifies the 20 order limit for batch cancellation
// per Kalshi API spec: max 20 orders per batch
func TestBatchCancelOrdersLimitValidation(t *testing.T) {
	client := createTestClient(t, "http://localhost")

	// Create 21 order IDs (exceeds limit)
	orderIDs := make([]string, 21)
	for i := range orderIDs {
		orderIDs[i] = "order-" + string(rune('a'+i))
	}

	_, err := client.BatchCancelOrders(context.Background(), models.BatchCancelOrdersRequest{
		OrderIDs: orderIDs,
	})
	if err == nil {
		t.Error("expected error for exceeding batch limit of 20 orders")
	}
}

// TestDecreaseOrderValidatesPositiveAmount ensures reduce_by must be positive
func TestDecreaseOrderValidatesPositiveAmount(t *testing.T) {
	client := createTestClient(t, "http://localhost")

	_, err := client.DecreaseOrder(context.Background(), "order-123", 0)
	if err == nil {
		t.Error("expected error for zero reduce_by amount")
	}

	_, err = client.DecreaseOrder(context.Background(), "order-123", -5)
	if err == nil {
		t.Error("expected error for negative reduce_by amount")
	}
}

// TestAmendOrderValidatesAtLeastOneField ensures at least price or count is specified
func TestAmendOrderValidatesAtLeastOneField(t *testing.T) {
	client := createTestClient(t, "http://localhost")

	_, err := client.AmendOrder(context.Background(), "order-123", models.AmendOrderRequest{})
	if err == nil {
		t.Error("expected error when neither price nor count is specified")
	}
}
