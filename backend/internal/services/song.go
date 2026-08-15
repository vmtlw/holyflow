package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/holyflow/backend/internal/models"
)

type SongService struct {
	db         *gorm.DB
	storageSvc *StorageService
}

func NewSongService(db *gorm.DB, storageSvc *StorageService) *SongService {
	return &SongService{
		db:         db,
		storageSvc: storageSvc,
	}
}

func (s *SongService) CreateSong(ctx context.Context, title, artist, text, rhythm string, createdBy uuid.UUID, mp3File, coverFile multipart.File) (*models.Song, error) {
	song := &models.Song{
		Title:       title,
		Artist:      artist,
		Text:        text,
		Rhythm:      rhythm,
		CreatedByID: createdBy,
		IsPublic:    true,
	}

	// Upload MP3 file if provided
	if mp3File != nil {
		// In a real implementation, you would get the filename and content type from the request
		filename := s.storageSvc.GenerateFileName("song.mp3", "mp3")
		url, err := s.storageSvc.UploadFile(ctx, mp3File, filename, "audio/mpeg")
		if err != nil {
			return nil, fmt.Errorf("failed to upload mp3 file: %w", err)
		}
		song.Mp3URL = url
	}

	// Upload cover image if provided
	if coverFile != nil {
		// In a real implementation, you would get the filename and content type from the request
		filename := s.storageSvc.GenerateFileName("cover.jpg", "cover")
		url, err := s.storageSvc.UploadFile(ctx, coverFile, filename, "image/jpeg")
		if err != nil {
			return nil, fmt.Errorf("failed to upload cover file: %w", err)
		}
		song.CoverURL = url
	}

	// Save song to database
	if err := s.db.Create(song).Error; err != nil {
		return nil, fmt.Errorf("failed to create song: %w", err)
	}

	return song, nil
}

func (s *SongService) GetSongByID(id uuid.UUID) (*models.Song, error) {
	var song models.Song
	if err := s.db.Preload("CreatedBy").First(&song, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get song: %w", err)
	}
	return &song, nil
}

func (s *SongService) ListSongs(page, limit int) ([]models.Song, int64, error) {
	var songs []models.Song
	var total int64

	offset := (page - 1) * limit

	// Get total count
	if err := s.db.Model(&models.Song{}).Where("is_public = ?", true).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count songs: %w", err)
	}

	// Get songs
	if err := s.db.Preload("CreatedBy").
		Where("is_public = ?", true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&songs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list songs: %w", err)
	}

	return songs, total, nil
}

func (s *SongService) UpdateSong(id uuid.UUID, updates map[string]interface{}) (*models.Song, error) {
	var song models.Song
	if err := s.db.First(&song, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("song not found: %w", err)
	}

	// Remove protected fields from updates
	delete(updates, "id")
	delete(updates, "created_by_id")

	// Explicitly update each field to handle nullable fields correctly
	if title, ok := updates["title"]; ok {
		song.Title = title.(string)
	}
	if artist, ok := updates["artist"]; ok {
		song.Artist = artist.(string)
	}
	if text, ok := updates["text"]; ok {
		song.Text = text.(string)
	}
	if rhythm, ok := updates["rhythm"]; ok {
		song.Rhythm = rhythm.(string)
	}
	if mp3Url, ok := updates["mp3_url"]; ok {
		song.Mp3URL = mp3Url.(string)
	}
	if coverUrl, ok := updates["cover_url"]; ok {
		song.CoverURL = coverUrl.(string)
	}
	if isPublic, ok := updates["is_public"]; ok {
		song.IsPublic = isPublic.(bool)
	}

	// Save the updated song
	if err := s.db.Save(&song).Error; err != nil {
		return nil, fmt.Errorf("failed to update song: %w", err)
	}

	return &song, nil
}

func (s *SongService) DeleteSong(id uuid.UUID) error {
	if err := s.db.Delete(&models.Song{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}
	return nil
}

// UploadFile uploads a file using the storage service
func (s *SongService) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	return s.storageSvc.UploadFile(ctx, file, filename, contentType)
}

// GenerateFileName generates a unique filename
func (s *SongService) GenerateFileName(originalName string, prefix string) string {
	return s.storageSvc.GenerateFileName(originalName, prefix)
}
