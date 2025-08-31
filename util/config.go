package util

import (
	"Badminton-Hub/internal/core/domain"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (domain.InternalConfig, error) {
	godotenv.Load()

	config := domain.InternalConfig{
		DBName:          getEnv("DB_Name", "default_db_name"),
		MongoDBURL:      getEnv("MongoDB_URL", "mongodb://localhost:27017"),
		KeyHashPassword: getEnv("Key_Hash_Password", "default_hash_key"),
	}

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
