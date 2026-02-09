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

var marketsCmd = &cobra.Command{
	Use:   "markets",
	Short: "Manage and view markets",
	Long:  `Commands for listing, viewing, and analyzing prediction markets.`,
}

var marketsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List markets",
	Long:  `List markets with optional filtering by status and series.`,
	RunE:  runMarketsList,
}

var marketsGetCmd = &cobra.Command{
	Use:   "get <ticker>",
	Short: "Get market details",
	Long:  `Get detailed information about a specific market.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMarketsGet,
}

var marketsOrderbookCmd = &cobra.Command{
	Use:   "orderbook <ticker>",
	Short: "Get market orderbook",
	Long:  `Get the orderbook for a specific market with visual display.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMarketsOrderbook,
}

var marketsTradesCmd = &cobra.Command{
	Use:   "trades <ticker>",
	Short: "Get market trades",
	Long:  `Get recent trades for a specific market.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMarketsTrades,
}

var marketsCandlesticksCmd = &cobra.Command{
	Use:   "candlesticks <ticker>",
	Short: "Get market candlesticks",
	Long: `Get candlestick (OHLCV) data for a specific market.

Supported periods: 1m, 5m, 15m, 1h, 4h, 1d`,
	Args: cobra.ExactArgs(1),
	RunE: runMarketsCandlesticks,
}

var seriesCmd = &cobra.Command{
	Use:   "series",
	Short: "Manage and view series",
	Long:  `Commands for listing and viewing market series.`,
}

var seriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List series",
	Long:  `List market series with optional category filtering.`,
	RunE:  runSeriesList,
}

var seriesGetCmd = &cobra.Command{
	Use:   "get <ticker>",
	Short: "Get series details",
	Long:  `Get detailed information about a specific series.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSeriesGet,
}

// Command flags
var (
	marketStatus   string
	marketLimit    int
	seriesTicker   string
	tradesLimit    int
	candlePeriod   string
	seriesCategory string
	seriesLimit    int
)

func init() {
	marketsListCmd.Flags().StringVar(&marketStatus, "status", "", "filter by status (open, closed, settled)")
	marketsListCmd.Flags().IntVar(&marketLimit, "limit", 50, "maximum number of markets to return")
	marketsListCmd.Flags().StringVar(&seriesTicker, "series", "", "filter by series ticker")

	marketsTradesCmd.Flags().IntVar(&tradesLimit, "limit", 100, "maximum number of trades to return")

	marketsCandlesticksCmd.Flags().StringVar(&candlePeriod, "period", "1h", "candlestick period (1m, 5m, 15m, 1h, 4h, 1d)")

	seriesListCmd.Flags().StringVar(&seriesCategory, "category", "", "filter by category")
	seriesListCmd.Flags().IntVar(&seriesLimit, "limit", 50, "maximum number of series to return")

	seriesCmd.AddCommand(seriesListCmd)
	seriesCmd.AddCommand(seriesGetCmd)

	marketsCmd.AddCommand(marketsListCmd)
	marketsCmd.AddCommand(marketsGetCmd)
	marketsCmd.AddCommand(marketsOrderbookCmd)
	marketsCmd.AddCommand(marketsTradesCmd)
	marketsCmd.AddCommand(marketsCandlesticksCmd)
	marketsCmd.AddCommand(seriesCmd)

	rootCmd.AddCommand(marketsCmd)
}

func runMarketsList(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := api.ListMarketsParams{
		Status:       marketStatus,
		SeriesTicker: seriesTicker,
		Limit:        marketLimit,
	}

	result, err := client.ListMarkets(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list markets: %w", err)
	}

	return outputMarketsList(result.Markets)
}

func outputMarketsList(markets []models.Market) error {
	format := GetOutputFormat()

	tableFunc := func() {
		headers := []string{"Ticker", "Title", "Status", "Yes Bid", "Yes Ask", "Volume"}
		var rows [][]string

		for _, m := range markets {
			title := truncateMarketString(m.Title, 50)
			rows = append(rows, []string{
				m.Ticker,
				title,
				formatMarketStatus(m.Status),
				formatCents(m.YesBid),
				formatCents(m.YesAsk),
				fmt.Sprintf("%d", m.Volume),
			})
		}

		ui.RenderTable(headers, rows)
	}

	plainFunc := func() {
		for _, m := range markets {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%d\n",
				m.Ticker,
				m.Title,
				m.Status,
				formatCents(m.YesBid),
				formatCents(m.YesAsk),
				m.Volume,
			)
		}
	}

	return ui.Output(format, tableFunc, markets, plainFunc)
}

func runMarketsGet(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	market, err := client.GetMarket(ctx, ticker)
	if err != nil {
		return fmt.Errorf("failed to get market: %w", err)
	}

	return outputMarketDetails(market)
}

func outputMarketDetails(market *models.Market) error {
	format := GetOutputFormat()

	tableFunc := func() {
		pairs := [][]string{
			{"Ticker", market.Ticker},
			{"Title", market.Title},
			{"Subtitle", market.Subtitle},
			{"Status", formatMarketStatus(market.Status)},
			{"Category", market.Category},
			{"Yes Bid", formatCents(market.YesBid)},
			{"Yes Ask", formatCents(market.YesAsk)},
			{"No Bid", formatCents(market.NoBid)},
			{"No Ask", formatCents(market.NoAsk)},
			{"Last Price", formatCents(market.LastPrice)},
			{"Volume", fmt.Sprintf("%d", market.Volume)},
			{"Volume 24h", fmt.Sprintf("%d", market.Volume24H)},
			{"Open Interest", fmt.Sprintf("%d", market.OpenInterest)},
			{"Open Time", formatMarketTime(market.OpenTime)},
			{"Close Time", formatMarketTime(market.CloseTime)},
			{"Expiration", formatMarketTime(market.ExpirationTime)},
		}

		if market.Result != "" {
			pairs = append(pairs, []string{"Result", market.Result})
		}

		ui.RenderKeyValue(pairs)
	}

	plainFunc := func() {
		fmt.Printf("Ticker: %s\n", market.Ticker)
		fmt.Printf("Title: %s\n", market.Title)
		fmt.Printf("Status: %s\n", market.Status)
		fmt.Printf("Yes Bid/Ask: %s / %s\n", formatCents(market.YesBid), formatCents(market.YesAsk))
		fmt.Printf("No Bid/Ask: %s / %s\n", formatCents(market.NoBid), formatCents(market.NoAsk))
		fmt.Printf("Last Price: %s\n", formatCents(market.LastPrice))
		fmt.Printf("Volume: %d\n", market.Volume)
	}

	return ui.Output(format, tableFunc, market, plainFunc)
}

func runMarketsOrderbook(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	orderbook, err := client.GetOrderbook(ctx, ticker)
	if err != nil {
		return fmt.Errorf("failed to get orderbook: %w", err)
	}

	return outputOrderbook(orderbook)
}

func outputOrderbook(ob *models.Orderbook) error {
	format := GetOutputFormat()

	tableFunc := func() {
		fmt.Printf("\n%s Orderbook for %s\n\n", ui.TitleStyle.Render("YES"), ob.Ticker)

		// YES side - Bids on left, Asks on right
		fmt.Println(ui.HeaderStyle.Render("         BIDS                    ASKS"))
		fmt.Println(ui.MutedStyle.Render("   Qty    Price           Price    Qty"))
		fmt.Println(strings.Repeat("-", 50))

		maxRows := maxInt(len(ob.YesBids), len(ob.YesAsks))
		for i := 0; i < maxRows; i++ {
			bidStr := "                    "
			askStr := "                    "

			if i < len(ob.YesBids) {
				bid := ob.YesBids[i]
				bidStr = fmt.Sprintf("%5d    %s", bid.Quantity, formatCents(bid.Price))
				bidStr = ui.PriceUpStyle.Render(bidStr)
			}

			if i < len(ob.YesAsks) {
				ask := ob.YesAsks[i]
				askStr = fmt.Sprintf("%s    %5d", formatCents(ask.Price), ask.Quantity)
				askStr = ui.PriceDownStyle.Render(askStr)
			}

			fmt.Printf("%s       %s\n", bidStr, askStr)
		}

		fmt.Println()
	}

	plainFunc := func() {
		fmt.Printf("Ticker: %s\n", ob.Ticker)
		fmt.Println("YES BIDS:")
		for _, bid := range ob.YesBids {
			fmt.Printf("  %s x %d\n", formatCents(bid.Price), bid.Quantity)
		}
		fmt.Println("YES ASKS:")
		for _, ask := range ob.YesAsks {
			fmt.Printf("  %s x %d\n", formatCents(ask.Price), ask.Quantity)
		}
	}

	return ui.Output(format, tableFunc, ob, plainFunc)
}

func runMarketsTrades(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := api.GetTradesParams{
		Ticker: ticker,
		Limit:  tradesLimit,
	}

	result, err := client.GetTrades(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to get trades: %w", err)
	}

	return outputTrades(result.Trades)
}

func outputTrades(trades []models.Trade) error {
	format := GetOutputFormat()

	tableFunc := func() {
		headers := []string{"Time", "Price", "Quantity", "Side"}
		var rows [][]string

		for _, t := range trades {
			side := formatTradeSide(t.TakerSide)
			rows = append(rows, []string{
				formatMarketTime(t.CreatedTime),
				formatCents(t.Price),
				fmt.Sprintf("%d", t.Count),
				side,
			})
		}

		ui.RenderTable(headers, rows)
	}

	plainFunc := func() {
		for _, t := range trades {
			fmt.Printf("%s\t%s\t%d\t%s\n",
				t.CreatedTime.Format(time.RFC3339),
				formatCents(t.Price),
				t.Count,
				t.TakerSide,
			)
		}
	}

	return ui.Output(format, tableFunc, trades, plainFunc)
}

func runMarketsCandlesticks(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := api.GetCandlesticksParams{
		Ticker: ticker,
		Period: candlePeriod,
	}

	result, err := client.GetCandlesticks(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to get candlesticks: %w", err)
	}

	return outputCandlesticks(result.Candlesticks)
}

func outputCandlesticks(candles []models.Candlestick) error {
	format := GetOutputFormat()

	tableFunc := func() {
		headers := []string{"Time", "Open", "High", "Low", "Close", "Volume"}
		var rows [][]string

		for _, c := range candles {
			rows = append(rows, []string{
				formatMarketTime(c.PeriodEnd),
				formatCents(c.Open),
				formatCents(c.High),
				formatCents(c.Low),
				formatCents(c.Close),
				fmt.Sprintf("%d", c.Volume),
			})
		}

		ui.RenderTable(headers, rows)
	}

	plainFunc := func() {
		for _, c := range candles {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%d\n",
				c.PeriodEnd.Format(time.RFC3339),
				formatCents(c.Open),
				formatCents(c.High),
				formatCents(c.Low),
				formatCents(c.Close),
				c.Volume,
			)
		}
	}

	return ui.Output(format, tableFunc, candles, plainFunc)
}

func runSeriesList(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := api.ListSeriesParams{
		Category: seriesCategory,
		Limit:    seriesLimit,
	}

	result, err := client.ListSeries(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to list series: %w", err)
	}

	return outputSeriesList(result.Series)
}

func outputSeriesList(series []models.Series) error {
	format := GetOutputFormat()

	tableFunc := func() {
		headers := []string{"Ticker", "Title", "Category", "Frequency"}
		var rows [][]string

		for _, s := range series {
			title := truncateMarketString(s.Title, 50)
			rows = append(rows, []string{
				s.Ticker,
				title,
				s.Category,
				s.Frequency,
			})
		}

		ui.RenderTable(headers, rows)
	}

	plainFunc := func() {
		for _, s := range series {
			fmt.Printf("%s\t%s\t%s\t%s\n",
				s.Ticker,
				s.Title,
				s.Category,
				s.Frequency,
			)
		}
	}

	return ui.Output(format, tableFunc, series, plainFunc)
}

func runSeriesGet(cmd *cobra.Command, args []string) error {
	ticker := args[0]

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	series, err := client.GetSeries(ctx, ticker)
	if err != nil {
		return fmt.Errorf("failed to get series: %w", err)
	}

	return outputSeriesDetails(series)
}

func outputSeriesDetails(series *models.Series) error {
	format := GetOutputFormat()

	tableFunc := func() {
		pairs := [][]string{
			{"Ticker", series.Ticker},
			{"Title", series.Title},
			{"Category", series.Category},
			{"Frequency", series.Frequency},
			{"Tags", strings.Join(series.Tags, ", ")},
		}

		ui.RenderKeyValue(pairs)
	}

	plainFunc := func() {
		fmt.Printf("Ticker: %s\n", series.Ticker)
		fmt.Printf("Title: %s\n", series.Title)
		fmt.Printf("Category: %s\n", series.Category)
		fmt.Printf("Frequency: %s\n", series.Frequency)
		fmt.Printf("Tags: %s\n", strings.Join(series.Tags, ", "))
	}

	return ui.Output(format, tableFunc, series, plainFunc)
}

// Helper functions for markets commands



func formatTradeSide(side string) string {
	switch strings.ToLower(side) {
	case "yes":
		return ui.PriceUpStyle.Render("YES")
	case "no":
		return ui.PriceDownStyle.Render("NO")
	default:
		return side
	}
}

func formatMarketTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

func truncateMarketString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
