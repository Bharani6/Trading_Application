package wallet

import (
	"net/http"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	svc WalletService
}

func NewWalletController() *WalletController {
	return &WalletController{svc: NewWalletService()}
}

func (c *WalletController) GetBalance(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context", nil)
		return
	}

	res, err := c.svc.GetBalance(userID.(string))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "WALLET_ERROR", err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Wallet balance retrieved", res)
}

func (c *WalletController) AddFunds(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req FundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input", err.Error())
		return
	}

	err := c.svc.AddFunds(userID.(string), req.Amount)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "WALLET_ERROR", err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Funds added successfully", nil)
}

func (c *WalletController) WithdrawFunds(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req FundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input", err.Error())
		return
	}

	err := c.svc.WithdrawFunds(userID.(string), req.Amount)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "WALLET_ERROR", err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Funds withdrawn successfully", nil)
}

func (c *WalletController) GetTransactions(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context", nil)
		return
	}

	transactions, err := c.svc.GetTransactions(userID.(string))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "TRANSACTION_ERROR", err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Transactions retrieved", transactions)
}
