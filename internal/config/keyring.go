package config

import (
	"encoding/json"
	"fmt"

	"github.com/99designs/keyring"
)

const (
	serviceName = "kalshi-cli"
	credsKey    = "credentials"
)

type Credentials struct {
	APIKeyID   string `json:"api_key_id"`
	PrivateKey string `json:"private_key"`
}

type KeyringStore struct {
	ring keyring.Keyring
}

func NewKeyringStore() (*KeyringStore, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
		KeychainTrustApplication: true,
		KeychainAccessibleWhenUnlocked: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	return &KeyringStore{ring: ring}, nil
}

func (k *KeyringStore) SaveCredentials(creds Credentials) error {
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	err = k.ring.Set(keyring.Item{
		Key:  credsKey,
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

func (k *KeyringStore) GetCredentials() (*Credentials, error) {
	item, err := k.ring.Get(credsKey)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(item.Data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

func (k *KeyringStore) DeleteCredentials() error {
	err := k.ring.Remove(credsKey)
	if err != nil && err != keyring.ErrKeyNotFound {
		return fmt.Errorf("failed to delete credentials: %w", err)
	}
	return nil
}

func (k *KeyringStore) HasCredentials() bool {
	_, err := k.ring.Get(credsKey)
	return err == nil
}
