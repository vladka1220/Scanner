package exchangerequests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"basis_go/types"
	"basis_go/utils"
)

type GateSpot struct{}

func (g *GateSpot) Name() string {
	return "Gate.io Spot"
}

func (g *GateSpot) IsFutures() bool {
	return false
}

func GetSpotPrices() map[string]types.ExchangePrice {
	resp, err := http.Get("https://api.gateio.ws/api/v4/spot/tickers")
	if err != nil {
		fmt.Println("Ошибка запроса Gate.io Spot:", err)
		return nil
	}
	defer resp.Body.Close()

	var data []struct {
		CurrencyPair string `json:"currency_pair"`
		Last         string `json:"last"`
		BaseVolume   string `json:"base_volume"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Ошибка декодирования Gate.io Spot:", err)
		return nil
	}

	prices := make(map[string]types.ExchangePrice)
	for _, item := range data {
		if !strings.HasSuffix(item.CurrencyPair, "USDT") && !strings.HasSuffix(item.CurrencyPair, "USDC") {
			continue
		}
		symbol := utils.NormalizeSymbol(item.CurrencyPair)
		price, err1 := utils.ParseFloat(item.Last)
		volume, err2 := utils.ParseFloat(item.BaseVolume)
		if err1 != nil || err2 != nil || price <= 0 || volume <= 0 {
			continue
		}

		prices[symbol] = types.ExchangePrice{
			Price:     price,
			Volume:    volume,
			IsFutures: false,
		}
	}
	return prices
}

func (g *GateSpot) GetPrices() map[string]types.ExchangePrice {
	return GetSpotPrices()
}
