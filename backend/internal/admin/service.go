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
	GetUserDetails(userID string) (*auth.UserResponse, error)
	ApproveUser(userID string) error
	RejectUser(userID string) error
	BlockUser(userID string) error
	CloseAccount(userID string) error
	RejectClosure(userID string) error
	UploadShares(file multipart.File) error
	DeleteAllShares() error
}

type adminService struct {
	repo AdminRepository
}

func NewAdminService() AdminService {
	return &adminService{repo: NewAdminRepository()}
}

func (s *adminService) UploadShares(file multipart.File) error {
	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var shares []trade.Share
	for _, record := range records {
		if len(record) < 5 {
			continue // skip invalid rows
		}

		name := record[0]
		symbol := record[0] // Using Name as Symbol since Series is usually 'EQ' which causes unique constraint violation
		segmentName := record[2]
		price, _ := strconv.ParseFloat(record[3], 64)
		// We ignore record[1] (Series) and record[4] (Changes) for now
		qty := 1000000 // Hardcoded large qty

		segment, err := s.repo.GetOrCreateSegment(segmentName)
		if err != nil {
			return err
		}

		shares = append(shares, trade.Share{
			ID:              uuid.New(),
			Symbol:          symbol,
			Name:            name,
			Price:           price,
			SegmentID:       segment.ID,
			TotalShares:     qty,
			AvailableShares: qty,
		})
	}

	return s.repo.BulkInsertShares(shares)
}

func (s *adminService) DeleteAllShares() error {
	return s.repo.DeleteAllShares()
}

func (s *adminService) GetUsers() ([]auth.UserResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	walletRepo := wallet.NewWalletRepository()
	var res []auth.UserResponse
	for _, u := range users {
		w, _ := walletRepo.GetWallet(u.ID.String())
		var wb float64 = 0
		if w != nil {
			wb = w.WalletBalance
		}

		res = append(res, auth.UserResponse{
			ID:            u.ID.String(),
			Name:          u.Name,
			Email:         u.Email,
			Mobile:        u.Mobile,
			Role:          u.Role,
			Status:        u.Status,
			DOB:           u.DOB.Format("2006-01-02"),
			Address:       u.Address,
			PAN:           u.PAN,
			Aadhaar:       u.Aadhaar,
			IncomeRange:   u.IncomeRange,
			IPVPhoto:      u.IPVPhoto,
			WalletBalance: wb,
			CreatedAt:     u.CreatedAt.Format("2006-01-02"),
			UpdatedAt:     u.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}

func (s *adminService) GetUserDetails(userID string) (*auth.UserResponse, error) {
	u, personal, banks, nominees, err := s.repo.GetUserDetails(userID)
	if err != nil {
		return nil, err
	}

	walletRepo := wallet.NewWalletRepository()
	w, _ := walletRepo.GetWallet(userID)
	var wb float64 = 0
	if w != nil {
		wb = w.WalletBalance
	}

	res := &auth.UserResponse{
		ID:            u.ID.String(),
		Name:          u.Name,
		Email:         u.Email,
		Mobile:        u.Mobile,
		Role:          u.Role,
		Status:        u.Status,
		DOB:           u.DOB.Format("2006-01-02"),
		Address:       u.Address,
		PAN:           u.PAN,
		Aadhaar:       u.Aadhaar,
		IncomeRange:   u.IncomeRange,
		Occupation:    u.Occupation,
		IPVPhoto:      u.IPVPhoto,
		IPVLatitude:   u.IPVLatitude,
		IPVLongitude:  u.IPVLongitude,
		WalletBalance: wb,
		CreatedAt:     u.CreatedAt.Format("2006-01-02"),
		UpdatedAt:     u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if personal != nil {
		res.FatherName = personal.FatherName
		res.MotherName = personal.MotherName
		res.Country = personal.Country
		res.State = personal.State
		res.City = personal.City
		res.Pincode = personal.Pincode
	}

	for _, b := range banks {
		res.BankAccounts = append(res.BankAccounts, profile.BankAccountDTO{
			AccountType:   b.AccountType,
			IFSC:          b.IFSC,
			BankName:      b.BankName,
			Branch:        b.Branch,
			AccountNumber: b.AccountNumber,
		})
	}

	for _, n := range nominees {
		res.Nominees = append(res.Nominees, profile.NomineeDTO{
			Name:                 n.Name,
			DOB:                  n.DOB,
			PAN:                  n.PAN,
			Relationship:         n.Relationship,
			GuardianName:         n.GuardianName,
			GuardianRelationship: n.GuardianRelationship,
			GuardianPAN:          n.GuardianPAN,
			GuardianDOB:          n.GuardianDOB,
			Percentage:           n.Percentage,
		})
	}

	return res, nil
}

func (s *adminService) ApproveUser(userID string) error {
	return s.repo.UpdateUserStatus(userID, "active")
}

func (s *adminService) RejectUser(userID string) error {
	return s.repo.UpdateUserStatus(userID, "rejected")
}

func (s *adminService) BlockUser(userID string) error {
	return s.repo.UpdateUserStatus(userID, "blocked")
}

func (s *adminService) CloseAccount(userID string) error {
	userRepo := user.NewUserRepository()
	userRepo.DeleteAllSessions(userID)
	return s.repo.UpdateUserStatus(userID, "closed")
}

func (s *adminService) RejectClosure(userID string) error {
	return s.repo.UpdateUserStatus(userID, "active")
}
