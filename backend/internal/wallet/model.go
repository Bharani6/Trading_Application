package wallet

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	UserID           uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	WalletBalance    float64   `json:"wallet_balance" gorm:"type:decimal(15,2);default:0.00"`
	BlockedBalance   float64   `json:"blocked_balance" gorm:"type:decimal(15,2);default:0.00"`
	AvailableBalance float64   `json:"available_balance" gorm:"type:decimal(15,2);default:0.00"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Transaction struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:char(36);index"`
	Type        string    `json:"type" gorm:"type:varchar(50)"`
	Amount      float64   `json:"amount" gorm:"type:decimal(15,2)"`
	ReferenceID string    `json:"reference_id" gorm:"type:varchar(100);uniqueIndex"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'completed'"`
	CreatedAt   time.Time `json:"created_at"`
}
