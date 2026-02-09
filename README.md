# kalshi-cli

A comprehensive command-line interface for the [Kalshi](https://kalshi.com) prediction market exchange. Trade event contracts, monitor positions, and stream real-time market data from your terminal.

Built for traders, developers, and automated trading systems. All commands support `--json` output for machine parsing and `--yes` to skip confirmations.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Global Flags](#global-flags)
- [Command Reference](#command-reference)
  - [Markets](#markets)
  - [Events](#events)
  - [Orders](#orders)
  - [Portfolio](#portfolio)
  - [Order Groups](#order-groups)
  - [RFQ/Quotes](#rfqquotes)
  - [Exchange](#exchange)
  - [Watch (WebSocket)](#watch-websocket)
  - [Config](#config)
- [Bot Integration](#bot-integration)
- [JSON Output Schemas](#json-output-schemas)
- [Error Handling](#error-handling)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Full Kalshi API coverage** - Markets, events, orders, portfolio, RFQs, order groups, and exchange info
- **Real-time WebSocket streaming** - Live price tickers, orderbook updates, trades, fills, and position changes
- **Secure authentication** - RSA-SHA256 signatures with credentials stored in your OS keyring (Keychain, Secret Service, Credential Manager)
- **Bot-friendly** - JSON output, no-confirmation mode, structured exit codes, idempotent operations
- **Demo-first** - Defaults to Kalshi's demo environment so you never accidentally trade real money
- **Cross-platform** - macOS (Intel + Apple Silicon), Linux, and Windows
- **Batch operations** - Create multiple orders from a JSON file in one call
- **Order groups** - Cap total fills across multiple orders
- **Block trading** - RFQ (Request for Quotes) workflow for large positions

## Installation

### Homebrew (macOS / Linux)

```bash
brew install 6missedcalls/tap/kalshi-cli
```

### Go Install

Requires Go 1.21+:

```bash
go install github.com/6missedcalls/kalshi-cli/cmd/kalshi-cli@latest
```

### Binary Download

Download the latest release for your platform:

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_darwin_arm64.tar.gz
tar -xzf kalshi-cli_darwin_arm64.tar.gz
sudo mv kalshi-cli /usr/local/bin/

# macOS (Intel)
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_darwin_amd64.tar.gz
tar -xzf kalshi-cli_darwin_amd64.tar.gz
sudo mv kalshi-cli /usr/local/bin/

# Linux (x86_64)
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_linux_amd64.tar.gz
tar -xzf kalshi-cli_linux_amd64.tar.gz
sudo mv kalshi-cli /usr/local/bin/

# Linux (ARM64)
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_linux_arm64.tar.gz
tar -xzf kalshi-cli_linux_arm64.tar.gz
sudo mv kalshi-cli /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/6missedcalls/kalshi-cli.git
cd kalshi-cli
go build -o kalshi-cli ./cmd/kalshi-cli
```

## Quick Start

```bash
# 1. Authenticate (interactive - generates RSA key pair, stores in system keyring)
kalshi-cli auth login

# 2. Check exchange status
kalshi-cli exchange status

# 3. Browse open markets
kalshi-cli markets list --status open

# 4. View your balance
kalshi-cli portfolio balance

# 5. Place an order on demo (default)
kalshi-cli orders create --market TICKER --side yes --qty 10 --price 50

# 6. Stream live prices
kalshi-cli watch ticker TICKER

# 7. When ready for production (real money)
kalshi-cli --prod orders create --market TICKER --side yes --qty 10 --price 50
```

## Authentication

### Interactive Login

The CLI generates a 4096-bit RSA key pair and stores credentials securely in your OS keyring:

```bash
kalshi-cli auth login
```

Follow the prompts to:
1. Copy the displayed public key
2. Add it to your Kalshi account at [kalshi.com/account/api-keys](https://kalshi.com/account/api-keys)
3. Enter the API Key ID when prompted

### Non-Interactive Login (Bots / CI)

For automated systems, pass credentials directly:

```bash
# Via flags
kalshi-cli auth login --api-key-id YOUR_KEY_ID --private-key-file /path/to/key.pem

# Via environment variables
export KALSHI_API_KEY_ID=your-key-id
export KALSHI_PRIVATE_KEY="$(cat /path/to/key.pem)"
kalshi-cli auth login
```

### Auth Status

```bash
kalshi-cli auth status --json
```

```json
{
  "logged_in": true,
  "api_key_id": "abc123...",
  "environment": "demo",
  "authenticated": true,
  "exchange_active": true,
  "trading_active": true
}
```

### API Key Management

```bash
kalshi-cli auth keys list          # List API keys
kalshi-cli auth keys create        # Create new API key
kalshi-cli auth keys delete KEY_ID # Delete API key
```

### Credential Storage

Credentials are stored in your OS keyring - never in plaintext files:

| OS | Backend |
|----|---------|
| macOS | Keychain |
| Linux | Secret Service (GNOME Keyring) |
| Windows | Credential Manager |

### Logout

```bash
kalshi-cli auth logout
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | | Output as JSON (recommended for scripts) |
| `--plain` | | Plain text output for piping |
| `--yes` | `-y` | Skip all confirmation prompts |
| `--prod` | | Use production API (default: demo) |
| `--verbose` | `-v` | Verbose output for debugging |
| `--config` | | Custom config file path |

**Standard command pattern for automation:**

```bash
kalshi-cli --json --yes [command] [subcommand] [flags]
```

## Command Reference

### Markets

```bash
# List markets with optional filters
kalshi-cli markets list [--status open|closed|settled] [--series TICKER] [--limit 50]

# Get details for a specific market
kalshi-cli markets get TICKER

# View the orderbook
kalshi-cli markets orderbook TICKER

# View recent trades
kalshi-cli markets trades TICKER [--limit 100]

# Get OHLCV candlestick data
kalshi-cli markets candlesticks TICKER [--period 1m|5m|15m|1h|4h|1d]

# List and view market series
kalshi-cli markets series list [--category CATEGORY]
kalshi-cli markets series get TICKER
```

### Events

```bash
# List events
kalshi-cli events list [--status active] [--limit 50]

# Get event details
kalshi-cli events get TICKER

# Event candlesticks
kalshi-cli events candlesticks TICKER [--period 1h]

# Multivariate events
kalshi-cli events multivariate list
kalshi-cli events multivariate get TICKER
```

### Orders

```bash
# List orders with filters
kalshi-cli orders list [--status resting|executed|canceled] [--market TICKER]

# Get a specific order
kalshi-cli orders get ORDER_ID

# Create an order
kalshi-cli orders create \
  --market TICKER \
  --side yes|no \
  --qty 10 \
  --price 50 \
  [--action buy|sell] \
  [--type limit|market]

# Amend an order (change price or quantity)
kalshi-cli orders amend ORDER_ID --qty 20 [--price 55]

# Cancel a single order
kalshi-cli orders cancel ORDER_ID

# Cancel all orders (optionally filtered by market)
kalshi-cli orders cancel-all [--market TICKER]

# Batch create from a JSON file
kalshi-cli orders batch-create --file orders.json

# Check queue position
kalshi-cli orders queue ORDER_ID
```

**Order parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `--market` | Yes | Market ticker |
| `--side` | Yes | `yes` or `no` |
| `--qty` | Yes | Number of contracts |
| `--price` | Limit only | Price in cents (1-99) |
| `--action` | No | `buy` (default) or `sell` |
| `--type` | No | `limit` (default) or `market` |

**Batch file format (`orders.json`):**

```json
[
  { "ticker": "MARKET1", "side": "yes", "action": "buy", "type": "limit", "count": 10, "yes_price": 50 },
  { "ticker": "MARKET2", "side": "no", "action": "buy", "type": "limit", "count": 5, "no_price": 30 }
]
```

### Portfolio

```bash
# Account balance (values in cents)
kalshi-cli portfolio balance

# List positions
kalshi-cli portfolio positions [--market TICKER]

# Trade fills
kalshi-cli portfolio fills [--limit 100]

# Settlements
kalshi-cli portfolio settlements [--limit 50]

# Subaccounts
kalshi-cli portfolio subaccounts list
kalshi-cli portfolio subaccounts create
kalshi-cli portfolio subaccounts transfer --from 0 --to 1 --amount 1000
```

### Order Groups

Order groups let you cap total fills across multiple orders.

```bash
# List groups (alias: og)
kalshi-cli order-groups list [--status active]
kalshi-cli og list

# Get group details
kalshi-cli order-groups get GROUP_ID

# Create with a fill limit
kalshi-cli order-groups create --limit 100

# Delete a group
kalshi-cli order-groups delete GROUP_ID

# Reset filled count
kalshi-cli order-groups reset GROUP_ID

# Trigger execution
kalshi-cli order-groups trigger GROUP_ID
```

### RFQ/Quotes

Request for Quotes workflow for block trading:

```bash
# RFQ commands
kalshi-cli rfq list [--status open]
kalshi-cli rfq get RFQ_ID
kalshi-cli rfq create --market TICKER --side yes --qty 1000
kalshi-cli rfq delete RFQ_ID

# Quote commands
kalshi-cli quotes list [--rfq-id RFQ_ID]
kalshi-cli quotes create --rfq RFQ_ID --price 50
kalshi-cli quotes accept QUOTE_ID
kalshi-cli quotes confirm QUOTE_ID
```

### Exchange

```bash
kalshi-cli exchange status         # Exchange status and trading hours
kalshi-cli exchange schedule       # Trading schedule
kalshi-cli exchange announcements  # Platform announcements
```

### Watch (WebSocket)

Stream real-time data via WebSocket. Output is newline-delimited JSON with `--json`.

```bash
# Public channels (no auth required)
kalshi-cli watch ticker TICKER      # Live price updates
kalshi-cli watch orderbook TICKER   # Orderbook changes
kalshi-cli watch trades             # Public trade feed

# Authenticated channels
kalshi-cli watch orders             # Your order updates
kalshi-cli watch fills              # Your fill notifications
kalshi-cli watch positions          # Your position changes
```

WebSocket features:
- Automatic reconnection with exponential backoff (1s-60s)
- Ping/pong keepalive (10-second intervals)
- Subscription persistence across reconnects
- Thread-safe concurrent access

### Config

```bash
kalshi-cli config show              # Display current configuration
kalshi-cli config get <key>         # Get a specific value
kalshi-cli config set <key> <value> # Set a value
```

Available config keys:

| Key | Default | Description |
|-----|---------|-------------|
| `api.production` | `false` | Use production environment |
| `api.timeout` | `30s` | HTTP request timeout |
| `output.format` | `table` | Output format: `table`, `json`, `plain` |
| `output.color` | `true` | Colorized terminal output |
| `defaults.limit` | `50` | Default result limit for list commands |

## Bot Integration

### Automation Flags

| Flag | Purpose |
|------|---------|
| `--json` | Machine-parseable structured output |
| `--yes` | Skip all interactive confirmations |
| `--plain` | Unformatted text for piping |
| `--prod` | Target production (omit for safe demo trading) |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication error |
| 3 | Validation error |
| 4 | API error |
| 5 | Network error |

### Example: Market Making Bot

```bash
#!/bin/bash
MARKET="BTC-100K-2025"

# Get current orderbook
BOOK=$(kalshi-cli markets orderbook $MARKET --json)

# Extract best bid/ask
BEST_BID=$(echo $BOOK | jq '.yes_bids[0].price // 0')
BEST_ASK=$(echo $BOOK | jq '.yes_asks[0].price // 100')

# Place orders at the spread
kalshi-cli orders create --market $MARKET --side yes --action buy  --qty 10 --price $((BEST_BID + 1)) --yes --json
kalshi-cli orders create --market $MARKET --side yes --action sell --qty 10 --price $((BEST_ASK - 1)) --yes --json
```

### Example: Position Monitor

```bash
#!/bin/bash
kalshi-cli portfolio positions --json | jq '.[] | {ticker, position, pnl}'
kalshi-cli portfolio balance --json | jq '{available: .available_balance, total: .total_balance}'
```

### Example: Real-Time Price Feed

```bash
#!/bin/bash
kalshi-cli watch ticker BTC-100K --json | while read line; do
  PRICE=$(echo $line | jq '.yes_price')
  echo "Current price: $PRICE"
done
```

### Example: Multi-Market Scanner

```bash
#!/bin/bash
kalshi-cli markets list --status open --limit 200 --json | \
  jq '.[] | select(.volume > 10000) | {ticker, yes_bid, yes_ask, volume}'
```

## JSON Output Schemas

### Order

```json
{
  "order_id": "string",
  "ticker": "string",
  "side": "yes|no",
  "action": "buy|sell",
  "type": "limit|market",
  "status": "resting|executed|canceled|pending",
  "yes_price": 50,
  "no_price": 50,
  "initial_quantity": 10,
  "remaining_quantity": 5,
  "created_time": "2024-01-01T00:00:00Z"
}
```

### Position

```json
{
  "ticker": "string",
  "position": 10,
  "market_exposure": 500,
  "realized_pnl": 100,
  "total_cost": 400
}
```

### Balance

```json
{
  "available_balance": 10000,
  "portfolio_value": 5000,
  "total_balance": 15000
}
```

All monetary values are in **cents**.

### Market

```json
{
  "ticker": "string",
  "title": "string",
  "status": "open|closed|settled",
  "yes_bid": 45,
  "yes_ask": 55,
  "last_price": 50,
  "volume": 10000
}
```

## Error Handling

Errors are returned as JSON when using `--json`:

```json
{
  "error": "API error [401] UNAUTHORIZED: Invalid API key",
  "code": "UNAUTHORIZED"
}
```

| Error Code | Description |
|------------|-------------|
| `UNAUTHORIZED` | Invalid or expired credentials |
| `FORBIDDEN` | Insufficient permissions |
| `NOT_FOUND` | Resource not found |
| `BAD_REQUEST` | Invalid parameters |
| `RATE_LIMITED` | Too many requests (auto-retried) |
| `INSUFFICIENT_BALANCE` | Not enough funds |

The CLI automatically retries rate-limited and transient server errors with exponential backoff (100ms initial, 10s max, 5 retries).

## Configuration

Configuration is stored at `~/.kalshi/config.yaml` (created automatically on first run).

```yaml
api:
  production: false
  timeout: 30s
output:
  format: table
  color: true
defaults:
  limit: 50
```

### Environment Variables

All config values can be overridden with environment variables:

```bash
export KALSHI_API_PRODUCTION=true
export KALSHI_API_TIMEOUT=60s
export KALSHI_OUTPUT_FORMAT=json
export KALSHI_API_KEY_ID=your-key-id
export KALSHI_PRIVATE_KEY="$(cat /path/to/key.pem)"
```

### Demo vs Production

| | Demo | Production |
|--|------|-----------|
| **Flag** | (default) | `--prod` |
| **API** | `demo-api.kalshi.co` | `api.elections.kalshi.com` |
| **WebSocket** | `wss://demo-api.kalshi.co/trade-api/ws/v2` | `wss://api.elections.kalshi.com/trade-api/ws/v2` |
| **Purpose** | Testing, paper trading | Real money trading |

## Architecture

```
kalshi-cli/
├── cmd/kalshi-cli/        # Entry point
│   └── main.go
├── internal/
│   ├── api/               # HTTP client, auth signing, all API methods
│   ├── cmd/               # Cobra command definitions
│   ├── config/            # Viper config + keyring credential store
│   ├── ui/                # Table and output formatting
│   └── websocket/         # WebSocket client, channels, message routing
├── pkg/
│   └── models/            # Shared request/response types
├── .goreleaser.yml        # Cross-platform release builds
├── go.mod
└── go.sum
```

**Key design decisions:**
- **Cobra + Viper** for CLI framework and configuration
- **Resty** HTTP client with automatic retry and rate-limit handling
- **nhooyr.io/websocket** for WebSocket streaming with auto-reconnect
- **OS keyring** for credential storage (never plaintext)
- **RSA-SHA256 signatures** for API authentication (keys never leave your machine)
- **Demo-first** - production requires explicit `--prod` flag

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Run tests: `go test ./...`
5. Run vet: `go vet ./...`
6. Commit with conventional commits (`feat:`, `fix:`, `refactor:`, etc.)
7. Open a pull request

### Building

```bash
go build -o kalshi-cli ./cmd/kalshi-cli
```

### Testing

```bash
go test ./...
go test -cover ./...
```

## License

MIT License - see [LICENSE](LICENSE) for details.
