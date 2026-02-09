package api

import (
	"context"
	"fmt"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// GetExchangeSchedule retrieves the exchange schedule
func (c *Client) GetExchangeSchedule(ctx context.Context, result *models.ExchangeScheduleResponse) error {
	path := TradeAPIPrefix + "/exchange/schedule"

	if err := c.DoRequest(ctx, "GET", path, nil, result); err != nil {
		return fmt.Errorf("failed to get exchange schedule: %w", err)
	}

	return nil
}

// GetAnnouncements retrieves the exchange announcements
func (c *Client) GetAnnouncements(ctx context.Context, result *models.AnnouncementsResponse) error {
	path := TradeAPIPrefix + "/exchange/announcements"

	if err := c.DoRequest(ctx, "GET", path, nil, result); err != nil {
		return fmt.Errorf("failed to get announcements: %w", err)
	}

	return nil
}

// GetFeeChanges retrieves the upcoming fee changes
// Deprecated: Use GetSeriesFeeChanges instead - this uses the wrong API path
func (c *Client) GetFeeChanges(ctx context.Context) (*models.FeeChangesResponse, error) {
	path := TradeAPIPrefix + "/exchange/fee-changes"

	var result models.FeeChangesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get fee changes: %w", err)
	}

	return &result, nil
}

// GetSeriesFeeChanges retrieves the upcoming series fee changes per Kalshi API spec
func (c *Client) GetSeriesFeeChanges(ctx context.Context) (*models.SeriesFeeChangesResponse, error) {
	path := TradeAPIPrefix + "/exchange/series-fee-changes"

	var result models.SeriesFeeChangesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get series fee changes: %w", err)
	}

	return &result, nil
}

// GetUserDataTimestamp retrieves the timestamp indicating when user data was last updated
func (c *Client) GetUserDataTimestamp(ctx context.Context) (*models.UserDataTimestampResponse, error) {
	path := TradeAPIPrefix + "/exchange/user-data-timestamp"

	var result models.UserDataTimestampResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get user data timestamp: %w", err)
	}

	return &result, nil
}
