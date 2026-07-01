package main

import (
	"fmt"
	"log"
	userpkg "stock-trading/internal/user"

	"stock-trading/internal/config"
	"stock-trading/internal/database"
	"stock-trading/internal/logger"
	"stock-trading/internal/response"
	"stock-trading/internal/routes"
	"stock-trading/internal/trade"

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
	database.DB.AutoMigrate(&userpkg.User{})
	// database.InitRedis()

	// 4. Setup Services & Background Workers
	tradeSvc := trade.NewTradeService()
	trade.StartTradeWorker(tradeSvc)

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
	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, 200, "Pong", nil)
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

