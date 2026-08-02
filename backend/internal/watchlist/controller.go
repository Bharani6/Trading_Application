package watchlist

import (
	"net/http"
	"stock-trading/internal/response"

	"github.com/gin-gonic/gin"
)

type WatchlistController struct {
	service *WatchlistService
}

func NewWatchlistController() *WatchlistController {
	return &WatchlistController{
		service: NewWatchlistService(NewWatchlistRepository()),
	}
}

func (pController *WatchlistController) AddStock(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	if !lExists {
		response.Error(lCtx, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	var lReq AddWatchlistRequest
	if lErr := lCtx.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", lErr.Error())
		return
	}

	if lErr := pController.service.AddStock(lUserID.(string), lReq); lErr != nil {
		if lErr.Error() == "Stock already in watchlist" {
			response.Error(lCtx, http.StatusConflict, "CONFLICT", lErr.Error(), nil)
			return
		}
		response.Error(lCtx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add stock to watchlist", lErr.Error())
		return
	}

	response.Success(lCtx, http.StatusCreated, "Stock added to watchlist successfully", nil)
}

func (pController *WatchlistController) RemoveStock(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	if !lExists {
		response.Error(lCtx, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	lID := lCtx.Param("id")
	if lID == "" {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Watchlist item ID is required", nil)
		return
	}

	if lErr := pController.service.RemoveStock(lUserID.(string), lID); lErr != nil {
		response.Error(lCtx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove stock from watchlist", lErr.Error())
		return
	}

	response.Success(lCtx, http.StatusOK, "Removed from watchlist", nil)
}

func (pController *WatchlistController) GetWatchlist(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	if !lExists {
		response.Error(lCtx, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	lWatchlists, lErr := pController.service.GetUserWatchlist(lUserID.(string))
	if lErr != nil {
		response.Error(lCtx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch watchlist", lErr.Error())
		return
	}

	response.Success(lCtx, http.StatusOK, "Watchlist fetched successfully", lWatchlists)
}

func (pController *WatchlistController) UpdateFavorite(lCtx *gin.Context) {
	lUserID, lExists := lCtx.Get("userID")
	if !lExists {
		response.Error(lCtx, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	lID := lCtx.Param("id")
	if lID == "" {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Watchlist item ID is required", nil)
		return
	}

	var lReq UpdateFavoriteRequest
	if lErr := lCtx.ShouldBindJSON(&lReq); lErr != nil {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", lErr.Error())
		return
	}

	if lReq.IsFavorite == nil {
		response.Error(lCtx, http.StatusBadRequest, "BAD_REQUEST", "isFavorite field is required", nil)
		return
	}

	if lErr := pController.service.UpdateFavorite(lUserID.(string), lID, *lReq.IsFavorite); lErr != nil {
		response.Error(lCtx, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update favorite status", lErr.Error())
		return
	}

	response.Success(lCtx, http.StatusOK, "Favorite status updated", nil)
}
