package core

import (
	"testing"

	"basis_go/types"
)

func TestCollect(t *testing.T) {
	sources := map[string]func() map[string]types.ExchangePrice{
		"ex1": func() map[string]types.ExchangePrice {
			return map[string]types.ExchangePrice{
				"BTC/USDT": {Price: 1},
			}
		},
		"ex2": func() map[string]types.ExchangePrice {
			return map[string]types.ExchangePrice{
				"BTC/USDT": {Price: 2},
				"ETH/USDT": {Price: 3},
			}
		},
	}

	result := collect(sources)

	if len(result) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(result))
	}
	if len(result["BTC/USDT"]) != 2 {
		t.Fatalf("expected both exchanges for BTC/USDT, got %v", result["BTC/USDT"])
	}
	if price := result["ETH/USDT"]["ex2"].Price; price != 3 {
		t.Fatalf("expected price 3, got %v", price)
	}
}
