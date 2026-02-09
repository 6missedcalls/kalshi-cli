package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var exchangeCmd = &cobra.Command{
	Use:   "exchange",
	Short: "Exchange information and status",
	Long:  `Get exchange status, schedule, and announcements.`,
}

var exchangeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get exchange status",
	Long:  `Get the current exchange status including trading activity and environment.`,
	RunE:  runExchangeStatus,
}

var exchangeScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Get exchange schedule",
	Long:  `Get the exchange trading schedule.`,
	RunE:  runExchangeSchedule,
}

var exchangeAnnouncementsCmd = &cobra.Command{
	Use:   "announcements",
	Short: "Get exchange announcements",
	Long:  `Get the latest exchange announcements.`,
	RunE:  runExchangeAnnouncements,
}

func init() {
	rootCmd.AddCommand(exchangeCmd)
	exchangeCmd.AddCommand(exchangeStatusCmd)
	exchangeCmd.AddCommand(exchangeScheduleCmd)
	exchangeCmd.AddCommand(exchangeAnnouncementsCmd)
}

func createAPIClient() (*api.Client, error) {
	keyStore, err := config.NewKeyringStore()
	if err != nil {
		return nil, fmt.Errorf("failed to access keyring: %w", err)
	}

	creds, err := keyStore.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds == nil {
		return nil, fmt.Errorf("no credentials found - run 'kalshi-cli auth login' first")
	}

	signer, err := api.NewSignerFromPEM(creds.APIKeyID, creds.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	cfg := GetConfig()
	client := api.NewClient(signer, api.WithBaseURL(cfg.BaseURL()))

	return client, nil
}

func runExchangeStatus(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, err := client.GetExchangeStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get exchange status: %w", err)
	}

	cfg := GetConfig()
	outputFmt := GetOutputFormat()

	return ui.Output(
		outputFmt,
		func() { renderExchangeStatusTable(status, cfg) },
		status,
		func() { renderExchangeStatusPlain(status, cfg) },
	)
}

func renderExchangeStatusTable(status *api.ExchangeStatusResponse, cfg *config.Config) {
	exchangeActive := formatStatusBool(status.ExchangeActive)
	tradingActive := formatStatusBool(status.TradingActive)
	environment := formatEnvironment(cfg.API.Production)

	pairs := [][]string{
		{ui.BoldStyle.Render("Exchange Active:"), exchangeActive},
		{ui.BoldStyle.Render("Trading Active:"), tradingActive},
		{ui.BoldStyle.Render("Environment:"), environment},
	}

	ui.RenderKeyValue(pairs)
}

func renderExchangeStatusPlain(status *api.ExchangeStatusResponse, cfg *config.Config) {
	exchangeActive := boolToYesNo(status.ExchangeActive)
	tradingActive := boolToYesNo(status.TradingActive)
	environment := cfg.Environment()

	fmt.Printf("exchange_active=%s\n", exchangeActive)
	fmt.Printf("trading_active=%s\n", tradingActive)
	fmt.Printf("environment=%s\n", environment)
}

func formatStatusBool(active bool) string {
	if active {
		return ui.SuccessStyle.Render("Yes")
	}
	return ui.ErrorStyle.Render("No")
}

func formatEnvironment(isProd bool) string {
	if isProd {
		return ui.ProdStyle.Render("Production")
	}
	return ui.DemoStyle.Render("Demo")
}

func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func runExchangeSchedule(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var schedule models.ExchangeScheduleResponse
	if err := client.GetExchangeSchedule(ctx, &schedule); err != nil {
		return fmt.Errorf("failed to get exchange schedule: %w", err)
	}

	outputFmt := GetOutputFormat()

	return ui.Output(
		outputFmt,
		func() { renderScheduleTable(schedule) },
		schedule,
		func() { renderSchedulePlain(schedule) },
	)
}

func renderScheduleTable(schedule models.ExchangeScheduleResponse) {
	if len(schedule.Schedule.ScheduleEntries) == 0 {
		fmt.Println(ui.MutedStyle.Render("No schedule entries found."))
		return
	}

	headers := []string{"Start Time", "End Time", "Maintenance"}
	var rows [][]string

	for _, entry := range schedule.Schedule.ScheduleEntries {
		maintenance := "No"
		if entry.Maintenance {
			maintenance = ui.WarningStyle.Render("Yes")
		}

		rows = append(rows, []string{
			formatTime(entry.StartTime),
			formatTime(entry.EndTime),
			maintenance,
		})
	}

	ui.RenderTable(headers, rows)
}

func renderSchedulePlain(schedule models.ExchangeScheduleResponse) {
	for i, entry := range schedule.Schedule.ScheduleEntries {
		maintenance := "no"
		if entry.Maintenance {
			maintenance = "yes"
		}
		fmt.Printf("entry_%d_start=%s\n", i, entry.StartTime.Format(time.RFC3339))
		fmt.Printf("entry_%d_end=%s\n", i, entry.EndTime.Format(time.RFC3339))
		fmt.Printf("entry_%d_maintenance=%s\n", i, maintenance)
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Local().Format("2006-01-02 15:04 MST")
}

func runExchangeAnnouncements(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var announcements models.AnnouncementsResponse
	if err := client.GetAnnouncements(ctx, &announcements); err != nil {
		return fmt.Errorf("failed to get announcements: %w", err)
	}

	outputFmt := GetOutputFormat()

	return ui.Output(
		outputFmt,
		func() { renderAnnouncementsTable(announcements) },
		announcements,
		func() { renderAnnouncementsPlain(announcements) },
	)
}

func renderAnnouncementsTable(announcements models.AnnouncementsResponse) {
	if len(announcements.Announcements) == 0 {
		fmt.Println(ui.MutedStyle.Render("No announcements found."))
		return
	}

	headers := []string{"Title", "Type", "Status", "Delivery Time"}
	var rows [][]string

	for _, ann := range announcements.Announcements {
		status := formatAnnouncementStatus(ann.Status)

		rows = append(rows, []string{
			truncateString(ann.Title, 50),
			ann.Type,
			status,
			formatTime(ann.DeliveryTime),
		})
	}

	ui.RenderTable(headers, rows)
}

func renderAnnouncementsPlain(announcements models.AnnouncementsResponse) {
	for i, ann := range announcements.Announcements {
		fmt.Printf("announcement_%d_id=%s\n", i, ann.ID)
		fmt.Printf("announcement_%d_title=%s\n", i, ann.Title)
		fmt.Printf("announcement_%d_type=%s\n", i, ann.Type)
		fmt.Printf("announcement_%d_status=%s\n", i, ann.Status)
		fmt.Printf("announcement_%d_delivery_time=%s\n", i, ann.DeliveryTime.Format(time.RFC3339))
	}
}

func formatAnnouncementStatus(status string) string {
	switch status {
	case "active":
		return ui.SuccessStyle.Render(status)
	case "pending":
		return ui.WarningStyle.Render(status)
	case "expired":
		return ui.MutedStyle.Render(status)
	default:
		return status
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
