# Watch Commands Reference

## `kalshi-cli watch`

Stream real-time data from Kalshi via WebSocket. All watch commands require authentication. Press Ctrl+C to stop.

## `kalshi-cli watch ticker <market-ticker>`

Live price updates for a market. Output includes bid/ask prices, volume, and open interest.

```bash
kalshi-cli watch ticker INXD-25FEB07-B5523.99
kalshi-cli watch ticker INXD-25FEB07-B5523.99 --json
kalshi-cli watch ticker INXD-25FEB07-B5523.99 --plain
```

## `kalshi-cli watch orderbook <market-ticker>`

Live orderbook delta updates. Shows best bid/ask, depth, and changes as they occur.

```bash
kalshi-cli watch orderbook INXD-25FEB07-B5523.99
kalshi-cli watch orderbook INXD-25FEB07-B5523.99 --json
```

## `kalshi-cli watch trades`

Public trades feed across all markets with optional filtering.

| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | Filter trades by market ticker |

```bash
kalshi-cli watch trades
kalshi-cli watch trades --market INXD-25FEB07-B5523.99
kalshi-cli watch trades --json
```

## `kalshi-cli watch orders`

Your order status changes (fills, cancellations, status transitions).

```bash
kalshi-cli watch orders
kalshi-cli watch orders --json
```

## `kalshi-cli watch fills`

Your fill notifications with price, count, and taker/maker status.

```bash
kalshi-cli watch fills
kalshi-cli watch fills --json
```

## `kalshi-cli watch positions`

Your position changes with realized PnL, exposure, and total cost.

```bash
kalshi-cli watch positions
kalshi-cli watch positions --json
```

## Internal WebSocket channels (not exposed as CLI subcommands)

These handlers exist in code but are not registered as CLI subcommands:

- `market_ticker_v2` - Incremental delta messages with delta_type, yes/no prices
- `market_lifecycle` - Market status transitions (old_status -> new_status)
- `order_group_updates` - Order group status changes with total/filled counts
- `communications` - RFQ/quote messages with type, ticker, quantity, price, side

## WebSocket URLs

- Demo: `wss://demo-api.kalshi.co/trade-api/ws/v2`
- Production: `wss://api.elections.kalshi.com/trade-api/ws/v2`
