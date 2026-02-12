# Auth Commands Reference

## `kalshi-cli auth`

Manage authentication credentials and API keys for the Kalshi API.

## `kalshi-cli auth login`

Authenticate with Kalshi using API credentials provisioned from the Kalshi dashboard.

**Credential resolution order**: flags > environment variables > interactive prompt.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--api-key-id` | string | "" | API Key ID (or env `KALSHI_API_KEY_ID`) |
| `--private-key` | string | "" | Private key PEM content (or env `KALSHI_PRIVATE_KEY`) |
| `--private-key-file` | string | "" | Path to private key PEM file |

**Environment variables**: `KALSHI_API_KEY_ID`, `KALSHI_PRIVATE_KEY`

```bash
# Interactive login
kalshi-cli auth login

# Non-interactive with file
kalshi-cli auth login --api-key-id <id> --private-key-file /path/to/key.pem

# Non-interactive with PEM content
kalshi-cli auth login --api-key-id <id> --private-key "$(cat key.pem)"

# Via environment variables
export KALSHI_API_KEY_ID=your-key-id
export KALSHI_PRIVATE_KEY="$(cat key.pem)"
kalshi-cli auth login
```

Credentials are stored in the system keyring after successful login.

## `kalshi-cli auth logout`

Remove stored API credentials from the system keyring. Prompts for confirmation (bypass with `--yes`).

```bash
kalshi-cli auth logout
kalshi-cli auth logout --yes
```

## `kalshi-cli auth status`

Display current authentication status and environment.

**Output fields** (JSON): `logged_in`, `api_key_id`, `environment`, `authenticated`, `exchange_active`, `trading_active`.

```bash
kalshi-cli auth status
kalshi-cli auth status --json
```

## `kalshi-cli auth keys list`

List all API keys associated with your account.

**Output columns**: ID, Name, Created, Expires, Scopes.

```bash
kalshi-cli auth keys list
kalshi-cli auth keys list --json
```

## `kalshi-cli auth keys create`

| Flag | Type | Description |
|------|------|-------------|
| `--name` | string | Name for the new API key |

```bash
kalshi-cli auth keys create --name "trading-bot"
```

## `kalshi-cli auth keys delete <id>`

Delete an API key by its ID. Prompts for confirmation (bypass with `--yes`).

```bash
kalshi-cli auth keys delete abc123
kalshi-cli auth keys delete abc123 --yes
```

## Internal functions

- `getAuthenticatedClient()` - Retrieves credentials from keyring and creates authenticated API client
- `createAuthenticatedClient(creds)` - Creates API client from credentials using `api.NewSignerFromPEM` and `api.NewClient`
- `resolveLoginCredentials(keyring)` - Resolves credentials: flags > env vars > interactive input
- `readPrivateKeyInput(reader)` - Reads multi-line PEM from stdin, or treats first line as file path
