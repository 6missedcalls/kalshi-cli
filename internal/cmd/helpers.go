package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/config"
)

// Common helper functions shared across commands

// createClient creates an API client using stored credentials
func createClient() (*api.Client, error) {
	keyring, err := config.NewKeyringStore()
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	creds, err := keyring.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds == nil {
		return nil, fmt.Errorf("not logged in. Run 'kalshi-cli auth login' first")
	}

	signer, err := api.NewSignerFromPEM(creds.APIKeyID, creds.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	return api.NewClient(cfg, signer), nil
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
