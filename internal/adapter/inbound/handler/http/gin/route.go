package gin

import (
	"Badminton-Hub/internal/core/port"
	"Badminton-Hub/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MainRoute struct {
	MiddlewareController MiddlewareController
	MemberController     MemberController
	Engine               *gin.Engine
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
	member := m.Engine.Group("/member")
	{
		member.POST("/register", m.MemberController.RegisterMember)
		member.POST("/login", m.MemberController.Login)
		member.GET("/google/login", m.MemberController.GoogleLogin)
		member.GET("/auth/google/login/callback", m.MiddlewareController.GoogleLoginCallback, TestLogin)
		member.GET("/google/register", m.MemberController.GoogleRegister)
		member.GET("/auth/google/register/callback", m.MemberController.GoogleRegisterCallback)
		// member.GET("/authenticate", m.MiddlewareController.Authenticate, TestFunc())
	}
}

func (m *MainRoute) Start() {
	m.Engine = gin.Default()
}

func (m *MainRoute) Run() {
	srv := util.HttpServer(m.Engine)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Listen error:", err)
		}
	}()
}
