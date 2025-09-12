package port

import "Badminton-Hub/internal/core/domain"

type CallbackService interface {
	MiddlewareCallback
}

type MiddlewareCallback interface {
	Authenticate(info domain.AuthInfo) (int, domain.RespAuth)
}
