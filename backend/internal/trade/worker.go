package trade

import (
	"log"
	"time"
)

func StartTradeWorker(pTradeSvc TradeService) {
	lTicker := time.NewTicker(1 * time.Minute)
	go func() {
		log.Println("Trade worker started, checking pending trades every minute.")
		for range lTicker.C {
			if IsMarketOpen() {
				lErr := pTradeSvc.ExecutePendingTrades()
				if lErr != nil {
					log.Printf("Error executing pending trades: %v\n", lErr)
				}
			}
		}
	}()
}
