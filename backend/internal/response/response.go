package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, code string, message string, err interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   err,
	})
}

func InternalServerError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred", err.Error())
}
