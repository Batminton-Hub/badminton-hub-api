package gin

import (
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MainRoute struct {
	Engine               *gin.Engine
	MiddlewareController MiddlewareController
	MemberController     MemberControllerGroup
	RedirectController   RedirectController
}

func NewGinMainRoute(
	middleware port.MiddlewareUtil,
	memberUtil port.MemberUtilGroup,
	redirectUtil port.RedirectUtil,
) *MainRoute {
	return &MainRoute{
		MiddlewareController: &MiddlewareControllerImpl{middleware},
		RedirectController:   &RedirectControllerImpl{redirectUtil},
		MemberController:     &MemberControllerImpl{memberUtil},
	}
}

type MemberControllerGroup interface {
	MemberController
	ProfileController
}

type MemberController interface {
	RegisterMember(c *gin.Context)
	Login(c *gin.Context)
}

type ProfileController interface {
	GetProfile(c *gin.Context)
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
