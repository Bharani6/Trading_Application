package watchlist

import (
	"stock-trading/internal/database"

	"gorm.io/gorm"
)

type WatchlistRepository struct {
	db *gorm.DB
}

func NewWatchlistRepository() *WatchlistRepository {
	return &WatchlistRepository{db: database.DB}
}

func (pRepo *WatchlistRepository) Create(pWatchlist *Watchlist) error {
	// Check if a soft-deleted record already exists
	var lExisting Watchlist
	lErr := pRepo.db.Unscoped().Where("user_id = ? AND share_id = ?", pWatchlist.UserID, pWatchlist.ShareID).First(&lExisting).Error
	if lErr == nil && lExisting.DeletedAt.Valid {
		// Restore it
		lExisting.DeletedAt = gorm.DeletedAt{Valid: false}
		lExisting.IsFavorite = false
		return pRepo.db.Unscoped().Save(&lExisting).Error
	}
	return pRepo.db.Create(pWatchlist).Error
}

func (pRepo *WatchlistRepository) Delete(pID, pUserID string) error {
	lResult := pRepo.db.Unscoped().Where("id = ? AND user_id = ?", pID, pUserID).Delete(&Watchlist{})
	if lResult.Error != nil {
		return lResult.Error
	}
	if lResult.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (pRepo *WatchlistRepository) GetByUserID(pUserID string) ([]Watchlist, error) {
	var lWatchlists []Watchlist
	lErr := pRepo.db.Where("user_id = ?", pUserID).Find(&lWatchlists).Error
	return lWatchlists, lErr
}

func (pRepo *WatchlistRepository) GetUserWatchlistWithDetails(pUserID string) ([]WatchlistResponse, error) {
	var lResponses []WatchlistResponse
	
	lErr := pRepo.db.Table("watchlists").
		Select("watchlists.id, watchlists.share_id as stock_id, shares.name as stock_name, shares.symbol, shares.price as current_price, shares.previous_price as previous_price, watchlists.is_favorite").
		Joins("JOIN shares ON shares.id = watchlists.share_id").
		Where("watchlists.user_id = ? AND watchlists.deleted_at IS NULL", pUserID).
		Scan(&lResponses).Error
		
	return lResponses, lErr
}

func (pRepo *WatchlistRepository) UpdateFavorite(pID, pUserID string, pIsFavorite bool) error {
	lResult := pRepo.db.Model(&Watchlist{}).Where("id = ? AND user_id = ?", pID, pUserID).Update("is_favorite", pIsFavorite)
	if lResult.Error != nil {
		return lResult.Error
	}
	if lResult.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
