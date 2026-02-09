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
	Long:  `Commands for listing, viewing, and managing Kalshi events.`,
}

var eventsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List events",
	Long:  `List events with optional filtering by status.`,
	RunE:  runEventsList,
}

var eventsGetCmd = &cobra.Command{
	Use:   "get <ticker>",
	Short: "Get event details",
	Long:  `Get detailed information about a specific event by ticker.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEventsGet,
}

var eventsCandlesticksCmd = &cobra.Command{
	Use:   "candlesticks <ticker>",
	Short: "Get event candlesticks",
	Long:  `Get candlestick (OHLCV) data for an event.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runEventsCandlesticks,
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

	eventsCandlesticksCmd.Flags().StringVar(&candlesticksPeriod, "period", "1h", "candlestick period (1m, 5m, 15m, 1h, 4h, 1d)")
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

	params := api.CandlesticksParams{
		Ticker: ticker,
		Period: candlesticksPeriod,
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
	headers := []string{"Ticker", "Title", "Category", "Status", "Markets"}
	rows := make([][]string, 0, len(events))

	for _, e := range events {
		rows = append(rows, []string{
			e.EventTicker,
			truncateEventString(e.Title, 40),
			e.Category,
			formatEventStatus(e.Status),
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
		{"Subtitle", event.Subtitle},
		{"Category", event.Category},
		{"Status", formatEventStatus(event.Status)},
		{"Mutually Exclusive", formatEventBool(event.MutuallyExclusive)},
		{"Strike Date", formatEventTime(event.StrikeDate)},
		{"Expected Expiration", formatEventTime(event.ExpectedExpiration)},
		{"Created", formatEventTime(event.CreatedTime)},
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
		fmt.Printf("%s\t%s\t%s\t%s\t%d\n",
			e.EventTicker, e.Title, e.Category, e.Status, len(e.Markets))
	}
}

func renderEventPlain(event *models.Event) {
	fmt.Printf("ticker=%s\n", event.EventTicker)
	fmt.Printf("series=%s\n", event.SeriesTicker)
	fmt.Printf("title=%s\n", event.Title)
	fmt.Printf("category=%s\n", event.Category)
	fmt.Printf("status=%s\n", event.Status)
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
