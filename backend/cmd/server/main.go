package main

import (
	"fmt"
	"log"
	"net/http"
	userpkg "stock-trading/internal/user"

	"stock-trading/internal/config"
	"stock-trading/internal/database"
	"stock-trading/internal/logger"
	"stock-trading/internal/market"
	"stock-trading/internal/market/service"
	"stock-trading/internal/routes"
	"stock-trading/internal/trade"
	"stock-trading/internal/profile"
	"stock-trading/internal/support"
	"stock-trading/internal/wallet"
	watchlistpkg "stock-trading/internal/watchlist"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Load config
	config.LoadConfig("configs")

	// 2. Init logger
	logger.InitLogger(config.App.App.Env)
	defer logger.Log.Sync()

	zap.L().Info("Starting Stock Trading API Server...")

	// 3. Init Database & Redis
	database.ConnectDB()
	if err := database.DB.AutoMigrate(
		&userpkg.User{},
		&userpkg.PasswordResetToken{},
		&watchlistpkg.Watchlist{},
		&userpkg.Session{},
		&userpkg.BankDetails{},
		&userpkg.NomineeDetails{},
		&userpkg.PersonalDetails{},
		&wallet.Wallet{},
		&wallet.Transaction{},
		&trade.Segment{},
		&trade.Share{},
		&trade.Trade{},
		&profile.FileUpload{},
		&support.SupportMessage{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
	// database.InitRedis()

	// 4. Setup Services & Background Workers
	tradeSvc := trade.NewTradeService()
	trade.StartTradeWorker(tradeSvc)

	marketSvc := service.NewYahooFinanceService()
	market.StartMarketDataWorker(database.DB, marketSvc)

	// 5. Setup Router
	if config.App.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		dbStatus := "OK"
		if err := database.DB.Exec("SELECT 1").Error; err != nil {
			dbStatus = "Error: " + err.Error()
		}

		yahooStatus := "OK"
		// Make a quick HEAD/GET request to Yahoo to check connectivity
		if resp, err := http.Get("https://query1.finance.yahoo.com/v7/finance/quote?symbols=AAPL"); err != nil || resp.StatusCode != 200 {
			yahooStatus = "Error or Unreachable"
		}

		workerStatus := market.WorkerStatus
		version := "v1.0.0"

		c.JSON(200, gin.H{
			"database": dbStatus,
			"yahoo":    yahooStatus,
			"worker":   workerStatus,
			"version":  version,
		})
	})

	// Register API Routes
	routes.SetupRouter(r)

	// 5. Start Server
	port := fmt.Sprintf(":%d", config.App.App.Port)
	zap.L().Info("Server running on port", zap.String("port", port))
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

