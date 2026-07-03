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

func (ctrl *WatchlistController) AddStock(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	var req AddWatchlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
		return
	}

	if err := ctrl.service.AddStock(userID.(string), req); err != nil {
		if err.Error() == "Stock already in watchlist" {
			response.Error(c, http.StatusConflict, "CONFLICT", err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add stock to watchlist", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Stock added to watchlist successfully", nil)
}

func (ctrl *WatchlistController) RemoveStock(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Watchlist item ID is required", nil)
		return
	}

	if err := ctrl.service.RemoveStock(userID.(string), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove stock from watchlist", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Removed from watchlist", nil)
}

func (ctrl *WatchlistController) GetWatchlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	watchlists, err := ctrl.service.GetUserWatchlist(userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch watchlist", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Watchlist fetched successfully", watchlists)
}

func (ctrl *WatchlistController) UpdateFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Watchlist item ID is required", nil)
		return
	}

	var req UpdateFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
		return
	}

	if req.IsFavorite == nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "isFavorite field is required", nil)
		return
	}

	if err := ctrl.service.UpdateFavorite(userID.(string), id, *req.IsFavorite); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update favorite status", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Favorite status updated", nil)
}
