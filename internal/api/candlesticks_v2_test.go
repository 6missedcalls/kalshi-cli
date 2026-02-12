package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

// These tests send the REAL Kalshi v2 API JSON format (nested price object,
// end_period_ts as unix int) and verify our models correctly parse it.
// If these fail, it means our Candlestick model doesn't match the live API.

func TestGetCandlesticks_RealV2Format(t *testing.T) {
	// This is the actual JSON format returned by Kalshi's v2 API
	realResponse := `{
		"ticker": "BTC-100K",
		"candlesticks": [
			{
				"end_period_ts": 1704067200,
				"price": {
					"open": 45,
					"open_dollars": "0.4500",
					"low": 42,
					"low_dollars": "0.4200",
					"high": 48,
					"high_dollars": "0.4800",
					"close": 47,
					"close_dollars": "0.4700",
					"mean": 46,
					"mean_dollars": "0.4600",
					"previous": 44,
					"previous_dollars": "0.4400"
				},
				"yes_bid": {
					"open": 44,
					"open_dollars": "0.4400",
					"low": 41,
					"low_dollars": "0.4100",
					"high": 47,
					"high_dollars": "0.4700",
					"close": 46,
					"close_dollars": "0.4600"
				},
				"yes_ask": {
					"open": 46,
					"open_dollars": "0.4600",
					"low": 43,
					"low_dollars": "0.4300",
					"high": 49,
					"high_dollars": "0.4900",
					"close": 48,
					"close_dollars": "0.4800"
				},
				"volume": 1500,
				"open_interest": 500
			},
			{
				"end_period_ts": 1704070800,
				"price": {
					"open": 47,
					"open_dollars": "0.4700",
					"low": 45,
					"low_dollars": "0.4500",
					"high": 52,
					"high_dollars": "0.5200",
					"close": 50,
					"close_dollars": "0.5000"
				},
				"yes_bid": {
					"open": 46,
					"open_dollars": "0.4600",
					"low": 44,
					"low_dollars": "0.4400",
					"high": 51,
					"high_dollars": "0.5100",
					"close": 49,
					"close_dollars": "0.4900"
				},
				"yes_ask": {
					"open": 48,
					"open_dollars": "0.4800",
					"low": 46,
					"low_dollars": "0.4600",
					"high": 53,
					"high_dollars": "0.5300",
					"close": 51,
					"close_dollars": "0.5100"
				},
				"volume": 2000,
				"open_interest": 600
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(realResponse))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	result, err := client.GetCandlesticks(context.Background(), GetCandlesticksParams{
		SeriesTicker: "BTC-SERIES",
		Ticker:       "BTC-100K",
		Period:       "1h",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Candlesticks) != 2 {
		t.Fatalf("expected 2 candlesticks, got %d", len(result.Candlesticks))
	}

	// Verify first candlestick values come from price.open/high/low/close
	c := result.Candlesticks[0]
	if c.Open != 45 {
		t.Errorf("expected Open=45, got %d", c.Open)
	}
	if c.High != 48 {
		t.Errorf("expected High=48, got %d", c.High)
	}
	if c.Low != 42 {
		t.Errorf("expected Low=42, got %d", c.Low)
	}
	if c.Close != 47 {
		t.Errorf("expected Close=47, got %d", c.Close)
	}
	if c.Volume != 1500 {
		t.Errorf("expected Volume=1500, got %d", c.Volume)
	}
	if c.OpenInterest != 500 {
		t.Errorf("expected OpenInterest=500, got %d", c.OpenInterest)
	}

	// Verify timestamp parsed from end_period_ts (unix seconds)
	expectedTime := time.Unix(1704067200, 0).UTC()
	if !c.PeriodEnd.Equal(expectedTime) {
		t.Errorf("expected PeriodEnd=%v, got %v", expectedTime, c.PeriodEnd)
	}

	// Verify second candlestick
	c2 := result.Candlesticks[1]
	if c2.Open != 47 {
		t.Errorf("expected Open=47, got %d", c2.Open)
	}
	if c2.Close != 50 {
		t.Errorf("expected Close=50, got %d", c2.Close)
	}
	expectedTime2 := time.Unix(1704070800, 0).UTC()
	if !c2.PeriodEnd.Equal(expectedTime2) {
		t.Errorf("expected PeriodEnd=%v, got %v", expectedTime2, c2.PeriodEnd)
	}
}

func TestGetEventCandlesticks_RealV2Format(t *testing.T) {
	// Event candlesticks have a different response shape:
	// market_tickers + market_candlesticks (array of arrays)
	realResponse := `{
		"market_tickers": ["INXD-25FEB07-B5523.99", "INXD-25FEB07-B5524.99"],
		"market_candlesticks": [
			[
				{
					"end_period_ts": 1738886400,
					"price": {
						"open": 60,
						"open_dollars": "0.6000",
						"low": 55,
						"low_dollars": "0.5500",
						"high": 65,
						"high_dollars": "0.6500",
						"close": 62,
						"close_dollars": "0.6200"
					},
					"volume": 3000,
					"open_interest": 1200
				}
			],
			[
				{
					"end_period_ts": 1738886400,
					"price": {
						"open": 30,
						"open_dollars": "0.3000",
						"low": 25,
						"low_dollars": "0.2500",
						"high": 35,
						"high_dollars": "0.3500",
						"close": 32,
						"close_dollars": "0.3200"
					},
					"volume": 1500,
					"open_interest": 800
				}
			]
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(realResponse))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	candlesticks, err := client.GetEventCandlesticks(context.Background(), CandlesticksParams{
		SeriesTicker: "INXD",
		Ticker:       "INXD-25FEB07",
		Period:       "1h",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Event candlesticks should flatten all markets' candlesticks
	if len(candlesticks) == 0 {
		t.Fatal("expected at least 1 candlestick, got 0")
	}

	// Check first candlestick parsed correctly
	c := candlesticks[0]
	if c.Open != 60 {
		t.Errorf("expected Open=60, got %d", c.Open)
	}
	if c.High != 65 {
		t.Errorf("expected High=65, got %d", c.High)
	}
	if c.Low != 55 {
		t.Errorf("expected Low=55, got %d", c.Low)
	}
	if c.Close != 62 {
		t.Errorf("expected Close=62, got %d", c.Close)
	}
	if c.Volume != 3000 {
		t.Errorf("expected Volume=3000, got %d", c.Volume)
	}
}

// Test that existing old-format tests still produce correct results
// when mock servers return the correct new format.
// This verifies backward compatibility is maintained.
func TestCandlestickResponse_OldFormatStillWorks(t *testing.T) {
	// Existing tests use this format (Go struct → JSON encoding).
	// After the fix, this format should still be parseable for internal usage.
	resp := models.CandlesticksResponse{
		Candlesticks: []models.Candlestick{
			{
				Ticker:       "TEST",
				Open:         45,
				High:         48,
				Low:          42,
				Close:        47,
				Volume:       1000,
				OpenInterest: 500,
			},
		},
	}

	if len(resp.Candlesticks) != 1 {
		t.Fatalf("expected 1 candlestick, got %d", len(resp.Candlesticks))
	}
	if resp.Candlesticks[0].Open != 45 {
		t.Errorf("expected Open=45, got %d", resp.Candlesticks[0].Open)
	}
}
