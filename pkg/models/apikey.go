package models

import "time"

// APIKey represents an API key
type APIKey struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedTime time.Time `json:"created_time"`
	ExpiresTime time.Time `json:"expires_time,omitempty"`
	Scopes      []string  `json:"scopes"`
}

// APIKeysResponse is the API response for API keys
type APIKeysResponse struct {
	APIKeys []APIKey `json:"api_keys"`
}

// CreateAPIKeyRequest is the request to create an API key
type CreateAPIKeyRequest struct {
	Name string `json:"name,omitempty"`
}

// CreateAPIKeyResponse is the response from creating an API key
type CreateAPIKeyResponse struct {
	APIKey     APIKey `json:"api_key"`
	PrivateKey string `json:"private_key"`
}

// CreateAPIKeyWithPublicKeyRequest creates API key with user's public key
type CreateAPIKeyWithPublicKeyRequest struct {
	Name      string `json:"name,omitempty"`
	PublicKey string `json:"public_key"`
}

// CreateAPIKeyWithPublicKeyResponse is the response
type CreateAPIKeyWithPublicKeyResponse struct {
	APIKey APIKey `json:"api_key"`
}

// APILimits represents account API limits
type APILimits struct {
	RateLimit        int `json:"rate_limit"`
	MaxOrdersPerCall int `json:"max_orders_per_call"`
}

// APILimitsResponse is the response for API limits
type APILimitsResponse struct {
	APILimits APILimits `json:"api_limits"`
}
