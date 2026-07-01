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

func (c *AdminController) GetUsers(ctx *gin.Context) {
	users, err := c.svc.GetUsers()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to fetch users", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Users retrieved", users)
}

func (c *AdminController) ApproveUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.ApproveUser(id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to approve user", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "User approved", nil)
}

func (c *AdminController) RejectUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.RejectUser(id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to reject user", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "User rejected", nil)
}

func (c *AdminController) BlockUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.BlockUser(id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to block user", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "User blocked", nil)
}

func (c *AdminController) UploadShares(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "UPLOAD_ERROR", "Failed to parse file", err.Error())
		return
	}
	defer file.Close()

	if err := c.svc.UploadShares(file); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "UPLOAD_ERROR", "Failed to process shares upload", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Shares uploaded successfully", nil)
}

func (c *AdminController) DeleteAllShares(ctx *gin.Context) {
	if err := c.svc.DeleteAllShares(); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "ADMIN_ERROR", "Failed to delete shares", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "All shares deleted successfully", nil)
}
