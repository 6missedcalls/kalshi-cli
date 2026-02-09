package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

const orderGroupsBasePath = TradeAPIPrefix + "/portfolio/order_groups"

// OrderGroupsOptions contains options for listing order groups
type OrderGroupsOptions struct {
	Status string
	Cursor string
	Limit  int
}

// toQueryParams converts OrderGroupsOptions to query parameters
func (o OrderGroupsOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
	if o.Status != "" {
		params["status"] = o.Status
	}
	if o.Cursor != "" {
		params["cursor"] = o.Cursor
	}
	if o.Limit > 0 {
		params["limit"] = strconv.Itoa(o.Limit)
	}
	return params
}

// GetOrderGroups returns a list of order groups based on the provided options
func (c *Client) GetOrderGroups(ctx context.Context, opts OrderGroupsOptions) (*models.OrderGroupsResponse, error) {
	path := orderGroupsBasePath + BuildQueryString(opts.toQueryParams())

	var result models.OrderGroupsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOrderGroup returns a single order group by ID
func (c *Client) GetOrderGroup(ctx context.Context, groupID string) (*models.OrderGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("order group ID is required")
	}

	path := orderGroupsBasePath + "/" + groupID

	var result models.OrderGroupResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrderGroup creates a new order group with the specified limit
func (c *Client) CreateOrderGroup(ctx context.Context, req models.CreateOrderGroupRequest) (*models.CreateOrderGroupResponse, error) {
	var result models.CreateOrderGroupResponse
	if err := c.PostJSON(ctx, orderGroupsBasePath, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrderGroupLimit updates the fill limit for an order group
// Per Kalshi API spec: PATCH /order-groups/{group_id}/limit
func (c *Client) UpdateOrderGroupLimit(ctx context.Context, groupID string, newLimit int) (*models.OrderGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("order group ID is required")
	}

	path := orderGroupsBasePath + "/" + groupID + "/limit"
	req := models.UpdateOrderGroupLimitRequest{Limit: newLimit}

	var result models.OrderGroupResponse
	if err := c.DoRequest(ctx, "PATCH", path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteOrderGroup deletes an order group
func (c *Client) DeleteOrderGroup(ctx context.Context, groupID string) error {
	if groupID == "" {
		return fmt.Errorf("order group ID is required")
	}

	path := orderGroupsBasePath + "/" + groupID
	return c.DeleteJSON(ctx, path, nil)
}

// ResetOrderGroup resets an order group's filled count
func (c *Client) ResetOrderGroup(ctx context.Context, groupID string) (*models.OrderGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("order group ID is required")
	}

	path := orderGroupsBasePath + "/" + groupID + "/reset"

	var result models.OrderGroupResponse
	if err := c.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TriggerOrderGroup triggers an order group to execute
func (c *Client) TriggerOrderGroup(ctx context.Context, groupID string) (*models.OrderGroupResponse, error) {
	if groupID == "" {
		return nil, fmt.Errorf("order group ID is required")
	}

	path := orderGroupsBasePath + "/" + groupID + "/trigger"

	var result models.OrderGroupResponse
	if err := c.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
