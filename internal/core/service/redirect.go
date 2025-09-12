package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
)

type Redirect struct {
	thirdParty port.RedirectService
}

func NewRedirect(
	thirdParty port.RedirectService,
) *Redirect {
	return &Redirect{
		thirdParty: thirdParty,
	}
}

func (m *Redirect) Login(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	return m.thirdParty.Login(info)
}

func (m *Redirect) Register(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	return m.thirdParty.Register(info)
}
