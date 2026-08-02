package profile

import (
	"errors"
	"regexp"
	"stock-trading/internal/user"
	"stock-trading/internal/utils"

	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfileService interface {
	SubmitKYC(userID string, req KYCSubmitRequest) error
	RequestClosure(userID string) error
	ChangePassword(userID string, req ChangePasswordRequest) error
}

type profileService struct {
	repo ProfileRepository
}

func NewProfileService() ProfileService {
	return &profileService{repo: NewProfileRepository()}
}

func (pService *profileService) SubmitKYC(pUserID string, pReq KYCSubmitRequest) error {
	lUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return errors.New("invalid user id")
	}

	var lDetails user.PersonalDetails
	if lExisting, lErr := pService.repo.GetPersonalDetails(lUID.String()); lErr == nil {
		lDetails = *lExisting
	} else {
		lDetails = user.PersonalDetails{UserID: lUID}
	}

	if pReq.FatherName != "" {
		lDetails.FatherName = pReq.FatherName
	}
	if pReq.MotherName != "" {
		lDetails.MotherName = pReq.MotherName
	}
	if pReq.Country != "" {
		lDetails.Country = pReq.Country
	}
	if pReq.State != "" {
		lDetails.State = pReq.State
	}
	if pReq.City != "" {
		lDetails.City = pReq.City
	}
	if pReq.Address != "" {
		lDetails.Address = pReq.Address
	}
	if pReq.Pincode != "" {
		lDetails.Pincode = pReq.Pincode
	}
	var lBankAccounts []user.BankDetails
	lIfscRegex := regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)

	for _, lB := range pReq.BankAccounts {
		if !lIfscRegex.MatchString(lB.IFSC) {
			return errors.New("invalid IFSC format")
		}
		lBankAccounts = append(lBankAccounts, user.BankDetails{
			UserID:        lUID,
			AccountType:   lB.AccountType,
			IFSC:          lB.IFSC,
			BankName:      lB.BankName,
			Branch:        lB.Branch,
			AccountNumber: lB.AccountNumber,
			IncomeRange:   pReq.IncomeRange,
		})
	}

	var lNominees []user.NomineeDetails
	var lTotalPercentage float64 = 0

	for _, lN := range pReq.Nominees {
		if lN.GuardianName != "" {
			if lN.GuardianDOB != "" {
				lDob, lErr := time.Parse("2006-01-02", lN.GuardianDOB)
				if lErr != nil || time.Since(lDob).Hours() < 18*365*24 {
					return errors.New("guardian must be at least 18 years old")
				}
			} else {
				return errors.New("guardian DOB is required")
			}
		}

		lNominees = append(lNominees, user.NomineeDetails{
			UserID:               lUID,
			Name:                 lN.Name,
			DOB:                  lN.DOB,
			PAN:                  lN.PAN,
			Relationship:         lN.Relationship,
			Percentage:           lN.Percentage,
			GuardianName:         lN.GuardianName,
			GuardianRelationship: lN.GuardianRelationship,
			GuardianPAN:          lN.GuardianPAN,
			GuardianDOB:          lN.GuardianDOB,
		})
		lTotalPercentage += lN.Percentage
	}

	if len(lNominees) > 0 && lTotalPercentage != 100 {
		return errors.New("total nominee percentage allocation must equal exactly 100")
	}

	lErr = pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		if lTxErr := pTx.Save(&lDetails).Error; lTxErr != nil {
			return lTxErr
		}

		// Replace bank details only if they are provided
		if len(lBankAccounts) > 0 {
			if lTxErr := pTx.Where("user_id = ?", lUID).Delete(&user.BankDetails{}).Error; lTxErr != nil {
				return lTxErr
			}
			for i := range lBankAccounts {
				if lTxErr := pTx.Create(&lBankAccounts[i]).Error; lTxErr != nil {
					return lTxErr
				}
			}
		}

		// Replace nominees only if they are provided
		if len(lNominees) > 0 {
			if lTxErr := pTx.Where("user_id = ?", lUID).Delete(&user.NomineeDetails{}).Error; lTxErr != nil {
				return lTxErr
			}
			for i := range lNominees {
				if lTxErr := pTx.Create(&lNominees[i]).Error; lTxErr != nil {
					return lTxErr
				}
			}
		}

		var lExistingUser user.User
		lUpdateData := map[string]interface{}{}
		if lTxErr := pTx.Where("id = ?", pUserID).First(&lExistingUser).Error; lTxErr == nil {
			if lExistingUser.Role != "admin" {
				lUpdateData["status"] = "pending_approval"
			}
		} else {
			lUpdateData["status"] = "pending_approval"
		}
		
		if pReq.Mobile != "" {
			lUpdateData["mobile"] = pReq.Mobile
		}
		if pReq.IncomeRange != "" {
			lUpdateData["income_range"] = pReq.IncomeRange
		}
		if pReq.Occupation != "" {
			lUpdateData["occupation"] = pReq.Occupation
		}
		if lTxErr := pTx.Model(&user.User{}).Where("id = ?", pUserID).Updates(lUpdateData).Error; lTxErr != nil {
			return lTxErr
		}
		return nil
	})

	if lErr != nil {
		zap.L().Error("KYC submission failed", zap.Error(lErr))
		return errors.New("failed to submit KYC details")
	}

	return nil
}

func (pService *profileService) RequestClosure(pUserID string) error {
	return pService.repo.UpdateUserStatus(pUserID, "closure_requested")
}

func (pService *profileService) ChangePassword(pUserID string, pReq ChangePasswordRequest) error {
	lUserRepo := user.NewUserRepository()
	lU, lErr := lUserRepo.GetUserByID(pUserID)
	if lErr != nil {
		return errors.New("invalid user")
	}

	if !utils.CheckPasswordHash(pReq.CurrentPassword, lU.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	if len(pReq.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}
	lHasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pReq.NewPassword)
	lHasLower := regexp.MustCompile(`[a-z]`).MatchString(pReq.NewPassword)
	lHasDigit := regexp.MustCompile(`[0-9]`).MatchString(pReq.NewPassword)
	if !lHasUpper || !lHasLower || !lHasDigit {
		return errors.New("new password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	lHashedPassword, lErr := utils.HashPassword(pReq.NewPassword)
	if lErr != nil {
		return errors.New("internal server error")
	}

	return lUserRepo.UpdatePassword(pUserID, lHashedPassword)
}

