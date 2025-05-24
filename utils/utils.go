package utils

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// FetchJSON делает GET-запрос и возвращает тело ответа как []byte
func FetchJSON(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// NormalizeSymbol возвращает токен без стейбла: BTC/USDT → BTC
func NormalizeSymbol(symbol string) string {
	symbol = strings.ToUpper(symbol)
	symbol = strings.ReplaceAll(symbol, "_", "")
	if strings.HasSuffix(symbol, "USDT") {
		return strings.TrimSuffix(symbol, "USDT") + "/USDT"
	}
	if strings.HasSuffix(symbol, "USDC") {
		return strings.TrimSuffix(symbol, "USDC") + "/USDC"
	}
	return symbol
}

func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
