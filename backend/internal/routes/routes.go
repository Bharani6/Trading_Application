package routes

import (
	"stock-trading/internal/admin"
	"stock-trading/internal/auth"
	"stock-trading/internal/middleware"
	"stock-trading/internal/profile"
	"stock-trading/internal/trade"
	"stock-trading/internal/utils"
	"stock-trading/internal/wallet"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Controllers
	authController := auth.NewAuthController()
	profileController := profile.NewProfileController()
	walletController := wallet.NewWalletController()
	tradeController := trade.NewTradeController()
	adminController := admin.NewAdminController()
	utilsController := utils.NewUtilsController()

	api := r.Group("/api/v1")
	{
		// Authentication routes (Public) - rate limited to 10 req/60s per IP
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimiter(10, 60))
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		utils := api.Group("/utils")
		{
			utils.GET("/ifsc/:code", utilsController.FetchIFSC)
			utils.GET("/pincode/:code", utilsController.FetchPincode)
		}

		// Protected routes
		secure := api.Group("/")
		secure.Use(middleware.AuthMiddleware())
		{
			secure.GET("/users/me", authController.GetMe)
			secure.POST("/users/kyc", profileController.SubmitKYC)

			wallet := secure.Group("/wallet")
			{
				wallet.GET("/balance", walletController.GetBalance)
				wallet.POST("/add-fund", walletController.AddFunds)
				wallet.POST("/withdraw", walletController.WithdrawFunds)
				wallet.GET("/transactions", walletController.GetTransactions)
			}

			shares := secure.Group("/shares")
			{
				shares.GET("", tradeController.GetShares)
			}

			trades := secure.Group("/trades")
			{
				trades.POST("/buy", tradeController.BuyShare)
				trades.POST("/sell", tradeController.SellShare)
				trades.GET("/history", tradeController.GetUserTrades)
				trades.POST("/:id/cancel", tradeController.CancelTrade)
			}

			// Admin only routes
			admin := secure.Group("/admin")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/users", adminController.GetUsers)
				admin.PUT("/users/:id/approve", adminController.ApproveUser)
				admin.PUT("/users/:id/reject", adminController.RejectUser)
				admin.PUT("/users/:id/block", adminController.BlockUser)
				admin.POST("/shares/upload", adminController.UploadShares)
				admin.DELETE("/shares", adminController.DeleteAllShares)
			}
		}
	}
}
