package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and API keys",
	Long: `Manage authentication credentials and API keys for the Kalshi API.

The auth commands allow you to:
  - Log in with a new RSA key pair
  - Log out and clear stored credentials
  - Check authentication status
  - Manage API keys`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Kalshi",
	Long: `Authenticate with Kalshi by generating an RSA key pair.

This command will:
  1. Check if you're already logged in
  2. Generate a new 4096-bit RSA key pair
  3. Display the public key for you to add to your Kalshi dashboard
  4. Prompt you for the API key ID after you've added the key
  5. Store your credentials securely in the system keyring
  6. Test the authentication`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	Long:  `Remove stored API credentials from the system keyring.`,
	RunE:  runLogout,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display the current authentication status and environment.`,
	RunE:  runStatus,
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage API keys",
	Long:  `List, create, and delete API keys for your Kalshi account.`,
}

var keysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Long:  `List all API keys associated with your Kalshi account.`,
	RunE:  runKeysList,
}

var keysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Long:  `Create a new API key for your Kalshi account.`,
	RunE:  runKeysCreate,
}

var keysDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an API key",
	Long:  `Delete an API key by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runKeysDelete,
}

var keyName string

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(keysCmd)

	keysCmd.AddCommand(keysListCmd)
	keysCmd.AddCommand(keysCreateCmd)
	keysCmd.AddCommand(keysDeleteCmd)

	keysCreateCmd.Flags().StringVar(&keyName, "name", "", "name for the new API key")
}

func runLogin(cmd *cobra.Command, args []string) error {
	keyring, err := config.NewKeyringStore()
	if err != nil {
		return fmt.Errorf("failed to access keyring: %w", err)
	}

	if keyring.HasCredentials() {
		existingCreds, err := keyring.GetCredentials()
		if err == nil && existingCreds != nil {
			fmt.Println(ui.WarningStyle.Render("You are already logged in."))
			fmt.Printf("API Key ID: %s\n", existingCreds.APIKeyID)
			fmt.Println()

			if !SkipConfirmation() {
				fmt.Print("Do you want to log out and create new credentials? [y/N]: ")
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					fmt.Println("Login cancelled.")
					return nil
				}
			}

			if err := keyring.DeleteCredentials(); err != nil {
				return fmt.Errorf("failed to clear existing credentials: %w", err)
			}
		}
	}

	fmt.Println(ui.TitleStyle.Render("Generating RSA Key Pair"))
	fmt.Println()

	privateKey, err := api.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	publicKeyPEM, err := api.EncodePublicKeyPEM(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %w", err)
	}

	fmt.Println("A new RSA key pair has been generated.")
	fmt.Println()
	fmt.Println(ui.HeaderStyle.Render("Public Key:"))
	fmt.Println(publicKeyPEM)
	fmt.Println()
	fmt.Println(ui.BoldStyle.Render("Next Steps:"))
	fmt.Println("1. Copy the public key above")
	fmt.Println("2. Go to your Kalshi dashboard: https://kalshi.com/account/api")
	fmt.Println("3. Click 'Add API Key' and paste the public key")
	fmt.Println("4. Copy the API Key ID provided by Kalshi")
	fmt.Println()

	fmt.Print("Enter your API Key ID: ")
	reader := bufio.NewReader(os.Stdin)
	apiKeyID, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	apiKeyID = strings.TrimSpace(apiKeyID)

	if apiKeyID == "" {
		return fmt.Errorf("API Key ID is required")
	}

	privateKeyPEM := api.EncodePrivateKeyPEM(privateKey)

	creds := config.Credentials{
		APIKeyID:   apiKeyID,
		PrivateKey: privateKeyPEM,
	}

	if err := keyring.SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Println()
	fmt.Println("Testing authentication...")

	client, err := createAuthenticatedClient(creds)
	if err != nil {
		if deleteErr := keyring.DeleteCredentials(); deleteErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to clean up credentials: %v\n", deleteErr)
		}
		return fmt.Errorf("authentication failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := client.GetExchangeStatus(ctx)
	if err != nil {
		if deleteErr := keyring.DeleteCredentials(); deleteErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to clean up credentials: %v\n", deleteErr)
		}
		return fmt.Errorf("authentication test failed: %w", err)
	}

	fmt.Println()
	PrintSuccess("Authentication successful!")
	fmt.Printf("Exchange Active: %v\n", status.ExchangeActive)
	fmt.Printf("Trading Active: %v\n", status.TradingActive)
	fmt.Printf("Environment: %s\n", cfg.Environment())

	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	keyring, err := config.NewKeyringStore()
	if err != nil {
		return fmt.Errorf("failed to access keyring: %w", err)
	}

	if !keyring.HasCredentials() {
		fmt.Println("You are not logged in.")
		return nil
	}

	if !SkipConfirmation() {
		fmt.Print("Are you sure you want to log out? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Logout cancelled.")
			return nil
		}
	}

	if err := keyring.DeleteCredentials(); err != nil {
		return fmt.Errorf("failed to delete credentials: %w", err)
	}

	PrintSuccess("Successfully logged out.")
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	keyring, err := config.NewKeyringStore()
	if err != nil {
		return fmt.Errorf("failed to access keyring: %w", err)
	}

	creds, err := keyring.GetCredentials()
	if err != nil {
		return fmt.Errorf("failed to get credentials: %w", err)
	}

	statusData := authStatusData{
		LoggedIn:    creds != nil,
		Environment: cfg.Environment(),
	}

	if creds != nil {
		statusData.APIKeyID = creds.APIKeyID

		client, err := createAuthenticatedClient(*creds)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			exchangeStatus, err := client.GetExchangeStatus(ctx)
			if err == nil {
				statusData.ExchangeActive = exchangeStatus.ExchangeActive
				statusData.TradingActive = exchangeStatus.TradingActive
				statusData.Authenticated = true
			}
		}
	}

	return ui.Output(
		outputFmt,
		func() { renderStatusTable(statusData) },
		statusData,
		func() { renderStatusPlain(statusData) },
	)
}

type authStatusData struct {
	LoggedIn       bool   `json:"logged_in"`
	APIKeyID       string `json:"api_key_id,omitempty"`
	Environment    string `json:"environment"`
	Authenticated  bool   `json:"authenticated"`
	ExchangeActive bool   `json:"exchange_active"`
	TradingActive  bool   `json:"trading_active"`
}

func renderStatusTable(data authStatusData) {
	var pairs [][]string

	loggedInStatus := ui.ErrorStyle.Render("No")
	if data.LoggedIn {
		loggedInStatus = ui.SuccessStyle.Render("Yes")
	}
	pairs = append(pairs, []string{"Logged In", loggedInStatus})

	if data.APIKeyID != "" {
		pairs = append(pairs, []string{"API Key ID", data.APIKeyID})
	}

	pairs = append(pairs, []string{"Environment", data.Environment})

	if data.LoggedIn {
		authStatus := ui.ErrorStyle.Render("Failed")
		if data.Authenticated {
			authStatus = ui.SuccessStyle.Render("Valid")
		}
		pairs = append(pairs, []string{"Authentication", authStatus})

		if data.Authenticated {
			exchangeStatus := "Inactive"
			if data.ExchangeActive {
				exchangeStatus = "Active"
			}
			pairs = append(pairs, []string{"Exchange", exchangeStatus})

			tradingStatus := "Inactive"
			if data.TradingActive {
				tradingStatus = "Active"
			}
			pairs = append(pairs, []string{"Trading", tradingStatus})
		}
	}

	ui.RenderKeyValue(pairs)
}

func renderStatusPlain(data authStatusData) {
	if data.LoggedIn {
		fmt.Printf("logged_in=true\n")
		fmt.Printf("api_key_id=%s\n", data.APIKeyID)
		fmt.Printf("environment=%s\n", data.Environment)
		fmt.Printf("authenticated=%v\n", data.Authenticated)
		if data.Authenticated {
			fmt.Printf("exchange_active=%v\n", data.ExchangeActive)
			fmt.Printf("trading_active=%v\n", data.TradingActive)
		}
	} else {
		fmt.Printf("logged_in=false\n")
		fmt.Printf("environment=%s\n", data.Environment)
	}
}

func runKeysList(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keys, err := client.ListAPIKeys(ctx)
	if err != nil {
		return fmt.Errorf("failed to list API keys: %w", err)
	}

	return ui.Output(
		outputFmt,
		func() { renderKeysTable(keys) },
		keys,
		func() { renderKeysPlain(keys) },
	)
}

func renderKeysTable(keys []api.APIKey) {
	headers := []string{"ID", "Name", "Created", "Expires", "Scopes"}
	var rows [][]string

	for _, key := range keys {
		expires := "-"
		if !key.ExpiresTime.IsZero() {
			expires = key.ExpiresTime.Format("2006-01-02")
		}
		scopes := strings.Join(key.Scopes, ", ")
		if scopes == "" {
			scopes = "-"
		}
		rows = append(rows, []string{
			key.ID,
			key.Name,
			key.CreatedTime.Format("2006-01-02"),
			expires,
			scopes,
		})
	}

	ui.RenderTable(headers, rows)
}

func renderKeysPlain(keys []api.APIKey) {
	for _, key := range keys {
		fmt.Printf("%s\t%s\t%s\n", key.ID, key.Name, key.CreatedTime.Format("2006-01-02"))
	}
}

func runKeysCreate(cmd *cobra.Command, args []string) error {
	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := api.CreateAPIKeyRequest{
		Name: keyName,
	}

	resp, err := client.CreateAPIKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	return ui.Output(
		outputFmt,
		func() { renderKeyCreatedTable(resp) },
		resp,
		func() { renderKeyCreatedPlain(resp) },
	)
}

func renderKeyCreatedTable(resp *api.CreateAPIKeyResponse) {
	PrintSuccess("API key created successfully!")
	fmt.Println()

	pairs := [][]string{
		{"ID", resp.APIKey.ID},
		{"Name", resp.APIKey.Name},
		{"Created", resp.APIKey.CreatedTime.Format("2006-01-02 15:04:05")},
	}
	ui.RenderKeyValue(pairs)

	fmt.Println()
	fmt.Println(ui.WarningStyle.Render("IMPORTANT: Save the private key below. It will not be shown again!"))
	fmt.Println()
	fmt.Println(resp.PrivateKey)
}

func renderKeyCreatedPlain(resp *api.CreateAPIKeyResponse) {
	fmt.Printf("id=%s\n", resp.APIKey.ID)
	fmt.Printf("name=%s\n", resp.APIKey.Name)
	fmt.Printf("private_key=%s\n", resp.PrivateKey)
}

func runKeysDelete(cmd *cobra.Command, args []string) error {
	keyID := args[0]

	if !SkipConfirmation() {
		fmt.Printf("Are you sure you want to delete API key '%s'? [y/N]: ", keyID)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Delete cancelled.")
			return nil
		}
	}

	client, err := getAuthenticatedClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.DeleteAPIKey(ctx, keyID); err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	result := map[string]interface{}{
		"deleted": true,
		"id":      keyID,
	}

	return ui.Output(
		outputFmt,
		func() { PrintSuccess(fmt.Sprintf("API key '%s' deleted successfully.", keyID)) },
		result,
		func() { fmt.Printf("deleted=%s\n", keyID) },
	)
}

func getAuthenticatedClient() (*api.Client, error) {
	keyring, err := config.NewKeyringStore()
	if err != nil {
		return nil, fmt.Errorf("failed to access keyring: %w", err)
	}

	creds, err := keyring.GetCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if creds == nil {
		return nil, fmt.Errorf("not logged in. Run 'kalshi-cli auth login' first")
	}

	return createAuthenticatedClient(*creds)
}

func createAuthenticatedClient(creds config.Credentials) (*api.Client, error) {
	signer, err := api.NewSignerFromPEM(creds.APIKeyID, creds.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	client := api.NewClient(cfg, signer)
	return client, nil
}
