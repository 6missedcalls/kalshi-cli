package models

import "time"

// RFQ represents a Request for Quote
type RFQ struct {
	RFQID       string    `json:"rfq_id"`
	Ticker      string    `json:"ticker"`
	Side        string    `json:"side"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	CreatedTime time.Time `json:"created_time"`
	ExpiresTime time.Time `json:"expires_time"`
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
	Ticker   string `json:"ticker"`
	Side     string `json:"side"`
	Quantity int    `json:"quantity"`
}

// Quote represents a quote on an RFQ
type Quote struct {
	QuoteID     string    `json:"quote_id"`
	RFQID       string    `json:"rfq_id"`
	Ticker      string    `json:"ticker"`
	Side        string    `json:"side"`
	Price       int       `json:"price"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	CreatedTime time.Time `json:"created_time"`
	ExpiresTime time.Time `json:"expires_time"`
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
	RFQID    string `json:"rfq_id"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity,omitempty"`
}

// CommunicationsID represents the communications ID
type CommunicationsID struct {
	ID string `json:"id"`
}

// CommunicationsIDResponse is the API response for communications ID
type CommunicationsIDResponse struct {
	CommunicationsID string `json:"communications_id"`
}
