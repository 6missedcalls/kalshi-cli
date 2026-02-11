package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/spf13/viper"
)

// Common helper functions shared across commands

// createClient creates an API client using stored credentials.
// It tries config file credentials first (api_key_id + private_key_path),
// then environment variables, then the keyring as a last resort.
// Config/env is checked first because keyring can hang in headless environments.
func createClient() (*api.Client, error) {
	// Try config file first (fast, no GUI prompts)
	apiKeyID := viper.GetString("api_key_id")
	privateKeyPath := viper.GetString("private_key_path")

	// Also check env vars
	if apiKeyID == "" {
		apiKeyID = os.Getenv("KALSHI_API_KEY_ID")
	}
	if privateKeyPath == "" {
		privateKeyPath = os.Getenv("KALSHI_PRIVATE_KEY_FILE")
	}

	if apiKeyID != "" && privateKeyPath != "" {
		pemData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file %s: %w", privateKeyPath, err)
		}

		signer, err := api.NewSignerFromPEM(apiKeyID, string(pemData))
		if err != nil {
			return nil, fmt.Errorf("failed to create signer from key file: %w", err)
		}

		return api.NewClient(cfg, signer), nil
	}

	// Also support KALSHI_PRIVATE_KEY env var (PEM content directly)
	privateKeyPEM := os.Getenv("KALSHI_PRIVATE_KEY")
	if apiKeyID != "" && privateKeyPEM != "" {
		signer, err := api.NewSignerFromPEM(apiKeyID, privateKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to create signer from env var: %w", err)
		}

		return api.NewClient(cfg, signer), nil
	}

	// Last resort: try keyring (may hang in headless environments)
	keyring, err := config.NewKeyringStore()
	if err == nil {
		creds, err := keyring.GetCredentials()
		if err == nil && creds != nil {
			signer, err := api.NewSignerFromPEM(creds.APIKeyID, creds.PrivateKey)
			if err == nil {
				return api.NewClient(cfg, signer), nil
			}
		}
	}

	return nil, fmt.Errorf("not logged in. Set api_key_id + private_key_path in ~/.kalshi/config.yaml, or run 'kalshi-cli auth login'")
}

// formatTimeStr formats a time.Time for display
func formatTimeStr(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04")
}

// truncateStr truncates a string to max length with ellipsis
func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// formatMarketStatus formats a market status for display
func formatMarketStatus(status string) string {
	switch status {
	case "open":
		return "Open"
	case "closed":
		return "Closed"
	case "settled":
		return "Settled"
	default:
		return status
	}
}

// formatSideStr formats a side (yes/no) for display
func formatSideStr(side string) string {
	switch side {
	case "yes":
		return "Yes"
	case "no":
		return "No"
	default:
		return side
	}
}

// formatCents formats cents as dollars
func formatCents(cents int) string {
	return fmt.Sprintf("$%.2f", float64(cents)/100)
}

// confirmAction prompts for confirmation unless --yes flag is set
func confirmAction(action string) bool {
	if yesFlag {
		return true
	}

	fmt.Printf("Are you sure you want to %s? [y/N]: ", action)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y" || response == "yes" || response == "Yes"
}

// withTimeout returns a context with the configured timeout
func withTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	timeout := 30 * time.Second
	if cfg != nil && cfg.API.Timeout > 0 {
		timeout = cfg.API.Timeout
	}
	return context.WithTimeout(parent, timeout)
}
