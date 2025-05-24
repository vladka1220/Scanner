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

func TestBybitSpotGetPrices(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/v5/market/tickers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("category") != "spot" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`{"retCode":0,"result":{"list":[{"symbol":"BTCUSDT","bid1Price":"50000","ask1Price":"50100"}]}}`))
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

func TestBybitFuturesGetPrices(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/v5/market/tickers", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("category") != "linear" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`{"retCode":0,"result":{"list":[{"symbol":"BTCUSDT","ask1Price":"50100","bid1Price":"50000"}]}}`))
	})
	h.HandleFunc("/v5/market/funding/history", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("category") != "linear" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`{"retCode":0,"result":{"list":[{"symbol":"BTCUSDT","fundingRate":"0.02","fundingRateTimestamp":1700000000000}]}}`))
	})
	cleanup := withTestServer(t, h)
	defer cleanup()

	prices := (&BybitFutures{}).GetPrices()
	p, ok := prices["BTC/USDT"]
	if !ok || !p.IsFutures {
		t.Fatalf("bad price: %+v", p)
	}
	if p.Price != 50100 || p.Volume != 50000 {
		t.Fatalf("unexpected values: %+v", p)
	}
	if p.FundingRate != 0.02 || !p.NextFundingTime.Equal(time.UnixMilli(1700000000000)) {
		t.Fatalf("unexpected funding: %+v", p)
	}
}
