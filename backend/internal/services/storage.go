package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/holyflow/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	client *minio.Client
	bucket string
	cfg    *config.Config
}

func NewStorageService(cfg *config.Config) *StorageService {
	// If S3 configuration is not provided, return a mock service
	if cfg.S3Endpoint == "" || cfg.S3AccessKey == "" || cfg.S3SecretKey == "" || cfg.S3Bucket == "" {
		return &StorageService{
			client: nil,
			bucket: "",
			cfg:    cfg,
		}
	}

	// Initialize MinIO client
	endpoint := cfg.S3Endpoint
	if endpoint == "" {
		endpoint = "localhost:9000" // Default MinIO endpoint
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		Secure: false, // For local development
	})
	if err != nil {
		// In production, you might want to handle this error differently
		fmt.Printf("Failed to initialize S3 client: %v\n", err)
		return &StorageService{
			client: nil,
			bucket: "",
			cfg:    cfg,
		}
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	err = client.MakeBucket(ctx, cfg.S3Bucket, minio.MakeBucketOptions{Region: cfg.S3Region})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, cfg.S3Bucket)
		if errBucketExists == nil && exists {
			fmt.Printf("We already own %s\n", cfg.S3Bucket)
		} else {
			fmt.Printf("Failed to create bucket: %v\n", err)
		}
	} else {
		fmt.Printf("Successfully created bucket %s\n", cfg.S3Bucket)
	}

	return &StorageService{
		client: client,
		bucket: cfg.S3Bucket,
		cfg:    cfg,
	}
}

func (s *StorageService) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// If S3 is not configured, return a mock URL
	if s.client == nil {
		return fmt.Sprintf("/uploads/%s", filename), nil
	}

	// Upload file to S3
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := s.client.PutObject(ctx, s.bucket, filename, file, -1, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the URL of the uploaded file
	return fmt.Sprintf("%s/%s/%s", s.cfg.S3Endpoint, s.bucket, filename), nil
}

func (s *StorageService) GenerateFileName(originalName string, prefix string) string {
	extension := filepath.Ext(originalName)
	filename := fmt.Sprintf("%s_%d%s", prefix, getFileTimestamp(), extension)
	return filename
}

func getFileTimestamp() int64 {
	return time.Now().Unix()
}
