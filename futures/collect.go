package futures

import (
	"fmt"
	"sync"

	binance "basis_go/exchangerequests/binance"
	bybit "basis_go/exchangerequests/bybit"
	gate "basis_go/exchangerequests/gate"
	mexc "basis_go/exchangerequests/mexc"
	"basis_go/types"
)

func CollectFuturesPrices() types.TokenPrices {
	all := make(types.TokenPrices)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	sources := map[string]func() map[string]types.ExchangePrice{
		"Binance Futures": func() map[string]types.ExchangePrice {
			return (&binance.BinanceFutures{}).GetPrices()
		},
		"MEXC Futures": func() map[string]types.ExchangePrice {
			return (&mexc.MEXCFutures{}).GetPrices()
		},
		"Gate Futures": func() map[string]types.ExchangePrice {
			return (&gate.GateFutures{}).GetPrices()
		},
		"Bybit Futures": func() map[string]types.ExchangePrice {
			return (&bybit.BybitFutures{}).GetPrices()
		},
	}

	wg.Add(len(sources))
	for name, fn := range sources {
		go func(exchange string, fetch func() map[string]types.ExchangePrice) {
			defer wg.Done()
			data := fetch()
			if len(data) == 0 {
				fmt.Printf("❌ [%s] ФЬЮЧЕРС НЕ ОТДАЁТ ДАННЫХ\n", exchange)
				return
			}
			fmt.Printf("✅ [%s] загружено токенов: %d\n", exchange, len(data))
			mu.Lock()
			for token, price := range data {
				if _, ok := all[token]; !ok {
					all[token] = map[string]types.ExchangePrice{}
				}
				all[token][exchange] = price
			}
			mu.Unlock()
		}(name, fn)
	}
	wg.Wait()
	return all
}
