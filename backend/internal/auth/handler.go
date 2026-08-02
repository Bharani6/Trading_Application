package auth

import (
	"fmt"
	"net/http"
	"strings"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{svc: NewAuthService()}
}

// ========================== REGISTER ==========================

func CollectRegister(pCtx *gin.Context, pReq *RegisterRequest) error {
	fmt.Println("CollectRegister (+)")
	lErr := pCtx.ShouldBindJSON(pReq)
	fmt.Println("CollectRegister (-)")
	return lErr
}

func ConstructRegister(pSvc AuthService, pReq RegisterRequest) (interface{}, error) {
	fmt.Println("ConstructRegister (+)")
	lUser, lErr := pSvc.Register(pReq)
	fmt.Println("ConstructRegister (-)")
	return lUser, lErr
}

func CommunicateRegister(pCtx *gin.Context, pUser interface{}) error {
	fmt.Println("CommunicateRegister (+)")
	response.Success(pCtx, http.StatusCreated, "User registered successfully", pUser)
	fmt.Println("CommunicateRegister (-)")
	return nil
}

func CompleteRegister(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	fmt.Println("CompleteRegister (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	fmt.Println("CompleteRegister (-)")
}

func (pController *AuthController) Register(pCtx *gin.Context) {
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq RegisterRequest
	var lUser interface{}

	lErr = CollectRegister(pCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input parameters"
		lDetails = lErr.Error()
		goto Complete
	}

	lUser, lErr = ConstructRegister(pController.svc, lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "REGISTRATION_FAILED"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateRegister(pCtx, lUser)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteRegister(pCtx, lErr, lStatus, lCode, lMsg, lDetails)
}

// ========================== LOGIN ==========================

func CollectLogin(pCtx *gin.Context, pReq *LoginRequest) error {
	fmt.Println("CollectLogin (+)")
	lErr := pCtx.ShouldBindJSON(pReq)
	fmt.Println("CollectLogin (-)")
	return lErr
}

func ConstructLogin(pCtx *gin.Context, pSvc AuthService, pReq LoginRequest) (*LoginResponse, error) {
	fmt.Println("ConstructLogin (+)")
	lTokens, lErr := pSvc.Login(pReq, pCtx.ClientIP(), pCtx.GetHeader("User-Agent"))
	fmt.Println("ConstructLogin (-)")
	return lTokens, lErr
}

func CommunicateLogin(pCtx *gin.Context, pTokens *LoginResponse) error {
	fmt.Println("CommunicateLogin (+)")
	pCtx.SetCookie("access_token", pTokens.AccessToken, 86400, "/", "", false, true)
	pCtx.SetCookie("refresh_token", pTokens.RefreshToken, 86400*7, "/", "", false, true)
	response.Success(pCtx, http.StatusOK, "Login successful", pTokens)
	fmt.Println("CommunicateLogin (-)")
	return nil
}

func CompleteLogin(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	fmt.Println("CompleteLogin (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	fmt.Println("CompleteLogin (-)")
}

func (pController *AuthController) Login(pCtx *gin.Context) {
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq LoginRequest
	var lTokens *LoginResponse

	lErr = CollectLogin(pCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input parameters"
		lDetails = lErr.Error()
		goto Complete
	}

	lTokens, lErr = ConstructLogin(pCtx, pController.svc, lReq)
	if lErr != nil {
		lStatus = http.StatusUnauthorized
		lCode = "LOGIN_FAILED"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateLogin(pCtx, lTokens)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteLogin(pCtx, lErr, lStatus, lCode, lMsg, lDetails)
}

// ========================== GET ME ==========================

func CollectGetMe(pCtx *gin.Context) (string, error) {
	fmt.Println("CollectGetMe (+)")
	lUserID, lExists := pCtx.Get("userID")
	if !lExists {
		fmt.Println("CollectGetMe (-)")
		return "", fmt.Errorf("User not found")
	}
	fmt.Println("CollectGetMe (-)")
	return lUserID.(string), nil
}

func ConstructGetMe(pSvc AuthService, pUserID string) (interface{}, error) {
	fmt.Println("ConstructGetMe (+)")
	lUser, lErr := pSvc.GetMe(pUserID)
	fmt.Println("ConstructGetMe (-)")
	return lUser, lErr
}

func CommunicateGetMe(pCtx *gin.Context, pUser interface{}) error {
	fmt.Println("CommunicateGetMe (+)")
	response.Success(pCtx, http.StatusOK, "User fetched", pUser)
	fmt.Println("CommunicateGetMe (-)")
	return nil
}

func CompleteGetMe(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	fmt.Println("CompleteGetMe (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	fmt.Println("CompleteGetMe (-)")
}

func (pController *AuthController) GetMe(pCtx *gin.Context) {
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lUserID string
	var lUser interface{}

	lUserID, lErr = CollectGetMe(pCtx)
	if lErr != nil {
		lStatus = http.StatusUnauthorized
		lCode = "UNAUTHORIZED"
		lMsg = "User not found"
		goto Complete
	}

	lUser, lErr = ConstructGetMe(pController.svc, lUserID)
	if lErr != nil {
		lStatus = http.StatusInternalServerError
		lCode = "USER_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateGetMe(pCtx, lUser)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteGetMe(pCtx, lErr, lStatus, lCode, lMsg, lDetails)
}

// ========================== FORGOT PASSWORD ==========================

func (pController *AuthController) ForgotPassword(pContext *gin.Context) {
	var lReq ForgotPasswordRequest
	if lErr := pContext.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(pContext, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", lErr.Error())
		return
	}

	lTokenStr, lErr := pController.svc.ForgotPassword(lReq)
	if lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "FORGOT_PASSWORD_ERROR", lErr.Error(), nil)
		return
	}

	// Token string is sent back to mock email functionality
	response.Success(pContext, http.StatusOK, "Password reset link has been generated", gin.H{"mock_token": lTokenStr})
}

// ========================== VERIFY RESET TOKEN ==========================

func (pController *AuthController) VerifyResetToken(pContext *gin.Context) {
	var lReq VerifyResetTokenRequest
	if lErr := pContext.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(pContext, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", lErr.Error())
		return
	}

	if lErr := pController.svc.VerifyResetToken(lReq); lErr != nil {
		response.Error(pContext, http.StatusBadRequest, "INVALID_TOKEN", lErr.Error(), nil)
		return
	}

	response.Success(pContext, http.StatusOK, "Token is valid", gin.H{"valid": true})
}

// ========================== RESET PASSWORD ==========================

func (pController *AuthController) ResetPassword(pContext *gin.Context) {
	var lReq ResetPasswordRequest
	if lErr := pContext.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(pContext, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", lErr.Error())
		return
	}

	if lErr := pController.svc.ResetPassword(lReq); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "RESET_PASSWORD_ERROR", lErr.Error(), nil)
		return
	}

	response.Success(pContext, http.StatusOK, "Password reset successful", nil)
}

func (pController *AuthController) GetSessions(pContext *gin.Context) {
	lUserID, lExists := pContext.Get("userID")
	if !lExists {
		response.Error(pContext, http.StatusUnauthorized, "UNAUTHORIZED", "User not logged in", nil)
		return
	}

	lTokenString := ""
	lAuthHeader := pContext.GetHeader("Authorization")
	if lAuthHeader != "" && strings.HasPrefix(lAuthHeader, "Bearer ") {
		lTokenString = strings.TrimPrefix(lAuthHeader, "Bearer ")
	}
	if lTokenString == "" {
		lCookie, lErr := pContext.Cookie("access_token")
		if lErr == nil {
			lTokenString = lCookie
		}
	}

	lSessions, lErr := pController.svc.GetActiveSessions(lUserID.(string), lTokenString)
	if lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "SESSIONS_ERROR", "Failed to retrieve sessions", lErr.Error())
		return
	}

	response.Success(pContext, http.StatusOK, "Sessions retrieved successfully", lSessions)
}

func (pController *AuthController) RevokeSession(pContext *gin.Context) {
	lSessionID := pContext.Param("id")
	if lErr := pController.svc.RevokeSession(lSessionID); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "REVOKE_SESSION_ERROR", "Failed to revoke session", lErr.Error())
		return
	}

	response.Success(pContext, http.StatusOK, "Session revoked successfully", nil)
}
