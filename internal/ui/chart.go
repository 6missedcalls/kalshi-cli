package ui

import (
	"fmt"
	"math"
	"strings"
)

// CandleData holds OHLCV data for chart rendering.
// Uses int cents to avoid importing models package.
type CandleData struct {
	Label                        string
	Open, High, Low, Close       int
	Volume                       int
}

const (
	chartHeight    = 16
	maxChartCandles = 40
)

var volumeBars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// RenderCandlestickChart prints an ASCII candlestick chart to stdout.
func RenderCandlestickChart(candles []CandleData, title string) {
	if len(candles) == 0 {
		fmt.Println(MutedStyle.Render("  No candlestick data to chart."))
		return
	}

	// Trim to most recent candles if too many
	visible := candles
	if len(visible) > maxChartCandles {
		visible = visible[len(visible)-maxChartCandles:]
	}

	priceMin, priceMax := priceBounds(visible)
	if priceMin == priceMax {
		priceMax = priceMin + 1
	}

	// Summary header
	fmt.Println()
	fmt.Print("  " + TitleStyle.Render(title))
	lastClose := visible[len(visible)-1].Close
	firstOpen := visible[0].Open
	change := lastClose - firstOpen
	changePct := 0.0
	if firstOpen != 0 {
		changePct = float64(change) / float64(firstOpen) * 100
	}
	summary := fmt.Sprintf("  Last: %s", FormatPrice(lastClose))
	if change >= 0 {
		summary += "  " + PriceUpStyle.Render(fmt.Sprintf("+%s (%.1f%%)", FormatPrice(change), changePct))
	} else {
		summary += "  " + PriceDownStyle.Render(fmt.Sprintf("%s (%.1f%%)", FormatPrice(change), changePct))
	}
	fmt.Println(summary)
	fmt.Println()

	// Build chart grid
	grid := buildGrid(visible, priceMin, priceMax)

	// Render rows with y-axis labels
	labelInterval := labelStep(chartHeight)
	for row := 0; row < chartHeight; row++ {
		price := rowToPrice(row, priceMin, priceMax)
		if row == 0 || row == chartHeight-1 || row%labelInterval == 0 {
			fmt.Printf("  %7s │", FormatPrice(price))
		} else {
			fmt.Print("          │")
		}
		for col := 0; col < len(visible); col++ {
			fmt.Print(grid[row][col])
		}
		fmt.Println()
	}

	// X-axis line
	fmt.Print("          └")
	fmt.Println(strings.Repeat("─", len(visible)*2))

	// X-axis labels
	renderXLabels(visible)

	// Volume sparkline
	renderVolumeLine(visible)
	fmt.Println()
}

func priceBounds(candles []CandleData) (int, int) {
	lo := math.MaxInt
	hi := math.MinInt
	for _, c := range candles {
		if c.Low < lo {
			lo = c.Low
		}
		if c.High > hi {
			hi = c.High
		}
	}
	// Add small padding (1 cent each side)
	if lo > 0 {
		lo--
	}
	hi++
	return lo, hi
}

func buildGrid(candles []CandleData, priceMin, priceMax int) [][]string {
	grid := make([][]string, chartHeight)
	for r := range grid {
		grid[r] = make([]string, len(candles))
		for c := range grid[r] {
			grid[r][c] = "  "
		}
	}

	for col, candle := range candles {
		highRow := priceToRow(candle.High, priceMin, priceMax)
		lowRow := priceToRow(candle.Low, priceMin, priceMax)

		openRow := priceToRow(candle.Open, priceMin, priceMax)
		closeRow := priceToRow(candle.Close, priceMin, priceMax)

		// Ensure body top <= body bottom (row 0 = top)
		bodyTop := openRow
		bodyBot := closeRow
		if bodyTop > bodyBot {
			bodyTop, bodyBot = bodyBot, bodyTop
		}

		bullish := candle.Close >= candle.Open
		bodyStyle := PriceDownStyle
		if bullish {
			bodyStyle = PriceUpStyle
		}

		for row := highRow; row <= lowRow; row++ {
			if row >= bodyTop && row <= bodyBot {
				if bodyTop == bodyBot {
					// Doji / flat candle
					grid[row][col] = bodyStyle.Render("─ ")
				} else {
					grid[row][col] = bodyStyle.Render("┃ ")
				}
			} else {
				// Wick
				grid[row][col] = MutedStyle.Render("│ ")
			}
		}
	}

	return grid
}

func priceToRow(price, priceMin, priceMax int) int {
	priceRange := priceMax - priceMin
	if priceRange == 0 {
		return chartHeight / 2
	}
	// row 0 = priceMax (top), row chartHeight-1 = priceMin (bottom)
	ratio := float64(priceMax-price) / float64(priceRange)
	row := int(math.Round(ratio * float64(chartHeight-1)))
	if row < 0 {
		return 0
	}
	if row >= chartHeight {
		return chartHeight - 1
	}
	return row
}

func rowToPrice(row, priceMin, priceMax int) int {
	priceRange := priceMax - priceMin
	if chartHeight <= 1 {
		return priceMin
	}
	return priceMax - (row * priceRange / (chartHeight - 1))
}

func labelStep(height int) int {
	if height <= 4 {
		return 1
	}
	return height / 4
}

func renderXLabels(candles []CandleData) {
	if len(candles) == 0 {
		return
	}

	type labelPos struct {
		col   int
		label string
	}
	var labels []labelPos

	if len(candles) == 1 {
		labels = append(labels, labelPos{0, candles[0].Label})
	} else if len(candles) <= 5 {
		labels = append(labels, labelPos{0, candles[0].Label})
		labels = append(labels, labelPos{len(candles) - 1, candles[len(candles)-1].Label})
	} else {
		labels = append(labels, labelPos{0, candles[0].Label})
		mid := len(candles) / 2
		labels = append(labels, labelPos{mid, candles[mid].Label})
		labels = append(labels, labelPos{len(candles) - 1, candles[len(candles)-1].Label})
	}

	// Extra space after last candle for label overflow
	maxLabelLen := 0
	for _, lp := range labels {
		if len(lp.label) > maxLabelLen {
			maxLabelLen = len(lp.label)
		}
	}
	totalWidth := len(candles)*2 + 11 + maxLabelLen
	buf := make([]byte, totalWidth)
	for i := range buf {
		buf[i] = ' '
	}

	// Place labels, skipping if they'd overlap a previous one
	lastEnd := 0
	for _, lp := range labels {
		offset := 11 + lp.col*2
		lbl := lp.label
		end := offset + len(lbl)
		if end > totalWidth {
			end = totalWidth
			lbl = lbl[:end-offset]
		}
		if offset < lastEnd {
			continue // skip overlapping label
		}
		copy(buf[offset:end], lbl)
		lastEnd = end + 1
	}

	fmt.Println(string(buf))
}

func renderVolumeLine(candles []CandleData) {
	if len(candles) == 0 {
		return
	}

	maxVol := 0
	for _, c := range candles {
		if c.Volume > maxVol {
			maxVol = c.Volume
		}
	}

	fmt.Print("  " + MutedStyle.Render("Vol") + "     ")
	for _, c := range candles {
		bar := volumeBar(c.Volume, maxVol)
		if c.Close >= c.Open {
			fmt.Print(PriceUpStyle.Render(string(bar)) + " ")
		} else {
			fmt.Print(PriceDownStyle.Render(string(bar)) + " ")
		}
	}
	fmt.Println()
}

func volumeBar(vol, maxVol int) rune {
	if maxVol == 0 || vol == 0 {
		return volumeBars[0]
	}
	idx := int(float64(vol) / float64(maxVol) * float64(len(volumeBars)-1))
	if idx < 0 {
		return volumeBars[0]
	}
	if idx >= len(volumeBars) {
		return volumeBars[len(volumeBars)-1]
	}
	return volumeBars[idx]
}
