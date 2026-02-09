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

// FeeChange represents a fee change
type FeeChange struct {
	Ticker      string    `json:"ticker"`
	OldFee      int       `json:"old_fee"`
	NewFee      int       `json:"new_fee"`
	EffectiveAt time.Time `json:"effective_at"`
}

// FeeChangesResponse is the API response for fee changes
type FeeChangesResponse struct {
	FeeChanges []FeeChange `json:"fee_changes"`
}
