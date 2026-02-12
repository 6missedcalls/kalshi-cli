# Orders Commands Reference

## `kalshi-cli orders`

Manage trading orders on the Kalshi exchange.

## `kalshi-cli orders list`

| Flag | Type | Description |
|------|------|-------------|
| `--status` | string | Filter: resting, canceled, executed, pending |
| `--market` | string | Filter by market ticker |

```bash
kalshi-cli orders list
kalshi-cli orders list --status resting
kalshi-cli orders list --market INXD-25FEB07-B5523.99 --json
```

## `kalshi-cli orders get <order-id>`

Get detailed information about a specific order.

```bash
kalshi-cli orders get abc123-def456-ghi789
```

## `kalshi-cli orders create`

Create a new limit order. Shows order preview before submission. Requires confirmation unless `--yes` is set.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--market` | string | **required** | Market ticker |
| `--side` | string | **required** | yes or no |
| `--qty` | int | **required** | Quantity |
| `--price` | int | **required** | Price in cents (1-99) |
| `--action` | string | buy | buy or sell |
| `--type` | string | limit | limit or market |

**Validation**:
- Price must be 1-99 cents
- Side must be "yes" or "no"
- Action must be "buy" or "sell"
- Type must be "limit" or "market"
- Quantity must be positive
- Shows PRODUCTION warning when using `--prod`

```bash
kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side yes --qty 10 --price 50
kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side no --qty 5 --price 30 --action sell
kalshi-cli orders create --market INXD-25FEB07-B5523.99 --side yes --qty 10 --price 50 --yes
```

## `kalshi-cli orders cancel <order-id>`

Cancel a resting order by ID. Prompts for confirmation.

```bash
kalshi-cli orders cancel abc123
kalshi-cli orders cancel abc123 --yes
```

## `kalshi-cli orders cancel-all`

Cancel all resting orders, optionally filtered by market.

| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | Filter by market ticker |

```bash
kalshi-cli orders cancel-all
kalshi-cli orders cancel-all --market INXD-25FEB07-B5523.99
kalshi-cli orders cancel-all --yes
```

## `kalshi-cli orders amend <order-id>`

Amend an existing order's quantity and/or price. At least one of `--qty` or `--price` must be specified.

| Flag | Type | Description |
|------|------|-------------|
| `--qty` | int | New quantity |
| `--price` | int | New price in cents (1-99) |

```bash
kalshi-cli orders amend abc123 --price 55
kalshi-cli orders amend abc123 --qty 20 --price 60
```

## `kalshi-cli orders batch-create`

Create multiple orders from a JSON file. Shows batch preview and requires confirmation.

| Flag | Type | Description |
|------|------|-------------|
| `--file` | string | **required** - Path to JSON file |

**JSON format**:
```json
[
  {
    "ticker": "INXD-25FEB07-B5523.99",
    "side": "yes",
    "action": "buy",
    "type": "limit",
    "count": 10,
    "yes_price": 50
  }
]
```

**Validation per order**: ticker required, side must be yes/no, count must be positive, prices 1-99.

```bash
kalshi-cli orders batch-create --file orders.json
kalshi-cli orders batch-create --file orders.json --yes
```

## `kalshi-cli orders queue <order-id>`

Get the queue position for a resting order.

```bash
kalshi-cli orders queue abc123
```
