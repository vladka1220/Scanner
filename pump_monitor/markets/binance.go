package markets

import (
	"basis_go/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BinanceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"lastPrice"`
	Volume string `json:"quoteVolume"`
}

func FetchBinanceTickers() (map[string]types.PriceInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.binance.com/api/v3/ticker/24hr")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tickers []BinanceTicker
	err = json.NewDecoder(resp.Body).Decode(&tickers)
	if err != nil {
		return nil, err
	}

	result := make(map[string]types.PriceInfo)
	for _, t := range tickers {
		/*if strings.HasSuffix(t.Symbol, "USDT") || strings.HasSuffix(t.Symbol, "USDC") {     //Проверка какиен пары мы получаем с бинанс для проверки торгов на USDT & USDC
			fmt.Println("✔ Binance символ:", t.Symbol)
		}*/
		var price, volume float64
		fmt.Sscanf(t.Price, "%f", &price)
		fmt.Sscanf(t.Volume, "%f", &volume)
		result[t.Symbol] = types.PriceInfo{Price: price, Volume: volume, QuoteVolume: volume}
	}
	return result, nil
}
