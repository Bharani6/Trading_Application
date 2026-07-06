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

func (sc *SupportController) SubmitMessage(c *gin.Context) {
	var req SupportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format or missing required fields", "details": err.Error()})
		return
	}

	msg := SupportMessage{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
		Status:  "Open",
	}

	if err := database.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save support message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Support message submitted successfully",
		"data":    msg,
	})
}

func (sc *SupportController) GetMessages(c *gin.Context) {
	var messages []SupportMessage
	if err := database.DB.Order("created_at desc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch support messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
