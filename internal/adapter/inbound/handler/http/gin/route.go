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
	{
		member.POST("/register", m.MemberController.RegisterMember)
		member.POST("/login", m.MemberController.Login)
		// member.GET("/authenticate", m.MiddlewareController.Authenticate, TestFunc())
	}

	r.Run()
}

// func TestFunc() func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		fmt.Println("Test function called")
// 		c.JSON(200, gin.H{"message": "Test function called"})
// 	}
// }
