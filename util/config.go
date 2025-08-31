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
		ServerPort:      getEnv("Server_Port", "8080"),
		MongoDBURL:      getEnv("MongoDB_URL", "mongodb://localhost:27017"),
		KeyHashAuth:     getEnv("Key_Hash_Auth", "default_hash_key"),
		KeyHashMember:   getEnv("Key_Hash_Member", "default_hash_key"),
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
