package domain

import "golang.org/x/oauth2"

type InternalConfig struct {
	DBName     string
	ServerPort string
	MongoDBURL string

	// Key
	KeyHashAuth     string
	KeyHashMember   string
	KeyHashPassword string

	// Google OAuth
	GoogleLoginRedirectURL    string
	GoogleRegisterRedirectURL string
	GoogleClinentID           string
	GoogleClientSecret        string

	// Redis Cache
	RedisCacheAddr     string
	RedisCachePassword string
	RedisCacheDB       int
}

type GoogleOAuth struct {
	Config *oauth2.Config
	State  string
}
