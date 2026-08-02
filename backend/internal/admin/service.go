package admin

import (
	"encoding/csv"
	"mime/multipart"
	"strconv"

	"stock-trading/internal/auth"
	"stock-trading/internal/profile"
	"stock-trading/internal/trade"
	"stock-trading/internal/user"
	"stock-trading/internal/wallet"

	"github.com/google/uuid"
)

type AdminService interface {
	GetUsers() ([]auth.UserResponse, error)
	GetUserDetails(pUserID string) (*auth.UserResponse, error)
	ApproveUser(pUserID string) error
	RejectUser(pUserID string) error
	BlockUser(pUserID string) error
	CloseAccount(pUserID string) error
	RejectClosure(pUserID string) error
	UploadShares(pFile multipart.File) error
	DeleteAllShares() error
}

type adminService struct {
	repo AdminRepository
}

func NewAdminService() AdminService {
	return &adminService{repo: NewAdminRepository()}
}

func (pService *adminService) UploadShares(pFile multipart.File) error {
	lReader := csv.NewReader(pFile)
	// Skip header
	if _, lErr := lReader.Read(); lErr != nil {
		return lErr
	}

	lRecords, lErr := lReader.ReadAll()
	if lErr != nil {
		return lErr
	}

	var lShares []trade.Share
	for _, lRecord := range lRecords {
		if len(lRecord) < 5 {
			continue // skip invalid rows
		}

		lName := lRecord[0]
		lSymbol := lRecord[0] // Using Name as Symbol since Series is usually 'EQ' which causes unique constraint violation
		lSegmentName := lRecord[2]
		lPrice, _ := strconv.ParseFloat(lRecord[3], 64)
		// We ignore lRecord[1] (Series) and lRecord[4] (Changes) for now
		lQty := 1000000 // Hardcoded large qty

		lSegment, lSegErr := pService.repo.GetOrCreateSegment(lSegmentName)
		if lSegErr != nil {
			return lSegErr
		}

		lShares = append(lShares, trade.Share{
			ID:              uuid.New(),
			Symbol:          lSymbol,
			Name:            lName,
			Price:           lPrice,
			SegmentID:       lSegment.ID,
			TotalShares:     lQty,
			AvailableShares: lQty,
		})
	}

	return pService.repo.BulkInsertShares(lShares)
}

func (pService *adminService) DeleteAllShares() error {
	return pService.repo.DeleteAllShares()
}

func (pService *adminService) GetUsers() ([]auth.UserResponse, error) {
	lUsers, lErr := pService.repo.GetAllUsers()
	if lErr != nil {
		return nil, lErr
	}

	lWalletRepo := wallet.NewWalletRepository()
	var lRes []auth.UserResponse
	for _, lUser := range lUsers {
		lWallet, _ := lWalletRepo.GetWallet(lUser.ID.String())
		var lWalletBalance float64 = 0
		if lWallet != nil {
			lWalletBalance = lWallet.WalletBalance
		}

		lRes = append(lRes, auth.UserResponse{
			ID:            lUser.ID.String(),
			Name:          lUser.Name,
			Email:         lUser.Email,
			Mobile:        lUser.Mobile,
			Role:          lUser.Role,
			Status:        lUser.Status,
			DOB:           lUser.DOB.Format("2006-01-02"),
			Address:       lUser.Address,
			PAN:           lUser.PAN,
			Aadhaar:       lUser.Aadhaar,
			IncomeRange:   lUser.IncomeRange,
			IPVPhoto:      lUser.IPVPhoto,
			WalletBalance: lWalletBalance,
			CreatedAt:     lUser.CreatedAt.Format("2006-01-02"),
			UpdatedAt:     lUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return lRes, nil
}

func (pService *adminService) GetUserDetails(pUserID string) (*auth.UserResponse, error) {
	lUser, lPersonal, lBanks, lNominees, lErr := pService.repo.GetUserDetails(pUserID)
	if lErr != nil {
		return nil, lErr
	}

	lWalletRepo := wallet.NewWalletRepository()
	lWallet, _ := lWalletRepo.GetWallet(pUserID)
	var lWalletBalance float64 = 0
	if lWallet != nil {
		lWalletBalance = lWallet.WalletBalance
	}

	lRes := &auth.UserResponse{
		ID:            lUser.ID.String(),
		Name:          lUser.Name,
		Email:         lUser.Email,
		Mobile:        lUser.Mobile,
		Role:          lUser.Role,
		Status:        lUser.Status,
		DOB:           lUser.DOB.Format("2006-01-02"),
		Address:       lUser.Address,
		PAN:           lUser.PAN,
		Aadhaar:       lUser.Aadhaar,
		IncomeRange:   lUser.IncomeRange,
		Occupation:    lUser.Occupation,
		IPVPhoto:      lUser.IPVPhoto,
		IPVLatitude:   lUser.IPVLatitude,
		IPVLongitude:  lUser.IPVLongitude,
		WalletBalance: lWalletBalance,
		CreatedAt:     lUser.CreatedAt.Format("2006-01-02"),
		UpdatedAt:     lUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if lPersonal != nil {
		lRes.FatherName = lPersonal.FatherName
		lRes.MotherName = lPersonal.MotherName
		lRes.Country = lPersonal.Country
		lRes.State = lPersonal.State
		lRes.City = lPersonal.City
		lRes.Pincode = lPersonal.Pincode
	}

	for _, lBank := range lBanks {
		lRes.BankAccounts = append(lRes.BankAccounts, profile.BankAccountDTO{
			AccountType:   lBank.AccountType,
			IFSC:          lBank.IFSC,
			BankName:      lBank.BankName,
			Branch:        lBank.Branch,
			AccountNumber: lBank.AccountNumber,
		})
	}

	for _, lNominee := range lNominees {
		lRes.Nominees = append(lRes.Nominees, profile.NomineeDTO{
			Name:                 lNominee.Name,
			DOB:                  lNominee.DOB,
			PAN:                  lNominee.PAN,
			Relationship:         lNominee.Relationship,
			GuardianName:         lNominee.GuardianName,
			GuardianRelationship: lNominee.GuardianRelationship,
			GuardianPAN:          lNominee.GuardianPAN,
			GuardianDOB:          lNominee.GuardianDOB,
			Percentage:           lNominee.Percentage,
		})
	}

	return lRes, nil
}

func (pService *adminService) ApproveUser(pUserID string) error {
	return pService.repo.UpdateUserStatus(pUserID, "active")
}

func (pService *adminService) RejectUser(pUserID string) error {
	return pService.repo.UpdateUserStatus(pUserID, "rejected")
}

func (pService *adminService) BlockUser(pUserID string) error {
	return pService.repo.UpdateUserStatus(pUserID, "blocked")
}

func (pService *adminService) CloseAccount(pUserID string) error {
	lUserRepo := user.NewUserRepository()
	lUserRepo.DeleteAllSessions(pUserID)
	return pService.repo.UpdateUserStatus(pUserID, "closed")
}

func (pService *adminService) RejectClosure(pUserID string) error {
	return pService.repo.UpdateUserStatus(pUserID, "active")
}
