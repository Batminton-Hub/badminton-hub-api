package port

import "Badminton-Hub/internal/core/domain"

type MiddlewareUtil interface {
	Authenticate(token string) (int, domain.AuthResponse)
}
