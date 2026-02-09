package api

import (
	"context"
	"strconv"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

const portfolioBasePath = TradeAPIPrefix + "/portfolio"

// PositionsOptions contains options for listing positions
type PositionsOptions struct {
	Ticker           string
	EventTicker      string
	Cursor           string
	Limit            int
	SettlementStatus string
	CountFilter      string
	SubaccountID     int
}

// toQueryParams converts PositionsOptions to query parameters
func (o PositionsOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
	if o.Ticker != "" {
		params["ticker"] = o.Ticker
	}
	if o.EventTicker != "" {
		params["event_ticker"] = o.EventTicker
	}
	if o.Cursor != "" {
		params["cursor"] = o.Cursor
	}
	if o.Limit > 0 {
		params["limit"] = strconv.Itoa(o.Limit)
	}
	if o.SettlementStatus != "" {
		params["settlement_status"] = o.SettlementStatus
	}
	if o.CountFilter != "" {
		params["count_filter"] = o.CountFilter
	}
	if o.SubaccountID > 0 {
		params["subaccount_id"] = strconv.Itoa(o.SubaccountID)
	}
	return params
}

// FillsOptions contains options for listing fills
type FillsOptions struct {
	Ticker       string
	OrderID      string
	Cursor       string
	Limit        int
	MinTS        int64
	MaxTS        int64
	SubaccountID int
}

// toQueryParams converts FillsOptions to query parameters
func (o FillsOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
	if o.Ticker != "" {
		params["ticker"] = o.Ticker
	}
	if o.OrderID != "" {
		params["order_id"] = o.OrderID
	}
	if o.Cursor != "" {
		params["cursor"] = o.Cursor
	}
	if o.Limit > 0 {
		params["limit"] = strconv.Itoa(o.Limit)
	}
	if o.MinTS > 0 {
		params["min_ts"] = strconv.FormatInt(o.MinTS, 10)
	}
	if o.MaxTS > 0 {
		params["max_ts"] = strconv.FormatInt(o.MaxTS, 10)
	}
	if o.SubaccountID > 0 {
		params["subaccount_id"] = strconv.Itoa(o.SubaccountID)
	}
	return params
}

// SettlementsOptions contains options for listing settlements
type SettlementsOptions struct {
	Cursor       string
	Limit        int
	SubaccountID int
}

// toQueryParams converts SettlementsOptions to query parameters
func (o SettlementsOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
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

// GetBalance returns the account balance
func (c *Client) GetBalance(ctx context.Context) (*models.BalanceResponse, error) {
	path := portfolioBasePath + "/balance"

	var result models.BalanceResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPositions returns market positions based on the provided options
func (c *Client) GetPositions(ctx context.Context, opts PositionsOptions) (*models.PositionsResponse, error) {
	path := portfolioBasePath + "/positions" + BuildQueryString(opts.toQueryParams())

	var result models.PositionsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFills returns trade fills based on the provided options
func (c *Client) GetFills(ctx context.Context, opts FillsOptions) (*models.FillsResponse, error) {
	path := portfolioBasePath + "/fills" + BuildQueryString(opts.toQueryParams())

	var result models.FillsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSettlements returns settlements based on the provided options
func (c *Client) GetSettlements(ctx context.Context, opts SettlementsOptions) (*models.SettlementsResponse, error) {
	path := portfolioBasePath + "/settlements" + BuildQueryString(opts.toQueryParams())

	var result models.SettlementsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSubaccounts returns all subaccounts
func (c *Client) GetSubaccounts(ctx context.Context) (*models.SubaccountsResponse, error) {
	path := portfolioBasePath + "/subaccounts"

	var result models.SubaccountsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateSubaccount creates a new subaccount
func (c *Client) CreateSubaccount(ctx context.Context) (*models.Subaccount, error) {
	path := portfolioBasePath + "/subaccounts"

	var result models.Subaccount
	if err := c.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTransfers returns all subaccount transfers
func (c *Client) GetTransfers(ctx context.Context) (*models.TransfersResponse, error) {
	path := portfolioBasePath + "/subaccounts/transfers"

	var result models.TransfersResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Transfer creates a transfer between subaccounts
func (c *Client) Transfer(ctx context.Context, req models.TransferRequest) (*models.Transfer, error) {
	path := portfolioBasePath + "/subaccounts/transfers"

	var result models.Transfer
	if err := c.PostJSON(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSubaccountBalances returns all subaccount balances
func (c *Client) GetSubaccountBalances(ctx context.Context) (*models.SubaccountBalancesResponse, error) {
	path := portfolioBasePath + "/subaccounts/balances"

	var result models.SubaccountBalancesResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRestingOrderValue returns the total resting order value (FCM only)
func (c *Client) GetRestingOrderValue(ctx context.Context) (*models.RestingOrderValueResponse, error) {
	path := portfolioBasePath + "/resting-order-value"

	var result models.RestingOrderValueResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
