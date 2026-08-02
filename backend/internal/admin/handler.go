package admin

import (
	"net/http"

	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	svc AdminService
}

func NewAdminController() *AdminController {
	return &AdminController{svc: NewAdminService()}
}

func (pController *AdminController) GetUsers(pContext *gin.Context) {
	lUsers, lErr := pController.svc.GetUsers()
	if lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to fetch users", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "Users retrieved", lUsers)
}

func (pController *AdminController) GetUserDetails(pContext *gin.Context) {
	lID := pContext.Param("id")
	lUserDetails, lErr := pController.svc.GetUserDetails(lID)
	if lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to fetch user details", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "User details retrieved", lUserDetails)
}

func (pController *AdminController) ApproveUser(pContext *gin.Context) {
	lID := pContext.Param("id")
	if lErr := pController.svc.ApproveUser(lID); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to approve user", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "User approved", nil)
}

func (pController *AdminController) RejectUser(pContext *gin.Context) {
	lID := pContext.Param("id")
	if lErr := pController.svc.RejectUser(lID); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to reject user", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "User rejected", nil)
}

func (pController *AdminController) BlockUser(pContext *gin.Context) {
	lID := pContext.Param("id")
	if lErr := pController.svc.BlockUser(lID); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to block user", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "User blocked", nil)
}

func (pController *AdminController) CloseAccount(pContext *gin.Context) {
	lID := pContext.Param("id")
	if lErr := pController.svc.CloseAccount(lID); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to close account", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "Account closed successfully", nil)
}

func (pController *AdminController) RejectClosure(pContext *gin.Context) {
	lID := pContext.Param("id")
	if lErr := pController.svc.RejectClosure(lID); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to reject closure", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "Account closure rejected", nil)
}

func (pController *AdminController) UploadShares(pContext *gin.Context) {
	lFile, _, lErr := pContext.Request.FormFile("file")
	if lErr != nil {
		response.Error(pContext, http.StatusBadRequest, "UPLOAD_ERROR", "Failed to parse file", lErr.Error())
		return
	}
	defer lFile.Close()

	if lErr := pController.svc.UploadShares(lFile); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to process shares upload", lErr.Error())
		return
	}

	response.Success(pContext, http.StatusOK, "Shares uploaded successfully", nil)
}

func (pController *AdminController) DeleteAllShares(pContext *gin.Context) {
	if lErr := pController.svc.DeleteAllShares(); lErr != nil {
		response.Error(pContext, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to delete shares", lErr.Error())
		return
	}
	response.Success(pContext, http.StatusOK, "All shares deleted successfully", nil)
}
