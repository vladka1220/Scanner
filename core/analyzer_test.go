package core

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"basis_go/types"
	"time"
)

func TestCompareAll(t *testing.T) {
	prices := types.TokenPrices{
		"BTC/USDT": {
			"ex1": {Price: 100},
			"ex2": {Price: 103},
		},
	}

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	CompareAll(prices)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "BTC/USDT") {
		t.Fatalf("output missing token: %s", out)
	}
	if !strings.Contains(out, "Спред") {
		t.Fatalf("output missing spread info: %s", out)
	}
}

func TestCompareAll_SortingFunding(t *testing.T) {
	now := time.Now()
	prices := types.TokenPrices{
		"ETH/USDT": {
			"SpotEx": {Price: 1000},
			"FutEx": {
				Price:            1100,
				IsFutures:        true,
				FundingRate:      0.001,
				FundingIntervalH: 8,
				NextFundingTime:  now.Add(time.Hour),
			},
		},
		"BTC/USDT": {
			"FutEx": {
				Price:            21000,
				IsFutures:        true,
				FundingRate:      -0.002,
				FundingIntervalH: 8,
				NextFundingTime:  now.Add(2 * time.Hour),
			},
			"SpotEx": {Price: 20000},
		},
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	CompareAll(prices)
	w.Close()
	os.Stdout = old
	io.Copy(&buf, r)
	out := buf.String()

	btc := strings.Index(out, "[BTC/USDT]")
	eth := strings.Index(out, "[ETH/USDT]")
	if btc == -1 || eth == -1 || btc > eth {
		t.Fatalf("tokens not sorted: %s", out)
	}
	if !strings.Contains(out, "Фьючерс") || !strings.Contains(out, "Фандинг") {
		t.Fatalf("futures funding info missing: %s", out)
	}
}
