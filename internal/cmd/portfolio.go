package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "Manage your portfolio and account",
	Long: `View and manage your Kalshi portfolio including balance, positions,
fills, settlements, and subaccounts.`,
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show account balance",
	Long:  `Display your current account balance including available balance, portfolio value, and total balance.`,
	RunE:  runBalance,
}

var positionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "List positions",
	Long:  `List your current market positions with details including average cost, P&L, and exposure.`,
	RunE:  runPositions,
}

var fillsCmd = &cobra.Command{
	Use:   "fills",
	Short: "List fills",
	Long:  `List your trade fills showing executed orders and their details.`,
	RunE:  runFills,
}

var settlementsCmd = &cobra.Command{
	Use:   "settlements",
	Short: "List settlements",
	Long:  `List your market settlements showing resolved positions and their outcomes.`,
	RunE:  runSettlements,
}

var subaccountsCmd = &cobra.Command{
	Use:   "subaccounts",
	Short: "Manage subaccounts",
	Long:  `List, create, and manage subaccounts for your Kalshi account.`,
}

var subaccountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List subaccounts",
	Long:  `List all subaccounts associated with your account.`,
	RunE:  runSubaccountsList,
}

var subaccountsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create subaccount",
	Long:  `Create a new subaccount.`,
	RunE:  runSubaccountsCreate,
}

var subaccountsTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer between subaccounts",
	Long:  `Transfer funds between subaccounts. Requires --from, --to, and --amount flags.`,
	RunE:  runSubaccountsTransfer,
}

var (
	positionsMarket   string
	fillsLimit        int
	settlementsLimit  int
	transferFrom      int
	transferTo        int
	transferAmount    int
)

func init() {
	rootCmd.AddCommand(portfolioCmd)

	portfolioCmd.AddCommand(balanceCmd)
	portfolioCmd.AddCommand(positionsCmd)
	portfolioCmd.AddCommand(fillsCmd)
	portfolioCmd.AddCommand(settlementsCmd)
	portfolioCmd.AddCommand(subaccountsCmd)

	subaccountsCmd.AddCommand(subaccountsListCmd)
	subaccountsCmd.AddCommand(subaccountsCreateCmd)
	subaccountsCmd.AddCommand(subaccountsTransferCmd)

	positionsCmd.Flags().StringVar(&positionsMarket, "market", "", "filter by market ticker")

	fillsCmd.Flags().IntVar(&fillsLimit, "limit", 100, "maximum number of fills to return")

	settlementsCmd.Flags().IntVar(&settlementsLimit, "limit", 50, "maximum number of settlements to return")

	subaccountsTransferCmd.Flags().IntVar(&transferFrom, "from", 0, "source subaccount ID")
	subaccountsTransferCmd.Flags().IntVar(&transferTo, "to", 0, "destination subaccount ID")
	subaccountsTransferCmd.Flags().IntVar(&transferAmount, "amount", 0, "amount to transfer in cents")
	subaccountsTransferCmd.MarkFlagRequired("from")
	subaccountsTransferCmd.MarkFlagRequired("to")
	subaccountsTransferCmd.MarkFlagRequired("amount")
}

func runBalance(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	balance, err := client.GetBalance(ctx)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderBalanceTable(balance) },
		balance,
		func() { renderBalancePlain(balance) },
	)
}

func renderBalanceTable(balance *models.BalanceResponse) {
	pairs := [][]string{
		{ui.BoldStyle.Render("Available Balance:"), ui.FormatPrice(balance.Balance)},
		{ui.BoldStyle.Render("Portfolio Value:"), ui.FormatPrice(balance.PortfolioValue)},
		{ui.BoldStyle.Render("Total Balance:"), ui.FormatPrice(balance.Balance + balance.PortfolioValue)},
	}

	ui.RenderKeyValue(pairs)
}

func renderBalancePlain(balance *models.BalanceResponse) {
	ui.PrintPlain("Available Balance: %s", ui.FormatPrice(balance.Balance))
	ui.PrintPlain("Portfolio Value: %s", ui.FormatPrice(balance.PortfolioValue))
	ui.PrintPlain("Total Balance: %s", ui.FormatPrice(balance.Balance+balance.PortfolioValue))
}

func runPositions(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	opts := api.PositionsOptions{
		Ticker: positionsMarket,
	}

	positions, err := client.GetPositions(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to get positions: %w", err)
	}

	if len(positions.Positions) == 0 {
		PrintWarning("No positions found")
		return nil
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderPositionsTable(positions.Positions) },
		positions,
		func() { renderPositionsPlain(positions.Positions) },
	)
}

func renderPositionsTable(positions []models.MarketPosition) {
	headers := []string{"Market", "Position", "Avg Cost", "P&L", "Exposure"}
	rows := make([][]string, 0, len(positions))

	for _, p := range positions {
		avgCost := calculateAvgCost(p)
		pnl := p.RealizedPnl

		pnlStr := ui.FormatPriceStyled(pnl, pnl >= 0)

		rows = append(rows, []string{
			p.Ticker,
			formatPosition(p.Position),
			ui.FormatPrice(avgCost),
			pnlStr,
			ui.FormatPrice(p.MarketExposure),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderPositionsPlain(positions []models.MarketPosition) {
	for _, p := range positions {
		avgCost := calculateAvgCost(p)
		ui.PrintPlain("%s\t%s\t%s\t%s\t%s",
			p.Ticker,
			formatPosition(p.Position),
			ui.FormatPrice(avgCost),
			ui.FormatPrice(p.RealizedPnl),
			ui.FormatPrice(p.MarketExposure),
		)
	}
}

func calculateAvgCost(p models.MarketPosition) int {
	if p.Position == 0 {
		return 0
	}
	return p.TotalTraded / abs(p.Position)
}

func formatPosition(position int) string {
	if position == 0 {
		return "0"
	}
	if position > 0 {
		return fmt.Sprintf("+%d", position)
	}
	return fmt.Sprintf("-%d", abs(position))
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func runFills(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	opts := api.FillsOptions{
		Limit: fillsLimit,
	}

	fills, err := client.GetFills(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to get fills: %w", err)
	}

	if len(fills.Fills) == 0 {
		PrintWarning("No fills found")
		return nil
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderFillsTable(fills.Fills) },
		fills,
		func() { renderFillsPlain(fills.Fills) },
	)
}

func renderFillsTable(fills []models.Fill) {
	headers := []string{"Time", "Ticker", "Side", "Action", "Count", "Price", "Taker"}
	rows := make([][]string, 0, len(fills))

	for _, f := range fills {
		price := f.YesPrice
		if f.Side == "no" {
			price = f.NoPrice
		}

		takerStr := "No"
		if f.IsTaker {
			takerStr = "Yes"
		}

		rows = append(rows, []string{
			f.CreatedTime.Format("2006-01-02 15:04"),
			f.Ticker,
			strings.ToUpper(f.Side),
			strings.ToUpper(f.Action),
			strconv.Itoa(f.Count),
			ui.FormatPrice(price),
			takerStr,
		})
	}

	ui.RenderTable(headers, rows)
}

func renderFillsPlain(fills []models.Fill) {
	for _, f := range fills {
		price := f.YesPrice
		if f.Side == "no" {
			price = f.NoPrice
		}
		ui.PrintPlain("%s\t%s\t%s\t%s\t%d\t%s",
			f.CreatedTime.Format("2006-01-02T15:04:05Z"),
			f.Ticker,
			f.Side,
			f.Action,
			f.Count,
			ui.FormatPrice(price),
		)
	}
}

func runSettlements(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	opts := api.SettlementsOptions{
		Limit: settlementsLimit,
	}

	settlements, err := client.GetSettlements(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to get settlements: %w", err)
	}

	if len(settlements.Settlements) == 0 {
		PrintWarning("No settlements found")
		return nil
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderSettlementsTable(settlements.Settlements) },
		settlements,
		func() { renderSettlementsPlain(settlements.Settlements) },
	)
}

func renderSettlementsTable(settlements []models.Settlement) {
	headers := []string{"Settled", "Ticker", "Result", "Revenue", "Yes Qty", "No Qty"}
	rows := make([][]string, 0, len(settlements))

	for _, s := range settlements {
		revenueStr := ui.FormatPriceStyled(s.Revenue, s.Revenue >= 0)

		rows = append(rows, []string{
			s.SettledTime.Format("2006-01-02"),
			s.Ticker,
			strings.ToUpper(s.MarketResult),
			revenueStr,
			strconv.Itoa(s.YesCount),
			strconv.Itoa(s.NoCount),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderSettlementsPlain(settlements []models.Settlement) {
	for _, s := range settlements {
		ui.PrintPlain("%s\t%s\t%s\t%s\t%d\t%d",
			s.SettledTime.Format("2006-01-02"),
			s.Ticker,
			s.MarketResult,
			ui.FormatPrice(s.Revenue),
			s.YesCount,
			s.NoCount,
		)
	}
}

func runSubaccountsList(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	subaccounts, err := client.GetSubaccounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get subaccounts: %w", err)
	}

	if len(subaccounts.Subaccounts) == 0 {
		PrintWarning("No subaccounts found")
		return nil
	}

	return ui.Output(
		GetOutputFormat(),
		func() { renderSubaccountsTable(subaccounts.Subaccounts) },
		subaccounts,
		func() { renderSubaccountsPlain(subaccounts.Subaccounts) },
	)
}

func renderSubaccountsTable(subaccounts []models.Subaccount) {
	headers := []string{"ID", "Balance", "Available"}
	rows := make([][]string, 0, len(subaccounts))

	for _, s := range subaccounts {
		rows = append(rows, []string{
			strconv.Itoa(s.SubaccountID),
			ui.FormatPrice(s.Balance),
			ui.FormatPrice(s.AvailableBalance),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderSubaccountsPlain(subaccounts []models.Subaccount) {
	for _, s := range subaccounts {
		ui.PrintPlain("%d\t%s\t%s",
			s.SubaccountID,
			ui.FormatPrice(s.Balance),
			ui.FormatPrice(s.AvailableBalance),
		)
	}
}

func runSubaccountsCreate(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	subaccount, err := client.CreateSubaccount(ctx)
	if err != nil {
		return fmt.Errorf("failed to create subaccount: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() {
			PrintSuccess(fmt.Sprintf("Subaccount created successfully (ID: %d)", subaccount.SubaccountID))
		},
		subaccount,
		func() {
			ui.PrintPlain("subaccount_id=%d", subaccount.SubaccountID)
		},
	)
}

func runSubaccountsTransfer(cmd *cobra.Command, args []string) error {
	if transferAmount <= 0 {
		return fmt.Errorf("amount must be positive (in cents)")
	}

	if !SkipConfirmation() {
		confirmed, err := confirmTransfer(transferFrom, transferTo, transferAmount)
		if err != nil {
			return err
		}
		if !confirmed {
			PrintWarning("Transfer cancelled")
			return nil
		}
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	request := &models.TransferRequest{
		FromSubaccount: transferFrom,
		ToSubaccount:   transferTo,
		Amount:         transferAmount,
	}

	transfer, err := client.Transfer(ctx, *request)
	if err != nil {
		return fmt.Errorf("failed to transfer: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() {
			PrintSuccess(fmt.Sprintf("Transfer complete: %s from subaccount %d to %d",
				ui.FormatPrice(transfer.Amount),
				transfer.FromSubaccount,
				transfer.ToSubaccount,
			))
		},
		transfer,
		func() {
			ui.PrintPlain("transfer_id=%s amount=%d from=%d to=%d",
				transfer.TransferID,
				transfer.Amount,
				transfer.FromSubaccount,
				transfer.ToSubaccount,
			)
		},
	)
}

func confirmTransfer(from, to, amount int) (bool, error) {
	fmt.Printf("\nTransfer Details:\n")
	fmt.Printf("  From Subaccount: %d\n", from)
	fmt.Printf("  To Subaccount:   %d\n", to)
	fmt.Printf("  Amount:          %s\n\n", ui.FormatPrice(amount))
	fmt.Print("Confirm transfer? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read confirmation: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}
