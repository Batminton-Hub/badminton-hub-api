package gin

import (
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MainRoute struct {
	MiddlewareController MiddlewareController
	MemberController     MemberController
	RedirectController   RedirectController
	Engine               *gin.Engine
}

func NewGinMainRoute(
	middleware port.MiddlewareUtil,
	memberUtil port.MemberUtil,
	redirectUtil port.RedirectUtil,
) *MainRoute {
	return &MainRoute{
		MiddlewareController: &MiddlewareControllerImpl{middleware},
		MemberController:     &MemberControllerImpl{memberUtil},
		RedirectController:   &RedirectControllerImpl{redirectUtil},
	}
}

type MemberController interface {
	RegisterMember(c *gin.Context)
	Login(c *gin.Context)
}

type MiddlewareController interface {
	Authenticate(c *gin.Context)
	GoogleLoginCallback(c *gin.Context)
	GoogleRegisterCallback(c *gin.Context)
}

type RedirectController interface {
	GoogleLogin(c *gin.Context)
	GoogleRegister(c *gin.Context)
}

type ProfileController interface {
	GetProfile(c *gin.Context)
}
