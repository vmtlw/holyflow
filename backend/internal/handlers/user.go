package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holyflow/backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type AddFavoriteRequest struct {
	SongID uuid.UUID `json:"songId" binding:"required"`
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hide sensitive information
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	user, err := h.userService.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hide sensitive information
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetFavorites(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	songs, err := h.userService.GetFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

func (h *UserHandler) AddFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	songID, err := uuid.Parse(c.Param("songId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	if err := h.userService.AddFavorite(userID, songID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite added successfully"})
}

func (h *UserHandler) RemoveFavorite(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	songID, err := uuid.Parse(c.Param("songId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	if err := h.userService.RemoveFavorite(userID, songID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}

func (h *UserHandler) GetUserSongs(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	songs, err := h.userService.GetUserSongs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}
