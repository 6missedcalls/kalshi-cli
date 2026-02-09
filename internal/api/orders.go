package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

const ordersBasePath = TradeAPIPrefix + "/portfolio/orders"

// OrdersOptions contains options for listing orders
type OrdersOptions struct {
	Ticker       string
	EventTicker  string
	Status       string
	Cursor       string
	Limit        int
	SubaccountID int
}

// toQueryParams converts OrdersOptions to query parameters
func (o OrdersOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
	if o.Ticker != "" {
		params["ticker"] = o.Ticker
	}
	if o.EventTicker != "" {
		params["event_ticker"] = o.EventTicker
	}
	if o.Status != "" {
		params["status"] = o.Status
	}
	if o.Cursor != "" {
		params["cursor"] = o.Cursor
	}
	if o.Limit > 0 {
		params["limit"] = strconv.Itoa(o.Limit)
	}
	if o.SubaccountID > 0 {
		params["subaccount_id"] = strconv.Itoa(o.SubaccountID)
	}
	return params
}

// GetOrders returns a list of orders based on the provided options
func (c *Client) GetOrders(ctx context.Context, opts OrdersOptions) (*models.OrdersResponse, error) {
	path := ordersBasePath + BuildQueryString(opts.toQueryParams())

	var result models.OrdersResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrder returns a single order by ID
func (c *Client) GetOrder(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	path := ordersBasePath + "/" + orderID

	var result models.OrderResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(ctx context.Context, req models.CreateOrderRequest) (*models.CreateOrderResponse, error) {
	var result models.CreateOrderResponse
	if err := c.PostJSON(ctx, ordersBasePath, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelOrder cancels an existing order
func (c *Client) CancelOrder(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	path := ordersBasePath + "/" + orderID

	var result models.OrderResponse
	if err := c.DeleteJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AmendOrder amends an existing order's price or count
func (c *Client) AmendOrder(ctx context.Context, orderID string, req models.AmendOrderRequest) (*models.OrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	path := ordersBasePath + "/" + orderID

	var result models.OrderResponse
	if err := c.PostJSON(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DecreaseOrder decreases an order's quantity
func (c *Client) DecreaseOrder(ctx context.Context, orderID string, reduceBy int) (*models.OrderResponse, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	path := ordersBasePath + "/" + orderID + "/decrease"
	req := models.DecreaseOrderRequest{ReduceBy: reduceBy}

	var result models.OrderResponse
	if err := c.PostJSON(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchCreateOrders creates multiple orders in a single request
func (c *Client) BatchCreateOrders(ctx context.Context, orders []models.CreateOrderRequest) (*models.BatchCreateOrdersResponse, error) {
	path := ordersBasePath + "/batched"
	req := models.BatchCreateOrdersRequest{Orders: orders}

	var result models.BatchCreateOrdersResponse
	if err := c.PostJSON(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BatchCancelOrders cancels multiple orders in a single request
func (c *Client) BatchCancelOrders(ctx context.Context, req models.BatchCancelOrdersRequest) (*models.BatchCancelOrdersResponse, error) {
	var result models.BatchCancelOrdersResponse
	if err := c.DeleteWithBody(ctx, ordersBasePath, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQueuePosition returns the queue position for a specific order
func (c *Client) GetQueuePosition(ctx context.Context, orderID string) (*models.QueuePosition, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	path := ordersBasePath + "/" + orderID + "/position"

	var result models.QueuePosition
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllQueuePositions returns queue positions for all resting orders
func (c *Client) GetAllQueuePositions(ctx context.Context) (*models.QueuePositionsResponse, error) {
	path := ordersBasePath + "/positions"

	var result models.QueuePositionsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
