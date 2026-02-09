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

func TestGetOrderGroups(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedGroups := []models.OrderGroup{
		{
			GroupID:        "group-1",
			Status:         "active",
			Limit:          100,
			FilledCount:    50,
			OrderCount:     3,
			OrderIDs:       []string{"order-1", "order-2", "order-3"},
			CreatedTime:    now,
			LastUpdateTime: now,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/order_groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.OrderGroupsResponse{
			OrderGroups: expectedGroups,
			Cursor:      "next-cursor",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetOrderGroups(context.Background(), OrderGroupsOptions{})
	if err != nil {
		t.Fatalf("GetOrderGroups failed: %v", err)
	}

	if len(result.OrderGroups) != 1 {
		t.Errorf("expected 1 order group, got %d", len(result.OrderGroups))
	}
	if result.OrderGroups[0].GroupID != "group-1" {
		t.Errorf("expected group ID 'group-1', got '%s'", result.OrderGroups[0].GroupID)
	}
}

func TestGetOrderGroupsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if query.Get("status") != "active" {
			t.Errorf("expected status 'active', got '%s'", query.Get("status"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("expected limit '10', got '%s'", query.Get("limit"))
		}

		resp := models.OrderGroupsResponse{OrderGroups: []models.OrderGroup{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetOrderGroups(context.Background(), OrderGroupsOptions{
		Status: "active",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("GetOrderGroups failed: %v", err)
	}
}

func TestGetOrderGroup(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expectedGroup := models.OrderGroup{
		GroupID:        "group-123",
		Status:         "active",
		Limit:          50,
		FilledCount:    25,
		OrderCount:     2,
		CreatedTime:    now,
		LastUpdateTime: now,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/order_groups/group-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := models.OrderGroupResponse{OrderGroup: expectedGroup}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.GetOrderGroup(context.Background(), "group-123")
	if err != nil {
		t.Fatalf("GetOrderGroup failed: %v", err)
	}

	if result.OrderGroup.GroupID != "group-123" {
		t.Errorf("expected group ID 'group-123', got '%s'", result.OrderGroup.GroupID)
	}
}

func TestCreateOrderGroup(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/order_groups" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.CreateOrderGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Limit != 100 {
			t.Errorf("expected limit 100, got %d", req.Limit)
		}

		resp := models.CreateOrderGroupResponse{
			OrderGroup: models.OrderGroup{
				GroupID:        "new-group-id",
				Limit:          req.Limit,
				Status:         "active",
				CreatedTime:    now,
				LastUpdateTime: now,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.CreateOrderGroup(context.Background(), models.CreateOrderGroupRequest{
		Limit: 100,
	})
	if err != nil {
		t.Fatalf("CreateOrderGroup failed: %v", err)
	}

	if result.OrderGroup.GroupID != "new-group-id" {
		t.Errorf("expected group ID 'new-group-id', got '%s'", result.OrderGroup.GroupID)
	}
	if result.OrderGroup.Limit != 100 {
		t.Errorf("expected limit 100, got %d", result.OrderGroup.Limit)
	}
}

func TestUpdateOrderGroupLimit(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/order_groups/group-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req models.UpdateOrderGroupLimitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Limit != 200 {
			t.Errorf("expected limit 200, got %d", req.Limit)
		}

		resp := models.OrderGroupResponse{
			OrderGroup: models.OrderGroup{
				GroupID:        "group-123",
				Limit:          req.Limit,
				Status:         "active",
				CreatedTime:    now,
				LastUpdateTime: now,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	result, err := client.UpdateOrderGroupLimit(context.Background(), "group-123", 200)
	if err != nil {
		t.Fatalf("UpdateOrderGroupLimit failed: %v", err)
	}

	if result.OrderGroup.Limit != 200 {
		t.Errorf("expected limit 200, got %d", result.OrderGroup.Limit)
	}
}

func TestDeleteOrderGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/trade-api/v2/portfolio/order_groups/group-to-delete" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	err := client.DeleteOrderGroup(context.Background(), "group-to-delete")
	if err != nil {
		t.Fatalf("DeleteOrderGroup failed: %v", err)
	}
}

func TestOrderGroupsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Order group not found",
			"code":  "ORDER_GROUP_NOT_FOUND",
		})
	}))
	defer server.Close()

	client := createTestClient(t, server.URL)
	_, err := client.GetOrderGroup(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent order group")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
}
