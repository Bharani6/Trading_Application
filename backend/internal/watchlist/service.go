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

func (s *WatchlistService) AddStock(userID string, req AddWatchlistRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	shareUUID, err := uuid.Parse(req.StockID)
	if err != nil {
		return err
	}

	watchlist := &Watchlist{
		UserID:  userUUID,
		ShareID: shareUUID,
	}
	
	err = s.repo.Create(watchlist)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "idx_user_share") || strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("Stock already in watchlist")
		}
		return err
	}
	return nil
}

func (s *WatchlistService) RemoveStock(userID string, id string) error {
	return s.repo.Delete(id, userID)
}

func (s *WatchlistService) UpdateFavorite(userID string, id string, isFavorite bool) error {
	return s.repo.UpdateFavorite(id, userID, isFavorite)
}

func (s *WatchlistService) GetUserWatchlist(userID string) ([]WatchlistResponse, error) {
	// Let's implement joining in repository, or do it here with db.
	// We will create a specific query in the repo for this to keep it clean.
	return s.repo.GetUserWatchlistWithDetails(userID)
}
