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
	lLoc, lErr := time.LoadLocation("Asia/Kolkata")
	if lErr != nil {
		lLoc = time.FixedZone("IST", 5*3600+1800)
	}
	lNow := time.Now().In(lLoc)

	if lNow.Weekday() == time.Saturday || lNow.Weekday() == time.Sunday {
		return false
	}

	lHour := lNow.Hour()
	lMin := lNow.Minute()

	// Market opens at 09:15
	if lHour < 9 || (lHour == 9 && lMin < 15) {
		return false
	}

	// Market closes at 15:30
	if lHour > 15 || (lHour == 15 && lMin > 30) {
		return false
	}

	return true
}

func StartMarketDataWorker(pDB *gorm.DB, pMarketSvc service.MarketDataService) {
	lIntervalStr := config.App.App.MarketUpdateInterval
	if lIntervalStr == "" {
		lIntervalStr = "1m" // default to 1 minute
	}

	lInterval, lErr := time.ParseDuration(lIntervalStr)
	if lErr != nil {
		log.Printf("Invalid MARKET_UPDATE_INTERVAL %s, defaulting to 1m", lIntervalStr)
		lInterval = 1 * time.Minute
	}

	lTicker := time.NewTicker(lInterval)
	WorkerStatus = "Running"
	go func() {
		log.Printf("Market Data Worker started, interval: %s\n", lInterval.String())
		
		// Fetch prices immediately on startup so they aren't 0.0
		if lErr := updateStockPrices(pDB, pMarketSvc); lErr != nil {
			log.Printf("Failed to initial update stock prices: %v\n", lErr)
		}

		for range lTicker.C {
			if !IsIndianMarketOpen() {
				// Market is closed, do nothing
				continue
			}

			if lErr := updateStockPrices(pDB, pMarketSvc); lErr != nil {
				log.Printf("Failed to update stock prices: %v\n", lErr)
			}
		}
	}()
}

func updateStockPrices(pDB *gorm.DB, pMarketSvc service.MarketDataService) error {
	var lShares []trade.Share
	if lErr := pDB.Find(&lShares).Error; lErr != nil {
		return lErr
	}

	if len(lShares) == 0 {
		return nil
	}

	var lSymbols []string
	for _, lS := range lShares {
		lSymbols = append(lSymbols, lS.Symbol)
	}

	lPrices, lErr := pMarketSvc.GetLatestPrices(lSymbols)
	if lErr != nil {
		return lErr
	}

	if len(lPrices) == 0 {
		return nil
	}

	// Begin transaction to update all prices
	return pDB.Transaction(func(pTx *gorm.DB) error {
		for _, lS := range lShares {
			if lNewPrice, lOk := lPrices[lS.Symbol]; lOk && lNewPrice.Current > 0 {
				if lErr := pTx.Model(&trade.Share{}).Where("id = ?", lS.ID).Updates(map[string]interface{}{
					"price":          lNewPrice.Current,
					"previous_price": lNewPrice.Previous,
					"updated_at":     time.Now(),
				}).Error; lErr != nil {
					return lErr
				}
			}
		}
		log.Printf("Successfully updated %d stock prices from market data provider", len(lPrices))
		return nil
	})
}
