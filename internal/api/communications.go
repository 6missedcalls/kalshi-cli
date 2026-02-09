package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

const (
	// Per Kalshi API docs, RFQs and Quotes are under /communications
	rfqsBasePath           = TradeAPIPrefix + "/communications/rfqs"
	quotesBasePath         = TradeAPIPrefix + "/communications/quotes"
	communicationsBasePath = TradeAPIPrefix + "/communications"
)

// RFQsOptions contains options for listing RFQs
type RFQsOptions struct {
	Ticker string
	Status string
	Cursor string
	Limit  int
}

// toQueryParams converts RFQsOptions to query parameters
func (o RFQsOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
	if o.Ticker != "" {
		params["ticker"] = o.Ticker
	}
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

// QuotesOptions contains options for listing quotes
type QuotesOptions struct {
	RFQID  string
	Status string
	Cursor string
	Limit  int
}

// toQueryParams converts QuotesOptions to query parameters
func (o QuotesOptions) toQueryParams() map[string]string {
	params := make(map[string]string)
	if o.RFQID != "" {
		params["rfq_id"] = o.RFQID
	}
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

// GetRFQs returns a list of RFQs based on the provided options
func (c *Client) GetRFQs(ctx context.Context, opts RFQsOptions) (*models.RFQsResponse, error) {
	path := rfqsBasePath + BuildQueryString(opts.toQueryParams())

	var result models.RFQsResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRFQ returns a single RFQ by ID
func (c *Client) GetRFQ(ctx context.Context, rfqID string) (*models.RFQResponse, error) {
	if rfqID == "" {
		return nil, fmt.Errorf("RFQ ID is required")
	}

	path := rfqsBasePath + "/" + rfqID

	var result models.RFQResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRFQ creates a new RFQ
func (c *Client) CreateRFQ(ctx context.Context, req models.CreateRFQRequest) (*models.RFQResponse, error) {
	var result models.RFQResponse
	if err := c.PostJSON(ctx, rfqsBasePath, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelRFQ cancels an existing RFQ
func (c *Client) CancelRFQ(ctx context.Context, rfqID string) error {
	if rfqID == "" {
		return fmt.Errorf("RFQ ID is required")
	}

	path := rfqsBasePath + "/" + rfqID
	return c.DeleteJSON(ctx, path, nil)
}

// GetQuotes returns a list of quotes based on the provided options
func (c *Client) GetQuotes(ctx context.Context, opts QuotesOptions) (*models.QuotesResponse, error) {
	path := quotesBasePath + BuildQueryString(opts.toQueryParams())

	var result models.QuotesResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQuote returns a single quote by ID
func (c *Client) GetQuote(ctx context.Context, quoteID string) (*models.QuoteResponse, error) {
	if quoteID == "" {
		return nil, fmt.Errorf("quote ID is required")
	}

	path := quotesBasePath + "/" + quoteID

	var result models.QuoteResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateQuote creates a new quote on an RFQ
func (c *Client) CreateQuote(ctx context.Context, req models.CreateQuoteRequest) (*models.QuoteResponse, error) {
	var result models.QuoteResponse
	if err := c.PostJSON(ctx, quotesBasePath, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AcceptQuote accepts a quote
func (c *Client) AcceptQuote(ctx context.Context, quoteID string) (*models.QuoteResponse, error) {
	if quoteID == "" {
		return nil, fmt.Errorf("quote ID is required")
	}

	path := quotesBasePath + "/" + quoteID + "/accept"

	var result models.QuoteResponse
	if err := c.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelQuote cancels an existing quote
func (c *Client) CancelQuote(ctx context.Context, quoteID string) error {
	if quoteID == "" {
		return fmt.Errorf("quote ID is required")
	}

	path := quotesBasePath + "/" + quoteID
	return c.DeleteJSON(ctx, path, nil)
}

// ConfirmQuote confirms a quote (quoter confirms their own quote after RFQ creator accepts)
func (c *Client) ConfirmQuote(ctx context.Context, quoteID string) (*models.QuoteResponse, error) {
	if quoteID == "" {
		return nil, fmt.Errorf("quote ID is required")
	}

	path := quotesBasePath + "/" + quoteID + "/confirm"

	var result models.QuoteResponse
	if err := c.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCommunicationsID returns the user's communications ID for websocket subscriptions
func (c *Client) GetCommunicationsID(ctx context.Context) (*models.CommunicationsIDResponse, error) {
	path := communicationsBasePath + "/id"

	var result models.CommunicationsIDResponse
	if err := c.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
