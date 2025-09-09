package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/util"

	"github.com/gin-gonic/gin"
)

func (m *MainRoute) Start() {
	m.engine = gin.Default()
}

func (m *MainRoute) Run() {
	util.RunServer(m.engine)
}

func (m *MainRoute) RouteAuthenticationSystem() {
	r := m.engine
	authentication := r.Group("/authentication")
	authentication.POST("/login", m.authentication.Login)
	authentication.POST("/register", m.authentication.Register)
}

func (m *MainRoute) RouteRedirect() {
	r := m.engine
	redirect := r.Group("/redirect")
	redirect.GET("/:platform/login", m.redirect.Login)
	redirect.GET("/:platform/register", m.redirect.Register)
}

func (m *MainRoute) RouteCallback() {
	r := m.engine
	callback := r.Group("/callback")
	callback.GET("/:platform/login", m.authentication.MiddleWare(domain.LOGIN), m.authentication.Login)
	callback.GET("/:platform/register", m.authentication.MiddleWare(domain.REGISTER), m.authentication.Register)
}
