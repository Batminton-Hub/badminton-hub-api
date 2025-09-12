package third_party

import (
	"Badminton-Hub/internal/adapter/outbound/3rdParty/google"
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
)

type ThirdPartyMiddlewareImpl struct {
	Google port.AuthenticateUtil
}

func New3rdPartyMiddleware(
	cache port.CacheUtil,
) *ThirdPartyMiddlewareImpl {
	return &ThirdPartyMiddlewareImpl{
		Google: google.NewGoogleMiddleware(cache),
	}
}

func (t *ThirdPartyMiddlewareImpl) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	response := domain.RespAuth{}
	switch info.Platform {
	case domain.GOOGLE:
		return t.Google.Authenticate(info)
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}
}
