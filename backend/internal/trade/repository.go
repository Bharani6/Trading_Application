package trade

import (
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
