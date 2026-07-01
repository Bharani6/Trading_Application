package main

import (
	"log"
	"stock-trading/internal/config"
	"stock-trading/internal/database"
	"stock-trading/internal/trade"
)

func main() {
	config.LoadConfig("configs")
	database.ConnectDB()

	res := database.DB.Unscoped().Where("1 = 1").Delete(&trade.Share{})
	if res.Error != nil {
		log.Fatalf("Failed to delete shares: %v", res.Error)
	}

	log.Printf("Successfully deleted %d shares", res.RowsAffected)
}
