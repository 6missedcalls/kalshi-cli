package models

import "time"

// MultivariateCollection represents a multivariate collection
type MultivariateCollection struct {
	Ticker       string   `json:"ticker"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Status       string   `json:"status"`
	LookupType   string   `json:"lookup_type"`
	LookupTable  []string `json:"lookup_table,omitempty"`
	CategoryPath []string `json:"category_path,omitempty"`
}

// MultivariateCollectionsResponse is the API response for multiple collections
type MultivariateCollectionsResponse struct {
	Collections []MultivariateCollection `json:"multivariate_collections"`
	Cursor      string                   `json:"cursor,omitempty"`
}

// MultivariateCollectionResponse is the API response for a single collection
type MultivariateCollectionResponse struct {
	Collection MultivariateCollection `json:"multivariate_collection"`
}

// LookupHistoryEntry represents a single lookup history entry
type LookupHistoryEntry struct {
	Ticker      string    `json:"ticker"`
	LookupValue string    `json:"lookup_value"`
	CreatedTime time.Time `json:"created_time"`
}

// LookupHistoryResponse is the API response for lookup history
type LookupHistoryResponse struct {
	History []LookupHistoryEntry `json:"history"`
	Cursor  string               `json:"cursor,omitempty"`
}

// CreateCollectionMarketResponse is the API response for creating a market
type CreateCollectionMarketResponse struct {
	MarketTicker string `json:"market_ticker"`
	Created      bool   `json:"created"`
}

// LookupCollectionMarketResponse is the API response for looking up a market
type LookupCollectionMarketResponse struct {
	MarketTicker string `json:"market_ticker"`
	Found        bool   `json:"found"`
}
