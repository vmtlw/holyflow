package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holyflow/backend/internal/services"
)

type SongHandler struct {
	songService *services.SongService
}

func NewSongHandler(songService *services.SongService) *SongHandler {
	return &SongHandler{
		songService: songService,
	}
}

type CreateSongRequest struct {
	Title    string `json:"title" binding:"required"`
	Artist   string `json:"artist" binding:"required"`
	Text     string `json:"text" binding:"required"`
	Rhythm   string `json:"rhythm" binding:"required"`
	Mp3Url   string `json:"mp3Url"`
	CoverUrl string `json:"coverUrl"`
	IsPublic *bool  `json:"isPublic"`
}

type UpdateSongRequest struct {
	Title    string `json:"title,omitempty"`
	Artist   string `json:"artist,omitempty"`
	Text     string `json:"text,omitempty"`
	Rhythm   string `json:"rhythm,omitempty"`
	IsPublic *bool  `json:"isPublic,omitempty"`
}

func (h *SongHandler) ListSongs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	songs, total, err := h.songService.ListSongs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songs": songs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *SongHandler) GetSong(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song, err := h.songService.GetSongByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"song": song})
}

func (h *SongHandler) CreateSong(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req CreateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, you would get the files from the request
	// For now, we'll pass nil for mp3File and coverFile
	song, err := h.songService.CreateSong(c.Request.Context(), req.Title, req.Artist, req.Text, req.Rhythm, userID, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update song with additional fields if provided
	updates := make(map[string]interface{})
	if req.Mp3Url != "" {
		updates["mp3_url"] = req.Mp3Url
	}
	if req.CoverUrl != "" {
		updates["cover_url"] = req.CoverUrl
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) > 0 {
		updatedSong, err := h.songService.UpdateSong(song.ID, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		song = updatedSong
	}

	c.JSON(http.StatusCreated, gin.H{"song": song})
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var req UpdateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Artist != "" {
		updates["artist"] = req.Artist
	}
	if req.Text != "" {
		updates["text"] = req.Text
	}
	if req.Rhythm != "" {
		updates["rhythm"] = req.Rhythm
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	song, err := h.songService.UpdateSong(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"song": song})
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	// Get user ID from context
	userID := c.MustGet("userID").(uuid.UUID)

	// Check if user is owner of the song
	song, err := h.songService.GetSongByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	if song.CreatedByID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this song"})
		return
	}

	if err := h.songService.DeleteSong(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// UploadMP3 handles MP3 file upload for a song
func (h *SongHandler) UploadMP3(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	// Generate filename
	filename := h.songService.GenerateFileName(file.Filename, "mp3")

	// Upload file to storage
	url, err := h.songService.UploadFile(c.Request.Context(), src, filename, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Update song with MP3 URL
	updates := map[string]interface{}{
		"mp3_url": url,
	}

	song, err := h.songService.UpdateSong(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update song: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MP3 uploaded successfully", "song": song})
}

// UploadCover handles cover image upload for a song
func (h *SongHandler) UploadCover(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	// Generate filename
	filename := h.songService.GenerateFileName(file.Filename, "cover")

	// Upload file to storage
	url, err := h.songService.UploadFile(c.Request.Context(), src, filename, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Update song with cover URL
	updates := map[string]interface{}{
		"cover_url": url,
	}

	song, err := h.songService.UpdateSong(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update song: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cover uploaded successfully", "song": song})
}
