package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/internal/websocket"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var (
	watchMarketFlag string
)

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.AddCommand(watchTickerCmd)
	watchCmd.AddCommand(watchOrderbookCmd)
	watchCmd.AddCommand(watchTradesCmd)
	watchCmd.AddCommand(watchOrdersCmd)
	watchCmd.AddCommand(watchFillsCmd)
	watchCmd.AddCommand(watchPositionsCmd)

	watchTradesCmd.Flags().StringVar(&watchMarketFlag, "market", "", "filter trades by market ticker")
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch live market data and account updates",
	Long: `Stream real-time data from Kalshi via WebSocket.

All watch commands require authentication (API credentials).
Press Ctrl+C to stop watching.

Available streams:
  ticker      Live price updates for a market (requires <market-ticker>)
  orderbook   Orderbook delta updates for a market (requires <market-ticker>)
  trades      Public trades feed (optional --market filter)
  orders      Your order status changes
  fills       Your fill notifications
  positions   Your position changes`,
	Example: `  kalshi-cli watch ticker INXD-25FEB07-B5523.99
  kalshi-cli watch orderbook INXD-25FEB07-B5523.99
  kalshi-cli watch trades --market INXD-25FEB07-B5523.99
  kalshi-cli watch orders
  kalshi-cli watch fills --json
  kalshi-cli watch positions`,
}

var watchTickerCmd = &cobra.Command{
	Use:   "ticker <market-ticker>",
	Short: "Watch live price updates for a market",
	Long: `Stream real-time price updates for a specific market.

Output includes bid/ask prices, volume, and open interest.
Use 'kalshi-cli markets list' to find available market tickers.`,
	Example: `  kalshi-cli watch ticker INXD-25FEB07-B5523.99
  kalshi-cli watch ticker INXD-25FEB07-B5523.99 --json
  kalshi-cli watch ticker INXD-25FEB07-B5523.99 --plain`,
	Args: cobra.ExactArgs(1),
	RunE: runWatchTicker,
}

var watchOrderbookCmd = &cobra.Command{
	Use:   "orderbook <market-ticker>",
	Short: "Watch live orderbook updates for a market",
	Long: `Stream real-time orderbook delta updates for a specific market.

Shows best bid/ask, depth, and orderbook changes as they occur.
Use 'kalshi-cli markets list' to find available market tickers.`,
	Example: `  kalshi-cli watch orderbook INXD-25FEB07-B5523.99
  kalshi-cli watch orderbook INXD-25FEB07-B5523.99 --json`,
	Args: cobra.ExactArgs(1),
	RunE: runWatchOrderbook,
}

var watchTradesCmd = &cobra.Command{
	Use:   "trades",
	Short: "Watch public trades feed",
	Long: `Stream real-time public trades across all markets.

Optionally filter to a single market using the --market flag.`,
	Example: `  kalshi-cli watch trades
  kalshi-cli watch trades --market INXD-25FEB07-B5523.99
  kalshi-cli watch trades --json`,
	RunE: runWatchTrades,
}

var watchOrdersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Watch your order updates",
	Long: `Stream real-time updates for your orders.

Shows order status changes, fills, and cancellations as they happen.`,
	Example: `  kalshi-cli watch orders
  kalshi-cli watch orders --json`,
	RunE: runWatchOrders,
}

var watchFillsCmd = &cobra.Command{
	Use:   "fills",
	Short: "Watch your fill notifications",
	Long: `Stream real-time fill notifications for your orders.

Shows each individual fill as it occurs, including price, count, and taker/maker status.`,
	Example: `  kalshi-cli watch fills
  kalshi-cli watch fills --json`,
	RunE: runWatchFills,
}

var watchPositionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "Watch your position changes",
	Long: `Stream real-time position updates.

Shows changes to your positions including realized PnL, exposure, and total cost.`,
	Example: `  kalshi-cli watch positions
  kalshi-cli watch positions --json`,
	RunE: runWatchPositions,
}

func runWatchTicker(_ *cobra.Command, args []string) error {
	ticker := args[0]
	params := map[string]string{"market_tickers": ticker}
	return runWatch(websocket.ChannelMarketTicker, params)
}

func runWatchOrderbook(_ *cobra.Command, args []string) error {
	ticker := args[0]
	params := map[string]string{"market_tickers": ticker}
	return runWatch(websocket.ChannelOrderbook, params)
}

func runWatchTrades(_ *cobra.Command, _ []string) error {
	params := make(map[string]string)
	if watchMarketFlag != "" {
		params["market_tickers"] = watchMarketFlag
	}
	return runWatch(websocket.ChannelPublicTrades, params)
}

func runWatchOrders(_ *cobra.Command, _ []string) error {
	return runWatch(websocket.ChannelUserOrders, nil)
}

func runWatchFills(_ *cobra.Command, _ []string) error {
	return runWatch(websocket.ChannelUserFills, nil)
}

func runWatchPositions(_ *cobra.Command, _ []string) error {
	return runWatch(websocket.ChannelMarketPositions, nil)
}

func runWatch(channel websocket.Channel, params map[string]string) error {
	return runWatchMultiple([]websocket.Channel{channel}, params)
}

func runWatchMultiple(channels []websocket.Channel, params map[string]string) error {
	cfg := GetConfig()

	opts, err := buildClientOptions(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		if IsVerbose() {
			fmt.Fprintln(os.Stderr, "\nShutting down...")
		}
		cancel()
	}()

	client := websocket.NewClient(opts)

	client.OnReconnect(func() {
		if IsVerbose() {
			fmt.Fprintln(os.Stderr, "Reconnected")
		}
	})

	client.OnError(func(err error) {
		if IsVerbose() {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	})

	registerHandlers(client, channels)

	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	if IsVerbose() {
		fmt.Fprintf(os.Stderr, "Connected to %s\n", cfg.Environment())
	}

	for _, ch := range channels {
		if err := client.Subscribe(ctx, ch, params); err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", ch, err)
		}
	}

	if IsVerbose() {
		channelNames := make([]string, len(channels))
		for i, ch := range channels {
			channelNames[i] = string(ch)
		}
		fmt.Fprintf(os.Stderr, "Subscribed to: %s\n", strings.Join(channelNames, ", "))
	}

	<-ctx.Done()
	return nil
}

func buildClientOptions(cfg *config.Config) (websocket.ClientOptions, error) {
	opts := websocket.ClientOptions{
		URL: cfg.WebSocketURL(),
	}

	signer, err := getSigner(cfg)
	if err != nil {
		return opts, fmt.Errorf("authentication required for WebSocket connection: %w", err)
	}

	timestamp := time.Now().UTC()
	signature, err := signer.Sign(timestamp, "GET", "/trade-api/ws/v2")
	if err != nil {
		return opts, fmt.Errorf("failed to sign request: %w", err)
	}

	opts.APIKeyID = signer.APIKeyID()
	opts.Signature = signature
	opts.Timestamp = api.TimestampHeader(timestamp)

	return opts, nil
}

func registerHandlers(client *websocket.Client, channels []websocket.Channel) {
	outputFormat := GetOutputFormat()

	for _, ch := range channels {
		switch ch {
		case websocket.ChannelMarketTicker:
			client.RegisterHandler(ch, &tickerHandler{format: outputFormat})
		case websocket.ChannelMarketTickerV2:
			client.RegisterHandler(ch, &tickerV2Handler{format: outputFormat})
		case websocket.ChannelOrderbook:
			client.RegisterHandler(ch, &orderbookHandler{format: outputFormat})
		case websocket.ChannelPublicTrades:
			client.RegisterHandler(ch, &tradesHandler{format: outputFormat, filterTicker: watchMarketFlag})
		case websocket.ChannelUserOrders:
			client.RegisterHandler(ch, &ordersHandler{format: outputFormat})
		case websocket.ChannelUserFills:
			client.RegisterHandler(ch, &fillsHandler{format: outputFormat})
		case websocket.ChannelMarketPositions:
			client.RegisterHandler(ch, &positionsHandler{format: outputFormat})
		case websocket.ChannelMarketLifecycle:
			client.RegisterHandler(ch, &lifecycleHandler{format: outputFormat})
		case websocket.ChannelOrderGroupUpdates:
			client.RegisterHandler(ch, &orderGroupHandler{format: outputFormat})
		case websocket.ChannelCommunications:
			client.RegisterHandler(ch, &communicationsHandler{format: outputFormat})
		}
	}
}

func requiresAuth(channels []websocket.Channel) bool {
	for _, ch := range channels {
		if websocket.ChannelRequiresAuth(ch) {
			return true
		}
	}
	return false
}

func getSigner(_ *config.Config) (*api.Signer, error) {
	// Try config file first (fast, no GUI prompts)
	apiKeyID := viper.GetString("api_key_id")
	privateKeyPath := viper.GetString("private_key_path")

	// Check env vars
	if apiKeyID == "" {
		apiKeyID = os.Getenv("KALSHI_API_KEY_ID")
	}
	if privateKeyPath == "" {
		privateKeyPath = os.Getenv("KALSHI_PRIVATE_KEY_FILE")
	}

	if apiKeyID != "" && privateKeyPath != "" {
		pemData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		return api.NewSignerFromPEM(apiKeyID, string(pemData))
	}

	// Check KALSHI_PRIVATE_KEY env var (PEM content directly)
	privateKeyPEM := os.Getenv("KALSHI_PRIVATE_KEY")
	if apiKeyID != "" && privateKeyPEM != "" {
		return api.NewSignerFromPEM(apiKeyID, privateKeyPEM)
	}

	// Last resort: keyring
	keyringStore, err := config.NewKeyringStore()
	if err != nil {
		return nil, fmt.Errorf("failed to access keyring: %w", err)
	}
	creds, err := keyringStore.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}
	if creds == nil {
		return nil, fmt.Errorf("no credentials configured")
	}
	return api.NewSignerFromPEM(creds.APIKeyID, creds.PrivateKey)
}

func formatTimestamp() string {
	return time.Now().Format("15:04:05")
}


func formatVolume(vol int) string {
	if vol >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(vol)/1000000)
	}
	if vol >= 1000 {
		return fmt.Sprintf("%.1fK", float64(vol)/1000)
	}
	return fmt.Sprintf("%d", vol)
}

// tickerHandler handles market ticker messages
type tickerHandler struct {
	format ui.OutputFormat
}

func (h *tickerHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.TickerData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse ticker data: %w", err)
	}

	return h.output(data)
}

func (h *tickerHandler) output(data websocket.TickerData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s %s yes=%d no=%d vol=%d oi=%d\n",
			formatTimestamp(), data.Ticker, data.YesPrice, data.NoPrice, data.Volume, data.OpenInterest)
	default:
		spread := ""
		if data.YesBid > 0 && data.YesAsk > 0 {
			spread = fmt.Sprintf("Yes %s / %s", formatCents(data.YesBid), formatCents(data.YesAsk))
		} else {
			spread = fmt.Sprintf("Yes %s", formatCents(data.YesPrice))
		}
		fmt.Printf("[%s] %s: %s | Vol: %s\n",
			formatTimestamp(), data.Ticker, spread, formatVolume(data.Volume))
	}
	return nil
}

// orderbookHandler handles orderbook messages
type orderbookHandler struct {
	format ui.OutputFormat
}

func (h *orderbookHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.OrderbookData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse orderbook data: %w", err)
	}

	return h.output(data)
}

func (h *orderbookHandler) output(data websocket.OrderbookData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		bids := formatLevels(data.YesBids, 3)
		asks := formatLevels(data.YesAsks, 3)
		fmt.Printf("%s %s bids=[%s] asks=[%s]\n",
			formatTimestamp(), data.Ticker, bids, asks)
	default:
		bestBid := "-"
		bestAsk := "-"
		bidDepth := 0
		askDepth := 0

		if len(data.YesBids) > 0 {
			bestBid = formatCents(data.YesBids[0].Price)
			for _, l := range data.YesBids {
				bidDepth += l.Quantity
			}
		}
		if len(data.YesAsks) > 0 {
			bestAsk = formatCents(data.YesAsks[0].Price)
			for _, l := range data.YesAsks {
				askDepth += l.Quantity
			}
		}

		fmt.Printf("[%s] %s: Bid %s (%d) | Ask %s (%d)\n",
			formatTimestamp(), data.Ticker, bestBid, bidDepth, bestAsk, askDepth)
	}
	return nil
}

func formatLevels(levels []websocket.OrderbookLevel, max int) string {
	if len(levels) == 0 {
		return "-"
	}

	count := len(levels)
	if count > max {
		count = max
	}

	parts := make([]string, count)
	for i := 0; i < count; i++ {
		parts[i] = fmt.Sprintf("%d@%d", levels[i].Quantity, levels[i].Price)
	}
	return strings.Join(parts, ",")
}

// tradesHandler handles public trades messages
type tradesHandler struct {
	format       ui.OutputFormat
	filterTicker string
}

func (h *tradesHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.TradeData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse trade data: %w", err)
	}

	if h.filterTicker != "" && data.Ticker != h.filterTicker {
		return nil
	}

	return h.output(data)
}

func (h *tradesHandler) output(data websocket.TradeData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s %s %s price=%d count=%d\n",
			formatTimestamp(), data.Ticker, data.TakerSide, data.Price, data.Count)
	default:
		side := data.TakerSide
		if side == "yes" {
			side = ui.PriceUpStyle.Render("BUY")
		} else {
			side = ui.PriceDownStyle.Render("SELL")
		}
		fmt.Printf("[%s] %s: %s %d @ %s\n",
			formatTimestamp(), data.Ticker, side, data.Count, formatCents(data.Price))
	}
	return nil
}

// ordersHandler handles user order messages
type ordersHandler struct {
	format ui.OutputFormat
}

func (h *ordersHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.OrderUpdateData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse order data: %w", err)
	}

	return h.output(data)
}

func (h *ordersHandler) output(data websocket.OrderUpdateData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		orderID := truncateID(data.OrderID, 8)
		fmt.Printf("%s order=%s ticker=%s status=%s side=%s action=%s qty=%d/%d\n",
			formatTimestamp(), orderID, data.Ticker, data.Status,
			data.Side, data.Action, data.FilledQuantity, data.InitialQuantity)
	default:
		orderID := truncateID(data.OrderID, 8)
		status := formatOrderStatus(data.Status)
		price := data.YesPrice
		if data.Side == "no" {
			price = data.NoPrice
		}
		fmt.Printf("[%s] Order %s: %s %s %s @ %s | %s (%d/%d filled)\n",
			formatTimestamp(), orderID, strings.ToUpper(data.Action),
			data.Ticker, strings.ToUpper(data.Side), formatCents(price),
			status, data.FilledQuantity, data.InitialQuantity)
	}
	return nil
}

func formatOrderStatus(status string) string {
	switch status {
	case string(models.OrderStatusResting):
		return ui.StatusActiveStyle.Render("RESTING")
	case string(models.OrderStatusExecuted):
		return ui.SuccessStyle.Render("EXECUTED")
	case string(models.OrderStatusCanceled):
		return ui.MutedStyle.Render("CANCELED")
	case string(models.OrderStatusPending):
		return ui.WarningStyle.Render("PENDING")
	default:
		return status
	}
}

// fillsHandler handles user fill messages
type fillsHandler struct {
	format ui.OutputFormat
}

func (h *fillsHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.FillData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse fill data: %w", err)
	}

	return h.output(data)
}

func (h *fillsHandler) output(data websocket.FillData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fillID := truncateID(data.FillID, 8)
		orderID := truncateID(data.OrderID, 8)
		fmt.Printf("%s fill=%s order=%s ticker=%s side=%s action=%s price=%d count=%d taker=%v\n",
			formatTimestamp(), fillID, orderID, data.Ticker,
			data.Side, data.Action, data.YesPrice, data.Count, data.IsTaker)
	default:
		takerMaker := "maker"
		if data.IsTaker {
			takerMaker = "taker"
		}
		price := data.YesPrice
		if data.Side == "no" {
			price = data.NoPrice
		}
		fmt.Printf("[%s] FILL: %s %s %s @ %s x%d (%s)\n",
			formatTimestamp(), strings.ToUpper(data.Action), data.Ticker,
			strings.ToUpper(data.Side), formatCents(price), data.Count, takerMaker)
	}
	return nil
}

func truncateID(id string, length int) string {
	if len(id) <= length {
		return id
	}
	return id[:length]
}

func printJSONLine(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// tickerV2Handler handles market_ticker_v2 incremental delta messages
type tickerV2Handler struct {
	format ui.OutputFormat
}

func (h *tickerV2Handler) HandleMessage(msg websocket.Message) error {
	var data websocket.TickerV2Data
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse ticker v2 data: %w", err)
	}

	return h.output(data)
}

func (h *tickerV2Handler) output(data websocket.TickerV2Data) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s %s delta_type=%s yes=%d no=%d delta=%d\n",
			formatTimestamp(), data.Ticker, data.DeltaType, data.YesPrice, data.NoPrice, data.Delta)
	default:
		fmt.Printf("[%s] %s: %s (delta: %+d) Yes %s / No %s\n",
			formatTimestamp(), data.Ticker, data.DeltaType, data.Delta,
			formatCents(data.YesPrice), formatCents(data.NoPrice))
	}
	return nil
}

// positionsHandler handles market_positions messages
type positionsHandler struct {
	format ui.OutputFormat
}

func (h *positionsHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.PositionData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse position data: %w", err)
	}

	return h.output(data)
}

func (h *positionsHandler) output(data websocket.PositionData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s ticker=%s position=%d cost=%d pnl=%d exposure=%d\n",
			formatTimestamp(), data.Ticker, data.Position, data.TotalCost, data.RealizedPnl, data.Exposure)
	default:
		pnlStyle := ui.MutedStyle
		if data.RealizedPnl > 0 {
			pnlStyle = ui.PriceUpStyle
		} else if data.RealizedPnl < 0 {
			pnlStyle = ui.PriceDownStyle
		}
		fmt.Printf("[%s] %s: Position %d | Cost %s | PnL %s | Exposure %s\n",
			formatTimestamp(), data.Ticker, data.Position,
			formatCents(data.TotalCost), pnlStyle.Render(formatCents(data.RealizedPnl)),
			formatCents(data.Exposure))
	}
	return nil
}

// lifecycleHandler handles market_lifecycle messages
type lifecycleHandler struct {
	format ui.OutputFormat
}

func (h *lifecycleHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.MarketLifecycleData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse lifecycle data: %w", err)
	}

	return h.output(data)
}

func (h *lifecycleHandler) output(data websocket.MarketLifecycleData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s ticker=%s status=%s old_status=%s\n",
			formatTimestamp(), data.Ticker, data.Status, data.OldStatus)
	default:
		fmt.Printf("[%s] %s: %s -> %s\n",
			formatTimestamp(), data.Ticker, data.OldStatus, data.Status)
	}
	return nil
}

// orderGroupHandler handles order_group_updates messages
type orderGroupHandler struct {
	format ui.OutputFormat
}

func (h *orderGroupHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.OrderGroupUpdateData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse order group data: %w", err)
	}

	return h.output(data)
}

func (h *orderGroupHandler) output(data websocket.OrderGroupUpdateData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s order_group=%s status=%s total=%d filled=%d\n",
			formatTimestamp(), data.OrderGroupID, data.Status, data.TotalOrders, data.FilledOrders)
	default:
		fmt.Printf("[%s] Order Group %s: %s (%d/%d filled)\n",
			formatTimestamp(), truncateID(data.OrderGroupID, 8), data.Status,
			data.FilledOrders, data.TotalOrders)
	}
	return nil
}

// communicationsHandler handles communications (RFQ/quote) messages
type communicationsHandler struct {
	format ui.OutputFormat
}

func (h *communicationsHandler) HandleMessage(msg websocket.Message) error {
	var data websocket.CommunicationData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to parse communication data: %w", err)
	}

	return h.output(data)
}

func (h *communicationsHandler) output(data websocket.CommunicationData) error {
	switch h.format {
	case ui.FormatJSON:
		return printJSONLine(data)
	case ui.FormatPlain:
		fmt.Printf("%s type=%s ticker=%s qty=%d price=%d side=%s\n",
			formatTimestamp(), data.Type, data.Ticker, data.Quantity, data.Price, data.Side)
	default:
		fmt.Printf("[%s] %s: %s %s %d @ %s\n",
			formatTimestamp(), strings.ToUpper(data.Type), data.Ticker,
			strings.ToUpper(data.Side), data.Quantity, formatCents(data.Price))
	}
	return nil
}
