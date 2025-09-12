package google

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/util"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuth = &domain.GoogleOAuth{}

func GoogleConfig(typeRedirect string) (*domain.GoogleOAuth, error) {
	config := util.LoadConfig()
	var callbackURL string
	switch typeRedirect {
	case domain.LOGIN:
		callbackURL = config.GoogleCallbackLoginURL
	case domain.REGISTER:
		callbackURL = config.GoogleCallbackRegisterURL
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
