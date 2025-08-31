package gin

import (
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MiddlewareController interface {
	Authenticate(c *gin.Context)
}
type MiddlewareControllerImpl struct {
	port.MiddlewareUtil
}

func (m *MiddlewareControllerImpl) Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	code, resp := m.MiddlewareUtil.Authenticate(token)
	if code != 200 {
		c.AbortWithStatusJSON(code, resp)
		return
	}

	c.Next()
}
