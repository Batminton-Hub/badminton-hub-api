package port

import "Badminton-Hub/internal/core/domain"

type MiddlewareUtil interface {
	// Authentication
	Encryptetion() Encryption
	Authenticate(token string) (int, domain.AuthResponse)
	GoogleLoginCallback(state, code string) (int, domain.ResponseGoogleLoginCallback)
	GoogleRegisterCallback(state, code string) (int, domain.ResponseGoogleRegisterCallback)
}

// type Authentication interface {
// 	Authenticate(token string) (int, domain.AuthResponse)
// }
