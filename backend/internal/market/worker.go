package market

import (
	"log"
	"time"

	"stock-trading/internal/config"
	"stock-trading/internal/market/service"
	"stock-trading/internal/trade"

	"gorm.io/gorm"
)

var WorkerStatus = "Stopped"

func IsIndianMarketOpen() bool {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*3600+1800)
	}
	now := time.Now().In(loc)

	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	hour := now.Hour()
	min := now.Minute()

	// Market opens at 09:15
	if hour < 9 || (hour == 9 && min < 15) {
		return false
	}

	// Market closes at 15:30
	if hour > 15 || (hour == 15 && min > 30) {
		return false
	}

	return true
}

func StartMarketDataWorker(db *gorm.DB, marketSvc service.MarketDataService) {
	intervalStr := config.App.App.MarketUpdateInterval
	if intervalStr == "" {
		intervalStr = "1m" // default to 1 minute
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Printf("Invalid MARKET_UPDATE_INTERVAL %s, defaulting to 1m", intervalStr)
		interval = 1 * time.Minute
	}

	ticker := time.NewTicker(interval)
	WorkerStatus = "Running"
	go func() {
		log.Printf("Market Data Worker started, interval: %s\n", interval.String())
		for range ticker.C {
			if !IsIndianMarketOpen() {
				// Market is closed, do nothing
				continue
			}

			if err := updateStockPrices(db, marketSvc); err != nil {
				log.Printf("Failed to update stock prices: %v\n", err)
			}
		}
	}()
}

func updateStockPrices(db *gorm.DB, marketSvc service.MarketDataService) error {
	var shares []trade.Share
	if err := db.Find(&shares).Error; err != nil {
		return err
	}

	if len(shares) == 0 {
		return nil
	}

	var symbols []string
	for _, s := range shares {
		symbols = append(symbols, s.Symbol)
	}

	prices, err := marketSvc.GetLatestPrices(symbols)
	if err != nil {
		return err
	}

	if len(prices) == 0 {
		return nil
	}

	// Begin transaction to update all prices
	return db.Transaction(func(tx *gorm.DB) error {
		for _, s := range shares {
			if newPrice, ok := prices[s.Symbol]; ok && newPrice.Current > 0 {
				if err := tx.Model(&trade.Share{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
					"price":          newPrice.Current,
					"previous_price": newPrice.Previous,
					"updated_at":     time.Now(),
				}).Error; err != nil {
					return err
				}
			}
		}
		log.Printf("Successfully updated %d stock prices from market data provider", len(prices))
		return nil
	})
}
