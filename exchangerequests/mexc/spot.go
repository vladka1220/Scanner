package exchangerequests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"basis_go/types"
	"basis_go/utils"
)

type MEXCSpot struct{}

func (m *MEXCSpot) Name() string {
	return "MEXC Spot"
}

func (m *MEXCSpot) IsFutures() bool {
	return false
}

func GetSpotPrices() map[string]types.ExchangePrice {
	resp, err := http.Get("https://api.mexc.com/api/v3/ticker/bookTicker")
	if err != nil {
		fmt.Println("Ошибка запроса MEXC Spot:", err)
		return nil
	}
	defer resp.Body.Close()

	var data []struct {
		Symbol string `json:"symbol"`
		Bid    string `json:"bidPrice"`
		Ask    string `json:"askPrice"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Ошибка декодирования MEXC Spot:", err)
		return nil
	}

	prices := make(map[string]types.ExchangePrice)
	for _, item := range data {
		if !strings.HasSuffix(item.Symbol, "USDT") && !strings.HasSuffix(item.Symbol, "USDC") {
			continue
		}
		symbol := utils.NormalizeSymbol(item.Symbol)
		ask, err1 := utils.ParseFloat(item.Ask)
		bid, err2 := utils.ParseFloat(item.Bid)
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

func (m *MEXCSpot) GetPrices() map[string]types.ExchangePrice {
	return GetSpotPrices()
}
