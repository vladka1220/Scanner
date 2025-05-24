package markets

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchBinanceTickers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"symbol":"BTCUSDT","lastPrice":"1.5","quoteVolume":"10"}]`))
	}))
	defer srv.Close()
	old := binanceTickersURL
	binanceTickersURL = srv.URL
	defer func() { binanceTickersURL = old }()

	prices, err := FetchBinanceTickers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, ok := prices["BTCUSDT"]
	if !ok || p.Price != 1.5 || p.Volume != 10 || p.QuoteVolume != 10 {
		t.Fatalf("unexpected result: %+v", prices)
	}
}

func TestFetchBinanceTickers_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bad"))
	}))
	defer srv.Close()
	old := binanceTickersURL
	binanceTickersURL = srv.URL
	defer func() { binanceTickersURL = old }()

	if _, err := FetchBinanceTickers(); err == nil {
		t.Fatalf("expected error")
	}
}
