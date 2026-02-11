package models

// RFQ represents a Request for Quote
type RFQ struct {
	ID                  string `json:"id"`
	CreatorID           string `json:"creator_id"`
	MarketTicker        string `json:"market_ticker"`
	Contracts           int    `json:"contracts"`
	ContractsFP         string `json:"contracts_fp"`
	TargetCostDollars   string `json:"target_cost_dollars"`
	Status              string `json:"status"`
	CreatedTs           string `json:"created_ts"`
	UpdatedTs           string `json:"updated_ts"`
	MveCollectionTicker string `json:"mve_collection_ticker,omitempty"`
	RestRemainder       bool   `json:"rest_remainder"`
	CancellationReason  string `json:"cancellation_reason,omitempty"`
	CreatorUserID       string `json:"creator_user_id"`
	CancelledTs         string `json:"cancelled_ts,omitempty"`
}

// RFQsResponse is the API response for RFQs
type RFQsResponse struct {
	RFQs   []RFQ  `json:"rfqs"`
	Cursor string `json:"cursor"`
}

// RFQResponse is the API response for a single RFQ
type RFQResponse struct {
	RFQ RFQ `json:"rfq"`
}

// CreateRFQRequest is the request to create an RFQ
type CreateRFQRequest struct {
	MarketTicker string `json:"market_ticker"`
	Contracts    int    `json:"contracts"`
}

// Quote represents a quote on an RFQ
type Quote struct {
	ID                 string `json:"id"`
	RFQID              string `json:"rfq_id"`
	CreatorID          string `json:"creator_id"`
	RFQCreatorID       string `json:"rfq_creator_id"`
	MarketTicker       string `json:"market_ticker"`
	Contracts          int    `json:"contracts"`
	ContractsFP        string `json:"contracts_fp"`
	YesBid             int    `json:"yes_bid"`
	NoBid              int    `json:"no_bid"`
	YesBidDollars      string `json:"yes_bid_dollars"`
	NoBidDollars       string `json:"no_bid_dollars"`
	CreatedTs          string `json:"created_ts"`
	UpdatedTs          string `json:"updated_ts"`
	Status             string `json:"status"`
	AcceptedSide       string `json:"accepted_side,omitempty"`
	AcceptedTs         string `json:"accepted_ts,omitempty"`
	ConfirmedTs        string `json:"confirmed_ts,omitempty"`
	ExecutedTs         string `json:"executed_ts,omitempty"`
	CancelledTs        string `json:"cancelled_ts,omitempty"`
	RestRemainder      bool   `json:"rest_remainder"`
	CancellationReason string `json:"cancellation_reason,omitempty"`
	CreatorUserID      string `json:"creator_user_id"`
}

// QuotesResponse is the API response for quotes
type QuotesResponse struct {
	Quotes []Quote `json:"quotes"`
	Cursor string  `json:"cursor"`
}

// QuoteResponse is the API response for a single quote
type QuoteResponse struct {
	Quote Quote `json:"quote"`
}

// CreateQuoteRequest is the request to create a quote
type CreateQuoteRequest struct {
	RFQID  string `json:"rfq_id"`
	YesBid int    `json:"yes_bid"`
	NoBid  int    `json:"no_bid,omitempty"`
}

// CommunicationsID represents the communications ID
type CommunicationsID struct {
	ID string `json:"id"`
}

// CommunicationsIDResponse is the API response for communications ID
type CommunicationsIDResponse struct {
	CommunicationsID string `json:"communications_id"`
}
