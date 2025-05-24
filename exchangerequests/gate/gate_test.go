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

func TestGateSpotGetPrices(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/api/v4/spot/tickers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"currency_pair":"BTC_USDT","last":"50100","base_volume":"50000"}]`))
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

func TestGateFuturesGetPrices(t *testing.T) {
	h := http.NewServeMux()
	h.HandleFunc("/api/v4/futures/usdt/tickers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"contract":"BTC_USDT","last":"50100","volume_24h_quote":"50000"}]`))
	})
	h.HandleFunc("/api/v4/futures/usdt/funding_rates", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"contract":"BTC_USDT","funding_rate":"0.03","funding_time":"1700000000000"}]}`))
	})
	cleanup := withTestServer(t, h)
	defer cleanup()

	prices := GetGateFuturesPrices()
	p, ok := prices["BTC/USDT"]
	if !ok || !p.IsFutures {
		t.Fatalf("bad price: %+v", p)
	}
	if p.Price != 50100 || p.Volume != 50000 {
		t.Fatalf("unexpected values: %+v", p)
	}
	if p.FundingRate != 0.03 || !p.NextFundingTime.Equal(time.UnixMilli(1700000000000)) {
		t.Fatalf("unexpected funding: %+v", p)
	}
}
