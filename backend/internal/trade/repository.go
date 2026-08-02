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
	FirstOrCreateSegment(segment *Segment) error
	UpdateShare(tx *gorm.DB, share *Share) error
	UpdateShareWithVersion(tx *gorm.DB, share *Share) error
	CreateTrade(tx *gorm.DB, trade *Trade) error
	GetTradesByUser(userID string) ([]Trade, error)
	RunInTransaction(fn func(tx *gorm.DB) error) error
	GetSegmentByName(name string) (*Segment, error)
}

type tradeRepository struct {
	db *gorm.DB
}

func NewTradeRepository() TradeRepository {
	return &tradeRepository{db: database.DB}
}

func (pRepo *tradeRepository) GetAllShares(pSearch string) ([]Share, error) {
	var lShares []Share = make([]Share, 0)
	lQuery := pRepo.db.Preload("Segment")
	
	pSearch = strings.TrimSpace(pSearch)
	if pSearch != "" {
		lSearchTerm := "%" + pSearch + "%"
		lQuery = lQuery.Where("name ILIKE ? OR symbol ILIKE ?", lSearchTerm, lSearchTerm)
	}
	lErr := lQuery.Find(&lShares).Error
	return lShares, lErr
}

func (pRepo *tradeRepository) GetSegmentByName(pName string) (*Segment, error) {
	var lSegment Segment
	lErr := pRepo.db.Where("name = ?", pName).First(&lSegment).Error
	return &lSegment, lErr
}

func (pRepo *tradeRepository) GetShareForUpdate(pTx *gorm.DB, pShareID string) (*Share, error) {
	var lShare Share
	lErr := pTx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", pShareID).First(&lShare).Error
	return &lShare, lErr
}

func (pRepo *tradeRepository) FirstOrCreateShare(pShare *Share) error {
	return pRepo.db.Where("symbol = ?", pShare.Symbol).FirstOrCreate(pShare).Error
}

func (pRepo *tradeRepository) FirstOrCreateSegment(pSegment *Segment) error {
	return pRepo.db.Where("name = ?", pSegment.Name).FirstOrCreate(pSegment).Error
}

func (pRepo *tradeRepository) UpdateShare(pTx *gorm.DB, pShare *Share) error {
	return pTx.Save(pShare).Error
}

func (pRepo *tradeRepository) UpdateShareWithVersion(pTx *gorm.DB, pShare *Share) error {
	lResult := pTx.Model(pShare).Where("version = ?", pShare.Version).Updates(map[string]interface{}{
		"available_shares": pShare.AvailableShares,
		"version":          pShare.Version + 1,
	})
	if lResult.Error != nil {
		return lResult.Error
	}
	if lResult.RowsAffected == 0 {
		return errors.New("optimistic lock failed for share")
	}
	pShare.Version++
	return nil
}

func (pRepo *tradeRepository) CreateTrade(pTx *gorm.DB, pTrade *Trade) error {
	return pTx.Create(pTrade).Error
}

func (pRepo *tradeRepository) GetTradesByUser(pUserID string) ([]Trade, error) {
	var lTrades []Trade
	lErr := pRepo.db.Preload("Share").Where("user_id = ?", pUserID).Order("created_at desc").Find(&lTrades).Error
	return lTrades, lErr
}

func (pRepo *tradeRepository) RunInTransaction(pFn func(pTx *gorm.DB) error) error {
	lTx := pRepo.db.Begin()
	defer func() {
		if lRec := recover(); lRec != nil {
			lTx.Rollback()
			panic(lRec) // re-throw panic after rollback
		}
	}()

	if lErr := pFn(lTx); lErr != nil {
		lTx.Rollback()
		return lErr
	}

	return lTx.Commit().Error
}
