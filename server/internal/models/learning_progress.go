package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearningProgress struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	ResourceID uuid.UUID `gorm:"type:uuid;index" json:"resourceId"`
	Progress   int       `gorm:"default:0" json:"progress"`             // Percentage 0-100 or seconds
	Status     string    `gorm:"size:32;default:started" json:"status"` // started, completed
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (lp *LearningProgress) BeforeCreate(_ *gorm.DB) error {
	if lp.ID == uuid.Nil {
		lp.ID = uuid.New()
	}
	return nil
}
