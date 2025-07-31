package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MemberController interface {
	RegisterMember(c *gin.Context)
	Login(c *gin.Context)
}
type MemberControllerImpl struct {
	MemberUtil port.MemberUtil
}

func (m *MemberControllerImpl) RegisterMember(c *gin.Context) {
	registerForm := domain.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	httpStatus, response := m.MemberUtil.RegisterMember(registerForm)

	c.JSON(httpStatus, response)
}

func (m *MemberControllerImpl) Login(c *gin.Context) {
	loginForm := domain.LoginForm{}
	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	httpStatus, response := m.MemberUtil.Login(loginForm)

	c.JSON(httpStatus, response)
}
