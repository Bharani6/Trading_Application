package trade

import (
	"errors"
	"stock-trading/internal/database"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TradeRepository interface {
	GetAllShares(search string) ([]Share, error)
	GetShareForUpdate(tx *gorm.DB, shareID string) (*Share, error)
	FirstOrCreateShare(share *Share) error
	UpdateShare(tx *gorm.DB, share *Share) error
	UpdateShareWithVersion(tx *gorm.DB, share *Share) error
	CreateTrade(tx *gorm.DB, trade *Trade) error
	GetTradesByUser(userID string) ([]Trade, error)
	RunInTransaction(fn func(tx *gorm.DB) error) error
}

type tradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository() TradeRepository {
	return &tradeRepository{db: database.DB}
}

func (r *tradeRepository) GetAllShares(search string) ([]Share, error) {
	var shares []Share = make([]Share, 0)
	query := r.db.Preload("Segment")
	
	search = strings.TrimSpace(search)
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name ILIKE ? OR symbol ILIKE ?", searchTerm, searchTerm)
	}
	err := query.Find(&shares).Error
	return shares, err
}

func (r *tradeRepository) GetShareForUpdate(tx *gorm.DB, shareID string) (*Share, error) {
	var share Share
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", shareID).First(&share).Error
	return &share, err
}

func (r *tradeRepository) FirstOrCreateShare(share *Share) error {
	return r.db.Where("symbol = ?", share.Symbol).FirstOrCreate(share).Error
}

func (r *tradeRepository) UpdateShare(tx *gorm.DB, share *Share) error {
	return tx.Save(share).Error
}

func (r *tradeRepository) UpdateShareWithVersion(tx *gorm.DB, share *Share) error {
	result := tx.Model(share).Where("version = ?", share.Version).Updates(map[string]interface{}{
		"available_shares": share.AvailableShares,
		"version":          share.Version + 1,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("optimistic lock failed for share")
	}
	share.Version++
	return nil
}

func (r *tradeRepository) CreateTrade(tx *gorm.DB, trade *Trade) error {
	return tx.Create(trade).Error
}

func (r *tradeRepository) GetTradesByUser(userID string) ([]Trade, error) {
	var trades []Trade
	err := r.db.Preload("Share").Where("user_id = ?", userID).Order("created_at desc").Find(&trades).Error
	return trades, err
}

func (r *tradeRepository) RunInTransaction(fn func(tx *gorm.DB) error) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			panic(rec) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
