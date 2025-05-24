package spot

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"testing"

	"basis_go/types"
)

func TestCompareSpotPricesOrdering(t *testing.T) {
	prices := types.TokenPrices{
		"BTC/USDT": {
			"Binance Spot": {Price: 100},
			"MEXC Spot":    {Price: 105},
			"Gate.io Spot": {Price: 103},
		},
	}

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	CompareSpotPrices(prices)
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
