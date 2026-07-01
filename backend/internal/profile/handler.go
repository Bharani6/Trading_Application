package profile

import (
	"fmt"
	"net/http"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	svc ProfileService
}

func NewProfileController() *ProfileController {
	return &ProfileController{svc: NewProfileService()}
}

func CollectSubmitKYC(pCtx *gin.Context, pReq *KYCSubmitRequest) error {
	fmt.Println("CollectSubmitKYC (+)")
	err := pCtx.ShouldBindJSON(pReq)
	fmt.Println("CollectSubmitKYC (-)")
	return err
}

func ValidateSubmitKYC(pCtx *gin.Context, pReq *KYCSubmitRequest) error {
	fmt.Println("ValidateSubmitKYC (+)")
	fmt.Println("ValidateSubmitKYC (-)")
	return nil
}

func ConstructSubmitKYC(pCtx *gin.Context, pSvc ProfileService, pUserID string, pReq KYCSubmitRequest) error {
	fmt.Println("ConstructSubmitKYC (+)")
	err := pSvc.SubmitKYC(pUserID, pReq)
	fmt.Println("ConstructSubmitKYC (-)")
	return err
}

func CommunicateSubmitKYC(pCtx *gin.Context) error {
	fmt.Println("CommunicateSubmitKYC (+)")
	response.Success(pCtx, http.StatusOK, "KYC details submitted successfully. Pending admin approval.", nil)
	fmt.Println("CommunicateSubmitKYC (-)")
	return nil
}

func CompleteSubmitKYC(pCtx *gin.Context, pErr error, pStatus int, pCode, pMsg string, pDetails interface{}) {
	fmt.Println("CompleteSubmitKYC (+)")
	if pErr != nil {
		response.Error(pCtx, pStatus, pCode, pMsg, pDetails)
	}
	fmt.Println("CompleteSubmitKYC (-)")
}

func (c *ProfileController) SubmitKYC(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	var lErr error
	var lStatus int
	var lCode, lMsg string
	var lDetails interface{}
	var lReq KYCSubmitRequest

	if !lExists {
		lErr = fmt.Errorf("User not found in context")
		lStatus = http.StatusUnauthorized
		lCode = "UNAUTHORIZED"
		lMsg = "User not found in context"
		goto Complete
	}

	lErr = CollectSubmitKYC(lCtx, &lReq)
	if lErr != nil {
		lStatus = http.StatusBadRequest
		lCode = "VALIDATION_ERROR"
		lMsg = "Invalid input parameters"
		lDetails = lErr.Error()
		goto Complete
	}

	lErr = ValidateSubmitKYC(lCtx, &lReq)
	if lErr != nil {
		goto Complete
	}

	lErr = ConstructSubmitKYC(lCtx, c.svc, lUserID.(string), lReq)
	if lErr != nil {
		lStatus = http.StatusInternalServerError
		lCode = "KYC_FAILED"
		lMsg = lErr.Error()
		goto Complete
	}

	lErr = CommunicateSubmitKYC(lCtx)
	if lErr != nil {
		goto Complete
	}

Complete:
	CompleteSubmitKYC(lCtx, lErr, lStatus, lCode, lMsg, lDetails)
}
