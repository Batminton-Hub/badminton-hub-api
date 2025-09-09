package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
)

type CallbackService struct {
	Google port.CallbackService
}

func NewCallbackService(
	google port.CallbackService,
) *CallbackService {
	return &CallbackService{
		Google: google,
	}
}

func (c *CallbackService) Authenticate(info domain.AuthInfo) (int, domain.RespAuth) {
	var response domain.RespAuth
	var callback port.CallbackService
	switch info.Platform {
	case domain.GOOGLE:
		callback = c.Google
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}
	return callback.Authenticate(info)
}
