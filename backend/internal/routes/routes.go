package routes

import (
	"stock-trading/internal/admin"
	"stock-trading/internal/auth"
	"stock-trading/internal/middleware"
	"stock-trading/internal/profile"
	"stock-trading/internal/support"
	"stock-trading/internal/trade"
	"stock-trading/internal/utils"
	"stock-trading/internal/wallet"
	"stock-trading/internal/watchlist"
	market_controller "stock-trading/internal/market/controller"
	market_service "stock-trading/internal/market/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Controllers
	authController := auth.NewAuthController()
	profileController := profile.NewProfileController()
	walletController := wallet.NewWalletController()
	adminController := admin.NewAdminController()
	utilsController := utils.NewUtilsController()
	supportController := support.NewSupportController()

	// Trade Dependencies
	tradeRepo := trade.NewTradeRepository()
	walletRepo := wallet.NewWalletRepository()
	marketSvc := market_service.NewYahooFinanceService()
	tradeService := trade.NewTradeService(tradeRepo, walletRepo, marketSvc)
	tradeController := trade.NewTradeController(tradeService)
	marketController := market_controller.NewMarketController(marketSvc)

	api := r.Group("/api/v1")
	{
		// Authentication routes (Public) - rate limited to 10 req/60s per IP
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimiter(10, 60))
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/forgot-password", authController.ForgotPassword)
			auth.POST("/verify-reset-token", authController.VerifyResetToken)
			auth.POST("/reset-password", authController.ResetPassword)
			auth.POST("/send-otp", authController.SendOTP)
			auth.POST("/verify-otp", authController.VerifyOTP)
			auth.POST("/send-email-otp", authController.SendEmailOTP)
			auth.POST("/verify-email-otp", authController.VerifyEmailOTP)
		}

		utils := api.Group("/utils")
		{
			utils.GET("/ifsc/:code", utilsController.FetchIFSC)
			utils.GET("/pincode/:code", utilsController.FetchPincode)
		}

		market := api.Group("/market")
		{
			market.GET("/indices", marketController.GetIndices)
		}

		api.POST("/support", supportController.SubmitMessage)

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
				admin.GET("/support", supportController.GetMessages)
			}
			
			watchlistController := watchlist.NewWatchlistController()
			secure.POST("/watchlist", watchlistController.AddStock)
			secure.GET("/watchlist", watchlistController.GetWatchlist)
			secure.DELETE("/watchlist/:id", watchlistController.RemoveStock)
			secure.PUT("/watchlist/:id/favorite", watchlistController.UpdateFavorite)
		}
	}
}
