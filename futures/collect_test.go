package futures

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"basis_go/types"
)

func TestCollectFuturesPrices_WithFetchers(t *testing.T) {
	fetchers := map[string]func() map[string]types.ExchangePrice{
		"f1": func() map[string]types.ExchangePrice {
			return map[string]types.ExchangePrice{
				"BTC/USDT": {Price: 1},
			}
		},
		"f2": func() map[string]types.ExchangePrice {
			return map[string]types.ExchangePrice{
				"BTC/USDT": {Price: 2},
				"ETH/USDT": {Price: 3},
			}
		},
		"empty": func() map[string]types.ExchangePrice { return nil },
	}

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	res := CollectFuturesPrices(fetchers)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "[empty] ФЬЮЧЕРС НЕ ОТДАЁТ ДАННЫХ") {
		t.Fatalf("expected warning for empty fetcher, got: %s", out)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(res))
	}
	if len(res["BTC/USDT"]) != 2 {
		t.Fatalf("expected BTC/USDT on 2 exchanges, got %v", res["BTC/USDT"])
	}
	if p := res["ETH/USDT"]["f2"].Price; p != 3 {
		t.Fatalf("expected ETH price 3, got %v", p)
	}
}
