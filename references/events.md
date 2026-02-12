# Events Commands Reference

## `kalshi-cli events`

Commands for listing, viewing, and managing Kalshi events. An event groups related markets (e.g., "S&P 500 close on Feb 7" has multiple price-bracket markets under it).

## `kalshi-cli events list`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | "" | Filter: active, closed, settled |
| `--limit` | int | 50 | Max results |
| `--cursor` | string | "" | Pagination cursor |

```bash
kalshi-cli events list
kalshi-cli events list --status active --limit 20
kalshi-cli events list --json
```

## `kalshi-cli events get <event-ticker>`

Get detailed information about a specific event by ticker.

```bash
kalshi-cli events get INXD-25FEB07
```

## `kalshi-cli events candlesticks <event-ticker>`

Get candlestick (OHLCV) data for an event. The `--series` flag is optional; if omitted, the series ticker is auto-resolved from the event.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--series` | string | "" | Series ticker (auto-resolved from event if omitted) |
| `--period` | string | 1h | Candlestick period: 1m, 1h, 1d |
| `--start` | string | "" | Start time (RFC3339 format) |
| `--end` | string | "" | End time (RFC3339 format) |

```bash
kalshi-cli events candlesticks INXD-25FEB07 --start 2025-02-01T00:00:00Z --end 2025-02-07T00:00:00Z
kalshi-cli events candlesticks INXD-25FEB07 --period 1d --start 2025-01-01T00:00:00Z --end 2025-02-01T00:00:00Z
kalshi-cli events candlesticks INXD-25FEB07 --series INXD --period 1h --start 2025-02-06T00:00:00Z --end 2025-02-07T00:00:00Z
```

If the event has no series ticker and `--series` is omitted, an error is returned asking the user to provide it explicitly.

## `kalshi-cli events multivariate list`

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | "" | Filter by status |
| `--limit` | int | 50 | Max results |
| `--cursor` | string | "" | Pagination cursor |

```bash
kalshi-cli events multivariate list
kalshi-cli events multivariate list --status active
```

## `kalshi-cli events multivariate get <ticker>`

Get detailed information about a specific multivariate event.

```bash
kalshi-cli events multivariate get INXD-MV
```
