package markets

import "testing"

func TestFormatSymbol(t *testing.T) {
	if got := formatSymbol("btc_usdt"); got != "BTCUSDT" {
		t.Fatalf("unexpected result: %s", got)
	}
}
