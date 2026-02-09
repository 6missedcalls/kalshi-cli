package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

var orderGroupsCmd = &cobra.Command{
	Use:     "order-groups",
	Aliases: []string{"og"},
	Short:   "Manage order groups",
	Long: `Order groups allow you to group multiple orders together and manage them
as a single unit. You can set limits on the number of contracts that can
be filled across all orders in a group.`,
}

var orderGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List order groups",
	Long:  `List all order groups for the authenticated user.`,
	RunE:  runOrderGroupsList,
}

var orderGroupsGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get order group details",
	Long:  `Get detailed information about a specific order group.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOrderGroupsGet,
}

var orderGroupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new order group",
	Long: `Create a new order group with a specified contract limit.

The limit specifies the maximum number of contracts that can be filled
across all orders in the group.`,
	RunE: runOrderGroupsCreate,
}

var orderGroupsDeleteCmd = &cobra.Command{
	Use:   "delete <group-id>",
	Short: "Delete an order group",
	Long:  `Delete an order group. All orders in the group will be canceled.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOrderGroupsDelete,
}

var orderGroupsResetCmd = &cobra.Command{
	Use:   "reset <group-id>",
	Short: "Reset an order group",
	Long:  `Reset an order group's filled count to zero, allowing more orders to fill.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOrderGroupsReset,
}

var orderGroupsTriggerCmd = &cobra.Command{
	Use:   "trigger <group-id>",
	Short: "Trigger an order group",
	Long:  `Trigger an order group to execute its orders.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runOrderGroupsTrigger,
}

var (
	orderGroupLimit  int
	orderGroupStatus string
)

func init() {
	rootCmd.AddCommand(orderGroupsCmd)

	orderGroupsCmd.AddCommand(orderGroupsListCmd)
	orderGroupsCmd.AddCommand(orderGroupsGetCmd)
	orderGroupsCmd.AddCommand(orderGroupsCreateCmd)
	orderGroupsCmd.AddCommand(orderGroupsDeleteCmd)
	orderGroupsCmd.AddCommand(orderGroupsResetCmd)
	orderGroupsCmd.AddCommand(orderGroupsTriggerCmd)

	orderGroupsListCmd.Flags().StringVar(&orderGroupStatus, "status", "", "filter by status")

	orderGroupsCreateCmd.Flags().IntVar(&orderGroupLimit, "limit", 0, "maximum contracts to fill (required)")
	orderGroupsCreateCmd.MarkFlagRequired("limit")
}

// getAPIClient uses createClient from helpers.go

func runOrderGroupsList(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}

	opts := api.OrderGroupsOptions{
		Status: orderGroupStatus,
	}

	result, err := client.GetOrderGroups(context.Background(), opts)
	if err != nil {
		return err
	}

	return outputOrderGroupsList(result.OrderGroups)
}

func runOrderGroupsGet(cmd *cobra.Command, args []string) error {
	groupID := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	result, err := client.GetOrderGroup(context.Background(), groupID)
	if err != nil {
		return err
	}

	return outputOrderGroupDetails(&result.OrderGroup)
}

func runOrderGroupsCreate(cmd *cobra.Command, args []string) error {
	if orderGroupLimit <= 0 {
		return fmt.Errorf("limit must be a positive integer")
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	req := models.CreateOrderGroupRequest{Limit: orderGroupLimit}
	result, err := client.CreateOrderGroup(context.Background(), req)
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Created order group: %s", result.OrderGroup.GroupID))
	return outputOrderGroupDetails(&result.OrderGroup)
}

func runOrderGroupsDelete(cmd *cobra.Command, args []string) error {
	groupID := args[0]

	if !SkipConfirmation() {
		fmt.Printf("Are you sure you want to delete order group %s? (y/N): ", groupID)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			PrintWarning("Deletion canceled")
			return nil
		}
	}

	client, err := createClient()
	if err != nil {
		return err
	}

	if err := client.DeleteOrderGroup(context.Background(), groupID); err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Deleted order group: %s", groupID))
	return nil
}

func runOrderGroupsReset(cmd *cobra.Command, args []string) error {
	groupID := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	result, err := client.ResetOrderGroup(context.Background(), groupID)
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Reset order group: %s", groupID))
	return outputOrderGroupDetails(&result.OrderGroup)
}

func runOrderGroupsTrigger(cmd *cobra.Command, args []string) error {
	groupID := args[0]

	client, err := createClient()
	if err != nil {
		return err
	}

	result, err := client.TriggerOrderGroup(context.Background(), groupID)
	if err != nil {
		return err
	}

	PrintSuccess(fmt.Sprintf("Triggered order group: %s", groupID))
	return outputOrderGroupDetails(&result.OrderGroup)
}

func outputOrderGroupsList(groups []models.OrderGroup) error {
	format := GetOutputFormat()

	tableFunc := func() {
		headers := []string{"Group ID", "Status", "Limit", "Filled", "Order Count"}
		rows := make([][]string, 0, len(groups))

		for _, g := range groups {
			rows = append(rows, []string{
				g.GroupID,
				g.Status,
				strconv.Itoa(g.Limit),
				strconv.Itoa(g.FilledCount),
				strconv.Itoa(g.OrderCount),
			})
		}

		ui.RenderTable(headers, rows)
	}

	plainFunc := func() {
		for _, g := range groups {
			ui.PrintPlain("%s\t%s\t%d\t%d\t%d",
				g.GroupID, g.Status, g.Limit, g.FilledCount, g.OrderCount)
		}
	}

	return ui.Output(format, tableFunc, groups, plainFunc)
}

func outputOrderGroupDetails(group *models.OrderGroup) error {
	format := GetOutputFormat()

	tableFunc := func() {
		pairs := [][]string{
			{"Group ID", group.GroupID},
			{"Status", group.Status},
			{"Limit", strconv.Itoa(group.Limit)},
			{"Filled Count", strconv.Itoa(group.FilledCount)},
			{"Order Count", strconv.Itoa(group.OrderCount)},
			{"Created", group.CreatedTime.Format("2006-01-02 15:04:05")},
			{"Last Updated", group.LastUpdateTime.Format("2006-01-02 15:04:05")},
		}

		if len(group.OrderIDs) > 0 {
			pairs = append(pairs, []string{"Order IDs", fmt.Sprintf("%v", group.OrderIDs)})
		}

		ui.RenderKeyValue(pairs)
	}

	plainFunc := func() {
		ui.PrintPlain("group_id=%s status=%s limit=%d filled=%d orders=%d",
			group.GroupID, group.Status, group.Limit, group.FilledCount, group.OrderCount)
	}

	return ui.Output(format, tableFunc, group, plainFunc)
}
