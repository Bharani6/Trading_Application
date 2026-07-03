package watchlist

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Watchlist struct {
	ID         uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:char(36);uniqueIndex:idx_user_share;not null"`
	ShareID    uuid.UUID      `json:"share_id" gorm:"type:char(36);uniqueIndex:idx_user_share;not null"`
	IsFavorite bool           `json:"is_favorite" gorm:"default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (w *Watchlist) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}

type AddWatchlistRequest struct {
	StockID string `json:"stockId" binding:"required"`
}

type UpdateFavoriteRequest struct {
	IsFavorite *bool `json:"isFavorite" binding:"required"`
}

type WatchlistResponse struct {
	ID           string  `json:"id"`
	StockID      string  `json:"stockId"`
	StockName    string  `json:"stockName"`
	Symbol       string  `json:"symbol"`
	CurrentPrice  float64 `json:"currentPrice"`
	PreviousPrice float64 `json:"previousPrice"`
	IsFavorite    bool    `json:"isFavorite"`
}
