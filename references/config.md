# Config Commands Reference

## `kalshi-cli config`

Manage CLI configuration. Config stored at `~/.kalshi/config.yaml`.

## Available configuration keys

| Key | Type | Valid Values | Default | Description |
|-----|------|-------------|---------|-------------|
| `output.format` | string | table, json, plain | table | Default output format |
| `output.color` | bool | true, false | true | Enable colored output |
| `defaults.limit` | int | Any positive integer | 50 | Default limit for list commands |

## `kalshi-cli config show`

Display all current configuration settings with keys, values, and descriptions. Also shows config file path.

```bash
kalshi-cli config show
kalshi-cli config show --json
```

## `kalshi-cli config get <key>`

Get the value of a specific configuration key.

```bash
kalshi-cli config get output.format
kalshi-cli config get defaults.limit --json
```

## `kalshi-cli config set <key> <value>`

Set a configuration value. Validates before saving.

```bash
kalshi-cli config set output.format json
kalshi-cli config set output.color false
kalshi-cli config set defaults.limit 100
```

## Environment variables

All config keys can be set via environment variables with the `KALSHI_` prefix:
- `KALSHI_API_PRODUCTION=true` maps to `api.production`
- `KALSHI_OUTPUT_FORMAT=json` maps to `output.format`

## API configuration

| Config Key | Default | Description |
|------------|---------|-------------|
| `api.production` | false | Use production API |
| `api.timeout` | 30s | API request timeout |

## API URLs

| Environment | Base URL | WebSocket URL |
|-------------|----------|---------------|
| Demo | `https://demo-api.kalshi.co` | `wss://demo-api.kalshi.co/trade-api/ws/v2` |
| Production | `https://api.elections.kalshi.com` | `wss://api.elections.kalshi.com/trade-api/ws/v2` |
