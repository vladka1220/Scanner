package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	data, err := FetchJSON(srv.URL)
	if err != nil {
		t.Fatalf("FetchJSON returned error: %v", err)
	}
	if string(data) != "{\"ok\":true}" {
		t.Fatalf("unexpected body: %s", string(data))
	}
}

func TestNormalizeSymbol(t *testing.T) {
	tests := map[string]string{
		"BTCUSDT":   "BTC/USDT",
		"eth_usdc":  "ETH/USDC",
		"DOGE/USDT": "DOGE/USDT",
	}
	for in, exp := range tests {
		if got := NormalizeSymbol(in); got != exp {
			t.Errorf("NormalizeSymbol(%s)=%s, want %s", in, got, exp)
		}
	}
}

func TestParseFloatInt64(t *testing.T) {
	f, err := ParseFloat("1.23")
	if err != nil || f != 1.23 {
		t.Errorf("ParseFloat failed, got %v err %v", f, err)
	}
	if _, err := ParseFloat("bad"); err == nil {
		t.Error("expected error for ParseFloat('bad')")
	}

	i, err := ParseInt64("42")
	if err != nil || i != 42 {
		t.Errorf("ParseInt64 failed, got %v err %v", i, err)
	}
	if _, err := ParseInt64("bad"); err == nil {
		t.Error("expected error for ParseInt64('bad')")
	}
}
