package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Favorite marks a resource saved by a user.
type Favorite struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	ResourceID uuid.UUID `gorm:"type:uuid;index" json:"resourceId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (f *Favorite) BeforeCreate(_ *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
