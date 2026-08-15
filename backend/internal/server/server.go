package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/holyflow/backend/internal/config"
	"github.com/holyflow/backend/internal/database"
	"github.com/holyflow/backend/internal/handlers"
	"github.com/holyflow/backend/internal/middleware"
	"github.com/holyflow/backend/internal/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     *config.Config
	db         *database.DB
	router     *gin.Engine
	authSvc    *services.AuthService
	userSvc    *services.UserService
	songSvc    *services.SongService
	storageSvc *services.StorageService
	emailSvc   *services.EmailService
}

func NewServer(cfg *config.Config) *Server {
	// Initialize database connection
	db, err := database.NewDB(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize services
	storageSvc := services.NewStorageService(cfg)
	emailSvc := services.NewEmailService(cfg)
	authSvc := services.NewAuthService(db.DB, emailSvc, cfg.JWTSecret, cfg)
	userSvc := services.NewUserService(db.DB)
	songSvc := services.NewSongService(db.DB, storageSvc)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authSvc, userSvc)
	userHandler := handlers.NewUserHandler(userSvc)
	songHandler := handlers.NewSongHandler(songSvc)

	server := &Server{
		config:     cfg,
		db:         db,
		authSvc:    authSvc,
		userSvc:    userSvc,
		songSvc:    songSvc,
		storageSvc: storageSvc,
		emailSvc:   emailSvc,
	}

	// Setup router
	server.setupRouter(authHandler, userHandler, songHandler)

	return server
}

func (s *Server) setupRouter(authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, songHandler *handlers.SongHandler) {
	// Set release mode for production
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.GET("/gitlab/login", authHandler.GitLabLogin)
			auth.GET("/gitlab/callback", authHandler.GitLabCallback)
		}

		// Public song routes
		publicSongs := v1.Group("/songs")
		{
			publicSongs.GET("/", songHandler.ListSongs)
			publicSongs.GET("/:id", songHandler.GetSong)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(s.authSvc))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PUT("/me", userHandler.UpdateMe)
				users.GET("/me/favorites", userHandler.GetFavorites)
				users.GET("/me/songs", userHandler.GetUserSongs)
				users.POST("/me/favorites/:songId", userHandler.AddFavorite)
				users.DELETE("/me/favorites/:songId", userHandler.RemoveFavorite)
			}

			// Protected song routes
			songs := protected.Group("/songs")
			{
				songs.POST("/", songHandler.CreateSong)
				songs.PUT("/:id", songHandler.UpdateSong)
				songs.DELETE("/:id", songHandler.DeleteSong)
				// File upload routes
				songs.POST("/:id/upload-mp3", songHandler.UploadMP3)
				songs.POST("/:id/upload-cover", songHandler.UploadCover)
			}
		}

	}

	s.router = router
}

func (s *Server) Run() error {
	// Run server in a goroutine
	srv := &http.Server{
		Addr:    s.config.GetServerAddress(),
		Handler: s.router,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server is running on port %s", s.config.ServerPort)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server exited")
	return nil
}
