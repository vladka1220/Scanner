package core

import (
	"fmt"
	"sort"
	"time"

	"basis_go/types"
	"basis_go/utils"
)

func typeName(p types.ExchangePrice) string {
	if p.IsFutures {
		return "Фьючерс"
	}
	return "Спот"
}

func fundingInfo(p types.ExchangePrice) string {
	if !p.IsFutures {
		return ""
	}
	return fmt.Sprintf(" | Фандинг: %.6f%%, раз в %dч, через: %s",
		p.FundingRate*100,
		p.FundingIntervalH,
		time.Until(p.NextFundingTime).Truncate(time.Second),
	)
}

func CompareAll(prices types.TokenPrices) {
	fmt.Println("=== АРБИТРАЖНЫЕ СВЯЗКИ (спред > 2% и < 100%) ===")
	tokens := make([]string, 0, len(prices))
	for token := range prices {
		tokens = append(tokens, token)
	}
	sort.Strings(tokens)

	for _, token := range tokens {
		data := prices[token]
		var printed bool

		for ex1, p1 := range data {
			for ex2, p2 := range data {
				if ex1 == ex2 || p1.Price <= 0 || p2.Price <= 0 {
					continue
				}
				spread := ((p2.Price - p1.Price) / p1.Price) * 100
				if spread >= 2 && spread <= 100 {
					if !printed {
						fmt.Printf("\n[%s]\n", token)
						printed = true
					}

					var fundingDiff float64
					if p1.IsFutures || p2.IsFutures {
						fundingDiff = (p2.FundingRate - p1.FundingRate) * 100
					}
					net := utils.NetSpread(spread, utils.DefaultMakerFeePercent, utils.DefaultTakerFeePercent, fundingDiff)
					fmt.Printf("- %s (%s) → %s (%s) | %.6f → %.6f | Спред: %.2f%% (Чистый: %.2f%%)%s%s\n",
						ex1, typeName(p1),
						ex2, typeName(p2),
						p1.Price, p2.Price, spread, net,
						fundingInfo(p1), fundingInfo(p2),
					)
				}
			}
		}
	}
	fmt.Println("==========================================")
}
