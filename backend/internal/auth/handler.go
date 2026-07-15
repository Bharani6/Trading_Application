package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"sync"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

var (
	otpStore      = make(map[string]string)
	otpMutex      sync.RWMutex
	emailOtpStore = make(map[string]string)
	emailOtpMutex sync.RWMutex
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

func (c *AuthController) Register(lCtx *gin.Context) {
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq RegisterRequest
	var lUser interface{}

	lErr = CollectRegister(lCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input parameters"
		lDetails = lErr.Error()
		goto Complete
	}

	lUser, lErr = ConstructRegister(c.svc, lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "REGISTRATION_FAILED"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateRegister(lCtx, lUser)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteRegister(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
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

func (c *AuthController) Login(lCtx *gin.Context) {
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq LoginRequest
	var lTokens *LoginResponse

	lErr = CollectLogin(lCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input parameters"
		lDetails = lErr.Error()
		goto Complete
	}

	lTokens, lErr = ConstructLogin(lCtx, c.svc, lReq)
	if lErr != nil {
		lStatus = http.StatusUnauthorized
		lCode = "LOGIN_FAILED"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateLogin(lCtx, lTokens)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteLogin(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
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

func (c *AuthController) GetMe(lCtx *gin.Context) {
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lUserID string
	var lUser interface{}

	lUserID, lErr = CollectGetMe(lCtx)
	if lErr != nil {
		lStatus = http.StatusUnauthorized
		lCode = "UNAUTHORIZED"
		lMsg = "User not found"
		goto Complete
	}

	lUser, lErr = ConstructGetMe(c.svc, lUserID)
	if lErr != nil {
		lStatus = http.StatusInternalServerError
		lCode = "USER_ERROR"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateGetMe(lCtx, lUser)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteGetMe(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
}

// ========================== FORGOT PASSWORD ==========================

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	tokenStr, err := c.svc.ForgotPassword(req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "FORGOT_PASSWORD_ERROR", err.Error(), nil)
		return
	}

	// Token string is sent back to mock email functionality
	response.Success(ctx, http.StatusOK, "Password reset link has been generated", gin.H{"mock_token": tokenStr})
}

// ========================== VERIFY RESET TOKEN ==========================

func (c *AuthController) VerifyResetToken(ctx *gin.Context) {
	var req VerifyResetTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	if err := c.svc.VerifyResetToken(req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "INVALID_TOKEN", err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Token is valid", gin.H{"valid": true})
}

// ========================== RESET PASSWORD ==========================

func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	if err := c.svc.ResetPassword(req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "RESET_PASSWORD_ERROR", err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "Password updated successfully", nil)
}

// ========================== OTP METHODS ==========================

func (c *AuthController) SendOTP(ctx *gin.Context) {
	var req SendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "OTP_GENERATION_FAILED", "Failed to generate secure OTP", err.Error())
		return
	}
	otp := fmt.Sprintf("%06d", n.Int64())

	otpMutex.Lock()
	otpStore[req.Mobile] = otp
	otpMutex.Unlock()

	// In production, an SMS gateway would be integrated here.
	response.Success(ctx, http.StatusOK, "OTP sent successfully", gin.H{"mock_otp": otp})
}

func (c *AuthController) VerifyOTP(ctx *gin.Context) {
	var req VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	otpMutex.RLock()
	storedOTP, exists := otpStore[req.Mobile]
	otpMutex.RUnlock()

	if !exists || storedOTP != req.OTP {
		response.Error(ctx, http.StatusBadRequest, "INVALID_OTP", "The OTP entered is incorrect or expired", nil)
		return
	}

	// Remove OTP after verification
	otpMutex.Lock()
	delete(otpStore, req.Mobile)
	otpMutex.Unlock()

	response.Success(ctx, http.StatusOK, "OTP verified successfully", nil)
}

func (c *AuthController) SendEmailOTP(ctx *gin.Context) {
	var req SendEmailOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "OTP_GENERATION_FAILED", "Failed to generate secure OTP", err.Error())
		return
	}
	otp := fmt.Sprintf("%06d", n.Int64())

	emailOtpMutex.Lock()
	emailOtpStore[req.Email] = otp
	emailOtpMutex.Unlock()

	// In production, an email service would be integrated here (e.g. SendGrid, AWS SES).
	response.Success(ctx, http.StatusOK, "OTP sent successfully to email", gin.H{"mock_otp": otp})
}

func (c *AuthController) VerifyEmailOTP(ctx *gin.Context) {
	var req VerifyEmailOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input parameters", err.Error())
		return
	}

	emailOtpMutex.RLock()
	storedOTP, exists := emailOtpStore[req.Email]
	emailOtpMutex.RUnlock()

	if !exists || storedOTP != req.OTP {
		response.Error(ctx, http.StatusBadRequest, "INVALID_OTP", "The OTP entered is incorrect or expired", nil)
		return
	}

	// Remove OTP after verification
	emailOtpMutex.Lock()
	delete(emailOtpStore, req.Email)
	emailOtpMutex.Unlock()

	response.Success(ctx, http.StatusOK, "Email OTP verified successfully", nil)
}
