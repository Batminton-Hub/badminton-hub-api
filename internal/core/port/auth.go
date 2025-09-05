package port

import "Badminton-Hub/internal/core/domain"

type MiddlewareUtil interface {
	Authentication
	MiddlewareEncryption
	MiddlewareCallback
}

type Authentication interface {
	Authenticate(token string) (int, domain.AuthResponse)
}

type MiddlewareEncryption interface {
	Encryption() Encryption
}

type MiddlewareCallback interface {
	GoogleLoginCallback(state, code string) (int, domain.ResponseGoogleLoginCallback)
	GoogleRegisterCallback(state, code string) (int, domain.ResponseGoogleRegisterCallback)
}
