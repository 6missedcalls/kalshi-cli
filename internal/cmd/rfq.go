package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var (
	rfqStatus    string
	rfqMarket    string
	rfqSide      string
	rfqQuantity  int
	quoteRFQID   string
	quotePrice   int
)

var rfqCmd = &cobra.Command{
	Use:   "rfq",
	Short: "Manage RFQs (Request for Quotes)",
	Long:  `Create, list, and manage RFQs for block trading on Kalshi.`,
}

var rfqListCmd = &cobra.Command{
	Use:   "list",
	Short: "List RFQs",
	Long:  `List all RFQs, optionally filtered by status.`,
	Example: `  kalshi-cli rfq list
  kalshi-cli rfq list --status open`,
	RunE: runRFQList,
}

var rfqGetCmd = &cobra.Command{
	Use:   "get <rfq-id>",
	Short: "Get RFQ details",
	Long:  `Get detailed information about a specific RFQ.`,
	Args:  cobra.ExactArgs(1),
	Example: `  kalshi-cli rfq get rfq_abc123`,
	RunE: runRFQGet,
}

var rfqCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new RFQ",
	Long:  `Create a new Request for Quote for block trading.`,
	Example: `  kalshi-cli rfq create --market INXD-25FEB07 --side yes --qty 1000`,
	RunE: runRFQCreate,
}

var rfqDeleteCmd = &cobra.Command{
	Use:   "delete <rfq-id>",
	Short: "Delete an RFQ",
	Long:  `Delete an existing RFQ by ID.`,
	Args:  cobra.ExactArgs(1),
	Example: `  kalshi-cli rfq delete rfq_abc123`,
	RunE: runRFQDelete,
}

var quotesCmd = &cobra.Command{
	Use:   "quotes",
	Short: "Manage quotes on RFQs",
	Long:  `Create, list, accept, and confirm quotes on RFQs.`,
}

var quotesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List quotes",
	Long:  `List all quotes, optionally filtered by RFQ ID.`,
	Example: `  kalshi-cli quotes list
  kalshi-cli quotes list --rfq-id rfq_abc123`,
	RunE: runQuotesList,
}

var quotesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a quote on an RFQ",
	Long:  `Create a new quote on an existing RFQ.`,
	Example: `  kalshi-cli quotes create --rfq rfq_abc123 --price 65`,
	RunE: runQuotesCreate,
}

var quotesAcceptCmd = &cobra.Command{
	Use:   "accept <quote-id>",
	Short: "Accept a quote",
	Long:  `Accept a quote that was offered on your RFQ.`,
	Args:  cobra.ExactArgs(1),
	Example: `  kalshi-cli quotes accept quote_xyz789`,
	RunE: runQuotesAccept,
}

var quotesConfirmCmd = &cobra.Command{
	Use:   "confirm <quote-id>",
	Short: "Confirm a quote",
	Long:  `Confirm a quote after it has been accepted.`,
	Args:  cobra.ExactArgs(1),
	Example: `  kalshi-cli quotes confirm quote_xyz789`,
	RunE: runQuotesConfirm,
}

func init() {
	// RFQ list flags
	rfqListCmd.Flags().StringVar(&rfqStatus, "status", "", "Filter by status (e.g., open, closed)")

	// RFQ create flags
	rfqCreateCmd.Flags().StringVar(&rfqMarket, "market", "", "Market ticker (required)")
	rfqCreateCmd.Flags().StringVar(&rfqSide, "side", "", "Side: yes or no (required)")
	rfqCreateCmd.Flags().IntVar(&rfqQuantity, "qty", 0, "Quantity (required)")
	rfqCreateCmd.MarkFlagRequired("market")
	rfqCreateCmd.MarkFlagRequired("side")
	rfqCreateCmd.MarkFlagRequired("qty")

	// Quotes list flags
	quotesListCmd.Flags().StringVar(&quoteRFQID, "rfq-id", "", "Filter by RFQ ID")

	// Quotes create flags
	quotesCreateCmd.Flags().StringVar(&quoteRFQID, "rfq", "", "RFQ ID (required)")
	quotesCreateCmd.Flags().IntVar(&quotePrice, "price", 0, "Price in cents (required)")
	quotesCreateCmd.MarkFlagRequired("rfq")
	quotesCreateCmd.MarkFlagRequired("price")

	// Add subcommands to rfq
	rfqCmd.AddCommand(rfqListCmd)
	rfqCmd.AddCommand(rfqGetCmd)
	rfqCmd.AddCommand(rfqCreateCmd)
	rfqCmd.AddCommand(rfqDeleteCmd)

	// Add subcommands to quotes
	quotesCmd.AddCommand(quotesListCmd)
	quotesCmd.AddCommand(quotesCreateCmd)
	quotesCmd.AddCommand(quotesAcceptCmd)
	quotesCmd.AddCommand(quotesConfirmCmd)

	// Register with root
	rootCmd.AddCommand(rfqCmd)
	rootCmd.AddCommand(quotesCmd)
}

func runRFQList(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.ListRFQs(ctx, api.ListRFQsParams{
		Status: rfqStatus,
	})
	if err != nil {
		return err
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderRFQsTable(result.RFQs) },
		result.RFQs,
		func() { renderRFQsPlain(result.RFQs) },
	)
}

func runRFQGet(cmd *cobra.Command, args []string) error {
	rfqID := args[0]

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rfq, err := client.GetRFQ(ctx, rfqID)
	if err != nil {
		return err
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderRFQDetail(rfq) },
		rfq,
		func() { renderRFQDetailPlain(rfq) },
	)
}

func runRFQCreate(cmd *cobra.Command, args []string) error {
	side := strings.ToLower(rfqSide)
	if side != "yes" && side != "no" {
		return fmt.Errorf("side must be 'yes' or 'no', got: %s", rfqSide)
	}

	if rfqQuantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rfq, err := client.CreateRFQ(ctx, models.CreateRFQRequest{
		Ticker:   rfqMarket,
		Side:     side,
		Quantity: rfqQuantity,
	})
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("RFQ created: %s", rfq.RFQID))

	return ui.Output(
		GetOutputFormat(),
		func() { renderRFQDetail(rfq) },
		rfq,
		func() { renderRFQDetailPlain(rfq) },
	)
}

func runRFQDelete(cmd *cobra.Command, args []string) error {
	rfqID := args[0]

	if !confirmAction(fmt.Sprintf("Delete RFQ %s?", rfqID)) {
		PrintWarning("Cancelled")
		return nil
	}

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.DeleteRFQ(ctx, rfqID); err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("RFQ %s deleted", rfqID))
	return nil
}

func runQuotesList(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.ListQuotes(ctx, api.ListQuotesParams{
		RFQID: quoteRFQID,
	})
	if err != nil {
		return err
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderQuotesTable(result.Quotes) },
		result.Quotes,
		func() { renderQuotesPlain(result.Quotes) },
	)
}

func runQuotesCreate(cmd *cobra.Command, args []string) error {
	if quotePrice <= 0 || quotePrice >= 100 {
		return fmt.Errorf("price must be between 1 and 99 cents, got: %d", quotePrice)
	}

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	quote, err := client.CreateQuote(ctx, models.CreateQuoteRequest{
		RFQID: quoteRFQID,
		Price: quotePrice,
	})
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Quote created: %s", quote.QuoteID))

	return ui.Output(
		GetOutputFormat(),
		func() { renderQuoteDetail(quote) },
		quote,
		func() { renderQuoteDetailPlain(quote) },
	)
}

func runQuotesAccept(cmd *cobra.Command, args []string) error {
	quoteID := args[0]

	if !confirmAction(fmt.Sprintf("Accept quote %s?", quoteID)) {
		PrintWarning("Cancelled")
		return nil
	}

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	quote, err := client.AcceptQuote(ctx, quoteID)
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Quote %s accepted", quote.QuoteID))

	return ui.Output(
		GetOutputFormat(),
		func() { renderQuoteDetail(quote) },
		quote,
		func() { renderQuoteDetailPlain(quote) },
	)
}

func runQuotesConfirm(cmd *cobra.Command, args []string) error {
	quoteID := args[0]

	if !confirmAction(fmt.Sprintf("Confirm quote %s?", quoteID)) {
		PrintWarning("Cancelled")
		return nil
	}

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	quote, err := client.ConfirmQuote(ctx, quoteID)
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Quote %s confirmed", quote.QuoteID))

	return ui.Output(
		GetOutputFormat(),
		func() { renderQuoteDetail(quote) },
		quote,
		func() { renderQuoteDetailPlain(quote) },
	)
}

// RFQ rendering functions

func renderRFQsTable(rfqs []models.RFQ) {
	if len(rfqs) == 0 {
		fmt.Println("No RFQs found")
		return
	}

	headers := []string{"ID", "Ticker", "Side", "Quantity", "Status", "Created", "Expires"}
	var rows [][]string

	for _, rfq := range rfqs {
		rows = append(rows, []string{
			rfq.RFQID,
			rfq.Ticker,
			formatSide(rfq.Side),
			fmt.Sprintf("%d", rfq.Quantity),
			formatStatus(rfq.Status),
			rfq.CreatedTime.Format("2006-01-02 15:04"),
			rfq.ExpiresTime.Format("2006-01-02 15:04"),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderRFQsPlain(rfqs []models.RFQ) {
	for _, rfq := range rfqs {
		fmt.Printf("%s\t%s\t%s\t%d\t%s\n",
			rfq.RFQID, rfq.Ticker, rfq.Side, rfq.Quantity, rfq.Status)
	}
}

func renderRFQDetail(rfq *models.RFQ) {
	pairs := [][]string{
		{"RFQ ID", rfq.RFQID},
		{"Ticker", rfq.Ticker},
		{"Side", formatSide(rfq.Side)},
		{"Quantity", fmt.Sprintf("%d", rfq.Quantity)},
		{"Status", formatStatus(rfq.Status)},
		{"Created", rfq.CreatedTime.Format("2006-01-02 15:04:05")},
		{"Expires", rfq.ExpiresTime.Format("2006-01-02 15:04:05")},
	}
	ui.RenderKeyValue(pairs)
}

func renderRFQDetailPlain(rfq *models.RFQ) {
	fmt.Printf("rfq_id=%s ticker=%s side=%s qty=%d status=%s\n",
		rfq.RFQID, rfq.Ticker, rfq.Side, rfq.Quantity, rfq.Status)
}

// Quote rendering functions

func renderQuotesTable(quotes []models.Quote) {
	if len(quotes) == 0 {
		fmt.Println("No quotes found")
		return
	}

	headers := []string{"Quote ID", "RFQ ID", "Ticker", "Side", "Price", "Qty", "Status", "Created"}
	var rows [][]string

	for _, quote := range quotes {
		rows = append(rows, []string{
			quote.QuoteID,
			quote.RFQID,
			quote.Ticker,
			formatSide(quote.Side),
			ui.FormatPrice(quote.Price),
			fmt.Sprintf("%d", quote.Quantity),
			formatStatus(quote.Status),
			quote.CreatedTime.Format("2006-01-02 15:04"),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderQuotesPlain(quotes []models.Quote) {
	for _, quote := range quotes {
		fmt.Printf("%s\t%s\t%s\t%s\t%d\t%d\t%s\n",
			quote.QuoteID, quote.RFQID, quote.Ticker, quote.Side,
			quote.Price, quote.Quantity, quote.Status)
	}
}

func renderQuoteDetail(quote *models.Quote) {
	pairs := [][]string{
		{"Quote ID", quote.QuoteID},
		{"RFQ ID", quote.RFQID},
		{"Ticker", quote.Ticker},
		{"Side", formatSide(quote.Side)},
		{"Price", ui.FormatPrice(quote.Price)},
		{"Quantity", fmt.Sprintf("%d", quote.Quantity)},
		{"Status", formatStatus(quote.Status)},
		{"Created", quote.CreatedTime.Format("2006-01-02 15:04:05")},
		{"Expires", quote.ExpiresTime.Format("2006-01-02 15:04:05")},
	}
	ui.RenderKeyValue(pairs)
}

func renderQuoteDetailPlain(quote *models.Quote) {
	fmt.Printf("quote_id=%s rfq_id=%s ticker=%s side=%s price=%d qty=%d status=%s\n",
		quote.QuoteID, quote.RFQID, quote.Ticker, quote.Side,
		quote.Price, quote.Quantity, quote.Status)
}

// Helper functions

func formatSide(side string) string {
	switch strings.ToLower(side) {
	case "yes":
		return ui.SuccessStyle.Render("YES")
	case "no":
		return ui.ErrorStyle.Render("NO")
	default:
		return side
	}
}

func formatStatus(status string) string {
	switch strings.ToLower(status) {
	case "open", "active":
		return ui.StatusOpenStyle.Render(strings.ToUpper(status))
	case "closed", "expired", "cancelled":
		return ui.StatusClosedStyle.Render(strings.ToUpper(status))
	case "accepted", "confirmed":
		return ui.StatusActiveStyle.Render(strings.ToUpper(status))
	default:
		return status
	}
}
