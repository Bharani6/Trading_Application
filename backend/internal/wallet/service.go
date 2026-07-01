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

func (s *walletService) GetBalance(userID string) (*WalletResponse, error) {
	wallet, err := s.repo.GetWallet(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &WalletResponse{WalletBalance: 0, BlockedBalance: 0, AvailableBalance: 0}, nil
		}
		return nil, errors.New("failed to retrieve wallet balance")
	}

	return &WalletResponse{
		WalletBalance:    wallet.WalletBalance,
		BlockedBalance:   wallet.BlockedBalance,
		AvailableBalance: wallet.AvailableBalance,
	}, nil
}

func (s *walletService) AddFunds(userID string, amount float64) error {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return s.repo.RunInTransaction(func(tx *gorm.DB) error {
		wallet, err := s.repo.GetWalletForUpdate(tx, userID)

		// If wallet doesn't exist, create it
		if errors.Is(err, gorm.ErrRecordNotFound) {
			wallet = &Wallet{
				UserID:           uID,
				WalletBalance:    amount,
				AvailableBalance: amount,
				BlockedBalance:   0,
			}
			if err := s.repo.CreateWallet(tx, wallet); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// Update existing wallet
			wallet.WalletBalance += amount
			wallet.AvailableBalance += amount
			if err := s.repo.UpdateWallet(tx, wallet); err != nil {
				return err
			}
		}

		// Create transaction log
		transaction := &Transaction{
			ID:          uuid.New(),
			UserID:      uID,
			Type:        "add_fund",
			Amount:      amount,
			ReferenceID: fmt.Sprintf("DEP-%d", time.Now().UnixNano()),
			Description: "Deposit to wallet",
			Status:      "completed",
		}

		return s.repo.CreateTransaction(tx, transaction)
	})
}

func (s *walletService) WithdrawFunds(userID string, amount float64) error {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return s.repo.RunInTransaction(func(tx *gorm.DB) error {
		wallet, err := s.repo.GetWalletForUpdate(tx, userID)
		if err != nil {
			return errors.New("wallet not found")
		}

		if wallet.AvailableBalance < amount {
			return errors.New("insufficient available balance")
		}

		wallet.WalletBalance -= amount
		wallet.AvailableBalance -= amount

		if err := s.repo.UpdateWallet(tx, wallet); err != nil {
			return err
		}

		transaction := &Transaction{
			ID:          uuid.New(),
			UserID:      uID,
			Type:        "withdraw",
			Amount:      amount,
			ReferenceID: fmt.Sprintf("WDL-%d", time.Now().UnixNano()),
			Description: "Withdrawal from wallet",
			Status:      "completed",
		}

		return s.repo.CreateTransaction(tx, transaction)
	})
}

func (s *walletService) GetTransactions(userID string) ([]Transaction, error) {
	return s.repo.GetTransactionsByUser(userID)
}
