package gin

import (
	"Badminton-Hub/util"

	"github.com/gin-gonic/gin"
)

type RunServer func()

var engine *gin.Engine

func (m *MainRoute) Start() RunServer {
	engine = gin.Default()
	return runServer()
}

func runServer() RunServer {
	return func() {
		util.RunServer(engine)
	}
}

func (m *MainRoute) RouteAuthenticationSystem() {
	authentication := engine.Group("/authentication")
	authentication.POST("/login", m.authentication.Login)
	authentication.POST("/register", m.authentication.Register)
}

func (m *MainRoute) RouteRedirect() {
	redirect := engine.Group("/redirect")
	redirect.GET("/:platform/login", m.redirect.Login)
	redirect.GET("/:platform/register", m.redirect.Register)
}

func (m *MainRoute) RouteCallback() {
	callback := engine.Group("/callback")
	callback.GET("/:platform/login", m.authentication.MiddleWare, m.authentication.Login)
	callback.GET("/:platform/register", m.authentication.MiddleWare, m.authentication.Register)
}

func (m *MainRoute) RouteMember() {
	member := engine.Group("/member")
	member.Use(m.authentication.MiddleWare)
	member.GET("/profile", m.member.GetProfile)
	member.PATCH("/profile", m.member.UpdateProfile)
}

func (m *MainRoute) RouteObservability() {
	observability := engine.Group("/")
	observability.GET("/metrics", m.observability.Metrics)
}
