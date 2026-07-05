package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Name         string         `json:"name" gorm:"type:varchar(100);not null"`
	Mobile       string         `json:"mobile" gorm:"type:varchar(20);unique;not null"`
	Email        string         `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash string         `json:"-" gorm:"type:varchar(255);not null"`
	PAN          string         `json:"pan" gorm:"type:varchar(10);unique;not null"`
	Aadhaar      string         `json:"aadhaar" gorm:"type:varchar(12);unique;not null"`
	Address      string         `json:"address" gorm:"type:text;not null"`
	IncomeRange  string         `json:"income_range" gorm:"type:varchar(50);not null"`
	DOB          time.Time      `json:"dob" gorm:"type:date;not null"`
	Role         string         `json:"role" gorm:"type:varchar(20);default:'user'"`
	Status       string         `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type Session struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	AccessToken  string    `json:"access_token" gorm:"type:text;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	IPAddress    string    `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent    string    `json:"user_agent" gorm:"type:varchar(255)"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type BankDetails struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:char(36);index"`
	AccountType   string    `json:"account_type" gorm:"type:varchar(20)"` // "primary" or "secondary"
	IFSC          string    `json:"ifsc" gorm:"type:varchar(20)"`
	BankName      string    `json:"bank_name" gorm:"type:varchar(100)"`
	Branch        string    `json:"branch" gorm:"type:varchar(100)"`
	AccountNumber string    `json:"account_number" gorm:"type:varchar(50)"`
	IncomeRange   string    `json:"income_range" gorm:"type:varchar(50)"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type NomineeDetails struct {
	ID                   uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID               uuid.UUID `json:"user_id" gorm:"type:char(36);index"`
	Name                 string    `json:"name" gorm:"type:varchar(100);not null"`
	DOB                  string    `json:"dob" gorm:"type:varchar(20);not null"`
	PAN                  string    `json:"pan" gorm:"type:varchar(20)"`
	Relationship         string    `json:"relationship" gorm:"type:varchar(50);not null"`
	GuardianName         string    `json:"guardian_name" gorm:"type:varchar(100)"`
	GuardianRelationship string    `json:"guardian_relationship" gorm:"type:varchar(50)"`
	GuardianPAN          string    `json:"guardian_pan" gorm:"type:varchar(20)"`
	GuardianDOB          string    `json:"guardian_dob" gorm:"type:varchar(20)"`
	Percentage           float64   `json:"percentage" gorm:"type:decimal(5,2);not null;default:100.00"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
type PersonalDetails struct {
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	FatherName string    `json:"father_name" gorm:"type:varchar(100)"`
	MotherName string    `json:"mother_name" gorm:"type:varchar(100)"`
	Country    string    `json:"country" gorm:"type:varchar(100)"`
	State      string    `json:"state" gorm:"type:varchar(100)"`
	City       string    `json:"city" gorm:"type:varchar(100)"`
	Address    string    `json:"address" gorm:"type:text"`
	Pincode    string    `json:"pincode" gorm:"type:varchar(20)"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (b *BankDetails) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
func (n *NomineeDetails) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return
}

type PasswordResetToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Token     string    `json:"token" gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
