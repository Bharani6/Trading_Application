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

func (pController *WalletController) GetBalance(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	if !lExists {
		response.Error(lCtx, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context", nil)
		return
	}

	lRes, lErr := pController.svc.GetBalance(lUserID.(string))
	if lErr != nil {
		response.Error(lCtx, http.StatusInternalServerError, "WALLET_ERROR", lErr.Error(), nil)
		return
	}

	response.Success(lCtx, http.StatusOK, "Wallet balance retrieved", lRes)
}

func (pController *WalletController) AddFunds(lCtx *gin.Context) {
	lUserID, _ := lCtx.Get("userID")

	var lReq FundRequest
	if lErr := lCtx.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(lCtx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input", lErr.Error())
		return
	}

	lErr := pController.svc.AddFunds(lUserID.(string), lReq.Amount)
	if lErr != nil {
		response.Error(lCtx, http.StatusBadRequest, "WALLET_ERROR", lErr.Error(), nil)
		return
	}

	response.Success(lCtx, http.StatusOK, "Funds added successfully", nil)
}

func (pController *WalletController) WithdrawFunds(lCtx *gin.Context) {
	lUserID, _ := lCtx.Get("userID")

	var lReq FundRequest
	if lErr := lCtx.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(lCtx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input", lErr.Error())
		return
	}

	lErr := pController.svc.WithdrawFunds(lUserID.(string), lReq.Amount)
	if lErr != nil {
		response.Error(lCtx, http.StatusBadRequest, "WALLET_ERROR", lErr.Error(), nil)
		return
	}

	response.Success(lCtx, http.StatusOK, "Funds withdrawn successfully", nil)
}

func (pController *WalletController) GetTransactions(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	if !lExists {
		response.Error(lCtx, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context", nil)
		return
	}

	lTransactions, lErr := pController.svc.GetTransactions(lUserID.(string))
	if lErr != nil {
		response.Error(lCtx, http.StatusInternalServerError, "TRANSACTION_ERROR", lErr.Error(), nil)
		return
	}

	response.Success(lCtx, http.StatusOK, "Transactions retrieved", lTransactions)
}
