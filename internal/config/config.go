package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const (
	// Base URLs WITHOUT the /trade-api/v2 prefix (added by API methods)
	DemoBaseURL = "https://demo-api.kalshi.co"
	ProdBaseURL = "https://api.elections.kalshi.com"

	// WebSocket URLs (these need the full path)
	DemoWSURL = "wss://demo-api.kalshi.co/trade-api/ws/v2"
	ProdWSURL = "wss://api.elections.kalshi.com/trade-api/ws/v2"
)

type Config struct {
	API     APIConfig     `mapstructure:"api"`
	Output  OutputConfig  `mapstructure:"output"`
	Defaults DefaultsConfig `mapstructure:"defaults"`
}

type APIConfig struct {
	Production bool          `mapstructure:"production"`
	Timeout    time.Duration `mapstructure:"timeout"`
}

type OutputConfig struct {
	Format string `mapstructure:"format"`
	Color  bool   `mapstructure:"color"`
}

type DefaultsConfig struct {
	Limit int `mapstructure:"limit"`
}

func (c *Config) BaseURL() string {
	if c.API.Production {
		return ProdBaseURL
	}
	return DemoBaseURL
}

func (c *Config) WebSocketURL() string {
	if c.API.Production {
		return ProdWSURL
	}
	return DemoWSURL
}

func (c *Config) Environment() string {
	if c.API.Production {
		return "production"
	}
	return "demo"
}

func Load(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(home, ".kalshi")
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	setDefaults()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("KALSHI")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("api.production", false)
	viper.SetDefault("api.timeout", 30*time.Second)
	viper.SetDefault("output.format", "table")
	viper.SetDefault("output.color", true)
	viper.SetDefault("defaults.limit", 50)
}

func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".kalshi", "config.yaml")

	viper.Set("api.production", cfg.API.Production)
	viper.Set("api.timeout", cfg.API.Timeout)
	viper.Set("output.format", cfg.Output.Format)
	viper.Set("output.color", cfg.Output.Color)
	viper.Set("defaults.limit", cfg.Defaults.Limit)

	return viper.WriteConfigAs(configPath)
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kalshi"), nil
}
