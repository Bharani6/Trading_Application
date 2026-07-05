package user

import (
	"stock-trading/internal/database"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByMobile(mobile string) (*User, error)
	GetUserByPAN(pan string) (*User, error)
	GetUserByAadhaar(aadhaar string) (*User, error)
	GetUserByID(id string) (*User, error)
	CreateSession(session *Session) error
	GetSessionByRefreshToken(token string) (*Session, error)
	DeleteSession(token string) error
	DeleteAllSessions(userID string) error
	GetUserBanks(userID string) ([]BankDetails, error)
	GetUserNominees(userID string) ([]NomineeDetails, error)
	GetUserPersonalDetails(userID string) (*PersonalDetails, error)
	RunInTransaction(fn func(tx *gorm.DB) error) error
	CreatePasswordResetToken(token *PasswordResetToken) error
	GetPasswordResetToken(tokenStr string) (*PasswordResetToken, error)
	DeletePasswordResetToken(tokenStr string) error
	UpdatePassword(userID string, hashedPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: database.DB}
}

func (r *userRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByMobile(mobile string) (*User, error) {
	var user User
	if err := r.db.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByPAN(pan string) (*User, error) {
	var user User
	if err := r.db.Where("pan = ?", pan).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByAadhaar(aadhaar string) (*User, error) {
	var user User
	if err := r.db.Where("aadhaar = ?", aadhaar).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserBanks(userID string) ([]BankDetails, error) {
	var banks []BankDetails
	err := r.db.Where("user_id = ?", userID).Find(&banks).Error
	return banks, err
}

func (r *userRepository) GetUserNominees(userID string) ([]NomineeDetails, error) {
	var nominees []NomineeDetails
	err := r.db.Where("user_id = ?", userID).Find(&nominees).Error
	return nominees, err
}

func (r *userRepository) GetUserPersonalDetails(userID string) (*PersonalDetails, error) {
	var pd PersonalDetails
	err := r.db.Where("user_id = ?", userID).First(&pd).Error
	return &pd, err
}

func (r *userRepository) CreateSession(session *Session) error {
	return r.db.Create(session).Error
}

func (r *userRepository) GetSessionByRefreshToken(token string) (*Session, error) {
	var session Session
	err := r.db.Where("refresh_token = ?", token).First(&session).Error
	return &session, err
}

func (r *userRepository) DeleteSession(token string) error {
	return r.db.Where("refresh_token = ?", token).Delete(&Session{}).Error
}

func (r *userRepository) DeleteAllSessions(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&Session{}).Error
}

func (r *userRepository) RunInTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *userRepository) CreatePasswordResetToken(token *PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *userRepository) GetPasswordResetToken(tokenStr string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	if err := r.db.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *userRepository) DeletePasswordResetToken(tokenStr string) error {
	return r.db.Where("token = ?", tokenStr).Delete(&PasswordResetToken{}).Error
}

func (r *userRepository) UpdatePassword(userID string, hashedPassword string) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("password_hash", hashedPassword).Error
}
