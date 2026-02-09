# kalshi-cli

A comprehensive command-line interface for the [Kalshi](https://kalshi.com) prediction market exchange. Built for automated trading systems and bots.

**Designed for OpenClaw.ai bots** - All commands support `--json` output for machine parsing and `--yes` to skip confirmations.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Global Flags](#global-flags)
- [Authentication](#authentication)
- [Bot-Friendly Features](#bot-friendly-features)
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
- [Common Bot Workflows](#common-bot-workflows)
- [JSON Output Schemas](#json-output-schemas)
- [Error Handling](#error-handling)
- [Exit Codes](#exit-codes)

## Installation

### Go Install
```bash
go install github.com/6missedcalls/kalshi-cli/cmd/kalshi-cli@latest
```

### Homebrew (macOS/Linux)
```bash
brew install 6missedcalls/tap/kalshi-cli
```

### Binary Download
```bash
# macOS (Apple Silicon)
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_darwin_arm64.tar.gz
tar -xzf kalshi-cli_darwin_arm64.tar.gz && sudo mv kalshi-cli /usr/local/bin/

# Linux (x86_64)
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_linux_amd64.tar.gz
tar -xzf kalshi-cli_linux_amd64.tar.gz && sudo mv kalshi-cli /usr/local/bin/
```

## Quick Start

```bash
# 1. Authenticate (generates RSA key pair, stores in system keyring)
kalshi-cli auth login

# 2. Check exchange status
kalshi-cli exchange status --json

# 3. List open markets
kalshi-cli markets list --status open --json

# 4. Get your balance
kalshi-cli portfolio balance --json

# 5. Place an order (demo by default)
kalshi-cli orders create --market TICKER --side yes --qty 10 --price 50 --yes --json

# 6. Use production (real money)
kalshi-cli --prod orders create --market TICKER --side yes --qty 10 --price 50 --yes --json
```

## Global Flags

| Flag | Description | Bot Usage |
|------|-------------|-----------|
| `--json` | Output as JSON | **Always use for bots** |
| `--plain` | Plain text output | For piping to other tools |
| `--yes`, `-y` | Skip confirmation prompts | **Required for automation** |
| `--prod` | Use production API | Default is demo (safe) |
| `--verbose`, `-v` | Verbose output | Debugging |
| `--config` | Custom config file path | Multi-account setups |

**Bot Command Pattern:**
```bash
kalshi-cli --json --yes [command] [subcommand] [flags]
```

## Authentication

### Initial Setup
```bash
# Interactive login - generates 4096-bit RSA key pair
kalshi-cli auth login
# Follow prompts:
# 1. Copy displayed public key
# 2. Add to Kalshi dashboard (kalshi.com/account/api-keys)
# 3. Enter API Key ID when prompted
```

### Check Auth Status
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

### Credentials Storage
Credentials are stored in system keyring:
- **macOS**: Keychain
- **Linux**: Secret Service (GNOME Keyring)
- **Windows**: Windows Credential Manager

## Bot-Friendly Features

1. **JSON Output**: All commands support `--json` for structured parsing
2. **No Prompts**: Use `--yes` to skip all confirmations
3. **Exit Codes**: Non-zero on errors for scripting
4. **Idempotent**: Safe to retry failed commands
5. **Demo Default**: Safe testing without `--prod` flag

## Command Reference

### Markets

#### List Markets
```bash
kalshi-cli markets list [--status open|closed|settled] [--series TICKER] [--limit 50] --json
```

#### Get Market Details
```bash
kalshi-cli markets get TICKER --json
```

#### Get Orderbook
```bash
kalshi-cli markets orderbook TICKER --json
```
Returns bids and asks with price/quantity.

#### Get Trades
```bash
kalshi-cli markets trades TICKER [--limit 100] --json
```

#### Get Candlesticks
```bash
kalshi-cli markets candlesticks TICKER [--period 1m|5m|15m|1h|4h|1d] --json
```

#### List/Get Series
```bash
kalshi-cli markets series list [--category CATEGORY] --json
kalshi-cli markets series get TICKER --json
```

### Events

#### List Events
```bash
kalshi-cli events list [--status active] [--limit 50] --json
```

#### Get Event Details
```bash
kalshi-cli events get TICKER --json
```

#### Get Event Candlesticks
```bash
kalshi-cli events candlesticks TICKER [--period 1h] --json
```

#### Multivariate Events
```bash
kalshi-cli events multivariate list --json
kalshi-cli events multivariate get TICKER --json
```

### Orders

#### List Orders
```bash
kalshi-cli orders list [--status resting|executed|canceled] [--market TICKER] --json
```

#### Get Order
```bash
kalshi-cli orders get ORDER_ID --json
```

#### Create Order
```bash
kalshi-cli orders create \
  --market TICKER \
  --side yes|no \
  --qty 10 \
  --price 50 \
  [--action buy|sell] \
  [--type limit|market] \
  --yes --json
```

**Parameters:**
- `--market`: Market ticker (required)
- `--side`: `yes` or `no` (required)
- `--qty`: Number of contracts (required)
- `--price`: Price in cents, 1-99 (required for limit orders)
- `--action`: `buy` (default) or `sell`
- `--type`: `limit` (default) or `market`

#### Cancel Order
```bash
kalshi-cli orders cancel ORDER_ID --yes --json
```

#### Cancel All Orders
```bash
kalshi-cli orders cancel-all [--market TICKER] --yes --json
```

#### Amend Order
```bash
kalshi-cli orders amend ORDER_ID --qty 20 [--price 55] --yes --json
```

#### Batch Create Orders
```bash
kalshi-cli orders batch-create --file orders.json --yes --json
```

**orders.json format:**
```json
[
  {"ticker": "MARKET1", "side": "yes", "action": "buy", "type": "limit", "count": 10, "yes_price": 50},
  {"ticker": "MARKET2", "side": "no", "action": "buy", "type": "limit", "count": 5, "no_price": 30}
]
```

#### Get Queue Position
```bash
kalshi-cli orders queue ORDER_ID --json
```

### Portfolio

#### Get Balance
```bash
kalshi-cli portfolio balance --json
```
```json
{
  "available_balance": 10000,
  "portfolio_value": 5000,
  "total_balance": 15000
}
```
*Values in cents*

#### List Positions
```bash
kalshi-cli portfolio positions [--market TICKER] --json
```

#### List Fills
```bash
kalshi-cli portfolio fills [--limit 100] --json
```

#### List Settlements
```bash
kalshi-cli portfolio settlements [--limit 50] --json
```

#### Subaccounts
```bash
kalshi-cli portfolio subaccounts list --json
kalshi-cli portfolio subaccounts create --yes --json
kalshi-cli portfolio subaccounts transfer --from 0 --to 1 --amount 1000 --yes --json
```

### Order Groups

Order groups limit total fills across multiple orders.

```bash
# List groups
kalshi-cli order-groups list [--status active] --json
# Alias: kalshi-cli og list --json

# Get group
kalshi-cli order-groups get GROUP_ID --json

# Create group with 100 contract limit
kalshi-cli order-groups create --limit 100 --json

# Delete group
kalshi-cli order-groups delete GROUP_ID --yes --json

# Reset filled count
kalshi-cli order-groups reset GROUP_ID --json

# Trigger execution
kalshi-cli order-groups trigger GROUP_ID --json
```

### RFQ/Quotes

Request for Quotes - for block trading.

#### RFQ Commands
```bash
# List RFQs
kalshi-cli rfq list [--status open] --json

# Get RFQ
kalshi-cli rfq get RFQ_ID --json

# Create RFQ
kalshi-cli rfq create --market TICKER --side yes --qty 1000 --json

# Delete RFQ
kalshi-cli rfq delete RFQ_ID --yes --json
```

#### Quote Commands
```bash
# List quotes
kalshi-cli quotes list [--rfq-id RFQ_ID] --json

# Create quote
kalshi-cli quotes create --rfq RFQ_ID --price 50 --json

# Accept quote
kalshi-cli quotes accept QUOTE_ID --yes --json

# Confirm quote
kalshi-cli quotes confirm QUOTE_ID --yes --json
```

### Exchange

```bash
# Exchange status
kalshi-cli exchange status --json

# Trading schedule
kalshi-cli exchange schedule --json

# Announcements
kalshi-cli exchange announcements --json
```

### Watch (WebSocket)

Real-time streaming data. Output is newline-delimited JSON when using `--json`.

```bash
# Price updates for a market
kalshi-cli watch ticker MARKET_TICKER --json

# Orderbook updates
kalshi-cli watch orderbook MARKET_TICKER --json

# Public trades
kalshi-cli watch trades [--market TICKER] --json

# Your order updates (auth required)
kalshi-cli watch orders --json

# Your fill notifications (auth required)
kalshi-cli watch fills --json

# Your position changes (auth required)
kalshi-cli watch positions --json
```

### Config

```bash
# Show config
kalshi-cli config show --json

# Get value
kalshi-cli config get output.format

# Set value
kalshi-cli config set output.format json
kalshi-cli config set defaults.limit 100
```

## Common Bot Workflows

### 1. Market Making Bot
```bash
#!/bin/bash
MARKET="BTC-100K-2024"

# Get current orderbook
BOOK=$(kalshi-cli markets orderbook $MARKET --json)

# Extract best bid/ask
BEST_BID=$(echo $BOOK | jq '.yes_bids[0].price // 0')
BEST_ASK=$(echo $BOOK | jq '.yes_asks[0].price // 100')

# Place orders at spread
kalshi-cli orders create --market $MARKET --side yes --action buy --qty 10 --price $((BEST_BID + 1)) --yes --json
kalshi-cli orders create --market $MARKET --side yes --action sell --qty 10 --price $((BEST_ASK - 1)) --yes --json
```

### 2. Position Monitor
```bash
#!/bin/bash
# Check positions and P&L
kalshi-cli portfolio positions --json | jq '.[] | {ticker, position, pnl}'

# Check balance
kalshi-cli portfolio balance --json | jq '{available: .available_balance, total: .total_balance}'
```

### 3. Order Management
```bash
#!/bin/bash
# Cancel all resting orders for a market
kalshi-cli orders cancel-all --market BTC-100K --yes --json

# Check order status
ORDER_ID="order_abc123"
kalshi-cli orders get $ORDER_ID --json | jq '.status'
```

### 4. Real-time Price Feed
```bash
#!/bin/bash
# Stream prices and process with jq
kalshi-cli watch ticker BTC-100K --json | while read line; do
  PRICE=$(echo $line | jq '.yes_price')
  echo "Current price: $PRICE"
  # Add trading logic here
done
```

### 5. Multi-Market Scanner
```bash
#!/bin/bash
# Get all open markets and filter by volume
kalshi-cli markets list --status open --limit 200 --json | \
  jq '.[] | select(.volume > 10000) | {ticker, yes_bid, yes_ask, volume}'
```

## JSON Output Schemas

### Order Response
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

### Position Response
```json
{
  "ticker": "string",
  "position": 10,
  "market_exposure": 500,
  "realized_pnl": 100,
  "total_cost": 400
}
```

### Balance Response
```json
{
  "available_balance": 10000,
  "portfolio_value": 5000,
  "total_balance": 15000
}
```

### Market Response
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

Common error codes:
- `UNAUTHORIZED` - Invalid or expired credentials
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource not found
- `BAD_REQUEST` - Invalid parameters
- `RATE_LIMITED` - Too many requests (retry with backoff)
- `INSUFFICIENT_BALANCE` - Not enough funds

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication error |
| 3 | Validation error |
| 4 | API error |
| 5 | Network error |

## Environment Variables

```bash
export KALSHI_API_PRODUCTION=true      # Use production API
export KALSHI_API_TIMEOUT=60s          # Request timeout
export KALSHI_OUTPUT_FORMAT=json       # Default output format
```

## Rate Limits

Kalshi API has rate limits. The CLI implements automatic retry with exponential backoff:
- Initial delay: 100ms
- Max delay: 10s
- Max retries: 5

For high-frequency bots, implement your own rate limiting.

## Demo vs Production

| Environment | Flag | Base URL | Purpose |
|-------------|------|----------|---------|
| Demo | (default) | demo-api.kalshi.co | Testing, paper trading |
| Production | `--prod` | api.kalshi.com | Real money trading |

**Always test in demo first!**

## License

MIT License - see [LICENSE](LICENSE) for details.

---

Built for [OpenClaw.ai](https://openclaw.ai) trading bots.
