package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr          string
	MongoURI          string
	MongoDatabase     string
	MongoCollectionItems string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI is required (set it in .env or the environment)")
	}

	db := getenvDefault("MONGODB_DATABASE", "restapi")
	itemsColl := getenvDefault("MONGODB_COLLECTION_ITEMS", "items")
	addr := getenvDefault("HTTP_ADDR", ":8080")

	return &Config{
		HTTPAddr:               addr,
		MongoURI:               mongoURI,
		MongoDatabase:          db,
		MongoCollectionItems:   itemsColl,
	}, nil
}

func getenvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
