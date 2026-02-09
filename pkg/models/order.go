package models

import "time"

// OrderSide represents buy/sell side
type OrderSide string

const (
	OrderSideYes OrderSide = "yes"
	OrderSideNo  OrderSide = "no"
)

// OrderType represents order type
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusResting   OrderStatus = "resting"
	OrderStatusCanceled  OrderStatus = "canceled"
	OrderStatusExecuted  OrderStatus = "executed"
	OrderStatusPending   OrderStatus = "pending"
)

// OrderAction represents order action
type OrderAction string

const (
	OrderActionBuy  OrderAction = "buy"
	OrderActionSell OrderAction = "sell"
)

// Order represents a trading order
type Order struct {
	OrderID              string      `json:"order_id"`
	UserID               string      `json:"user_id"`
	Ticker               string      `json:"ticker"`
	Status               OrderStatus `json:"status"`
	YesPrice             int         `json:"yes_price"`
	NoPrice              int         `json:"no_price"`
	Type                 OrderType   `json:"type"`
	Side                 OrderSide   `json:"side"`
	Action               OrderAction `json:"action"`
	InitialQuantity      int         `json:"initial_quantity"`
	RemainingQuantity    int         `json:"remaining_quantity"`
	FilledQuantity       int         `json:"filled_quantity"`
	AverageFillPrice     int         `json:"average_fill_price"`
	ExpirationTime       *time.Time  `json:"expiration_time,omitempty"`
	CreatedTime          time.Time   `json:"created_time"`
	LastUpdateTime       time.Time   `json:"last_update_time"`
	OrderGroupID         string      `json:"order_group_id,omitempty"`
	DecreaseQty          int         `json:"decrease_quantity,omitempty"`
	TakerFillCount       int         `json:"taker_fill_count"`
	TakerFillCost        int         `json:"taker_fill_cost"`
	TakerFees            int         `json:"taker_fees"`
	MakerFillCount       int         `json:"maker_fill_count"`
	MakerFillCost        int         `json:"maker_fill_cost"`
	MakerFees            int         `json:"maker_fees"`
	ClientOrderID        string      `json:"client_order_id,omitempty"`
	SubaccountID         int         `json:"subaccount_id,omitempty"`
}

// OrderResponse is the API response for a single order
type OrderResponse struct {
	Order Order `json:"order"`
}

// OrdersResponse is the API response for multiple orders
type OrdersResponse struct {
	Orders []Order `json:"orders"`
	Cursor string  `json:"cursor"`
}

// CreateOrderRequest is the request to create an order
type CreateOrderRequest struct {
	Ticker        string      `json:"ticker"`
	Side          OrderSide   `json:"side"`
	Action        OrderAction `json:"action"`
	Type          OrderType   `json:"type"`
	Count         int         `json:"count"`
	YesPrice      int         `json:"yes_price,omitempty"`
	NoPrice       int         `json:"no_price,omitempty"`
	ExpirationTs  int64       `json:"expiration_ts,omitempty"`
	ClientOrderID string      `json:"client_order_id,omitempty"`
	OrderGroupID  string      `json:"order_group_id,omitempty"`
	SubaccountID  int         `json:"subaccount_id,omitempty"`
	SellPositionFloor int     `json:"sell_position_floor,omitempty"`
	BuyMaxCost    int         `json:"buy_max_cost,omitempty"`
}

// CreateOrderResponse is the response from creating an order
type CreateOrderResponse struct {
	Order Order `json:"order"`
}

// AmendOrderRequest is the request to amend an order
type AmendOrderRequest struct {
	Price int `json:"price,omitempty"`
	Count int `json:"count,omitempty"`
}

// DecreaseOrderRequest is the request to decrease an order
type DecreaseOrderRequest struct {
	ReduceBy int `json:"reduce_by"`
}

// BatchCreateOrdersRequest is for batch order creation
type BatchCreateOrdersRequest struct {
	Orders []CreateOrderRequest `json:"orders"`
}

// BatchCreateOrdersResponse is the response from batch order creation
type BatchCreateOrdersResponse struct {
	Orders []Order `json:"orders"`
}

// BatchCancelOrdersRequest is for batch order cancellation
type BatchCancelOrdersRequest struct {
	OrderIDs []string `json:"order_ids,omitempty"`
	Ticker   string   `json:"ticker,omitempty"`
}

// BatchCancelOrdersResponse is the response from batch cancellation
type BatchCancelOrdersResponse struct {
	Orders []Order `json:"orders"`
}

// QueuePosition represents an order's queue position
type QueuePosition struct {
	OrderID       string `json:"order_id"`
	QueuePosition int    `json:"queue_position"`
}

// QueuePositionsResponse is the response for queue positions
type QueuePositionsResponse struct {
	Positions []QueuePosition `json:"queue_positions"`
}
