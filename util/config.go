package util

import (
	"Badminton-Hub/internal/core/domain"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuth = &domain.GoogleOAuth{}
var config = &domain.InternalConfig{}

const (
	MODE                 = "MODE"
	SERVER_PORT          = "SERVER_PORT"
	DB_NAME              = "DB_NAME"
	MONGO_DB_URL         = "MONGO_DB_URL"
	KEY_BEARER_TOKEN     = "KEY_BEARER_TOKEN"
	KEY_HASH_AUTH        = "KEY_HASH_AUTH"
	KEY_HASH_MEMBER      = "KEY_HASH_MEMBER"
	KEY_HASH_PASSWORD    = "KEY_HASH_PASSWORD"
	GOOGLE_LOGIN_URL     = "GOOGLE_LOGIN_URL"
	GOOGLE_REGISTER_URL  = "GOOGLE_REGISTER_URL"
	GOOGLE_CLIENT_ID     = "GOOGLE_CLIENT_ID"
	GOOGLE_CLIENT_SECRET = "GOOGLE_CLIENT_SECRET"
	REDIS_CACHE_ADDR     = "REDIS_CACHE_ADDR"
	REDIS_CACHE_PASSWORD = "REDIS_CACHE_PASSWORD"
	REDIS_CACHE_DB       = "REDIS_CACHE_DB"
	DEFAULT_AES_IV       = "DEFAULT_AES_IV"
	DEFAULT_GOOGLE_STATE = "DEFAULT_GOOGLE_STATE"
	BEARER_TOKEN_EXP     = "BEARER_TOKEN_EXP"
)

// Server Config
func SetConfig() error {
	godotenv.Load()

	config = &domain.InternalConfig{
		// Mode
		Mode: getEnv(MODE, "DEVERLOP"), // DEVERLOP, UAT , PRODUCTION

		// Server
		ServerPort: getEnv(SERVER_PORT, "8080"),

		// Mongo
		DBName:     getEnv(DB_NAME, "default_db_name"),
		MongoDBURL: getEnv(MONGO_DB_URL, "mongodb://localhost:27017"),

		// Key
		KeyBearerToken:  getEnv(KEY_BEARER_TOKEN, "0123456789ABCDEF"),
		KeyHashAuth:     getEnv(KEY_HASH_AUTH, "default_hash_key"),
		KeyHashMember:   getEnv(KEY_HASH_MEMBER, "default_hash_key"),
		KeyHashPassword: getEnv(KEY_HASH_PASSWORD, "default_hash_key"),

		// Google OAuth
		// GoogleLoginRedirectURL:    getEnv(GOOGLE_LOGIN_URL, "http://localhost:8080/member/auth/google/callback/login"),
		// GoogleRegisterRedirectURL: getEnv(GOOGLE_REGISTER_URL, "http://localhost:8080/member/auth/google/callback/register"),
		GoogleLoginRedirectURL:    getEnv(GOOGLE_LOGIN_URL, "http://localhost:8080/callback/google/login"),
		GoogleRegisterRedirectURL: getEnv(GOOGLE_REGISTER_URL, "http://localhost:8080/callback/google/register"),
		GoogleClinentID:           getEnv(GOOGLE_CLIENT_ID, "1030829763252-hngbodu9d2vqu2c82n80f86gl8urtq5n.apps.googleusercontent.com"),
		GoogleClientSecret:        getEnv(GOOGLE_CLIENT_SECRET, "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"),

		// Redis Cache
		RedisCacheAddr:     getEnv(REDIS_CACHE_ADDR, "localhost:6379"),
		RedisCachePassword: getEnv(REDIS_CACHE_PASSWORD, ""),
		RedisCacheDB:       getEnv(REDIS_CACHE_DB, 0),

		// RandomFunc
		DefaultAESIV:       getEnv(DEFAULT_AES_IV, []byte("0123456789ABCDEF")), // 16 bytes
		DefaultGoogleState: getEnv(DEFAULT_GOOGLE_STATE, "0123456789ABCDEF"),

		// Token
		BearerTokenExp: getEnv(BEARER_TOKEN_EXP, 5*time.Minute),
	}

	fmt.Println("config.TokenExpired ", config.BearerTokenExp)
	fmt.Println("config.DefaultAESIV ", config.DefaultAESIV)
	fmt.Println("config.RedisCacheDB ", config.RedisCacheDB)
	fmt.Println("config.DefaultGoogleState ", config.DefaultGoogleState)
	fmt.Println("config.GoogleLoginRedirectURL ", config.GoogleLoginRedirectURL)
	fmt.Println("config.GoogleRegisterRedirectURL ", config.GoogleRegisterRedirectURL)

	return nil
}

func LoadConfig() domain.InternalConfig {
	return *config
}

// Google Config
func GoogleConfig(typeRedirect string) (*domain.GoogleOAuth, error) {
	config := LoadConfig()
	var callbackURL string
	switch typeRedirect {
	case domain.LOGIN:
		callbackURL = config.GoogleLoginRedirectURL
	case domain.REGISTER:
		callbackURL = config.GoogleRegisterRedirectURL
	default:
		return nil, domain.ErrActionNotSupport.Err
	}
	fmt.Println("Callback URL:", callbackURL)
	googleOAuth.Config = &oauth2.Config{
		RedirectURL:  callbackURL,
		ClientID:     config.GoogleClinentID,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{
			domain.GoogleUserInfoEmail,
			domain.GoogleUserInfoProfile,
		},
		Endpoint: google.Endpoint,
	}
	return googleOAuth, nil
}

// Other Function
type TypeEnv interface {
	string | []byte | int | time.Duration
}

func getEnv[T any](keyEnv string, defaultVal T) T {
	value := os.Getenv(keyEnv)
	if value != "" {
		switch any(defaultVal).(type) {
		case string:
			return any(value).(T)
		case []byte:
			return any([]byte(value)).(T)
		case int:
			num, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Setting env key[%s] error : %s", keyEnv, err.Error())
			}
			return any(num).(T)
		case time.Duration:
			num, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Setting env key[%s] error : %s", keyEnv, err.Error())
			}
			return any(time.Duration(num) * time.Minute).(T)
		}
	}
	return defaultVal
}
