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

func main() {
	err := godotenv.Load("/Users/vladyslav/Documents/Проекты /basis_go/.env")
	if err != nil {
		fmt.Println("❌❌❌ Проблема оповещений❌❌❌", err)
	} else {
		fmt.Println("✅✅✅ Телеграм подключен✅✅✅")
	}
	mode := flag.String("mode", "all", "Режим запуска: spot, futures, spotfutures, pump, all")
	flag.Parse()

	for {
		start := time.Now()

		switch *mode {
		case "spot":
			spotPrices := spot.CollectSpotPrices()
			spot.CompareSpotPrices(spotPrices)

		case "futures":
			futuresPrices := futures.CollectFuturesPrices()
			futures.ComparePrices(futuresPrices)

		case "spotfutures":
			spotPrices := spot.CollectSpotPrices()
			futuresPrices := futures.CollectFuturesPrices()
			comparison_price.CompareSpotFutures(spotPrices, futuresPrices)

		case "pump":
			pump_monitor.MonitorPumps()
			return

		case "all":
			spotPrices := spot.CollectSpotPrices()
			spot.CompareSpotPrices(spotPrices)

			futuresPrices := futures.CollectFuturesPrices()
			futures.ComparePrices(futuresPrices)

			comparison_price.CompareSpotFutures(spotPrices, futuresPrices)

		default:
			fmt.Println("❌ Неизвестный режим. Используйте: -mode=spot, -mode=futures, -mode=spotfutures, -mode=pump или -mode=all")
			return
		}

		fmt.Printf("⏱ Цикл занял: %.2fs\n", time.Since(start).Seconds())
		time.Sleep(15 * time.Second)
	}
}
