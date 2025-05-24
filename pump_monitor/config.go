package pump_monitor

// тестовый быстрый просчет
/*
var (
	MinPriceGrowth     float64 = 3  // минимальный рост цены в %
	MinVolumeGrowth    float64 = 5  // минимальный рост объема в %
	IntervalSec        int     = 2  // интервал между запросами в секундах
	CompareIntervalSec int     = 30 // интервал, за который сравнивается рост в секндах

	// Настройки анализа сделок
	MinTrades             int     = 10   // минимальное количество сделок
	TradeFetchLimit       int     = 30   // кол-во последних сделок для анализа
	MinBuyPercentage      float64 = 35.0 // минимальный процент покупок
	MinBuyVolumeUSDT      float64 = 500  // минимальный объём покупок в USDT
	MinPriceChangePercent float64 = 3    // минимальное изменение цены за минуту в %
	MinQuoteVolumeUSDT    float64 = 500  // Минимальный суточный объём USDT

)
*/

// нужно подобрать настройки
var (
	MinPriceGrowth     float64 = 0.5 // Цена должна вырасти минимум на 2.5%
	MinVolumeGrowth    float64 = 0.5 // Объем должен вырасти минимум на 3%
	IntervalSec        int     = 5   // Проверяем каждые 5 секунд
	CompareIntervalSec int     = 60  // Сравниваем с данными минутной давности

	// Настройки анализа сделок (глубокий анализ)
	MinTrades             int     = 40   // Нужно минимум 40 сделок за 60 сек
	TradeFetchLimit       int     = 100  // Загружаем последние 100 сделок
	MinBuyPercentage      float64 = 40.0 // Минимум 75% покупок
	MinBuyVolumeUSDT      float64 = 500  // Объем покупок за минуту от $5000
	MinPriceChangePercent float64 = 0.5  // Рост цены внутри сделок должен быть ≥ 1.2%
	MinQuoteVolumeUSDT    float64 = 500  // Минимальный суточный объём USDT
)

/*
// High Confidence Pump
var (
	MinPriceGrowth     float64 = 5.0 // Цена должна вырасти минимум на 5%
	MinVolumeGrowth    float64 = 8.0 // Объем должен вырасти минимум на 10%
	IntervalSec        int     = 5   // Проверяем каждые 5 секунд
	CompareIntervalSec int     = 60  // Сравниваем с минутной давностью

	// Глубокий анализ
	MinTrades             int     = 30   // Нужно минимум 80 сделок
	TradeFetchLimit       int     = 80   // Загружаем 120 последних сделок
	MinBuyPercentage      float64 = 70.0 // Минимум 85% покупок
	MinBuyVolumeUSDT      float64 = 500  // Объем покупок от $10,000
	MinPriceChangePercent float64 = 2.5  // Рост цены внутри сделок ≥ 2.5%
	MinQuoteVolumeUSDT    float64 = 5000 // Минимальный суточный объём USDT
)
*/

/*
// Ultra Clean Pump Only
var (
	MinPriceGrowth     float64 = 8.0  // Цена должна вырасти минимум на 8%
	MinVolumeGrowth    float64 = 20.0 // Объем должен вырасти минимум на 20%
	IntervalSec        int     = 5
	CompareIntervalSec int     = 60

	// Глубокий анализ
	MinTrades             int     = 100   // Минимум 100 сделок за минуту
	TradeFetchLimit       int     = 150   // Последние 150 сделок
	MinBuyPercentage      float64 = 90.0  // 90% сделок — покупки
	MinBuyVolumeUSDT      float64 = 15000 // От $15,000
	MinPriceChangePercent float64 = 3.5   // Рост ≥ 3.5%
	MinQuoteVolumeUSDT    float64 = 5000 // Минимальный суточный объём USDT
)
*/
