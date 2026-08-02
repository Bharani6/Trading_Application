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

func SetupRouter(pRouter *gin.Engine) {
	// Controllers
	lAuthController := auth.NewAuthController()
	lProfileController := profile.NewProfileController()
	lWalletController := wallet.NewWalletController()
	lAdminController := admin.NewAdminController()
	lUtilsController := utils.NewUtilsController()
	lSupportController := support.NewSupportController()

	// Trade Dependencies
	lTradeRepo := trade.NewTradeRepository()
	lWalletRepo := wallet.NewWalletRepository()
	lMarketSvc := market_service.NewYahooFinanceService()
	lTradeService := trade.NewTradeService(lTradeRepo, lWalletRepo, lMarketSvc)
	lTradeController := trade.NewTradeController(lTradeService)
	lMarketController := market_controller.NewMarketController(lMarketSvc)

	lAPI := pRouter.Group("/api/v1")
	{
		// Authentication routes (Public) - rate limited to 10 req/60s per IP
		lAuth := lAPI.Group("/auth")
		lAuth.Use(middleware.RateLimiter(10, 60))
		{
			lAuth.POST("/register", lAuthController.Register)
			lAuth.POST("/login", lAuthController.Login)
			lAuth.POST("/forgot-password", lAuthController.ForgotPassword)
			lAuth.POST("/verify-reset-token", lAuthController.VerifyResetToken)
			lAuth.POST("/reset-password", lAuthController.ResetPassword)
		}

		lUtils := lAPI.Group("/utils")
		{
			lUtils.GET("/ifsc/:code", lUtilsController.FetchIFSC)
			lUtils.GET("/pincode/:code", lUtilsController.FetchPincode)
		}

		lMarket := lAPI.Group("/market")
		{
			lMarket.GET("/indices", lMarketController.GetIndices)
		}

		lAPI.POST("/support", lSupportController.SubmitMessage)

		// Protected routes
		lSecure := lAPI.Group("/")
		lSecure.Use(middleware.AuthMiddleware())
		{
			lSecure.GET("/users/me", lAuthController.GetMe)
			lSecure.POST("/users/kyc", lProfileController.SubmitKYC)
			lSecure.POST("/users/closure", lProfileController.RequestClosure)
			lSecure.POST("/users/change-password", lProfileController.ChangePassword)
			
			lSecure.GET("/auth/sessions", lAuthController.GetSessions)
			lSecure.DELETE("/auth/sessions/:id", lAuthController.RevokeSession)

			lWalletGroup := lSecure.Group("/wallet")
			{
				lWalletGroup.GET("/balance", lWalletController.GetBalance)
				lWalletGroup.POST("/add-fund", lWalletController.AddFunds)
				lWalletGroup.POST("/withdraw", lWalletController.WithdrawFunds)
				lWalletGroup.GET("/transactions", lWalletController.GetTransactions)
			}

			lShares := lSecure.Group("/shares")
			{
				lShares.GET("", lTradeController.GetShares)
			}

			lTrades := lSecure.Group("/trades")
			{
				lTrades.POST("/buy", lTradeController.BuyShare)
				lTrades.POST("/sell", lTradeController.SellShare)
				lTrades.GET("/history", lTradeController.GetUserTrades)
				lTrades.POST("/:id/cancel", lTradeController.CancelTrade)
			}

			// Admin only routes
			lAdmin := lSecure.Group("/admin")
			lAdmin.Use(middleware.RoleMiddleware("admin"))
			{
				lAdmin.GET("/users", lAdminController.GetUsers)
				lAdmin.GET("/users/:id/details", lAdminController.GetUserDetails)
				lAdmin.PUT("/users/:id/approve", lAdminController.ApproveUser)
				lAdmin.PUT("/users/:id/reject", lAdminController.RejectUser)
				lAdmin.PUT("/users/:id/block", lAdminController.BlockUser)
				lAdmin.PUT("/users/:id/close_account", lAdminController.CloseAccount)
				lAdmin.PUT("/users/:id/reject_closure", lAdminController.RejectClosure)
				lAdmin.POST("/shares/upload", lAdminController.UploadShares)
				lAdmin.DELETE("/shares", lAdminController.DeleteAllShares)
				lAdmin.GET("/support", lSupportController.GetMessages)
				lAdmin.PUT("/support/:id/status", lSupportController.UpdateStatus)
			}
			
			lWatchlistController := watchlist.NewWatchlistController()
			lSecure.POST("/watchlist", lWatchlistController.AddStock)
			lSecure.GET("/watchlist", lWatchlistController.GetWatchlist)
			lSecure.DELETE("/watchlist/:id", lWatchlistController.RemoveStock)
			lSecure.PUT("/watchlist/:id/favorite", lWatchlistController.UpdateFavorite)
		}
	}
}
