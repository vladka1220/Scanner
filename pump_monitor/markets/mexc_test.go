package markets

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchMEXCTickers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"symbol":"BTCUSDT","lastPrice":"5","volume":"6"}]`))
	}))
	defer srv.Close()
	old := mexcTickersURL
	mexcTickersURL = srv.URL
	defer func() { mexcTickersURL = old }()

	prices, err := FetchMEXCTickers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := prices["BTCUSDT"]
	if !ok || p.Price != 5 || p.Volume != 6 {
		t.Fatalf("unexpected prices: %+v", prices)
	}
}

func TestFetchMEXCTickers_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := mexcTickersURL
	mexcTickersURL = srv.URL
	defer func() { mexcTickersURL = old }()

	if _, err := FetchMEXCTickers(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestFetchRecentTrades(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"price":"1","qty":"2","isBuyerMaker":false,"time":1}]`))
	}))
	defer srv.Close()
	old := mexcRecentTradesURL
	mexcRecentTradesURL = srv.URL + "?symbol=%s"
	defer func() { mexcRecentTradesURL = old }()

	trades, err := FetchRecentTrades("BTCUSDT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trades) != 1 || trades[0].Price != "1" {
		t.Fatalf("unexpected trades: %+v", trades)
	}
}

func TestFetchRecentTrades_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := mexcRecentTradesURL
	mexcRecentTradesURL = srv.URL + "?symbol=%s"
	defer func() { mexcRecentTradesURL = old }()

	if _, err := FetchRecentTrades("BTCUSDT"); err == nil {
		t.Fatalf("expected error")
	}
}
