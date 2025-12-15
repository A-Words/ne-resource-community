package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Request struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Bounty      int       `gorm:"default:0" json:"bounty"`
	Status      string    `gorm:"size:32;default:open" json:"status"` // open, fulfilled, closed
	UserID      uuid.UUID `gorm:"type:uuid" json:"userId"`
	User        User      `json:"user"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (r *Request) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
