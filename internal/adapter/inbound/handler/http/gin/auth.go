package gin

import (
	"Badminton-Hub/internal/core/port"
	"fmt"

	"github.com/gin-gonic/gin"
)

type MiddlewareController interface {
	Authenticate(c *gin.Context)
}
type MiddlewareControllerImpl struct {
	// MiddlewareUtil port.MiddlewareUtil
	port.MiddlewareUtil
}

func (m *MiddlewareControllerImpl) Authenticate(c *gin.Context) {
	fmt.Println("Authenticating request...")
	m.MiddlewareUtil.Authenticate("")
}
