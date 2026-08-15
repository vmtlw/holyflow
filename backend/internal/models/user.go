package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Username        string         `gorm:"uniqueIndex;not null" json:"username"`
	Email           string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string         `gorm:"not null" json:"-"`
	FirstName       string         `json:"firstName"`
	LastName        string         `json:"lastName"`
	AvatarURL       string         `json:"avatarUrl"`
	IsEmailVerified bool           `json:"isEmailVerified"`
	IsActive        bool           `gorm:"default:true" json:"isActive"`
	LastLoginAt     *time.Time     `json:"lastLoginAt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

type Favorite struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	SongID    uuid.UUID `gorm:"type:uuid;not null;index" json:"songId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	f.ID = uuid.New()
	return nil
}
