package profile

import (
	"errors"
	"regexp"
	"stock-trading/internal/user"

	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfileService interface {
	SubmitKYC(userID string, req KYCSubmitRequest) error
}

type profileService struct {
	repo ProfileRepository
}

func NewProfileService() ProfileService {
	return &profileService{repo: NewProfileRepository()}
}

func (s *profileService) SubmitKYC(userID string, req KYCSubmitRequest) error {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	var details user.PersonalDetails
	if existing, err := s.repo.GetPersonalDetails(uID.String()); err == nil {
		details = *existing
	} else {
		details = user.PersonalDetails{UserID: uID}
	}

	if req.FatherName != "" {
		details.FatherName = req.FatherName
	}
	if req.MotherName != "" {
		details.MotherName = req.MotherName
	}
	if req.Country != "" {
		details.Country = req.Country
	}
	if req.State != "" {
		details.State = req.State
	}
	if req.City != "" {
		details.City = req.City
	}
	if req.Address != "" {
		details.Address = req.Address
	}
	if req.Pincode != "" {
		details.Pincode = req.Pincode
	}
	var bankAccounts []user.BankDetails
	ifscRegex := regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)

	for _, b := range req.BankAccounts {
		if !ifscRegex.MatchString(b.IFSC) {
			return errors.New("invalid IFSC format")
		}
		bankAccounts = append(bankAccounts, user.BankDetails{
			UserID:        uID,
			AccountType:   b.AccountType,
			IFSC:          b.IFSC,
			BankName:      b.BankName,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
			IncomeRange:   req.IncomeRange,
		})
	}

	var nominees []user.NomineeDetails
	var totalPercentage float64 = 0

	for _, n := range req.Nominees {
		if n.GuardianName != "" {
			if n.GuardianDOB != "" {
				dob, err := time.Parse("2006-01-02", n.GuardianDOB)
				if err != nil || time.Since(dob).Hours() < 18*365*24 {
					return errors.New("guardian must be at least 18 years old")
				}
			} else {
				return errors.New("guardian DOB is required")
			}
		}

		nominees = append(nominees, user.NomineeDetails{
			UserID:               uID,
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

	err = s.repo.RunInTransaction(func(tx *gorm.DB) error {
		if err := tx.Save(&details).Error; err != nil {
			return err
		}

		// Replace bank details only if they are provided
		if len(bankAccounts) > 0 {
			if err := tx.Where("user_id = ?", uID).Delete(&user.BankDetails{}).Error; err != nil {
				return err
			}
			for i := range bankAccounts {
				if err := tx.Create(&bankAccounts[i]).Error; err != nil {
					return err
				}
			}
		}

		// Replace nominees only if they are provided
		if len(nominees) > 0 {
			if err := tx.Where("user_id = ?", uID).Delete(&user.NomineeDetails{}).Error; err != nil {
				return err
			}
			for i := range nominees {
				if err := tx.Create(&nominees[i]).Error; err != nil {
					return err
				}
			}
		}

		var existingUser user.User
		updateData := map[string]interface{}{}
		if err := tx.Where("id = ?", userID).First(&existingUser).Error; err == nil {
			if existingUser.Role != "admin" {
				updateData["status"] = "pending_approval"
			}
		} else {
			updateData["status"] = "pending_approval"
		}
		
		if req.Mobile != "" {
			updateData["mobile"] = req.Mobile
		}
		if req.IncomeRange != "" {
			updateData["income_range"] = req.IncomeRange
		}
		if req.Occupation != "" {
			updateData["occupation"] = req.Occupation
		}
		if err := tx.Model(&user.User{}).Where("id = ?", userID).Updates(updateData).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		zap.L().Error("KYC submission failed", zap.Error(err))
		return errors.New("failed to submit KYC details")
	}

	return nil
}

