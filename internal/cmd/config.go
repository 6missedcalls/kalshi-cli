package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
)

var validConfigKeys = map[string]struct {
	description string
	validate    func(string) error
}{
	"output.format": {
		description: "Output format (table, json, plain)",
		validate:    validateOutputFormat,
	},
	"output.color": {
		description: "Enable colored output (true, false)",
		validate:    validateBool,
	},
	"defaults.limit": {
		description: "Default limit for list commands (number)",
		validate:    validatePositiveInt,
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Manage kalshi-cli configuration settings.

Configuration is stored in ~/.kalshi/config.yaml.

Available configuration keys:
  output.format   Output format (table, json, plain)
  output.color    Enable colored output (true, false)
  defaults.limit  Default limit for list commands (number)`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display all current configuration settings.`,
	RunE:  runConfigShow,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long: `Get the value of a specific configuration key.

Available keys:
  output.format   Output format (table, json, plain)
  output.color    Enable colored output (true, false)
  defaults.limit  Default limit for list commands`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys and values:
  output.format   table, json, plain
  output.color    true, false
  defaults.limit  Any positive integer`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	currentConfig := GetConfig()

	configData := map[string]interface{}{
		"output.format":  currentConfig.Output.Format,
		"output.color":   currentConfig.Output.Color,
		"defaults.limit": currentConfig.Defaults.Limit,
	}

	configPath, err := config.ConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() {
			renderConfigTable(configData, configPath)
		},
		configData,
		func() {
			printConfigPlain(configData)
		},
	)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	if _, valid := validConfigKeys[key]; !valid {
		return fmt.Errorf("unknown configuration key: %s\n\nValid keys: %s", key, getValidKeysList())
	}

	currentConfig := GetConfig()
	value := getConfigValue(currentConfig, key)

	return ui.Output(
		GetOutputFormat(),
		func() {
			ui.RenderKeyValue([][]string{
				{key, fmt.Sprintf("%v", value)},
			})
		},
		map[string]interface{}{key: value},
		func() {
			ui.PrintPlain("%v", value)
		},
	)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	keyConfig, valid := validConfigKeys[key]
	if !valid {
		return fmt.Errorf("unknown configuration key: %s\n\nValid keys: %s", key, getValidKeysList())
	}

	if err := keyConfig.validate(value); err != nil {
		return fmt.Errorf("invalid value for %s: %w", key, err)
	}

	currentConfig := GetConfig()
	updatedConfig := applyConfigValue(currentConfig, key, value)

	if err := config.Save(updatedConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return ui.Output(
		GetOutputFormat(),
		func() {
			PrintSuccess(fmt.Sprintf("Set %s = %s", key, value))
		},
		map[string]interface{}{
			"key":   key,
			"value": value,
		},
		func() {
			ui.PrintPlain("%s=%s", key, value)
		},
	)
}

func validateOutputFormat(value string) error {
	validFormats := []string{"table", "json", "plain"}
	for _, format := range validFormats {
		if value == format {
			return nil
		}
	}
	return fmt.Errorf("must be one of: %s", strings.Join(validFormats, ", "))
}

func validateBool(value string) error {
	if value == "true" || value == "false" {
		return nil
	}
	return fmt.Errorf("must be true or false")
}

func validatePositiveInt(value string) error {
	n, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("must be a valid number")
	}
	if n <= 0 {
		return fmt.Errorf("must be a positive number")
	}
	return nil
}

func getValidKeysList() string {
	keys := make([]string, 0, len(validConfigKeys))
	for key := range validConfigKeys {
		keys = append(keys, key)
	}
	return strings.Join(keys, ", ")
}

func getConfigValue(cfg *config.Config, key string) interface{} {
	switch key {
	case "output.format":
		return cfg.Output.Format
	case "output.color":
		return cfg.Output.Color
	case "defaults.limit":
		return cfg.Defaults.Limit
	default:
		return viper.Get(key)
	}
}

func applyConfigValue(cfg *config.Config, key string, value string) *config.Config {
	return &config.Config{
		API: cfg.API,
		Output: applyOutputConfigValue(cfg.Output, key, value),
		Defaults: applyDefaultsConfigValue(cfg.Defaults, key, value),
	}
}

func applyOutputConfigValue(output config.OutputConfig, key string, value string) config.OutputConfig {
	switch key {
	case "output.format":
		return config.OutputConfig{
			Format: value,
			Color:  output.Color,
		}
	case "output.color":
		return config.OutputConfig{
			Format: output.Format,
			Color:  value == "true",
		}
	default:
		return output
	}
}

func applyDefaultsConfigValue(defaults config.DefaultsConfig, key string, value string) config.DefaultsConfig {
	switch key {
	case "defaults.limit":
		limit, _ := strconv.Atoi(value)
		return config.DefaultsConfig{
			Limit: limit,
		}
	default:
		return defaults
	}
}

func renderConfigTable(configData map[string]interface{}, configPath string) {
	fmt.Printf("Configuration file: %s/config.yaml\n\n", configPath)

	rows := [][]string{
		{"output.format", fmt.Sprintf("%v", configData["output.format"]), validConfigKeys["output.format"].description},
		{"output.color", fmt.Sprintf("%v", configData["output.color"]), validConfigKeys["output.color"].description},
		{"defaults.limit", fmt.Sprintf("%v", configData["defaults.limit"]), validConfigKeys["defaults.limit"].description},
	}

	ui.RenderTable([]string{"Key", "Value", "Description"}, rows)
}

func printConfigPlain(configData map[string]interface{}) {
	ui.PrintPlain("output.format=%v", configData["output.format"])
	ui.PrintPlain("output.color=%v", configData["output.color"])
	ui.PrintPlain("defaults.limit=%v", configData["defaults.limit"])
}
