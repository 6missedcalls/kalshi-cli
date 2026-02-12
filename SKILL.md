---
name: kalshi-cli
description: Comprehensive CLI for the Kalshi prediction market exchange. Provides trading, portfolio management, market data, real-time WebSocket streaming, and account management. Use when building trading workflows, querying markets/events, managing orders/positions, or automating Kalshi API interactions.
metadata:
  author: 6missedcalls
  version: "1.0"
compatibility: Requires Go 1.21+. Uses system keyring for credential storage.
---

# kalshi-cli

A Go CLI (Cobra + Viper) for the Kalshi prediction market exchange API v2. Binary: `kalshi-cli`.

## When to use

Use this skill when:
- Creating, amending, or canceling trading orders on Kalshi
- Querying market data, events, series, or orderbooks
- Managing portfolio positions, fills, settlements, or subaccounts
- Streaming real-time data via WebSocket (ticker, orderbook, trades, orders, fills, positions)
- Managing authentication credentials and API keys
- Working with RFQs (Request for Quotes) and block trading
- Managing order groups for grouped order execution
- Checking exchange status, schedule, or announcements
- Configuring CLI output format and defaults

## Prerequisites

1. Get API credentials from https://kalshi.com/account/api (prod) or https://demo.kalshi.com/account/api (demo)
2. Login: `kalshi-cli auth login --api-key-id <id> --private-key-file /path/to/key.pem`
3. By default all commands use the **demo** environment. Add `--prod` for production.

## Global flags

Available on ALL commands:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | | string | `~/.kalshi/config.yaml` | Config file path |
| `--prod` | | bool | false | Use production API |
| `--json` | | bool | false | Output as JSON |
| `--plain` | | bool | false | Plain text output (for pipes/scripts) |
| `--yes` | `-y` | bool | false | Skip confirmation prompts |
| `--verbose` | `-v` | bool | false | Verbose output |

## Command tree

```
kalshi-cli
├── auth                          # Manage authentication
│   ├── login                     # Authenticate with Kalshi
│   ├── logout                    # Clear stored credentials
│   ├── status                    # Show auth status
│   └── keys                      # Manage API keys
│       ├── list                  # List API keys
│       ├── create                # Create new API key
│       └── delete <id>           # Delete an API key
├── markets                       # Market data
│   ├── list                      # List markets
│   ├── get <ticker>              # Get market details
│   ├── orderbook <ticker>        # Visual orderbook display
│   ├── trades <ticker>           # Recent trades
│   ├── candlesticks <ticker>     # OHLCV candlestick data
│   └── series                    # Market series
│       ├── list                  # List series
│       └── get <ticker>          # Get series details
├── events                        # Event data
│   ├── list                      # List events
│   ├── get <ticker>              # Get event details
│   ├── candlesticks <ticker>     # Event OHLCV data
│   └── multivariate              # Multivariate events
│       ├── list                  # List multivariate events
│       └── get <ticker>          # Get multivariate event
├── orders                        # Order management
│   ├── list                      # List orders
│   ├── get <order-id>            # Get order details
│   ├── create                    # Create new order
│   ├── cancel <order-id>         # Cancel an order
│   ├── cancel-all                # Cancel all resting orders
│   ├── amend <order-id>          # Amend order qty/price
│   ├── batch-create              # Create orders from JSON file
│   └── queue <order-id>          # Get queue position
├── portfolio                     # Portfolio management
│   ├── balance                   # Show account balance
│   ├── positions                 # List positions
│   ├── fills                     # List trade fills
│   ├── settlements               # List settlements
│   └── subaccounts               # Subaccount management
│       ├── list                  # List subaccounts
│       ├── create                # Create subaccount
│       └── transfer              # Transfer between subaccounts
├── order-groups (alias: og)      # Grouped order management
│   ├── list                      # List order groups
│   ├── get <group-id>            # Get group details
│   ├── create                    # Create order group
│   ├── delete <group-id>         # Delete order group
│   ├── reset <group-id>          # Reset filled count
│   ├── trigger <group-id>        # Trigger order group
│   └── update-limit <group-id>   # Update contract limit
├── rfq                           # Request for Quotes
│   ├── list                      # List RFQs
│   ├── get <rfq-id>              # Get RFQ details
│   ├── create                    # Create new RFQ
│   └── delete <rfq-id>           # Delete an RFQ
├── quotes                        # Quote management (top-level)
│   ├── list                      # List quotes
│   ├── create                    # Create quote on RFQ
│   ├── accept <quote-id>         # Accept a quote
│   └── confirm <quote-id>        # Confirm a quote
├── exchange                      # Exchange information
│   ├── status                    # Exchange status
│   ├── schedule                  # Trading schedule
│   └── announcements             # Exchange announcements
├── watch                         # Real-time WebSocket streams
│   ├── ticker <ticker>           # Live price updates
│   ├── orderbook <ticker>        # Orderbook deltas
│   ├── trades                    # Public trades feed
│   ├── orders                    # Your order updates
│   ├── fills                     # Your fill notifications
│   └── positions                 # Your position changes
├── config                        # CLI configuration
│   ├── show                      # Show all settings
│   ├── get <key>                 # Get a config value
│   └── set <key> <value>         # Set a config value
└── version                       # Print version info
```

## Quick reference: All command-specific flags

### auth login
| Flag | Type | Description |
|------|------|-------------|
| `--api-key-id` | string | API Key ID (or env `KALSHI_API_KEY_ID`) |
| `--private-key` | string | Private key PEM content (or env `KALSHI_PRIVATE_KEY`) |
| `--private-key-file` | string | Path to private key PEM file |

### auth keys create
| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | Name for the new API key |

### markets list
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | | Filter: open, closed, settled |
| `--limit` | int | 50 | Max results |
| `--series` | string | | Filter by series ticker |

### markets trades
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 100 | Max trades to return |

### markets candlesticks
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--series` | string | **required** | Series ticker |
| `--period` | string | 1h | Period: 1m, 1h, 1d |

### markets series list
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--category` | string | | Filter by category |
| `--limit` | int | 50 | Max results |

### events list
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | | Filter: active, closed, settled |
| `--limit` | int | 50 | Max results |
| `--cursor` | string | | Pagination cursor |

### events candlesticks
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--series` | string | | Series ticker (auto-resolved if omitted) |
| `--period` | string | 1h | Period: 1m, 1h, 1d |
| `--start` | string | | Start time (RFC3339) |
| `--end` | string | | End time (RFC3339) |

### events multivariate list
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | | Filter by status |
| `--limit` | int | 50 | Max results |
| `--cursor` | string | | Pagination cursor |

### orders list
| Flag | Type | Description |
|------|------|-------------|
| `--status` | string | Filter: resting, canceled, executed, pending |
| `--market` | string | Filter by market ticker |

### orders create
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--market` | string | **required** | Market ticker |
| `--side` | string | **required** | yes or no |
| `--qty` | int | **required** | Quantity |
| `--price` | int | **required** | Price in cents (1-99) |
| `--action` | string | buy | buy or sell |
| `--type` | string | limit | limit or market |

### orders cancel-all
| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | Filter by market ticker |

### orders amend
| Flag | Type | Description |
|------|------|-------------|
| `--qty` | int | New quantity (at least one of qty/price required) |
| `--price` | int | New price in cents (at least one of qty/price required) |

### orders batch-create
| Flag | Type | Description |
|------|------|-------------|
| `--file` | string | **required** - Path to JSON file with order array |

Batch JSON format: `[{"ticker":"...","side":"yes","action":"buy","type":"limit","count":10,"yes_price":50}]`

### portfolio positions
| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | Filter by market ticker |

### portfolio fills
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 100 | Max fills to return |

### portfolio settlements
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Max settlements to return |

### portfolio subaccounts transfer
| Flag | Type | Description |
|------|------|-------------|
| `--from` | int | **required** - Source subaccount ID |
| `--to` | int | **required** - Destination subaccount ID |
| `--amount` | int | **required** - Amount in cents |

### order-groups list
| Flag | Type | Description |
|------|------|-------------|
| `--status` | string | Filter by status |

### order-groups create
| Flag | Type | Description |
|------|------|-------------|
| `--limit` | int | **required** - Max contracts to fill across group |

### order-groups update-limit
| Flag | Type | Description |
|------|------|-------------|
| `--limit` | int | **required** - New max contracts to fill |

### rfq list
| Flag | Type | Description |
|------|------|-------------|
| `--status` | string | Filter: open, closed |

### rfq create
| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | **required** - Market ticker |
| `--qty` | int | **required** - Quantity (must be > 0) |

### quotes list
| Flag | Type | Description |
|------|------|-------------|
| `--rfq-id` | string | Filter by RFQ ID |

### quotes create
| Flag | Type | Description |
|------|------|-------------|
| `--rfq` | string | **required** - RFQ ID |
| `--price` | int | **required** - Price in cents (1-99) |

### watch trades
| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | Filter trades by market ticker |

### config set / config get
Valid keys: `output.format` (table/json/plain), `output.color` (true/false), `defaults.limit` (positive int)

## Common patterns

**All prices are in cents** (1-99 for contract prices, larger for balances). Display helpers convert to dollars.

**Output formats**: Every command supports `--json`, `--plain`, and table (default). Use `--plain` for scripting.

**Pagination**: List commands accept `--limit` and some accept `--cursor` for cursor-based pagination.

**Confirmation prompts**: Destructive actions (cancel, delete, transfer) prompt for confirmation. Use `--yes` to bypass.

**Credential resolution** (auth login): flags > env vars (`KALSHI_API_KEY_ID`, `KALSHI_PRIVATE_KEY`) > interactive prompt.

**Environment**: Demo by default. Add `--prod` for production. Config: `~/.kalshi/config.yaml`.

## Detailed references

- [Auth commands](references/auth.md) - Login flows, credential storage, API key management
- [Markets commands](references/markets.md) - Market data, orderbook, trades, candlesticks, series
- [Events commands](references/events.md) - Events, multivariate events, event candlesticks
- [Orders commands](references/orders.md) - Order lifecycle, batch creation, amendments
- [Portfolio commands](references/portfolio.md) - Balance, positions, fills, settlements, subaccounts
- [Order groups commands](references/order-groups.md) - Grouped order management
- [RFQ and quotes commands](references/rfq-quotes.md) - Block trading workflow
- [Exchange commands](references/exchange.md) - Exchange status, schedule, announcements
- [Watch commands](references/watch.md) - WebSocket real-time streaming
- [Config commands](references/config.md) - CLI configuration management
- [Data models](references/models.md) - API response structures and types
