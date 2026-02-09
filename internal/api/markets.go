package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// ListMarketsParams contains parameters for listing markets
type ListMarketsParams struct {
	Status       string
	SeriesTicker string
	EventTicker  string
	Tickers      []string
	MinCloseTs   int64
	MaxCloseTs   int64
	Limit        int
	Cursor       string
}

// ListMarkets retrieves a list of markets
func (c *Client) ListMarkets(ctx context.Context, params ListMarketsParams) (*models.MarketsResponse, error) {
	queryParams := map[string]string{}

	if params.Status != "" {
		queryParams["status"] = params.Status
	}
	if params.SeriesTicker != "" {
		queryParams["series_ticker"] = params.SeriesTicker
	}
	if params.EventTicker != "" {
		queryParams["event_ticker"] = params.EventTicker
	}
	if len(params.Tickers) > 0 {
		queryParams["tickers"] = strings.Join(params.Tickers, ",")
	}
	if params.MinCloseTs > 0 {
		queryParams["min_close_ts"] = strconv.FormatInt(params.MinCloseTs, 10)
	}
	if params.MaxCloseTs > 0 {
		queryParams["max_close_ts"] = strconv.FormatInt(params.MaxCloseTs, 10)
	}
	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Cursor != "" {
		queryParams["cursor"] = params.Cursor
	}

	path := TradeAPIPrefix + "/markets" + BuildQueryString(queryParams)

	var result models.MarketsResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to list markets: %w", err)
	}

	return &result, nil
}

// GetMarket retrieves a single market by ticker
func (c *Client) GetMarket(ctx context.Context, ticker string) (*models.Market, error) {
	path := TradeAPIPrefix + "/markets/" + ticker

	var result models.MarketResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get market: %w", err)
	}

	return &result.Market, nil
}

// GetOrderbook retrieves the orderbook for a market
func (c *Client) GetOrderbook(ctx context.Context, ticker string) (*models.Orderbook, error) {
	path := TradeAPIPrefix + "/markets/" + ticker + "/orderbook"

	var result models.OrderbookResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get orderbook: %w", err)
	}

	return &result.Orderbook, nil
}

// GetOrderbookWithDepth retrieves the orderbook for a market with a specific depth
func (c *Client) GetOrderbookWithDepth(ctx context.Context, ticker string, depth int) (*models.Orderbook, error) {
	queryParams := map[string]string{}
	if depth > 0 {
		queryParams["depth"] = strconv.Itoa(depth)
	}

	path := TradeAPIPrefix + "/markets/" + ticker + "/orderbook" + BuildQueryString(queryParams)

	var result models.OrderbookResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get orderbook: %w", err)
	}

	return &result.Orderbook, nil
}

// GetTradesParams contains parameters for getting trades
type GetTradesParams struct {
	Ticker string
	Limit  int
	Cursor string
	MinTs  int64
	MaxTs  int64
}

// GetTrades retrieves trades for a market
func (c *Client) GetTrades(ctx context.Context, params GetTradesParams) (*models.TradesResponse, error) {
	queryParams := map[string]string{}

	if params.Ticker != "" {
		queryParams["ticker"] = params.Ticker
	}
	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Cursor != "" {
		queryParams["cursor"] = params.Cursor
	}
	if params.MinTs > 0 {
		queryParams["min_ts"] = strconv.FormatInt(params.MinTs, 10)
	}
	if params.MaxTs > 0 {
		queryParams["max_ts"] = strconv.FormatInt(params.MaxTs, 10)
	}

	path := TradeAPIPrefix + "/markets/trades" + BuildQueryString(queryParams)

	var result models.TradesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}

	return &result, nil
}

// GetCandlesticksParams contains parameters for getting candlesticks
type GetCandlesticksParams struct {
	Ticker     string
	Period     string
	StartTime  int64
	EndTime    int64
}

// GetCandlesticks retrieves candlestick data for a market
func (c *Client) GetCandlesticks(ctx context.Context, params GetCandlesticksParams) (*models.CandlesticksResponse, error) {
	queryParams := map[string]string{}

	if params.Period != "" {
		queryParams["period_interval"] = periodToInterval(params.Period)
	}
	if params.StartTime > 0 {
		queryParams["start_ts"] = strconv.FormatInt(params.StartTime, 10)
	}
	if params.EndTime > 0 {
		queryParams["end_ts"] = strconv.FormatInt(params.EndTime, 10)
	}

	path := TradeAPIPrefix + "/markets/" + params.Ticker + "/candlesticks" + BuildQueryString(queryParams)

	var result models.CandlesticksResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get candlesticks: %w", err)
	}

	return &result, nil
}

// periodToInterval converts user-friendly period strings to API intervals
func periodToInterval(period string) string {
	switch period {
	case "1m":
		return "1"
	case "5m":
		return "5"
	case "15m":
		return "15"
	case "1h":
		return "60"
	case "4h":
		return "240"
	case "1d":
		return "1440"
	default:
		return period
	}
}

// ListSeriesParams contains parameters for listing series
type ListSeriesParams struct {
	Category string
	Limit    int
	Cursor   string
}

// ListSeries retrieves a list of series
func (c *Client) ListSeries(ctx context.Context, params ListSeriesParams) (*models.SeriesResponse, error) {
	queryParams := map[string]string{}

	if params.Category != "" {
		queryParams["category"] = params.Category
	}
	if params.Limit > 0 {
		queryParams["limit"] = strconv.Itoa(params.Limit)
	}
	if params.Cursor != "" {
		queryParams["cursor"] = params.Cursor
	}

	path := TradeAPIPrefix + "/series" + BuildQueryString(queryParams)

	var result models.SeriesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to list series: %w", err)
	}

	return &result, nil
}

// GetSeriesResponse is the API response for a single series
type GetSeriesResponse struct {
	Series models.Series `json:"series"`
}

// GetSeries retrieves a single series by ticker
func (c *Client) GetSeries(ctx context.Context, ticker string) (*models.Series, error) {
	path := TradeAPIPrefix + "/series/" + ticker

	var result GetSeriesResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	return &result.Series, nil
}

// GetBatchCandlesticksParams contains parameters for getting batch candlesticks
type GetBatchCandlesticksParams struct {
	Tickers   []string
	Period    string
	StartTime int64
	EndTime   int64
}

// MaxBatchCandlesticksTickers is the maximum number of tickers allowed in batch request
const MaxBatchCandlesticksTickers = 100

// GetBatchCandlesticks retrieves candlestick data for multiple markets (up to 100 tickers)
func (c *Client) GetBatchCandlesticks(ctx context.Context, params GetBatchCandlesticksParams) (*models.BatchCandlesticksResponse, error) {
	if len(params.Tickers) > MaxBatchCandlesticksTickers {
		return nil, fmt.Errorf("batch candlesticks request exceeds maximum of %d tickers", MaxBatchCandlesticksTickers)
	}

	queryParams := map[string]string{}

	if len(params.Tickers) > 0 {
		queryParams["tickers"] = strings.Join(params.Tickers, ",")
	}
	if params.Period != "" {
		queryParams["period_interval"] = periodToInterval(params.Period)
	}
	if params.StartTime > 0 {
		queryParams["start_ts"] = strconv.FormatInt(params.StartTime, 10)
	}
	if params.EndTime > 0 {
		queryParams["end_ts"] = strconv.FormatInt(params.EndTime, 10)
	}

	path := TradeAPIPrefix + "/markets/candlesticks/batch" + BuildQueryString(queryParams)

	var result models.BatchCandlesticksResponse
	if err := c.DoRequest(ctx, "GET", path, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get batch candlesticks: %w", err)
	}

	return &result, nil
}
