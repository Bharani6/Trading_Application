package watchlist

import (
	"errors"
	"strings"
	"stock-trading/internal/trade"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WatchlistService struct {
	repo      *WatchlistRepository
	tradeRepo *trade.TradeRepository // I'll use GORM directly to get share
}

func NewWatchlistService(repo *WatchlistRepository) *WatchlistService {
	return &WatchlistService{repo: repo}
}

func (pService *WatchlistService) AddStock(pUserID string, pReq AddWatchlistRequest) error {
	lUserUUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return lErr
	}
	lShareUUID, lErr := uuid.Parse(pReq.StockID)
	if lErr != nil {
		return lErr
	}

	lWatchlist := &Watchlist{
		UserID:  lUserUUID,
		ShareID: lShareUUID,
	}
	
	lErr = pService.repo.Create(lWatchlist)
	if lErr != nil {
		if errors.Is(lErr, gorm.ErrDuplicatedKey) || strings.Contains(lErr.Error(), "idx_user_share") || strings.Contains(lErr.Error(), "duplicate key value") {
			return errors.New("Stock already in watchlist")
		}
		return lErr
	}
	return nil
}

func (pService *WatchlistService) RemoveStock(pUserID string, pID string) error {
	return pService.repo.Delete(pID, pUserID)
}

func (pService *WatchlistService) UpdateFavorite(pUserID string, pID string, pIsFavorite bool) error {
	return pService.repo.UpdateFavorite(pID, pUserID, pIsFavorite)
}

func (pService *WatchlistService) GetUserWatchlist(pUserID string) ([]WatchlistResponse, error) {
	// Let's implement joining in repository, or do it here with db.
	// We will create a specific query in the repo for this to keep it clean.
	return pService.repo.GetUserWatchlistWithDetails(pUserID)
}
