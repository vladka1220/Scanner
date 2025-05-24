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

type BybitFutures struct{}

func (b *BybitFutures) Name() string {
	return "Bybit Futures"
}

func (b *BybitFutures) IsFutures() bool {
	return true
}

func (b *BybitFutures) GetPrices() map[string]types.ExchangePrice {
	tickerURL := "https://api.bybit.com/v5/market/tickers?category=linear"
	fundingURL := "https://api.bybit.com/v5/market/funding/history?category=linear"

	// ----------- TICKER -----------
	var tickerResp struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string `json:"category"`
			List     []struct {
				Symbol   string `json:"symbol"`
				AskPrice string `json:"ask1Price"`
				BidPrice string `json:"bid1Price"`
			} `json:"list"`
		} `json:"result"`
	}

	tResp, err := http.Get(tickerURL)
	if err != nil {
		fmt.Println("Ошибка запроса Bybit Futures ticker:", err)
		return nil
	}
	defer tResp.Body.Close()

	if err := json.NewDecoder(tResp.Body).Decode(&tickerResp); err != nil {
		fmt.Println("Ошибка декодирования тикеров Bybit:", err)
		return nil
	}

	// ----------- FUNDING -----------
	var fundingResp struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string `json:"category"`
			List     []struct {
				Symbol               string `json:"symbol"`
				FundingRate          string `json:"fundingRate"`
				FundingRateTimestamp int64  `json:"fundingRateTimestamp"`
			} `json:"list"`
		} `json:"result"`
	}

	fResp, err := http.Get(fundingURL)
	if err != nil {
		fmt.Println("Ошибка запроса Bybit Futures funding:", err)
		return nil
	}
	defer fResp.Body.Close()

	if err := json.NewDecoder(fResp.Body).Decode(&fundingResp); err != nil {
		fmt.Println("Ошибка декодирования фандинга Bybit:", err)
		return nil
	}

	fundingMap := make(map[string]struct {
		Rate float64
		Time time.Time
	})
	for _, item := range fundingResp.Result.List {
		rate, _ := utils.ParseFloat(item.FundingRate)
		timeUnix := time.UnixMilli(item.FundingRateTimestamp)
		fundingMap[item.Symbol] = struct {
			Rate float64
			Time time.Time
		}{Rate: rate, Time: timeUnix}
	}

	prices := make(map[string]types.ExchangePrice)
	for _, item := range tickerResp.Result.List {
		if !strings.HasSuffix(item.Symbol, "USDT") && !strings.HasSuffix(item.Symbol, "USDC") {
			continue
		}

		ask, err1 := utils.ParseFloat(item.AskPrice)
		bid, err2 := utils.ParseFloat(item.BidPrice)
		if err1 != nil || err2 != nil || ask <= 0 || bid <= 0 {
			continue
		}

		symbol := utils.NormalizeSymbol(item.Symbol)
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
