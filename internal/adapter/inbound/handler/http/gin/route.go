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

func (m *MainRoute) RouteHealthCheck() {
	healthCheck := engine.Group("/")
	healthCheck.GET("/health-check", m.healthCheck.HealthCheck)
}

func (m *MainRoute) RouteGang() {
	gang := engine.Group("/gang")
	gang.Use(m.authentication.MiddleWare)
	// gang.GET("/:gang_id")                   // ดึงข้อมูลก๊วน
	gang.POST("/create", m.gang.CreateGang) // สร้างก๊วน
	// gang.GET("/search/:search_value")       // ค้นหาก๊วน
	// gang.POST("/join-request/:gang_id")     // ส่งคำขอเข้าร่วมก๊วน
	// gang.GET("/join-status/:gang_id")       // ดึงสถานะการเข้าร่วมก๊วน *polling
	// gang.PATCH("/update/:gang_id")          // อัพเดทก๊วน
	// gang.PUT("/allow-permisison")           // ให้สิทธิ์การแก้ไขก๊วน

	// gang.PATCH("/leave/:gang_id")     // ออกจากก๊วน
	// gang.PATCH("/dismiss-permission") // ยกเลิกสิทธิ์การแก้ไขก๊วน

}
