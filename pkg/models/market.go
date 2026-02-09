package models

import "time"

// Market represents a Kalshi prediction market
type Market struct {
	Ticker              string    `json:"ticker"`
	EventTicker         string    `json:"event_ticker"`
	MarketType          string    `json:"market_type"`
	Title               string    `json:"title"`
	Subtitle            string    `json:"subtitle"`
	Status              string    `json:"status"`
	YesBid              int       `json:"yes_bid"`
	YesAsk              int       `json:"yes_ask"`
	NoBid               int       `json:"no_bid"`
	NoAsk               int       `json:"no_ask"`
	LastPrice           int       `json:"last_price"`
	PreviousYesBid      int       `json:"previous_yes_bid"`
	PreviousYesAsk      int       `json:"previous_yes_ask"`
	PreviousPrice       int       `json:"previous_price"`
	Volume              int       `json:"volume"`
	Volume24H           int       `json:"volume_24h"`
	OpenInterest        int       `json:"open_interest"`
	DollarVolume        int       `json:"dollar_volume"`
	DollarOpenInterest  int       `json:"dollar_open_interest"`
	Result              string    `json:"result"`
	ExpirationTime      time.Time `json:"expiration_time"`
	LatestExpirationTime time.Time `json:"latest_expiration_time"`
	CloseTime           time.Time `json:"close_time"`
	OpenTime            time.Time `json:"open_time"`
	CreatedTime         time.Time `json:"created_time"`
	CanCloseEarly       bool      `json:"can_close_early"`
	RiskLimitCents      int       `json:"risk_limit_cents"`
	NotionalValue       int       `json:"notional_value"`
	TickSize            int       `json:"tick_size"`
	YesBidFee           int       `json:"yes_bid_fee"`
	NoBidFee            int       `json:"no_bid_fee"`
	YesAskFee           int       `json:"yes_ask_fee"`
	NoAskFee            int       `json:"no_ask_fee"`
	Category            string    `json:"category"`
	Rules               string    `json:"rules"`
	RulesSecondary      string    `json:"rules_secondary"`
	SettlementTimerSeconds int    `json:"settlement_timer_seconds"`
}

// MarketResponse is the API response for a single market
type MarketResponse struct {
	Market Market `json:"market"`
}

// MarketsResponse is the API response for multiple markets
type MarketsResponse struct {
	Markets []Market `json:"markets"`
	Cursor  string   `json:"cursor"`
}

// Orderbook represents a market orderbook
type Orderbook struct {
	Ticker   string           `json:"ticker"`
	YesBids  []OrderbookLevel `json:"yes_bids"`
	YesAsks  []OrderbookLevel `json:"yes_asks"`
	NoBids   []OrderbookLevel `json:"no_bids"`
	NoAsks   []OrderbookLevel `json:"no_asks"`
}

// OrderbookLevel represents a single price level in the orderbook
type OrderbookLevel struct {
	Price    int `json:"price"`
	Quantity int `json:"quantity"`
}

// OrderbookResponse is the API response for orderbook
type OrderbookResponse struct {
	Orderbook Orderbook `json:"orderbook"`
}

// Trade represents a public trade
type Trade struct {
	TradeID    string    `json:"trade_id"`
	Ticker     string    `json:"ticker"`
	Price      int       `json:"price"`
	Count      int       `json:"count"`
	TakerSide  string    `json:"taker_side"`
	CreatedTime time.Time `json:"created_time"`
}

// TradesResponse is the API response for trades
type TradesResponse struct {
	Trades []Trade `json:"trades"`
	Cursor string  `json:"cursor"`
}

// Candlestick represents OHLCV data
type Candlestick struct {
	Ticker     string    `json:"ticker"`
	Open       int       `json:"open"`
	High       int       `json:"high"`
	Low        int       `json:"low"`
	Close      int       `json:"close"`
	Volume     int       `json:"volume"`
	OpenInterest int     `json:"open_interest"`
	PeriodEnd  time.Time `json:"period_end"`
}

// CandlesticksResponse is the API response for candlesticks
type CandlesticksResponse struct {
	Candlesticks []Candlestick `json:"candlesticks"`
}

// Series represents a market series
type Series struct {
	Ticker      string `json:"ticker"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Frequency   string `json:"frequency"`
	Tags        []string `json:"tags"`
}

// SeriesResponse is the API response for series
type SeriesResponse struct {
	Series []Series `json:"series"`
	Cursor string   `json:"cursor"`
}
