package support

import (
	"net/http"
	"stock-trading/internal/database"

	"github.com/gin-gonic/gin"
)

type SupportController struct{}

func NewSupportController() *SupportController {
	return &SupportController{}
}

type SupportRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

func (pController *SupportController) SubmitMessage(lCtx *gin.Context) {
	var lReq SupportRequest
	if lErr := lCtx.ShouldBindJSON(&lReq); lErr != nil {
		lCtx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format or missing required fields", "details": lErr.Error()})
		return
	}

	lMsg := SupportMessage{
		Name:    lReq.Name,
		Email:   lReq.Email,
		Message: lReq.Message,
		Status:  "Open",
	}

	if lErr := database.DB.Create(&lMsg).Error; lErr != nil {
		lCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save support message"})
		return
	}

	lCtx.JSON(http.StatusCreated, gin.H{
		"message": "Support message submitted successfully",
		"data":    lMsg,
	})
}

func (pController *SupportController) GetMessages(lCtx *gin.Context) {
	var lMessages []SupportMessage
	if lErr := database.DB.Order("created_at desc").Find(&lMessages).Error; lErr != nil {
		lCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch support messages"})
		return
	}

	lCtx.JSON(http.StatusOK, gin.H{"messages": lMessages})
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (pController *SupportController) UpdateStatus(lCtx *gin.Context) {
	lID := lCtx.Param("id")
	var lReq UpdateStatusRequest
	if lErr := lCtx.ShouldBindJSON(&lReq); lErr != nil {
		lCtx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": lErr.Error()})
		return
	}

	var lMsg SupportMessage
	if lErr := database.DB.First(&lMsg, lID).Error; lErr != nil {
		lCtx.JSON(http.StatusNotFound, gin.H{"error": "Support message not found"})
		return
	}

	lMsg.Status = lReq.Status
	if lErr := database.DB.Save(&lMsg).Error; lErr != nil {
		lCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	lCtx.JSON(http.StatusOK, gin.H{
		"message": "Status updated successfully",
		"data":    lMsg,
	})
}
