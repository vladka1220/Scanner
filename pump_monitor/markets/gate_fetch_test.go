package markets

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchGateTickers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"currency_pair":"BTC_USDT","last":"3","base_volume":"4","quote_volume":"4"}]`))
	}))
	defer srv.Close()
	old := gateTickersURL
	gateTickersURL = srv.URL
	defer func() { gateTickersURL = old }()

	prices, err := FetchGateTickers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := prices["BTCUSDT"]
	if !ok || p.Price != 3 || p.Volume != 4 || p.QuoteVolume != 4 {
		t.Fatalf("unexpected prices: %+v", prices)
	}
}

func TestFetchGateTickers_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := gateTickersURL
	gateTickersURL = srv.URL
	defer func() { gateTickersURL = old }()

	if _, err := FetchGateTickers(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestFetchRecentGateTrades(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"price":"1","amount":"2","create_time_ms":1,"side":"buy"}]`))
	}))
	defer srv.Close()
	old := gateRecentTradesURL
	gateRecentTradesURL = srv.URL + "?currency_pair=%s"
	defer func() { gateRecentTradesURL = old }()

	trades, err := FetchRecentGateTrades("BTC_USDT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trades) != 1 || trades[0].Price != "1" {
		t.Fatalf("unexpected trades: %+v", trades)
	}
}

func TestFetchRecentGateTrades_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := gateRecentTradesURL
	gateRecentTradesURL = srv.URL + "?currency_pair=%s"
	defer func() { gateRecentTradesURL = old }()

	if _, err := FetchRecentGateTrades("BTC_USDT"); err == nil {
		t.Fatalf("expected error")
	}
}
