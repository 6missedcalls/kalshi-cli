package models

import "time"

// MarketPosition represents a market position from the API
type MarketPosition struct {
	Ticker                string `json:"ticker"`
	TotalTraded           int    `json:"total_traded"`
	TotalTradedDollars    string `json:"total_traded_dollars"`
	Position              int    `json:"position"`
	PositionFP            string `json:"position_fp"`
	MarketExposure        int    `json:"market_exposure"`
	MarketExposureDollars string `json:"market_exposure_dollars"`
	RealizedPnl           int    `json:"realized_pnl"`
	RealizedPnlDollars    string `json:"realized_pnl_dollars"`
	RestingOrdersCount    int    `json:"resting_orders_count"`
	FeesPaid              int    `json:"fees_paid"`
	FeesPaidDollars       string `json:"fees_paid_dollars"`
	LastUpdatedTs         string `json:"last_updated_ts"`
}

// Position is an alias for MarketPosition for backward compatibility
type Position = MarketPosition

// EventPosition represents an event-level position from the API
type EventPosition struct {
	EventTicker          string `json:"event_ticker"`
	TotalCost            int    `json:"total_cost"`
	TotalCostDollars     string `json:"total_cost_dollars"`
	TotalCostShares      int    `json:"total_cost_shares"`
	TotalCostSharesFP    string `json:"total_cost_shares_fp"`
	EventExposure        int    `json:"event_exposure"`
	EventExposureDollars string `json:"event_exposure_dollars"`
	RealizedPnl          int    `json:"realized_pnl"`
	RealizedPnlDollars   string `json:"realized_pnl_dollars"`
	FeesPaid             int    `json:"fees_paid"`
	FeesPaidDollars      string `json:"fees_paid_dollars"`
}

// PositionsResponse is the API response for positions
type PositionsResponse struct {
	Positions      []MarketPosition `json:"market_positions"`
	EventPositions []EventPosition  `json:"event_positions,omitempty"`
	Cursor         string           `json:"cursor"`
}

// BalanceResponse is the API response for balance
type BalanceResponse struct {
	Balance        int   `json:"balance"`
	PortfolioValue int   `json:"portfolio_value"`
	UpdatedTs      int64 `json:"updated_ts"`
}

// Fill represents a trade fill
type Fill struct {
	TradeID     string    `json:"trade_id"`
	OrderID     string    `json:"order_id"`
	Ticker      string    `json:"ticker"`
	Side        string    `json:"side"`
	Action      string    `json:"action"`
	Type        string    `json:"type"`
	YesPrice    int       `json:"yes_price"`
	NoPrice     int       `json:"no_price"`
	Count       int       `json:"count"`
	IsTaker     bool      `json:"is_taker"`
	CreatedTime time.Time `json:"created_time"`
}

// FillsResponse is the API response for fills
type FillsResponse struct {
	Fills  []Fill `json:"fills"`
	Cursor string `json:"cursor"`
}

// Settlement represents a market settlement
type Settlement struct {
	Ticker       string    `json:"ticker"`
	MarketResult string    `json:"market_result"`
	NoTotalCost  int       `json:"no_total_cost"`
	YesTotalCost int       `json:"yes_total_cost"`
	NoCount      int       `json:"no_count"`
	YesCount     int       `json:"yes_count"`
	Revenue      int       `json:"revenue"`
	SettledTime  time.Time `json:"settled_time"`
}

// SettlementsResponse is the API response for settlements
type SettlementsResponse struct {
	Settlements []Settlement `json:"settlements"`
	Cursor      string       `json:"cursor"`
}

// Subaccount represents a subaccount
type Subaccount struct {
	SubaccountID     int `json:"subaccount_id"`
	Balance          int `json:"balance"`
	AvailableBalance int `json:"available_balance"`
}

// SubaccountsResponse is the API response for subaccounts
type SubaccountsResponse struct {
	Subaccounts []Subaccount `json:"subaccounts"`
}

// Transfer represents a subaccount transfer
type Transfer struct {
	TransferID     string    `json:"transfer_id"`
	FromSubaccount int       `json:"from_subaccount"`
	ToSubaccount   int       `json:"to_subaccount"`
	Amount         int       `json:"amount"`
	CreatedTime    time.Time `json:"created_time"`
}

// TransfersResponse is the API response for transfers
type TransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

// TransferRequest is the request to transfer between subaccounts
type TransferRequest struct {
	FromSubaccount int `json:"from_subaccount_id"`
	ToSubaccount   int `json:"to_subaccount_id"`
	Amount         int `json:"amount"`
}

// SubaccountBalance represents a subaccount balance entry
type SubaccountBalance struct {
	SubaccountID     int `json:"subaccount_id"`
	Balance          int `json:"balance"`
	AvailableBalance int `json:"available_balance"`
}

// SubaccountBalancesResponse is the API response for subaccount balances
type SubaccountBalancesResponse struct {
	Balances []SubaccountBalance `json:"balances"`
}

// RestingOrderValueResponse is the API response for resting order value
type RestingOrderValueResponse struct {
	RestingOrderValue int `json:"resting_order_value"`
}
