# Portfolio Commands Reference

## `kalshi-cli portfolio`

View and manage your Kalshi portfolio including balance, positions, fills, settlements, and subaccounts.

## `kalshi-cli portfolio balance`

Display current account balance including available balance, portfolio value, and total balance. All values in cents.

```bash
kalshi-cli portfolio balance
kalshi-cli portfolio balance --json
```

## `kalshi-cli portfolio positions`

List current market positions with average cost, P&L, and exposure.

| Flag | Type | Description |
|------|------|-------------|
| `--market` | string | Filter by market ticker |

```bash
kalshi-cli portfolio positions
kalshi-cli portfolio positions --market INXD-25FEB07-B5523.99
```

## `kalshi-cli portfolio fills`

List trade fills showing executed orders and details.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 100 | Max fills to return |

```bash
kalshi-cli portfolio fills
kalshi-cli portfolio fills --limit 20
```

## `kalshi-cli portfolio settlements`

List market settlements showing resolved positions and outcomes.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Max settlements to return |

```bash
kalshi-cli portfolio settlements
kalshi-cli portfolio settlements --limit 10
```

## `kalshi-cli portfolio subaccounts list`

List all subaccounts associated with your account.

```bash
kalshi-cli portfolio subaccounts list
```

## `kalshi-cli portfolio subaccounts create`

Create a new subaccount.

```bash
kalshi-cli portfolio subaccounts create
```

## `kalshi-cli portfolio subaccounts transfer`

Transfer funds between subaccounts. Prompts for confirmation.

| Flag | Type | Description |
|------|------|-------------|
| `--from` | int | **required** - Source subaccount ID |
| `--to` | int | **required** - Destination subaccount ID |
| `--amount` | int | **required** - Amount in cents (must be positive) |

```bash
kalshi-cli portfolio subaccounts transfer --from 1 --to 2 --amount 10000
kalshi-cli portfolio subaccounts transfer --from 1 --to 2 --amount 10000 --yes
```
