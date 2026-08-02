package wallet

import (
	"errors"
	"stock-trading/internal/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetWallet(userID string) (*Wallet, error)
	GetWalletForUpdate(tx *gorm.DB, userID string) (*Wallet, error)
	CreateWallet(tx *gorm.DB, wallet *Wallet) error
	UpdateWallet(tx *gorm.DB, wallet *Wallet) error
	UpdateWalletWithVersion(tx *gorm.DB, wallet *Wallet) error
	CreateTransaction(tx *gorm.DB, transaction *Transaction) error
	GetTransactionsByUser(userID string) ([]Transaction, error)
	RunInTransaction(fn func(tx *gorm.DB) error) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository() WalletRepository {
	return &walletRepository{db: database.DB}
}

func (pRepo *walletRepository) GetWallet(pUserID string) (*Wallet, error) {
	var lWallet Wallet
	lErr := pRepo.db.Where("user_id = ?", pUserID).First(&lWallet).Error
	return &lWallet, lErr
}

func (pRepo *walletRepository) GetWalletForUpdate(pTx *gorm.DB, pUserID string) (*Wallet, error) {
	var lWallet Wallet
	// SELECT * FROM wallets WHERE user_id = ? FOR UPDATE (Row level lock)
	lErr := pTx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", pUserID).First(&lWallet).Error
	return &lWallet, lErr
}

func (pRepo *walletRepository) CreateWallet(pTx *gorm.DB, pWallet *Wallet) error {
	return pTx.Create(pWallet).Error
}

func (pRepo *walletRepository) UpdateWallet(pTx *gorm.DB, pWallet *Wallet) error {
	return pTx.Save(pWallet).Error
}

func (pRepo *walletRepository) UpdateWalletWithVersion(pTx *gorm.DB, pWallet *Wallet) error {
	lResult := pTx.Model(pWallet).Where("version = ?", pWallet.Version).Updates(map[string]interface{}{
		"wallet_balance":    pWallet.WalletBalance,
		"blocked_balance":   pWallet.BlockedBalance,
		"available_balance": pWallet.AvailableBalance,
		"version":           pWallet.Version + 1,
	})
	if lResult.Error != nil {
		return lResult.Error
	}
	if lResult.RowsAffected == 0 {
		return errors.New("optimistic lock failed for wallet")
	}
	pWallet.Version++
	return nil
}

func (pRepo *walletRepository) CreateTransaction(pTx *gorm.DB, pTransaction *Transaction) error {
	return pTx.Create(pTransaction).Error
}

func (pRepo *walletRepository) GetTransactionsByUser(pUserID string) ([]Transaction, error) {
	var lTransactions []Transaction
	lErr := pRepo.db.Where("user_id = ?", pUserID).Order("created_at desc").Find(&lTransactions).Error
	return lTransactions, lErr
}

func (pRepo *walletRepository) RunInTransaction(pFn func(pTx *gorm.DB) error) error {
	return pRepo.db.Transaction(pFn)
}
