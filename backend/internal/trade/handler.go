package trade

import (
	"fmt"
	"log"
	"net/http"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type TradeController struct {
	svc TradeService
}

func NewTradeController() *TradeController {
	return &TradeController{svc: NewTradeService()}
}

// ========================== GET SHARES ==========================

func ConstructGetShares(pSvc TradeService) (interface{}, error) {
	log.Println("ConstructGetShares (+)")
	lShares, lErr := pSvc.GetAllShares()
	log.Println("ConstructGetShares (-)")
	return lShares, lErr
}

func CommunicateGetShares(pCtx *gin.Context, pShares interface{}) error {
	log.Println("CommunicateGetShares (+)")
	response.Success(pCtx, http.StatusOK, "Shares retrieved", pShares)
	log.Println("CommunicateGetShares (-)")
	return nil
}

func CompleteGetShares(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	log.Println("CompleteGetShares (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	log.Println("CompleteGetShares (-)")
}

func (c *TradeController) GetShares(lCtx *gin.Context) {
	log.Println("GetShares (+)")
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lShares interface{}

	lShares, lErr = ConstructGetShares(c.svc)
	if lErr != nil {
		lStatus = http.StatusInternalServerError
		lCode = "FETCH_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateGetShares(lCtx, lShares)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteGetShares(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
	log.Println("GetShares (-)")
}

// ========================== MARKET OPEN UTILS ==========================

// isMarketOpen logic moved to service layer

// ========================== BUY SHARE ==========================

func CollectBuyShare(pCtx *gin.Context, pReq *TradeRequest) (string, error) {
	log.Println("CollectBuyShare (+)")
	lUserID, _ := pCtx.Get("userID")
	lErr := pCtx.ShouldBindJSON(pReq)
	log.Println("CollectBuyShare (-)")
	return lUserID.(string), lErr
}

func ValidateBuyShare(pReq TradeRequest) (bool, error) {
	log.Println("ValidateBuyShare (+)")
	lIsPending := !IsMarketOpen()
	log.Println("ValidateBuyShare (-)")
	return lIsPending, nil
}

func ConstructBuyShare(pSvc TradeService, pUserID string, pReq TradeRequest, pIsPending bool) error {
	log.Println("ConstructBuyShare (+)")
	lErr := pSvc.BuyShare(pUserID, pReq, pIsPending)
	log.Println("ConstructBuyShare (-)")
	return lErr
}

func CommunicateBuyShare(pCtx *gin.Context, pIsPending bool) error {
	log.Println("CommunicateBuyShare (+)")
	msg := "Share bought successfully"
	if pIsPending {
		msg = "Order placed (Pending)"
	}
	response.Success(pCtx, http.StatusOK, msg, nil)
	log.Println("CommunicateBuyShare (-)")
	return nil
}

func CompleteBuyShare(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	log.Println("CompleteBuyShare (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	log.Println("CompleteBuyShare (-)")
}

func (c *TradeController) BuyShare(lCtx *gin.Context) {
	log.Println("BuyShare (+)")
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq TradeRequest
	var lUserID string

	var lIsPending bool

	lUserID, lErr = CollectBuyShare(lCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input"
		lDetails = lErr.Error()
		goto Complete
	}

	lIsPending, lErr = ValidateBuyShare(lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = ConstructBuyShare(c.svc, lUserID, lReq, lIsPending)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "TRADE_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateBuyShare(lCtx, lIsPending)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteBuyShare(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
	log.Println("BuyShare (-)")
}

// ========================== SELL SHARE ==========================

func CollectSellShare(pCtx *gin.Context, pReq *TradeRequest) (string, error) {
	log.Println("CollectSellShare (+)")
	lUserID, _ := pCtx.Get("userID")
	lErr := pCtx.ShouldBindJSON(pReq)
	log.Println("CollectSellShare (-)")
	return lUserID.(string), lErr
}

func ValidateSellShare(pReq TradeRequest) (bool, error) {
	log.Println("ValidateSellShare (+)")
	lIsPending := !IsMarketOpen()
	log.Println("ValidateSellShare (-)")
	return lIsPending, nil
}

func ConstructSellShare(pSvc TradeService, pUserID string, pReq TradeRequest, pIsPending bool) error {
	log.Println("ConstructSellShare (+)")
	lErr := pSvc.SellShare(pUserID, pReq, pIsPending)
	log.Println("ConstructSellShare (-)")
	return lErr
}

func CommunicateSellShare(pCtx *gin.Context, pIsPending bool) error {
	log.Println("CommunicateSellShare (+)")
	msg := "Share sold successfully"
	if pIsPending {
		msg = "Order placed (Pending)"
	}
	response.Success(pCtx, http.StatusOK, msg, nil)
	log.Println("CommunicateSellShare (-)")
	return nil
}

func CompleteSellShare(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	log.Println("CompleteSellShare (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	log.Println("CompleteSellShare (-)")
}

func (c *TradeController) SellShare(lCtx *gin.Context) {
	log.Println("SellShare (+)")
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq TradeRequest
	var lUserID string

	var lIsPending bool

	lUserID, lErr = CollectSellShare(lCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input"
		lDetails = lErr.Error()
		goto Complete
	}

	lIsPending, lErr = ValidateSellShare(lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = ConstructSellShare(c.svc, lUserID, lReq, lIsPending)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "TRADE_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateSellShare(lCtx, lIsPending)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteSellShare(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
	log.Println("SellShare (-)")
}

// ========================== GET USER TRADES ==========================

func CollectGetUserTrades(pCtx *gin.Context) (string, error) {
	log.Println("CollectGetUserTrades (+)")
	lUserID, lExists := pCtx.Get("userID")
	if !lExists {
		log.Println("CollectGetUserTrades (-)")
		return "", fmt.Errorf("User not found in context")
	}
	log.Println("CollectGetUserTrades (-)")
	return lUserID.(string), nil
}

func ConstructGetUserTrades(pSvc TradeService, pUserID string) (interface{}, error) {
	log.Println("ConstructGetUserTrades (+)")
	lTrades, lErr := pSvc.GetUserTrades(pUserID)
	log.Println("ConstructGetUserTrades (-)")
	return lTrades, lErr
}

func CommunicateGetUserTrades(pCtx *gin.Context, pTrades interface{}) error {
	log.Println("CommunicateGetUserTrades (+)")
	response.Success(pCtx, http.StatusOK, "User trades retrieved", pTrades)
	log.Println("CommunicateGetUserTrades (-)")
	return nil
}

func CompleteGetUserTrades(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	log.Println("CompleteGetUserTrades (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	log.Println("CompleteGetUserTrades (-)")
}

func (c *TradeController) GetUserTrades(lCtx *gin.Context) {
	log.Println("GetUserTrades (+)")
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lUserID string
	var lTrades interface{}

	lUserID, lErr = CollectGetUserTrades(lCtx)
	if lErr != nil {
		lStatus = http.StatusUnauthorized
		lCode = "UNAUTHORIZED"
		lMsg = lErr.Error()
		goto Complete
	}

	lTrades, lErr = ConstructGetUserTrades(c.svc, lUserID)
	if lErr != nil {
		lStatus = http.StatusInternalServerError
		lCode = "TRADE_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateGetUserTrades(lCtx, lTrades)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteGetUserTrades(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
	log.Println("GetUserTrades (-)")
}

// ========================== CANCEL TRADE ==========================

func CollectCancelTrade(pCtx *gin.Context) (string, string, error) {
	log.Println("CollectCancelTrade (+)")
	lUserID, lExists := pCtx.Get("userID")
	if !lExists {
		log.Println("CollectCancelTrade (-)")
		return "", "", fmt.Errorf("User not found in context")
	}
	lTradeID := pCtx.Param("id")
	if lTradeID == "" {
		log.Println("CollectCancelTrade (-)")
		return "", "", fmt.Errorf("Trade ID is required")
	}
	log.Println("CollectCancelTrade (-)")
	return lUserID.(string), lTradeID, nil
}

func ValidateCancelTrade() error {
	return nil
}

func ConstructCancelTrade(pSvc TradeService, pUserID, pTradeID string) error {
	log.Println("ConstructCancelTrade (+)")
	lErr := pSvc.CancelTrade(pUserID, pTradeID)
	log.Println("ConstructCancelTrade (-)")
	return lErr
}

func CommunicateCancelTrade(pCtx *gin.Context) error {
	log.Println("CommunicateCancelTrade (+)")
	response.Success(pCtx, http.StatusOK, "Trade cancelled successfully", nil)
	log.Println("CommunicateCancelTrade (-)")
	return nil
}

func CompleteCancelTrade(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	log.Println("CompleteCancelTrade (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	log.Println("CompleteCancelTrade (-)")
}

func (c *TradeController) CancelTrade(lCtx *gin.Context) {
	log.Println("CancelTrade (+)")
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lUserID, lTradeID string

	lUserID, lTradeID, lErr = CollectCancelTrade(lCtx)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = ValidateCancelTrade()
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = ConstructCancelTrade(c.svc, lUserID, lTradeID)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "TRADE_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateCancelTrade(lCtx)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteCancelTrade(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
	log.Println("CancelTrade (-)")
}
