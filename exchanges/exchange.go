package exchanges

import (
	binance "basis_go/exchangerequests/binance"
	bybit "basis_go/exchangerequests/bybit"
	gate "basis_go/exchangerequests/gate"
	mexc "basis_go/exchangerequests/mexc"
	"basis_go/types"
)

type Exchange interface {
	Name() string
	IsFutures() bool
	GetPrices() map[string]types.ExchangePrice
}

var AllExchanges = []Exchange{
	&binance.BinanceSpot{},
	&binance.BinanceFutures{},
	&mexc.MEXCSpot{},
	&mexc.MEXCFutures{},
	&gate.GateSpot{},
	&gate.GateFutures{},
	&bybit.BybitSpot{},
	&bybit.BybitFutures{},
}
