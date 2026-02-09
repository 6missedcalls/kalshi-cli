# kalshi-cli

A comprehensive command-line interface for the [Kalshi](https://kalshi.com) prediction market exchange. Access all API endpoints, real-time WebSocket streaming, and enjoy a first-class trading experience from your terminal.

## Installation

### Homebrew (macOS and Linux)

```bash
brew install 6missedcalls/tap/kalshi-cli
```

### Go Install

```bash
go install github.com/6missedcalls/kalshi-cli/cmd/kalshi-cli@latest
```

### Binary Download

Download the latest release for your platform from the [Releases](https://github.com/6missedcalls/kalshi-cli/releases) page.

#### macOS (Apple Silicon)
```bash
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_darwin_arm64.tar.gz
tar -xzf kalshi-cli_darwin_arm64.tar.gz
sudo mv kalshi-cli /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_darwin_amd64.tar.gz
tar -xzf kalshi-cli_darwin_amd64.tar.gz
sudo mv kalshi-cli /usr/local/bin/
```

#### Linux (x86_64)
```bash
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_linux_amd64.tar.gz
tar -xzf kalshi-cli_linux_amd64.tar.gz
sudo mv kalshi-cli /usr/local/bin/
```

#### Linux (ARM64)
```bash
curl -LO https://github.com/6missedcalls/kalshi-cli/releases/latest/download/kalshi-cli_linux_arm64.tar.gz
tar -xzf kalshi-cli_linux_arm64.tar.gz
sudo mv kalshi-cli /usr/local/bin/
```

#### Windows
Download `kalshi-cli_windows_amd64.zip` from the releases page and add the executable to your PATH.

## Quick Start

1. **Configure your API credentials**

   Generate API keys from your [Kalshi account settings](https://kalshi.com/account/api-keys) and configure the CLI:

   ```bash
   kalshi-cli auth login
   ```

2. **Explore available commands**

   ```bash
   kalshi-cli --help
   ```

3. **Start with the demo environment**

   By default, all commands use the demo API. This is perfect for testing:

   ```bash
   kalshi-cli markets list
   ```

4. **Switch to production when ready**

   ```bash
   kalshi-cli --prod markets list
   ```

## Command Overview

```
kalshi-cli [command] [flags]

Global Flags:
  --config string   Config file (default: $HOME/.kalshi/config.yaml)
  --prod            Use production API (default: demo)
  --json            Output as JSON
  --plain           Output as plain text (for pipes)
  -y, --yes         Skip confirmation prompts
  -v, --verbose     Verbose output
```

## Configuration

Configuration is stored in `~/.kalshi/config.yaml`. You can also use environment variables with the `KALSHI_` prefix.

### Example Configuration

```yaml
api:
  production: false
  timeout: 30s

output:
  format: table
  color: true

defaults:
  limit: 50
```

### Environment Variables

```bash
export KALSHI_API_PRODUCTION=true
export KALSHI_API_TIMEOUT=60s
export KALSHI_OUTPUT_FORMAT=json
```

### Credential Storage

API credentials are stored securely using your system's keyring:
- **macOS**: Keychain
- **Linux**: Secret Service (GNOME Keyring, KWallet)
- **Windows**: Windows Credential Manager

## Output Formats

The CLI supports multiple output formats:

| Format | Flag | Description |
|--------|------|-------------|
| Table | (default) | Human-readable tables with colors |
| JSON | `--json` | Machine-readable JSON output |
| Plain | `--plain` | Plain text for piping to other tools |

## Demo vs Production

Kalshi provides a demo environment for testing. The CLI uses demo by default:

| Environment | Flag | API Base URL |
|-------------|------|--------------|
| Demo | (default) | `demo-api.kalshi.co` |
| Production | `--prod` | `api.kalshi.com` |

## License

MIT License - see [LICENSE](LICENSE) for details.
