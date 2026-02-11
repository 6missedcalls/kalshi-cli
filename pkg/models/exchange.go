package models

import "time"

// ExchangeStatusResponse is the API response for exchange status
type ExchangeStatusResponse struct {
	ExchangeActive bool `json:"exchange_active"`
	TradingActive  bool `json:"trading_active"`
}

// ExchangeSchedule represents the exchange schedule
type ExchangeSchedule struct {
	StandardHours      []WeeklySchedule    `json:"standard_hours"`
	MaintenanceWindows []MaintenanceWindow `json:"maintenance_windows"`
}

// WeeklySchedule represents a weekly schedule block
type WeeklySchedule struct {
	StartTime string          `json:"start_time"`
	EndTime   string          `json:"end_time"`
	Monday    []DailySchedule `json:"monday"`
	Tuesday   []DailySchedule `json:"tuesday"`
	Wednesday []DailySchedule `json:"wednesday"`
	Thursday  []DailySchedule `json:"thursday"`
	Friday    []DailySchedule `json:"friday"`
	Saturday  []DailySchedule `json:"saturday"`
	Sunday    []DailySchedule `json:"sunday"`
}

// DailySchedule represents open/close times for a day
type DailySchedule struct {
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

// MaintenanceWindow represents a scheduled maintenance window
type MaintenanceWindow struct {
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
}

// ExchangeScheduleResponse is the API response for schedule
type ExchangeScheduleResponse struct {
	Schedule ExchangeSchedule `json:"schedule"`
}

// Announcement represents an exchange announcement
type Announcement struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	CreatedTime  time.Time `json:"created_time"`
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
