package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Port string
}

func Load() (ApiConfig, error) {
	if err := godotenv.Load(); err != nil {
		return ApiConfig{}, err
	}

	port := os.Getenv("PORT")
	cfg := ApiConfig{
		Port: port,
	}
	return cfg, nil
}
