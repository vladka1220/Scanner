package markets

import (
	"basis_go/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	bybitTickersURL      = "https://api.bybit.com/v5/market/tickers?category=spot"
	bybitRecentTradesURL = "https://api.bybit.com/v5/market/recent-trade?category=spot&symbol=%s"
)

type BybitSpotTicker struct {
	Symbol      string `json:"symbol"`
	LastPrice   string `json:"lastPrice"`
	QuoteVolume string `json:"turnover24h"`
	Volume      string `json:"volume24h"`
}

type BybitTrade struct {
	Price     string `json:"price"`
	Qty       string `json:"qty"`
	Side      string `json:"side"`
	Timestamp int64  `json:"time"`
}

func FetchBybitTickers() (map[string]types.PriceInfo, error) {
	url := bybitTickersURL
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса Bybit Spot: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string            `json:"category"`
			List     []BybitSpotTicker `json:"list"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	prices := make(map[string]types.PriceInfo)
	for _, t := range result.Result.List {
		if !strings.HasSuffix(t.Symbol, "USDT") && !strings.HasSuffix(t.Symbol, "USDC") {
			continue
		}
		var price, volume, quoteVol float64
		fmt.Sscanf(t.LastPrice, "%f", &price)
		fmt.Sscanf(t.Volume, "%f", &volume)
		fmt.Sscanf(t.QuoteVolume, "%f", &quoteVol)

		if price <= 0 || quoteVol <= 0 {
			continue
		}

		symbol := strings.ToUpper(t.Symbol)
		prices[symbol] = types.PriceInfo{
			Price:       price,
			Volume:      volume,
			QuoteVolume: quoteVol,
		}
	}
	return prices, nil
}

func FetchRecentBybitTrades(symbol string) ([]BybitTrade, error) {
	url := fmt.Sprintf(bybitRecentTradesURL, symbol)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса последних трейдов: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string       `json:"category"`
			List     []BybitTrade `json:"list"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования трейдов Bybit: %w", err)
	}

	return result.Result.List, nil
}
