package services

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/holyflow/backend/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Remove sensitive fields from updates
	delete(updates, "id")
	delete(updates, "password_hash")
	delete(updates, "email")

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (s *UserService) AddFavorite(userID, songID uuid.UUID) error {
	// Check if favorite already exists
	var existing models.Favorite
	if err := s.db.Where("user_id = ? AND song_id = ?", userID, songID).First(&existing).Error; err == nil {
		// Favorite already exists
		return nil
	}

	// Create new favorite
	favorite := &models.Favorite{
		UserID: userID,
		SongID: songID,
	}

	if err := s.db.Create(favorite).Error; err != nil {
		return fmt.Errorf("failed to add favorite: %w", err)
	}

	return nil
}

func (s *UserService) RemoveFavorite(userID, songID uuid.UUID) error {
	if err := s.db.Where("user_id = ? AND song_id = ?", userID, songID).Delete(&models.Favorite{}).Error; err != nil {
		return fmt.Errorf("failed to remove favorite: %w", err)
	}
	return nil
}

func (s *UserService) GetFavorites(userID uuid.UUID) ([]models.Song, error) {
	var songs []models.Song
	if err := s.db.Joins("JOIN favorites ON favorites.song_id = songs.id").
		Where("favorites.user_id = ?", userID).
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("failed to get favorites: %w", err)
	}
	return songs, nil
}

func (s *UserService) GetUserSongs(userID uuid.UUID) ([]models.Song, error) {
	var songs []models.Song
	if err := s.db.Where("created_by_id = ?", userID).
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("failed to get user songs: %w", err)
	}
	return songs, nil
}
