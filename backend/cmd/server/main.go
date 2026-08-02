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
	
	// Automatically seed the database on startup
	if err := SeedData(database.DB); err != nil {
		zap.L().Warn("Failed to seed initial data", zap.Error(err))
	}
	// database.InitRedis()

	// 4. Setup Services & Background Workers
	lTradeRepo := trade.NewTradeRepository()
	lWalletRepo := wallet.NewWalletRepository()
	lMarketService := service.NewYahooFinanceService()
	
	lTradeService := trade.NewTradeService(lTradeRepo, lWalletRepo, lMarketService)
	trade.StartTradeWorker(lTradeService)

	market.StartMarketDataWorker(database.DB, lMarketService)

	// 5. Setup Router
	if config.App.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	lRouter := gin.Default()

	// CORS middleware
	lRouter.Use(func(pContext *gin.Context) {
		pContext.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		pContext.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		pContext.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		pContext.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if pContext.Request.Method == "OPTIONS" {
			pContext.AbortWithStatus(204)
			return
		}

		pContext.Next()
	})

	// Health check
	lRouter.GET("/health", func(pContext *gin.Context) {
		lDBStatus := "OK"
		if lErr := database.DB.Exec("SELECT 1").Error; lErr != nil {
			lDBStatus = "Error: " + lErr.Error()
		}

		lYahooStatus := "OK"
		// Make a quick HEAD/GET request to Yahoo to check connectivity
		if lResp, lErr := http.Get("https://query1.finance.yahoo.com/v7/finance/quote?symbols=AAPL"); lErr != nil || lResp.StatusCode != 200 {
			lYahooStatus = "Error or Unreachable"
		}

		lWorkerStatus := market.WorkerStatus
		lVersion := "v1.0.0"

		pContext.JSON(200, gin.H{
			"database": lDBStatus,
			"yahoo":    lYahooStatus,
			"worker":   lWorkerStatus,
			"version":  lVersion,
		})
	})

	// Register API Routes
	routes.SetupRouter(lRouter)

	// 5. Start Server
	lPort := fmt.Sprintf(":%d", config.App.App.Port)
	zap.L().Info("Server running on port", zap.String("port", lPort))
	if lErr := lRouter.Run(lPort); lErr != nil {
		log.Fatalf("Server failed to start: %v", lErr)
	}
}

