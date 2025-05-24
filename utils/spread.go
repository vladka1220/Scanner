package utils

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
