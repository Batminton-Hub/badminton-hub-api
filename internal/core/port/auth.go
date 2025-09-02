package port

import "Badminton-Hub/internal/core/domain"

type MiddlewareUtil interface {
	// Authentication
	Authenticate(token string) (int, domain.AuthResponse)
	Encryptetion() Encryption
}

// type Authentication interface {
// 	Authenticate(token string) (int, domain.AuthResponse)
// }
