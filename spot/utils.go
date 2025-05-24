package spot

import (
	"basis_go/types"
	"basis_go/utils"
	"fmt"
	"sort"
)

func typeName(p types.ExchangePrice) string {
	if p.IsFutures {
		return "Фьючерс"
	}
	return "Спот"
}

func CompareSpotPrices(prices types.TokenPrices) {
	fmt.Println("=== АРБИТРАЖ SPOT ===")
	type Opportunity struct {
		token    string
		ex1, ex2 string
		p1, p2   float64
		t1, t2   string
		spread   float64
	}

	var results []Opportunity

	for token, data := range prices {
		for ex1, p1 := range data {
			for ex2, p2 := range data {
				if ex1 == ex2 || p1.Price <= 0 || p2.Price <= 0 {
					continue
				}
				//fmt.Printf("✅ [Спот] Токен %s найден на %d биржах\n", token, len(data))
				spread := utils.CalculateSpotSpread(p1.Price, p2.Price)
				if spread >= 0.25 && spread <= 100 {
					results = append(results, Opportunity{
						token:  token,
						ex1:    ex1,
						ex2:    ex2,
						p1:     p1.Price,
						p2:     p2.Price,
						t1:     typeName(p1),
						t2:     typeName(p2),
						spread: spread,
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
		net := utils.NetSpread(r.spread, utils.DefaultMakerFeePercent, utils.DefaultTakerFeePercent, 0)
		fmt.Printf("[%s]\n- %s (%s) → %s (%s) | %.6f → %.6f | Спред: %.2f%% (Чистый: %.2f%%)\n",
			r.token, r.ex1, r.t1, r.ex2, r.t2, r.p1, r.p2, r.spread, net)
	}
	fmt.Println("====================================")
}
