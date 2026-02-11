package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type OutputFormat int

const (
	FormatTable OutputFormat = iota
	FormatJSON
	FormatPlain
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	successColor   = lipgloss.Color("#10B981")
	warningColor   = lipgloss.Color("#F59E0B")
	errorColor     = lipgloss.Color("#EF4444")
	mutedColor     = lipgloss.Color("#6B7280")
	accentColor    = lipgloss.Color("#3B82F6")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(errorColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	BoldStyle = lipgloss.NewStyle().
			Bold(true)

	// Price styles
	PriceUpStyle = lipgloss.NewStyle().
			Foreground(successColor)

	PriceDownStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Status styles
	StatusOpenStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	StatusClosedStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	StatusActiveStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// Environment indicator
	DemoStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FEF3C7")).
			Foreground(lipgloss.Color("#92400E")).
			Padding(0, 1)

	ProdStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FEE2E2")).
			Foreground(lipgloss.Color("#991B1B")).
			Padding(0, 1).
			Bold(true)
)

func FormatPrice(cents int) string {
	if cents < 0 {
		return fmt.Sprintf("-$%.2f", float64(-cents)/100.0)
	}
	return fmt.Sprintf("$%.2f", float64(cents)/100.0)
}

func FormatPriceStyled(cents int, positive bool) string {
	absCents := cents
	if absCents < 0 {
		absCents = -absCents
	}
	dollars := float64(absCents) / 100.0
	style := PriceDownStyle
	prefix := "-"
	if positive {
		style = PriceUpStyle
		prefix = "+"
	}
	return style.Render(fmt.Sprintf("%s$%.2f", prefix, dollars))
}

func FormatPercent(value float64) string {
	return fmt.Sprintf("%.1f%%", value*100)
}

func FormatQuantity(qty int) string {
	return fmt.Sprintf("%d", qty)
}
