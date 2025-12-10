package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Resource represents a shared asset in the repository.
type Resource struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string     `gorm:"size:255;not null" json:"title"`
	Description   string     `gorm:"type:text" json:"description"`
	Type          string     `gorm:"size:64" json:"type"` // tool/template/document/course
	Vendor        string     `gorm:"size:128" json:"vendor"`
	DeviceModel   string     `gorm:"size:128" json:"deviceModel"`
	Protocol      string     `gorm:"size:128" json:"protocol"`
	Scenario      string     `gorm:"size:128" json:"scenario"`
	Tags          string     `gorm:"size:512" json:"tags"` // comma-separated
	FilePath      string     `gorm:"size:512" json:"filePath"`
	FileName      string     `gorm:"size:255" json:"fileName"`
	ContentType   string     `gorm:"size:128" json:"contentType"`
	FileHash      string     `gorm:"size:64;index" json:"fileHash"`         // SHA256
	ExternalLink  string     `gorm:"size:512" json:"externalLink"`          // Optional external link
	Status        string     `gorm:"size:32;default:pending" json:"status"` // pending, approved, rejected
	RejectReason  string     `gorm:"size:255" json:"rejectReason"`
	DownloadCount int64      `gorm:"default:0" json:"downloadCount"`
	RatingAverage float64    `gorm:"default:0" json:"ratingAverage"`
	RatingCount   int64      `gorm:"default:0" json:"ratingCount"`
	ParentID      *uuid.UUID `gorm:"type:uuid;index" json:"parentId"` // Points to previous version
	Version       string     `gorm:"size:32;default:'1.0'" json:"version"`
	UploaderID    uuid.UUID  `gorm:"type:uuid" json:"uploaderId"`
	Uploader      User       `json:"uploader"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	SearchVector  string     `gorm:"type:tsvector" json:"-"`
}

func (r *Resource) BeforeCreate(_ *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
