package gin

import (
	"Badminton-Hub/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *MainRoute) RouteMember() {
	member := m.Engine.Group("/member")
	{
		member.POST("/register", m.MemberController.RegisterMember)
		member.POST("/login", m.MemberController.Login)
		member.GET("/google/login", m.RedirectController.GoogleLogin)
		member.GET("/google/register", m.RedirectController.GoogleRegister)
		member.GET("/auth/google/callback/login", m.MiddlewareController.GoogleLoginCallback, m.MemberController.Login)
		member.GET("/auth/google/callback/register", m.MiddlewareController.GoogleRegisterCallback, m.MemberController.RegisterMember)

		member.GET("/profile", m.MiddlewareController.Authenticate, m.MemberController.GetProfile)
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
