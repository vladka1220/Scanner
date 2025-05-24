package main

import (
	"basis_go/comparison_price"
	"basis_go/futures"
	"basis_go/pump_monitor"
	"basis_go/spot"
	"flag"
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

// function variables allow tests to stub expensive operations
var (
	collectSpotPrices    = spot.CollectSpotPrices
	compareSpotPrices    = spot.CompareSpotPrices
	collectFuturesPrices = futures.CollectFuturesPrices
	compareFuturesPrices = futures.ComparePrices
	compareSpotFutures   = comparison_price.CompareSpotFutures
	monitorPumps         = pump_monitor.MonitorPumps
)

// runOnce controls whether runMode executes only a single iteration. It is
// intended for tests.
var runOnce bool

// runMode executes the logic for the provided mode. If once is true, only a
// single iteration is executed.
func runMode(mode string, once bool) {
	for {
		start := time.Now()

		switch mode {
		case "spot":
			sp := collectSpotPrices()
			compareSpotPrices(sp)

		case "futures":
			fp := collectFuturesPrices()
			compareFuturesPrices(fp)

		case "spotfutures":
			sp := collectSpotPrices()
			fp := collectFuturesPrices()
			compareSpotFutures(sp, fp)

		case "pump":
			monitorPumps()
			return

		case "all":
			sp := collectSpotPrices()
			compareSpotPrices(sp)

			fp := collectFuturesPrices()
			compareFuturesPrices(fp)

			compareSpotFutures(sp, fp)

		default:
			fmt.Println("❌ Неизвестный режим. Используйте: -mode=spot, -mode=futures, -mode=spotfutures, -mode=pump или -mode=all")
			return
		}

		fmt.Printf("⏱ Цикл занял: %.2fs\n", time.Since(start).Seconds())
		if once {
			return
		}
		time.Sleep(15 * time.Second)
	}
}

func main() {
	err := godotenv.Load("/Users/vladyslav/Documents/Проекты /basis_go/.env")
	if err != nil {
		fmt.Println("❌❌❌ Проблема оповещений❌❌❌", err)
	} else {
		fmt.Println("✅✅✅ Телеграм подключен✅✅✅")
	}
	mode := flag.String("mode", "all", "Режим запуска: spot, futures, spotfutures, pump, all")
	flag.Parse()

	runMode(*mode, runOnce)
}
