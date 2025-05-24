package exchangerequests

import (
	"basis_go/types"
	"basis_go/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GateFutures struct{}

func (g *GateFutures) Name() string {
	return "Gate Futures"
}

func (g *GateFutures) IsFutures() bool {
	return true
}

func (g *GateFutures) GetPrices() map[string]types.ExchangePrice {
	return GetGateFuturesPrices()
}

func GetGateFuturesPrices() map[string]types.ExchangePrice {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// --------- 1. Тикеры -----------
	req, _ := http.NewRequest("GET", "https://fx-api.gateio.ws/api/v4/futures/usdt/tickers", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Bot/1.0)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Ошибка запроса Gate Futures тикеров:", err)
		return nil
	}
	defer resp.Body.Close()

	var tickerData []struct {
		Symbol      string `json:"contract"`
		Last        string `json:"last"`
		QuoteVolume string `json:"volume_24h_quote"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tickerData); err != nil {
		fmt.Println("❌ Ошибка декодирования тикеров Gate Futures:", err)
		return nil
	}

	// --------- 2. Funding ----------
	freq, _ := http.NewRequest("GET", "https://fx-api.gateio.ws/api/v4/futures/usdt/funding_rates", nil)
	freq.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Bot/1.0)")
	fResp, err := client.Do(freq)
	if err != nil {
		fmt.Println("❌ Ошибка запроса funding_rates Gate:", err)
		return nil
	}
	defer fResp.Body.Close()

	var fundingResponse struct {
		Data []struct {
			Contract    string `json:"contract"`
			FundingRate string `json:"funding_rate"`
			FundingTime string `json:"funding_time"`
		} `json:"data"`
	}

	if err := json.NewDecoder(fResp.Body).Decode(&fundingResponse); err != nil {
		fmt.Println("❌ Ошибка декодирования funding_rates Gate:", err)
		return nil
	}

	fundingMap := make(map[string]struct {
		Rate float64
		Time time.Time
	})

	for _, item := range fundingResponse.Data {
		rate, _ := utils.ParseFloat(item.FundingRate)
		ts, _ := utils.ParseInt64(item.FundingTime)
		fundingMap[item.Contract] = struct {
			Rate float64
			Time time.Time
		}{
			Rate: rate,
			Time: time.UnixMilli(ts),
		}
	}

	// --------- 4. Объединённый результат -----------
	result := make(map[string]types.ExchangePrice)
	for _, item := range tickerData {
		if !strings.HasSuffix(item.Symbol, "_USDT") && !strings.HasSuffix(item.Symbol, "_USDC") {
			continue
		}

		symbol := utils.NormalizeSymbol(item.Symbol)
		price, err1 := utils.ParseFloat(item.Last)
		volume, err2 := utils.ParseFloat(item.QuoteVolume)
		if err1 != nil || err2 != nil || price <= 0 {
			continue
		}

		f := fundingMap[item.Symbol]

		result[symbol] = types.ExchangePrice{
			Price:            price,
			Volume:           volume,
			IsFutures:        true,
			FundingRate:      f.Rate,
			FundingIntervalH: 8,
			NextFundingTime:  f.Time,
		}
	}

	return result
}
