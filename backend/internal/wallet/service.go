package wallet

import (
	"errors"
	"fmt"
	"time"

	"stock-trading/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService interface {
	GetBalance(userID string) (*WalletResponse, error)
	AddFunds(userID string, amount float64) error
	WithdrawFunds(userID string, amount float64) error
	GetTransactions(userID string) ([]Transaction, error)
}

type walletService struct {
	repo     WalletRepository
	userRepo user.UserRepository
}

func NewWalletService() WalletService {
	return &walletService{
		repo:     NewWalletRepository(),
		userRepo: user.NewUserRepository(),
	}
}

func (pService *walletService) GetBalance(pUserID string) (*WalletResponse, error) {
	lWallet, lErr := pService.repo.GetWallet(pUserID)
	if lErr != nil {
		if errors.Is(lErr, gorm.ErrRecordNotFound) {
			return &WalletResponse{WalletBalance: 0, BlockedBalance: 0, AvailableBalance: 0}, nil
		}
		return nil, errors.New("failed to retrieve wallet balance")
	}

	return &WalletResponse{
		WalletBalance:    lWallet.WalletBalance,
		BlockedBalance:   lWallet.BlockedBalance,
		AvailableBalance: lWallet.AvailableBalance,
	}, nil
}

func (pService *walletService) AddFunds(pUserID string, pAmount float64) error {
	lUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return errors.New("invalid user id")
	}
	if pAmount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	lU, lErr := pService.userRepo.GetUserByID(pUserID)
	if lErr != nil {
		return errors.New("invalid user")
	}
	if lU.Status == "closure_requested" {
		return errors.New("account closure requested, action not permitted")
	}

	return pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		lWallet, lTxErr := pService.repo.GetWalletForUpdate(pTx, pUserID)

		// If wallet doesn't exist, create it
		if errors.Is(lTxErr, gorm.ErrRecordNotFound) {
			lWallet = &Wallet{
				UserID:           lUID,
				WalletBalance:    pAmount,
				AvailableBalance: pAmount,
				BlockedBalance:   0,
			}
			if lErr := pService.repo.CreateWallet(pTx, lWallet); lErr != nil {
				return lErr
			}
		} else if lTxErr != nil {
			return lTxErr
		} else {
			// Update existing wallet
			lWallet.WalletBalance += pAmount
			lWallet.AvailableBalance += pAmount
			if lErr := pService.repo.UpdateWallet(pTx, lWallet); lErr != nil {
				return lErr
			}
		}

		// Create transaction log
		lTransaction := &Transaction{
			ID:          uuid.New(),
			UserID:      lUID,
			Type:        "add_fund",
			Amount:      pAmount,
			ReferenceID: fmt.Sprintf("DEP-%d", time.Now().UnixNano()),
			Description: "Deposit to wallet",
			Status:      "completed",
		}

		return pService.repo.CreateTransaction(pTx, lTransaction)
	})
}

func (pService *walletService) WithdrawFunds(pUserID string, pAmount float64) error {
	lUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return errors.New("invalid user id")
	}
	if pAmount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	lU, lErr := pService.userRepo.GetUserByID(pUserID)
	if lErr != nil {
		return errors.New("invalid user")
	}
	if lU.Status == "closure_requested" {
		return errors.New("account closure requested, action not permitted")
	}

	return pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		lWallet, lTxErr := pService.repo.GetWalletForUpdate(pTx, pUserID)
		if lTxErr != nil {
			return errors.New("wallet not found")
		}

		if lWallet.AvailableBalance < pAmount {
			return errors.New("insufficient available balance")
		}

		lWallet.WalletBalance -= pAmount
		lWallet.AvailableBalance -= pAmount

		if lErr := pService.repo.UpdateWallet(pTx, lWallet); lErr != nil {
			return lErr
		}

		lTransaction := &Transaction{
			ID:          uuid.New(),
			UserID:      lUID,
			Type:        "withdraw",
			Amount:      pAmount,
			ReferenceID: fmt.Sprintf("WDL-%d", time.Now().UnixNano()),
			Description: "Withdrawal from wallet",
			Status:      "completed",
		}

		return pService.repo.CreateTransaction(pTx, lTransaction)
	})
}

func (pService *walletService) GetTransactions(pUserID string) ([]Transaction, error) {
	return pService.repo.GetTransactionsByUser(pUserID)
}
