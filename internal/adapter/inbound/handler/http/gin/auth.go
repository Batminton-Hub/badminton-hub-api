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
	token = token[len("Bearer "):] // Remove "Bearer " prefix
	if err := m.MiddlewareUtil.Authenticate(token); err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}

	c.Next()
}
