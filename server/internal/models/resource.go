package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Resource represents a shared asset in the repository.
type Resource struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string    `gorm:"size:255;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Type          string    `gorm:"size:64" json:"type"`          // tool/template/document/course
	Vendor        string    `gorm:"size:128" json:"vendor"`
	DeviceModel   string    `gorm:"size:128" json:"deviceModel"`
	Protocol      string    `gorm:"size:128" json:"protocol"`
	Scenario      string    `gorm:"size:128" json:"scenario"`
	Tags          string    `gorm:"size:512" json:"tags"`          // comma-separated
	FilePath      string    `gorm:"size:512" json:"filePath"`
	FileName      string    `gorm:"size:255" json:"fileName"`
	ContentType   string    `gorm:"size:128" json:"contentType"`
	DownloadCount int64     `gorm:"default:0" json:"downloadCount"`
	RatingAverage float64   `gorm:"default:0" json:"ratingAverage"`
	RatingCount   int64     `gorm:"default:0" json:"ratingCount"`
	UploaderID    uuid.UUID `gorm:"type:uuid" json:"uploaderId"`
	Uploader      User      `json:"uploader"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	SearchVector  string    `gorm:"type:tsvector" json:"-"`
}

func (r *Resource) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
