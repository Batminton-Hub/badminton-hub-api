package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MemberController interface {
	RegisterMember(c *gin.Context)
}
type MemberControllerImpl struct {
	// Middleware port.MiddlewareUtil
	MemberUtil port.MemberUtil
	// port.MiddlewareUtil
	// port.MemberUtil
}

func (m *MemberControllerImpl) RegisterMember(c *gin.Context) {
	registerForm := domain.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := m.MemberUtil.RegisterMember(registerForm); err != nil {
		c.JSON(500, gin.H{"error": "Failed to register member"})
		return
	}

	c.JSON(200, gin.H{"message": "Member registered successfully"})
}
