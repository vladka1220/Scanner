package markets

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBybitTickers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"retCode":0,"result":{"category":"spot","list":[{"symbol":"BTCUSDT","lastPrice":"2","turnover24h":"20","volume24h":"1"}]}}`))
	}))
	defer srv.Close()
	old := bybitTickersURL
	bybitTickersURL = srv.URL
	defer func() { bybitTickersURL = old }()

	prices, err := FetchBybitTickers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := prices["BTCUSDT"]
	if !ok || p.Price != 2 || p.Volume != 1 || p.QuoteVolume != 20 {
		t.Fatalf("unexpected result: %+v", prices)
	}
}

func TestFetchBybitTickers_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := bybitTickersURL
	bybitTickersURL = srv.URL
	defer func() { bybitTickersURL = old }()
	if _, err := FetchBybitTickers(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestFetchRecentBybitTrades(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"retCode":0,"result":{"category":"spot","list":[{"price":"1","qty":"0.5","side":"Buy","time":1}]}}`))
	}))
	defer srv.Close()
	old := bybitRecentTradesURL
	bybitRecentTradesURL = srv.URL + "?symbol=%s"
	defer func() { bybitRecentTradesURL = old }()

	trades, err := FetchRecentBybitTrades("BTCUSDT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trades) != 1 || trades[0].Price != "1" {
		t.Fatalf("unexpected trades: %+v", trades)
	}
}

func TestFetchRecentBybitTrades_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := bybitRecentTradesURL
	bybitRecentTradesURL = srv.URL + "?symbol=%s"
	defer func() { bybitRecentTradesURL = old }()

	if _, err := FetchRecentBybitTrades("BTCUSDT"); err == nil {
		t.Fatalf("expected error")
	}
}
