package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/go-resty/resty/v2"
)

const (
	DefaultBaseURL    = "https://api.elections.kalshi.com"
	TradeAPIPrefix    = "/trade-api/v2"
	defaultTimeout    = 30 * time.Second
	maxRetries        = 5
	baseRetryDelay    = 100 * time.Millisecond
	maxRetryDelay     = 10 * time.Second
	retryMultiplier   = 2.0
	headerTimestamp   = "KALSHI-ACCESS-TIMESTAMP"
	headerAuth        = "Authorization"
)

// Client handles HTTP requests to the Kalshi API
type Client struct {
	resty   *resty.Client
	signer  *Signer
	baseURL string
	timeout time.Duration
}

// ClientOption is a functional option for configuring the client (legacy support)
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL (legacy support)
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
		c.resty.SetBaseURL(baseURL)
	}
}

// NewClientLegacy creates a new API client using the legacy functional options pattern
// Deprecated: Use NewClient(cfg, signer) instead
func NewClientLegacy(signer *Signer, opts ...ClientOption) *Client {
	client := &Client{
		resty:   resty.New(),
		signer:  signer,
		baseURL: DefaultBaseURL,
		timeout: defaultTimeout,
	}

	client.resty.SetBaseURL(DefaultBaseURL)
	client.resty.SetTimeout(defaultTimeout)
	client.resty.SetHeader("Content-Type", "application/json")
	client.resty.SetHeader("Accept", "application/json")

	// Add request signing middleware
	client.resty.OnBeforeRequest(client.signRequest)

	// Add retry configuration for rate limiting with exponential backoff
	client.resty.SetRetryCount(maxRetries)
	client.resty.SetRetryWaitTime(baseRetryDelay)
	client.resty.SetRetryMaxWaitTime(maxRetryDelay)
	client.resty.AddRetryCondition(func(resp *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		return IsRateLimitError(resp.StatusCode()) || IsServerError(resp.StatusCode())
	})
	client.resty.SetRetryAfter(func(c *resty.Client, resp *resty.Response) (time.Duration, error) {
		return calculateBackoff(resp)
	})

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// NewClient creates a new API client
func NewClient(cfg *config.Config, signer *Signer) *Client {
	baseURL := config.DemoBaseURL
	timeout := defaultTimeout

	if cfg != nil {
		baseURL = cfg.BaseURL()
		if cfg.API.Timeout > 0 {
			timeout = cfg.API.Timeout
		}
	}

	client := &Client{
		resty:   resty.New(),
		signer:  signer,
		baseURL: baseURL,
		timeout: timeout,
	}

	client.resty.SetBaseURL(baseURL)
	client.resty.SetTimeout(timeout)
	client.resty.SetHeader("Content-Type", "application/json")
	client.resty.SetHeader("Accept", "application/json")

	// Add request signing middleware
	client.resty.OnBeforeRequest(client.signRequest)

	// Add retry configuration for rate limiting with exponential backoff
	client.resty.SetRetryCount(maxRetries)
	client.resty.SetRetryWaitTime(baseRetryDelay)
	client.resty.SetRetryMaxWaitTime(maxRetryDelay)
	client.resty.AddRetryCondition(func(resp *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		return IsRateLimitError(resp.StatusCode()) || IsServerError(resp.StatusCode())
	})
	client.resty.SetRetryAfter(func(c *resty.Client, resp *resty.Response) (time.Duration, error) {
		return calculateBackoff(resp)
	})

	return client
}

// signRequest adds authentication headers to requests
func (c *Client) signRequest(client *resty.Client, req *resty.Request) error {
	if c.signer == nil {
		return nil
	}

	timestamp := time.Now().UTC()
	path := req.URL

	// Get the body as string for signing
	var body string
	if req.Body != nil {
		switch v := req.Body.(type) {
		case string:
			body = v
		case []byte:
			body = string(v)
		default:
			bodyBytes, err := json.Marshal(req.Body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			body = string(bodyBytes)
		}
	}

	signature, err := c.signer.Sign(timestamp, req.Method, path, body)
	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	req.SetHeader(headerTimestamp, TimestampHeader(timestamp))
	req.SetHeader(headerAuth, c.signer.AuthHeader(signature))

	return nil
}

// calculateBackoff determines the retry delay using exponential backoff
func calculateBackoff(resp *resty.Response) (time.Duration, error) {
	// Check for Retry-After header
	if retryAfter := resp.Header().Get("Retry-After"); retryAfter != "" {
		seconds, err := strconv.Atoi(retryAfter)
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second, nil
		}
	}

	// Use exponential backoff
	attempt := resp.Request.Attempt
	delay := float64(baseRetryDelay) * math.Pow(retryMultiplier, float64(attempt-1))
	if delay > float64(maxRetryDelay) {
		delay = float64(maxRetryDelay)
	}

	return time.Duration(delay), nil
}

// BaseURL returns the base URL of the API
func (c *Client) BaseURL() string {
	return c.baseURL
}

// SetBaseURL updates the base URL (useful for testing)
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
	c.resty.SetBaseURL(url)
}

// SetDebug enables or disables debug logging
func (c *Client) SetDebug(enabled bool) {
	c.resty.SetDebug(enabled)
}

// Get performs a GET request and returns the raw resty.Response
func (c *Client) Get(ctx context.Context, path string) (*resty.Response, error) {
	return c.resty.R().
		SetContext(ctx).
		Get(path)
}

// Post performs a POST request with a JSON body and returns the raw resty.Response
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*resty.Response, error) {
	return c.resty.R().
		SetContext(ctx).
		SetBody(body).
		Post(path)
}

// Put performs a PUT request with a JSON body and returns the raw resty.Response
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*resty.Response, error) {
	return c.resty.R().
		SetContext(ctx).
		SetBody(body).
		Put(path)
}

// DeleteRaw performs a DELETE request and returns the raw resty.Response
func (c *Client) DeleteRaw(ctx context.Context, path string) (*resty.Response, error) {
	return c.resty.R().
		SetContext(ctx).
		Delete(path)
}

// Legacy method signatures for backward compatibility with existing codebase

// GetJSON performs a GET request and unmarshals response into result
func (c *Client) GetJSON(ctx context.Context, path string, result interface{}) error {
	return c.DoRequest(ctx, http.MethodGet, path, nil, result)
}

// PostJSON performs a POST request and unmarshals response into result
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.DoRequest(ctx, http.MethodPost, path, body, result)
}

// PutJSON performs a PUT request and unmarshals response into result
func (c *Client) PutJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.DoRequest(ctx, http.MethodPut, path, body, result)
}

// DeleteJSON performs a DELETE request and unmarshals response into result
func (c *Client) DeleteJSON(ctx context.Context, path string, result interface{}) error {
	return c.DoRequest(ctx, http.MethodDelete, path, nil, result)
}

// DeleteWithBody performs a DELETE request with a request body
func (c *Client) DeleteWithBody(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.DoRequest(ctx, http.MethodDelete, path, body, result)
}

// APIError represents an error response from the Kalshi API
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error [%d] %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("API error [%d]: %s", e.StatusCode, e.Message)
}

// ParseAPIError extracts an APIError from a response
func ParseAPIError(resp *resty.Response) *APIError {
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}

	var apiErr APIError
	if err := json.Unmarshal(resp.Body(), &apiErr); err != nil {
		return &APIError{
			Code:       "UNKNOWN",
			Message:    string(resp.Body()),
			StatusCode: resp.StatusCode(),
		}
	}

	apiErr.StatusCode = resp.StatusCode()
	return &apiErr
}

// IsRateLimitError checks if the status code indicates rate limiting
func IsRateLimitError(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests
}

// IsServerError checks if the status code indicates a server error
func IsServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// BuildQueryString builds a query string from a map of parameters
func BuildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for k, v := range params {
		if v != "" {
			values.Set(k, v)
		}
	}

	if len(values) == 0 {
		return ""
	}

	return "?" + values.Encode()
}

// DoRequest performs an authenticated HTTP request (for backward compatibility)
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var resp *resty.Response
	var err error

	req := c.resty.R().SetContext(ctx)

	if body != nil {
		req.SetBody(body)
	}

	if result != nil {
		req.SetResult(result)
	}

	switch method {
	case http.MethodGet:
		resp, err = req.Get(path)
	case http.MethodPost:
		resp, err = req.Post(path)
	case http.MethodPut:
		resp, err = req.Put(path)
	case http.MethodPatch:
		resp, err = req.Patch(path)
	case http.MethodDelete:
		resp, err = req.Delete(path)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		apiErr := ParseAPIError(resp)
		if apiErr != nil {
			return apiErr
		}
	}

	return nil
}

// GetExchangeStatus returns the current exchange status
func (c *Client) GetExchangeStatus(ctx context.Context) (*ExchangeStatusResponse, error) {
	var result ExchangeStatusResponse
	if err := c.DoRequest(ctx, "GET", "/exchange/status", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExchangeStatusResponse is the response for exchange status
type ExchangeStatusResponse struct {
	ExchangeActive bool `json:"exchange_active"`
	TradingActive  bool `json:"trading_active"`
}

// ListAPIKeys returns all API keys for the authenticated user
func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	var result apiKeysResponse
	if err := c.DoRequest(ctx, "GET", "/api-keys", nil, &result); err != nil {
		return nil, err
	}
	return result.APIKeys, nil
}

// APIKey represents an API key
type APIKey struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	CreatedTime JSONTime `json:"created_time"`
	ExpiresTime JSONTime `json:"expires_time,omitempty"`
	Scopes      []string `json:"scopes"`
}

type apiKeysResponse struct {
	APIKeys []APIKey `json:"api_keys"`
}

// CreateAPIKeyRequest is the request to create an API key
type CreateAPIKeyRequest struct {
	Name string `json:"name,omitempty"`
}

// CreateAPIKeyResponse is the response from creating an API key
type CreateAPIKeyResponse struct {
	APIKey     APIKey `json:"api_key"`
	PrivateKey string `json:"private_key"`
}

// CreateAPIKey creates a new API key
func (c *Client) CreateAPIKey(ctx context.Context, req CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	var result CreateAPIKeyResponse
	if err := c.DoRequest(ctx, "POST", "/api-keys", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAPIKey deletes an API key by ID
func (c *Client) DeleteAPIKey(ctx context.Context, keyID string) error {
	return c.DoRequest(ctx, "DELETE", "/api-keys/"+keyID, nil, nil)
}

// JSONTime is a time.Time that can be unmarshaled from JSON
type JSONTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" || s == `""` {
		return nil
	}
	s = s[1 : len(s)-1]

	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return err
		}
	}
	t.Time = parsed
	return nil
}

// MarshalJSON implements json.Marshaler
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}

// Format formats the time
func (t JSONTime) Format(layout string) string {
	return t.Time.Format(layout)
}

// IsZero returns whether the time is zero
func (t JSONTime) IsZero() bool {
	return t.Time.IsZero()
}
