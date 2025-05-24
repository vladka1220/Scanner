package pump_monitor

import (
	"fmt"
)

type Trade struct {
	Price     float64
	Quantity  float64
	IsBuyer   bool // true если это покупка
	Timestamp int64
	Symbol    string
}

type PumpAnalysis struct {
	IsPump      bool
	Reason      string
	BuyPercent  float64
	BuyVolume   float64
	AvgBuySize  float64
	PriceChange float64
}

func AnalyzeTrades(trades []Trade, priceChange float64) PumpAnalysis {
	var buyVolume, totalVolume float64
	var totalBuyCount int
	for _, t := range trades {
		vol := t.Price * t.Quantity
		totalVolume += vol
		if t.IsBuyer {
			buyVolume += vol
			totalBuyCount++
		}
	}

	totalTrades := len(trades)
	if totalTrades == 0 {
		return PumpAnalysis{IsPump: false, Reason: "❌ Нет сделок"}
	}

	buyPercent := (float64(totalBuyCount) / float64(totalTrades)) * 100
	avgBuySize := 0.0
	if totalBuyCount > 0 {
		avgBuySize = buyVolume / float64(totalBuyCount)
	}

	// поэтапная проверка
	if totalTrades < MinTrades {
		return PumpAnalysis{IsPump: false, Reason: fmt.Sprintf("❌ Мало сделок: %d < %d", totalTrades, MinTrades)}
	}
	if buyPercent < MinBuyPercentage {
		return PumpAnalysis{IsPump: false, Reason: fmt.Sprintf("❌ Мало покупок: %.1f%% < %.1f%%", buyPercent, MinBuyPercentage)}
	}
	if buyVolume < MinBuyVolumeUSDT {
		return PumpAnalysis{IsPump: false, Reason: fmt.Sprintf("❌ Недостаточный объем покупок: $%.0f < $%.0f", buyVolume, MinBuyVolumeUSDT)}
	}
	if priceChange < MinPriceChangePercent {
		return PumpAnalysis{IsPump: false, Reason: fmt.Sprintf("❌ Слабый рост цены в сделках: %.2f%% < %.2f%%", priceChange, MinPriceChangePercent)}
	}

	// Если все прошло
	reason := fmt.Sprintf(
		"💥 Памп на %s: цена +%.2f%%, %.1f%% покупок, $%.0f объём, %d сделок, средняя покупка $%.2f",
		trades[0].Symbol, priceChange, buyPercent, buyVolume, totalTrades, avgBuySize,
	)

	return PumpAnalysis{
		IsPump:      true,
		Reason:      reason,
		BuyPercent:  buyPercent,
		BuyVolume:   buyVolume,
		AvgBuySize:  avgBuySize,
		PriceChange: priceChange,
	}
}
