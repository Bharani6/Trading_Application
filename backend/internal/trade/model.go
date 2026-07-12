package trade

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Segment struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(20);unique"`
}

type Share struct {
	ID              uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Symbol          string         `json:"symbol" gorm:"type:varchar(50);uniqueIndex;not null"`
	Name            string         `json:"name" gorm:"type:varchar(150);not null"`
	Price           float64        `json:"price" gorm:"type:decimal(15,2);not null"`
	PreviousPrice   float64        `json:"previous_price" gorm:"type:decimal(15,2);default:0"`
	SegmentID       uint           `json:"segment_id"`
	Segment         Segment        `json:"segment" gorm:"foreignKey:SegmentID"`
	TotalShares     int            `json:"total_shares"`
	AvailableShares int            `json:"available_shares"`
	Version         int            `json:"version" gorm:"default:1"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (s *Share) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type Trade struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);index"`
	ShareID   uuid.UUID `json:"share_id" gorm:"type:char(36);index"`
	Share     Share     `json:"share" gorm:"foreignKey:ShareID"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"type:decimal(15,2);not null"`
	Type      string    `json:"type" gorm:"type:varchar(20)"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:'completed'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Trade) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
