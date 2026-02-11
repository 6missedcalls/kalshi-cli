package models

import "time"

// Event represents a Kalshi event
type Event struct {
	EventTicker          string     `json:"event_ticker"`
	SeriesTicker         string     `json:"series_ticker"`
	Title                string     `json:"title"`
	SubTitle             string     `json:"sub_title"`
	Category             string     `json:"category"`
	MutuallyExclusive    bool       `json:"mutually_exclusive"`
	CollateralReturnType string     `json:"collateral_return_type"`
	StrikeDate           *time.Time `json:"strike_date,omitempty"`
	StrikePeriod         string     `json:"strike_period,omitempty"`
	AvailableOnBrokers   bool       `json:"available_on_brokers"`
	Markets              []string   `json:"markets,omitempty"`
}

// EventResponse is the API response for a single event
type EventResponse struct {
	Event Event `json:"event"`
}

// EventsResponse is the API response for multiple events
type EventsResponse struct {
	Events []Event `json:"events"`
	Cursor string  `json:"cursor"`
}

// EventMetadata contains additional event information
type EventMetadata struct {
	EventTicker string            `json:"event_ticker"`
	Metadata    map[string]string `json:"metadata"`
}

// EventMetadataResponse is the API response for event metadata
type EventMetadataResponse struct {
	EventMetadata EventMetadata `json:"event_metadata"`
}

// ForecastPercentilePoint represents a single point in forecast history
type ForecastPercentilePoint struct {
	Timestamp time.Time `json:"timestamp"`
	P10       int       `json:"p10"`
	P25       int       `json:"p25"`
	P50       int       `json:"p50"`
	P75       int       `json:"p75"`
	P90       int       `json:"p90"`
}

// ForecastPercentileHistoryResponse is the API response for forecast history
type ForecastPercentileHistoryResponse struct {
	History []ForecastPercentilePoint `json:"history"`
}

// MultivariateEvent represents a multivariate event
type MultivariateEvent struct {
	Ticker      string   `json:"ticker"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	LookupTable []string `json:"lookup_table"`
	LookupType  string   `json:"lookup_type"`
}

// MultivariateEventsResponse is the API response for multivariate events
type MultivariateEventsResponse struct {
	Events []MultivariateEvent `json:"multivariate_events"`
	Cursor string              `json:"cursor"`
}
