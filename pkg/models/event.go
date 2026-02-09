package models

import "time"

// Event represents a Kalshi event
type Event struct {
	EventTicker     string    `json:"event_ticker"`
	SeriesTicker    string    `json:"series_ticker"`
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle"`
	Category        string    `json:"category"`
	Status          string    `json:"status"`
	Markets         []string  `json:"markets"`
	MutuallyExclusive bool    `json:"mutually_exclusive"`
	StrikeDate      time.Time `json:"strike_date"`
	ExpectedExpiration time.Time `json:"expected_expiration"`
	CreatedTime     time.Time `json:"created_time"`
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

// MultivariateEvent represents a multivariate event
type MultivariateEvent struct {
	Ticker          string   `json:"ticker"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Status          string   `json:"status"`
	LookupTable     []string `json:"lookup_table"`
	LookupType      string   `json:"lookup_type"`
}

// MultivariateEventsResponse is the API response for multivariate events
type MultivariateEventsResponse struct {
	Events []MultivariateEvent `json:"multivariate_events"`
	Cursor string              `json:"cursor"`
}
