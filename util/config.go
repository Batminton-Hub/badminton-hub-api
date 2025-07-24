package util

import (
	"Badminton-Hub/internal/core/domain"
	"os"
)

func LoadConfig() (domain.InternalConfig, error) {
	config := domain.InternalConfig{}
	if dbName := os.Getenv("DB_Name"); dbName != "" {
		config.DBName = dbName
	} else {
		config.DBName = "default_db_name" // Default value if not set
	}

	return config, nil
}
