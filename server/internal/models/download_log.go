package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DownloadLog tracks user download history for learning space analytics.
type DownloadLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	ResourceID uuid.UUID `gorm:"type:uuid;index" json:"resourceId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (d *DownloadLog) BeforeCreate(_ *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
