package auth

import (
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
}

type authService struct {
	repo user.UserRepository
}

func NewAuthService() AuthService {
	return &authService{repo: user.NewUserRepository()}
}

func (s *authService) Register(req RegisterRequest) (*UserResponse, error) {
	// Trim all string inputs - reject whitespace-only values
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Mobile = strings.TrimSpace(req.Mobile)
	req.PAN = strings.TrimSpace(strings.ToUpper(req.PAN))
	req.Aadhaar = strings.TrimSpace(req.Aadhaar)
	req.Address = strings.TrimSpace(req.Address)
	req.IncomeRange = strings.TrimSpace(req.IncomeRange)

	if req.Name == "" {
		return nil, errors.New("Name is required")
	}
	if req.Email == "" {
		return nil, errors.New("Email is required")
	}
	if req.Mobile == "" {
		return nil, errors.New("Mobile is required")
	}
	if req.PAN == "" {
		return nil, errors.New("PAN is required")
	}
	if req.Aadhaar == "" {
		return nil, errors.New("Aadhaar is required")
	}
	if req.Address == "" {
		return nil, errors.New("Address is required")
	}
	if req.IncomeRange == "" {
		return nil, errors.New("Income range is required")
	}

	if len(req.Name) < 3 {
		return nil, errors.New("Name is invalid (minimum 3 letters should be there)")
	}

	if len(req.Mobile) != 10 {
		return nil, errors.New("Phone number must be exactly 10 digits")
	}
	if !regexp.MustCompile(`^[6-9]`).MatchString(req.Mobile) {
		return nil, errors.New("Phone number must start with a digit between 6 to 9")
	}

	// Validate Password complexity
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(req.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(req.Password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(req.Password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(req.Password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return nil, errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	if !regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]{1}$`).MatchString(req.PAN) {
		return nil, errors.New("PAN Number should be like ABCDE1234F")
	}

	if !regexp.MustCompile(`^\d{12}$`).MatchString(req.Aadhaar) {
		return nil, errors.New("Aadhaar Number must be exactly 12 digits")
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(req.Name) {
		return nil, errors.New("invalid Name format. Only letters and spaces are allowed")
	}

	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return nil, errors.New("invalid date of birth format, expected YYYY-MM-DD")
	}

	// Calculate age
	age := time.Now().Year() - dob.Year()
	if time.Now().YearDay() < dob.YearDay() {
		age--
	}
	if age < 18 {
		return nil, errors.New("you must be at least 18 years old to register")
	}

	// Check existing user
	_, err = s.repo.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already in use")
	}

	_, err = s.repo.GetUserByMobile(req.Mobile)
	if err == nil {
		return nil, errors.New("mobile number already in use")
	}

	_, err = s.repo.GetUserByPAN(req.PAN)
	if err == nil {
		return nil, errors.New("PAN already in use")
	}

	_, err = s.repo.GetUserByAadhaar(req.Aadhaar)
	if err == nil {
		return nil, errors.New("Aadhaar already in use")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	user := &user.User{
		Name:         req.Name,
		Mobile:       req.Mobile,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		PAN:          req.PAN,
		Aadhaar:      req.Aadhaar,
		Address:      req.Address,
		IncomeRange:  req.IncomeRange,
		DOB:          dob,
		Role:         "user",
		Status:       "pending",
	}

	var bankAccounts []userpkg.BankDetails
	var nominees []userpkg.NomineeDetails
	var totalPercentage float64 = 0

	err = s.repo.RunInTransaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		ifscRegex := regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)
		// Insert bank accounts
		for _, b := range req.BankAccounts {
			if !ifscRegex.MatchString(b.IFSC) {
				return errors.New("invalid IFSC format")
			}
			bankAccounts = append(bankAccounts, userpkg.BankDetails{
				UserID:        user.ID,
				AccountType:   b.AccountType,
				IFSC:          b.IFSC,
				BankName:      b.BankName,
				Branch:        b.Branch,
				AccountNumber: b.AccountNumber,
				IncomeRange:   req.IncomeRange,
			})
		}
		for i := range bankAccounts {
			if err := tx.Create(&bankAccounts[i]).Error; err != nil {
				return err
			}
		}

		// Insert nominees
		for _, n := range req.Nominees {
			if n.GuardianName != "" {
				if n.GuardianDOB != "" {
					gDob, err := time.Parse("2006-01-02", n.GuardianDOB)
					if err != nil || time.Since(gDob).Hours() < 18*365*24 {
						return errors.New("guardian must be at least 18 years old")
					}
				} else {
					return errors.New("guardian DOB is required")
				}
			}

			nominees = append(nominees, userpkg.NomineeDetails{
				UserID:               user.ID,
				Name:                 n.Name,
				DOB:                  n.DOB,
				PAN:                  n.PAN,
				Relationship:         n.Relationship,
				Percentage:           n.Percentage,
				GuardianName:         n.GuardianName,
				GuardianRelationship: n.GuardianRelationship,
				GuardianPAN:          n.GuardianPAN,
				GuardianDOB:          n.GuardianDOB,
			})
			totalPercentage += n.Percentage
		}

		if len(nominees) > 0 && totalPercentage != 100 {
			return errors.New("total nominee percentage allocation must equal exactly 100")
		}

		for i := range nominees {
			if err := tx.Create(&nominees[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		zap.L().Error("Failed to create user and details", zap.Error(err))
		return nil, errors.New("failed to register user: " + err.Error())
	}

	return &UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		Mobile:      user.Mobile,
		Role:        user.Role,
		Status:      user.Status,
		DOB:         user.DOB.Format("2006-01-02"),
		Address:     user.Address,
		PAN:         user.PAN,
		Aadhaar:     user.Aadhaar,
		IncomeRange: user.IncomeRange,
	}, nil
}

func (s *authService) Login(req LoginRequest, ip string, userAgent string) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if user.Status == "blocked" {
		return nil, errors.New("account blocked")
	}

	access, refresh, err := utils.GenerateTokens(user.ID, user.Role)
	if err != nil {
		zap.L().Error("Token generation failed", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	session := &userpkg.Session{
		UserID:       user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Add(time.Duration(config.App.JWT.ExpirationHours) * time.Hour * 24 * 7),
		IPAddress:    ip,
		UserAgent:    userAgent,
	}

	if err := s.repo.CreateSession(session); err != nil {
		zap.L().Error("Failed to save session", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	return &LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User: &UserResponse{
			ID:          user.ID.String(),
			Name:        user.Name,
			Email:       user.Email,
			Mobile:      user.Mobile,
			Role:        user.Role,
			Status:      user.Status,
			DOB:         user.DOB.Format("2006-01-02"),
			Address:     user.Address,
			PAN:         user.PAN,
			Aadhaar:     user.Aadhaar,
			IncomeRange: user.IncomeRange,
		},
	}, nil
}

func (s *authService) GetMe(userID string) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	banks, _ := s.repo.GetUserBanks(userID)

	var bankDTOs []profile.BankAccountDTO
	for _, b := range banks {
		bankDTOs = append(bankDTOs, profile.BankAccountDTO{
			AccountType:   b.AccountType,
			IFSC:          b.IFSC,
			BankName:      b.BankName,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
		})
	}

	nominees, _ := s.repo.GetUserNominees(userID)
	var nomineeDTOs []profile.NomineeDTO
	for _, n := range nominees {
		nomineeDTOs = append(nomineeDTOs, profile.NomineeDTO{
			Name:                 n.Name,
			DOB:                  n.DOB,
			PAN:                  n.PAN,
			Relationship:         n.Relationship,
			Percentage:           n.Percentage,
			GuardianName:         n.GuardianName,
			GuardianRelationship: n.GuardianRelationship,
			GuardianPAN:          n.GuardianPAN,
			GuardianDOB:          n.GuardianDOB,
		})
	}

	pd, _ := s.repo.GetUserPersonalDetails(userID)
	if pd == nil {
		pd = &userpkg.PersonalDetails{}
	}

	return &UserResponse{
		ID:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		Mobile:       user.Mobile,
		Role:         user.Role,
		Status:       user.Status,
		DOB:          user.DOB.Format("2006-01-02"),
		Address:      pd.Address,
		PAN:          user.PAN,
		Aadhaar:      user.Aadhaar,
		IncomeRange:  user.IncomeRange,
		FatherName:   pd.FatherName,
		MotherName:   pd.MotherName,
		Country:      pd.Country,
		State:        pd.State,
		City:         pd.City,
		Pincode:      pd.Pincode,
		BankAccounts: bankDTOs,
		Nominees:     nomineeDTOs,
	}, nil
}

