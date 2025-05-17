package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	Database string
	HTTPPort string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system environment")
	}

	cfg := &Config{
		MongoURI: os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DB_NAME"),
		HTTPPort: os.Getenv("HTTP_PORT"),
	}

	if cfg.MongoURI == "" || cfg.Database == "" || cfg.HTTPPort == "" {
		log.Fatal("Missing Mongo configuration in environment")
	}

	return cfg
}
