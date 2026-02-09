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

// ListStructuredTargetsParams contains parameters for listing structured targets
type ListStructuredTargetsParams struct {
	Limit  int
	Cursor string
}

// toQueryParams converts ListStructuredTargetsParams to query parameters
func (p ListStructuredTargetsParams) toQueryParams() map[string]string {
	params := make(map[string]string)
	if p.Limit > 0 {
		params["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Cursor != "" {
		params["cursor"] = p.Cursor
	}
	return params
}

// GetSportsFilters retrieves sports filtering options
func (c *Client) GetSportsFilters(ctx context.Context) (*models.SportsFiltersResponse, error) {
	path := TradeAPIPrefix + "/search/sports/filters"

	var result models.SportsFiltersResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get sports filters: %w", err)
	}

	return &result, nil
}

// GetSearchTags retrieves series categories to tags mapping
func (c *Client) GetSearchTags(ctx context.Context) (*models.TagsResponse, error) {
	path := TradeAPIPrefix + "/search/tags"

	var result models.TagsResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get search tags: %w", err)
	}

	return &result, nil
}

// GetStructuredTarget retrieves a single structured target by ID
func (c *Client) GetStructuredTarget(ctx context.Context, targetID string) (*models.StructuredTarget, error) {
	if targetID == "" {
		return nil, fmt.Errorf("target_id is required")
	}

	path := TradeAPIPrefix + "/structured-targets/" + targetID

	var result models.StructuredTargetResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get structured target: %w", err)
	}

	return &result.Target, nil
}

// ListStructuredTargets retrieves structured targets with pagination
// Limit must be between 1 and 2000 per page
func (c *Client) ListStructuredTargets(ctx context.Context, params ListStructuredTargetsParams) (*models.StructuredTargetsResponse, error) {
	// Enforce limit bounds per API spec
	if params.Limit < 1 {
		params.Limit = 100 // Default
	}
	if params.Limit > 2000 {
		params.Limit = 2000
	}

	path := TradeAPIPrefix + "/structured-targets" + BuildQueryString(params.toQueryParams())

	var result models.StructuredTargetsResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to list structured targets: %w", err)
	}

	return &result, nil
}

// GetIncentives retrieves available rewards programs
func (c *Client) GetIncentives(ctx context.Context) (*models.IncentivesResponse, error) {
	path := TradeAPIPrefix + "/incentives"

	var result models.IncentivesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get incentives: %w", err)
	}

	return &result, nil
}
