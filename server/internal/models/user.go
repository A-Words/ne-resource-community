package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a platform user (uploader, reviewer, or admin).
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	DisplayName  string    `gorm:"size:255" json:"displayName"`
	Role         string    `gorm:"size:32;default:user" json:"role"`
	Points       int       `gorm:"default:0" json:"points"`
	Level        int       `gorm:"-" json:"level"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CalculateLevel()
	return nil
}

func (u *User) AfterFind(_ *gorm.DB) error {
	u.CalculateLevel()
	return nil
}

func (u *User) CalculateLevel() {
	u.Level = 1 + u.Points/100
}

// SetPassword hashes and stores password.
func (u *User) SetPassword(raw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword compares raw password with stored hash.
func (u *User) CheckPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(raw)) == nil
}
