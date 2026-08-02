package admin

import (
	"stock-trading/internal/database"
	"stock-trading/internal/trade"
	userpkg "stock-trading/internal/user"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetAllUsers() ([]userpkg.User, error)
	GetUserDetails(pUserID string) (*userpkg.User, *userpkg.PersonalDetails, []userpkg.BankDetails, []userpkg.NomineeDetails, error)
	UpdateUserStatus(pUserID string, pStatus string) error
	GetOrCreateSegment(pName string) (*trade.Segment, error)
	BulkInsertShares(pShares []trade.Share) error
	DeleteAllShares() error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository() AdminRepository {
	return &adminRepository{db: database.DB}
}

func (pRepo *adminRepository) GetAllUsers() ([]userpkg.User, error) {
	var lUsers []userpkg.User
	lErr := pRepo.db.Find(&lUsers).Error
	return lUsers, lErr
}

func (pRepo *adminRepository) GetUserDetails(pUserID string) (*userpkg.User, *userpkg.PersonalDetails, []userpkg.BankDetails, []userpkg.NomineeDetails, error) {
	var lUser userpkg.User
	if lErr := pRepo.db.Where("id = ?", pUserID).First(&lUser).Error; lErr != nil {
		return nil, nil, nil, nil, lErr
	}

	var lPersonal userpkg.PersonalDetails
	pRepo.db.Where("user_id = ?", pUserID).First(&lPersonal) // ignore error as they might not exist yet

	var lBanks []userpkg.BankDetails
	pRepo.db.Where("user_id = ?", pUserID).Find(&lBanks)

	var lNominees []userpkg.NomineeDetails
	pRepo.db.Where("user_id = ?", pUserID).Find(&lNominees)

	return &lUser, &lPersonal, lBanks, lNominees, nil
}

func (pRepo *adminRepository) UpdateUserStatus(pUserID string, pStatus string) error {
	return pRepo.db.Model(&userpkg.User{}).Where("id = ?", pUserID).Update("status", pStatus).Error
}

func (pRepo *adminRepository) GetOrCreateSegment(pName string) (*trade.Segment, error) {
	var lSegment trade.Segment
	lErr := pRepo.db.Where("name = ?", pName).FirstOrCreate(&lSegment, trade.Segment{Name: pName}).Error
	return &lSegment, lErr
}

func (pRepo *adminRepository) BulkInsertShares(pShares []trade.Share) error {
	return pRepo.db.CreateInBatches(pShares, 100).Error
}

func (pRepo *adminRepository) DeleteAllShares() error {
	return pRepo.db.Unscoped().Where("1 = 1").Delete(&trade.Share{}).Error
}
