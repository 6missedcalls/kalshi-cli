# kalshi-cli

A command-line interface for the [Kalshi](https://kalshi.com) prediction market exchange. Trade event contracts, monitor positions, stream real-time market data, and view ASCII candlestick charts from your terminal.

All commands support `--json` output for machine parsing, `--plain` for piping, and `--yes` to skip confirmations. Defaults to the demo API so you never accidentally trade real money.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Global Flags](#global-flags)
- [Commands](#commands)
  - [auth](#auth)
  - [markets](#markets)
  - [events](#events)
  - [orders](#orders)
  - [portfolio](#portfolio)
  - [order-groups](#order-groups)
  - [rfq](#rfq)
  - [quotes](#quotes)
  - [exchange](#exchange)
  - [watch](#watch)
  - [config](#config)
  - [version](#version)
  - [completion](#completion)
- [Configuration](#configuration)
- [Bot Integration](#bot-integration)
- [Architecture](#architecture)
- [License](#license)

## Installation

### Homebrew (macOS / Linux)

```bash
brew install 6missedcalls/tap/kalshi-cli
```

### Go Install

Requires Go 1.25+:

```bash
go install github.com/6missedcalls/kalshi-cli/cmd/kalshi-cli@latest
```

### Build from Source

```bash
git clone https://github.com/6missedcalls/kalshi-cli.git
cd kalshi-cli
go build -o kalshi-cli ./cmd/kalshi-cli
```

## Quick Start

```bash
# 1. Authenticate
kalshi-cli auth login

# 2. Check exchange status
kalshi-cli exchange status

# 3. Browse markets
kalshi-cli markets list --status open --limit 20

# 4. View your balance
kalshi-cli portfolio balance

# 5. View an orderbook
kalshi-cli markets orderbook KXBTC-26FEB12-B97000

# 6. View candlestick chart for an event
kalshi-cli events candlesticks KXINXU-26FEB11H1600 --series KXINXU --period 1h \
  --start 2026-02-10T00:00:00Z --end 2026-02-11T23:00:00Z

# 7. Place an order (demo)
kalshi-cli orders create --market KXBTC-26FEB12-B97000 --side yes --qty 10 --price 50

# 8. Stream live prices
kalshi-cli watch ticker KXBTC-26FEB12-B97000

# 9. When ready for production (real money)
kalshi-cli --prod orders create --market KXBTC-26FEB12-B97000 --side yes --qty 10 --price 50
```

## Authentication

### Interactive Login

```bash
kalshi-cli auth login
```

Follow the prompts to:
1. Copy the displayed public key
2. Add it to your Kalshi account at [kalshi.com/account/api-keys](https://kalshi.com/account/api-keys)
3. Enter the API Key ID when prompted

### Non-Interactive Login (Bots / CI)

```bash
# Via flags
kalshi-cli auth login --api-key-id YOUR_KEY_ID --private-key-file /path/to/key.pem

# Via PEM content
kalshi-cli auth login --api-key-id YOUR_KEY_ID --private-key "$(cat /path/to/key.pem)"

# Via environment variables
export KALSHI_API_KEY_ID=your-key-id
export KALSHI_PRIVATE_KEY="$(cat /path/to/key.pem)"
kalshi-cli auth login
```

### Config File Credentials

Add to `~/.kalshi/config.yaml`:

```yaml
api_key_id: your-key-id
private_key_path: /path/to/key.pem
```

Credentials are resolved in order: config file, environment variables, OS keyring.

### Credential Storage

| OS | Backend |
|----|---------|
| macOS | Keychain |
| Linux | Secret Service (GNOME Keyring) |
| Windows | Credential Manager |

## Global Flags

Every command accepts these flags:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--json` | | `false` | Output as JSON (for scripts and automation) |
| `--plain` | | `false` | Plain text output (for piping) |
| `--yes` | `-y` | `false` | Skip all confirmation prompts |
| `--prod` | | `false` | Use production API (default: demo) |
| `--verbose` | `-v` | `false` | Verbose output for debugging |
| `--config` | | `~/.kalshi/config.yaml` | Path to config file |

## Commands

---

### auth

Manage authentication credentials and API keys.

#### `auth login`

Authenticate with Kalshi using API credentials.

```
kalshi-cli auth login [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--api-key-id` | No | | API Key ID (or set `KALSHI_API_KEY_ID` env var) |
| `--private-key` | No | | Private key PEM content (or set `KALSHI_PRIVATE_KEY` env var) |
| `--private-key-file` | No | | Path to private key PEM file |

If no flags are provided, runs in interactive mode.

#### `auth logout`

Remove stored API credentials from the system keyring.

```
kalshi-cli auth logout
```

No additional flags.

#### `auth status`

Display the current authentication status and environment.

```
kalshi-cli auth status
```

No additional flags. Use `--json` to get machine-readable output:

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

#### `auth keys`

Manage API keys for your Kalshi account.

#### `auth keys list`

List all API keys associated with your account.

```
kalshi-cli auth keys list
```

No additional flags.

#### `auth keys create`

Create a new API key.

```
kalshi-cli auth keys create [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--name` | No | | Name for the new API key |

#### `auth keys delete`

Delete an API key by its ID.

```
kalshi-cli auth keys delete <id>
```

Positional argument: the API key ID to delete.

---

### markets

Commands for listing, viewing, and analyzing prediction markets.

#### `markets list`

List markets with optional filtering.

```
kalshi-cli markets list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--status` | No | | Filter by status: `open`, `closed`, `settled` |
| `--series` | No | | Filter by series ticker |
| `--limit` | No | `50` | Maximum number of markets to return |

```bash
kalshi-cli markets list --status open --limit 20
kalshi-cli markets list --series KXBTC --json
```

#### `markets get`

Get detailed information about a specific market.

```
kalshi-cli markets get <market-ticker>
```

Positional argument: the market ticker.

```bash
kalshi-cli markets get KXBTC-26FEB12-B97000
```

#### `markets orderbook`

Get the orderbook for a market with visual bid/ask display.

```
kalshi-cli markets orderbook <market-ticker>
```

Positional argument: the market ticker.

```bash
kalshi-cli markets orderbook KXBTC-26FEB12-B97000
kalshi-cli markets orderbook KXBTC-26FEB12-B97000 --json
```

#### `markets trades`

Get recent trades for a market.

```
kalshi-cli markets trades <market-ticker> [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--limit` | No | `100` | Maximum number of trades to return |

```bash
kalshi-cli markets trades KXBTC-26FEB12-B97000 --limit 20
```

#### `markets candlesticks`

Get candlestick (OHLCV) data for a market. Displays an ASCII candlestick chart above a data table.

```
kalshi-cli markets candlesticks <market-ticker> [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--series` | **Yes** | | Series ticker (e.g., `KXBTC`) |
| `--period` | No | `1h` | Candlestick period: `1m`, `1h`, `1d` |

```bash
kalshi-cli markets candlesticks KXBTC-26FEB12-B97000 --series KXBTC
kalshi-cli markets candlesticks KXBTC-26FEB12-B97000 --series KXBTC --period 1d
```

#### `markets series list`

List market series with optional category filtering.

```
kalshi-cli markets series list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--category` | No | | Filter by category (e.g., `Economics`, `Crypto`, `Politics`) |
| `--limit` | No | `50` | Maximum number of series to return |

```bash
kalshi-cli markets series list --category Economics
```

#### `markets series get`

Get details for a specific series.

```
kalshi-cli markets series get <series-ticker>
```

Positional argument: the series ticker.

```bash
kalshi-cli markets series get KXBTC
```

---

### events

Commands for listing, viewing, and managing events. An event groups related markets (e.g., "Bitcoin price range on Feb 12" has multiple strike-bracket markets under it).

#### `events list`

List events with optional filtering.

```
kalshi-cli events list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--status` | No | | Filter by status: `active`, `closed`, `settled` |
| `--limit` | No | `50` | Maximum number of events to return |
| `--cursor` | No | | Pagination cursor from a previous response |

```bash
kalshi-cli events list --limit 20
```

#### `events get`

Get detailed information about a specific event.

```
kalshi-cli events get <event-ticker>
```

Positional argument: the event ticker.

```bash
kalshi-cli events get KXBTC-26FEB12
```

#### `events candlesticks`

Get candlestick (OHLCV) data for an event across all its markets. Displays an ASCII candlestick chart above a data table.

```
kalshi-cli events candlesticks <event-ticker> [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--series` | No | Auto-resolved from event | Series ticker |
| `--period` | No | `1h` | Candlestick period: `1m`, `1h`, `1d` |
| `--start` | No | | Start time in RFC3339 format |
| `--end` | No | | End time in RFC3339 format |

```bash
kalshi-cli events candlesticks KXINXU-26FEB11H1600 \
  --start 2026-02-10T00:00:00Z --end 2026-02-11T23:00:00Z

kalshi-cli events candlesticks KXINXU-26FEB11H1600 --period 1d \
  --series KXINXU --start 2026-02-01T00:00:00Z --end 2026-02-11T00:00:00Z
```

The chart output looks like:

```
  Event Candlesticks  Last: $0.03  -$0.28 (-90.3%)

    $0.99 │              │
          │              ┃
          │              ┃
          │    │         ┃
    $0.73 │  │ ┃         ┃
          │  ┃ ┃         ┃
          │  ┃ │         ┃
          │  ┃ │         ┃     ┃
    $0.47 │  ┃   ┃       ┃     ┃ ─   ─ ─ ─       ┃ │
          │  ┃   ┃       ┃   ─ ┃                 ┃ ┃
          │─ ┃   ┃       ┃                       ┃ ┃
          ││ ┃ │   ┃ │ │ ┃                       ┃ ┃
    $0.20 │      ┃ ┃ │   ┃                     ─ ┃ ┃
          │      │ ┃ ┃ ┃ ┃         ─       ─ ─ │     ┃
          │          ┃ ┃ ┃ │                         ┃
    $0.00 │            │ │ ─                             ─
          └──────────────────────────────────────────────
           02/11 11:00         02/11 15:00       02/11 21:00
  Vol     ▁ ▂ ▁ ▁ ▁ █ ▂ ▂ ▄ ▁ ▂ ▁ ▁ ▁ ▁ ▁ ▁ ▁ ▁ ▆ ▇ ▅
```

- Green `┃` = bullish candle (close >= open)
- Red `┃` = bearish candle (close < open)
- Gray `│` = wicks (high/low beyond body)
- `─` = doji (open == close at same row)
- Bottom row: volume sparkline colored per candle direction

#### `events multivariate list`

List multivariate events.

```
kalshi-cli events multivariate list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--status` | No | | Filter by status |
| `--limit` | No | `50` | Maximum number of events to return |
| `--cursor` | No | | Pagination cursor |

#### `events multivariate get`

Get details for a multivariate event.

```
kalshi-cli events multivariate get <ticker>
```

Positional argument: the multivariate event ticker.

---

### orders

Manage trading orders on the Kalshi exchange.

#### `orders list`

List orders with optional filters.

```
kalshi-cli orders list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--status` | No | | Filter by status: `resting`, `canceled`, `executed`, `pending` |
| `--market` | No | | Filter by market ticker |

```bash
kalshi-cli orders list --status resting
kalshi-cli orders list --market KXBTC-26FEB12-B97000 --json
```

#### `orders get`

Get details for a specific order.

```
kalshi-cli orders get <order-id>
```

Positional argument: the order ID.

#### `orders create`

Create a new order. Shows a preview before submission (skip with `--yes`).

```
kalshi-cli orders create [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--market` | **Yes** | | Market ticker |
| `--side` | **Yes** | | `yes` or `no` |
| `--qty` | **Yes** | | Number of contracts |
| `--price` | **Yes** (limit) | | Price in cents (1-99) |
| `--action` | No | `buy` | `buy` or `sell` |
| `--type` | No | `limit` | `limit` or `market` |

```bash
kalshi-cli orders create --market KXBTC-26FEB12-B97000 --side yes --qty 10 --price 50
kalshi-cli orders create --market KXBTC-26FEB12-B97000 --side no --qty 5 --price 30 --action sell
kalshi-cli orders create --market KXBTC-26FEB12-B97000 --side yes --qty 10 --price 50 --yes --json
```

#### `orders amend`

Amend an existing order's quantity and/or price. At least one of `--qty` or `--price` must be specified.

```
kalshi-cli orders amend <order-id> [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--qty` | No | | New quantity |
| `--price` | No | | New price in cents |

#### `orders cancel`

Cancel a resting order by its ID.

```
kalshi-cli orders cancel <order-id>
```

#### `orders cancel-all`

Cancel all resting orders. Optionally filter by market.

```
kalshi-cli orders cancel-all [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--market` | No | | Only cancel orders for this market ticker |

#### `orders batch-create`

Create multiple orders from a JSON file.

```
kalshi-cli orders batch-create [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--file` | **Yes** | | Path to JSON file containing orders |

The JSON file should contain an array of order objects:

```json
[
  { "ticker": "MARKET1", "side": "yes", "action": "buy", "type": "limit", "count": 10, "yes_price": 50 },
  { "ticker": "MARKET2", "side": "no", "action": "buy", "type": "limit", "count": 5, "no_price": 30 }
]
```

#### `orders queue`

Get the queue position for a resting order.

```
kalshi-cli orders queue <order-id>
```

---

### portfolio

View and manage your Kalshi portfolio.

#### `portfolio balance`

Display account balance. All values are in cents.

```
kalshi-cli portfolio balance
```

No additional flags.

#### `portfolio positions`

List current market positions.

```
kalshi-cli portfolio positions [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--market` | No | | Filter by market ticker |

#### `portfolio fills`

List trade fills.

```
kalshi-cli portfolio fills [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--limit` | No | `100` | Maximum number of fills to return |

#### `portfolio settlements`

List market settlements.

```
kalshi-cli portfolio settlements [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--limit` | No | `50` | Maximum number of settlements to return |

#### `portfolio subaccounts list`

List all subaccounts.

```
kalshi-cli portfolio subaccounts list
```

No additional flags.

#### `portfolio subaccounts create`

Create a new subaccount.

```
kalshi-cli portfolio subaccounts create
```

No additional flags.

#### `portfolio subaccounts transfer`

Transfer funds between subaccounts.

```
kalshi-cli portfolio subaccounts transfer [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--from` | **Yes** | | Source subaccount ID |
| `--to` | **Yes** | | Destination subaccount ID |
| `--amount` | **Yes** | | Amount to transfer in cents |

---

### order-groups

Order groups cap total fills across multiple orders. Alias: `og`.

#### `order-groups list`

List all order groups.

```
kalshi-cli order-groups list [flags]
kalshi-cli og list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--status` | No | | Filter by status |

#### `order-groups get`

Get details for an order group.

```
kalshi-cli order-groups get <group-id>
```

#### `order-groups create`

Create a new order group with a contract limit.

```
kalshi-cli order-groups create [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--limit` | **Yes** | | Maximum contracts to fill across all orders in the group |

#### `order-groups delete`

Delete an order group. All orders in the group will be canceled.

```
kalshi-cli order-groups delete <group-id>
```

#### `order-groups reset`

Reset an order group's filled count to zero.

```
kalshi-cli order-groups reset <group-id>
```

#### `order-groups trigger`

Trigger an order group to execute its orders.

```
kalshi-cli order-groups trigger <group-id>
```

#### `order-groups update-limit`

Update the contract limit for an order group. If the new limit is lower than the current filled count, the group is triggered.

```
kalshi-cli order-groups update-limit <group-id> [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--limit` | **Yes** | | New maximum contracts to fill |

---

### rfq

Request for Quotes for block trading.

#### `rfq list`

List all RFQs.

```
kalshi-cli rfq list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--status` | No | | Filter by status (e.g., `open`, `closed`) |

#### `rfq get`

Get details for a specific RFQ.

```
kalshi-cli rfq get <rfq-id>
```

#### `rfq create`

Create a new RFQ.

```
kalshi-cli rfq create [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--market` | **Yes** | | Market ticker |
| `--qty` | **Yes** | | Quantity |

```bash
kalshi-cli rfq create --market KXBTC-26FEB12-B97000 --qty 1000
```

#### `rfq delete`

Delete an RFQ.

```
kalshi-cli rfq delete <rfq-id>
```

---

### quotes

Manage quotes on RFQs.

#### `quotes list`

List all quotes, optionally filtered by RFQ.

```
kalshi-cli quotes list [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--rfq-id` | No | | Filter by RFQ ID |

#### `quotes create`

Create a quote on an existing RFQ.

```
kalshi-cli quotes create [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--rfq` | **Yes** | | RFQ ID |
| `--price` | **Yes** | | Price in cents |

```bash
kalshi-cli quotes create --rfq rfq_abc123 --price 65
```

#### `quotes accept`

Accept a quote offered on your RFQ.

```
kalshi-cli quotes accept <quote-id>
```

#### `quotes confirm`

Confirm an accepted quote.

```
kalshi-cli quotes confirm <quote-id>
```

---

### exchange

Get exchange status, schedule, and announcements.

#### `exchange status`

Get current exchange status including trading activity and environment.

```
kalshi-cli exchange status
```

No additional flags.

#### `exchange schedule`

Get the exchange trading schedule.

```
kalshi-cli exchange schedule
```

No additional flags.

#### `exchange announcements`

Get the latest exchange announcements.

```
kalshi-cli exchange announcements
```

No additional flags.

---

### watch

Stream real-time data via WebSocket. All watch commands require authentication. Press `Ctrl+C` to stop. Use `--json` for newline-delimited JSON output.

Features:
- Automatic reconnection with exponential backoff (1s-60s)
- Ping/pong keepalive (10-second intervals)
- Subscription persistence across reconnects

#### `watch ticker`

Stream live price updates for a market.

```
kalshi-cli watch ticker <market-ticker>
```

Positional argument: the market ticker.

```bash
kalshi-cli watch ticker KXBTC-26FEB12-B97000
kalshi-cli watch ticker KXBTC-26FEB12-B97000 --json
```

#### `watch orderbook`

Stream orderbook delta updates for a market.

```
kalshi-cli watch orderbook <market-ticker>
```

Positional argument: the market ticker.

#### `watch trades`

Stream public trades. Optionally filter to a single market.

```
kalshi-cli watch trades [flags]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--market` | No | | Filter trades by market ticker |

```bash
kalshi-cli watch trades
kalshi-cli watch trades --market KXBTC-26FEB12-B97000 --json
```

#### `watch orders`

Stream your order status changes.

```
kalshi-cli watch orders
```

No additional flags.

#### `watch fills`

Stream your fill notifications.

```
kalshi-cli watch fills
```

No additional flags.

#### `watch positions`

Stream your position changes.

```
kalshi-cli watch positions
```

No additional flags.

---

### config

Manage configuration settings stored in `~/.kalshi/config.yaml`.

#### `config show`

Display all current configuration settings.

```
kalshi-cli config show
```

#### `config get`

Get the value of a specific configuration key.

```
kalshi-cli config get <key>
```

Available keys:

| Key | Default | Description |
|-----|---------|-------------|
| `output.format` | `table` | Output format: `table`, `json`, `plain` |
| `output.color` | `true` | Enable colored output |
| `defaults.limit` | `50` | Default result limit for list commands |

#### `config set`

Set a configuration value.

```
kalshi-cli config set <key> <value>
```

```bash
kalshi-cli config set output.format json
kalshi-cli config set defaults.limit 100
```

---

### version

Print version information.

```
kalshi-cli version
```

---

### completion

Generate shell autocompletion scripts.

```
kalshi-cli completion bash
kalshi-cli completion zsh
kalshi-cli completion fish
kalshi-cli completion powershell
```

## Configuration

Configuration file: `~/.kalshi/config.yaml` (created on first run).

```yaml
api:
  production: false
  timeout: 30s
api_key_id: ""
private_key_path: ""
output:
  format: table
  color: true
defaults:
  limit: 50
```

### Environment Variables

All config values can be overridden:

| Variable | Description |
|----------|-------------|
| `KALSHI_API_PRODUCTION` | `true` for production |
| `KALSHI_API_TIMEOUT` | HTTP request timeout (e.g., `60s`) |
| `KALSHI_OUTPUT_FORMAT` | Default output format |
| `KALSHI_API_KEY_ID` | API Key ID |
| `KALSHI_PRIVATE_KEY` | Private key PEM content |
| `KALSHI_PRIVATE_KEY_FILE` | Path to private key PEM file |

### Demo vs Production

| | Demo | Production |
|--|------|-----------|
| **Flag** | (default) | `--prod` |
| **API** | `demo-api.kalshi.co` | `api.elections.kalshi.com` |
| **WebSocket** | `wss://demo-api.kalshi.co/trade-api/ws/v2` | `wss://api.elections.kalshi.com/trade-api/ws/v2` |

## Bot Integration

### Automation Flags

```bash
kalshi-cli --json --yes [command] [subcommand] [flags]
```

| Flag | Purpose |
|------|---------|
| `--json` | Machine-parseable structured output |
| `--yes` | Skip all interactive confirmations |
| `--plain` | Unformatted text for piping |
| `--prod` | Target production |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Authentication error |
| 3 | Validation error |
| 4 | API error |
| 5 | Network error |

### JSON Output Schemas

**Order:**
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

**Position:**
```json
{
  "ticker": "string",
  "position": 10,
  "market_exposure": 500,
  "realized_pnl": 100,
  "total_cost": 400
}
```

**Balance:**
```json
{
  "available_balance": 10000,
  "portfolio_value": 5000,
  "total_balance": 15000
}
```

**Candlestick:**
```json
{
  "open": 50,
  "high": 60,
  "low": 40,
  "close": 55,
  "volume": 100,
  "open_interest": 200,
  "period_end": "2026-02-11T16:00:00Z"
}
```

All monetary values are in **cents**.

### Example: Market Making Bot

```bash
#!/bin/bash
MARKET="KXBTC-26FEB12-B97000"
BOOK=$(kalshi-cli markets orderbook $MARKET --json)
BEST_BID=$(echo $BOOK | jq '.yes_bids[0].price // 0')
BEST_ASK=$(echo $BOOK | jq '.yes_asks[0].price // 100')

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
kalshi-cli watch ticker KXBTC-26FEB12-B97000 --json | while read line; do
  PRICE=$(echo $line | jq '.yes_price')
  echo "Current price: $PRICE"
done
```

## Architecture

```
kalshi-cli/
├── cmd/kalshi-cli/        # Entry point
│   └── main.go
├── internal/
│   ├── api/               # HTTP client, RSA-PSS auth signing, all API methods
│   ├── cmd/               # Cobra command definitions
│   ├── config/            # Viper config + keyring credential store
│   ├── ui/                # Table formatting, ASCII candlestick charts, output routing
│   └── websocket/         # WebSocket client, channel subscriptions, auto-reconnect
├── pkg/
│   └── models/            # Shared request/response types
├── .goreleaser.yaml       # Cross-platform release builds
├── go.mod
└── go.sum
```

**Key design decisions:**
- **Cobra + Viper** for CLI framework and configuration
- **Resty** HTTP client with automatic retry and rate-limit handling
- **nhooyr.io/websocket** for WebSocket streaming with auto-reconnect
- **OS keyring** for credential storage (never plaintext)
- **RSA-PSS signatures** (`timestamp_ms + METHOD + path`) for API authentication
- **Demo-first** - production requires explicit `--prod` flag
- **lipgloss** for terminal styling (green/red price coloring, chart rendering)

## License

MIT License - see [LICENSE](LICENSE) for details.
