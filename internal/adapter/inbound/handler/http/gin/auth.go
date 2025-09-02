package gin

import (
	"Badminton-Hub/internal/core/port"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MiddlewareController interface {
	Authenticate(c *gin.Context)
	GoogleLoginCallback(c *gin.Context)
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

func (m *MiddlewareControllerImpl) GoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	httpStatus, response := m.MiddlewareUtil.GoogleLoginCallback(state, code)
	if httpStatus != http.StatusOK {
		fmt.Println("GoogleLoginCallback error:", response.Error)
		c.AbortWithStatus(httpStatus)
		return
	}

	c.Set("response", response)
	c.Set("type_login", "google")

	c.Next()
}
