package core

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"basis_go/types"
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
