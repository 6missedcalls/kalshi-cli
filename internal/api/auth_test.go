package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strconv"
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

	signature, err := signer.Sign(timestamp, method, path)
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
	signature, err := signer.Sign(timestamp, method, path)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if signature == "" {
		t.Fatal("Sign returned empty signature")
	}
}

func TestBuildAuthMessage(t *testing.T) {
	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	tsMs := strconv.FormatInt(ts.UnixMilli(), 10)

	tests := []struct {
		name      string
		timestamp time.Time
		method    string
		path      string
		expected  string
	}{
		{
			name:      "GET request",
			timestamp: ts,
			method:    "GET",
			path:      "/trade-api/v2/markets",
			expected:  tsMs + "GET/trade-api/v2/markets",
		},
		{
			name:      "POST request",
			timestamp: ts,
			method:    "POST",
			path:      "/trade-api/v2/orders",
			expected:  tsMs + "POST/trade-api/v2/orders",
		},
		{
			name:      "DELETE request",
			timestamp: ts,
			method:    "DELETE",
			path:      "/trade-api/v2/orders/abc123",
			expected:  tsMs + "DELETE/trade-api/v2/orders/abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := BuildAuthMessage(tt.timestamp, tt.method, tt.path)
			if msg != tt.expected {
				t.Errorf("expected message:\n%s\ngot:\n%s", tt.expected, msg)
			}
		})
	}
}


func generateTestKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
