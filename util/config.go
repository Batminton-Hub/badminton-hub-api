package util

import (
	"Badminton-Hub/internal/core/domain"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (domain.InternalConfig, error) {
	godotenv.Load()

	var config domain.InternalConfig

	config.DBName = getEnv("DB_Name", "default_db_name")

	return config, nil
}

func getEnv(keyEnv, defaultVal string) string {
	value := os.Getenv(keyEnv)
	if value != "" {
		return value
	} else {
		return defaultVal
	}
}
