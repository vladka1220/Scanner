package exchangerequests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"basis_go/types"
	"basis_go/utils"
)

type BinanceFutures struct{}

func (b *BinanceFutures) Name() string {
	return "Binance Futures"
}

func (b *BinanceFutures) IsFutures() bool {
	return true
}

func GetFuturesPrices() map[string]types.ExchangePrice {
	tickerURL := "https://fapi.binance.com/fapi/v1/ticker/bookTicker"
	fundingURL := "https://fapi.binance.com/fapi/v1/premiumIndex"

	tickerResp, err := http.Get(tickerURL)
	if err != nil {
		fmt.Println("Binance futures ticker error:", err)
		return nil
	}
	defer tickerResp.Body.Close()

	fundingResp, err := http.Get(fundingURL)
	if err != nil {
		fmt.Println("Binance futures funding error:", err)
		return nil
	}
	defer fundingResp.Body.Close()

	var tickers []struct {
		Symbol string `json:"symbol"`
		Ask    string `json:"askPrice"`
		Bid    string `json:"bidPrice"`
	}
	var fundingData []struct {
		Symbol          string `json:"symbol"`
		LastFundingRate string `json:"lastFundingRate"`
		NextFundingTime int64  `json:"nextFundingTime"`
	}

	if err := json.NewDecoder(tickerResp.Body).Decode(&tickers); err != nil {
		fmt.Println("Binance futures decode ticker error:", err)
		return nil
	}
	if err := json.NewDecoder(fundingResp.Body).Decode(&fundingData); err != nil {
		fmt.Println("Binance futures decode funding error:", err)
		return nil
	}

	// Map funding by symbol
	fundingMap := make(map[string]struct {
		Rate float64
		Time time.Time
	})
	for _, item := range fundingData {
		rate, err := utils.ParseFloat(item.LastFundingRate)
		if err != nil {
			continue
		}
		fundingMap[item.Symbol] = struct {
			Rate float64
			Time time.Time
		}{
			Rate: rate,
			Time: time.UnixMilli(item.NextFundingTime),
		}
	}

	prices := make(map[string]types.ExchangePrice)
	for _, item := range tickers {
		if !strings.HasSuffix(item.Symbol, "USDT") && !strings.HasSuffix(item.Symbol, "USDC") {
			continue
		}
		token := utils.NormalizeSymbol(item.Symbol)
		ask, err1 := utils.ParseFloat(item.Ask)
		bid, err2 := utils.ParseFloat(item.Bid)
		if err1 != nil || err2 != nil || ask <= 0 || bid <= 0 {
			continue
		}

		funding := fundingMap[item.Symbol]
		prices[token] = types.ExchangePrice{
			Price:            ask,
			Volume:           bid,
			IsFutures:        true,
			FundingRate:      funding.Rate,
			FundingIntervalH: 8, // Binance по умолчанию — 8ч, можно сделать динамическим
			NextFundingTime:  funding.Time,
		}
	}
	return prices
}
func (b *BinanceFutures) GetPrices() map[string]types.ExchangePrice {
	return GetFuturesPrices()
}
