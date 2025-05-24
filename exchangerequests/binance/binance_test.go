package exchangerequests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type rewriteTransport struct {
	base *url.URL
	rt   http.RoundTripper
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.base.Scheme
	req.URL.Host = t.base.Host
	return t.rt.RoundTrip(req)
}

func withTestServer(t *testing.T, handler http.Handler) func() {
	srv := httptest.NewServer(handler)
	u, _ := url.Parse(srv.URL)
	old := http.DefaultTransport
	http.DefaultTransport = rewriteTransport{base: u, rt: old}
	return func() {
		http.DefaultTransport = old
		srv.Close()
	}
}

func TestBinanceSpotGetPrices(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/api/v3/ticker/bookTicker", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"symbol":"BTCUSDT","bidPrice":"50000","askPrice":"50100"}]`))
	})
	cleanup := withTestServer(t, h)
	defer cleanup()

	prices := GetSpotPrices()
	p, ok := prices["BTC/USDT"]
	if !ok || p.IsFutures {
		t.Fatalf("bad price: %+v", p)
	}
	if p.Price != 50100 || p.Volume != 50000 {
		t.Fatalf("unexpected values: %+v", p)
	}
}

func TestBinanceFuturesGetPrices(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/fapi/v1/ticker/bookTicker", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"symbol":"BTCUSDT","askPrice":"50100","bidPrice":"50000"}]`))
	})
	h.HandleFunc("/fapi/v1/premiumIndex", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"symbol":"BTCUSDT","lastFundingRate":"0.01","nextFundingTime":1700000000000}]`))
	})
	cleanup := withTestServer(t, h)
	defer cleanup()

	prices := GetFuturesPrices()
	p, ok := prices["BTC/USDT"]
	if !ok || !p.IsFutures {
		t.Fatalf("bad price: %+v", p)
	}
	if p.Price != 50100 || p.Volume != 50000 {
		t.Fatalf("unexpected values: %+v", p)
	}
	if p.FundingRate != 0.01 || !p.NextFundingTime.Equal(time.UnixMilli(1700000000000)) {
		t.Fatalf("unexpected funding: %+v", p)
	}
}
