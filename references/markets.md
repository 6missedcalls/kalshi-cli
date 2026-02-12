# Markets Commands Reference

## `kalshi-cli markets`

Commands for listing, viewing, and analyzing prediction markets.

## `kalshi-cli markets list`

List markets with optional filtering.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | "" | Filter: open, closed, settled |
| `--limit` | int | 50 | Max results |
| `--series` | string | "" | Filter by series ticker |

**Output columns**: Ticker, Title, Status, Yes Bid, Yes Ask, Volume.

```bash
kalshi-cli markets list
kalshi-cli markets list --status open --limit 20
kalshi-cli markets list --series INXD --json
```

## `kalshi-cli markets get <market-ticker>`

Get detailed information about a specific market.

**Output fields**: Ticker, Title, Subtitle, Status, Category, Yes Bid, Yes Ask, No Bid, No Ask, Last Price, Volume, Volume 24h, Open Interest, Open Time, Close Time, Expiration, Result.

```bash
kalshi-cli markets get INXD-25FEB07-B5523.99
kalshi-cli markets get INXD-25FEB07-B5523.99 --json
```

## `kalshi-cli markets orderbook <market-ticker>`

Visual orderbook display with YES bids and asks at each price level.

```bash
kalshi-cli markets orderbook INXD-25FEB07-B5523.99
kalshi-cli markets orderbook INXD-25FEB07-B5523.99 --json
```

## `kalshi-cli markets trades <market-ticker>`

Get recent trades for a specific market.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 100 | Max trades to return |

**Output columns**: Time, Price, Quantity, Side.

```bash
kalshi-cli markets trades INXD-25FEB07-B5523.99
kalshi-cli markets trades INXD-25FEB07-B5523.99 --limit 20
```

## `kalshi-cli markets candlesticks <market-ticker>`

Get candlestick (OHLCV) data. Renders ASCII candlestick chart followed by data table.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--series` | string | **required** | Series ticker |
| `--period` | string | 1h | Period: 1m, 1h, 1d |

**Output**: ASCII chart + table with columns: Time, Open, High, Low, Close, Volume.

```bash
kalshi-cli markets candlesticks INXD-25FEB07-B5523.99 --series INXD
kalshi-cli markets candlesticks INXD-25FEB07-B5523.99 --series INXD --period 1d
```

## `kalshi-cli markets series list`

List market series with optional category filtering.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--category` | string | "" | Filter by category |
| `--limit` | int | 50 | Max results |

**Output columns**: Ticker, Title, Category, Frequency.

```bash
kalshi-cli markets series list
kalshi-cli markets series list --category economics
```

## `kalshi-cli markets series get <series-ticker>`

Get detailed information about a specific series.

**Output fields**: Ticker, Title, Category, Frequency, Tags.

```bash
kalshi-cli markets series get INXD
```
