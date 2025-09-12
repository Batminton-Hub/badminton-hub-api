package domain

import (
	"time"

	"golang.org/x/oauth2"
)

type InternalConfig struct {
	// Mode
	Mode string

	// Server
	ServerPort string

	// Database
	DBName     string
	MongoDBURL string

	// Key
	KeyBearerToken  string
	KeyHashAuth     string
	KeyHashMember   string
	KeyHashPassword string

	// Google OAuth
	GoogleCallbackLoginURL    string
	GoogleCallbackRegisterURL string
	GoogleClinentID           string
	GoogleClientSecret        string

	// Redis Cache
	RedisCacheAddr     string
	RedisCachePassword string
	RedisCacheDB       int

	// RandomFunc
	DefaultAESIV       []byte // 16 bytes
	DefaultGoogleState string

	// Token
	BearerTokenExp time.Duration
}

type GoogleOAuth struct {
	Config *oauth2.Config
	State  string
}
