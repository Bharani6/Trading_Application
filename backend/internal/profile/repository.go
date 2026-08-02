package profile

import (
	"stock-trading/internal/database"
	user "stock-trading/internal/user" 

	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetPersonalDetails(userID string) (*user.PersonalDetails, error)
	SavePersonalDetails(details *user.PersonalDetails) error
	SaveBankDetails(details *user.BankDetails) error
	UpdateUserStatus(userID string, status string) error
	RunInTransaction(fn func(tx *gorm.DB) error) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{db: database.DB}
}

func (pRepo *profileRepository) GetPersonalDetails(pUserID string) (*user.PersonalDetails, error) {
	var lDetails user.PersonalDetails
	lErr := pRepo.db.Where("user_id = ?", pUserID).First(&lDetails).Error
	return &lDetails, lErr
}

func (pRepo *profileRepository) SavePersonalDetails(pDetails *user.PersonalDetails) error {
	return pRepo.db.Save(pDetails).Error
}

func (pRepo *profileRepository) SaveBankDetails(pDetails *user.BankDetails) error {
	return pRepo.db.Save(pDetails).Error
}

func (pRepo *profileRepository) UpdateUserStatus(pUserID string, pStatus string) error {
	return pRepo.db.Model(&user.User{}).Where("id = ?", pUserID).Update("status", pStatus).Error
}

func (pRepo *profileRepository) RunInTransaction(pFn func(pTx *gorm.DB) error) error {
	return pRepo.db.Transaction(pFn)
}


