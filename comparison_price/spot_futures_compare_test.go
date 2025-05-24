package comparison_price

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"basis_go/types"
)

func TestCompareSpotFutures(t *testing.T) {
	spotPrices := types.TokenPrices{
		"BTC/USDT": {
			"ex1": {Price: 100},
		},
	}
	futuresPrices := types.TokenPrices{
		"BTC/USDT": {
			"ex2": {Price: 102, IsFutures: true, FundingRate: 0.001, FundingTime: time.Now().Add(1 * time.Hour).UnixMilli(), FundingIntervalH: 8},
		},
	}

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	CompareSpotFutures(spotPrices, futuresPrices)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "BTC/USDT") {
		t.Fatalf("output missing token: %s", out)
	}
	if !strings.Contains(out, "0.100%") {
		t.Fatalf("output missing funding info: %s", out)
	}
}
