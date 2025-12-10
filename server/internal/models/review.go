package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Review captures rating and comment for a resource.
type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ResourceID uuid.UUID `gorm:"type:uuid;index" json:"resourceId"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	Score      int       `gorm:"check:score >= 1 AND score <= 5" json:"score"`
	Comment    string    `gorm:"type:text" json:"comment"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (r *Review) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
