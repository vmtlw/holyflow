package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/holyflow/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hide sensitive information
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Hide sensitive information
	user.PasswordHash = ""

	resp := LoginResponse{
		User:  user,
		Token: token,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// In a real implementation, you would validate the refresh token and generate a new access token
	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.SendPasswordResetEmail(req.Email); err != nil {
		// Don't reveal if user exists or not
		c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset instructions have been sent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset instructions have been sent"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *AuthHandler) GitLabLogin(c *gin.Context) {
	// Redirect to GitLab OAuth endpoint
	oauthState := uuid.New().String()

	// Store state in session or temporary storage for validation in callback
	// For simplicity, we'll pass it as a parameter (not recommended for production)

	url := fmt.Sprintf(
		"https://gitlab.vmtlw.ru/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=read_user",
		h.authService.GetGitLabClientID(),
		h.authService.GetGitLabRedirectURL(),
		oauthState,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GitLabCallback(c *gin.Context) {
	// Get code from query parameters
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	// Handle GitLab callback
	user, token, err := h.authService.HandleGitLabCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to handle GitLab callback: %v", err)})
		return
	}

	// Hide sensitive information
	user.PasswordHash = ""

	resp := LoginResponse{
		User:  user,
		Token: token,
	}

	c.JSON(http.StatusOK, resp)
}
