# Order Groups Commands Reference

## `kalshi-cli order-groups` (alias: `og`)

Order groups allow you to group multiple orders and manage them as a single unit with a shared contract fill limit.

## `kalshi-cli order-groups list`

| Flag | Type | Description |
|------|------|-------------|
| `--status` | string | Filter by status |

**Output columns**: Group ID, Status, Limit, Filled, Order Count.

```bash
kalshi-cli order-groups list
kalshi-cli og list --status active --json
```

## `kalshi-cli order-groups get <group-id>`

**Output fields**: Group ID, Status, Limit, Filled Count, Order Count, Created, Last Updated, Order IDs.

```bash
kalshi-cli order-groups get abc-123
kalshi-cli og get abc-123 --json
```

## `kalshi-cli order-groups create`

| Flag | Type | Description |
|------|------|-------------|
| `--limit` | int | **required** - Max contracts to fill across group (must be > 0) |

```bash
kalshi-cli order-groups create --limit 100
kalshi-cli og create --limit 50 --json
```

## `kalshi-cli order-groups delete <group-id>`

Delete an order group. All orders in the group will be canceled. Prompts for confirmation.

```bash
kalshi-cli order-groups delete abc-123
kalshi-cli og delete abc-123 --yes
```

## `kalshi-cli order-groups reset <group-id>`

Reset an order group's filled count to zero, allowing more orders to fill.

```bash
kalshi-cli order-groups reset abc-123
```

## `kalshi-cli order-groups trigger <group-id>`

Trigger an order group to execute its orders.

```bash
kalshi-cli order-groups trigger abc-123
```

## `kalshi-cli order-groups update-limit <group-id>`

Update the max contract limit. If new limit < current filled count, the group will be triggered.

| Flag | Type | Description |
|------|------|-------------|
| `--limit` | int | **required** - New max contracts (>= 0) |

```bash
kalshi-cli order-groups update-limit abc-123 --limit 200
kalshi-cli og update-limit abc-123 --limit 0
```
