# Exchange Commands Reference

## `kalshi-cli exchange`

Get exchange status, schedule, and announcements. All subcommands require authentication.

## `kalshi-cli exchange status`

Get current exchange status including trading activity and environment.

**Output fields**: Exchange Active (Yes/No), Trading Active (Yes/No), Environment (Production/Demo).

```bash
kalshi-cli exchange status
kalshi-cli exchange status --prod
kalshi-cli exchange status --json
kalshi-cli exchange status --plain
# Plain output: exchange_active=yes, trading_active=yes, environment=demo
```

## `kalshi-cli exchange schedule`

Get the exchange trading schedule.

**Output**: Standard hours (per-day open/close times grouped by week), Maintenance windows (start/end datetimes).

```bash
kalshi-cli exchange schedule
kalshi-cli exchange schedule --json
kalshi-cli exchange schedule --plain
# Plain output: week_0_start=..., week_0_monday_0_open=..., maintenance_0_start=...
```

## `kalshi-cli exchange announcements`

Get the latest exchange announcements.

**Output columns**: Title (truncated 50 chars), Type, Status (color-coded: green=active, yellow=pending, gray=expired), Delivery Time.

```bash
kalshi-cli exchange announcements
kalshi-cli exchange announcements --json
kalshi-cli exchange announcements --plain
# Plain output: announcement_0_id=..., announcement_0_title=..., announcement_0_status=...
```
