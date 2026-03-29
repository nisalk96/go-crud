package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr              string
	MongoURI              string
	MongoDatabase         string
	MongoCollectionMovies string
	UploadDir             string
	MaxUploadBytes        int64
	APIToken string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI is required (set it in .env or the environment)")
	}

	db := getenvDefault("MONGODB_DATABASE", "restapi")
	moviesColl := getenvDefault("MONGODB_COLLECTION_MOVIES", "movies")
	addr := getenvDefault("HTTP_ADDR", ":8080")
	uploadDir := getenvDefault("UPLOAD_DIR", "data/covers")

	maxMB := 10
	if s := os.Getenv("MAX_UPLOAD_MB"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			maxMB = v
		}
	}
	maxBytes := int64(maxMB) * 1024 * 1024

	apiToken := strings.TrimSpace(os.Getenv("API_TOKEN"))
	if apiToken == "" {
		return nil, fmt.Errorf("API_TOKEN is required (set it in .env or the environment)")
	}

	return &Config{
		HTTPAddr:              addr,
		MongoURI:              mongoURI,
		MongoDatabase:         db,
		MongoCollectionMovies: moviesColl,
		UploadDir:             uploadDir,
		MaxUploadBytes:        maxBytes,
		APIToken:              apiToken,
	}, nil
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
