package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
)

type Redirect struct {
	Google port.AuthenticationRedirect
}

func NewRedirect(
	google port.RedirectService,
) *Redirect {
	return &Redirect{
		Google: google,
	}
}

func (m *Redirect) Login(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	var response domain.RespRedirect
	var redirect port.AuthenticationRedirect
	switch info.Platform {
	case domain.GOOGLE:
		redirect = m.Google
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}

	return redirect.Login(info)
}

func (m *Redirect) Register(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	var response domain.RespRedirect
	var redirect port.AuthenticationRedirect
	switch info.Platform {
	case domain.GOOGLE:
		redirect = m.Google
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}

	return redirect.Register(info)
}
