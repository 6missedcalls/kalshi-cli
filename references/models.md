# Data Models Reference

API response structures and types used across commands. All prices are in cents (int).

## Market

| Field | Type | Description |
|-------|------|-------------|
| Ticker | string | Market identifier |
| EventTicker | string | Parent event |
| Title | string | Market title |
| Subtitle | string | Market subtitle |
| Status | string | open, closed, settled |
| Category | string | Market category |
| YesBid, YesAsk | int | Current yes bid/ask in cents |
| NoBid, NoAsk | int | Current no bid/ask in cents |
| LastPrice | int | Last trade price in cents |
| Volume, Volume24H | int | Total and 24h volume |
| OpenInterest | int | Open interest |
| Result | string | Settlement result (if settled) |
| OpenTime, CloseTime, ExpirationTime | time.Time | Market timing |

## Order

**Constants**: Side: yes/no. Type: limit/market. Status: resting/canceled/executed/pending. Action: buy/sell.

| Field | Type | Description |
|-------|------|-------------|
| OrderID | string | Unique order ID |
| Ticker | string | Market ticker |
| Status | OrderStatus | Current status |
| Side | OrderSide | yes or no |
| Action | OrderAction | buy or sell |
| Type | OrderType | limit or market |
| YesPrice, NoPrice | int | Price in cents |
| InitialCount | int | Original quantity |
| RemainingCount | int | Unfilled quantity |
| FillCount | int | Filled quantity |
| CreatedTime | time.Time | Creation timestamp |

## Event

| Field | Type | Description |
|-------|------|-------------|
| EventTicker | string | Event identifier |
| SeriesTicker | string | Parent series |
| Title | string | Event title |
| Category | string | Category |
| MutuallyExclusive | bool | Markets are mutually exclusive |
| Markets | []string | List of market tickers |

## Position (MarketPosition)

| Field | Type | Description |
|-------|------|-------------|
| Ticker | string | Market ticker |
| Position | int | Net position |
| TotalTraded | int | Total contracts traded |
| MarketExposure | int | Current exposure in cents |
| RealizedPnl | int | Realized P&L in cents |
| FeesPaid | int | Total fees in cents |

## Balance (BalanceResponse)

| Field | Type | Description |
|-------|------|-------------|
| Balance | int | Available balance in cents |
| PortfolioValue | int | Portfolio value in cents |

## Fill

| Field | Type | Description |
|-------|------|-------------|
| TradeID | string | Trade identifier |
| OrderID | string | Parent order |
| Ticker | string | Market ticker |
| Side | string | yes or no |
| Action | string | buy or sell |
| YesPrice, NoPrice | int | Execution price in cents |
| Count | int | Fill quantity |
| IsTaker | bool | Taker or maker |
| CreatedTime | time.Time | Fill timestamp |

## OrderGroup

| Field | Type | Description |
|-------|------|-------------|
| GroupID | string | Group identifier |
| Status | string | Group status |
| Limit | int | Max contracts to fill |
| FilledCount | int | Current filled count |
| OrderCount | int | Number of orders in group |
| OrderIDs | []string | Order identifiers |

## RFQ

| Field | Type | Description |
|-------|------|-------------|
| ID | string | RFQ identifier |
| MarketTicker | string | Market ticker |
| Status | string | RFQ status |
| Contracts | int | Requested quantity |

## Quote

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Quote identifier |
| RFQID | string | Parent RFQ |
| MarketTicker | string | Market ticker |
| Status | string | Quote status |
| Contracts | int | Quantity |
| YesBid, NoBid | int | Bid prices in cents |

## Candlestick

| Field | Type | Description |
|-------|------|-------------|
| Open, High, Low, Close | int | OHLC prices in cents |
| Volume | int | Period volume |
| Timestamp | time.Time | Candle start time |

## API response wrappers

- Single items: `{"object": {...}}`
- Lists: `{"objects": [...], "cursor": "..."}`
- Use `cursor` value in `--cursor` flag for pagination

## API error format

```json
{"code": "not_found", "message": "Market not found", "status_code": 404}
```

Automatic retry on 429 (rate limit) and 5xx errors with exponential backoff (100ms base, 10s max, 5 retries).
