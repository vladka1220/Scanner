package comparison_price

import (
	"basis_go/funding"
	"basis_go/types"
	"fmt"
	"sort"
)

func typeName(p types.ExchangePrice) string {
	if p.IsFutures {
		return "Фьючерс"
	}
	return "Спот"
}

func formatFunding(p types.ExchangePrice) string {
	if !p.IsFutures {
		return ""
	}
	return funding.FormatFunding(
		p.FundingRate,
		p.FundingTime,
		p.Price,
		1.0, // без плеча, просто отображаем чистую ставку
		int64(p.FundingIntervalH),
	)
}

func CompareSpotFutures(spotPrices, futuresPrices types.TokenPrices) {
	fmt.Println("✅✅✅✅✅=== АРБИТРАЖ SPOT → FUTURES ===✅✅✅✅✅✅")
	type Opportunity struct {
		token    string
		ex1, ex2 string
		p1, p2   float64
		t1, t2   string
		spread   float64
		f1, f2   types.ExchangePrice
	}

	var results []Opportunity

	for token, spotData := range spotPrices {
		futuresData, ok := futuresPrices[token]
		if !ok {
			continue
		}
		for ex1, p1 := range spotData {
			for ex2, p2 := range futuresData {
				if p1.Price <= 0 || p2.Price <= 0 {
					continue
				}
				spread := ((p2.Price - p1.Price) / p1.Price) * 100
				if spread >= 1 && spread <= 100 {
					results = append(results, Opportunity{
						token:  token,
						ex1:    ex1,
						ex2:    ex2,
						p1:     p1.Price,
						p2:     p2.Price,
						t1:     typeName(p1),
						t2:     typeName(p2),
						spread: spread,
						f1:     p1,
						f2:     p2,
					})
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].spread > results[j].spread
	})

	limit := 50
	if len(results) < limit {
		limit = len(results)
	}

	for _, r := range results[:limit] {
		fmt.Printf("[%s]\n- %s (%s, %s) → %s (%s, %s) | %.6f → %.6f | Спред: %.2f%%\n",
			r.token,
			r.ex1, r.t1, formatFunding(r.f1),
			r.ex2, r.t2, formatFunding(r.f2),
			r.p1, r.p2, r.spread)
	}
	fmt.Println("====================================")
}
