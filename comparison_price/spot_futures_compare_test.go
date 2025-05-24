package comparison_price

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"basis_go/funding"
	"basis_go/types"
)

func TestCompareSpotFutures(t *testing.T) {
	now := time.Now()

	spotPrices := types.TokenPrices{
		"BTC/USDT": {
			"ex_spot": {Price: 100},
		},
		"ETH/USDT": {
			"ex_spot": {Price: 50},
		},
	}

	btcFundingTime := now.Add(2 * time.Hour).UnixMilli()
	ethFundingTime := now.Add(1 * time.Hour).UnixMilli()

	futuresPrices := types.TokenPrices{
		"BTC/USDT": {
			"ex_fut": {
				Price:            110,
				IsFutures:        true,
				FundingRate:      0.001,
				FundingTime:      btcFundingTime,
				FundingIntervalH: 8,
			},
		},
		"ETH/USDT": {
			"ex_fut": {
				Price:            65,
				IsFutures:        true,
				FundingRate:      0.002,
				FundingTime:      ethFundingTime,
				FundingIntervalH: 8,
			},
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

	expectedFundingETH := funding.FormatFunding(0.002, ethFundingTime, 65, 1.0, 8)
	expectedFundingBTC := funding.FormatFunding(0.001, btcFundingTime, 110, 1.0, 8)

	if !strings.Contains(out, expectedFundingETH) {
		t.Fatalf("ETH funding formatting mismatch: %s", out)
	}
	if !strings.Contains(out, expectedFundingBTC) {
		t.Fatalf("BTC funding formatting mismatch: %s", out)
	}

	ethIndex := strings.Index(out, "[ETH/USDT]")
	btcIndex := strings.Index(out, "[BTC/USDT]")
	if ethIndex == -1 || btcIndex == -1 {
		t.Fatalf("missing token output: %s", out)
	}
	if ethIndex > btcIndex {
		t.Fatalf("results not sorted by spread: %s", out)
	}
}
