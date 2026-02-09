package api

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
)

// Signer handles RSA signature generation for Kalshi API authentication
type Signer struct {
	apiKeyID   string
	privateKey *rsa.PrivateKey
}

// NewSigner creates a new signer with the given API key ID and private key
func NewSigner(apiKeyID string, privateKey *rsa.PrivateKey) (*Signer, error) {
	if apiKeyID == "" {
		return nil, errors.New("API key ID is required")
	}
	if privateKey == nil {
		return nil, errors.New("private key is required")
	}
	return &Signer{
		apiKeyID:   apiKeyID,
		privateKey: privateKey,
	}, nil
}

// NewSignerFromPEM creates a new signer from a PEM-encoded private key
func NewSignerFromPEM(apiKeyID string, pemData string) (*Signer, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey
	var err error

	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", parseErr)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return NewSigner(apiKeyID, privateKey)
}

// APIKeyID returns the API key ID
func (s *Signer) APIKeyID() string {
	return s.apiKeyID
}

// Sign generates a signature for the given request parameters
func (s *Signer) Sign(timestamp time.Time, method, path, body string) (string, error) {
	message := BuildAuthMessage(timestamp, method, path, body)

	hash := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// BuildAuthMessage constructs the message to be signed
// Format: timestamp + method + path + body
func BuildAuthMessage(timestamp time.Time, method, path, body string) string {
	ts := timestamp.UTC().Format("2006-01-02T15:04:05Z")
	return ts + method + path + body
}

// AuthHeader returns the Authorization header value
func (s *Signer) AuthHeader(signature string) string {
	return fmt.Sprintf("KALSHI-API-KEY %s:%s", s.apiKeyID, signature)
}

// TimestampHeader returns the timestamp in the format required by Kalshi
func TimestampHeader(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

// GenerateKeyPair generates a new RSA key pair for API authentication
func GenerateKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 4096)
}

// EncodePrivateKeyPEM encodes an RSA private key to PEM format
func EncodePrivateKeyPEM(key *rsa.PrivateKey) string {
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return string(pem.EncodeToMemory(block))
}

// EncodePublicKeyPEM encodes an RSA public key to PEM format
func EncodePublicKeyPEM(key *rsa.PublicKey) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(block)), nil
}
