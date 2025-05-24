package core

import (
	"testing"

	"basis_go/types"
)

// TestCollectMergesTokenPrices verifies that collect merges prices from multiple sources
// and returns a combined map of TokenPrices.
func TestCollectMergesTokenPrices(t *testing.T) {
	sources := map[string]func() map[string]types.ExchangePrice{
		"ex1": func() map[string]types.ExchangePrice {
			return map[string]types.ExchangePrice{
				"BTC/USDT": {Price: 1},
				"ETH/USDT": {Price: 5},
			}
		},
		"ex2": func() map[string]types.ExchangePrice {
			return map[string]types.ExchangePrice{
				"ETH/USDT": {Price: 6},
				"XRP/USDT": {Price: 7},
			}
		},
	}

	result := collect(sources)

	if len(result) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(result))
	}

	if price := result["BTC/USDT"]["ex1"].Price; price != 1 {
		t.Fatalf("expected BTC/USDT price 1 from ex1, got %v", price)
	}

	eth, ok := result["ETH/USDT"]
	if !ok || len(eth) != 2 {
		t.Fatalf("expected ETH/USDT prices from both exchanges, got %v", eth)
	}

	if price := eth["ex1"].Price; price != 5 {
		t.Fatalf("expected ETH/USDT price 5 from ex1, got %v", price)
	}
	if price := eth["ex2"].Price; price != 6 {
		t.Fatalf("expected ETH/USDT price 6 from ex2, got %v", price)
	}

	if price := result["XRP/USDT"]["ex2"].Price; price != 7 {
		t.Fatalf("expected XRP/USDT price 7 from ex2, got %v", price)
	}
}
