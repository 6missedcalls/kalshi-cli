package models

import "time"

// Milestone represents a milestone in the Kalshi system
type Milestone struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	Category       string    `json:"category"`
	TargetDate     time.Time `json:"target_date"`
	ResolutionDate time.Time `json:"resolution_date,omitempty"`
	CreatedTime    time.Time `json:"created_time"`
}

// MilestoneResponse is the API response for a single milestone
type MilestoneResponse struct {
	Milestone Milestone `json:"milestone"`
}

// MilestonesResponse is the API response for multiple milestones
type MilestonesResponse struct {
	Milestones []Milestone `json:"milestones"`
	Cursor     string      `json:"cursor,omitempty"`
}

// LiveData represents current data for a milestone
type LiveData struct {
	MilestoneID string    `json:"milestone_id"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp"`
}

// LiveDataResponse is the API response for live data
type LiveDataResponse struct {
	Data LiveData `json:"live_data"`
}

// BatchLiveDataResponse is the API response for multiple live data items
type BatchLiveDataResponse struct {
	Data []LiveData `json:"live_data"`
}
