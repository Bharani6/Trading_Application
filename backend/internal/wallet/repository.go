package wallet

import (
	"stock-trading/internal/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetWallet(userID string) (*Wallet, error)
	GetWalletForUpdate(tx *gorm.DB, userID string) (*Wallet, error)
	CreateWallet(tx *gorm.DB, wallet *Wallet) error
	UpdateWallet(tx *gorm.DB, wallet *Wallet) error
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

func (r *walletRepository) GetWallet(userID string) (*Wallet, error) {
	var wallet Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *walletRepository) GetWalletForUpdate(tx *gorm.DB, userID string) (*Wallet, error) {
	var wallet Wallet
	// SELECT * FROM wallets WHERE user_id = ? FOR UPDATE (Row level lock)
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *walletRepository) CreateWallet(tx *gorm.DB, wallet *Wallet) error {
	return tx.Create(wallet).Error
}

func (r *walletRepository) UpdateWallet(tx *gorm.DB, wallet *Wallet) error {
	return tx.Save(wallet).Error
}

func (r *walletRepository) CreateTransaction(tx *gorm.DB, transaction *Transaction) error {
	return tx.Create(transaction).Error
}

func (r *walletRepository) GetTransactionsByUser(userID string) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
	return transactions, err
}

func (r *walletRepository) RunInTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
