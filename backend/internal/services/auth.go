package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/holyflow/backend/internal/config"
	"github.com/holyflow/backend/internal/models"
)

type GitLabUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type GitLabTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
}

type AuthService struct {
	db        *gorm.DB
	emailSvc  *EmailService
	jwtSecret string
	cfg       *config.Config
}

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthService(db *gorm.DB, emailSvc *EmailService, jwtSecret string, cfg *config.Config) *AuthService {
	return &AuthService{
		db:        db,
		emailSvc:  emailSvc,
		jwtSecret: jwtSecret,
		cfg:       cfg,
	}
}

func (s *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("user with this username or email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) LoginUser(username, password string) (*models.User, string, error) {
	// Find user
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", fmt.Errorf("invalid credentials")
		}
		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.db.Save(&user).Error; err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return &user, token, nil
}

func (s *AuthService) generateToken(user models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *AuthService) SendPasswordResetEmail(email string) error {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		// Don't reveal if user exists or not
		return nil
	}

	// Generate reset token (in a real app, you'd store this in the database)
	resetToken := uuid.New().String()

	// Send email
	subject := "Password Reset"
	body := fmt.Sprintf("Click the link to reset your password: https://yourdomain.com/reset-password?token=%s", resetToken)

	if s.emailSvc != nil {
		return s.emailSvc.SendEmail(context.Background(), email, subject, body)
	}

	// Mock email sending
	fmt.Printf("Mock email sent to %s with subject '%s'\n", email, subject)
	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	// In a real implementation, you would validate the token and update the password
	// For now, we'll just hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// In a real implementation, you would find the user by token and update their password
	fmt.Printf("Password would be reset to hashed value: %s\n", string(hashedPassword))
	return nil
}

func (s *AuthService) GetGitLabClientID() string {
	if s.cfg != nil {
		return s.cfg.GitLabClientID
	}
	return ""
}

func (s *AuthService) GetGitLabRedirectURL() string {
	if s.cfg != nil {
		return s.cfg.GitLabRedirectURL
	}
	return ""
}

func (s *AuthService) HandleGitLabCallback(code string) (*models.User, string, error) {
	// Exchange code for access token
	tokenResp, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from GitLab
	oidcUser, err := s.getGitLabUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get GitLab user info: %w", err)
	}

	// Find or create user in our database
	user, err := s.findOrCreateGitLabUser(oidcUser)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(*user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) exchangeCodeForToken(code string) (*GitLabTokenResponse, error) {
	// Prepare request
	data := url.Values{}
	data.Set("client_id", s.GetGitLabClientID())
	data.Set("client_secret", s.cfg.GitLabClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", s.GetGitLabRedirectURL())

	// Use GitLabBaseURL from configuration
	tokenURL := fmt.Sprintf("%s/oauth/token", s.cfg.GitLabBaseURL)

	// Make request to GitLab
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to GitLab: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp GitLabTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

func (s *AuthService) getGitLabUserInfo(accessToken string) (*GitLabUser, error) {
	// Use GitLabBaseURL from configuration
	apiURL := fmt.Sprintf("%s/api/v4/user", s.cfg.GitLabBaseURL)

	// Make request to GitLab API
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to GitLab API: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var oidcUser GitLabUser
	if err := json.Unmarshal(body, &oidcUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &oidcUser, nil
}

func (s *AuthService) findOrCreateGitLabUser(oidcUser *GitLabUser) (*models.User, error) {
	// Try to find existing user by GitLab ID or email
	var user models.User
	if err := s.db.Where("email = ?", oidcUser.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create new one
			user = models.User{
				Username: oidcUser.Username,
				Email:    oidcUser.Email,
				// For GitLab users, we'll set a random password since they'll login via OAuth
				PasswordHash: "", // Will be set to a random value
			}

			// Generate a random password for the user (they won't use it)
			randomPassword := uuid.New().String()
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash random password: %w", err)
			}
			user.PasswordHash = string(hashedPassword)

			// Create user
			if err := s.db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to query database: %w", err)
		}
	}

	return &user, nil
}
