package exchangerequests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"basis_go/types"
	"basis_go/utils"
)

type BybitSpot struct{}

func (b *BybitSpot) Name() string {
	return "Bybit Spot"
}

func (b *BybitSpot) IsFutures() bool {
	return false
}

func (b *BybitSpot) GetPrices() map[string]types.ExchangePrice {
	return GetSpotPrices()
}

func GetSpotPrices() map[string]types.ExchangePrice {
	resp, err := http.Get("https://api.bybit.com/v5/market/tickers?category=spot")
	if err != nil {
		fmt.Println("Ошибка запроса Bybit Spot:", err)
		return nil
	}
	defer resp.Body.Close()

	var parsed struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			List []struct {
				Symbol   string `json:"symbol"`
				BidPrice string `json:"bid1Price"`
				AskPrice string `json:"ask1Price"`
			} `json:"list"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		fmt.Println("Ошибка декодирования Bybit Spot:", err)
		return nil
	}

	prices := make(map[string]types.ExchangePrice)
	for _, item := range parsed.Result.List {
		if !strings.HasSuffix(item.Symbol, "USDT") && !strings.HasSuffix(item.Symbol, "USDC") {
			continue
		}

		symbol := utils.NormalizeSymbol(item.Symbol)
		ask, err1 := utils.ParseFloat(item.AskPrice)
		bid, err2 := utils.ParseFloat(item.BidPrice)
		if err1 != nil || err2 != nil || ask <= 0 || bid <= 0 {
			continue
		}

		prices[symbol] = types.ExchangePrice{
			Price:     ask,
			Volume:    bid,
			IsFutures: false,
		}
	}

	return prices
}
