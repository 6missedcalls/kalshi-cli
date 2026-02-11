package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/6missedcalls/kalshi-cli/internal/ui"
)

var (
	cfgFile   string
	useProd   bool
	jsonOut   bool
	plainOut  bool
	yesFlag   bool
	verbose   bool
	cfg       *config.Config
	outputFmt ui.OutputFormat

	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "kalshi-cli",
	Short: "CLI for the Kalshi prediction market exchange",
	Long: `kalshi-cli is a comprehensive command-line interface for the Kalshi
prediction market exchange. It provides access to all API endpoints,
real-time WebSocket streaming, and a first-class trading experience.

By default, commands use the demo API. Use --prod for production.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kalshi/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&useProd, "prod", false, "use production API (default: demo)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVar(&plainOut, "plain", false, "output as plain text (for pipes)")
	rootCmd.PersistentFlags().BoolVarP(&yesFlag, "yes", "y", false, "skip confirmation prompts")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	viper.BindPFlag("api.production", rootCmd.PersistentFlags().Lookup("prod"))
	viper.BindPFlag("output.json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("output.plain", rootCmd.PersistentFlags().Lookup("plain"))

	rootCmd.AddCommand(versionCmd)
}

func initConfig() error {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if useProd {
		cfg.API.Production = true
	}

	switch {
	case jsonOut:
		outputFmt = ui.FormatJSON
	case plainOut:
		outputFmt = ui.FormatPlain
	default:
		outputFmt = ui.FormatTable
	}

	return nil
}

func GetConfig() *config.Config {
	return cfg
}

func GetOutputFormat() ui.OutputFormat {
	return outputFmt
}

func IsVerbose() bool {
	return verbose
}

func SkipConfirmation() bool {
	return yesFlag
}

func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "%s %s\n", ui.ErrorStyle.Render("Error:"), err.Error())
}

func PrintSuccess(msg string) {
	fmt.Println(ui.SuccessStyle.Render(msg))
}

func PrintWarning(msg string) {
	fmt.Println(ui.WarningStyle.Render(msg))
}

// SetVersionInfo is called from main to inject build-time variables.
func SetVersionInfo(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date
	rootCmd.Version = version
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kalshi-cli %s\n", buildVersion)
		fmt.Printf("  commit:  %s\n", buildCommit)
		fmt.Printf("  built:   %s\n", buildDate)
	},
}
