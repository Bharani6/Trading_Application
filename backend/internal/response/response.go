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

func Success(pCtx *gin.Context, pStatusCode int, pMessage string, pData interface{}) {
	pCtx.JSON(pStatusCode, SuccessResponse{
		Success: true,
		Message: pMessage,
		Data:    pData,
	})
}

func Error(pCtx *gin.Context, pStatusCode int, pCode string, pMessage string, pErr interface{}) {
	pCtx.JSON(pStatusCode, ErrorResponse{
		Success: false,
		Code:    pCode,
		Message: pMessage,
		Error:   pErr,
	})
}

func InternalServerError(pCtx *gin.Context, pErr error) {
	Error(pCtx, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred", pErr.Error())
}
