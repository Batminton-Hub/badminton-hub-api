package port

import "Badminton-Hub/internal/core/domain"

type AuthenticationSystem interface {
	MiddlewareService
	AuthenticationService
}

type AuthenticationService interface {
	Login(loginInfo domain.LoginInfo) (int, domain.RespLogin)
	Register(registerInfo domain.RegisterInfo) (int, domain.RespRegister)
}
type MiddlewareService interface {
	AuthenticateUtil
	MiddlewareUtil
}

type AuthenticateUtil interface {
	Authenticate(authInfo domain.AuthInfo) (int, domain.RespAuth)
}

type MiddlewareUtil interface {
	GenBearerToken(hashBody domain.HashAuth) (domain.BearerToken, error)
	ValidateBearerToken(token domain.BearerToken) (domain.AuthBody, error)
}
