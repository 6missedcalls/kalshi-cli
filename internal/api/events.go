package api

import (
	"context"
	"strconv"
	"time"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// ListEventsParams contains parameters for listing events
type ListEventsParams struct {
	Status string
	Limit  int
	Cursor string
}

// CandlesticksParams contains parameters for getting candlesticks
type CandlesticksParams struct {
	SeriesTicker string
	Ticker       string
	Period       string
	StartTime    *time.Time
	EndTime      *time.Time
}

// ListMultivariateParams contains parameters for listing multivariate events
type ListMultivariateParams struct {
	Status string
	Limit  int
	Cursor string
}

// ForecastPercentileHistoryParams contains parameters for getting forecast history
type ForecastPercentileHistoryParams struct {
	Ticker    string
	StartTime *time.Time
	EndTime   *time.Time
}

// ListEvents retrieves a list of events with optional filtering
func (c *Client) ListEvents(ctx context.Context, params ListEventsParams) ([]models.Event, string, error) {
	queryParams := make(map[string]string)

	if params.Status != "" {
		queryParams["status"] = params.Status
	}
	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Cursor != "" {
		queryParams["cursor"] = params.Cursor
	}

	path := TradeAPIPrefix + "/events" + BuildQueryString(queryParams)

	var resp models.EventsResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, "", err
	}

	return resp.Events, resp.Cursor, nil
}

// GetEvent retrieves a single event by ticker
func (c *Client) GetEvent(ctx context.Context, ticker string) (*models.Event, error) {
	path := TradeAPIPrefix + "/events/" + ticker

	var resp models.EventResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Event, nil
}

// GetEventCandlesticks retrieves candlestick data for an event
func (c *Client) GetEventCandlesticks(ctx context.Context, params CandlesticksParams) ([]models.Candlestick, error) {
	queryParams := make(map[string]string)

	if params.Period != "" {
		queryParams["period_interval"] = periodToInterval(params.Period)
	}
	if params.StartTime != nil {
		queryParams["start_ts"] = strconv.FormatInt(params.StartTime.Unix(), 10)
	}
	if params.EndTime != nil {
		queryParams["end_ts"] = strconv.FormatInt(params.EndTime.Unix(), 10)
	}

	path := TradeAPIPrefix + "/series/" + params.SeriesTicker + "/events/" + params.Ticker + "/candlesticks" + BuildQueryString(queryParams)

	var resp models.CandlesticksResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Candlesticks, nil
}

// ListMultivariateEvents retrieves a list of multivariate events
func (c *Client) ListMultivariateEvents(ctx context.Context, params ListMultivariateParams) ([]models.MultivariateEvent, string, error) {
	queryParams := make(map[string]string)

	if params.Status != "" {
		queryParams["status"] = params.Status
	}
	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Cursor != "" {
		queryParams["cursor"] = params.Cursor
	}

	path := TradeAPIPrefix + "/events/multivariate" + BuildQueryString(queryParams)

	var resp models.MultivariateEventsResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, "", err
	}

	return resp.Events, resp.Cursor, nil
}

// MultivariateEventResponse is the API response for a single multivariate event
type MultivariateEventResponse struct {
	Event models.MultivariateEvent `json:"multivariate_event"`
}

// GetMultivariateEvent retrieves a single multivariate event by ticker
func (c *Client) GetMultivariateEvent(ctx context.Context, ticker string) (*models.MultivariateEvent, error) {
	path := TradeAPIPrefix + "/events/multivariate/" + ticker

	var resp MultivariateEventResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Event, nil
}

// GetEventMetadata retrieves metadata for a specific event
func (c *Client) GetEventMetadata(ctx context.Context, ticker string) (*models.EventMetadata, error) {
	path := TradeAPIPrefix + "/events/" + ticker + "/metadata"

	var resp models.EventMetadataResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.EventMetadata, nil
}

// GetForecastPercentileHistory retrieves historical forecast percentile data for an event
func (c *Client) GetForecastPercentileHistory(ctx context.Context, params ForecastPercentileHistoryParams) ([]models.ForecastPercentilePoint, error) {
	queryParams := make(map[string]string)

	if params.StartTime != nil {
		queryParams["start_ts"] = strconv.FormatInt(params.StartTime.Unix(), 10)
	}
	if params.EndTime != nil {
		queryParams["end_ts"] = strconv.FormatInt(params.EndTime.Unix(), 10)
	}

	path := TradeAPIPrefix + "/events/" + params.Ticker + "/forecast-percentile-history" + BuildQueryString(queryParams)

	var resp models.ForecastPercentileHistoryResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}

	return resp.History, nil
}
