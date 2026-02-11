package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "Manage trading orders",
	Long: `Manage trading orders on the Kalshi exchange.

Commands for listing, creating, canceling, and amending orders.`,
	Example: `  kalshi-cli orders list --status resting
  kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side yes --qty 10 --price 50
  kalshi-cli orders cancel ORDER_ID`,
}

var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List orders",
	Long:  `List orders with optional filters for status and market ticker.`,
	Example: `  kalshi-cli orders list
  kalshi-cli orders list --status resting
  kalshi-cli orders list --market INXD-25FEB07-B5523.99 --json`,
	RunE: runOrdersList,
}

var ordersGetCmd = &cobra.Command{
	Use:     "get <order-id>",
	Short:   "Get order details",
	Long:    `Get detailed information about a specific order.`,
	Example: `  kalshi-cli orders get abc123-def456-ghi789`,
	Args:    cobra.ExactArgs(1),
	RunE:    runOrdersGet,
}

var ordersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new order",
	Long: `Create a new limit order on a market.

The order preview will be shown before submission. You must confirm
unless the --yes flag is set.

Price must be between 1-99 cents.`,
	Example: `  kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side yes --qty 10 --price 50
  kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side no --qty 5 --price 30 --action sell
  kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side yes --qty 10 --price 50 --yes`,
	RunE: runOrdersCreate,
}

var ordersCancelCmd = &cobra.Command{
	Use:   "cancel <order-id>",
	Short: "Cancel an order",
	Long:  `Cancel a resting order by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOrdersCancel,
}

var ordersCancelAllCmd = &cobra.Command{
	Use:   "cancel-all",
	Short: "Cancel all orders",
	Long:  `Cancel all resting orders, optionally filtered by market ticker.`,
	RunE:  runOrdersCancelAll,
}

var ordersAmendCmd = &cobra.Command{
	Use:   "amend <order-id>",
	Short: "Amend an existing order",
	Long: `Amend an existing order's quantity and/or price.

At least one of --qty or --price must be specified.`,
	Args: cobra.ExactArgs(1),
	RunE: runOrdersAmend,
}

var ordersBatchCreateCmd = &cobra.Command{
	Use:   "batch-create",
	Short: "Create multiple orders from a file",
	Long: `Create multiple orders from a JSON file.

The JSON file should contain an array of order objects with fields:
- ticker: Market ticker (required)
- side: "yes" or "no" (required)
- action: "buy" or "sell" (required)
- type: "limit" or "market" (required)
- count: Quantity (required)
- yes_price: Price in cents for yes side (optional)
- no_price: Price in cents for no side (optional)`,
	RunE: runOrdersBatchCreate,
}

var ordersQueueCmd = &cobra.Command{
	Use:   "queue <order-id>",
	Short: "Get queue position for an order",
	Long:  `Get the queue position for a resting order.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOrdersQueue,
}

// Flags
var (
	orderStatusFilter   string
	orderMarketFilter   string
	orderCreateMarket   string
	orderCancelAllMarket string
	orderSide           string
	orderCreateQty      int
	orderCreatePrice    int
	orderAmendQty       int
	orderAmendPrice     int
	orderAction         string
	orderType           string
	batchFile           string
)

func init() {
	rootCmd.AddCommand(ordersCmd)

	ordersCmd.AddCommand(ordersListCmd)
	ordersCmd.AddCommand(ordersGetCmd)
	ordersCmd.AddCommand(ordersCreateCmd)
	ordersCmd.AddCommand(ordersCancelCmd)
	ordersCmd.AddCommand(ordersCancelAllCmd)
	ordersCmd.AddCommand(ordersAmendCmd)
	ordersCmd.AddCommand(ordersBatchCreateCmd)
	ordersCmd.AddCommand(ordersQueueCmd)

	// List flags
	ordersListCmd.Flags().StringVar(&orderStatusFilter, "status", "", "filter by status (resting, canceled, executed, pending)")
	ordersListCmd.Flags().StringVar(&orderMarketFilter, "market", "", "filter by market ticker")

	// Create flags
	ordersCreateCmd.Flags().StringVar(&orderCreateMarket, "market", "", "market ticker (required)")
	ordersCreateCmd.Flags().StringVar(&orderSide, "side", "", "order side: yes or no (required)")
	ordersCreateCmd.Flags().IntVar(&orderCreateQty, "qty", 0, "quantity (required)")
	ordersCreateCmd.Flags().IntVar(&orderCreatePrice, "price", 0, "price in cents 1-99 (required)")
	ordersCreateCmd.Flags().StringVar(&orderAction, "action", "buy", "order action: buy or sell (default: buy)")
	ordersCreateCmd.Flags().StringVar(&orderType, "type", "limit", "order type: limit or market (default: limit)")
	ordersCreateCmd.MarkFlagRequired("market")
	ordersCreateCmd.MarkFlagRequired("side")
	ordersCreateCmd.MarkFlagRequired("qty")
	ordersCreateCmd.MarkFlagRequired("price")

	// Cancel all flags
	ordersCancelAllCmd.Flags().StringVar(&orderCancelAllMarket, "market", "", "filter by market ticker")

	// Amend flags
	ordersAmendCmd.Flags().IntVar(&orderAmendQty, "qty", 0, "new quantity")
	ordersAmendCmd.Flags().IntVar(&orderAmendPrice, "price", 0, "new price in cents")

	// Batch create flags
	ordersBatchCreateCmd.Flags().StringVar(&batchFile, "file", "", "path to JSON file containing orders (required)")
	ordersBatchCreateCmd.MarkFlagRequired("file")
}

// createAPIClient is defined in helpers.go

func getEnvironmentLabel() string {
	cfg := GetConfig()
	if cfg.API.Production {
		return ui.ProdStyle.Render(" PRODUCTION ")
	}
	return ui.DemoStyle.Render(" DEMO ")
}


func runOrdersList(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := make(map[string]string)
	if orderStatusFilter != "" {
		params["status"] = orderStatusFilter
	}
	if orderMarketFilter != "" {
		params["ticker"] = orderMarketFilter
	}

	var response models.OrdersResponse
	path := "/trade-api/v2/portfolio/orders"
	if len(params) > 0 {
		path += api.BuildQueryString(params)
	}

	if err := client.GetJSON(ctx, path, &response); err != nil {
		return fmt.Errorf("failed to list orders: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrdersTable(response.Orders) },
		response.Orders,
		func() { renderOrdersPlain(response.Orders) },
	)
}

func runOrdersGet(cmd *cobra.Command, args []string) error {
	orderID := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var response models.OrderResponse
	path := fmt.Sprintf("/trade-api/v2/portfolio/orders/%s", orderID)

	if err := client.GetJSON(ctx, path, &response); err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrderDetails(response.Order) },
		response.Order,
		func() { renderOrderPlain(response.Order) },
	)
}

func runOrdersCreate(cmd *cobra.Command, args []string) error {
	// Validate price range
	if orderCreatePrice < 1 || orderCreatePrice > 99 {
		return fmt.Errorf("price must be between 1 and 99 cents, got %d", orderCreatePrice)
	}

	// Validate side
	side := strings.ToLower(orderSide)
	if side != "yes" && side != "no" {
		return fmt.Errorf("side must be 'yes' or 'no', got '%s'", orderSide)
	}

	// Validate action
	action := strings.ToLower(orderAction)
	if action != "buy" && action != "sell" {
		return fmt.Errorf("action must be 'buy' or 'sell', got '%s'", orderAction)
	}

	// Validate type
	oType := strings.ToLower(orderType)
	if oType != "limit" && oType != "market" {
		return fmt.Errorf("type must be 'limit' or 'market', got '%s'", orderType)
	}

	// Validate quantity
	if orderCreateQty <= 0 {
		return fmt.Errorf("quantity must be positive, got %d", orderCreateQty)
	}

	// Build order request
	orderReq := models.CreateOrderRequest{
		Ticker: orderCreateMarket,
		Side:   models.OrderSide(side),
		Action: models.OrderAction(action),
		Type:   models.OrderType(oType),
		Count:  orderCreateQty,
	}

	if side == "yes" {
		orderReq.YesPrice = orderCreatePrice
	} else {
		orderReq.NoPrice = orderCreatePrice
	}

	// Show order preview
	fmt.Println()
	fmt.Println(ui.HeaderStyle.Render("Order Preview"))
	fmt.Println()
	fmt.Printf("  Environment:  %s\n", getEnvironmentLabel())
	fmt.Printf("  Market:       %s\n", orderReq.Ticker)
	fmt.Printf("  Side:         %s\n", strings.ToUpper(side))
	fmt.Printf("  Action:       %s\n", strings.ToUpper(action))
	fmt.Printf("  Type:         %s\n", strings.ToUpper(oType))
	fmt.Printf("  Quantity:     %d contracts\n", orderReq.Count)
	fmt.Printf("  Price:        %d cents\n", orderCreatePrice)

	// Calculate potential cost/payout
	potentialCost := orderCreateQty * orderCreatePrice
	potentialPayout := orderCreateQty * 100

	if action == "buy" {
		fmt.Printf("  Max Cost:     %s\n", ui.FormatPrice(potentialCost))
		fmt.Printf("  Max Payout:   %s\n", ui.FormatPrice(potentialPayout))
	} else {
		fmt.Printf("  Max Credit:   %s\n", ui.FormatPrice(potentialCost))
	}
	fmt.Println()

	// Confirm unless --yes flag
	cfg := GetConfig()
	envWarning := ""
	if cfg.API.Production {
		envWarning = " (PRODUCTION - real money)"
	}

	if !confirmAction(fmt.Sprintf("Submit this order%s?", envWarning)) {
		PrintWarning("Order cancelled")
		return nil
	}

	// Submit order
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var response models.CreateOrderResponse
	if err := client.PostJSON(ctx, "/trade-api/v2/portfolio/orders", orderReq, &response); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	PrintSuccess("Order created successfully!")
	fmt.Printf("Order ID: %s\n", response.Order.OrderID)

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrderDetails(response.Order) },
		response.Order,
		func() { renderOrderPlain(response.Order) },
	)
}

func runOrdersCancel(cmd *cobra.Command, args []string) error {
	orderID := args[0]

	if !confirmAction(fmt.Sprintf("Cancel order %s?", orderID)) {
		PrintWarning("Cancellation aborted")
		return nil
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var response models.OrderResponse
	path := fmt.Sprintf("/trade-api/v2/portfolio/orders/%s", orderID)

	if err := client.DeleteJSON(ctx, path, &response); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	PrintSuccess("Order cancelled successfully!")

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrderDetails(response.Order) },
		response.Order,
		func() { renderOrderPlain(response.Order) },
	)
}

func runOrdersCancelAll(cmd *cobra.Command, args []string) error {
	ticker := orderCancelAllMarket

	prompt := "Cancel ALL resting orders?"
	if ticker != "" {
		prompt = fmt.Sprintf("Cancel all resting orders for market %s?", ticker)
	}

	if !confirmAction(prompt) {
		PrintWarning("Cancellation aborted")
		return nil
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := models.BatchCancelOrdersRequest{}
	if ticker != "" {
		req.Ticker = ticker
	}

	var response models.BatchCancelOrdersResponse
	if err := client.DeleteWithBody(ctx, "/trade-api/v2/portfolio/orders", req, &response); err != nil {
		return fmt.Errorf("failed to cancel orders: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Cancelled %d orders", len(response.Orders)))

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrdersTable(response.Orders) },
		response.Orders,
		func() { renderOrdersPlain(response.Orders) },
	)
}

func runOrdersAmend(cmd *cobra.Command, args []string) error {
	orderID := args[0]

	if orderAmendQty == 0 && orderAmendPrice == 0 {
		return fmt.Errorf("at least one of --qty or --price must be specified")
	}

	if orderAmendPrice != 0 && (orderAmendPrice < 1 || orderAmendPrice > 99) {
		return fmt.Errorf("price must be between 1 and 99 cents, got %d", orderAmendPrice)
	}

	// Build amend request
	amendReq := models.AmendOrderRequest{}
	if orderAmendQty > 0 {
		amendReq.Count = orderAmendQty
	}
	if orderAmendPrice > 0 {
		amendReq.Price = orderAmendPrice
	}

	// Show amendment preview
	fmt.Println()
	fmt.Println(ui.HeaderStyle.Render("Amend Order Preview"))
	fmt.Println()
	fmt.Printf("  Environment:  %s\n", getEnvironmentLabel())
	fmt.Printf("  Order ID:     %s\n", orderID)
	if orderAmendQty > 0 {
		fmt.Printf("  New Quantity: %d contracts\n", orderAmendQty)
	}
	if orderAmendPrice > 0 {
		fmt.Printf("  New Price:    %d cents\n", orderAmendPrice)
	}
	fmt.Println()

	if !confirmAction("Amend this order?") {
		PrintWarning("Amendment cancelled")
		return nil
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var response models.OrderResponse
	path := fmt.Sprintf("/trade-api/v2/portfolio/orders/%s", orderID)

	if err := client.PatchJSON(ctx, path, amendReq, &response); err != nil {
		return fmt.Errorf("failed to amend order: %w", err)
	}

	PrintSuccess("Order amended successfully!")

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrderDetails(response.Order) },
		response.Order,
		func() { renderOrderPlain(response.Order) },
	)
}

func runOrdersBatchCreate(cmd *cobra.Command, args []string) error {
	// Read and parse the JSON file
	data, err := os.ReadFile(batchFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var orders []models.CreateOrderRequest
	if err := json.Unmarshal(data, &orders); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(orders) == 0 {
		return fmt.Errorf("no orders found in file")
	}

	// Validate all orders
	for i, order := range orders {
		if order.Ticker == "" {
			return fmt.Errorf("order %d: ticker is required", i+1)
		}
		if order.Side != models.OrderSideYes && order.Side != models.OrderSideNo {
			return fmt.Errorf("order %d: side must be 'yes' or 'no'", i+1)
		}
		if order.Count <= 0 {
			return fmt.Errorf("order %d: count must be positive", i+1)
		}
		if order.YesPrice > 0 && (order.YesPrice < 1 || order.YesPrice > 99) {
			return fmt.Errorf("order %d: yes_price must be between 1 and 99", i+1)
		}
		if order.NoPrice > 0 && (order.NoPrice < 1 || order.NoPrice > 99) {
			return fmt.Errorf("order %d: no_price must be between 1 and 99", i+1)
		}
	}

	// Show preview
	fmt.Println()
	fmt.Println(ui.HeaderStyle.Render("Batch Order Preview"))
	fmt.Println()
	fmt.Printf("  Environment:  %s\n", getEnvironmentLabel())
	fmt.Printf("  Total Orders: %d\n", len(orders))
	fmt.Println()

	// Show each order
	for i, order := range orders {
		price := order.YesPrice
		if order.Side == models.OrderSideNo {
			price = order.NoPrice
		}
		fmt.Printf("  %d. %s %s %s @ %d cents x %d\n",
			i+1,
			strings.ToUpper(string(order.Action)),
			strings.ToUpper(string(order.Side)),
			order.Ticker,
			price,
			order.Count,
		)
	}
	fmt.Println()

	cfg := GetConfig()
	envWarning := ""
	if cfg.API.Production {
		envWarning = " (PRODUCTION - real money)"
	}

	if !confirmAction(fmt.Sprintf("Submit %d orders%s?", len(orders), envWarning)) {
		PrintWarning("Batch order cancelled")
		return nil
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	batchReq := models.BatchCreateOrdersRequest{Orders: orders}
	var response models.BatchCreateOrdersResponse

	if err := client.PostJSON(ctx, "/trade-api/v2/portfolio/orders/batched", batchReq, &response); err != nil {
		return fmt.Errorf("failed to create batch orders: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Created %d orders successfully!", len(response.Orders)))

	return ui.Output(
		GetOutputFormat(),
		func() { renderOrdersTable(response.Orders) },
		response.Orders,
		func() { renderOrdersPlain(response.Orders) },
	)
}

func runOrdersQueue(cmd *cobra.Command, args []string) error {
	orderID := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var response models.QueuePositionsResponse
	path := fmt.Sprintf("/trade-api/v2/portfolio/orders/%s/queue-position", orderID)

	if err := client.GetJSON(ctx, path, &response); err != nil {
		return fmt.Errorf("failed to get queue position: %w", err)
	}

	if len(response.Positions) == 0 {
		return fmt.Errorf("no queue position found for order %s", orderID)
	}

	position := response.Positions[0]

	return ui.Output(
		GetOutputFormat(),
		func() {
			fmt.Println()
			fmt.Println(ui.HeaderStyle.Render("Queue Position"))
			fmt.Println()
			fmt.Printf("  Order ID:  %s\n", position.OrderID)
			fmt.Printf("  Position:  %d\n", position.QueuePosition)
			fmt.Println()
		},
		position,
		func() {
			fmt.Printf("%s\t%d\n", position.OrderID, position.QueuePosition)
		},
	)
}

// Render functions

func renderOrdersTable(orders []models.Order) {
	if len(orders) == 0 {
		fmt.Println("No orders found")
		return
	}

	headers := []string{"Order ID", "Market", "Side", "Price", "Qty", "Status", "Created"}
	rows := make([][]string, 0, len(orders))

	for _, order := range orders {
		price := order.YesPrice
		if order.Side == models.OrderSideNo {
			price = order.NoPrice
		}

		rows = append(rows, []string{
			truncateOrderID(order.OrderID),
			order.Ticker,
			strings.ToUpper(string(order.Side)),
			fmt.Sprintf("%d", price),
			fmt.Sprintf("%d/%d", order.RemainingCount, order.InitialCount),
			formatOrderStatusModel(order.Status),
			order.CreatedTime.Format("2006-01-02 15:04"),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderOrdersPlain(orders []models.Order) {
	for _, order := range orders {
		price := order.YesPrice
		if order.Side == models.OrderSideNo {
			price = order.NoPrice
		}
		fmt.Printf("%s\t%s\t%s\t%d\t%d\t%s\n",
			order.OrderID,
			order.Ticker,
			order.Side,
			price,
			order.RemainingCount,
			order.Status,
		)
	}
}

func renderOrderDetails(order models.Order) {
	price := order.YesPrice
	if order.Side == models.OrderSideNo {
		price = order.NoPrice
	}

	pairs := [][]string{
		{"Order ID", order.OrderID},
		{"Market", order.Ticker},
		{"Status", formatOrderStatusModel(order.Status)},
		{"Side", strings.ToUpper(string(order.Side))},
		{"Action", strings.ToUpper(string(order.Action))},
		{"Type", strings.ToUpper(string(order.Type))},
		{"Price", fmt.Sprintf("%d cents", price)},
		{"Initial Qty", fmt.Sprintf("%d", order.InitialCount)},
		{"Remaining Qty", fmt.Sprintf("%d", order.RemainingCount)},
		{"Filled Qty", fmt.Sprintf("%d", order.FillCount)},
		{"Created", order.CreatedTime.Format("2006-01-02 15:04:05")},
		{"Last Updated", order.LastUpdateTime.Format("2006-01-02 15:04:05")},
	}

	if order.TakerFillCost > 0 || order.TakerFees > 0 {
		pairs = append(pairs, []string{"Taker Fills", fmt.Sprintf("%d", order.TakerFillCount)})
		pairs = append(pairs, []string{"Taker Cost", ui.FormatPrice(order.TakerFillCost)})
		pairs = append(pairs, []string{"Taker Fees", ui.FormatPrice(order.TakerFees)})
	}

	if order.MakerFillCost > 0 || order.MakerFees > 0 {
		pairs = append(pairs, []string{"Maker Fills", fmt.Sprintf("%d", order.MakerFillCount)})
		pairs = append(pairs, []string{"Maker Cost", ui.FormatPrice(order.MakerFillCost)})
		pairs = append(pairs, []string{"Maker Fees", ui.FormatPrice(order.MakerFees)})
	}

	if order.ClientOrderID != "" {
		pairs = append(pairs, []string{"Client Order ID", order.ClientOrderID})
	}

	if order.OrderGroupID != "" {
		pairs = append(pairs, []string{"Order Group ID", order.OrderGroupID})
	}

	fmt.Println()
	ui.RenderKeyValue(pairs)
	fmt.Println()
}

func renderOrderPlain(order models.Order) {
	price := order.YesPrice
	if order.Side == models.OrderSideNo {
		price = order.NoPrice
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%d\t%d/%d\t%s\n",
		order.OrderID,
		order.Ticker,
		order.Side,
		order.Status,
		price,
		order.RemainingCount,
		order.InitialCount,
		order.CreatedTime.Format("2006-01-02T15:04:05Z"),
	)
}

func formatOrderStatusModel(status models.OrderStatus) string {
	switch status {
	case models.OrderStatusResting:
		return ui.StatusActiveStyle.Render("RESTING")
	case models.OrderStatusExecuted:
		return ui.SuccessStyle.Render("EXECUTED")
	case models.OrderStatusCanceled:
		return ui.MutedStyle.Render("CANCELED")
	case models.OrderStatusPending:
		return ui.WarningStyle.Render("PENDING")
	default:
		return string(status)
	}
}

func truncateOrderID(id string) string {
	if len(id) > 12 {
		return id[:12] + "..."
	}
	return id
}
