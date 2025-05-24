package pump_monitor

import "testing"

func TestAnalyzeTrades_NoTrades(t *testing.T) {
	res := AnalyzeTrades(nil, 0)
	if res.IsPump || res.Reason != "❌ Нет сделок" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestAnalyzeTrades_ValidPump(t *testing.T) {
	// override thresholds
	oldTrades := MinTrades
	oldBuyPerc := MinBuyPercentage
	oldBuyVol := MinBuyVolumeUSDT
	oldPriceChange := MinPriceChangePercent
	MinTrades = 2
	MinBuyPercentage = 50
	MinBuyVolumeUSDT = 1
	MinPriceChangePercent = 0.1
	defer func() {
		MinTrades = oldTrades
		MinBuyPercentage = oldBuyPerc
		MinBuyVolumeUSDT = oldBuyVol
		MinPriceChangePercent = oldPriceChange
	}()

	trades := []Trade{
		{Price: 1, Quantity: 1, IsBuyer: true, Symbol: "BTCUSDT"},
		{Price: 1, Quantity: 1, IsBuyer: true, Symbol: "BTCUSDT"},
	}
	res := AnalyzeTrades(trades, 1.0)
	if !res.IsPump {
		t.Fatalf("expected pump, got %+v", res)
	}
}

func TestGetSpotTradeLink(t *testing.T) {
	if l := GetSpotTradeLink("Binance", "BTCUSDT"); l != "https://www.binance.com/en/trade/BTC_USDT" {
		t.Errorf("binance link unexpected: %s", l)
	}
	if l := GetSpotTradeLink("MEXC", "BTCUSDT"); l != "https://www.mexc.com/exchange/BTC_USDT" {
		t.Errorf("mexc link unexpected: %s", l)
	}
	if l := GetSpotTradeLink("Gate", "BTCUSDT"); l != "https://www.gate.io/ru/trade/BTC_USDT" {
		t.Errorf("gate link unexpected: %s", l)
	}
	if l := GetSpotTradeLink("Bybit", "BTCUSDT"); l != "https://www.bybit.com/trade/spot/BTC/USDT" {
		t.Errorf("bybit link unexpected: %s", l)
	}
}
