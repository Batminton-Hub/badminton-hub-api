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
	observability        port.Observability
}

func (a *AuthenticationSystem) Login(c *gin.Context) {
	counter := domain.MetricsCounter{
		Name: "login_request",
		Help: "Number of login requests",
	}
	countLogin := a.observability.Metrics().Counter(counter)
	countLogin.Inc()

	platform := getPlatform(c)
	loginInfo := domain.LoginInfo{
		Platform:     platform,
		PlatformData: getPlatformData(c),
		TypeSystem:   getTypeSystem(platform),
	}

	loginForm := domain.LoginForm{}
	if err := c.ShouldBindJSON(&loginForm); err != nil && loginInfo.PlatformData == nil {
		RespAuth(c, http.StatusBadRequest, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg, "")
		return
	}

	loginInfo.LoginForm = loginForm
	code, response := a.authenticationSystem.Login(loginInfo)

	switch response.Resp.Status {
	case domain.SUCCESS:
		counter = domain.MetricsCounter{
			Name: "login_request_success",
			Help: "Number of login requests success",
		}
		countLogin = a.observability.Metrics().Counter(counter)
		countLogin.Inc()
	case domain.ERROR:
		counter = domain.MetricsCounter{
			Name: "login_request_error",
			Help: "Number of login requests error",
		}
		countLogin = a.observability.Metrics().Counter(counter)
		countLogin.Inc()
	}

	RespAuth(c, code, response.Resp.Code, response.Resp.Msg, response.BearerToken)
}

func (a *AuthenticationSystem) Register(c *gin.Context) {
	platform := getPlatform(c)
	registerInfo := domain.RegisterInfo{
		Platform:     platform,
		PlatformData: getPlatformData(c),
		TypeSystem:   getTypeSystem(platform),
	}
	registerForm := domain.RegisterForm{}
	err := c.ShouldBindJSON(&registerForm)
	if err != nil && registerInfo.PlatformData == nil {
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
		State:      getState(c),
		Code:       getCode(c),
		Action:     getAction(c),
		TypeSystem: getTypeSystem(platform),
		Platform:   platform,
	}

	authInfo.TypeSystem = getTypeSystem(platform)

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
