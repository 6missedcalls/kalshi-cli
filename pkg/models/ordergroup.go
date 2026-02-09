package models

import "time"

// OrderGroup represents an order group
type OrderGroup struct {
	GroupID         string    `json:"order_group_id"`
	Status          string    `json:"status"`
	Limit           int       `json:"limit"`
	FilledCount     int       `json:"filled_count"`
	OrderCount      int       `json:"order_count"`
	OrderIDs        []string  `json:"order_ids"`
	CreatedTime     time.Time `json:"created_time"`
	LastUpdateTime  time.Time `json:"last_update_time"`
}

// OrderGroupsResponse is the API response for order groups
type OrderGroupsResponse struct {
	OrderGroups []OrderGroup `json:"order_groups"`
	Cursor      string       `json:"cursor"`
}

// OrderGroupResponse is the API response for a single order group
type OrderGroupResponse struct {
	OrderGroup OrderGroup `json:"order_group"`
}

// CreateOrderGroupRequest is the request to create an order group
type CreateOrderGroupRequest struct {
	Limit int `json:"limit"`
}

// CreateOrderGroupResponse is the response from creating an order group
type CreateOrderGroupResponse struct {
	OrderGroup OrderGroup `json:"order_group"`
}

// UpdateOrderGroupLimitRequest is the request to update order group limit
type UpdateOrderGroupLimitRequest struct {
	Limit int `json:"limit"`
}
