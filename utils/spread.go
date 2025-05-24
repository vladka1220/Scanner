package utils

// Default estimated fees in percent (e.g. 0.1 means 0.1%).
const (
	DefaultMakerFeePercent = 0.02
	DefaultTakerFeePercent = 0.07
)

// FeeInfo holds maker/taker fees for a particular exchange.
type FeeInfo struct {
	Maker float64
	Taker float64
}

// ExchangeFees lists maker and taker fees by exchange name. The names should
// match what the Exchange implementations return from the Name() method.
var ExchangeFees = map[string]FeeInfo{
	"Binance Spot":    {Maker: 0.02, Taker: 0.07},
	"Binance Futures": {Maker: 0.02, Taker: 0.07},
	"MEXC Spot":       {Maker: 0.02, Taker: 0.07},
	"MEXC Futures":    {Maker: 0.02, Taker: 0.07},
	"Gate.io Spot":    {Maker: 0.02, Taker: 0.07},
	"Gate Futures":    {Maker: 0.02, Taker: 0.07},
	"Bybit Spot":      {Maker: 0.02, Taker: 0.07},
	"Bybit Futures":   {Maker: 0.02, Taker: 0.07},
}

// GetFees returns the fee info for the given exchange. If the exchange is not
// found in the table, the default fees are returned as a fallback.
func GetFees(exchange string) FeeInfo {
	if f, ok := ExchangeFees[exchange]; ok {
		return f
	}
	return FeeInfo{Maker: DefaultMakerFeePercent, Taker: DefaultTakerFeePercent}
}

// CalculateSpotSpread — расчёт спреда для спотового рынка
func CalculateSpotSpread(p1, p2 float64) float64 {
	if p1 <= 0 || p2 <= 0 {
		return 0
	}
	return ((p2 - p1) / p1) * 100
}

// CalculateFuturesSpread — расчёт спреда для фьючерсного рынка
func CalculateFuturesSpread(p1, p2 float64) float64 {
	if p1 <= 0 || p2 <= 0 {
		return 0
	}
	return ((p2 - p1) / p1) * 100
}

// NetSpread subtracts maker/taker fees and funding costs from the spread.
// makerFeePercent and takerFeePercent should be provided as percentages
// (e.g. 0.1 for 0.1%). fundingDiffPercent represents the expected funding
// payment difference in percent.
func NetSpread(spreadPercent, makerFeePercent, takerFeePercent, fundingDiffPercent float64) float64 {
	return spreadPercent - makerFeePercent - takerFeePercent - fundingDiffPercent
}
