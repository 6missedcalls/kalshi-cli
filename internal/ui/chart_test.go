package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRenderCandlestickChart_Empty(t *testing.T) {
	out := captureOutput(func() {
		RenderCandlestickChart(nil, "Test")
	})

	if !strings.Contains(out, "No candlestick data") {
		t.Errorf("expected 'No candlestick data' message, got: %s", out)
	}
}

func TestRenderCandlestickChart_SingleCandle(t *testing.T) {
	candles := []CandleData{
		{Label: "12:00", Open: 50, High: 60, Low: 40, Close: 55, Volume: 100},
	}

	out := captureOutput(func() {
		RenderCandlestickChart(candles, "Single")
	})

	if out == "" {
		t.Fatal("expected output, got empty string")
	}
	if !strings.Contains(out, "Single") {
		t.Error("expected title in output")
	}
	if !strings.Contains(out, "$0.55") {
		t.Errorf("expected last close $0.55 in output, got: %s", out)
	}
}

func TestRenderCandlestickChart_MultipleCandles(t *testing.T) {
	candles := []CandleData{
		{Label: "10:00", Open: 40, High: 50, Low: 35, Close: 45, Volume: 200},
		{Label: "10:15", Open: 45, High: 52, Low: 42, Close: 48, Volume: 180},
		{Label: "10:30", Open: 48, High: 55, Low: 44, Close: 52, Volume: 250},
		{Label: "10:45", Open: 52, High: 58, Low: 48, Close: 55, Volume: 310},
		{Label: "11:00", Open: 55, High: 60, Low: 50, Close: 53, Volume: 280},
		{Label: "11:15", Open: 53, High: 57, Low: 47, Close: 50, Volume: 220},
		{Label: "11:30", Open: 50, High: 54, Low: 45, Close: 48, Volume: 190},
		{Label: "11:45", Open: 48, High: 56, Low: 46, Close: 54, Volume: 340},
		{Label: "12:00", Open: 54, High: 62, Low: 52, Close: 60, Volume: 400},
		{Label: "12:15", Open: 60, High: 65, Low: 55, Close: 58, Volume: 150},
	}

	out := captureOutput(func() {
		RenderCandlestickChart(candles, "Multi")
	})

	if !strings.Contains(out, "Multi") {
		t.Error("expected title in output")
	}
	if !strings.Contains(out, "$") {
		t.Error("expected dollar-formatted price labels")
	}
	if !strings.Contains(out, "Vol") {
		t.Error("expected volume sparkline row")
	}
	if !strings.Contains(out, "10:00") {
		t.Error("expected first x-axis label")
	}
	if !strings.Contains(out, "12:15") {
		t.Error("expected last x-axis label")
	}
}

func TestRenderCandlestickChart_FlatCandle(t *testing.T) {
	candles := []CandleData{
		{Label: "09:00", Open: 50, High: 50, Low: 50, Close: 50, Volume: 10},
	}

	out := captureOutput(func() {
		RenderCandlestickChart(candles, "Flat")
	})

	if out == "" {
		t.Fatal("expected output for flat candle")
	}
}

func TestPriceToRow(t *testing.T) {
	tests := []struct {
		name     string
		price    int
		min, max int
		wantRow  int
	}{
		{"at max", 100, 0, 100, 0},
		{"at min", 0, 0, 100, chartHeight - 1},
		{"midpoint", 50, 0, 100, chartHeight / 2},
		{"above max clamps", 110, 0, 100, 0},
		{"below min clamps", -5, 0, 100, chartHeight - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := priceToRow(tt.price, tt.min, tt.max)
			if got != tt.wantRow {
				t.Errorf("priceToRow(%d, %d, %d) = %d, want %d", tt.price, tt.min, tt.max, got, tt.wantRow)
			}
		})
	}
}

func TestRowToPrice(t *testing.T) {
	// Top row should be max price
	got := rowToPrice(0, 0, 100)
	if got != 100 {
		t.Errorf("rowToPrice(0, 0, 100) = %d, want 100", got)
	}

	// Bottom row should be min price
	got = rowToPrice(chartHeight-1, 0, 100)
	if got != 0 {
		t.Errorf("rowToPrice(%d, 0, 100) = %d, want 0", chartHeight-1, got)
	}
}

func TestVolumeBar(t *testing.T) {
	tests := []struct {
		vol, maxVol int
		want        rune
	}{
		{0, 100, '▁'},
		{100, 100, '█'},
		{0, 0, '▁'},
		{50, 100, '▄'},
	}

	for _, tt := range tests {
		got := volumeBar(tt.vol, tt.maxVol)
		if got != tt.want {
			t.Errorf("volumeBar(%d, %d) = %c, want %c", tt.vol, tt.maxVol, got, tt.want)
		}
	}
}

func TestPriceBounds(t *testing.T) {
	candles := []CandleData{
		{Low: 30, High: 50},
		{Low: 20, High: 60},
		{Low: 25, High: 55},
	}

	lo, hi := priceBounds(candles)
	// min Low = 20, padding -1 = 19
	if lo != 19 {
		t.Errorf("priceBounds min = %d, want 19", lo)
	}
	// max High = 60, padding +1 = 61
	if hi != 61 {
		t.Errorf("priceBounds max = %d, want 61", hi)
	}
}

func TestRenderCandlestickChart_ExceedsMaxCandles(t *testing.T) {
	candles := make([]CandleData, 60)
	for i := range candles {
		candles[i] = CandleData{
			Label:  "T",
			Open:   50 + i,
			High:   60 + i,
			Low:    40 + i,
			Close:  55 + i,
			Volume: 100,
		}
	}

	out := captureOutput(func() {
		RenderCandlestickChart(candles, "Many")
	})

	if out == "" {
		t.Fatal("expected output for many candles")
	}
}
