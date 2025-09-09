package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"fmt"
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
	fmt.Println("CallbackService Authenticate Start")
	var response domain.RespAuth
	var callback port.CallbackService
	switch info.Platform {
	case domain.GOOGLE:
		fmt.Println("CallbackService Google Start")
		callback = c.Google
	default:
		response.Resp = domain.ErrPlatformNotSupport
		return response.Resp.HttpStatus, response
	}
	return callback.Authenticate(info)
}
