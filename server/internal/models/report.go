package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Report represents a user report on a resource.
type Report struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	User       User      `json:"user"`
	ResourceID uuid.UUID `gorm:"type:uuid;index" json:"resourceId"`
	Resource   Resource  `json:"resource"`
	Reason     string    `gorm:"size:512" json:"reason"`
	Status     string    `gorm:"size:32;default:pending" json:"status"` // pending, resolved
	CreatedAt  time.Time `json:"createdAt"`
}

func (r *Report) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
