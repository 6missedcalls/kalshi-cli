package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// =============================================================================
// TDD Step 2: Implement to make tests pass (GREEN)
// =============================================================================

// ListMilestonesParams contains parameters for listing milestones
type ListMilestonesParams struct {
	MinDate time.Time
	MaxDate time.Time
	Limit   int
	Cursor  string
}

// toQueryParams converts ListMilestonesParams to query parameters
func (p ListMilestonesParams) toQueryParams() map[string]string {
	params := make(map[string]string)
	if !p.MinDate.IsZero() {
		params["min_date"] = p.MinDate.Format("2006-01-02")
	}
	if !p.MaxDate.IsZero() {
		params["max_date"] = p.MaxDate.Format("2006-01-02")
	}
	if p.Limit > 0 {
		params["limit"] = strconv.Itoa(p.Limit)
	}
	if p.Cursor != "" {
		params["cursor"] = p.Cursor
	}
	return params
}

// GetMilestone retrieves a single milestone by ID
func (c *Client) GetMilestone(ctx context.Context, milestoneID string) (*models.Milestone, error) {
	if milestoneID == "" {
		return nil, fmt.Errorf("milestone_id is required")
	}

	path := TradeAPIPrefix + "/milestones/" + milestoneID

	var result models.MilestoneResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get milestone: %w", err)
	}

	return &result.Milestone, nil
}

// ListMilestones retrieves milestones with optional date filtering
func (c *Client) ListMilestones(ctx context.Context, params ListMilestonesParams) (*models.MilestonesResponse, error) {
	path := TradeAPIPrefix + "/milestones" + BuildQueryString(params.toQueryParams())

	var result models.MilestonesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to list milestones: %w", err)
	}

	return &result, nil
}

// GetLiveData retrieves current data for a specific milestone
func (c *Client) GetLiveData(ctx context.Context, milestoneID string) (*models.LiveData, error) {
	if milestoneID == "" {
		return nil, fmt.Errorf("milestone_id is required")
	}

	path := TradeAPIPrefix + "/live-data/" + milestoneID

	var result models.LiveDataResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get live data: %w", err)
	}

	return &result.Data, nil
}

// GetBatchLiveData retrieves current data for multiple milestones
func (c *Client) GetBatchLiveData(ctx context.Context, milestoneIDs []string) (*models.BatchLiveDataResponse, error) {
	params := make(map[string]string)
	if len(milestoneIDs) > 0 {
		params["milestone_ids"] = strings.Join(milestoneIDs, ",")
	}

	path := TradeAPIPrefix + "/live-data" + BuildQueryString(params)

	var result models.BatchLiveDataResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get batch live data: %w", err)
	}

	return &result, nil
}
