package profile

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileUpload struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:char(36);index"`
	FileName      string    `json:"file_name" gorm:"type:varchar(255)"`
	FilePath      string    `json:"file_path" gorm:"type:text"`
	FileSize      int64     `json:"file_size"`
	FileType      string    `json:"file_type" gorm:"type:enum('pan_card', 'address_proof', 'bank_proof')"`
	FileExtension string    `json:"file_extension" gorm:"type:varchar(10)"`
	CreatedAt     time.Time `json:"created_at"`
}

func (f *FileUpload) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
