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

func (r *profileRepository) GetPersonalDetails(userID string) (*user.PersonalDetails, error) {
	var details user.PersonalDetails
	err := r.db.Where("user_id = ?", userID).First(&details).Error
	return &details, err
}

func (r *profileRepository) SavePersonalDetails(details *user.PersonalDetails) error {
	return r.db.Save(details).Error
}

func (r *profileRepository) SaveBankDetails(details *user.BankDetails) error {
	return r.db.Save(details).Error
}

func (r *profileRepository) UpdateUserStatus(userID string, status string) error {
	return r.db.Model(&user.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *profileRepository) RunInTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}


