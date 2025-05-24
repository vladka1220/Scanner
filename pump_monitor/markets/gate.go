package markets

import (
	"basis_go/types"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	gateTickersURL      = "https://api.gateio.ws/api/v4/spot/tickers"
	gateRecentTradesURL = "https://api.gate.io/api/v4/spot/trades?currency_pair=%s"
)

type GateSpotTicker struct {
	Symbol      string `json:"currency_pair"`
	Last        string `json:"last"`
	BaseVolume  string `json:"base_volume"`
	QuoteVolume string `json:"quote_volume"`
}

type GateTrade struct {
	Price        string `json:"price"`
	Amount       string `json:"amount"`
	CreateTimeMS int64  `json:"create_time_ms"`
	Side         string `json:"side"`
}

func FetchGateTickers() (map[string]types.PriceInfo, error) {
	url := gateTickersURL

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // необходимо, если ошибка TLS
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes := make([]byte, 4096)
	n, _ := resp.Body.Read(bodyBytes)
	if strings.Contains(string(bodyBytes[:n]), "<html>") {
		return nil, fmt.Errorf("некорректный JSON: %s", string(bodyBytes[:n]))
	}

	resp.Body.Close()
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickers []GateSpotTicker
	err = json.NewDecoder(resp.Body).Decode(&tickers)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	result := make(map[string]types.PriceInfo)
	for _, t := range tickers {
		var price, volume, quoteVol float64
		fmt.Sscanf(t.Last, "%f", &price)
		fmt.Sscanf(t.BaseVolume, "%f", &volume)
		fmt.Sscanf(t.QuoteVolume, "%f", &quoteVol)

		if price <= 0 || quoteVol <= 0 {
			continue
		}

		symbol := formatSymbol(t.Symbol)
		result[symbol] = types.PriceInfo{
			Price:       price,
			Volume:      volume,
			QuoteVolume: quoteVol,
		}
	}
	return result, nil
}

func FetchRecentGateTrades(symbol string) ([]GateTrade, error) {
	url := fmt.Sprintf(gateRecentTradesURL, symbol)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trades []GateTrade
	err = json.NewDecoder(resp.Body).Decode(&trades)
	if err != nil {
		return nil, err
	}

	return trades, nil
}

func formatSymbol(s string) string {
	return strings.ToUpper(strings.Replace(s, "_", "", -1))
}
