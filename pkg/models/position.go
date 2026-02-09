package models

import "time"

// Position represents a market position
type Position struct {
	Ticker                string    `json:"ticker"`
	EventTicker           string    `json:"event_ticker"`
	EventExposure         int       `json:"event_exposure"`
	Exposure              int       `json:"exposure"`
	RestingBuyExposure    int       `json:"resting_buy_exposure"`
	RestingSellExposure   int       `json:"resting_sell_exposure"`
	MarketExposure        int       `json:"market_exposure"`
	Position              int       `json:"position"`
	TotalCost             int       `json:"total_cost"`
	RealizedPnl           int       `json:"realized_pnl"`
	Fees                  int       `json:"fees"`
	SettlementValue       int       `json:"settlement_value"`
	Direction             string    `json:"direction"`
	TotalBought           int       `json:"total_bought"`
	TotalSold             int       `json:"total_sold"`
	LastUpdateTime        time.Time `json:"last_update_time"`
}

// PositionsResponse is the API response for positions
type PositionsResponse struct {
	Positions           []Position `json:"market_positions"`
	Cursor              string     `json:"cursor"`
	EventPositions      []Position `json:"event_positions,omitempty"`
}

// Balance represents account balance
type Balance struct {
	Balance             int `json:"balance"`
	AvailableBalance    int `json:"available_balance"`
	PortfolioValue      int `json:"portfolio_value"`
	BonusCashBalance    int `json:"bonus_cash_balance"`
	TotalRestingOrdersValue int `json:"total_resting_orders_value"`
	PayoutBalance       int `json:"payout_balance"`
}

// BalanceResponse is the API response for balance
type BalanceResponse struct {
	Balance int `json:"balance"`
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
	Ticker         string    `json:"ticker"`
	MarketResult   string    `json:"market_result"`
	NoTotalCost    int       `json:"no_total_cost"`
	YesTotalCost   int       `json:"yes_total_cost"`
	NoCount        int       `json:"no_count"`
	YesCount       int       `json:"yes_count"`
	Revenue        int       `json:"revenue"`
	SettledTime    time.Time `json:"settled_time"`
}

// SettlementsResponse is the API response for settlements
type SettlementsResponse struct {
	Settlements []Settlement `json:"settlements"`
	Cursor      string       `json:"cursor"`
}

// Subaccount represents a subaccount
type Subaccount struct {
	SubaccountID     int    `json:"subaccount_id"`
	Balance          int    `json:"balance"`
	AvailableBalance int    `json:"available_balance"`
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
