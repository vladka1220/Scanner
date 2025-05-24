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

type MEXCFutures struct{}

func (m *MEXCFutures) Name() string {
	return "MEXC Futures"
}

func (m *MEXCFutures) IsFutures() bool {
	return true
}

func (m *MEXCFutures) GetPrices() map[string]types.ExchangePrice {
	tickerURL := "https://contract.mexc.com/api/v1/contract/ticker"
	fundingURL := "https://contract.mexc.com/api/v1/contract/funding_rate"

	// ---------- TICKER ----------
	var tickerRespParsed struct {
		Data []struct {
			Symbol string  `json:"symbol"`
			Ask    float64 `json:"ask1"`
			Bid    float64 `json:"bid1"`
		} `json:"data"`
	}
	tickerResp, err := http.Get(tickerURL)
	if err != nil {
		fmt.Println("Ошибка запроса MEXC Futures ticker:", err)
		return nil
	}
	defer tickerResp.Body.Close()
	if err := json.NewDecoder(tickerResp.Body).Decode(&tickerRespParsed); err != nil {
		fmt.Println("Ошибка декодирования тикеров MEXC:", err)
		return nil
	}

	// ---------- FUNDING ----------
	var fundingRespParsed struct {
		Data []struct {
			Symbol          string  `json:"symbol"`
			FundingRate     float64 `json:"fundingRate"`
			NextFundingTime int64   `json:"nextFundingTime"`
		} `json:"data"`
	}
	fundingResp, err := http.Get(fundingURL)
	if err != nil {
		fmt.Println("Ошибка запроса MEXC Futures funding:", err)
		return nil
	}
	defer fundingResp.Body.Close()
	if err := json.NewDecoder(fundingResp.Body).Decode(&fundingRespParsed); err != nil {
		fmt.Println("Ошибка декодирования фандинга MEXC:", err)
		return nil
	}

	// карта фандинга по символам
	fundingMap := make(map[string]struct {
		Rate float64
		Time time.Time
	})
	for _, item := range fundingRespParsed.Data {
		rate := item.FundingRate
		fundingMap[item.Symbol] = struct {
			Rate float64
			Time time.Time
		}{
			Rate: rate,
			Time: time.UnixMilli(item.NextFundingTime),
		}
	}

	prices := make(map[string]types.ExchangePrice)
	for _, item := range tickerRespParsed.Data {
		if !strings.HasSuffix(item.Symbol, "USDT") {
			continue
		}
		symbol := utils.NormalizeSymbol(item.Symbol)
		ask := item.Ask
		bid := item.Bid
		if ask <= 0 || bid <= 0 {
			continue
		}
		f := fundingMap[item.Symbol]
		prices[symbol] = types.ExchangePrice{
			Price:            ask,
			Volume:           bid,
			IsFutures:        true,
			FundingRate:      f.Rate,
			FundingIntervalH: 8,
			NextFundingTime:  f.Time,
		}
	}
	return prices
}
