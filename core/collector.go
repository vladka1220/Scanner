package core

import (
	binance "basis_go/exchangerequests/binance"
	bybit "basis_go/exchangerequests/bybit"
	gate "basis_go/exchangerequests/gate"
	mexc "basis_go/exchangerequests/mexc"
	"basis_go/types"
	"sync"
)

func CollectSpotPrices() types.TokenPrices {
	sources := map[string]func() map[string]types.ExchangePrice{
		"MEXC":    mexc.GetSpotPrices,
		"Binance": binance.GetSpotPrices,
		"Gate":    (&gate.GateSpot{}).GetPrices,
		"Bybit":   (&bybit.BybitSpot{}).GetPrices,
	}
	return collect(sources)
}

func CollectFuturesPrices() types.TokenPrices {
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
	return collect(sources)
}

func collect(sources map[string]func() map[string]types.ExchangePrice) types.TokenPrices {
	all := make(types.TokenPrices)
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(sources))
	for name, fetch := range sources {
		go func(exchange string, fn func() map[string]types.ExchangePrice) {
			defer wg.Done()
			data := fn()
			if len(data) == 0 {
				return
			}
			mu.Lock()
			for token, price := range data {
				if _, ok := all[token]; !ok {
					all[token] = make(map[string]types.ExchangePrice)
				}
				all[token][exchange] = price
			}
			mu.Unlock()
		}(name, fetch)
	}
	wg.Wait()
	return all
}
