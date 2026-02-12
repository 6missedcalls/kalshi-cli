# RFQ and Quotes Commands Reference

Two top-level command groups for block trading: `rfq` and `quotes`.

## `kalshi-cli rfq list`

| Flag | Type | Description |
|------|------|-------------|
| `--status` | string | Filter: open, closed |

```bash
kalshi-cli rfq list
kalshi-cli rfq list --status open
```

## `kalshi-cli rfq get <rfq-id>`

```bash
kalshi-cli rfq get rfq_abc123
```

## `kalshi-cli rfq create`

| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | **required** - Market ticker |
| `--qty` | int | **required** - Quantity (must be > 0) |

```bash
kalshi-cli rfq create --market INXD-25FEB07 --qty 1000
```

## `kalshi-cli rfq delete <rfq-id>`

Prompts for confirmation before deleting.

```bash
kalshi-cli rfq delete rfq_abc123
```

## `kalshi-cli quotes list`

| Flag | Type | Description |
|------|------|-------------|
| `--rfq-id` | string | Filter by RFQ ID |

```bash
kalshi-cli quotes list
kalshi-cli quotes list --rfq-id rfq_abc123
```

## `kalshi-cli quotes create`

| Flag | Type | Description |
|------|------|-------------|
| `--rfq` | string | **required** - RFQ ID |
| `--price` | int | **required** - Price in cents (1-99) |

```bash
kalshi-cli quotes create --rfq rfq_abc123 --price 65
```

## `kalshi-cli quotes accept <quote-id>`

Accept a quote offered on your RFQ. Prompts for confirmation.

```bash
kalshi-cli quotes accept quote_xyz789
```

## `kalshi-cli quotes confirm <quote-id>`

Confirm a quote after acceptance. Prompts for confirmation.

```bash
kalshi-cli quotes confirm quote_xyz789
```

## Block trading workflow

1. Create an RFQ: `kalshi-cli rfq create --market TICKER --qty 1000`
2. Wait for quotes: `kalshi-cli quotes list --rfq-id RFQ_ID`
3. Accept best quote: `kalshi-cli quotes accept QUOTE_ID`
4. Confirm the trade: `kalshi-cli quotes confirm QUOTE_ID`
