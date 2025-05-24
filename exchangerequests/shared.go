package exchanges

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type ExchangePrice struct {
	Exchange  string
	Price     float64
	IsFutures bool
}

type TokenPrices map[string]map[string]ExchangePrice

var client = &http.Client{}

func fetchJSON(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
