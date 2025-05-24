package funding

import (
	"fmt"
	"time"
)

// FormatFunding возвращает фандинг в процентах и время до следующего начисления
func FormatFunding(funding float64, nextTime int64, price float64, leverage float64, hoursPerFunding int64) string {
	// преобразуем в проценты и округляем
	percent := funding * 100

	now := time.Now().UnixMilli()

	// если время устарело — пересчитай на следующее
	intervalMs := hoursPerFunding * 60 * 60 * 1000
	for nextTime < now {
		nextTime += intervalMs
	}

	remaining := nextTime - now
	dur := time.Duration(remaining) * time.Millisecond
	hours := int(dur.Hours())
	minutes := int(dur.Minutes()) % 60

	// расчёт дохода/убытка по фандингу
	fundingAmount := price * funding * leverage

	return fmt.Sprintf("%.3f%% (через %dh %dm | %.6f USDT)", percent, hours, minutes, fundingAmount)
}

func FormatNextFundingTime(timestampMs int64) string {
	remaining := timestampMs - time.Now().UnixMilli()
	if remaining < 0 {
		remaining = 0
	}
	duration := time.Duration(remaining) * time.Millisecond
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", h, m)
}
