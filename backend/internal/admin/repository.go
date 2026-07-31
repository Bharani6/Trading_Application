package admin

import (
	"stock-trading/internal/database"
	"stock-trading/internal/trade"
	userpkg "stock-trading/internal/user"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetAllUsers() ([]userpkg.User, error)
	GetUserDetails(userID string) (*userpkg.User, *userpkg.PersonalDetails, []userpkg.BankDetails, []userpkg.NomineeDetails, error)
	UpdateUserStatus(userID string, status string) error
	GetOrCreateSegment(name string) (*trade.Segment, error)
	BulkInsertShares(shares []trade.Share) error
	DeleteAllShares() error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository() AdminRepository {
	return &adminRepository{db: database.DB}
}

func (r *adminRepository) GetAllUsers() ([]userpkg.User, error) {
	var users []userpkg.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *adminRepository) GetUserDetails(userID string) (*userpkg.User, *userpkg.PersonalDetails, []userpkg.BankDetails, []userpkg.NomineeDetails, error) {
	var user userpkg.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	var personal userpkg.PersonalDetails
	r.db.Where("user_id = ?", userID).First(&personal) // ignore error as they might not exist yet

	var banks []userpkg.BankDetails
	r.db.Where("user_id = ?", userID).Find(&banks)

	var nominees []userpkg.NomineeDetails
	r.db.Where("user_id = ?", userID).Find(&nominees)

	return &user, &personal, banks, nominees, nil
}

func (r *adminRepository) UpdateUserStatus(userID string, status string) error {
	return r.db.Model(&userpkg.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (r *adminRepository) GetOrCreateSegment(name string) (*trade.Segment, error) {
	var segment trade.Segment
	err := r.db.Where("name = ?", name).FirstOrCreate(&segment, trade.Segment{Name: name}).Error
	return &segment, err
}

func (r *adminRepository) BulkInsertShares(shares []trade.Share) error {
	return r.db.CreateInBatches(shares, 100).Error
}

func (r *adminRepository) DeleteAllShares() error {
	return r.db.Unscoped().Where("1 = 1").Delete(&trade.Share{}).Error
}
