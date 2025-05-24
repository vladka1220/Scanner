package pump_monitor

import (
	"basis_go/notifier"
	"basis_go/pump_monitor/markets"
	"basis_go/types"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PricePoint struct {
	Price     float64
	Volume    float64
	Timestamp time.Time
}

var debugMode = true // ← отключить когда тесты прошли
var lastShown = make(map[string]types.PriceInfo)
var history = make(map[string][]PricePoint)

var telegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")
var telegramChatID = os.Getenv("TELEGRAM_CHAT_ID")

func sendTelegramMessage(message string) {
	if telegramToken == "" || telegramChatID == "" {
		return
	}
	sendURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)
	resp, err := http.PostForm(sendURL, url.Values{
		"chat_id": {telegramChatID},
		"text":    {message},
	})
	if err != nil {
		fmt.Println("❌ Ошибка отправки в Telegram:", err)
		return
	}
	defer resp.Body.Close()
}

func MonitorPumps() {
	for {
		allData := make(map[string]types.PriceInfo)
		now := time.Now()

		binanceData, err := markets.FetchBinanceTickers()
		if err == nil {
			fmt.Printf("🔄 Получено %d тикеров с Binance\n", len(binanceData))
			for k, v := range binanceData {
				allData["Binance|"+k] = v
			}
		} else {
			fmt.Println("❌ Ошибка Binance:", err)
		}

		mexcData, err := markets.FetchMEXCTickers()
		if err == nil {
			fmt.Printf("🔄 Получено %d тикеров с MEXC\n", len(mexcData))
			for k, v := range mexcData {
				allData["MEXC|"+k] = v
			}
		} else {
			fmt.Println("❌ Ошибка MEXC:", err)
		}

		gateData, err := markets.FetchGateTickers()
		if err == nil {
			fmt.Printf("🔄 Получено %d тикеров с Gate.io\n", len(gateData))
			for k, v := range gateData {
				allData["Gate|"+k] = v
			}
		} else {
			fmt.Println("❌ Ошибка Gate.io:", err)
		}

		bybitData, err := markets.FetchBybitTickers()
		if err == nil {
			fmt.Printf("🔄 Получено %d тикеров с Bybit\n", len(bybitData))
			for k, v := range bybitData {
				allData["Bybit|"+k] = v
			}
		} else {
			fmt.Println("❌ Ошибка Bybit:", err)
		}

		var found []string

		for key, current := range allData {
			parts := strings.SplitN(key, "|", 2)
			if len(parts) != 2 {
				continue
			}

			exchange, symbol := parts[0], parts[1]
			if current.QuoteVolume < MinQuoteVolumeUSDT {
				continue
			}

			if !strings.HasSuffix(symbol, "USDT") && !strings.HasSuffix(symbol, "USDC") {
				continue
			}

			history[key] = append(history[key], PricePoint{
				Price:     current.Price,
				Volume:    current.Volume,
				Timestamp: now,
			})

			cutoff := now.Add(-time.Duration(CompareIntervalSec) * time.Second)
			var recent []PricePoint
			for _, pt := range history[key] {
				if pt.Timestamp.After(cutoff) {
					recent = append(recent, pt)
				}
			}
			history[key] = recent

			if len(recent) < 2 {
				continue
			}
			old := recent[0]

			if old.Price > 0 && old.Volume > 0 {
				priceChange := ((current.Price - old.Price) / old.Price) * 100
				volumeChange := ((current.Volume - old.Volume) / old.Volume) * 100

				last, wasShown := lastShown[key]
				recentPriceChange := 0.0
				if wasShown && last.Price > 0 {
					recentPriceChange = ((current.Price - last.Price) / last.Price) * 100
				}

				if priceChange >= MinPriceGrowth && volumeChange >= MinVolumeGrowth {
					if !wasShown || recentPriceChange > 0.5 {
						tradesRaw, err := markets.FetchRecentTrades(symbol)
						if err == nil {
							converted := make([]Trade, 0, len(tradesRaw))
							for _, tr := range tradesRaw {
								price, _ := strconv.ParseFloat(tr.Price, 64)
								qty, _ := strconv.ParseFloat(tr.Qty, 64)
								isBuy := !tr.IsBuyerMaker
								converted = append(converted, Trade{Price: price, Quantity: qty, IsBuyer: isBuy, Timestamp: tr.Time})
							}
							result := AnalyzeTrades(converted, priceChange)

							if result.IsPump {
								printDebugPass(exchange, symbol, priceChange, volumeChange, len(converted), result, current.QuoteVolume)

								link := GetSpotTradeLink(exchange, symbol)
								msg := fmt.Sprintf("🚨 %s [%s][%s]", result.Reason, exchange, symbol)
								linkLine := fmt.Sprintf("🟢 Торговать: %s", link)
								combined := fmt.Sprintf("%s\n%s", msg, linkLine)
								found = append(found, combined)
								sendTelegramMessage(combined)
								lastShown[key] = current
								continue
							}
						}
						// ❌ Никаких fallback-логов, если фильтр не прошёл
					}
				}
			}
		}

		sort.Strings(found)
		for _, msg := range found {
			fmt.Println(msg)
			err := notifier.SendTelegramMessage(msg)
			if err != nil {
				fmt.Println("BOT_TOKEN:", os.Getenv("TELEGRAM_BOT_TOKEN"))
				fmt.Println("CHAT_ID:", os.Getenv("TELEGRAM_CHAT_ID"))
				fmt.Println("❌ Ошибка отправки в Telegram:", err)
			}
		}

		time.Sleep(time.Duration(IntervalSec) * time.Second)
	}
}

func printDebugPass(
	exchange, symbol string,
	priceChange, volumeChange float64,
	tradeCount int,
	result PumpAnalysis,
	quoteVolume float64,
) {
	if !debugMode {
		return
	}
	fmt.Printf(
		"🚨🚨🚨🚨🚨🚨 Памп [%s][%s]\n"+
			"✅ Прошёл фильтры:\n"+
			"   • Цена: %.2f%% ≥ %.2f%%\n"+
			"   • Объём: %.2f%% ≥ %.2f%%\n"+
			"   • Сделки: %d ≥ %d\n"+
			"   • %% покупок: %.1f%% ≥ %.1f%%\n"+
			"   • Объём покупок: $%.0f ≥ $%.0f\n"+
			"   • Рост в сделках: %.2f%% ≥ %.2f%%\n"+
			"   • Объём за 24ч: $%.0f ≥ $%.0f\n\n",
		exchange, symbol,
		priceChange, MinPriceGrowth,
		volumeChange, MinVolumeGrowth,
		tradeCount, MinTrades,
		result.BuyPercent, MinBuyPercentage,
		result.BuyVolume, MinBuyVolumeUSDT,
		result.PriceChange, MinPriceChangePercent,
		quoteVolume, MinQuoteVolumeUSDT,
	)
}

func GetSpotTradeLink(exchange, symbol string) string {
	switch exchange {
	case "MEXC":
		return fmt.Sprintf("https://www.mexc.com/exchange/%s", strings.Replace(symbol, "USDT", "_USDT", 1))
	case "Binance":
		return fmt.Sprintf("https://www.binance.com/en/trade/%s_USDT", strings.ToUpper(strings.Replace(symbol, "USDT", "", 1)))
	case "Gate":
		return fmt.Sprintf("https://www.gate.io/ru/trade/%s", strings.Replace(symbol, "USDT", "_USDT", 1))
	case "Bybit":
		return fmt.Sprintf("https://www.bybit.com/trade/spot/%s", strings.Replace(symbol, "USDT", "/USDT", 1))
	default:
		return ""
	}
}
