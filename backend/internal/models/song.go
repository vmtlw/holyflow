package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Song struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Artist      string         `json:"artist"`
	Text        string         `json:"text"`
	Rhythm      string         `json:"rhythm"`
	Mp3URL      string         `json:"mp3Url"`
	CoverURL    string         `json:"coverUrl"`
	CreatedByID uuid.UUID      `gorm:"type:uuid;not null" json:"createdById"`
	CreatedBy   User           `gorm:"foreignKey:CreatedByID" json:"-"`
	IsPublic    bool           `gorm:"default:true" json:"isPublic"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Song) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}
