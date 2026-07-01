package trade

import (
	"stock-trading/internal/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TradeRepository interface {
	GetAllShares() ([]Share, error)
	GetShareForUpdate(tx *gorm.DB, shareID string) (*Share, error)
	UpdateShare(tx *gorm.DB, share *Share) error
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

func (r *tradeRepository) GetAllShares() ([]Share, error) {
	var shares []Share
	err := r.db.Preload("Segment").Find(&shares).Error
	return shares, err
}

func (r *tradeRepository) GetShareForUpdate(tx *gorm.DB, shareID string) (*Share, error) {
	var share Share
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", shareID).First(&share).Error
	return &share, err
}

func (r *tradeRepository) UpdateShare(tx *gorm.DB, share *Share) error {
	return tx.Save(share).Error
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
	return r.db.Transaction(fn)
}
