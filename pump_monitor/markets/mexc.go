package markets

import (
	"basis_go/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MEXCSpotTicker struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
}

type Trade struct {
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	Time         int64  `json:"time"`
}

func FetchMEXCTickers() (map[string]types.PriceInfo, error) {
	url := "https://api.mexc.com/api/v3/ticker/24hr"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickers []MEXCSpotTicker
	err = json.NewDecoder(resp.Body).Decode(&tickers)
	if err != nil {
		return nil, err
	}

	result := make(map[string]types.PriceInfo)
	for _, t := range tickers {
		var price, volume float64
		_, err1 := fmt.Sscanf(t.LastPrice, "%f", &price)
		_, err2 := fmt.Sscanf(t.Volume, "%f", &volume)

		if err1 != nil || err2 != nil {
			continue
		}

		result[t.Symbol] = types.PriceInfo{Price: price, Volume: volume}
	}

	return result, nil
}

func FetchRecentTrades(symbol string) ([]Trade, error) {
	url := fmt.Sprintf("https://api.mexc.com/api/v3/trades?symbol=%s", symbol)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trades []Trade
	err = json.NewDecoder(resp.Body).Decode(&trades)
	if err != nil {
		return nil, err
	}

	return trades, nil
}
