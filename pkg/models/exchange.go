package models

import "time"

// ExchangeStatus represents the exchange status
type ExchangeStatus struct {
	TradingActive  bool   `json:"trading_active"`
	ExchangeActive bool   `json:"exchange_active"`
}

// ExchangeStatusResponse is the API response for exchange status
type ExchangeStatusResponse struct {
	ExchangeActive bool `json:"exchange_active"`
	TradingActive  bool `json:"trading_active"`
}

// ExchangeSchedule represents the exchange schedule
type ExchangeSchedule struct {
	ScheduleEntries []ScheduleEntry `json:"schedule_entries"`
}

// ScheduleEntry represents a schedule entry
type ScheduleEntry struct {
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Maintenance   bool      `json:"maintenance"`
}

// ExchangeScheduleResponse is the API response for schedule
type ExchangeScheduleResponse struct {
	Schedule ExchangeSchedule `json:"schedule"`
}

// Announcement represents an exchange announcement
type Announcement struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	CreatedTime time.Time `json:"created_time"`
	DeliveryTime time.Time `json:"delivery_time"`
}

// AnnouncementsResponse is the API response for announcements
type AnnouncementsResponse struct {
	Announcements []Announcement `json:"announcements"`
}

// FeeChange represents a fee change (deprecated - use SeriesFeeChange)
type FeeChange struct {
	Ticker      string    `json:"ticker"`
	OldFee      int       `json:"old_fee"`
	NewFee      int       `json:"new_fee"`
	EffectiveAt time.Time `json:"effective_at"`
}

// FeeChangesResponse is the API response for fee changes (deprecated - use SeriesFeeChangesResponse)
type FeeChangesResponse struct {
	FeeChanges []FeeChange `json:"fee_changes"`
}

// SeriesFeeChange represents a fee change for a series per Kalshi API spec
type SeriesFeeChange struct {
	SeriesTicker  string    `json:"series_ticker"`
	OldFeeRate    float64   `json:"old_fee_rate"`
	NewFeeRate    float64   `json:"new_fee_rate"`
	EffectiveDate time.Time `json:"effective_date"`
	AnnouncedDate time.Time `json:"announced_date"`
}

// SeriesFeeChangesResponse is the API response for series fee changes
type SeriesFeeChangesResponse struct {
	SeriesFeeChanges []SeriesFeeChange `json:"series_fee_changes"`
}

// UserDataTimestampResponse is the API response for user data timestamp
type UserDataTimestampResponse struct {
	Timestamp time.Time `json:"timestamp"`
}
