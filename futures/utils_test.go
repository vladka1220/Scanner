package futures

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"basis_go/types"
)

func TestComparePricesOrdering(t *testing.T) {
	now := time.Now()
	prices := types.TokenPrices{
		"BTC/USDT": {
			"Binance Futures": {Price: 100, IsFutures: true, FundingRate: 0.001, FundingTime: now.Add(1 * time.Hour).UnixMilli()},
			"MEXC Futures":    {Price: 105, IsFutures: true, FundingRate: 0.002, FundingTime: now.Add(1 * time.Hour).UnixMilli()},
			"Gate Futures":    {Price: 103, IsFutures: true, FundingRate: 0.0015, FundingTime: now.Add(1 * time.Hour).UnixMilli()},
		},
	}

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	ComparePrices(prices)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()

	re := regexp.MustCompile(`Спред: ([0-9.]+)%`)
	matches := re.FindAllStringSubmatch(out, -1)
	if len(matches) != 3 {
		t.Fatalf("expected 3 spreads, got %d", len(matches))
	}
	prev, _ := strconv.ParseFloat(matches[0][1], 64)
	if prev < 0.25 || prev > 100 {
		t.Fatalf("spread out of range: %v", prev)
	}
	for i := 1; i < len(matches); i++ {
		cur, _ := strconv.ParseFloat(matches[i][1], 64)
		if cur < 0.25 || cur > 100 {
			t.Fatalf("spread %d out of range: %v", i, cur)
		}
		if cur > prev {
			t.Fatalf("spreads not sorted desc: %v then %v", prev, cur)
		}
		prev = cur
	}
}
