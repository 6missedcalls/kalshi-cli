package models

// SportsFilter represents a sports filtering option
type SportsFilter struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sport    string `json:"sport"`
	League   string `json:"league"`
	Category string `json:"category"`
}

// SportsFiltersResponse is the API response for sports filters
type SportsFiltersResponse struct {
	Filters []SportsFilter `json:"filters"`
}

// TagMapping represents a category to tags mapping
type TagMapping struct {
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

// TagsResponse is the API response for tags
type TagsResponse struct {
	Mappings []TagMapping `json:"tag_mappings"`
}

// StructuredTarget represents a structured target
type StructuredTarget struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// StructuredTargetResponse is the API response for a single target
type StructuredTargetResponse struct {
	Target StructuredTarget `json:"structured_target"`
}

// StructuredTargetsResponse is the API response for multiple targets
type StructuredTargetsResponse struct {
	Targets []StructuredTarget `json:"structured_targets"`
	Cursor  string             `json:"cursor,omitempty"`
}

// Incentive represents a rewards program
type Incentive struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
	Status      string  `json:"status"`
}

// IncentivesResponse is the API response for incentives
type IncentivesResponse struct {
	Incentives []Incentive `json:"incentives"`
}
