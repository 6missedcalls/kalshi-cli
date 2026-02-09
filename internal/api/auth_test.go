package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"
)

func TestNewSigner(t *testing.T) {
	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("test-api-key-id", privateKey)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	if signer == nil {
		t.Fatal("NewSigner returned nil signer")
	}

	if signer.APIKeyID() != "test-api-key-id" {
		t.Errorf("expected API key ID 'test-api-key-id', got '%s'", signer.APIKeyID())
	}
}

func TestNewSignerFromPEM(t *testing.T) {
	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	signer, err := NewSignerFromPEM("test-api-key-id", string(pemBytes))
	if err != nil {
		t.Fatalf("NewSignerFromPEM failed: %v", err)
	}

	if signer == nil {
		t.Fatal("NewSignerFromPEM returned nil signer")
	}
}

func TestNewSignerFromPEM_InvalidPEM(t *testing.T) {
	_, err := NewSignerFromPEM("test-api-key-id", "invalid-pem-data")
	if err == nil {
		t.Fatal("expected error for invalid PEM data")
	}
}

func TestSign(t *testing.T) {
	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("test-api-key-id", privateKey)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	timestamp := time.Now().UTC()
	method := "GET"
	path := "/trade-api/v2/markets"

	signature, err := signer.Sign(timestamp, method, path, "")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if signature == "" {
		t.Fatal("Sign returned empty signature")
	}

	// Signature should be base64 encoded
	if len(signature) < 100 {
		t.Errorf("signature seems too short: %d chars", len(signature))
	}
}

func TestSignWithBody(t *testing.T) {
	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("test-api-key-id", privateKey)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	timestamp := time.Now().UTC()
	method := "POST"
	path := "/trade-api/v2/orders"
	body := `{"ticker":"BTC-100K","side":"yes","count":10}`

	signature, err := signer.Sign(timestamp, method, path, body)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if signature == "" {
		t.Fatal("Sign returned empty signature")
	}
}

func TestBuildAuthMessage(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		method    string
		path      string
		body      string
		expected  string
	}{
		{
			name:      "GET request without body",
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			method:    "GET",
			path:      "/trade-api/v2/markets",
			body:      "",
			expected:  "2024-01-15T12:00:00ZGET/trade-api/v2/markets",
		},
		{
			name:      "POST request with body",
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			method:    "POST",
			path:      "/trade-api/v2/orders",
			body:      `{"ticker":"TEST"}`,
			expected:  "2024-01-15T12:00:00ZPOST/trade-api/v2/orders{\"ticker\":\"TEST\"}",
		},
		{
			name:      "DELETE request",
			timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			method:    "DELETE",
			path:      "/trade-api/v2/orders/abc123",
			body:      "",
			expected:  "2024-01-15T12:00:00ZDELETE/trade-api/v2/orders/abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := BuildAuthMessage(tt.timestamp, tt.method, tt.path, tt.body)
			if msg != tt.expected {
				t.Errorf("expected message:\n%s\ngot:\n%s", tt.expected, msg)
			}
		})
	}
}

func TestAuthHeader(t *testing.T) {
	privateKey, err := generateTestKey()
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	signer, err := NewSigner("my-key-id", privateKey)
	if err != nil {
		t.Fatalf("NewSigner failed: %v", err)
	}

	header := signer.AuthHeader("base64signature")
	expected := "KALSHI-API-KEY my-key-id:base64signature"
	if header != expected {
		t.Errorf("expected header '%s', got '%s'", expected, header)
	}
}

func generateTestKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
