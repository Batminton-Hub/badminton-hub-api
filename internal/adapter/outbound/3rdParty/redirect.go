package third_party

import (
	"Badminton-Hub/internal/adapter/outbound/3rdParty/google"
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
)

type ThirdPartyRedirectImpl struct {
	Google port.RedirectService
}

func New3rdPartyRedirect(
	cache port.CacheUtil,
) *ThirdPartyRedirectImpl {
	return &ThirdPartyRedirectImpl{
		Google: google.NewGoogleRedirect(cache),
	}
}

func (t *ThirdPartyRedirectImpl) Login(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	response := domain.RespRedirect{}
	switch info.Platform {
	case domain.GOOGLE:
		return t.Google.Login(info)
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}
}

func (t *ThirdPartyRedirectImpl) Register(info domain.RedirectLoginInfo) (int, domain.RespRedirect) {
	response := domain.RespRedirect{}
	switch info.Platform {
	case domain.GOOGLE:
		return t.Google.Register(info)
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}
}
