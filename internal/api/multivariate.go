package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// =============================================================================
// TDD Step 2: Implement to make tests pass (GREEN)
// =============================================================================

// ListMultivariateCollectionsParams contains parameters for listing multivariate collections
type ListMultivariateCollectionsParams struct {
	Status string
	Limit  int
	Cursor string
}

// toQueryParams converts ListMultivariateCollectionsParams to query parameters
func (p ListMultivariateCollectionsParams) toQueryParams() map[string]string {
	params := make(map[string]string)
	if p.Status != "" {
		params["status"] = p.Status
	}
	if p.Limit > 0 {
		params["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Cursor != "" {
		params["cursor"] = p.Cursor
	}
	return params
}

// LookupHistoryParams contains parameters for getting lookup history
type LookupHistoryParams struct {
	Limit  int
	Cursor string
}

// toQueryParams converts LookupHistoryParams to query parameters
func (p LookupHistoryParams) toQueryParams() map[string]string {
	params := make(map[string]string)
	if p.Limit > 0 {
		params["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Cursor != "" {
		params["cursor"] = p.Cursor
	}
	return params
}

// CreateCollectionMarketRequest is the request to create a market in a collection
type CreateCollectionMarketRequest struct {
	LookupValue string `json:"lookup_value"`
}

// ListMultivariateCollections retrieves all multivariate collections
func (c *Client) ListMultivariateCollections(ctx context.Context, params ListMultivariateCollectionsParams) (*models.MultivariateCollectionsResponse, error) {
	path := TradeAPIPrefix + "/multivariate-collections" + BuildQueryString(params.toQueryParams())

	var result models.MultivariateCollectionsResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to list multivariate collections: %w", err)
	}

	return &result, nil
}

// GetMultivariateCollection retrieves a single multivariate collection by ticker
func (c *Client) GetMultivariateCollection(ctx context.Context, ticker string) (*models.MultivariateCollection, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}

	path := TradeAPIPrefix + "/multivariate-collections/" + ticker

	var result models.MultivariateCollectionResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get multivariate collection: %w", err)
	}

	return &result.Collection, nil
}

// GetCollectionLookupHistory retrieves the lookup history for a multivariate collection
func (c *Client) GetCollectionLookupHistory(ctx context.Context, ticker string, params LookupHistoryParams) (*models.LookupHistoryResponse, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}

	path := TradeAPIPrefix + "/multivariate-collections/" + ticker + "/lookup-history" + BuildQueryString(params.toQueryParams())

	var result models.LookupHistoryResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get lookup history: %w", err)
	}

	return &result, nil
}

// CreateCollectionMarket creates a market in a multivariate collection
// This must be called before trading on a specific lookup value
func (c *Client) CreateCollectionMarket(ctx context.Context, ticker string, req CreateCollectionMarketRequest) (*models.CreateCollectionMarketResponse, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}
	if req.LookupValue == "" {
		return nil, fmt.Errorf("lookup_value is required")
	}

	path := TradeAPIPrefix + "/multivariate-collections/" + ticker + "/markets"

	var result models.CreateCollectionMarketResponse
	if err := c.DoRequest(ctx, "POST", path, req, &result); err != nil {
		return nil, fmt.Errorf("failed to create collection market: %w", err)
	}

	return &result, nil
}

// LookupCollectionMarket looks up a market ticker by lookup value
// Returns 404 if the market was never created
func (c *Client) LookupCollectionMarket(ctx context.Context, ticker string, lookupValue string) (*models.LookupCollectionMarketResponse, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is required")
	}
	if lookupValue == "" {
		return nil, fmt.Errorf("lookup_value is required")
	}

	params := map[string]string{
		"lookup_value": lookupValue,
	}
	path := TradeAPIPrefix + "/multivariate-collections/" + ticker + "/markets/lookup" + BuildQueryString(params)

	var result models.LookupCollectionMarketResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to lookup collection market: %w", err)
	}

	return &result, nil
}
