package port

import "Badminton-Hub/internal/core/domain"

type RedirectService interface {
	AuthenticationRedirect
}

type AuthenticationRedirect interface {
	Login(info domain.RedirectLoginInfo) (int, domain.RespRedirect)
	Register(info domain.RedirectLoginInfo) (int, domain.RespRedirect)
}
