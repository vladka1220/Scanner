package types

import "time"

type ExchangePrice struct {
	Price            float64   // Цена
	Volume           float64   // Объём
	IsFutures        bool      // Является ли фьючерсом
	FundingRate      float64   // Ставка фандинга (если есть)
	FundingIntervalH int       // Интервал начисления в часах (например, 8)
	NextFundingTime  time.Time // Время следующего фандинга
	FundingTime      int64     `json:"fundingTime"` // Время следующего начисления фандинга в миллисекундах
}

type TokenPrices map[string]map[string]ExchangePrice

type FundingInfo struct {
	Rate         float64
	IntervalMins int
	NextFunding  int64 // Unix timestamp
}

type OrderBook struct {
	Ask float64
	Bid float64
}

type PriceInfo struct {
	Price       float64
	Volume      float64
	QuoteVolume float64
}
