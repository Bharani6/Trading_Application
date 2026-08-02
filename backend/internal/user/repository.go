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
	GetSessionsByUserID(userID string) ([]Session, error)
	DeleteSessionByID(sessionID string) error
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

func (pRepo *userRepository) CreateUser(pUser *User) error {
	return pRepo.db.Create(pUser).Error
}

func (pRepo *userRepository) GetUserByEmail(pEmail string) (*User, error) {
	var lUser User
	if lErr := pRepo.db.Where("email = ?", pEmail).First(&lUser).Error; lErr != nil {
		return nil, lErr
	}
	return &lUser, nil
}

func (pRepo *userRepository) GetUserByMobile(pMobile string) (*User, error) {
	var lUser User
	if lErr := pRepo.db.Where("mobile = ?", pMobile).First(&lUser).Error; lErr != nil {
		return nil, lErr
	}
	return &lUser, nil
}

func (pRepo *userRepository) GetUserByPAN(pPan string) (*User, error) {
	var lUser User
	if lErr := pRepo.db.Where("pan = ?", pPan).First(&lUser).Error; lErr != nil {
		return nil, lErr
	}
	return &lUser, nil
}

func (pRepo *userRepository) GetUserByAadhaar(pAadhaar string) (*User, error) {
	var lUser User
	if lErr := pRepo.db.Where("aadhaar = ?", pAadhaar).First(&lUser).Error; lErr != nil {
		return nil, lErr
	}
	return &lUser, nil
}

func (pRepo *userRepository) GetUserByID(pID string) (*User, error) {
	var lUser User
	lErr := pRepo.db.Where("id = ?", pID).First(&lUser).Error
	return &lUser, lErr
}

func (pRepo *userRepository) GetUserBanks(pUserID string) ([]BankDetails, error) {
	var lBanks []BankDetails
	lErr := pRepo.db.Where("user_id = ?", pUserID).Find(&lBanks).Error
	return lBanks, lErr
}

func (pRepo *userRepository) GetUserNominees(pUserID string) ([]NomineeDetails, error) {
	var lNominees []NomineeDetails
	lErr := pRepo.db.Where("user_id = ?", pUserID).Find(&lNominees).Error
	return lNominees, lErr
}

func (pRepo *userRepository) GetUserPersonalDetails(pUserID string) (*PersonalDetails, error) {
	var lPd PersonalDetails
	lErr := pRepo.db.Where("user_id = ?", pUserID).First(&lPd).Error
	return &lPd, lErr
}

func (pRepo *userRepository) CreateSession(pSession *Session) error {
	return pRepo.db.Create(pSession).Error
}

func (pRepo *userRepository) GetSessionByRefreshToken(pToken string) (*Session, error) {
	var lSession Session
	lErr := pRepo.db.Where("refresh_token = ?", pToken).First(&lSession).Error
	return &lSession, lErr
}

func (pRepo *userRepository) DeleteSession(pToken string) error {
	return pRepo.db.Where("refresh_token = ?", pToken).Delete(&Session{}).Error
}

func (pRepo *userRepository) DeleteAllSessions(pUserID string) error {
	return pRepo.db.Where("user_id = ?", pUserID).Delete(&Session{}).Error
}

func (pRepo *userRepository) GetSessionsByUserID(pUserID string) ([]Session, error) {
	var lSessions []Session
	lErr := pRepo.db.Where("user_id = ?", pUserID).Find(&lSessions).Error
	return lSessions, lErr
}

func (pRepo *userRepository) DeleteSessionByID(pSessionID string) error {
	return pRepo.db.Where("id = ?", pSessionID).Delete(&Session{}).Error
}

func (pRepo *userRepository) RunInTransaction(pFn func(pTx *gorm.DB) error) error {
	return pRepo.db.Transaction(pFn)
}

func (pRepo *userRepository) CreatePasswordResetToken(pToken *PasswordResetToken) error {
	return pRepo.db.Create(pToken).Error
}

func (pRepo *userRepository) GetPasswordResetToken(pTokenStr string) (*PasswordResetToken, error) {
	var lToken PasswordResetToken
	if lErr := pRepo.db.Where("token = ?", pTokenStr).First(&lToken).Error; lErr != nil {
		return nil, lErr
	}
	return &lToken, nil
}

func (pRepo *userRepository) DeletePasswordResetToken(pTokenStr string) error {
	return pRepo.db.Where("token = ?", pTokenStr).Delete(&PasswordResetToken{}).Error
}

func (pRepo *userRepository) UpdatePassword(pUserID string, pHashedPassword string) error {
	return pRepo.db.Model(&User{}).Where("id = ?", pUserID).Update("password_hash", pHashedPassword).Error
}
