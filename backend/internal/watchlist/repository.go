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

func (r *WatchlistRepository) Create(watchlist *Watchlist) error {
	// Check if a soft-deleted record already exists
	var existing Watchlist
	err := r.db.Unscoped().Where("user_id = ? AND share_id = ?", watchlist.UserID, watchlist.ShareID).First(&existing).Error
	if err == nil && existing.DeletedAt.Valid {
		// Restore it
		existing.DeletedAt = gorm.DeletedAt{Valid: false}
		existing.IsFavorite = false
		return r.db.Unscoped().Save(&existing).Error
	}
	return r.db.Create(watchlist).Error
}

func (r *WatchlistRepository) Delete(id, userID string) error {
	result := r.db.Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&Watchlist{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *WatchlistRepository) GetByUserID(userID string) ([]Watchlist, error) {
	var watchlists []Watchlist
	err := r.db.Where("user_id = ?", userID).Find(&watchlists).Error
	return watchlists, err
}

func (r *WatchlistRepository) GetUserWatchlistWithDetails(userID string) ([]WatchlistResponse, error) {
	var responses []WatchlistResponse
	
	err := r.db.Table("watchlists").
		Select("watchlists.id, watchlists.share_id as stock_id, shares.name as stock_name, shares.symbol, shares.price as current_price, shares.previous_price as previous_price, watchlists.is_favorite").
		Joins("JOIN shares ON shares.id = watchlists.share_id").
		Where("watchlists.user_id = ? AND watchlists.deleted_at IS NULL", userID).
		Scan(&responses).Error
		
	return responses, err
}

func (r *WatchlistRepository) UpdateFavorite(id, userID string, isFavorite bool) error {
	result := r.db.Model(&Watchlist{}).Where("id = ? AND user_id = ?", id, userID).Update("is_favorite", isFavorite)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
