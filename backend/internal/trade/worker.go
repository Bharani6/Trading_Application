package trade

import (
	"log"
	"time"
)

func StartTradeWorker(tradeSvc TradeService) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		log.Println("Trade worker started, checking pending trades every minute.")
		for range ticker.C {
			if IsMarketOpen() {
				err := tradeSvc.ExecutePendingTrades()
				if err != nil {
					log.Printf("Error executing pending trades: %v\n", err)
				}
			}
		}
	}()
}
