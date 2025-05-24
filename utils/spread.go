package utils

// Default estimated fees in percent (e.g. 0.1 means 0.1%).
const (
	DefaultMakerFeePercent = 0.02
	DefaultTakerFeePercent = 0.07
)

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
