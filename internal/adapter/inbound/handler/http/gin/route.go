package gin

import (
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MainRoute struct {
	MiddlewareController MiddlewareController
	MemberController     MemberController
	// MiddlewareController
	// MemberController
}

func NewGinMainRoute(
	middleware port.MiddlewareUtil,
	memberUtil port.MemberUtil,
) *MainRoute {
	return &MainRoute{
		MiddlewareController: &MiddlewareControllerImpl{middleware},
		MemberController:     &MemberControllerImpl{memberUtil},
	}
}

func (m *MainRoute) RouteMember() {
	r := gin.Default()
	member := r.Group("/member")
	member.Use(m.MiddlewareController.Authenticate)
	{
		member.POST("/register", m.MemberController.RegisterMember)
	}

	r.Run()
}
