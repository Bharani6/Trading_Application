package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"stock-trading/internal/config"
	"stock-trading/internal/user"

	"stock-trading/internal/profile"
	userpkg "stock-trading/internal/user"
	"stock-trading/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req RegisterRequest) (*UserResponse, error)
	Login(req LoginRequest, ip string, userAgent string) (*LoginResponse, error)
	GetMe(userID string) (*UserResponse, error)
	ForgotPassword(req ForgotPasswordRequest) (string, error)
	VerifyResetToken(req VerifyResetTokenRequest) error
	ResetPassword(req ResetPasswordRequest) error
	GetActiveSessions(userID string, currentToken string) ([]SessionDTO, error)
	RevokeSession(sessionID string) error
}

type authService struct {
	repo user.UserRepository
}

func NewAuthService() AuthService {
	return &authService{repo: user.NewUserRepository()}
}

func (pService *authService) Register(pReq RegisterRequest) (*UserResponse, error) {
	// Trim all string inputs - reject whitespace-only values
	pReq.Name = strings.TrimSpace(pReq.Name)
	pReq.Email = strings.TrimSpace(strings.ToLower(pReq.Email))
	pReq.Mobile = strings.TrimSpace(pReq.Mobile)
	pReq.PAN = strings.TrimSpace(strings.ToUpper(pReq.PAN))
	pReq.Aadhaar = strings.TrimSpace(pReq.Aadhaar)
	pReq.Address = strings.TrimSpace(pReq.Address)
	pReq.IncomeRange = strings.TrimSpace(pReq.IncomeRange)
	pReq.Occupation = strings.TrimSpace(pReq.Occupation)

	if pReq.Name == "" {
		return nil, errors.New("Name is required")
	}
	if pReq.Email == "" {
		return nil, errors.New("Email is required")
	}
	if pReq.Mobile == "" {
		return nil, errors.New("Mobile is required")
	}
	if pReq.PAN == "" {
		return nil, errors.New("PAN is required")
	}
	if pReq.Aadhaar == "" {
		return nil, errors.New("Aadhaar is required")
	}
	if pReq.Address == "" {
		return nil, errors.New("Address is required")
	}
	if pReq.IncomeRange == "" {
		return nil, errors.New("Income range is required")
	}
	if pReq.Occupation == "" {
		return nil, errors.New("Occupation is required")
	}
	if pReq.IPVPhoto == "" {
		return nil, errors.New("IPV Photo is required")
	}
	if pReq.IPVLatitude == "" || pReq.IPVLongitude == "" {
		return nil, errors.New("IPV Location coordinates are required")
	}

	if len(pReq.Name) < 3 {
		return nil, errors.New("Name is invalid (minimum 3 letters should be there)")
	}

	if len(pReq.Mobile) != 10 {
		return nil, errors.New("Phone number must be exactly 10 digits")
	}
	if !regexp.MustCompile(`^[6-9]`).MatchString(pReq.Mobile) {
		return nil, errors.New("Phone number must start with a digit between 6 to 9")
	}

	// Validate Password complexity
	if len(pReq.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	lHasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pReq.Password)
	lHasLower := regexp.MustCompile(`[a-z]`).MatchString(pReq.Password)
	lHasDigit := regexp.MustCompile(`[0-9]`).MatchString(pReq.Password)
	lHasSpecial := regexp.MustCompile(`[\W_]`).MatchString(pReq.Password)

	if !lHasUpper || !lHasLower || !lHasDigit || !lHasSpecial {
		return nil, errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	if !regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`).MatchString(pReq.PAN) {
		return nil, errors.New("PAN Number should be like ABCDE1234F")
	}

	if !regexp.MustCompile(`^\d{12}$`).MatchString(pReq.Aadhaar) {
		return nil, errors.New("Aadhaar Number must be exactly 12 digits")
	}

	lNameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !lNameRegex.MatchString(pReq.Name) {
		return nil, errors.New("invalid Name format. Only letters and spaces are allowed")
	}

	lDOB, lErr := time.Parse("2006-01-02", pReq.DOB)
	if lErr != nil {
		return nil, errors.New("invalid date of birth format, expected YYYY-MM-DD")
	}

	// Calculate age
	lAge := time.Now().Year() - lDOB.Year()
	if time.Now().YearDay() < lDOB.YearDay() {
		lAge--
	}
	if lAge < 18 {
		return nil, errors.New("you must be at least 18 years old to register")
	}

	// Check existing user
	_, lErr = pService.repo.GetUserByEmail(pReq.Email)
	if lErr == nil {
		return nil, errors.New("email already in use")
	}

	_, lErr = pService.repo.GetUserByMobile(pReq.Mobile)
	if lErr == nil {
		return nil, errors.New("mobile number already in use")
	}

	_, lErr = pService.repo.GetUserByPAN(pReq.PAN)
	if lErr == nil {
		return nil, errors.New("PAN already in use")
	}

	_, lErr = pService.repo.GetUserByAadhaar(pReq.Aadhaar)
	if lErr == nil {
		return nil, errors.New("Aadhaar already in use")
	}

	lHashedPassword, lErr := utils.HashPassword(pReq.Password)
	if lErr != nil {
		zap.L().Error("Failed to hash password", zap.Error(lErr))
		return nil, errors.New("internal server error")
	}

	lUser := &user.User{
		Name:         pReq.Name,
		Mobile:       pReq.Mobile,
		Email:        pReq.Email,
		PasswordHash: lHashedPassword,
		PAN:          pReq.PAN,
		Aadhaar:      pReq.Aadhaar,
		Address:      pReq.Address,
		IncomeRange:  pReq.IncomeRange,
		Occupation:   pReq.Occupation,
		IPVPhoto:     pReq.IPVPhoto,
		IPVLatitude:  pReq.IPVLatitude,
		IPVLongitude: pReq.IPVLongitude,
		DOB:          lDOB,
		Role:         "user",
		Status:       "pending",
	}

	var lBankAccounts []userpkg.BankDetails
	var lNominees []userpkg.NomineeDetails
	var lTotalPercentage float64 = 0

	lErr = pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		if lTxErr := pTx.Create(lUser).Error; lTxErr != nil {
			return lTxErr
		}

		lIfscRegex := regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)
		// Insert bank accounts
		for _, lBank := range pReq.BankAccounts {
			if !lIfscRegex.MatchString(lBank.IFSC) {
				return errors.New("invalid IFSC format")
			}
			lBankAccounts = append(lBankAccounts, userpkg.BankDetails{
				UserID:        lUser.ID,
				AccountType:   lBank.AccountType,
				IFSC:          lBank.IFSC,
				BankName:      lBank.BankName,
				Branch:        lBank.Branch,
				AccountNumber: lBank.AccountNumber,
				IncomeRange:   pReq.IncomeRange,
			})
		}
		for lIdx := range lBankAccounts {
			if lTxErr := pTx.Create(&lBankAccounts[lIdx]).Error; lTxErr != nil {
				return lTxErr
			}
		}

		// Insert nominees
		for _, lNominee := range pReq.Nominees {
			if lNominee.GuardianName != "" {
				if lNominee.GuardianDOB != "" {
					lGDob, lGDobErr := time.Parse("2006-01-02", lNominee.GuardianDOB)
					if lGDobErr != nil || time.Since(lGDob).Hours() < 18*365*24 {
						return errors.New("guardian must be at least 18 years old")
					}
				} else {
					return errors.New("guardian DOB is required")
				}
			}

			lNominees = append(lNominees, userpkg.NomineeDetails{
				UserID:               lUser.ID,
				Name:                 lNominee.Name,
				DOB:                  lNominee.DOB,
				PAN:                  lNominee.PAN,
				Relationship:         lNominee.Relationship,
				Percentage:           lNominee.Percentage,
				GuardianName:         lNominee.GuardianName,
				GuardianRelationship: lNominee.GuardianRelationship,
				GuardianPAN:          lNominee.GuardianPAN,
				GuardianDOB:          lNominee.GuardianDOB,
			})
			lTotalPercentage += lNominee.Percentage
		}

		if len(lNominees) > 0 && lTotalPercentage != 100 {
			return errors.New("total nominee percentage allocation must equal exactly 100")
		}

		for lIdx := range lNominees {
			if lTxErr := pTx.Create(&lNominees[lIdx]).Error; lTxErr != nil {
				return lTxErr
			}
		}

		return nil
	})

	if lErr != nil {
		zap.L().Error("Failed to create user and details", zap.Error(lErr))
		return nil, errors.New("failed to register user: " + lErr.Error())
	}

	return pService.GetMe(lUser.ID.String())
}

func (pService *authService) Login(pReq LoginRequest, pIP string, pUserAgent string) (*LoginResponse, error) {
	lUser, lErr := pService.repo.GetUserByEmail(pReq.Email)
	if lErr != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(pReq.Password, lUser.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if lUser.Status == "closed" {
		return nil, errors.New("invalid credentials")
	}

	if lUser.Status == "blocked" {
		return nil, errors.New("account blocked")
	}

	lAccess, lRefresh, lErr := utils.GenerateTokens(lUser.ID, lUser.Role)
	if lErr != nil {
		zap.L().Error("Token generation failed", zap.Error(lErr))
		return nil, errors.New("internal server error")
	}

	lSession := &userpkg.Session{
		UserID:       lUser.ID,
		AccessToken:  lAccess,
		RefreshToken: lRefresh,
		ExpiresAt:    time.Now().Add(time.Duration(config.App.JWT.ExpirationHours) * time.Hour * 24 * 7),
		IPAddress:    pIP,
		UserAgent:    pUserAgent,
	}

	if lErr := pService.repo.CreateSession(lSession); lErr != nil {
		zap.L().Error("Failed to save session", zap.Error(lErr))
		return nil, errors.New("internal server error")
	}

	lUserResp, lErr := pService.GetMe(lUser.ID.String())
	if lErr != nil {
		return nil, errors.New("failed to fetch user details")
	}

	return &LoginResponse{
		AccessToken:  lAccess,
		RefreshToken: lRefresh,
		User:         lUserResp,
	}, nil
}

func (pService *authService) GetMe(pUserID string) (*UserResponse, error) {
	lUser, lErr := pService.repo.GetUserByID(pUserID)
	if lErr != nil {
		return nil, errors.New("user not found")
	}

	lBanks, _ := pService.repo.GetUserBanks(pUserID)

	var lBankDTOs []profile.BankAccountDTO
	for _, lBank := range lBanks {
		lBankDTOs = append(lBankDTOs, profile.BankAccountDTO{
			AccountType:   lBank.AccountType,
			IFSC:          lBank.IFSC,
			BankName:      lBank.BankName,
			Branch:        lBank.Branch,
			AccountNumber: lBank.AccountNumber,
		})
	}

	lNominees, _ := pService.repo.GetUserNominees(pUserID)
	var lNomineeDTOs []profile.NomineeDTO
	for _, lNominee := range lNominees {
		lNomineeDTOs = append(lNomineeDTOs, profile.NomineeDTO{
			Name:                 lNominee.Name,
			DOB:                  lNominee.DOB,
			PAN:                  lNominee.PAN,
			Relationship:         lNominee.Relationship,
			Percentage:           lNominee.Percentage,
			GuardianName:         lNominee.GuardianName,
			GuardianRelationship: lNominee.GuardianRelationship,
			GuardianPAN:          lNominee.GuardianPAN,
			GuardianDOB:          lNominee.GuardianDOB,
		})
	}

	lPD, _ := pService.repo.GetUserPersonalDetails(pUserID)
	if lPD == nil {
		lPD = &userpkg.PersonalDetails{}
	}

	return &UserResponse{
		ID:           lUser.ID.String(),
		Name:         lUser.Name,
		Email:        lUser.Email,
		Mobile:       lUser.Mobile,
		Role:         lUser.Role,
		Status:       lUser.Status,
		DOB:          lUser.DOB.Format("2006-01-02"),
		Address:      lPD.Address,
		PAN:          lUser.PAN,
		Aadhaar:      lUser.Aadhaar,
		IncomeRange:  lUser.IncomeRange,
		Occupation:   lUser.Occupation,
		IPVPhoto:     lUser.IPVPhoto,
		IPVLatitude:  lUser.IPVLatitude,
		IPVLongitude: lUser.IPVLongitude,
		FatherName:   lPD.FatherName,
		MotherName:   lPD.MotherName,
		Country:      lPD.Country,
		State:        lPD.State,
		City:         lPD.City,
		Pincode:      lPD.Pincode,
		BankAccounts: lBankDTOs,
		Nominees:     lNomineeDTOs,
	}, nil
}

func generateCryptoToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (pService *authService) ForgotPassword(pReq ForgotPasswordRequest) (string, error) {
	pReq.Email = strings.TrimSpace(strings.ToLower(pReq.Email))
	if pReq.Email == "" {
		return "", errors.New("Email is required")
	}

	lUser, lErr := pService.repo.GetUserByEmail(pReq.Email)
	if lErr != nil {
		// As per security best practices, do not reveal if email exists.
		return "", nil 
	}

	lTokenStr, lErr := generateCryptoToken()
	if lErr != nil {
		return "", errors.New("failed to generate reset token")
	}

	lResetToken := &userpkg.PasswordResetToken{
		UserID:    lUser.ID,
		Token:     lTokenStr,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if lErr := pService.repo.CreatePasswordResetToken(lResetToken); lErr != nil {
		return "", errors.New("failed to save reset token")
	}

	// Returning the token here to mock the email sending process
	return lTokenStr, nil
}

func (pService *authService) VerifyResetToken(pReq VerifyResetTokenRequest) error {
	lTokenStr := strings.TrimSpace(pReq.Token)
	if lTokenStr == "" {
		return errors.New("Token is required")
	}

	lToken, lErr := pService.repo.GetPasswordResetToken(lTokenStr)
	if lErr != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(lToken.ExpiresAt) {
		pService.repo.DeletePasswordResetToken(lTokenStr)
		return errors.New("invalid or expired token")
	}

	return nil
}

func (pService *authService) ResetPassword(pReq ResetPasswordRequest) error {
	if lErr := pService.VerifyResetToken(VerifyResetTokenRequest{Token: pReq.Token}); lErr != nil {
		return lErr
	}

	lToken, _ := pService.repo.GetPasswordResetToken(pReq.Token)

	// Validate Password complexity
	if len(pReq.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	lHasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pReq.Password)
	lHasLower := regexp.MustCompile(`[a-z]`).MatchString(pReq.Password)
	lHasDigit := regexp.MustCompile(`[0-9]`).MatchString(pReq.Password)
	lHasSpecial := regexp.MustCompile(`[\W_]`).MatchString(pReq.Password)

	if !lHasUpper || !lHasLower || !lHasDigit || !lHasSpecial {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	lHashedPassword, lErr := utils.HashPassword(pReq.Password)
	if lErr != nil {
		return errors.New("internal server error")
	}

	if lErr := pService.repo.UpdatePassword(lToken.UserID.String(), lHashedPassword); lErr != nil {
		return errors.New("failed to reset password")
	}

	// Invalidate the token
	pService.repo.DeletePasswordResetToken(pReq.Token)
    
	// invalidate all existing sessions
	pService.repo.DeleteAllSessions(lToken.UserID.String())

	return nil
}

func (pService *authService) GetActiveSessions(pUserID string, pCurrentToken string) ([]SessionDTO, error) {
	lSessions, lErr := pService.repo.GetSessionsByUserID(pUserID)
	if lErr != nil {
		return nil, lErr
	}

	var lDTOs []SessionDTO
	for _, lSession := range lSessions {
		lDTOs = append(lDTOs, SessionDTO{
			ID:        lSession.ID.String(),
			IPAddress: lSession.IPAddress,
			UserAgent: lSession.UserAgent,
			CreatedAt: lSession.CreatedAt.Format("2006-01-02 15:04:05"),
			IsCurrent: lSession.AccessToken == pCurrentToken,
		})
	}
	return lDTOs, nil
}

func (pService *authService) RevokeSession(pSessionID string) error {
	return pService.repo.DeleteSessionByID(pSessionID)
}
