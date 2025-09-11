package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationSystemController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	MiddleWare(c *gin.Context)
}

type AuthenticationSystem struct {
	authenticationSystem port.AuthenticationSystem
}

func (a *AuthenticationSystem) Login(c *gin.Context) {
	loginInfo := domain.LoginInfo{
		Platform:     getPlatform(c),
		PlatformData: getPlatformData(c),
	}

	loginForm := domain.LoginForm{}
	if err := c.ShouldBindJSON(&loginForm); err != nil && loginInfo.PlatformData == nil {
		RespAuth(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg, "")
		return
	}

	loginInfo.LoginForm = loginForm
	code, response := a.authenticationSystem.Login(loginInfo)

	RespAuth(c, code, response.Resp.Code, response.Resp.Msg, response.BearerToken)
}

func (a *AuthenticationSystem) Register(c *gin.Context) {
	registerInfo := domain.RegisterInfo{
		Platform:     getPlatform(c),
		PlatformData: getPlatformData(c),
	}
	registerForm := domain.RegisterForm{}
	err := c.ShouldBindJSON(&registerForm)
	if err != nil {
		RespAuth(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg, "")
		return
	}

	registerInfo.RegisterForm = registerForm
	code, response := a.authenticationSystem.Register(registerInfo)

	RespAuth(c, code, response.Resp.Code, response.Resp.Msg, response.BearerToken)
}

func (a *AuthenticationSystem) MiddleWare(c *gin.Context) {
	platform := getPlatformParam(c)
	authInfo := domain.AuthInfo{
		BearerToken: domain.BearerToken{
			Token: getBearerToken(c),
		},
		State:    getState(c),
		Code:     getCode(c),
		Action:   getAction(c),
		Platform: platform,
	}
	httpStatus, response := a.authenticationSystem.Authenticate(authInfo)
	if response.Resp.Status == domain.ERROR {
		RespMiddleWare(c, httpStatus, response.Resp.Code, response.Resp.Msg)
		return
	}

	c.Set(domain.PlatformData, response.PlatformData)
	c.Set(domain.Platform, platform)
	c.Set(domain.UserID, response.AuthBody.Data.UserID)

	c.Next()
}
