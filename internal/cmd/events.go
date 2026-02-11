package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Manage events",
	Long: `Commands for listing, viewing, and managing Kalshi events.

An event groups related markets (e.g., "S&P 500 close on Feb 7" has
multiple price-bracket markets under it).`,
	Example: `  kalshi-cli events list --status active
  kalshi-cli events get INXD-25FEB07
  kalshi-cli events candlesticks INXD-25FEB07 --start 2025-02-01T00:00:00Z --end 2025-02-07T00:00:00Z`,
}

var eventsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List events",
	Long:  `List events with optional filtering by status.`,
	Example: `  kalshi-cli events list
  kalshi-cli events list --status active --limit 20
  kalshi-cli events list --json`,
	RunE: runEventsList,
}

var eventsGetCmd = &cobra.Command{
	Use:   "get <event-ticker>",
	Short: "Get event details",
	Long: `Get detailed information about a specific event by ticker.

Use 'kalshi-cli events list' to find event tickers.`,
	Example: `  kalshi-cli events get INXD-25FEB07`,
	Args:    cobra.ExactArgs(1),
	RunE:    runEventsGet,
}

var eventsCandlesticksCmd = &cobra.Command{
	Use:   "candlesticks <event-ticker>",
	Short: "Get event candlesticks",
	Long: `Get candlestick (OHLCV) data for an event.

The --series flag is optional; if omitted, the series ticker is
auto-resolved from the event. Requires --start and --end timestamps.

Supported periods: 1m, 1h, 1d`,
	Example: `  kalshi-cli events candlesticks INXD-25FEB07 --start 2025-02-01T00:00:00Z --end 2025-02-07T00:00:00Z
  kalshi-cli events candlesticks INXD-25FEB07 --period 1d --start 2025-01-01T00:00:00Z --end 2025-02-01T00:00:00Z
  kalshi-cli events candlesticks INXD-25FEB07 --series INXD --period 1h --start 2025-02-06T00:00:00Z --end 2025-02-07T00:00:00Z`,
	Args: cobra.ExactArgs(1),
	RunE: runEventsCandlesticks,
}

var multivariateCmd = &cobra.Command{
	Use:   "multivariate",
	Short: "Manage multivariate events",
	Long:  `Commands for listing and viewing multivariate events.`,
}

var multivariateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List multivariate events",
	Long:  `List all multivariate events.`,
	RunE:  runMultivariateList,
}

var multivariateGetCmd = &cobra.Command{
	Use:   "get <ticker>",
	Short: "Get multivariate event details",
	Long:  `Get detailed information about a specific multivariate event.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMultivariateGet,
}

var (
	eventsStatus          string
	eventsLimit           int
	eventsCursor          string
	eventSeriesTicker     string
	candlesticksPeriod    string
	candlesticksStartTime string
	candlesticksEndTime   string
	multivariateStatus    string
	multivariateLimit     int
	multivariateCursor    string
)

func init() {
	rootCmd.AddCommand(eventsCmd)

	eventsListCmd.Flags().StringVar(&eventsStatus, "status", "", "filter by status (active, closed, settled)")
	eventsListCmd.Flags().IntVar(&eventsLimit, "limit", 50, "maximum number of events to return")
	eventsListCmd.Flags().StringVar(&eventsCursor, "cursor", "", "pagination cursor")

	eventsCandlesticksCmd.Flags().StringVar(&eventSeriesTicker, "series", "", "series ticker (auto-resolved from event if not provided)")
	eventsCandlesticksCmd.Flags().StringVar(&candlesticksPeriod, "period", "1h", "candlestick period (1m, 1h, 1d)")
	eventsCandlesticksCmd.Flags().StringVar(&candlesticksStartTime, "start", "", "start time (RFC3339 format)")
	eventsCandlesticksCmd.Flags().StringVar(&candlesticksEndTime, "end", "", "end time (RFC3339 format)")

	multivariateListCmd.Flags().StringVar(&multivariateStatus, "status", "", "filter by status")
	multivariateListCmd.Flags().IntVar(&multivariateLimit, "limit", 50, "maximum number of events to return")
	multivariateListCmd.Flags().StringVar(&multivariateCursor, "cursor", "", "pagination cursor")

	multivariateCmd.AddCommand(multivariateListCmd)
	multivariateCmd.AddCommand(multivariateGetCmd)

	eventsCmd.AddCommand(eventsListCmd)
	eventsCmd.AddCommand(eventsGetCmd)
	eventsCmd.AddCommand(eventsCandlesticksCmd)
	eventsCmd.AddCommand(multivariateCmd)
}

func runEventsList(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	params := api.ListEventsParams{
		Status: eventsStatus,
		Limit:  eventsLimit,
		Cursor: eventsCursor,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	events, cursor, err := client.ListEvents(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}

	outputFormat := GetOutputFormat()

	return ui.Output(
		outputFormat,
		func() { renderEventsTable(events, cursor) },
		createEventsResponse(events, cursor),
		func() { renderEventsPlain(events) },
	)
}

func runEventsGet(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	event, err := client.GetEvent(ctx, ticker)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	outputFormat := GetOutputFormat()

	return ui.Output(
		outputFormat,
		func() { renderEventDetails(event) },
		event,
		func() { renderEventPlain(event) },
	)
}

func runEventsCandlesticks(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seriesTicker, err := resolveSeriesTicker(ctx, client, ticker, eventSeriesTicker)
	if err != nil {
		return err
	}

	params := api.CandlesticksParams{
		Ticker:       ticker,
		SeriesTicker: seriesTicker,
		Period:       candlesticksPeriod,
	}

	if candlesticksStartTime != "" {
		t, parseErr := time.Parse(time.RFC3339, candlesticksStartTime)
		if parseErr != nil {
			return fmt.Errorf("invalid start time format: %w", parseErr)
		}
		params.StartTime = &t
	}

	if candlesticksEndTime != "" {
		t, parseErr := time.Parse(time.RFC3339, candlesticksEndTime)
		if parseErr != nil {
			return fmt.Errorf("invalid end time format: %w", parseErr)
		}
		params.EndTime = &t
	}

	candlesticks, err := client.GetEventCandlesticks(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to get candlesticks: %w", err)
	}

	outputFormat := GetOutputFormat()

	return ui.Output(
		outputFormat,
		func() { renderCandlesticksTable(candlesticks) },
		candlesticks,
		func() { renderCandlesticksPlain(candlesticks) },
	)
}

// resolveSeriesTicker returns the series ticker for the candlesticks API call.
// If explicitSeries is provided (via --series flag), it is returned directly.
// Otherwise, the event is fetched to extract its SeriesTicker field.
func resolveSeriesTicker(ctx context.Context, client *api.Client, ticker string, explicitSeries string) (string, error) {
	if explicitSeries != "" {
		return explicitSeries, nil
	}

	event, err := client.GetEvent(ctx, ticker)
	if err != nil {
		return "", fmt.Errorf("failed to resolve series ticker for event %s: %w", ticker, err)
	}

	if event.SeriesTicker == "" {
		return "", fmt.Errorf("event %s has no series ticker; please provide --series explicitly", ticker)
	}

	return event.SeriesTicker, nil
}

func runMultivariateList(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	params := api.ListMultivariateParams{
		Status: multivariateStatus,
		Limit:  multivariateLimit,
		Cursor: multivariateCursor,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	events, cursor, err := client.ListMultivariateEvents(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list multivariate events: %w", err)
	}

	outputFormat := GetOutputFormat()

	return ui.Output(
		outputFormat,
		func() { renderMultivariateEventsTable(events, cursor) },
		createMultivariateResponse(events, cursor),
		func() { renderMultivariateEventsPlain(events) },
	)
}

func runMultivariateGet(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	event, err := client.GetMultivariateEvent(ctx, ticker)
	if err != nil {
		return fmt.Errorf("failed to get multivariate event: %w", err)
	}

	outputFormat := GetOutputFormat()

	return ui.Output(
		outputFormat,
		func() { renderMultivariateEventDetails(event) },
		event,
		func() { renderMultivariateEventPlain(event) },
	)
}

// Table Rendering

func renderEventsTable(events []models.Event, cursor string) {
	headers := []string{"Ticker", "Title", "Category", "Markets"}
	rows := make([][]string, 0, len(events))

	for _, e := range events {
		rows = append(rows, []string{
			e.EventTicker,
			truncateEventString(e.Title, 40),
			e.Category,
			strconv.Itoa(len(e.Markets)),
		})
	}

	ui.RenderTable(headers, rows)

	if cursor != "" {
		fmt.Printf("\nMore results available. Use --cursor %s to continue.\n", cursor)
	}
}

func renderEventDetails(event *models.Event) {
	pairs := [][]string{
		{"Ticker", event.EventTicker},
		{"Series", event.SeriesTicker},
		{"Title", event.Title},
		{"Subtitle", event.SubTitle},
		{"Category", event.Category},
		{"Mutually Exclusive", formatEventBool(event.MutuallyExclusive)},
		{"Strike Date", formatEventTimePtr(event.StrikeDate)},
		{"Markets Count", strconv.Itoa(len(event.Markets))},
	}

	ui.RenderKeyValue(pairs)

	if len(event.Markets) > 0 {
		fmt.Println("\nMarkets:")
		for _, m := range event.Markets {
			fmt.Printf("  - %s\n", m)
		}
	}
}

func renderCandlesticksTable(candlesticks []models.Candlestick) {
	headers := []string{"Time", "Open", "High", "Low", "Close", "Volume", "OI"}
	rows := make([][]string, 0, len(candlesticks))

	for _, c := range candlesticks {
		rows = append(rows, []string{
			c.PeriodEnd.Format("2006-01-02 15:04"),
			ui.FormatPrice(c.Open),
			ui.FormatPrice(c.High),
			ui.FormatPrice(c.Low),
			ui.FormatPrice(c.Close),
			strconv.Itoa(c.Volume),
			strconv.Itoa(c.OpenInterest),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderMultivariateEventsTable(events []models.MultivariateEvent, cursor string) {
	headers := []string{"Ticker", "Title", "Status", "Lookup Type"}
	rows := make([][]string, 0, len(events))

	for _, e := range events {
		rows = append(rows, []string{
			e.Ticker,
			truncateEventString(e.Title, 40),
			formatEventStatus(e.Status),
			e.LookupType,
		})
	}

	ui.RenderTable(headers, rows)

	if cursor != "" {
		fmt.Printf("\nMore results available. Use --cursor %s to continue.\n", cursor)
	}
}

func renderMultivariateEventDetails(event *models.MultivariateEvent) {
	pairs := [][]string{
		{"Ticker", event.Ticker},
		{"Title", event.Title},
		{"Description", event.Description},
		{"Status", formatEventStatus(event.Status)},
		{"Lookup Type", event.LookupType},
	}

	ui.RenderKeyValue(pairs)

	if len(event.LookupTable) > 0 {
		fmt.Println("\nLookup Table:")
		for i, item := range event.LookupTable {
			fmt.Printf("  %d. %s\n", i+1, item)
		}
	}
}

// Plain Rendering

func renderEventsPlain(events []models.Event) {
	for _, e := range events {
		fmt.Printf("%s\t%s\t%s\t%d\n",
			e.EventTicker, e.Title, e.Category, len(e.Markets))
	}
}

func renderEventPlain(event *models.Event) {
	fmt.Printf("ticker=%s\n", event.EventTicker)
	fmt.Printf("series=%s\n", event.SeriesTicker)
	fmt.Printf("title=%s\n", event.Title)
	fmt.Printf("category=%s\n", event.Category)
	fmt.Printf("markets_count=%d\n", len(event.Markets))
	if len(event.Markets) > 0 {
		fmt.Printf("markets=%s\n", strings.Join(event.Markets, ","))
	}
}

func renderCandlesticksPlain(candlesticks []models.Candlestick) {
	for _, c := range candlesticks {
		fmt.Printf("%s\t%d\t%d\t%d\t%d\t%d\t%d\n",
			c.PeriodEnd.Format(time.RFC3339),
			c.Open, c.High, c.Low, c.Close, c.Volume, c.OpenInterest)
	}
}

func renderMultivariateEventsPlain(events []models.MultivariateEvent) {
	for _, e := range events {
		fmt.Printf("%s\t%s\t%s\t%s\n", e.Ticker, e.Title, e.Status, e.LookupType)
	}
}

func renderMultivariateEventPlain(event *models.MultivariateEvent) {
	fmt.Printf("ticker=%s\n", event.Ticker)
	fmt.Printf("title=%s\n", event.Title)
	fmt.Printf("description=%s\n", event.Description)
	fmt.Printf("status=%s\n", event.Status)
	fmt.Printf("lookup_type=%s\n", event.LookupType)
}

// Response Helpers

type eventsListResponse struct {
	Events []models.Event `json:"events"`
	Cursor string         `json:"cursor,omitempty"`
}

type multivariateListResponse struct {
	Events []models.MultivariateEvent `json:"multivariate_events"`
	Cursor string                     `json:"cursor,omitempty"`
}

func createEventsResponse(events []models.Event, cursor string) eventsListResponse {
	return eventsListResponse{
		Events: events,
		Cursor: cursor,
	}
}

func createMultivariateResponse(events []models.MultivariateEvent, cursor string) multivariateListResponse {
	return multivariateListResponse{
		Events: events,
		Cursor: cursor,
	}
}

// Formatting Helpers (prefixed with "Event" to avoid conflicts with other cmd files)

func truncateEventString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatEventStatus(status string) string {
	return strings.ToUpper(status)
}

func formatEventBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func formatEventTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatEventTimePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return formatEventTime(*t)
}
