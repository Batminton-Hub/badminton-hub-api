package util

import (
	"Badminton-Hub/internal/core/domain"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuth = &domain.GoogleOAuth{}
var config = &domain.InternalConfig{}

// Server Config
func SetConfig() error {
	godotenv.Load()

	config = &domain.InternalConfig{
		// Mode
		Mode: getEnv("Mode", "DEVERLOP"), // DEVERLOP, UAT , PRODUCTION

		// Server
		ServerPort: getEnv("Server_Port", "8080"),

		// Mongo
		DBName:     getEnv("DB_Name", "default_db_name"),
		MongoDBURL: getEnv("MongoDB_URL", "mongodb://localhost:27017"),

		// Key
		KeyHashAuth:     getEnv("Key_Hash_Auth", "default_hash_key"),
		KeyHashMember:   getEnv("Key_Hash_Member", "default_hash_key"),
		KeyHashPassword: getEnv("Key_Hash_Password", "default_hash_key"),

		// Google OAuth
		GoogleLoginRedirectURL:    getEnv("Google_Login_URL", "http://localhost:8080/member/auth/google/callback/login"),
		GoogleRegisterRedirectURL: getEnv("Google_Register_URL", "http://localhost:8080/member/auth/google/callback/register"),
		GoogleClinentID:           getEnv("Google_Client_ID", "1030829763252-hngbodu9d2vqu2c82n80f86gl8urtq5n.apps.googleusercontent.com"),
		GoogleClientSecret:        getEnv("Google_Client_Secret", "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"),

		// Redis Cache
		RedisCacheAddr:     getEnv("Redis_Cache_Addr", "localhost:6379"),
		RedisCachePassword: getEnv("Redis_Cache_Password", ""),
		RedisCacheDB:       getEnv("Redis_Cache_DB", 0),

		// RandomFunc
		DefaultAESIV:       getEnv("Default_AES_IV", []byte("0123456789ABCDEF")), // 16 bytes
		DefaultGoogleState: getEnv("Default_Google_State", "0123456789ABCDEF"),
	}

	return nil
}

func LoadConfig() domain.InternalConfig {
	return *config
}

// Google Config
func GoogleConfig(typeRedirect string) (*domain.GoogleOAuth, error) {
	config := LoadConfig()
	var redirectURL string
	switch typeRedirect {
	case "LOGIN":
		redirectURL = config.GoogleLoginRedirectURL
	case "REGISTER":
		redirectURL = config.GoogleRegisterRedirectURL
	default:
		redirectURL = strings.ToUpper(typeRedirect)
	}
	fmt.Println("Redirect URL:", redirectURL)
	googleOAuth.Config = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     config.GoogleClinentID,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	// googleOAuth.State = "randomstate"
	return googleOAuth, nil
}

// Other Function
func getEnv[T any](keyEnv string, defaultVal T) T {
	value := os.Getenv(keyEnv)
	// if value != "" {
	// 	return value
	// } else {
	// 	return defaultVal
	// }
	if value != "" {
		return any(value).(T)
	}
	return defaultVal
}
